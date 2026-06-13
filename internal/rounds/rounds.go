package rounds

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"roundfix/internal/reviewsource"
)

const (
	StatusPending    = "pending"
	StatusValid      = "valid"
	StatusInvalid    = "invalid"
	StatusResolved   = "resolved"
	StatusDuplicated = "duplicated"
	StatusFailed     = "failed"
)

type PersistRequest struct {
	ArtifactDir    string
	Source         string
	PRNumber       string
	HeadRepository string
	HeadBranch     string
	HeadSHA        string
	Round          int
	CreatedAt      time.Time
	Items          []reviewsource.ReviewItem
}

type PersistResult struct {
	Round      int
	RoundDir   string
	IssuePaths []string
	Reused     bool
}

type SelectRequest struct {
	ArtifactDir    string
	PRNumber       string
	HeadRepository string
	HeadBranch     string
	Round          int
}

type SelectResult struct {
	Issues []Issue
	Rounds []int
}

type BatchRequest struct {
	Issues    []Issue
	BatchSize int
}

type BatchPlan struct {
	Batches    []Batch
	Duplicates []DuplicateAssociation
}

type Batch struct {
	Number int
	Issues []Issue
}

type DuplicateAssociation struct {
	Fingerprint string
	Older       Issue
	Newest      Issue
}

type Issue struct {
	Path                    string
	Title                   string
	Source                  string
	PRNumber                string
	Round                   int
	RoundCreatedAt          time.Time
	Status                  string
	HeadRepository          string
	HeadBranch              string
	HeadSHA                 string
	File                    string
	Line                    int
	Severity                string
	Author                  string
	SourceRef               string
	ReviewHash              string
	DuplicateOf             string
	SourceReviewID          string
	SourceReviewSubmittedAt time.Time
}

type NoCompatibleArtifactsError struct {
	PRNumber       string
	HeadRepository string
	HeadBranch     string
	Round          int
}

type AmbiguousDuplicateError struct {
	Fingerprint string
	Issues      []Issue
}

func (err AmbiguousDuplicateError) Error() string {
	paths := make([]string, 0, len(err.Issues))
	for _, issue := range err.Issues {
		paths = append(paths, issue.Path)
	}
	return fmt.Sprintf("ambiguous newest Review Issue for fingerprint %q; fix artifact timestamps before running Roundfix again: %s", err.Fingerprint, strings.Join(paths, ", "))
}

func (err NoCompatibleArtifactsError) Error() string {
	roundScope := "all"
	if err.Round > 0 {
		roundScope = fmt.Sprintf("%03d", err.Round)
	}
	return fmt.Sprintf("no downloaded Unresolved Review Issues in Compatible Artifacts for Open Pull Request #%s, Head Repository %q, PR Head Branch %q, Round %s; run `roundfix fetch --source coderabbit --pr %s` or use `roundfix watch`", err.PRNumber, err.HeadRepository, err.HeadBranch, roundScope, err.PRNumber)
}

type issueFrontmatter struct {
	Source                  string `yaml:"source"`
	PR                      string `yaml:"pr"`
	Round                   int    `yaml:"round"`
	RoundCreatedAt          string `yaml:"round_created_at"`
	Status                  string `yaml:"status"`
	HeadRepository          string `yaml:"head_repository"`
	HeadBranch              string `yaml:"head_branch"`
	HeadSHA                 string `yaml:"head_sha"`
	File                    string `yaml:"file"`
	Line                    int    `yaml:"line"`
	Severity                string `yaml:"severity"`
	Author                  string `yaml:"author"`
	SourceRef               string `yaml:"source_ref"`
	ReviewHash              string `yaml:"review_hash"`
	DuplicateOf             string `yaml:"duplicate_of"`
	SourceReviewID          string `yaml:"source_review_id"`
	SourceReviewSubmittedAt string `yaml:"source_review_submitted_at"`
}

