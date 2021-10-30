package contract

type IIOPath interface {
	GetRoot() string
	Join(paths ...string) string
}
