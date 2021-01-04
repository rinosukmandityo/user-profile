package helper

import (
	"net/http"
	"net/url"
	"time"
)

const (
	SESSION_COOKIE_KEY = "profileSessionID"
	EMAIL_COOKIE_KEY   = "profileEmail"
)

func SetCookie(w http.ResponseWriter, r *http.Request, name, value string, expiresAfter time.Duration) *http.Cookie {
	c := &http.Cookie{Name: name, Value: value, HttpOnly: true}
	c.Path = "/"
	u, e := url.Parse(r.URL.String())
	if e == nil {
		c.Expires = time.Now().Add(expiresAfter)
		c.Domain = u.Host
	}
	http.SetCookie(w, c)

	return c
}
