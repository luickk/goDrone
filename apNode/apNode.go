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
	InfoLogger = log.New(os.Stdout, "[apNode] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[apNode] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[apNode] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// disableing debug information
	// InfoLogger.SetOutput(ioutil.Discard)
	// ErrorLogger.SetOutput(ioutil.Discard)
	// WarningLogger.SetOutput(ioutil.Discard)

	// creating node instance object which contains node struct in which all intern comm channels and topic/ action data maps are contained
	nodeInstance := rcfNode.Create(32)

	// initiating node by opening tcp server on node id
	// strarting action and topic handlers
	rcfNode.Init(nodeInstance)

	rcfNode.ServiceCreate(nodeInstance, "flyMission", func(params []byte, n rcfNode.Node) []byte {

		InfoLogger.Println("flying mission ")
		println(string(params))
		return []byte("")
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
