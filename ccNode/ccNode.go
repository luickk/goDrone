package main

import (
	utils "goDrone/utils/utils"
	naza "goNazaV2Interface/goNazaV2Interface"
	rcfNodeClient "rcf/rcfNodeClient"
	ellipsoid "goDrone/utils/ellipsoid"
	"time"
	"log"
	"os"
	"strconv"
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

	var interfaceConf naza.InterfaceConfig

	interfaceConf.StickDir = make(map[int]string)

	interfaceConf.StickDir[naza.Achannel] = "rev"
	interfaceConf.StickDir[naza.Echannel] = "norm"
	interfaceConf.StickDir[naza.Tchannel] = "norm"
	interfaceConf.StickDir[naza.Rchannel] = "rev"

	
	interfaceConf.LeftStickMaxPos = make(map[int]int)
	interfaceConf.NeutralStickPos = make(map[int]int)
	interfaceConf.RightStickMaxPos = make(map[int]int)

	// key: channel, value: stick max pos
	interfaceConf.LeftStickMaxPos[naza.Achannel] = 400
	interfaceConf.NeutralStickPos[naza.Achannel] = 315
	interfaceConf.RightStickMaxPos[naza.Achannel] = 320

	interfaceConf.LeftStickMaxPos[naza.Echannel] = 400
	interfaceConf.NeutralStickPos[naza.Echannel] = 315
	interfaceConf.RightStickMaxPos[naza.Echannel] = 220

	interfaceConf.LeftStickMaxPos[naza.Tchannel] = 400
	interfaceConf.NeutralStickPos[naza.Tchannel] = 0
	interfaceConf.RightStickMaxPos[naza.Tchannel] = 230

	interfaceConf.LeftStickMaxPos[naza.Rchannel] = 400
	interfaceConf.NeutralStickPos[naza.Rchannel] = 315
	interfaceConf.RightStickMaxPos[naza.Rchannel] = 220

	interfaceConf.GpsModeFlipSwitchDutyCycle = 390
	interfaceConf.FailsafeModeFlipSwitchDutyCycle = 350
	interfaceConf.SelectableModeFlipSwitchDutyCycle = 250

	if !naza.InitPCA9685(&interfaceConf) {
		ErrorLogger.Println("failed to init PCA9685")
	} else {
		if !naza.InitNaza(&interfaceConf) {
			ErrorLogger.Println("failed to init naza")
		}	
	}
	
	gpsClient, gpsConnected := rcfNodeClient.NodeOpenConn(31)
	if !gpsConnected {
		ErrorLogger.Println("could not connect to gps node")
	}

	// arming motors 
	rcfNode.ActionCreate(nodeInstance, "armmotors", func(params []byte, n rcfNode.Node) {
		naza.ArmMotors(&interfaceConf)
	})


	// initiating fly to latitude longitude service
	// service args: alt, int32
	rcfNode.ServiceCreate(nodeInstance, "flytolatlon", func(params []byte, n rcfNode.Node) []byte {
		targetLat, targetLon, _ := utils.DecodeLatLonAlt(params)
		if targetLat != 0 && targetLon != 0 {
			InfoLogger.Println("flying to lat/ lon:", targetLat, targetLon)
		} else {
			WarningLogger.Println("flytolatlon deconding err")
		}
		
		liveDiff, targetHeading := 0, 0
		targetDiffAccuracy := 30 
		targetDistanceAccuracy := 10
		yawTargetLocked, distanceTargetReached, isOriented := false, false, false
		targetDistance, bearing := 0.0, 0.0
			
		naza.SetYaw(&interfaceConf, 70)

		// Create Ellipsoid object with WGS84-ellipsoid,
		// angle units are degrees, distance units are meter.
		geo1 := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
		
		for !distanceTargetReached {
			currentLat, _ := strconv.ParseFloat(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["lat"], 64)
			currentLon, _ := strconv.ParseFloat(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["lon"], 64)
			
			// Calculate the distance and bearing from SFO to LAX.
			targetDistance, bearing = geo1.To(currentLat, currentLon, targetLat, targetLon)
			targetHeading = utils.BearingToHeadding(int(bearing))
			InfoLogger.Printf("Distance = %v Bearing = %v target Heading = %v \n", targetDistance, bearing, targetHeading)

			yawTargetLocked = false
			
			for !yawTargetLocked {
				currentDir, _ := strconv.Atoi(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["heading"])
				liveDiff = utils.CalcDiff(currentDir, int(targetHeading))
				if liveDiff <= targetDiffAccuracy {
					if !isOriented {
						isOriented = true
						yawTargetLocked = true
						naza.SetYaw(&interfaceConf, 50)

						naza.SetPitch(&interfaceConf, 70)

						InfoLogger.Println("turnto oriented to to target deg")	
					} else {
						yawTargetLocked = true
					}
				} else {
					isOriented = false
					if !isOriented {
						InfoLogger.Println("turnto angel diff(orienting): ", liveDiff)
						naza.SetPitch(&interfaceConf, 50)
						time.Sleep(time.Second * 1)
						naza.SetYaw(&interfaceConf, 70)			
					}
				}
			}	


			if int(targetDistance) <= targetDistanceAccuracy {
				distanceTargetReached = true
				InfoLogger.Println("flytolatlon reached target")
				naza.SetPitch(&interfaceConf, 50)
			} else {
				InfoLogger.Println("flytolatlon approaching target")
			}
			time.Sleep(1* time.Second)
		}
		
		return []byte("flew to given lat lon")
	})

	// initiating service to take off with the drone
	rcfNode.ServiceCreate(nodeInstance, "takeoff", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			alt := utils.ByteArrayToInt(params)
			InfoLogger.Println("taking off to height ", alt)

			naza.ArmMotors(&interfaceConf)
			naza.SetThrottle(&interfaceConf, 50)
			time.Sleep(5*time.Second)
			naza.SetThrottle(&interfaceConf, 60)
			time.Sleep(3*time.Second)
			
			if (naza.SetPitch(&interfaceConf, 0) && naza.SetRoll(&interfaceConf, 0) && naza.SetYaw(&interfaceConf, 0) && naza.SetThrottle(&interfaceConf, 50)) {
				InfoLogger.Println("set stick pos to neutral")
			} else {
				InfoLogger.Println("failed to set stick pos to neutral")
			}
			naza.SetFlightMode(&interfaceConf, "gps")

			// todo, alt loop
		}
		return []byte("taken off")
	})

	// initiating service to land drone
	rcfNode.ServiceCreate(nodeInstance, "land", func(params []byte, n rcfNode.Node) []byte {

		naza.SetFlightMode(&interfaceConf, "failsafe")

		InfoLogger.Println("land landing")
		return []byte("landed")
	})

	// initiating service to turn drone
	rcfNode.ServiceCreate(nodeInstance, "turnto", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			bearing := utils.ByteArrayToInt(params)
			InfoLogger.Println("turning to ", bearing)
			
			naza.SetYaw(&interfaceConf, 70)

			yawTargetLocked := false
			liveDiff := 360
			targetDiffAccuracy := 30 

			for !yawTargetLocked {
				currentDir, _ := strconv.Atoi(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["heading"])
				liveDiff = utils.CalcDiff(currentDir, int(bearing))
				if liveDiff <= targetDiffAccuracy {
					yawTargetLocked = true
					InfoLogger.Println("turnto turned to to target deg")
				}

				time.Sleep(1* time.Second)
			}
		}
		return []byte("turned")
	})

	// initiating service to change altitude
	rcfNode.ServiceCreate(nodeInstance, "changealt", func(params []byte, n rcfNode.Node) []byte {
		if len(params) == 8 {
			targetAlt := utils.ByteArrayToInt(params)
			InfoLogger.Println("changing alt to ", targetAlt)

			currentAlt, _ := strconv.Atoi(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["alt"])
			if currentAlt > int(targetAlt) {
				naza.SetThrottle(&interfaceConf, 40)
			} else if currentAlt < int(targetAlt) {
				naza.SetThrottle(&interfaceConf, 60)
			}
			
			altTargetLocked := false
			liveAlt := 360
			targetAltAccuracy := 30 

			for !altTargetLocked {
				currentAlt, _ := strconv.Atoi(rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")[0]["alt"])
				liveAlt = utils.CalcDiff(currentAlt, int(targetAlt))
				if liveAlt <= targetAltAccuracy {
					altTargetLocked = true
					InfoLogger.Println("changealt reached target alt")
				}
			}
		}
		InfoLogger.Println("changing alt")
		println(string(params))
		return []byte("changed alt")
	})

	// initiating action to hold current drones position
	rcfNode.ActionCreate(nodeInstance, "holdpos", func(params []byte, n rcfNode.Node) {
		InfoLogger.Println("holding pos")

		if (naza.SetPitch(&interfaceConf, 0) && naza.SetRoll(&interfaceConf, 0) && naza.SetYaw(&interfaceConf, 0) && naza.SetThrottle(&interfaceConf, 50)) {
			InfoLogger.Println("set stick pos to neutral")
		} else {
			InfoLogger.Println("failed to set stick pos to neutral")
		}
		naza.SetFlightMode(&interfaceConf, "gps")
	})

	// initiating action to set all channels to neutral position
	rcfNode.ActionCreate(nodeInstance, "setneutral", func(params []byte, n rcfNode.Node) {
		InfoLogger.Println("set stick positions to neutral")
		
		naza.SetNeutral(&interfaceConf)
		naza.SetFlightMode(&interfaceConf, "gps")
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
