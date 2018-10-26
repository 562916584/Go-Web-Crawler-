package scheduler

import "WebSpider/crawler/engine"

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) Submit(r engine.Request) {
	// 送 request给 worker channel
	// 启动concurrentEngine 开始工作
	go func() {
		s.workerChan <- r
	}()
}

func (s *SimpleScheduler) ConfigureMusterWorkerChan(c chan engine.Request) {
	s.workerChan = c
}

func (s *SimpleScheduler) WorkerReady(chan engine.Request) {
}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}
