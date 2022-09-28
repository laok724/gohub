package main

import (
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"gohub/bootstrap"
	btsConfig "gohub/config"
	"gohub/pkg/config"
)

func init() {
	btsConfig.Initialize()
}
func main() {
	var env string
	flag.StringVar(&env, "env", "", "加载.env文件，如--env=testing加载的是.env.testing文件")
	flag.Parse()
	config.InitConfig(env)

	router := gin.New()
	bootstrap.SetupRoute(router)
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}

}
