package worktree

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

// MergedHead records the head a Spec's pull request merged at.
type MergedHead struct {
	SpecSlug     string
	TargetBranch string
	Head         string
	MergeCommit  string
	PullRequest  string
}

// InspectTerminalRunMerged classifies a terminal Run using ancestry, QA
// supersession, and, when those proofs miss, the represented work at a merged
// head.
func InspectTerminalRunMerged(
	ctx context.Context,
	run store.Run,
	merged []MergedHead,
) (RunWorktreeReconciliation, error) {
	return inspectTerminalRunMerged(ctx, execGitRunner{}, run, merged)
}

type mergedHeadSource struct {
	head          string
	ref           string
	label         string
	defaultBranch bool
	archived      bool
}

type mergedHeadManifest struct {
	Graph struct {
		Nodes []struct {
			ID   string `yaml:"id"`
			File string `yaml:"file"`
		} `yaml:"nodes"`
	} `yaml:"graph"`
}

func chooseMergedHead(
	ctx context.Context,
	runner gitRunner,
	run store.Run,
	gitRoot string,
	merged []MergedHead,
) (mergedHeadSource, bool) {
	cleanSlug, err := cleanPathSegment(run.SpecSlug)
	if err != nil || cleanSlug != run.SpecSlug {
		return mergedHeadSource{}, false
	}
	matching := make([]MergedHead, 0)
	distinct := make(map[string]struct{})
	for _, record := range merged {
		if record.SpecSlug != run.SpecSlug {
			continue
		}
		matching = append(matching, record)
		distinct[record.Head] = struct{}{}
	}
	if len(distinct) == 1 && len(matching) != 0 {
		head := matching[0].Head
		if validMergedHeadRevision(head) {
			if _, err := runner.Run(ctx, gitRoot, "cat-file", "-e", head+"^{commit}"); err == nil {
				return mergedHeadSource{
					head:  head,
					ref:   head,
					label: mergedHeadRecordLabel(matching, head),
				}, true
			}
		}
	}

	defaultBranch, defaultHead, resolved := resolveDefaultBranchHead(ctx, runner, gitRoot)
	if !resolved {
		return mergedHeadSource{}, false
	}
	_, archiveErr := newestQAReportAtHeadInDirectories(
		ctx,
		runner,
		gitRoot,
		defaultHead,
		[]string{archivedQAReportDirectory(run.SpecSlug)},
	)
	return mergedHeadSource{
		head:          defaultHead,
		ref:           defaultBranch,
		label:         fmt.Sprintf("default branch %q", defaultBranch),
		defaultBranch: true,
		archived:      archiveErr == nil,
	}, true
}

func validMergedHeadRevision(head string) bool {
	if len(head) != 40 && len(head) != 64 {
		return false
	}
	for _, char := range head {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return false
		}
	}
	return true
}

func mergedHeadRecordLabel(records []MergedHead, head string) string {
	pullRequest := ""
	for _, record := range records {
		if record.Head != head {
			continue
		}
		candidate := strings.TrimPrefix(strings.TrimSpace(record.PullRequest), "#")
		if candidate != "" && decimalString(candidate) {
			pullRequest = candidate
			break
		}
	}
	shortHead := head
	if len(shortHead) > 12 {
		shortHead = shortHead[:12]
	}
	if pullRequest != "" {
		return fmt.Sprintf("pull request #%s merged head %s", pullRequest, shortHead)
	}
	return "merged head " + shortHead
}

func decimalString(value string) bool {
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return value != ""
}

