package cli

import (
	"runtime/debug"
	"strings"
	"sync"
)

// buildIdentity says how this binary was built (feature 004,
// specs/004-release-engineering/contracts/build-identity.md).
//
//   - release: stamped by the release workflow with -X ...cli.Kind=release.
//   - module: built by `go install module@version` from a downloaded,
//     checksum-verified module (build info carries Main.Sum).
//   - development: anything else, including builds from a source checkout,
//     which Go stamps with a VCS-derived version but no module checksum.
type buildIdentity struct {
	kind         string
	version      string
	commit       string
	label        string
	heading      string
	helpSentence string
}

func resolveBuild(kind, version, commit string, info *debug.BuildInfo, ok bool) buildIdentity {
	if kind == "release" {
		return buildIdentity{
			kind: "release", version: version, commit: commit,
			label:        "release build",
			heading:      "release " + version,
			helpSentence: "This is release " + version + ".",
		}
	}
	if ok && info != nil && info.Main.Sum != "" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		v := strings.TrimPrefix(info.Main.Version, "v")
		return buildIdentity{
			kind: "module", version: v, commit: "unknown",
			label:        "module build (built from the module version, not a release archive)",
			heading:      "module build " + v,
			helpSentence: "This is a module build of " + v + ", not a release archive.",
		}
	}
	return buildIdentity{
		kind: "development", version: version, commit: commit,
		label:        "development build, not a release",
		heading:      "development build",
		helpSentence: "This is a development build.",
	}
}

var (
	buildOnce     sync.Once
	buildResolved buildIdentity
	// buildOverride lets tests inject an identity; nil in production.
	buildOverride *buildIdentity
)

func currentBuild() buildIdentity {
	if buildOverride != nil {
		return *buildOverride
	}
	buildOnce.Do(func() {
		info, ok := debug.ReadBuildInfo()
		buildResolved = resolveBuild(Kind, Version, Commit, info, ok)
	})
	return buildResolved
}
