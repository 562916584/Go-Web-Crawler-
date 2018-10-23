package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/zhenai/parser"
)

//爬去网页 并转码为utf-8
func main() {
	engine.Run(engine.Request{
		Url:       "http://www.zhenai.com/zhenghun",
		ParseFunc: parser.ParseCityList,
	})
}
