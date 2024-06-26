package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"regexp"
)

func Find() ([]string, error) {
	var plugins []string
	home, err := os.UserHomeDir()
	if err != nil {
		return plugins, err
	}

	gobin := filepath.Join(home, "go", "bin")
	bins, err := os.ReadDir(gobin)
	if err != nil {
		return plugins, err
	}

	r := regexp.MustCompile("^go-llm-plugin-.*$")

	for _, f := range bins {
		if r.MatchString(f.Name()) {
			plugins = append(plugins, filepath.Join(gobin, f.Name()+".so"))
		}
	}

	return plugins, nil
}

type F func()

type Plugin struct {
	V int
	F F
}

func Load(path string) (Plugin, error) {
	fmt.Println("loading plugin: ", path)
	var p Plugin
	_, err := plugin.Open(path)
	if err != nil {
		return p, err
	}

	// v, err := so.Lookup("V")
	// if err != nil {
	// 	return p, err
	// }

	// f, err := so.Lookup("F")
	// if err != nil {
	// 	return p, err
	// }

	// p = Plugin{
	// 	V: v.(int),
	// 	F: f.(F),
	// }

	return p, nil
}
