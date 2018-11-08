package client

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/rpcsupport"
	"log"
)

func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 1
		for {
			item := <-out
			log.Printf("Saver item :#%d:  %v\n", itemCount, item)
			itemCount++

			//err := Save(item, client, index)
			result := ""
			err1 := client.Call(config.ItemSaverRpc, item, &result)
			if err1 != nil {
				log.Printf("Item Saver : error"+
					"saving item  %v: %v", item, err1)
			}
		}
	}()
	return out, nil
}
