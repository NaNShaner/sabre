// Package cmdline
/*
命令行的设置
*/
package cmdline

import (
	"awesomeProject/pkg/getdeploypkg"
	"awesomeProject/pkg/sabstruct"
	"awesomeProject/pkg/util/aboutuser"
	"awesomeProject/pkg/yamlfmt"
	"fmt"
	"os/user"
)

//var rootCmd = &cobra.Command{
//	Use:   "sabrectl",
//	Short: "sabrectl controls the sabrectl cluster manager",
//	Long: `A longer description that spans multiple lines and likely contains
//examples and usage of using your application. For example:
//Cobra is a CLI library for Go that empowers applications.
//This application is a tool to generate the needed files
//to quickly create a Cobra application.`,
//	Run: func(cmdline *cobra.Command, args []string) {
//		fmt.Println("OK")
//	},
//}
//
//// AddCommand 封装了两个函数
//func AddCommand(cmdline *cobra.Command) {
//	rootCmd.AddCommand(cmdline)
//}
//
//func Execute() error {
//	return rootCmd.Execute()
//}

// GetYamlFmt  从命令行解析yaml文件
func GetYamlFmt() *getdeploypkg.Basest {
	file := "t.yaml"
	yamlFmt, err := yamlfmt.YamlFmt(file, sabstruct.Config{})
	if err != nil {
		return nil
	}
	return (*getdeploypkg.Basest)(yamlFmt)
}

// DeployTomcat 下载安装包
func DeployTomcat() (string, error) {
	// TODO： 将文件从输入结构体中读取
	// url := "https://dlcdn.apache.org/tomcat/tomcat-8/v8.5.78/bin/apache-tomcat-8.5.78.tar.gz"
	d := GetYamlFmt()
	// 需要下载文件的url
	//d.PkgFromUrl = url
	// 下载文件
	fmt.Printf("%+v", *d)
	pkgFromUrl, err := (*getdeploypkg.Basest).GetDeployPkgFromUrl(d)
	if err != nil {
		return "", err
	}

	u := user.User{
		Username: d.Spec.User.Name,
	}

	isUserExist, err := aboutuser.IsUserExist(u)
	if err != nil {
		return "", err
	}
	// 判断用户yaml文件中需求的用户在服务器本地是否存在
	if isUserExist.Username == d.Name {

	}

	return pkgFromUrl, nil
}
