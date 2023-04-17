package util

import (
	"github.com/v2fly/v2ray-core/v5/common/net"
	"golang.org/x/exp/slices"
	"math/rand"
	"regexp"
)

func GetNoUsePort(from int, to int, num int, exPort []int) []int {
	ports := make([]int, 0)
	randValues := make([]int, 0)
	for i := from; i <= to; i++ {
		if slices.Contains(exPort, i) {
			continue
		}
		randValues = append(randValues, i)
	}
	for {
		ln := len(randValues)
		if ln > 0 {
			portIndex := rand.Intn(len(randValues))
			port := randValues[portIndex]
			if CheckPort(port) {
				ports = append(ports, port)
			}
			if len(ports) == num {
				break
			} else {
				randValues = append(randValues[:portIndex], randValues[portIndex+1:]...)
			}
		} else {
			break
		}

	}
	return ports
}

func CheckPort(port int) bool {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.AnyIP.IP(),
		Port: port,
	})
	if err == nil {
		listener.Close()
		return true
	}
	return false
}
func IsIP(ip string) bool {
	matchString, _ := regexp.MatchString("^((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:){6}((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^::([\\da-fA-F]{1,4}:){0,4}((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:):([\\da-fA-F]{1,4}:){0,3}((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:){2}:([\\da-fA-F]{1,4}:){0,2}((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:){3}:([\\da-fA-F]{1,4}:){0,1}((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:){4}:((25[0-5]|2[0-4]\\d|[01]?\\d\\d?)\\.){3}(25[0-5]|2[0-4]\\d|[01]?\\d\\d?)$|^([\\da-fA-F]{1,4}:){7}[\\da-fA-F]{1,4}$|^:((:[\\da-fA-F]{1,4}){1,6}|:)$|^[\\da-fA-F]{1,4}:((:[\\da-fA-F]{1,4}){1,5}|:)$|^([\\da-fA-F]{1,4}:){2}((:[\\da-fA-F]{1,4}){1,4}|:)$|^([\\da-fA-F]{1,4}:){3}((:[\\da-fA-F]{1,4}){1,3}|:)$|^([\\da-fA-F]{1,4}:){4}((:[\\da-fA-F]{1,4}){1,2}|:)$|^([\\da-fA-F]{1,4}:){5}:([\\da-fA-F]{1,4})?$|^([\\da-fA-F]{1,4}:){6}:$\n", ip)
	return matchString
}
