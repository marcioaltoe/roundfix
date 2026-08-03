package baseline

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type batchObjectContentReader interface {
	Read(string) ([]byte, error)
	Close() error
}

type batchObjectGitRunner interface {
	Run(context.Context, ...string) ([]byte, error)
	OpenBatch(context.Context, ...string) (batchObjectContentReader, error)
}

type execBatchObjectGitRunner struct{}

func (execBatchObjectGitRunner) Run(ctx context.Context, args ...string) ([]byte, error) {
	return runRestoreGit(ctx, args...)
}

func (execBatchObjectGitRunner) OpenBatch(
	ctx context.Context,
	args ...string,
) (batchObjectContentReader, error) {
	return newBatchObjectReader(ctx, args...)
}

type batchObjectReader struct {
	stdin    io.WriteCloser
	stdout   *bufio.Reader
	cmd      *exec.Cmd
	stderr   *bytes.Buffer
	maxBytes int64
	broken   bool
	closed   bool
	closeErr error
}

func newBatchObjectReader(ctx context.Context, gitArgs ...string) (*batchObjectReader, error) {
	args := append(append([]string(nil), gitArgs...), "cat-file", "--batch")
	command := exec.CommandContext(ctx, "git", args...)
	command.Env = restoreGitEnvironment()
	return startBatchObjectReader(command)
}

func startBatchObjectReader(command *exec.Cmd) (*batchObjectReader, error) {
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open batch object stdin: %w", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		closeErr := stdin.Close()
		return nil, errors.Join(fmt.Errorf("open batch object stdout: %w", err), closeErr)
	}
	stderr := &bytes.Buffer{}
	command.Stderr = stderr
	if err := command.Start(); err != nil {
		stdinErr := stdin.Close()
		stdoutErr := stdout.Close()
		return nil, errors.Join(fmt.Errorf("start batch object process: %w", err), stdinErr, stdoutErr)
	}
	return &batchObjectReader{
		stdin:    stdin,
		stdout:   bufio.NewReader(stdout),
		cmd:      command,
		stderr:   stderr,
		maxBytes: restoreMaxBytes,
	}, nil
}

func (reader *batchObjectReader) Read(object string) ([]byte, error) {
	if reader.closed {
		return nil, errors.New("batch object reader is closed")
	}
	if reader.broken {
		return nil, errors.New("batch object reader is desynchronized")
	}
	if object == "" || strings.ContainsAny(object, "\r\n") {
		return nil, fmt.Errorf("invalid batch object name %q", object)
	}
	if _, err := fmt.Fprintln(reader.stdin, object); err != nil {
		reader.broken = true
		return nil, fmt.Errorf("write batch object request: %w", err)
	}
	header, err := reader.stdout.ReadString('\n')
	if err != nil {
		reader.broken = true
		return nil, fmt.Errorf("read batch object header: %w", err)
	}
	fields := strings.Fields(strings.TrimSuffix(header, "\n"))
	if len(fields) == 2 && fields[1] == "missing" {
		return nil, fmt.Errorf("batch object %s is missing", object)
	}
	if len(fields) != 3 {
		reader.broken = true
		return nil, fmt.Errorf("malformed batch object header %q", strings.TrimSuffix(header, "\n"))
	}
	if fields[1] != "blob" {
		reader.broken = true
		return nil, fmt.Errorf("batch object %s has type %s, want blob", object, fields[1])
	}
	size, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil || size < 0 {
		reader.broken = true
		return nil, fmt.Errorf("invalid batch object size %q", fields[2])
	}
	if reader.maxBytes > 0 && size > reader.maxBytes {
		reader.broken = true
		return nil, fmt.Errorf(
			"batch object %s size %d exceeds the read limit %d",
			object,
			size,
			reader.maxBytes,
		)
	}
	content := make([]byte, int(size))
	if _, err := io.ReadFull(reader.stdout, content); err != nil {
		reader.broken = true
		return nil, fmt.Errorf("read batch object content: %w", err)
	}
	terminator, err := reader.stdout.ReadByte()
	if err != nil {
		reader.broken = true
		return nil, fmt.Errorf("read batch object terminator: %w", err)
	}
	if terminator != '\n' {
		reader.broken = true
		return nil, fmt.Errorf("invalid batch object terminator %#x", terminator)
	}
	return content, nil
}

func (reader *batchObjectReader) Close() error {
	if reader.closed {
		return reader.closeErr
	}
	reader.closed = true
	stdinErr := reader.stdin.Close()
	waitErr := reader.cmd.Wait()
	if waitErr != nil {
		detail := strings.TrimSpace(reader.stderr.String())
		if detail == "" {
			detail = waitErr.Error()
		}
		waitErr = fmt.Errorf(
			"batch object process %s: %s: %w",
			strings.Join(reader.cmd.Args, " "),
			detail,
			waitErr,
		)
	}
	reader.closeErr = errors.Join(stdinErr, waitErr)
	return reader.closeErr
}

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
	gitRunner := execBatchObjectGitRunner{}
	for _, contract := range contracts {
		files, err := readRestoreGitTree(ctx, objectStore, contract, gitRunner)
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
	gitRunner batchObjectGitRunner,
) ([]restoreFile, error) {
	output, err := gitRunner.Run(
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
	type treeEntry struct {
		relative string
		object   string
	}
	entries := []treeEntry{}
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
		entries = append(entries, treeEntry{relative: relative, object: string(parts[2])})
	}
	if len(entries) == 0 {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.empty-tree",
			fmt.Sprintf("The declared source path %s contains no restorable files.", contract.Source.Path),
			"Fix the Setup Snapshot source path and immutable revision.",
			errors.New("source subtree is empty"),
		)
	}
	reader, err := gitRunner.OpenBatch(ctx, "--git-dir", objectStore)
	if err != nil {
		return nil, restoreGitReadError(entries[0].relative, err)
	}
	files := make([]restoreFile, 0, len(entries))
	totalBytes := 0
	for _, entry := range entries {
		content, err := reader.Read(entry.object)
		if err != nil {
			if closeErr := reader.Close(); closeErr != nil {
				err = errors.Join(err, closeErr)
			}
			return nil, restoreGitReadError(entry.relative, err)
		}
		files = append(files, restoreFile{Path: entry.relative, Content: content})
		totalBytes += len(content)
		if len(files) > restoreMaxFiles || totalBytes > restoreMaxBytes {
			if err := reader.Close(); err != nil {
				return nil, restoreGitReadError(entry.relative, err)
			}
			return nil, restoreError(
				SkillsRestoreExecution,
				"source.limit-exceeded",
				"Acquired skill content exceeds the restoration file-count or byte limit.",
				"Reduce the declared immutable source subtree before retrying.",
				errors.New("source subtree limit exceeded"),
			)
		}
	}
	if err := reader.Close(); err != nil {
		return nil, restoreGitReadError(entries[len(entries)-1].relative, err)
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

func restoreGitReadError(relative string, err error) error {
	return restoreError(
		SkillsRestoreExecution,
		"source.read-failed",
		fmt.Sprintf("Could not read acquired file %s: %v.", relative, err),
		"Verify the offline object store or declared GitHub commit and retry.",
		err,
	)
}

func runRestoreGit(ctx context.Context, args ...string) ([]byte, error) {
	command := exec.CommandContext(ctx, "git", args...)
	command.Env = restoreGitEnvironment()
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

func restoreGitEnvironment() []string {
	return append(os.Environ(),
		"GIT_TERMINAL_PROMPT=0",
		"GCM_INTERACTIVE=Never",
		"GIT_ASKPASS=",
		"SSH_ASKPASS=",
	)
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
