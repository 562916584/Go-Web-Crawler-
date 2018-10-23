package engine

import (
	"WebSpider/crawler/fetche"
	"log"
)

func Run(seeds ...Request) {
	var requests []Request
	for _, r := range seeds {
		requests = append(requests, r)
	}
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		log.Printf("Fetching: %s", r.Url)
		body, err := fetche.Fetche(r.Url)
		if err != nil {
			log.Printf("Fecher: error URL :%s err: %s", r.Url, err)
			continue
		}
		parserResult := r.ParseFunc(body)
		// 将一系列的 parserResult.Request 全部加入到requests
		requests = append(requests, parserResult.Requests...)
		for _, item := range parserResult.Items {
			log.Printf("Got item %v", item)
		}
	}
}
