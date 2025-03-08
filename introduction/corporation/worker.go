package corporation

import (
	config "factorySimulator/configuration"
	"fmt"
	"time"
)

// worker represents a corporation-worker
// products is the fan-in where pre product is sent.
// task_server is the resource where worker can get new tasks.
func worker(workerID int, taskRequests chan<- taskRequest, products chan<- product) {
	for {
		//var taskToDo task

		// Prepare and send new request
		request := taskRequest{response: make(chan task)}
		taskRequests <- request

		// Check response
		taskToDo := <-request.response
		if taskToDo.operation == nil {
			continue
		}

		// Make new product
		val := taskToDo.operation.Execute(taskToDo.args...)
		if val == nil {
			fmt.Printf("Worker %d is destroying invalid task: %d %c\n", workerID, taskToDo.args, taskToDo.operation.Signature)
			continue
		}
		newProduct := product{value: val}

		if config.IsVerboseModeOn {
			fmt.Printf("\u001b[32mWorker\u001b[0m %d made product: %d %c = %d\n", workerID, taskToDo.args,
				taskToDo.operation.Signature, val)
		}
		products <- newProduct

		time.Sleep(config.WorkerDelay)
	}
}
