package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// RecordInstalledOwnedSkills records installed folder hashes only for existing
// lock entries named by the owned bundle, preserving all other document bytes.
func RecordInstalledOwnedSkills(ctx context.Context, root string, names []string) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	filename := filepath.Join(root, "skills-lock.json")
	document, err := loadSkillsLock(filename)
	if err != nil {
		return nil, err
	}
	if document.before == nil {
		return nil, nil
	}
	fields, err := ownedLockObjectFields(document.before)
	if err != nil {
		return nil, fmt.Errorf("read skills-lock.json fields: %w", err)
	}
	skillsField := fields["skills"]
	entries, err := ownedLockObjectFields(document.before[skillsField.start:skillsField.end])
	if err != nil {
		return nil, fmt.Errorf("read skills-lock.json skills: %w", err)
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	var updated []string
	type replacement struct {
		start, end int
		data       []byte
	}
	var edits []replacement
	for _, name := range sorted {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if len(updated) > 0 && updated[len(updated)-1] == name {
			continue
		}
		if _, ok := document.objectEntry("skills", name); !ok {
			continue
		}
		files, exists, err := inspectRestoreTarget(root, ".agents/skills/"+name)
		if err != nil {
			return nil, fmt.Errorf("inspect installed skill %q for skills-lock.json: %w", name, err)
		}
		if !exists {
			continue
		}
		digest := externalSkillsLockDigest(files)
		entry := entries[name]
		offset := skillsField.start + entry.start
		raw := document.before[offset : skillsField.start+entry.end]
		values, err := ownedLockObjectFields(raw)
		if err != nil {
			return nil, fmt.Errorf("read skills-lock.json entry %q: %w", name, err)
		}
		encoded, err := json.Marshal(digest)
		if err != nil {
			return nil, fmt.Errorf("encode skills-lock.json hash: %w", err)
		}
		if hash, ok := values["computedHash"]; ok {
			var previous string
			if json.Unmarshal(raw[hash.start:hash.end], &previous) == nil && previous == digest {
				continue
			}
			edits = append(edits, replacement{offset + hash.start, offset + hash.end, encoded})
		} else {
			insertion := []byte(`"computedHash":`)
			if len(values) > 0 {
				insertion = append([]byte(","), insertion...)
			}
			insertion = append(insertion, encoded...)
			end := offset + len(raw) - 1
			edits = append(edits, replacement{end, end, insertion})
		}
		updated = append(updated, name)
	}
	if len(updated) == 0 {
		return nil, nil
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	data := append([]byte(nil), document.before...)
	for _, edit := range edits {
		data = append(append(append([]byte(nil), data[:edit.start]...), edit.data...), data[edit.end:]...)
	}
	info, err := os.Lstat(filename)
	if err != nil {
		return nil, fmt.Errorf("stat skills-lock.json: %w", err)
	}
	file, err := os.CreateTemp(root, ".skills-lock-*")
	if err != nil {
		return nil, fmt.Errorf("create skills-lock.json temporary file: %w", err)
	}
	defer os.Remove(file.Name())
	if err := file.Chmod(info.Mode()); err != nil {
		file.Close()
		return nil, fmt.Errorf("preserve skills-lock.json mode: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return nil, fmt.Errorf("write skills-lock.json: %w", err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("close skills-lock.json: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := os.Rename(file.Name(), filename); err != nil {
		return nil, fmt.Errorf("replace skills-lock.json: %w", err)
	}
	return updated, nil
}

type ownedLockField struct{ start, end int }

// The validated ordered document supplies semantics; offsets let the writer
// replace just the selected JSON value without normalizing unrelated bytes.
func ownedLockObjectFields(data []byte) (map[string]ownedLockField, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	fields := make(map[string]ownedLockField)
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		name, ok := token.(string)
		if !ok {
			return nil, fmt.Errorf("object key is not a string")
		}
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			return nil, err
		}
		end := int(decoder.InputOffset())
		fields[name] = ownedLockField{end - len(raw), end}
	}
	return fields, nil
}
