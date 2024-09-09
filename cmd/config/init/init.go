package init

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/msm/constants"
	"github.com/msm/structure"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Use this command to initialize your configuration",
		Long:  "This command will create a folder called \".msm\" in your settings directory. You can define one or leave the default value",
		Run:   run,
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(`You executed the config init command, we are going to start the configuration of some parameters.`)
	fmt.Println("Default values will appear between parenthesis. i.e: Option (default value) : ...")
	fmt.Println("\nFirst we will set the path of your settings.xml files directory. It must be an ABSOLUTE path")
	setHomeDir()
	askCompleteMode()
	readSettingsFiles()
	askCurrentAppName()
	readHomeDir()
	fmt.Println("Settings found: ")
	fileNames := strings.Split(viper.GetString(constants.SETTINGS_FILES_NAMES), ";")
	for _, v := range fileNames {
		fmt.Println("\t", v)
	}
	createConfigFile()
}

func askCompleteMode() {
	fmt.Println("You will be able to run all the features if you configure the \"Complete Mode\"")
	fmt.Println("WARNING: the structure of your different settings.xml files will be modified and you won't be able to use comments in them")
	fmt.Print("Configure complete mode (y/n): ")
	completeMode := make([]byte, 10)
	inputLen, err := os.Stdin.Read([]byte(completeMode))
	if err != nil {
		panic(err)
	}
	completeMode = []byte(strings.TrimSpace(string(completeMode[:inputLen])))
	strCompleteMode := string(completeMode)

	var option, out bool
	for !out {
		switch strCompleteMode {
		case "y":
			option = true
			out = true
		case "n":
			option = false
			out = true
		default:
			fmt.Println(`Enter "y" or "n" option`)
			out = false

		}
	}

	viper.Set(constants.COMPLETE_MODE, option)
}

func readHomeDir() {
	fileNames := strings.Split(viper.GetString(constants.SETTINGS_FILES_NAMES), ";")
	homePath := viper.GetString(constants.HOME_SETTINGS_DIR)
	for _, v := range fileNames {
		fileName := homePath + "/" + v
		fileInfo, err := os.Stat(fileName)
		if err != nil || fileInfo.IsDir() {
			fmt.Printf("Could not modify %v\n", v)
		}
		if viper.GetBool(constants.COMPLETE_MODE) {
			setAppName(fileName, v)
		} else {
			populateMap(fileName, v, &v)
		}
	}
}

func populateMap(absoluteFileName, appFileName string, appName *string) {
	regex := regexp.MustCompile(`-([^.]+)\.`)
	strFound := regex.FindString(appFileName)

	if strFound != "" {
		*appName = strings.TrimSpace(strFound[1 : utf8.RuneCountInString(strFound)-1])
	} else {
		*appName = viper.GetString(constants.CURRENT_APP)
	}

	if strFound == "" {
		newAbsoluteName := strings.Split(absoluteFileName, ".xml")
		absoluteFileName = newAbsoluteName[0] + " - " + *appName + ".xml"
	}

	mapFileNameAppName := viper.GetStringMapString(constants.MAP_FILE_APPNAME)
	mapFileNameAppName[*appName] = absoluteFileName
	viper.Set(constants.MAP_FILE_APPNAME, mapFileNameAppName)
}

func setAppName(absoluteFileName, appFileName string) {
	data, err := os.ReadFile(absoluteFileName)
	if err != nil {
		panic(err)
	}
	settings := &structure.Settings{}
	err = xml.Unmarshal(data, settings)
	if err != nil {
		panic(err)
	}

	var appName string
	populateMap(absoluteFileName, appFileName, &appName)

	settings.Comment = "%APP:" + appName + "%"
	settingsData, _ := xml.MarshalIndent(settings, "", "\t")
	fixedSettingsData := utils.FixXMLData(settingsData)
	err = os.WriteFile(absoluteFileName, fixedSettingsData, os.ModeAppend)
	if err != nil {
		panic(err)
	}

}

func askCurrentAppName() {
	if viper.GetBool(constants.COMPLETE_MODE) {
		fmt.Println("\nThe name of each application is stored in a tag \"app\" inside \"settings\" in your settings.xml.")
		fmt.Println("We will set the tag in every settings file under the Home Directory that you just set.")
	} else {
		fmt.Println(`The command isn't running in "Complete Mode", your settings.xml files won't be modified `)
	}
	fmt.Println("By default the pattern for the settings files is \"settings-APPNAME.xml\", make sure that your current settings files follow this pattern")
	fmt.Print("\nSet a name for your current application: ")
	currentApp := make([]byte, 100)
	inputLen, err := os.Stdin.Read([]byte(currentApp))
	if err != nil {
		panic(err)
	}
	currentApp = []byte(strings.TrimSpace(string(currentApp[:inputLen])))
	viper.Set(constants.CURRENT_APP, string(currentApp))
}

func setHomeDir() {
	var home string
	for {
		home, _ = os.UserHomeDir()
		home = strings.Replace(home, "\\", "/", -1) + "/.m2"
		fmt.Printf("Settings path (%v) : ", home)
		userHome := make([]byte, 20)
		inputLen, err := os.Stdin.Read(userHome)
		if err != nil {
			panic(err)
		}
		userHomeString := strings.TrimSpace(string(userHome[0:inputLen]))
		if userHomeString != "" {
			home = userHomeString
		}
		isValid, err := isValidPath(home)
		if isValid {
			break
		}
		if err != nil {
			fmt.Println("An error happened ", err.Error())
		}

		fmt.Println("Invalid path, enter a valid one")
	}

	viper.Set(constants.HOME_SETTINGS_DIR, home)

	fmt.Printf("Settings path set to: %v\n", home)
}

func isValidPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createConfigFile() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configPath := homePath+"/.msm" 
	err = os.Mkdir(configPath, os.ModeDevice)
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("Old configuration found, deleting and trying to create it again...")
			err := os.RemoveAll(configPath)
			if err != nil {
				panic(err)
			}
			createConfigFile()
			return
		}
		fmt.Println("An error ocurred creating the config file")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	configFilePath := configPath + "/config.json"
	_, err = os.Create(configFilePath)
	if err != nil {
		fmt.Println("An error ocurred creating the config file")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	viper.SetConfigFile(configFilePath)
	viper.WriteConfig()
}

func readSettingsFiles() {
	regex := regexp.MustCompile(`settings.*\.xml$`)

	homeDir := viper.GetString(constants.HOME_SETTINGS_DIR)
	entries, err := os.ReadDir(homeDir)
	if err != nil {
		panic(err)
	}
	fileNames := make([]string, 1)
	for _, v := range entries {
		if !v.IsDir() {
			fileName := strings.TrimSpace(v.Name())

			if regex.Match([]byte(fileName)) {
				fileNames = append(fileNames, fileName)
			}
		}
	}
	fileNamesJoined := strings.Join(fileNames, ";")[1:]
	viper.Set(constants.SETTINGS_FILES_NAMES, fileNamesJoined)
}
