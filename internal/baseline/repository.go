package baseline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	RepositoryIdentitySchemaVersion = "roundfix/repository-identity/v1"
	RepositorySnapshotSchemaVersion = "roundfix/repository-snapshot/v1"

	repositoryIdentityDigestDomain = RepositoryIdentitySchemaVersion + "\x00"
	repositorySnapshotDigestDomain = RepositorySnapshotSchemaVersion + "\x00"

	maxInventoryCarriers  = 256
	maxInventoryFileBytes = 2 * 1024 * 1024
	maxInventoryBytes     = 8 * 1024 * 1024
)

type PreimageKind string

const (
	PreimageMissing   PreimageKind = "missing"
	PreimageRegular   PreimageKind = "regular"
	PreimageSymlink   PreimageKind = "symlink"
	PreimageDirectory PreimageKind = "directory"
	PreimageSpecial   PreimageKind = "special"
)

type CarrierKind string

const (
	CarrierRegular CarrierKind = "regular"
	CarrierAlias   CarrierKind = "alias"
	CarrierOpaque  CarrierKind = "opaque"
)

// GitRunner executes local, read-only Git inspection commands.
type GitRunner interface {
	RunGit(context.Context, string, ...string) (string, error)
}

// ExecGitRunner executes Git without optional locks, prompts, or fsmonitor.
type ExecGitRunner struct{}

func (ExecGitRunner) RunGit(ctx context.Context, workDir string, args ...string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if workDir == "" {
		workDir = "."
	}
	gitArgs := append([]string{"-C", workDir, "-c", "core.fsmonitor=false"}, args...)
	command := exec.CommandContext(ctx, "git", gitArgs...)
	command.Env = append(os.Environ(), "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", ctxErr
		}
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = strings.TrimSpace(stdout.String())
		}
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RepositoryIdentity is a clone-stable identity derived only from Git
// lineage, never from a checkout path, branch, upstream, or dirty state.
type RepositoryIdentity struct {
	SchemaVersion string   `json:"schemaVersion"`
	ObjectFormat  string   `json:"objectFormat"`
	RootCommits   []string `json:"rootCommits"`
	Digest        string   `json:"digest"`
}

// Preimage records one normalized repository path consulted by inspection or
// reserved as a possible mutation target.
type Preimage struct {
	Path            string       `json:"path"`
	Exists          bool         `json:"exists"`
	Kind            PreimageKind `json:"kind"`
	Mode            uint32       `json:"mode"`
	LinkTarget      string       `json:"linkTarget,omitempty"`
	Bytes           int64        `json:"bytes"`
	ContentIdentity string       `json:"contentIdentity,omitempty"`
}

// SourceEvidence records trusted bytes once, even when multiple instruction
// aliases select the same in-repository target.
type SourceEvidence struct {
	Path            string `json:"path"`
	Bytes           int64  `json:"bytes"`
	ContentIdentity string `json:"contentIdentity"`
}

// InstructionCarrier records one bounded instruction or agent-document
// carrier and the trusted source bytes it selects, when safe.
type InstructionCarrier struct {
	Path            string      `json:"path"`
	Scope           string      `json:"scope"`
	Kind            CarrierKind `json:"kind"`
	TargetPath      string      `json:"targetPath,omitempty"`
	ContentIdentity string      `json:"contentIdentity,omitempty"`
}

// Finding is one deterministic warning or apply-blocking repository fact.
type Finding struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

// RepositorySnapshot is the bounded, path-relative repository preimage used
// by Baseline planning.
type RepositorySnapshot struct {
	SchemaVersion string               `json:"schemaVersion"`
	Preimages     []Preimage           `json:"preimages"`
	Sources       []SourceEvidence     `json:"sources"`
	Carriers      []InstructionCarrier `json:"carriers"`
	Warnings      []Finding            `json:"warnings"`
	Blocking      []Finding            `json:"blocking"`
	Digest        string               `json:"digest"`
}

// RepositoryInspection combines the local root with portable identity and
// snapshot data. Root is intentionally excluded from JSON output.
type RepositoryInspection struct {
	Root     string             `json:"-"`
	Identity RepositoryIdentity `json:"repository"`
	Snapshot RepositorySnapshot `json:"snapshot"`
}

// InventoryRequest identifies extra normalized paths a later planning stage
// may mutate. Root instruction carriers are always included.
type InventoryRequest struct {
	MutablePaths []string
}

