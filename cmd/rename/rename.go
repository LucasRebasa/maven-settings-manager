package rename

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/msm/constants"
	"github.com/msm/structure"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RenameCmd = &cobra.Command{
	Use:   "rename [current] [new]",
	Short: "Rename your application",
	Args:  cobra.ExactArgs(2),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	apps := utils.GetAppsFromMap(appsMap)
	if !slices.Contains(apps, args[0]) {
		fmt.Printf("App %v not found. \n", args[0])
		utils.PrintAppsListFromMap(appsMap)
		return
	}
	if slices.Contains(apps, args[1]) {
		fmt.Printf("App %v already exists. \n", args[1])
		utils.PrintAppsListFromMap(appsMap)
		return
	}
	renameFiles(appsMap, args[0], args[1])
	isCompleteMode := viper.GetBool(constants.COMPLETE_MODE)
	if isCompleteMode {
		renameAppTagXml(args[1])
	}
	fmt.Printf(`App "%v" renamed to "%v"`, args[0], args[1])
}

func renameFiles(appsMap map[string]string, oldName, newName string) {
	oldAppFileName := appsMap[oldName]
	splittedOldAppFileName := strings.Split(oldAppFileName, "/")
	currentApp := viper.GetString(constants.CURRENT_APP)
	var newAppFileName string
	if currentApp == oldName {
		viper.Set(constants.CURRENT_APP, newName)
		newAppFileName = oldAppFileName
	} else {
		newAppFileName = strings.Join(splittedOldAppFileName[0:len(splittedOldAppFileName)-1], "/") + "/settings - " + newName + ".xml"
		err := os.Rename(oldAppFileName, newAppFileName)
		if err != nil {
			panic(err)
		}
	}
	delete(appsMap, oldName)
	appsMap[newName] = newAppFileName
	viper.Set(constants.MAP_FILE_APPNAME, appsMap)
	viper.WriteConfig()
}

func renameAppTagXml(newName string) {
	appsMap := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	fileName := appsMap[newName]
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	settings := structure.Settings{}
	err = xml.Unmarshal(data, &settings)
	if err != nil {
		panic(err)
	}
	appTag := fmt.Sprintf(`%%APP:%v%%`, newName)
	settings.Comment = appTag
	data,err = xml.MarshalIndent(settings, "", "\t")
	if err != nil {
		panic(err)
	}
	fixedData := utils.FixXMLData(data)
	err = os.WriteFile(fileName, fixedData, os.ModeAppend)
	if err != nil {
		panic(err)
	}
}