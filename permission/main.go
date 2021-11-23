package main

import (
	"flag"
	"os"

	"github.com/ahl5esoft/go-skill/permission/api"
	"github.com/ahl5esoft/go-skill/permission/contract"
	relationaloperator "github.com/ahl5esoft/go-skill/permission/model/enum/relational-operator"
	valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"
	"github.com/ahl5esoft/go-skill/permission/model/global"
	"github.com/ahl5esoft/go-skill/permission/service/ginsvc"
	"github.com/ahl5esoft/go-skill/permission/service/usersvc"
	"github.com/ahl5esoft/lite-go/service/apisvc"
	lgginsvc "github.com/ahl5esoft/lite-go/service/ginsvc"
	"github.com/ahl5esoft/lite-go/service/mongosvc"
	"github.com/ahl5esoft/lite-go/service/ossvc"
	"github.com/ahl5esoft/lite-go/service/pathsvc"
)

func main() {
	wd, _ := os.Getwd()
	ioPath := pathsvc.NewIOPath(wd, "permission")
	ioFactory := ossvc.NewIOFactory(ioPath)

	// 生成api
	// go build github.com/ahl5esoft/go-skill/permission && permission --api
	generateFlag := flag.Bool("api", false, "generate api metadata.go")
	flag.Parse()
	if *generateFlag {
		if err := apisvc.GenerateMetadata(ioFactory, ioPath, "github.com"); err != nil {
			panic(err)
		}
		return
	}

	apiFactory := apisvc.NewFactory()
	api.Register(apiFactory)

	dbFactory := mongosvc.New(nil, "permission", "mongodb://localhost:27017")

	uow := dbFactory.Uow()
	routePermissionDb := dbFactory.Db(global.RoutePermission{}, uow)
	var routePermissions []global.RoutePermission
	if err := routePermissionDb.Query().ToArray(&routePermissions); err != nil {
		panic(err)
	}

	if len(routePermissions) == 0 {
		_ = routePermissionDb.Add(global.RoutePermission{
			Conditions: []global.ValueCondition{
				{
					Count:     12,
					Op:        relationaloperator.LT,
					ValueType: valuetype.Age,
				},
			},
			ID:     "enter",
			Routes: []string{"/mobile/enter", "/mobile/play", "/mobile/quick-play"},
			Response: map[string]interface{}{
				"err": 1,
				"msg": "小于12岁禁止进入",
			},
		})

		maxOnelineTime := int64(60 * 60 * 8)
		_ = routePermissionDb.Add(global.RoutePermission{
			Conditions: []global.ValueCondition{
				{
					Count:     18,
					Op:        relationaloperator.LE,
					ValueType: valuetype.Age,
				},
				{
					Count:     maxOnelineTime,
					Op:        relationaloperator.GE,
					ValueType: valuetype.OnlineTime,
				},
			},
			ID:     "play",
			Routes: []string{"/mobile/play"},
			Response: map[string]interface{}{
				"err": 1,
				"msg": "已超过累积游戏时间",
			},
		})

		_ = routePermissionDb.Add(global.RoutePermission{
			Conditions: []global.ValueCondition{
				{
					Count:     2,
					Op:        relationaloperator.LT,
					ValueType: valuetype.Vip,
				},
			},
			ID:     "quick-play",
			Routes: []string{"/mobile/quick-play"},
			Response: map[string]interface{}{
				"err": 1,
				"msg": "vip大于2级才能使用快捷功能",
			},
		})

		userDb := dbFactory.Db(global.UserValue{}, uow)

		_ = userDb.Add(global.UserValue{
			ID: "user-one",
			Values: map[valuetype.Value]int64{
				valuetype.Age:        17,
				valuetype.OnlineTime: maxOnelineTime,
			},
		})

		_ = userDb.Add(global.UserValue{
			ID: "user-two",
			Values: map[valuetype.Value]int64{
				valuetype.Age:        17,
				valuetype.OnlineTime: maxOnelineTime - 100,
			},
		})

		_ = userDb.Add(global.UserValue{
			ID: "user-three",
			Values: map[valuetype.Value]int64{
				valuetype.Age:        20,
				valuetype.OnlineTime: maxOnelineTime,
				valuetype.Vip:        1,
			},
		})

		_ = userDb.Add(global.UserValue{
			ID: "user-four",
			Values: map[valuetype.Value]int64{
				valuetype.Age:        20,
				valuetype.OnlineTime: maxOnelineTime,
				valuetype.Vip:        2,
			},
		})

		if err := uow.Commit(); err != nil {
			panic(err)
		}
	}

	lgginsvc.Listen(
		ginsvc.NewPostOption(apiFactory, dbFactory, func(uid string) contract.IUserValueService {
			return usersvc.NewValueService(dbFactory, uid)
		}),
		lgginsvc.NewPortOpion(10002),
	)
}
