package puppetcontent

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TemplateData struct {
	TemplatesPath string

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
	logrus.Debugf("Parsed: %+v", tmpl)

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
		TargetOutput, _ = filepath.Split(TargetOutput)
		TargetOutput = filepath.Clean(TargetOutput)
	}

	contentDir := filepath.Join(LocalTemplateCache, SelectedTemplate, "content")
	data := buildTemplateData(SelectedTemplate, LocalTemplateCache, TargetOutput, TargetName)

	logrus.Debugf("Name: %s", TargetName)
	logrus.Debugf("Output: %s", TargetOutput)
	logrus.Debugf("ContentDir: %+v", contentDir)
	logrus.Tracef("Data: %+v", data)

	var templateFiles []ContentTemplateFile
	if err := filepath.WalkDir(contentDir,
		func(path string, info os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			logrus.Tracef("Processing: %s", path)

			replacer := strings.NewReplacer(
				contentDir, TargetOutput,
				"__REPLACE__", TargetName,
				".tmpl", "",
			)
			targetFile := replacer.Replace(path)
			logrus.Tracef("Targetfile: %s", path)

			dir, file := filepath.Split(targetFile)
			i := ContentTemplateFile{
				TemplatePath:   path,
				TargetFilePath: targetFile,
				TargetDir:      dir,
				TargetFile:     file,
				IsDirectory:    info.IsDir(),
			}
			logrus.Tracef("Processed: %+v", i)

			templateFiles = append(templateFiles, i)
			return nil
		}); err != nil {
		log.Println(err)
	}

	var errs []error
	var deployedFiles []string

	for _, templateFile := range templateFiles {
		logrus.Tracef("Deploying: %+v", templateFile)

		if templateFile.IsDirectory {
			if _, err := os.Stat(templateFile.TargetFilePath); os.IsNotExist(err) {
				logrus.Tracef("Creating: %s", templateFile.TargetFilePath)
				err := os.MkdirAll(templateFile.TargetFilePath, os.ModePerm)
				if err != nil {
					logrus.Tracef("Created: %s", templateFile.TargetFilePath)
					errs = append(errs, err)
				}
			}
		} else {
			logrus.Tracef("Creating: %s", templateFile.TargetDir)
			err := os.MkdirAll(templateFile.TargetDir, os.ModePerm)
			if err != nil {
				logrus.Tracef("Created: %s", templateFile.TargetDir)
				errs = append(errs, err)
			}
			logrus.Tracef("Parsing: %s", templateFile.TemplatePath)
			t, err := template.ParseFiles(templateFile.TemplatePath)
			if err != nil {
				logrus.Tracef("Parsed: %s", templateFile.TemplatePath)
				errs = append(errs, err)
			}

			logrus.Tracef("Creating: %s", templateFile.TargetFilePath)
			actualTargetFile, err := os.Create(templateFile.TargetFilePath)
			if err != nil {
				logrus.Tracef("Created: %s", templateFile.TargetFilePath)
				errs = append(errs, err)
			}
			logrus.Tracef("Templating: %s", templateFile.TargetFilePath)
			err = t.Execute(actualTargetFile, data)
			if err != nil {
				logrus.Tracef("Templed: %s", templateFile.TargetFilePath)
				errs = append(errs, err)
			}
		}

		deployedFiles = append(deployedFiles, templateFile.TargetFilePath)
	}

	return deployedFiles, errs
}

func buildTemplateData(SelectedTemplate string, LocalTemplateCache string, TargetOutput string, TargetName string) TemplateData {
	data := TemplateData{
		TemplatesPath: LocalTemplateCache,
		TemplateName:  SelectedTemplate,
		ProjectName:   filepath.Base(TargetOutput),
		ItemName:      TargetName,

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
