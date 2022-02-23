package config

import "fmt"

var (
	MySQLSource = "root:123456@tcp(127.0.0.1:3306)/filestore?charset=utf8"
)

func UploadDBHost(host string) {
	MySQLSource = fmt.Sprintf("root:123456@tcp(%s)/filestore?charset=utf8", host)
}
