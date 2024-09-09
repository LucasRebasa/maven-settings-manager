package list

import (
	"github.com/msm/constants"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the applications",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	utils.PrintAppsListFromMap(appsMap)
}
