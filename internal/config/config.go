package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	dirName    = ".recipe_box"
	configFile = "config.json"
)

var ErrKeyNotFound = errors.New("key not found")

type Config struct {
	path string
	data map[string]string
}

func New() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, dirName)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, configFile)
	c := &Config{
		path: configPath,
		data: make(map[string]string),
	}

	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return c, nil
}

func (c *Config) Get(key string) (string, error) {
	value, ok := c.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (c *Config) Set(key, value string) error {
	c.data[key] = value
	return c.save()
}

func (c *Config) load() error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(&c.data)
}

func (c *Config) save() error {
	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(c.data)
}
