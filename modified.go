package main

import (
	"fmt"
	"slices"

	"github.com/haashemi/mypac/internal/arch"
	"github.com/haashemi/mypac/internal/config"
)

func runModified(args []string) error {
	fs, configPath := newFlags("modified")
	paths := fs.Bool("paths", false, "print only the paths, e.g. for `tar -T -`")
	fs.Parse(args)

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load your config: %w", err)
	}

	files, err := arch.ModifiedFiles()
	if err != nil {
		return fmt.Errorf("failed to query pacman: %w", err)
	}

	lastPkg := ""
	for _, f := range files {
		if slices.Contains(cfg.Modified.IgnorePackages, f.Package) || slices.Contains(cfg.Modified.IgnoreFiles, f.Path) {
			continue
		}

		if *paths {
			fmt.Println(f.Path)
			continue
		}

		// pacman lists packages in order, so a package's files are consecutive.
		if f.Package != lastPkg {
			boldGreenP.Println(f.Package)
			lastPkg = f.Package
		}
		blueP.Printf("%s- %s\n", prefixSpace, f.Path)
	}
	return nil
}
