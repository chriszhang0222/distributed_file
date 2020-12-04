package main

import (
	"distributed_file/account/handler"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main(){

	pwd,_ := os.Getwd()
	fmt.Println(pwd)
	http.Handle("/static/", http.FileServer(http.Dir(filepath.Join(pwd, "./"))))
	http.HandleFunc("/user/signup", handler.SignUpHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	err := http.ListenAndServe(":7998", nil)
	if err != nil{
		fmt.Println(err.Error())
	}

}
