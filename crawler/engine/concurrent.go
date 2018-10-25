package engine

import "log"

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
}
type Scheduler interface {
	Submit(Request)
	ConfigureMusterWorkerChan(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {
	// 准备工作
	in := make(chan Request)
	out := make(chan ParseResult)
	// 将in 这个channel与workerChannel关联
	e.Scheduler.ConfigureMusterWorkerChan(in)

	for i := 0; i < e.WorkerCount; i++ {
		createWorker(in, out)
	}
	// 传入request
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		for _, item := range result.Items {
			log.Printf("Got item :  %v\n", item)
		}
		for _, request := range result.Requests {
			e.Scheduler.Submit(request)
		}
	}
}

// 开启线程 将request通过channel读取 然后调用worker函数解析
func createWorker(in chan Request, out chan ParseResult) {
	go func() {
		for {
			request := <-in
			result, err := worker(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}
