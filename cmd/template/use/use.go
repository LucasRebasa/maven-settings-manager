package use

import (
	"encoding/xml"
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

var UseCmd = &cobra.Command{
	Use:   "use [template]",
	//Args:  cobra.ExactArgs(1),
	Short: "Use the template passing the corresponding values",
	Run:   run,
	DisableFlagParsing: true,
}


func run(cmd *cobra.Command, args []string) {
	parseCustomFlags(cmd,args,args[0])
	fmt.Println(cmd.Flags().GetString("param"))
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
			templates = append(templates, v.Name())
		}
	}
	if !slices.Contains(templates, args[0]) {
		fmt.Printf(`"%v template does not exist"`, args[0])
		os.Exit(1)
	}

}

func readTemplate(xmlFileName string) {
	//TODO Crear un map[string]string que asocie templates a sus parametros

	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		panic(err)
	}
	
	decoder := xml.NewDecoder(xmlFile)
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
			for _,v := range findings {
				strings.Split(strings.Split(v, "%")[0], ":")
			}
			fmt.Println(findings)
		default:
		}
	}
}

func parseCustomFlags(cmd *cobra.Command, args []string,templateName string) {
	templateMap := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
	params := templateMap[templateName] 
	for _,v := range params {
		cmd.Flags().String(v, "", "Flag for setting the param \"" + v + "\"")
	}
	cmd.DisableFlagParsing = false
	err := cmd.ParseFlags(args)
	if err != nil {
		fmt.Println("Wrong parameters.")
		fmt.Println("Usage:")
		for _,v := range params {
			fmt.Println("--",v,"value")
		}
	}
	if cmd.Flag("help").Value.String() == "true" {
		cmd.Help()
		fmt.Println()
		os.Exit(0)
	}
}