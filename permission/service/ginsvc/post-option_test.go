package ginsvc

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahl5esoft/go-skill/permission/contract"
	relationaloperator "github.com/ahl5esoft/go-skill/permission/model/enum/relational-operator"
	"github.com/ahl5esoft/go-skill/permission/model/global"
	lgcontract "github.com/ahl5esoft/lite-go/contract"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_NewPostOption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		app := gin.New()

		mockApiFactory := lgcontract.NewMockIApiFactory(ctrl)
		mockDbFactory := lgcontract.NewMockIDbFactory(ctrl)
		mockUserValueService := contract.NewMockIUserValueService(ctrl)
		uid := "user-id"
		NewPostOption(mockApiFactory, mockDbFactory, func(id string) contract.IUserValueService {
			assert.Equal(t, uid, id)
			return mockUserValueService
		})(app)

		mockApi := lgcontract.NewMockIApi(ctrl)
		mockApiFactory.EXPECT().Build("mobile", "api").Return(mockApi)

		mockDbRepo := lgcontract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.RoutePermission{}).Return(mockDbRepo)

		mockDbQuery := lgcontract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		route := "/mobile/api"
		mockDbQuery.EXPECT().Where(bson.M{
			"Routes": bson.M{
				"$in": route,
			},
		}).Return(mockDbQuery)

		routePermission := global.RoutePermission{
			Conditions: []global.ValueCondition{
				{
					Count:     11,
					Op:        relationaloperator.GT,
					ValueType: 1,
				},
			},
			Routes: []string{route},
		}
		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		).SetArg(0, []global.RoutePermission{routePermission})

		mockUserValueService.EXPECT().MustCheckConditions(routePermission.Conditions).Return(true)

		mockApi.EXPECT().Call().Return("ok", nil)

		req := httptest.NewRequest(
			"POST",
			route,
			strings.NewReader(``),
		)
		req.Header.Add("uid", uid)
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)

		res, err := ioutil.ReadAll(
			resp.Result().Body,
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			string(res),
			`{"data":"ok","err":0}`,
		)
	})

	t.Run("权限不足", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		app := gin.New()

		mockApiFactory := lgcontract.NewMockIApiFactory(ctrl)
		mockDbFactory := lgcontract.NewMockIDbFactory(ctrl)
		mockUserValueService := contract.NewMockIUserValueService(ctrl)
		uid := "user-id"
		NewPostOption(mockApiFactory, mockDbFactory, func(id string) contract.IUserValueService {
			assert.Equal(t, uid, id)
			return mockUserValueService
		})(app)

		mockApi := lgcontract.NewMockIApi(ctrl)
		mockApiFactory.EXPECT().Build("mobile", "api").Return(mockApi)

		mockDbRepo := lgcontract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.RoutePermission{}).Return(mockDbRepo)

		mockDbQuery := lgcontract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		route := "/mobile/api"
		mockDbQuery.EXPECT().Where(bson.M{
			"Routes": bson.M{
				"$in": route,
			},
		}).Return(mockDbQuery)

		routePermission := global.RoutePermission{
			Conditions: []global.ValueCondition{
				{
					Count:     11,
					Op:        relationaloperator.GT,
					ValueType: 1,
				},
			},
			Response: make(map[string]interface{}),
			Routes:   []string{route},
		}
		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		).SetArg(0, []global.RoutePermission{routePermission})

		mockUserValueService.EXPECT().MustCheckConditions(routePermission.Conditions).Return(false)

		req := httptest.NewRequest(
			"POST",
			route,
			strings.NewReader(``),
		)
		req.Header.Add("uid", uid)
		resp := httptest.NewRecorder()
		app.ServeHTTP(resp, req)

		res, err := ioutil.ReadAll(
			resp.Result().Body,
		)
		assert.NoError(t, err)
		assert.Equal(
			t,
			string(res),
			`{"data":{},"err":1}`,
		)
	})
}
