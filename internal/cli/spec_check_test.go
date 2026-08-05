// Suite: Spec commands
// Invariant: the CLI exposes Spec checks and close audits through stable streams, schemas, and exit codes.
// Boundary IN: public Run dispatch, configured Spec Root resolution, command renderers, and real Git and Run Database fixtures
// Boundary OUT: checker and audit classification correctness, owned by internal/speccheck and internal/specaudit
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/specaudit"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

func TestRunSpecCheckCleanText(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "Spec clean\nNo findings.\n") {
		t.Fatalf("stdout = %q, want clean report", stdout.String())
	}
	if strings.Contains(stdout.String(), "[error]") || strings.Contains(stdout.String(), "[gap]") {
		t.Fatalf("clean stdout contains a finding:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecCheckErrorText(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "tooling-unauthorized"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{
		"[error] " + speccheck.CodeToolingUnauthorized,
		"docs/specs/tooling-unauthorized/_prd.md:",
		"docs/workflow/authorizations/other-spec.md:",
		"  fix: ",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecCheckGapStrictPromotion(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecCheckWorkspace(t, "citation-dirty")
	techSpecPath := filepath.Join(repoDir, "docs", "specs", "citation-dirty", "_techspec.md")
	techSpec := strings.ReplaceAll(mustRead(t, techSpecPath), "\n## Decisions\n\nThe implementation also follows ADR-0103.\n", "")
	mustWrite(t, techSpecPath, techSpec)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"spec", "check", "citation-dirty"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("gap-only spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[gap] "+speccheck.CodeADRRelated) {
		t.Fatalf("gap-only stdout does not contain gap finding:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "check", "citation-dirty", "--strict"}, &stdout, &stderr)
	if code != exitRunFailed {
		t.Fatalf("strict gap exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	if !strings.Contains(stdout.String(), "[error] "+speccheck.CodeADRRelated) {
		t.Fatalf("strict stdout does not contain promoted error:\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "[gap] "+speccheck.CodeADRRelated) {
		t.Fatalf("strict stdout retained gap severity:\n%s", stdout.String())
	}
}

func TestRunSpecCheckJSONWritesOneObjectPerSpec(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean", "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "check", "clean", "tooling-unauthorized", "--format", "json",
	}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("multi-Spec JSON exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("JSON lines = %d, want 2; stdout=%q", len(lines), stdout.String())
	}
	for index, wantSlug := range []string{"clean", "tooling-unauthorized"} {
		var document struct {
			Schema string `json:"schema"`
			Slug   string `json:"slug"`
		}
		if err := json.Unmarshal([]byte(lines[index]), &document); err != nil {
			t.Fatalf("line %d is not JSON: %v; line=%q", index+1, err, lines[index])
		}
		if document.Schema != speccheck.SchemaVersion || document.Slug != wantSlug {
			t.Errorf("line %d document = %#v, want schema %q slug %q", index+1, document, speccheck.SchemaVersion, wantSlug)
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecCheckUnknownSlugIsUsageError(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "no-such-slug"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown slug exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no partial report", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no-such-slug") {
		t.Fatalf("stderr does not name unknown slug: %q", stderr.String())
	}
}

func TestRunSpecCheckWithoutSlugChecksEveryActiveSpec(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean", "tooling-unauthorized")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("all-active exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, slug := range []string{"clean", "tooling-unauthorized"} {
		count := 0
		for _, line := range strings.Split(stdout.String(), "\n") {
			if line == "Spec "+slug {
				count++
			}
		}
		if count != 1 {
			t.Errorf("stdout does not contain exactly one report for %q:\n%s", slug, stdout.String())
		}
	}
}

func TestRunSpecCheckUnreadableRootIsUsageError(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecCheckWorkspace(t, "clean")
	missingRoot := filepath.Join(repoDir, "missing-spec-root")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), "specs:\n  root: "+missingRoot+"\n")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unreadable root exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no report", stdout.String())
	}
	if !strings.Contains(stderr.String(), missingRoot) || !strings.Contains(stderr.String(), "does not exist") {
		t.Fatalf("stderr does not name unreadable Spec Root: %q", stderr.String())
	}
}

func TestRunSpecCheckHelpAppearsInTopLevelUsageAndCommandList(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("top-level help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{
		"roundfix spec check [<slug> ...] [--format <text|json>] [--strict]",
		"spec       Check Spec artifact consistency",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("top-level help does not contain %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "check", "--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("spec check help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{"--format", "--strict", "0  no errors", "1  at least one error", "2  usage error"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("spec check help does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("help stderr = %q, want empty", stderr.String())
	}
}

func TestRunSpecAuditCleanText(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	archiveSpecAuditFixture(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	want := "Spec audit " + specAuditFixtureSlug + "\nNo residue or undelivered work.\n"
	if stdout.String() != want {
		t.Fatalf("stdout = %q, want clean report %q", stdout.String(), want)
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditResidueText(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	const branch = "ma/spec-audit-residue"
	gitImplement(t, repoDir, "checkout", "-b", branch)
	mustWrite(t, filepath.Join(repoDir, "residue.txt"), "residue\n")
	gitImplement(t, repoDir, "add", "residue.txt")
	gitImplement(
		t,
		repoDir,
		"commit",
		"-m", "feat: add residue fixture",
		"-m", "Roundfix-Spec: "+specAuditFixtureSlug,
	)
	gitImplement(t, repoDir, "checkout", "main")
	gitImplement(t, repoDir, "merge", "--squash", branch)
	gitImplement(t, repoDir, "commit", "-m", "feat: merge residue fixture")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{
		"residue branch " + branch,
		"evidence: survivor content is fully represented",
		"reclaim: git branch -D -- '" + branch + "'",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditUndeliveredTextNamesHoldingBranch(t *testing.T) {
	t.Parallel()
	_, repoDir := newSpecAuditWorkspace(t)
	const branch = "ma/spec-audit-archive"
	gitImplement(t, repoDir, "checkout", "-b", branch)
	archivedPath := filepath.ToSlash(filepath.Join("docs", "specs", "_archived", specAuditFixtureSlug))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(repoDir, archivedPath)), 0o755); err != nil {
		t.Fatalf("create archived Spec parent: %v", err)
	}
	gitImplement(t, repoDir, "mv", filepath.ToSlash(filepath.Join("docs", "specs", specAuditFixtureSlug)), archivedPath)
	gitImplement(
		t,
		repoDir,
		"commit",
		"-m", "docs: archive Spec audit fixture",
		"-m", "Roundfix-Spec: "+specAuditFixtureSlug,
	)
	gitImplement(t, repoDir, "checkout", "main")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "audit", specAuditFixtureSlug}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitRunFailed, stderr.String())
	}
	for _, want := range []string{archivedPath, "held by: " + branch} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRunSpecAuditJSONWritesOneObject(t *testing.T) {
	t.Parallel()
	_, _ = newSpecAuditWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "audit", specAuditFixtureSlug, "--format", "json",
	}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec audit exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	decoder := json.NewDecoder(strings.NewReader(stdout.String()))
	var document struct {
		Schema        string `json:"schema"`
		SchemaVersion string `json:"schemaVersion"`
		Type          string `json:"type"`
		OK            bool   `json:"ok"`
		Slug          string `json:"slug"`
	}
	if err := decoder.Decode(&document); err != nil {
		t.Fatalf("stdout is not JSON: %v; stdout=%q", err, stdout.String())
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("stdout contains more than one JSON object: %v; stdout=%q", err, stdout.String())
	}
	if document.Schema != specAuditSchemaVersion ||
		document.SchemaVersion != specAuditSchemaVersion ||
		document.Type != specAuditDocumentType ||
		!document.OK ||
		document.Slug != specAuditFixtureSlug {
		t.Fatalf(
			"document = %#v, want schema %q type %q ok true slug %q",
			document,
			specAuditSchemaVersion,
			specAuditDocumentType,
			specAuditFixtureSlug,
		)
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want no diagnostics", stderr.String())
	}
}

func TestRenderSpecAuditJSONReportsAttention(t *testing.T) {
	t.Parallel()
	data, err := renderSpecAuditJSON(specaudit.Result{
		Slug: specAuditFixtureSlug,
		Survivors: []specaudit.Survivor{{
			Name: "ma/spec-audit-residue",
			Kind: specaudit.KindResidue,
		}},
	})
	if err != nil {
		t.Fatalf("render Spec audit JSON: %v", err)
	}
	var document struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		t.Fatalf("decode Spec audit JSON: %v", err)
	}
	if document.OK {
		t.Fatal("attention-requiring Spec audit JSON ok = true, want false")
	}
}

func TestRunSpecAuditUnknownSlugIsUsageError(t *testing.T) {
	t.Parallel()
	_, _ = newSpecAuditWorkspace(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{
		"spec", "audit", "no-such-slug", "--format", "json",
	}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("unknown slug exit = %d, want %d", code, exitPreflight)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want no partial JSON", stdout.String())
	}
	if !strings.Contains(stderr.String(), "no-such-slug") {
		t.Fatalf("stderr does not name unknown slug: %q", stderr.String())
	}
}

func TestRunSpecAuditHelpAppearsInUsageAndCommandList(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("top-level help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{
		"roundfix spec audit <slug> [--format <text|json>]",
		"spec       Check Spec artifact consistency; audit Spec delivery",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("top-level help does not contain %q:\n%s", want, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"spec", "audit", "--help"}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("spec audit help exit = %d, want %d", code, exitOK)
	}
	for _, want := range []string{"--format", "0  no residue", "1  residue or undelivered work", "2  usage error or unknown Spec slug"} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("spec audit help does not contain %q:\n%s", want, stdout.String())
		}
	}
	if stderr.String() != "" {
		t.Fatalf("help stderr = %q, want empty", stderr.String())
	}
}

func TestRunSpecAuditPreservesSpecCheckBehavior(t *testing.T) {
	t.Parallel()
	_, _ = newSpecCheckWorkspace(t, "clean")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"spec", "check", "clean"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("spec check exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if !strings.HasPrefix(stdout.String(), "Spec clean\nNo findings.\n") {
		t.Fatalf("spec check stdout = %q, want unchanged clean report prefix", stdout.String())
	}
	if strings.Contains(stdout.String(), "[error]") || strings.Contains(stdout.String(), "[gap]") {
		t.Fatalf("spec check clean stdout contains a finding:\n%s", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("spec check stderr = %q, want empty", stderr.String())
	}
}

const specAuditFixtureSlug = "0068-spec-close-audit"

func newSpecAuditWorkspace(t *testing.T) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	gittest.InitRepo(t, repoDir, "--initial-branch=main")
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = roundfix-test@example.com\n[commit]\n\tgpgsign = false\n")
	if err := os.MkdirAll(filepath.Join(repoDir, "docs", "specs", specAuditFixtureSlug), 0o755); err != nil {
		t.Fatalf("create fixture Spec directory: %v", err)
	}
	mustWrite(
		t,
		filepath.Join(repoDir, "docs", "specs", specAuditFixtureSlug, "_prd.md"),
		"---\nspec: "+specAuditFixtureSlug+"\nstatus: active\n---\n",
	)
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "docs: seed Spec audit fixture")
	runStore, err := store.Open(context.Background(), homeDir)
	if err != nil {
		t.Fatalf("open fixture Run Database: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close fixture Run Database: %v", err)
	}
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve repo dir: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, resolved)
	return homeDir, resolved
}

func archiveSpecAuditFixture(t *testing.T, repoDir string) {
	t.Helper()
	activePath := filepath.ToSlash(filepath.Join("docs", "specs", specAuditFixtureSlug))
	archivedPath := filepath.ToSlash(filepath.Join("docs", "specs", "_archived", specAuditFixtureSlug))
	if err := os.MkdirAll(filepath.Dir(filepath.Join(repoDir, archivedPath)), 0o755); err != nil {
		t.Fatalf("create archived Spec parent: %v", err)
	}
	gitImplement(t, repoDir, "mv", activePath, archivedPath)
	gitImplement(t, repoDir, "commit", "-m", "docs: archive Spec audit fixture")
}

func newSpecCheckWorkspace(t *testing.T, slugs ...string) (string, string) {
	t.Helper()
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	fixtureRepo := filepath.Join(cliTestRepoRoot(t), "internal", "speccheck", "testdata", "repo")
	for _, relative := range []string{
		filepath.Join("docs", "agents"),
		filepath.Join("docs", "adr"),
		filepath.Join("docs", "workflow"),
	} {
		if err := copyDir(filepath.Join(fixtureRepo, relative), filepath.Join(repoDir, relative)); err != nil {
			t.Fatalf("copy fixture %s: %v", relative, err)
		}
	}
	for _, slug := range slugs {
		relative := filepath.Join("docs", "specs", slug)
		if err := copyDir(filepath.Join(fixtureRepo, relative), filepath.Join(repoDir, relative)); err != nil {
			t.Fatalf("copy fixture Spec %s: %v", slug, err)
		}
	}

	gittest.InitRepo(t, repoDir, "--initial-branch=main")
	gittest.AppendConfig(t, repoDir, "[user]\n\tname = Roundfix Test\n\temail = roundfix-test@example.com\n[commit]\n\tgpgsign = false\n")
	gitImplement(t, repoDir, "add", "-A")
	gitImplement(t, repoDir, "commit", "-m", "seed Spec Consistency Check fixtures")
	gitImplement(t, repoDir, "checkout", "-b", "ma/spec-check")
	resolved, err := filepath.EvalSymlinks(repoDir)
	if err != nil {
		t.Fatalf("resolve repo dir: %v", err)
	}
	setCommandEnvironmentForTest(t, homeDir, resolved)
	return homeDir, resolved
}
