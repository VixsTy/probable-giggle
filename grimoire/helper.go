package grimoire

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

var (
	goFiles    = getGoFiles()
	goSrcFiles = getGoSrcFiles()
)

// -----------------------------------------------------------------------------

func getGoFiles() []string {
	var goFiles []string

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "vendor/") {
			return filepath.SkipDir
		}
		if strings.Contains(path, "tools/") {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		goFiles = append(goFiles, path)
		return nil
	})

	return goFiles
}

func getGoSrcFiles() []string {
	var goSrcFiles []string

	for _, path := range goFiles {
		if !strings.HasSuffix(path, "_test.go") {
			continue
		}

		goSrcFiles = append(goSrcFiles, path)
	}

	return goSrcFiles
}

// Tag returns the git tag for the current branch or "" if none.
func Tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// Hash returns the git hash for the current repo or "" if none.
func Hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// Branch returns the git branch for current repo
func Branch() string {
	hash, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	return hash
}

func mustStr(r string, err error) string {
	if err != nil {
		panic(err)
	}
	return r
}

func mustGoGenerate(txt, name string) {
	fmt.Printf(" > %s [%s]\n", txt, name)
	err := sh.RunV("go", "generate", name)
	if err != nil {
		panic(err)
	}
}

func runIntegrationTest(txt, name string) {
	fmt.Printf(" > %s [%s]\n", txt, name)
	err := sh.RunV("gotestsum", "--junitfile", fmt.Sprintf("test-results/junit/integration-%s.xml", strings.ToLower(txt)), name, "--", "-tags=integration", "-race")
	if err != nil {
		panic(err)
	}
}

// NormalizePath turns a path into an absolute path and removes symlinks
func NormalizePath(name string) string {
	absPath := mustStr(filepath.Abs(name))
	return absPath
}
