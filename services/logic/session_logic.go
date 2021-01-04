package logic

import (
	"fmt"
	"time"

	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"
	svc "github.com/rinosukmandityo/user-profile/services"
)

var (
	_expiredduration = time.Duration(time.Minute * 30)
)

type sessionService struct {
	sessionRepo repo.SessionRepository
}

func NewSessionService(sessionRepo repo.SessionRepository) svc.SessionService {
	return &sessionService{
		sessionRepo,
	}
}

func (u *sessionService) Authenticate(email, password string) (found bool, user m.User, e error) {
	return u.sessionRepo.Authenticate(email, password)
}

func (u *sessionService) GetBy(filter map[string]interface{}) (m.Session, error) {
	return u.sessionRepo.GetBy(filter)
}

func (u *sessionService) CreateNewSession(user m.User) (tSession m.Session, e error) {
	tSession, _ = u.sessionRepo.GetBy(map[string]interface{}{"Email": user.Email})

	if tSession.ID == "" || tSession.Expired.Before(time.Now().UTC()) {
		tSession.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
		tSession.UserID = user.ID
		tSession.Email = user.Email
		tSession.Created = time.Now().UTC()
		tSession.Expired = time.Now().UTC().Add(_expiredduration)

		e = u.sessionRepo.Store(&tSession)
	} else {
		tSession.Expired = time.Now().UTC().Add(_expiredduration)

		_, e = u.sessionRepo.Update(map[string]interface{}{"Expired": tSession.Expired}, tSession.ID)
	}

	return
}

func (u *sessionService) IsSessionActive(id string) (stat bool, sess m.Session, e error) {
	sess = m.Session{}

	sess, e = u.sessionRepo.GetBy(map[string]interface{}{"ID": id})
	if e != nil {
		return
	}
	if sess.Expired.Before(time.Now().UTC()) {
		return
	}
	stat = true
	return
}

func (u *sessionService) Logout(data m.Session) error {
	_, e := u.sessionRepo.Update(map[string]interface{}{"Expired": time.Now().UTC()}, data.ID)
	return e

}
