// Package skillcoverage compares recorded behavior with its authored skill coverage.
package skillcoverage

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"
)

const (
	MapPath      = "docs/references/skill-coverage.json"
	RecordPath   = "docs/references/behavior-surfaces.json"
	MapSchema    = "roundfix/skill-coverage/v1"
	RecordSchema = "roundfix/behavior-surfaces/v1"
)

type Surface struct {
	ID        string   `json:"id"`
	Skills    []string `json:"skills,omitempty"`
	Uncovered string   `json:"uncovered,omitempty"`
	Sources   []string `json:"sources,omitempty"`
	Review    string   `json:"review,omitempty"`
}
type Map struct {
	SchemaVersion string    `json:"schemaVersion"`
	Surfaces      []Surface `json:"surfaces"`
}
type Record struct {
	SchemaVersion string            `json:"schemaVersion"`
	Surfaces      map[string]string `json:"surfaces"`
}
type Snapshot struct {
	Map    *Map
	Record *Record
}
type Change struct {
	ID      string
	Kind    string
	Skills  []string
	Outcome string
}

// ParseMap refuses invalid coverage entries rather than silently ignoring them.
func ParseMap(data []byte) (Map, error) {
	var raw struct {
		SchemaVersion string            `json:"schemaVersion"`
		Surfaces      []json.RawMessage `json:"surfaces"`
	}
	if err := decode(data, &raw); err != nil {
		return Map{}, fmt.Errorf("parse map: %w", err)
	}
	if raw.SchemaVersion != MapSchema {
		return Map{}, fmt.Errorf("schemaVersion: expected %q", MapSchema)
	}
	if raw.Surfaces == nil {
		return Map{}, fmt.Errorf("surfaces: expected array")
	}
	m := Map{SchemaVersion: raw.SchemaVersion, Surfaces: make([]Surface, 0, len(raw.Surfaces))}
	seen := make(map[string]bool)
	for i, entry := range raw.Surfaces {
		var identity struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(entry, &identity); err != nil {
			return Map{}, fmt.Errorf("surfaces[%d] id: %w", i, err)
		}
		var s Surface
		if err := decode(entry, &s); err != nil {
			return Map{}, fmt.Errorf("surface %q: %w", identity.ID, err)
		}
		if strings.TrimSpace(s.ID) == "" {
			return Map{}, fmt.Errorf("surfaces[%d] id: must be non-empty", i)
		}
		if seen[s.ID] {
			return Map{}, fmt.Errorf("surface %q id: duplicate", s.ID)
		}
		if i > 0 && m.Surfaces[i-1].ID > s.ID {
			return Map{}, fmt.Errorf("surface %q id: surfaces must be sorted", s.ID)
		}
		seen[s.ID] = true
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(entry, &fields); err != nil {
			return Map{}, fmt.Errorf("surface %q: %w", s.ID, err)
		}
		_, skillsPresent := fields["skills"]
		_, uncoveredPresent := fields["uncovered"]
		if skillsPresent == uncoveredPresent || (skillsPresent && len(s.Skills) == 0) || (uncoveredPresent && strings.TrimSpace(s.Uncovered) == "") {
			return Map{}, fmt.Errorf("surface %q: exactly one of non-empty skills or non-blank uncovered is required", s.ID)
		}
		if _, reviewPresent := fields["review"]; reviewPresent && !validCoverageReview(s.Review) {
			return Map{}, fmt.Errorf("surface %q review: want %q, got %q", s.ID, "YYYY-MM-DD — <reason>", s.Review)
		}
		for _, skill := range s.Skills {
			if !cleanPath(skill, false) {
				return Map{}, fmt.Errorf("surface %q skills: invalid path %q", s.ID, skill)
			}
		}
		for _, source := range s.Sources {
			if !cleanPath(source, true) {
				return Map{}, fmt.Errorf("surface %q sources: invalid path %q", s.ID, source)
			}
			if _, err := path.Match(strings.TrimSuffix(source, "/"), ""); err != nil {
				return Map{}, fmt.Errorf("surface %q sources pattern %q: %w", s.ID, source, err)
			}
		}
		m.Surfaces = append(m.Surfaces, s)
	}
	return m, nil
}

func cleanPath(value string, directory bool) bool {
	if strings.Contains(value, "\\") {
		return false
	}
	if directory {
		value = strings.TrimSuffix(value, "/")
	}
	if value == "" || value == "." || strings.HasPrefix(value, "/") || path.Clean(value) != value {
		return false
	}
	for _, part := range strings.Split(value, "/") {
		if part == ".." {
			return false
		}
	}
	// Drive paths are absolute on Windows even though path.IsAbs uses slash grammar.
	return !(len(value) >= 2 && value[1] == ':')
}

