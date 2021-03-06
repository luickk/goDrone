package main

import (
	utils "goDrone/utils/utils"
	"log"
	"os"
	"net/http"
	"net/url"
	rcfNodeClient "rcf/rcfNodeClient"
	"strconv"
)

var ccConnected, gpsConnected bool
var gpsClient, ccClient rcfNodeClient.Client

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

	InfoLogger = log.New(os.Stdout, "[WEB CLIENT] INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, "[WEB CLIENT] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "[WEB CLIENT] ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	ccClient, ccConnected = rcfNodeClient.NodeOpenConn(1050)
	gpsClient, gpsConnected = rcfNodeClient.NodeOpenConn(1051)

	if !ccConnected {
		ErrorLogger.Println("cc conn failed")
	}
	if !gpsConnected {
		ErrorLogger.Println("gps conn failed")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./userInterface/webInterface/static")))
	mux.HandleFunc("/reconnect", reconnectHandler)
	mux.HandleFunc("/takeOff", takeOffHandler)
	mux.HandleFunc("/land", landHandler)
	mux.HandleFunc("/markHomePos", markHomePos)
	mux.HandleFunc("/turnto", turntoHandler)
	mux.HandleFunc("/flytolatlon", flytolatlonHandler)
	mux.HandleFunc("/listTopics", listTopicsHandler)
	mux.HandleFunc("/setNeutral", setNeutralHandler)
	mux.HandleFunc("/setState", setStateHandler)
	mux.HandleFunc("/getState", getStateHandler)
	mux.HandleFunc("/getGpsPos", getGpsPosHandler)
	mux.HandleFunc("/changeAlt", changeAltHandler)
	mux.HandleFunc("/endcom", endcomHandler)
	
	http.ListenAndServe(":80", mux)
}

func reconnectHandler(w http.ResponseWriter, r *http.Request) {
	gpsClient, gpsConnected = rcfNodeClient.NodeOpenConn(1051)
	ccClient, ccConnected = rcfNodeClient.NodeOpenConn(1050)
	if !ccConnected {
		ErrorLogger.Println("cc conn failed")
		w.Write([]byte("could not connect to cc Node"))
	} else if !gpsConnected {
		ErrorLogger.Println("gps conn failed")
		w.Write([]byte("could not connect to gps Node"))
	}
}
func takeOffHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	intAlt, err := strconv.Atoi(params.Get("alt"))

	if err != nil {
		w.Write([]byte("take off alt conv err"))
		WarningLogger.Println("take off alt conv err")
		return
	}
	if !airborne && !STATELESS {
		rcfNodeClient.ServiceExec(ccClient, "takeoff", utils.IntToByteArray(int64(intAlt)))
		airborne = true
		InfoLogger.Println("taken off")
		w.Write([]byte("taken off to: " + strconv.Itoa(intAlt)))
	} else {
		WarningLogger.Println("can not take off if airborne")
		w.Write([]byte("can not take off if airborne"))
	}
}

func landHandler(w http.ResponseWriter, r *http.Request) {
	if airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "land", []byte(""))
		w.Write([]byte("set control mode to recovery"))
		InfoLogger.Println("set control mode to recovery")
		airborne = false
	} else {
		WarningLogger.Println("can only land if airborne")
		w.Write([]byte("can only land if airborne"))
		airborne = true
	}
}

func markHomePos(w http.ResponseWriter, r *http.Request) {
	if !airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "markhomepos", []byte(""))
		w.Write([]byte("set home pos"))
		InfoLogger.Println("set home pos")
	} else {
		WarningLogger.Println("cannot set home pos if airborne")
		w.Write([]byte("cannot set gome pos if airborne"))
	}
}

func turntoHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	heading, err := strconv.Atoi(params.Get("heading"))

	if err != nil {
		WarningLogger.Println("missing arg heading for service turnto")
		w.Write([]byte("missing arg heading for service turnto"))
		return
	}

	if err == nil {
		if airborne && !STATELESS {
			rcfNodeClient.ServiceExec(ccClient, "turnto", utils.IntToByteArray(int64(heading)))
			w.Write([]byte("turned to given heading"))
			InfoLogger.Println("turned to given heading")
		} else {
			InfoLogger.Println("can only rotate drone if airbrone")
			w.Write([]byte("can only rotate drone if airbrone"))
		}
	} else {
		WarningLogger.Println("turnto heading conv error")
		w.Write([]byte("turnto heading conv error"))
	}

}

func changeAltHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	alt, err := strconv.Atoi(params.Get("alt"))

	if err != nil {
		WarningLogger.Println("missing arg alt for service change alt")
		w.Write([]byte("missing arg alt for service change alt"))
		return
	}

	if err == nil {
		if airborne && !STATELESS {
			rcfNodeClient.ServiceExec(ccClient, "changealt", utils.IntToByteArray(int64(alt)))
			w.Write([]byte("reached given alt"))
			InfoLogger.Println("reached given alt")
		} else {
			InfoLogger.Println("can only change alt if airbrone")
			w.Write([]byte("can only change alt if airbrone"))
		}
	} else {
		WarningLogger.Println("change alt conv error")
		w.Write([]byte("change alt conv error"))
	}

}

func flytolatlonHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	lat, errLon := strconv.ParseFloat(params.Get("lat"), 64)
	lon, errLat := strconv.ParseFloat(params.Get("lon"), 64)

	if errLat == nil || errLon == nil {
		WarningLogger.Println("missing arg for service flyToLatLon")
		w.Write([]byte("missing arg for service flyToLatLon"))
		return
	}

	if airborne && !STATELESS {
		rcfNodeClient.ServiceExec(ccClient, "flytolatlon", utils.EncodeLatLonAlt(lat, lon, 0))
		InfoLogger.Println("flew to lat lon")
	} else {
		InfoLogger.Println("can only change location if airborne")
		w.Write([]byte("can only change location if airborne"))
	}
}

func listTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// data_map := make(map[string]string)
	// data_map["cli"] = args[2]
	// rcfNodeClient.TopicPublishGlobData(ccClient, args[1], data_map)

}
func setNeutralHandler(w http.ResponseWriter, r *http.Request) {
	if !airborne && !STATELESS {
		rcfNodeClient.ActionExec(ccClient, "setneutral", []byte(""))
		w.Write([]byte("set stick pos to neutral"))
		InfoLogger.Println("set stick pos to neutral")
	} else {
		InfoLogger.Println("can only set to neutral if not airborne")
		w.Write([]byte("can only set to neutral if not airborne"))
	}
}

func setStateHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	state := params.Get("state")
	val := params.Get("val")

	if state == "airborne" {
		if val == "true" {
			airborne = true
			w.Write([]byte("set state airborne to true"))
			InfoLogger.Println("set state airborne to true")
		} else if val == "false" {
			airborne = false
			w.Write([]byte("set state airborne to false"))
			InfoLogger.Println("set state airborne to false")
		}
	}
}
func getStateHandler(w http.ResponseWriter, r *http.Request) {
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		ErrorLogger.Println("Internal server error")
		return
	}

	params := parsedURL.Query()
	state := params.Get("state")

	if state == "airborne" {
		InfoLogger.Println("airborne: ", airborne)
		w.Write([]byte(strconv.FormatBool(airborne)))
	} else if state == "gpsconnected" {
		InfoLogger.Println("gpsConnected: ", gpsConnected)
		w.Write([]byte(strconv.FormatBool(gpsConnected)))
	} else if state == "ccconnected" {
		InfoLogger.Println("ccConnected: ", ccConnected)
		w.Write([]byte(strconv.FormatBool(ccConnected)))
	}
}
func getGpsPosHandler(w http.ResponseWriter, r *http.Request) {
	gpsSlice := rcfNodeClient.TopicPullGlobData(gpsClient, 1, "gpsData")
	InfoLogger.Println(gpsSlice)
}
func endcomHandler(w http.ResponseWriter, r *http.Request) {
	rcfNodeClient.NodeCloseConn(ccClient)
	rcfNodeClient.NodeCloseConn(gpsClient)
}
