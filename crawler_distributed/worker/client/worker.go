package client

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/rpcsupport"
	"WebSpider/crawler_distributed/worker"
	"fmt"
)

// 返回整个worker的分布式函数
func CreateProcessor() (engine.Processor, error) {
	client, err := rpcsupport.NewClient(
		fmt.Sprintf(":%d", config.WorkerPort0))
	if err != nil {
		return nil, err
	}
	return func(r engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerialzeRequest(r)
		var sResult worker.ParseResult
		err := client.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult)
	}, nil
}
