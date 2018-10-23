package parser

import (
	"WebSpider/crawler/engine"
	"log"
	"regexp"
)

const cityList = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" [^>]*>([^<]+)</a>`

func ParseCityList(contents []byte) engine.ParseResult {
	//<a href="http://www.zhenai.com/zhenghun/dadukou" data-v-0c63b635="">大渡口</a>
	// `` 里面不存在转义
	re, err := regexp.Compile(cityList)
	if err != nil {
		log.Printf("ParseCityList failed : %s", err)
	}
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, v := range matches {
		result.Items = append(result.Items, string(v[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url:       string(v[1]),
			ParseFunc: engine.NilParser,
		})
		//fmt.Printf("城市: %s  url: %s  \n",v[2],v[1])
	}
	return result
}
