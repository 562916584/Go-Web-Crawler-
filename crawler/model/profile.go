package model

import "encoding/json"

type Profile struct {
	Name       string
	Gender     string // 性别
	Age        int
	Height     int
	Weight     int
	Income     string // 收入
	Marriage   string // 婚姻状况
	Education  string
	Occupation string
	Hokou      string // 户口
	Xinzuo     string // 星座
	House      string
	Car        string
	Job        string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}
	err = json.Unmarshal(s, &profile)
	return profile, err
}
