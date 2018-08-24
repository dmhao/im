package storage

type App struct {
	AppId      int `gorm:"primary_key"`
	SecretId   string
	SecretKey  string
	AppName    string
	CreateTime int64
	Status     int8
}

const ShowAppStatus = 1

func GetAppBySecret(secretId string, secretKey string) *App {
	app := &App{}
	dbClient.Model(app).Where(" secret_id=? and secret_key=? and status=?", secretId, secretKey, ShowAppStatus).
		First(app)
	return app
}

func GetAppById(appId int) *App {
	app := &App{}
	cacheApp, err := GetAppInfoCacheByAppId(appId)
	if err == nil && cacheApp.AppId != 0 {
		return cacheApp
	}
	dbClient.Model(app).Where(" app_id=? and status=? ", appId, ShowAppStatus).
		First(app)
	if app.AppId != 0 {
		SetAppInfoCacheByAppId(appId, app)
	}
	return app
}
