package baseline

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	derivedOwnershipYML  = "_ownership.yml"
	derivedOwnershipYAML = "_ownership.yaml"

	sanctionedBaselineDigestCommand = "make baseline-digests"
	baselineDigestRegenerationHint  = "run '" + sanctionedBaselineDigestCommand + "'"
)

type derivedOwnershipOwner string

const (
	derivedOwnerSanctioned derivedOwnershipOwner = "sanctioned"
	derivedOwnerDedicated  derivedOwnershipOwner = "dedicated"
	derivedOwnerFrozen     derivedOwnershipOwner = "frozen"
)

type derivedOwnershipRecord struct {
	Owner      derivedOwnershipOwner       `yaml:"owner"`
	Command    string                      `yaml:"command,omitempty"`
	Reason     string                      `yaml:"reason"`
	Exceptions []derivedOwnershipException `yaml:"exceptions,omitempty"`
}

type derivedOwnershipException struct {
	Path    string                `yaml:"path"`
	Owner   derivedOwnershipOwner `yaml:"owner"`
	Command string                `yaml:"command,omitempty"`
	Reason  string                `yaml:"reason,omitempty"`
}

type derivedOwnershipExceptionClaim struct {
	RecordPath string
	Owner      derivedOwnershipOwner
	Command    string
	Reason     string
}

// OutputsFor returns the repository-relative paths the named regeneration
// command owns, read from the ownership records in the tree. ADR-0129 explains
// why a grant names the command rather than its consequences.
func OutputsFor(repoRoot string, command string) ([]string, error) {
	root, err := cleanRepositoryRoot(repoRoot)
	if err != nil {
		return nil, err
	}
	scanRoots, err := readDerivedDigestScanRoots(root)
	if err != nil {
		return nil, err
	}

	const baselineRoot = "internal/baseline"
	fileSystem := os.DirFS(filepath.Join(root, filepath.FromSlash(baselineRoot)))
	resolved, err := validateDerivedOwnership(fileSystem, scanRoots)
	if err != nil {
		return nil, fmt.Errorf("resolve derived ownership: %w", err)
	}

	command = strings.TrimSpace(command)
	outputs := make([]string, 0)
	for artifactPath, record := range resolved {
		ownedCommand, ownsOutput := commandForDerivedOwnership(record)
		if !ownsOutput || ownedCommand != command {
			continue
		}
		info, err := fs.Stat(fileSystem, artifactPath)
		if err != nil {
			return nil, fmt.Errorf("stat owned derived path %q: %w", artifactPath, err)
		}
		if !info.Mode().IsRegular() {
			continue
		}
		outputs = append(outputs, path.Join(baselineRoot, artifactPath))
	}
	sort.Strings(outputs)
	return outputs, nil
}

func readDerivedDigestScanRoots(repoRoot string) ([]string, error) {
	content, err := os.ReadFile(filepath.Join(repoRoot, "Makefile"))
	if err != nil {
		return nil, fmt.Errorf("read Makefile DERIVED_DIGEST_PATHS: %w", err)
	}

	const (
		assignment     = "DERIVED_DIGEST_PATHS :="
		baselinePrefix = "internal/baseline/"
	)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, assignment) {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(strings.TrimPrefix(line, assignment)))
		if len(fields) == 0 {
			return nil, errors.New("Makefile DERIVED_DIGEST_PATHS is empty")
		}
		roots := make([]string, 0, len(fields))
		seen := make(map[string]struct{}, len(fields))
		for _, field := range fields {
			clean := filepath.ToSlash(filepath.Clean(field))
			if clean != field || !strings.HasPrefix(clean, baselinePrefix) {
				return nil, fmt.Errorf("derived scan path %q is outside internal/baseline", field)
			}
			root := strings.TrimPrefix(clean, baselinePrefix)
			if root == "" {
				return nil, errors.New("derived scan path must name a path below internal/baseline")
			}
			if _, duplicate := seen[root]; duplicate {
				return nil, fmt.Errorf("derived scan path %q is declared more than once", field)
			}
			seen[root] = struct{}{}
			roots = append(roots, root)
		}
		return roots, nil
	}
	return nil, errors.New("Makefile has no DERIVED_DIGEST_PATHS assignment")
}

func commandForDerivedOwnership(record derivedOwnershipRecord) (string, bool) {
	switch record.Owner {
	case derivedOwnerSanctioned:
		return sanctionedBaselineDigestCommand, true
	case derivedOwnerDedicated:
		return record.Command, true
	default:
		return "", false
	}
}

