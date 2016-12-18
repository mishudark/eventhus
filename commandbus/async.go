package async

import "cqrs"

var workerPool = make(chan chan interface{})

//Worker contains the basic info to manage commands
type Worker struct {
	WorkerPool     chan chan interface{}
	JobChannel     chan interface{}
	CommandHandler cqrs.CommandHandlerRegister
}

//Bus stores the command handler
type Bus struct {
	CommandHandler cqrs.CommandHandle
}

//Start initialize a worker ready to receive jobs
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				handler, err := w.CommandHandler.Get(job)
				if err != nil {
					continue
				}

				if err = handler.Handle(job); err != nil {
					//TODO: log the error
				}
			}
		}
	}()
}

//NewWorker initialize the values of worker and start it
func NewWorker(commandHandler cqrs.CommandHandlerRegister) {
	w := Worker{
		WorkerPool: workerPool,
	}

	w.Start()
}

//Add a job to the queue
func (b *Bus) Add(command interface{}) {
	go func(c interface{}) {
		workerJobQueue := <-workerPool
		workerJobQueue <- c
	}(command)
}
