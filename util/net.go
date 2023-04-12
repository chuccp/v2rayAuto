package util

import "github.com/v2fly/v2ray-core/v5/common/net"

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
