package cli

import (
	"runtime/debug"
	"testing"
)

// Build identity (feature 004, contracts/build-identity.md).

func info(version, sum string) *debug.BuildInfo {
	return &debug.BuildInfo{Main: debug.Module{Path: "github.com/tosin2013/spec-status", Version: version, Sum: sum}}
}

func TestResolveBuildKinds(t *testing.T) {
	const sha = "0123456789abcdef0123456789abcdef01234567"
	cases := []struct {
		name                  string
		kind, version, commit string
		bi                    *debug.BuildInfo
		ok                    bool
		wantKind, wantVersion string
		wantCommit            string
	}{
		{"release stamp wins over module info", "release", "0.1.0", sha, info("v9.9.9", "h1:x="), true, "release", "0.1.0", sha},
		{"release candidate", "release", "0.1.0-rc.1", sha, nil, false, "release", "0.1.0-rc.1", sha},
		{"go install tagged module", "", "0.0.0-dev", "unknown", info("v0.1.0", "h1:abc="), true, "module", "0.1.0", "unknown"},
		{"go install pseudo-version", "", "0.0.0-dev", "unknown", info("v0.0.0-20261112040751-b228f6020c7b", "h1:abc="), true, "module", "0.0.0-20261112040751-b228f6020c7b", "unknown"},
		{"checkout build with VCS pseudo-version", "", "0.0.0-dev", "unknown", info("v0.0.0-20261112040751-b228f6020c7b+dirty", ""), true, "development", "0.0.0-dev", "unknown"},
		{"clean checkout of a tag", "", "0.0.0-dev", "unknown", info("v0.1.0", ""), true, "development", "0.0.0-dev", "unknown"},
		{"devel", "", "0.0.0-dev", "unknown", info("(devel)", "h1:abc="), true, "development", "0.0.0-dev", "unknown"},
		{"no build info", "", "0.0.0-dev", "unknown", nil, false, "development", "0.0.0-dev", "unknown"},
		{"stamped commit without release kind stays development", "", "0.0.0-dev", sha, nil, false, "development", "0.0.0-dev", sha},
	}
	for _, c := range cases {
		got := resolveBuild(c.kind, c.version, c.commit, c.bi, c.ok)
		if got.kind != c.wantKind || got.version != c.wantVersion || got.commit != c.wantCommit {
			t.Errorf("%s: got %s %s %s, want %s %s %s", c.name, got.kind, got.version, got.commit, c.wantKind, c.wantVersion, c.wantCommit)
		}
	}
}

func TestBuildLabels(t *testing.T) {
	rel := resolveBuild("release", "0.1.0", "c", nil, false)
	mod := resolveBuild("", "0.0.0-dev", "unknown", info("v0.1.0", "h1:a="), true)
	dev := resolveBuild("", "0.0.0-dev", "unknown", nil, false)
	for _, c := range []struct {
		b                    buildIdentity
		label, heading, help string
	}{
		{rel, "release build", "release 0.1.0", "This is release 0.1.0."},
		{mod, "module build (built from the module version, not a release archive)", "module build 0.1.0", "This is a module build of 0.1.0, not a release archive."},
		{dev, "development build, not a release", "development build", "This is a development build."},
	} {
		if c.b.label != c.label || c.b.heading != c.heading || c.b.helpSentence != c.help {
			t.Errorf("%s: got %q %q %q", c.b.kind, c.b.label, c.b.heading, c.b.helpSentence)
		}
	}
}
