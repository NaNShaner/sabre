package res

import "net"

type Hosts struct {
	HostName     string
	IPAddr       net.IPAddr
	Online       bool
	OnlineStatus []string
	BelongTo     string
}

type HostRegister interface {
	ServerRegister(*Hosts) (string, error)
}

func (h *Hosts) ServerRegister() (string, error) {
	return "", nil
}
