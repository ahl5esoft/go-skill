package ossvc

import (
	"os"
	"testing"

	"github.com/ahl5esoft/go-skill/easy-api/service/pathsvc"

	"github.com/stretchr/testify/assert"
)

func Test_ioNode_GetName(t *testing.T) {
	t.Run("dir", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, "a", "b", "c").GetName()
		assert.Equal(t, res, "c")
	})

	t.Run("file", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, "a", "b", "c.txt").GetName()
		assert.Equal(t, res, "c.txt")
	})
}

func Test_ioNode_GetParent(t *testing.T) {
	t.Run("dir", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, "a", "b").GetParent()
		assert.Equal(
			t,
			res.GetPath(),
			"a",
		)
	})

	t.Run("file", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, "a", "b.txt").GetParent()
		assert.Equal(
			t,
			res.GetPath(),
			"a",
		)
	})
}

func Test_ioNode_IsExist_F(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, "a", "b", "c").IsExist()
		assert.False(t, res)
	})

	t.Run("false", func(t *testing.T) {
		wd, err := os.Getwd()
		assert.NoError(t, err)

		ioPath := pathsvc.NewIOPath()
		res := newIONode(ioPath, wd, "io-node.go").IsExist()
		assert.True(t, res)
	})
}
