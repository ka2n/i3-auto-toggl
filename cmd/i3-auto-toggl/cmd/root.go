/*
Copyright © 2023 ka2n
*/
package cmd

import (
	"os"

	i3autotoggl "github.com/ka2n/i3-auto-toggle"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "i3-auto-toggl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cflg, err := cmd.Flags().GetString("c")
		if err != nil {
			return err
		}
		return i3autotoggl.StartDaemonCmd(ctx, cflg)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.i3-auto-toggl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().String("c", "", "Configuration file, default is $XDG_CONFIG_HOME/i3-auto-toggl.yaml/i3-auto-toggl.yaml")
}
