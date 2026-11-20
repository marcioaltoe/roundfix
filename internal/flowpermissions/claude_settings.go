package flowpermissions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

var (
	ErrNotStrictJSON   = errors.New("settings are not strict JSON")
	ErrUnexpectedShape = errors.New("settings have an unexpected shape")
)

type jsonSpan struct{ start, end int }
type settingsDocument struct {
	root, permissions jsonSpan
	fields            map[string]jsonSpan
	allow, ask        []string
}

// ReadClaudeSettings reads only the two rule lists; nil denotes an absent file.
func ReadClaudeSettings(data []byte) (allow, ask []string, err error) {
	doc, err := readSettings(data)
	if err != nil {
		return nil, nil, err
	}
	return doc.allow, doc.ask, nil
}

func readSettings(data []byte) (settingsDocument, error) {
	var doc settingsDocument
	if data == nil {
		return doc, nil
	}
	if !json.Valid(data) {
		return doc, ErrNotStrictJSON
	}
	doc.root = jsonSpan{len(data) - len(bytes.TrimLeft(data, " \t\r\n")), len(bytes.TrimRight(data, " \t\r\n"))}
	fields, err := objectSpans(data[doc.root.start:doc.root.end], doc.root.start)
	if err != nil {
		return doc, fmt.Errorf("read settings object: %w", err)
	}
	permissions, ok := fields["permissions"]
	if !ok {
		return doc, nil
	}
	doc.permissions = permissions
	doc.fields, err = objectSpans(data[permissions.start:permissions.end], permissions.start)
	if err != nil {
		return doc, fmt.Errorf("read permissions object: %w", err)
	}
	for _, name := range []string{"allow", "ask"} {
		span, ok := doc.fields[name]
		if !ok {
			continue
		}
		values, _, err := stringSpans(data[span.start:span.end])
		if err != nil {
			return doc, fmt.Errorf("read permissions.%s: %w", name, err)
		}
		if name == "allow" {
			doc.allow = values
		} else {
			doc.ask = values
		}
	}
	return doc, nil
}

// Decoder offsets locate raw values without normalizing any surrounding bytes.
// The input was validated before these helpers run.
func objectSpans(data []byte, offset int) (map[string]jsonSpan, error) {
	if len(data) == 0 || data[0] != '{' {
		return nil, ErrUnexpectedShape
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	if _, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("read object opening: %w", err)
	}
	fields := make(map[string]jsonSpan)
	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("read object key: %w", err)
		}
		name, ok := token.(string)
		if !ok {
			return nil, ErrUnexpectedShape
		}
		if _, exists := fields[name]; exists {
			return nil, fmt.Errorf("duplicate key %q: %w", name, ErrUnexpectedShape)
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, fmt.Errorf("read object value: %w", err)
		}
		end := offset + int(dec.InputOffset())
		fields[name] = jsonSpan{end - len(raw), end}
	}
	return fields, nil
}

func stringSpans(data []byte) ([]string, []jsonSpan, error) {
	if len(data) == 0 || data[0] != '[' {
		return nil, nil, ErrUnexpectedShape
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	if _, err := dec.Token(); err != nil {
		return nil, nil, fmt.Errorf("read list opening: %w", err)
	}
	var values []string
	var spans []jsonSpan
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, nil, fmt.Errorf("read list element: %w", err)
		}
		if len(raw) == 0 || raw[0] != '"' {
			return nil, nil, ErrUnexpectedShape
		}
		var value string
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, nil, fmt.Errorf("read string: %w", err)
		}
		end := int(dec.InputOffset())
		values = append(values, value)
		spans = append(spans, jsonSpan{end - len(raw), end})
	}
	return values, spans, nil
}

// EditClaudeSettings inserts or removes string spans, preserving unrelated bytes.
// A missing file, object or list is created with two-space indentation.
func EditClaudeSettings(data []byte, addAllow, addAsk, remove []string) ([]byte, error) {
	return editClaudeSettings(data, addAllow, addAsk, remove, remove)
}

