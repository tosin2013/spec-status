package status

import (
	"os"
	"path/filepath"

	"github.com/tosin2013/spec-status/internal/safeio"
)

// Option selects one read-only evaluation.
type Option struct {
	Root    string
	Roadmap string
	FailOn  string
	CWD     string
}

// Evaluate reads the roadmap and returns one report. It does not write.
func Evaluate(opt Option) Report {
	failOn := opt.FailOn
	if failOn == "" {
		failOn = "error"
	}
	root, _, ferr := resolveRoot(opt.CWD, opt.Root)
	if ferr != nil {
		r := fatalReport(ferr.Code, ferr.Message, "", failOn)
		return r
	}
	fs, err := safeio.New(root)
	if err != nil {
		return fatalReport("SS_ROOT_UNAVAILABLE", "Selected root is unavailable.", "", failOn)
	}
	specEnt, err := fs.StatRel(".specify")
	if err != nil || !specEnt.Exists || !specEnt.Dir {
		return fatalReport("SS_ROOT_UNAVAILABLE", "Selected root has no .specify directory.", "", failOn)
	}
	roadmap := ".specify/memory/roadmap.md"
	if opt.Roadmap != "" {
		var rerr error
		roadmap, rerr = fs.Relativize(opt.Roadmap)
		if rerr != nil {
			return fatalReport("SS_PATH_UNSAFE", "Roadmap path is outside the selected root.", "", failOn)
		}
	}
	if safeio.Sensitive(roadmap) {
		return fatalReport("SS_REFERENCE_SENSITIVE_SKIPPED", "Roadmap path matched the sensitive-name matcher and was not read.", "", failOn)
	}
	data, err := fs.ReadRel(roadmap)
	if err != nil {
		if se, ok := err.(*safeio.Error); ok {
			switch se.Code {
			case safeio.CodeNotFound:
				return fatalReport("SS_ROADMAP_UNAVAILABLE", "Required roadmap is unavailable.", roadmap, failOn)
			case safeio.CodeSensitive:
				return fatalReport("SS_REFERENCE_SENSITIVE_SKIPPED", "Roadmap path matched the sensitive-name matcher and was not read.", "", failOn)
			case safeio.CodeLimit:
				return fatalReport("SS_INPUT_LIMIT", "Input budget exceeded.", roadmap, failOn)
			case safeio.CodeChanged:
				return fatalReport("SS_SOURCE_CHANGED", "Input changed while it was being read.", roadmap, failOn)
			case safeio.CodeUnsafe:
				return fatalReport("SS_PATH_UNSAFE", "Roadmap path is outside the selected root or uses a denied symlink.", "", failOn)
			default:
				return fatalReport("SS_ROADMAP_UNAVAILABLE", "Required roadmap is unavailable.", roadmap, failOn)
			}
		}
		return fatalReport("SS_ROADMAP_UNAVAILABLE", "Required roadmap is unavailable.", roadmap, failOn)
	}
	items, _, perr := parseRoadmap(roadmap, data)
	if perr != nil {
		r := fatalReport(perr.Code, perr.Message, roadmap, failOn)
		if perr.Line > 0 {
			r.Diagnostics[0].Line = perr.Line
			r.Fatal.Line = perr.Line
		}
		return r
	}
	rep := build(roadmap, items, fs)
	if rep.Outcome == "fatal" {
		if rep.FailOn == "" {
			rep.FailOn = failOn
		}
		if rep.Roadmap == "" {
			rep.Roadmap = roadmap
		}
		return rep
	}
	rep.FailOn = failOn
	return rep
}

type rootErr struct {
	Code, Message string
}

func resolveRoot(cwd, explicit string) (string, string, *rootErr) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return "", "", &rootErr{"SS_ROOT_UNAVAILABLE", "Selected root is unavailable."}
		}
	}
	var candidate string
	if explicit != "" {
		if filepath.IsAbs(explicit) {
			candidate = explicit
		} else {
			candidate = filepath.Join(cwd, explicit)
		}
		info, err := os.Stat(candidate)
		if err != nil || !info.IsDir() {
			return "", "", &rootErr{"SS_ROOT_UNAVAILABLE", "Selected root is unavailable."}
		}
		return candidate, explicit, nil
	}
	dir, err := filepath.Abs(cwd)
	if err != nil {
		return "", "", &rootErr{"SS_ROOT_UNAVAILABLE", "Selected root is unavailable."}
	}
	for {
		info, statErr := os.Lstat(filepath.Join(dir, ".specify"))
		if statErr == nil {
			if info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
				return dir, "", nil
			}
			// A .specify that is a symlink or a file is the nearest project
			// marker. Stop here instead of silently reporting an outer project.
			return "", "", &rootErr{"SS_ROOT_UNAVAILABLE", "Nearest .specify is not a usable directory."}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", "", &rootErr{"SS_ROOT_UNAVAILABLE", "Selected root is unavailable."}
}

// ExitCode applies the threshold. Fatal reports always return 2.
func ExitCode(rep Report) int {
	if rep.Outcome == "fatal" || !rep.Complete {
		return 2
	}
	if rep.FailOn == "none" {
		return 0
	}
	need := map[string]int{"info": 1, "warning": 2, "error": 3}[rep.FailOn]
	rank := map[string]int{"info": 1, "warning": 2, "error": 3}
	for _, d := range rep.Diagnostics {
		if rank[d.Severity] >= need && need > 0 {
			return 1
		}
	}
	return 0
}
