package services

import (
	m "github.com/rinosukmandityo/user-profile/models"
	"time"
)

type TokenService interface {
	IsTokenValid(userid, tokenid string) (bool, error)
	ResetPasswordByMail(email string, duration time.Duration) (m.User, string, error)
	ChangePasswordToken(userID, passwd, tokenid string) error
}
