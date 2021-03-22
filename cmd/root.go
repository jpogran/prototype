package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	// ConfigFileName is the name of the config file we accept
	ConfigFileName = ".prototype"
	// EnvironmentVariablePrefix is the prefix for env variables we accept
	EnvironmentVariablePrefix = "PROTO"
)

var (
	cfgFile string
	Data    appconfiguration
)

type appconfiguration struct {
	LogLevel string

	LocalTemplateCache string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = createRootCommand()

func createRootCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:        "prototype",
		Aliases:    []string{},
		SuggestFor: []string{},
		Short:      "PDK prototype content template system",
		Long:       `PDK prototype content template system`,
		Version:    "3.0.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetOutput(os.Stdout)
			lvl, err := logrus.ParseLevel(Data.LogLevel)
			if err != nil {
				return err
			}
			logrus.SetLevel(lvl)
			logrus.Debug("root PersistentPreRunE")
			logrus.Debugf("Using config file: %v", viper.ConfigFileUsed())
			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debug("root PreRunE")
			bindAllFlagsForCommand(cmd)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debug("root RunE")
			return nil
		},
		PostRun: func(cmd *cobra.Command, args []string) {
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debug("root PostRunE")
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Debug("root PersistentPostRunE")
			return nil
		},
	}
	tmp.PersistentFlags().StringVar(&Data.LogLevel, "log-level", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic")
	tmp.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		var levels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
		return find(levels, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	tmp.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.prototype.yaml)")
	return tmp
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))

		viper.SetConfigName(".prototype")
		viper.SetConfigType("yml")
		viper.SetConfigType("yaml")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error using config file:", viper.ConfigFileUsed(), err)
	}

	viper.SetEnvPrefix(EnvironmentVariablePrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.Unmarshal(&Data)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

// contains checks if a string is present in a slice
func find(s []string, str string) []string {
	var matches []string
	if contains(s, str) {
		matches = append(matches, str)
	}
	return matches
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func bindAllFlagsForCommand(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to PROTO_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, fmt.Sprintf("%s_%s", EnvironmentVariablePrefix, envVarSuffix)) //nolint:errcheck
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			// fmt.Printf("name: %v changed:%v set:%v\n", f.Name, !f.Changed, v.IsSet(f.Name))
			cobra.CompDebugln(fmt.Sprintf("name: %v changed:%v set:%v\n", f.Name, !f.Changed, viper.IsSet(f.Name)), false)
			// fmt.Fprintf(os.Stderr, "DEBU[0000] Setting flag name: %v changed:%v set:%v\n", f.Name, !f.Changed, viper.IsSet(f.Name))
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)) //nolint:errcheck
		}
	})
}