// InspectRepository performs the complete read-only Baseline preflight.
func InspectRepository(ctx context.Context, workDir string, runner GitRunner) (RepositoryInspection, error) {
	if err := ctx.Err(); err != nil {
		return RepositoryInspection{}, err
	}
	root, identity, err := inspectRepositoryIdentity(ctx, workDir, runner)
	if err != nil {
		return RepositoryInspection{}, err
	}
	snapshot, err := inspectRepositorySnapshot(root, InventoryRequest{})
	if err != nil {
		return RepositoryInspection{}, fmt.Errorf("inspect bounded repository state: %w", err)
	}
	return RepositoryInspection{Root: root, Identity: identity, Snapshot: snapshot}, nil
}

func inspectRepositoryIdentity(ctx context.Context, workDir string, runner GitRunner) (string, RepositoryIdentity, error) {
	if runner == nil {
		runner = ExecGitRunner{}
	}
	rootText, err := runner.RunGit(ctx, workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("detect Git worktree root: %w", err)
	}
	root, err := filepath.Abs(strings.TrimSpace(rootText))
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("normalize Git worktree root: %w", err)
	}
	root = filepath.Clean(root)
	rootInfo, err := os.Lstat(root)
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("inspect Git worktree root: %w", err)
	}
	if rootInfo.Mode()&fs.ModeSymlink != 0 || !rootInfo.IsDir() {
		return "", RepositoryIdentity{}, fmt.Errorf("inspect Git worktree root: %q must be a real directory", root)
	}

	if _, err := runner.RunGit(ctx, root, "rev-parse", "--verify", "HEAD^{commit}"); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", RepositoryIdentity{}, ctxErr
		}
		return "", RepositoryIdentity{}, fmt.Errorf("Baseline repository requires at least one commit: %w", err)
	}
	objectFormat, err := runner.RunGit(ctx, root, "rev-parse", "--show-object-format")
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("detect Git object format: %w", err)
	}
	objectFormat = strings.TrimSpace(objectFormat)
	objectLength := 0
	switch objectFormat {
	case "sha1":
		objectLength = 40
	case "sha256":
		objectLength = 64
	default:
		return "", RepositoryIdentity{}, fmt.Errorf("detect Git object format: unsupported format %q", objectFormat)
	}
	rootOutput, err := runner.RunGit(ctx, root, "rev-list", "--max-parents=0", "HEAD")
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("detect Git root commits: %w", err)
	}
	rootCommits := strings.Fields(rootOutput)
	if len(rootCommits) == 0 {
		return "", RepositoryIdentity{}, errors.New("detect Git root commits: repository has no root commit")
	}
	sort.Strings(rootCommits)
	rootCommits = compactStrings(rootCommits)
	for _, commit := range rootCommits {
		if len(commit) != objectLength || !isLowerHex(commit) {
			return "", RepositoryIdentity{}, fmt.Errorf("detect Git root commits: invalid %s object ID %q", objectFormat, commit)
		}
	}

	payload := struct {
		SchemaVersion string   `json:"schemaVersion"`
		ObjectFormat  string   `json:"objectFormat"`
		RootCommits   []string `json:"rootCommits"`
	}{
		SchemaVersion: RepositoryIdentitySchemaVersion,
		ObjectFormat:  objectFormat,
		RootCommits:   rootCommits,
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return "", RepositoryIdentity{}, fmt.Errorf("serialize repository identity: %w", err)
	}
	sum := sha256.Sum256(append([]byte(repositoryIdentityDigestDomain), normalized...))
	return root, RepositoryIdentity{
		SchemaVersion: RepositoryIdentitySchemaVersion,
		ObjectFormat:  objectFormat,
		RootCommits:   append([]string(nil), rootCommits...),
		Digest:        "sha256:" + hex.EncodeToString(sum[:]),
	}, nil
}

type inventoryBuilder struct {
	root       *os.Root
	preimages  map[string]Preimage
	sources    map[string]SourceEvidence
	carriers   map[string]InstructionCarrier
	warnings   []Finding
	blocking   []Finding
	totalBytes int64
}

