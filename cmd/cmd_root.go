package cmd

import (
	"fmt"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configFile = ".prototype"
	envPrefix  = "PROTO"
)

var rootCmd = createRootCommand()

func createRootCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:     "prototype",
		Short:   "PDK prototype content template system",
		Long:    `PDK prototype content template system`,
		Version: "3.0.0",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			// out := cmd.OutOrStdout()
			// fmt.Fprintln(out, "The param is:", param)
		},
	}
	return tmp
}

func init() {}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(configFile)

	home, err := homedir.Dir()
	cobra.CheckErr(err)
	v.AddConfigPath(home)
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to PROTO_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
	return nil
}
