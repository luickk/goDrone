package main

import (
	utils "goDrone/utils"
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
	// service args: alt, int32
	rcfNode.ServiceCreate(nodeInstance, "flyToLatLon", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 16 {
			lat := utils.Float64frombytes(params[0:7])
			lon := utils.Float64frombytes(params[8:15])

			InfoLogger.Println("flying to lat %f lon %f", lat, lon)
		}
		return []byte("")
	})

	// initiating service to take off with the drone
	rcfNode.ServiceCreate(nodeInstance, "takeOff", func(params []byte, n rcfNode.Node) []byte {

		if len(params) == 8 {
			alt := utils.ByteArrayToInt(params)

			InfoLogger.Println("taking off to height %d", alt)
		}

		InfoLogger.Println("taking off")
		println(string(params))
		return []byte("")
	})

	// initiating service to land drone
	rcfNode.ServiceCreate(nodeInstance, "land", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("landing")
		return []byte("")
	})

	// initiating service to turn drone
	rcfNode.ServiceCreate(nodeInstance, "turnTo", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			deg := utils.ByteArrayToInt(params)
			InfoLogger.Println("turning to %d", deg)
		}
		return []byte("")
	})

	// initiating service to change altitude
	rcfNode.ServiceCreate(nodeInstance, "changeAlt", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			alt := utils.ByteArrayToInt(params)
			InfoLogger.Println("changing alt to %d", alt)
		}
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
