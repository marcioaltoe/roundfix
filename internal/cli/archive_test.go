package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

func archiveTestPath(kind spec.ArchiveKind, elements ...string) string {
	parts := append([]string{filepath.FromSlash(spec.ArchiveDir(kind))}, elements...)
	return filepath.ToSlash(filepath.Join(parts...))
}

func archiveTestRepositoryPath(repoRoot string, kind spec.ArchiveKind, elements ...string) string {
	return filepath.Join(repoRoot, filepath.FromSlash(archiveTestPath(kind, elements...)))
}

func TestRunArchiveMovesCompletedSpecAndStampsMetadata(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Wire the widget API", status: string(spec.StatusCompleted), needs: []string{"task_01"}},
	})
	writeArchiveQAReport(t, repoDir, spec.VerdictPass, "rows_blocked_environment: 3")
	withNoEngineCollaborators(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected archive exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	wantStdout := "archived " + implementTestSlug + " -> " + archiveTestPath(spec.ArchiveKindSpec, implementTestSlug) + "\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout %q, got %q", wantStdout, stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected archive stderr empty, got %q", stderr.String())
	}
	assertPathMissing(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	archivedDir := archiveTestRepositoryPath(repoDir, spec.ArchiveKindSpec, implementTestSlug)
	assertPathExists(t, archivedDir)
	assertPathExists(t, filepath.Join(archivedDir, "task_01.md"))
	prd := mustRead(t, filepath.Join(archivedDir, "_prd.md"))
	for _, want := range []string{
		"status: archived",
		"source_slug: " + implementTestSlug,
		"archived:",
	} {
		if !strings.Contains(prd, want) {
			t.Fatalf("expected archived PRD to contain %q, got:\n%s", want, prd)
		}
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveDeclaredUnreachableContract(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		verdict            string
		blockedDeclared    int
		blockedFinding     int
		blockedEnvironment int
		declarations       []string
		wantUnproven       []string
		wantStderr         []string
	}{
		{
			name:         "passing report remains unchanged",
			verdict:      spec.VerdictPass,
			declarations: []string{"a maintainer publishes the tagged release"},
		},
		{
			name:            "declared-only partial report archives",
			verdict:         spec.VerdictPartial,
			blockedDeclared: 2,
			declarations: []string{
				"a maintainer publishes the tagged release",
				"a maintainer records the production identity exchange",
			},
			wantUnproven: []string{
				"a maintainer publishes the tagged release",
				"a maintainer records the production identity exchange",
			},
		},
		{
			name:            "surplus declarations still cover declared rows",
			verdict:         spec.VerdictPartial,
			blockedDeclared: 1,
			declarations: []string{
				"a maintainer publishes the tagged release",
				"a maintainer records the production identity exchange",
			},
			wantUnproven: []string{
				"a maintainer publishes the tagged release",
				"a maintainer records the production identity exchange",
			},
		},
		{
			name:            "finding-blocked partial report refuses",
			verdict:         spec.VerdictPartial,
			blockedDeclared: 1,
			blockedFinding:  2,
			declarations:    []string{"a maintainer publishes the tagged release"},
			wantStderr:      []string{"rows_blocked_finding is 2", "expected 0"},
		},
		{
			name:               "environment-blocked partial report refuses",
			verdict:            spec.VerdictPartial,
			blockedDeclared:    1,
			blockedEnvironment: 3,
			declarations:       []string{"a maintainer publishes the tagged release"},
			wantStderr:         []string{"rows_blocked_environment is 3", "expected 0"},
		},
		{
			name:            "declaration shortfall refuses",
			verdict:         spec.VerdictPartial,
			blockedDeclared: 3,
			declarations:    []string{"a maintainer publishes the tagged release"},
			wantStderr: []string{
				"rows_blocked_declared is 3",
				"Spec declares 1 unreachable acceptance",
				"shortfall is 2",
			},
		},
		{
			name:    "failing report refuses exactly as before",
			verdict: spec.VerdictFail,
			wantStderr: []string{
				"no passing QA verdict",
				"newest QA Report verdict is \"fail\"",
				"expected \"pass\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
			})
			if len(tt.declarations) > 0 {
				appendArchiveUnreachableDeclarations(t, repoDir, tt.declarations)
			}
			writeArchiveQAReport(t, repoDir, tt.verdict,
				fmt.Sprintf("rows_blocked_declared: %d", tt.blockedDeclared),
				fmt.Sprintf("rows_blocked_finding: %d", tt.blockedFinding),
				fmt.Sprintf("rows_blocked_environment: %d", tt.blockedEnvironment),
			)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

			activeDir := filepath.Join(repoDir, "docs", "specs", implementTestSlug)
			archivedDir := archiveTestRepositoryPath(repoDir, spec.ArchiveKindSpec, implementTestSlug)
			if len(tt.wantStderr) > 0 {
				if code != exitPreflight {
					t.Fatalf("expected archive refusal exit %d, got %d stderr=%q stdout=%q", exitPreflight, code, stderr.String(), stdout.String())
				}
				if stdout.String() != "" {
					t.Fatalf("expected refusal stdout empty, got %q", stdout.String())
				}
				for _, want := range tt.wantStderr {
					if !strings.Contains(stderr.String(), want) {
						t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
					}
				}
				assertPathExists(t, activeDir)
				assertPathMissing(t, archivedDir)
				assertNoRunDatabase(t, homeDir)
				return
			}

			if code != exitOK {
				t.Fatalf("expected archive exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
			}
			if stderr.String() != "" {
				t.Fatalf("expected archive stderr empty, got %q", stderr.String())
			}
			assertPathMissing(t, activeDir)
			assertPathExists(t, archivedDir)
			prdPath := filepath.Join(archivedDir, "_prd.md")
			gotUnproven := readArchivedUnproven(t, prdPath)
			if len(gotUnproven) != len(tt.wantUnproven) {
				t.Fatalf("archived PRD unproven = %q, want %q", gotUnproven, tt.wantUnproven)
			}
			for index := range tt.wantUnproven {
				if gotUnproven[index] != tt.wantUnproven[index] {
					t.Fatalf("archived PRD unproven[%d] = %q, want %q", index, gotUnproven[index], tt.wantUnproven[index])
				}
			}
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunArchiveUsesConfiguredExternalSpecRoot(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Internal fixture should stay active", status: string(spec.StatusCompleted)},
	})
	externalRoot := filepath.Join(t.TempDir(), "external-specs")
	writeImplementSpecAtRoot(t, externalRoot, implementTestSlug, []implementSeed{
		{id: "task_01", title: "Archive external task", status: string(spec.StatusCompleted)},
	})
	reportPath := filepath.Join(externalRoot, implementTestSlug, "qa", "qa-report-2026-07-06.md")
	mustMkdir(t, filepath.Dir(reportPath))
	mustWrite(t, reportPath, implementQAReport(spec.VerdictPass))
	configureExternalSpecsRoot(t, repoDir, externalRoot)
	withNoEngineCollaborators(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected archive exit 0, got %d stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	archivedDir := filepath.Join(spec.ArchiveSpecRoot(externalRoot, false), implementTestSlug)
	wantStdout := "archived " + implementTestSlug + " -> " + filepath.ToSlash(archivedDir) + "\n"
	if stdout.String() != wantStdout {
		t.Fatalf("expected stdout %q, got %q", wantStdout, stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected archive stderr empty, got %q", stderr.String())
	}
	assertPathMissing(t, filepath.Join(externalRoot, implementTestSlug))
	assertPathExists(t, filepath.Join(archivedDir, "task_01.md"))
	assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveRefusesIncompleteTask(t *testing.T) {
	t.Parallel()
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
		{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
		{id: "task_02", title: "Wire the widget API", status: string(spec.StatusPending), needs: []string{"task_01"}},
	})
	writeArchiveQAReport(t, repoDir, spec.VerdictPass)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected archive refusal exit %d, got %d", exitPreflight, code)
	}
	if stdout.String() != "" {
		t.Fatalf("expected refusal stdout empty, got %q", stdout.String())
	}
	for _, want := range []string{
		"Task \"task_02\" is \"pending\"",
		"archive requires every Task to be \"completed\"",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
		}
	}
	assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
	assertPathMissing(t, archiveTestRepositoryPath(repoDir, spec.ArchiveKindSpec, implementTestSlug))
	prd := mustRead(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_prd.md"))
	if strings.Contains(prd, "status: archived") || strings.Contains(prd, "source_slug: "+implementTestSlug) {
		t.Fatalf("expected active PRD left unstamped, got:\n%s", prd)
	}
	assertNoRunDatabase(t, homeDir)
}

func TestRunArchiveRefusesMissingOrNonPassingQA(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		verdict          string
		extraFrontmatter []string
		wantStderr       []string
	}{
		{
			name:    "missing QA Report",
			verdict: "",
			wantStderr: []string{
				"no passing QA verdict",
				"no QA Report found",
			},
		},
		{
			name:    "failing QA Report",
			verdict: spec.VerdictFail,
			wantStderr: []string{
				"no passing QA verdict",
				"newest QA Report verdict is \"fail\"",
				"expected \"pass\"",
			},
		},
		{
			name:    "partial QA Report",
			verdict: spec.VerdictPartial,
			wantStderr: []string{
				"no passing QA verdict",
				"newest QA Report verdict is \"partial\"",
				"expected \"pass\"",
			},
		},
		{
			name:             "finding-blocked pass QA Report",
			verdict:          spec.VerdictPass,
			extraFrontmatter: []string{"rows_blocked_finding: 1"},
			wantStderr: []string{
				"no passing QA verdict",
				"QA Report",
				"is unreadable",
				"rows_blocked_finding must be zero when verdict is \"pass\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := newImplementWorkspace(t, []implementSeed{
				{id: "task_01", title: "Build the widget core", status: string(spec.StatusCompleted)},
			})
			if tt.verdict != "" {
				writeArchiveQAReport(t, repoDir, tt.verdict, tt.extraFrontmatter...)
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("expected archive refusal exit %d, got %d", exitPreflight, code)
			}
			if stdout.String() != "" {
				t.Fatalf("expected refusal stdout empty, got %q", stdout.String())
			}
			for _, want := range tt.wantStderr {
				if !strings.Contains(stderr.String(), want) {
					t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
				}
			}
			assertPathExists(t, filepath.Join(repoDir, "docs", "specs", implementTestSlug))
			assertPathMissing(t, archiveTestRepositoryPath(repoDir, spec.ArchiveKindSpec, implementTestSlug))
			assertNoRunDatabase(t, homeDir)
		})
	}
}

