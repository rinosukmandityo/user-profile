package logic

import (
	"fmt"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"
	svc "github.com/rinosukmandityo/user-profile/services"

	errs "github.com/pkg/errors"
)

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) svc.UserService {
	return &userService{
		userRepo,
	}
}

func (u *userService) GetById(id string) (m.User, error) {
	res, e := u.userRepo.GetBy(map[string]interface{}{"ID": id})
	if e != nil {
		return res, e
	}

	return res, nil

}
func (u *userService) Store(data *m.User) error {
	if data.ID == "" {
		data.ID = fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}
	if isFound, _, _ := u.GetByEmail(data.Email); isFound {
		return errs.Wrap(helper.ErrEmailDuplicate, "service.User.Store")
	}
	if data.Password != "" {
		data.Password = repo.EncryptPassword(data.Password)
	} else {
		data.IsGoogleAuth = true
	}
	return u.userRepo.Store(data)

}
func (u *userService) Update(data map[string]interface{}, id string) (m.User, error) {
	user := m.User{}
	var e error

	if isFound, usr, _ := u.GetByEmail(data["Email"].(string)); isFound {
		if usr.ID != id {
			return user, errs.Wrap(helper.ErrEmailDuplicate, "service.User.Update")
		}
	}
	if _, hasData := data["Password"]; hasData {
		if data["Password"].(string) != "" {
			data["Password"] = repo.EncryptPassword(data["Password"].(string))
		}
	}
	user, e = u.userRepo.Update(data, id)
	if e != nil {
		return user, errs.Wrap(e, "service.User.Update")
	}
	return user, nil

}

func (u *userService) GetByEmail(email string) (bool, m.User, error) {
	res, e := u.userRepo.GetBy(map[string]interface{}{"Email": email})
	if e != nil {
		return false, res, e
	}

	return true, res, nil
}
