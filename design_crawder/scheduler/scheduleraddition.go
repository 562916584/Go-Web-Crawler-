package scheduler

import (
	"WebSpider/design_crawder/engine"
	"WebSpider/design_crawder/model"
)

// 引擎策略模式 可以方便切换引擎

// 统一接口方法
type Operator interface {
	Submit(r engine.Request)
	WorkerReady(w chan engine.Request)
	ConfigureMusterWorkerChan(c chan engine.Request)
	WorkerChan() chan engine.Request
	Run()
}

type Operation struct {
	Operator Operator
}

// 统一的接口启动引擎
func (this *Operation) OperateSubmit(r engine.Request) {
	this.Operator.Submit(r)
}
func (this *Operation) OperateWorkerReady(w chan engine.Request) {
	this.Operator.WorkerReady(w)
}
func (this *Operation) OperateConfigureMusterWorkerChan(c chan engine.Request) {
	this.Operator.ConfigureMusterWorkerChan(c)
}
func (this *Operation) OperateWorkerChan() chan engine.Request {
	return this.Operator.WorkerChan()
}
func (this *Operation) OperateRun() {
	this.Operator.Run()
}

type OperationCreator struct {
}

func (this *OperationCreator) Create() model.Entity {
	s := new(Operation)
	return s
}
