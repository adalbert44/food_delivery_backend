package controllers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strings"
)

func getDBConnection() (*sqlx.DB, error){
	connParams := strings.Join([]string{
		"parseTime=true",
		"interpolateParams=true",
		"timeout=10s",
		"collation_server=utf8_general_ci",
		"sql_select_limit=18446744073709551615",
		"compile_only=false",
		"enable_auto_profile=false",
		"sql_mode='STRICT_ALL_TABLES,ONLY_FULL_GROUP_BY'",
	}, "&")

	defaultConfigString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/information_schema?%s",
		os.Getenv("MEMSQL_USER"),
		os.Getenv("MEMSQL_PASSWORD"),
		os.Getenv("MEMSQL_HOST"),
		os.Getenv("MEMSQL_PORT"),
		connParams,
	)

	return sqlx.Open("mysql", defaultConfigString)
}
