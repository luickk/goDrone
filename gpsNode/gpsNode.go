package main

import (
	"goDjiNazaGpsDecoder/djiNazaGpsDecoder"
	utils "goDrone/utils/utils"
	"log"
	"os"
	rcfNode "rcf/rcfNode"
	rcfUtil "rcf/rcfUtil"
	"strconv"
	"time"
)

var POSEMULATION bool

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

	POSEMULATION = false

	// disableing debug information
	// InfoLogger.SetOutput(ioutil.Discard)
	// ErrorLogger.SetOutput(ioutil.Discard)
	// WarningLogger.SetOutput(ioutil.Discard)

	// creating node instance object which contains node struct in which all intern comm channels and topic/ action data maps are contained
	nodeInstance := rcfNode.Create(1051)

	// initiating node by opening tcp server on node id
	// strarting action and topic handlers
	rcfNode.Init(nodeInstance)

	// creating topic by sending cmd to node
	rcfNode.TopicCreate(nodeInstance, "gpsData")

	type homePosition struct {
		Lat float64
		Lon float64
		Alt int
	}

	// declaring state holding vars

	// declare serial reader struct & decoding info struct
	var sRead djiNazaGpsDecoder.SerialRead
	var decInfo djiNazaGpsDecoder.DecodedInformation

	var homePos homePosition

	djiNazaGpsDecoder.OpenSerial(&sRead, "/dev/serial0")

	// loop to create sample data which is pushed to topic
	for {
		djiNazaGpsDecoder.ReadByte(&sRead, &decInfo)

		// fmt.Printf("Sats: %d \n", int(decInfo.Satellites))
		// fmt.Printf("Heading: %d \n", int(decInfo.Heading))
		// fmt.Printf("Alt: %d \n", int(decInfo.Altitude))
		// fmt.Printf("Speed: %d \n", int(decInfo.Speed))
		// fmt.Printf("Lat: %e \n", decInfo.Latitude)
		// fmt.Printf("Lon: %e \n", decInfo.Longitude)
		// fmt.Printf("Time: %d, %d, %d, %d, %d, %d \n", int(decInfo.Year), int(decInfo.Month), int(decInfo.Day), int(decInfo.Hour), int(decInfo.Minute), int(decInfo.Second))
		// fmt.Printf("HW Version: %d, SW Version: %d \n", int(decInfo.HardwareVersion.Version), int(decInfo.FirmwareVersion.Version))

		// putting decoded data into map
		dataMap := make(map[string]string)
		dataMap["lat"] = strconv.FormatFloat(float64(decInfo.Latitude), 'f', 5, 64)
		dataMap["lon"] = strconv.FormatFloat(float64(decInfo.Longitude), 'f', 5, 64)
		dataMap["heading"] = strconv.Itoa(int(decInfo.Heading))
		dataMap["alt"] = strconv.Itoa(int(decInfo.Altitude))
		dataMap["speed"] = strconv.Itoa(int(decInfo.Speed))
		dataMap["sats"] = strconv.Itoa(int(decInfo.Satellites))

		// fmt.Println(decInfo.Latitude)

		encodedData, err := rcfUtil.GlobMapEncode(dataMap)
		encodedDataSlice := []byte(encodedData.Bytes())

		// generating sample data for emulation
		sampleDataMap := make(map[string]string)
		sampleDataMap["lat"] = "49.45300997697536"
		sampleDataMap["lon"] = "10.96558038704124"
		sampleDataMap["heading"] = strconv.Itoa(5)
		sampleDataMap["alt"] = strconv.Itoa(60)
		sampleDataMap["speed"] = strconv.Itoa(0)
		sampleDataMap["sats"] = strconv.Itoa(8)

		encodedSampleData, _ := rcfUtil.GlobMapEncode(sampleDataMap)
		sampleDataMapSlice := []byte(encodedSampleData.Bytes())

		if err != nil {
			WarningLogger.Println("GlobMapEncode encoding error")
			WarningLogger.Println(err)
		} else {
			if POSEMULATION {
				// pushing emulated data map node
				rcfNode.TopicPublishData(nodeInstance, "gpsData", sampleDataMapSlice)
			} else {
				// pushing encoded data map to node
				rcfNode.TopicPublishData(nodeInstance, "gpsData", encodedDataSlice)
			}
		}

		// euals 10 Hz
		time.Sleep(100 * time.Millisecond)
	}

	rcfNode.ActionCreate(nodeInstance, "markhome", func(params []byte, node rcfNode.Node) {
		homePos.Lat, homePos.Lon, homePos.Alt = utils.DecodeLatLonAlt(params)
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
