package models

import (
	"database/sql"
	"strconv"
	"wastecontrol/db"
)

func CheckCount(rows *sql.Rows) int {
	var count int
	for rows.Next() {
		count++
	}
	return count
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Bin2hex(str string) (string, error) {
	i, err := strconv.ParseInt(str, 2, 0)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(i, 16), nil
}

func GetAPILogins() map[string]string {
	m := make(map[string]string)

	rows, _ := db.Conn.Query("SELECT username, password FROM endpoint")
	for rows.Next() {
		var username string
		var password string
		rows.Scan(&username, &password)
		m[username] = password
	}

	defer rows.Close()

	return m
}
