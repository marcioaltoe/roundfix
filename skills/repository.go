package skills

import (
	"bytes"
	"context"
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
	Owned            []OwnedSkillReadiness
	MissingExternal  []string
	OutdatedExternal []string
}

func (readiness RepositoryReadiness) Ready() bool {
	if len(readiness.MissingOwned) != 0 ||
		len(readiness.MissingExternal) != 0 ||
		len(readiness.OutdatedExternal) != 0 {
		return false
	}
	for _, owned := range readiness.Owned {
		if owned.State != ReadinessSatisfies && owned.State != ReadinessUnversioned {
			return false
		}
	}
	return true
}

type RepositoryReadinessError struct {
	Ownership RepositoryOwnership
	Operation string
	Path      string
	Err       error
	// MissingExternal names every required external skill the lock does not
	// declare, so a caller can render one remediation per skill instead of
	// stopping at the first gap.
	MissingExternal []string
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

func CheckRepository(ctx context.Context, root string) (RepositoryReadiness, error) {
	return CheckRepositoryWithExternal(ctx, root, Recommended())
}

// CheckRepositoryWithExternal checks the repository against the caller's
// external skill requirement and the running binary's embedded owned bundle.
func CheckRepositoryWithExternal(ctx context.Context, root string, external []string) (RepositoryReadiness, error) {
	if err := ctx.Err(); err != nil {
		return RepositoryReadiness{}, fmt.Errorf("check repository skill set: %w", err)
	}
	owned := Names()
	if err := validateRequiredSkillNames(owned, external); err != nil {
		return RepositoryReadiness{
			OwnedRequired:    len(owned),
			ExternalRequired: len(external),
		}, err
	}
	return checkRepositoryWithReadiness(ctx, root, owned, external, Readiness)
}

type readinessComparison func(SkillVersion, string) (ReadinessState, error)

func checkRepositoryWithReadiness(ctx context.Context, root string, owned []string, external []string, compare readinessComparison) (readiness RepositoryReadiness, returnErr error) {
	readiness = RepositoryReadiness{
		OwnedRequired:    len(owned),
		ExternalRequired: len(external),
	}
	if err := ctx.Err(); err != nil {
		return readiness, repositoryReadinessError("", "inspect repository skill set", root, err)
	}
	if err := validateRequiredSkillNames(owned, external); err != nil {
		return readiness, err
	}
	root = strings.TrimSpace(root)
	if root == "" {
		return readiness, repositoryReadinessError("", "validate repository root", "", errors.New("path is empty"))
	}

	if err := ctx.Err(); err != nil {
		return readiness, repositoryReadinessError("", "open repository root", root, err)
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

	var externalHashes map[string]string
	if len(external) > 0 {
		const lockPath = "skills-lock.json"
		lockDisplayPath := repositoryDisplayPath(root, lockPath)
		lock, err := readLocalSkillsLock(ctx, repositoryRoot, lockPath, lockDisplayPath)
		if err != nil {
			return readiness, err
		}
		externalHashes, err = requiredExternalHashes(lockDisplayPath, lock, external)
		if err != nil {
			return readiness, err
		}
	}

	for _, sharedDirectory := range []string{".agents", ".agents/skills"} {
		missing, err := inspectSharedSkillDirectory(ctx, repositoryRoot, sharedDirectory, repositoryDisplayPath(root, sharedDirectory))
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
		if err := ctx.Err(); err != nil {
			return readiness, repositoryReadinessError(RepositoryOwnershipOwned, "inspect Roundfix-owned skill", repositoryDisplayPath(root, skillsRoot), err)
		}
		skillRoot := path.Join(skillsRoot, name)
		minimum, ok := ownedSkillMinimumVersions[name]
		if !ok {
			return readiness, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"resolve Roundfix-owned skill minimum",
				name,
				errors.New("minimum version is not declared"),
			)
		}
		missing, ownedReadiness, err := inspectOwnedSkillReadiness(
			ctx,
			repositoryRoot,
			name,
			skillRoot,
			repositoryDisplayPath(root, skillRoot),
			minimum,
			compare,
		)
		if err != nil {
			return readiness, err
		}
		if missing {
			readiness.MissingOwned = append(readiness.MissingOwned, name)
		} else {
			readiness.Owned = append(readiness.Owned, ownedReadiness)
		}
	}
	for _, name := range external {
		if err := ctx.Err(); err != nil {
			return readiness, repositoryReadinessError(RepositoryOwnershipExternal, "inspect external skill", repositoryDisplayPath(root, skillsRoot), err)
		}
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
		hash, err := skillFolderHash(ctx, repositoryRoot.FS(), skillRoot, skillDisplayPath)
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

func inspectSharedSkillDirectory(ctx context.Context, repositoryRoot *os.Root, relative string, display string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, repositoryReadinessError("", "inspect repository skill directory", display, err)
	}
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

func repositoryDisplayPath(root string, relativePath string) string {
	return filepath.Join(root, filepath.FromSlash(relativePath))
}

func sortRepositoryReadiness(readiness *RepositoryReadiness) {
	sort.Strings(readiness.MissingOwned)
	sort.Slice(readiness.Owned, func(i, j int) bool {
		return readiness.Owned[i].Skill < readiness.Owned[j].Skill
	})
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

func readLocalSkillsLock(ctx context.Context, repositoryRoot *os.Root, lockPath string, displayPath string) (localSkillsLock, error) {
	if err := ctx.Err(); err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"inspect skills lock",
			displayPath,
			err,
		)
	}
	info, err := repositoryRoot.Lstat(lockPath)
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
			RepositoryOwnershipExternal,
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

	if err := ctx.Err(); err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"read skills lock",
			displayPath,
			err,
		)
	}
	data, err := repositoryRoot.ReadFile(lockPath)
	if err != nil {
		return localSkillsLock{}, repositoryReadinessError(
			RepositoryOwnershipExternal,
			"read skills lock",
			displayPath,
			err,
		)
	}
	if err := ctx.Err(); err != nil {
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

func requiredExternalHashes(lockPath string, lock localSkillsLock, external []string) (map[string]string, error) {
	hashes := make(map[string]string, len(external))
	var missing []string
	for _, name := range external {
		entry, exists := lock.Skills[name]
		if !exists {
			missing = append(missing, name)
			continue
		}
		if !lowerHexSHA256(entry.ComputedHash) {
			return nil, repositoryReadinessError(
				RepositoryOwnershipExternal,
				"validate required skills lock hash",
				lockPath,
				fmt.Errorf("external skill %q computedHash must be 64 lowercase hexadecimal characters", name),
			)
		}
		hashes[name] = entry.ComputedHash
	}
	if len(missing) > 0 {
		return nil, &RepositoryReadinessError{
			Ownership:       RepositoryOwnershipExternal,
			Operation:       "validate required skills lock entry",
			Path:            lockPath,
			Err:             fmt.Errorf("external skills are missing: %s", strings.Join(missing, ", ")),
			MissingExternal: append([]string(nil), missing...),
		}
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

func inspectOwnedSkillReadiness(
	ctx context.Context,
	repositoryRoot *os.Root,
	name string,
	root string,
	displayRoot string,
	minimum string,
	compare readinessComparison,
) (missing bool, readiness OwnedSkillReadiness, err error) {
	source := path.Join(root, "SKILL.md")
	displaySource := filepath.Join(displayRoot, "SKILL.md")
	readiness = OwnedSkillReadiness{
		Skill:   name,
		Minimum: minimum,
		Source:  displaySource,
	}
	if err := ctx.Err(); err != nil {
		return false, readiness, repositoryReadinessError(
			RepositoryOwnershipOwned,
			"inspect Roundfix-owned skill",
			displayRoot,
			err,
		)
	}
	info, lstatErr := repositoryRoot.Lstat(root)
	if errors.Is(lstatErr, fs.ErrNotExist) {
		return true, readiness, nil
	}
	if lstatErr != nil {
		return false, readiness, repositoryReadinessError(
			RepositoryOwnershipOwned,
			"inspect Roundfix-owned skill",
			displayRoot,
			lstatErr,
		)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return classifyUnversionedOwnedSkill(readiness, "skill source is unreachable", minimum, compare)
	}

	if err := ctx.Err(); err != nil {
		return false, readiness, repositoryReadinessError(
			RepositoryOwnershipOwned,
			"inspect Roundfix-owned skill version",
			displaySource,
			err,
		)
	}
	info, lstatErr = repositoryRoot.Lstat(source)
	if lstatErr != nil || info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
		detail := "skill version source is unreachable"
		if lstatErr != nil && !errors.Is(lstatErr, fs.ErrNotExist) {
			detail = lstatErr.Error()
		}
		return classifyUnversionedOwnedSkill(readiness, detail, minimum, compare)
	}
	data, readErr := repositoryRoot.ReadFile(source)
	if readErr != nil {
		if err := ctx.Err(); err != nil {
			return false, readiness, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"read Roundfix-owned skill version",
				displaySource,
				err,
			)
		}
		return classifyUnversionedOwnedSkill(readiness, readErr.Error(), minimum, compare)
	}
	metadata, ok := parseSkillFrontmatter(string(data))
	declared := ""
	if ok {
		declared = strings.TrimSpace(metadata.Version)
	} else {
		declared = "invalid-frontmatter"
	}
	state, compareErr := compare(SkillVersion{Declared: declared, Source: displaySource}, minimum)
	readiness.Found = strings.TrimSpace(metadata.Version)
	readiness.State = state
	if compareErr != nil {
		if !errors.Is(compareErr, ErrSkillVersionUnresolvable) {
			return false, readiness, repositoryReadinessError(
				RepositoryOwnershipOwned,
				"compare Roundfix-owned skill version",
				displaySource,
				compareErr,
			)
		}
		readiness.Detail = compareErr.Error()
	}
	return false, readiness, nil
}

func classifyUnversionedOwnedSkill(
	readiness OwnedSkillReadiness,
	detail string,
	minimum string,
	compare readinessComparison,
) (bool, OwnedSkillReadiness, error) {
	state, compareErr := compare(SkillVersion{}, minimum)
	if state != ReadinessUnversioned || !errors.Is(compareErr, ErrSkillVersionUnresolvable) {
		if compareErr == nil {
			compareErr = fmt.Errorf("comparison returned state %q", state)
		}
		return false, readiness, repositoryReadinessError(
			RepositoryOwnershipOwned,
			"compare unreachable Roundfix-owned skill version",
			readiness.Source,
			compareErr,
		)
	}
	readiness.State = state
	readiness.Detail = detail
	return false, readiness, nil
}

func repositoryReadinessError(ownership RepositoryOwnership, operation string, path string, err error) error {
	return &RepositoryReadinessError{
		Ownership: ownership,
		Operation: operation,
		Path:      path,
		Err:       err,
	}
}
