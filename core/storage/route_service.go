package storage

type RouteService struct {
	Id         int
	AppId      int
	Addr       string
	CreateTime int64
}

const ShowRouteService = 1

func GetAllRouteService() []*RouteService {
	var routeServices []*RouteService
	dbClient.Model(&RouteService{}).
		Where("status=? and app_id=0", ShowRouteService).Find(&routeServices)
	return routeServices
}

func GetAppRouteService(appId int) []*RouteService {
	cacheRouteService, err := GetAppRouteServiceCacheByAppId(appId)
	if err == nil && len(cacheRouteService) > 0 {
		return cacheRouteService
	}
	var routeServices []*RouteService
	dbClient.Model(&RouteService{}).
		Where("status=? and app_id=?", ShowRouteService, appId).Find(&routeServices)
	if len(routeServices) > 0 {
		SetAppRouteServiceCacheByAppId(appId, routeServices)
	}
	return routeServices
}
