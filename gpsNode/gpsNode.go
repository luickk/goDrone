package main

import (
	"fmt"
	"goDjiNazaGpsDecoder/djiNazaGpsDecoder"
	ellipsoid "goDrone/utils/ellipsoid"
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

	POSEMULATION = true

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
		dataMap["lat"] = string(utils.Float64bytes(float64(decInfo.Latitude)))
		dataMap["lon"] = string(utils.Float64bytes(float64(decInfo.Longitude)))
		dataMap["heading"] = strconv.Itoa(int(decInfo.Heading))
		dataMap["alt"] = strconv.Itoa(int(decInfo.Altitude))
		dataMap["speed"] = strconv.Itoa(int(decInfo.Speed))
		dataMap["sats"] = strconv.Itoa(int(decInfo.Satellites))

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

	rcfNode.ServiceCreate(nodeInstance, "calcDist", func(params []byte, node rcfNode.Node) []byte {
		lat1, lon1 := 37.619002, -122.374843 //SFO
		lat2, lon2 := 33.942536, -118.408074 //LAX

		// Create Ellipsoid object with WGS84-ellipsoid,
		// angle units are degrees, distance units are meter.
		geo1 := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)

		// Calculate the distance and bearing from SFO to LAX.
		distance, bearing := geo1.To(lat1, lon1, lat2, lon2)
		fmt.Printf("Distance = %v Bearing = %v\n", distance, bearing)

		// Calculate where you are when going from SFO in
		// direction 45.0 deg. for 20000 meters.
		lat3, lon3 := geo1.At(lat1, lon1, 20000.0, 45.0)
		fmt.Printf("lat3 = %v lon3 = %v\n", lat3, lon3)

		// Convert Lat-Lon-Alt to ECEF.
		lat4, lon4, alt4 := 39.197807, -77.108574, 55.0 // Some location near Baltimore
		// that the author of the Perl module geo-ecef used. I reused the coords of the tests.
		x, y, z := geo1.ToECEF(lat4, lon4, alt4)
		fmt.Printf("x = %v \ny = %v \nz = %v\n", x, y, z)

		// Convert ECEF to Lat-Lon-Alt.
		x1, y1, z1 := 1.1042590709397183e+06, -4.824765955871677e+06, 4.0093940281868847e+06
		lat5, lon5, alt5 := geo1.ToLLA(x1, y1, z1)
		fmt.Printf("lat5 = %v lon5 = %v alt3 = %v\n", lat5, lon5, alt5)

		return []byte("calced")
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
