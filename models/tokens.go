package models

import (
	"time"
)

type Token struct {
	ID        string    `json:"ID" bson:"ID" msgpack:"ID" db:"ID"`
	UserID    string    `json:"UserID" bson:"UserID" msgpack:"UserID" db:"UserID"`
	Created   time.Time `json:"Created" bson:"Created" msgpack:"Created" db:"Created"`
	Expired   time.Time `json:"Expired" bson:"Expired" msgpack:"Expired" db:"Expired"`
	IsClaimed bool      `json:"IsClaimed" bson:"IsClaimed" msgpack:"IsClaimed" db:"IsClaimed"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (s *Token) SplitByField() []interface{} {
	return []interface{}{
		s.ID,
		s.UserID,
		s.Created,
		s.Expired,
		s.IsClaimed,
	}
}
