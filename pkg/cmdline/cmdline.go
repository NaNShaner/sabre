// Package cmdline
/*
命令行的设置
*/
package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sabrectl",
	Short: "sabrectl controls the sabrectl cluster manager",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmdline *cobra.Command, args []string) {
		fmt.Println("OK")
	},
}

// AddCommand 封装了两个函数
func AddCommand(cmdline *cobra.Command) {
	rootCmd.AddCommand(cmdline)
}

func Execute() error {
	return rootCmd.Execute()
}