func inspectRepositorySnapshot(root string, request InventoryRequest) (RepositorySnapshot, error) {
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("open repository root: %w", err)
	}
	defer anchored.Close()

	builder := &inventoryBuilder{
		root:      anchored,
		preimages: make(map[string]Preimage),
		sources:   make(map[string]SourceEvidence),
		carriers:  make(map[string]InstructionCarrier),
	}
	if err := builder.inspectAgentDocumentRoots(); err != nil {
		return RepositorySnapshot{}, err
	}
	if err := fs.WalkDir(anchored.FS(), ".", builder.walk); err != nil {
		return RepositorySnapshot{}, fmt.Errorf("walk bounded carriers: %w", err)
	}
	mutablePaths := append([]string{"AGENTS.md", "CLAUDE.md"}, request.MutablePaths...)
	for _, mutablePath := range mutablePaths {
		relative, err := normalizeRepositoryPath(mutablePath)
		if err != nil {
			return RepositorySnapshot{}, fmt.Errorf("normalize mutable path %q: %w", mutablePath, err)
		}
		if _, exists := builder.preimages[relative]; exists {
			continue
		}
		builder.recordPath(relative)
	}
	if len(builder.carriers) > maxInventoryCarriers {
		builder.blocking = append(builder.blocking, Finding{
			Code:    "baseline.inventory.limit.carriers",
			Path:    ".",
			Message: fmt.Sprintf("carrier count %d exceeds %d", len(builder.carriers), maxInventoryCarriers),
		})
	}

	snapshot := RepositorySnapshot{
		SchemaVersion: RepositorySnapshotSchemaVersion,
		Preimages:     sortedPreimages(builder.preimages),
		Sources:       sortedSources(builder.sources),
		Carriers:      sortedCarriers(builder.carriers),
		Warnings:      sortedFindings(builder.warnings),
		Blocking:      sortedFindings(builder.blocking),
	}
	digestPayload := snapshot
	digestPayload.Digest = ""
	normalized, err := json.Marshal(digestPayload)
	if err != nil {
		return RepositorySnapshot{}, fmt.Errorf("serialize repository snapshot: %w", err)
	}
	sum := sha256.Sum256(append([]byte(repositorySnapshotDigestDomain), normalized...))
	snapshot.Digest = "sha256:" + hex.EncodeToString(sum[:])
	return snapshot, nil
}

func (builder *inventoryBuilder) inspectAgentDocumentRoots() error {
	for _, relative := range []string{"docs", "docs/agents"} {
		info, err := builder.root.Lstat(filepath.FromSlash(relative))
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			builder.blocking = append(builder.blocking, Finding{
				Code:    "baseline.inventory.agent-directory-unreadable",
				Path:    relative,
				Message: "bounded agent-document directory cannot be inspected",
			})
			return nil
		}
		if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
			builder.preimages[relative] = preimageFromInfo(relative, info, "")
			builder.blocking = append(builder.blocking, Finding{
				Code:    "baseline.inventory.unsafe-agent-directory",
				Path:    relative,
				Message: "bounded agent-document directory must be a real directory",
			})
			return nil
		}
	}
	return nil
}

func (builder *inventoryBuilder) walk(relative string, entry fs.DirEntry, walkErr error) error {
	if walkErr != nil {
		normalized := filepath.ToSlash(relative)
		builder.blocking = append(builder.blocking, Finding{
			Code:    "baseline.inventory.path-unreadable",
			Path:    normalized,
			Message: "bounded repository path cannot be inspected",
		})
		if entry != nil && entry.IsDir() {
			return fs.SkipDir
		}
		return nil
	}
	if relative == "." {
		return nil
	}
	normalized := filepath.ToSlash(relative)
	if inventoryPathIgnored(normalized) {
		if entry.IsDir() {
			return fs.SkipDir
		}
		return nil
	}
	if entry.IsDir() && path.Base(normalized) != "AGENTS.md" && path.Base(normalized) != "CLAUDE.md" {
		return nil
	}
	if !isBoundedCarrierPath(normalized) {
		return nil
	}

	info, err := entry.Info()
	if err != nil {
		builder.blocking = append(builder.blocking, Finding{
			Code:    "baseline.inventory.carrier-unreadable",
			Path:    normalized,
			Message: "bounded carrier metadata cannot be read",
		})
		return nil
	}
	scope := "nested"
	if !strings.Contains(normalized, "/") {
		scope = "root"
	}
	switch {
	case info.Mode().IsRegular():
		builder.addRegularCarrier(normalized, normalized, scope)
	case info.Mode()&fs.ModeSymlink != 0:
		builder.addAliasCarrier(normalized, scope)
	default:
		builder.preimages[normalized] = preimageFromInfo(normalized, info, "")
		builder.carriers[normalized] = InstructionCarrier{Path: normalized, Scope: scope, Kind: CarrierOpaque}
		if scope == "root" {
			builder.blocking = append(builder.blocking, Finding{
				Code:    "baseline.inventory.special-carrier",
				Path:    normalized,
				Message: "bounded carrier must be a regular file or safe in-repository alias",
			})
		}
	}
	if scope == "nested" {
		builder.warnings = append(builder.warnings, Finding{
			Code:    "baseline.inventory.nested-carrier-conflict",
			Path:    normalized,
			Message: "nested carrier remains unchanged and requires conflict review",
		})
	}
	if entry.IsDir() {
		return fs.SkipDir
	}
	return nil
}

