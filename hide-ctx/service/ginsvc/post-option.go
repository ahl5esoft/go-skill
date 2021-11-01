package ginsvc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
	errorcode "github.com/ahl5esoft/go-skill/hide-ctx/model/enum/error-code"
	"github.com/ahl5esoft/go-skill/hide-ctx/model/response"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/errorsvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/iocsvc"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
						resp.Data = err.Error()
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

				if err = jsoniter.Unmarshal(bodyBytes, &api); err != nil {
					return
				}

				if err = validate.Struct(api); err != nil {
					err = errorsvc.Newf(errorcode.Verify, "")
					return
				}
			}

			iocsvc.Inject(api, func(v reflect.Value) reflect.Value {
				if w, ok := v.Interface().(contract.IContextWrapper); ok {
					return reflect.ValueOf(
						w.WithContext(ctx),
					)
				}
				return v
			})

			resp.Data, err = api.Call()
		})
	}
}
