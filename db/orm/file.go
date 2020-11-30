package orm

import (
	"database/sql"
	"log"
	mydb "distributed_file/db/conn"
)

func OnFileUploadFinished(filehash, filename string,
	filesize int64, fileaddr string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file" +
			"(`file_sha1`, `file_name`, `file_size`" +
			"`file_addr`, `status`) values (?,?,?,?,1)")
	if err != nil{
		log.Println("Failed to prepare statement, err:" + err.Error())
		res.Suc = false
		return
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		return
	}
	if rf, err := ret.RowsAffected(); err == nil{
		if rf <= 0{
			log.Printf("File with hash:%s has been uploaded before", filehash)
		}
		res.Suc = true
		return
	}
	res.Suc = false
	return
}

func GetFileMeta(filehash string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file \" +\n\t\t\t\"where file_sha1=? and status=1 limit 1")
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			// 查不到对应记录， 返回参数及错误均为nil
			res.Suc = true
			res.Data = nil
			return
		} else {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = tfile
	return

}

func GetFileMetaList(limit int64)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file" +
			"where status=1 limit ?")
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	var tfiles []TableFile
	for i:=0 ; rows.Next();i++{
		tfile := TableFile{}
		err = rows.Scan(&tfile.FileHash, &tfile.FileAddr,
			&tfile.FileName, &tfile.FileSize)
		if err != nil{
			log.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	res.Suc = true
	res.Data = tfiles
	return

}

func UpdateFileLocation(filehash, fileaddr string)(res ExecResult){
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_file set`file_addr`=? where `file_sha1`=? limit 1")
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileaddr, filehash)
	if err != nil{
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	if rf, err := ret.RowsAffected();err == nil{
		if rf <= 0{
			log.Printf("更新文件location失败, filehash:%s", filehash)
			res.Suc = false
			res.Msg = "No record updated"
		}
		res.Suc = true
		return
	}
	res.Suc = false
	res.Msg = err.Error()
	return
}