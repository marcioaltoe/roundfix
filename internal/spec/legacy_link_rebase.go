package spec

import (
	"bytes"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

// LegacyLinkRebase names the surviving files of a legacy folder conversion.
// All paths and promotion keys are repository-relative and slash-separated.
type LegacyLinkRebase struct {
	Folder, RecordPath string
	Files              []string
	Promoted           map[string]string
}

// Rebase keeps relative destinations naming the same target after text moves,
// redirecting removed folder targets to their promoted copy or the record.
// Bounded path text uses the Archive Relink pass, including text inside code.
func (r LegacyLinkRebase) Rebase(content []byte, from, to string) []byte {
	moves := ArchiveCitationMoves(r.Folder, r.RecordPath, r.Files, r.Promoted)
	sources := make([]string, 0, len(moves))
	for source := range moves {
		if source != "" {
			sources = append(sources, source)
		}
	}
	sort.Slice(sources, func(i, j int) bool {
		if len(sources[i]) != len(sources[j]) {
			return len(sources[i]) > len(sources[j])
		}
		return sources[i] < sources[j]
	})
	var out bytes.Buffer
	end := 0
	for _, span := range scanArchiveMarkdown(content) {
		text, _ := relinkArchivePaths(content[end:span.start], sources, moves)
		out.Write(text)
		destination := content[span.start:span.end]
		target, suffix, relative := archiveLinkTarget(string(destination), filepath.Dir(from))
		target = filepath.ToSlash(target)
		if relative && !filepath.IsAbs(target) && target != ".." && !strings.HasPrefix(target, "../") {
			if target == strings.TrimSuffix(r.Folder, "/") || archivePathUnder(target, r.Folder) {
				if promoted, ok := r.Promoted[target]; ok {
					target = promoted
				} else {
					target = r.RecordPath
				}
			}
			if rel, err := filepath.Rel(filepath.Dir(to), filepath.FromSlash(target)); err == nil {
				out.WriteString((&url.URL{Path: filepath.ToSlash(rel)}).EscapedPath() + suffix)
			} else {
				out.Write(destination)
			}
		} else {
			out.Write(destination)
		}
		end = span.end
	}
	text, _ := relinkArchivePaths(content[end:], sources, moves)
	out.Write(text)
	return out.Bytes()
}
