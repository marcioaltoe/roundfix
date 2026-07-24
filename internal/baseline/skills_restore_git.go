package baseline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func acquireRestoreGroup(
	ctx context.Context,
	provenance restoreProvenance,
	contracts []restoreSkillContract,
	sourceDir string,
) (map[string][]restoreFile, error) {
	if provenance.Provider != "github" {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.provider-unsupported",
			fmt.Sprintf("Skill source provider %q is not supported.", provenance.Provider),
			"Use an embedded snapshot whose external source provider is github.",
			errors.New("source provider is unsupported"),
		)
	}
	temporary, err := os.MkdirTemp("", "roundfix-baseline-skills-restore-")
	if err != nil {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.acquire-failed",
			fmt.Sprintf("Could not create the immutable source staging area: %v.", err),
			"Fix temporary-directory permissions and retry.",
			err,
		)
	}
	defer os.RemoveAll(temporary)
	objectStore := filepath.Join(temporary, "objects.git")
	if _, err := runRestoreGit(ctx, "init", "--bare", objectStore); err != nil {
		return nil, classifyRestoreGitError(err, provenance)
	}
	origin := "https://github.com/" + provenance.Repository + ".git"
	if sourceDir != "" {
		origin = sourceDir
	}
	if _, err := runRestoreGit(
		ctx,
		"--git-dir", objectStore,
		"fetch", "--no-tags", "--depth=1",
		origin, provenance.Ref,
	); err != nil {
		return nil, classifyRestoreGitError(err, provenance)
	}
	resolved, err := runRestoreGit(
		ctx,
		"--git-dir", objectStore,
		"rev-parse", "--verify", "FETCH_HEAD^{commit}",
	)
	if err != nil {
		return nil, classifyRestoreGitError(err, provenance)
	}
	if strings.TrimSpace(string(resolved)) != provenance.Ref {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.commit-mismatch",
			fmt.Sprintf(
				"Acquired %s resolved to %s, not the declared commit %s.",
				provenance.Repository,
				strings.TrimSpace(string(resolved)),
				provenance.Ref,
			),
			"Use an offline source that contains the exact declared immutable commit.",
			errors.New("acquired commit identity mismatch"),
		)
	}

	result := make(map[string][]restoreFile, len(contracts))
	for _, contract := range contracts {
		files, err := readRestoreGitTree(ctx, objectStore, contract)
		if err != nil {
			return nil, err
		}
		result[contract.Name] = files
	}
	return result, nil
}

func readRestoreGitTree(
	ctx context.Context,
	objectStore string,
	contract restoreSkillContract,
) ([]restoreFile, error) {
	output, err := runRestoreGit(
		ctx,
		"--git-dir", objectStore,
		"ls-tree", "-r", "-z", "--full-tree",
		contract.Source.Ref,
		"--", contract.Source.Path,
	)
	if err != nil {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.commit-unavailable",
			fmt.Sprintf(
				"Could not inspect %s@%s: %v.",
				contract.Source.Repository,
				contract.Source.Ref,
				err,
			),
			"Verify the exact commit and source path in the embedded Setup Snapshot.",
			err,
		)
	}
	prefix := []byte(contract.Source.Path + "/")
	files := []restoreFile{}
	totalBytes := 0
	for _, raw := range bytes.Split(output, []byte{0}) {
		if len(raw) == 0 {
			continue
		}
		header, rawPath, ok := bytes.Cut(raw, []byte{'\t'})
		parts := bytes.Fields(header)
		if !ok || len(parts) != 3 {
			return nil, unsafeRestoreSource("the acquired Git tree contains a malformed entry")
		}
		if !bytes.HasPrefix(rawPath, prefix) {
			return nil, unsafeRestoreSource(fmt.Sprintf(
				"Git tree entry escapes %s: %s",
				contract.Source.Path,
				string(rawPath),
			))
		}
		relative := string(rawPath[len(prefix):])
		if !safeRestoreRelative(relative) || path.Clean(relative) != relative {
			return nil, unsafeRestoreSource(fmt.Sprintf(
				"the acquired skill subtree contains an unsafe path: %q",
				relative,
			))
		}
		if string(parts[1]) != "blob" ||
			(string(parts[0]) != "100644" && string(parts[0]) != "100755") {
			return nil, unsafeRestoreSource(fmt.Sprintf(
				"the acquired skill subtree contains a link, device, or unsupported entry: %s",
				relative,
			))
		}
		content, err := runRestoreGit(
			ctx,
			"--git-dir", objectStore,
			"cat-file", "blob", string(parts[2]),
		)
		if err != nil {
			return nil, restoreError(
				SkillsRestoreExecution,
				"source.read-failed",
				fmt.Sprintf("Could not read acquired file %s: %v.", relative, err),
				"Verify the offline object store or declared GitHub commit and retry.",
				err,
			)
		}
		files = append(files, restoreFile{Path: relative, Content: content})
		totalBytes += len(content)
		if len(files) > restoreMaxFiles || totalBytes > restoreMaxBytes {
			return nil, restoreError(
				SkillsRestoreExecution,
				"source.limit-exceeded",
				"Acquired skill content exceeds the restoration file-count or byte limit.",
				"Reduce the declared immutable source subtree before retrying.",
				errors.New("source subtree limit exceeded"),
			)
		}
	}
	if len(files) == 0 {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.empty-tree",
			fmt.Sprintf("The declared source path %s contains no restorable files.", contract.Source.Path),
			"Fix the Setup Snapshot source path and immutable revision.",
			errors.New("source subtree is empty"),
		)
	}
	sortRestoreFiles(files)
	observed := portableRestoreDigest(files)
	if observed != contract.TreeDigest {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.digest-mismatch",
			fmt.Sprintf(
				"Acquired skill %s digest %s does not match %s.",
				contract.Name,
				observed,
				contract.TreeDigest,
			),
			"Regenerate the embedded Setup Snapshot from the exact committed skill subtree.",
			errors.New("source tree digest mismatch"),
		)
	}
	return files, nil
}

func runRestoreGit(ctx context.Context, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Env = append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GCM_INTERACTIVE=Never",
		"GIT_ASKPASS=",
		"SSH_ASKPASS=",
	)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if contextErr := ctx.Err(); contextErr != nil {
			return nil, contextErr
		}
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	return stdout.Bytes(), nil
}

func classifyRestoreGitError(err error, provenance restoreProvenance) error {
	if errors.Is(err, exec.ErrNotFound) {
		return restoreError(
			SkillsRestoreExecution,
			"source.git-missing",
			"Git is required for immutable external skill acquisition.",
			"Install Git and rerun the same restoration preview.",
			err,
		)
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return restoreError(
		SkillsRestoreExecution,
		"source.commit-unavailable",
		fmt.Sprintf(
			"Could not acquire %s@%s: %v.",
			provenance.Repository,
			provenance.Ref,
			err,
		),
		"Verify the exact commit and source path in the embedded Setup Snapshot.",
		err,
	)
}

func unsafeRestoreSource(message string) error {
	return restoreError(
		SkillsRestoreExecution,
		"source.unsafe-tree",
		message+".",
		"Remove unsafe entries from the immutable source and publish a new snapshot.",
		errors.New(message),
	)
}
