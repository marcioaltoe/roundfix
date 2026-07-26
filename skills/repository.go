package skills

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type RepositoryOwnership string

const (
	RepositoryOwnershipOwned    RepositoryOwnership = "Roundfix-owned"
	RepositoryOwnershipExternal RepositoryOwnership = "external"
)

type RepositoryReadiness struct {
	OwnedRequired    int
	ExternalRequired int
	MissingOwned     []string
	OutdatedOwned    []string
	MissingExternal  []string
	OutdatedExternal []string
}

func (readiness RepositoryReadiness) Ready() bool {
	return len(readiness.MissingOwned) == 0 &&
		len(readiness.OutdatedOwned) == 0 &&
		len(readiness.MissingExternal) == 0 &&
		len(readiness.OutdatedExternal) == 0
}

type RepositoryReadinessError struct {
	Ownership RepositoryOwnership
	Operation string
	Path      string
	Err       error
}

func (err *RepositoryReadinessError) Error() string {
	if err.Path == "" {
		return fmt.Sprintf("%s: %v", err.Operation, err.Err)
	}
	return fmt.Sprintf("%s %q: %v", err.Operation, err.Path, err.Err)
}

func (err *RepositoryReadinessError) Unwrap() error {
	return err.Err
}

type localSkillsLock struct {
	Version int                       `json:"version"`
	Skills  map[string]localLockSkill `json:"skills"`
}

type localLockSkill struct {
	ComputedHash string `json:"computedHash"`
}

func CheckRepository(root string) (RepositoryReadiness, error) {
	owned := Names()
	external := Recommended()
	if err := validateRequiredSkillNames(owned, external); err != nil {
		return RepositoryReadiness{
			OwnedRequired:    len(owned),
			ExternalRequired: len(external),
		}, err
	}
	files, err := Files()
	if err != nil {
		return RepositoryReadiness{
			OwnedRequired:    len(owned),
			ExternalRequired: len(external),
		}, repositoryReadinessError(RepositoryOwnershipOwned, "read embedded skill bundle", "", err)
	}
	return checkRepository(root, owned, external, files)
}

func checkRepository(root string, owned []string, external []string, files []File) (readiness RepositoryReadiness, returnErr error) {
	readiness = RepositoryReadiness{
		OwnedRequired:    len(owned),
		ExternalRequired: len(external),
	}
	if err := validateRequiredSkillNames(owned, external); err != nil {
		return readiness, err
	}
	root = strings.TrimSpace(root)
	if root == "" {
		return readiness, repositoryReadinessError("", "validate repository root", "", errors.New("path is empty"))
	}

	repositoryRoot, err := os.OpenRoot(root)
	if err != nil {
		return readiness, repositoryReadinessError("", "open repository root", root, err)
	}
	defer func() {
		if closeErr := repositoryRoot.Close(); closeErr != nil {
			wrapped := repositoryReadinessError("", "close repository root", root, closeErr)
			if returnErr == nil {
				returnErr = wrapped
				return
			}
			returnErr = errors.Join(returnErr, wrapped)
		}
	}()

	expected, err := expectedOwnedFiles(owned, files)
	if err != nil {
		return readiness, err
	}
	const lockPath = "skills-lock.json"
	lockDisplayPath := repositoryDisplayPath(root, lockPath)
	lock, err := readLocalSkillsLock(repositoryRoot, lockPath, lockDisplayPath)
	if err != nil {
		return readiness, err
	}
	externalHashes, err := requiredExternalHashes(lockDisplayPath, lock, external)
	if err != nil {
		return readiness, err
	}

	for _, sharedDirectory := range []string{".agents", ".agents/skills"} {
		missing, err := inspectSharedSkillDirectory(repositoryRoot, sharedDirectory, repositoryDisplayPath(root, sharedDirectory))
		if err != nil {
			return readiness, err
		}
		if missing {
			readiness.MissingOwned = append(readiness.MissingOwned, owned...)
			readiness.MissingExternal = append(readiness.MissingExternal, external...)
			sortRepositoryReadiness(&readiness)
			return readiness, nil
		}
	}

	const skillsRoot = ".agents/skills"
	for _, name := range owned {
		skillRoot := path.Join(skillsRoot, name)
		missing, outdated, err := compareOwnedSkill(
			repositoryRoot,
			skillRoot,
			repositoryDisplayPath(root, skillRoot),
			expected[name],
		)
		if err != nil {
			return readiness, err
		}
		if missing {
			readiness.MissingOwned = append(readiness.MissingOwned, name)
		} else if outdated {
			readiness.OutdatedOwned = append(readiness.OutdatedOwned, name)
		}
	}
	for _, name := range external {
		skillRoot := path.Join(skillsRoot, name)
		skillDisplayPath := repositoryDisplayPath(root, skillRoot)
		info, err := repositoryRoot.Lstat(skillRoot)
		if errors.Is(err, fs.ErrNotExist) {
			readiness.MissingExternal = append(readiness.MissingExternal, name)
			continue
		}
		if err != nil {
			return readiness, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"inspect external skill",
				skillDisplayPath,
				err,
			)
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return readiness, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"inspect external skill",
				skillDisplayPath,
				errors.New("symbolic links are not supported"),
			)
		}
		if !info.IsDir() {
			return readiness, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"inspect external skill",
				skillDisplayPath,
				errors.New("skill root is not a directory"),
			)
		}
		hash, err := skillFolderHash(repositoryRoot.FS(), skillRoot, skillDisplayPath)
		if err != nil {
			return readiness, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"hash external skill",
				skillDisplayPath,
				err,
			)
		}
		if hash != externalHashes[name] {
			readiness.OutdatedExternal = append(readiness.OutdatedExternal, name)
		}
	}

	sortRepositoryReadiness(&readiness)
	return readiness, nil
}

