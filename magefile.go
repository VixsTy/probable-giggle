// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	// mage:import
	"github.com/VixsTy/probable-giggle/grimoire"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

var (
	Default = Build
)

var curDir = func() string {
	name, _ := os.Getwd()
	return name
}()

// Calculate file paths
var toolsBinDir = grimoire.NormalizePath(path.Join(curDir, "tools", "bin"))

func init() {
	time.Local = time.UTC

	// Add local bin in PATH
	err := os.Setenv("PATH", fmt.Sprintf("%s:%s", toolsBinDir, os.Getenv("PATH")))
	if err != nil {
		panic(err)
	}
}

func Build() {
	banner := figure.NewFigure("probable-giggle", "", true)
	banner.Print()

	fmt.Println("")
	color.Red("# Build Info ---------------------------------------------------------------")
	fmt.Printf("Go version : %s\n", runtime.Version())
	fmt.Printf("Git revision : %s\n", grimoire.Hash())
	fmt.Printf("Git branch : %s\n", grimoire.Branch())
	fmt.Printf("Tag : %s\n", grimoire.Tag())

	fmt.Println("")

	color.Red("# Core packages ------------------------------------------------------------")
	mg.SerialDeps(grimoire.Go.Deps, grimoire.Go.License, grimoire.Go.Generate, grimoire.Go.Format, grimoire.Go.Lint, grimoire.Go.Test)

	fmt.Println("")
	color.Red("# Artifacts ----------------------------------------------------------------")
	mg.Deps(grimoire.Bin.ProbableGiggle)
}
