package ossvc

import "github.com/ahl5esoft/go-skill/easy-api/contract"

type ioFactory struct {
	ioPath contract.IIOPath
}

func (m ioFactory) BuildDirectory(paths ...string) contract.IIODirectory {
	return NewIODirectory(m.ioPath, paths...)
}

func (m ioFactory) BuildFile(paths ...string) contract.IIOFile {
	return NewIOFile(m.ioPath, paths...)
}

func NewIOFactory(ioPath contract.IIOPath) contract.IIOFactory {
	return &ioFactory{
		ioPath: ioPath,
	}
}
