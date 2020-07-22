package main

import (
	"bufio"
	"fmt"
	utils "goDrone/utils/utils"
	"os"
	"log"
	rcfNodeClient "rcf/rcfNodeClient"
	"strconv"
	"strings"
)

// defines wether client stores states and blocks possibly dangerous service/action executions
var STATELESS bool

// drone states
var (
	airborne bool
)

// basic logger declarations
var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	STATELESS = false

	airborne = false

	InfoLogger = log.New(os.Stdout, "[CONSOLE CLIENT] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[CONSOLE CLIENT] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[CONSOLE CLIENT] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	ccClient, ccConnected := rcfNodeClient.NodeOpenConn(1050)
	gpsClient, gpsConnected := rcfNodeClient.NodeOpenConn(1051)

	if !ccConnected {
		ErrorLogger.Println("cc conn failed")
	}
	if !gpsConnected {
		ErrorLogger.Println("gps conn failed")
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		cmd_txt, _ := reader.ReadString('\n')
		cmd_txt = strings.Replace(cmd_txt, "\n", "", -1)
		args := strings.Split(cmd_txt, " ")
		if string(args[0]) == "reconnect" {
			if len(args) == 2 {
				if args[1] == "gps" {
					gpsClient, gpsConnected = rcfNodeClient.NodeOpenConn(1051)
				} else if args[1] == "cc" {
					ccClient, ccConnected = rcfNodeClient.NodeOpenConn(1050)
				}
				if !ccConnected {
					ErrorLogger.Println("cc conn failed")
				} else if !gpsConnected {
					ErrorLogger.Println("gps conn failed")
				}
			}
		} else if string(args[0]) == "takeoff" && ccConnected {
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					if !airborne && !STATELESS {
						rcfNodeClient.ServiceExec(ccClient, "takeoff", utils.IntToByteArray(int64(intAlt)))
						airborne = true
						InfoLogger.Println("taken off")
					} else {
						InfoLogger.Println("can not take of if airborne")
					}
				} else {
					WarningLogger.Println("takoff alt conv error")
				}
			} else {
				WarningLogger.Println("missing arg alt for service takeoff")
			}
		} else if string(args[0]) == "land" && ccConnected {
			if airborne && !STATELESS {
				rcfNodeClient.ActionExec(ccClient, "land", []byte(""))
				InfoLogger.Println("set control mode to recovery")
			} else {
				InfoLogger.Println("can only land if airborne")
				airborne = false
			}
		} else if string(args[0]) == "markhomepos" && gpsConnected {
			if !airborne && !STATELESS {
				rcfNodeClient.ActionExec(ccClient, "markhomepos", []byte(""))
				InfoLogger.Println("marked home pos")
			} else {
				InfoLogger.Println("cannot set home pos if airborne")
			}
		} else if string(args[0]) == "turnto" && ccConnected {
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					if airborne && !STATELESS {
						rcfNodeClient.ServiceExec(ccClient, "turnto", utils.IntToByteArray(int64(intAlt)))
						InfoLogger.Println("turned")
					} else {
						InfoLogger.Println("can only rotate drone if airbrone")
					}
				} else {
					WarningLogger.Println("turnto heading conv error")
				}
			} else {
				WarningLogger.Println("missing arg heading for service turnto")
			}

		} else if string(args[0]) == "flytolatlon" && ccConnected {
			if len(args) == 3 {
				lat, latErr := strconv.ParseFloat(args[1], 64)
				lon, lonErr := strconv.ParseFloat(args[2], 64)
				if latErr == nil && lonErr == nil {
					if airborne && !STATELESS {
						rcfNodeClient.ServiceExec(ccClient, "flytolatlon", utils.EncodeLatLonAlt(lat, lon, 0))
						InfoLogger.Println("flew to lat lon")
					} else {
						InfoLogger.Println("can only change location if airborne")
					}
				} else {
					WarningLogger.Println("takoff alt conv error")
				}
			} else {
				WarningLogger.Println("missing arfs lat lon for service flytolatlon")
			}
		} else if string(args[0]) == "listtopics" && ccConnected {
			if len(args) >= 2 {
				data_map := make(map[string]string)
				data_map["cli"] = args[2]
				rcfNodeClient.TopicPublishGlobData(ccClient, args[1], data_map)
			}
		} else if string(args[0]) == "setneutral" && ccConnected {
			if !airborne && !STATELESS {
				rcfNodeClient.ActionExec(ccClient, "setneutral", []byte(""))
				InfoLogger.Println("set stick pos to neutral")
			} else {
				InfoLogger.Println("can only set to neutral if on ground")
			}
		} else if string(args[0]) == "setstate" && ccConnected {
			if len(args) >= 2 {
				if args[1] == "airborne" {
					if args[2] == "true" {
						airborne = true
						InfoLogger.Println("set state airborne to true")
					} else if args[2] == "false" {
						airborne = false
						InfoLogger.Println("set state airborne to false")
					}
				}
			}
		} else if string(args[0]) == "getstates" {
			InfoLogger.Println("ariborne: ", airborne)
			InfoLogger.Println("gpsConnected: ", gpsConnected)
			InfoLogger.Println("ccConnected: ", ccConnected)
		} else if string(args[0]) == "getgps" && gpsConnected {
			gpsSlice := rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")
			InfoLogger.Println(gpsSlice)
		} else if string(args[0]) == "endcom" {
			if len(args) >= 0 {
				rcfNodeClient.NodeCloseConn(ccClient)
				rcfNodeClient.NodeCloseConn(gpsClient)
				InfoLogger.Println("closed all conns and quit console")
				return
			}
		} else {
			InfoLogger.Println("command not known")
		}
	}

	rcfNodeClient.NodeCloseConn(ccClient)
}
