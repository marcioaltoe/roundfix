package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"

	"roundfix/internal/suiteguardcontract"
)

const ArchiveRecordSchema = "roundfix/archive-record/v1"
const ArchiveRecordTargetBytes = 2048

type ArchiveDisposition string

const (
	ArchivePass       ArchiveDisposition = "pass"
	ArchivePartial    ArchiveDisposition = "partial"
	ArchiveQAOverride ArchiveDisposition = "qa-override"
	ArchiveFailedQA   ArchiveDisposition = "failed-qa"
	ArchiveSuperseded ArchiveDisposition = "superseded"
	ArchiveNoQA       ArchiveDisposition = "no-qa"
)

type ArchiveRegeneration struct {
	Command string   `yaml:"command"`
	Outputs []string `yaml:"outputs,omitempty"`
}
type QAArchiveOverrideRecord struct{ Approval, Reason, QAOutcome, QATaskStatus, Revision string }
type ArchiveRecord struct {
	Spec, Title, Created, Archived string
	Disposition                    ArchiveDisposition
	Source, SourceRevision         string
	QATask, QAReport, QAVerdict    string
	Unproven                       []string
	QAOverride                     *QAArchiveOverrideRecord
	SupersededBy                   string
	ADRs, Sources, Promoted        []string
	Regeneration                   []ArchiveRegeneration
	PullRequest, DeliveryCommit    string
	Outcome                        string
}
type ArchiveRecordInput struct {
	Legacy                                          bool
	SpecDir, Slug, Source, SourceRevision, Archived string
	Unproven                                        []string
	QAOverride                                      *QAArchiveOverrideRecord
	Promoted                                        []string
}

func ArchiveRecordPath(root, slug string) string { return filepath.Join(root, slug+".md") }

