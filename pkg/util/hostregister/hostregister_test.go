package hostregister

import (
	"net"
	"testing"
)

func TestHosts_Ping(t *testing.T) {
	var h Hosts
	ping, err := h.Ping(net.ParseIP("180.101.49.12"))
	if err != nil {
		t.Error(err)
	}
	t.Log(ping)
}
