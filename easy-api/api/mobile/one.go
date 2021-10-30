package mobile

type OneApi struct {
}

func (m OneApi) Call() (interface{}, error) {
	return "this is 1", nil
}
