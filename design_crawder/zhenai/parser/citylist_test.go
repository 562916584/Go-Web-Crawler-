package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseCityList(t *testing.T) {
	contents, err := ioutil.ReadFile("citylist_data.html")
	if err != nil {
		panic(err)
	}
	result := ParseCityList(contents, "")
	const resultSize = 470
	expectUrls := []string{
		"http://www.zhenai.com/zhenghun/aba", "http://www.zhenai.com/zhenghun/akesu", "http://www.zhenai.com/zhenghun/alashanmeng",
	}

	if len(result.Requests) != resultSize {
		t.Errorf("result shuold have %d"+
			"requests ,but had %d", resultSize, len(result.Requests))
	}
	for i, url := range expectUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d :%s ; but"+
				"was %s", i, url, result.Requests[i].Url)
		}
	}

	// verify result
	//fmt.Printf("%s \n", contents)
}
