package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ahl5esoft/go-skill/hide-ctx/api"
	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/apisvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/ginsvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/goredissvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/iocsvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/ossvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/pathsvc"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	wd, _ := os.Getwd()
	ioPath := pathsvc.NewIOPath(wd, "hide-ctx")
	ioFactory := ossvc.NewIOFactory(ioPath)

	generateFlag := flag.Bool("api", false, "generate api metadata.go")
	flag.Parse()
	if *generateFlag {
		if err := apisvc.GenerateMetadata(ioFactory, ioPath, "github.com"); err != nil {
			panic(err)
		}
		return
	}

	redis := goredissvc.NewRedis(redis.Options{
		Addr: "127.0.0.1:6379",
	})
	iocsvc.Set(
		new(contract.IRedis),
		redis,
	)

	apiFactory := apisvc.NewFactory()
	api.Register(apiFactory)

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	ginsvc.NewPostOption(apiFactory)(app)
	fmt.Println("gin启动")
	err := app.Run(":10001")
	if err != nil {
		panic(err)
	}
}
