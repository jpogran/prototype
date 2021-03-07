package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ConfigFileName is the name of the config file we accept
	ConfigFileName = ".prototype"
	// EnvironmentVariablePrefix is the prefix for env variables we accept
	EnvironmentVariablePrefix = "PROTO"
)

// PrototypeAppConfig is the main app config variable
type PrototypeAppConfig struct {
	LocalTemplatePath string
	TargetName        string
	TargetOutput      string
	ListTemplates     bool
	SelectedTemplate  string
	OutputFormat      string
	Author            string
	Summary           string
	License           string
	Source            string
}

var (
	// Config is the main app config variable
	Config = PrototypeAppConfig{
		LocalTemplatePath: "",
		TargetName:        "",
		TargetOutput:      "",
		SelectedTemplate:  "",
		ListTemplates:     false,
		OutputFormat:      "",
		Author:            "",
		Summary:           "",
		License:           "",
		Source:            "",
	}
)

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(ConfigFileName)
	v.SetConfigType("yml")
	v.SetConfigType("yaml")

	loadDefaultSettingsFor(v)

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	v.AddConfigPath(home)
	v.AddConfigPath(filepath.Join(home, ".config"))
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	} else {
		fmt.Printf("Config: %v :: %v\n", v.ConfigFileUsed(), v.AllSettings())
	}

	v.SetEnvPrefix(EnvironmentVariablePrefix)
	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to PROTO_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", EnvironmentVariablePrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			fmt.Printf("name: %v changed:%v set:%v\n", f.Name, !f.Changed, v.IsSet(f.Name))
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})

	err = v.Unmarshal(&Config)
	return nil
}

func loadDefaultSettingsFor(v *viper.Viper) error {
	home, _ := homedir.Dir()
	v.SetDefault("templates", filepath.Join(home, ".pdk", "templates"))
	// cwd, _ := os.Getwd()
	// v.SetDefault("name", filepath.Base(cwd))
	// v.SetDefault("output", cwd)
	v.SetDefault("format", "table")

	var currentUser string
	u, _ := user.Current()
	if strings.Contains(u.Username, "\\") {
		currentUser = strings.Split(u.Username, "\\")[1]
	} else {
		currentUser = u.Username
	}
	viper.SetDefault("author", currentUser)
	viper.SetDefault("license", "Apache2")
	viper.SetDefault("source", "")
	return nil
}

// osExit is a copy of `os.Exit` to ease the "exit status" test.
// See: https://stackoverflow.com/a/40801733/8367711
var osExit = os.Exit

// EchoStdErrIfError is an STDERR wrappter and returns 0(zero) or 1.
// It does nothing if the error is nil and returns 0.
func EchoStdErrIfError(err error) int {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		return 1
	}

	return 0
}

// Execute is the main function of `cmd` package.
func Execute() error {
	return rootCmd.Execute()
}
