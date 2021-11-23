package global

import valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"

type UserValue struct {
	ID     string `alias:"" bson:"_id" db:"_id"`
	Values map[valuetype.Value]int64
}

func (m UserValue) GetID() string {
	return m.ID
}
