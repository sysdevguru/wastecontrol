package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

const (
	STAGE_DB_HOST = "tcp(e7qyahb3d90mletd.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306)"
	STAGE_DB_NAME = "cmbwghn06te6f39v"
	STAGE_DB_USER = "dhw2xxcwstxp24sc"
	STAGE_DB_PASS = "w8kgfw59y8hg7xrt"

	PROD_DB_HOST = "tcp(remj7niuux921pss.chr7pe7iynqr.eu-west-1.rds.amazonaws.com:3306)"
	PROD_DB_NAME = "primary_app_db"
	PROD_DB_USER = "nlc3e37fdmimmvl5"
	PROD_DB_PASS = "xoe3sre8hw90o9ej"
)

var Conn = Connect()

/**
 * [Connect is the function who makes the connection to the Database]
 * @return 	sql.DB
 */
func Connect() *sql.DB {
	dsn := PROD_DB_USER + ":" + PROD_DB_PASS + "@" + PROD_DB_HOST + "/" + PROD_DB_NAME + "?charset=utf8"

	if os.Getenv("ENVIRONMENT") == "STAGE" {
		dsn = STAGE_DB_USER + ":" + STAGE_DB_PASS + "@" + STAGE_DB_HOST + "/" + STAGE_DB_NAME + "?charset=utf8"
	}

	conn, _ := sql.Open("mysql", dsn)

	return conn
}
