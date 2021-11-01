package mobile

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/ginsvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/goredissvc"
	"github.com/ahl5esoft/go-skill/hide-ctx/service/iocsvc"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_DemoApi_Call_It(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	redis := goredissvc.NewRedis(redis.Options{
		Addr: "127.0.0.1:6379",
	})
	iocsvc.Set(
		new(contract.IRedis),
		redis,
	)

	t.Run("第一次", func(t *testing.T) {
		defer redis.Del("test-once")

		gin.SetMode(gin.ReleaseMode)
		app := gin.New()

		mockApiFactory := contract.NewMockIApiFactory(ctrl)
		ginsvc.NewPostOption(mockApiFactory)(app)

		mockApiFactory.EXPECT().Build("mobile", "api").Return(
			new(DemoApi),
		)

		req := httptest.NewRequest(
			"POST",
			"/mobile/api",
			strings.NewReader(`{"key":"test-once","value":"a"}`),
		)
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)

		res, err := ioutil.ReadAll(
			resp.Result().Body,
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			string(res),
			`{"data":"a","err":0}`,
		)
	})

	t.Run("多次", func(t *testing.T) {
		defer redis.Del("test-multi")

		gin.SetMode(gin.ReleaseMode)
		app := gin.New()

		mockApiFactory := contract.NewMockIApiFactory(ctrl)
		ginsvc.NewPostOption(mockApiFactory)(app)

		mockApiFactory.EXPECT().Build("mobile", "api").Return(
			new(DemoApi),
		).AnyTimes()

		req := httptest.NewRequest(
			"POST",
			"/mobile/api",
			strings.NewReader(`{"key":"test-multi","value":"a"}`),
		)
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)

		req = httptest.NewRequest(
			"POST",
			"/mobile/api",
			strings.NewReader(`{"key":"test-multi","value":"b"}`),
		)
		resp = httptest.NewRecorder()
		app.ServeHTTP(resp, req)

		res, err := ioutil.ReadAll(
			resp.Result().Body,
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			string(res),
			`{"data":"a,b","err":0}`,
		)
	})
}
