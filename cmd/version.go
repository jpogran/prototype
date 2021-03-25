package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func NewCmdVersion(version string, buildDate string, commit string) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(Format(version, buildDate, commit))
		},
	}

	return cmd
}

func Format(version string, buildDate string, commit string) string {
	version = strings.TrimPrefix(version, "v")

	var dateStr string
	if buildDate != "" {
		d, _ := time.Parse(time.RFC3339, buildDate)
		dateStr = d.Format("2006-01-02")
	}

	commitSHA := substr(commit, 0, 7)

	return fmt.Sprintf("prototype version %s %s %s\n%s\n", version, commitSHA, dateStr, changelogURL(version))
}

func changelogURL(version string) string {
	path := "https://github.com/jpogran/prototype"
	r := regexp.MustCompile(`^v?\d+\.\d+\.\d+(-[\w.]+)?$`)
	if !r.MatchString(version) {
		return fmt.Sprintf("%s/releases/latest", path)
	}

	url := fmt.Sprintf("%s/releases/tag/v%s", path, strings.TrimPrefix(version, "v"))
	return url
}

// NOTE: this isn't multi-Unicode-codepoint aware, like specifying skintone or
//       gender of an emoji: https://unicode.org/emoji/charts/full-emoji-modifiers.html
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
