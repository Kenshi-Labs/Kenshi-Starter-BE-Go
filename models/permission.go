package models

type Permission struct {
	ID string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Resource string `json:"resource" bson:"resource"`
	Actions []string `json:"actions" bson:"actions"`
}