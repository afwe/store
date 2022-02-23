package mysqldb

import (
	"database/sql"
	"deck/service/db/conn"
	"fmt"
)

func OnFileUploadFinish(filehash, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := conn.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`,`file_addr`,`status`)" +
		" values (?,?,?,?,1)")
	if err != nil {
		fmt.Print("db/file.go" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Print("db/file.go" + err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("file with hash %s exited", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := conn.DBConn().Prepare("select file_sha1,file_name,file_size,file_addr from tbl_file where file_sha1=? and status=1 limit 1")
	if err != nil {
		fmt.Print("cannot get file in mysqldb/file" + err.Error())
		return nil, err
	}
	defer stmt.Close()
	tflie := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tflie.FileHash, &tflie.FileName, &tflie.FileSize, &tflie.FileAddr)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return &tflie, nil
}
