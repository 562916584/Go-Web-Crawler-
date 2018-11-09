package worker

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/zhenai/parser"
	"WebSpider/crawler_distributed/config"
	"fmt"
	"github.com/pkg/errors"
	"log"
)

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
func DeserializeRequest(r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}
func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		return engine.NewFuncParser(parser.ParseCityList, config.ParseCityList), nil
	case config.ParseCity:
		return engine.NewFuncParser(parser.ParseCity, config.ParseCity), nil
	case config.ParseProfile:
		username, ok := p.Args.(string)
		if ok == true {
			return parser.NewProfileParser(username), nil
		} else {
			return nil, fmt.Errorf("invaild args : %v", p.Args)
		}
	case config.NilParser:
		return engine.NilParser{}, nil
	default:
		return nil, errors.New("No method finding!")
	}
}
func DeserializeResult(r ParseResult) (engine.ParseResult, error) {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		enginerequest, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error desearialize request : %v",
				err)
			continue
		}
		result.Requests = append(result.Requests,
			enginerequest)
	}
	return result, nil
}
