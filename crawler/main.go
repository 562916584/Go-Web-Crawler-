package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/persist"
	"WebSpider/crawler/scheduler"
	"WebSpider/crawler/zhenai/parser"
)

//爬去网页 并转码为utf-8
func main() {

	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    persist.ItemSaver(),
	}
	// 并发爬虫入口
	//e.Run(engine.Request{
	//	Url:       "http://www.zhenai.com/zhenghun",
	//	ParseFunc: parser.ParseCityList,
	//})

	//上海城市测试
	e.Run(engine.Request{
		Url:       "http://www.zhenai.com/zhenghun/shanghai",
		ParseFunc: parser.ParseCity,
	})
}
