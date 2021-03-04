package helpers

import (
	"fmt"
	"math/rand"
	"database/sql"
	def "../definitions"
)

func OnDatabase(isSelect bool, query *string) ([]map[string]interface{}, error) {
	/* creating data source string */
	dataSource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		def.DBUser, def.DBPass, def.DBHost, def.DBPort, def.DBName)
	/* creating connection */
	connection, err := sql.Open(def.DBDriver, dataSource)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	/* executing query and getting result */
	rows, err := connection.Query(*query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	/* if query type is not select*/
	if !isSelect {
		return nil, nil
	}

	/* getting columns names of the query result */
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var count int = len(columns) // column count
	values := make([]map[string]interface{}, 0)
	keys := make([]interface{}, count)
	ptrKeys := make([]interface{}, count)

	/* iterating through rows of result */
	for rows.Next() {
		for index := 0; index < count; index++ {
			ptrKeys[index] = &keys[index]
		}
		rows.Scan(ptrKeys...)

		each := make(map[string]interface{})
		for col, key := range columns {
			var value interface{}
			bytes, ok := keys[col].([]byte)
			if ok {
				value = string(bytes)
			} else {
				value = keys[col]
			}
			each[key] = value
		}
		values = append(values, each)
	}
	return values, nil
}


func sendPlainEmail(emailTo, message string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	auth := smtp.PlainAuth("", def.AppEmail, def.AppPass, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort,
		auth,
		def.AppEmail,
		[]string{
			emailTo
		},
		[]byte(message)
	return err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func getRandomStr() string {
	bytes := make([]byte, def.CodeLen)
	for i := range bytes {
		bytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}

