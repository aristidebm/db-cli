package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"example.com/db/internal/shutil"
	"example.com/db/internal/types"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Sources map[string]types.Source `toml:"sources"`
	Schemes map[string]types.Scheme `toml:"schemes"`
}

func (c *Config) Load() error {
	configPath := c.getConfigPath()

	if err := shutil.CreateDirIfNotExists(filepath.Dir(configPath)); err != nil {
		return fmt.Errorf("%w: failed to create config directory", err)
	}

	if !shutil.FileExists(configPath) {
		c.Sources = map[string]types.Source{}
		c.Schemes = map[string]types.Scheme{}
		return c.Save()
	}

	if _, err := toml.DecodeFile(configPath, c); err != nil {
		return fmt.Errorf("%w: failed to load config file", err)
	}

	if c.Sources == nil {
		c.Sources = map[string]types.Source{}
	}

	if c.Schemes == nil {
		c.Schemes = map[string]types.Scheme{}
	}

	return c.validate()
}

func (c *Config) Save() error {

	configPath := c.getConfigPath()

	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("%w: failed to create config file", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("%w: failed to encode config", err)
	}
	return nil
}

func (c *Config) AddSource(name string, url string) error {
	u, err := shutil.ParseURL(url)
	if err != nil {
		return err
	}

	if !c.isSchemeSupported(u.Scheme) {
		return fmt.Errorf("%w: support for '%s' is coming soon", UnsupportedScheme, u.Scheme)
	}

	if c.Sources == nil {
		c.Sources = map[string]types.Source{}
	}
	if _, ok := c.Sources[name]; ok {
		return fmt.Errorf("%w: source '%s' already exist", InvalidSource, name)
	}
	c.Sources[name] = types.Source{URL: url}
	return c.Save()
}

func (c *Config) RemoveSource(name string) error {
	if c.Sources == nil {
		return fmt.Errorf("%w: source '%s' not found", InvalidSource, name)
	}
	if _, ok := c.Sources[name]; !ok {

		return fmt.Errorf("%w: source '%s' not found", InvalidSource, name)
	}
	delete(c.Sources, name)
	return c.Save()
}

func (c *Config) GetSource(name string) (types.Source, error) {
	source, ok := c.Sources[name]
	if !ok {
		return types.Source{}, fmt.Errorf("%w: source '%s' not found", InvalidSource, name)
	}
	return source, nil
}

func (c *Config) GetScheme(name string) (types.Scheme, error) {
	scheme, ok := c.Schemes[name]
	if !ok {
		return types.Scheme{}, fmt.Errorf("%w: scheme '%s' not found", UnsupportedScheme, name)
	}
	return scheme, nil
}

func (c *Config) Edit() error {
	configPath := c.getConfigPath()
	editor := shutil.GetEditor()

	cmd := exec.Command(editor, configPath)
	return shutil.Run(cmd,
		shutil.WithStdin(os.Stdin),
		shutil.WithStdout(os.Stdout),
		shutil.WithStderr(os.Stderr),
	)
}

func (c *Config) ListSources() error {
	for src, _ := range c.Sources {
		fmt.Println(src)
	}
	return nil
}

func (c *Config) ListSchemes() error {
	for drv, _ := range c.Schemes {
		fmt.Println(drv)
	}
	return nil
}

func (c *Config) getConfigPath() string {
	return filepath.Join(shutil.GetConfigDir(), "config.toml")
}

func (c *Config) validate() error {

	for name, source := range c.Sources {
		if source.URL == "" {
			return fmt.Errorf("%w: missing url in source %s",
				InvalidSource, name)
		}

		if source.Interactive != "" && !c.isExecutable(source.Interactive) {
			return fmt.Errorf("%w: client '%s' in source %s is not executable",
				InvalidClient, source.Interactive, name)
		}
	}

	for name, scheme := range c.Schemes {
		if scheme.Interactive != "" && !c.isExecutable(scheme.Interactive) {
			return fmt.Errorf("%w: client '%s' in scheme %s is not executable",
				InvalidClient, scheme.Interactive, name)
		}
	}

	return nil
}

func (c *Config) isExecutable(command string) bool {
	if command == "" {
		return false
	}
	return shutil.IsCommandInstalled(strings.Fields(command)[0])
}

func (c *Config) isSchemeSupported(scheme string) bool {
	supported := []string{
		"sqlite3",
		"postgres",
		"mysql",
		"redis",
	}
	return slices.Contains(supported, scheme)
}
