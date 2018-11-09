package worker

import "WebSpider/crawler/engine"

type CrawlService struct {
}

func (c CrawlService) Process(
	req Request, result *ParseResult) error {
	engineReq, err := DeserializeRequest(req)
	if err != nil {
		return err
	}
	engineResult, err1 := engine.Worker(engineReq)
	if err1 != nil {
		return err1
	}
	*result = SerialzeResult(engineResult)
	return nil
}
