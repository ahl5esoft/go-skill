package ossvc

import (
	"io/ioutil"
	"os"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"

	underscore "github.com/ahl5esoft/golang-underscore"
)

type ioDirectory struct {
	contract.IIONode

	ioPath contract.IIOPath
}

func (m ioDirectory) Create() error {
	if m.IsExist() {
		return nil
	}

	dirPath := m.GetPath()
	return os.MkdirAll(dirPath, os.ModePerm)
}

func (m ioDirectory) FindDirectories() []contract.IIODirectory {
	children := make([]contract.IIODirectory, 0)
	return m.findNodes(children, func(r os.FileInfo, filePath string) interface{} {
		if r.IsDir() {
			children = append(
				children,
				NewIODirectory(m.ioPath, filePath),
			)
		}
		return children
	}).([]contract.IIODirectory)
}

func (m ioDirectory) FindFiles() []contract.IIOFile {
	children := make([]contract.IIOFile, 0)
	return m.findNodes(children, func(r os.FileInfo, filePath string) interface{} {
		if !r.IsDir() {
			children = append(
				children,
				NewIOFile(m.ioPath, filePath),
			)
		}
		return children
	}).([]contract.IIOFile)
}

func (m ioDirectory) findNodes(defaultValue interface{}, handleNodeFunc func(r os.FileInfo, nodePath string) interface{}) interface{} {
	dirPath := m.GetPath()
	nodes, err := ioutil.ReadDir(dirPath)
	if err != nil || len(nodes) == 0 {
		return defaultValue
	}

	var res interface{}
	underscore.Chain(nodes).Each(func(r os.FileInfo, _ int) {
		nodePath := m.ioPath.Join(
			dirPath,
			r.Name(),
		)
		res = handleNodeFunc(r, nodePath)
	})
	return res
}

func NewIODirectory(ioPath contract.IIOPath, paths ...string) contract.IIODirectory {
	return &ioDirectory{
		IIONode: newIONode(ioPath, paths...),
		ioPath:  ioPath,
	}
}
