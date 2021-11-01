package contract

type IIOFactory interface {
	BuildDirectory(pathArgs ...string) IIODirectory
	BuildFile(pathArgs ...string) IIOFile
}
