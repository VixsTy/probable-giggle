package grimoire

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Bin is a mage namespace which will manage binaries build actions
type Bin mg.Namespace

// ProbableGiggle build the ProbableGiggle binarie
func (Bin) ProbableGiggle() error {
	return goBuild("github.com/VixsTy/probable-giggle/cli/probablegiggle", "probablegiggle")
}

func goBuild(packageName, out string) error {
	fmt.Printf(" > Building %s [%s]\n", out, packageName)

	varsSetByLinker := map[string]string{
		"github.com/VixsTy/probable-giggle/internal/version.Version":   Tag(),
		"github.com/VixsTy/probable-giggle/internal/version.Revision":  Hash(),
		"github.com/VixsTy/probable-giggle/internal/version.Branch":    Branch(),
		"github.com/VixsTy/probable-giggle/internal/version.BuildUser": os.Getenv("USER"),
		"github.com/VixsTy/probable-giggle/internal/version.BuildDate": time.Now().Format(time.RFC3339),
		"github.com/VixsTy/probable-giggle/internal/version.GoVersion": runtime.Version(),
	}
	var linkerArgs []string
	for name, value := range varsSetByLinker {
		linkerArgs = append(linkerArgs, "-X", fmt.Sprintf("%s=%s", name, value))
	}
	linkerArgs = append(linkerArgs, "-s", "-w")

	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
	}, "go", "build", "-ldflags", strings.Join(linkerArgs, " "), "-mod=vendor", "-o", fmt.Sprintf("bin/%s", out), packageName)
}
