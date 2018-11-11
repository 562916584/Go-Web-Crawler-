package client

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"WebSpider/crawler_distributed/rpcsupport"
	"log"
)

// 连接itemSaver服务器 参数是服务地址  (客户端不需要管elastic配置，配置放在服务器)
// 返回向存item服务 输送item的通道out
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
			// 客户端呼叫服务
			err1 := client.Call(config.ItemSaverRpc, item, &result)
			if err1 != nil {
				log.Printf("Item Saver : error"+
					"saving item  %v: %v", item, err1)
			}
		}
	}()
	return out, nil
}
