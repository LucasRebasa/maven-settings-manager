package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/msm/constants"
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
			if paramKeyValue[0] != "APP" {
				params = append(params, paramKeyValue[0])
				descriptions = append(descriptions, paramKeyValue[1])
			}
		}
	}
	return params, descriptions
}
