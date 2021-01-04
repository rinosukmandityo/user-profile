package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	svc "github.com/rinosukmandityo/user-profile/services"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type LoginHandler interface {
	Auth(http.ResponseWriter, *http.Request)
	GoogleLogin(http.ResponseWriter, *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type loginHandler struct {
	sessionService svc.SessionService
	userService    svc.UserService
}

const (
	GOOGLE_CALLBACK_URL = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_LOGIN_REDIRECT"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = ""
)

func init() {
	if oauthStateString == "" {
		timestamp := time.Now()
		id := sha256.Sum256([]byte(fmt.Sprintf("%s", timestamp.String())))
		oauthStateString = hex.EncodeToString(id[:])
	}
}

func NewLoginHandler(sessionService svc.SessionService, userService svc.UserService) LoginHandler {
	return &loginHandler{sessionService, userService}
}

func (u *loginHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (u *loginHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		log.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, e := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
	if e != nil {
		log.Printf("Could not get token: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, e := http.Get(fmt.Sprintf("%s%s", GOOGLE_CALLBACK_URL, token.AccessToken))
	if e != nil {
		log.Printf("Could not create get request: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	contents, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Printf("Could not parse response: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	respMap, e := GetSerializer(ContentTypeJson).DecodeMap(contents)
	if e != nil {
		log.Printf("Could not decode response: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	email := respMap["email"].(string)
	name := ""
	if nameInt, hasData := respMap["name"]; hasData {
		name = nameInt.(string)
	}
	user := m.User{
		Name:  name,
		Email: email,
	}
	if _, _, e = u.userService.GetByEmail(user.Email); e != nil {
		if errors.Cause(e) == helper.ErrUserNotFound {
			log.Println(helper.ErrAuthEmailMsg)
			http.Redirect(w, r, "/notregistered", http.StatusTemporaryRedirect)
			return
		}
	}

	tSession, e := u.sessionService.CreateNewSession(user)
	if e != nil {
		log.Printf("Could not create new session: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	helper.SetCookie(w, r, helper.SESSION_COOKIE_KEY, tSession.ID, time.Until(tSession.Expired))
	helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, tSession.Email, time.Until(tSession.Expired))

	http.Redirect(w, r, "/updateprofile", http.StatusTemporaryRedirect)
}

func (u *loginHandler) Auth(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	result := helper.NewResult(nil)
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	user, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		return
	}
	if _, _, e = u.sessionService.Authenticate(user.Email, user.Password); e != nil {
		if errors.Cause(e) == helper.ErrUserNotFound {
			ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrAuthEmailMsg), http.StatusNotFound)
		} else if errors.Cause(e) == helper.ErrPasswordDoesNotMatch {
			ResponseWithResult(w, contentType, result.SetErrMsg(helper.ErrAuthPasswordMsg), http.StatusNotFound)
		} else {
			ResponseWithResult(w, contentType, result.SetError(e), http.StatusBadRequest)
		}
		return
	}
	tSession, e := u.sessionService.CreateNewSession(*user)
	if e != nil {
		ResponseWithResult(w, contentType, result.SetError(e), http.StatusInternalServerError)
		return
	}
	helper.SetCookie(w, r, helper.SESSION_COOKIE_KEY, tSession.ID, time.Until(tSession.Expired))
	helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, tSession.Email, time.Until(tSession.Expired))

	ResponseWithResult(w, contentType, result.SetMessage(helper.SuccessLogin), http.StatusOK)
}

func (u *loginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	sess, ok := ctx.Value("session").(m.Session)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if e := u.sessionService.Logout(sess); e != nil {
		log.Printf("Error on session logout %s\n", e.Error())
	}
	helper.SetCookie(w, r, helper.SESSION_COOKIE_KEY, sess.ID, time.Minute*-30)
	helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, sess.Email, time.Minute*-30)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
