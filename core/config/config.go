package config

import (
	"gopkg.in/yaml.v1"
	"im/core/tools"
	"io/ioutil"
	"path/filepath"
)

var confName = "config.yaml"

var commonConf *CommonConf
var clientConf *ClientConf
var routeConf *RouteConf

func GetRouteConf() *RouteConf {
	return routeConf
}

func GetClientConf() *ClientConf {
	return clientConf
}

func GetCommonConf() *CommonConf {
	return commonConf
}

func init() {
	commonConf = &CommonConf{}
	clientConf = &ClientConf{}
	routeConf = &RouteConf{}
}

func InitConfig() error {
	wd, err := tools.GetWorkDir()
	if err != nil {
		return err
	}
	confFile := filepath.Join(wd, confName)
	confData, err := ioutil.ReadFile(confFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(confData, commonConf)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(confData, clientConf)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(confData, routeConf)
	if err != nil {
		return err
	}
	return nil
}
