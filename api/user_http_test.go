package api_test

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	. "github.com/rinosukmandityo/user-profile/api"
	"github.com/rinosukmandityo/user-profile/api/caller"
	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
)

func TestUpdateUserSuccess(t *testing.T) {
	testData := []m.User{{
		Name:      "User 01",
		Password:  "Password.1",
		ID:        "userid01",
		Email:     "usermail01@gmail.com",
		Telephone: "08123456789",
		Address:   "User Address 01",
		IsActive:  false,
	}}

	defer func() {
		if e := cleanupUserData(); e != nil {
			t.Fatal(e)
		}
	}()

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}

	updatedData := testData[0]
	updatedData.Telephone = updatedData.Telephone + "UPDATED"

	testTable := []struct {
		name               string
		validSessionCookie bool
		expectedStatusCode int
		expectedData       m.User
		expectedResp       interface{}
	}{
		{
			name:               "updated_successfully",
			validSessionCookie: true,
			expectedStatusCode: http.StatusOK,
			expectedData:       updatedData,
			expectedResp:       &helper.ResultInfo{},
		},
		{
			name:               "no_data_updated",
			validSessionCookie: true,
			expectedStatusCode: http.StatusOK,
			expectedData:       testData[0],
			expectedResp:       &helper.ResultInfo{},
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			_data := tt.expectedData
			dataBytes, e := getBytes(_data)
			if e != nil {
				t.Errorf("cannot get data byte: %s", e.Error())
			}

			tSession, e := sessionSvc.CreateNewSession(_data)
			if e != nil {
				t.Errorf("cannot create a new session: %s", e.Error())
			}
			sessionID := ""
			if tt.validSessionCookie {
				sessionID = tSession.ID
			}
			cookieSession := fmt.Sprintf("%s=%s; %s=%s", helper.SESSION_COOKIE_KEY, sessionID, helper.EMAIL_COOKIE_KEY, _data.Email)

			req, e := http.NewRequest("PUT", "/user", bytes.NewReader(dataBytes))
			if e != nil {
				t.Errorf("failed to create a mew request: %s ", e.Error())
			}

			resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(tt.expectedResp).
				SetHeader("Cookie", cookieSession).
				SetHeader("Content-Type", ContentTypeJson).Exec()
			if e != nil {
				t.Errorf("failed to call an API: %s ", e.Error())
			}

			result := respBody.(*helper.ResultInfo)
			respData := result.Data.(map[string]interface{})
			actualData := m.FromMapToUser(respData)

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("status response is not correct, want %d, got %d", tt.expectedStatusCode, resp.Code)
			}
			if reflect.DeepEqual(tt.expectedData, actualData) {
				t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", tt.expectedData, actualData)
			}
			if e := sessionRepo.DeleteAll(); e != nil {
				t.Errorf("unable to delete session data")
			}
		})
	}
}

func TestUpdateUserFailed(t *testing.T) {
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

	testTable := []struct {
		name               string
		validSessionCookie bool
		expectedStatusCode int
		expectedData       m.User
		expectedResp       interface{}
	}{
		{
			name:               "email_does_not_exists_redirected",
			validSessionCookie: true,
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedData:       m.User{},
			expectedResp:       nil,
		},
		{
			name:               "session_invalid_redirected",
			validSessionCookie: false,
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedData:       testData[0],
			expectedResp:       nil,
		},
		{
			name:               "no_data_updated",
			validSessionCookie: true,
			expectedStatusCode: http.StatusOK,
			expectedData:       testData[0],
			expectedResp:       nil,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			_data := tt.expectedData
			dataBytes, e := getBytes(_data)
			if e != nil {
				t.Errorf("cannot get data byte: %s", e.Error())
			}

			tSession, e := sessionSvc.CreateNewSession(_data)
			if e != nil {
				t.Errorf("cannot create a new session: %s", e.Error())
			}
			sessionID := ""
			if tt.validSessionCookie {
				sessionID = tSession.ID
			}
			cookieSession := fmt.Sprintf("%s=%s; %s=%s", helper.SESSION_COOKIE_KEY, sessionID, helper.EMAIL_COOKIE_KEY, _data.Email)

			req, e := http.NewRequest("PUT", "/user", bytes.NewReader(dataBytes))
			if e != nil {
				t.Errorf("failed to create a mew request: %s ", e.Error())
			}

			resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(tt.expectedResp).
				SetHeader("Cookie", cookieSession).
				SetHeader("Content-Type", ContentTypeJson).Exec()
			if e != nil {
				t.Errorf("failed to call an API: %s ", e.Error())
			}

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("status response is not correct, want %d, got %d", tt.expectedStatusCode, resp.Code)
			}
			if respBody != tt.expectedResp {
				t.Errorf("response body should be nil")
			}
			if e := sessionRepo.DeleteAll(); e != nil {
				t.Errorf("unable to delete session data")
			}
		})
	}
}

