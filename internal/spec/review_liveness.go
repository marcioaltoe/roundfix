package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReviewLiveness is what local Git can prove about an orphan Review Artifact's
// Pull Request. Undecidable remains live for downstream decisions.
type ReviewLiveness string

const (
	ReviewFinished    ReviewLiveness = "finished"
	ReviewLive        ReviewLiveness = "live"
	ReviewUndecidable ReviewLiveness = "undecidable"
)

type reviewRoundMetadata struct {
	HeadBranch string `yaml:"head_branch"`
	HeadSHA    string `yaml:"head_sha"`
}

// ClassifyReview reads the newest Round's recorded head and classifies the
// orphan Review Artifact from local Git only. An undecidable answer is not an
// error: callers leave that Review Artifact live and can report reason.
func ClassifyReview(repoRoot, reviewDir string) (ReviewLiveness, string, error) {
	metadata, roundName, err := newestReviewRoundMetadata(reviewDir)
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("newest Round metadata cannot be read: %v", err), nil
	}

	head := strings.TrimSpace(metadata.HeadSHA)
	if head == "" {
		return ReviewUndecidable, fmt.Sprintf("newest Round %q records no head_sha", roundName), nil
	}
	if !validReviewHead(head) {
		return ReviewUndecidable, fmt.Sprintf("newest Round %q records invalid head_sha %q", roundName, head), nil
	}

	defaultRef, err := reviewDefaultBranch(repoRoot)
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("local Git cannot identify the default branch: %v", err), nil
	}

	resolvedHead, _, err := runReviewGit(repoRoot, "rev-parse", "--verify", head+"^{commit}")
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("local Git cannot resolve recorded head %s: %v", head, err), nil
	}
	resolvedHead = strings.TrimSpace(resolvedHead)

	_, exitCode, ancestryErr := runReviewGit(repoRoot, "merge-base", "--is-ancestor", resolvedHead, defaultRef)
	switch exitCode {
	case 0:
		return ReviewFinished, fmt.Sprintf("recorded head %s is an ancestor of default branch %s", resolvedHead, defaultRef), nil
	case 1:
		// Git answered that the head is not an ancestor. Reachability decides
		// whether the Review Artifact is still live or was abandoned.
	default:
		return ReviewUndecidable, fmt.Sprintf("local Git cannot compare recorded head %s with default branch %s: %v", resolvedHead, defaultRef, ancestryErr), nil
	}

	containingRefs, _, err := runReviewGit(repoRoot, "for-each-ref", "--format=%(refname)", "--contains="+resolvedHead)
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("local Git cannot inspect refs containing recorded head %s: %v", resolvedHead, err), nil
	}
	refs := nonEmptyReviewLines(containingRefs)
	if len(refs) > 0 {
		return ReviewLive, fmt.Sprintf("recorded head %s is reachable from %d local ref(s) and is not an ancestor of default branch %s", resolvedHead, len(refs), defaultRef), nil
	}

	branch := strings.TrimSpace(metadata.HeadBranch)
	if branch == "" {
		return ReviewUndecidable, fmt.Sprintf("recorded head %s is unreachable, but newest Round %q records no head_branch", resolvedHead, roundName), nil
	}
	branchExists, err := reviewBranchExists(repoRoot, branch)
	if err != nil {
		return ReviewUndecidable, fmt.Sprintf("local Git cannot determine whether PR Head Branch %q exists: %v", branch, err), nil
	}
	if branchExists {
		return ReviewUndecidable, fmt.Sprintf("recorded head %s is unreachable, but PR Head Branch %q still exists", resolvedHead, branch), nil
	}
	return ReviewFinished, fmt.Sprintf("recorded head %s is unreachable from local refs and PR Head Branch %q no longer exists", resolvedHead, branch), nil
}