// ParseRecord preserves duplicate-key detection before decoding fingerprints.
func ParseRecord(data []byte) (Record, error) {
	var raw struct {
		SchemaVersion string          `json:"schemaVersion"`
		Surfaces      json.RawMessage `json:"surfaces"`
	}
	if err := decode(data, &raw); err != nil {
		return Record{}, fmt.Errorf("parse record: %w", err)
	}
	if raw.SchemaVersion != RecordSchema {
		return Record{}, fmt.Errorf("schemaVersion: expected %q", RecordSchema)
	}
	d := json.NewDecoder(bytes.NewReader(raw.Surfaces))
	token, err := d.Token()
	if err != nil {
		return Record{}, fmt.Errorf("surfaces: %w", err)
	}
	if token != json.Delim('{') {
		return Record{}, fmt.Errorf("surfaces: expected object")
	}
	r := Record{SchemaVersion: raw.SchemaVersion, Surfaces: make(map[string]string)}
	for d.More() {
		token, err := d.Token()
		if err != nil {
			return Record{}, fmt.Errorf("surfaces id: %w", err)
		}
		id, ok := token.(string)
		if !ok || strings.TrimSpace(id) == "" {
			return Record{}, fmt.Errorf("surfaces id: must be non-empty")
		}
		if _, exists := r.Surfaces[id]; exists {
			return Record{}, fmt.Errorf("surface %q id: duplicate", id)
		}
		var fingerprint string
		if err := d.Decode(&fingerprint); err != nil {
			return Record{}, fmt.Errorf("surface %q fingerprint: %w", id, err)
		}
		if !validFingerprint(fingerprint) {
			return Record{}, fmt.Errorf("surface %q fingerprint: expected sha256: and 64 lowercase hex digits", id)
		}
		r.Surfaces[id] = fingerprint
	}
	if _, err := d.Token(); err != nil {
		return Record{}, fmt.Errorf("surfaces: %w", err)
	}
	return r, nil
}

func validFingerprint(value string) bool {
	if len(value) != len("sha256:")+64 || !strings.HasPrefix(value, "sha256:") {
		return false
	}
	for _, c := range value[len("sha256:"):] {
		if !(c >= '0' && c <= '9') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func decode(data []byte, target any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return fmt.Errorf("decode fields: %w", err)
	}
	var extra json.RawMessage
	if err := d.Decode(&extra); !errors.Is(err, io.EOF) {
		if err != nil {
			return fmt.Errorf("trailing JSON: %w", err)
		}
		return fmt.Errorf("trailing JSON: expected one object")
	}
	return nil
}

// EncodeRecord gives equal records identical bytes, including one final newline.
func EncodeRecord(record Record) []byte {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		panic(err)
	} // Record contains only strings and a string map.
	return append(data, '\n')
}

func Fingerprint(content []byte) string { return fmt.Sprintf("sha256:%x", sha256.Sum256(content)) }

func (m Map) SurfacesForPath(filePath string) []Surface {
	var result []Surface
	for _, s := range m.Surfaces {
		for _, source := range s.Sources {
			match := false
			if strings.HasSuffix(source, "/") {
				match = strings.HasPrefix(filePath, source)
			} else {
				var err error
				match, err = path.Match(source, filePath)
				if err != nil {
					continue
				} // Unparsed maps may contain invalid patterns.
			}
			if match {
				result = append(result, s)
				break
			}
		}
	}
	return result
}

// Compare judges only changed fingerprints, using base coverage for removals.
func Compare(base, target Snapshot, changedPaths map[string]bool) []Change {
	var baseRecords, targetRecords map[string]string
	if base.Record != nil {
		baseRecords = base.Record.Surfaces
	}
	if target.Record != nil {
		targetRecords = target.Record.Surfaces
	}
	ids := make(map[string]bool)
	for id := range baseRecords {
		ids[id] = true
	}
	for id := range targetRecords {
		ids[id] = true
	}
	baseEntries, targetEntries := entries(base.Map), entries(target.Map)
	var changes []Change
	for id := range ids {
		before, had := baseRecords[id]
		after, has := targetRecords[id]
		if had && has && before == after {
			continue
		}
		c := Change{ID: id, Kind: "changed", Outcome: "lagging"}
		entry, found := targetEntries[id]
		if !had {
			c.Kind = "added"
		} else if !has {
			c.Kind = "removed"
			entry, found = baseEntries[id]
		}
		if found {
			c.Skills = append([]string(nil), entry.Skills...)
			if strings.TrimSpace(entry.Uncovered) != "" {
				c.Outcome = "uncovered"
			} else {
				for _, skill := range entry.Skills {
					if changedPaths[skill] {
						c.Outcome = "described"
						break
					}
				}
				if c.Outcome == "lagging" && c.Kind != "removed" && entry.Review != "" && entry.Review != baseEntries[id].Review {
					c.Outcome = "reviewed"
				}
			}
		}
		changes = append(changes, c)
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].ID < changes[j].ID })
	return changes
}

func entries(m *Map) map[string]Surface {
	result := make(map[string]Surface)
	if m != nil {
		for _, s := range m.Surfaces {
			result[s.ID] = s
		}
	}
	return result
}

// validCoverageReview reports whether a Coverage Review reads
// "YYYY-MM-DD — <reason>" with a real calendar date and a non-blank reason.
func validCoverageReview(review string) bool {
	date, reason, found := strings.Cut(review, " — ")
	if !found || strings.TrimSpace(reason) == "" {
		return false
	}
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}
