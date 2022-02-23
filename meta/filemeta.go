package meta

import (
	"deck/mysqldb"
	"fmt"
)

type FileMeta struct {
	FileShal string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileShal] = fmeta
}
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mysqldb.OnFileUploadFinish(fmeta.FileShal, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}
func GetFileMeta(fileShal string) FileMeta {
	return fileMetas[fileShal]
}
func GetFileMetaDB(filesha1 string) (FileMeta, error) {
	tfile, err := mysqldb.GetFileMeta(filesha1)
	if err != nil {
		fmt.Print(err.Error())
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileShal: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}