func newestReviewRoundMetadata(reviewDir string) (reviewRoundMetadata, string, error) {
	entries, err := os.ReadDir(reviewDir)
	if err != nil {
		return reviewRoundMetadata{}, "", fmt.Errorf("read Review Artifact directory %q: %w", reviewDir, err)
	}

	newestNumber := -1
	newestName := ""
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "round-") {
			continue
		}
		number, err := strconv.Atoi(strings.TrimPrefix(entry.Name(), "round-"))
		if err != nil || number < 1 {
			continue
		}
		if number == newestNumber && entry.Name() != newestName {
			return reviewRoundMetadata{}, "", fmt.Errorf("Review Artifact directory %q has ambiguous Round %d directories %q and %q", reviewDir, number, newestName, entry.Name())
		}
		if number > newestNumber {
			newestNumber = number
			newestName = entry.Name()
		}
	}
	if newestName == "" {
		return reviewRoundMetadata{}, "", fmt.Errorf("Review Artifact directory %q has no Round metadata", reviewDir)
	}

	metadataPath := filepath.Join(reviewDir, newestName, "round.md")
	content, err := os.ReadFile(metadataPath)
	if err != nil {
		return reviewRoundMetadata{}, newestName, fmt.Errorf("read Round metadata %q: %w", metadataPath, err)
	}
	frontmatter, _, err := splitFrontmatter(content)
	if err != nil {
		return reviewRoundMetadata{}, newestName, fmt.Errorf("parse Round metadata %q: %w", metadataPath, err)
	}
	var metadata reviewRoundMetadata
	if err := yaml.Unmarshal(frontmatter, &metadata); err != nil {
		return reviewRoundMetadata{}, newestName, fmt.Errorf("parse Round metadata frontmatter %q: %w", metadataPath, err)
	}
	return metadata, newestName, nil
}

func reviewDefaultBranch(repoRoot string) (string, error) {
	target, exitCode, err := runReviewGit(repoRoot, "symbolic-ref", "-q", "refs/remotes/origin/HEAD")
	if exitCode == 0 {
		target = strings.TrimSpace(target)
		if branch, ok := strings.CutPrefix(target, "refs/remotes/origin/"); ok && branch != "" {
			localRef := "refs/heads/" + branch
			localExists, localErr := reviewRefExists(repoRoot, localRef)
			if localErr != nil {
				return "", localErr
			}
			if localExists {
				return localRef, nil
			}

			exists, existsErr := reviewRefExists(repoRoot, target)
			if existsErr != nil {
				return "", existsErr
			}
			if exists {
				return target, nil
			}
		}
	} else if exitCode != 1 && exitCode != 128 {
		return "", err
	}

	for _, candidate := range []string{"refs/heads/main", "refs/heads/master"} {
		exists, err := reviewRefExists(repoRoot, candidate)
		if err != nil {
			return "", err
		}
		if exists {
			return candidate, nil
		}
	}
	return "", errors.New("neither origin/HEAD nor a local main or master branch resolves")
}

func reviewRefExists(repoRoot string, ref string) (bool, error) {
	_, exitCode, err := runReviewGit(repoRoot, "show-ref", "--verify", "--quiet", ref)
	switch exitCode {
	case 0:
		return true, nil
	case 1:
		return false, nil
	default:
		return false, err
	}
}

func reviewBranchExists(repoRoot string, branch string) (bool, error) {
	output, _, err := runReviewGit(repoRoot, "for-each-ref", "--format=%(refname)", "refs/heads", "refs/remotes")
	if err != nil {
		return false, err
	}
	for _, ref := range nonEmptyReviewLines(output) {
		if strings.TrimPrefix(ref, "refs/heads/") == branch && strings.HasPrefix(ref, "refs/heads/") {
			return true, nil
		}
		remoteRef, ok := strings.CutPrefix(ref, "refs/remotes/")
		if !ok {
			continue
		}
		_, remoteBranch, ok := strings.Cut(remoteRef, "/")
		if ok && remoteBranch == branch {
			return true, nil
		}
	}
	return false, nil
}

func validReviewHead(head string) bool {
	if len(head) < 4 || len(head) > 64 {
		return false
	}
	for _, char := range head {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}

func nonEmptyReviewLines(output string) []string {
	lines := make([]string, 0)
	for _, line := range strings.Split(output, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func runReviewGit(repoRoot string, args ...string) (string, int, error) {
	gitArgs := append([]string{"-C", repoRoot, "-c", "core.fsmonitor=false"}, args...)
	command := exec.Command("git", gitArgs...)
	command.Env = reviewGitEnv()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	text := strings.TrimRight(stdout.String(), "\n")
	if err == nil {
		return text, 0, nil
	}
	exitCode := -1
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
	}
	detail := strings.TrimSpace(stderr.String())
	if detail == "" {
		detail = strings.TrimSpace(text)
		if detail == "" {
			detail = err.Error()
		}
	}
	return "", exitCode, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
}

func reviewGitEnv() []string {
	env := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "GIT_OPTIONAL_LOCKS=") {
			continue
		}
		env = append(env, entry)
	}
	return append(env, "GIT_OPTIONAL_LOCKS=0")
}