type roundFrontmatter struct {
	Source         string `yaml:"source"`
	PR             string `yaml:"pr"`
	Round          int    `yaml:"round"`
	RoundCreatedAt string `yaml:"round_created_at"`
	HeadRepository string `yaml:"head_repository"`
	HeadBranch     string `yaml:"head_branch"`
	HeadSHA        string `yaml:"head_sha"`
	IssueCount     int    `yaml:"issue_count"`
}

func PersistRound(ctx context.Context, req PersistRequest) (PersistResult, error) {
	if err := ctx.Err(); err != nil {
		return PersistResult{}, err
	}
	if err := validatePersistRequest(req); err != nil {
		return PersistResult{}, err
	}
	createdAt := req.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	roundNumber := req.Round
	if roundNumber == 0 {
		existing, found, err := findMatchingRound(ctx, req)
		if err != nil {
			return PersistResult{}, err
		}
		if found {
			return existing, nil
		}
		next, err := NextRoundNumber(req.ArtifactDir, req.PRNumber)
		if err != nil {
			return PersistResult{}, err
		}
		roundNumber = next
	}

	roundDir := filepath.Join(req.ArtifactDir, "reviews", "pr-"+req.PRNumber, fmt.Sprintf("round-%03d", roundNumber))
	if _, err := os.Stat(roundDir); err == nil {
		return PersistResult{}, fmt.Errorf("Round artifact directory %q already exists", roundDir)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return PersistResult{}, fmt.Errorf("stat Round artifact directory %q: %w", roundDir, err)
	}
	if err := os.MkdirAll(roundDir, 0o755); err != nil {
		return PersistResult{}, fmt.Errorf("create Round artifact directory %q: %w", roundDir, err)
	}

	if err := writeRoundMetadata(roundDir, req, roundNumber, createdAt); err != nil {
		return PersistResult{}, err
	}
	issuePaths := make([]string, 0, len(req.Items))
	for index, item := range req.Items {
		path := filepath.Join(roundDir, fmt.Sprintf("issue_%03d.md", index+1))
		if err := writeIssue(path, req, roundNumber, createdAt, index+1, item); err != nil {
			return PersistResult{}, err
		}
		issuePaths = append(issuePaths, path)
	}
	return PersistResult{
		Round:      roundNumber,
		RoundDir:   roundDir,
		IssuePaths: issuePaths,
	}, nil
}

func SelectCompatibleIssues(ctx context.Context, req SelectRequest) (SelectResult, error) {
	if err := ctx.Err(); err != nil {
		return SelectResult{}, err
	}
	if err := validateSelectRequest(req); err != nil {
		return SelectResult{}, err
	}

	roundDirs, err := compatibleRoundDirs(req)
	if err != nil {
		return SelectResult{}, err
	}

	result := SelectResult{}
	roundSet := map[int]bool{}
	for _, roundDir := range roundDirs {
		if err := ctx.Err(); err != nil {
			return SelectResult{}, err
		}
		issuePaths, err := filepath.Glob(filepath.Join(roundDir, "issue_*.md"))
		if err != nil {
			return SelectResult{}, fmt.Errorf("find Review Issue artifacts in %q: %w", roundDir, err)
		}
		sort.Strings(issuePaths)
		for _, issuePath := range issuePaths {
			issue, err := ParseIssue(issuePath)
			if err != nil {
				return SelectResult{}, err
			}
			if !matchesSelectRequest(issue, req) {
				continue
			}
			if !AllowedStatus(issue.Status) {
				return SelectResult{}, fmt.Errorf("Review Issue artifact %q has unsupported status %q", issuePath, issue.Status)
			}
			if IsTerminalStatus(issue.Status) {
				continue
			}
			result.Issues = append(result.Issues, issue)
			if !roundSet[issue.Round] {
				roundSet[issue.Round] = true
				result.Rounds = append(result.Rounds, issue.Round)
			}
		}
	}
	sort.Ints(result.Rounds)
	if len(result.Issues) == 0 {
		return SelectResult{}, NoCompatibleArtifactsError{
			PRNumber:       req.PRNumber,
			HeadRepository: req.HeadRepository,
			HeadBranch:     req.HeadBranch,
			Round:          req.Round,
		}
	}
	return result, nil
}

