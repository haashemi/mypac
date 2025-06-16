package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type PackageGroup struct {
	Name     string          `yaml:"name"`
	Packages []string        `yaml:"packages"`
	Groups   []*PackageGroup `yaml:"groups"`
}

// loadPackageGroups returns a slice of [PackageGroup] read from the specified
// YAML file.
func loadPackageGroups(path string) ([]*PackageGroup, error) {
	packagesFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var data []*PackageGroup
	err = yaml.Unmarshal(packagesFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// loadInstalledPackages returns a [PackageGroup] of explicitly installed
// packages by executing the `pacman -Qqe` command.
func loadInstalledPackages(groupName string) (*PackageGroup, error) {
	cmd := exec.Command("pacman", "-Qqe")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	group := &PackageGroup{Name: groupName}
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		pkgName := scanner.Text()
		if pkgName != "" {
			group.Packages = append(group.Packages, pkgName)
		}
	}

	return group, nil
}
