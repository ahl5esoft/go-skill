package ossvc

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ahl5esoft/go-skill/hide-ctx/service/pathsvc"

	"github.com/stretchr/testify/assert"
)

func Test_ioFile_GetExt(t *testing.T) {
	ioPath := pathsvc.NewIOPath()
	res := NewIOFile(ioPath, "a.txt").GetExt()
	assert.Equal(t, res, ".txt")
}

func Test_ioFile_Write(t *testing.T) {
	t.Run("bytes.Buffer", func(t *testing.T) {
		wd, err := os.Getwd()
		assert.NoError(t, err)

		ioPath := pathsvc.NewIOPath()
		file := NewIOFile(ioPath, wd, "write-bytes.Buffer.txt")
		defer os.RemoveAll(
			file.GetPath(),
		)

		ioutil.WriteFile(
			file.GetPath(),
			[]byte("aa"),
			os.ModePerm,
		)

		var bf bytes.Buffer
		bf.WriteString("byte-slice")
		err = file.Write(bf)
		assert.NoError(t, err)

		f, err := file.(*ioFile).GetFile()
		assert.NoError(t, err)

		defer f.Close()

		res, err := ioutil.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(
			t,
			string(res),
			bf.String(),
		)
	})
}
