package main

import (
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/rpcsupport"
	"WebSpider/crawler_distributed/worker"
	"fmt"
	"log"
)

func main() {
	// 起分布式爬虫 服务器
	log.Fatal(rpcsupport.ServeRpc(fmt.Sprintf(":%d", config.WorkerPort0),
		worker.CrawlService{}))
}
