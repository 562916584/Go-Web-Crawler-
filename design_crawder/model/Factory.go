package model

// 抽象工厂类

type Entity interface {
}
type WorkerCreator interface {
	Create() Entity
}
