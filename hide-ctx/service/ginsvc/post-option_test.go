package ginsvc

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func Test_NewPostOption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gin.SetMode(gin.TestMode)
	app := gin.New()

	mockApiFactory := contract.NewMockIApiFactory(ctrl)

	mockApi := contract.NewMockIApi(ctrl)
	mockApiFactory.EXPECT().Build("mobile", "api").Return(mockApi)

	mockApi.EXPECT().Call().Return("ok", nil)

	NewPostOption(mockApiFactory)(app)

	req := httptest.NewRequest(
		"POST",
		"/mobile/api",
		strings.NewReader(``),
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
		`{"data":"ok","err":0}`,
	)
}
