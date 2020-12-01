package meta

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

func GetFileMeta(key string)FileMeta{
	return fileMetas[key]
}

func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}