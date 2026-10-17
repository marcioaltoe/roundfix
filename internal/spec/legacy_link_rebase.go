package spec

import "strings"

// LegacyLinkRebase preserves targets when text leaves a removed Spec folder.
// All paths and promotion keys are repository-relative and slash-separated.
type LegacyLinkRebase struct {
	Folder, RecordPath string
	Files              []string
	Promoted           map[string]string
	Locate             func(path string) (GitLocator, bool, error)
}

// Rebase points at promoted files and locates every other removed folder path.
func (r LegacyLinkRebase) Rebase(content []byte, from, to string) ([]byte, error) {
	locate := func(path string) (GitLocator, bool, error) {
		if path == r.RecordPath {
			path = strings.TrimSuffix(r.Folder, "/")
		}
		if path != strings.TrimSuffix(r.Folder, "/") && !archivePathUnder(path, r.Folder) {
			return GitLocator{}, false, nil
		}
		if r.Locate != nil {
			return r.Locate(path)
		}
		return GitLocator{Path: path}, true, nil
	}
	rewritten, _, err := (CitationRewrite{Moves: r.Promoted, Locate: locate}).Apply(content, from, to)
	return rewritten, err
}
