package speccheck

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

// SourceCondition names the repository fact that withheld authored command
// execution. It stays separate from the shared refusal code so operators can
// distinguish the repair without parsing prose.
type SourceCondition string

const (
	SourceConditionOutOfTree         SourceCondition = "out-of-tree-spec-root"
	SourceConditionUntrackedArtifact SourceCondition = "untracked-artifact"
	SourceConditionModifiedArtifact  SourceCondition = "modified-artifact"
	SourceConditionUnresolved        SourceCondition = "unresolved-source"
)

// AuthoredCommandSourceRequest identifies the Task artifact and exact command
// text an execution entry point is about to pass to a shell. Revision defaults
// to HEAD in the delivery target repository.
type AuthoredCommandSourceRequest struct {
	RepoRoot  string
	SpecsRoot string
	SpecSlug  string
	Artifact  string
	Revision  string
	Commands  []string
}

// SourceUntrustedError is the shared refusal emitted at every authored-command
// execution entry point.
type SourceUntrustedError struct {
	Code      string
	Condition SourceCondition
	Artifact  string
	Detail    string
}

func (err *SourceUntrustedError) Error() string {
	if err == nil {
		return ""
	}
	detail := strings.TrimSpace(err.Detail)
	if detail == "" {
		detail = "authored command source did not satisfy committed provenance"
	}
	return fmt.Sprintf("%s: %s for %s: %s", err.Code, err.Condition, err.Artifact, detail)
}

// AuthoredCommandDigest returns the representation execution approvals store
// for one exact command string.
func AuthoredCommandDigest(command string) string {
	digest := sha256.Sum256([]byte(command))
	return "sha256:" + hex.EncodeToString(digest[:])
}

type authoredCommandSource struct {
	condition         SourceCondition
	detail            string
	targetRoot        string
	targetRevision    string
	authorizationPath string
	repositoryRoot    string
	repositoryID      string
	revision          string
	artifact          string
	commands          []string
}

// AuthorizeAuthoredCommands permits a shell boundary only for committed
// authored provenance or an exact execution approval already committed in the
// delivery target's ancestry. It performs repository reads only.
func AuthorizeAuthoredCommands(ctx context.Context, request AuthoredCommandSourceRequest) error {
	source, err := inspectAuthoredCommandSource(ctx, request)
	if err != nil {
		return err
	}
	if source.condition == "" {
		return nil
	}
	if source.condition != SourceConditionUnresolved {
		approved, approvalErr := committedExecutionApproval(ctx, source)
		if approvalErr != nil {
			return approvalErr
		}
		if approved {
			return nil
		}
	}
	return sourceRefusal(source.condition, request.Artifact, source.detail)
}

