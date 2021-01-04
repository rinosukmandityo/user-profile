package services_test

import (
	"net/http/httptest"

	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"
	rh "github.com/rinosukmandityo/user-profile/repositories/helper"
	. "github.com/rinosukmandityo/user-profile/services"
	"github.com/rinosukmandityo/user-profile/services/logic"
)

var (
	userRepo     repo.UserRepository
	sessionRepo  repo.SessionRepository
	tokenRepo    repo.TokenRepository
	sessionSvc   SessionService
	userService  UserService
	tokenService TokenService
	ts           *httptest.Server
)

func init() {
	userRepo, sessionRepo, tokenRepo = rh.ChooseRepo()

	sessionSvc = logic.NewSessionService(sessionRepo)
	userService = logic.NewUserService(userRepo)
	tokenService = logic.NewTokenService(tokenRepo, userRepo)
}

func seedUserData(testdata []m.User) error {
	var err error
	for _, data := range testdata {
		if err = userRepo.Store(&data); err != nil {
			return err
		}
	}
	return nil
}

func cleanupUserData() error {
	return userRepo.DeleteAll()
}

func cleanupSessionData() error {
	return sessionRepo.DeleteAll()
}

func cleanupTokenData() error {
	return tokenRepo.DeleteAll()
}