func inspectRunAtMergedHead(
	ctx context.Context,
	runner gitRunner,
	run store.Run,
	gitRoot string,
	result RunWorktreeReconciliation,
	worktreePresent bool,
	runBranchPresent bool,
	merged []MergedHead,
	source mergedHeadSource,
) RunWorktreeReconciliation {
	commitsOutput, err := runner.Run(ctx, gitRoot, "rev-list", result.RunHead, "^"+source.head)
	if err != nil {
		result.State = ReconciliationUnintegrated
		result.Reason = boundedReconciliationReason(fmt.Sprintf(
			"Run Branch representation could not be inspected at %s",
			source.label,
		))
		return result
	}
	commits := strings.Fields(commitsOutput)
	if len(commits) == 0 {
		result.State = ReconciliationSafe
		result.Reason = boundedReconciliationReason("Run Branch is contained in " + source.label)
		result.evidence = newMergedHeadReconciliationEvidence(
			run,
			gitRoot,
			result,
			worktreePresent,
			runBranchPresent,
			merged,
			source,
		)
		return result
	}

	archived := specArchivedAtMergedHead(ctx, runner, gitRoot, source.head, run.SpecSlug)
	specPaths := make(map[string]struct{})
	taskCommits := 0
	qaCommits := 0
	otherPaths := make(map[string]struct{})
	for _, commit := range commits {
		message, messageErr := runner.Run(ctx, gitRoot, "show", "-s", "--format=%B", commit)
		specSlug, specErr := commitTrailer(ctx, runner, gitRoot, commit, "Roundfix-Spec")
		taskID, taskErr := commitTrailer(ctx, runner, gitRoot, commit, "Roundfix-Task")
		if messageErr != nil || specErr != nil || taskErr != nil {
			return unrepresentedMergedHeadCommit(
				result,
				commit,
				source,
				"its commit metadata could not be inspected",
			)
		}

		if taskID != "" {
			if specSlug != run.SpecSlug {
				return unrepresentedMergedHeadCommit(
					result,
					commit,
					source,
					fmt.Sprintf("it belongs to Spec %s", specSlug),
				)
			}
			if !taskCompletedAtMergedHead(ctx, runner, gitRoot, source.head, run.SpecSlug, taskID) {
				return unrepresentedMergedHeadCommit(
					result,
					commit,
					source,
					fmt.Sprintf("Task %s is not completed", taskID),
				)
			}
			taskCommits++
			continue
		}

		paths, pathsOK := mergedHeadCommitPaths(ctx, runner, gitRoot, commit)
		if !pathsOK {
			return unrepresentedMergedHeadCommit(
				result,
				commit,
				source,
				"its changed files could not be inspected",
			)
		}
		if matchesQAReportCommitMessage(message, run.SpecSlug) &&
			len(paths) != 0 && pathsUnderGitDirectories(paths, qaReportDirectories(run.SpecSlug)) {
			report, proven := supersedingQAReportAfterQAOnly(
				ctx,
				runner,
				gitRoot,
				source.head,
				commit,
				run.SpecSlug,
			)
			if !proven {
				return unrepresentedMergedHeadCommit(
					result,
					commit,
					source,
					"its QA Report is not superseded",
				)
			}
			qaCommits++
			result.SupersedingReport = report
			continue
		}
		for _, changedPath := range paths {
			if archived && pathUnderAnyGitDirectory(changedPath, mergedSpecDirectories(run.SpecSlug)) {
				specPaths[changedPath] = struct{}{}
				continue
			}
			otherPaths[changedPath] = struct{}{}
		}
	}

	paths := make([]string, 0, len(otherPaths))
	for changedPath := range otherPaths {
		paths = append(paths, changedPath)
	}
	sort.Strings(paths)
	runOnly, differingShared, retainedRunDeletions, proven := compareRunContentAtHead(
		ctx,
		runner,
		gitRoot,
		result.RunHead,
		source.head,
		paths,
	)
	if !proven {
		result.State = ReconciliationUnintegrated
		result.Reason = boundedReconciliationReason(fmt.Sprintf(
			"Run Branch content comparison could not prove integration against %s",
			source.label,
		))
		return result
	}
	if runOnly != 0 || differingShared != 0 || retainedRunDeletions != 0 {
		result.State = ReconciliationUnintegrated
		result.Reason = mergedHeadContentDifferenceReason(
			source.label,
			runOnly,
			differingShared,
			retainedRunDeletions,
		)
		return result
	}

	if taskCommits != 0 || qaCommits != 0 || len(specPaths) != 0 {
		result.State = ReconciliationSuperseded
		reason := fmt.Sprintf(
			"Run work is superseded at %s: %d Task commit%s completed, %d QA Report commit%s superseded",
			source.label,
			taskCommits,
			pluralSuffix(taskCommits),
			qaCommits,
			pluralSuffix(qaCommits),
		)
		if len(specPaths) != 0 {
			reason = boundedReconciliationReasonWithSuffix(reason, fmt.Sprintf(", %d Spec-directory path(s) archived", len(specPaths)))
		}
		result.Reason = boundedReconciliationReason(reason)
	} else {
		result.State = ReconciliationSafe
		result.Reason = boundedReconciliationReason(
			"Run Branch content is fully represented on " + source.label,
		)
	}
	result.evidence = newMergedHeadReconciliationEvidence(
		run,
		gitRoot,
		result,
		worktreePresent,
		runBranchPresent,
		merged,
		source,
	)
	return result
}

