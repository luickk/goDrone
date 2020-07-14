package main

import (
	"bufio"
	"fmt"
	utils "goDrone/utils/utils"
	"os"
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

func main() {
	STATELESS = false

	airborne = false

	ccClient, ccConnected := rcfNodeClient.NodeOpenConn(1050)
	gpsClient, gpsConnected := rcfNodeClient.NodeOpenConn(1051)

	if !ccConnected {
		println("cc conn failed")
	}  else if !gpsConnected {
		println("gps conn failed")
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
					gpsClient, gpsConnected = rcfNodeClient.NodeOpenConn(31)
				} else if args[1] == "cc" {
					ccClient, ccConnected = rcfNodeClient.NodeOpenConn(30)
				}
				if !ccConnected {
					println("cc conn failed")
				}  else if !gpsConnected {
					println("gps conn failed")
				}
			}
		} else if string(args[0]) == "takeoff" && ccConnected{
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					if !airborne && !STATELESS {
						result := rcfNodeClient.ServiceExec(ccClient, "takeoff", utils.IntToByteArray(int64(intAlt)))
						airborne = true
						fmt.Println(string(result))	
					} else {
						println("can not take of if airborne")
					}
				} else {
					println("takoff alt conv error")
				}
			} else {
				println("missing arg alt for service takeoff")
			}
		} else if string(args[0]) == "land"  && ccConnected{
			if airborne && !STATELESS {	
				rcfNodeClient.ActionExec(ccClient, "land", []byte(""))
			} else {
				println("can only land if airborne")
				airborne = false
			}
		} else if string(args[0]) == "markhomepos"  && gpsConnected{
			if !airborne && !STATELESS {	
				rcfNodeClient.ActionExec(ccClient, "markhomepos", []byte(""))
			} else {
				println("cannot set home pos if airborne")
			}
		} else if string(args[0]) == "turnto" && ccConnected {
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					if airborne && !STATELESS {	
						result := rcfNodeClient.ServiceExec(ccClient, "turnto", utils.IntToByteArray(int64(intAlt)))
						fmt.Println(string(result))
					} else {
						println("can only rotate drone if airbrone")
					}
				} else {
					println("turnto deg conv error")
				}
			} else {
				println("missing arg deg for service turnto")
			}

		} else if string(args[0]) == "flytolatlon" && ccConnected{
			if len(args) == 3 {
				lat, latErr := strconv.ParseFloat(args[1], 64)
				lon, lonErr := strconv.ParseFloat(args[2], 64)
				if latErr == nil && lonErr == nil {
					if airborne && !STATELESS {	
						result := rcfNodeClient.ServiceExec(ccClient, "flytolatlon", utils.EncodeLatLonAlt(lat, lon, 0))
						fmt.Println(string(result))
					} else {
						println("can only change location if airborne")
					}
				} else {
					println("takoff alt conv error")
				}
			} else {
				println("missing arfs lat lon for service flytolatlon")
			}
		} else if string(args[0]) == "listtopics"  && ccConnected{
			if len(args) >= 2 {
				data_map := make(map[string]string)
				data_map["cli"] = args[2]
				rcfNodeClient.TopicPublishGlobData(ccClient, args[1], data_map)
			}
		} else if string(args[0]) == "setneutral"  && ccConnected{
			if !airborne && !STATELESS {	
				rcfNodeClient.ActionExec(ccClient, "setneutral", []byte(""))
			} else {
				println("can only set to neutral if on ground")
			}
		} else if string(args[0]) == "setstate"  && ccConnected{
			if len(args) >= 2 {
				if args[1] == "airborne" {
					if args[2] == "true" {
						airborne = true
					} else if args[2] == "false" {
						airborne = false
					}
				}
			}
		} else if string(args[0]) == "getstates" {
			fmt.Println("ariborne: ", airborne)
			fmt.Println("gpsConnected: ", gpsConnected)
			fmt.Println("ccConnected: ", ccConnected)
		} else if string(args[0]) == "getgps"  && gpsConnected{
			elements := rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")
			fmt.Println(elements)
		}  else if string(args[0]) == "endcom" {
			if len(args) >= 0 {
				rcfNodeClient.NodeCloseConn(ccClient)
				rcfNodeClient.NodeCloseConn(gpsClient)
				return
			}
		} else {
			fmt.Println("command not known")	
		}
	}

	rcfNodeClient.NodeCloseConn(ccClient)
}
