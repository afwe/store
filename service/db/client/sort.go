package client

import (
	"deck/meta"
	"time"
)

const baseFormat = "2006-01-02 15:04:05"

type ByUploadTime []meta.FileMeta

func (a ByUploadTime) Len() int {
	return len(a)
}
func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt)
	return iTime.UnixNano() > jTime.Unix()
}
