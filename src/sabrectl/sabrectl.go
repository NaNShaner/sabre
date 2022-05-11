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
	//	解析命令行的如参数
	// sabrectl create tomcat -f tomcat.yaml
	// action := flag.String("create", "", "部署类型")
	//installPkgPath := flag.String("lp", installPkgPath, "安装文件")
	//installLocalPath := flag.String("lp", installLocalPath, "安装路径")
	//mdUser := flag.String("lp", mdUser, "启动用户")
	//mdUserGroup := flag.String("lp", mdUserGroup, "启动用户组")
	//flag.Parse()
	//fmt.Println()

	var cmdCreate = &cobra.Command{
		Use:   "create [string to echo]",
		Short: "Echo anything to the screen",
		Long:  `echo is for echoing anything back. Echo works a lot like print, except it has a child command.`,
		Args:  cobra.MinimumNArgs(1),
		//Example: "sabrectl create www.baidu.com",
		//Run: func(cmd *cobra.Command, args []string) {
		//	fmt.Printf("%s\n", args)
		//},
		// 在执行Run之前，如果错误则无法继续执行，直接报错。类似django的handler的post_save

	}

	var cmdDeployTomcat = &cobra.Command{
		Use:   "tomcat [deploy tomcat middleware]",
		Short: "deploy tomcat middleware",
		Long:  `download pkg from server, and deploy it.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// f := "/Users/bijingrui/sabre/pkg/getdeploypkg/tomcatInstll.yaml"
			for _, f := range args {
				yamlFmt, err := yamlfmt.YamlFmt(f, sabstruct.Config{})
				printResultJson, err := yamlfmt.PrintResultJson((*commontools.Basest)(yamlFmt))
				if err != nil {
					return
				}
				fmt.Printf("%s\n", printResultJson)
				if err != nil {
					return
				}
				_, err = Ti.TomcatInstall((*commontools.Basest)(yamlFmt))
				if err != nil {
					fmt.Printf("install fail %s\n", err)
					os.Exit(-1)
				}
				fmt.Printf("tomcat install done\n")
			}
		},
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
	cmdDeployTomcat.Flags().StringVarP(&yamlfile, "yaml", "f", "default", "所需yaml文件")
	// cmdTimes.Flags().StringVar(&tomcat, "q", "default", "times to echo the input")

	var rootCmd = &cobra.Command{Use: "sabrectl"}
	// rootCmd.AddCommand 是根命令
	rootCmd.AddCommand(cmdCreate)
	rootCmd.PersistentFlags().StringVar(&yamlfile, "yaml", "", "config file (default is $HOME/.cobra_exp1.yaml)")
	// 如果在根命令上再执行AddCommand，则为该根命令的子命令
	cmdCreate.AddCommand(cmdDeployTomcat)

	rootCmd.Execute()

}
