package api

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	svc "github.com/rinosukmandityo/user-profile/services"

	"github.com/pkg/errors"
)

type UserHandler interface {
	UserCtx(http.Handler) http.Handler
	Get(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
}

type userHandler struct {
	userService    svc.UserService
	sessionService svc.SessionService
}

func NewUserHandler(userService svc.UserService, sessionService svc.SessionService) UserHandler {
	return &userHandler{userService, sessionService}
}

func (u *userHandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "no-cache, no-store, max-age=0, must-revalidate")
		w.Header().Set("pragma", "no-cache")
		w.Header().Set("expires", "Sat, 01 Jan 1990 00:00:00 GMT")

		sessCookie, e := r.Cookie(helper.SESSION_COOKIE_KEY)
		if e != nil || sessCookie == nil {
			log.Printf("Error on get session ID on cookies: %s\n", e.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		isActive, _, e := u.sessionService.IsSessionActive(sessCookie.Value)
		if e != nil {
			log.Printf("Error on get session: %s\n", e.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		if !isActive {
			log.Printf("Session expired\n")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		emailCookie, e := r.Cookie(helper.EMAIL_COOKIE_KEY)
		if e != nil || emailCookie == nil {
			log.Printf("Error on get email on cookies: %s\n", e.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		found, user, e := u.userService.GetByEmail(emailCookie.Value)
		if e != nil {
			if errors.Cause(e) == helper.ErrUserNotFound || !found {
				log.Printf("User not found\n")
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "user", &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	result := helper.NewResult(nil)
	contentType := r.Header.Get("Content-Type")
	ctx := r.Context()
	data, ok := ctx.Value("user").(*m.User)
	if !ok {
		ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrUserNotFound.Error()), http.StatusBadRequest)
		return
	}
	data.Password = ""
	ResponseWithResult(w, contentType, result.SetData(data), http.StatusFound)
}

func (u *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	result := helper.NewResult(nil)
	contentType := r.Header.Get("Content-Type")
	ctx := r.Context()
	existingData, ok := ctx.Value("user").(*m.User)
	if !ok {
		ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrUserNotFound.Error()), http.StatusBadRequest)
		return
	}
	id := existingData.ID
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).DecodeMap(requestBody)
	if e != nil {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	user, e := u.userService.Update(data, id)
	if e != nil && errors.Cause(e) != helper.ErrUserNotFound {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusInternalServerError)
		return
	}
	helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, user.Email, time.Duration(time.Minute*30))
	user.Password = ""
	ResponseWithResult(w, contentType, result.SetData(data), http.StatusOK)

}