func (builder *inventoryBuilder) addRegularCarrier(carrierPath, sourcePath, scope string) {
	contentIdentity, ok := builder.readTrustedSource(sourcePath, scope)
	if !ok {
		builder.carriers[carrierPath] = InstructionCarrier{Path: carrierPath, Scope: scope, Kind: CarrierOpaque}
		return
	}
	builder.carriers[carrierPath] = InstructionCarrier{
		Path:            carrierPath,
		Scope:           scope,
		Kind:            CarrierRegular,
		TargetPath:      sourcePath,
		ContentIdentity: contentIdentity,
	}
}

func (builder *inventoryBuilder) addAliasCarrier(carrierPath, scope string) {
	info, err := builder.root.Lstat(filepath.FromSlash(carrierPath))
	if err != nil {
		builder.blockUnsafeAlias(carrierPath, scope, "alias metadata cannot be read")
		return
	}
	linkTarget, err := builder.root.Readlink(filepath.FromSlash(carrierPath))
	if err != nil {
		builder.preimages[carrierPath] = preimageFromInfo(carrierPath, info, "")
		builder.blockUnsafeAlias(carrierPath, scope, "alias target cannot be read")
		return
	}
	builder.preimages[carrierPath] = preimageFromInfo(carrierPath, info, filepath.ToSlash(linkTarget))
	targetPath, err := builder.resolveAlias(carrierPath, linkTarget)
	if err != nil {
		builder.carriers[carrierPath] = InstructionCarrier{Path: carrierPath, Scope: scope, Kind: CarrierOpaque}
		builder.blockUnsafeAlias(carrierPath, scope, err.Error())
		return
	}
	contentIdentity, ok := builder.readTrustedSource(targetPath, scope)
	if !ok {
		builder.carriers[carrierPath] = InstructionCarrier{Path: carrierPath, Scope: scope, Kind: CarrierOpaque}
		builder.blockUnsafeAlias(carrierPath, scope, "alias target is not trusted regular-file evidence")
		return
	}
	builder.carriers[carrierPath] = InstructionCarrier{
		Path:            carrierPath,
		Scope:           scope,
		Kind:            CarrierAlias,
		TargetPath:      targetPath,
		ContentIdentity: contentIdentity,
	}
}

func (builder *inventoryBuilder) blockUnsafeAlias(carrierPath, scope, message string) {
	if scope != "root" {
		return
	}
	builder.blocking = append(builder.blocking, Finding{
		Code:    "baseline.inventory.unsafe-alias",
		Path:    carrierPath,
		Message: message,
	})
}

