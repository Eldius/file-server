package cmd

import (
	"fmt"
	"github.com/eldius/file-server/internal/config"
	"github.com/eldius/file-server/internal/logger"
	"github.com/eldius/file-server/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the file server",
	Long:  `Starts the file server.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := config.GetServerPort()
		bp := "."
		if len(args) > 0 {
			bp = strings.Join(args, " ")
		}
		if err := server.Start(p, bp); err != nil {
			err = fmt.Errorf("starting server: %w", err)
			logger.GetLogger().With("error", err).Error("failed to start server")
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntP("port", "p", 8080, "Port to listen to")
	if err := viper.BindPFlag(config.ServerPortKey, startCmd.Flags().Lookup("port")); err != nil {
		err = fmt.Errorf("binding server port flag to viper property: %w", err)
		panic(err)
	}
}
