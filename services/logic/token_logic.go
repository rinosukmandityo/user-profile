package logic

import (
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"
	svc "github.com/rinosukmandityo/user-profile/services"
)

type tokenService struct {
	tokenRepo repo.TokenRepository
	userRepo  repo.UserRepository
}

func NewTokenService(tokenRepo repo.TokenRepository, userRepo repo.UserRepository) svc.TokenService {
	return &tokenService{
		tokenRepo, userRepo,
	}
}

func (u *tokenService) ResetPasswordByMail(email string, expired time.Duration) (user m.User, tokenID string, e error) {
	user, e = u.userRepo.GetBy(map[string]interface{}{"Email": email})
	if e != nil {
		return
	}

	tToken, e := u.tokenRepo.GetLatestToken(user.ID)
	if time.Now().UTC().After(tToken.Expired) {
		e = fmt.Errorf(helper.ErrTokenExpired)
	} else if tToken.IsClaimed {
		e = fmt.Errorf(helper.ErrTokenClaimed)
	}

	tokenID = tToken.ID
	if tokenID != "" && e == nil {
		return
	}
	// token expired or never been created
	tToken, e = u.CreateNewToken(user.ID, expired)
	if e != nil {
		e = fmt.Errorf("Reset password failed to create token: %s", e.Error())
	}
	tokenID = tToken.ID

	return
}

func (u *tokenService) ChangePassword(userID, passwd string) (e error) {
	if _, e = u.userRepo.Update(map[string]interface{}{"Password": passwd}, userID); e != nil {
		if errors.Cause(e).Error() == helper.ErrUserNotFound.Error() {
			e = nil
		}
		return e
	}

	return
}

func (u *tokenService) ChangePasswordToken(userID, passwd, tokenID string) (e error) {
	gToken, e := u.tokenRepo.GetLatestToken(userID)
	if e != nil {
		return
	}
	if gToken.ID != tokenID {
		return helper.ErrTokenNotMatch
	}
	if e = u.ChangePassword(userID, repo.EncryptPassword(passwd)); e == nil {
		e = u.Claim(gToken)
	}

	return
}

func (u *tokenService) IsTokenValid(userid, tokenID string) (isValid bool, e error) {
	if _, e = u.userRepo.GetBy(map[string]interface{}{"ID": userid}); e != nil {
		if errors.Cause(e).Error() == helper.ErrUserNotFound.Error() {
			e = helper.ErrUserNotFound
		}
		return
	}
	tToken, e := u.tokenRepo.GetLatestToken(userid)
	if e != nil {
		if errors.Cause(e).Error() == helper.ErrTokenNotFound.Error() {
			e = helper.ErrTokenNotFound
		}
		return
	}

	if tToken.ID == tokenID {
		if time.Now().UTC().After(tToken.Expired) {
			e = fmt.Errorf(helper.ErrTokenExpired)
		} else if tToken.IsClaimed {
			e = fmt.Errorf(helper.ErrTokenClaimed)
		}
	} else {
		e = helper.ErrTokenNotFound
	}
	if e != nil {
		return
	}

	return true, nil
}

func (u *tokenService) CreateNewToken(userid string, validity time.Duration) (token m.Token, e error) {
	token = m.Token{
		ID:        fmt.Sprintf("%d", time.Now().UTC().UnixNano()),
		UserID:    userid,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(validity),
		IsClaimed: false,
	}
	e = u.tokenRepo.Store(&token)
	return
}

func (u *tokenService) Claim(token m.Token) error {
	return u.tokenRepo.Update(map[string]interface{}{"IsClaimed": true}, token.ID)
}
