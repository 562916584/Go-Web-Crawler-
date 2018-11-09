package worker

import "WebSpider/crawler/engine"

// 序列化函数用于网络传递
type SerializedParser struct {
	// 函数名字
	Name string
	// 函数参数
	Args interface{}
}

// {'ParseCity',nil}  {'ParseCityList',nil} {'ProfileParser',Username}
type Request struct {
	Url    string
	Parser SerializedParser
}
type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

// 将request序列化
func SerialzeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

// 将ParseResult 序列化
func SerialzeResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests,
			SerialzeRequest(req))
	}
	return result
}

// 将序列化Request 转化为engine的request
func DeserializeRequest(r Request) engine.Request {
	return engine.Request{
		Url:    r.Url,
		Parser: deserializeParser(r.Parser),
	}
}
func deserializeParser(p SerializedParser) engine.Parser {
}
func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests,
			DeserializeRequest(req))
	}
	return result
}

type CrawlService struct {
}

func (CrawlService) Procress(
	req engine.Request, result *engine.ParseResult) error {
	panic("aaa")
}
