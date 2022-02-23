package handler

import (
	rPool "deck/cache/redis"
	"deck/mysqldb"
	"deck/util"
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type MultiPartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, err.Error(), nil).JSONBytes())
		fmt.Print(err)
		return
	}
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	mpinfo := MultiPartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   fmt.Sprintf("**%d", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024,
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}
	rConn.Do("HSET", "MP_"+mpinfo.UploadID, "chunkcount", mpinfo.ChunkCount)
	rConn.Do("HSET", "MP_"+mpinfo.UploadID, "filehash", mpinfo.FileHash)
	rConn.Do("HSET", "MP_"+mpinfo.UploadID, "filesize", mpinfo.FileSize)
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uploadID := r.Form.Get("filehash")
	chunkIndex := r.Form.Get("index")
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	filepath := "./tmp/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(filepath), 0744)
	fd, err := os.Create(filepath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part falied", nil).JSONBytes())
	}
	defer fd.Close()
	rConn.Do("HSET", "MP_"+uploadID, "index_"+chunkIndex, 1)
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	upid := r.Form.Get("filename")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.ParseInt(r.Form.Get("filesize"), 10, 64)
	if err != nil {
		w.Write(util.NewRespMsg(-2, "cantresolvestrtoint64", nil).JSONBytes())
		return
	}
	filename := r.Form.Get("filename")
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "fail to req mpupload status", nil).JSONBytes())
		fmt.Print(err.Error())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {

		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "shangchuanchucuoo", nil).JSONBytes())
		return
	}
	mysqldb.OnFileUploadFinish(filehash, filename, filesize, "")
}
