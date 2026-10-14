// Package config finds and parses the user's mypac.{yml,yaml} file.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// Packages is the user's hand-written categorization of explicitly
	// installed packages.
	Packages []*PackageGroup `yaml:"packages"`

	// LostFiles maps a path (a file, or a directory covering everything
	// under it) to the reason it exists without being owned by a package.
	LostFiles map[string]string `yaml:"lostfiles"`

	Modified struct {
		IgnorePackages []string `yaml:"ignore_packages"`
		IgnoreFiles    []string `yaml:"ignore_files"`
	} `yaml:"modified"`
}

type PackageGroup struct {
	Name     string          `yaml:"name"`
	Packages []string        `yaml:"packages"`
	Groups   []*PackageGroup `yaml:"groups"`
}

// SearchPaths returns the locations checked for a config file, in order.
func SearchPaths() []string {
	var paths []string
	for _, dir := range []string{"/etc", "/etc/mypac"} {
		paths = append(paths, filepath.Join(dir, "mypac.yml"), filepath.Join(dir, "mypac.yaml"))
	}
	return paths
}

// Load parses the config at path, or the first one found in [SearchPaths]
// when path is empty.
func Load(path string) (*Config, error) {
	if path == "" {
		for _, p := range SearchPaths() {
			if _, err := os.Stat(p); err == nil {
				path = p
				break
			}
		}
		if path == "" {
			return nil, errors.New("no config file found, searched:\n  " + strings.Join(SearchPaths(), "\n  "))
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

func parse(data []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		// Legacy format: the whole file is a bare list of package groups.
		if yaml.Unmarshal(data, &c.Packages) != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
	}
	return &c, nil
}

// LostFileReason returns the configured path covering file and its reason.
// The most specific (longest) matching path wins.
func (c *Config) LostFileReason(file string) (path, reason string, ok bool) {
	for p, r := range c.LostFiles {
		trimmed := strings.TrimSuffix(p, "/")
		if (file == trimmed || strings.HasPrefix(file, trimmed+"/")) && len(p) > len(path) {
			path, reason, ok = p, r, true
		}
	}
	return path, reason, ok
}
