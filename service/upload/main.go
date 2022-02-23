package main

import (
	"deck/common"
	"deck/config"
	"deck/mq"
	dbproxy "deck/service/db/client"
	cfg "deck/service/upload/config"
	upProto "deck/service/upload/proto"
	"deck/service/upload/route"
	upRpc "deck/service/upload/rpc"
	"deck/util"
	"fmt"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.Flags(common.CustomFlags...),
		micro.Registry(util.ConsulRegistry()),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)
	// 初始化mq client
	mq.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
	go startPrometheus()
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.UploadServiceHost)
	// service := web.NewService(
	// 	web.Name("go.micro.web.upload"),
	// 	web.Handler(router),
	// 	web.RegisterTTL(10*time.Second),
	// 	web.RegisterInterval(5*time.Second),
	// )
	// if err := service.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := service.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	os.MkdirAll(config.TempLocalRootDir, 0777)
	os.MkdirAll(config.TempPartRootDir, 0777)

	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("127.0.0.1:8884", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
