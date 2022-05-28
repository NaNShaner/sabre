package res

import (
	"fmt"
	"net"
	"sabre/pkg/util/commontools"
	"testing"
)

func HRegister(h HostRegister) error {
	host := h.ReturnHost()
	ipFmt, err := h.IpFmt(host.IPAddr.String())
	if err != nil {
		return fmt.Errorf("输入的%s地址不是正确的IP地址格式，格式例如192.168.1.1\n", host.IPAddr.String())
	}

	// 判断地址是为当前服务器的地址
	hostInfo, GetHostIPErr := h.GetOsInfo()
	if GetHostIPErr != nil {
		return fmt.Errorf("输入的%s地址不是当前服务的地址，请确认\n", host.IPAddr.String())
	}
	fmt.Printf("hostInfo: %+v\n", hostInfo)

	isLocalHostIP, isLocalHostIPErr := commontools.GetHostIP(ipFmt)
	if isLocalHostIPErr != nil {
		return fmt.Errorf("%s\n", isLocalHostIPErr)
	}
	fmt.Printf("是否当前主机: %t\n", isLocalHostIP)
	ping, PingErr := h.Ping(ipFmt)
	if PingErr != nil {
		return fmt.Errorf("当前主机ping状态, %t\n", PingErr)
	}
	fmt.Printf("当前主机ping状态: %t\n", ping)

	getMem, getMemErr := h.GetMem()
	if getMemErr != nil {
		return fmt.Errorf("获取当前主机的内存信息失败, %s\n", getMemErr)
	}
	fmt.Printf("当前主机mem状态: %+v\n", getMem)

	return nil

}

func TestRegister(t *testing.T) {
	ip := net.ParseIP("192.168.3.152")
	H := Hosts{
		IPAddr: ip,
	}

	err := HRegister(&H)
	if err != nil {
		t.Error(err.Error())
	}

	t.Log("done\n")
}

func TestHosts_Ping(t *testing.T) {
	var h Hosts
	ping, err := h.Ping(net.ParseIP("180.101.49.12"))
	if err != nil {
		t.Error(err)
	}
	t.Log(ping)
}
