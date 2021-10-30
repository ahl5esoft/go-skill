package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ahl5esoft/go-skill/easy-api/api"
	"github.com/ahl5esoft/go-skill/easy-api/service/apisvc"
	"github.com/ahl5esoft/go-skill/easy-api/service/ginsvc"
	"github.com/ahl5esoft/go-skill/easy-api/service/ossvc"
	"github.com/ahl5esoft/go-skill/easy-api/service/pathsvc"

	"github.com/gin-gonic/gin"
)

func main() {
	wd, _ := os.Getwd()
	ioPath := pathsvc.NewIOPath(wd, "one")
	ioFactory := ossvc.NewIOFactory(ioPath)

	generateFlag := flag.Bool("api", false, "generate api metadata.go")
	flag.Parse()
	if *generateFlag {
		if err := apisvc.GenerateMetadata(ioFactory, ioPath, "github.com"); err != nil {
			panic(err)
		}
	}

	apiFactory := apisvc.NewFactory()
	api.Register(apiFactory)

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	ginsvc.NewPostOption(apiFactory)(app)
	fmt.Println("gin启动")
	err := app.Run(":9999")
	if err != nil {
		panic(err)
	}
}