func inspectSharedSkillDirectory(repositoryRoot *os.Root, relative string, display string) (bool, error) {
	info, err := repositoryRoot.Lstat(relative)
	if errors.Is(err, fs.ErrNotExist) {
		return true, nil
	}
	if err != nil {
		return false, repositoryReadinessError("", "inspect repository skill directory", display, err)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return false, repositoryReadinessError(
			"",
			"inspect repository skill directory",
			display,
			errors.New("symbolic links are not supported"),
		)
	}
	if !info.IsDir() {
		return false, repositoryReadinessError(
			"",
			"inspect repository skill directory",
			display,
			errors.New("path is not a directory"),
		)
	}
	return false, nil
}

func repositoryDisplayPath(root string, relative string) string {
	return filepath.Join(root, filepath.FromSlash(relative))
}

func sortRepositoryReadiness(readiness *RepositoryReadiness) {
	sort.Strings(readiness.MissingOwned)
	sort.Strings(readiness.OutdatedOwned)
	sort.Strings(readiness.MissingExternal)
	sort.Strings(readiness.OutdatedExternal)
}

func validateRequiredSkillNames(owned []string, external []string) error {
	seenOwned := make(map[string]struct{}, len(owned))
	for _, name := range owned {
		if !safeSkillName(name) {
			return repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate Roundfix-owned skill name",
				name,
				errors.New("unsafe skill name"),
			)
		}
		if _, exists := seenOwned[name]; exists {
			return repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate Roundfix-owned skill name",
				name,
				errors.New("duplicate skill name"),
			)
		}
		seenOwned[name] = struct{}{}
	}

	seenExternal := make(map[string]struct{}, len(external))
	for _, name := range external {
		if !safeSkillName(name) {
			return repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate external skill name",
				name,
				errors.New("unsafe skill name"),
			)
		}
		if _, exists := seenExternal[name]; exists {
			return repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate external skill name",
				name,
				errors.New("duplicate skill name"),
			)
		}
		if _, overlaps := seenOwned[name]; overlaps {
			return repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate external skill ownership",
				name,
				errors.New("skill overlaps the Roundfix-owned set"),
			)
		}
		seenExternal[name] = struct{}{}
	}
	return nil
}

