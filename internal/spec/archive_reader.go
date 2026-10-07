package spec

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

// ReadArchivedSpecAt resolves archive metadata in an immutable tree. read must
// return os.ErrNotExist for absent paths. Legacy folders retain their existing
// readers; only the PRD is needed here to identify that form.
func ReadArchivedSpecAt(root, slug string, read func(string) ([]byte, error)) (ArchivedSpec, error) {
	if slug == "" || path.Base(slug) != slug || slug == "." || slug == ".." || strings.Contains(slug, "\\") {
		return ArchivedSpec{}, fmt.Errorf("unsafe Spec slug %q", slug)
	}
	recordPath := path.Join(root, slug+".md")
	folder := path.Join(root, slug)
	content, recordErr := read(recordPath)
	_, folderErr := read(path.Join(folder, "_prd.md"))
	if recordErr == nil && folderErr == nil {
		return ArchivedSpec{}, fmt.Errorf("both archive record and folder exist for %q", slug)
	}
	for _, err := range []error{recordErr, folderErr} {
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return ArchivedSpec{}, err
		}
	}
	if recordErr == nil {
		record, err := ParseArchiveRecord(content)
		if err != nil {
			return ArchivedSpec{}, err
		}
		if record.Spec != slug {
			return ArchivedSpec{}, fmt.Errorf("archive record spec %q differs from file stem %q", record.Spec, slug)
		}
		if !validArchiveSource(record.Source, record.SourceRevision) {
			return ArchivedSpec{}, fmt.Errorf("archive record %q has unsafe source or source_revision", slug)
		}
		return ArchivedSpec{Slug: slug, Form: ArchivedRecord, Path: recordPath, Record: record}, nil
	}
	if folderErr == nil {
		return ArchivedSpec{Slug: slug, Form: ArchivedFolder, Path: folder}, nil
	}
	return ArchivedSpec{}, ErrNotArchived
}

func validArchiveSource(source, revision string) bool {
	if source == "" || source == "." || path.IsAbs(source) || path.Clean(source) != source || source == ".." || strings.HasPrefix(source, "../") || strings.ContainsAny(source, "\\\x00") {
		return false
	}
	if len(revision) != 40 && len(revision) != 64 {
		return false
	}
	for _, c := range revision {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}
	return true
}

// ArchivedTaskCompleted is the record's completion contract. It does not
// infer membership: callers bind Task identity through commit metadata.
func ArchivedTaskCompleted(record ArchiveRecord, task string) bool {
	return task != "" && !(task == record.QATask && record.Disposition == ArchiveQAOverride && record.QAOverride != nil && record.QAOverride.QATaskStatus != "")
}
