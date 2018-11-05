package main

import (
	"WebSpider/crawler/fronted/controller"
	"net/http"
)

func main() {
	// golang 相对路径是基于当前执行命令目录来说的  当前执行目录 mygo/
	// 表示以下访问内容都在 此目录下
	// 便于加载css js等配置文件
	http.Handle("/", http.FileServer(
		http.Dir("src/WebSpider/crawler/fronted/view")))
	http.Handle("/search", controller.CreateSearchResultHandler(
		"src/WebSpider/crawler/fronted/view/template.html"))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