func validateDerivedOwnership(
	fileSystem fs.FS,
	scanRoots []string,
) (map[string]derivedOwnershipRecord, error) {
	records, err := readDerivedOwnershipRecords(fileSystem, scanRoots)
	if err != nil {
		return nil, err
	}
	exceptionClaims, err := derivedOwnershipExceptionClaims(records)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]derivedOwnershipRecord)
	usedRecords := make(map[string]struct{}, len(records))
	usedExceptions := make(map[string]struct{})
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
			recordPath, record, exceptionPath, err := resolveDerivedOwnershipRecord(
				records,
				exceptionClaims,
				recordPaths[0],
				artifactPath,
			)
			if err != nil {
				return err
			}
			resolved[artifactPath] = record
			usedRecords[recordPath] = struct{}{}
			if exceptionPath != "" {
				usedExceptions[exceptionPath] = struct{}{}
			}
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
		for _, exception := range records[recordPath].Exceptions {
			exceptionPath, err := derivedOwnershipExceptionPath(recordPath, exception.Path)
			if err != nil {
				return nil, err
			}
			if _, used := usedExceptions[exceptionPath]; !used {
				return nil, fmt.Errorf(
					"ownership exception %q in record %q governs no derived path",
					exception.Path,
					recordPath,
				)
			}
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
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return derivedOwnershipRecord{}, fmt.Errorf("decode ownership record %q: multiple YAML documents", recordPath)
		}
		return derivedOwnershipRecord{}, fmt.Errorf("decode ownership record %q: %w", recordPath, err)
	}

	record.Command = strings.TrimSpace(record.Command)
	record.Reason = strings.TrimSpace(record.Reason)
	for index := range record.Exceptions {
		record.Exceptions[index].Path = strings.TrimSpace(record.Exceptions[index].Path)
		record.Exceptions[index].Command = strings.TrimSpace(record.Exceptions[index].Command)
		record.Exceptions[index].Reason = strings.TrimSpace(record.Exceptions[index].Reason)
	}
	if err := record.validate(recordPath); err != nil {
		return derivedOwnershipRecord{}, fmt.Errorf("validate ownership record %q: %w", recordPath, err)
	}
	return record, nil
}

func (record derivedOwnershipRecord) validate(recordPath string) error {
	if err := validateDerivedOwnershipOwner(record.Owner); err != nil {
		return err
	}
	if record.Reason == "" {
		return errors.New("reason is required")
	}
	if record.Owner == derivedOwnerDedicated && record.Command == "" {
		return errors.New("command is required for dedicated ownership")
	}

	seenExceptions := make(map[string]struct{}, len(record.Exceptions))
	for _, exception := range record.Exceptions {
		if err := validateDerivedOwnershipOwner(exception.Owner); err != nil {
			return fmt.Errorf("exception path %q: %w", exception.Path, err)
		}
		if exception.Owner == derivedOwnerDedicated && exception.Command == "" {
			return fmt.Errorf(
				"command is required for dedicated ownership exception %q",
				exception.Path,
			)
		}
		if exception.Owner == derivedOwnerFrozen && exception.Reason == "" {
			return fmt.Errorf(
				"reason is required for frozen ownership exception %q",
				exception.Path,
			)
		}
		exceptionPath, err := derivedOwnershipExceptionPath(recordPath, exception.Path)
		if err != nil {
			return err
		}
		if _, duplicate := seenExceptions[exceptionPath]; duplicate {
			return fmt.Errorf(
				"derived path %q resolves to more than one ownership exception",
				exceptionPath,
			)
		}
		seenExceptions[exceptionPath] = struct{}{}
	}
	return nil
}

func validateDerivedOwnershipOwner(owner derivedOwnershipOwner) error {
	switch owner {
	case derivedOwnerSanctioned, derivedOwnerDedicated, derivedOwnerFrozen:
		return nil
	default:
		return fmt.Errorf("owner %q is not sanctioned, dedicated, or frozen", owner)
	}
}

