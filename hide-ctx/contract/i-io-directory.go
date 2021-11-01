package contract

type IIODirectory interface {
	IIONode

	FindDirectories() []IIODirectory
	FindFiles() []IIOFile
}
