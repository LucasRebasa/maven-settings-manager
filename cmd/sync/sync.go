package sync

import (
	"fmt"
	"os"
	"strings"

	"github.com/msm/constants"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Updates and syncs the templates and apps. Use it when adding new templates or apps manually",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	syncAppsSettings()
	syncTemplates()
	viper.WriteConfig()
}

func syncTemplates() {
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		panic(err)
	}
	templatesList := make(map[string][]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			templateName := strings.Split(fileName, ".")[0]
			file, err := os.Open(templatesDir+"/"+fileName)
			if err != nil {
				panic(err)
			}
			params, _ := utils.GetParamsFromTemplate(file)
			templatesList[templateName] = params
			fmt.Println("Found template -->", templatesDir+"/"+fileName)
		}
	}
	viper.Set(constants.TEMPLATES_MAP, templatesList)
}

func syncAppsSettings() {
	settingsDir := viper.GetString(constants.HOME_SETTINGS_DIR)
	entries, err := os.ReadDir(settingsDir)
	if err != nil {
		panic(err)
	}
	settingsList := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			splittedFileName := strings.Split(fileName, "-")
			if strings.TrimSpace(splittedFileName[0]) == "settings" {
				appName := strings.Split(splittedFileName[1], ".")[0]
				settingsList[appName] = settingsDir + "/" + fileName
				fmt.Println("Found setting -->", settingsDir + "/" + fileName)
			}
		}
	}
	viper.Set(constants.MAP_FILE_APPNAME, settingsList)
}