func TestGetDataSuccess(t *testing.T) {
	testData := []m.User{{
		Name:         "User 01",
		ID:           "userid01",
		Email:        "usermail01@gmail.com",
		Address:      "User Address 01",
		IsActive:     false,
		IsGoogleAuth: true,
	}}

	defer func() {
		if err := cleanupUserData(); err != nil {
			t.Fatal(err)
		}
	}()

	if err := seedUserData(testData); err != nil {
		t.Fatal(err)
	}

	expectedData := testData[0]
	dataBytes, e := getBytes(expectedData)
	if e != nil {
		t.Errorf("cannot get data byte: %s", e.Error())
	}
	tSession, e := sessionSvc.CreateNewSession(expectedData)
	if e != nil {
		t.Errorf("cannot create a new session: %s", e.Error())
	}
	cookieSession := fmt.Sprintf("%s=%s; %s=%s", helper.SESSION_COOKIE_KEY, tSession.ID, helper.EMAIL_COOKIE_KEY, expectedData.Email)

	req, e := http.NewRequest("GET", "/user", bytes.NewReader(dataBytes))
	if e != nil {
		t.Errorf("failed to create a mew request: %s ", e.Error())
	}

	resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(&helper.ResultInfo{}).
		SetHeader("Cookie", cookieSession).
		SetHeader("Content-Type", ContentTypeJson).Exec()
	if e != nil {
		t.Errorf("failed to call an API: %s ", e.Error())
	}
	result := respBody.(*helper.ResultInfo)
	respData := result.Data.(map[string]interface{})
	actualData := m.FromMapToUser(respData)

	if resp.Code != http.StatusFound {
		t.Errorf("status response is not correct, want %d, got %d", http.StatusFound, resp.Code)
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
	if e := sessionRepo.DeleteAll(); e != nil {
		t.Errorf("unable to delete session data")
	}
}

func TestGetDataByIdFailed(t *testing.T) {
	testData := []m.User{{
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

	if err := seedUserData(testData); err != nil {
		t.Fatal(err)
	}

	dataNotFound := testData[0]
	dataNotFound.Email = "not found"

	testTable := []struct {
		name               string
		validSessionCookie bool
		expectedStatusCode int
		expectedData       m.User
		expectedResp       interface{}
	}{
		{
			name:               "email_not_found_redirected",
			validSessionCookie: true,
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedData:       dataNotFound,
			expectedResp:       nil,
		},
		{
			name:               "invalid_cookie_redirected",
			validSessionCookie: false,
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedData:       testData[0],
			expectedResp:       nil,
		},
	}

	for _, _tt := range testTable {
		tt := _tt
		t.Run(tt.name, func(t *testing.T) {
			_data := tt.expectedData
			tSession, e := sessionSvc.CreateNewSession(_data)
			if e != nil {
				t.Errorf("cannot create a new session: %s", e.Error())
			}
			sessionID := ""
			if tt.validSessionCookie {
				sessionID = tSession.ID
			}
			cookieSession := fmt.Sprintf("%s=%s; %s=%s", helper.SESSION_COOKIE_KEY, sessionID, helper.EMAIL_COOKIE_KEY, _data.Email)

			req, e := http.NewRequest("GET", "/user", nil)
			if e != nil {
				t.Errorf("failed to create a mew request: %s ", e.Error())
			}

			resp, respBody, e := caller.New(r).SetRequest(req).SetResponse(tt.expectedResp).
				SetHeader("Cookie", cookieSession).
				SetHeader("Content-Type", ContentTypeJson).Exec()
			if e != nil {
				t.Errorf("failed to call an API: %s ", e.Error())
			}

			if resp.Code != tt.expectedStatusCode {
				t.Errorf("status response is not correct, want %d, got %d", tt.expectedStatusCode, resp.Code)
			}
			if respBody != nil {
				t.Errorf("response body should be nil")
			}
			if e := sessionRepo.DeleteAll(); e != nil {
				t.Errorf("unable to delete session data")
			}
		})
	}
}
