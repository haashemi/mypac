package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/haashemi/mypac/internal/arch"
	"github.com/haashemi/mypac/internal/config"
)

func runLostFiles(args []string) error {
	fs, configPath := newFlags("lostfiles")
	strict := fs.Bool("strict", false, "run lostfiles in strict mode (-s)")
	paths := fs.Bool("paths", false, "print only the paths, e.g. for `tar -T -`")
	fs.Parse(args)

	cfg, err := config.Load(*configPath)
	if err != nil {
		return fmt.Errorf("failed to load your config: %w", err)
	}

	files, err := arch.LostFiles(*strict)
	if err != nil {
		return fmt.Errorf("failed to run lostfiles: %w", err)
	}

	if *paths {
		for _, f := range files {
			fmt.Println(f)
		}
		return nil
	}

	// reason -> configured path -> number of lost files under it
	explained := map[string]map[string]int{}
	var unexplained []string
	for _, f := range files {
		path, reason, ok := cfg.LostFileReason(f)
		if !ok {
			unexplained = append(unexplained, f)
			continue
		}
		if explained[reason] == nil {
			explained[reason] = map[string]int{}
		}
		explained[reason][path]++
	}

	for _, reason := range slices.Sorted(maps.Keys(explained)) {
		boldGreenP.Println(reason)
		for _, path := range slices.Sorted(maps.Keys(explained[reason])) {
			blueP.Printf("%s- %s (%d entries)\n", prefixSpace, path, explained[reason][path])
		}
		fmt.Println()
	}

	boldYellowP.Println("Unexplained")
	for _, f := range unexplained {
		blueP.Printf("%s- %s\n", prefixSpace, f)
	}
	return nil
}
