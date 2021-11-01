package ossvc

import (
	"os"
	"testing"

	"github.com/ahl5esoft/go-skill/hide-ctx/service/pathsvc"

	"github.com/stretchr/testify/assert"
)

func Test_ioDirectory_FindDirectories(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		cwd, err := os.Getwd()
		assert.NoError(t, err)

		ioPath := pathsvc.NewIOPath()
		childDirPath := ioPath.Join(cwd, "dir")
		err = os.Mkdir(childDirPath, os.ModePerm)
		assert.NoError(t, err)

		defer os.Remove(childDirPath)

		res := NewIODirectory(ioPath, cwd).FindDirectories()
		assert.Len(t, res, 1)
	})

	t.Run("NotExists", func(t *testing.T) {
		ioPath := pathsvc.NewIOPath()
		res := NewIODirectory(ioPath, "a", "b", "c").FindDirectories()
		assert.Len(t, res, 0)
	})
}

func Test_ioDirectory_FindFiles(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	ioPath := pathsvc.NewIOPath()
	childDirPath := ioPath.Join(cwd, "files")
	err = os.Mkdir(childDirPath, os.ModePerm)
	assert.NoError(t, err)

	defer os.Remove(childDirPath)

	res := NewIODirectory(ioPath, childDirPath).FindFiles()
	assert.Len(t, res, 0)
}

func Test_ioDirectory_GetName(t *testing.T) {
	ioPath := pathsvc.NewIOPath()
	res := NewIODirectory(ioPath, "a", "b", "c").GetName()
	assert.Equal(t, res, "c")
}
