package engine

var visitedUrls = make(map[string]bool)

// 爬虫引擎
type ConcurrentEngine struct {
	// 并发爬虫的调度器接口
	Scheduler Scheduler
	// 并发worker的数目
	WorkerCount int
	// 向elastic 存入数据的通道
	ItemChan chan Item
	// 定义worker函数 两种 一种单机 一种分布式
	RequestProcessor Processor
}

// 定义 worker函数
type Processor func(r Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	// 将request传入调度器中
	Submit(Request)
	// 返回一个输送request的通道
	WorkerChan() chan Request
	// 配置调度器的输送request通道
	ConfigureMusterWorkerChan(chan Request)
	// 调度器开始工作
	Run()
}

type ReadyNotifier interface {
	// worker的request输入通道，通过调度器workerChan，送入worker队列中
	WorkerReady(chan Request)
}

// 并发爬虫引擎开始
func (e *ConcurrentEngine) Run(seeds ...Request) {
	// 准备工作
	//in := make(chan Request)
	out := make(chan ParseResult)
	// 将in 这个channel与workerChannel关联
	//e.Scheduler.ConfigureMusterWorkerChan(in)

	e.Scheduler.Run()
	// 产生worker
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
			// 判断是不是重复的url
			if isDuplicate(request.Url) {
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

// 开启线程 将request通过channel读取 然后调用worker函数 从网上去抓取数据
// 传入两个参数 in--输入request的通道 out--输出Parseresult的通道
func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// 告诉scheduler 我准备好了
			ready.WorkerReady(in)
			// 等待调度器调度，然后in输出request
			request := <-in
			//result, err := Worker(request)
			// 替换为分布式worker --rpc
			// worker 调用fetcher函数 和 Parse函数 返回result
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