func (builder *inventoryBuilder) resolveAlias(aliasPath, linkTarget string) (string, error) {
	if filepath.IsAbs(linkTarget) || path.IsAbs(filepath.ToSlash(linkTarget)) {
		return "", errors.New("alias target must be repository-relative")
	}
	pending := path.Join(path.Dir(aliasPath), filepath.ToSlash(linkTarget))
	if !repositoryPathIsSafe(pending) {
		return "", errors.New("alias target escapes the repository")
	}
	seen := map[string]struct{}{aliasPath: {}}
	for resolutions := 0; resolutions <= maxInventoryCarriers; resolutions++ {
		parts := strings.Split(pending, "/")
		restarted := false
		for index := range parts {
			candidate := path.Join(parts[:index+1]...)
			info, err := builder.root.Lstat(filepath.FromSlash(candidate))
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					builder.preimages[candidate] = Preimage{Path: candidate, Kind: PreimageMissing}
					return "", fmt.Errorf("alias target %q does not exist", candidate)
				}
				return "", fmt.Errorf("alias target %q cannot be inspected", candidate)
			}
			link := ""
			if info.Mode()&fs.ModeSymlink != 0 {
				link, err = builder.root.Readlink(filepath.FromSlash(candidate))
				if err != nil {
					return "", fmt.Errorf("alias target %q cannot be read", candidate)
				}
			}
			builder.preimages[candidate] = preimageFromInfo(candidate, info, filepath.ToSlash(link))
			if info.Mode()&fs.ModeSymlink != 0 {
				if _, cycle := seen[candidate]; cycle {
					return "", fmt.Errorf("alias target cycle at %q", candidate)
				}
				seen[candidate] = struct{}{}
				if filepath.IsAbs(link) || path.IsAbs(filepath.ToSlash(link)) {
					return "", fmt.Errorf("alias target %q is absolute", candidate)
				}
				target := path.Join(path.Dir(candidate), filepath.ToSlash(link))
				if index+1 < len(parts) {
					target = path.Join(target, path.Join(parts[index+1:]...))
				}
				if !repositoryPathIsSafe(target) {
					return "", fmt.Errorf("alias target %q escapes the repository", candidate)
				}
				pending = target
				restarted = true
				break
			}
			if index+1 < len(parts) && !info.IsDir() {
				return "", fmt.Errorf("alias path component %q is not a directory", candidate)
			}
			if index+1 == len(parts) && !info.Mode().IsRegular() {
				return "", fmt.Errorf("alias target %q is not a regular file", candidate)
			}
		}
		if !restarted {
			return pending, nil
		}
	}
	return "", errors.New("alias target has too many symbolic-link resolutions")
}

func (builder *inventoryBuilder) readTrustedSource(relative, scope string) (string, bool) {
	if source, exists := builder.sources[relative]; exists {
		return source.ContentIdentity, true
	}
	info, err := builder.root.Lstat(filepath.FromSlash(relative))
	if err != nil {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.carrier-unreadable",
			Path:    relative,
			Message: "bounded carrier cannot be inspected",
		})
		return "", false
	}
	if !info.Mode().IsRegular() {
		builder.preimages[relative] = preimageFromInfo(relative, info, "")
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.special-carrier",
			Path:    relative,
			Message: "trusted carrier source must be a regular file",
		})
		return "", false
	}
	if info.Size() > maxInventoryFileBytes {
		builder.preimages[relative] = preimageFromInfo(relative, info, "")
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.limit.file-bytes",
			Path:    relative,
			Message: fmt.Sprintf("carrier size %d exceeds %d", info.Size(), maxInventoryFileBytes),
		})
		return "", false
	}
	file, err := builder.root.Open(filepath.FromSlash(relative))
	if err != nil {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.carrier-unreadable",
			Path:    relative,
			Message: "bounded carrier bytes cannot be read",
		})
		return "", false
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.carrier-changed",
			Path:    relative,
			Message: "bounded carrier changed while it was inspected",
		})
		return "", false
	}
	data, err := io.ReadAll(io.LimitReader(file, maxInventoryFileBytes+1))
	if err != nil {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.carrier-unreadable",
			Path:    relative,
			Message: "bounded carrier bytes cannot be read",
		})
		return "", false
	}
	if len(data) > maxInventoryFileBytes {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.limit.file-bytes",
			Path:    relative,
			Message: fmt.Sprintf("carrier bytes exceed %d", maxInventoryFileBytes),
		})
		return "", false
	}
	builder.totalBytes += int64(len(data))
	if builder.totalBytes > maxInventoryBytes {
		builder.addCarrierBlocking(scope, Finding{
			Code:    "baseline.inventory.limit.total-bytes",
			Path:    relative,
			Message: fmt.Sprintf("carrier bytes exceed %d", maxInventoryBytes),
		})
		return "", false
	}
	sum := sha256.Sum256(data)
	identity := "sha256:" + hex.EncodeToString(sum[:])
	builder.preimages[relative] = Preimage{
		Path:            relative,
		Exists:          true,
		Kind:            PreimageRegular,
		Mode:            uint32(info.Mode()),
		Bytes:           int64(len(data)),
		ContentIdentity: identity,
	}
	builder.sources[relative] = SourceEvidence{
		Path:            relative,
		Bytes:           int64(len(data)),
		ContentIdentity: identity,
	}
	return identity, true
}