func editClaudeSettings(data []byte, addAllow, addAsk, removeAllow, removeAsk []string) ([]byte, error) {
	doc, err := readSettings(data)
	if err != nil {
		return nil, err
	}
	wantAllow := changedList(doc.allow, addAllow, removeAllow)
	wantAsk := changedList(doc.ask, addAsk, removeAsk)
	if slices.Equal(wantAllow, doc.allow) && slices.Equal(wantAsk, doc.ask) {
		return bytes.Clone(data), nil
	}
	if data == nil {
		data = []byte("{}\n")
		doc, err = readSettings(data)
		if err != nil {
			return nil, err
		}
	}
	nl := "\n"
	if bytes.Contains(data, []byte("\r\n")) {
		nl = "\r\n"
	}
	var edits []settingsEdit
	if doc.permissions.end == 0 {
		members := []string{}
		if len(wantAllow) > 0 {
			members = append(members, `"allow": `+newList(wantAllow, "    ", nl))
		}
		if len(wantAsk) > 0 {
			members = append(members, `"ask": `+newList(wantAsk, "    ", nl))
		}
		value := "{" + nl + "    " + strings.Join(members, ","+nl+"    ") + nl + "  }"
		fields, err := objectSpans(data[doc.root.start:doc.root.end], doc.root.start)
		if err != nil {
			return nil, err
		}
		edits = append(edits, insertMembers(doc.root, fields, []string{`"permissions": ` + value}, "", nl))
	} else {
		var members []string
		for _, item := range []struct {
			name              string
			add, remove, want []string
		}{
			{"allow", addAllow, removeAllow, wantAllow}, {"ask", addAsk, removeAsk, wantAsk},
		} {
			span, exists := doc.fields[item.name]
			if !exists {
				if len(item.want) > 0 {
					members = append(members, quoteJSON(item.name)+": "+newList(item.want, "    ", nl))
				}
				continue
			}
			value, err := editList(data[span.start:span.end], item.add, item.remove)
			if err != nil {
				return nil, err
			}
			if !bytes.Equal(value, data[span.start:span.end]) {
				edits = append(edits, settingsEdit{span, value})
			}
		}
		if len(members) > 0 {
			edits = append(edits, insertMembers(doc.permissions, doc.fields, members, "  ", nl))
		}
	}
	result := applySettingsEdits(data, edits)
	allow, ask, err := ReadClaudeSettings(result)
	if err != nil {
		return nil, fmt.Errorf("validate edited settings: %w", err)
	}
	if !slices.Equal(allow, wantAllow) || !slices.Equal(ask, wantAsk) {
		return nil, fmt.Errorf("validate edited rule lists: %w", ErrUnexpectedShape)
	}
	return result, nil
}

func changedList(values, add, remove []string) []string {
	result := make([]string, 0, len(values)+len(add))
	for _, value := range values {
		if !slices.Contains(remove, value) {
			result = append(result, value)
		}
	}
	for _, value := range add {
		if !slices.Contains(result, value) {
			result = append(result, value)
		}
	}
	return result
}

type settingsEdit struct {
	jsonSpan
	data []byte
}

func applySettingsEdits(data []byte, edits []settingsEdit) []byte {
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	result := bytes.Clone(data)
	for _, edit := range edits {
		next := make([]byte, 0, len(result)-(edit.end-edit.start)+len(edit.data))
		next = append(next, result[:edit.start]...)
		next = append(next, edit.data...)
		next = append(next, result[edit.end:]...)
		result = next
	}
	return result
}

func editList(data []byte, add, remove []string) ([]byte, error) {
	values, spans, err := stringSpans(data)
	if err != nil {
		return nil, err
	}
	var kept []int
	for i, value := range values {
		if !slices.Contains(remove, value) {
			kept = append(kept, i)
		}
	}
	var result []byte
	if len(kept) == 0 {
		if len(values) > 0 {
			result = []byte("[]")
		} else {
			result = bytes.Clone(data)
		}
	} else {
		var edits []settingsEdit
		if kept[0] > 0 {
			edits = append(edits, settingsEdit{jsonSpan{spans[0].start, spans[kept[0]].start}, nil})
		}
		for i, index := range kept {
			next := len(spans)
			if i+1 < len(kept) {
				next = kept[i+1]
			}
			if next > index+1 {
				edits = append(edits, settingsEdit{jsonSpan{spans[index].end, spans[next-1].end}, nil})
			}
		}
		result = applySettingsEdits(data, edits)
	}
	values, spans, err = stringSpans(result)
	if err != nil {
		return nil, err
	}
	for _, value := range add {
		if slices.Contains(values, value) {
			continue
		}
		if len(spans) == 0 {
			result = applySettingsEdits(result, []settingsEdit{{jsonSpan{1, 1}, []byte(quoteJSON(value))}})
		} else {
			last := spans[len(spans)-1]
			separator := ", "
			if len(spans) > 1 {
				separator = string(result[spans[len(spans)-2].end:last.start])
			} else if last.start > 1 {
				separator = "," + string(result[1:last.start])
			}
			result = applySettingsEdits(result, []settingsEdit{{jsonSpan{last.end, last.end}, []byte(separator + quoteJSON(value))}})
		}
		values, spans, err = stringSpans(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func quoteJSON(value string) string {
	// Encoding a string cannot fail.
	data, _ := json.Marshal(value)
	return string(data)
}

func newList(values []string, indent, nl string) string {
	encoded := make([]string, len(values))
	for i, value := range values {
		encoded[i] = quoteJSON(value)
	}
	return "[" + nl + indent + "  " + strings.Join(encoded, ","+nl+indent+"  ") + nl + indent + "]"
}

func insertMembers(object jsonSpan, fields map[string]jsonSpan, members []string, indent, nl string) settingsEdit {
	last := object.start + 1
	for _, span := range fields {
		if span.end > last {
			last = span.end
		}
	}
	prefix := nl + indent + "  "
	if len(fields) > 0 {
		prefix = "," + prefix
	}
	text := prefix + strings.Join(members, ","+nl+indent+"  ")
	// Existing whitespace before the closing brace stays exactly where it was.
	if last == object.end-1 {
		text += nl + indent
	}
	return settingsEdit{jsonSpan{last, last}, []byte(text)}
}
