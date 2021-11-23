package global

type RoutePermission struct {
	Conditions []ValueCondition
	ID         string `alias:"" bson:"_id" db:"_id"`
	Response   map[string]interface{}
	Routes     []string
}

func (m RoutePermission) GetID() string {
	return m.ID
}
