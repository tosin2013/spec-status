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
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, imp := range f.Imports {
			p := strings.Trim(imp.Path.Value, `"`)
			if banned[p] || strings.HasPrefix(p, "net/") {
				t.Errorf("%s imports %s", path, p)
			}
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(body)
		for _, needle := range []string{"os.WriteFile", "os.Create", "os.Remove", "os.Mkdir", "os.Rename", "exec.Command", "http.Get"} {
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
