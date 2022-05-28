package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	Ji "sabre/pkg/util/jdk/install"
	"sabre/pkg/yamlfmt"
)

var yamlfile string

var cmdDeployJdk = &cobra.Command{
	Use:   "jdk [deploy jdk]",
	Short: "deploy jdk , config env",
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

			// 执行安装操作
			jdkInstall := Ji.JdkInstall{}
			jdkInstallPathErr := jdkInstall.JdkInstallForApp((*commontools.Basest)(yamlFmt))
			if jdkInstallPathErr != nil {
				fmt.Printf("%s\n", err)
				os.Exit(-1)
			}
			fmt.Printf("jdk deployed to %s\n", jdkInstallPathErr)
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
	cmdDeployJdk.Flags().StringVarP(&yamlfile, "action", "a", "", "控制jdk安装动作")
	cmdCreate.AddCommand(cmdDeployJdk)
}
