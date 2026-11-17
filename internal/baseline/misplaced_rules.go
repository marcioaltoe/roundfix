package baseline

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// SpecificRepositoryPath is the repository-owned destination for local rules.
const SpecificRepositoryPath = specificRepositoryPath

// MisplacedRepositoryRule names a passage without moving its bytes.
type MisplacedRepositoryRule struct {
	Path      string `json:"path"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	FirstLine string `json:"firstLine"`
	ManagedID string `json:"managedId,omitempty"`
	MoveTo    string `json:"moveTo"`
}

// MisplacedRepositoryRuleReport describes the carriers of an adopted repository.
type MisplacedRepositoryRuleReport struct {
	ManifestFound    bool
	ExtensionEnabled bool
	Rules            []MisplacedRepositoryRule
}

// MisplacedRepositoryRules reads the Setup Manifest and names local passages.
// It never moves rules or changes any repository file.
func MisplacedRepositoryRules(root string) (MisplacedRepositoryRuleReport, error) {
	var report MisplacedRepositoryRuleReport
	manifestBytes, err := readOptionalRegular(root, manifestPath)
	if err != nil {
		return report, fmt.Errorf("read misplaced-rule Setup Manifest: %w", err)
	}
	if manifestBytes == nil {
		return report, nil
	}
	report.ManifestFound = true
	manifest, valid := parseManagedSetupManifest(manifestBytes)
	if !valid {
		return report, fmt.Errorf("read misplaced-rule Setup Manifest: %w", errors.New(
			"baseline.preservation.managed-refresh.manifest.invalid",
		))
	}
	report.ExtensionEnabled = manifest.Decisions["repository.extension.enabled"].Value == true
	sourceByPath := manifestArtifactsByPath(manifest.ManagedArtifacts)
	contents := make(map[string][]byte)
	allCarriersEligible := true
	for _, relative := range sortedKeys(sourceByPath) {
		if !misplacedRuleCarrier(relative) {
			allCarriersEligible = false
			continue
		}
		if relative == "CLAUDE.md" {
			info, err := os.Lstat(filepath.Join(root, relative))
			if err != nil {
				return report, fmt.Errorf("inspect misplaced-rule carrier %q: %w", relative, err)
			}
			if !info.Mode().IsRegular() {
				allCarriersEligible = false
				continue
			}
		}
		content, err := readOptionalRegular(root, relative)
		if err != nil {
			return report, fmt.Errorf("read misplaced-rule carrier %q: %w", relative, err)
		}
		if content == nil {
			return report, fmt.Errorf("read misplaced-rule carrier %q: %w", relative, fs.ErrNotExist)
		}
		contents[relative] = content
	}
	// A regular root CLAUDE.md is read even when the manifest does not list
	// it. A symbolic link is skipped, so its target is reported only once.
	if _, listed := sourceByPath["CLAUDE.md"]; !listed {
		info, err := os.Lstat(filepath.Join(root, "CLAUDE.md"))
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return report, fmt.Errorf("inspect misplaced-rule carrier %q: %w", "CLAUDE.md", err)
		case info.Mode().IsRegular():
			content, err := readOptionalRegular(root, "CLAUDE.md")
			if err != nil {
				return report, fmt.Errorf("read misplaced-rule carrier %q: %w", "CLAUDE.md", err)
			}
			if content != nil {
				contents["CLAUDE.md"] = content
			}
		}
	}

	// Managed refresh reads every artifact. Use its classifier only when that
	// does not cross this detector's narrower read boundary. Otherwise the
	// same shared digest matcher operates on the eligible carriers below.
	var unrecorded []UnrecordedManagedRegion
	if allCarriersEligible {
		unrecorded, _, err = classifyManagedRegions(root)
		if err != nil {
			return report, fmt.Errorf("classify misplaced-rule carriers: %w", err)
		}
	}
	mismatches := make(map[string]map[string]bool)
	for _, region := range unrecorded {
		if region.Reason == UnrecordedManagedRegionReasonDigestMismatch {
			if mismatches[region.Path] == nil {
				mismatches[region.Path] = make(map[string]bool)
			}
			mismatches[region.Path][region.ManagedID] = true
		}
	}
	for _, relative := range sortedKeys(contents) {
		content := contents[relative]
		entries := partitionRootSource(relative, content)
		managedCounts := make(map[string]int)
		for _, entry := range entries {
			if entry.Kind == "managed-block" {
				id, _ := entry.StructuralProvenance["managedId"].(string)
				managedCounts[id]++
			}
		}
		for _, entry := range entries {
			if entry.Kind == "unmarked-span" {
				report.Rules = append(report.Rules, misplacedUnmarkedPassages(relative, content, entry)...)
				continue
			}
			id, _ := entry.StructuralProvenance["managedId"].(string)
			mismatch := mismatches[relative][id]
			if !allCarriersEligible {
				artifact, recorded := sourceByPath[relative][id]
				mismatch = recorded && managedCounts[id] == 1 && !managedEntryMatchesManifest(entry, artifact)
			}
			if !mismatch {
				continue
			}
			rule := misplacedRuleAt(relative, content, entry.Start, entry.End)
			rule.ManagedID = id
			body := entry.SourceBytes
			beginEnd := bytes.Index(body, []byte("-->")) + len("-->")
			endStart := bytes.LastIndex(body, []byte("<!--"))
			rule.FirstLine = ""
			for _, line := range strings.Split(string(body[beginEnd:endStart]), "\n") {
				if strings.TrimSpace(line) != "" {
					rule.FirstLine = misplacedFirstLine(line)
					break
				}
			}
			report.Rules = append(report.Rules, rule)
		}
	}
	sort.Slice(report.Rules, func(i, j int) bool {
		if report.Rules[i].Path != report.Rules[j].Path {
			return report.Rules[i].Path < report.Rules[j].Path
		}
		return report.Rules[i].StartLine < report.Rules[j].StartLine
	})
	return report, nil
}

func misplacedRuleCarrier(relative string) bool {
	if relative == SpecificRepositoryPath || path.Ext(relative) != ".md" {
		return false
	}
	// Nested instruction carriers have their own scope. Managed guides under
	// docs/agents remain eligible when their paths are in the manifest.
	base := path.Base(relative)
	return relative == base || base != "AGENTS.md" && base != "CLAUDE.md"
}

var misplacedForeignMarker = regexp.MustCompile(`^\s*<!--\s*(BEGIN|END):([^\s]+?)\s*-->\s*$`)

func misplacedUnmarkedPassages(relative string, content []byte, entry ReadoptionSourceEntry) []MisplacedRepositoryRule {
	lines := bytes.SplitAfter(entry.SourceBytes, []byte("\n"))
	offsets := make([]int, len(lines)+1)
	offsets[0] = entry.Start
	for i, line := range lines {
		offsets[i+1] = offsets[i] + len(line)
	}
	var rules []MisplacedRepositoryRule
	appendPassage := func(start, end int) {
		for start < end && len(bytes.TrimSpace(lines[start])) == 0 {
			start++
		}
		for start < end && len(bytes.TrimSpace(lines[end-1])) == 0 {
			end--
		}
		if start < end {
			rules = append(rules, misplacedRuleAt(relative, content, offsets[start], offsets[end]))
		}
	}
	start := 0
	for i := 0; i < len(lines); i++ {
		begin := misplacedForeignMarker.FindSubmatch(bytes.TrimRight(lines[i], "\r\n"))
		if begin == nil || string(begin[1]) != "BEGIN" {
			continue
		}
		for j := i + 1; j < len(lines); j++ {
			end := misplacedForeignMarker.FindSubmatch(bytes.TrimRight(lines[j], "\r\n"))
			if end == nil || string(end[1]) != "END" || !bytes.Equal(begin[2], end[2]) {
				continue
			}
			appendPassage(start, i)
			start = j + 1
			i = j
			break
		}
	}
	appendPassage(start, len(lines))
	return rules
}

func misplacedRuleAt(relative string, content []byte, start, end int) MisplacedRepositoryRule {
	// End is exclusive; do not count the newline following the last line.
	if end > start && content[end-1] == '\n' {
		end--
	}
	first := content[start:end]
	if newline := bytes.IndexByte(first, '\n'); newline >= 0 {
		first = first[:newline]
	}
	return MisplacedRepositoryRule{
		Path: relative, StartLine: bytes.Count(content[:start], []byte("\n")) + 1,
		EndLine:   bytes.Count(content[:end], []byte("\n")) + 1,
		FirstLine: misplacedFirstLine(string(first)), MoveTo: SpecificRepositoryPath,
	}
}

func misplacedFirstLine(line string) string {
	line = strings.TrimSpace(line)
	if len(line) <= 120 {
		return line
	}
	end := 120
	for !utf8.RuneStart(line[end]) {
		end--
	}
	return line[:end] + "..."
}
