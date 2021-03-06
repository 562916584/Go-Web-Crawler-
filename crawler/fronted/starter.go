package main

import (
	"net/http"
)

// MVC 框架
func main() {
	// golang 相对路径是基于当前执行命令目录来说的  当前执行目录 mygo/
	// 表示以下访问内容都在 此目录下
	// 便于加载css js等配置文件
	// 自动寻找到目录下的index.html打开
	http.Handle("/", http.FileServer(
		http.Dir("src/WebSpider/crawler/fronted/view")))
	http.Handle("/search", controller.CreateSearchResultHandler(
		"src/WebSpider/crawler/fronted/view/template.html"))
	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		panic(err)
	}
}
