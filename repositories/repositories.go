package repositories

import (
	m "github.com/rinosukmandityo/user-profile/models"
)

type UserRepository interface {
	GetBy(filter map[string]interface{}) (m.User, error)
	Store(data *m.User) error
	Update(data map[string]interface{}, id string) (m.User, error)
	DeleteAll() error
}

type SessionRepository interface {
	GetBy(filter map[string]interface{}) (m.Session, error)
	Store(data *m.Session) error
	Update(data map[string]interface{}, id string) (m.Session, error)
	Authenticate(email, password string) (bool, m.User, error)
	DeleteAll() error
}

type TokenRepository interface {
	GetLatestToken(userid string) (m.Token, error)
	Store(data *m.Token) error
	Update(data map[string]interface{}, id string) error
	DeleteAll() error
}
