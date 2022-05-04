package main

import (
	"awesomeProject/pkg/cmdline"
	"fmt"
	"github.com/spf13/cobra"
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

	var cmdEcho = &cobra.Command{
		Use:     "create [string to echo]",
		Short:   "Echo anything to the screen",
		Long:    `echo is for echoing anything back. Echo works a lot like print, except it has a child command.`,
		Args:    cobra.MinimumNArgs(1),
		Example: "sabrectl create www.baidu.com",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s\n", args)
		},
		// 在执行Run之前，如果错误则无法继续执行，直接报错。类似django的handler的post_save
		//PersistentPreRunE: func(cmdline *cobra.Command, args []string) error {
		//	flagName := cmdline.Flags().Name()
		//	if flagName != "tomcat" {
		//		return fmt.Errorf("must specify one of -f and -k")
		//	}
		//	return nil
		//},

	}

	var cmdTimes = &cobra.Command{
		Use:   "tomcat [deploy tomcat middleware]",
		Short: "deploy tomcat middleware",
		Long:  `download pkg from server, and deploy it.`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			//for i := 0; i < echoTimes; i++ {
			//	fmt.Println("Echo: " + strings.Join(args, " "))
			//}
			//fmt.Println("Echo: " + strings.Join(args, " "))

			tomcat, err := cmdline.DeployTomcat()
			if err != nil {
				return
			}
			fmt.Printf(tomcat)
		},
	}
	// var tomcat string
	// echoTimes变量可以作为 cmdTimes 函数作为入参
	// cmdTimes.Flags().IntVarP(&echoTimes, "tomcat", "f", 10, "times to echo the input")
	// cmdTimes.Flags().StringVarP(&tomcat, "tomcat", "q", "default", "times to echo the input")
	// cmdTimes.Flags().StringVar(&tomcat, "q", "default", "times to echo the input")

	var rootCmd = &cobra.Command{Use: "sabrectl"}
	// rootCmd.AddCommand 是根命令
	rootCmd.AddCommand(cmdEcho)
	// 如果在根命令上再执行AddCommand，则为该根命令的子命令
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()

}
