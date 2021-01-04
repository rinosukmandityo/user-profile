package services

import (
	m "github.com/rinosukmandityo/user-profile/models"
)

type SessionService interface {
	Authenticate(email, password string) (bool, m.User, error)
	GetBy(filter map[string]interface{}) (m.Session, error)
	IsSessionActive(id string) (bool, m.Session, error)
	CreateNewSession(user m.User) (m.Session, error)
	Logout(data m.Session) error
}
