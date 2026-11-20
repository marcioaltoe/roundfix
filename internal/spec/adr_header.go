package spec

import (
	"bytes"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ADRHeader is the documentation metadata of a decision record.
type ADRHeader struct {
	Frontmatter bool
	Status      string
	Knowledge   *bool
	AmendedBy   []string
}

// adrHeaderLines finds a leading frontmatter block without rewriting its bytes.
func adrHeaderLines(content []byte) ([][]byte, int) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	if len(lines) == 0 || strings.TrimRight(string(lines[0]), "\r\n") != "---" {
		return lines, -1
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(string(lines[i]), "\r\n") == "---" {
			return lines, i
		}
	}
	return lines, -1
}

var adrIdentifier = regexp.MustCompile(`^ADR-[0-9]+$`)

// ReadADRHeader reads fields wherever they occur in the header. It never encodes YAML.
func ReadADRHeader(content []byte) ADRHeader {
	lines, end := adrHeaderLines(content)
	h := ADRHeader{Frontmatter: end >= 0, AmendedBy: []string{}}
	if end < 0 {
		return h
	}
	var document yaml.Node
	if yaml.Unmarshal(bytes.Join(lines[1:end], nil), &document) != nil || len(document.Content) == 0 {
		return h
	}
	fields := document.Content[0]
	if fields.Kind != yaml.MappingNode {
		return h
	}
	for i := 0; i+1 < len(fields.Content); i += 2 {
		key, value := fields.Content[i].Value, fields.Content[i+1]
		switch key {
		case "status":
			h.Status = strings.ToLower(strings.TrimSpace(value.Value))
		case "knowledge":
			if value.Tag == "!!bool" {
				var mark bool
				if value.Decode(&mark) == nil {
					h.Knowledge = &mark
				}
			}
		case "amended_by":
			if value.Kind == yaml.SequenceNode {
				for _, id := range value.Content {
					if adrIdentifier.MatchString(id.Value) {
						h.AmendedBy = append(h.AmendedBy, id.Value)
					}
				}
			}
		}
	}
	return h
}

var adrNotice = regexp.MustCompile(`(?:Amended by|Narrowed by|Extended by|Superseded in part by)\s+((?:\[ADR-[0-9]+\]\([^\n)]*\)|ADR-[0-9]+\b)(?:(?:\s*,\s*(?:and\s+)?|\s+and\s+)(?:\[ADR-[0-9]+\]\([^\n)]*\)|ADR-[0-9]+\b))*)`)
var adrNoticeIdentifier = regexp.MustCompile(`ADR-[0-9]+`)
var adrLink = regexp.MustCompile(`\[(ADR-[0-9]+)\]\([^\n)]*\)`)

// ADRAmendmentNotices reads the four amendment forms outside Markdown code.
func ADRAmendmentNotices(content []byte) []string {
	lines, end := adrHeaderLines(content)
	if end >= 0 {
		lines = lines[end+1:]
	}
	var prose strings.Builder
	fence := byte(0)
	fenceLength := 0
	inline := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if len(trimmed) >= 3 && (trimmed[0] == '`' || trimmed[0] == '~') {
			n := 0
			for n < len(trimmed) && trimmed[n] == trimmed[0] {
				n++
			}
			if n >= 3 {
				if fence == 0 {
					fence = trimmed[0]
					fenceLength = n
					prose.WriteByte('\n')
					continue
				}
				if fence == trimmed[0] && n >= fenceLength && strings.TrimSpace(trimmed[n:]) == "" {
					fence = 0
					continue
				}
			}
		}
		if fence != 0 {
			continue
		}
		for i := 0; i < len(line); {
			if line[i] == '`' {
				n := 1
				for i+n < len(line) && line[i+n] == '`' {
					n++
				}
				if inline == 0 {
					inline = n
					prose.WriteByte(' ')
				} else if inline == n {
					inline = 0
					prose.WriteByte(' ')
				}
				i += n
				continue
			}
			if inline == 0 {
				prose.WriteByte(line[i])
			}
			i++
		}
	}
	ids := []string{}
	for _, match := range adrNotice.FindAllStringSubmatch(prose.String(), -1) {
		// Link destinations are not identifiers in the notice's text.
		text := adrLink.ReplaceAllString(match[1], "$1")
		ids = append(ids, adrNoticeIdentifier.FindAllString(text, -1)...)
	}
	return orderedADRIdentifiers(ids)
}

// adrNumber is the number an identifier names, so ADR-55, ADR-055 and 0055 agree.
func adrNumber(id string) string {
	return strings.TrimLeft(strings.TrimPrefix(id, "ADR-"), "0")
}

func orderedADRIdentifiers(ids []string) []string {
	unique := map[string]bool{}
	result := []string{}
	for _, id := range ids {
		if adrIdentifier.MatchString(id) && !unique[id] {
			unique[id] = true
			result = append(result, id)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		a, b := adrNumber(result[i]), adrNumber(result[j])
		if len(a) != len(b) {
			return len(a) < len(b)
		}
		if a != b {
			return a < b
		}
		return result[i] < result[j]
	})
	return result
}

// WithADRAmendedBy inserts or replaces only the Amendment List line, keeping all
// other bytes (including timestamps and line endings) as they were.
func WithADRAmendedBy(content []byte, ids []string) []byte {
	return WithADRAmendedByExcept(content, ids, "")
}

// WithADRAmendedByExcept writes the union of the identifiers the header already
// lists and ids, without an identifier whose number equals own's. An ADR never
// lists itself, and every other identifier a person wrote stays. With an empty
// own it is WithADRAmendedBy.
func WithADRAmendedByExcept(content []byte, ids []string, own string) []byte {
	h := ReadADRHeader(content)
	combined := append(append([]string{}, h.AmendedBy...), ids...)
	kept := make([]string, 0, len(combined))
	for _, id := range combined {
		if own != "" && adrNumber(id) == adrNumber(own) {
			continue
		}
		kept = append(kept, id)
	}
	all := orderedADRIdentifiers(kept)
	lines, end := adrHeaderLines(content)
	existing, status, knowledge := -1, -1, -1
	if end >= 0 {
		for i := 1; i < end; i++ {
			field, _, ok := strings.Cut(string(lines[i]), ":")
			if !ok {
				continue
			}
			switch field {
			case "amended_by":
				if existing < 0 {
					existing = i
				}
			case "status":
				status = i
			case "knowledge":
				knowledge = i
			}
		}
	}
	// A header with no Amendment List and no identifier to write stays as it was;
	// an ADR that drops only its own identifier still rewrites an existing list.
	if len(all) == 0 && (own == "" || existing < 0) {
		return content
	}
	newline := "\n"
	if bytes.Contains(content, []byte("\r\n")) {
		newline = "\r\n"
	}
	line := []byte("amended_by: [" + strings.Join(all, ", ") + "]" + newline)
	if end < 0 {
		return append([]byte("---"+newline+string(line)+"---"+newline), content...)
	}
	if existing >= 0 {
		if !bytes.HasSuffix(lines[existing], []byte("\n")) {
			line = bytes.TrimSuffix(line, []byte(newline))
		}
		output := append([][]byte(nil), lines...)
		output[existing] = line
		return bytes.Join(output, nil)
	}
	insert := 1
	if knowledge >= 0 {
		insert = knowledge + 1
	} else if status >= 0 {
		insert = status + 1
	}
	result := bytes.Join(lines[:insert], nil)
	result = append(result, line...)
	return append(result, bytes.Join(lines[insert:], nil)...)
}