func commitTrailer(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	commit string,
	key string,
) (string, error) {
	output, err := runner.Run(
		ctx,
		gitRoot,
		"show",
		"-s",
		"--format=%(trailers:key="+key+",valueonly,unfold)",
		commit,
	)
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(output)
	if strings.ContainsAny(value, "\x00\r\n") {
		return "", fmt.Errorf("commit %s has ambiguous %s trailers", commit, key)
	}
	return value, nil
}

func taskCompletedAtMergedHead(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	head string,
	slug string,
	taskID string,
) bool {
	cleanSlug, err := cleanPathSegment(slug)
	if err != nil || cleanSlug != slug {
		return false
	}
	roots := []string{
		filepath.ToSlash(filepath.Join("docs", "specs", slug)),
		filepath.ToSlash(filepath.Join(filepath.FromSlash(spec.ArchiveDir(spec.ArchiveKindSpec)), slug)),
	}
	completed := 0
	for _, root := range roots {
		manifestPath := path.Join(root, "_tasks.md")
		manifest, readErr := gitBlobAtHead(ctx, runner, gitRoot, head, manifestPath)
		if readErr != nil {
			continue
		}
		taskFile, found := mergedHeadTaskFile(manifest, taskID)
		if !found {
			continue
		}
		taskPath := path.Join(root, taskFile)
		taskContent, readErr := gitBlobAtHead(ctx, runner, gitRoot, head, taskPath)
		if readErr != nil {
			continue
		}
		status, statusErr := spec.CarryForwardStatus(taskPath, taskContent)
		if statusErr == nil && status == spec.StatusCompleted {
			completed++
		}
	}
	return completed == 1
}

func gitBlobAtHead(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	head string,
	file string,
) ([]byte, error) {
	output, err := runner.Run(ctx, gitRoot, "show", head+":"+file)
	if err != nil {
		return nil, err
	}
	return []byte(output), nil
}

func mergedHeadTaskFile(content []byte, taskID string) (string, bool) {
	frontmatter, ok := markdownFrontmatter(content)
	if !ok {
		return "", false
	}
	var manifest mergedHeadManifest
	if err := yaml.Unmarshal(frontmatter, &manifest); err != nil {
		return "", false
	}
	file := ""
	for _, node := range manifest.Graph.Nodes {
		if node.ID != taskID {
			continue
		}
		if file != "" {
			return "", false
		}
		file = strings.TrimSpace(node.File)
	}
	if file == "" || path.IsAbs(file) || path.Clean(file) != file || file == "." || strings.HasPrefix(file, "../") {
		return "", false
	}
	return file, true
}

func markdownFrontmatter(content []byte) ([]byte, bool) {
	text := string(content)
	if !strings.HasPrefix(text, "---\n") {
		return nil, false
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return nil, false
	}
	return []byte(text[4 : 4+end]), true
}

func mergedHeadCommitPaths(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	commit string,
) ([]string, bool) {
	output, err := runner.Run(
		ctx,
		gitRoot,
		"diff-tree",
		"-m",
		"--root",
		"--no-commit-id",
		"--name-only",
		"-r",
		"-z",
		"--no-renames",
		commit,
		"--",
	)
	if err != nil {
		return nil, false
	}
	paths := nonEmptyNULTerms(output)
	slices.Sort(paths)
	return slices.Compact(paths), true
}

func pathsUnderGitDirectories(paths []string, directories []string) bool {
	for _, changedPath := range paths {
		if !pathUnderAnyGitDirectory(changedPath, directories) {
			return false
		}
	}
	return true
}

func unrepresentedMergedHeadCommit(
	result RunWorktreeReconciliation,
	commit string,
	source mergedHeadSource,
	cause string,
) RunWorktreeReconciliation {
	shortCommit := commit
	if len(shortCommit) > 12 {
		shortCommit = shortCommit[:12]
	}
	result.State = ReconciliationUnintegrated
	result.Reason = boundedReconciliationReason(fmt.Sprintf(
		"Run commit %s is not represented at %s: %s",
		shortCommit,
		source.label,
		cause,
	))
	return result
}

