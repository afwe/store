package handler

import (
	"deck/meta"
	"deck/store/ceph"
	"deck/util"
	"encoding/json"
	"fmt"
	"gopkg.in/amz.v1/s3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("File get file%s\n", err.Error())
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("fail to create file%s", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("fail to save data into file%s\n", err.Error())
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileShal = util.FileSha1(newFile)
		newFile.Seek(0, 0)
		data, _ := ioutil.ReadAll(newFile)
		bucket := ceph.GetCephBucket("userfile")
		cephPath := "/ceph/" + fileMeta.FileShal
		_ = bucket.Put(cephPath, data, "octet-stream", s3.PublicRead)
		fileMeta.Location = cephPath
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		fmt.Println(fileMeta.FileShal)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finish")
}
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	//filename := r.Form.Get("filename")
	//filesize := r.Form.Get("filesize")
	filemeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if filemeta == (meta.FileMeta{}) {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "秒传成功",
	}
	w.Write(resp.JSONBytes())
}