func derivedOwnershipExceptionClaims(
	records map[string]derivedOwnershipRecord,
) (map[string]derivedOwnershipExceptionClaim, error) {
	claims := make(map[string]derivedOwnershipExceptionClaim)
	for recordPath, record := range records {
		for _, exception := range record.Exceptions {
			exceptionPath, err := derivedOwnershipExceptionPath(recordPath, exception.Path)
			if err != nil {
				return nil, err
			}
			if priorClaim, duplicate := claims[exceptionPath]; duplicate {
				return nil, fmt.Errorf(
					"derived path %q resolves to more than one ownership exception from %q and %q",
					exceptionPath,
					priorClaim.RecordPath,
					recordPath,
				)
			}
			claims[exceptionPath] = derivedOwnershipExceptionClaim{
				RecordPath: recordPath,
				Owner:      exception.Owner,
				Command:    exception.Command,
				Reason:     exception.Reason,
			}
		}
	}
	return claims, nil
}

func derivedOwnershipExceptionPath(recordPath string, exceptionPath string) (string, error) {
	if exceptionPath == "" {
		return "", errors.New("exception path is required")
	}
	if path.IsAbs(exceptionPath) {
		return "", fmt.Errorf("exception path %q is outside record directory", exceptionPath)
	}

	recordDirectory := path.Dir(recordPath)
	resolvedPath := path.Join(recordDirectory, path.Clean(exceptionPath))
	if resolvedPath == recordDirectory || !pathWithinDerivedScanRoot(resolvedPath, recordDirectory) {
		return "", fmt.Errorf("exception path %q is outside record directory %q", exceptionPath, recordDirectory)
	}
	return resolvedPath, nil
}

func resolveDerivedOwnershipRecord(
	records map[string]derivedOwnershipRecord,
	exceptionClaims map[string]derivedOwnershipExceptionClaim,
	baseRecordPath string,
	artifactPath string,
) (string, derivedOwnershipRecord, string, error) {
	claim, ok := exceptionClaims[artifactPath]
	if !ok {
		return baseRecordPath, records[baseRecordPath], "", nil
	}
	if claim.RecordPath != baseRecordPath && isDerivedOwnershipSidecar(baseRecordPath, artifactPath) {
		return "", derivedOwnershipRecord{}, "", fmt.Errorf(
			"derived path %q resolves to more than one ownership record: %s, %s",
			artifactPath,
			baseRecordPath,
			claim.RecordPath,
		)
	}
	record := records[claim.RecordPath]
	record.Owner = claim.Owner
	record.Command = claim.Command
	record.Reason = claim.Reason
	return claim.RecordPath, record, artifactPath, nil
}

func isDerivedOwnershipSidecar(recordPath string, artifactPath string) bool {
	return recordPath == artifactPath+derivedOwnershipYML ||
		recordPath == artifactPath+derivedOwnershipYAML
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

func derivedArtifactRemediation(
	fileSystem fs.FS,
	scanRoot string,
	artifactPath string,
) (string, error) {
	scanRoot = path.Clean(scanRoot)
	artifactPath = path.Clean(artifactPath)
	if !pathWithinDerivedScanRoot(artifactPath, scanRoot) {
		return "", fmt.Errorf(
			"derived artifact %q is outside scan root %q",
			artifactPath,
			scanRoot,
		)
	}

	recordPaths, err := resolveDerivedOwnershipRecordPaths(
		fileSystem,
		scanRoot,
		artifactPath,
		false,
	)
	if err != nil {
		return "", fmt.Errorf("resolve ownership for derived artifact %q: %w", artifactPath, err)
	}
	if len(recordPaths) != 1 {
		return "", derivedOwnershipResolutionError(artifactPath, recordPaths)
	}

	records, err := readDerivedOwnershipRecords(fileSystem, []string{scanRoot})
	if err != nil {
		return "", err
	}
	exceptionClaims, err := derivedOwnershipExceptionClaims(records)
	if err != nil {
		return "", err
	}
	recordPath, record, _, err := resolveDerivedOwnershipRecord(
		records,
		exceptionClaims,
		recordPaths[0],
		artifactPath,
	)
	if err != nil {
		return "", fmt.Errorf("resolve ownership for derived artifact %q: %w", artifactPath, err)
	}
	switch record.Owner {
	case derivedOwnerSanctioned:
		return baselineDigestRegenerationHint, nil
	case derivedOwnerDedicated:
		return fmt.Sprintf(
			"run the exact command declared by ownership record %q: %s",
			recordPath,
			record.Command,
		), nil
	case derivedOwnerFrozen:
		return fmt.Sprintf(
			"nothing regenerates this artifact; ownership record %q records why: %s",
			recordPath,
			record.Reason,
		), nil
	default:
		return "", fmt.Errorf("ownership record %q has unsupported owner %q", recordPath, record.Owner)
	}
}
