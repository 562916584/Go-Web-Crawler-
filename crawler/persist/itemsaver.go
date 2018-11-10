package persist

import (
	"WebSpider/crawler/engine"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func ItemSaver(index string) (chan engine.Item, error) {
	client, err := elastic.NewClient(
		// 运行在docker上 是内网访问
		elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	go func() {
		itemCount := 1
		for {
			// 等待数据从out中送出 然后存入Elastic search中
			item := <-out
			log.Printf("Saver item :#%d:  %v\n", itemCount, item)
			itemCount++

			err := Save(item, client, index)
			if err != nil {
				log.Printf("Item Saver : error"+
					"saving item  %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(item engine.Item, client *elastic.Client, index string) (err error) {
	if err != nil {
		return err
	}
	if item.Type == "" {
		return errors.New("must supply type")
	}
	// index--数据库名字
	// Type--表名
	// ID--数据ID
	indexService := client.Index().Index(index).
		Type(item.Type).
		BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err1 := indexService.Do(context.Background())
	if err1 != nil {
		return err1
	}
	//fmt.Println("%+v", resp)
	return nil
}
