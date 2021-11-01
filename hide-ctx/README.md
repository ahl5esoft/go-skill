# 前言

当我们使用`go`进行开发的时候,由于`context.Context`的关键作用,大部分第三方库都需要传入`context.Context`(全链路跟踪、数据库、redis、MQ等),当业务比较复杂涉及的模块比较多的时候,就需要将`context.Context`一层层传递过去,那么今天就来分享如何将`context.Context`隐藏起来从而减少开发人员开发的难度(毕竟对于开发人员而言他们并不关注`context.Context`的作用而且用到它的机会也很少).本文将以`redis`使用为例并在[简化api接口开发(基于gin实现)](https://my.oschina.net/ahl5esoft/blog/5291633 "简化api接口开发(基于gin实现)")基础上进行扩展.

# redis

由于示例的关系这里就只定义了`get`和`set`的接口,代码如下:
```go
type IRedis interface {
	Get(k string) (string, error)
	Set(k, v string, expires time.Duration) (bool, error)
}
```
然后使用`github.com/go-redis/redis`实现接口,代码如下:
```go
var adapterMutex sync.Mutex

type adapter struct {
	cfg    *redis.Options
	client redis.Cmdable
	ctx    context.Context
}

func (m *adapter) Get(k string) (string, error) {
	res, err := m.getClient().Get(m.ctx, k).Result()
	if err != nil && err == redis.Nil {
		return "", nil
	}

	return res, err
}

func (m *adapter) Set(k, v string, expires time.Duration) (bool, error) {
	res, err := m.getClient().Set(m.ctx, k, v, expires).Result()
	if err != nil {
		return false, err
	}

	return res == "OK", nil
}

func (m *adapter) getClient() redis.Cmdable {
	if m.client == nil {
		adapterMutex.Lock()
		defer adapterMutex.Unlock()

		if m.client == nil {
			m.client = redis.NewClient(m.cfg)
		}
	}
	return m.client
}

func NewRedis(cfg redis.Options) contract.IRedis {
	return &adapter{
		ctx: context.Background(),
		cfg: &cfg,
	}
}
```

# 注入

有了`redis`的接口和实现类之后就可以通过依赖注入(参考[go - 依赖注入](https://my.oschina.net/ahl5esoft/blog/4922378 "go - 依赖注入")或使用第三方库)的方式将实例注入到具体的业务中去,但是这里有一个问题,由于`redis`实例是一个单例,而每次请求的`context.Context`都是不一样的,因此需要生成新的`redis`实例,由于本文中只以`redis`为例,而实际业务中还会引用其他的需要`context.Context`的模块,因此需要定义一个重新返回实例的接口,接口如下:
```go
type IContextWrapper interface {
	WithContext(context.Context) interface{}
}
```
接下来只需要调用`api接口`的方法前加入注入代码即可,大致代码如下:
```go
func NewPostOption(apiFactory contract.IApiFactory) Option {
	return func(app *gin.Engine) {
		// 略
		app.POST("/:endpoint/:api", func(ctx *gin.Context) {
			// 略

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
```

# api接口

有了以上的准备之后,开发人员只需要实现`contract.IApi`并定义需要注入的字段即可,示例接口代码如下:
```go
type DemoApi struct {
	Redis contract.IRedis `inject:""`

	Key   string
	Value string
}

func (m DemoApi) Call() (res interface{}, err error) {
	var s string
	if s, err = m.Redis.Get(m.Key); err != nil {
		return
	}

	if s != "" {
		s += ","
	}

	if _, err = m.Redis.Set(m.Key, s+m.Value, 20*time.Second); err != nil {
		return
	}

	res = s + m.Value
	return
}
```

# 结束语

今天就到这里了,如果有任何问题或者不对的地方请大家告诉我,谢谢.