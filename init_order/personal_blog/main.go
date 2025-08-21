package main

import (
	"personal_blog/config"
	"personal_blog/router"
)

func main() {
	config.InitConfig()
	r := router.SetRouter()
	r.Run(":9001")
}
