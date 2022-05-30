package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/apiserver"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	Ti "sabre/pkg/util/tomcat/install"
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
			// TODO 便于调试后续删除
			printResultJson, err := yamlfmt.PrintResultJson((*commontools.Basest)(yamlFmt))
			if err != nil {
				return
			}
			fmt.Printf("%s\n", printResultJson)
			// 信息入库
			setInfoToDB, setInfoToDBErr := apiserver.HttpReq((*apiserver.Basest)(yamlFmt))
			if setInfoToDBErr != nil {
				fmt.Printf("%s\n", setInfoToDBErr)
				os.Exit(-1)
			}
			fmt.Printf("Tomcat information warehousing succeeded，%s\n", setInfoToDB)
			// 执行安装操作
			_, err = Ti.Deploy((*commontools.Basest)(yamlFmt))
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(-1)
			}
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