package use

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
	err := utils.ParseCustomFlagsFromTemplate(cmd, args, args[0], true)
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
	utils.UseTemplate(args[0], cmd, false)
}

