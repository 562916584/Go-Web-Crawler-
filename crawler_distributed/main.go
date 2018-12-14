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

// 大体架构： 将itemSave和work两个动作放在远端服务器，然后通过jsonRPC来调用
// 连接池作用：因为有两个work的服务器程序，为了均衡负载，起先就建立起两个服务器的连接并返回客户端
// 			   然后通过通道不停的循环送出客户端，客户端程序需要调用work函数时，得到一个客户端并开始呼叫
// 客户端/服务器并发问题： 客户端并发work函数，work函数得到服务端的连接client，然后开始呼叫服务端上的work函数
// 			   服务器端：当把方法注册了以后，每当客户端访问，连接成功，并发jsonRPC的serveConn，做到能并发服务多个客户端
//				的呼叫
// 服务器运行原理： 当你把服务端的方法注册以后，客户端发送序列化函数请求，服务端根据请求和服务上注册的方法，自动调用方法
// 				【函数序列化，反序列化，序列化返回值并返回等写在服务器注册方法里】，完成一次调用。
