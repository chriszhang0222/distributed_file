package orm

import "database/sql"

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

type TableUser struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

type ExecResult struct {
	Suc bool `json:"suc"`
	Code int  `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`

}
