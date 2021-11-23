package contract

import (
	valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"
	"github.com/ahl5esoft/go-skill/permission/model/global"
)

type IUserValueService interface {
	MustCheckConditions([]global.ValueCondition) bool
	MustGetCount(valueType valuetype.Value) int64
}
