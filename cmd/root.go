package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/msm/cmd/config"
	"github.com/msm/cmd/list"
	"github.com/msm/cmd/rename"
	"github.com/msm/cmd/set"
	"github.com/msm/cmd/sync"
	"github.com/msm/cmd/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "msm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//Checks if .msm directory exists
		currentCommand := strings.Split(cmd.CommandPath(), " ")
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		viper.SetConfigFile(home + "/.msm/config.json")

		err = viper.ReadInConfig()
		//If config file doesn't exist and the current command is not the "config init" command
		if err != nil && (!slices.Contains(currentCommand, "config") || (slices.Contains(currentCommand, "config") && !slices.Contains(currentCommand, "init"))) {
			fmt.Println(err)
			fmt.Println("You have to run the initialization command first \n", "Run \"msm config init\" to initialize your configuration")
			os.Exit(1)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(set.SetCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(rename.RenameCmd)
	rootCmd.AddCommand(template.TemplateCmd)
	rootCmd.AddCommand(sync.SyncCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.msm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
