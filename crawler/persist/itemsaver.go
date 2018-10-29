package persist

import (
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func ItemSaver() chan interface{} {
	out := make(chan interface{})
	go func() {
		itemCount := 1
		for {
			item := <-out
			log.Printf("Saver item :#%d:  %v\n", itemCount, item)
			itemCount++

			save(item)
		}
	}()
	return out
}

func save(item interface{}) (id string, err error) {
	client, err := elastic.NewClient(
		// 运行在docker上 是内网访问
		elastic.SetSniff(false))
	if err != nil {
		return "", err
	}
	resp, err := client.Index().Index("dating_profile").
		Type("zhenai").BodyJson(item).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	//fmt.Println("%+v", resp)
	return resp.Id, nil
}
