package model

import (
	"encoding/json"
	"reflect"
	"sync"
)

// 全局类名反射工厂对象且不导出包外
var relectFactory *CreateFactory

// 同步控制锁
var singletonMutex sync.Mutex

func init() {
	relectFactory = &CreateFactory{}
}

// 获取单例模式下的工厂
func NewSingleton() *CreateFactory {
	if relectFactory == nil {
		singletonMutex.Lock()
		if relectFactory == nil {
			relectFactory = &CreateFactory{}
		}
		singletonMutex.Unlock()
	}
	return relectFactory
}

//  类名反射
type CreateFactory struct {
	ReflectAdd map[string]reflect.Type
}

// 注册类名反射
func (c *CreateFactory) Register() {
	c.ReflectAdd = make(map[string]reflect.Type)
	c.ReflectAdd["Profile"] = reflect.TypeOf(Profile{})
}

// 依据名来创建对象
func (c *CreateFactory) Create(className string) interface{} {
	return reflect.New(c.ReflectAdd[className]).Interface()
}

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
