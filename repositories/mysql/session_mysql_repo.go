package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/rinosukmandityo/user-profile/helper"
	m "github.com/rinosukmandityo/user-profile/models"
	repo "github.com/rinosukmandityo/user-profile/repositories"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type sessionMySQLRepository struct {
	url     string
	timeout time.Duration
}

func newSessionClient(URL string) (*sql.DB, error) {
	db, e := sql.Open("mysql", URL)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		return nil, e
	}

	return db, e
}

func (r *sessionMySQLRepository) createNewTable() error {
	tablename := new(m.Session).TableName()
	schema := `CREATE TABLE ` + tablename + ` (
		ID VARCHAR(30) NOT NULL UNIQUE,
		UserID VARCHAR(30),
		Email VARCHAR(50),
		Created TIMESTAMP,
		Expired TIMESTAMP DEFAULT '1970-01-01 00:00:01'
	);`
	db, e := newSessionClient(r.url)

	if e != nil {
		return errors.Wrap(e, "repository.Session.CreateTable")
	}
	defer db.Close()
	res, e := db.Exec(schema)

	if res != nil && e == nil {
		fmt.Printf("Table '%s' created\n", tablename)
	}
	return nil
}

func NewSessionRepository(URL, DB string, timeout int) (repo.SessionRepository, error) {
	repo := &sessionMySQLRepository{
		url:     fmt.Sprintf("%s?parseTime=true", URL),
		timeout: time.Duration(timeout) * time.Second,
	}
	repo.createNewTable()
	return repo, nil
}

func (r *sessionMySQLRepository) GetBy(filter map[string]interface{}) (m.Session, error) {
	res := m.Session{}
	db, e := newSessionClient(r.url)
	if e != nil {
		return res, errors.Wrap(e, "repository.Session.GetBy")
	}
	defer db.Close()
	q, dataFields := constructSessionGetBy(filter)

	if e = db.QueryRow(q, dataFields...).Scan(&res.ID, &res.UserID, &res.Email, &res.Created, &res.Expired); e != nil {
		if e == sql.ErrNoRows {
			return res, errors.Wrap(helper.ErrUserNotFound, "repository.Session.GetBy")
		}
		return res, errors.Wrap(e, "repository.Session.GetBy")
	}

	return res, nil

}
func (r *sessionMySQLRepository) Store(data *m.Session) error {
	db, e := newSessionClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.Session.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.Session.Store")
	}

	q, dataField := constructSessionStoreQuery(data)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.Session.Store")
	}
	defer stmt.Close()
	if _, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.Session.Store")
	}

	return nil

}

func (r *sessionMySQLRepository) Update(data map[string]interface{}, id string) (m.Session, error) {
	user := m.Session{}
	var e error
	db, e := newSessionClient(r.url)
	if e != nil {
		return user, errors.Wrap(e, "repository.Session.Update")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return user, errors.Wrap(e, "repository.Session.Update")
	}
	defer conn.Close()

	filter := map[string]interface{}{"ID": id}
	q, dataField := constructSessionUpdateQuery(data, filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return user, errors.Wrap(e, "repository.Session.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataField...); e != nil {
		return user, errors.Wrap(e, "repository.Session.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return user, errors.Wrap(e, "repository.Session.Update")
		}
		if count == 0 {
			return user, errors.Wrap(helper.ErrUserNotFound, "repository.Session.Update")
		}
	}
	user, e = r.GetBy(filter)
	if e != nil {
		return user, errors.Wrap(e, "repository.Session.Update")
	}

	return user, nil

}

func (r *sessionMySQLRepository) Authenticate(email, password string) (bool, m.User, error) {
	res := m.User{}
	db, e := newUserClient(r.url)
	if e != nil {
		return false, res, errors.Wrap(e, "repository.Session.Authenticate")
	}
	defer db.Close()
	q, dataFields := constructAuth(map[string]interface{}{"Email": email})

	if e = db.QueryRow(q, dataFields...).Scan(&res.ID, &res.Name, &res.Email, &res.Password, &res.Telephone, &res.Address, &res.IsActive, &res.IsGoogleAuth); e != nil {
		if e == sql.ErrNoRows {
			return false, res, errors.Wrap(helper.ErrUserNotFound, "repository.Session.Authenticate")
		}
		return false, res, errors.Wrap(e, "repository.Session.Authenticate")
	}

	if res.ID == "" {
		return false, res, errors.Wrap(helper.ErrUserNotFound, "repository.Session.Authenticate")
	}

	if !repo.IsPasswordMatch(password, res.Password) {
		return false, res, errors.Wrap(helper.ErrPasswordDoesNotMatch, "repository.Session.Authenticate")
	}

	return true, res, nil
}

func (r *sessionMySQLRepository) DeleteAll() error {
	db, e := newSessionClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.Session.DeleteAll")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.Session.DeleteAll")
	}
	defer conn.Close()

	stmt, e := conn.PrepareContext(ctx, fmt.Sprintf("DELETE FROM %s ", new(m.Session).TableName()))
	if e != nil {
		return errors.Wrap(e, "repository.Session.DeleteAll")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(); e != nil {
		return errors.Wrap(e, "repository.Session.DeleteAll")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.Session.DeleteAll")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrSessionNotFound, "repository.Session.DeleteAll")
		}
	}

	return nil

}
