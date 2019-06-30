package config

const (
	// Worker数目
	WorkCount = 100
	// 爬去网页地址
	URL = "http://www.zhenai.com/zhenghun"
	// Service Port
	WorkerPort0   = 9000
	ItemSaverPort = 1234

	// ElasitcSearch
	ElasticIndex = "dating_profile"

	// RPC Endpoints
	ItemSaverRpc    = "ItemSaverService.Save"
	CrawlServiceRpc = "CrawlService.Process"

	//Parser names
	ParseCity     = "ParseCity"
	ParseCityList = "ParseCityList"
	ParseProfile  = "ParseProfile"
	NilParser     = "NilParser"
)
