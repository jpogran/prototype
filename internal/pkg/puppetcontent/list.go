package puppetcontent

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

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

func filterFiles(ss []ContentTemplateConfig, test func(ContentTemplateConfig) bool) (ret []ContentTemplateConfig) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
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
