package persist

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/model"
	"encoding/json"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"testing"
)

func TestSave(t *testing.T) {
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
	// TODO: Try to start up elasticsearch
	// here using docker go client
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	// Save expected item
	const index = "dating_test"
	err1 := Save(expected, client, index)
	if err1 != nil {
		panic(err1)
	}
	// Fetch saved item
	resp, err := client.Get().Index(index).
		Type(expected.Type).
		Id(expected.Id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//t.Logf("%s", resp.Source)
	var actual engine.Item
	err = json.Unmarshal([]byte(*resp.Source), &actual)
	if err != nil {
		panic(err)
	}

	actualprofile, err := model.FromJsonObj(actual.Payload)
	actual.Payload = actualprofile
	// Verify result
	if expected != actual {
		t.Errorf("got %v; expected %v", actual, expected)
	}
}
