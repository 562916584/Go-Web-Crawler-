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
	"fmt"
	"log"
	"net/rpc"
	"strings"
)

// 命令行参数
var (
	itemSaverHost = flag.String("itemSaver_host", "",
		"itemSaver host")
	workerHosts = flag.String("worker_hosts", "",
		"worker hosts(comma separated)")
)

//分布式爬虫 总客户端
func main() {
	flag.Parse()
	if *itemSaverHost == "" || *workerHosts == "" {
		fmt.Print("itemSaverHost,workerHosts are nil")
	}
	// 开启存item客户端 连接存Item服务
	itemChan, err := Itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}
	// 建立连接池 返回多个连接到不同worker服务器的rpcClient
	pool := createClientPool(strings.Split(*workerHosts, ","))
	// 配置分布式worker 函数
	processor, err1 := worker.CreateProcessor(pool)
	if err1 != nil {
		panic(err1)
	}
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		// 存服务
		ItemChan: itemChan,
		// worker服务
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

// 建立多个Client来连接worker服务器 参数为服务器地址Slice
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
	// 创建通道 分发rpc.client
	out := make(chan *rpc.Client)
	go func() {
		// 一直循环顺序发送*rpc.client
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
