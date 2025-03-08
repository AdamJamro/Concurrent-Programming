package corporation

import (
	config "factorySimulator/configuration"
	"fmt"
	"time"
)

// boss generates new task's for the workers.
// tasks channel is a fan-out where tasks are being sent.
func boss(tasks chan<- task) {
	for {
		firstArg := config.GetRandomIntArgument()
		secondArg := config.GetRandomIntArgument()
		operation := config.GetRandomOperation()
		newTask := task{
			args:      []any{firstArg, secondArg},
			operation: operation,
		}

		if config.IsVerboseModeOn {
			fmt.Printf("\u001b[31mBoss\u001b[0m added new task %d %c %d\n",
				firstArg, operation.Signature, secondArg)
		}

		tasks <- newTask // blocks until task queue server is able to receive new task
		// only after the task is accepted by server boss procrastinates
		time.Sleep(config.GetBossDelay())
	}
}