func ParseIssue(path string) (Issue, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Issue{}, fmt.Errorf("read Review Issue artifact %q: %w", path, err)
	}
	frontmatterBytes, body, err := splitMarkdown(content)
	if err != nil {
		return Issue{}, fmt.Errorf("parse Review Issue artifact %q: %w", path, err)
	}
	var frontmatter issueFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return Issue{}, fmt.Errorf("parse Review Issue frontmatter %q: %w", path, err)
	}
	roundCreatedAt, err := parseOptionalTime(frontmatter.RoundCreatedAt)
	if err != nil {
		return Issue{}, fmt.Errorf("parse round_created_at in %q: %w", path, err)
	}
	sourceReviewSubmittedAt, err := parseOptionalTime(frontmatter.SourceReviewSubmittedAt)
	if err != nil {
		return Issue{}, fmt.Errorf("parse source_review_submitted_at in %q: %w", path, err)
	}
	return Issue{
		Path:                    path,
		Title:                   parseIssueTitle(body),
		Source:                  frontmatter.Source,
		PRNumber:                frontmatter.PR,
		Round:                   frontmatter.Round,
		RoundCreatedAt:          roundCreatedAt,
		Status:                  frontmatter.Status,
		HeadRepository:          frontmatter.HeadRepository,
		HeadBranch:              frontmatter.HeadBranch,
		HeadSHA:                 frontmatter.HeadSHA,
		File:                    frontmatter.File,
		Line:                    frontmatter.Line,
		Severity:                frontmatter.Severity,
		Author:                  frontmatter.Author,
		SourceRef:               frontmatter.SourceRef,
		ReviewHash:              frontmatter.ReviewHash,
		DuplicateOf:             frontmatter.DuplicateOf,
		SourceReviewID:          frontmatter.SourceReviewID,
		SourceReviewSubmittedAt: sourceReviewSubmittedAt,
	}, nil
}

func PlanBatches(req BatchRequest) (BatchPlan, error) {
	if req.BatchSize < 1 {
		return BatchPlan{}, errors.New("Batch size must be greater than 0")
	}
	byFingerprint := make(map[string][]Issue)
	for _, issue := range req.Issues {
		fingerprint, err := ReviewIssueFingerprint(issue)
		if err != nil {
			return BatchPlan{}, err
		}
		byFingerprint[fingerprint] = append(byFingerprint[fingerprint], issue)
	}

	fingerprints := make([]string, 0, len(byFingerprint))
	for fingerprint := range byFingerprint {
		fingerprints = append(fingerprints, fingerprint)
	}
	sort.Strings(fingerprints)

	assigned := make([]Issue, 0, len(fingerprints))
	duplicates := []DuplicateAssociation{}
	for _, fingerprint := range fingerprints {
		group := byFingerprint[fingerprint]
		newest, err := newestIssue(fingerprint, group)
		if err != nil {
			return BatchPlan{}, err
		}
		assigned = append(assigned, newest)
		for _, issue := range group {
			if issue.Path == newest.Path {
				continue
			}
			duplicates = append(duplicates, DuplicateAssociation{
				Fingerprint: fingerprint,
				Older:       issue,
				Newest:      newest,
			})
		}
	}
	sort.SliceStable(assigned, func(i, j int) bool {
		if assigned[i].Round != assigned[j].Round {
			return assigned[i].Round < assigned[j].Round
		}
		return assigned[i].Path < assigned[j].Path
	})
	sort.SliceStable(duplicates, func(i, j int) bool {
		if duplicates[i].Fingerprint != duplicates[j].Fingerprint {
			return duplicates[i].Fingerprint < duplicates[j].Fingerprint
		}
		return duplicates[i].Older.Path < duplicates[j].Older.Path
	})

	return BatchPlan{
		Batches:    makeBatches(assigned, req.BatchSize),
		Duplicates: duplicates,
	}, nil
}

