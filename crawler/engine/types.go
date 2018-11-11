package engine

// 解析函数ParserFunc
type ParserFunc func(
	contents []byte, url string) ParseResult

// 网页解析器接口
type Parser interface {
	// 解析函数 参数：contents为html读取字节slice url为当前访问的url  返回值：解析结果结构体
	Parse(contents []byte, url string) ParseResult
	// 序列化函数 返回值： 函数名字  函数参数 --用于分布式rpc传输
	Serialize() (name string, args interface{})
}

// 定义的request请求结构体
// 由 url 和 网页解析器（每个解析函数都是西安了接口方法，这样确定该url用什么解析函数）
type Request struct {
	Url string
	// 网页解析器接口
	Parser Parser
}

// 解析返回结构体
// 由[]request 要访问的request列表  []item 抓取的人信息
type ParseResult struct {
	Requests []Request
	Items    []Item
}

// 抓取的人信息结构体
type Item struct {
	// 存入elastic的表名
	Type string
	// 此人的url主页地址
	Url string
	// 此人的网页用户ID 用于elastic ID
	Id string
	// 此人的个人信息
	Payload interface{}
}

// 空方法
type NilParser struct {
}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}

func (NilParser) Serialize() (
	name string, args interface{}) {
	return "NilParser", nil
}

// 实现网页解析器接口的结构体
type FuncParser struct {
	// 解析函数
	parser ParserFunc
	// 函数名字
	Name string
}

func (f *FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.parser(contents, url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.Name, nil
}

// 创建一个FuncParser
func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		Name:   name,
	}
}
