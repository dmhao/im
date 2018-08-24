package server

import (
	"context"
	"im/core"
	"im/core/log"
	"im/core/storage"
	"im/core/tools"
	"im/mp"
	"im/mp/common"
	"math/rand"
	"net"
	"sync"
	"time"
)

var ApiMessagesCh chan []*core.ApiMessage

func init() {
	ApiMessagesCh = make(chan []*core.ApiMessage, 100)
}

type imServer struct {
	ls        *net.TCPListener
	rt        *router
	allApp    *apps
	connCount *int32

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	once   *sync.Once
}

func NewClientServer(listener *net.TCPListener) *imServer {
	im := &imServer{
		wg:        &sync.WaitGroup{},
		once:      &sync.Once{},
		ls:        listener,
		connCount: new(int32),
	}
	im.allApp = &apps{
		rw:   &sync.RWMutex{},
		data: make(map[int]*appInfo),
	}
	im.ctx, im.cancel = context.WithCancel(context.Background())
	return im
}

func (im *imServer) Start() {
	im.wg.Add(1)
	go im.StartClientServer()

	im.wg.Add(1)
	go im.StartRouteConn()

	im.wg.Add(1)
	go im.StartReportTask(60 * time.Second)

	im.wg.Add(1)
	go im.StartApiListen()

	log.Infoln("im启动完成")
	im.wg.Wait()
}

//启动api服务
func (im *imServer) StartApiListen() {
	defer im.wg.Done()
	for {
		select {
		case <-im.ctx.Done():
			if len(ApiMessagesCh) == 0 {
				return
			}
		case apiMessages := <-ApiMessagesCh:
			if apiMessages == nil {
				return
			}
			im.handlerApiMessages(apiMessages)
		}
	}
}

//启动连接路由协程
func (im *imServer) StartRouteConn() {
	defer func() {
		im.Stop()
		log.Infoln("路由服务关闭完成， 关闭整体服务")
		recoverFunc()
		im.wg.Done()
	}()

	im.rt = &router{
		rw:   &sync.RWMutex{},
		data: make(map[string]*routeConn),
	}
	zeroRouteCount := 0
	for {
		select {
		case <-im.ctx.Done():
			log.Infoln("关闭路由服务-收到服务关闭信号")
			return
		default:
			routeServices := storage.GetAllRouteService()
			if len(routeServices) == 0 {
				zeroRouteCount++
				if zeroRouteCount >= 3 {
					log.Infoln("关闭路由服务-可连接路由数量为0")
					return
				}
			} else {
				zeroRouteCount = 0
			}
			im.rt.rw.Lock()
			for _, routeService := range routeServices {
				if _, ok := im.rt.data[routeService.Addr]; !ok {
					addr, err := net.ResolveTCPAddr("tcp", routeService.Addr)
					conn, err := net.DialTCP("tcp", nil, addr)
					if err != nil {
						log.Infoln("关闭路由服务-连接路由服务失败", routeService.Addr, err)
						im.rt.rw.Unlock()
						return
					}
					im.rt.data[routeService.Addr] = newRouteConn(im, conn, routeService.Addr)
					im.wg.Add(1)
					go im.rt.data[routeService.Addr].handlerConn()
				}
			}
			im.rt.rw.Unlock()
		}
		time.Sleep(10 * time.Second)
	}
}

//启动客户端连接服务
func (im *imServer) StartClientServer() {
	defer func() {
		im.Stop()
		log.Infoln("连接服务关闭完成，关闭整体服务")
		recoverFunc()
		im.wg.Done()
	}()

	for {
		select {
		case <-im.ctx.Done():
			log.Infoln("关闭连接服务-收到服务关闭信号")
			return
		default:
			if im.ls == nil {
				return
			}
			conn, err := im.ls.AcceptTCP()
			if err != nil {

				log.Infoln("客户端连接accept失败", err)
				continue
			}
			//包装客户端连接
			cc := newClientConn(im, conn)
			im.wg.Add(1)
			go cc.handlerConn()
		}
	}
}

//启动统计报告协程
func (im *imServer) StartReportTask(duration time.Duration) {
	defer im.wg.Done()
	for {
		t := time.NewTimer(duration)
		select {
		case <-im.ctx.Done():
			t.Stop()
			return
		case <-t.C:
			log.Infoln("当前连接用户数:", *im.connCount)
		}
	}

}

//查询当前所有的用户链接
func (im *imServer) findAllUserConn() []*clientConn {
	im.allApp.rw.RLock()
	var allApps []*appInfo
	var allUserConn []*clientConn
	for _, ai := range im.allApp.data {
		allApps = append(allApps, ai)
	}
	if len(allApps) > 0 {
		for _, ai := range allApps {
			ai.rw.RLock()
			for _, userConn := range ai.users {
				allUserConn = append(allUserConn, userConn)
			}
			ai.rw.RUnlock()
		}
	}
	im.allApp.rw.RUnlock()
	return allUserConn
}

//上报消息至某一台路由服务
func (im *imServer) pushOneRouteMessage(appId int, routeMsg *core.RouteMessage, routeAddr string) {
	var routeConn *routeConn
	im.rt.rw.RLock()
	if rc, ok := im.rt.data[routeAddr]; ok {
		routeConn = rc
	}
	im.rt.rw.RUnlock()
	if routeConn != nil {
		routeConn.sendMessage(routeMsg)
	}
}

func (im *imServer) pushRouteMessage(appId int, routeMsg *core.RouteMessage, allRoute bool) {
	//查询分配给当前App的路由，上报或转发数据
	routeServices := storage.GetAppRouteService(appId)
	if allRoute {
		im.rt.rw.RLock()
		var routeConns []*routeConn
		for _, rs := range routeServices {
			if rc, ok := im.rt.data[rs.Addr]; ok {
				routeConns = append(routeConns, rc)
			}
		}
		im.rt.rw.RUnlock()
		for _, rc := range routeConns {
			rc.sendMessage(routeMsg)
		}
	} else {
		var routeConn *routeConn
		im.rt.rw.RLock()
		for range routeServices {
			index := rand.Intn(len(routeServices))
			addr := routeServices[index].Addr
			if rc, ok := im.rt.data[addr]; ok {
				routeConn = rc
				break
			}
		}
		if routeConn == nil {
			for _, rs := range routeServices {
				if rc, ok := im.rt.data[rs.Addr]; ok {
					routeConn = rc
					break
				}
			}
		}
		im.rt.rw.RUnlock()
		if routeConn != nil {
			routeConn.sendMessage(routeMsg)
		}
	}
}

//推送单聊消息到路由服务
func (im *imServer) sendRouteImMsg(appId int, cmd int8, msg *common.MsgIm) {
	routeMsgIm := &mp.RouteMsgIm{
		SenderId: msg.GetSenderId(),
		ReceiverId: msg.GetReceiverId(),
		AppId: int32(appId),
		TransportData: msg,
	}

	routeMsg := &core.RouteMessage{
		Cmd: cmd,
		RouteData: routeMsgIm,
	}
	im.pushRouteMessage(appId, routeMsg, false)
	log.Infoln("消息上报路由完成", tools.JsonMarshal(msg))
}

func (im *imServer) Stop() {
	im.once.Do(func() {
		close(ApiMessagesCh)
		im.ls.Close()
		im.cancel()
	})
}
