package pool

import (
	"fmt"
	"sync"
	"time"
)

const DefaultTaskBufferSize = 10

type Task string

type Worker struct {
	id     int
	wg     *sync.WaitGroup
	stopCh chan struct{}
	taskCh chan Task
}

type WorkerPool struct {
	workers       map[int]*Worker
	scheduleIndex int
	nextID        int
	wg            sync.WaitGroup
	mu            sync.Mutex
	taskBufSize   int
}

func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Воркер %d: ошибка %v, завершение\n", w.id, r)
			}
		}()

		for {
			select {
			case <-w.stopCh:
				fmt.Printf("Воркер %d: получен сигнал остановки\n", w.id)
				return
			case task, ok := <-w.taskCh:
				if !ok {
					fmt.Printf("Воркер %d: канал задач закрыт, завершение \n", w.id)
					return
				}
				fmt.Printf("Воркер %d: обрабатывает задачу '%s'\n", w.id, task)
				time.Sleep(time.Second)
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.stopCh)
}

func NewWorkerPool(bufferSize int) *WorkerPool {
	if bufferSize <= 0 {
		bufferSize = DefaultTaskBufferSize
	}
	return &WorkerPool{
		workers:     make(map[int]*Worker),
		taskBufSize: bufferSize,
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	w := &Worker{
		id:     wp.nextID,
		wg:     &wp.wg,
		stopCh: make(chan struct{}),
		taskCh: make(chan Task, wp.taskBufSize),
	}

	wp.workers[w.id] = w
	w.Start()

	fmt.Printf("Добавлен воркер %d\n", w.id)
	wp.nextID++
}

func (wp *WorkerPool) RemoveWorker(id int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if worker, ok := wp.workers[id]; ok {
		close(worker.taskCh)
		worker.Stop()
		delete(wp.workers, id)
		fmt.Printf("Удален воркер %d (Остановлен по сигналу)\n", id)
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if len(wp.workers) == 0 {
		fmt.Println("Нет активных воркеров, задача пропущена:", task)
		return
	}

	ids := make([]int, 0, len(wp.workers))
	for id := range wp.workers {
		ids = append(ids, id)
	}

	workerID := ids[wp.scheduleIndex%len(ids)]
	wp.scheduleIndex++

	worker := wp.workers[workerID]
	worker.taskCh <- task

}

func (wp *WorkerPool) Shutdown() {
	wp.mu.Lock()
	for id, worker := range wp.workers {
		close(worker.taskCh)
		worker.Stop()
		delete(wp.workers, id)
	}
	wp.mu.Unlock()

	wp.wg.Wait()
	fmt.Println("Пул завершен. Все воркеры остановлены.")
}
