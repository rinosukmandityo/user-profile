package services_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rinosukmandityo/user-profile/helper"
	"github.com/rinosukmandityo/user-profile/repositories"
	"reflect"
	"testing"
	"time"

	m "github.com/rinosukmandityo/user-profile/models"
)

func TestResetPasswordByMailSuccess(t *testing.T) {
	userData := []m.User{
		{
			Name:     "User 01",
			Password: "Password.1",
			ID:       "userid01",
			Email:    "usermail01@gmail.com",
			Address:  "User Address 01",
			IsActive: false,
		},
		{
			Name:     "User 02",
			Password: "Password.1",
			ID:       "userid02",
			Email:    "usermail02@gmail.com",
			Address:  "User Address 02",
			IsActive: false,
		},
		{
			Name:     "User 03",
			Password: "Password.1",
			ID:       "userid03",
			Email:    "usermail03@gmail.com",
			Address:  "User Address 03",
			IsActive: false,
		},
	}

	tokenExpired := m.Token{
		ID:        "token02",
		UserID:    userData[1].ID,
		Created:   time.Now().UTC().Add(time.Minute * -30),
		Expired:   time.Now().UTC().Add(time.Minute * -1),
		IsClaimed: false,
	}
	tokenClaimed := m.Token{
		ID:        "token03",
		UserID:    userData[2].ID,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(time.Minute * 30),
		IsClaimed: true,
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}
	if e := tokenRepo.Store(&tokenExpired); e != nil {
		t.Fatal(e)
	}
	if e := tokenRepo.Store(&tokenClaimed); e != nil {
		t.Fatal(e)
	}

	defer func() {
		if err := cleanupUserData(); err != nil {
			t.Fatal(err)
		}
		if err := cleanupTokenData(); err != nil {
			t.Fatal(err)
		}
	}()

	testTable := []struct {
		name             string
		expectedUserData m.User
	}{
		{
			name:             "create_new_token",
			expectedUserData: userData[0],
		},
		{
			name:             "recreate_expired_token",
			expectedUserData: userData[1],
		},
		{
			name:             "recreate_claimed_token",
			expectedUserData: userData[2],
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			actualUserData, tokenID, e := tokenService.ResetPasswordByMail(tt.expectedUserData.Email, time.Minute*30)
			if e != nil {
				t.Errorf("unable to create a new session: %s", e.Error())
			}

			if !reflect.DeepEqual(tt.expectedUserData, actualUserData) {
				t.Errorf("user data is not correct, \nwant: \n%+v, \ngot: \n%+v", tt.expectedUserData, actualUserData)
			}
			if tokenID == "" {
				t.Errorf("token ID should not be empty")
			}
		})
	}
}

func TestResetPasswordByMailFailed(t *testing.T) {
	userData := []m.User{
		{
			Name:     "User 01",
			Password: "Password.1",
			ID:       "userid01",
			Email:    "usermail01@gmail.com",
			Address:  "User Address 01",
			IsActive: false,
		},
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}

	defer func() {
		if err := cleanupUserData(); err != nil {
			t.Fatal(err)
		}
	}()

	wrongEmail := userData[0]
	wrongEmail.Email = "wrong email"

	testTable := []struct {
		name             string
		expectedUserData m.User
		expectedErrorMsg string
	}{
		{
			name:             "wrong_email",
			expectedUserData: m.User{},
			expectedErrorMsg: helper.ErrUserNotFound.Error(),
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			actualUserData, tokenID, e := tokenService.ResetPasswordByMail(tt.expectedUserData.Email, time.Minute*30)
			if e == nil {
				t.Errorf("it should be error to create new token")
			}

			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is not correct, want %s, got %s", tt.expectedErrorMsg, e.Error())
			}

			if !reflect.DeepEqual(tt.expectedUserData, actualUserData) {
				t.Errorf("user data is not correct, \nwant: \n%+v, \ngot: \n%+v", tt.expectedUserData, actualUserData)
			}
			if tokenID != "" {
				t.Errorf("token ID should be empty")
			}
		})
	}
}

func TestIsTokenValidSuccess(t *testing.T) {
	userData := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}}

	defer func() {
		if e := cleanupUserData(); e != nil {
			t.Fatal(e)
		}
		if e := cleanupTokenData(); e != nil {
			t.Fatal(e)
		}
	}()

	tokenData := m.Token{
		ID:        "token01",
		UserID:    userData[0].ID,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(time.Minute * 30),
		IsClaimed: false,
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}

	if e := tokenRepo.Store(&tokenData); e != nil {
		t.Fatal(e)
	}

	testTable := []struct {
		name    string
		userID  string
		tokenID string
	}{
		{
			name:    "success",
			userID:  userData[0].ID,
			tokenID: tokenData.ID,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			found, e := tokenService.IsTokenValid(tt.userID, tt.tokenID)
			if e != nil {
				t.Errorf("failed to check token: %s", e.Error())
			}
			if !found {
				t.Errorf("token should be found")
			}
		})
	}
}

