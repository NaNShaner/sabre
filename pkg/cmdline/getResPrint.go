package cmdline

import (
	"fmt"
	"github.com/spf13/cobra"
	"sabre/pkg/util/getSomethingToPrint"
)

var (
	ns    string
	rType string
	rKind string
)

var cmdGetResPrint = &cobra.Command{
	Use:   "get [resource information]",
	Short: "Obtain resource information from the platform for presentation",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		getInfoFromCmdlineErr := getSomethingToPrint.PrintFmt(rType, ns, rKind)
		if getInfoFromCmdlineErr != nil {
			fmt.Printf("%s\n", getInfoFromCmdlineErr)
			return
		}
	},
}

func init() {
	cmdGetResPrint.Flags().StringVarP(&ns, "namespace", "n", "", "主机所属系统简称")
	cmdGetResPrint.Flags().StringVarP(&rType, "rsType", "t", "", "资源类型")
	cmdGetResPrint.Flags().StringVarP(&rKind, "rsKind", "k", "", "资源种类")
	namespaceErr := cmdGetResPrint.MarkFlagRequired("namespace")
	if namespaceErr != nil {
		fmt.Printf("%s，请输入-n 或者--namespace 输入", namespaceErr)
	}

	rsTypeErr := cmdGetResPrint.MarkFlagRequired("rsType")
	if rsTypeErr != nil {
		fmt.Printf("%s，请输入-n 或者--namespace 输入", rsTypeErr)
	}

	rsKindErr := cmdGetResPrint.MarkFlagRequired("rsKind")
	if rsKindErr != nil {
		fmt.Printf("%s，请输入-n 或者--namespace 输入", rsKindErr)
	}
	AddCommand(cmdGetResPrint)
}
