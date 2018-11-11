package client

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/worker"
	"net/rpc"
)

// 返回整个worker的分布式函数
func CreateProcessor(clientChan chan *rpc.Client) (engine.Processor, error) {
	return func(r engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerialzeRequest(r)
		var sResult worker.ParseResult
		WorkerClient := <-clientChan
		err := WorkerClient.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializeResult(sResult)
	}, nil
}
