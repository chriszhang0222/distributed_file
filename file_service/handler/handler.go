package handler

import (
	"distributed_file/file_service/meta"
	"distributed_file/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	db "distributed_file/db/orm"
)
import "io/ioutil"


func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil{
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST"{
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
		filemeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/Users/chris/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:00:00"),
		}

		newFile, err := os.Create("/Users/chris/tmp/" + head.Filename)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		defer newFile.Close()
		filemeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		newFile.Seek(0,0)
		filemeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMetaDB(filemeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}

}

func UploadSucceedhandler(w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Upload Completed!")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(fileHash)
	data, err := json.Marshal(fMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	fm, err := meta.GetFileMetaFromDB(fileHash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	location := fm.Location
	f, err := os.Open(location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data,err := ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment;filename=\"" + fm.FileName + "\"")
	w.Write(data)

}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newfilename := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newfilename
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filesha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}

func TryFastUploadhandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))
	fileMeta, err := meta.GetFileMetaFromDB(filehash)
	if err != nil{
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if fileMeta == nil{
		resp := util.RespMsg{
			Code: -1,
			Msg: "Failed to fast upload",
		}
		w.Write(resp.JSONBytes())
		return
	}
	suc := db.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc.Suc{
		resp := util.RespMsg{
			Code: 0,
			Msg: "Dast Upload Succeed",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg: "Fastupload failed",
	}
	fmt.Println(suc.Msg)
	w.Write(resp.JSONBytes())
	return
}

func DownLoadURLHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	row, err := db.GetFileMeta(filehash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tableData := row.Data.(db.TableFile)
	if strings.HasPrefix(tableData.FileAddr.String, "/tmp"){
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			r.Host, filehash, username, token)
		w.Write([]byte(tmpUrl))
	}else if strings.HasPrefix(tableData.FileAddr.String, "/ceph"){
		//TODO
	}else if strings.HasPrefix(tableData.FileAddr.String, "oss/"){
		//TODO
	}


}