func ReviewIssueFingerprint(issue Issue) (string, error) {
	sourceRef := strings.TrimSpace(issue.SourceRef)
	if sourceRef != "" {
		return "source_ref:" + sourceRef, nil
	}
	reviewHash := strings.TrimSpace(issue.ReviewHash)
	if reviewHash != "" {
		return "review_hash:" + reviewHash, nil
	}
	return "", fmt.Errorf("Review Issue artifact %q is missing source_ref and review_hash for fingerprinting", issue.Path)
}

func reviewItemFingerprint(item reviewsource.ReviewItem) (string, error) {
	sourceRef := strings.TrimSpace(item.SourceRef)
	if sourceRef != "" {
		return "source_ref:" + sourceRef, nil
	}
	reviewHash := strings.TrimSpace(item.ReviewHash)
	if reviewHash != "" {
		return "review_hash:" + reviewHash, nil
	}
	return "", errors.New("Review Source item is missing source_ref and review_hash for fingerprinting")
}

func findMatchingRound(ctx context.Context, req PersistRequest) (PersistResult, bool, error) {
	if err := ctx.Err(); err != nil {
		return PersistResult{}, false, err
	}
	incoming, err := reviewItemFingerprints(req.Items)
	if err != nil {
		return PersistResult{}, false, err
	}
	roundDirs, err := compatibleRoundDirs(SelectRequest{
		ArtifactDir:    req.ArtifactDir,
		PRNumber:       req.PRNumber,
		HeadRepository: req.HeadRepository,
		HeadBranch:     req.HeadBranch,
	})
	if err != nil {
		return PersistResult{}, false, err
	}

	var match PersistResult
	found := false
	for _, roundDir := range roundDirs {
		if err := ctx.Err(); err != nil {
			return PersistResult{}, false, err
		}
		result, ok, err := matchingRoundResult(roundDir, req, incoming)
		if err != nil {
			return PersistResult{}, false, err
		}
		if ok {
			match = result
			found = true
		}
	}
	return match, found, nil
}

func MarkDuplicatedAfterTerminal(ctx context.Context, associations []DuplicateAssociation) (int, error) {
	marked := 0
	for _, association := range associations {
		if err := ctx.Err(); err != nil {
			return marked, err
		}
		newest, err := ParseIssue(association.Newest.Path)
		if err != nil {
			return marked, err
		}
		if newest.Status != StatusResolved && newest.Status != StatusInvalid {
			continue
		}
		if err := SetIssueStatus(association.Older.Path, StatusDuplicated, association.Newest.Path); err != nil {
			return marked, err
		}
		marked++
	}
	return marked, nil
}

func SetIssueStatus(path string, status string, duplicateOf string) error {
	if !AllowedStatus(status) {
		return fmt.Errorf("Review Issue status %q is not allowed", status)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Review Issue artifact %q: %w", path, err)
	}
	frontmatterBytes, body, err := splitMarkdown(content)
	if err != nil {
		return fmt.Errorf("parse Review Issue artifact %q: %w", path, err)
	}
	var frontmatter issueFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return fmt.Errorf("parse Review Issue frontmatter %q: %w", path, err)
	}
	frontmatter.Status = status
	frontmatter.DuplicateOf = duplicateOf
	updated, err := renderMarkdown(frontmatter, string(body))
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("write Review Issue artifact %q: %w", path, err)
	}
	return nil
}

func NextRoundNumber(artifactDir string, prNumber string) (int, error) {
	root := filepath.Join(artifactDir, "reviews", "pr-"+prNumber)
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return 1, nil
	}
	if err != nil {
		return 0, fmt.Errorf("read pull request artifact directory %q: %w", root, err)
	}
	maxRound := 0
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "round-") {
			continue
		}
		number, err := strconv.Atoi(strings.TrimPrefix(entry.Name(), "round-"))
		if err != nil {
			continue
		}
		if number > maxRound {
			maxRound = number
		}
	}
	return maxRound + 1, nil
}

