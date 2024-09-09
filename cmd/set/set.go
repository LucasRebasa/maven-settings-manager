package set

import (
	"fmt"
	"os"
	"slices"

	"github.com/msm/constants"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SetCmd = &cobra.Command{
	Use:     "set [application]",
	Short:   "Set the settings.xml to the application that you choose",
	Aliases: []string{"s"},
	Args:    cobra.MatchAll(cobra.ExactArgs(1)),
	Run:     Run,
}

func Run(cmd *cobra.Command, args []string) {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	apps := utils.GetAppsFromMap(appsMap)
	if len(args) == 0 || !slices.Contains(apps, args[0]) {
		fmt.Println("Invalid app")
		utils.PrintAppsListFromMap(appsMap)
		return
	}
	setAppName(args[0])
	fmt.Printf("App set to: %v", args[0])
}

func setAppName(appName string) {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	newAppFile := appsMap[appName]
	currentAppName := viper.GetString(constants.CURRENT_APP)
	currentAppFile := appsMap[currentAppName]
	homeDir := viper.GetString(constants.HOME_SETTINGS_DIR)
	err := os.Rename(homeDir+"/settings.xml", currentAppFile)
	if err != nil {
		panic(err)
	}
	err = os.Rename(newAppFile, homeDir+"/settings.xml")
	if err != nil {
		panic(err)
	}
	viper.Set(constants.CURRENT_APP, appName)
	viper.WriteConfig()
}

func getValidArgs() []string {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	apps := make([]string, 0, len(appsMap))
	for k := range appsMap {
		apps = append(apps, k)
	}
	return apps
}
