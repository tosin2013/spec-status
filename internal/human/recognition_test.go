package human

import (
	"encoding/json"
	"github.com/tosin2013/spec-status/internal/status"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoadmapRecognitionFixtures(t *testing.T) {
	for _, name := range []string{
		"unrecognized-empty", "unrecognized-whitespace", "unrecognized-narrative",
		"planned-in-fence", "planned-wrong-form", "explicit-empty", "items-no-planned",
		"narrative-with-label", "no-roadmap",
	} {
		t.Run(name, func(t *testing.T) {
			root := twoModeRoot(name)
			raw, err := os.ReadFile(filepath.Join(root, "expected-status.json"))
			if err != nil {
				t.Fatal(err)
			}
			var want recognitionExpect
			if err := json.Unmarshal(raw, &want); err != nil {
				t.Fatal(err)
			}
			rep := status.Evaluate(status.Option{Root: root})
			if got := status.ExitCode(rep); got != want.Exit {
				t.Fatalf("exit %d, want %d (fatal %+v)", got, want.Exit, rep.Fatal)
			}
			var doc map[string]any
			if err := json.Unmarshal(status.JSON(rep), &doc); err != nil {
				t.Fatal(err)
			}
			if want.Code != "" {
				if rep.Fatal == nil || rep.Fatal.Code != want.Code {
					t.Fatalf("fatal %+v, want %s", rep.Fatal, want.Code)
				}
				if want.Message != nil && rep.Fatal.Message != *want.Message {
					t.Fatalf("message %q, want %q", rep.Fatal.Message, *want.Message)
				}
				errDoc, _ := doc["error"].(map[string]any)
				if errDoc == nil || errDoc["code"] != want.Code {
					t.Fatalf("JSON error %v", doc["error"])
				}
				text := Text(rep, View{Width: 80})
				if !strings.HasPrefix(text, "spec: "+want.Code+": ") {
					t.Fatalf("fatal text %q", text)
				}
				// A fatal is never masked by --fail-on none.
				none := status.Evaluate(status.Option{Root: root, FailOn: "none"})
				if status.ExitCode(none) != 2 {
					t.Fatal("--fail-on none masked a fatal")
				}
				return
			}
			if rep.Fatal != nil {
				t.Fatalf("unexpected fatal %+v", rep.Fatal)
			}
			summary, _ := doc["summary"].(map[string]any)
			if want.TotalItems != nil && int(summary["total_items"].(float64)) != *want.TotalItems {
				t.Fatalf("total_items %v, want %d", summary["total_items"], *want.TotalItems)
			}
			text := Text(rep, View{Width: 80})
			if want.ExplicitEmptyLine != nil && strings.Contains(text, explicitEmptyLine) != *want.ExplicitEmptyLine {
				t.Fatalf("explicit-empty line present=%v, want %v:\n%s", !*want.ExplicitEmptyLine, *want.ExplicitEmptyLine, text)
			}
		})
	}
}

func TestItemsHideExplicitEmptyLine(t *testing.T) {
	root := filepath.Join("..", "..", "testdata", "example-roadmap")
	text := Text(status.Evaluate(status.Option{Root: root}), View{Width: 80})
	if strings.Contains(text, explicitEmptyLine) {
		t.Fatal("explicit-empty line shown for a roadmap with items")
	}
}