func compareRunContentAtHead(
	ctx context.Context,
	runner gitRunner,
	gitRoot string,
	runHead string,
	mergedHead string,
	paths []string,
) (runOnly int, differingShared int, retainedRunDeletions int, proven bool) {
	if len(paths) == 0 {
		return 0, 0, 0, true
	}
	pathspecs := make([]string, 0, len(paths))
	for _, changedPath := range paths {
		pathspecs = append(pathspecs, ":(literal)"+changedPath)
	}
	diffArgs := func(filter string, left string, right string) []string {
		args := []string{"diff", "--name-only", "-z", "--no-renames", "--diff-filter=" + filter, left, right, "--"}
		return append(args, pathspecs...)
	}
	runOnlyOutput, err := runner.Run(ctx, gitRoot, diffArgs("D", runHead, mergedHead)...)
	if err != nil {
		return 0, 0, 0, false
	}
	differingSharedOutput, err := runner.Run(ctx, gitRoot, diffArgs("MT", runHead, mergedHead)...)
	if err != nil {
		return 0, 0, 0, false
	}
	mergeBaseOutput, err := runner.Run(ctx, gitRoot, "merge-base", runHead, mergedHead)
	if err != nil {
		return 0, 0, 0, false
	}
	mergeBase := strings.TrimSpace(mergeBaseOutput)
	if mergeBase == "" {
		return 0, 0, 0, false
	}
	runDeletedOutput, err := runner.Run(ctx, gitRoot, diffArgs("D", mergeBase, runHead)...)
	if err != nil {
		return 0, 0, 0, false
	}
	runDeleted := nonEmptyNULTerms(runDeletedOutput)
	if len(runDeleted) != 0 {
		args := []string{"ls-tree", "-r", "--name-only", "-z", mergedHead, "--"}
		args = append(args, pathspecs...)
		mergedTreeOutput, treeErr := runner.Run(ctx, gitRoot, args...)
		if treeErr != nil {
			return 0, 0, 0, false
		}
		mergedPaths := make(map[string]struct{})
		for _, mergedPath := range nonEmptyNULTerms(mergedTreeOutput) {
			mergedPaths[mergedPath] = struct{}{}
		}
		for _, deletedPath := range runDeleted {
			if _, retained := mergedPaths[deletedPath]; retained {
				retainedRunDeletions++
			}
		}
	}
	return len(nonEmptyNULTerms(runOnlyOutput)),
		len(nonEmptyNULTerms(differingSharedOutput)),
		retainedRunDeletions,
		true
}

func mergedHeadContentDifferenceReason(
	source string,
	runOnly int,
	differingShared int,
	retainedRunDeletions int,
) string {
	var evidence []string
	if runOnly != 0 {
		evidence = append(evidence, fmt.Sprintf("%d Run-only file%s", runOnly, pluralSuffix(runOnly)))
	}
	if differingShared != 0 {
		evidence = append(evidence, fmt.Sprintf(
			"%d differing shared file%s",
			differingShared,
			pluralSuffix(differingShared),
		))
	}
	if retainedRunDeletions != 0 {
		evidence = append(evidence, fmt.Sprintf(
			"%d Run-deleted file%s retained by default",
			retainedRunDeletions,
			pluralSuffix(retainedRunDeletions),
		))
	}
	return boundedReconciliationReason(fmt.Sprintf(
		"Run Branch content is not fully represented: %s against %s",
		strings.Join(evidence, ", "),
		source,
	))
}

func newMergedHeadReconciliationEvidence(
	run store.Run,
	gitRoot string,
	result RunWorktreeReconciliation,
	worktreePresent bool,
	runBranchPresent bool,
	merged []MergedHead,
	source mergedHeadSource,
) *terminalRunReconciliationEvidence {
	evidence := newTerminalRunReconciliationEvidence(
		run,
		gitRoot,
		result,
		worktreePresent,
		runBranchPresent,
	)
	evidence.merged = slices.Clone(merged)
	evidence.mergedSnapshot = slices.Clone(merged)
	evidence.proofHead = source.head
	evidence.proofRef = source.ref
	return evidence
}

