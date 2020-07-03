package main

import (
	// utils "goDrone/utils"
	"log"
	"math/rand"
	"os"
	rcfNode "rcf/rcfNode"
	rcfUtil "rcf/rcfUtil"
	"strconv"
	"time"
)

// basic logger declarations
var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	InfoLogger = log.New(os.Stdout, "[gpsNode] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[gpsNode] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[gpsNode] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// disableing debug information
	// InfoLogger.SetOutput(ioutil.Discard)
	// ErrorLogger.SetOutput(ioutil.Discard)
	// WarningLogger.SetOutput(ioutil.Discard)

	// creating node instance object which contains node struct in which all intern comm channels and topic/ action data maps are contained
	nodeInstance := rcfNode.Create(31)

	// initiating node by opening tcp server on node id
	// strarting action and topic handlers
	rcfNode.Init(nodeInstance)

	// creating topic by sending cmd to node
	rcfNode.TopicCreate(nodeInstance, "gpsData")

	// loop to create sample data which is pushed to topic
	for {
		// generating random int
		alt := rand.Intn(200)

		// putting sample data into map
		dataMap := make(map[string]string)
		dataMap["lat"] = strconv.Itoa(alt)
		dataMap["lon"] = strconv.Itoa(alt)
		dataMap["heading"] = strconv.Itoa(alt)
		dataMap["alt"] = strconv.Itoa(alt)
		dataMap["speed"] = strconv.Itoa(alt)

		encodedData, err := rcfUtil.GlobMapEncode(dataMap)
		encodedDataSlice := []byte(encodedData.Bytes())
		if err != nil {
			WarningLogger.Println("GlobMapEncode encoding error")
			WarningLogger.Println(err)
		} else {
			// pushing alt value to node, encoded as string. every sent string/ alt value represents one element/ msg in the topic
			rcfNode.TopicPublishData(nodeInstance, "gpsData", encodedDataSlice)
		}

		// euals 10 Hz
		time.Sleep(100 * time.Millisecond)
	}

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
