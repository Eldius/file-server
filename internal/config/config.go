package config

import (
	"github.com/spf13/viper"
	"os"
	"text/template"
)

var (
	BuildDate  string
	Version    string
	CommitDate string
	Commit     string
)

type versionInfo struct {
	BuildDate  string
	Version    string
	CommitDate string
	Commit     string
}

func GetDebugModeEnabled() bool {
	return viper.GetBool(LogDebugLevelKey)
}

func GetServerPort() int {
	return viper.GetInt(ServerPortKey)
}

func GetLogFormat() string {
	return viper.GetString(LogFormatKey)
}

func VersionInfo() {
	template.Must(template.New("version").Parse(`---
version:     {{.Version}}
commit:      {{.Commit}}
commit date: {{.CommitDate}}
build date:  {{.BuildDate}}
`)).Execute(os.Stdout, versionInfo{
		BuildDate:  BuildDate,
		Version:    Version,
		CommitDate: CommitDate,
		Commit:     Commit,
	})
}