func AllowedStatus(status string) bool {
	switch status {
	case StatusPending, StatusValid, StatusInvalid, StatusResolved, StatusDuplicated, StatusFailed:
		return true
	default:
		return false
	}
}

// IsSettledStatus reports whether a Review Issue needs no further Agent
// work in the current Batch: every Terminal status plus failed, which
// stays Unresolved and is retried in a later Round.
func IsSettledStatus(status string) bool {
	return IsTerminalStatus(status) || status == StatusFailed
}

func IsTerminalStatus(status string) bool {
	switch status {
	case StatusResolved, StatusInvalid, StatusDuplicated:
		return true
	default:
		return false
	}
}

func validatePersistRequest(req PersistRequest) error {
	required := map[string]string{
		"Artifact Directory": req.ArtifactDir,
		"source":             req.Source,
		"pull request":       req.PRNumber,
		"Head Repository":    req.HeadRepository,
		"PR Head Branch":     req.HeadBranch,
		"HEAD":               req.HeadSHA,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required to persist a Round", label)
		}
	}
	if req.Round < 0 {
		return errors.New("Round must be auto or greater than 0")
	}
	return nil
}

func validateSelectRequest(req SelectRequest) error {
	required := map[string]string{
		"Artifact Directory": req.ArtifactDir,
		"pull request":       req.PRNumber,
		"Head Repository":    req.HeadRepository,
		"PR Head Branch":     req.HeadBranch,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required to select Compatible Artifacts", label)
		}
	}
	if req.Round < 0 {
		return errors.New("Round selector must be all or greater than 0")
	}
	return nil
}

func compatibleRoundDirs(req SelectRequest) ([]string, error) {
	prDir := filepath.Join(req.ArtifactDir, "reviews", "pr-"+req.PRNumber)
	if req.Round > 0 {
		roundDir := filepath.Join(prDir, fmt.Sprintf("round-%03d", req.Round))
		if _, err := os.Stat(roundDir); errors.Is(err, os.ErrNotExist) {
			return nil, nil
		} else if err != nil {
			return nil, fmt.Errorf("stat Compatible Artifact Round %q: %w", roundDir, err)
		}
		return []string{roundDir}, nil
	}

	entries, err := os.ReadDir(prDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pull request artifact directory %q: %w", prDir, err)
	}
	roundDirs := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "round-") {
			continue
		}
		if _, err := strconv.Atoi(strings.TrimPrefix(entry.Name(), "round-")); err != nil {
			continue
		}
		roundDirs = append(roundDirs, filepath.Join(prDir, entry.Name()))
	}
	sort.Strings(roundDirs)
	return roundDirs, nil
}

func matchingRoundResult(roundDir string, req PersistRequest, incoming []string) (PersistResult, bool, error) {
	metadata, ok, err := readRoundMetadata(roundDir)
	if err != nil || !ok {
		return PersistResult{}, false, err
	}
	if metadata.Source != req.Source ||
		metadata.PR != req.PRNumber ||
		metadata.HeadRepository != req.HeadRepository ||
		metadata.HeadBranch != req.HeadBranch ||
		metadata.HeadSHA != req.HeadSHA {
		return PersistResult{}, false, nil
	}

	issuePaths, fingerprints, err := roundIssueFingerprints(roundDir)
	if err != nil {
		return PersistResult{}, false, err
	}
	if metadata.IssueCount != len(incoming) || !sameFingerprints(incoming, fingerprints) {
		return PersistResult{}, false, nil
	}
	return PersistResult{
		Round:      metadata.Round,
		RoundDir:   roundDir,
		IssuePaths: issuePaths,
		Reused:     true,
	}, true, nil
}

