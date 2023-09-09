package netutil

import "testing"

const MyIPV4 = "192.168.1.250"

func TestGetLocalIPV4(t *testing.T) {
	if GetLocalIPV4(false) != MyIPV4 {
		t.Error("failed GetLocalIPV4 function test")
	}
}
