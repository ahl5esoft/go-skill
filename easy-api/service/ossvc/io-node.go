package ossvc

import (
	"os"
	"path/filepath"

	"github.com/ahl5esoft/go-skill/easy-api/contract"
)

type ioNode struct {
	ioPath contract.IIOPath
	path   string
}

func (m ioNode) GetName() string {
	return filepath.Base(m.path)
}

func (m ioNode) GetParent() contract.IIODirectory {
	return NewIODirectory(
		m.ioPath,
		m.GetPath(),
		"..",
	)
}

func (m ioNode) GetPath() string {
	return m.path
}

func (m ioNode) IsExist() bool {
	_, err := os.Stat(m.path)
	return err == nil || os.IsExist(err)
}

func newIONode(ioPath contract.IIOPath, paths ...string) contract.IIONode {
	return ioNode{
		ioPath: ioPath,
		path:   ioPath.Join(paths...),
	}
}
