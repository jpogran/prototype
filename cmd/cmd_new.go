package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jpogran/prototype/internal/pkg/puppetcontent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = createNewCommand()

func init() {
	rootCmd.AddCommand(newCmd)
}

func createNewCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "new",
		Short: "creates a Puppet project or other artifact based on a template",
		Long: `creates a Puppet project or other artifact based on a template
		Project Level Variables
			{{.ProjectName}}
			{{.ItemName}}

		Global Variables
			{{.TemplatesPath}}
			{{.Author}}
			{{.Summary}}
			{{.License}}
			{{.Source}}
			{{.PuppetContentTemplate.URL}}
			{{.PuppetContentTemplate.Version}}

		Prototype Variables
			{{.Prototype.Version}}
			{{.Prototype.CommitHash}}
			{{.Prototype.BuildDate}}`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return puppetcontent.CompleteName(Config.LocalTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && !Config.ListTemplates {
				Config.ListTemplates = true
			}

			if Config.TargetName == "" && len(args) == 2 {
				Config.TargetName = args[1]
			}

			if len(args) >= 1 {
				Config.SelectedTemplate = args[0]
			}

			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if Config.ListTemplates {
				tmpls, err := puppetcontent.List(Config.LocalTemplateCache, Config.SelectedTemplate)
				if err != nil {
					return err
				}

				puppetcontent.Format(tmpls, Config.OutputFormat)

				return nil
			}

			deployed, errs := puppetcontent.Deploy(Config.SelectedTemplate, Config.LocalTemplateCache, Config.TargetOutput, Config.TargetName)
			for e := range errs {
				log.Printf("Error: %v", e)
			}

			switch Config.OutputFormat {
			case "table":
				for _, d := range deployed {
					log.Printf("Deployed: %v", d)
				}
			case "json":
				prettyJSON, _ := json.MarshalIndent(deployed, "", "  ")
				fmt.Printf("%s\n", string(prettyJSON))
			}

			return nil
		},
	}
	tmp.Flags().StringVarP(&Config.LocalTemplateCache, "templates", "t", "", "template install directory (default is $HOME/.pdk/templates)")
	viper.BindPFlag("templates", tmp.Flags().Lookup("templates")) //nolint:errcheck
	tmp.Flags().StringVarP(&Config.TargetName, "name", "n", "", "the name for the created output. (default is the name of the current directory)")
	tmp.Flags().StringVarP(&Config.TargetOutput, "output", "o", "", "location to place the generated output. (default is the current directory)")
	tmp.Flags().StringVarP(&Config.OutputFormat, "formatoutput", "f", "", "formating (default is table)")
	tmp.Flags().BoolVarP(&Config.ListTemplates, "list", "l", false, "list templates")
	tmp.RegisterFlagCompletionFunc("list", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return puppetcontent.CompleteName(Config.LocalTemplateCache, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})
	tmp.Flags().StringVar(&Config.Author, "author", "", "the author name. (default is the current user name)")
	tmp.Flags().StringVar(&Config.Summary, "summary", "", "the purpose of the content")
	tmp.Flags().StringVar(&Config.License, "license", "", "the license for the content (default is Apache2)")
	tmp.Flags().StringVar(&Config.Source, "source", "", "the source control repository link for this content")
	return tmp
}
