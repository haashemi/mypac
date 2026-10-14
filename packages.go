package main

import (
	"fmt"
	"slices"

	"github.com/haashemi/mypac/internal/arch"
	"github.com/haashemi/mypac/internal/config"
)

func runPackages(args []string) error {
	fs, configPath := newFlags("packages")
	fs.Parse(args)

	installed, err := arch.ExplicitPackages()
	if err != nil {
		return fmt.Errorf("failed to load installed packages: %w", err)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load your config: %w", err)
	}

	for _, group := range cfg.Packages {
		installed = printPackageGroup(group, installed, "")
		fmt.Println()
	}

	boldGreenP.Println("Not Categorized")
	for _, pkg := range installed {
		blueP.Printf("%s- %s\n", prefixSpace, pkg)
	}
	return nil
}

// printPackageGroup prints the group name and its packages and nested groups
// inside by recursively calling itself. It returns the installed packages
// left uncategorized, so each lookup searches a smaller slice and the rest
// can be printed as not categorized.
func printPackageGroup(group *config.PackageGroup, installed []string, prefix string) []string {
	if group == nil {
		return installed
	}

	boldGreenP.Println(prefix + group.Name)
	for _, pkg := range group.Packages {
		if i := slices.Index(installed, pkg); i != -1 {
			blueP.Printf("%s%s- %s\n", prefix, prefixSpace, pkg)
			installed = slices.Delete(installed, i, i+1)
		} else {
			boldRedP.Printf("%s%s- %s: (Not installed)\n", prefix, prefixSpace, pkg)
		}
	}

	for _, subGroup := range group.Groups {
		installed = printPackageGroup(subGroup, installed, prefix+prefixSpace)
	}
	return installed
}