func (builder *inventoryBuilder) addCarrierBlocking(scope string, finding Finding) {
	if scope == "root" {
		builder.blocking = append(builder.blocking, finding)
	}
}

func (builder *inventoryBuilder) recordPath(relative string) {
	info, err := builder.root.Lstat(filepath.FromSlash(relative))
	if errors.Is(err, fs.ErrNotExist) {
		builder.preimages[relative] = Preimage{Path: relative, Kind: PreimageMissing}
		return
	}
	if err != nil {
		builder.blocking = append(builder.blocking, Finding{
			Code:    "baseline.inventory.mutable-path-unreadable",
			Path:    relative,
			Message: "possible mutation target cannot be inspected",
		})
		return
	}
	linkTarget := ""
	if info.Mode()&fs.ModeSymlink != 0 {
		if target, linkErr := builder.root.Readlink(filepath.FromSlash(relative)); linkErr == nil {
			linkTarget = filepath.ToSlash(target)
		}
	}
	builder.preimages[relative] = preimageFromInfo(relative, info, linkTarget)
}

func preimageFromInfo(relative string, info fs.FileInfo, linkTarget string) Preimage {
	kind := PreimageSpecial
	switch {
	case info.Mode().IsRegular():
		kind = PreimageRegular
	case info.Mode()&fs.ModeSymlink != 0:
		kind = PreimageSymlink
	case info.IsDir():
		kind = PreimageDirectory
	}
	return Preimage{
		Path:       relative,
		Exists:     true,
		Kind:       kind,
		Mode:       uint32(info.Mode()),
		LinkTarget: linkTarget,
		Bytes:      info.Size(),
	}
}

func isBoundedCarrierPath(relative string) bool {
	base := path.Base(relative)
	if base == "AGENTS.md" || base == "CLAUDE.md" {
		return true
	}
	return strings.HasPrefix(relative, "docs/agents/")
}

func inventoryPathIgnored(relative string) bool {
	parts := strings.Split(relative, "/")
	for _, part := range parts {
		switch part {
		case ".git", ".hg", ".svn", ".venv", "__pycache__", "node_modules", "vendor", "venv":
			return true
		}
	}
	return relative == "skills" || strings.HasPrefix(relative, "skills/") ||
		relative == ".agents/skills" || strings.HasPrefix(relative, ".agents/skills/")
}

func normalizeRepositoryPath(value string) (string, error) {
	value = filepath.ToSlash(strings.TrimSpace(value))
	if !repositoryPathIsSafe(value) {
		return "", errors.New("path must be normalized and repository-relative")
	}
	return value, nil
}

func repositoryPathIsSafe(value string) bool {
	if value == "" || value == "." || path.IsAbs(value) || strings.Contains(value, `\`) {
		return false
	}
	cleaned := path.Clean(value)
	return cleaned == value && cleaned != ".." && !strings.HasPrefix(cleaned, "../")
}

func sortedPreimages(values map[string]Preimage) []Preimage {
	paths := sortedMapKeys(values)
	result := make([]Preimage, 0, len(paths))
	for _, relative := range paths {
		result = append(result, values[relative])
	}
	return result
}

func sortedSources(values map[string]SourceEvidence) []SourceEvidence {
	paths := sortedMapKeys(values)
	result := make([]SourceEvidence, 0, len(paths))
	for _, relative := range paths {
		result = append(result, values[relative])
	}
	return result
}

func sortedCarriers(values map[string]InstructionCarrier) []InstructionCarrier {
	paths := sortedMapKeys(values)
	result := make([]InstructionCarrier, 0, len(paths))
	for _, relative := range paths {
		result = append(result, values[relative])
	}
	return result
}

func sortedFindings(findings []Finding) []Finding {
	result := append([]Finding(nil), findings...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Path != result[j].Path {
			return result[i].Path < result[j].Path
		}
		if result[i].Code != result[j].Code {
			return result[i].Code < result[j].Code
		}
		return result[i].Message < result[j].Message
	})
	return result
}

func sortedMapKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func compactStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	write := 1
	for index := 1; index < len(values); index++ {
		if values[index] == values[write-1] {
			continue
		}
		values[write] = values[index]
		write++
	}
	return values[:write]
}

func isLowerHex(value string) bool {
	for _, character := range value {
		if character >= '0' && character <= '9' || character >= 'a' && character <= 'f' {
			continue
		}
		return false
	}
	return true
}
