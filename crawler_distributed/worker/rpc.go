package worker

import "WebSpider/crawler/engine"

type CrawlService struct {
	// 存放配置信息
}

func (c CrawlService) Process(
	req Request, result *ParseResult) error {
	// 将req 反序列化为engine.request
	engineReq, err := DeserializeRequest(req)
	if err != nil {
		return err
	}
	// 正常调用worker函数
	engineResult, err1 := engine.Worker(engineReq)
	if err1 != nil {
		return err1
	}
	// 序列化engine.ParseResult然后返回
	*result = SerialzeResult(engineResult)
	return nil
}
