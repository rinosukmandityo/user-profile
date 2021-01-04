package mysql

import (
	"fmt"
	"strings"

	m "github.com/rinosukmandityo/user-profile/models"
)

func constructUpdateQuery(data, filter map[string]interface{}) (string, []interface{}) {
	// 	"UPDATE <tablename> SET field1=?, field2=?  WHERE filter1=?"
	q := fmt.Sprintf("UPDATE %s SET", new(m.User).TableName())
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

func constructDeleteQuery(filter map[string]interface{}) (string, []interface{}) {
	// 	"DELETE <tablename> WHERE filter1=?"
	q := fmt.Sprintf("DELETE FROM %s WHERE", new(m.User).TableName())
	values := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructStoreQuery(data *m.User) (string, []interface{}) {
	// 	"INSERT INTO <tablename> VALUES(?, ?, ?, ?)"
	dataFields := data.SplitByField()
	fieldName := "("
	qVal := "VALUES("
	values := []interface{}{}
	for _, v := range dataFields {
		for key, val := range v {
			fieldName += key + ","
			values = append(values, val)
		}
		qVal += "?,"
	}
	fieldName = strings.TrimSuffix(fieldName, ",") + ")"
	qVal = strings.TrimSuffix(qVal, ",") + ")"
	q := fmt.Sprintf("INSERT INTO %s %s %s", data.TableName(), fieldName, qVal)

	return q, values
}

func constructGetBy(filter map[string]interface{}) (string, []interface{}) {
	// SELECT * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.User).TableName())
	dataFields := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		dataFields = append(dataFields, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, dataFields
}

func constructGetAll() string {
	return "select * from users"
}
