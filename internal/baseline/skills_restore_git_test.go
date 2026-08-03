// Suite: immutable skill source Git reads
// Invariant: one batch process returns exact blob bytes in request order and every read failure keeps the restore error contract.
// Boundary IN: skills-restore tree parsing, the cat-file batch protocol, and Git process lifetime.
// Boundary OUT: restoration planning and mutation, which skills_restore_test.go owns.
package baseline

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

const batchObjectDeathFixture = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func TestBatchObjectReaderReturnsMultipleObjectsInRequestOrder(t *testing.T) {
	t.Parallel()

	repo, objects := newBatchObjectFixture(t, map[string]string{
		"first.txt":  "first\n",
		"second.txt": "second\n",
	})
	reader, err := newBatchObjectReader(t.Context(), "-C", repo)
	if err != nil {
		t.Fatalf("start batch object reader: %v", err)
	}

	var got [][]byte
	for _, object := range []string{objects["second.txt"], objects["first.txt"]} {
		content, err := reader.Read(object)
		if err != nil {
			t.Fatalf("read object %s: %v", object, err)
		}
		got = append(got, content)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close batch object reader: %v", err)
	}
	if want := [][]byte{[]byte("second\n"), []byte("first\n")}; !reflect.DeepEqual(got, want) {
		t.Fatalf("batch contents = %q, want %q", got, want)
	}
}

func TestBatchObjectReaderReportsMissingObject(t *testing.T) {
	t.Parallel()

	repo, _ := newBatchObjectFixture(t, map[string]string{"present.txt": "present\n"})
	reader, err := newBatchObjectReader(t.Context(), "-C", repo)
	if err != nil {
		t.Fatalf("start batch object reader: %v", err)
	}
	missing := strings.Repeat("0", 40)
	if _, err := reader.Read(missing); err == nil || !strings.Contains(err.Error(), "missing") {
		t.Fatalf("missing object error = %v, want missing diagnostic", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close batch object reader: %v", err)
	}
}

func TestBatchObjectReaderReturnsZeroByteBlob(t *testing.T) {
	t.Parallel()

	repo, objects := newBatchObjectFixture(t, map[string]string{"empty.txt": ""})
	reader, err := newBatchObjectReader(t.Context(), "-C", repo)
	if err != nil {
		t.Fatalf("start batch object reader: %v", err)
	}
	content, err := reader.Read(objects["empty.txt"])
	if err != nil {
		t.Fatalf("read zero-byte object: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close batch object reader: %v", err)
	}
	if len(content) != 0 {
		t.Fatalf("zero-byte object length = %d, want 0", len(content))
	}
}

func TestBatchObjectReaderPreservesFramingDelimitersInContent(t *testing.T) {
	t.Parallel()

	want := []byte("header words\n" + batchObjectDeathFixture + " blob 7\nbody\x00tail\n")
	repo, objects := newBatchObjectFixture(t, map[string]string{"delimiters.bin": string(want)})
	reader, err := newBatchObjectReader(t.Context(), "-C", repo)
	if err != nil {
		t.Fatalf("start batch object reader: %v", err)
	}
	content, err := reader.Read(objects["delimiters.bin"])
	if err != nil {
		t.Fatalf("read delimiter-bearing object: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close batch object reader: %v", err)
	}
	if !bytes.Equal(content, want) {
		t.Fatalf("delimiter-bearing content = %q, want %q", content, want)
	}
}

func TestBatchObjectReaderReportsProcessDeathMidStream(t *testing.T) {
	t.Parallel()

	command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestBatchObjectReaderProcessHelper$")
	command.Env = append(os.Environ(), "ROUNDFIX_BATCH_OBJECT_HELPER=die-mid-stream")
	reader, err := startBatchObjectReader(command)
	if err != nil {
		t.Fatalf("start batch object helper: %v", err)
	}
	if _, err := reader.Read(batchObjectDeathFixture); err == nil {
		t.Fatal("mid-stream process death returned no read error")
	}
	if err := reader.Close(); err == nil {
		t.Fatal("mid-stream process death returned no process error")
	}
}

func TestBatchObjectReaderRejectsOversizedObjectsAndPoisonsProtocolFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reply   string
		wantErr string
	}{
		{name: "header read failure", wantErr: "read batch object header"},
		{name: "malformed header", reply: "malformed\n", wantErr: "malformed batch object header"},
		{
			name:    "oversized object",
			reply:   fmt.Sprintf("%s blob %d\n", batchObjectDeathFixture, restoreMaxBytes+1),
			wantErr: "exceeds the read limit",
		},
		{
			name:    "short content",
			reply:   fmt.Sprintf("%s blob 6\nabc", batchObjectDeathFixture),
			wantErr: "read batch object content",
		},
		{
			name:    "invalid terminator",
			reply:   fmt.Sprintf("%s blob 3\nabcX", batchObjectDeathFixture),
			wantErr: "invalid batch object terminator",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			stdin := &batchObjectWriteCloser{}
			reader := &batchObjectReader{
				stdin:    stdin,
				stdout:   bufio.NewReader(strings.NewReader(test.reply)),
				maxBytes: restoreMaxBytes,
			}
			if _, err := reader.Read(batchObjectDeathFixture); err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("first read error = %v, want diagnostic containing %q", err, test.wantErr)
			}
			if _, err := reader.Read(batchObjectDeathFixture); err == nil || !strings.Contains(err.Error(), "desynchronized") {
				t.Fatalf("second read error = %v, want desynchronized-reader diagnostic", err)
			}
			if got, want := stdin.String(), batchObjectDeathFixture+"\n"; got != want {
				t.Fatalf("batch object requests after protocol failure = %q, want %q", got, want)
			}
		})
	}
}

