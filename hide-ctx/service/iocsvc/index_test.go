package iocsvc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type iInterface interface {
	Test()
}

type derive struct{}

func (m derive) Test() {
	fmt.Println("set test")
}

type defaultTest struct {
	One iInterface `inject:""`
}

type composeTest struct {
	defaultTest

	Child iInterface `inject:""`
}

type composeInterfaceTest struct {
	iInterface

	Child iInterface `inject:""`
}

func Test_Get(t *testing.T) {
	t.Run("无效", func(t *testing.T) {
		ct := getInterfaceType(
			new(iInterface),
		)
		defer func() {
			rv := recover()
			assert.NotNil(t, rv)

			err, ok := rv.(error)
			assert.True(t, ok)
			assert.Equal(
				t,
				err,
				fmt.Errorf(invalidTypeFormat, ct),
			)
		}()

		Get(ct)
	})

	t.Run("默认", func(t *testing.T) {
		defer func() {
			assert.Nil(
				t,
				recover(),
			)
		}()

		ct := getInterfaceType(
			new(iInterface),
		)
		instanceValues[ct] = reflect.ValueOf(1)
		defer delete(instanceValues, ct)

		res := Get(ct)
		assert.Equal(t, res, 1)
	})
}

func Test_Inject(t *testing.T) {
	t.Run("默认", func(t *testing.T) {
		it := getInterfaceType(
			new(iInterface),
		)
		instanceValues[it] = reflect.ValueOf(
			new(derive),
		)

		var m defaultTest
		Inject(&m, func(v reflect.Value) reflect.Value {
			return v
		})

		assert.Equal(
			t,
			m.One,
			instanceValues[it].Interface(),
		)
	})

	t.Run("selectorFunc is nil", func(t *testing.T) {
		it := getInterfaceType(
			new(iInterface),
		)
		instanceValues[it] = reflect.ValueOf(
			new(derive),
		)

		var m defaultTest
		Inject(&m, nil)

		assert.Equal(
			t,
			m.One,
			instanceValues[it].Interface(),
		)
	})

	t.Run("继承", func(t *testing.T) {
		it := getInterfaceType(
			new(iInterface),
		)
		instanceValues[it] = reflect.ValueOf(
			new(derive),
		)

		var self composeTest
		Inject(&self, nil)

		assert.Equal(
			t,
			self.Child,
			instanceValues[it].Interface(),
		)
		assert.Equal(
			t,
			self.One,
			instanceValues[it].Interface(),
		)
	})

	t.Run("继承接口", func(t *testing.T) {
		it := getInterfaceType(
			new(iInterface),
		)
		instanceValues[it] = reflect.ValueOf(
			new(derive),
		)

		var self composeInterfaceTest
		Inject(&self, nil)

		assert.Equal(
			t,
			self.Child,
			instanceValues[it].Interface(),
		)
	})
}

func Test_Set(t *testing.T) {
	t.Run("默认", func(t *testing.T) {
		ct := reflect.TypeOf(
			new(iInterface),
		).Elem()
		defer delete(instanceValues, ct)

		Set(
			ct,
			new(derive),
		)

		_, ok := instanceValues[ct]
		assert.True(t, ok)
	})

	t.Run("非接口", func(t *testing.T) {
		it := reflect.TypeOf(1)
		defer func() {
			rv := recover()
			assert.NotNil(t, rv)

			err, ok := rv.(error)
			assert.True(t, ok)
			assert.Equal(
				t,
				err,
				fmt.Errorf(notInterfaceTypeFormat, it),
			)
		}()
		Set(
			it,
			new(derive),
		)
	})

	t.Run("没有继承", func(t *testing.T) {
		it := getInterfaceType(
			new(iInterface),
		)
		v := ""
		defer func() {
			rv := recover()
			assert.NotNil(t, rv)

			err, ok := rv.(error)
			assert.True(t, ok)
			assert.Equal(
				t,
				err,
				fmt.Errorf(notImplementsFormat, v, it),
			)
		}()
		Set(it, v)
	})
}
