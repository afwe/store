package main

import (
	"deck/common"
	"deck/config"
	"deck/mq"
	dbproxy "deck/service/db/client"
	"deck/service/transfer/process"
	"deck/util"
	"fmt"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"), // 服务名称
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

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startTranserService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")

	// 初始化mq client
	mq.Init()

	mq.StartConsumer(
		config.TransOSSQueueName,
		"transfer_oss",
		process.Transfer)
	go startPrometheus()
}

func main() {
	// 文件转移服务
	go startTranserService()

	// rpc 服务
	startRPCService()
}
func startPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	// 启动web服务，监听8085端口
	go func() {
		err := http.ListenAndServe("localhost:8883", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
}
