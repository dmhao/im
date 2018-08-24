package config

import "time"

type CommonConf struct {
	DBUser         	string        	`yaml:"dbUser"`
	DBPassword     	string        	`yaml:"dbPassword"`
	DBHost         	string        	`yaml:"dbHost"`
	DBPort         	string        	`yaml:"dbPort"`
	DBDatabase     	string        	`yaml:"dbDatabase"`
	RedisAddr      	string        	`yaml:"redisAddr"`
	RedisPassWord	string 			`yaml:"redisPassWord"`
	LogDir         	string        	`yaml:"logDir"`
	RotationFormat 	string        	`yaml:"rotationFormat"`
	LogExpire      	time.Duration 	`yaml:"logExpire"`
	RotationTime   	time.Duration 	`yaml:"rotationTime"`
}

type ClientConf struct {
	ClientServerAddr 	string 	`yaml:"clientServerAddr"`
	ApiPort          	string 	`yaml:"apiPort"`
	ApiPProf         	bool   	`yaml:"apiPProf"`
}

type RouteConf struct {
	RouteServerAddr 	string 	`yaml:"routeServerAddr"`
}
