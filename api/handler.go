package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	rh "github.com/rinosukmandityo/user-profile/repositories/helper"
	"github.com/rinosukmandityo/user-profile/services/logic"
)

var (
	sessionCheckPath = map[string]bool{
		"/updateprofile": true,
		"/mainprofile":   true,
		"/logout":        true,
	}
	pathPageMap = map[string]string{
		"/":               "index.html",
		"/updateprofile":  "update-profile.html",
		"/mainprofile":    "main-profile.html",
		"/forgotpassword": "forgot-password.html",
		"/notregistered":  "user-not-registered.html",
		"/emailduplicate": "email-duplicate.html",
		"/resetsuccess":   "reset-success.html",
		"/forgotsuccess":  "forgot-success.html",
	}
	baseViewURL = "ui/views"
)

func RegisterHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	userRepo, sessionRepo, tokenRepo := rh.ChooseRepo()
	pageHandler := NewPageHandler(logic.NewSessionService(sessionRepo))
	userHandler := NewUserHandler(logic.NewUserService(userRepo), logic.NewSessionService(sessionRepo))
	loginHandler := NewLoginHandler(logic.NewSessionService(sessionRepo), logic.NewUserService(userRepo))
	signupHandler := NewSignUpHandler(logic.NewSessionService(sessionRepo), logic.NewUserService(userRepo))
	changePasswordHandler := NewChangePassword(logic.NewTokenService(tokenRepo, userRepo))
	r.Use(pageHandler.CheckSession)

	registerUserHandler(r, userHandler)
	registerLoginHandler(r, loginHandler)
	registerSignupHandler(r, signupHandler)
	registerPageHandler(r, pageHandler)
	registerChangePasswordHandler(r, changePasswordHandler)

	return r
}

func registerUserHandler(r *chi.Mux, handler UserHandler) {
	r.Route("/user", func(r chi.Router) {
		r.Use(handler.UserCtx)
		r.Get("/", handler.Get)    // GET /user/
		r.Put("/", handler.Update) // PUT /user/
	})
}

func registerPageHandler(r *chi.Mux, handler PageHandler) {
	for pathURL := range pathPageMap {
		r.HandleFunc(pathURL, handler.CommonPageHandler)
	}
}

func registerLoginHandler(r *chi.Mux, handler LoginHandler) {
	r.Post("/auth", handler.Auth)
	r.HandleFunc("/googlelogin", handler.GoogleLogin)
	r.HandleFunc("/googlecallback", handler.GoogleCallback)
	r.HandleFunc("/logout", handler.Logout)
}

func registerSignupHandler(r *chi.Mux, handler SignupHandler) {
	r.Post("/dosignup", handler.DoSignUp)
	r.HandleFunc("/signup", handler.SignUp)
	r.HandleFunc("/googlesignup", handler.GoogleSignUp)
	r.HandleFunc("/googlesignupcallback", handler.GoogleSignUpCallback)
}

func registerChangePasswordHandler(r *chi.Mux, handler ChangePassword) {
	r.HandleFunc("/resetlink", handler.ResetLink)
	r.HandleFunc("/resetpassword", handler.ResetPassword)
	r.HandleFunc("/changepassword", handler.ChangePassword)
}
