package engine

import "WebSpider/design_crawder/model"

// 引擎策略模式 可以方便切换引擎

// 统一接口方法
type Operator interface {
	Run(seeds ...Request)
}

type Operation struct {
	Operator Operator
}

// 统一的接口启动引擎
func (this *Operation) Operate(seed Request) {
	this.Operator.Run(seed)
}

type OperationCreator struct {
}

func (this *OperationCreator) Create() model.Entity {
	s := new(Operation)
	return s
}
