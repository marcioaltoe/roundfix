// Suite: canonical Baseline asset synchronization
// Invariant: setup snapshots change only from a clean immutable Git source whose normalized tree is catalog-compatible, and failed refresh restores every preimage.
// Boundary IN: source provenance, setup normalization, drift detection, catalog validation, and recoverable refresh.
// Boundary OUT: public flag parsing and rendering, owned by internal/cli/baseline_assets_sync_test.go.
package baseline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/suiteguard"
)

func TestInspectAssetsSyncCheckoutCombinesResolutionQueries(t *testing.T) {
	t.Parallel()

	sourceDir := t.TempDir()
	root, err := filepath.EvalSymlinks(sourceDir)
	if err != nil {
		t.Fatalf("resolve temporary checkout root: %v", err)
	}
	revision := strings.Repeat("a", 40)
	var calls [][]string
	runner := assetsSyncSourceGitRunner(func(_ context.Context, args ...string) ([]byte, error) {
		calls = append(calls, append([]string(nil), args...))
		switch strings.Join(args, "\x00") {
		case "-C\x00" + sourceDir + "\x00rev-parse\x00--show-toplevel\x00--verify\x00HEAD^{commit}":
			return []byte(root + "\n" + revision + "\n"), nil
		case "-C\x00" + root + "\x00status\x00--porcelain=v1\x00--untracked-files=all":
			return nil, nil
		case "-C\x00" + root + "\x00remote\x00get-url\x00origin":
			return []byte("https://github.com/example/skills.git\n"), nil
		default:
			return nil, fmt.Errorf("unexpected Git command: %v", args)
		}
	})

	checkout, err := inspectAssetsSyncCheckoutWithRunner(t.Context(), sourceDir, runner)
	if err != nil {
		t.Fatalf("inspect assets sync checkout: %v", err)
	}
	if checkout.root != root || checkout.repository != "example/skills" || checkout.revision != revision {
		t.Fatalf("checkout = %+v, want root, repository, and revision from combined resolution", checkout)
	}
	if len(calls) != 3 {
		t.Fatalf("Git calls = %v, want one resolution plus status and remote", calls)
	}
	wantResolution := []string{"-C", sourceDir, "rev-parse", "--show-toplevel", "--verify", "HEAD^{commit}"}
	if !reflect.DeepEqual(calls[0], wantResolution) {
		t.Fatalf("resolution call = %v, want %v", calls[0], wantResolution)
	}
}

