package persist

import (
	"WebSpider/design_crawder/engine"
	"WebSpider/design_crawder/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

type ElasticSearch struct {
}

type ElasticSearchCreator struct {
}

func (this *ElasticSearchCreator) Create() model.Entity {
	s := new(ElasticSearch)
	return s
}

// 返回一个 向elasticSearch 存数据的 channel通道
func (this *ElasticSearch) ItemSaver(index string) (chan engine.Item, error) {
	// 创建一个连接elastic的客户端 --- 客户端使用默认的接口 192.168.99.100:9200
	client, err := elastic.NewClient(
		// 运行在docker上 是内网访问
		elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	out := make(chan engine.Item)
	// 常用做法，将out输出的处理方法 go func出去
	// 然后返回out通道给外界,以此运行方法
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
	// 向elastic存Item
	_, err1 := indexService.Do(context.Background())
	if err1 != nil {
		return err1
	}
	//fmt.Println("%+v", resp)
	return nil
}
