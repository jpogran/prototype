package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/jpogran/prototype/internal/pkg/puppetcontent"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	localTemplateCache string
	jsonOutput         bool

	selectedTemplate string
	listTemplates    bool

	targetName   string
	targetOutput string
)

func CreateNewCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "new",
		Short: "creates a Puppet project or other artifact based on a template",
		Long:  `creates a Puppet project or other artifact based on a template`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && !listTemplates {
				listTemplates = true
			}

			if targetName == "" && len(args) == 2 {
				targetName = args[1]
			}

			if len(args) >= 1 {
				selectedTemplate = args[0]
			}

			return nil
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			localTemplateCache = viper.GetString("templatepath")
			return puppetcontent.CompleteName(localTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logrus.Trace("new PreRunE")
			bindAllFlagsForCommand(cmd)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Trace("new Run")
			logrus.Tracef("Templatepath: %v", localTemplateCache)

			if listTemplates {
				tmpls, err := puppetcontent.List(localTemplateCache, "")
				if err != nil {
					return err
				}

				puppetcontent.Format(tmpls, jsonOutput)

				return nil
			}

			deployed, errs := puppetcontent.Deploy(
				selectedTemplate,
				localTemplateCache,
				targetOutput,
				targetName,
			)

			for e := range errs {
				logrus.Errorf("Error: %v", e)
			}

			if jsonOutput {
				prettyJSON, _ := json.MarshalIndent(deployed, "", "  ")
				fmt.Printf("%s\n", string(prettyJSON))
			} else {
				for _, d := range deployed {
					logrus.Infof("Deployed: %v", d)
				}
			}

			return nil
		},
	}
	home, _ := homedir.Dir()
	tmp.Flags().StringVar(&localTemplateCache, "templatepath", "", "Log level (debug, info, warn, error, fatal, panic")
	viper.SetDefault("templatepath", filepath.Join(home, ".pdk", "puppet-content-templates"))
	viper.BindPFlag("templatepath", tmp.Flags().Lookup("templatepath")) //nolint:errcheck

	tmp.Flags().BoolVar(&jsonOutput, "json", false, "json output")
	tmp.Flags().BoolVarP(&listTemplates, "list", "l", false, "list templates")
	tmp.RegisterFlagCompletionFunc("list", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return puppetcontent.CompleteName(localTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	tmp.Flags().StringVarP(&targetName, "name", "n", "", "the name for the created output. (default is the name of the current directory)")
	tmp.Flags().StringVarP(&targetOutput, "output", "o", "", "location to place the generated output. (default is the current directory)")
	return tmp
}
