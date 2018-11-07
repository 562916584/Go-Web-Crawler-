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
type SearchResultHandler struct {
	view   view.SearchResultView
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

// 获取搜索结果
func (h SearchResultHandler) getSearchResult(q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q
	resp, err := h.client.Search("dating_profile").
		Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).
		From(from).
		Do(context.Background())
	if err != nil {
		return result, err
	}
	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(reflect.TypeOf(engine.Item{}))
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

func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	return re.ReplaceAllString(q, "Payload.$1:")
}