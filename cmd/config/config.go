package config

import (
	"fmt"

	myinit "github.com/msm/cmd/config/init"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:     "config",
	Short:   "Initialization and Configuration of the tool",
	Aliases: []string{"c"},
	Run:     Run,
}

func init() {
	ConfigCmd.AddCommand(myinit.NewConfigInitCmd())
}

func Run(cmd *cobra.Command, args []string) {
	fmt.Println("config command")
}
