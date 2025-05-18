package main

import (
	"grid-travelers-v5/config"
	simulator "grid-travelers-v5/travelers_simulator"
)

func main() {
	config.ValidateConfiguration()
	simulator.RunSimulation()
}
