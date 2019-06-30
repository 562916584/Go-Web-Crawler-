package scheduler

import "WebSpider/design_crawder/engine"

// worker 是抽象出来的工作函数 他的职能是获取
//  并发爬虫调度器
//  形成worker队列 等待request传入 然后调用worker函数抓取数据
type QueuedScheduler struct {
	// 输出request的通道
	requestChan chan engine.Request
	// 输出worker 的 in通道（in 为 worker函数的request的通道）
	workerChan chan chan engine.Request
}

// 将request传入 requestChan
func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

// 将初化始的worker的in通道传入队列workerChan中，等待任务到来
func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

// 将外面输入request的通道与队列通道绑定   用于单个worker中并发使用
func (s *QueuedScheduler) ConfigureMusterWorkerChan(c chan engine.Request) {
	s.requestChan = c
}

// 创建一个 request的通道   用于创建100个worker的request通道
func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request)
}

// 并发式爬虫核心部分之一 调度器调度
// worker队列循环等待送入request 执行任务
func (s *QueuedScheduler) Run() {
	// requestChan 通过这个通道 从外界进入request
	// workerChan 通过这个通道 得到具体那个worker的in通道 （in 是worker函数中 输入request的通道）
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			// activeRequest 当前执行的request
			// activeWorker  当前执行的worker的in通道（in 是worker函数中 输入request的通道）
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			// 从队列中取出一个request和一个worker 的 in channel
			if len(requestQ) > 0 && len(workerQ) > 0 {
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			// selec为了满足同时存在三种下面分支同时都在执行的任务
			select {
			case r := <-s.requestChan:
				// send r to a request
				// request存入队列
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				// send a chan engine.request to w
				// worker的in通道 存入队列
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest:
				// 被选择 再从队列中剔除
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}
		}
	}()
}
