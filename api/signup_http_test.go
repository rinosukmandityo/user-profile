package api_test

import (
	"bytes"
	"net/http"
	"reflect"
	"strings"
	"testing"

	. "github.com/rinosukmandityo/user-profile/api"
	"github.com/rinosukmandityo/user-profile/api/caller"
	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	"github.com/rinosukmandityo/user-profile/repositories"
)

func TestSignUpSuccess(t *testing.T) {
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

	expectedData := testdata[0]
	dataBytes, err := getBytes(expectedData)
	if err != nil {
		t.Errorf("cannot get data byte: %s", err.Error())
	}
	req, e := http.NewRequest("POST", "/dosignup", bytes.NewReader(dataBytes))
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
	expectedMsg := helper.SuccessSignup

	actualData, e := userService.GetById("userid01")
	if e != nil {
		t.Errorf("unable to get data by ID")
	}

	if resp.Code != http.StatusOK {
		t.Errorf("status response is not correct, want %d, got %d", http.StatusOK, resp.Code)
	}

	if actualMsg != expectedMsg {
		t.Errorf("error message is not correct, want %s, got %s", expectedMsg, actualMsg)
	}

	expectedData.Password = repositories.EncryptPassword(expectedData.Password)
	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
}

func TestSignUpFailed(t *testing.T) {
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

	duplicateEmail := testdata[0]
	duplicateEmail.ID = "random ID"

	duplicateEntry := testdata[0]
	duplicateEntry.Email = "random email"

	testTable := []struct {
		name               string
		data               m.User
		expectedMsg        string
		expectedStatusCode int
		expectedResp       interface{}
	}{
		{
			name:               "duplicate_email",
			data:               duplicateEmail,
			expectedMsg:        helper.ErrDuplicateEmail,
			expectedStatusCode: http.StatusBadRequest,
			expectedResp:       &helper.ResultInfo{},
		},
		{
			name:               "duplicate_entry",
			data:               duplicateEntry,
			expectedMsg:        "Duplicate entry 'userid01' for key 'ID'",
			expectedStatusCode: http.StatusBadRequest,
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

			req, e := http.NewRequest("POST", "/dosignup", bytes.NewReader(dataBytes))
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

			if !strings.Contains(actualMsg, tt.expectedMsg) {
				t.Errorf("error messsage is not correct, want %s, got %s", tt.expectedMsg, actualMsg)
			}
		})
	}
}
