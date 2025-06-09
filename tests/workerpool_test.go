package tests

import (
	"testing"
	"time"

	"WorkerPoolVK/pkg/pool"
)

func TestAddAndRemoveWorker(t *testing.T) {
	p1 := pool.NewWorkerPool(10)

	p1.AddWorker()
	p1.AddWorker()

	p1.RemoveWorker(0)

	time.Sleep(100 * time.Millisecond)

	p1.Shutdown()
}

func TestTaskProcessing(t *testing.T) {
	p2 := pool.NewWorkerPool(10)

	p2.AddWorker()
	p2.AddWorker()

	p2.Submit("task1")
	p2.Submit("task2")
	p2.Submit("task3")

	time.Sleep(1500 * time.Millisecond)
	p2.Shutdown()
}

func TestSubmitWithoutWorkers(t *testing.T) {
	p3 := pool.NewWorkerPool(10)

	p3.Submit("task1")

	p3.Shutdown()
}

func TestDoubleShutdown(t *testing.T) {
	p4 := pool.NewWorkerPool(10)

	p4.AddWorker()
	p4.Submit("task1")
	time.Sleep(500 * time.Millisecond)

	p4.Shutdown()

	p4.Shutdown()
}
