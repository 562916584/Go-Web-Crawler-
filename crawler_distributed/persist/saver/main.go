package main

import (
	"WebSpider/crawler_distributed/persist"
	"WebSpider/crawler_distributed/rpcsupport"
	"gopkg.in/olivere/elastic.v5"
)

func main() {
	err := serveRpc(":1234", "dating_profile")
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

	rpcsupport.ServeRpc(":1234",
		persist.ItemSaverService{
			Client: client,
			Index:  "dating_profile",
		})
	return nil
}
