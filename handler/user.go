package handler

import (
	"distributed_file/util"
	"github.com/samtake/filestore-server/filestore-server-gin/db"
	"net/http"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt = "*#890"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5{
		w.Write([]byte("Invalid parameter"))
		return
	}
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	suc := db.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("success"))
	}else{
		w.Write([]byte("failed"))
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password+pwdSalt))
	pwdChecked := db.UserSignin(username, encPasswd)
	if !pwdChecked{
		w.Write([]byte("Failed"))
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
		},
	}
	w.Write(resp.JSONBytes())

}
