package models

import (
	"time"
)

type Session struct {
	ID      string    `json:"ID" bson:"ID" msgpack:"ID" db:"ID"`
	UserID  string    `json:"UserID" bson:"UserID" msgpack:"UserID" db:"UserID"`
	Email   string    `json:"Email" bson:"Email" msgpack:"Email" db:"Email"`
	Created time.Time `json:"Created" bson:"Created" msgpack:"Created" db:"Created"`
	Expired time.Time `json:"Expired" bson:"Expired" msgpack:"Expired" db:"Expired"`
}

func (s *Session) TableName() string {
	return "sessions"
}

func (s *Session) SplitByField() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Email,
		s.Created,
		s.Expired,
	}
}
