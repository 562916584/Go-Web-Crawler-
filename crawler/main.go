package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/scheduler"
	"WebSpider/crawler/zhenai/parser"
)

//爬去网页 并转码为utf-8
func main() {

	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 100,
	}
	e.Run(engine.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
}
