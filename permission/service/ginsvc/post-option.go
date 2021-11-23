package ginsvc

import (
	"errors"

	"github.com/ahl5esoft/go-skill/permission/contract"
	"github.com/ahl5esoft/go-skill/permission/model/global"
	lgcontract "github.com/ahl5esoft/lite-go/contract"
	contextkey "github.com/ahl5esoft/lite-go/model/enum/context-key"
	"github.com/ahl5esoft/lite-go/service/errorsvc"
	"github.com/ahl5esoft/lite-go/service/ginsvc"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/bson"
)

func NewPostOption(
	apiFactory lgcontract.IApiFactory,
	dbFactory lgcontract.IDbFactory,
	getUserValueService func(uid string) contract.IUserValueService,
) ginsvc.Option {
	return ginsvc.NewPostOption("/:endpoint/:api", func(ctx *gin.Context) (lgcontract.IApi, error) {
		uid := ctx.GetHeader("uid")
		if uid == "" {
			return nil, errors.New("未认证")
		}

		var routePermissions []global.RoutePermission
		err := dbFactory.Db(global.RoutePermission{}).Query().Where(bson.M{
			"Routes": ctx.Request.URL.String(),
		}).ToArray(&routePermissions)
		if err != nil {
			return nil, err
		}

		if len(routePermissions) > 0 {
			userValueService := getUserValueService(uid)
			for _, r := range routePermissions {
				if ok := userValueService.MustCheckConditions(r.Conditions); !ok {
					return nil, errorsvc.New(1, r.Response)
				}
			}
		}

		api := apiFactory.Build(
			ctx.Param("endpoint"),
			ctx.Param("api"),
		)

		httpBody, hasBody := ctx.Get(contextkey.HttpBody)
		if hasBody {
			jsoniter.Unmarshal(
				httpBody.([]byte),
				api,
			)
		}

		return api, nil
	})
}
