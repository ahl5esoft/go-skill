package mobile

type QuickPlayApi struct{}

func (m QuickPlayApi) Call() (interface{}, error) {
	return "start quick play", nil
}
