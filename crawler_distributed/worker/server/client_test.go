package main

import (
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/rpcsupport"
	"WebSpider/crawler_distributed/worker"
	"fmt"
	"testing"
	"time"
)

func TestCrawlService(t *testing.T) {
	const host = ":1234"
	go rpcsupport.ServeRpc(host, worker.CrawlService{})
	time.Sleep(3 * time.Second)

	client, err := rpcsupport.NewClient(host)
	if err != nil {
		panic(err)
	}
	req := worker.Request{
		Url: "http://album.zhenai.com/u/108906739",
		Parser: worker.SerializedParser{
			Name: config.ParseProfile,
			Args: "安静的雪",
		},
	}
	var result worker.ParseResult
	err1 := client.Call(config.CrawlServiceRpc, req, &result)
	if err1 != nil {
		t.Errorf("%s", err1)
	} else {
		fmt.Printf("%+v", result)
	}

}
