package mysql

import (
	"fmt"
	"strings"

	m "github.com/rinosukmandityo/user-profile/models"
)

func constructTokenUpdateQuery(data, filter map[string]interface{}) (string, []interface{}) {
	// 	"UPDATE <tablename> SET field1=?, field2=?  WHERE filter1=?"
	q := fmt.Sprintf("UPDATE %s SET", new(m.Token).TableName())
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

func constructTokenStoreQuery(data *m.Token) (string, []interface{}) {
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

func constructTokenGetBy(filter map[string]interface{}) (string, []interface{}) {
	// SELECT * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.Token).TableName())
	dataFields := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, dataFields
}

func constructTokenGetLatest(filter map[string]interface{}) (string, []interface{}) {
	// SELECT * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.Token).TableName())
	dataFields := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	q = strings.TrimSuffix(q, ",")
	q += " ORDER BY Expired DESC"

	return q, dataFields
}
