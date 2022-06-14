package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"sabre/pkg/util/hostregister"
)

var namespace string
var area string

var cmdHostRegister = &cobra.Command{
	Use:   "hosted [server ipaddr]",
	Short: "Register host to platform",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// 发起请求主机注册请求给到saberlet
		var h res.Hosts
		for _, f := range args {
			hostName, getHostNameErr := commontools.GetLocalServerName()
			if getHostNameErr != nil {
				fmt.Println(getHostNameErr)
				os.Exit(-1)
			}
			kName := hostregister.KeyName(namespace, hostName, area, f)
			valueName, err := hostregister.ValueName(&h, f, namespace, area)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(-1)
			}
			v := make(map[string]res.Hosts)
			v[kName] = valueName
			//注册主机信息
			registerHost, setHttpReqErr := hostregister.SetHttpReq(kName, valueName)
			if setHttpReqErr != nil {
				fmt.Printf("Failed to request API server when registering host information, %s\n", setHttpReqErr)
				os.Exit(-1)
			}
			//注册主机列表，便于sabrelet查询
			HostToListSaveToDBKey, getHostToListSaveToDBKeyErr := hostregister.AddHostToListSaveToDB(kName)
			if getHostToListSaveToDBKeyErr != nil {
				fmt.Printf("Get qurey host key err, %s", getHostToListSaveToDBKeyErr)
				os.Exit(-1)
			}

			if hostToListSaveToDBErr := hostregister.SetHostListInfoTODB(HostToListSaveToDBKey, kName); hostToListSaveToDBErr != nil {
				fmt.Printf("Failed to request API server while registering host list, %s\n", hostToListSaveToDBErr)
				os.Exit(-1)
			}
			fmt.Printf("%s\n", registerHost)

		}
	},
	// 命令执行前进行判断，类似django的post_save
	PersistentPreRunE: func(cmdline *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Register host please enter the IP address of the host\n")
		}
		return nil
	},
}

func init() {
	cmdHostRegister.Flags().StringVarP(&namespace, "namespace", "n", "", "主机所属系统简称")
	cmdHostRegister.Flags().StringVarP(&area, "area", "a", "", "主机所属网络安全域")
	namespaceErr := cmdHostRegister.MarkFlagRequired("namespace")
	if namespaceErr != nil {
		fmt.Printf("%s，请输入-n 或者--namespace 输入", namespaceErr)
	}
	areaErr := cmdHostRegister.MarkFlagRequired("area")
	if areaErr != nil {
		fmt.Printf("%s，请输入-a 或者--area 输入", areaErr)
	}
	AddCommand(cmdHostRegister)
}
