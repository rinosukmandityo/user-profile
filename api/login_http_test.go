package api_test

import (
	"bytes"
	"github.com/rinosukmandityo/user-profile/api/caller"
	"github.com/rinosukmandityo/user-profile/helper"
	"net/http"
	"testing"

	. "github.com/rinosukmandityo/user-profile/api"
	m "github.com/rinosukmandityo/user-profile/models"
)

func TestAuthenticateUserSuccess(t *testing.T) {
	testdata := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}}

	defer func() {
		if err := cleanupUserData(); err != nil {
			t.Fatal(err)
		}
		if err := cleanupSessionData(); err != nil {
			t.Fatal(err)
		}
	}()
	_data := testdata[0]

	if err := userService.Store(&_data); err != nil {
		t.Fatal(err)
	}

	dataBytes, err := getBytes(testdata[0])
	if err != nil {
		t.Errorf("cannot get data byte: %s", err.Error())
	}
	req, e := http.NewRequest("POST", "/auth", bytes.NewReader(dataBytes))
	if e != nil {
		t.Errorf("failed to create a mew request: %s ", e.Error())
	}

	resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(&helper.ResultInfo{}).
		SetHeader("Content-Type", ContentTypeJson).Exec()
	if e != nil {
		t.Errorf("failed to call an API: %s ", e.Error())
	}

	result := respBody.(*helper.ResultInfo)
	actualMsg := result.Message
	expectedMsg := helper.SuccessLogin

	if resp.Code != http.StatusOK {
		t.Errorf("status response is not correct, want %d, got %d", http.StatusOK, resp.Code)
	}
	if actualMsg != expectedMsg {
		t.Errorf("error message is not correct, want %s, got %s", expectedMsg, actualMsg)
	}
}

func TestAuthenticateUserFailed(t *testing.T) {
	testdata := []m.User{{
		Name:     "User 01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}}

	defer func() {
		if err := cleanupUserData(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := seedUserData(testdata); err != nil {
		t.Fatal(err)
	}

	wrongEmail := testdata[0]
	wrongEmail.Email = "wrong email"

	wrongPassword := testdata[0]
	wrongPassword.Password = "wrong password"

	testTable := []struct {
		name               string
		data               m.User
		expectedMsg        string
		expectedStatusCode int
		expectedResp       interface{}
	}{
		{
			name:               "email_does_not_exists",
			data:               wrongEmail,
			expectedMsg:        helper.ErrAuthEmailMsg,
			expectedStatusCode: http.StatusNotFound,
			expectedResp:       &helper.ResultInfo{},
		},
		{
			name:               "password_does_not_match",
			data:               wrongPassword,
			expectedMsg:        helper.ErrAuthPasswordMsg,
			expectedStatusCode: http.StatusNotFound,
			expectedResp:       &helper.ResultInfo{},
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			dataBytes, err := getBytes(tt.data)
			if err != nil {
				t.Errorf("cannot get data byte: %s", err.Error())
			}

			req, e := http.NewRequest("POST", "/auth", bytes.NewReader(dataBytes))
			if e != nil {
				t.Errorf("failed to create a mew request: %s ", e.Error())
			}

			resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(tt.expectedResp).
				SetHeader("Content-Type", ContentTypeJson).Exec()
			if e != nil {
				t.Errorf("failed to call an API: %s ", e.Error())
			}

			result := respBody.(*helper.ResultInfo)
			actualMsg := result.Message

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("status response is not correct, want %d, got %d", tt.expectedStatusCode, resp.Code)
			}

			if actualMsg != tt.expectedMsg {
				t.Errorf("error messsage is not correct, want %s, got %s", tt.expectedMsg, actualMsg)
			}
		})
	}
}
