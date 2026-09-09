// Package suiteguardcontract holds repository contracts shared by the suite
// guard and the audits that verify its installation and regeneration behavior.
package suiteguardcontract

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
	"time"

	"gopkg.in/yaml.v3"
)

const sanctionedRegenerationHeading = "Sanctioned regeneration"

const (
	legacyAuthorizationRoot = "docs/workflow/authorizations"
	specAuthorizationRoot   = "docs/specs"
	baselineRoot            = "internal/baseline"
	baselineDigestCommand   = "make baseline-digests"
	ownershipYML            = "_ownership.yml"
	ownershipYAML           = "_ownership.yaml"
)

// SanctionedRegeneration binds one declared command to any outputs the record
// still enumerates. An empty Outputs slice leaves output resolution to the
// command's repository-owned declaration.
type SanctionedRegeneration struct {
	Command string
	Outputs []string
}

// ReadSanctionedRegenerations reads every operative authorization record under
// root and resolves command-only declarations through the repository-owned
// derived-output declarations.
func ReadSanctionedRegenerations(root string) ([]SanctionedRegeneration, error) {
	var declarations []SanctionedRegeneration
	if err := walkAuthorizationMarkdown(
		filepath.Join(root, filepath.FromSlash(legacyAuthorizationRoot)),
		func(_ string, content []byte) {
			declarations = append(declarations, ParseSanctionedRegenerations(content)...)
		},
	); err != nil {
		return nil, err
	}
	if err := walkAuthorizationMarkdown(
		filepath.Join(root, filepath.FromSlash(specAuthorizationRoot)),
		func(relative string, content []byte) {
			parts := strings.Split(filepath.ToSlash(relative), "/")
			if len(parts) < 2 || !approvedSpecGrant(content, parts[0]) {
				return
			}
			declarations = append(declarations, ParseSanctionedRegenerations(content)...)
		},
	); err != nil {
		return nil, err
	}

	for index := range declarations {
		if declarations[index].Outputs != nil {
			continue
		}
		outputs, err := OutputsForCommand(root, declarations[index].Command)
		if err != nil {
			return nil, fmt.Errorf(
				"resolve sanctioned regeneration command %q: %w",
				declarations[index].Command,
				err,
			)
		}
		declarations[index].Outputs = outputs
	}

	sort.Slice(declarations, func(i, j int) bool {
		if declarations[i].Command != declarations[j].Command {
			return declarations[i].Command < declarations[j].Command
		}
		return strings.Join(declarations[i].Outputs, "\x00") <
			strings.Join(declarations[j].Outputs, "\x00")
	})
	return declarations, nil
}

func walkAuthorizationMarkdown(
	root string,
	visit func(relative string, content []byte),
) error {
	info, err := os.Lstat(root)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect authorization root %q: %w", root, err)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return nil
	}
	if !info.IsDir() {
		return fmt.Errorf("authorization root %q is not a directory", root)
	}

	err = filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("inspect authorization record %q: %w", filePath, walkErr)
		}
		if entry.IsDir() || entry.Type()&fs.ModeSymlink != 0 || filepath.Ext(entry.Name()) != ".md" {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("stat authorization record %q: %w", filePath, err)
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read authorization record %q: %w", filePath, err)
		}
		relative, err := filepath.Rel(root, filePath)
		if err != nil {
			return fmt.Errorf("make authorization record %q relative to %q: %w", filePath, root, err)
		}
		visit(relative, content)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk authorization records %q: %w", root, err)
	}
	return nil
}

func approvedSpecGrant(content []byte, specSlug string) bool {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return false
	}
	closing := 0
	for index := 1; index < len(lines); index++ {
		if lines[index] == "---" {
			closing = index
			break
		}
	}
	if closing == 0 {
		return false
	}

	var grant struct {
		Status    string   `yaml:"status"`
		Granted   string   `yaml:"granted"`
		Action    string   `yaml:"action"`
		Consuming string   `yaml:"consuming"`
		Paths     []string `yaml:"paths"`
	}
	if err := yaml.Unmarshal([]byte(strings.Join(lines[1:closing], "\n")), &grant); err != nil {
		return false
	}
	if strings.TrimSpace(grant.Status) != "approved" ||
		strings.TrimSpace(grant.Action) == "" ||
		strings.TrimSpace(grant.Consuming) != specSlug {
		return false
	}
	if _, err := time.Parse("2006-01-02", strings.TrimSpace(grant.Granted)); err != nil {
		return false
	}
	if len(grant.Paths) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(grant.Paths))
	for _, declared := range grant.Paths {
		clean := cleanRepositoryPath(declared)
		if clean == "" || clean != declared || strings.ContainsAny(clean, "*?") {
			return false
		}
		if _, duplicate := seen[clean]; duplicate {
			return false
		}
		seen[clean] = struct{}{}
	}
	return true
}

