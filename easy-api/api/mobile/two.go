package mobile

type TwoApi struct {
}

func (m TwoApi) Call() (interface{}, error) {
	return "this is 2", nil
}
