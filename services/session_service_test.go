package services_test

import (
	"github.com/pkg/errors"
	"github.com/rinosukmandityo/user-profile/helper"
	"reflect"
	"testing"
	"time"

	m "github.com/rinosukmandityo/user-profile/models"
)

func TestCreateNewSession(t *testing.T) {
	userData := m.User{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}

	defer func() {
		if err := cleanupSessionData(); err != nil {
			t.Fatal(err)
		}
	}()

	testTable := []struct {
		name     string
		userData m.User
	}{
		{
			name:     "success",
			userData: userData,
		},
		{
			name:     "already_exists",
			userData: userData,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			actualData, e := sessionSvc.CreateNewSession(userData)
			if e != nil {
				t.Errorf("unable to create a new session")
			}

			if actualData.UserID != tt.userData.ID {
				t.Errorf("data is not correct, want: %s, got: %s", tt.userData.ID, actualData.UserID)
			}
			if actualData.Email != tt.userData.Email {
				t.Errorf("data is not correct, want: %s, got: %s", tt.userData.Email, actualData.Email)
			}

			time.Sleep(time.Second)
		})
	}
}

func TestAuthenticateUserSuccess(t *testing.T) {
	testData := []m.User{{
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
	}()

	userData := testData[0]
	if e := userService.Store(&userData); e != nil {
		t.Fatal(e)
	}

	testTable := []struct {
		name     string
		email    string
		password string
	}{
		{
			name:     "success",
			email:    userData.Email,
			password: testData[0].Password,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			_, _, e := sessionSvc.Authenticate(tt.email, tt.password)
			if e != nil {
				t.Errorf("authentication failed: %s", e.Error())
			}
		})
	}
}

func TestAuthenticateUserFailed(t *testing.T) {
	testData := []m.User{{
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
	}()

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}
	userData := testData[0]

	testTable := []struct {
		name             string
		email            string
		password         string
		expectedData     m.User
		expectedErrorMsg string
	}{
		{
			name:             "email_does_not_exists",
			email:            "wrong email",
			password:         userData.Password,
			expectedData:     m.User{},
			expectedErrorMsg: helper.ErrUserNotFound.Error(),
		},
		{
			name:             "password_does_not_match",
			email:            userData.Email,
			password:         "wrong password",
			expectedData:     userData,
			expectedErrorMsg: helper.ErrPasswordDoesNotMatch.Error(),
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			authenticated, actualData, e := sessionSvc.Authenticate(tt.email, tt.password)
			if e == nil {
				t.Errorf("authentication should be failed")
			}
			if authenticated {
				t.Errorf("user should not be authenticated")
			}
			if !reflect.DeepEqual(tt.expectedData, actualData) {
				t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", tt.expectedData, actualData)
			}

			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is not correct, want %s, got %s", tt.expectedErrorMsg, e.Error())
			}
		})
	}
}
