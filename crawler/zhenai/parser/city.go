package parser

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"regexp"
)

var profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)" [^>]*>([^<]+)</a>`)
var cityUrlRe = regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/shanghai/[^"]+)">下一页</a>`)

// <a href="http://album.zhenai.com/u/1565829136" target="_blank">朗歌</a>
func ParseCity(contents []byte, _ string) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, v := range matches {
		//name := string(v[2])
		//result.Items = append(result.Items, "User"+name)
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(v[1]),
			Parser: NewProfileParser(string(v[2])),
		})
		//fmt.Printf("城市: %s  url: %s  \n",v[2],v[1])
	}
	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(m[1]),
			Parser: engine.NewFuncParser(ParseCity, config.ParseCity),
		})
	}
	return result
}