func inspectAuthoredCommandSource(ctx context.Context, request AuthoredCommandSourceRequest) (authoredCommandSource, error) {
	source := authoredCommandSource{commands: append([]string(nil), request.Commands...)}
	targetRoot, targetID, err := resolveGitRepository(ctx, request.RepoRoot)
	if err != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve delivery target repository: %v", err))
	}
	source.targetRoot = targetRoot
	revision := strings.TrimSpace(request.Revision)
	if revision == "" {
		revision = "HEAD"
	}
	source.targetRevision, err = resolveSourceRevision(ctx, targetRoot, revision)
	if err != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve delivery target revision %q: %v", revision, err))
	}
	if !validSourceSpecSlug(request.SpecSlug) {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, "consuming Spec slug is not a single canonical path component")
	}

	specsRoot, err := canonicalExistingPath(request.SpecsRoot)
	if err != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve Spec Root: %v", err))
	}
	artifactRelative, ok := cleanSourceArtifact(request.Artifact)
	if !ok {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, "carrying artifact is not relative to the Spec Root")
	}
	artifactPath, err := canonicalExistingPath(filepath.Join(specsRoot, filepath.FromSlash(artifactRelative)))
	if err != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("read carrying artifact: %v", err))
	}
	if !pathInside(artifactPath, specsRoot) {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, "carrying artifact resolves outside the Spec Root")
	}
	workingBytes, err := os.ReadFile(artifactPath)
	if err != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("read carrying artifact: %v", err))
	}
	if commands := verificationCommandsFromArtifact(workingBytes); !slices.Equal(commands, request.Commands) {
		source.condition = SourceConditionModifiedArtifact
		source.detail = "requested command text differs from the carrying artifact"
	}

	insideTarget := pathInside(specsRoot, targetRoot)
	if insideTarget {
		specsRelative, relErr := filepath.Rel(targetRoot, specsRoot)
		if relErr != nil {
			return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("make Spec Root repository-relative: %v", relErr))
		}
		source.authorizationPath = filepath.ToSlash(filepath.Join(specsRelative, request.SpecSlug, "_authorization.md"))
		relative, relErr := filepath.Rel(targetRoot, artifactPath)
		if relErr != nil || !pathInside(artifactPath, targetRoot) {
			return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, "make carrying artifact repository-relative")
		}
		source.repositoryRoot = targetRoot
		source.repositoryID = targetID
		source.artifact = filepath.ToSlash(relative)
		committedBytes, tracked, readErr := readSourceBlob(ctx, targetRoot, source.targetRevision, source.artifact)
		if readErr != nil {
			return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("read committed carrying artifact: %v", readErr))
		}
		if !tracked {
			source.condition = SourceConditionUntrackedArtifact
			source.detail = "carrying artifact is absent from the resolved revision"
			source.revision = source.targetRevision
			return source, nil
		}
		source.revision, err = sourceArtifactRevision(ctx, targetRoot, source.targetRevision, source.artifact)
		if err != nil {
			return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve carrying artifact revision: %v", err))
		}
		if source.condition == "" && !bytes.Equal(authoredProjection(workingBytes), authoredProjection(committedBytes)) {
			source.condition = SourceConditionModifiedArtifact
			source.detail = "authored projection differs from the resolved revision"
		}
		return source, nil
	}

	source.authorizationPath = filepath.ToSlash(filepath.Join("docs", "specs", request.SpecSlug, "_authorization.md"))
	source.condition = SourceConditionOutOfTree
	source.detail = "Spec Root resolves outside the delivery target Git tree"
	sourceRoot, sourceID, sourceErr := resolveGitRepository(ctx, filepath.Dir(artifactPath))
	if sourceErr != nil {
		return source, nil
	}
	relative, relErr := filepath.Rel(sourceRoot, artifactPath)
	if relErr != nil || !pathInside(artifactPath, sourceRoot) {
		return source, nil
	}
	source.repositoryRoot = sourceRoot
	source.repositoryID = sourceID
	source.artifact = filepath.ToSlash(relative)
	sourceHead, headErr := resolveSourceRevision(ctx, sourceRoot, "HEAD")
	if headErr != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve external source revision: %v", headErr))
	}
	_, tracked, readErr := readSourceBlob(ctx, sourceRoot, sourceHead, source.artifact)
	if readErr != nil {
		return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("read external source object: %v", readErr))
	}
	if tracked {
		source.revision, err = sourceArtifactRevision(ctx, sourceRoot, sourceHead, source.artifact)
		if err != nil {
			return source, sourceRefusal(SourceConditionUnresolved, request.Artifact, fmt.Sprintf("resolve external artifact revision: %v", err))
		}
	} else {
		source.revision = sourceHead
	}
	return source, nil
}

type executionApproval struct {
	Repository    string `yaml:"repository"`
	Revision      string `yaml:"revision"`
	Artifact      string `yaml:"artifact"`
	CommandDigest string `yaml:"command_digest"`
}