func readRoundMetadata(roundDir string) (roundFrontmatter, bool, error) {
	path := filepath.Join(roundDir, "round.md")
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return roundFrontmatter{}, false, nil
	}
	if err != nil {
		return roundFrontmatter{}, false, fmt.Errorf("read Round metadata %q: %w", path, err)
	}
	frontmatterBytes, err := extractFrontmatter(content)
	if err != nil {
		return roundFrontmatter{}, false, fmt.Errorf("parse Round metadata %q: %w", path, err)
	}
	var frontmatter roundFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return roundFrontmatter{}, false, fmt.Errorf("parse Round metadata frontmatter %q: %w", path, err)
	}
	return frontmatter, true, nil
}

func roundIssueFingerprints(roundDir string) ([]string, []string, error) {
	issuePaths, err := filepath.Glob(filepath.Join(roundDir, "issue_*.md"))
	if err != nil {
		return nil, nil, fmt.Errorf("find Review Issue artifacts in %q: %w", roundDir, err)
	}
	sort.Strings(issuePaths)
	fingerprints := make([]string, 0, len(issuePaths))
	for _, issuePath := range issuePaths {
		issue, err := ParseIssue(issuePath)
		if err != nil {
			return nil, nil, err
		}
		fingerprint, err := ReviewIssueFingerprint(issue)
		if err != nil {
			return nil, nil, err
		}
		fingerprints = append(fingerprints, fingerprint)
	}
	sort.Strings(fingerprints)
	return issuePaths, fingerprints, nil
}

func reviewItemFingerprints(items []reviewsource.ReviewItem) ([]string, error) {
	fingerprints := make([]string, 0, len(items))
	for _, item := range items {
		fingerprint, err := reviewItemFingerprint(item)
		if err != nil {
			return nil, err
		}
		fingerprints = append(fingerprints, fingerprint)
	}
	sort.Strings(fingerprints)
	return fingerprints, nil
}

func sameFingerprints(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func matchesSelectRequest(issue Issue, req SelectRequest) bool {
	if issue.PRNumber != req.PRNumber {
		return false
	}
	if issue.HeadRepository != req.HeadRepository {
		return false
	}
	if issue.HeadBranch != req.HeadBranch {
		return false
	}
	if req.Round > 0 && issue.Round != req.Round {
		return false
	}
	return true
}

func newestIssue(fingerprint string, issues []Issue) (Issue, error) {
	if len(issues) == 0 {
		return Issue{}, fmt.Errorf("empty duplicate group for fingerprint %q", fingerprint)
	}
	newest := issues[0]
	ambiguous := []Issue{newest}
	for _, issue := range issues[1:] {
		switch compareIssueNewness(issue, newest) {
		case 1:
			newest = issue
			ambiguous = []Issue{issue}
		case 0:
			ambiguous = append(ambiguous, issue)
		}
	}
	if len(ambiguous) > 1 {
		sort.SliceStable(ambiguous, func(i, j int) bool {
			return ambiguous[i].Path < ambiguous[j].Path
		})
		return Issue{}, AmbiguousDuplicateError{Fingerprint: fingerprint, Issues: ambiguous}
	}
	return newest, nil
}

func compareIssueNewness(left Issue, right Issue) int {
	if left.Round > right.Round {
		return 1
	}
	if left.Round < right.Round {
		return -1
	}
	if timeAfter(left.SourceReviewSubmittedAt, right.SourceReviewSubmittedAt) {
		return 1
	}
	if timeAfter(right.SourceReviewSubmittedAt, left.SourceReviewSubmittedAt) {
		return -1
	}
	if timeAfter(left.RoundCreatedAt, right.RoundCreatedAt) {
		return 1
	}
	if timeAfter(right.RoundCreatedAt, left.RoundCreatedAt) {
		return -1
	}
	return 0
}

func timeAfter(left time.Time, right time.Time) bool {
	if left.IsZero() || right.IsZero() {
		return false
	}
	return left.After(right)
}

func makeBatches(issues []Issue, batchSize int) []Batch {
	batches := make([]Batch, 0, (len(issues)+batchSize-1)/batchSize)
	for start := 0; start < len(issues); start += batchSize {
		end := start + batchSize
		if end > len(issues) {
			end = len(issues)
		}
		batches = append(batches, Batch{
			Number: len(batches) + 1,
			Issues: append([]Issue(nil), issues[start:end]...),
		})
	}
	return batches
}

func writeRoundMetadata(roundDir string, req PersistRequest, roundNumber int, createdAt time.Time) error {
	frontmatter := roundFrontmatter{
		Source:         req.Source,
		PR:             req.PRNumber,
		Round:          roundNumber,
		RoundCreatedAt: formatTime(createdAt),
		HeadRepository: req.HeadRepository,
		HeadBranch:     req.HeadBranch,
		HeadSHA:        req.HeadSHA,
		IssueCount:     len(req.Items),
	}
	content, err := renderMarkdown(frontmatter, fmt.Sprintf("# Round %03d\n\nFetched %d Review Issue(s).\n", roundNumber, len(req.Items)))
	if err != nil {
		return err
	}
	path := filepath.Join(roundDir, "round.md")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write Round metadata %q: %w", path, err)
	}
	return nil
}

