package engine

import (
	"WebSpider/crawler/fetche"
	"log"
)

// 并发worker 输入request 返回ParseResult
func Worker(r Request) (ParseResult, error) {
	//log.Printf("Fetching: %s", r.Url)
	// 访问url 转换成body字节slice
	body, err := fetche.Fetche(r.Url)
	if err != nil {
		log.Printf("Fecher: error URL :%s err: %s", r.Url, err)
		return ParseResult{}, err
	}
	return r.Parser.Parse(body, r.Url), nil
}
