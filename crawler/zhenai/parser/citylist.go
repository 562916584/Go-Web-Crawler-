package parser

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler_distributed/config"
	"log"
	"regexp"
)

const cityList = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" [^>]*>([^<]+)</a>`

// 解析城市列表
func ParseCityList(contents []byte, _ string) engine.ParseResult {
	//<a href="http://www.zhenai.com/zhenghun/dadukou" data-v-0c63b635="">大渡口</a>
	// `` 里面不存在转义
	re, err := regexp.Compile(cityList)
	if err != nil {
		log.Printf("ParseCityList failed : %s", err)
	}
	// 返回正则表达式匹配到的所有子串 返回[][][]byte 可以看做是一个[][]string
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	// 爬去前10个城市  不然 太慢了
	//limit := 10
	for _, v := range matches {
		//result.Items = append(result.Items, "City"+string(v[2]))
		// v[0]为匹配原本 后面依次按括号顺序返回匹配结果
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(v[1]),
			Parser: engine.NewFuncParser(ParseCity, config.ParseCity),
		})
		//fmt.Printf("城市: %s  url: %s  \n",v[2],v[1])
		// 单线程版本需要限制城市数目
		//limit--
		//if limit == 0 {
		//	break
		//}
	}
	return result
}
