package handler

import (
	"distributed_file/meta"
	"distributed_file/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"encoding/json"
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

		newFile, err := os.Create("/User/chris/tmp/" + head.Filename)
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
		meta.UpdateFileMeta(filemeta)
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
	fm := meta.GetFileMeta(fileHash)
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
