// Package cmd has the CLI commands implementations
package cmd

import (
	"fmt"
	"github.com/eldius/file-server/internal/config"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file-server",
	Short: "A simple tool to share files in a network",
	Long:  `A simple tool to share files in a network.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.Setup(cfgFile)
		return nil
	},
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.file-server/config.yaml)")

	rootCmd.PersistentFlags().Bool("debug", false, "Runs with debug log enabled")
	if err := viper.BindPFlag(config.LogDebugLevelKey, rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		err = fmt.Errorf("binding debug mode flag to viper property: %w", err)
		panic(err)
	}

	rootCmd.PersistentFlags().String("log-format", "text", "Output log format")
	if err := viper.BindPFlag(config.LogFormatKey, rootCmd.PersistentFlags().Lookup("log-format")); err != nil {
		err = fmt.Errorf("binding log format flag to viper property: %w", err)
		panic(err)
	}
}
