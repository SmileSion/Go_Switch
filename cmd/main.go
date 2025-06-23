package main

import (
	"edulimitrate/config"
	"edulimitrate/router"
	"fmt"
)

func main() {
	config.InitConfig()
	config.InitDB()
	config.InitRedis()
	r := router.InitRouter()
	r.Run(":" + fmt.Sprint(config.Conf.Server.Port))
}
