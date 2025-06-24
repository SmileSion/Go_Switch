package main

import (
	"edulimitrate/config"
	"edulimitrate/handler"
	"edulimitrate/router"
	"fmt"
)

func main() {
	config.InitConfig()
	config.InitDB()
	config.InitRedis()
	handler.InitDecryptedSecret()
	r := router.InitRouter()
	r.Run(":" + fmt.Sprint(config.Conf.Server.Port))
}
