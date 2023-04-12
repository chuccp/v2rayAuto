package util

import (
	"github.com/v2fly/v2ray-core/v5/common/net"
	"testing"
)

func TestName(t *testing.T) {
	pr := &net.PortRange{From: 8088, To: 8089}
	portRanges := GetNoUsePort(pr)
	for _, v := range portRanges {
		t.Log(v.FromPort(), v.ToPort())
	}
}
