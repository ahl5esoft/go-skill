package usersvc

import (
	"testing"

	relationaloperator "github.com/ahl5esoft/go-skill/permission/model/enum/relational-operator"
	valuetype "github.com/ahl5esoft/go-skill/permission/model/enum/value-type"
	"github.com/ahl5esoft/go-skill/permission/model/global"

	"github.com/ahl5esoft/lite-go/contract"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_valueService_MustCheckConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockDbFactory := contract.NewMockIDbFactory(ctrl)
		uid := "user-id"
		self := NewValueService(mockDbFactory, uid)

		mockDbRepo := contract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.UserValue{}).Return(mockDbRepo)

		mockDbQuery := contract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		mockDbQuery.EXPECT().Where(bson.M{
			"_id": uid,
		}).Return(mockDbQuery)

		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		).SetArg(0, []global.UserValue{
			{
				Values: map[valuetype.Value]int64{
					1: 11,
					2: 22,
				},
			},
		})

		res := self.MustCheckConditions([]global.ValueCondition{
			{
				Count:     11,
				Op:        relationaloperator.EQ,
				ValueType: 1,
			},
			{
				Count:     21,
				Op:        relationaloperator.GT,
				ValueType: 2,
			},
		})
		assert.True(t, res)
	})

	t.Run("global.UserValue不存在", func(t *testing.T) {
		mockDbFactory := contract.NewMockIDbFactory(ctrl)
		uid := "user-id"
		self := NewValueService(mockDbFactory, uid)

		mockDbRepo := contract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.UserValue{}).Return(mockDbRepo)

		mockDbQuery := contract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		mockDbQuery.EXPECT().Where(bson.M{
			"_id": uid,
		}).Return(mockDbQuery)

		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		)

		res := self.MustCheckConditions([]global.ValueCondition{
			{
				Count:     11,
				Op:        relationaloperator.EQ,
				ValueType: 1,
			},
			{
				Count:     21,
				Op:        relationaloperator.GT,
				ValueType: 2,
			},
		})
		assert.False(t, res)
	})
}

func Test_valueService_MustGetCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("ok", func(t *testing.T) {
		mockDbFactory := contract.NewMockIDbFactory(ctrl)
		uid := "user-id"
		self := NewValueService(mockDbFactory, uid)

		mockDbRepo := contract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.UserValue{}).Return(mockDbRepo)

		mockDbQuery := contract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		mockDbQuery.EXPECT().Where(bson.M{
			"_id": uid,
		}).Return(mockDbQuery)

		valueType := valuetype.Value(1)
		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		).SetArg(0, []global.UserValue{
			{
				Values: map[valuetype.Value]int64{
					valueType: 11,
				},
			},
		})

		res := self.MustGetCount(valueType)
		assert.Equal(
			t,
			res,
			int64(11),
		)
	})

	t.Run("global.UserValue不存在", func(t *testing.T) {
		mockDbFactory := contract.NewMockIDbFactory(ctrl)
		uid := "user-id"
		self := NewValueService(mockDbFactory, uid)

		mockDbRepo := contract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.UserValue{}).Return(mockDbRepo)

		mockDbQuery := contract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		mockDbQuery.EXPECT().Where(bson.M{
			"_id": uid,
		}).Return(mockDbQuery)

		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		)

		valueType := valuetype.Value(1)
		res := self.MustGetCount(valueType)
		assert.Equal(
			t,
			res,
			int64(0),
		)
	})

	t.Run("global.UserValue.Values[valuetype.Value]不存在", func(t *testing.T) {
		mockDbFactory := contract.NewMockIDbFactory(ctrl)
		uid := "user-id"
		self := NewValueService(mockDbFactory, uid)

		mockDbRepo := contract.NewMockIDbRepository(ctrl)
		mockDbFactory.EXPECT().Db(global.UserValue{}).Return(mockDbRepo)

		mockDbQuery := contract.NewMockIDbQuery(ctrl)
		mockDbRepo.EXPECT().Query().Return(mockDbQuery)

		mockDbQuery.EXPECT().Where(bson.M{
			"_id": uid,
		}).Return(mockDbQuery)

		mockDbQuery.EXPECT().ToArray(
			gomock.Any(),
		).SetArg(0, []global.UserValue{
			{
				Values: make(map[valuetype.Value]int64),
			},
		})

		valueType := valuetype.Value(1)
		res := self.MustGetCount(valueType)
		assert.Equal(
			t,
			res,
			int64(0),
		)
	})
}
