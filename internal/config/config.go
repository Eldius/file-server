package config

import (
	"github.com/spf13/viper"
)

func GetDebugModeEnabled() bool {
	return viper.GetBool(LogDebugLevelKey)
}

func GetServerPort() int {
	return viper.GetInt(ServerPortKey)
}

func GetLogFormat() string {
	return viper.GetString(LogFormatKey)
}
