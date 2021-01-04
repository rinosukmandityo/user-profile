package repohelper

import (
	"log"
	"os"
	"strconv"

	repo "github.com/rinosukmandityo/user-profile/repositories"
	mr "github.com/rinosukmandityo/user-profile/repositories/mysql"
)

const (
	MONGO_DRIVER = "mongo"
)

func ChooseRepo() (repo.UserRepository, repo.SessionRepository, repo.TokenRepository) {
	url := os.Getenv("url")
	db := os.Getenv("db")
	timeout, _ := strconv.Atoi(os.Getenv("timeout"))
	switch os.Getenv("driver") {
	case MONGO_DRIVER:
		// return mongo repo here
	default:
		if url == "" {
			url = "user:Password.1@tcp(127.0.0.1:3306)/users"
		}
		if db == "" {
			db = "users"
		}
		if timeout == 0 {
			timeout = 10
		}
		userRepo, e := mr.NewUserRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}
		sessionRepo, e := mr.NewSessionRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}
		tokenRepo, e := mr.NewTokenRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}

		return userRepo, sessionRepo, tokenRepo
	}
	return nil, nil, nil
}
