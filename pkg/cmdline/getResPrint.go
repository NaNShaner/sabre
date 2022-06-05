package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/sabstruct/res"
	"sabre/pkg/util/commontools"
	"sabre/pkg/util/hostregister"
)

//### 例如：获取erp系统的下的demo工程部署在哪些机器的Tomcat中
//key :/mid/ERP/Tomcat/{projectName}/{hostname/ipaddr}
//sabrectl get mid -t tomcat -a app -n erp

//输出：
//namespace	host		midType	projectName	port	version	monitor running	runningTime
//MNPP		127.0.0.1 	Tomcat 	demo		8099	7.0.78 	True	True	10d
//MNPP		127.0.0.2 	Tomcat 	demo		8099	7.0.78 	True	True	10d

var cmdGetResPrint = &cobra.Command{
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
			//json, err := yamlfmt.PrintResultJson(v)
			//if err != nil {
			//	return
			//}
			//fmt.Printf("%s\n", json)
			reqResp, setHttpReqErr := hostregister.SetHttpReq(kName, valueName)
			if setHttpReqErr != nil {
				fmt.Printf("请求sabrelet 失败,%s\n", setHttpReqErr)
				os.Exit(-1)
			}
			fmt.Printf("%s\n", reqResp)

		}
	},
	// 命令执行前进行判断，类似django的post_save
	PersistentPreRunE: func(cmdline *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("参数数量不正确，仅需要请输入yaml文件名称即可\n")
		}
		return nil
	},
}

func init() {
	cmdGetResPrint.Flags().StringVarP(&namespace, "namespace", "n", "", "主机所属系统简称")
	cmdGetResPrint.Flags().StringVarP(&area, "area", "a", "", "主机所属网络安全域")
	namespaceErr := cmdGetResPrint.MarkFlagRequired("namespace")
	if namespaceErr != nil {
		fmt.Printf("%s，请输入-n 或者--namespace 输入", namespaceErr)
	}
	areaErr := cmdHostRegister.MarkFlagRequired("area")
	if areaErr != nil {
		fmt.Printf("%s，请输入-a 或者--area 输入", areaErr)
	}
	AddCommand(cmdGetResPrint)
}
