package dbops

import (
	"database/sql"
	"fmt"

	// mysql 初始化
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	dbConn, err = sql.Open("mysql", "haibei:h9420x@tcp(120.79.57.219:3306)/videowebserver?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("dbConn +%v\n", dbConn)
}
