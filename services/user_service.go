package services

import (
	m "github.com/rinosukmandityo/user-profile/models"
)

type UserService interface {
	GetById(id string) (m.User, error)
	GetByEmail(email string) (bool, m.User, error)
	Store(data *m.User) error
	Update(data map[string]interface{}, id string) (m.User, error)
}
