package main

import (
	"WebSpider/design_crawder/config"
	"WebSpider/design_crawder/engine"
	"WebSpider/design_crawder/model"
	"WebSpider/design_crawder/persist"
	"WebSpider/design_crawder/scheduler"
	"WebSpider/design_crawder/zhenai/parser"
)

func main() {
	//creater := &model.CreateFactory{}
	//creater.Register()
	//tp := creater.Create("Profile")
	//profile := tp.(*model.Profile)
	//profile.Name = "123123"
	//profile.Age = 15
	//fmt.Println(profile)

	// 创建抽象工厂
	var c model.WorkerCreator
	c = new(persist.ElasticSearchCreator)
	persist := c.Create().(*persist.ElasticSearch)
	c = new(engine.ConcurrentEngineCreator)
	engines := c.Create().(*engine.ConcurrentEngine)
	c = new(engine.RequestCreator)
	request := c.Create().(*engine.Request)
	itemChan, err := persist.ItemSaver(config.ElasticIndex)
	if err != nil {
		panic(err)
	}
	engines.Prepare(itemChan, &scheduler.QueuedScheduler{}, config.WorkCount, engine.Worker)
	//并发爬虫入口 --- 配置第一个request的url和解析函数
	request.Prepare(config.URL,
		engine.NewFuncParser(parser.ParseCityList, config.ParseCityList))
	c = new(engine.OperationCreator)
	// 策略模式 (放入引擎然后开始工作)
	engineOperation := c.Create().(*engine.Operation)
	engineOperation.Operator = engines
	engineOperation.Operate(*request)
}
