package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/scheduler"
	"WebSpider/crawler/zhenai/parser"
	"WebSpider/crawler_distributed/config"
	Itemsaver "WebSpider/crawler_distributed/persist/client"
	"WebSpider/crawler_distributed/rpcsupport"
	worker "WebSpider/crawler_distributed/worker/client"
	"flag"
	"log"
	"net/rpc"
	"strings"
)

var (
	itemSaverHost = flag.String("itemSaver_host", "",
		"itemSaver host")
	workerHosts = flag.String("worker_hosts", "",
		"worker hosts(comma separated)")
)

//爬去网页 并转码为utf-8
func main() {
	flag.Parse()
	// 开启客户端 连接存Item服务
	itemChan, err := Itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	pool := createClientPool(strings.Split(*workerHosts, ","))
	// 配置分布式worker 函数
	processor, err1 := worker.CreateProcessor(pool)
	if err1 != nil {
		panic(err1)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}
	// 分布式爬虫入口
	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList),
	})

	//上海城市测试
	//e.Run(engine.Request{
	//	Url:       "http://www.zhenai.com/zhenghun/shanghai",
	//	ParseFunc: parser.ParseCity,
	//})
}

func createClientPool(host []string) chan *rpc.Client {
	// 建立连接池 连接所多个workerRpc
	var clients []*rpc.Client
	for _, h := range host {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Successful connecting to %s", h)
		} else {
			log.Printf("error connecting :%s", err)
		}
	}
	// 创建通道 分发rpc.client 一直循环顺序发送
	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
