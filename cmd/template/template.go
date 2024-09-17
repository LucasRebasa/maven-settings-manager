package template

import (
	"fmt"

	"github.com/msm/cmd/template/create"
	"github.com/msm/cmd/template/list"
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
	TemplateCmd.AddCommand(list.ListCmd)
	TemplateCmd.AddCommand(create.CreateCmd)
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(`Use "msm template list" or "msm template use [name]" to execute this command`)
}