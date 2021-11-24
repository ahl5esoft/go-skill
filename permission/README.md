# 前言

随着移动应用的普及,不管应用是否免费,只要应用内部包含购买功能的,货多多少都需要对应用的功能做一些限制,好能更好的区分付费用户和非付费用户的区别,或者由于国家出台的一系列政策,对于青少年在使用应用方面也会提出一些要求,因此功能限制也算是基础的功能了,今天将以围绕数值系统的方式来构建功能限制(功能权限).

# 流程

比如某些功能只能对大于12岁的人开放,那么这个`12`其实是一个数值,然后要访问这个功能的时候,需要通过网络的方式访问服务端,服务端通过会话或者令牌获取该用户的数值数据,从数值数据中获取年龄,如果年龄小于等于`12`则返回错误响应数据给客户端,客户端收到响应数据之后弹出对应的提示语告知用户.

# 数值类型

这里用一个数值枚举来存储用所有的数值类型, 定义如下:
```go
package valuetype

type Value int

const (
	Age        Value = 1 // Age is 年龄
)
```

# 用户数值模型

由于数值类型会随着应用的开发而增加,为了防止字段的增加而导致频繁修改结构,因此我们可以用`map[valuetype.Value]int64`来存储,结构如下:
```go
type UserValue struct {
	ID     string
	Values map[valuetype.Value]int64
}
```

# 条件

根据以上的业务需求,需要存储数值类型\数量以及运算符,如果条件不满足的情况下需要下发数据给客户端,为了让数据格式更灵活,那么可以用`map[string]interface{}`,其次就是定义哪些路由会触发条件判断,因此模型如下:
```go
package relationaloperator

type Value int

const (
	EQ Value = 0 // EQ is 等于
	GE Value = 1 // GE is 大于等于
	GT Value = 2 // GT is 大于
	LE Value = 3 // LE is 小于等于
	LT Value = 4 // LT is 小于
)

// 数值条件
type ValueCondition struct {
	Count     int64
	Op        relationaloperator.Value
	ValueType valuetype.Value
}

// 路由权限
type RoutePermission struct {
	Conditions []ValueCondition
	ID         string
	Response   map[string]interface{}
	Routes     []string
}
```
`Conditions`是以`and`方式运算的,如果要满足`or`则需要同一个路由定义多个`RoutePermission`即可

# 数据

有了以上接口,就可以定义业务需要的数据,数据如下:
```go
RoutePermission{
	Conditions: []ValueCondition{
		{
			Count:     12,
			Op:        relationaloperator.LT,
			ValueType: valuetype.Age,
		},
	},
	ID:     "id",
	Routes: []string{"路由"},
	Response: map[string]interface{}{
		"err": 1,
		"msg": "小于12岁禁止进入",
	},
}
```

# 用户数值接口

由于需要通过数值类型获取用户的数值,如果用户数值数据不存在或者error则应该返回默认值(0),否则则返回具体的值,如果直接在使用的时候获取用户数据那么需要频繁判断,由于条件判断也需要用到用户的数值,因此可以将判断条件的逻辑也放在该接口内,代码如下:
```go
// 接口
type IUserValueService interface {
	MustCheckConditions([]ValueCondition) bool
	MustGetCount(valueType valuetype.Value) int64
}

// 实现
type valueService struct {
	dbFactory IDbFactory
	rows      []UserValue
	uid       string
}

func (m *valueService) MustCheckConditions(conditions []ValueCondition) bool {
	return underscore.Chain(conditions).Any(func(r ValueCondition, _ int) bool {
		count := m.MustGetCount(r.ValueType)
		ok := (r.Op == relationaloperator.EQ && count == r.Count) ||
			(r.Op == relationaloperator.LE && count <= r.Count) ||
			(r.Op == relationaloperator.LT && count < r.Count) ||
			(r.Op == relationaloperator.GE && count >= r.Count) ||
			(r.Op == relationaloperator.GT && count > r.Count)
		return !ok
	})
}

func (m *valueService) MustGetCount(valueType valuetype.Value) int64 {
	entry, err := m.getEntry()
	if err != nil || entry == nil {
		return 0
	}

	if v, ok := entry.Values[valueType]; ok {
		return v
	}

	return 0
}

func (m *valueService) getEntry() (*UserValue, error) {
	// 延迟加载
	if m.rows == nil {
		err := m.dbFactory.Db(UserValue{}).Query().Where(bson.M{
			"_id": m.uid,
		}).ToArray(&(m.rows))
		if err != nil {
			return nil, err
		}
	}

	if len(m.rows) > 0 {
		return &m.rows[0], nil
	}

	return nil, nil
}

func NewValueService(
	dbFactory IDbFactory,
	uid string,
) contract.IUserValueService {
	return &valueService{
		dbFactory: dbFactory,
		uid:       uid,
	}
}
```

# 触发

接下来只需要在请求的时候,根据路由获取权限,如果权限数据存在则判断,只要其中一条权限不满足则下发`Response`,否则则调用api,并下发api结果.这里以`gin`为例,代码如下:
```go
app.POST("/:endpoint/:api", func(ctx *gin.Context) {
	// 其他略
	
	uid := ctx.GetHeader("uid")
	if uid == "" {
		err = errors.New("未认证")
		return
	}
	
	var routePermissions []RoutePermission
	err = dbFactory.Db(RoutePermission{}).Query().Where(bson.M{
		"Routes": ctx.Request.URL.String(),
	}).ToArray(&routePermissions)
	if err != nil {
		return
	}
	
	if len(routePermissions) > 0 {
		userValueService := NewValueService(dbFactory, uid)
		for _, r := range routePermissions {
			if ok := userValueService.MustCheckConditions(r.Conditions); !ok {
				err = fmt.Errorf("%+v", r.Response)
				return
			}
		}
	}
})
```

# 示例

以下提供2个示例:
1. 国家政策要求一定年龄的小孩只能玩一定时间
1. vip到达一定等级之后才能使用快捷功能
```go
// 年龄小于等于18岁只能在线8小时
RoutePermission{
	Conditions: []ValueCondition{
		{
			Count:     18,
			Op:        relationaloperator.LE,
			ValueType: valuetype.Age,
		},
		{
			Count:     60 * 60 * 8,
			Op:        relationaloperator.GE,
			ValueType: valuetype.OnlineTime,
		},
	},
	ID:     "id-online-time",
	Routes: []string{"路由"},
	Response: map[string]interface{}{
		"err": 1,
		"msg": "已超过累积游戏时间",
	},
}

// vip等级大于2才能使用快捷功能
RoutePermission{
	Conditions: []ValueCondition{
		{
			Count:     2,
			Op:        relationaloperator.LT,
			ValueType: valuetype.Vip,
		},
	},
	ID:     "id-quick-play",
	Routes: []string{"路由"},
	Response: map[string]interface{}{
		"err": 1,
		"msg": "vip大于2级才能使用快捷功能",
	},
}
```

# 结尾

那么今天的文章到这里就结束了,数值模块的实现可以[参考](https://my.oschina.net/ahl5esoft/blog/4301907 "参考")而`IDbFactory`则可以[参考](https://my.oschina.net/ahl5esoft/blog/4911231 "参考").

如果文章有任何问题或者对以上内容有任何疑问,请留言告诉我,谢谢.