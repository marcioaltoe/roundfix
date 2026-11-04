package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
)

// OwnedSkillsLockUpdate names the installed entries changed by the lock writer.
type OwnedSkillsLockUpdate struct {
	Hashed   []string
	Unpinned []string
}

// Names returns the sorted union of entries whose hash or ref changed.
func (update OwnedSkillsLockUpdate) Names() []string {
	names := append([]string(nil), update.Hashed...)
	names = append(names, update.Unpinned...)
	sort.Strings(names)
	return slices.Compact(names)
}

// RecordInstalledOwnedSkills records installed hashes and removes refs only for
// existing lock entries named by the owned bundle, preserving all other bytes.
func RecordInstalledOwnedSkills(ctx context.Context, root string, names []string) (OwnedSkillsLockUpdate, error) {
	if err := ctx.Err(); err != nil {
		return OwnedSkillsLockUpdate{}, err
	}
	filename := filepath.Join(root, "skills-lock.json")
	document, err := loadSkillsLock(filename)
	if err != nil {
		return OwnedSkillsLockUpdate{}, err
	}
	if document.before == nil {
		return OwnedSkillsLockUpdate{}, nil
	}
	fields, err := ownedLockObjectFields(document.before)
	if err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("read skills-lock.json fields: %w", err)
	}
	skillsField := fields["skills"]
	entries, err := ownedLockObjectFields(document.before[skillsField.start:skillsField.end])
	if err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("read skills-lock.json skills: %w", err)
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	var update OwnedSkillsLockUpdate
	type replacement struct {
		start, end int
		data       []byte
	}
	var edits []replacement
	for i, name := range sorted {
		if err := ctx.Err(); err != nil {
			return OwnedSkillsLockUpdate{}, err
		}
		if i > 0 && sorted[i-1] == name {
			continue
		}
		if _, ok := document.objectEntry("skills", name); !ok {
			continue
		}
		files, exists, err := inspectRestoreTarget(root, ".agents/skills/"+name)
		if err != nil {
			return OwnedSkillsLockUpdate{}, fmt.Errorf("inspect installed skill %q for skills-lock.json: %w", name, err)
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
			return OwnedSkillsLockUpdate{}, fmt.Errorf("read skills-lock.json entry %q: %w", name, err)
		}
		if ref, ok := values["ref"]; ok {
			edits = append(edits, replacement{offset + ref.removeStart, offset + ref.removeEnd, nil})
			update.Unpinned = append(update.Unpinned, name)
		}
		encoded, err := json.Marshal(digest)
		if err != nil {
			return OwnedSkillsLockUpdate{}, fmt.Errorf("encode skills-lock.json hash: %w", err)
		}
		if hash, ok := values["computedHash"]; ok {
			var previous string
			if json.Unmarshal(raw[hash.start:hash.end], &previous) == nil && previous == digest {
				continue
			}
			edits = append(edits, replacement{offset + hash.start, offset + hash.end, encoded})
		} else {
			insertion := []byte(`"computedHash":`)
			remaining := len(values)
			if _, ok := values["ref"]; ok {
				remaining--
			}
			if remaining > 0 {
				insertion = append([]byte(","), insertion...)
			}
			insertion = append(insertion, encoded...)
			end := offset + len(raw) - 1
			edits = append(edits, replacement{end, end, insertion})
		}
		update.Hashed = append(update.Hashed, name)
	}
	if len(edits) == 0 {
		return OwnedSkillsLockUpdate{}, nil
	}
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	data := append([]byte(nil), document.before...)
	for _, edit := range edits {
		data = append(append(append([]byte(nil), data[:edit.start]...), edit.data...), data[edit.end:]...)
	}
	info, err := os.Lstat(filename)
	if err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("stat skills-lock.json: %w", err)
	}
	file, err := os.CreateTemp(root, ".skills-lock-*")
	if err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("create skills-lock.json temporary file: %w", err)
	}
	defer os.Remove(file.Name())
	if err := file.Chmod(info.Mode()); err != nil {
		file.Close()
		return OwnedSkillsLockUpdate{}, fmt.Errorf("preserve skills-lock.json mode: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return OwnedSkillsLockUpdate{}, fmt.Errorf("write skills-lock.json: %w", err)
	}
	if err := file.Close(); err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("close skills-lock.json: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return OwnedSkillsLockUpdate{}, err
	}
	if err := os.Rename(file.Name(), filename); err != nil {
		return OwnedSkillsLockUpdate{}, fmt.Errorf("replace skills-lock.json: %w", err)
	}
	return update, nil
}

type ownedLockField struct {
	start, end             int
	removeStart, removeEnd int
}

// The validated ordered document supplies semantics; offsets let the writer
// replace a value or remove a member without normalizing unrelated bytes.
func ownedLockObjectFields(data []byte) (map[string]ownedLockField, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}
	fields := make(map[string]ownedLockField)
	previousEnd := 0
	var firstName string
	for decoder.More() {
		keyStart := int(decoder.InputOffset())
		for data[keyStart] != '"' {
			keyStart++
		}
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
		removeStart := previousEnd
		if len(fields) == 0 {
			firstName = name
			removeStart = keyStart
		} else if len(fields) == 1 {
			first := fields[firstName]
			first.removeEnd = keyStart
			fields[firstName] = first
		}
		fields[name] = ownedLockField{start: end - len(raw), end: end, removeStart: removeStart, removeEnd: end}
		previousEnd = end
	}
	return fields, nil
}
