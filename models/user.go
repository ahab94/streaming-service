package models

import (
	"time"

	"github.com/fatih/structs"
)

// User this struct holds end customer information
type User struct {
	UserID    string    `db:"user_id" structs:"user_id" json:"user_id"`
	Timestamp time.Time `db:"timestamp" structs:"timestamp" json:"timestamp"`
}

// Map instantiates a map of user struct
func (u *User) Map() map[string]interface{} {
	return structs.Map(u)
}

// Names instantiates a list of User struct fields
func (u *User) Names() []string {
	fields := structs.Fields(u)
	names := make([]string, len(fields))
	for i, field := range fields {
		name := field.Name()
		tagName := field.Tag(structs.DefaultTagName)
		if tagName != "" {
			name = tagName
		}
		names[i] = name
	}
	return names
}
