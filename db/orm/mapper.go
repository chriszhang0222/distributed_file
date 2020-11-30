package orm

import (
	"errors"
	"reflect"
)

var funcs = map[string]interface{}{
	"/file/OnFileUploadFinished": OnFileUploadFinished,
	"/file/GetFileMeta": GetFileMeta,
	"/file/GetFileMetaList": GetFileMetaList,
	"/file/UpdateFileLocation": UpdateFileLocation,

	"/user/UserSignUp": UserSignUp,
	"/user/UserSignIn": UserSignin,
	"/user/GetUserInfo": GetUserInfo,
	"/user/UserExist": UserExist,

	"/ufile/OnUserFileUploadFinished": OnUserFileUploadFinished,
	"/ufile/QueryUserFileMetas":       QueryUserFileMetas,
	"/ufile/DeleteUserFile":           DeleteUserFile,
	"/ufile/RenameFileName":           RenameFileName,
	"/ufile/QueryUserFileMeta":        QueryUserFileMeta,
}

func FuncCall(name string, params ...interface{})(result []reflect.Value, err error){
	if _, ok := funcs[name]; !ok{
		err = errors.New("Function does not exist")
		return
	}
	f := reflect.ValueOf(funcs[name])
	if len(params) != f.Type().NumIn(){
		err = errors.New("Numebr of Params Dest not match")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params{
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}
