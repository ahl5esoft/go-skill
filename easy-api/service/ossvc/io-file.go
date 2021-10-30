package ossvc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ahl5esoft/go-skill/easy-api/contract"
)

type ioFile struct {
	contract.IIONode
}

func (m ioFile) GetExt() string {
	filePath := m.GetPath()
	return filepath.Ext(filePath)
}

func (m ioFile) GetFile() (*os.File, error) {
	var file *os.File
	var err error
	filePath := m.GetPath()
	if m.IsExist() {
		file, err = os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
	} else {
		file, err = os.Create(filePath)
	}
	return file, err
}

func (m ioFile) Read(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("ossvc.ioFile.Read: v必须为指针")
	}

	f, err := m.GetFile()
	if err != nil {
		return err
	}

	defer f.Close()

	bf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	value = value.Elem()
	if value.Kind() == reflect.String {
		value.SetString(
			string(bf),
		)
		return nil
	} else if value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Uint8 {
		value.SetBytes(bf)
		return nil
	}

	return fmt.Errorf(
		"不支持ossvc.ioFile.Read(%s)",
		value.Type(),
	)
}

func (m ioFile) Write(data interface{}) error {
	if bf, ok := data.(bytes.Buffer); ok {
		return m.writeBytes(
			bf.Bytes(),
		)
	}

	return fmt.Errorf(
		"ossvc.ioFile.Write暂不支持%s",
		reflect.TypeOf(data).Name(),
	)
}

func (m ioFile) writeBytes(bf []byte) error {
	file, err := m.GetFile()
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(bf)
	return err
}

func NewIOFile(ioPath contract.IIOPath, paths ...string) contract.IIOFile {
	return &ioFile{
		IIONode: newIONode(ioPath, paths...),
	}
}
