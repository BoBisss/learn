package main

import (
	"personal_blog/config"
	"personal_blog/router"
)

func main() {
	config.InitConfig()
	r := router.SetRouter()
	port := config.AppConfig.Port
	r.Run(port)
}
