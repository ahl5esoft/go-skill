package mobile

type EnterApi struct{}

func (m EnterApi) Call() (interface{}, error) {
	return "hello", nil
}
