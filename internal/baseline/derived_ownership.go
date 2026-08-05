package baseline

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	derivedOwnershipYML  = "_ownership.yml"
	derivedOwnershipYAML = "_ownership.yaml"
)

type derivedOwnershipOwner string

const (
	derivedOwnerSanctioned derivedOwnershipOwner = "sanctioned"
	derivedOwnerDedicated  derivedOwnershipOwner = "dedicated"
	derivedOwnerFrozen     derivedOwnershipOwner = "frozen"
)

type derivedOwnershipRecord struct {
	Owner   derivedOwnershipOwner `yaml:"owner"`
	Command string                `yaml:"command,omitempty"`
	Reason  string                `yaml:"reason"`
}

func validateDerivedOwnership(
	fileSystem fs.FS,
	scanRoots []string,
) (map[string]derivedOwnershipRecord, error) {
	records, err := readDerivedOwnershipRecords(fileSystem, scanRoots)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]derivedOwnershipRecord)
	usedRecords := make(map[string]struct{}, len(records))
	for _, scanRoot := range scanRoots {
		scanRoot = path.Clean(scanRoot)
		rootInfo, err := fs.Stat(fileSystem, scanRoot)
		if err != nil {
			return nil, fmt.Errorf("stat derived scan root %q: %w", scanRoot, err)
		}
		if !rootInfo.IsDir() {
			return nil, fmt.Errorf("derived scan root %q is not a directory", scanRoot)
		}

		rootRecordPaths, err := derivedDirectoryRecordPaths(fileSystem, scanRoot)
		if err != nil {
			return nil, err
		}
		if len(rootRecordPaths) != 1 {
			return nil, derivedOwnershipResolutionError(scanRoot, rootRecordPaths)
		}
		resolved[scanRoot] = records[rootRecordPaths[0]]
		usedRecords[rootRecordPaths[0]] = struct{}{}

		err = fs.WalkDir(fileSystem, scanRoot, func(artifactPath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk derived path %q: %w", artifactPath, walkErr)
			}
			if artifactPath == scanRoot || isDerivedOwnershipRecordPath(artifactPath) {
				return nil
			}

			recordPaths, err := resolveDerivedOwnershipRecordPaths(
				fileSystem,
				scanRoot,
				artifactPath,
				entry.IsDir(),
			)
			if err != nil {
				return err
			}
			if len(recordPaths) != 1 {
				return derivedOwnershipResolutionError(artifactPath, recordPaths)
			}
			resolved[artifactPath] = records[recordPaths[0]]
			usedRecords[recordPaths[0]] = struct{}{}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("validate derived scan root %q: %w", scanRoot, err)
		}
	}

	for recordPath := range records {
		if _, used := usedRecords[recordPath]; !used {
			return nil, fmt.Errorf("ownership record %q governs no derived path", recordPath)
		}
	}
	return resolved, nil
}

func readDerivedOwnershipRecords(
	fileSystem fs.FS,
	scanRoots []string,
) (map[string]derivedOwnershipRecord, error) {
	records := make(map[string]derivedOwnershipRecord)
	for _, scanRoot := range scanRoots {
		scanRoot = path.Clean(scanRoot)
		err := fs.WalkDir(fileSystem, scanRoot, func(filePath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk derived path %q: %w", filePath, walkErr)
			}
			if entry.IsDir() || !isDerivedOwnershipRecordPath(filePath) {
				return nil
			}
			if _, duplicate := records[filePath]; duplicate {
				return fmt.Errorf("ownership record %q is present in more than one derived scan root", filePath)
			}
			record, err := readDerivedOwnershipRecord(fileSystem, filePath)
			if err != nil {
				return err
			}
			records[filePath] = record
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("read ownership records under %q: %w", scanRoot, err)
		}
	}
	return records, nil
}

