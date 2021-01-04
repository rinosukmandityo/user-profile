package api

import (
	"context"
	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	svc "github.com/rinosukmandityo/user-profile/services"
	"html/template"
	"log"
	"net/http"
	"path"
)

type PageHandler interface {
	CheckSession(http.Handler) http.Handler
	CommonPageHandler(w http.ResponseWriter, r *http.Request)
}

type pageHandler struct {
	sessionSvc svc.SessionService
}

func NewPageHandler(sessionSvc svc.SessionService) PageHandler {
	return &pageHandler{sessionSvc}
}

func (u *pageHandler) CheckSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "no-cache, no-store, max-age=0, must-revalidate")
		w.Header().Set("pragma", "no-cache")
		w.Header().Set("expires", "Sat, 01 Jan 1990 00:00:00 GMT")

		urlPath := r.URL.Path
		ctx := r.Context()
		var isSessionChecked bool

		if sessionCheckPath[urlPath] {
			sessCookie, e := r.Cookie(helper.SESSION_COOKIE_KEY)
			if e != nil || sessCookie == nil {
				log.Printf("Error on get session ID on cookies: %s\n", e.Error())
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			isActive, sess, e := u.sessionSvc.IsSessionActive(sessCookie.Value)
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
			ctx = context.WithValue(ctx, "session", sess)
			isSessionChecked = true
		}

		ctx = context.WithValue(ctx, "sessionChecked", isSessionChecked)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *pageHandler) CommonPageHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	ctx := r.Context()
	isSessionChecked, _ := ctx.Value("sessionChecked").(bool)
	if isSessionChecked {
		_, ok := ctx.Value("session").(m.Session)
		if !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
	}
	registerPage(w, path.Join(baseViewURL, pathPageMap[urlPath]), nil)
}

func registerPage(w http.ResponseWriter, fpath string, data map[string]interface{}) {
	var tmpl, e = template.ParseFiles(fpath)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	if e = tmpl.Execute(w, data); e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