func TestAssetsSyncCommittedTreeDigestReadsManyFilesThroughOneBatchProcess(t *testing.T) {
	t.Parallel()

	const sourcePath = "skills/example"
	contents := map[string][]byte{
		"first-object":  []byte("first\n"),
		"second-object": []byte("second\n"),
		"third-object":  []byte("third\n"),
	}
	runner := &recordingAssetsSyncObjectGitRunner{
		tree: restoreTreeFixture(sourcePath, []restoreTreeFixtureEntry{
			{path: "first.txt", object: "first-object"},
			{path: "second.txt", object: "second-object"},
			{path: "third.txt", object: "third-object"},
		}),
		reader: &recordingBatchObjectReader{contents: contents},
	}
	want := portableRestoreDigest([]restoreFile{
		{Path: "first.txt", Content: contents["first-object"]},
		{Path: "second.txt", Content: contents["second-object"]},
		{Path: "third.txt", Content: contents["third-object"]},
	})

	got, err := assetsSyncCommittedTreeDigestWithRunner(
		t.Context(),
		"source.git",
		strings.Repeat("a", 40),
		sourcePath,
		runner,
	)
	if err != nil {
		t.Fatalf("committed tree digest: %v", err)
	}
	if got != want {
		t.Fatalf("committed tree digest = %s, want %s", got, want)
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
	if wantReads := []string{"first-object", "second-object", "third-object"}; !reflect.DeepEqual(runner.reader.reads, wantReads) {
		t.Fatalf("batch object requests = %v, want %v", runner.reader.reads, wantReads)
	}
}

func TestAssetsSyncCommittedTreeDigestKeepsTreeAndReadErrors(t *testing.T) {
	t.Parallel()

	const sourcePath = "skills/example"
	readErr := errors.New("existing Git blob read failure")
	tests := []struct {
		name    string
		tree    []byte
		readErr error
		want    string
	}{
		{
			name: "unsafe path",
			tree: []byte(fmt.Sprintf(
				"100644 blob object\t%s/../escape%c",
				sourcePath,
				byte(0),
			)),
			want: `Git tree entry has unsafe path "../escape"`,
		},
		{
			name: "non-blob entry",
			tree: []byte(fmt.Sprintf(
				"040000 tree object\t%s/nested%c",
				sourcePath,
				byte(0),
			)),
			want: "Git tree entry is not a regular file: nested",
		},
		{
			name:    "read failure",
			tree:    restoreTreeFixture(sourcePath, []restoreTreeFixtureEntry{{path: "SKILL.md", object: "object"}}),
			readErr: readErr,
			want:    readErr.Error(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			runner := &recordingAssetsSyncObjectGitRunner{
				tree:   test.tree,
				reader: &recordingBatchObjectReader{readErr: test.readErr},
			}
			_, err := assetsSyncCommittedTreeDigestWithRunner(
				t.Context(),
				"source.git",
				strings.Repeat("a", 40),
				sourcePath,
				runner,
			)
			if err == nil || err.Error() != test.want {
				t.Fatalf("committed tree error = %v, want %q", err, test.want)
			}
			if test.readErr != nil && !errors.Is(err, test.readErr) {
				t.Fatalf("committed tree error = %v, want wrapped %v", err, test.readErr)
			}
		})
	}
}

type recordingAssetsSyncObjectGitRunner struct {
	tree        []byte
	runCalls    int
	batchStarts int
	reader      *recordingBatchObjectReader
}

func (runner *recordingAssetsSyncObjectGitRunner) Run(
	_ context.Context,
	_ ...string,
) ([]byte, error) {
	runner.runCalls++
	return runner.tree, nil
}

func (runner *recordingAssetsSyncObjectGitRunner) OpenBatch(
	_ context.Context,
	_ ...string,
) (batchObjectContentReader, error) {
	runner.batchStarts++
	return runner.reader, nil
}

func TestAssetsSyncCheckIsReadOnlyAndReportsDrift(t *testing.T) {
	t.Parallel()

	targetRepo, assetRoot := newAssetsSyncTarget(t)
	sourceDir, _ := newAssetsSyncSource(t, assetRoot)
	before := captureAssetsSyncTree(t, assetRoot)

	payload, err := syncAssets(context.Background(), AssetsSyncRequest{
		SourceDir: sourceDir,
		Check:     true,
	}, assetsSyncDependencies{
		repository: targetRepo,
		assetRoot:  assetRoot,
	})

	var syncErr *AssetsSyncError
	if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncExecution {
		t.Fatalf("check error = %T %v, want execution drift", err, err)
	}
	if payload.OK || payload.Summary.Errors != 3 {
		t.Fatalf("check payload = %+v", payload)
	}
	for _, finding := range payload.Findings {
		if finding.Code != "skills.setup-snapshot.drift" {
			t.Fatalf("finding = %+v", finding)
		}
	}
	if after := captureAssetsSyncTree(t, assetRoot); !reflect.DeepEqual(after, before) {
		t.Fatal("check mode changed canonical Baseline assets")
	}
}

func TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent(t *testing.T) {
	t.Parallel()

	targetRepo, assetRoot := newAssetsSyncTarget(t)
	sourceDir, revision := newAssetsSyncSource(t, assetRoot)
	beforeOwnedMinimum := assetsSyncOwnedMinimum(t, filepath.Join(assetRoot, "setups", "go-cli.json"))

	payload, err := syncAssets(context.Background(), AssetsSyncRequest{
		SourceDir: sourceDir,
	}, assetsSyncDependencies{
		repository: targetRepo,
		assetRoot:  assetRoot,
	})
	if err != nil {
		t.Fatalf("refresh: %v payload=%+v", err, payload)
	}
	if !payload.OK || payload.Summary.Info != 3 {
		t.Fatalf("refresh payload = %+v", payload)
	}
	catalog, err := LoadCatalog(os.DirFS(assetRoot))
	if err != nil {
		t.Fatalf("load refreshed catalog: %v", err)
	}
	if len(catalog.SetupIDs()) != 3 {
		t.Fatalf("refreshed setup IDs = %v", catalog.SetupIDs())
	}
	assertAssetsSyncOwnedSkillHasNoContentPin(t, filepath.Join(assetRoot, "setups", "go-cli.json"))
	if afterOwnedMinimum := assetsSyncOwnedMinimum(
		t,
		filepath.Join(assetRoot, "setups", "go-cli.json"),
	); afterOwnedMinimum != beforeOwnedMinimum {
		t.Fatalf(
			"Roundfix-owned minimum = %q, want preserved declaration %q",
			afterOwnedMinimum,
			beforeOwnedMinimum,
		)
	}
	for _, setupID := range catalog.SetupIDs() {
		var snapshot struct {
			Source struct {
				Repository string `json:"repository"`
				Ref        string `json:"ref"`
			} `json:"source"`
			Digest string `json:"digest"`
		}
		readAssetsSyncJSON(t, filepath.Join(assetRoot, "setups", setupID+".json"), &snapshot)
		if snapshot.Source.Repository != "example/skills" || snapshot.Source.Ref != revision {
			t.Fatalf("%s source = %+v", setupID, snapshot.Source)
		}
		if !lowercaseSHA256.MatchString(snapshot.Digest) {
			t.Fatalf("%s digest = %q", setupID, snapshot.Digest)
		}
	}

	after := captureAssetsSyncTree(t, assetRoot)
	repeated, err := syncAssets(context.Background(), AssetsSyncRequest{
		SourceDir: sourceDir,
	}, assetsSyncDependencies{
		repository: targetRepo,
		assetRoot:  assetRoot,
	})
	if err != nil || !repeated.OK || len(repeated.Findings) != 0 {
		t.Fatalf("repeat = %+v err=%v", repeated, err)
	}
	if got := captureAssetsSyncTree(t, assetRoot); !reflect.DeepEqual(got, after) {
		t.Fatal("idempotent refresh changed canonical assets")
	}
}

func TestAssetsSyncProvenanceAndPreMutationRefusals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(t *testing.T, sourceDir string)
		want   string
	}{
		{
			name: "dirty or untracked checkout",
			mutate: func(t *testing.T, sourceDir string) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(filepath.Dir(sourceDir), "untracked.txt"), []byte("dirty\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: "dirty or untracked",
		},
		{
			name: "non portable Git identity",
			mutate: func(t *testing.T, sourceDir string) {
				t.Helper()
				runAssetsSyncGit(
					t,
					filepath.Dir(sourceDir),
					"remote",
					"set-url",
					"origin",
					"file:///tmp/canonical-skills",
				)
			},
			want: "portable GitHub repository identity",
		},
		{
			name: "mutable declared ref",
			mutate: func(t *testing.T, sourceDir string) {
				t.Helper()
				path := filepath.Join(sourceDir, "rust-cli.json")
				var snapshot struct {
					Skills []map[string]any `json:"skills"`
				}
				target := filepath.Join(sourceDir, "rust-cli.txt")
				lines, err := os.ReadFile(target)
				if err != nil {
					t.Fatal(err)
				}
				for _, line := range strings.Split(string(lines), "\n") {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					snapshot.Skills = append(snapshot.Skills, map[string]any{"path": line})
				}
				snapshot.Skills[0]["source"] = map[string]any{
					"type": "github", "repository": "example/skills", "ref": "main",
					"path": snapshot.Skills[0]["path"],
				}
				writeAssetsSyncJSON(t, path, snapshot)
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "add", ".")
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "commit", "--quiet", "-m", "mutable ref")
			},
			want: "full immutable commit",
		},
		{
			name: "non portable skill path",
			mutate: func(t *testing.T, sourceDir string) {
				t.Helper()
				writeAssetsSyncJSON(t, filepath.Join(sourceDir, "go-cli.json"), map[string]any{
					"skills": []any{map[string]any{"path": "/tmp/unsafe"}},
				})
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "add", ".")
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "commit", "--quiet", "-m", "unsafe path")
			},
			want: "not portable",
		},
		{
			name: "incompatible empty setup",
			mutate: func(t *testing.T, sourceDir string) {
				t.Helper()
				writeAssetsSyncJSON(t, filepath.Join(sourceDir, "typescript-bun.json"), map[string]any{
					"skills": []any{},
				})
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "add", ".")
				runAssetsSyncGit(t, filepath.Dir(sourceDir), "commit", "--quiet", "-m", "empty setup")
			},
			want: "non-empty skills list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each case builds its own target repository and source
			// checkout, so the refusal cases overlap instead of queueing.
			t.Parallel()
			targetRepo, assetRoot := newAssetsSyncTarget(t)
			sourceDir, _ := newAssetsSyncSource(t, assetRoot)
			tt.mutate(t, sourceDir)
			before := captureAssetsSyncTree(t, assetRoot)

			payload, err := syncAssets(context.Background(), AssetsSyncRequest{
				SourceDir: sourceDir,
			}, assetsSyncDependencies{
				repository: targetRepo,
				assetRoot:  assetRoot,
			})

			var syncErr *AssetsSyncError
			if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncInvalid {
				t.Fatalf("error = %T %v payload=%+v", err, err, payload)
			}
			found := false
			for _, finding := range payload.Findings {
				found = found || strings.Contains(finding.Message, tt.want)
			}
			if !found {
				t.Fatalf("findings = %+v, want message containing %q", payload.Findings, tt.want)
			}
			if after := captureAssetsSyncTree(t, assetRoot); !reflect.DeepEqual(after, before) {
				t.Fatal("pre-mutation refusal changed canonical assets")
			}
		})
	}
}

