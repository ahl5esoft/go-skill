package contract

import "time"

type IRedis interface {
	Del(...string) (int64, error)
	Get(k string) (string, error)
	Set(k, v string, expires time.Duration) (bool, error)
}
