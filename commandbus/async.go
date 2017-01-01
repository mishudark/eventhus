package async

import "cqrs"

var workerPool = make(chan chan cqrs.Command)

//Worker contains the basic info to manage commands
type Worker struct {
	WorkerPool     chan chan cqrs.Command
	JobChannel     chan cqrs.Command
	CommandHandler cqrs.CommandHandlerRegister
}

//Bus stores the command handler
type Bus struct {
	CommandHandler cqrs.CommandHandlerRegister
	maxWorkers     int
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

				if !job.IsValid() {
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
		WorkerPool:     workerPool,
		CommandHandler: commandHandler,
		JobChannel:     make(chan cqrs.Command),
	}

	w.Start()
}

//HandleCommand ad a job to the queue
func (b *Bus) HandleCommand(command cqrs.Command) {
	go func(c cqrs.Command) {
		workerJobQueue := <-workerPool
		workerJobQueue <- c
	}(command)
}

//NewBus return a bus with command handler register
func NewBus(register cqrs.CommandHandlerRegister, maxWorkers int) *Bus {
	b := &Bus{
		CommandHandler: register,
		maxWorkers:     maxWorkers,
	}

	//start the bus
	b.Start()
	return b
}

//Start the bus
func (b *Bus) Start() {
	for i := 0; i < b.maxWorkers; i++ {
		NewWorker(b.CommandHandler)
	}
}
