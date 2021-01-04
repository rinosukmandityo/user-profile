package api_test

import (
	"github.com/go-chi/chi"
	. "github.com/rinosukmandityo/user-profile/api"
	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"
	rh "github.com/rinosukmandityo/user-profile/repositories/helper"
	"github.com/rinosukmandityo/user-profile/services"
	"github.com/rinosukmandityo/user-profile/services/logic"
	"net/http/httptest"
)

var (
	sessionRepo repo.SessionRepository
	sessionSvc  services.SessionService
	userRepo    repo.UserRepository
	userService services.UserService
	r           *chi.Mux
	ts          *httptest.Server
)

func init() {
	userRepo, sessionRepo, _ = rh.ChooseRepo()
	r = RegisterHandler()

	sessionSvc = logic.NewSessionService(sessionRepo)
	userService = logic.NewUserService(userRepo)
}

func getBytes(_data m.User) ([]byte, error) {
	dataBytes, e := GetSerializer(ContentTypeJson).Encode(&_data)
	if e != nil {
		return dataBytes, e
	}
	return dataBytes, nil
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
