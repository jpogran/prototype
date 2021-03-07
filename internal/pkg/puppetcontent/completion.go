package puppetcontent

import (
	"strings"
)

// CompleteName returns the template name matching the provided string
func CompleteName(templatePath string, match string) []string {
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
