package contract

type IApi interface {
	Call() (interface{}, error)
}
