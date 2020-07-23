package main

import (
	// utils "goDrone/utils"
	"log"
    "io/ioutil"
	"os" 
	"strings"
	rcfNode "rcf/rcfNode"
)

// basic logger declarations
var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

var MissionFolderPath string

func main() {
	InfoLogger = log.New(os.Stdout, "[Auto Pilot Node] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[Auto Pilot Node] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[Auto Pilot Node] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	MissionFolderPath = "/home/pi/missions/"

	// disableing debug information
	// InfoLogger.SetOutput(ioutil.Discard)
	// ErrorLogger.SetOutput(ioutil.Discard)
	// WarningLogger.SetOutput(ioutil.Discard)

	// creating node instance object which contains node struct in which all intern comm channels and topic/ action data maps are contained
	nodeInstance := rcfNode.Create(1052)

	// initiating node by opening tcp server on node id
	// strarting action and topic handlers
	rcfNode.Init(nodeInstance)
	
	// params: string name
	rcfNode.ServiceCreate(nodeInstance, "executeMission", func(params []byte, n rcfNode.Node) []byte {
		missionName := string(params)
		InfoLogger.Println("flying mission " + missionName)
		
		data, err := ioutil.ReadFile(MissionFolderPath+missionName)
		if err != nil {
			ErrorLogger.Println("File reading error", err)
		}
		InfoLogger.Println("Contents of file:", string(data))
		
		for instruction := range data {
			instructionArgs := strings.Split(string(instruction), " ")
			if len(instructionArgs) > 0 {
				if instructionArgs[0] == "takeoff" {
					InfoLogger.Println("taking off to " + instructionArgs[1])
				} else if instructionArgs[0] == "land" {
					InfoLogger.Println("ladning")

				} else if instructionArgs[0] == "turnto" {
					InfoLogger.Println("turning to " + instructionArgs[1])

				} else if instructionArgs[0] == "changealt" {
					InfoLogger.Println("changing alt to " + instructionArgs[1])

				} else if instructionArgs[0] == "holdpos" {
					InfoLogger.Println("holding position" + instructionArgs[1])

				} else if instructionArgs[0] == "wait" {
					InfoLogger.Println("waiting for " + instructionArgs[1])

				}
			}
		}

		return []byte("")
	})

	rcfNode.ServiceCreate(nodeInstance, "listMission", func(params []byte, n rcfNode.Node) []byte {
		InfoLogger.Println("listing missions in " + MissionFolderPath)
		
		file, err := os.Open(MissionFolderPath)
		if err != nil {
			log.Fatalf("failed opening directory: %s", err)
		}
		defer file.Close()
		fileList,_ := file.Readdirnames(0)
		

		return []byte(strings.Join(fileList, ","))
	})

	// halting node so it doesn't quit
	rcfNode.NodeHalt()
}
