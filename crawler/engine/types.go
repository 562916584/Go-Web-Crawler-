package engine

type ParserFunc func(
	contents []byte, url string) ParseResult

type Request struct {
	Url string
	// 解析返回butf-8文本 返回 url列表 和 item信息
	ParseFunc ParserFunc
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

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
