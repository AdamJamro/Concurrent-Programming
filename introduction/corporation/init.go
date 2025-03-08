package corporation

import (
	config "factorySimulator/configuration"
	"flag"
	"fmt"
	"os"
	"runtime"
)

func parseFlags() {
	silentMode := flag.Bool("s", false, "deactivates verbose mode and activates silent mode")
	displayInfo := flag.Bool("i", false, "displays information about the program")

	flag.Parse()

	if *silentMode {
		config.IsVerboseModeOn = false
	}
	if *displayInfo {
		fmt.Println("Configuration:")
		fmt.Println("Number of workers:", config.NumOfWorkers)
		fmt.Println("Number of clients:", config.NumOfClients)
		fmt.Println("Size of list of tasks:", config.SizeOfQueue)
		fmt.Println("Delay for client:", config.ClientDelay)
		fmt.Println("Delay for boss (upper bound):", config.BossDelayUpperBound)
		fmt.Println("Delay for worker:", config.WorkerDelay)
		fmt.Println("Size of storage:", config.SizeOfTaskStorage)
		fmt.Println("\u001b[31mTo run in silent mode set the -s flag!\u001b[0m")
		os.Exit(0)
	}
}

func printCommands() {
	fmt.Println("Usage:")
	fmt.Println("m - print list of stored products")
	fmt.Println("t - print list of tasks to do")
	fmt.Println("q - quit")
}

func Init() {
	parseFlags()

	fmt.Println("\u001b[32m####################################################")
	fmt.Println("\t\u001b[1mWelcome in corporation simulator\u001b[0m\u001b[32m")
	fmt.Println("####################################################\u001b[0m")
	fmt.Print("\n(Press enter to stop the simulation)\n\n")

	//Channel for new tasks from boss.
	bossNewTasksChannel := make(chan task)

	// deprecated:
	//taskQueueCapacitySemaphore := make(chan struct{}, config.SizeOfQueue)

	// Channels for worker
	workerTaskRequestsChannel := make(chan taskRequest)
	workerNewProductsChannel := make(chan product)

	// Channel for client
	clientPurchaseChannel := make(chan buyRequest)

	// Info channels
	tasksQueueServerInfoChannel := make(chan struct{})
	taskStorageServerInfoChannel := make(chan struct{})

	// run servers ordered by receivers first
	go taskQueueServer(
		workerTaskRequestsChannel,
		bossNewTasksChannel,
		tasksQueueServerInfoChannel,
	)
	go taskStorageServer(workerNewProductsChannel, clientPurchaseChannel, taskStorageServerInfoChannel)
	go boss(bossNewTasksChannel)
	for i := 0; i < config.NumOfWorkers; i++ {
		go worker(i, workerTaskRequestsChannel, workerNewProductsChannel)
	}
	for i := 0; i < config.NumOfClients; i++ {
		go client(i, clientPurchaseChannel)
	}

	if config.IsVerboseModeOn {
		// any action will terminate
		fmt.Scanln()
		fmt.Println("\u001b[33m################################################")
		fmt.Println("\t\u001b[1mSimulation is being stopped\u001b[0m\u001b[33m")
		fmt.Println("################################################\u001b[0m")
	} else { // silent mode
		runtime.LockOSThread()
		var cmd string
		printCommands()
		for {
			fmt.Scanln(&cmd)
			switch cmd {
			case "m":
				taskStorageServerInfoChannel <- struct{}{}
			case "t":
				tasksQueueServerInfoChannel <- struct{}{}
			case "h":
				printCommands()
			case "q":
				runtime.UnlockOSThread()
				return
			default:
				fmt.Println("Invalid command")
				fmt.Println("Type 'h' to see available commands")
			}
		}
	}
}
