package utils

import (
	"encoding/xml"
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

func GetParamsFromTemplate(templateFile *os.File) []string {
	var params []string
	decoder := xml.NewDecoder(templateFile)
	regex := regexp.MustCompile(`%[A-Za-z0-9]+:[A-Za-z0-9]+%`)
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			var content string
			decoder.DecodeElement(&content, &se)
			findings := regex.FindAllString(content, -1)
			for _, v := range findings {
				//Split %key:value% in slice [key:value] and then extracts value into [key, value]
				paramKey := strings.Split(strings.Split(v, "%")[0], ":")[0]
				params = append(params, paramKey)
			}
			fmt.Println(findings)
		default:
		}
	}
	return params
}
