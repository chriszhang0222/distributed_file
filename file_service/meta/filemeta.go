package meta

import "distributed_file/db/orm"
type FileMeta struct{
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(meta FileMeta){
	if meta.FileSha1 != "" {
		fileMetas[meta.FileSha1] = meta
	}
}
func UpdateFileMetaDB(meta FileMeta) orm.ExecResult {
	return orm.OnFileUploadFinished(meta.FileSha1, meta.FileName, meta.FileSize, meta.UploadAt)
}
func GetFileMeta(key string) FileMeta {
	return fileMetas[key]
}
func GetFileMetaFromDB(key string)(*FileMeta, error){
	data, err := orm.GetFileMeta(key)
	if data.Suc{
		tableData := data.Data.(orm.TableFile)
		return &FileMeta{
			FileSha1: key,
			FileName: tableData.FileName.String,
			FileSize: tableData.FileSize.Int64,
			Location: tableData.FileAddr.String,
		}, nil
	}
	return nil, err
}
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}