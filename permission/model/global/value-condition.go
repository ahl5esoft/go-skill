package global

import (
	relationaloperator "github.com/ahl5esoft/go-skill/permission/model/enum/relational-operator"
	valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"
)

type ValueCondition struct {
	Count     int64                    `json:"count"`
	Op        relationaloperator.Value `json:"op"`
	ValueType valuetype.Value          `json:"valueType"`
}
