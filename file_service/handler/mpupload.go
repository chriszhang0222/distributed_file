package handler

import (
	db "distributed_file/db/orm"
	rPool "distributed_file/redis"
	"distributed_file/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	chucksize = 5 * 1024 * 1024
)
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

func getChuckSize(size int) int{
	return int(math.Ceil(float64(size))) / chucksize
}

func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil{
		w.Write(util.NewRespMsg(-1, "invalid parameters", nil).JSONBytes())
	}
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	upInfo := MultipartUploadInfo{
		FileHash: filehash,
		FileSize: filesize,
		UploadID: username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: chucksize,
		ChunkCount: getChuckSize(filesize),
	}
	rConn.Do("HSET", "MP_" + upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_" + upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

func UploadPartHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	uploadId := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	fpath := "Users/chris/tmp/data/" + uploadId + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil{
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024 * 1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil{
			break
		}
	}
	rConn.Do("HSET", "MP_"+uploadId, "chkidx_"+chunkIndex, 1)
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())

 }

 func CompleteUploadHandler(w http.ResponseWriter, r *http.Request){
 	r.ParseForm()
 	uploadId := r.Form.Get("uploadid")
	 username := r.Form.Get("username")
	 filehash := r.Form.Get("filehash")
	 filesize := r.Form.Get("filesize")
	 filename := r.Form.Get("filename")

	 rConn := rPool.RedisPool().Get()
	 defer rConn.Close()

	 data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadId))
	 if err != nil{
	 	w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
	 	return
	 }
	 totalCount := 0
	 chunkCount := 0
	 for i:=0; i < len(data);i+=2{
	 	k := string(data[i].([]byte))
	 	v := string(data[i+1].([]byte))
	 	if k == "chunkcount" {
	 		totalCount , _ = strconv.Atoi(v)
		}else if strings.HasPrefix(k, "chkidx_") && v == "1"{
			chunkCount ++
		}
	 }

	 if totalCount != chunkCount{
		 w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		 return
	 }

	 //TODO combine chunks
	 fsize, _ := strconv.Atoi(filesize)
	 db.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	 db.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))
	 w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
 }

 func CancelUploadhandler(w http.ResponseWriter, r *http.Request){

 }