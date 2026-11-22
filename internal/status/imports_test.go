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
		styleSurface := rel == "internal/human/style.go"
		// The public library and its schema stay standard-library-only
		// (feature 012); they may import this module's internal packages.
		libSurface := strings.HasPrefix(rel, "internal/schema/") || strings.HasPrefix(rel, "pkg/status/")
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
				if libSurface && !strings.HasPrefix(p, "github.com/tosin2013/spec-status/internal/") {
					t.Errorf("%s imports %s; the library may import only internal packages of this module", rel, p)
				}
				if libSurface && (p == "github.com/tosin2013/spec-status/internal/cli" || p == "github.com/tosin2013/spec-status/internal/human") {
					t.Errorf("%s imports %s, which links third-party modules", rel, p)
				}
				continue
			}
			switch {
			case cliSurface && allowedCLI[p]:
			case styleSurface && allowedStyle[p]:
			case libSurface:
				t.Errorf("%s imports dependency %s; the library is standard library only", rel, p)
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

// TestLibraryDependencyClosureIsStandardLibrary checks that pkg/status and
// everything it imports from this module (transitively) import only the
// standard library and this module, so library users link no third-party
// code (feature 012, SC-005).
func TestLibraryDependencyClosureIsStandardLibrary(t *testing.T) {
	const module = "github.com/tosin2013/spec-status/"
	root := filepath.Join("..", "..")
	seen := map[string]bool{}
	var visit func(pkg string)
	visit = func(pkg string) {
		if seen[pkg] {
			return
		}
		seen[pkg] = true
		dir := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(pkg, module)))
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("%s: %v", pkg, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			f, err := parser.ParseFile(token.NewFileSet(), filepath.Join(dir, name), nil, parser.ImportsOnly)
			if err != nil {
				t.Fatal(err)
			}
			for _, imp := range f.Imports {
				p := strings.Trim(imp.Path.Value, `"`)
				switch {
				case strings.HasPrefix(p, module):
					visit(p)
				case strings.Contains(strings.SplitN(p, "/", 2)[0], "."):
					t.Errorf("%s/%s imports %s; pkg/status must link only the standard library and this module", pkg, name, p)
				}
			}
		}
	}
	visit(module + "pkg/status")
	if !seen[module+"internal/status"] || !seen[module+"internal/schema"] {
		t.Fatalf("walk did not reach the core: %v", seen)
	}
}