// BuildArchiveRecord derives archive metadata without modifying any source file.
func BuildArchiveRecord(in ArchiveRecordInput) (ArchiveRecord, error) {
	content, err := os.ReadFile(filepath.Join(in.SpecDir, "_prd.md"))
	if err != nil {
		return ArchiveRecord{}, fmt.Errorf("read archive PRD: %w", err)
	}
	fm, body, err := splitFrontmatter(content)
	if err != nil {
		return ArchiveRecord{}, err
	}
	var meta struct {
		Status                              string `yaml:"status"`
		Spec, Created, Archived, SourceSlug string
		Unproven                            []string `yaml:"unproven"`
		Override                            bool     `yaml:"qa_override"`
		Approval                            string   `yaml:"qa_override_approval"`
		Reason                              string   `yaml:"qa_override_reason"`
		Outcome                             string   `yaml:"qa_override_qa_outcome"`
		TaskStatus                          string   `yaml:"qa_override_qa_task_status"`
		Revision                            string   `yaml:"qa_override_revision"`
	}
	if in.Legacy {
		err = readLegacyArchivePRD(fm, &meta)
	} else {
		if err = yaml.Unmarshal(fm, &meta); err != nil {
			err = fmt.Errorf("parse archive PRD: %w", err)
		}
	}
	if err != nil {
		return ArchiveRecord{}, err
	}
	r := ArchiveRecord{Spec: in.Slug, Created: meta.Created, Archived: in.Archived, Source: in.Source, SourceRevision: in.SourceRevision, Unproven: in.Unproven, QAOverride: in.QAOverride, Promoted: in.Promoted}
	if r.Spec == "" {
		r.Spec = meta.Spec
	}
	if r.Archived == "" {
		r.Archived = meta.Archived
	}
	if r.Unproven == nil {
		r.Unproven = meta.Unproven
	}
	if r.QAOverride == nil && meta.Override {
		r.QAOverride = &QAArchiveOverrideRecord{meta.Approval, meta.Reason, meta.Outcome, meta.TaskStatus, meta.Revision}
	}
	for _, line := range strings.Split(string(body), "\n") {
		if strings.HasPrefix(line, "# ") {
			r.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			break
		}
	}
	r.Outcome = archiveParagraph(string(body), "# ")
	seen := map[string]bool{}
	for _, adr := range regexp.MustCompile(`ADR-[0-9]{4}`).FindAllString(archiveSection(string(body), "Decisions"), -1) {
		if !seen[adr] {
			r.ADRs = append(r.ADRs, adr)
			seen[adr] = true
		}
	}
	index, err := os.ReadFile(filepath.Join(in.SpecDir, "references", "_index.md"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return r, err
	}
	column := -1
	for _, line := range strings.Split(string(index), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "|") {
			continue
		}
		cells := strings.Split(strings.Trim(strings.TrimSpace(line), "|"), "|")
		if column < 0 {
			for i, c := range cells {
				if strings.EqualFold(strings.TrimSpace(c), "path") {
					column = i
				}
			}
			continue
		}
		if column >= len(cells) {
			continue
		}
		value := strings.Trim(strings.TrimSpace(cells[column]), "`")
		if value != "" && !strings.HasPrefix(value, "-") && !strings.HasPrefix(value, ":") {
			r.Sources = append(r.Sources, value)
		}
	}
	auth, err := os.ReadFile(filepath.Join(in.SpecDir, "_authorization.md"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return r, err
	}
	for _, regen := range suiteguardcontract.ParseSanctionedRegenerations(auth) {
		r.Regeneration = append(r.Regeneration, ArchiveRegeneration{regen.Command, regen.Outputs})
	}
	var graph CauseGraph
	if in.Legacy {
		graph, _, err = ReadLegacyCauseGraph(filepath.Dir(in.SpecDir), filepath.Base(in.SpecDir))
	} else {
		graph, err = ReadCauseGraph(filepath.Dir(in.SpecDir), filepath.Base(in.SpecDir))
	}
	if err == nil {
		r.QATask = graph.QATaskID
	} else if _, statErr := os.Stat(filepath.Join(in.SpecDir, "_tasks.md")); !errors.Is(statErr, os.ErrNotExist) {
		return r, err
	}
	report, err := ReadQAReport(in.SpecDir)
	if err == nil {
		reportPath, pathErr := NewestQAReport(in.SpecDir)
		if pathErr != nil {
			return r, pathErr
		}
		r.QAReport = filepath.Base(reportPath)
		r.QAVerdict = report.Verdict
		r.Disposition = ArchiveDisposition(report.Verdict)
		if report.Verdict == VerdictFail {
			r.Disposition = ArchiveFailedQA
		}
		text, readErr := os.ReadFile(reportPath)
		if readErr != nil {
			return r, readErr
		}
		outcome := archiveParagraph(archiveSection(string(text), "Outcome"), "")
		if outcome != "" {
			r.Outcome = outcome
		}
	} else if r.QAOverride == nil && !errors.Is(err, ErrNoQAReport) {
		return r, err
	}
	if r.QAOverride != nil {
		r.Disposition = ArchiveQAOverride
	}
	if _, err := os.Stat(filepath.Join(in.SpecDir, "_tasks.md")); errors.Is(err, os.ErrNotExist) {
		s, err := ReadSupersession(in.SpecDir)
		if err != nil {
			if !errors.Is(err, ErrNoSupersession) {
				return r, err
			}
		} else {
			r.Disposition = ArchiveSuperseded
			r.SupersededBy = s.SupersededBy
			r.Outcome = strings.Join(strings.Fields(s.Reason), " ")
		}
	}
	if r.Disposition == "" {
		r.Disposition = ArchiveNoQA
	}
	return r, nil
}

// Decode legacy unproven separately so scalar spellings survive YAML typing.
// The remaining PRD metadata uses the same decoder as active Specs.
func readLegacyArchivePRD(fm []byte, meta any) error {
	var document yaml.Node
	if err := yaml.Unmarshal(fm, &document); err != nil {
		return fmt.Errorf("parse archive PRD: %w", err)
	}
	if len(document.Content) > 0 && document.Content[0].Kind == yaml.MappingNode {
		mapping := document.Content[0]
		for i := 0; i < len(mapping.Content); i += 2 {
			if mapping.Content[i].Value == "unproven" {
				value := mapping.Content[i+1]
				var items []yaml.Node
				if err := value.Decode(&items); err != nil {
					return fmt.Errorf("parse archive PRD: %w", err)
				}
				var lines []string
				if items != nil {
					lines = make([]string, len(items))
				}
				for j := range items {
					line, ok := legacyUnprovenLine(&items[j])
					if !ok {
						return fmt.Errorf("legacy unproven item %d cannot be read as text", j+1)
					}
					lines[j] = line
				}
				if err := value.Encode(lines); err != nil {
					return fmt.Errorf("encode legacy unproven text: %w", err)
				}
			}
		}
	}
	if err := document.Decode(meta); err != nil {
		return fmt.Errorf("parse archive PRD: %w", err)
	}
	return nil
}

func legacyUnprovenLine(item *yaml.Node) (string, bool) {
	if item.Kind == yaml.ScalarNode {
		return item.Value, true
	}
	if item.Kind != yaml.MappingNode || len(item.Content) == 0 {
		return "", false
	}
	values := make(map[string]string)
	var keys []string
	for i := 0; i < len(item.Content); i += 2 {
		key := item.Content[i]
		value, ok := legacyUnprovenValue(item.Content[i+1])
		if key.Kind != yaml.ScalarNode || !ok {
			return "", false
		}
		// A Node keeps duplicate keys that decoding into a map would refuse;
		// refuse them too rather than keep only the last value.
		if _, duplicate := values[key.Value]; duplicate {
			return "", false
		}
		values[key.Value] = value
		if key.Value != "row" && key.Value != "claim" {
			keys = append(keys, key.Value)
		}
	}
	sort.Strings(keys)
	var parts []string
	for _, key := range keys {
		parts = append(parts, key+": "+values[key])
	}
	line := values["claim"]
	if len(parts) > 0 {
		other := strings.Join(parts, "; ")
		if line != "" {
			line += " (" + other + ")"
		} else {
			line = other
		}
	}
	if row, exists := values["row"]; exists {
		line = "row " + row + ": " + line
	}
	return line, true
}

func legacyUnprovenValue(value *yaml.Node) (string, bool) {
	switch value.Kind {
	case yaml.ScalarNode:
		return strings.Join(strings.Fields(value.Value), " "), true
	case yaml.SequenceNode:
		var parts []string
		for _, scalar := range value.Content {
			if scalar.Kind != yaml.ScalarNode {
				return "", false
			}
			parts = append(parts, strings.Join(strings.Fields(scalar.Value), " "))
		}
		return strings.Join(parts, ", "), true
	default:
		return "", false
	}
}

func archiveSection(body, heading string) string {
	lines := strings.Split(body, "\n")
	start := -1
	for i, line := range lines {
		if line == "## "+heading {
			start = i + 1
			continue
		}
		if start >= 0 && strings.HasPrefix(line, "## ") {
			return strings.Join(lines[start:i], "\n")
		}
	}
	if start >= 0 {
		return strings.Join(lines[start:], "\n")
	}
	return ""
}
func archiveParagraph(body, after string) string {
	started := after == ""
	var lines []string
	for _, line := range strings.Split(body, "\n") {
		if !started {
			if strings.HasPrefix(line, after) {
				started = true
			}
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if len(lines) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(strings.Fields(strings.Join(lines, " ")), " ")
}
func archiveRecordMapping(r ArchiveRecord) *yaml.Node {
	n := &yaml.Node{Kind: yaml.MappingNode}
	add := func(k string, v any) {
		var value yaml.Node
		_ = value.Encode(v)
		n.Content = append(n.Content, &yaml.Node{Kind: yaml.ScalarNode, Value: k, Tag: "!!str"}, &value)
	}
	add("schema", ArchiveRecordSchema)
	add("spec", r.Spec)
	add("title", r.Title)
	add("status", "archived")
	add("created", r.Created)
	add("archived", r.Archived)
	add("disposition", r.Disposition)
	add("source", r.Source)
	add("source_revision", r.SourceRevision)
	add("qa_task", r.QATask)
	add("qa_report", r.QAReport)
	add("qa_verdict", r.QAVerdict)
	list := func(v []string) []string {
		if v == nil {
			return []string{}
		}
		return v
	}
	add("unproven", list(r.Unproven))
	if r.QAOverride != nil {
		o := r.QAOverride
		add("qa_override", true)
		add("qa_override_approval", o.Approval)
		add("qa_override_reason", o.Reason)
		add("qa_override_qa_outcome", o.QAOutcome)
		add("qa_override_qa_task_status", o.QATaskStatus)
		add("qa_override_revision", o.Revision)
	}
	if r.Disposition == ArchiveSuperseded {
		add("superseded_by", r.SupersededBy)
	}
	add("adrs", list(r.ADRs))
	add("sources", list(r.Sources))
	regen := r.Regeneration
	if regen == nil {
		regen = []ArchiveRegeneration{}
	}
	regenNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, entry := range regen {
		mapping := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		mapping.Content = append(mapping.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "command"}, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: entry.Command})
		if entry.Outputs != nil {
			var outputs yaml.Node
			_ = outputs.Encode(entry.Outputs)
			mapping.Content = append(mapping.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "outputs"}, &outputs)
		}
		regenNode.Content = append(regenNode.Content, mapping)
	}
	n.Content = append(n.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "regeneration"}, regenNode)
	add("promoted", list(r.Promoted))
	if r.PullRequest != "" {
		add("pull_request", r.PullRequest)
	}
	if r.DeliveryCommit != "" {
		add("delivery_commit", r.DeliveryCommit)
	}
	return n
}
func RenderArchiveRecord(r ArchiveRecord) ([]byte, error) {
	var b bytes.Buffer
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(archiveRecordMapping(r)); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	prefix := "---\n" + b.String() + "---\n\n# " + r.Title + "\n\n"
	outcome := r.Outcome
	available := ArchiveRecordTargetBytes - len(prefix) - 1
	if len(outcome) > available && available >= 0 {
		budget := available - len("…")
		if budget < 0 {
			outcome = ""
		} else {
			cut := budget
			for cut > 0 && !utf8.RuneStart(outcome[cut]) {
				cut--
			}
			shortened := outcome[:cut]
			end := strings.LastIndexAny(shortened, ".!?")
			if end >= 0 {
				outcome = shortened[:end+1]
			} else {
				outcome = shortened + "…"
			}
		}
	}
	return []byte(prefix + outcome + "\n"), nil
}
func ParseArchiveRecord(content []byte) (ArchiveRecord, error) {
	fm, body, err := splitFrontmatter(content)
	if err != nil {
		return ArchiveRecord{}, err
	}
	var m map[string]yaml.Node
	if err := yaml.Unmarshal(fm, &m); err != nil {
		return ArchiveRecord{}, err
	}
	required := []string{"schema", "spec", "title", "status", "created", "archived", "disposition", "source", "source_revision", "qa_task", "qa_report", "qa_verdict", "unproven", "adrs", "sources", "regeneration", "promoted"}
	for _, k := range required {
		if _, ok := m[k]; !ok {
			// Failed QA fields use the disposition's specific validation below.
			if m["disposition"].Value == string(ArchiveFailedQA) && (k == "qa_report" || k == "qa_verdict") {
				continue
			}
			return ArchiveRecord{}, fmt.Errorf("archive record missing %s", k)
		}
	}
	get := func(k string) string { return m[k].Value }
	if get("schema") != ArchiveRecordSchema {
		return ArchiveRecord{}, fmt.Errorf("unknown archive record schema %q", get("schema"))
	}
	if get("status") != "archived" {
		return ArchiveRecord{}, errors.New("archive record status must be archived")
	}
	r := ArchiveRecord{Spec: get("spec"), Title: get("title"), Created: get("created"), Archived: get("archived"), Disposition: ArchiveDisposition(get("disposition")), Source: get("source"), SourceRevision: get("source_revision"), QATask: get("qa_task"), QAReport: get("qa_report"), QAVerdict: get("qa_verdict"), SupersededBy: get("superseded_by"), PullRequest: get("pull_request"), DeliveryCommit: get("delivery_commit"), Outcome: archiveParagraph(string(body), "# ")}
	for k, dst := range map[string]any{"unproven": &r.Unproven, "adrs": &r.ADRs, "sources": &r.Sources, "regeneration": &r.Regeneration, "promoted": &r.Promoted} {
		n := m[k]
		if err := n.Decode(dst); err != nil {
			return r, err
		}
	}
	if r.Disposition == ArchiveQAOverride {
		for _, k := range []string{"qa_override", "qa_override_approval", "qa_override_reason", "qa_override_qa_outcome", "qa_override_qa_task_status", "qa_override_revision"} {
			if _, ok := m[k]; !ok {
				return r, fmt.Errorf("archive record missing %s", k)
			}
		}
		if strings.TrimSpace(get("qa_override")) != "true" {
			return r, fmt.Errorf("archive record qa_override must be true for disposition %s", ArchiveQAOverride)
		}
		r.QAOverride = &QAArchiveOverrideRecord{get("qa_override_approval"), get("qa_override_reason"), get("qa_override_qa_outcome"), get("qa_override_qa_task_status"), get("qa_override_revision")}
	}
	if r.Disposition == ArchiveSuperseded && r.SupersededBy == "" {
		return r, errors.New("archive record missing superseded_by")
	}
	switch r.Disposition {
	case ArchiveFailedQA:
		if r.QAVerdict != VerdictFail {
			return r, errors.New("archive record disposition failed-qa requires qa_verdict fail")
		}
		if strings.TrimSpace(r.QAReport) == "" {
			return r, errors.New("archive record disposition failed-qa requires qa_report")
		}
		if _, ok := m["qa_override"]; ok {
			return r, errors.New("archive record disposition failed-qa cannot carry qa_override")
		}
	case ArchivePass, ArchivePartial, ArchiveQAOverride, ArchiveSuperseded, ArchiveNoQA:
	default:
		return r, fmt.Errorf("unknown archive disposition %q", r.Disposition)
	}
	return r, nil
}

