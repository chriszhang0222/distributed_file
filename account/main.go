package main

import (
	h1 "distributed_file/account/handler"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	h2 "distributed_file/file_service/handler"
)

func main(){

	pwd,_ := os.Getwd()
	fmt.Println(pwd)
	http.Handle("/static/", http.FileServer(http.Dir(filepath.Join(pwd, "./"))))
	http.HandleFunc("/user/signup", h1.SignUpHandler)
	http.HandleFunc("/user/signin", h1.SignInHandler)
	http.HandleFunc("/file/upload", h2.UploadHandler)
	http.HandleFunc("/file/fastupload", h2.TryFastUploadhandler)
	err := http.ListenAndServe(":7998", nil)
	if err != nil{
		fmt.Println(err.Error())
	}

}
