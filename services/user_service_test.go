package services_test

import (
	"reflect"
	"testing"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"

	"github.com/pkg/errors"
)

func TestInsertUserSuccess(t *testing.T) {
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

	expectedData := testData[0]

	if e := userService.Store(&expectedData); e != nil {
		t.Errorf("failed to save data: %s ", e.Error())
	}

	actualData, e := userService.GetById(expectedData.ID)
	if e != nil {
		t.Errorf("unable to get data: %s", e.Error())
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
}

func TestInsertUserFailed(t *testing.T) {

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

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}
	duplicateID := testData[0]
	duplicateID.Email = ""

	testTable := []struct {
		name             string
		data             m.User
		expectedErrorMsg string
	}{
		{
			name:             "duplicate_email",
			data:             testData[0],
			expectedErrorMsg: helper.ErrEmailDuplicate.Error(),
		},
		{
			name:             "duplicate_ID",
			data:             duplicateID,
			expectedErrorMsg: "Error 1062: Duplicate entry 'userid01' for key 'ID'",
		},
	}

	for _, _tt := range testTable {
		tt := _tt

		t.Run(tt.name, func(t *testing.T) {
			e := userService.Store(&tt.data)
			if e == nil {
				t.Error("duplicate validation is not working")
			}

			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is not correct, want %s, got %s", tt.expectedErrorMsg, e.Error())
			}
		})
	}
}

func TestGetUserSuccess(t *testing.T) {
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

	expectedData := testData[0]

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}

	actualData, e := userService.GetById(expectedData.ID)
	if e != nil {
		t.Errorf("unable to get data: %s", e.Error())
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
}

func TestGetUserFailed(t *testing.T) {
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

	expectedData := testData[0]

	if e := userService.Store(&expectedData); e != nil {
		t.Errorf("failed to save data: %s ", e.Error())
	}

	actualData, e := userService.GetById(expectedData.ID)
	if e != nil {
		t.Errorf("unable to get data: %s", e.Error())
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
}

func TestUpdateUserSuccess(t *testing.T) {
	testData := []m.User{{
		Name:     "User 01",
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

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}

	expectedData := testData[0]
	expectedData.Email += "UPDATED"
	dataMap := expectedData.GetMapFormat()

	actualData, e := userService.Update(dataMap, expectedData.ID)
	if e != nil {
		t.Errorf("failed to update data: %s ", e.Error())
	}

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", expectedData, actualData)
	}
}

func TestUpdateUserFailed(t *testing.T) {
	testData := []m.User{{
		Name:     "User 01",
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

	if e := seedUserData(testData); e != nil {
		t.Fatal(e)
	}

	duplicateEmail := testData[0]
	duplicateEmail.ID = "another ID"

	invalidID := testData[0]
	invalidID.ID = "invalid"
	invalidID.Email = "wrong email"

	testTable := []struct {
		name             string
		data             m.User
		expectedData     m.User
		expectedErrorMsg string
	}{
		{
			name:             "id_does_not_exists",
			data:             invalidID,
			expectedData:     m.User{},
			expectedErrorMsg: helper.ErrUserNotFound.Error(),
		},
		{
			name:             "duplicate_email",
			data:             duplicateEmail,
			expectedData:     m.User{},
			expectedErrorMsg: helper.ErrEmailDuplicate.Error(),
		},
		{
			name:             "no_row_affected",
			data:             testData[0],
			expectedData:     m.User{},
			expectedErrorMsg: helper.ErrUserNotFound.Error(),
		},
	}

	for _, _tt := range testTable {
		tt := _tt

		t.Run(tt.name, func(t *testing.T) {
			dataMap := tt.data.GetMapFormat()
			actualData, e := userService.Update(dataMap, tt.data.ID)
			if e == nil {
				t.Errorf("update should be failed")
			}

			if errors.Cause(e).Error() != tt.expectedErrorMsg {
				t.Errorf("error message is not correct, should be contain: %s, got: %s", tt.expectedErrorMsg, e.Error())
			}

			if !reflect.DeepEqual(tt.expectedData, actualData) {
				t.Errorf("data is not correct, \nwant: \n%+v, \ngot: \n%+v", tt.expectedData, actualData)
			}
		})
	}
}