func mergedHeadRecordsEqual(left []MergedHead, right []MergedHead) bool {
	return slices.Equal(left, right)
}

// The archived PRD is positive proof that the Spec directory was delivered.
func specArchivedAtMergedHead(ctx context.Context, runner gitRunner, root, head, slug string) bool {
	clean, err := cleanPathSegment(slug)
	if err != nil || clean != slug {
		return false
	}
	_, err = runner.Run(ctx, root, "cat-file", "-e", head+":"+path.Join(spec.ArchiveDir(spec.ArchiveKindSpec), slug, "_prd.md"))
	return err == nil
}

func dirtyPathsInArchivedSpec(ctx context.Context, runner gitRunner, root, head, slug string, dirty []string) bool {
	directories := mergedSpecDirectories(slug)
	if pathsUnderGitDirectories(dirty, directories) {
		return true
	}
	archive := path.Join(spec.ArchiveDir(spec.ArchiveKindSpec), slug)
	content, err := gitBlobAtHead(ctx, runner, root, head, path.Join(archive, "_tasks.md"))
	if err != nil {
		return false
	}
	frontmatter, ok := markdownFrontmatter(content)
	if !ok {
		return false
	}
	var manifest mergedHeadManifest
	if err := yaml.Unmarshal(frontmatter, &manifest); err != nil {
		return false
	}
	scoped := make(map[string]bool)
	for _, node := range manifest.Graph.Nodes {
		file, ok := mergedHeadTaskFile(content, node.ID)
		if !ok {
			return false
		}
		taskContent, err := gitBlobAtHead(ctx, runner, root, head, path.Join(archive, file))
		if err != nil {
			return false
		}
		// Validate the Task with the shared parser before reading its declarations.
		if _, err := spec.CarryForwardStatus(file, taskContent); err != nil {
			return false
		}
		inContext := false
		for _, line := range strings.Split(string(taskContent), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "# ") {
				inContext = false
			}
			if strings.HasPrefix(line, "## ") {
				inContext = line == "## Context"
				continue
			}
			if !inContext || !strings.HasPrefix(line, "- ") {
				continue
			}
			kind, value, ok := strings.Cut(strings.TrimPrefix(line, "- "), ":")
			if !ok || (strings.TrimSpace(kind) != "interface" && strings.TrimSpace(kind) != "creates" && strings.TrimSpace(kind) != "deletes") {
				continue
			}
			value = strings.TrimSpace(value)
			if start := strings.IndexByte(value, '`'); start >= 0 {
				if end := strings.IndexByte(value[start+1:], '`'); end >= 0 {
					value = value[start+1 : start+1+end]
				}
			}
			if validScopedGitPath(value) {
				scoped[value] = true
			}
		}
		for _, recorded := range spec.RecordedTaskPaths(taskContent) {
			if validScopedGitPath(recorded) {
				scoped[recorded] = true
			}
		}
	}
	for _, file := range dirty {
		if !pathUnderAnyGitDirectory(file, directories) && !scoped[file] {
			return false
		}
	}
	return true
}

func validScopedGitPath(file string) bool {
	return file != "" && file != "." && file != ".." && !path.IsAbs(file) && path.Clean(file) == file && !strings.HasPrefix(file, "../") && !strings.ContainsAny(file, "\\\x00")
}

// Porcelain -z keeps filenames literal; renames/copies carry both paths.
func terminalRunDirtyPaths(status string) ([]string, bool) {
	if status == "" {
		return nil, true
	}
	if !strings.HasSuffix(status, "\x00") {
		return nil, false
	}
	entries := strings.Split(strings.TrimSuffix(status, "\x00"), "\x00")
	paths := make([]string, 0, len(entries))
	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if len(entry) < 4 || entry[2] != ' ' {
			return nil, false
		}
		paths = append(paths, entry[3:])
		if strings.ContainsAny(entry[:2], "RC") {
			i++
			if i >= len(entries) || entries[i] == "" {
				return nil, false
			}
			paths = append(paths, entries[i])
		}
	}
	slices.Sort(paths)
	return slices.Compact(paths), true
}

func mergedSpecDirectories(slug string) []string {
	directories := qaReportDirectories(slug)
	for i := range directories {
		directories[i] = path.Dir(directories[i])
	}
	return directories
}
