package main

import (
	"deck/common"
	dbproxy "deck/service/db/client"
	cfg "deck/service/download/config"
	dlProto "deck/service/download/proto"
	"deck/service/download/route"
	dlRpc "deck/service/download/rpc"
	"deck/util"
	"fmt"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log"
	"net/http"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.Flags(common.CustomFlags...),
		micro.Registry(util.ConsulRegistry()),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
	go startPrometheus()
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("localhost:8082", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
