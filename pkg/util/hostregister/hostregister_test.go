package hostregister

import (
	"fmt"
	"os"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"testing"
)

func TestKeyName(t *testing.T) {
	// e.g. /hosts/erp/machine/app/hostname
	n := "erp"
	a := "app"
	h := "local"
	i := "192.168.3.58"

	t.Log(KeyName(n, h, a, i))

}

func TestAddHostToListSaveToDB(t *testing.T) {
	s := "/hosts/OICQ/machine/WEB/RuierMacBook-Pro.local/192.168.3.152"
	addHostToListSaveToDB, err := AddHostToListSaveToDB(s)
	if err != nil {
		t.Errorf(addHostToListSaveToDB)
		return
	}
	t.Log(addHostToListSaveToDB)
}

func TestSetHttpReq(t *testing.T) {
	var h res.Hosts
	hostName, getHostNameErr := commontools.GetLocalServerName()
	if getHostNameErr != nil {
		fmt.Println(getHostNameErr)
		os.Exit(-1)
	}
	kName := KeyName("OICQ", hostName, "WEB", "192.168.3.152")
	valueName, err := ValueName(&h, "192.168.3.152", "OICQ", "WEB")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(-1)
	}
	v := make(map[string]res.Hosts)
	v[kName] = valueName
	//注册主机信息
	registerHost, setHttpReqErr := SetHttpReq(kName, valueName)
	if setHttpReqErr != nil {
		fmt.Printf("Failed to request API server when registering host information, %s\n", setHttpReqErr)
		os.Exit(-1)
	}
	//注册主机列表，便于sabrelet查询
	HostToListSaveToDBKey, getHostToListSaveToDBKeyErr := AddHostToListSaveToDB(kName)
	if getHostToListSaveToDBKeyErr != nil {
		fmt.Printf("Get qurey host key err, %s", getHostToListSaveToDBKeyErr)
		os.Exit(-1)
	}
	_, hostToListSaveToDBErr := SetHttpReq(HostToListSaveToDBKey, kName)
	if hostToListSaveToDBErr != nil {
		fmt.Printf("Failed to request API server while registering host list, %s\n", hostToListSaveToDBErr)
		os.Exit(-1)
	}

	fmt.Printf("%s\n", registerHost)
}