// ParseSanctionedRegenerations reads the "Sanctioned regeneration" YAML
// blocks used by both the suite guard and the changed-path audit.
func ParseSanctionedRegenerations(content []byte) []SanctionedRegeneration {
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	lines := strings.Split(text, "\n")
	inSanctionedRegeneration := false
	var declarations []SanctionedRegeneration
	for index := 0; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "## ") {
			inSanctionedRegeneration = strings.EqualFold(
				strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")),
				sanctionedRegenerationHeading,
			)
			continue
		}
		if !inSanctionedRegeneration || trimmed != "```yaml" {
			continue
		}

		start := index + 1
		for index++; index < len(lines) && strings.TrimSpace(lines[index]) != "```"; index++ {
		}
		if index >= len(lines) {
			return declarations
		}
		var declaration struct {
			Command string    `yaml:"command"`
			Outputs *[]string `yaml:"outputs"`
		}
		if err := yaml.Unmarshal([]byte(strings.Join(lines[start:index], "\n")), &declaration); err != nil {
			continue
		}
		declaration.Command = strings.TrimSpace(declaration.Command)
		if declaration.Command == "" {
			continue
		}

		var outputs []string
		if declaration.Outputs != nil {
			outputs = make([]string, 0, len(*declaration.Outputs))
			seen := make(map[string]bool, len(*declaration.Outputs))
			for _, output := range *declaration.Outputs {
				output = strings.TrimSpace(output)
				clean := cleanRepositoryPath(output)
				if clean == "" || clean != output || strings.ContainsAny(clean, "*?") || seen[clean] {
					continue
				}
				seen[clean] = true
				outputs = append(outputs, clean)
			}
			if len(outputs) == 0 {
				continue
			}
			sort.Strings(outputs)
		}
		declarations = append(declarations, SanctionedRegeneration{
			Command: declaration.Command,
			Outputs: outputs,
		})
	}
	return declarations
}

type ownershipOwner string

const (
	ownerSanctioned ownershipOwner = "sanctioned"
	ownerDedicated  ownershipOwner = "dedicated"
	ownerFrozen     ownershipOwner = "frozen"
)

type ownershipRecord struct {
	Owner      ownershipOwner       `yaml:"owner"`
	Command    string               `yaml:"command,omitempty"`
	Reason     string               `yaml:"reason"`
	Exceptions []ownershipException `yaml:"exceptions,omitempty"`
}

type ownershipException struct {
	Path    string         `yaml:"path"`
	Owner   ownershipOwner `yaml:"owner"`
	Command string         `yaml:"command,omitempty"`
	Reason  string         `yaml:"reason,omitempty"`
}

type ownershipExceptionClaim struct {
	RecordPath string
	Owner      ownershipOwner
	Command    string
	Reason     string
}

// OutputsForCommand returns the sorted regular files owned by command through
// DERIVED_DIGEST_PATHS and the ownership records below internal/baseline.
func OutputsForCommand(repoRoot string, command string) ([]string, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat repository root %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("repository root %q is not a directory", root)
	}

	scanRoots, err := readDerivedScanRoots(root)
	if err != nil {
		return nil, err
	}
	derivedRoot := filepath.Join(root, filepath.FromSlash(baselineRoot))
	for _, scanRoot := range scanRoots {
		if err := rejectSymlinkComponents(derivedRoot, scanRoot); err != nil {
			return nil, err
		}
	}
	fileSystem := os.DirFS(derivedRoot)
	resolved, entries, err := resolveOwnership(fileSystem, scanRoots)
	if err != nil {
		return nil, fmt.Errorf("resolve derived ownership: %w", err)
	}

	command = strings.TrimSpace(command)
	outputs := make([]string, 0)
	for artifactPath, record := range resolved {
		entry := entries[artifactPath]
		if entry == nil || !entry.Type().IsRegular() {
			continue
		}
		ownedCommand, owned := commandForOwnership(record)
		if !owned || ownedCommand != command {
			continue
		}
		outputs = append(outputs, path.Join(baselineRoot, artifactPath))
	}
	sort.Strings(outputs)
	return outputs, nil
}

