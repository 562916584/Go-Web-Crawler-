package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/scheduler"
	"WebSpider/crawler/zhenai/parser"
	"WebSpider/crawler_distributed/config"
	Itemsaver "WebSpider/crawler_distributed/persist/client"
	worker "WebSpider/crawler_distributed/worker/client"
	"fmt"
)

//爬去网页 并转码为utf-8
func main() {
	// 开启客户端 连接上服务
	itemChan, err := Itemsaver.ItemSaver(fmt.Sprintf(":%d", config.ItemSaverPort))
	if err != nil {
		panic(err)
	}
	// 配置分布式worker 函数
	processor, err1 := worker.CreateProcessor()
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
