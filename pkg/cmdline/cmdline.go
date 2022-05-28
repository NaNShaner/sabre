// Package cmdline
/*
命令行的设置
*/
package cmdline

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sabrectl",
	Short: "sabrectl controls the sabrectl cluster manager",
}

// AddCommand 封装了两个函数
func AddCommand(cmdline *cobra.Command) {
	rootCmd.AddCommand(cmdline)
}

func Execute() error {
	return rootCmd.Execute()
}
