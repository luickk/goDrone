package main

import (
	"fmt"
	utils "goDrone/utils/utils"
	"net/http"
	rcfNodeClient "rcf/rcfNodeClient"
	"strconv"
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
	}
	if !gpsConnected {
		println("gps conn failed")
	}

	http.Handle("/", http.FileServer(http.Dir("./userInterface/webInterface/static")))
	mux := http.NewServeMux()
	mux.HandleFunc("/reconnect", reconnectHandler)
	mux.HandleFunc("/takeOffHandler", takeOffHandler)
	mux.HandleFunc("/landHandler", landHandler)
	mux.HandleFunc("/markHomePos", markHomePos)
	mux.HandleFunc("/turntoHandler", turntoHandler)
	mux.HandleFunc("/flytolatlonHandler", flytolatlonHandler)
	mux.HandleFunc("/listTopicsHandler", listTopicsHandler)
	mux.HandleFunc("/setNeutralHandler", setNeutralHandler)
	mux.HandleFunc("/setStateHandler", setStateHandler)
	mux.HandleFunc("/getGpsPosHandler", getGpsPosHandler)
	mux.HandleFunc("/endcomHandler", endcomHandler)

	http.ListenAndServe(":80", nil)
}

func reconnectHandler(w http.ResponseWriter, r *http.Request) {

	if string(args[0]) == "reconnect" {
		if len(args) == 2 {
			if args[1] == "gps" {
				gpsClient, gpsConnected = rcfNodeClient.NodeOpenConn(1051)
			} else if args[1] == "cc" {
				ccClient, ccConnected = rcfNodeClient.NodeOpenConn(1050)
			}
			if !ccConnected {
				println("cc conn failed")
			} else if !gpsConnected {
				println("gps conn failed")
			}
		}
	}
}
func takeOffHandler(w http.ResponseWriter, r *http.Request) {
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
}

func landHandler(w http.ResponseWriter, r *http.Request) {
	if airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "land", []byte(""))
	} else {
		println("can only land if airborne")
		airborne = false
	}
}

func markHomePos(w http.ResponseWriter, r *http.Request) {
	if !airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "markhomepos", []byte(""))
	} else {
		println("cannot set home pos if airborne")
	}
}

func turntoHandler(w http.ResponseWriter, r *http.Request) {
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

}

func flytolatlonHandler(w http.ResponseWriter, r *http.Request) {
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
}

func listTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if len(args) >= 2 {
		data_map := make(map[string]string)
		data_map["cli"] = args[2]
		rcfNodeClient.TopicPublishGlobData(ccClient, args[1], data_map)
	}
}
func setNeutralHandler(w http.ResponseWriter, r *http.Request) {
	if !airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "setneutral", []byte(""))
	} else {
		println("can only set to neutral if on ground")
	}
}
func setStateHandler(w http.ResponseWriter, r *http.Request) {
	if len(args) >= 2 {
		if args[1] == "airborne" {
			if args[2] == "true" {
				airborne = true
			} else if args[2] == "false" {
				airborne = false
			}
		}
	}
}
func listStatesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ariborne: ", airborne)
	fmt.Println("gpsConnected: ", gpsConnected)
	fmt.Println("ccConnected: ", ccConnected)
}
func getGpsPosHandler(w http.ResponseWriter, r *http.Request) {
	elements := rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")
	fmt.Println(elements)
}
func endcomHandler(w http.ResponseWriter, r *http.Request) {
	if len(args) >= 0 {
		rcfNodeClient.NodeCloseConn(ccClient)
		rcfNodeClient.NodeCloseConn(gpsClient)
		return
	}
}
