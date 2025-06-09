package main

import (
	"time"

	"WorkerPoolVK/pkg/pool"
)

func main() {
	p := pool.NewWorkerPool(10)

	p.AddWorker()
	p.AddWorker()

	p.Submit("первая")
	p.Submit("вторая")
	p.Submit("третья")

	time.Sleep(2 * time.Second)

	p.AddWorker()
	time.Sleep(200 * time.Millisecond)
	p.Submit("четвёртая")

	time.Sleep(1 * time.Second)

	p.RemoveWorker(0)
	p.Submit("пятая")

	time.Sleep(2 * time.Second)
	p.Shutdown()
}
