package mobile

import (
	"time"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"
)

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
