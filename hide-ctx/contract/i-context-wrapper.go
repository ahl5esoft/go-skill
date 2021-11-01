package contract

import "context"

type IContextWrapper interface {
	WithContext(context.Context) interface{}
}
