package fronted

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/fronted/model"
	common "WebSpider/crawler/model"
	"html/template"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	template := template.Must(
		template.ParseFiles("template.html"))
	out, err := os.Create("template.test.html")
	if err != nil {
		panic(err)
	}
	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: common.Profile{
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
	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}
	page.Start = 0
	println(len(page.Items))
	err1 := template.Execute(out, page)
	if err1 != nil {
		panic(err1)
	}
}
