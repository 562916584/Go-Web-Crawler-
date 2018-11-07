package main

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/model"
	"WebSpider/crawler_distributed/rpcsupport"
	"testing"
)

func TestItemSaver(t *testing.T) {
	// start ItemSaverServe
	go serveRpc(":1234", "test1")
	// start ItemSaverClient
	client, err := rpcsupport.NewClient(":1234")
	if err != nil {
		panic(err)
	}
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3000-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Xinzuo:     "牧羊座",
			Occupation: "人事/行政",
			Marriage:   "离异",
			House:      "已购房",
			Hokou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}
	// Call
	result := ""
	err := client.Call("ItemSaverService.Save", expected, result)
}
