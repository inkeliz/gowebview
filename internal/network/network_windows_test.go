package network

import (
	"testing"
)

func TestEnablePrivateConnections(t *testing.T) {
	if err := EnablePrivateConnections(); err != nil {
		t.Error(err)
	}

	if IsAllowedPrivateConnections() == false {
		t.Error("not enabled")
	}
}
