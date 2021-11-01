# 前言

在互联网软件项目中,开发人员不仅要开发面向客户的移动、web`api接口`,还得开发面向运营或产品的后台`api接口`,因此在整个项目中,`api接口`的数据量是很多的,虽然现在市面上也有不少成型好用的框架,当然学习这些框架是需要成本的,因此今天就来分享如何让开发人员不需要关注框架即可开发`api接口`的经验.

# 接口

首先我们需要定义一个接口(编程定义),不管是什么样的`api接口`,传入的参数是不固定的,但是结果返回值却可以是固定的一个,这个固定的返回值可以是任何类型,因此我们用`interface{}`类型,另外一个返回值则是`error`,因此接口定义如下:
```go
package contract

type IApi interface {
	Call() (interface{}, error)
}
```

# 接口工厂

为了尽量让开发人员直接调用具体的api接口,因此我们需要一个工厂来帮助创建`contract.IApi`,由于`go`语言是将静态语言,无法像`C#``Java`或者动态语言那样直接通过名字字符串就就可以动态调用包内部的成员,因此工厂需要开放一个用来注册`contract.IApi`的方法,工厂实现如下:
```go
package apisvc

import (
	"reflect"

	"github.com/ahl5esoft/go-skill/one/contract"
	errorcode "github.com/ahl5esoft/go-skill/one/model/enum/error-code"
	"github.com/ahl5esoft/go-skill/one/service/errorsvc"
)

var (
	errNilApi = errorsvc.Newf(errorcode.API, "")
	nilApiPtr = &nilApi{}
)

type factory map[string]map[string]reflect.Type

func (m factory) Build(endpoint, api string) interface{} {
	if apiTypes, ok := m[endpoint]; ok {
		if apiType, ok := apiTypes[api]; ok {
			return reflect.New(apiType).Interface()
		}
	}

	return nilApiPtr
}

func (m factory) Register(endpoint, api string, apiInstance interface{}) {
	if _, ok := m[endpoint]; !ok {
		m[endpoint] = make(map[string]reflect.Type)
	}

	apiType := reflect.TypeOf(apiInstance)
	if apiType.Kind() == reflect.Ptr {
		apiType = apiType.Elem()
	}
	m[endpoint][api] = apiType
}

type nilApi struct{}

func (m nilApi) Call() (interface{}, error) {
	return nil, errNilApi
}

func NewFactory() contract.IApiFactory {
	return make(factory)
}
```
以上还使用了`Martin Fowler`在`<<重构>>`中提到的`NullObject模式`,当工厂创建的时候找不到注册的`contract.IApi`时则返回它,那么在业务中调用的时候,则不需要对返回的`contract.IApi`进行`nil`判断,直接使用即可.

# 代码生成

由于`go`语言的特性,因此我们只能在编写`api接口`的时候,通过主动编码的方式来注册api接口(比如`init`函数),但是其实这些代码都是相似的,因此我们可以定义一些规范,比如`api接口`存放的目录、`api接口`命名规范等,让开发人员遵守,然后我们在项目编译前,根据定义的规范扫描目录将所有符合条件的`api接口`生成代码并注册到`api工厂`中,大致代码如下:
```go
const (
	metadataTpl = `package api

import (
    "github.com/ahl5esoft/go-skill/easy-api/contract"
    {{- range .packages }}
    {{ .Name }} "{{ $.workspace }}/{{ join .RelativePathParts "/" }}"
    {{- end }}
)

func Register(apiFactory contract.IApiFactory) {
	{{- range $i, $r := .packages }}{{ range $ci, $cr := $r.Apis }}
    apiFactory.Register("{{ $r.Endpoint }}", "{{ $cr.Route }}", {{ $r.Name }}.{{ $cr.Struct }}Api{}){{ end }}{{ end }}
}`
	metadataFilename = "metadata.go"
)

var (
	regApi   = regexp.MustCompile(`type\s(\w+)Api`)
	tplFuncs = template.FuncMap{
		"join": func(elems []string, sep string) string {
			return strings.Join(elems, sep)
		},
	}
)

type apiData struct {
	Struct string
	Route  string
}

type packageData struct {
	Apis              []apiData
	Endpoint          string
	Name              string
	RelativePathParts []string
}

func GenerateMetadata(ioFactory contract.IIOFactory, ioPath contract.IIOPath, workspace string) (err error) {
	packages := make([]packageData, 0)
	apiDir := ioFactory.BuildDirectory(
		ioPath.GetRoot(),
		"api",
	)
	err = readGoFiles(apiDir, &packages, workspace)
	if err != nil {
		return
	}

	var tpl *template.Template
	if tpl, err = template.New("").Funcs(tplFuncs).Parse(metadataTpl); err != nil {
		return
	}

	var bf bytes.Buffer
	err = tpl.Execute(&bf, map[string]interface{}{
		"packages":  packages,
		"workspace": workspace,
	})
	if err != nil {
		return
	}

	err = ioFactory.BuildFile(
		apiDir.GetPath(),
		metadataFilename,
	).Write(bf)
	return
}

func readGoFiles(dir contract.IIODirectory, packages *[]packageData, workspace string) (err error) {
	files := dir.FindFiles()

	apis := make([]apiData, 0)
	for _, r := range files {
		if r.GetExt() != ".go" || r.GetName() == metadataFilename {
			continue
		}

		isTest := strings.Contains(
			r.GetName(),
			"_test",
		)
		if isTest {
			continue
		}

		api := apiData{
			Route: strings.Replace(
				r.GetName(),
				r.GetExt(),
				"",
				1,
			),
		}

		var text string
		if err = r.Read(&text); err != nil {
			return
		}

		matches := regApi.FindStringSubmatch(text)
		if len(matches) == 0 {
			continue
		}

		api.Struct = matches[1]
		apis = append(apis, api)
	}

	if len(apis) > 0 {
		pkg := packageData{
			Apis:              apis,
			RelativePathParts: make([]string, 0),
		}
		var temp contract.IIODirectory
		for {
			if len(pkg.RelativePathParts) == 0 {
				temp = dir
			} else {
				temp = temp.GetParent()
			}

			if temp.GetName() == workspace {
				break
			}

			pkg.RelativePathParts = append([]string{
				temp.GetName(),
			}, pkg.RelativePathParts...)
		}

		if pkg.RelativePathParts[len(pkg.RelativePathParts)-2] == "api" {
			pkg.Endpoint = pkg.RelativePathParts[len(pkg.RelativePathParts)-1]
			pkg.Name = pkg.RelativePathParts[len(pkg.RelativePathParts)-1]
		} else {
			pkg.Endpoint = strings.Join(
				pkg.RelativePathParts[len(pkg.RelativePathParts)-2:],
				"/",
			)
			pkg.Name = strings.Join(
				pkg.RelativePathParts[len(pkg.RelativePathParts)-2:],
				"",
			)
		}
		pkg.Name = strings.Replace(pkg.Name, "-", "", -1)

		*packages = append(*packages, pkg)
	}

	childDirs := dir.FindDirectories()
	if len(childDirs) == 0 {
		return
	}

	for _, r := range childDirs {
		readGoFiles(r, packages, workspace)
	}
	return
}
```

# gin

接口工厂、api接口元数据文件都准备好了,那么就只剩下请求入口了,这里用`gin`来实现,以post为例,大致代码如下:
```go
gin.SetMode(gin.ReleaseMode)
app := gin.New()
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
```

# 结束语

以上就是基本的实现思路,如果有任何疑问或者优化方案欢迎告诉我,谢谢.