package use

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/msm/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const VIPER_APP_NAME = "app_name"

var UseCmd = &cobra.Command{
	Use: "use [template]",
	//Args:  cobra.ExactArgs(1),
	Short:              "Use the template passing the corresponding values",
	Run:                run,
	DisableFlagParsing: true,
}

func run(cmd *cobra.Command, args []string) {
 	err :=	parseCustomFlags(cmd, args, args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("There's no templates")
			return
		} else {
			panic(err)
		}
	}
	var templates []string

	for _, v := range entries {
		if !v.IsDir() {
			templates = append(templates, strings.Split(v.Name(), ".")[0])
		}
	}
	if !slices.Contains(templates, args[0]) {
		fmt.Printf(`"%v template does not exist"`, args[0])
		os.Exit(1)
	}
	useTemplate(args[0], cmd)
}

func useTemplate(templateName string, cmd *cobra.Command) {
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
	appName, _ := cmd.Flags().GetString("app")
	data := getFileWithValuesSet(xmlFile, mapParamValue)
	xmlFile.Close()

	createSettingsFile(data, appName)
}

func createSettingsFile(data []string, appName string) {
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	file, err := os.Create(templatesDir + "/settings - " + appName + ".xml")
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
	file.Close()
}

func getFileWithValuesSet(file *os.File, mapParamValue map[string]string) []string {
	regex := regexp.MustCompile(`%[A-Za-z0-9]+:[A-Za-z0-9]+%`)
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentText := scanner.Text()
		findings := regex.FindAllString(currentText, -1)
		for _, finding := range findings {
			paramKeyValue := strings.Split(strings.Split(finding, "%")[0], ":")
			value, ok := mapParamValue[paramKeyValue[0]]
			if ok {
				fmt.Printf(`Value "%v" set to "%v"`, paramKeyValue[0], value)
				modifiedLine := strings.Replace(currentText, finding, value, -1)
				currentText = modifiedLine
			}

		}
		lines = append(lines, currentText)
	}
	return lines
}

func parseCustomFlags(cmd *cobra.Command, args []string, templateName string) error {
	templateMap := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
	params := templateMap[templateName]
	for _, v := range params {
		cmd.Flags().String(v, "", "Flag for setting the param \""+v+"\"")
	}
	cmd.Flags().String("app", "", "Flag for setting the app name")
	params = append(params, "app")
	cmd.DisableFlagParsing = false
	err := cmd.ParseFlags(args)
	if err != nil {
		fmt.Println("Wrong parameters.")
		fmt.Println("Usage:")
		for _, v := range params {
			fmt.Println("--", v, "value")
		}
	}
	if cmd.Flag("help").Value.String() == "true" {
		cmd.Help()
		fmt.Println()
		os.Exit(0)
	}
	for _, param := range params {
		flagValue,err := cmd.Flags().GetString(param)
		if err != nil || flagValue == ""{
			return errors.New(fmt.Sprintf("Flag \"%v\" missing", param))
		}
	}
	return nil
}
