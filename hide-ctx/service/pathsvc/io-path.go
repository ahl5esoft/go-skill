package pathsvc

import (
	"path/filepath"

	"github.com/ahl5esoft/go-skill/hide-ctx/contract"

	underscore "github.com/ahl5esoft/golang-underscore"
)

type ioPath struct {
	root string
}

func (m ioPath) GetRoot() string {
	return m.root
}

func (m ioPath) Join(paths ...string) string {
	var res string
	underscore.Chain(paths).Aggregate(func(memo string, r string, _ int) string {
		if memo == "" {
			return r
		}

		if r == ".." {
			return filepath.Dir(memo)
		}

		return filepath.Join(memo, r)
	}, "").Value(&res)
	return res
}

// NewIOPath is 路径实例
func NewIOPath(rootArgs ...string) contract.IIOPath {
	p := new(ioPath)
	p.root = p.Join(rootArgs...)
	return p
}
