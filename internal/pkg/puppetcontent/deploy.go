package puppetcontent

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

type TemplateData struct {
	TemplateName string

	ProjectName string
	ItemName    string

	Author  string
	Summary string
	License string
	Source  string

	Prototype PrototypeInfo

	TemplateConfig ContentTemplateConfig
}

type PrototypeInfo struct {
	Version    string
	CommitHash string
	BuildDate  string
}

// ContentTemplateFile contains the path information for a given file inside a ContentTemplate
type ContentTemplateFile struct {
	TemplatePath   string
	TargetFilePath string
	TargetDir      string
	TargetFile     string
	IsDirectory    bool
}

// Deploy a Puppet Content Template to a given output folder
func Deploy(SelectedTemplate string, LocalTemplateCache string, TargetOutput string, TargetName string) ([]string, []error) {

	file := filepath.Join(LocalTemplateCache, SelectedTemplate, "templateconfig.yml")
	tmpl, _ := read(file)

	// prototype new foo-foo
	if TargetName == "" && TargetOutput == "" {
		cwd, _ := os.Getwd()
		TargetName = filepath.Base(cwd)
		TargetOutput = cwd
	}

	// prototype new foo-foo -n wakka
	if TargetName != "" && TargetOutput == "" {
		cwd, _ := os.Getwd()
		TargetOutput = filepath.Join(cwd, TargetName)
	}

	// prototype new foo-foo -o /foo/bar/baz
	if TargetName == "" && TargetOutput != "" {
		TargetName = filepath.Base(TargetOutput)
	}

	// prototype new foo-foo
	if TargetName == "" {
		cwd, _ := os.Getwd()
		TargetName = filepath.Base(cwd)
	}

	// prototype new foo-foo
	// prototype new foo-foo -n wakka
	// prototype new foo-foo -n wakka -o c:/foo
	// prototype new foo-foo -n wakka -o c:/foo/wakka
	switch tmpl.TemplateType {
	case "project":
		if TargetOutput == "" {
			cwd, _ := os.Getwd()
			TargetOutput = cwd
		} else if strings.HasSuffix(TargetOutput, TargetName) {
			// user has specified outputpath with the targetname in it
		} else {
			TargetOutput = filepath.Join(TargetOutput, TargetName)
		}
	case "item":
		if TargetOutput == "" {
			cwd, _ := os.Getwd()
			TargetOutput = cwd
		}
		//  else if strings.HasSuffix(TargetOutput, TargetName) {
		// 	// user has specified outputpath with the targetname in it
		// } else {
		// 	// use what the user tells us
		// }
	}

	data := buildTemplateData(SelectedTemplate, LocalTemplateCache, TargetOutput, TargetName)

	contentDir := filepath.Join(LocalTemplateCache, SelectedTemplate, "content")

	log.Printf("TemplateDir: %s\n", LocalTemplateCache)
	log.Printf("Output: %s\n", TargetOutput)
	log.Printf("Name: %s\n", TargetName)

	var templateFiles []ContentTemplateFile
	if err := filepath.WalkDir(contentDir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			replacer := strings.NewReplacer(
				contentDir, TargetOutput,
				"__REPLACE__", TargetName,
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
				err := os.MkdirAll(templateFile.TargetFilePath, os.ModePerm)
				if err != nil {
					errs = append(errs, err)
				}
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

func buildTemplateData(SelectedTemplate string, LocalTemplateCache string, TargetOutput string, TargetName string) TemplateData {
	data := TemplateData{
		TemplateName: SelectedTemplate,
		ProjectName:  filepath.Base(TargetOutput),
		ItemName:     TargetName,

		Author:  viper.GetString("author"),
		Summary: viper.GetString("summary"),
		License: viper.GetString("license"),
		Source:  viper.GetString("source"),
	}

	localTemplatePath := filepath.Join(LocalTemplateCache, SelectedTemplate)
	tmplConfig, _ := read(filepath.Join(localTemplatePath, "templateconfig.yml"))
	data.TemplateConfig = tmplConfig

	data.Prototype = PrototypeInfo{
		Version:    viper.GetString("version"),
		CommitHash: viper.GetString("commithash"),
		BuildDate:  viper.GetString("builddate"),
	}

	return data
}