func committedExecutionApproval(ctx context.Context, source authoredCommandSource) (bool, error) {
	if source.repositoryRoot == "" || source.repositoryID == "" || source.revision == "" || source.artifact == "" {
		return false, nil
	}
	resolution := spec.ReadAuthorization(ctx, spec.AuthorizationReadRequest{
		RepoRoot:   source.targetRoot,
		RecordPath: source.authorizationPath,
		Revision:   source.targetRevision,
		Role:       spec.AuthorizationRoleSpec,
		AskingSpec: filepath.Base(filepath.Dir(source.authorizationPath)),
	})
	if resolution.Outcome != spec.AuthorizationGranted {
		return false, nil
	}
	content, tracked, err := readSourceBlob(ctx, source.targetRoot, source.targetRevision, source.authorizationPath)
	if err != nil {
		return false, sourceRefusal(SourceConditionUnresolved, source.artifact, fmt.Sprintf("read committed execution approval: %v", err))
	}
	if !tracked {
		return false, nil
	}
	frontmatter, err := sourceFrontmatter(content)
	if err != nil {
		return false, nil
	}
	var document struct {
		ExecutionApprovals []executionApproval `yaml:"execution_approvals"`
	}
	if err := yaml.Unmarshal(frontmatter, &document); err != nil {
		return false, nil
	}
	for _, command := range source.commands {
		approved, approvalErr := commandHasExecutionApproval(ctx, source, document.ExecutionApprovals, command)
		if approvalErr != nil {
			return false, approvalErr
		}
		if !approved {
			return false, nil
		}
	}
	return true, nil
}

func commandHasExecutionApproval(ctx context.Context, source authoredCommandSource, approvals []executionApproval, command string) (bool, error) {
	wantDigest := AuthoredCommandDigest(command)
	for _, approval := range approvals {
		approvalArtifact, canonical := cleanSourceArtifact(approval.Artifact)
		if !canonical || approvalArtifact != approval.Artifact || approvalArtifact != source.artifact || strings.TrimSpace(approval.CommandDigest) != wantDigest {
			continue
		}
		approvalRoot := strings.TrimSpace(approval.Repository)
		if approvalRoot == "." {
			approvalRoot = source.targetRoot
		} else if approvalRoot != "" && !filepath.IsAbs(approvalRoot) {
			approvalRoot = filepath.Join(source.targetRoot, filepath.FromSlash(approvalRoot))
		}
		_, approvalID, err := resolveGitRepository(ctx, approvalRoot)
		if err != nil || approvalID != source.repositoryID {
			continue
		}
		approvedRevision, err := resolveSourceRevision(ctx, source.repositoryRoot, approval.Revision)
		if err != nil {
			return false, sourceRefusal(SourceConditionUnresolved, source.artifact, fmt.Sprintf("resolve approved source revision %q: %v", approval.Revision, err))
		}
		if approvedRevision == source.revision {
			return true, nil
		}
	}
	return false, nil
}

func authoredProjection(content []byte) []byte {
	lines := bytes.SplitAfter(content, []byte("\n"))
	projection := make([]byte, 0, len(content))
	inFrontmatter := len(lines) > 0 && strings.TrimSpace(string(lines[0])) == "---"
	for index, line := range lines {
		trimmed := strings.TrimSpace(string(line))
		if index > 0 && inFrontmatter && trimmed == "---" {
			inFrontmatter = false
			projection = append(projection, line...)
			continue
		}
		if inFrontmatter && strings.HasPrefix(trimmed, "status:") {
			continue
		}
		if !inFrontmatter && trimmed == "## Result" {
			break
		}
		projection = append(projection, line...)
	}
	projection = bytes.TrimRight(projection, "\r\n")
	return append(projection, '\n')
}

func verificationCommandsFromArtifact(content []byte) []string {
	var commands []string
	inVerification := false
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "## Result" {
			break
		}
		if strings.HasPrefix(trimmed, "# ") {
			inVerification = false
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			inVerification = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Verification"
			continue
		}
		if !inVerification || !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		start := strings.IndexByte(trimmed, '`')
		if start < 0 {
			continue
		}
		end := strings.IndexByte(trimmed[start+1:], '`')
		if end >= 0 && end > 0 {
			commands = append(commands, trimmed[start+1:start+1+end])
		}
	}
	return commands
}

func readSourceBlob(ctx context.Context, repoRoot, revision, artifact string) ([]byte, bool, error) {
	entry, err := runSourceGit(ctx, repoRoot, "ls-tree", revision, "--", artifact)
	if err != nil {
		return nil, false, err
	}
	if len(bytes.TrimSpace(entry)) == 0 {
		return nil, false, nil
	}
	metadata, _, found := bytes.Cut(entry, []byte{'\t'})
	if !found {
		return nil, false, fmt.Errorf("parse Git tree entry for %q", artifact)
	}
	fields := strings.Fields(string(metadata))
	if len(fields) < 2 || fields[0] == "120000" || fields[1] != "blob" {
		return nil, false, fmt.Errorf("carrying artifact %q is not a regular Git blob", artifact)
	}
	content, err := runSourceGit(ctx, repoRoot, "cat-file", "blob", revision+":"+artifact)
	if err != nil {
		return nil, false, err
	}
	return content, true, nil
}

