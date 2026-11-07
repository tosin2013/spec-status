package status

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProductCodeDoesNotImportSideEffectPackages(t *testing.T) {
	root := filepath.Join("..", "..")
	banned := map[string]bool{"os/exec": true, "net": true, "net/http": true, "plugin": true}
	allowedCLI := map[string]bool{
		"github.com/spf13/cobra": true,
		"github.com/spf13/pflag": true,
		"golang.org/x/sys/unix":  true,
	}
	allowedStyle := map[string]bool{
		"github.com/charmbracelet/lipgloss": true,
		"github.com/muesli/termenv":         true,
		"github.com/rivo/uniseg":            true,
	}
	bannedNames := []string{"bubbletea", "glamour", "lipgloss/tree", "huh"}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		cliSurface := strings.HasPrefix(rel, "internal/cli/") || strings.HasPrefix(rel, "cmd/")
		styleSurface := rel == "internal/status/style.go"
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if banned[p] || strings.HasPrefix(p, "net/") {
				t.Errorf("%s imports %s", path, p)
			}
			for _, name := range bannedNames {
				if strings.Contains(p, name) {
					t.Errorf("%s imports %s", rel, p)
				}
			}
			if !strings.Contains(p, ".") {
				continue
			}
			if strings.HasPrefix(p, "github.com/tosin2013/spec-status/") {
				continue
			}
			switch {
			case cliSurface && allowedCLI[p]:
			case styleSurface && allowedStyle[p]:
			default:
				t.Errorf("%s imports dependency %s outside the approved command or style exception", rel, p)
			}
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(body)
		for _, needle := range []string{"os.WriteFile", "os.Create", "os.OpenFile", "os.Remove", "os.Mkdir", "os.Rename", "exec.Command", "http.Get"} {
			if strings.Contains(text, needle) {
				t.Errorf("%s contains %s", path, needle)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
