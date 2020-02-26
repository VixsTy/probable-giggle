// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"

	// mage:import
	grimoire "github.com/VixsTy/grimoire"
	figure "github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
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
	banner := figure.NewFigure(grimoire.MainDirectoryName(), "", true)
	banner.Print()

	fmt.Println("")
	color.Red("# Build Info ---------------------------------------------------------------")
	fmt.Printf("Go version : %s\n", runtime.Version())
	fmt.Printf("Git revision : %s\n", grimoire.Hash())
	fmt.Printf("Git branch : %s\n", grimoire.Branch())
	fmt.Printf("Tag : %s\n", grimoire.Tag())

	fmt.Println("")

	color.Red("# Core packages ------------------------------------------------------------")
	mg.SerialDeps(grimoire.Go.Deps, grimoire.Go.License, Go.Generate, grimoire.Go.Format, grimoire.Go.Lint, Go.Test)

	fmt.Println("")
	color.Red("# Artifacts ----------------------------------------------------------------")
	mg.Deps(Bin.ProbableGiggle)
}

// Bob is a mage namespace which will manage binaries build actions
type Bin mg.Namespace

// ProbableGiggle build the ProbableGiggle binarie
func (Bin) ProbableGiggle() error {
	return grimoire.Build{}.Binary("hello", "string2")
}

// Go is a mage namespace which will manage golang actions
type Go mg.Namespace

// Generate go code
func (Go) Generate() error {
	color.Cyan("## Generate code")
	// mg.SerialDeps(Gen.Protobuf, Gen.Mocks, Gen.Migrations, Gen.Decorators, Gen.Wire)
	mg.SerialDeps(Gen.Protobuf)
	return nil
}

// Test run all test
func (Go) Test() error {
	color.Cyan("## Running unit tests")
	sh.Run("mkdir", "-p", "test-results/junit")
	return sh.RunV("gotestsum", "--junitfile", "test-results/junit/unit-tests.xml", "--", "-short", "-race", "-cover", "./...")
}

// IntegrationTest run test tagged integration
func (Go) IntegrationTest() {
	color.Cyan("## Running integration tests")
	sh.Run("mkdir", "-p", "test-results/junit")

	grimoire.RunIntegrationTest("Repositories", "go.zenithar.org/spotigraph/internal/repositories/test/integration")
}

// Gen is a mage namespace which will manage generation actions
type Gen mg.Namespace

// Generate initializers
// func (Gen) Wire() {
// 	color.Blue("### Wiring dispatchers")
// 	mustGoGenerate("gRPC", "go.zenithar.org/spotigraph/cli/spotigraph/internal/dispatchers/grpc")
// }

// // Generate mocks for tests
// func (Gen) Mocks() {
// 	color.Blue("### Mocks")

// 	mustGoGenerate("Repositories", "go.zenithar.org/spotigraph/internal/repositories")
// 	mustGoGenerate("Services", "go.zenithar.org/spotigraph/internal/services")
// }

// Generate mocks for tests
// func (Gen) Decorators() {
// 	color.Blue("### Decorators")

// 	mustGoGenerate("Repositories", "go.zenithar.org/spotigraph/internal/repositories/pkg/...")
// 	mustGoGenerate("Services", "go.zenithar.org/spotigraph/internal/services/pkg/...")
// }

// // Generate initializers
// func (Gen) Migrations() {
// 	color.Blue("### Database migrations")

// 	mustGoGenerate("PostgreSQL", "go.zenithar.org/spotigraph/internal/repositories/pkg/postgresql")
// }

// Protobuf will Generate protobuf files
func (Gen) Protobuf() error {
	color.Blue("### Protobuf")

	return sh.RunV("prototool", "all", "--fix", "pkg/protocol")
}
