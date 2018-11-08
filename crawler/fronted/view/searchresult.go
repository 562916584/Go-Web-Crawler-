package view

import (
	"WebSpider/crawler/fronted/model"
	"html/template"
	"io"
)

type SearchResultView struct {
	// 返回生成的模板html
	template *template.Template
}

// 返回生成的模板html
func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(
			template.ParseFiles(filename)),
	}
}

// 向模板html 写入数据
func (s SearchResultView) Render(w io.Writer, data model.SearchResult) error {
	return s.template.Execute(w, data)
}
