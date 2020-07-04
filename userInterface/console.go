package main

import (
	"bufio"
	"fmt"
	"goDrone/utils"
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
			if len(args) == 1 {
				intAlt, err := strconv.Atoi(args[1])
				if err == nil {
					result := rcfNodeClient.ServiceExec(client, "takeoff", utils.IntToByteArray(int64(intAlt)))
					fmt.Println(string(result))
				} else {
					println("takoff alt conv error")
				}
			}
		} else if string(args[0]) == "land" {

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
