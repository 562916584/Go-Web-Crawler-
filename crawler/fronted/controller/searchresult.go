package controller

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/fronted/model"
	"WebSpider/crawler/fronted/view"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// 连接搜索结果网页 和 elasticsearch服务
// Handler实现了Handler接口 serverHTTP函数
type SearchResultHandler struct {
	// 返回搜索模板heml
	view view.SearchResultView
	// elastic 客户端
	client *elastic.Client
}

// 初始化handler 有view和client
func CreateSearchResultHandler(template string) SearchResultHandler {
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

// localhost:8888/search? q=男 已购房&from=20
// 需要查询和翻页
func (h SearchResultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.FormValue("q"))
	from, err := strconv.Atoi(req.FormValue("from"))
	if err != nil {
		from = 0
	}
	var page model.SearchResult
	// 访问elastic 获得结果
	page, err1 := h.getSearchResult(q, from)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
	}
	// 将结果写入模板html 并返回给客户端
	err2 := h.view.Render(w, page)
	if err2 != nil {
		http.Error(w, err2.Error(), http.StatusBadRequest)
	}
}

// 输入q=搜搜条件  from 分页
// 获取搜索结果
func (h SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q
	// 返回结果
	resp, err := h.client.Search("dating_profile").
		Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).
		From(from).
		Do(context.Background())
	if err != nil {
		return result, err
	}
	result.Hits = resp.TotalHits()
	result.Start = from
	// 将搜索结果换成engine.Item{}的格式
	// result.Items = resp.Each(reflect.TypeOf(engine.Item{}))
	// 类型断言
	for _, v := range resp.Each(reflect.TypeOf(engine.Item{})) {
		item, ok := v.(engine.Item)
		if ok {
			result.Items = append(result.Items, item)
		}
	}
	result.PrevFrom = result.Start - len(result.Items)
	result.NextFrom = result.Start + len(result.Items)
	return result, err
}

// 搜索q=男 Age(<20)
// 变成 a=男 Payload.Age(<20)
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}
