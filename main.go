package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/fatih/color"
)

const (
	// prefixSpace is the space used for indentation in the output
	prefixSpace = "  "
)

var (
	filePath = flag.String("path", "packages.yml", "Path to the package groups file")
)

var (
	boldGreenP = color.New(color.FgGreen, color.Bold)
	boldRedP   = color.New(color.FgHiRed, color.Bold)
	blueP      = color.New(color.FgBlue)
)

func main() {
	flag.Parse()

	installed, err := loadInstalledPackages("Not Categorized")
	if err != nil {
		boldRedP.Println("Failed to load installed packages:", err)
		return
	}

	packageGroups, err := loadPackageGroups(*filePath)
	if err != nil {
		boldRedP.Println("Failed to load your package groups:", err)
		return
	}

	for _, group := range packageGroups {
		printPackageGroup(group, installed, "")
		fmt.Println()
	}
	printPackageGroup(installed, installed, "")
}

// printPackageGroup prints the group name and its packages and nested groups
// inside by recursively calling itself. It updates the installed [PackageGroup]
// so we could find the packages that are not categorized and have a smaller
// array to lookup for each package.
func printPackageGroup(group *PackageGroup, installed *PackageGroup, prefix string) {
	if group == nil || installed == nil {
		return
	}

	boldGreenP.Println(prefix + group.Name)
	for _, pkg := range group.Packages {
		pkgIndex := slices.Index(installed.Packages, pkg)
		isInstalled := pkgIndex != -1

		if isInstalled {
			blueP.Printf("%s%s- %s\n", prefix, prefixSpace, pkg)
			installed.Packages = slices.Delete(installed.Packages, pkgIndex, pkgIndex+1)
		} else {
			boldRedP.Printf("%s%s- %s: (Not installed)\n", prefix, prefixSpace, pkg)
		}
	}

	for _, subGroup := range group.Groups {
		printPackageGroup(subGroup, installed, prefix+prefixSpace)
	}
}
