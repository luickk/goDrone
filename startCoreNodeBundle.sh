sleep 1
go run /home/pi/go/src/goDrone/gpsNode/gpsNode.go &
P1=$!

sleep 1
go run /home/pi/go/src/goDrone/ccNode/ccNode.go &
P2=$!

sleep 1
go run /home/pi/go/src/goDrone/apNode/apNode.go &
P3=$!

wait $P1 $P2 $P3
