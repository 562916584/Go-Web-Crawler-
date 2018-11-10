package engine

var visitedUrls = make(map[string]bool)

// 爬虫引擎
type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	// 存入Item数据的管道  在main函数配置
	// 当这个channel 有值输入那么便会进入到 SvaeItem中的gorutine里存入Item
	ItemChan chan Item
	// 定义worker函数
	RequestProcessor Processor
}

// 定义 worker函数
type Processor func(r Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	ConfigureMusterWorkerChan(chan Request)
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	// 准备工作
	//in := make(chan Request)
	out := make(chan ParseResult)
	// 将in 这个channel与workerChannel关联
	//e.Scheduler.ConfigureMusterWorkerChan(in)

	e.Scheduler.Run()

	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}
	// 传入request
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}
	//itemCount := 1
	for {
		result := <-out
		if result.Items != nil {
			for _, item := range result.Items {
				//log.Printf("Got item :#%d:  %v\n", itemCount, item)
				//itemCount++
				go func() {
					e.ItemChan <- item
				}()
			}
		}
		for _, request := range result.Requests {
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

// 开启线程 将request通过channel读取 然后调用worker函数
// 传入两个参数 in--输入request的通道 out--输出Parseresult的通道
// 从网上去抓取数据
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// 告诉scheduler 我准备好了
			ready.WorkerReady(in)
			request := <-in
			//result, err := Worker(request)
			// 替换为分布式worker --rpc
			result, err := e.RequestProcessor(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

// 利用map 映射  去除重复的url
func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}
	visitedUrls[url] = true
	return false
}
