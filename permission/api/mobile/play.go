package mobile

type PlayApi struct{}

func (m PlayApi) Call() (interface{}, error) {
	return "start play", nil
}
