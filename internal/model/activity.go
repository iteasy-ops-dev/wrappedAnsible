package model

// TODO: 무엇을 관리할 것인가.
type Activiry struct {
	ID      string `bson:"_id,omitempty"`
	Email   string `bson:"email"`
	Name    string `bson:"name"`
	RunType string `bson:"runType"`

	AtDate int64 `bson:"atDate"`
}
