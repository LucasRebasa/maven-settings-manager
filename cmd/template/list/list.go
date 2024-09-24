package list

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/msm/constants"
	"github.com/msm/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your templates",
	Run:   run,
}

func init() {
	ListCmd.Flags().BoolP("description", "d", false, "Flag to print the description of the different params for each template")
}

func run(cmd *cobra.Command, args []string) {
	templates := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
	templatesDir := viper.GetString(constants.TEMPLATES_DIR)
	descFlag, _ := cmd.Flags().GetBool("description")
	fmt.Println("Templates directory: ", templatesDir)
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("There's no templates")
			return
		} else {
			panic(err)
		}
	}
	if len(entries) == 0 {
		fmt.Println("There's no templates")
		return
	}
	fmt.Println("Templates: ")
	for _, v := range entries {
		if !v.IsDir() {
			name := strings.Split(v.Name(), ".")[0]
			if !descFlag {
				fmt.Println("\t", name)
			} else {
				file, err := os.Open(templatesDir + "/" + name + ".xml")
				if err != nil {
					fmt.Println("Could not load params for template")
				} else {
					params, description := utils.GetParamsFromTemplate(file)
					for i := range params {
						fmt.Printf("\t%v --> Param: %v, Description: %v \n", name, params[i], description[i])
					}
				}
			}
			if !slices.Contains(maps.Keys(templates), name) {
				fmt.Println("New template found! -> ", v.Name())
				setParamsForTemplate(templates, v.Name(), name)
			}
		}
	}

}

func setParamsForTemplate(templates map[string][]string, fileName string, name string) {
	templateDir := viper.GetString(constants.TEMPLATES_DIR)
	file, err := os.Open(templateDir + "/" + fileName)
	if err != nil {
		panic(err)
	}

	params, descriptions := utils.GetParamsFromTemplate(file)
	templates[name] = params
	fmt.Println("With params: ")
	for i := range params {
		fmt.Printf("\tParam: %v, Description: %v \n", params[i], descriptions[i])
	}
	viper.Set(constants.TEMPLATES_MAP, templates)
	viper.WriteConfig()
}
