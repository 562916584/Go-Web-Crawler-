package main

import (
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/persist"
	"WebSpider/crawler_distributed/rpcsupport"
	"flag"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
)

// go run main.go --port=1234
var port = flag.Int("port", 0, "the port foe me to listen on")

// 起服务 存item的服务
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

// 传入服务器地址host和elastic的数据库名字index
func serveRpc(host, index string) error {
	// 初始化elastic的客户端
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	// 绑定服务--itemSaverService  并且在服务器端进行配置
	rpcsupport.ServeRpc(host,
		&persist.ItemSaverService{
			Client: client,
			Index:  index,
		})
	return nil
}
