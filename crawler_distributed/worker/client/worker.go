package client

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/worker"
	"net/rpc"
)

// 返回 worker的分布式函数
func CreateProcessor(clientChan chan *rpc.Client) (engine.Processor, error) {
	// 通过返回匿名函数，将clientChan（输出*rpc.client的通道） 配置到函数中去
	return func(r engine.Request) (engine.ParseResult, error) {
		// worker函数 序列化request 客户端呼叫worker服务 传入序列化request和result
		// 返回的result 再反序列化回来 返回出去
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
