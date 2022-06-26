package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/apiserver"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/callsabrelet"
	"sabre/pkg/util/commontools"
	"sabre/pkg/yamlfmt"
)

var cmdDeployTomcat = &cobra.Command{
	Use:   "tomcat [deploy tomcat middleware]",
	Short: "deploy tomcat middleware",
	Long:  `download pkg from server, and deploy it.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		for _, f := range args {
			// 获取并解析命令行输入的yaml文件
			yamlFmt, getYamlFmtErr := yamlfmt.YamlFmt(f, sabstruct.Config{})
			if getYamlFmtErr != nil {
				fmt.Printf("%s\n", getYamlFmtErr)
				os.Exit(-1)
			}
			u := (*commontools.Basest)(yamlFmt)

			//判断yaml文件中期望部署的服务器是否属于当前namespace
			CheckInstallServerBelongToNSErr := u.CheckInstallServerBelongToNS()
			if CheckInstallServerBelongToNSErr != nil {
				fmt.Printf("%s\n", CheckInstallServerBelongToNSErr)
				return
			}
			// 信息入库
			setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(yamlFmt))
			if setInfoToDBErr != nil {
				fmt.Printf("%s\n", setInfoToDBErr)
				return
			}
			fmt.Printf("Tomcat information warehousing succeeded %s\n", setInfoToDB)

			b := (*callsabrelet.Basest)(yamlFmt)
			callsabrelet.CallFaceOfSabrelet(b, u.DeployHost)
			fmt.Printf("tomcat install done\n")
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
	cmdCreate.AddCommand(cmdDeployTomcat)
}
