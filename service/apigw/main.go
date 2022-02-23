package main

import (
	"deck/service/apigw/route"
)

//创建监控项,并且用标签的形式区分

func main() {
	r := route.Router()
	r.Run(":8180")
}
