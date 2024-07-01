package plugins

import (
	"os"
	"path/filepath"
	"regexp"
)

type PluginListItem struct {
	Name string
	Path string
}

func Find() ([]PluginListItem, error) {
	var plugins []PluginListItem
	home, err := os.UserHomeDir()
	if err != nil {
		return plugins, err
	}

	gobin := filepath.Join(home, "go", "bin")
	bins, err := os.ReadDir(gobin)
	if err != nil {
		return plugins, err
	}

	r := regexp.MustCompile("^go-llm-plugin-(.*)$")

	for _, f := range bins {
		if r.MatchString(f.Name()) {
			plugins = append(plugins, PluginListItem{
				Name: r.FindStringSubmatch(f.Name())[1],
				Path: filepath.Join(gobin, f.Name()),
			})
		}
	}

	return plugins, nil
}
