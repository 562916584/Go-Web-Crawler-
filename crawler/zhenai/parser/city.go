package parser

import (
	"WebSpider/crawler/engine"
	"log"
	"regexp"
)

const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]+)" [^>]*>([^<]+)</a>`

// <a href="http://album.zhenai.com/u/1565829136" target="_blank">朗歌</a>
func ParseCity(contents []byte) engine.ParseResult {
	re, err := regexp.Compile(cityRe)
	if err != nil {
		log.Printf("ParseCityList failed : %s", err)
	}
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, v := range matches {
		name := string(v[2])
		result.Items = append(result.Items, "User"+name)
		result.Requests = append(result.Requests, engine.Request{
			Url: string(v[1]),
			ParseFunc: func(bytes []byte) engine.ParseResult {
				return ParseProfile(bytes, name)
			},
		})
		//fmt.Printf("城市: %s  url: %s  \n",v[2],v[1])
	}
	return result
}
