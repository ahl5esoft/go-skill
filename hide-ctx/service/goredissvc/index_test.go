package goredissvc

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

var (
	clientCfg = redis.Options{
		Addr: "127.0.0.1:6379",
	}
	client = redis.NewClient(&clientCfg)
)

func Test_adapter_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		k := "test-get-ok"
		v := "ok"
		_, err := client.Set(ctx, k, v, 0).Result()
		assert.NoError(t, err)

		defer client.Del(ctx, k)

		res, err := NewRedis(clientCfg).Get(k)
		assert.NoError(t, err)
		assert.Equal(t, res, v)
	})

	t.Run("过期", func(t *testing.T) {
		k := "test-get-expire"
		v := "ok"
		_, err := client.Set(ctx, k, v, time.Second).Result()
		assert.NoError(t, err)

		defer client.Del(ctx, k)

		time.Sleep(time.Second)

		res, err := NewRedis(clientCfg).Get(k)
		assert.NoError(t, err)
		assert.Empty(t, res)
	})
}

func Test_adapter_Set(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		k := "test-set-ok"
		v := "ok"
		res, err := NewRedis(clientCfg).Set(k, v, time.Second)
		assert.NoError(t, err)
		assert.True(t, res)

		defer client.Del(ctx, k)

		res2, err := client.Get(ctx, k).Result()
		assert.NoError(t, err)
		assert.Equal(t, res2, v)
	})

	t.Run("过期", func(t *testing.T) {
		k := "test-set-expire"
		v := "ok"
		res, err := NewRedis(clientCfg).Set(k, v, time.Second)
		assert.NoError(t, err)
		assert.True(t, res)

		defer client.Del(ctx, k)

		time.Sleep(time.Second)

		res2, err := client.Exists(ctx, k).Result()
		assert.NoError(t, err)
		assert.Equal(
			t,
			res2,
			int64(0),
		)
	})
}