type ArchivedForm string

const (
	ArchivedRecord ArchivedForm = "record"
	ArchivedFolder ArchivedForm = "folder"
)

type ArchivedSpec struct {
	Slug   string
	Form   ArchivedForm
	Path   string
	Record ArchiveRecord
}

var ErrNotArchived = errors.New("Spec is not archived")

func ReadArchivedSpec(root, slug string) (ArchivedSpec, error) {
	if filepath.Base(slug) != slug || slug == "." || slug == ".." {
		return ArchivedSpec{}, fmt.Errorf("unsafe Spec slug %q", slug)
	}
	recordPath := ArchiveRecordPath(root, slug)
	folder := filepath.Join(root, slug)
	_, re := os.Lstat(recordPath)
	_, fe := os.Lstat(folder)
	if re == nil && fe == nil {
		return ArchivedSpec{}, fmt.Errorf("both archive record and folder exist for %q", slug)
	}
	if re != nil && !errors.Is(re, os.ErrNotExist) {
		return ArchivedSpec{}, re
	}
	if fe != nil && !errors.Is(fe, os.ErrNotExist) {
		return ArchivedSpec{}, fe
	}
	if re == nil {
		content, err := os.ReadFile(recordPath)
		if err != nil {
			return ArchivedSpec{}, err
		}
		r, err := ParseArchiveRecord(content)
		if err != nil {
			return ArchivedSpec{}, err
		}
		if r.Spec != slug {
			return ArchivedSpec{}, fmt.Errorf("archive record spec %q differs from file stem %q", r.Spec, slug)
		}
		return ArchivedSpec{slug, ArchivedRecord, recordPath, r}, nil
	}
	if fe == nil {
		status, err := ReadPRDStatus(folder)
		if err != nil {
			return ArchivedSpec{}, err
		}
		if status != "archived" {
			if _, err := ReadSupersession(folder); err != nil {
				return ArchivedSpec{}, fmt.Errorf("legacy Spec %q frontmatter status is %q; expected %q", slug, status, "archived")
			}
		}
		r, err := BuildArchiveRecord(ArchiveRecordInput{SpecDir: folder, Slug: slug, Legacy: true})
		if err != nil {
			return ArchivedSpec{}, err
		}
		return ArchivedSpec{slug, ArchivedFolder, folder, r}, nil
	}
	return ArchivedSpec{}, ErrNotArchived
}
func ArchivedSpecSlugs(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	for _, entry := range entries {
		if entry.IsDir() {
			if _, err := os.Stat(filepath.Join(root, entry.Name(), "_prd.md")); err == nil {
				seen[entry.Name()] = true
			}
		} else if strings.HasSuffix(entry.Name(), ".md") {
			seen[strings.TrimSuffix(entry.Name(), ".md")] = true
		}
	}
	slugs := make([]string, 0, len(seen))
	for slug := range seen {
		slugs = append(slugs, slug)
	}
	sort.Strings(slugs)
	return slugs, nil
}
