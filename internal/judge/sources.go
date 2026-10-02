package judge

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type adoptionRow struct{ kind, name string }

// Only a table with the adoption contract's header grants source ownership.
func adoptionRows(dir string) ([]adoptionRow, error) {
	if _, err := os.Lstat(dir); errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if !groupingDirectory(dir) {
		return nil, errors.New("not a regular file in its directory")
	}
	text, err := readRegular(filepath.Join(dir, "_index.md"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rows []adoptionRow
	active := false
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			active = false
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		if len(cells) != 5 {
			active = false
			continue
		}
		if strings.Join(cells, "|") == "source|type|owner|adopted date|path" {
			active = true
			continue
		}
		if !active || strings.HasPrefix(cells[0], "---") {
			continue
		}
		name := strings.Trim(cells[4], "`")
		if strings.HasPrefix(name, "[") {
			_, dest, ok := strings.Cut(name, "](")
			if ok {
				name = strings.TrimSuffix(dest, ")")
			}
		}
		name = strings.TrimPrefix(name, "./")
		// A repository-relative path is accepted only when it names this directory.
		if strings.HasPrefix(name, "docs/") {
			suffix := filepath.ToSlash(filepath.Clean(dir))
			prefix := strings.TrimSuffix(name, "/"+filepath.Base(name))
			if strings.HasSuffix(suffix, "/"+prefix) {
				name = filepath.Base(name)
			}
		}
		rows = append(rows, adoptionRow{cells[1], name})
	}
	return rows, nil
}

func groupingDirectory(dir string) bool {
	// Reject linked directory components, including a linked repository root.
	for path := filepath.Clean(dir); ; path = filepath.Dir(path) {
		info, err := os.Lstat(path)
		if err != nil || !info.IsDir() {
			return false
		}
		parent := filepath.Dir(path)
		if filepath.Base(path) == "docs" {
			root, err := os.Lstat(parent)
			return err == nil && root.IsDir()
		}
		if parent == path {
			break
		}
	}
	return true
}

// readGroupingSource is the third bounded constructor for Source.
func readGroupingSource(dir, name string) (src Source, ok bool, err error) {
	if name == "_index.md" || filepath.Base(name) != name || !strings.HasSuffix(name, ".md") || !groupingDirectory(dir) {
		return Source{}, false, errors.New("not a regular file in its directory")
	}
	q, err := Load()
	if err != nil {
		return Source{}, false, err
	}
	adopted := filepath.Base(dir) == "references"
	if adopted {
		rows, err := adoptionRows(dir)
		if err != nil {
			return Source{}, false, err
		}
		accepted := false
		for _, row := range rows {
			if row.name == name && slices.Contains(q.Grouping.AdoptedSourceTypes, row.kind) {
				accepted = true
				break
			}
		}
		if !accepted {
			return Source{}, false, errors.New("only Findings and Backlog Entries are sent")
		}
	} else if filepath.Base(filepath.Dir(dir)) != "docs" || (filepath.Base(dir) != "backlog" && filepath.Base(dir) != "findings") {
		return Source{}, false, errors.New("only Findings and Backlog Entries are sent")
	}
	path := filepath.Join(dir, name)
	text, err := readRegular(path)
	if err != nil {
		return Source{}, false, errors.New("not a regular file in its directory")
	}
	if !adopted {
		front, _, _, err := splitFrontMatter(text)
		if err != nil || front == "" {
			return Source{}, false, nil
		}
		var lifecycle struct {
			Status string `yaml:"status"`
		}
		if yaml.Unmarshal([]byte(front), &lifecycle) != nil {
			return Source{}, false, nil
		}
		statuses := q.Grouping.OpenBacklogStatuses
		if filepath.Base(dir) == "findings" {
			statuses = q.Grouping.UnresolvedFindingStatuses
		}
		if !slices.Contains(statuses, lifecycle.Status) {
			return Source{}, false, nil
		}
	}
	if !q.Language.isEnglish(text) {
		return Source{}, false, errors.New("not English")
	}
	return Source{path: path, text: text}, true, nil
}

func collectGroupingSource(dir, name string, sources *[]Source, skipped *[]SkippedArtifact) {
	src, ok, err := readGroupingSource(dir, name)
	if err != nil {
		*skipped = append(*skipped, SkippedArtifact{filepath.Join(dir, name), err.Error()})
	}
	if ok {
		*sources = append(*sources, src)
	}
}

func adoptedSources(q Questions, specDir string) ([]Source, []SkippedArtifact, error) {
	dir := filepath.Join(specDir, "references")
	rows, err := adoptionRows(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read adopted sources: %w", err)
	}
	var sources []Source
	var skipped []SkippedArtifact
	for _, row := range rows {
		if !slices.Contains(q.Grouping.AdoptedSourceTypes, row.kind) {
			skipped = append(skipped, SkippedArtifact{filepath.Join(dir, row.name), "only Findings and Backlog Entries are sent"})
			continue
		}
		collectGroupingSource(dir, row.name, &sources, &skipped)
	}
	return sources, skipped, nil
}

func openSources(q Questions, repoRoot string) ([]Source, []SkippedArtifact, error) {
	var sources []Source
	var skipped []SkippedArtifact
	for _, family := range []string{"backlog", "findings"} {
		dir := filepath.Join(repoRoot, "docs", family)
		if _, err := os.Lstat(dir); errors.Is(err, os.ErrNotExist) {
			continue
		}
		if !groupingDirectory(dir) {
			skipped = append(skipped, SkippedArtifact{dir, "not a regular file in its directory"})
			continue
		}
		entries, err := os.ReadDir(dir)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("list grouping sources: %w", err)
		}
		for _, entry := range entries {
			collectGroupingSource(dir, entry.Name(), &sources, &skipped)
		}
	}
	return sources, skipped, nil
}

func prepareSource(q Questions, src Source) string {
	text := src.text
	if strings.HasPrefix(text, "---\n") {
		if end := strings.Index(text[3:], "\n---\n"); end >= 0 {
			text = text[3+end+5:]
		}
	}
	return cut(q.Grouping.SourceScrub.ReplaceAllString(text, ""), q.Grouping.SourceMaxChars)
}
