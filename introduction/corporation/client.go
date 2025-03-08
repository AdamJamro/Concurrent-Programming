package corporation

import (
	config "factorySimulator/configuration"
	"fmt"
	"time"
)

// client makes the requests for buying new products
func client(clientID int, requests chan<- buyRequest) {
	for {
		request := buyRequest{response: make(chan product)}
		requests <- request

		response := <-request.response

		if response == (product{}) {
			continue
		}

		if config.IsVerboseModeOn {
			fmt.Printf("\u001b[34mClient\u001b[0m %d bought product with value %d\n", clientID, response.value)
		}

		time.Sleep(config.ClientDelay)
	}
}
