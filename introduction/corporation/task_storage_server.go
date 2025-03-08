package corporation

import (
	config "factorySimulator/configuration"
	"fmt"
	"sync"
)

// this cloud be well omitted and the taskStorageServer could be simulated by a buffered channel
// so it's only for demonstrating the use of mutexes and semaphores
type syncTaskStorage struct {
	products  []product
	semaphore chan struct{} // semaphore shall ensure only a set amount of products is stored at once
	mutex     sync.Mutex    // mutex shall synchronize each write and read from products
}

// taskStorageServer implements storing and trading products in the task-storehouse.
// accepts new products 	 	 via manufacturedProducts channel.
// handles purchase requests 	 via purchaseRequests channel.
// accepts display info requests via info channel.
func taskStorageServer(manufacturedProducts <-chan product, purchaseRequests <-chan buyRequest, info <-chan struct{}) {
	// Optionally lock the current goroutine to its current OS thread
	// opting either way shouldn't affect the program in an observable degree
	//runtime.LockOSThread()
	//defer runtime.UnlockOSThread()

	syncStorage := syncTaskStorage{
		make([]product, 0), // don't set capacity and use it as if it was a deque for simplicity
		make(chan struct{}, config.SizeOfTaskStorage), // allow only #SizeOfTaskStorage products to be stored at once
		sync.Mutex{}, // allow only one goroutine to access storedProducts at a time
	}

	// fetcherService is not part of select multiplexer due to the semaphore acquiring
	// fetching and store newly manufactured products on a separate goroutine
	go productFetcherService(&syncStorage, &manufacturedProducts)

	// in each iteration let select occasionally check print-info requests
	for {
		select {
		case request := <-purchaseRequests:
			handlePurchaseRequest(&syncStorage, &request)
		case <-info:
			go printInfo(&syncStorage)
		}
	}
}

func productFetcherService(syncStorage *syncTaskStorage, manufacturedProducts *<-chan product) {
	for {
		syncStorage.semaphore <- struct{}{} // if semaphore is full, storage capacity is full, so it needs to block
		newProduct := <-*manufacturedProducts
		syncStorage.mutex.Lock()
		syncStorage.products = append(syncStorage.products, newProduct)
		syncStorage.mutex.Unlock()

		// optionally print synchronized result:
		//if config.IsVerboseModeOn {
		//	fmt.Printf("\u001b[33mStorage\u001b[0m stored product with value %d\n", newProduct.value)
		//}
	}
}

func handlePurchaseRequest(syncStorage *syncTaskStorage, request *buyRequest) {
	if len(syncStorage.products) == 0 {
		request.response <- product{} // send empty product to indicate failure
		return
	}

	// this should be abstracted away but this
	syncStorage.mutex.Lock()
	request.response <- syncStorage.products[0]
	syncStorage.products = syncStorage.products[1:]
	syncStorage.mutex.Unlock()

	// release semaphore to indicate freed storage space
	<-syncStorage.semaphore
}

func printInfo(syncStorage *syncTaskStorage) {
	if len(syncStorage.products) == 0 {
		fmt.Println("List of products is empty!")
	} else {
		// optionally lock syncStorage to prevent race conditions...
		// ...even though we just want to output the "current" state
		// whatever "current" would actually mean
		syncStorage.mutex.Lock()
		fmt.Println("Stored products: ")
		for i := range syncStorage.products {
			fmt.Printf("\u001b[35mProduct\u001b[0m %d with value %d\n", i, syncStorage.products[i].value)
		}
		syncStorage.mutex.Unlock()
	}
}
