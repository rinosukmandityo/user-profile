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
	"path"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	svc "github.com/rinosukmandityo/user-profile/services"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type SignupHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	DoSignUp(http.ResponseWriter, *http.Request)
	GoogleSignUp(http.ResponseWriter, *http.Request)
	GoogleSignUpCallback(w http.ResponseWriter, r *http.Request)
}

type signupHandler struct {
	sessionService svc.SessionService
	userService    svc.UserService
}

var (
	googleSignupOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_SIGNUP_REDIRECT"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)

func init() {
	if oauthStateString == "" {
		timestamp := time.Now()
		id := sha256.Sum256([]byte(fmt.Sprintf("%s", timestamp.String())))
		oauthStateString = hex.EncodeToString(id[:])
	}
}

func NewSignUpHandler(sessionService svc.SessionService, userService svc.UserService) SignupHandler {
	return &signupHandler{sessionService, userService}
}

func (u *signupHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	registerPage(w, path.Join("ui/views", "signup.html"), nil)
}

func (u *signupHandler) GoogleSignUp(w http.ResponseWriter, r *http.Request) {
	url := googleSignupOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (u *signupHandler) DoSignUp(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	result := helper.NewResult(nil)
	requestBody, e := ioutil.ReadAll(r.Body)
	statusCode := http.StatusOK
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	user, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := u.userService.Store(user); e != nil {
		if errors.Cause(e) == helper.ErrEmailDuplicate {
			result.SetErrMsg(helper.ErrDuplicateEmail)
		} else {
			result.SetErrMsg(e.Error())
		}
		statusCode = http.StatusBadRequest
	}
	if result.Success {
		tSession, e := u.sessionService.CreateNewSession(*user)
		if e != nil {
			result.SetError(e)
		}
		helper.SetCookie(w, r, helper.SESSION_COOKIE_KEY, tSession.ID, time.Until(tSession.Expired))
		helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, tSession.Email, time.Until(tSession.Expired))

		result.SetMessage(helper.SuccessSignup)
	}
	respBody, e := GetSerializer(contentType).EncodeResult(result)
	if e != nil {
		statusCode = http.StatusBadRequest
	}
	SetupResponse(w, contentType, respBody, statusCode)
}

func (u *signupHandler) GoogleSignUpCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthStateString {
		log.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, e := googleSignupOauthConfig.Exchange(context.Background(), r.FormValue("code"))
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
	user := &m.User{
		Name:  name,
		Email: email,
	}
	if e := u.userService.Store(user); e != nil {
		if errors.Cause(e) == helper.ErrEmailDuplicate {
			log.Println(helper.ErrEmailDuplicate)
			http.Redirect(w, r, "/emailduplicate", http.StatusTemporaryRedirect)
			return
		}
	}

	tSession, e := u.sessionService.CreateNewSession(*user)
	if e != nil {
		log.Printf("Could not create new session: %s\n", e.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	helper.SetCookie(w, r, helper.SESSION_COOKIE_KEY, tSession.ID, time.Until(tSession.Expired))
	helper.SetCookie(w, r, helper.EMAIL_COOKIE_KEY, tSession.Email, time.Until(tSession.Expired))

	http.Redirect(w, r, "/updateprofile", http.StatusTemporaryRedirect)
}
