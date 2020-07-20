# Go Drone

A drone project in which the drone is controlled by a raspberry controlled DJI Naza V2 flight controller. The project builds up on three other projets off mine, the [robot communication interface](https://github.com/luickk/rcf), [goNazaV2Interface](https://github.com/luickk/goNazaV2Interface) and the [goDjiNazaGpsDecoder](https://github.com/luickk/goDjiNazaGpsDecoder).

The goDrone project aims to build an open source drone platform which can keep up with the latest drone technologies and as such supports the same functionality and reliability. Having that in mind, building a open source, self developed flight controller that can keep up with the newest tech respectively industry standrads is unrealistic. 

The alternative is to use existing hardware and make use of the work that's already been done by the big tech giants. To enable the use of flight controllers that don't offer access via. a computer I built the [go Naza interface](https://github.com/luickk/goNazaV2Interface) which when used with a Raspberry Pi emulates PWM signals, similar to those from a standard receiver.


## User Interfaces

### Console 

The drone can bw controlled via. a command line tool. It can be found in `./userInterface/console.go` and can be started with `go run userInterface/console.go`. <br>
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

Ends connection to goDrone core nodes