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

type userMySQLRepository struct {
	url     string
	timeout time.Duration
}

func newUserClient(URL string) (*sql.DB, error) {
	db, e := sql.Open("mysql", URL)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		return nil, e
	}

	return db, e
}

func (r *userMySQLRepository) createNewTable() error {
	tablename := new(m.User).TableName()
	schema := `CREATE TABLE ` + tablename + ` (
		ID VARCHAR(30) NOT NULL UNIQUE,
		Name VARCHAR(30),
		Email VARCHAR(50),
		Password VARCHAR(50),
		Telephone VARCHAR(30),
		Address VARCHAR(50),
		IsActive boolean,
		IsGoogleAuth boolean
	);`
	db, e := newUserClient(r.url)

	if e != nil {
		return errors.Wrap(e, "repository.User.CreateTable")
	}
	defer db.Close()
	res, e := db.Exec(schema)
	if res != nil && e == nil {
		fmt.Printf("Table '%s' created\n", tablename)
	}
	return nil
}

func NewUserRepository(URL, DB string, timeout int) (repo.UserRepository, error) {
	repo := &userMySQLRepository{
		url:     fmt.Sprintf("%s?parseTime=true", URL),
		timeout: time.Duration(timeout) * time.Second,
	}
	repo.createNewTable()
	return repo, nil
}

func (r *userMySQLRepository) GetBy(filter map[string]interface{}) (m.User, error) {
	res := m.User{}
	db, e := newUserClient(r.url)
	if e != nil {
		return res, errors.Wrap(e, "repository.User.GetBy")
	}
	defer db.Close()
	q, dataFields := constructGetBy(filter)

	if e = db.QueryRow(q, dataFields...).Scan(&res.ID, &res.Name, &res.Email, &res.Password, &res.Telephone, &res.Address, &res.IsActive, &res.IsGoogleAuth); e != nil {
		if e == sql.ErrNoRows {
			return res, errors.Wrap(helper.ErrUserNotFound, "repository.User.GetBy")
		}
		return res, errors.Wrap(e, "repository.User.GetBy")
	}

	return res, nil

}
func (r *userMySQLRepository) Store(data *m.User) error {
	db, e := newUserClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	q, dataField := constructStoreQuery(data)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}
	defer stmt.Close()
	if _, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.User.Store")
	}

	return nil

}

func (r *userMySQLRepository) Update(data map[string]interface{}, id string) (m.User, error) {
	user := m.User{}
	var e error
	db, e := newUserClient(r.url)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	}
	defer conn.Close()

	filter := map[string]interface{}{"ID": id}
	q, dataField := constructUpdateQuery(data, filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataField...); e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return user, errors.Wrap(e, "repository.User.Update")
		}
		if count == 0 {
			return user, errors.Wrap(helper.ErrUserNotFound, "repository.User.Update")
		}
	}
	user, e = r.GetBy(filter)
	if e != nil {
		return user, errors.Wrap(e, "repository.User.Update")
	}

	return user, nil

}

func (r *userMySQLRepository) DeleteAll() error {
	db, e := newUserClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.User.DeleteAll")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.User.DeleteAll")
	}
	defer conn.Close()

	stmt, e := conn.PrepareContext(ctx, fmt.Sprintf("DELETE FROM %s ", new(m.User).TableName()))
	if e != nil {
		return errors.Wrap(e, "repository.User.DeleteAll")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(); e != nil {
		return errors.Wrap(e, "repository.User.DeleteAll")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.User.DeleteAll")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrUserNotFound, "repository.User.DeleteAll")
		}
	}

	return nil
}
