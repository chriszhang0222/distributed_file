package orm

import (
	mydb "distributed_file/db/conn"
	"log"
)

func UserSignUp(username, password string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil{
		log.Println("Failed to insert, err:" + err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rowsAffected, err := ret.RowsAffected();nil == err && rowsAffected >= 0{
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = "User already exists"
	return

}

func UserSignin(username, passwd string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare("" +
		"select * from tbl_user where user_name=? limit 1")
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil{
		res.Suc = false
		res.Msg = err.Error()
		return
	}else if rows == nil{
		res.Suc = false
		res.Msg = "用户名未注册"
	}
	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == passwd{
		res.Suc = true
		res.Data = true
		return
	}
	res.Suc = false
	res.Msg = "用户名/密码不匹配"
	return
}

func GetUserInfo(username string)(res ExecResult){
	user := TableUser{}
	stmt, err := mydb.DBConn().Prepare(
		"select user_name,email, signup_at from tbl_user where user_name=? limit 1")
	if err != nil{
		log.Println(err.Error())
		// error不为nil, 返回时user应当置为nil
		//return user, err
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(username).Scan(&user.Username, &user.Email, &user.SignupAt)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	res.Data = user
	return
}

func UserExist(user string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tbl_user where user_name=? limit 1")
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(user)
	res.Suc = true
	res.Data = map[string]bool{
		"exists": rows.Next(),
	}
	return
}

