package engine

type ParserFunc func(
	contents []byte, url string) ParseResult

// 网页解析器接口
type Parser interface {
	// 解析函数 参数：contents为html读取字节slice url为当前访问的url  返回值：解析结果结构体
	Parse(contents []byte, url string) ParseResult
	// 序列化函数 返回值： 函数名字  函数参数 --用于分布式rpc传输
	Serialize() (name string, args interface{})
}
type Request struct {
	Url string
	// 解析返回butf-8文本 返回 url列表 和 item信息
	Parser Parser
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}
type Item struct {
	Type    string
	Url     string
	Id      string
	Payload interface{}
}

type NilParser struct {
}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}

func (NilParser) Serialize() (
	name string, args interface{}) {
	return "NilParser", nil
}

type FuncParser struct {
	parser ParserFunc
	Name   string
}

func (f *FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.Name, nil
}

func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		Name:   name,
	}
}