type batchObjectWriteCloser struct {
	bytes.Buffer
}

func (*batchObjectWriteCloser) Close() error {
	return nil
}

func TestBatchObjectReaderProcessHelper(t *testing.T) {
	if os.Getenv("ROUNDFIX_BATCH_OBJECT_HELPER") != "die-mid-stream" {
		return
	}
	if _, err := fmt.Fprintf(os.Stdout, "%s blob 6\nabc", batchObjectDeathFixture); err != nil {
		os.Exit(18)
	}
	os.Exit(17)
}

func TestSkillsRestoreReadsManyFilesThroughOneBatchProcess(t *testing.T) {
	t.Parallel()

	contents := map[string][]byte{
		"first-object":  []byte("first\n"),
		"second-object": []byte("second\n"),
		"third-object":  []byte("third\n"),
	}
	runner := &recordingRestoreObjectGitRunner{
		tree: restoreTreeFixture("skills/example", []restoreTreeFixtureEntry{
			{path: "first.txt", object: "first-object"},
			{path: "second.txt", object: "second-object"},
			{path: "third.txt", object: "third-object"},
		}),
		reader: &recordingBatchObjectReader{contents: contents},
	}
	files := []restoreFile{
		{Path: "first.txt", Content: contents["first-object"]},
		{Path: "second.txt", Content: contents["second-object"]},
		{Path: "third.txt", Content: contents["third-object"]},
	}
	contract := restoreSkillContract{
		Name: "example",
		Source: RestoreSource{
			Repository: "example/skills",
			Ref:        batchObjectDeathFixture,
			Path:       "skills/example",
		},
		TreeDigest: portableRestoreDigest(files),
	}

	got, err := readRestoreGitTree(t.Context(), "objects.git", contract, runner)
	if err != nil {
		t.Fatalf("read restore tree: %v", err)
	}
	if !reflect.DeepEqual(got, files) {
		t.Fatalf("restore files = %+v, want %+v", got, files)
	}
	if runner.batchStarts != 1 {
		t.Fatalf("batch process starts = %d, want 1", runner.batchStarts)
	}
	if runner.runCalls != 1 {
		t.Fatalf("non-batch Git calls = %d, want one ls-tree call", runner.runCalls)
	}
	if runner.reader.closeCalls != 1 {
		t.Fatalf("batch process closes = %d, want 1", runner.reader.closeCalls)
	}
	if want := []string{"first-object", "second-object", "third-object"}; !reflect.DeepEqual(runner.reader.reads, want) {
		t.Fatalf("batch object requests = %v, want %v", runner.reader.reads, want)
	}
}

