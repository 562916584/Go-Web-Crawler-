package engine

var visitedUrls = make(map[string]bool)

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan    chan Item
}
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
		createWorker(e.Scheduler.WorkerChan(), out, e.Scheduler)
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
func createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			// 告诉scheduler 我准备好了
			ready.WorkerReady(in)
			request := <-in
			result, err := worker(request)
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
