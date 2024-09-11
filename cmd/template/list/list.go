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

func run(cmd *cobra.Command, args []string) {
	fmt.Println(`Templates: `)
	templates := viper.GetStringMapStringSlice(constants.TEMPLATES_MAP)
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
	for _, v := range entries {
		if !v.IsDir() {
			name := strings.Split(v.Name(), ".")[0]
			if !slices.Contains(maps.Keys(templates), name) {
				fmt.Println("New template found! -> ", v.Name())
				setParamsForTemplate(templates, v.Name(), name)
			}
			fmt.Println(name)
		}
	}
}

func setParamsForTemplate(templates map[string][]string, fileName string, name string) {
	templateDir := viper.GetString(constants.TEMPLATES_DIR)
	file, err := os.Open(templateDir+"/"+fileName)
	if err != nil {
		panic(err)
	}
	
	params := utils.GetParamsFromTemplate(file)
	templates[name] = params
}

