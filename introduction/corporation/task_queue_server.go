package corporation

import (
	"factorySimulator/configuration"
	"fmt"
)

// taskQueueServer handles list of pending task.
// sends tasks to workers 	   via taskRequests channel.
// receives boss' tasks 	   via tasks channel.
// shows list of pending tasks via info channel.
func taskQueueServer(
	taskRequests <-chan taskRequest,
	tasks <-chan task,
	info <-chan struct{},
) {
	// Optionally lock the current goroutine to its current OS thread
	// opting either way shouldn't affect the program in an observable degree
	//runtime.LockOSThread()
	//defer runtime.UnlockOSThread()

	// List of tasks to do.
	tasksToDo := make([]task, 0)

	// Infinite loop of worker.
	for {
		select {
		case request := <-taskRequests:
			// If queue of tasks is empty indicate it by nil value
			if len(tasksToDo) == 0 {
				request.response <- task{}
			} else {
				// Otherwise send first task from list
				request.response <- tasksToDo[0]
				tasksToDo = tasksToDo[1:]
			}
		case newTask := <-tasks:
			// If boss send new task add it to list of tasks to do.
			if len(tasksToDo) >= configuration.SizeOfQueue {
				if configuration.IsVerboseModeOn {
					fmt.Println("List of tasks is full!")
				}
			} else {
				tasksToDo = append(tasksToDo, newTask)
			}
		case <-info:
			// If user sends request show list of tasks
			if len(tasksToDo) == 0 {
				fmt.Println("List of tasks is empty!")
			} else {
				fmt.Println("Tasks waiting for workers:")
				for i := range tasksToDo {
					fmt.Printf("\u001b[36mTask\u001b[0m %d: %d %c\n", i, tasksToDo[i].args, tasksToDo[i].operation.Signature)
				}
			}
		}
	}
}