func sourceArtifactRevision(ctx context.Context, repoRoot, revision, artifact string) (string, error) {
	output, err := runSourceGit(ctx, repoRoot, "log", "-1", "--format=%H", revision, "--", artifact)
	if err != nil {
		return "", err
	}
	resolved := strings.TrimSpace(string(output))
	if resolved == "" || strings.Contains(resolved, "\n") {
		return "", errors.New("Git returned no artifact revision")
	}
	return resolved, nil
}

func resolveSourceRevision(ctx context.Context, repoRoot, revision string) (string, error) {
	revision = strings.TrimSpace(revision)
	if revision == "" || strings.ContainsAny(revision, "\x00\r\n") {
		return "", errors.New("revision is empty or contains a control character")
	}
	output, err := runSourceGit(ctx, repoRoot, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		return "", err
	}
	resolved := strings.TrimSpace(string(output))
	if resolved == "" || strings.Contains(resolved, "\n") {
		return "", fmt.Errorf("Git returned invalid revision %q", resolved)
	}
	return resolved, nil
}

func resolveGitRepository(ctx context.Context, directory string) (string, string, error) {
	directory = strings.TrimSpace(directory)
	if directory == "" {
		return "", "", errors.New("repository path is empty")
	}
	rootOutput, err := runSourceGit(ctx, directory, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", err
	}
	root, err := canonicalExistingPath(strings.TrimSpace(string(rootOutput)))
	if err != nil {
		return "", "", err
	}
	commonOutput, err := runSourceGit(ctx, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", "", err
	}
	common, err := canonicalExistingPath(strings.TrimSpace(string(commonOutput)))
	if err != nil {
		return "", "", err
	}
	return root, filepath.ToSlash(common), nil
}

func runSourceGit(ctx context.Context, repoRoot string, args ...string) ([]byte, error) {
	gitArgs := append([]string{"-C", repoRoot, "-c", "core.fsmonitor=false"}, args...)
	command := exec.CommandContext(ctx, "git", gitArgs...)
	command.Env = sourceGitEnvironment()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return append([]byte(nil), stdout.Bytes()...), nil
}

func sourceGitEnvironment() []string {
	environment := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "GIT_OPTIONAL_LOCKS=") {
			continue
		}
		environment = append(environment, entry)
	}
	return append(environment, "GIT_OPTIONAL_LOCKS=0")
}

func sourceFrontmatter(content []byte) ([]byte, error) {
	normalized := bytes.ReplaceAll(content, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(normalized, []byte("---\n")) {
		return nil, errors.New("missing authorization frontmatter")
	}
	rest := normalized[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---"))
	if end < 0 {
		return nil, errors.New("authorization frontmatter is not closed")
	}
	return rest[:end], nil
}

func cleanSourceArtifact(value string) (string, bool) {
	if strings.TrimSpace(value) != value || value == "" || filepath.IsAbs(value) || strings.ContainsAny(value, "\\\x00\r\n") {
		return "", false
	}
	clean := filepath.Clean(filepath.FromSlash(value))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.ToSlash(clean), true
}

func validSourceSpecSlug(value string) bool {
	return strings.TrimSpace(value) == value && value != "" && filepath.Base(value) == value && !strings.ContainsAny(value, "\\/\x00\r\n")
}

func canonicalExistingPath(value string) (string, error) {
	absolute, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(absolute)
}

func pathInside(candidate, root string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func sourceRefusal(condition SourceCondition, artifact, detail string) *SourceUntrustedError {
	return &SourceUntrustedError{
		Code:      CodeSourceUntrusted,
		Condition: condition,
		Artifact:  filepath.ToSlash(artifact),
		Detail:    detail,
	}
}
