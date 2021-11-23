package usersvc

import (
	"github.com/ahl5esoft/go-skill/permission/contract"
	relationaloperator "github.com/ahl5esoft/go-skill/permission/model/enum/relational-operator"
	valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"
	"github.com/ahl5esoft/go-skill/permission/model/global"
	underscore "github.com/ahl5esoft/golang-underscore"
	lgcontract "github.com/ahl5esoft/lite-go/contract"

	"go.mongodb.org/mongo-driver/bson"
)

type valueService struct {
	dbFactory lgcontract.IDbFactory
	rows      []global.UserValue
	uid       string
}

func (m *valueService) MustCheckConditions(conditions []global.ValueCondition) bool {
	return underscore.Chain(conditions).Any(func(r global.ValueCondition, _ int) bool {
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

func (m *valueService) getEntry() (*global.UserValue, error) {
	if m.rows == nil {
		err := m.dbFactory.Db(global.UserValue{}).Query().Where(bson.M{
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
	dbFactory lgcontract.IDbFactory,
	uid string,
) contract.IUserValueService {
	return &valueService{
		dbFactory: dbFactory,
		uid:       uid,
	}
}
