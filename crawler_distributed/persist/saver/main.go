package main

import (
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/persist"
	"WebSpider/crawler_distributed/rpcsupport"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
)

// 起服务 存服务
func main() {
	err := serveRpc(fmt.Sprintf(":%d", config.ItemSaverPort),
		config.ElasticIndex)
	if err != nil {
		panic(err)
	}
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client: client,
			Index:  index,
		})
	return nil
}
