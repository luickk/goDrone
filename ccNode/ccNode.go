package main

import (
	// utils "goDrone/utils"
	"log"
	"os"
	rcfNode "rcf/rcfNode"
)

// basic logger declarations
var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	InfoLogger = log.New(os.Stdout, "[ccNode] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[ccNode] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[ccNode] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// disableing debug information
	// InfoLogger.SetOutput(ioutil.Discard)
	// ErrorLogger.SetOutput(ioutil.Discard)
	// WarningLogger.SetOutput(ioutil.Discard)

	// creating node instance object which contains node struct in which all intern comm channels and topic/ action data maps are contained
	nodeInstance := rcfNode.Create(30)

	// initiating node by opening tcp server on node id
	// strarting action and topic handlers
	rcfNode.Init(nodeInstance)

	// initiating fly to latitude longitude service
	rcfNode.ServiceCreate(nodeInstance, "flyToLatLon", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("flying to lat lon ")
		println(string(params))
		return []byte("")
	})

	// initiating service to take off with the drone
	rcfNode.ServiceCreate(nodeInstance, "takeOff", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("taking off")
		println(string(params))
		return []byte("")
	})

	// initiating service to land drone
	rcfNode.ServiceCreate(nodeInstance, "land", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("landing")
		println(string(params))
		return []byte("")
	})

	// initiating service to turn drone
	rcfNode.ServiceCreate(nodeInstance, "turnTo", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("turning to")
		println(string(params))
		return []byte("")
	})

	// initiating service to change altitude
	rcfNode.ServiceCreate(nodeInstance, "changeAlt", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("changing alt")
		println(string(params))
		return []byte("")
	})

	// initiating service to hold current drones position
	rcfNode.ServiceCreate(nodeInstance, "holdPos", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("holding pos")
		println(string(params))
		return []byte("")
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
