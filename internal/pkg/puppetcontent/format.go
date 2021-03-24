package puppetcontent

import (
	"fmt"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
)

// Format displays an array of Puppet Content Template config entries
func Format(configs []ContentTemplateConfig, json bool) {
	if json {
		// fmt.Printf("%+v", configs)
		// prettyJSON, err := json.MarshalIndent(&configs, "", "  ")
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		prettyJSON, err := json.Marshal(&configs)
		if err != nil {
			log.Printf("Error converting to json: %v", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	} else {
		fmt.Println("")
		if len(configs) == 1 {
			fmt.Printf("DisplayName:     %v\n", configs[0].DisplayName)
			fmt.Printf("Name:            %v\n", configs[0].Name)
			fmt.Printf("Context:         %v\n", configs[0].Context)
			fmt.Printf("Tags:            %v\n", configs[0].Tags)
			fmt.Printf("TemplateType:    %v\n", configs[0].TemplateType)
			fmt.Printf("TemplateURL:     %v\n", configs[0].TemplateURL)
			fmt.Printf("TemplateVersion: %v\n", configs[0].TemplateVersion)
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"DisplayName", "Name", "Type"})
			table.SetBorder(false)
			for _, v := range configs {
				table.Append([]string{v.DisplayName, v.Name, v.TemplateType})
			}
			table.Render()
		}
	}
}
