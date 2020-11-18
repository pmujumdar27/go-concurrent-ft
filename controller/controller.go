package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Initializing Controller")
	if len(os.Args) != 3 {
		fmt.Printf("USAGE: %s <controllerPort> <numServers>\n", os.Args[0])
		os.Exit(1)
	}
	controllerPort = os.Args[1]
	controllerAddr := "localhost:" + controllerPort
	numServers, err := strconv.Atoi(os.Args[2])
	handleError(err, "Invalid numServers")

	controller, err := net.Listen("tcp", controllerAddr)
	handleError(err, "Listening error")

	defer controller.Close()

	fmt.Println("Controller listening on:", controllerAddr)
	fmt.Printf("Waiting for %d servers\n", numServers)

	runController(controller, controllerAddr, int64(numServers))
}
