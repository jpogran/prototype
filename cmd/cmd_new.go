package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var newCmd = createNewCommand()

func createNewCommand() *cobra.Command {
	tmp := &cobra.Command{
		Use:   "new",
		Short: "creates a Puppet project or other artifact based on a template",
		Long:  `creates a Puppet project or other artifact based on a template`,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) != 0 {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return Complete(Config.LocalTemplatePath, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		},
		Args: func(cmd *cobra.Command, args []string) error {
			log.Printf("args: %v", args)

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
			// prototype new foo-foo
			if Config.TargetName == "" && Config.TargetOutput == "" {
				cwd, _ := os.Getwd()
				Config.TargetName = filepath.Base(cwd)
				Config.TargetOutput = cwd
			}

			// prototype new foo-foo -n wakka
			if Config.TargetName != "" && Config.TargetOutput == "" {
				cwd, _ := os.Getwd()
				Config.TargetOutput = filepath.Join(cwd, Config.TargetName)
			}

			// prototype new foo-foo -o /foo/bar/baz
			if Config.TargetName == "" && Config.TargetOutput != "" {
				Config.TargetName = filepath.Base(Config.TargetOutput)
			}

			// prototype new foo-foo
			if Config.TargetName == "" {
				cwd, _ := os.Getwd()
				Config.TargetName = filepath.Base(cwd)
			}

			// prototype new foo-foo
			// prototype new foo-foo -n wakka
			// prototype new foo-foo -n wakka -o c:/foo
			// prototype new foo-foo -n wakka -o c:/foo/wakka
			if Config.TargetOutput == "" {
				cwd, _ := os.Getwd()
				Config.TargetOutput = cwd
			} else if strings.HasSuffix(Config.TargetOutput, Config.TargetName) {
				// user has specified outputpath with the targetname in it
			} else {
				Config.TargetOutput = filepath.Join(Config.TargetOutput, Config.TargetName)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "templatepath", Config.LocalTemplatePath)
			fmt.Fprintln(out, "name", Config.TargetName)
			fmt.Fprintln(out, "output", Config.TargetOutput)
			fmt.Fprintln(out, "list", Config.ListTemplates)
			fmt.Fprintln(out, "tmpl", Config.SelectedTemplate)

			if Config.ListTemplates {
				tmpls, err := List(Config.LocalTemplatePath, Config.SelectedTemplate)
				if err != nil {
					return err
				}

				Format(tmpls, Config.OutputFormat)

				return nil
			}

			data := map[string]interface{}{}
			data["ProjectName"] = filepath.Base(Config.TargetOutput)
			data["ItemName"] = Config.TargetName
			data["Author"] = viper.GetString("author")
			data["Summary"] = viper.GetString("summary")
			data["License"] = viper.GetString("license")
			data["Source"] = viper.GetString("source")
			data["Prototype"] = map[string]interface{}{
				"Version":    "0.1.0",
				"CommitHash": "a4b89ba",
				"BuildDate":  "2020-06-27",
			}

			deployed, errs := Deploy(data, Config.TargetName, Config.TargetOutput, Config.SelectedTemplate, Config.LocalTemplatePath)
			for e := range errs {
				log.Printf("Error: %v", e)
			}

			switch Config.OutputFormat {
			case "table":
				fmt.Println("")
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
	tmp.Flags().
		StringVarP(
			&Config.LocalTemplatePath,
			"templates", "t",
			"",
			"template install directory (default is $HOME/.pdk/templates)",
		)
	tmp.Flags().
		StringVarP(
			&Config.TargetName,
			"name", "n",
			"",
			"the name for the created output. (default is the name of the current directory)",
		)
	tmp.Flags().
		StringVarP(
			&Config.TargetOutput,
			"output", "o",
			"",
			"location to place the generated output. (default is the current directory)",
		)
	tmp.Flags().
		StringVarP(
			&Config.OutputFormat,
			"format", "f",
			"",
			"formating (default is table)",
		)
	tmp.Flags().
		BoolVarP(
			&Config.ListTemplates,
			"list", "l",
			false,
			"list templates",
		)
	tmp.RegisterFlagCompletionFunc("list", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint:errcheck
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return Complete(Config.LocalTemplatePath, toComplete), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	})

	tmp.Flags().StringVar(&Config.Author, "author", "", "the author name. (default is the current user name)")
	tmp.Flags().StringVar(&Config.Summary, "summary", "", "the purpose of the content")
	tmp.Flags().StringVar(&Config.License, "license", "", "the license for the content (default is Apache2)")
	tmp.Flags().StringVar(&Config.Source, "source", "", "the source control repository link for this content")
	return tmp
}

func init() {
	rootCmd.AddCommand(newCmd)
}

// Complete returns the template name matching the provided string
func Complete(templatePath string, match string) []string {
	tmpls, _ := List(templatePath, "")
	var names []string
	for _, tmpl := range tmpls {
		if strings.HasPrefix(tmpl.Name, match) {
			m := tmpl.Name + "\t" + tmpl.DisplayName
			names = append(names, m)
		}
	}
	return names
}

// List gets all installed Puppet Content Templates
func List(templatePath string, templateName string) ([]ContentTemplateConfig, error) {
	var tmpls []ContentTemplateConfig
	if err := filepath.WalkDir(templatePath,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			file := filepath.Join(templatePath, info.Name(), "templateconfig.yml")
			if !pathExists(file) {
				return nil
			}

			tmpl, err := read(file)
			if err != nil {
				// TODO print stderr here not stdout
				log.Println("Error: ", err)
			}

			tmpls = append(tmpls, tmpl)

			return nil
		}); err != nil {
		return []ContentTemplateConfig{}, err
	}

	if templateName != "" {
		tmpls = filterFiles(tmpls, func(f ContentTemplateConfig) bool { return f.Name == templateName })
	}

	return tmpls, nil
}

// Deploy a Puppet Content Template to a given output folder
func Deploy(
	data map[string]interface{},
	name string,
	output string,
	templateName string,
	mainTemplatePath string,
) ([]string, []error) {

	mainTemplatePath = filepath.Join(mainTemplatePath, templateName)
	contentDir := filepath.Join(mainTemplatePath, "content")
	config, _ := read(filepath.Join(mainTemplatePath, "templateconfig.yml"))

	log.Printf("TemplateDir: %s\n", mainTemplatePath)
	log.Printf("Output: %s\n", output)
	log.Printf("Name: %s\n", name)

	data["ContentTemplateConfig"] = map[string]interface{}{
		"URL":     config.TemplateURL,
		"Version": config.TemplateVersion,
	}

	var templateFiles []ContentTemplateFile
	if err := filepath.WalkDir(contentDir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			replacer := strings.NewReplacer(
				contentDir, output,
				"__REPLACE__", name,
				".tmpl", "",
			)
			targetFile := replacer.Replace(path)

			dir, file := filepath.Split(targetFile)
			i := ContentTemplateFile{
				TemplatePath:   path,
				TargetFilePath: targetFile,
				TargetDir:      dir,
				TargetFile:     file,
				IsDirectory:    info.IsDir(),
			}

			templateFiles = append(templateFiles, i)
			return nil
		}); err != nil {
		log.Println(err)
	}

	var errs []error
	var deployedFiles []string

	for _, templateFile := range templateFiles {
		if templateFile.IsDirectory {
			if _, err := os.Stat(templateFile.TargetFilePath); os.IsNotExist(err) {
				os.Mkdir(templateFile.TargetFilePath, os.ModePerm)
			}
		} else {
			err := os.MkdirAll(templateFile.TargetDir, os.ModePerm)
			if err != nil {
				errs = append(errs, err)
			}
			t, err := template.ParseFiles(templateFile.TemplatePath)
			if err != nil {
				errs = append(errs, err)
			}

			actualTargetFile, err := os.Create(templateFile.TargetFilePath)
			if err != nil {
				errs = append(errs, err)
			}
			err = t.Execute(actualTargetFile, data)
			if err != nil {
				errs = append(errs, err)
			}
		}

		deployedFiles = append(deployedFiles, templateFile.TargetFilePath)
	}

	return deployedFiles, errs
}