func TestIsTokenValidFailed(t *testing.T) {
	userData := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}, {
		Name:     "User 02",
		Password: "Password.1",
		ID:       "userid02",
		Email:    "usermail02@gmail.com",
		Address:  "User Address 02",
		IsActive: false,
	}, {
		Name:     "User 03",
		Password: "Password.1",
		ID:       "userid03",
		Email:    "usermail03@gmail.com",
		Address:  "User Address 03",
		IsActive: false,
	}}

	defer func() {
		if e := cleanupUserData(); e != nil {
			t.Fatal(e)
		}
		if e := cleanupTokenData(); e != nil {
			t.Fatal(e)
		}
	}()

	tokenExpired := m.Token{
		ID:        "token01",
		UserID:    userData[0].ID,
		Created:   time.Now().UTC().Add(time.Minute * -30),
		Expired:   time.Now().UTC().Add(time.Minute * -1),
		IsClaimed: false,
	}
	tokenClaimed := m.Token{
		ID:        "token02",
		UserID:    userData[1].ID,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(time.Minute * 30),
		IsClaimed: true,
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}
	if e := tokenRepo.Store(&tokenExpired); e != nil {
		t.Fatal(e)
	}
	if e := tokenRepo.Store(&tokenClaimed); e != nil {
		t.Fatal(e)
	}

	testTable := []struct {
		name             string
		userID           string
		tokenID          string
		expectedErrorMsg string
	}{
		{
			name:             "user_not_found",
			userID:           "not found",
			tokenID:          "",
			expectedErrorMsg: helper.ErrUserNotFound.Error(),
		},
		{
			name:             "token_not_found",
			userID:           userData[0].ID,
			tokenID:          "not found",
			expectedErrorMsg: helper.ErrTokenNotFound.Error(),
		},
		{
			name:             "token_not_created_yet",
			userID:           userData[2].ID,
			tokenID:          "not found",
			expectedErrorMsg: helper.ErrTokenNotFound.Error(),
		},
		{
			name:             "token_expired",
			userID:           userData[0].ID,
			tokenID:          tokenExpired.ID,
			expectedErrorMsg: helper.ErrTokenExpired,
		},
		{
			name:             "token_claimed",
			userID:           userData[1].ID,
			tokenID:          tokenClaimed.ID,
			expectedErrorMsg: helper.ErrTokenClaimed,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			found, e := tokenService.IsTokenValid(tt.userID, tt.tokenID)
			if e == nil {
				t.Errorf("token validation should be failed")
			}
			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is not correct, want %s, got %s", tt.expectedErrorMsg, e.Error())
			}
			if found {
				t.Errorf("token should be not found")
			}
		})
	}
}

func TestChangePasswordTokenSuccess(t *testing.T) {
	userData := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}}

	defer func() {
		if e := cleanupUserData(); e != nil {
			t.Fatal(e)
		}
		if e := cleanupTokenData(); e != nil {
			t.Fatal(e)
		}
	}()

	tokenData := m.Token{
		ID:        "token01",
		UserID:    userData[0].ID,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(time.Minute * 30),
		IsClaimed: false,
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}

	if e := tokenRepo.Store(&tokenData); e != nil {
		t.Fatal(e)
	}

	testTable := []struct {
		name        string
		userID      string
		newPassword string
		tokenID     string
	}{
		{
			name:        "success",
			newPassword: "new password",
			userID:      userData[0].ID,
			tokenID:     tokenData.ID,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			e := tokenService.ChangePasswordToken(tt.userID, tt.newPassword, tt.tokenID)
			if e != nil {
				t.Errorf("failed to change password with token: %s", e.Error())
			}
			actualUser, e := userRepo.GetBy(map[string]interface{}{"ID": tt.userID})
			if e != nil {
				t.Errorf("unable to get actual user: %s", e.Error())
			}
			expectedPassword := repositories.EncryptPassword(tt.newPassword)
			if expectedPassword != actualUser.Password {
				t.Errorf("password is incorrect, want %s, got %s", expectedPassword, actualUser.Password)
			}
			actualToken, e := tokenRepo.GetLatestToken(tt.userID)
			if e != nil {
				t.Errorf("unable to get actual token: %s", e.Error())
			}
			if !actualToken.IsClaimed {
				t.Errorf("token should be claimed")
			}
		})
	}
}

func TestChangePasswordTokenFailed(t *testing.T) {
	userData := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}, {
		Name:     "User 02",
		Password: "Password.1",
		ID:       "userid02",
		Email:    "usermail02@gmail.com",
		Address:  "User Address 02",
		IsActive: false,
	}}

	defer func() {
		if e := cleanupUserData(); e != nil {
			t.Fatal(e)
		}
		if e := cleanupTokenData(); e != nil {
			t.Fatal(e)
		}
	}()

	tokenSuccess := m.Token{
		ID:        "token01",
		UserID:    userData[0].ID,
		Created:   time.Now().UTC(),
		Expired:   time.Now().UTC().Add(time.Minute * 30),
		IsClaimed: true,
	}

	if e := seedUserData(userData); e != nil {
		t.Fatal(e)
	}

	if e := tokenRepo.Store(&tokenSuccess); e != nil {
		t.Fatal(e)
	}

	testTable := []struct {
		name             string
		userID           string
		newPassword      string
		tokenID          string
		expectedErrorMsg string
	}{
		{
			name:             "user_not_found",
			userID:           "not found",
			newPassword:      "",
			tokenID:          "",
			expectedErrorMsg: helper.ErrTokenNotFound.Error(),
		},
		{
			name:             "token_not_found",
			userID:           userData[0].ID,
			newPassword:      "",
			tokenID:          "not found",
			expectedErrorMsg: helper.ErrTokenNotMatch.Error(),
		},
		{
			name:             "token_not_created_yet",
			userID:           userData[1].ID,
			newPassword:      "",
			tokenID:          "not found",
			expectedErrorMsg: helper.ErrTokenNotFound.Error(),
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			e := tokenService.ChangePasswordToken(tt.userID, tt.newPassword, tt.tokenID)
			if e == nil {
				t.Errorf("should be failed to change password with token")
			}
			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is incorrect, want %s, got %s", e.Error(), tt.expectedErrorMsg)
			}
			fmt.Println(e)
		})
	}
}
