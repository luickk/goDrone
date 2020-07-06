package utils

import (
	"encoding/binary"
	"math"
	"unsafe"
)

func EncodeLatLonAlt(lat float64, lon float64, alt float64) []byte {
	return append(append(Float64bytes(lat), Float64bytes(lon)...), IntToByteArray(int64(alt))...)
}

func DecodeLatLonAlt(encodedData []byte) (float64, float64, int) {
	lat, lon, alt := 0.0, 0.0, 0
	// parsing lat lon from service params
	// params{lat float64, lon float64}
	if len(encodedData) == 24 {
		lat = Float64frombytes(encodedData[0:8])
		lon = Float64frombytes(encodedData[8:16])
		alt = int(ByteArrayToInt(encodedData[16:24]))
	}
	return lat, lon, alt
}

//by https://gist.github.com/ecoshub/5be18dc63ac64f3792693bb94f00662f
func IntToByteArray(num int64) []byte {
	size := int(unsafe.Sizeof(num))
	arr := make([]byte, size)
	for i := 0; i < size; i++ {
		byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
		arr[i] = byt
	}
	return arr
}

func ByteArrayToInt(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0; i < size; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

// by https://stackoverflow.com/questions/22491876/convert-byte-slice-uint8-to-float64-in-golang/22492518#22492518
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