// Format displays an array of Puppet Content Template config entries
func Format(tmpls []ContentTemplateConfig, format string) {
	switch format {
	case "table":
		fmt.Println("")
		if len(tmpls) == 1 {
			fmt.Printf("DisplayName:     %v\n", tmpls[0].DisplayName)
			fmt.Printf("Name:            %v\n", tmpls[0].Name)
			fmt.Printf("Context:         %v\n", tmpls[0].Context)
			fmt.Printf("Tags:            %v\n", tmpls[0].Tags)
			fmt.Printf("TemplateType:    %v\n", tmpls[0].TemplateType)
			fmt.Printf("TemplateURL:     %v\n", tmpls[0].TemplateURL)
			fmt.Printf("TemplateVersion: %v\n", tmpls[0].TemplateVersion)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"DisplayName", "Name", "Type"})
			table.SetBorder(false)
			for _, v := range tmpls {
				table.Append([]string{v.DisplayName, v.Name, v.TemplateType})
			}
			table.Render()
		}
	case "json":
		prettyJSON, _ := json.MarshalIndent(tmpls, "", "  ")
		fmt.Printf("%s\n", string(prettyJSON))
	}
}

// ContentTemplateConfig is the config for a Puppet Content Template module
type ContentTemplateConfig struct {
	Name            string   `yaml:"name"`
	DisplayName     string   `yaml:"display_name"`
	Context         string   `yaml:"context"`
	Tags            []string `yaml:"tags"`
	TemplateType    string   `yaml:"template_type"`
	TemplateURL     string   `yaml:"template_url"`
	TemplateVersion string   `yaml:"template_version"`
}

// ContentTemplateFile contains the path information for a given file inside a ContentTemplate
type ContentTemplateFile struct {
	TemplatePath   string
	TargetFilePath string
	TargetDir      string
	TargetFile     string
	IsDirectory    bool
}

func read(path string) (ContentTemplateConfig, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return ContentTemplateConfig{}, err
	}

	var config ContentTemplateConfig
	if err := config.parse(yamlFile); err != nil {
		return ContentTemplateConfig{}, err
	}

	return config, nil
}

func (template *ContentTemplateConfig) parse(data []byte) error {
	if err := yaml.Unmarshal(data, template); err != nil {
		return err
	}
	return nil
}

func filterFiles(ss []ContentTemplateConfig, test func(ContentTemplateConfig) bool) (ret []ContentTemplateConfig) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return false
}
