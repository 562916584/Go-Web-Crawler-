package persist

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

// rpc 服务实体--存item服务实体 配置完然后绑定在服务器地址上
type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

// 服务方法 满足jsonRpc标准规范 一个参数 一个返回值 一个错误返回
func (s *ItemSaverService) Save(item engine.Item, result *string) error {
	// 调用存item函数
	err := persist.Save(item, s.Client, s.Index)
	log.Printf("Saved Item : %+v ", item)
	if err == nil {
		*result = "ok"
	}
	return err
}