func TestAssetsSyncRollbackRestoresCompleteAssetPreimage(t *testing.T) {
	t.Parallel()

	targetRepo, assetRoot := newAssetsSyncTarget(t)
	sourceDir, _ := newAssetsSyncSource(t, assetRoot)
	before := captureAssetsSyncTree(t, assetRoot)
	replacements := 0

	payload, err := syncAssets(context.Background(), AssetsSyncRequest{
		SourceDir: sourceDir,
	}, assetsSyncDependencies{
		repository: targetRepo,
		assetRoot:  assetRoot,
		transactionHook: func(point transactionFaultPoint) error {
			if point.Phase == transactionPhaseReplacing {
				replacements++
				if replacements == 2 {
					return errors.New("injected asset refresh failure")
				}
			}
			return nil
		},
	})

	var syncErr *AssetsSyncError
	if !errors.As(err, &syncErr) || syncErr.Category != AssetsSyncExecution {
		t.Fatalf("rollback error = %T %v payload=%+v", err, err, payload)
	}
	if after := captureAssetsSyncTree(t, assetRoot); !reflect.DeepEqual(after, before) {
		t.Fatal("failed refresh did not restore the complete asset preimage")
	}
}

func TestAssetsSyncCompatibilityMatchesMaintainedPythonContract(t *testing.T) {
	t.Parallel()

	targetRepo, assetRoot := newAssetsSyncTarget(t)
	sourceDir, revision := newAssetsSyncSource(t, assetRoot)

	payload, err := syncAssets(context.Background(), AssetsSyncRequest{
		SourceDir: sourceDir,
	}, assetsSyncDependencies{
		repository: targetRepo,
		assetRoot:  assetRoot,
	})
	if err != nil {
		t.Fatal(err)
	}
	if payload.SchemaVersion != AssetsSyncSchemaVersion ||
		payload.Summary != (AssetsSyncSummary{Info: 3}) ||
		len(payload.PlannedChanges) != 0 {
		t.Fatalf("maintained result shape = %+v", payload)
	}
	for _, finding := range payload.Findings {
		if finding.Code != "skills.setup-snapshot.updated" ||
			finding.Severity != "info" ||
			!strings.HasPrefix(finding.ManagedID, "setup.") {
			t.Fatalf("maintained finding shape = %+v", finding)
		}
	}

	var fixture struct {
		Classification string   `json:"classification"`
		SourceTests    []string `json:"sourceTests"`
		Manifest       struct {
			Setups []map[string]any `json:"setups"`
		} `json:"manifest"`
	}
	readAssetsSyncJSON(
		t,
		filepath.Join("testdata", "parity-corpus", "v1", "fixtures", "asset-sync.json"),
		&fixture,
	)
	if fixture.Classification != "exact" ||
		!reflect.DeepEqual(fixture.SourceTests, []string{"tests/test_sync_setups.py"}) {
		t.Fatalf("asset-sync parity row = %+v", fixture)
	}
	for _, expected := range fixture.Manifest.Setups {
		setupID, ok := expected["id"].(string)
		if !ok {
			t.Fatalf("fixture setup id = %#v", expected["id"])
		}
		source, ok := expected["source"].(map[string]any)
		if !ok {
			t.Fatalf("%s fixture source = %#v", setupID, expected["source"])
		}
		source["ref"] = revision
		skills, ok := expected["skills"].([]any)
		if !ok {
			t.Fatalf("%s fixture skills = %#v", setupID, expected["skills"])
		}
		for _, rawSkill := range skills {
			skill := rawSkill.(map[string]any)
			skillSource := skill["source"].(map[string]any)
			if skillSource["type"] == "github" {
				skillSource["ref"] = revision
			}
		}
		digestPayload := any(skills)
		if bundles, exists := expected["activationBundles"]; exists {
			digestPayload = map[string]any{
				"activationBundles": bundles,
				"skills":            skills,
			}
		}
		digest, err := canonicalSHA256(digestPayload)
		if err != nil {
			t.Fatal(err)
		}
		expected["digest"] = digest
		expectedDocument, err := json.Marshal(expected)
		if err != nil {
			t.Fatal(err)
		}
		var expectedSnapshot assetsSyncSnapshot
		if err := json.Unmarshal(expectedDocument, &expectedSnapshot); err != nil {
			t.Fatal(err)
		}
		expectedBytes, err := json.MarshalIndent(expectedSnapshot, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		expectedBytes = append(expectedBytes, '\n')
		actualBytes, err := os.ReadFile(filepath.Join(assetRoot, "setups", setupID+".json"))
		if err != nil {
			t.Fatal(err)
		}
		if string(actualBytes) != string(expectedBytes) {
			t.Fatalf("%s normalized bytes differ from maintained Python asset-sync row", setupID)
		}
	}

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("embedded catalog depends on installed skill state: %v", err)
	}
	if len(catalog.SetupIDs()) == 0 {
		t.Fatal("embedded catalog has no setup snapshots")
	}
}

