package engine

type Request struct {
	Url string
	// 解析返回butf-8文本 返回 url列表 和 item信息
	ParseFunc func([]byte) ParseResult
}

type ParseResult struct {
	Requests []Request
	Items    []interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
