package main

import (
	"flag"
	"fmt"
	"zero-demo/common/middleware"

	"zero-demo/user-api/internal/config"
	"zero-demo/user-api/internal/handler"
	"zero-demo/user-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 第三个参数代表会使用到环境变量中的配置，同时还需要在配置文件中将对应的参数设置为读取环境变量 ${XX}
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	// 全局中间件
	server.Use(middleware.NewGlobalMiddleware().Handle)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