// assetsSyncTemplate holds one built copy of each Assets Sync fixture.
//
// Building a fixture means copying the whole embedded asset tree and then
// running git init, config, add, and commit over hundreds of files — around
// seven subprocesses per call. Roughly a dozen tests each paid that, and the
// four heaviest were 68s of this package's 271s of serial work.
//
// The fixtures are identical every time, so they are built once and each test
// receives a directory copy, .git included. Copying a tree costs a fraction of
// rebuilding a repository, and every test still owns a private copy it is free
// to mutate.
var assetsSyncTemplate struct {
	once     sync.Once
	root     string
	revision string
	err      error
}

func assetsSyncTemplateRoot(t *testing.T) (string, string) {
	t.Helper()
	assetsSyncTemplate.once.Do(func() {
		root, err := os.MkdirTemp("", "assets-sync-template")
		if err != nil {
			assetsSyncTemplate.err = err
			return
		}
		assetsSyncTemplate.root = root

		repository := filepath.Join(root, "target")
		assetRoot := filepath.Join(repository, "internal", "baseline", "assets")
		copyAssetsSyncFS(t, embeddedAssets, assetRoot)
		runAssetsSyncGit(t, repository, "init", "--quiet")
		gittest.AppendConfig(t, repository, "[user]\n\temail = fixture@example.com\n\tname = Fixture\n[commit]\n\tgpgsign = false\n")
		runAssetsSyncGit(t, repository, "add", ".")
		runAssetsSyncGit(t, repository, "commit", "--quiet", "-m", "asset target")

		checkout := filepath.Join(root, "source")
		assetsSyncTemplate.revision = buildAssetsSyncSource(t, checkout, assetRoot)
	})
	if assetsSyncTemplate.err != nil {
		t.Fatal(assetsSyncTemplate.err)
	}
	return assetsSyncTemplate.root, assetsSyncTemplate.revision
}

