package fetche

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// UA代理
var userAgent = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/14.0.835.163 Safari/535.1",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; InfoPath.3)",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; GTB7.0)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E) QQBrowser/6.9.11079.201",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.472.33 Safari/534.3 SE 2.X MetaSr 1.0",
	"Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50",
	"Mozilla/5.0(Windows;U;WindowsNT6.1;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50",
	"Mozilla/4.0(compatible;MSIE7.0;WindowsNT5.1;TencentTraveler4.0)",
}

// 限制访问速度设置的访问延时
var rateLimiter = time.Tick(10 * time.Millisecond)

// 访问url 返回得到的html字节Slice
func Fetche(url string) ([]byte, error) {
	//resp, err := http.Get(url)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	// 通过限制访问速度 防止封锁
	<-rateLimiter
	log.Printf("Fetching url %s", url)
	// 得到一个http请求的客户端
	client := &http.Client{}
	// 返回一个get 请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// 随机产生UA  初始化种子
	rand.Seed(time.Now().Unix())
	req.Header.Set("User-Agent", userAgent[rand.Intn(10)])

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//fmt.Println("Error : status code", resp.StatusCode)
		return nil, fmt.Errorf("wrong status code : %d", resp.StatusCode)
	}
	//将GBK转为UTF-8 第一种特定转
	//utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	//自动发现网页编码 然后转换成UTF-8
	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(resp.Body,
		e.NewDecoder())
	//all, err := ioutil.ReadAll(utf8Reader)
	//if err != nil {
	//	panic(err)
	//}
	//return all, nil
	return ioutil.ReadAll(utf8Reader)
}

// 通过读入 1024字节的内容判断编码方式 并返回
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Feche error :%s", err)
		return unicode.UTF8
	}
	// 返回编码方式
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
