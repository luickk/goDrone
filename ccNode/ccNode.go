package main

import (
	utils "goDrone/utils/utils"
	// naza "goNazaV2Interface/goNazaV2Interface"
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

	// arming motors 
	rcfNode.ActionCreate(nodeInstance, "armmotors", func(params []byte, n rcfNode.Node) {
		
	})


	// initiating fly to latitude longitude service
	// service args: alt, int32
	rcfNode.ServiceCreate(nodeInstance, "flytolatlon", func(params []byte, n rcfNode.Node) []byte {
		lat, lon, _ := utils.DecodeLatLonAlt(params)
		if lat != 0 && lon != 0 {
			InfoLogger.Println("flying to lat/ lon:", lat, lon)
		} else {
			WarningLogger.Println("flytolatlon deconding err")
		}
		return []byte("flew to given lat lon")
	})

	// initiating service to take off with the drone
	rcfNode.ServiceCreate(nodeInstance, "takeoff", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			alt := utils.ByteArrayToInt(params)
			InfoLogger.Println("taking off to height ", alt)
		}
		return []byte("taken off")
	})

	// initiating service to land drone
	rcfNode.ServiceCreate(nodeInstance, "land", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("landing")
		return []byte("landed")
	})

	// initiating service to turn drone
	rcfNode.ServiceCreate(nodeInstance, "turnto", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			deg := utils.ByteArrayToInt(params)
			InfoLogger.Println("turning to ", deg)
		}
		return []byte("turned")
	})

	// initiating service to change altitude
	rcfNode.ServiceCreate(nodeInstance, "changealt", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			alt := utils.ByteArrayToInt(params)
			InfoLogger.Println("changing alt to ", alt)
		}
		InfoLogger.Println("changing alt")
		println(string(params))
		return []byte("changed alt")
	})

	// initiating service to hold current drones position
	rcfNode.ServiceCreate(nodeInstance, "holdpos", func(params []byte, n rcfNode.Node) []byte {
		InfoLogger.Println("holding pos")
		println(string(params))
		return []byte("holding pos")
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
