package create

import (
	"fmt"
	"os"

	"github.com/msm/utils"
	"github.com/spf13/cobra"
)

var CreateCmd = &cobra.Command{
	Use: "create template newSettingsName",
	Run: run,
	Args: cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
}

func run(cmd *cobra.Command, args []string) {
	//TODO: Actualizar el archivo de configuracion cuando se crea o se usa un template 
	err := utils.ParseCustomFlagsFromTemplate(cmd,args, args[0], false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	utils.UseTemplate(args[0], cmd, true)
	fmt.Println("create template newSettingsName")
}