// removeAssetsSyncTemplate drops the shared fixture directory once the
// package's tests are done with it.
func removeAssetsSyncTemplate() {
	if assetsSyncTemplate.root != "" {
		_ = os.RemoveAll(assetsSyncTemplate.root)
	}
}

func copyAssetsSyncDir(t *testing.T, source string, target string) {
	t.Helper()
	if err := os.CopyFS(target, os.DirFS(source)); err != nil {
		t.Fatalf("copy fixture %s: %v", source, err)
	}
}

func newAssetsSyncTarget(t *testing.T) (string, string) {
	t.Helper()
	root, _ := assetsSyncTemplateRoot(t)
	repository := t.TempDir()
	copyAssetsSyncDir(t, filepath.Join(root, "target"), repository)
	return repository, filepath.Join(repository, "internal", "baseline", "assets")
}

func newAssetsSyncSource(t *testing.T, assetRoot string) (string, string) {
	t.Helper()
	root, revision := assetsSyncTemplateRoot(t)
	checkout := t.TempDir()
	copyAssetsSyncDir(t, filepath.Join(root, "source"), checkout)
	return filepath.Join(checkout, "setups"), revision
}

func buildAssetsSyncSource(t *testing.T, checkout string, assetRoot string) string {
	t.Helper()
	sourceDir := filepath.Join(checkout, "setups")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	paths := map[string]struct{}{}
	entries, err := filepath.Glob(filepath.Join(assetRoot, "setups", "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, snapshotPath := range entries {
		var snapshot struct {
			ID     string `json:"id"`
			Skills []struct {
				Path string `json:"path"`
			} `json:"skills"`
		}
		readAssetsSyncJSON(t, snapshotPath, &snapshot)
		lines := []string{"# canonical setup"}
		for _, skill := range snapshot.Skills {
			lines = append(lines, skill.Path)
			paths[skill.Path] = struct{}{}
		}
		if err := os.WriteFile(
			filepath.Join(sourceDir, snapshot.ID+".txt"),
			[]byte(strings.Join(lines, "\n")+"\n"),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}
	for skillPath := range paths {
		target := filepath.Join(checkout, filepath.FromSlash(skillPath))
		if filepath.Base(target) != "SKILL.md" {
			target = filepath.Join(target, "SKILL.md")
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte("# "+filepath.Base(filepath.Dir(target))+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runAssetsSyncGit(t, checkout, "init", "--quiet")
	gittest.AppendConfig(t, checkout, "[user]\n\temail = fixture@example.com\n\tname = Fixture\n[commit]\n\tgpgsign = false\n")
	runAssetsSyncGit(t, checkout, "remote", "add", "origin", "https://github.com/example/skills.git")
	runAssetsSyncGit(t, checkout, "add", ".")
	runAssetsSyncGit(t, checkout, "commit", "--quiet", "-m", "canonical source")
	return strings.TrimSpace(runAssetsSyncGit(t, checkout, "rev-parse", "HEAD"))
}

func copyAssetsSyncFS(t *testing.T, source fs.FS, target string) {
	t.Helper()
	if err := fs.WalkDir(source, ".", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		destination := filepath.Join(target, filepath.FromSlash(name))
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		data, err := fs.ReadFile(source, name)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, data, 0o644)
	}); err != nil {
		t.Fatal(err)
	}
}

func captureAssetsSyncTree(t *testing.T, root string) map[string]string {
	t.Helper()
	state := map[string]string{}
	if err := filepath.WalkDir(root, func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		state[filepath.ToSlash(relative)] = string(data)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return state
}

func runAssetsSyncGit(t *testing.T, directory string, args ...string) string {
	t.Helper()
	gitArgs := append(
		[]string{
			"-C", directory,
			"-c", "core.fsmonitor=false",
			"-c", "maintenance.auto=false",
		},
		args...,
	)
	command := exec.Command("git", gitArgs...)
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func readAssetsSyncJSON(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

func writeAssetsSyncJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertAssetsSyncOwnedSkillHasNoContentPin(t *testing.T, snapshotPath string) {
	t.Helper()
	var snapshot struct {
		Skills []map[string]any `json:"skills"`
	}
	readAssetsSyncJSON(t, snapshotPath, &snapshot)
	for _, skill := range snapshot.Skills {
		if skill["name"] != "setup-context-driven" {
			continue
		}
		for _, field := range []string{"treeDigest", "contentDigest"} {
			if _, exists := skill[field]; exists {
				t.Errorf("setup-context-driven retains compatibility content pin %s", field)
			}
		}
		return
	}
	t.Fatal("setup-context-driven entry is missing")
}

func assetsSyncOwnedMinimum(t *testing.T, snapshotPath string) string {
	t.Helper()
	var snapshot struct {
		Skills []struct {
			Name           string `json:"name"`
			MinimumVersion string `json:"minimumVersion"`
		} `json:"skills"`
	}
	readAssetsSyncJSON(t, snapshotPath, &snapshot)
	for _, skill := range snapshot.Skills {
		if skill.Name == "setup-context-driven" {
			if skill.MinimumVersion == "" {
				t.Fatal("setup-context-driven minimum is missing")
			}
			return skill.MinimumVersion
		}
	}
	t.Fatal("setup-context-driven entry is missing")
	return ""
}

// TestMain drops the shared Assets Sync fixture after the package's tests
// finish. The template outlives every individual test that copies it, so no
// single test can own its cleanup.
const baselineDigestRegenerationCommand = "make baseline-digests"

func declareBaselineDigestRegeneration() {
	suiteguard.DeclareSanctionedRegeneration(baselineDigestRegenerationCommand)
}

func TestMain(m *testing.M) {
	code := suiteguard.Main(m, filepath.Join("..", ".."))
	removeAssetsSyncTemplate()
	os.Exit(code)
}
