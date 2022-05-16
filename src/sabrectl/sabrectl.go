package main

import (
	//"sabre/pkg/cmdline"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sabre/pkg/sabstruct"
	"sabre/pkg/util/commontools"
	Ti "sabre/pkg/util/tomcat/install"
	"sabre/pkg/yamlfmt"
)

func main() {

	var cmdCreate = &cobra.Command{
		Use:     "create [install or deploy middleware resources]",
		Short:   "Install anything to the host",
		Long:    `create  has a child command.`,
		Args:    cobra.MinimumNArgs(1),
		Example: "sabrectl create tomcat [path to yaml file]",
	}

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

				// 执行安装操作
				_, err = Ti.TomcatInstall((*commontools.Basest)(yamlFmt))
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
				tomcatInstallPath, tomcatInstallPathErr := Ti.TomcatInstall((*commontools.Basest)(yamlFmt))
				if tomcatInstallPathErr != nil {
					fmt.Printf("%s\n", err)
					os.Exit(-1)
				}
				fmt.Printf("tomcat deployed to %s\n", tomcatInstallPath)
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

	var yamlfile string
	// echoTimes变量可以作为 cmdTimes 函数作为入参
	// cmdTimes.Flags().IntVarP(&echoTimes, "tomcat", "f", 10, "times to echo the input")
	cmdDeployJdk.Flags().StringVarP(&yamlfile, "action", "a", "", "控制jdk安装动作")
	// cmdTimes.Flags().StringVar(&tomcat, "q", "default", "times to echo the input")

	var rootCmd = &cobra.Command{Use: "sabrectl"}
	// rootCmd.AddCommand 是根命令
	rootCmd.AddCommand(cmdCreate)
	// 全局命令的flags
	// rootCmd.PersistentFlags().StringVar(&yamlfile, "yaml", "", "config file (default is $HOME/.cobra_exp1.yaml)")
	// 如果在根命令上再执行AddCommand，则为该根命令的子命令
	cmdCreate.AddCommand(cmdDeployTomcat)
	cmdCreate.AddCommand(cmdDeployJdk)

	rootCmd.Execute()

}
