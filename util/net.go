package util

import (
	"github.com/v2fly/v2ray-core/v5/common/net"
	"regexp"
)

func GetNoUsePort(portRange *net.PortRange) []*net.PortRange {
	portRanges := make([]*net.PortRange, 0)
	pr := &net.PortRange{From: 0, To: 0}
	for i := portRange.FromPort(); i <= portRange.ToPort(); i++ {
		if CheckPort(int(i)) {
			if pr.From == 0 {
				pr.From = uint32(i)
			}
			if i == portRange.ToPort() {
				pr.To = uint32(i)
				portRanges = append(portRanges, pr)
			}
		} else {
			if pr.From != 0 {
				pr.To = uint32(i - 1)
				portRanges = append(portRanges, pr)
				pr = &net.PortRange{From: 0, To: 0}
			}
		}
	}
	return portRanges
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
