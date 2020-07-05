package main

import (
	"bufio"
	"fmt"
	utils "goDrone/utils"
	"os"
	rcfNodeClient "rcf/rcfNodeClient"
	"strconv"
	"strings"
)

func main() {
	client := rcfNodeClient.NodeOpenConn(30)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		cmd_txt, _ := reader.ReadString('\n')
		cmd_txt = strings.Replace(cmd_txt, "\n", "", -1)
		args := strings.Split(cmd_txt, " ")

		if string(args[0]) == "takeoff" {
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					result := rcfNodeClient.ServiceExec(client, "takeoff", utils.IntToByteArray(int64(intAlt)))
					fmt.Println(string(result))
				} else {
					println("takoff alt conv error")
				}
			} else {
				println("missing arg alt for service takeoff")
			}
		} else if string(args[0]) == "land" {

		} else if string(args[0]) == "turnto" {
			if len(args) == 2 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					result := rcfNodeClient.ServiceExec(client, "turnto", utils.IntToByteArray(int64(intAlt)))
					fmt.Println(string(result))
				} else {
					println("turnto deg conv error")
				}
			} else {
				println("missing arg deg for service turnto")
			}

		} else if string(args[0]) == "flytolatlon" {
			if len(args) == 3 {
				intLat, latErr := strconv.ParseFloat(args[1], 64)
				intLon, lonErr := strconv.ParseFloat(args[2], 64)
				if latErr == nil && lonErr == nil {
					result := rcfNodeClient.ServiceExec(client, "flytolatlon", append(utils.Float64bytes(intLat), utils.Float64bytes(intLon)...))
					fmt.Println(string(result))
				} else {
					println("takoff alt conv error")
				}
			} else {
				println("missing arfs lat lon for service flytolatlon")
			}
		} else if string(args[0]) == "listtopics" {
			if len(args) >= 2 {
				data_map := make(map[string]string)
				data_map["cli"] = args[2]
				rcfNodeClient.TopicPublishGlobData(client, args[1], data_map)
			}
		} else if string(args[0]) == "endcom" {
			if len(args) >= 0 {
				rcfNodeClient.NodeCloseConn(client)
				return
			}
		}
	}

	rcfNodeClient.NodeCloseConn(client)
}
