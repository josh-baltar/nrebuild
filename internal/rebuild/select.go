// Package rebuild decides which build targets must be re-run after a set of
// files changed since the last successful build.
package rebuild

import (
	"sort"
	"strings"
)

// Target is one buildable unit in the monorepo manifest.
type Target struct {
	// Name uniquely identifies the target.
	Name string
	// Dir is the directory the target owns; a changed file under Dir is one
	// of this target's inputs.
	Dir string
	// Generated lists path suffixes (relative to Dir) that this target writes
	// into its own directory, e.g. "gen/api.pb.go".
	Generated []string
	// Deps names the other targets this target consumes the outputs of.
	Deps []string
}

// SelectRebuilds returns the sorted names of the targets that must be rebuilt
// given changedPaths (repo-relative paths that differ from the last build).
//
// A target is selected when one of its input files changed. The result is
// deterministic (sorted) so CI can key a cache on it.
func SelectRebuilds(targets []Target, changedPaths []string) []string {
	byName := make(map[string]Target, len(targets))
	for _, t := range targets {
		byName[t.Name] = t
	}

	selected := make(map[string]bool)

	// Direct hits: a changed path under a target's Dir is an input change.
	for _, t := range targets {
		for _, p := range changedPaths {
			if underDir(p, t.Dir) {
				selected[t.Name] = true
				break
			}
		}
	}

	return sortedKeys(selected)
}

// underDir reports whether repo-relative path p lives under dir.
func underDir(p, dir string) bool {
	dir = strings.TrimSuffix(dir, "/")
	if dir == "" || dir == "." {
		return true
	}
	return p == dir || strings.HasPrefix(p, dir+"/")
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
