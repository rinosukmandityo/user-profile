package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	svc "github.com/rinosukmandityo/user-profile/services"

	"github.com/pkg/errors"
)

type ChangePassword interface {
	ResetPassword(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
	ResetLink(w http.ResponseWriter, r *http.Request)
}

type changePassword struct {
	tokenSvc svc.TokenService
}

func NewChangePassword(tokenSvc svc.TokenService) ChangePassword {
	return &changePassword{tokenSvc}
}

func (u *changePassword) ResetPassword(w http.ResponseWriter, r *http.Request) {
	urlParsed, e := url.Parse(r.RequestURI)
	if e != nil {
		return
	}
	q := urlParsed.Query()

	data := map[string]interface{}{}

	if q.Get("e") == "" {
		data["ErrorMessage"] = "Invalid Token"
		registerPage(w, path.Join(baseViewURL, "reset-error.html"), data)
		return
	}

	if q.Get("d") == "" {
		data["ErrorMessage"] = helper.ErrInvalidUserID
		registerPage(w, path.Join(baseViewURL, "reset-error.html"), data)
		return
	}

	isTokenValid, e := u.tokenSvc.IsTokenValid(q.Get("d"), q.Get("e"))
	if e != nil || !isTokenValid {
		data["ErrorMessage"] = e.Error()
		registerPage(w, path.Join(baseViewURL, "reset-error.html"), data)
		return
	}

	registerPage(w, path.Join(baseViewURL, "reset-password.html"), data)
}

func (u *changePassword) ChangePassword(w http.ResponseWriter, r *http.Request) {
	result := helper.NewResult(nil)
	statusCode := http.StatusOK
	contentType := r.Header.Get("Content-Type")

	urlParsed, e := url.Parse(r.Referer())

	if e != nil {
		return
	}
	q := urlParsed.Query()
	tokenID := q.Get("e")
	userID := q.Get("d")

	if tokenID == "" {
		ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrInvalidToken), http.StatusBadRequest)
		return
	}

	if userID == "" {
		ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrInvalidUserID), http.StatusBadRequest)
		return
	}

	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println("Error on reading the request body:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).DecodeMap(requestBody)
	if e != nil {
		log.Println("Error on decoding the request body:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	passwd := data["Password"].(string)

	if e = u.tokenSvc.ChangePasswordToken(userID, passwd, tokenID); e != nil {
		log.Println("Error on changing password with token:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	ResponseWithResult(w, contentType, result.SetData(data), statusCode)
}

func (u *changePassword) ResetLink(w http.ResponseWriter, r *http.Request) {
	result := helper.NewResult(nil)
	statusCode := http.StatusOK
	contentType := r.Header.Get("Content-Type")

	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Println("Error on reading the request body:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).DecodeMap(requestBody)
	if e != nil {
		log.Println("Error on decoding the request body:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	email := data["Email"].(string)
	user, tokenID, e := u.tokenSvc.ResetPasswordByMail(email, time.Minute*15)
	if e != nil {
		if errors.Cause(e) == helper.ErrUserNotFound {
			log.Println("User not found")
			ResponseWithResult(w, contentType, result.SetError(helper.ErrUserNotFound), http.StatusNotFound)
			return
		}
		log.Println("Error on getting token:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	hostURL := helper.ConstructEmailURL(r.Referer(), map[string]string{"suffix": "/resetpassword?e="})
	userMap := map[string]string{"id": user.ID, "fullname": user.Name, "email": user.Email}
	if e = helper.ResetPasswordMailContent(userMap, tokenID, hostURL); e != nil {
		log.Println("Error on sending reset password email:", e.Error())
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}

	ResponseWithResult(w, contentType, result.SetData(data), statusCode)
}
