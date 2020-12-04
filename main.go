package main

import (
	"distributed_file/db/orm"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)
func main(){

	exc := orm.GetUserInfo("chris")
	fmt.Println(exc.Data)


}