func writeIssue(path string, req PersistRequest, roundNumber int, createdAt time.Time, issueNumber int, item reviewsource.ReviewItem) error {
	frontmatter := issueFrontmatter{
		Source:                  req.Source,
		PR:                      req.PRNumber,
		Round:                   roundNumber,
		RoundCreatedAt:          formatTime(createdAt),
		Status:                  StatusPending,
		HeadRepository:          req.HeadRepository,
		HeadBranch:              req.HeadBranch,
		HeadSHA:                 req.HeadSHA,
		File:                    item.File,
		Line:                    item.Line,
		Severity:                item.Severity,
		Author:                  item.Author,
		SourceRef:               item.SourceRef,
		ReviewHash:              item.ReviewHash,
		DuplicateOf:             "",
		SourceReviewID:          item.SourceReviewID,
		SourceReviewSubmittedAt: formatOptionalTime(item.SourceReviewSubmittedAt),
	}
	title := item.Title
	if strings.TrimSpace(title) == "" {
		title = fmt.Sprintf("Issue %03d", issueNumber)
	}
	body := fmt.Sprintf("# Issue %03d: %s\n\n## Review Comment\n\n%s\n\n## Triage\n\n- Decision: `UNREVIEWED`\n- Notes:\n", issueNumber, title, item.Body)
	content, err := renderMarkdown(frontmatter, body)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("write Review Issue artifact %q: %w", path, err)
	}
	return nil
}

func renderMarkdown(frontmatter any, body string) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(frontmatter)
	if err != nil {
		return nil, fmt.Errorf("marshal Review Issue frontmatter: %w", err)
	}
	var builder strings.Builder
	builder.WriteString("---\n")
	builder.Write(yamlBytes)
	builder.WriteString("---\n\n")
	builder.WriteString(body)
	return []byte(builder.String()), nil
}

func extractFrontmatter(content []byte) ([]byte, error) {
	frontmatter, _, err := splitMarkdown(content)
	return frontmatter, err
}

func splitMarkdown(content []byte) ([]byte, []byte, error) {
	text := string(content)
	const opening = "---\n"
	if !strings.HasPrefix(text, opening) {
		return nil, nil, errors.New("missing YAML frontmatter opening marker")
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return nil, nil, errors.New("missing YAML frontmatter closing marker")
	}
	bodyStart := end + len("\n---")
	if strings.HasPrefix(rest[bodyStart:], "\n") {
		bodyStart++
	}
	return []byte(rest[:end]), []byte(rest[bodyStart:]), nil
}

func parseIssueTitle(body []byte) string {
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "# ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		if strings.HasPrefix(title, "Issue ") {
			if marker := strings.Index(title, ":"); marker >= 0 {
				return strings.TrimSpace(title[marker+1:])
			}
		}
		return title
	}
	return ""
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func formatOptionalTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return formatTime(value)
}

func parseOptionalTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
}
