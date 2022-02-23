package main

import (
	"deck/common"
	"deck/service/account/handler"
	"deck/service/account/proto"
	db "deck/service/db/client"
	"deck/util"
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	reg := util.ConsulRegistry()
	service := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Flags(common.CustomFlags...),
		micro.Registry(reg),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	service.Init()
	db.Init(service)
	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
	go startPrometheus()
}
func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("localhost:8887", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
