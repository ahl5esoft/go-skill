package ginsvc

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ahl5esoft/go-skill/easy-api/contract"
	errorcode "github.com/ahl5esoft/go-skill/easy-api/model/enum/error-code"
	"github.com/ahl5esoft/go-skill/easy-api/model/response"
	"github.com/ahl5esoft/go-skill/easy-api/service/errorsvc"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

func NewPostOption(apiFactory contract.IApiFactory) Option {
	return func(app *gin.Engine) {
		validate := validator.New()
		app.POST("/:endpoint/:api", func(ctx *gin.Context) {
			var resp response.Api
			defer func() {
				ctx.JSON(http.StatusOK, resp)
			}()

			var err error
			defer func() {
				if rv := recover(); rv != nil {
					var ok bool
					if err, ok = rv.(error); !ok {
						err = fmt.Errorf("%v", rv)
					}
				}

				if err != nil {
					if cErr, ok := err.(contract.IError); ok {
						resp.Error = cErr.GetCode()
						if cErr.GetData() != nil {
							resp.Data = cErr.GetData()
						} else {
							resp.Data = cErr.Error()
						}
					} else {
						resp.Error = errorcode.Panic
					}
				}
			}()

			api := apiFactory.Build(
				ctx.Param("endpoint"),
				ctx.Param("api"),
			)

			if ctx.Request.ContentLength > 0 {
				var bodyBytes []byte
				if bodyBytes, err = ioutil.ReadAll(ctx.Request.Body); err != nil {
					return
				}

				ctx.Set("body", bodyBytes)
				if err = jsoniter.Unmarshal(bodyBytes, &api); err != nil {
					return
				}

				if err = validate.Struct(api); err != nil {
					err = errorsvc.Newf(errorcode.Verify, "")
					return
				}
			}

			resp.Data, err = api.Call()
		})
	}
}