func safeSkillName(name string) bool {
	return name != "" &&
		name != "." &&
		name != ".." &&
		filepath.Base(name) == name &&
		!strings.ContainsAny(name, `/\`)
}

func expectedOwnedFiles(owned []string, files []File) (map[string]map[string][]byte, error) {
	expected := make(map[string]map[string][]byte, len(owned))
	for _, name := range owned {
		expected[name] = make(map[string][]byte)
	}
	for _, file := range files {
		if !safeSkillName(file.Skill) {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifact owner",
				file.Skill,
				errors.New("unsafe skill name"),
			)
		}
		skillFiles, required := expected[file.Skill]
		if !required {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifact owner",
				file.Skill,
				errors.New("skill is not in the required owned set"),
			)
		}
		prefix := file.Skill + "/"
		if !strings.HasPrefix(file.Path, prefix) {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifact path",
				file.Path,
				errors.New("path is outside its skill directory"),
			)
		}
		relative := strings.TrimPrefix(file.Path, prefix)
		if !safeRelativeArtifactPath(relative) {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifact path",
				file.Path,
				errors.New("unsafe artifact path"),
			)
		}
		if _, duplicate := skillFiles[relative]; duplicate {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifact path",
				file.Path,
				errors.New("duplicate artifact path"),
			)
		}
		skillFiles[relative] = file.Data
	}
	for _, name := range owned {
		if len(expected[name]) == 0 {
			return nil, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"validate embedded skill artifacts",
				name,
				errors.New("required skill has no embedded files"),
			)
		}
	}
	return expected, nil
}

func safeRelativeArtifactPath(path string) bool {
	if path == "" || strings.Contains(path, `\`) || filepath.IsAbs(path) {
		return false
	}
	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(path)))
	return cleaned == path && cleaned != "." && cleaned != ".." && !strings.HasPrefix(cleaned, "../")
}

func readLocalSkillsLock(repositoryRoot *os.Root, path string, displayPath string) (localSkillsLock, error) {
	info, err := repositoryRoot.Lstat(path)
	if err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"inspect skills lock",
			displayPath,
			err,
		)
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return localSkillsLock{}, repositoryReadinessError(
			"",
			"inspect skills lock",
			displayPath,
			errors.New("symbolic links are not supported"),
		)
	}
	if !info.Mode().IsRegular() {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"inspect skills lock",
			displayPath,
			errors.New("lock is not a regular file"),
		)
	}

	data, err := repositoryRoot.ReadFile(path)
	if err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"read skills lock",
			displayPath,
			err,
		)
	}
	var lock localSkillsLock
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&lock); err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"decode skills lock",
			displayPath,
			err,
		)
	}
	var trailing json.RawMessage
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("unexpected trailing JSON value")
		}
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"decode skills lock",
			displayPath,
			err,
		)
	}
	if lock.Version != 1 {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"validate skills lock version",
			displayPath,
			fmt.Errorf("got %d, want 1", lock.Version),
		)
	}
	return lock, nil
}

func requiredExternalHashes(path string, lock localSkillsLock, external []string) (map[string]string, error) {
	hashes := make(map[string]string, len(external))
	for _, name := range external {
		entry, exists := lock.Skills[name]
		if !exists {
			return nil, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate required skills lock entry",
				path,
				fmt.Errorf("external skill %q is missing", name),
			)
		}
		if !lowerHexSHA256(entry.ComputedHash) {
			return nil, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate required skills lock hash",
				path,
				fmt.Errorf("external skill %q computedHash must be 64 lowercase hexadecimal characters", name),
			)
		}
		hashes[name] = entry.ComputedHash
	}
	return hashes, nil
}

func lowerHexSHA256(value string) bool {
	if len(value) != sha256HexLength {
		return false
	}
	for _, char := range value {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') {
			return false
		}
	}
	return true
}

const sha256HexLength = 64

func compareOwnedSkill(
	repositoryRoot *os.Root,
	root string,
	displayRoot string,
	expected map[string][]byte,
) (missing bool, outdated bool, err error) {
	info, lstatErr := repositoryRoot.Lstat(root)
	if errors.Is(lstatErr, fs.ErrNotExist) {
		return true, false, nil
	}
	if lstatErr != nil {
		return false, false, repositoryReadinessError(
			RepositoryOwnershipOwned,
			"inspect Roundfix-owned skill",
			displayRoot,
			lstatErr,
		)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return false, true, nil
	}

	found := make(map[string]struct{}, len(expected))
	tree := repositoryRoot.FS()
	walkErr := fs.WalkDir(tree, root, func(entryPath string, entry fs.DirEntry, walkErr error) error {
		relative := strings.TrimPrefix(entryPath, root+"/")
		if entryPath == root {
			relative = "."
		}
		fullPath := filepath.Join(displayRoot, filepath.FromSlash(relative))
		if walkErr != nil {
			return repositoryReadinessError(RepositoryOwnershipOwned, "walk Roundfix-owned skill artifact", fullPath, walkErr)
		}
		if entryPath == root {
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			outdated = true
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			return repositoryReadinessError(RepositoryOwnershipOwned, "inspect Roundfix-owned skill artifact", fullPath, infoErr)
		}
		if !info.Mode().IsRegular() {
			outdated = true
			return nil
		}
		relative = filepath.ToSlash(relative)
		want, required := expected[relative]
		if !required {
			outdated = true
			return nil
		}
		data, readErr := fs.ReadFile(tree, entryPath)
		if readErr != nil {
			return repositoryReadinessError(RepositoryOwnershipOwned, "read Roundfix-owned skill artifact", fullPath, readErr)
		}
		found[relative] = struct{}{}
		if !bytes.Equal(data, want) {
			outdated = true
		}
		return nil
	})
	if walkErr != nil {
		return false, false, walkErr
	}
	for relative := range expected {
		if _, exists := found[relative]; !exists {
			outdated = true
		}
	}
	return false, outdated, nil
}

func repositoryReadinessError(ownership RepositoryOwnership, operation string, path string, err error) error {
	return &RepositoryReadinessError{
		Ownership: ownership,
		Operation: operation,
		Path:      path,
		Err:       err,
	}
}
