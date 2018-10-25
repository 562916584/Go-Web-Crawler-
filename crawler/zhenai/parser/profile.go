package parser

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/model"
	"regexp"
	"strconv"
)

var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var nameRe = regexp.MustCompile(`<a class="name fs24">([^<]+)</a>`)
var Gender = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var Height = regexp.MustCompile(`<td><span class="label">身高：</span><span field="">([\d]+)CM</span></td>`)
var Weight = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([\d]+)KG</span></td>`)
var Incom = regexp.MustCompile(` <td><span class="label">月收入：</span>([^<]+)</td>`)
var Education = regexp.MustCompile(`<td><span class="label">学历：</span>([^<]+)</td>`)
var Occupation = regexp.MustCompile(`<td><span class="label">工作地：</span>([^<]+)</td>`)
var Hokou = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var Xinzuo = regexp.MustCompile(`<td><span class="label">星座：</span>([^<]+)</td>`)
var House = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var Car = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var Job = regexp.MustCompile(` <td><span class="label">职业：</span><span field="">([^<]+)</span></td>`)

func ParseProfile(contents []byte, name string) engine.ParseResult {
	profile := model.Profile{}
	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Age = age
	}
	height, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Height = height
	}
	weight, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Weight = weight
	}
	profile.Marriage = extractString(contents, marriageRe)
	profile.Name = name
	profile.Gender = extractString(contents, Gender)
	profile.Income = extractString(contents, Incom)
	profile.Education = extractString(contents, Education)
	profile.Occupation = extractString(contents, Occupation)
	profile.Hokou = extractString(contents, Hokou)
	profile.Xinzuo = extractString(contents, Xinzuo)
	profile.House = extractString(contents, House)
	profile.Car = extractString(contents, Car)
	profile.Job = extractString(contents, Job)

	result := engine.ParseResult{}
	result.Items = append(result.Items, profile)
	return result
}

func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
