package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/msm/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// -- MAP UTILS -- //

func GetAppsFromMap(appsMap map[string]string) []string {
	apps := make([]string, 0, len(appsMap))
	for k := range appsMap {
		apps = append(apps, k)
	}
	return apps
}

func PrintAppsListFromMap(appsMap map[string]string) {
	fmt.Println("Choose an application to set from: ")
	for k := range appsMap {
		fmt.Println("\t" + k)
	}
}

// -- FILE UTILS -- //
func FixXMLData(data []byte) []byte {
	strSettingsData := string(data)
	strSettingsDataFixed := strings.ReplaceAll(strSettingsData[10:], `xmlns="http://maven.apache.org/SETTINGS/1.0.0"`, "")
	strSettingsData = strings.Replace(strSettingsData[:10]+strSettingsDataFixed, "settings", "settings "+constants.SCHEMA, 1)
	return []byte(strSettingsData)
}

func GetParamsFromTemplate(templateFile *os.File) ([]string, []string) {
	var params []string
	var descriptions []string
	regex := regexp.MustCompile(`%\s*[A-Za-z0-9 ]+\s*:\s*[A-Za-z0-9 ]+\s*%`)
	scanner := bufio.NewScanner(templateFile)
	for scanner.Scan() {
		currentText := scanner.Text()
		findings := regex.FindAllString(currentText, -1)
		for _, finding := range findings {
			paramKeyValue := strings.Split(strings.Split(finding, "%")[1], ":")
			if paramKeyValue[0] != "NAME" {
				params = append(params, paramKeyValue[0])
				descriptions = append(descriptions, paramKeyValue[1])
			}
		}
	}
	return params, descriptions
}

func CreateSettingsFile(data []string, appName string) {
	settingsDir := viper.GetString(constants.HOME_SETTINGS_DIR)
	fileName := settingsDir + "/settings - " + appName + ".xml"
	var tries int
	checkExistingFile(&fileName, settingsDir, appName, &tries)
	
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	fmt.Println("File created: ", fileName)
	file.Close()
}

func checkExistingFile(fileName *string, settingsDir, appName string, tries *int) {
	_,err := os.Stat(*fileName)
	if err == nil {
		*fileName = settingsDir + "/settings - " + appName + " - " + strconv.Itoa(*tries) + ".xml"
		fmt.Println("File Already exists, setting name to: ", *fileName)
		*tries++
		checkExistingFile(fileName, settingsDir, appName, tries)
	}
}

func ParseCustomFlagsFromTemplate(cmd *cobra.Command, args []string, templateName string, required bool) error {
	templateMap := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
	params := templateMap[templateName]
	for _, v := range params {
		cmd.Flags().String(v, "", "Flag for setting the param \""+v+"\"")
	}
	
	cmd.Flags().String("name", "", "Flag for setting the app or template name")
	
	params = append(params, "name")
	cmd.DisableFlagParsing = false
	err := cmd.ParseFlags(args)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Wrong parameters.")
		fmt.Println("Usage:")
		for _, v := range params {
			fmt.Printf("--%v %v", v, "value")
		}
		return errors.New("")
	}
	if cmd.Flag("help").Value.String() == "true" {
		cmd.Help()
		fmt.Println()
		os.Exit(0)
	}
	for _, param := range params {
		flagValue, err := cmd.Flags().GetString(param)
		if required && (err != nil || flagValue == "") {
			return fmt.Errorf("Flag \"%v\" missing", param)
		}
	}
	return nil
}

func UseTemplate(templateName string, cmd *cobra.Command, createTemplate bool) {
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	xmlFile, err := os.Open(templatesDir + "/" + templateName + ".xml")
	if err != nil {
		panic(err)
	}
	templateMap := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
	params := templateMap[templateName]
	var mapParamValue = make(map[string]string)
	for _, param := range params {
		value, _ := cmd.Flags().GetString(param)
		mapParamValue[param] = value
	}
	appName, _ := cmd.Flags().GetString("name")
	fmt.Println("NAME: ",appName)
	mapParamValue["name"] = appName
	data := getFileWithValuesSet(xmlFile, mapParamValue)
	xmlFile.Close()

	if createTemplate{
		createTemplateFile(data, appName)
	}else{
		CreateSettingsFile(data, appName)
	}
}

func getFileWithValuesSet(file *os.File, mapParamValue map[string]string) []string {
	regex := regexp.MustCompile(`%\s*[A-Za-z0-9 ]+\s*:\s*[A-Za-z0-9 ]+\s*%`)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentText := scanner.Text()
		findings := regex.FindAllString(currentText, -1)
		for _, finding := range findings {
			paramKeyValue := strings.Split(strings.Split(finding, "%")[1], ":")
			value, ok := mapParamValue[strings.ToLower(paramKeyValue[0])]
			if ok {
				fmt.Printf(`Value "%v" set to "%v"`, paramKeyValue[0], value)
				fmt.Println("")
				var modifiedLine string
				if paramKeyValue[0] == "NAME" {
					modifiedLine = strings.Replace(currentText, finding, fmt.Sprintf("NAME:%v", value), -1)
				} else{
					modifiedLine = strings.Replace(currentText, finding, value, -1)
				}
				currentText = modifiedLine
			}

		}
		lines = append(lines, currentText)
	}
	return lines
}

func createTemplateFile(data []string, appName string) {
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	fileName := templatesDir + "/" + appName + ".xml"
	var tries int
	checkExistingFile(&fileName, templatesDir, appName, &tries)
	
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(file)
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
	fmt.Println("File created: ", fileName)
	file.Close()
}

