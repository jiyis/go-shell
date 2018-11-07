package common

import (
	"fmt"
)

// 实现了handle方法的对象就可以入列
type Job interface {
	Handle(i interface{}) error
}

var JobQueue chan Job

// 声明worker
type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool //停止队列
}

// 负责操作入列分发
type Dispatcher struct {
	workerPool chan chan Job // 队列池
	maxWorkers int           // 队列池的个数
	jobQueue   chan Job
}

/**
创建一个task
*/
func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

/**
启动一个消费task
*/
func (w Worker) start() {
	go func() {
		for {
			// 消费完成了，就重新放入work pool中,等到下一个获取
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				// 取出job 开始消费
				job.Handle(w.id)

				fmt.Printf("worker%d: completed!\n", w.id)
			case <-w.quitChan:
				// 停止消费
				Log.Info("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

/**
停止task
*/
func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

/**
创建一个分发器
*/
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

/**
启动分发器
*/
func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

/**
分发任务
*/
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func(job Job) {
				Log.Info("read to dispatch job to workJobQueue")
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}(job)
		}
	}
}
