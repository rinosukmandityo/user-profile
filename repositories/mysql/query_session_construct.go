package mysql

import (
	"fmt"
	"strings"

	m "github.com/rinosukmandityo/user-profile/models"
)

func constructSessionUpdateQuery(data, filter map[string]interface{}) (string, []interface{}) {
	// 	"UPDATE <tablename> SET field1=?, field2=?  WHERE filter1=?"
	q := fmt.Sprintf("UPDATE %s SET", new(m.Session).TableName())
	values := []interface{}{}
	for k, v := range data {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")
	q += " WHERE"
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructSessionDeleteQuery(filter map[string]interface{}) (string, []interface{}) {
	// 	"DELETE <tablename> WHERE filter1=?"
	q := fmt.Sprintf("DELETE FROM %s WHERE", new(m.Session).TableName())
	values := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructSessionStoreQuery(data *m.Session) (string, []interface{}) {
	// 	"INSERT INTO <tablename> VALUES(?, ?, ?, ?)"
	dataFields := data.SplitByField()
	q := fmt.Sprintf("INSERT INTO %s VALUES(", data.TableName())
	values := []interface{}{}
	for _, v := range dataFields {
		q += "?,"
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",") + ")"

	return q, values
}

func constructSessionGetBy(filter map[string]interface{}) (string, []interface{}) {
	// SELECT * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.Session).TableName())
	dataFields := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, dataFields
}

func constructAuth(filter map[string]interface{}) (string, []interface{}) {
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.User).TableName())
	count := 0
	dataFields := []interface{}{}
	for k, v := range filter {
		if count == 0 {
			q += fmt.Sprintf(" %s=?", k)
		} else {
			q += fmt.Sprintf(" OR %s=?", k)
		}
		dataFields = append(dataFields, v)
		count++
	}
	return q, dataFields
}