func readDerivedScanRoots(repoRoot string) ([]string, error) {
	content, err := os.ReadFile(filepath.Join(repoRoot, "Makefile"))
	if err != nil {
		return nil, fmt.Errorf("read Makefile DERIVED_DIGEST_PATHS: %w", err)
	}
	const (
		assignment = "DERIVED_DIGEST_PATHS :="
		prefix     = baselineRoot + "/"
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
			if clean != field || !strings.HasPrefix(clean, prefix) {
				return nil, fmt.Errorf("derived scan path %q is outside %s", field, baselineRoot)
			}
			root := strings.TrimPrefix(clean, prefix)
			if root == "" {
				return nil, fmt.Errorf("derived scan path must name a path below %s", baselineRoot)
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

func rejectSymlinkComponents(root, relative string) error {
	current := root
	for _, component := range strings.Split(filepath.FromSlash(relative), string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if err != nil {
			return fmt.Errorf("stat derived scan path %q: %w", filepath.ToSlash(relative), err)
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("derived scan path %q contains symlink %q", filepath.ToSlash(relative), current)
		}
	}
	return nil
}

func resolveOwnership(
	fileSystem fs.FS,
	scanRoots []string,
) (map[string]ownershipRecord, map[string]fs.DirEntry, error) {
	records, err := readOwnershipRecords(fileSystem, scanRoots)
	if err != nil {
		return nil, nil, err
	}
	exceptionClaims, err := ownershipExceptionClaims(records)
	if err != nil {
		return nil, nil, err
	}

	resolved := make(map[string]ownershipRecord)
	entries := make(map[string]fs.DirEntry)
	usedRecords := make(map[string]struct{}, len(records))
	usedExceptions := make(map[string]struct{})
	for _, scanRoot := range scanRoots {
		scanRoot = path.Clean(scanRoot)
		rootRecordPaths, err := directoryRecordPaths(fileSystem, scanRoot)
		if err != nil {
			return nil, nil, err
		}
		if len(rootRecordPaths) != 1 {
			return nil, nil, ownershipResolutionError(scanRoot, rootRecordPaths)
		}

		err = fs.WalkDir(fileSystem, scanRoot, func(artifactPath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk derived path %q: %w", artifactPath, walkErr)
			}
			if entry.Type()&fs.ModeSymlink != 0 {
				return fmt.Errorf("derived path %q is a symlink", artifactPath)
			}
			if isOwnershipRecordPath(artifactPath) {
				return nil
			}

			recordPaths := rootRecordPaths
			if artifactPath != scanRoot {
				var err error
				recordPaths, err = resolveOwnershipRecordPaths(
					fileSystem,
					scanRoot,
					artifactPath,
					entry.IsDir(),
				)
				if err != nil {
					return err
				}
			}
			if len(recordPaths) != 1 {
				return ownershipResolutionError(artifactPath, recordPaths)
			}
			recordPath, record, exceptionPath, err := resolveOwnershipRecord(
				records,
				exceptionClaims,
				recordPaths[0],
				artifactPath,
			)
			if err != nil {
				return err
			}
			resolved[artifactPath] = record
			entries[artifactPath] = entry
			usedRecords[recordPath] = struct{}{}
			if exceptionPath != "" {
				usedExceptions[exceptionPath] = struct{}{}
			}
			return nil
		})
		if err != nil {
			return nil, nil, fmt.Errorf("validate derived scan root %q: %w", scanRoot, err)
		}
	}

	for recordPath, record := range records {
		if _, used := usedRecords[recordPath]; !used {
			return nil, nil, fmt.Errorf("ownership record %q governs no derived path", recordPath)
		}
		for _, exception := range record.Exceptions {
			exceptionPath, err := ownershipExceptionPath(recordPath, exception.Path)
			if err != nil {
				return nil, nil, err
			}
			if _, used := usedExceptions[exceptionPath]; !used {
				return nil, nil, fmt.Errorf(
					"ownership exception %q in record %q governs no derived path",
					exception.Path,
					recordPath,
				)
			}
		}
	}
	return resolved, entries, nil
}

func readOwnershipRecords(fileSystem fs.FS, scanRoots []string) (map[string]ownershipRecord, error) {
	records := make(map[string]ownershipRecord)
	for _, scanRoot := range scanRoots {
		err := fs.WalkDir(fileSystem, scanRoot, func(filePath string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk derived path %q: %w", filePath, walkErr)
			}
			if entry.IsDir() || !isOwnershipRecordPath(filePath) {
				return nil
			}
			if entry.Type()&fs.ModeSymlink != 0 {
				return fmt.Errorf("ownership record %q is a symlink", filePath)
			}
			if _, duplicate := records[filePath]; duplicate {
				return fmt.Errorf("ownership record %q is present in more than one derived scan root", filePath)
			}
			record, err := readOwnershipRecord(fileSystem, filePath)
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

func readOwnershipRecord(fileSystem fs.FS, recordPath string) (ownershipRecord, error) {
	content, err := fs.ReadFile(fileSystem, recordPath)
	if err != nil {
		return ownershipRecord{}, fmt.Errorf("read ownership record %q: %w", recordPath, err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	var record ownershipRecord
	if err := decoder.Decode(&record); err != nil {
		return ownershipRecord{}, fmt.Errorf("decode ownership record %q: %w", recordPath, err)
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return ownershipRecord{}, fmt.Errorf("decode ownership record %q: multiple YAML documents", recordPath)
		}
		return ownershipRecord{}, fmt.Errorf("decode ownership record %q: %w", recordPath, err)
	}

	record.Command = strings.TrimSpace(record.Command)
	record.Reason = strings.TrimSpace(record.Reason)
	for index := range record.Exceptions {
		record.Exceptions[index].Path = strings.TrimSpace(record.Exceptions[index].Path)
		record.Exceptions[index].Command = strings.TrimSpace(record.Exceptions[index].Command)
		record.Exceptions[index].Reason = strings.TrimSpace(record.Exceptions[index].Reason)
	}
	if err := record.validate(recordPath); err != nil {
		return ownershipRecord{}, fmt.Errorf("validate ownership record %q: %w", recordPath, err)
	}
	return record, nil
}

func (record ownershipRecord) validate(recordPath string) error {
	if err := validateOwnershipOwner(record.Owner); err != nil {
		return err
	}
	if record.Reason == "" {
		return errors.New("reason is required")
	}
	if record.Owner == ownerDedicated && record.Command == "" {
		return errors.New("command is required for dedicated ownership")
	}

	seen := make(map[string]struct{}, len(record.Exceptions))
	for _, exception := range record.Exceptions {
		if err := validateOwnershipOwner(exception.Owner); err != nil {
			return fmt.Errorf("exception path %q: %w", exception.Path, err)
		}
		if exception.Owner == ownerDedicated && exception.Command == "" {
			return fmt.Errorf("command is required for dedicated ownership exception %q", exception.Path)
		}
		if exception.Owner == ownerFrozen && exception.Reason == "" {
			return fmt.Errorf("reason is required for frozen ownership exception %q", exception.Path)
		}
		exceptionPath, err := ownershipExceptionPath(recordPath, exception.Path)
		if err != nil {
			return err
		}
		if _, duplicate := seen[exceptionPath]; duplicate {
			return fmt.Errorf("derived path %q resolves to more than one ownership exception", exceptionPath)
		}
		seen[exceptionPath] = struct{}{}
	}
	return nil
}

func validateOwnershipOwner(owner ownershipOwner) error {
	switch owner {
	case ownerSanctioned, ownerDedicated, ownerFrozen:
		return nil
	default:
		return fmt.Errorf("owner %q is not sanctioned, dedicated, or frozen", owner)
	}
}

func ownershipExceptionClaims(
	records map[string]ownershipRecord,
) (map[string]ownershipExceptionClaim, error) {
	claims := make(map[string]ownershipExceptionClaim)
	for recordPath, record := range records {
		for _, exception := range record.Exceptions {
			exceptionPath, err := ownershipExceptionPath(recordPath, exception.Path)
			if err != nil {
				return nil, err
			}
			if prior, duplicate := claims[exceptionPath]; duplicate {
				return nil, fmt.Errorf(
					"derived path %q resolves to more than one ownership exception from %q and %q",
					exceptionPath,
					prior.RecordPath,
					recordPath,
				)
			}
			claims[exceptionPath] = ownershipExceptionClaim{
				RecordPath: recordPath,
				Owner:      exception.Owner,
				Command:    exception.Command,
				Reason:     exception.Reason,
			}
		}
	}
	return claims, nil
}

func ownershipExceptionPath(recordPath, exceptionPath string) (string, error) {
	if exceptionPath == "" {
		return "", errors.New("exception path is required")
	}
	if path.IsAbs(exceptionPath) {
		return "", fmt.Errorf("exception path %q is outside record directory", exceptionPath)
	}
	recordDirectory := path.Dir(recordPath)
	resolvedPath := path.Join(recordDirectory, path.Clean(exceptionPath))
	if resolvedPath == recordDirectory || !pathWithinRoot(resolvedPath, recordDirectory) {
		return "", fmt.Errorf("exception path %q is outside record directory %q", exceptionPath, recordDirectory)
	}
	return resolvedPath, nil
}

func resolveOwnershipRecord(
	records map[string]ownershipRecord,
	exceptionClaims map[string]ownershipExceptionClaim,
	baseRecordPath string,
	artifactPath string,
) (string, ownershipRecord, string, error) {
	claim, ok := exceptionClaims[artifactPath]
	if !ok {
		return baseRecordPath, records[baseRecordPath], "", nil
	}
	if claim.RecordPath != baseRecordPath && isOwnershipSidecar(baseRecordPath, artifactPath) {
		return "", ownershipRecord{}, "", fmt.Errorf(
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

func resolveOwnershipRecordPaths(
	fileSystem fs.FS,
	scanRoot string,
	artifactPath string,
	isDirectory bool,
) ([]string, error) {
	if !isDirectory {
		sidecars, err := existingOwnershipRecordPaths(
			fileSystem,
			artifactPath+ownershipYML,
			artifactPath+ownershipYAML,
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
		recordPaths, err := directoryRecordPaths(fileSystem, directory)
		if err != nil || len(recordPaths) != 0 {
			return recordPaths, err
		}
		if directory == scanRoot {
			return nil, nil
		}
		parent := path.Dir(directory)
		if parent == directory || !pathWithinRoot(parent, scanRoot) {
			return nil, nil
		}
		directory = parent
	}
}

func directoryRecordPaths(fileSystem fs.FS, directory string) ([]string, error) {
	return existingOwnershipRecordPaths(
		fileSystem,
		path.Join(directory, ownershipYML),
		path.Join(directory, ownershipYAML),
	)
}

func existingOwnershipRecordPaths(fileSystem fs.FS, candidates ...string) ([]string, error) {
	var existing []string
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

func commandForOwnership(record ownershipRecord) (string, bool) {
	switch record.Owner {
	case ownerSanctioned:
		return baselineDigestCommand, true
	case ownerDedicated:
		return record.Command, true
	default:
		return "", false
	}
}

func isOwnershipRecordPath(filePath string) bool {
	base := path.Base(filePath)
	return strings.HasSuffix(base, ownershipYML) || strings.HasSuffix(base, ownershipYAML)
}

func isOwnershipSidecar(recordPath, artifactPath string) bool {
	return recordPath == artifactPath+ownershipYML || recordPath == artifactPath+ownershipYAML
}

func pathWithinRoot(candidate, root string) bool {
	return candidate == root || strings.HasPrefix(candidate, root+"/")
}

func ownershipResolutionError(artifactPath string, recordPaths []string) error {
	if len(recordPaths) == 0 {
		return fmt.Errorf("derived path %q resolves to zero ownership records", artifactPath)
	}
	return fmt.Errorf(
		"derived path %q resolves to more than one ownership record: %s",
		artifactPath,
		strings.Join(recordPaths, ", "),
	)
}

func cleanRepositoryPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, `\`) {
		return ""
	}
	clean := filepath.ToSlash(filepath.Clean(value))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return ""
	}
	return clean
}
