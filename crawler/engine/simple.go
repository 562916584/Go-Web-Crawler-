package engine

import (
	"WebSpider/crawler/fetche"
	"log"
)

type SimpleEngine struct{}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]
		// 从队列中访问URL 返回 ParseResult
		parserResult, err := worker(r)
		if err != nil {
			continue
		}
		// 将一系列的 parserResult.Request 全部加入到requests
		requests = append(requests, parserResult.Requests...)
		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}
	}
}

// 并发worker 输入request 返回requests和Items
func worker(r Request) (ParseResult, error) {
	log.Printf("Fetching: %s", r.Url)
	body, err := fetche.Fetche(r.Url)
	if err != nil {
		log.Printf("Fecher: error URL :%s err: %s", r.Url, err)
		return ParseResult{}, err
	}
	return r.ParseFunc(body), nil
}