func readDerivedOwnershipRecord(fileSystem fs.FS, recordPath string) (derivedOwnershipRecord, error) {
	content, err := fs.ReadFile(fileSystem, recordPath)
	if err != nil {
		return derivedOwnershipRecord{}, fmt.Errorf("read ownership record %q: %w", recordPath, err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	var record derivedOwnershipRecord
	if err := decoder.Decode(&record); err != nil {
		return derivedOwnershipRecord{}, fmt.Errorf("decode ownership record %q: %w", recordPath, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return derivedOwnershipRecord{}, fmt.Errorf("decode ownership record %q: multiple YAML documents", recordPath)
		}
		return derivedOwnershipRecord{}, fmt.Errorf("decode ownership record %q: %w", recordPath, err)
	}

	record.Command = strings.TrimSpace(record.Command)
	record.Reason = strings.TrimSpace(record.Reason)
	if err := record.validate(); err != nil {
		return derivedOwnershipRecord{}, fmt.Errorf("validate ownership record %q: %w", recordPath, err)
	}
	return record, nil
}

func (record derivedOwnershipRecord) validate() error {
	switch record.Owner {
	case derivedOwnerSanctioned, derivedOwnerDedicated, derivedOwnerFrozen:
	default:
		return fmt.Errorf("owner %q is not sanctioned, dedicated, or frozen", record.Owner)
	}
	if record.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if record.Owner == derivedOwnerDedicated && record.Command == "" {
		return fmt.Errorf("command is required for dedicated ownership")
	}
	return nil
}

func resolveDerivedOwnershipRecordPaths(
	fileSystem fs.FS,
	scanRoot string,
	artifactPath string,
	isDirectory bool,
) ([]string, error) {
	if !isDirectory {
		sidecars, err := existingDerivedOwnershipRecordPaths(
			fileSystem,
			artifactPath+derivedOwnershipYML,
			artifactPath+derivedOwnershipYAML,
		)
		if err != nil || len(sidecars) != 0 {
			return sidecars, err
		}
	}

	directory := artifactPath
	if !isDirectory {
		directory = path.Dir(artifactPath)
	}
	for {
		recordPaths, err := derivedDirectoryRecordPaths(fileSystem, directory)
		if err != nil || len(recordPaths) != 0 {
			return recordPaths, err
		}
		if directory == scanRoot {
			return nil, nil
		}
		parent := path.Dir(directory)
		if parent == directory || !pathWithinDerivedScanRoot(parent, scanRoot) {
			return nil, nil
		}
		directory = parent
	}
}

func derivedDirectoryRecordPaths(fileSystem fs.FS, directory string) ([]string, error) {
	return existingDerivedOwnershipRecordPaths(
		fileSystem,
		path.Join(directory, derivedOwnershipYML),
		path.Join(directory, derivedOwnershipYAML),
	)
}

func existingDerivedOwnershipRecordPaths(fileSystem fs.FS, candidates ...string) ([]string, error) {
	existing := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		info, err := fs.Stat(fileSystem, candidate)
		switch {
		case err == nil && info.IsDir():
			return nil, fmt.Errorf("ownership record %q is a directory", candidate)
		case err == nil:
			existing = append(existing, candidate)
		case !errors.Is(err, fs.ErrNotExist):
			return nil, fmt.Errorf("stat ownership record %q: %w", candidate, err)
		}
	}
	return existing, nil
}

func isDerivedOwnershipRecordPath(filePath string) bool {
	base := path.Base(filePath)
	return strings.HasSuffix(base, derivedOwnershipYML) ||
		strings.HasSuffix(base, derivedOwnershipYAML)
}

func pathWithinDerivedScanRoot(candidate string, scanRoot string) bool {
	return candidate == scanRoot || strings.HasPrefix(candidate, scanRoot+"/")
}

func derivedOwnershipResolutionError(artifactPath string, recordPaths []string) error {
	if len(recordPaths) == 0 {
		return fmt.Errorf("derived path %q resolves to zero ownership records", artifactPath)
	}
	return fmt.Errorf(
		"derived path %q resolves to more than one ownership record: %s",
		artifactPath,
		strings.Join(recordPaths, ", "),
	)
}
