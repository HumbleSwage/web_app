package main

import (
	"web_app/config"
	"web_app/pkg/utils"
	"web_app/repository/cache"
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
	utils.InitLog()
	dao.InitMySql()
	cache.InitRedis()
}
