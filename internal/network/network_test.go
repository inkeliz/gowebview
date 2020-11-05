package network

import (
	"testing"
)

func TestIsPrivateNetworkString(t *testing.T) {
	for _, ip := range []string{"http://192.168.2.1", "192.168.2.1", "http://192.168.3.243", "192.168.3.243", "http://127.0.0.1", "127.0.0.1", "::1", "http://::1:2031"} {
		if IsPrivateNetworkString(ip) == false {
			t.Error("invalid verification of Private Network")
		}
	}
}
