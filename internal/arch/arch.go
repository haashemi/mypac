// Package arch queries the system through pacman and lostfiles.
package arch

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// ExplicitPackages returns the output of `pacman -Qqe`.
func ExplicitPackages() ([]string, error) {
	out, err := exec.Command("pacman", "-Qqe").Output()
	if err != nil {
		return nil, err
	}
	return lines(out), nil
}

type ModifiedFile struct {
	Package string
	Path    string
}

// ModifiedFiles returns the package backup files (mostly under /etc) that
// differ from their packaged version, as reported by `pacman -Qii`.
func ModifiedFiles() ([]ModifiedFile, error) {
	cmd := exec.Command("pacman", "-Qii")
	cmd.Env = append(os.Environ(), "LC_ALL=C") // parsed output must not be localized
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseModified(out), nil
}

func parseModified(out []byte) []ModifiedFile {
	var pkg string
	var files []ModifiedFile
	for _, line := range lines(out) {
		if rest, ok := strings.CutPrefix(line, "Name "); ok {
			_, name, _ := strings.Cut(rest, ":")
			pkg = strings.TrimSpace(name)
		} else if rest, ok := strings.CutSuffix(line, " [modified]"); ok {
			if i := strings.IndexByte(rest, '/'); i >= 0 {
				files = append(files, ModifiedFile{Package: pkg, Path: rest[i:]})
			}
		}
	}
	return files
}

// LostFiles runs `lostfiles` (which requires root) and returns the files not
// owned by any package.
func LostFiles(strict bool) ([]string, error) {
	var args []string
	if strict {
		args = append(args, "-s")
	}

	cmd := exec.Command("lostfiles", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return lines(out), nil
}

func lines(out []byte) []string {
	var result []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			result = append(result, line)
		}
	}
	return result
}
