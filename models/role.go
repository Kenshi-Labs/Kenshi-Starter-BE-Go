package models

type Role struct {
	Name		string		`json:"" bson:""`
	Permissions	[]string	`json:"permissions" bson:"permissions"`
}