package template

import (
	"fmt"

	"github.com/msm/cmd/template/use"
	"github.com/spf13/cobra"
)

var TemplateCmd = &cobra.Command{
	Use: "template [option]",
	Short: "List and use templates",
	Run: run,
}

func init() {
	TemplateCmd.AddCommand(use.UseCmd)
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(`Use "msm template list" or "msm template use [name]" to execute this command`)
}