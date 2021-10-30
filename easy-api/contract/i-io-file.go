package contract

type IIOFile interface {
	IIONode

	GetExt() string
	Read(data interface{}) error
	Write(data interface{}) error
}
