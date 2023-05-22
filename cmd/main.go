package main

import (
	"web_app/config"
	"web_app/repository/db/dao"
	"web_app/router"
)

func main() {
	loading()
	r := router.NewRouter()
	_ = r.Run(config.Config.System.HttpPort)
}

func loading() {
	config.InitConfig()
	dao.InitMySql()
}
