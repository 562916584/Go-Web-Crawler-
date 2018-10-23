package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

//爬去网页 并转码为utf-8
func main() {
	resp, err := http.Get(
		"http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error : status code", resp.StatusCode)
		return
	}
	//将GBK转为UTF-8 第一种特定转
	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	//自动发现网页编码 然后转换成UTF-8
	e := determineEncoding(resp.Body)
	utf8Reader := transform.NewReader(resp.Body,
		e.NewDecoder())
	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s\n", all)
	printCityList(all)
}

func determineEncoding(r io.Reader) encoding.Encoding {
	bytes, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

func printCityList(reader []byte) {
	//<a href="http://www.zhenai.com/zhenghun/dadukou" data-v-0c63b635="">大渡口</a>
	re, err := regexp.Compile(`<a href="http://www.zhenai.com/zhenghun/[0-9a-z]+" [^>]*>[^<]+</a>`)
	if err != nil {
		panic(err)
	}
	matches := re.FindAll(reader, -1)
	for _, v := range matches {
		fmt.Printf("%s", v)
	}
}
