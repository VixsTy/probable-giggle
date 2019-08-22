// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	BinariesName = "probable-giggle"
	Default      = Build
	goFiles      = getGoFiles()
	goSrcFiles   = getGoSrcFiles()
)

var curDir = func() string {
	name, _ := os.Getwd()
	return name
}()

// Calculate file paths
var toolsBinDir = normalizePath(path.Join(curDir, "tools", "bin"))

func init() {
	time.Local = time.UTC

	// Add local bin in PATH
	err := os.Setenv("PATH", fmt.Sprintf("%s:%s", toolsBinDir, os.Getenv("PATH")))
	if err != nil {
		panic(err)
	}
}

func Build() {
	banner := figure.NewFigure(BinariesName, "", true)
	banner.Print()

	fmt.Println("")
	color.Red("# Build Info ---------------------------------------------------------------")
	fmt.Printf("Go version : %s\n", runtime.Version())
	fmt.Printf("Git revision : %s\n", hash())
	fmt.Printf("Git branch : %s\n", branch())
	fmt.Printf("Tag : %s\n", tag())

	fmt.Println("")

	color.Red("# Core packages ------------------------------------------------------------")
	mg.SerialDeps(Go.Deps, Go.License, Go.Generate, Go.Format, Go.Lint, Go.Test)

	fmt.Println("")
	color.Red("# Artifacts ----------------------------------------------------------------")
	mg.Deps(Bin.Ndsbuildingblock)
}

// -----------------------------------------------------------------------------

type Ci mg.Namespace

// Validate circleci configuration file (circleci/config.yml).
func (Ci) Validate() error {
	return sh.RunV("circleci-cli", "config", "validate")
}

// execute circleci job build on local.
func (ci Ci) Build() error {
	return ci.localExecute("build")
}

func (ci Ci) localExecute(job string) error {
	return sh.RunV("circleci-cli", "local", "execute", "--job", job)
}

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

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// branch returns the git branch for current repo
func branch() string {
	hash, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	return hash
}

func mustStr(r string, err error) string {
	if err != nil {
		panic(err)
	}
	return r
}

// normalizePath turns a path into an absolute path and removes symlinks
func normalizePath(name string) string {
	absPath := mustStr(filepath.Abs(name))
	return absPath
}
