package config

import (
	"fmt"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/bcdxn/go-llm/internal/plugins"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SelectedPlugin plugins.PluginListItem `yaml:"selectedPlugin"`
	SelectedModel  string                 `yaml:"selectedModel"`
	Log            LogConfig              `yaml:"log"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type CtxConfig struct{}

// Init creates the config file if it does not already exist and saves it and also returns the
// unmarshalled configuration (default if no config file was found)
func Init() (Config, error) {
	var cfg Config
	cfg, err := Load()
	if err != nil {
		return cfg, err
	}

	err = Persist(cfg)
	return cfg, err
}

// Load initializes the config file if it doesn't exist and returns the unmarshalled config data.
func Load() (Config, error) {
	var cfg Config
	cfgPath, err := createConfigDirectory()
	if err != nil {
		return cfg, err
	}

	cfg = Config{
		Log: LogConfig{
			Level: "error",
			File:  filepath.Join(cfgPath, "log"),
		},
	}

	cfgPath, err = createConfigFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	loadedCfg, err := readConfigFile(cfgPath)
	if err != nil {
		return cfg, err
	}

	var finalCfg Config

	err = mergo.Merge(&finalCfg, loadedCfg, mergo.WithOverride)

	return cfg, nil
}

// Save the configuration to file
func Persist(cfg Config) error {
	path, err := ConfigFilePath()
	if err != nil {
		return err
	}

	err = writeConfigFile(cfg, path)
	return err
}

// Get the absolute path to the folder where the configuration file lives
func ConfigDirPath() (string, error) {
	var cfgPath string
	home, err := os.UserHomeDir()
	if err != nil {
		return cfgPath, fmt.Errorf("unable to get user home directory location: %w", err)
	}
	cfgPath = filepath.Join(home, ".llm")
	return cfgPath, nil
}

// Get the absolute path to the configuration file
func ConfigFilePath() (string, error) {
	var cfgPath string
	dir, err := ConfigDirPath()
	if err != nil {
		return cfgPath, nil
	}

	cfgPath = filepath.Join(dir, "config")

	return cfgPath, nil
}

// createConfigDirectory creates the directory that holds the config file if it does not exist
func createConfigDirectory() (string, error) {
	cfgPath, err := ConfigDirPath()
	if err != nil {
		return cfgPath, nil
	}

	err = os.MkdirAll(cfgPath, os.ModePerm)
	if err != nil {
		return cfgPath, fmt.Errorf("unable to create ~/.llm directory: %w", err)
	}
	return cfgPath, nil
}

// createConfigFile opens and returns the file, creating the file if it does not exist
func createConfigFile(dirPath string) (string, error) {
	cfgPath := filepath.Join(dirPath, "config")
	f, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return cfgPath, fmt.Errorf("unable to open ~/.llm/config file: %w", err)
	}
	defer closeConfigFile(f)
	return cfgPath, nil
}

// readConfigFile returns the unmarshalled configuration from the config file
func readConfigFile(path string) (Config, error) {
	cfg := Config{}
	c, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(c, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func writeConfigFile(cfg Config, path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer closeConfigFile(f)

	d, err := yaml.Marshal(&cfg)
	if err != nil {
		return err
	}

	_, err = f.Write(d)
	return err
}

// closeConfigFile closes the config file and should be called via defer after createConfigFile is
// called.
func closeConfigFile(f *os.File) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}
