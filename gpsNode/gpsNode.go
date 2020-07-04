package main

import (
	"fmt"
	"goDjiNazaGpsDecoder/djiNazaGpsDecoder"
	utils "goDrone/utils"
	"log"
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

	var sRead djiNazaGpsDecoder.SerialRead
	var decInfo djiNazaGpsDecoder.DecodedInformation

	djiNazaGpsDecoder.OpenSerial(&sRead, "/dev/serial0")

	// loop to create sample data which is pushed to topic
	for {
		djiNazaGpsDecoder.ReadByte(&sRead, &decInfo)

		fmt.Printf("Sats: %d \n", int(decInfo.Satellites))
		fmt.Printf("Heading: %d \n", int(decInfo.Heading))
		fmt.Printf("Alt: %d \n", int(decInfo.Altitude))
		fmt.Printf("Speed: %d \n", int(decInfo.Speed))
		fmt.Printf("Lat: %e \n", decInfo.Latitude)
		fmt.Printf("Lon: %e \n", decInfo.Longitude)
		fmt.Printf("Time: %d, %d, %d, %d, %d, %d \n", int(decInfo.Year), int(decInfo.Month), int(decInfo.Day), int(decInfo.Hour), int(decInfo.Minute), int(decInfo.Second))
		fmt.Printf("HW Version: %d, SW Version: %d \n", int(decInfo.HardwareVersion.Version), int(decInfo.FirmwareVersion.Version))

		// putting sample data into map
		dataMap := make(map[string]string)
		dataMap["lat"] = string(utils.Float64bytes(decInfo.Latitude))
		dataMap["lon"] = string(utils.Float64bytes(decInfo.Longitude))

		dataMap["heading"] = strconv.Itoa(int(decInfo.Heading))
		dataMap["alt"] = strconv.Itoa(int(decInfo.Altitude))
		dataMap["speed"] = strconv.Itoa(int(decInfo.Speed))
		dataMap["sats"] = strconv.Itoa(int(decInfo.Satellites))

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