func TestSkillsRestoreBatchFailuresKeepReadErrorSurface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
	}{
		{name: "missing object", err: errors.New("batch object is missing")},
		{name: "process death mid-stream", err: errors.New("unexpected EOF")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &recordingRestoreObjectGitRunner{
				tree: restoreTreeFixture("skills/example", []restoreTreeFixtureEntry{{
					path: "SKILL.md", object: "object",
				}}),
				reader: &recordingBatchObjectReader{readErr: test.err},
			}
			contract := restoreSkillContract{
				Name: "example",
				Source: RestoreSource{
					Repository: "example/skills",
					Ref:        batchObjectDeathFixture,
					Path:       "skills/example",
				},
			}

			_, err := readRestoreGitTree(t.Context(), "objects.git", contract, runner)
			var restoreErr *SkillsRestoreError
			if !errors.As(err, &restoreErr) {
				t.Fatalf("restore error = %v, want SkillsRestoreError", err)
			}
			if restoreErr.Category != SkillsRestoreExecution || restoreErr.Finding.Code != "source.read-failed" {
				t.Fatalf("restore error = %+v, want execution/source.read-failed", restoreErr)
			}
		})
	}
}

type restoreTreeFixtureEntry struct {
	path   string
	object string
}

func restoreTreeFixture(root string, entries []restoreTreeFixtureEntry) []byte {
	var tree bytes.Buffer
	for _, entry := range entries {
		fmt.Fprintf(&tree, "100644 blob %s\t%s/%s%c", entry.object, root, entry.path, byte(0))
	}
	return tree.Bytes()
}

type recordingRestoreObjectGitRunner struct {
	tree        []byte
	runCalls    int
	batchStarts int
	reader      *recordingBatchObjectReader
}

func (runner *recordingRestoreObjectGitRunner) Run(
	_ context.Context,
	_ ...string,
) ([]byte, error) {
	runner.runCalls++
	return runner.tree, nil
}

func (runner *recordingRestoreObjectGitRunner) OpenBatch(
	_ context.Context,
	_ ...string,
) (batchObjectContentReader, error) {
	runner.batchStarts++
	return runner.reader, nil
}

type recordingBatchObjectReader struct {
	contents   map[string][]byte
	readErr    error
	reads      []string
	closeCalls int
}

func (reader *recordingBatchObjectReader) Read(object string) ([]byte, error) {
	reader.reads = append(reader.reads, object)
	if reader.readErr != nil {
		return nil, reader.readErr
	}
	return append([]byte(nil), reader.contents[object]...), nil
}

func (reader *recordingBatchObjectReader) Close() error {
	reader.closeCalls++
	return nil
}

func newBatchObjectFixture(t *testing.T, files map[string]string) (string, map[string]string) {
	t.Helper()
	repo := newInspectionRepository(t)
	for relative, content := range files {
		writeInspectionFile(t, repo, relative, content)
	}
	commitInspectionRepository(t, repo, "batch object fixture")
	objects := make(map[string]string, len(files))
	for relative := range files {
		object, err := (ExecGitRunner{}).RunGit(t.Context(), repo, "rev-parse", "HEAD:"+relative)
		if err != nil {
			t.Fatalf("resolve object for %s: %v", relative, err)
		}
		objects[relative] = object
	}
	return repo, objects
}
