package main

import (
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/persist"
	"WebSpider/crawler_distributed/rpcsupport"
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
)

// 起服务 存item的服务
// go run main.go --port=1234
var port = flag.Int("port", 0, "the port foe me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
	}
	err := serveRpc(fmt.Sprintf(":%d", *port),
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