func TestRunArchiveHelp(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(t, context.Background(), []string{"archive", "--help"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected help exit 0, got %d", code)
	}
	if stderr.String() != "" {
		t.Fatalf("expected help stderr empty, got %q", stderr.String())
	}
	for _, want := range []string{
		"Usage:",
		"roundfix archive <slug>",
		"covered only by declared Unreachable Acceptance",
		"archive creates no Run and",
		"never pushes",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected help to contain %q, got %q", want, stdout.String())
		}
	}
}

func writeArchiveQAReport(t *testing.T, repoDir string, verdict string, extraFrontmatter ...string) {
	t.Helper()
	reportPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "qa", "qa-report-2026-07-06.md")
	mustMkdir(t, filepath.Dir(reportPath))
	fields := append([]string{"verdict: " + verdict}, extraFrontmatter...)
	report := "---\n" + strings.Join(fields, "\n") + "\n---\n\n# QA Report\n"
	mustWrite(t, reportPath, report)
}

func appendArchiveUnreachableDeclarations(t *testing.T, repoDir string, actions []string) {
	t.Helper()
	prdPath := filepath.Join(repoDir, "docs", "specs", implementTestSlug, "_prd.md")
	var section strings.Builder
	section.WriteString("\n## Unreachable Acceptance\n")
	for index, action := range actions {
		section.WriteString(fmt.Sprintf("\n- criterion: acceptance criterion %d\n  reason: no hermetic Verification can reach it\n  satisfied-by: %s\n", index+1, action))
	}
	mustWrite(t, prdPath, mustRead(t, prdPath)+section.String())
}

func readArchivedUnproven(t *testing.T, prdPath string) []string {
	t.Helper()
	content := mustRead(t, prdPath)
	const opening = "---\n"
	if !strings.HasPrefix(content, opening) {
		t.Fatalf("archived PRD %q has no frontmatter", prdPath)
	}
	rest := content[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		t.Fatalf("archived PRD %q has no frontmatter closing marker", prdPath)
	}
	var frontmatter struct {
		Unproven []string `yaml:"unproven"`
	}
	if err := yaml.Unmarshal([]byte(rest[:end]), &frontmatter); err != nil {
		t.Fatalf("parse archived PRD %q frontmatter: %v", prdPath, err)
	}
	return frontmatter.Unproven
}

func assertPathExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected path %s to exist: %v", path, err)
	}
}

func assertPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected path %s to be missing", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("stat path %s: %v", path, err)
	}
}

func withNoEngineCollaborators(t *testing.T) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.newEngineCollaborators = func() engineCollaborators {
			t.Fatal("archive must not create Run engine collaborators")
			return engineCollaborators{}
		}
	})
}
