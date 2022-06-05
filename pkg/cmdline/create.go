package cmdline

import (
	"github.com/spf13/cobra"
)

var cmdCreate = &cobra.Command{
	Use:     "create [install or deploy middleware resources]",
	Short:   "Install anything to the host",
	Long:    `create  has a child command.`,
	Args:    cobra.MinimumNArgs(1),
	Example: "sabrectl create tomcat [path to yaml file]",
}

func init() {
	AddCommand(cmdCreate)
}
