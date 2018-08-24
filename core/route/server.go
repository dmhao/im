package route

import (
	"context"
	"im/core/log"
	"net"
	"sync"
)

type routeServer struct {
	ls  *net.TCPListener
	iss *imServers

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup
	once   *sync.Once
}

func NewRouteServer(listener *net.TCPListener) *routeServer {
	rs := &routeServer{
		wg:   &sync.WaitGroup{},
		once: &sync.Once{},
		ls:   listener,
	}
	rs.iss = &imServers{
		rw:   &sync.RWMutex{},
		data: make(map[*imServerConn]*apps),
	}
	rs.ctx, rs.cancel = context.WithCancel(context.Background())
	return rs
}

//添加会话 至 会话集合
func (rs *routeServer) insertImServerConn(isc *imServerConn) {
	rs.iss.rw.Lock()
	rs.iss.data[isc] = &apps{rw: &sync.RWMutex{}, data: make(map[int]*appInfo)}
	rs.iss.rw.Unlock()
}

//会话集合中 删除某个会话
func (rs *routeServer) removeImServerConn(isc *imServerConn) {
	rs.iss.rw.Lock()
	delete(rs.iss.data, isc)
	rs.iss.rw.Unlock()
}

//获取所有im服务的会话列表
func (rs *routeServer) getAllImServerConn() []*imServerConn {
	rs.iss.rw.RLock()

	var iscs []*imServerConn
	for isc, _ := range rs.iss.data {
		iscs = append(iscs, isc)
	}
	rs.iss.rw.RUnlock()
	return iscs
}

func (rs *routeServer) StartRouteServer() {
	defer func() {
		log.Infoln("路由服务关闭完成")
		recoverFunc()
		rs.wg.Done()
	}()

	for {
		select {
		case <-rs.ctx.Done():
			log.Infoln("路由服务准备关闭")
			return
		default:
			if rs.ls == nil {
				return
			}
			conn, err := rs.ls.AcceptTCP()
			if err != nil {
				log.Infoln(err)
				continue
			}
			isc := NewImServerConn(rs, conn)
			rs.insertImServerConn(isc)
			rs.wg.Add(1)
			go isc.handlerConn()
		}
	}
}

func (rs *routeServer) Start() {
	rs.wg.Add(1)
	go rs.StartRouteServer()
	log.Infoln("路由启动完成")
	rs.wg.Wait()
}

func (rs *routeServer) Stop() {
	rs.once.Do(func() {
		rs.ls.Close()
		rs.cancel()
	})
}
