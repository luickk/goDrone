# Go Drone

A drone project in which the drone is controlled by a raspberry pi controlled DJI Naza V2 flight controller. The project builds up on three other projets off mine, the [robot communication interface](https://github.com/luickk/rcf), [goNazaV2Interface](https://github.com/luickk/goNazaV2Interface) and the [goDjiNazaGpsDecoder](https://github.com/luickk/goDjiNazaGpsDecoder).

The goDrone project aims to build an open source drone platform which can keep up with the latest drone technologies and as such supports the same functionality and reliability. Having that in mind, building a open source, self developed flight controller that can keep up with the newest tech respectively industry standrads is unrealistic. 

The alternative is to use existing hardware and make use of the work that's already been done by the big tech giants. To enable the use of flight controllers that don't offer access via. a computer I built the [go Naza interface](https://github.com/luickk/goNazaV2Interface) which when used with a Raspberry Pi emulates PWM signals, similar to those from a standard receiver.

## Core Nodes

The core nodes compromise of 3 core nodes which are the command & control Node(ccNode), gps Node and the auto pilot Node(apNode). These core services offer the goDrone's main functionality.

### cc Node

The ccNodes offers the following services and actions: 
- `takeOff`
Let's the drone take off
- `land`
Set's the drone control mode to recovery which forces the drone to return to home and land immediately 
- `turnTo`
Let's the drone turn to given heading
- `flyToLatLon`
Let's the drone fly to given latitude and longitude 
- `armMotors`
Arms the drones motors
- `changeAlt`
Changes the drones to given altitude
- `holdPos`
Lets the drone hold it's position and cancel all other operations

### GPS Node

The gps Node publishes live gps data with a 10 Hz frequency on the topic `gpsData` as a glob encode map.
The map consists of the following fields: `lat, lon, heading, alt, speed, sats`.

### Auto Pilot Node

The Auto pilot node offer one service which is `executeMission`, which requires a file path to a mission file in which a script like plan of cc Node services/ action describes the mission and as such enables fully autonomous mission mode.

## User Interfaces

### Web Interface

The web interface offers a simple interface with which the drone can be controlled. To launch the interface, use `go run /userInterface/webInterface/webInterface.go`

### Console 

The drone can bw controlled via. a command line tool. To use the console tool, use `go run /userInterface/console.go`.

#### Commands: 

- `reconnect`
Reconnects the go program to the goDrone core nodes

- `takeoff`
Starts the takeoff service which runs the arming and taking off sequence

- `land`
Calls the landing action which set the drone to recovery mode

- `markhomepos`
Calls the markhomepos action which stores the current position as home pos

- `turnto`
Calls the turnto service which let's the drone turn to given direction

- `flytolatlon`
Calls the flytolatlon service which let's the drone fly to given lat lon

- `listtopics`
Lists active topics on connected nodes

- `setneutral`
Sets thew stick positions to neutral (no thrust)

- `setstate`
Sets given state to given value (only usefull if operating non STATELESS)

- `getstates` 
Lists all states (also connection states)

- `getgps`
Prints live GPS data

- `endcom` 
Ends communication with drone and quits console