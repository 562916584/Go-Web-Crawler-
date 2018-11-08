package persist

import (
	"WebSpider/crawler/engine"
	"WebSpider/crawler/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaverService) Save(item engine.Item, result *string) error {
	err := persist.Save(item, s.Client, s.Index)
	log.Printf("Saved Item : %+v ", item)
	if err == nil {
		*result = "ok"
	}
	return err
}
