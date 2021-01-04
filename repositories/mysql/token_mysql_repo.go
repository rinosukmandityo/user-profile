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

type tokenMySQLRepository struct {
	url     string
	timeout time.Duration
}

func newTokenClient(URL string) (*sql.DB, error) {
	db, e := sql.Open("mysql", URL)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		return nil, e
	}

	return db, e
}

func (r *tokenMySQLRepository) createNewTable() error {
	tablename := new(m.Token).TableName()
	schema := `CREATE TABLE ` + tablename + ` (
		ID VARCHAR(30) NOT NULL UNIQUE,
		UserID VARCHAR(30),
		Created TIMESTAMP,
		Expired TIMESTAMP  DEFAULT '1970-01-01 00:00:01',
		IsClaimed boolean
	);`

	db, e := newTokenClient(r.url)

	if e != nil {
		return errors.Wrap(e, "repository.Token.CreateTable")
	}
	defer db.Close()
	res, e := db.Exec(schema)
	if res != nil && e == nil {
		fmt.Printf("Table '%s' created\n", tablename)
	}
	return nil
}

func NewTokenRepository(URL, DB string, timeout int) (repo.TokenRepository, error) {
	repo := &tokenMySQLRepository{
		url:     fmt.Sprintf("%s?parseTime=true", URL),
		timeout: time.Duration(timeout) * time.Second,
	}
	repo.createNewTable()
	return repo, nil
}

func (r *tokenMySQLRepository) GetLatestToken(userid string) (m.Token, error) {
	res := []m.Token{}
	token := m.Token{}
	db, e := newTokenClient(r.url)
	if e != nil {
		return token, errors.Wrap(e, "repository.Token.GetLatestToken")
	}
	defer db.Close()
	q, dataFields := constructTokenGetLatest(map[string]interface{}{"UserID": userid})

	results, e := db.Query(q, dataFields...)
	if e != nil {
		return token, errors.Wrap(e, "repository.Token.GetLatestToken")
	}
	for results.Next() {
		var item m.Token
		if e := results.Scan(&item.ID, &item.UserID, &item.Created, &item.Expired, &item.IsClaimed); e != nil {
			return token, errors.Wrap(e, "repository.Token.GetLatestToken")
		}
		res = append(res, item)
	}
	if len(res) > 0 {
		token = res[0]
	} else {
		return token, helper.ErrTokenNotFound
	}
	return token, nil

}
func (r *tokenMySQLRepository) Store(data *m.Token) error {
	db, e := newTokenClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Store")
	}

	q, dataField := constructTokenStoreQuery(data)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Store")
	}
	defer stmt.Close()
	if _, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.Token.Store")
	}

	return nil

}

func (r *tokenMySQLRepository) Update(data map[string]interface{}, id string) error {
	var e error
	db, e := newTokenClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Update")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Update")
	}
	defer conn.Close()

	filter := map[string]interface{}{"ID": id}
	q, dataField := constructTokenUpdateQuery(data, filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.Token.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataField...); e != nil {
		return errors.Wrap(e, "repository.Token.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.Token.Update")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrUserNotFound, "repository.Token.Update")
		}
	}

	return nil

}

func (r *tokenMySQLRepository) DeleteAll() error {
	db, e := newTokenClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.Token.DeleteAll")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.Token.DeleteAll")
	}
	defer conn.Close()

	stmt, e := conn.PrepareContext(ctx, fmt.Sprintf("DELETE FROM %s ", new(m.Token).TableName()))
	if e != nil {
		return errors.Wrap(e, "repository.Token.DeleteAll")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(); e != nil {
		return errors.Wrap(e, "repository.Token.DeleteAll")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.Token.DeleteAll")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrTokenNotFound, "repository.Token.DeleteAll")
		}
	}

	return nil

}
