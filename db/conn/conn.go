package conn

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var (
	MYSQLSource = "root:31415926@tcp(192.168.0.10:3306)/fileserver?charset=utf8"
)

var db *sql.DB

func init(){
	db, _ = sql.Open("mysql", MYSQLSource)
	db.SetMaxOpenConns(100)
	err := db.Ping()
	if err != nil{
		fmt.Println("Failded to connect to MySQL :" + err.Error())
		os.Exit(1)
	}
}

func ParseRows(rows *sql.Rows)[]map[string]interface{}{
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j:= range values{
		scanArgs[j] = &values[j]
	}
	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)

	for rows.Next(){
		err := rows.Scan(scanArgs...)
		if err != nil{
			panic(err)
		}
		for i, col := range values{
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records

}

func DBConn() *sql.DB{
	return db
}