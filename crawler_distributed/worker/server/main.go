package main

import (
	"WebSpider/crawler_distributed/rpcsupport"
	"WebSpider/crawler_distributed/worker"
	"flag"
	"fmt"
	"log"
)

// 命令行参数
// 用法 go run main.go --port=9000 给port传值为9000
var port = flag.Int("port", 0, "the port foe me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	// 爬去网页的服务
	log.Fatal(rpcsupport.ServeRpc(fmt.Sprintf(":%d", *port),
		worker.CrawlService{}))
}
