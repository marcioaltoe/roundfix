package skills

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// Suite: Baseline skill distribution contract
// Invariant: canonical and shipped Baseline guidance is identical and invokes only the public CLI.
// Boundary IN: setup-context-driven and Roundfix owned skill guidance.
// Boundary OUT: executable setup-runtime removal and Baseline CLI behavior.

func TestBaselineSkillContract(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	skillNames := []string{"setup-context-driven", "roundfix"}
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	embeddedByPath := make(map[string][]byte, len(embeddedFiles))
	for _, file := range embeddedFiles {
		embeddedByPath[file.Path] = file.Data
	}

	for _, name := range skillNames {
		t.Run(name, func(t *testing.T) {
			canonicalPath := filepath.Join(repoRoot, ".agents", "skills", name, "SKILL.md")
			distributedPath := filepath.Join(repoRoot, "skills", name, "SKILL.md")
			canonical := readBaselineSkillContractFile(t, canonicalPath)
			distributed := readBaselineSkillContractFile(t, distributedPath)
			if !bytes.Equal(canonical, distributed) {
				t.Fatalf("%s canonical and distributed guidance differ; run make skills-sync", name)
			}
			embedded := embeddedByPath[filepath.ToSlash(filepath.Join(name, "SKILL.md"))]
			if !bytes.Equal(distributed, embedded) {
				t.Fatalf("%s distributed and embedded guidance differ", name)
			}
			body := baselineSkillBody(string(canonical))
			if strings.Contains(strings.ToLower(body), "python") ||
				strings.Contains(body, "context_setup.py") ||
				strings.Contains(body, "scripts/context") {
				t.Fatalf("%s skill body invokes a retired independent setup runtime", name)
			}
			for _, forbidden := range []string{"/Users/", "/home/", `C:\`} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("%s skill body contains environment-specific path %q", name, forbidden)
				}
			}
		})
	}

	setup := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven", "SKILL.md"),
	))
	for _, required := range []string{
		"only runtime authority",
		"roundfix baseline --repo . --format text",
		"roundfix baseline plan",
		"roundfix baseline apply",
		"roundfix/baseline-plan/v1",
		"roundfix/baseline-result/v1",
		"Repository Skill Set restoration",
		"Canonical asset synchronization",
		"behavioral fallback",
	} {
		if !strings.Contains(setup, required) {
			t.Fatalf("setup-context-driven skill missing %q", required)
		}
	}
}

func TestNoPythonBaselineRuntime(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	for _, root := range []string{
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven"),
		filepath.Join(repoRoot, "skills", "setup-context-driven"),
	} {
		var files []string
		if err := filepath.WalkDir(root, func(filePath string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(relative))
			return nil
		}); err != nil {
			t.Fatalf("walk setup skill %s: %v", root, err)
		}
		sort.Strings(files)
		if !reflect.DeepEqual(files, []string{"SKILL.md"}) {
			t.Fatalf("setup skill %s ships non-guidance files: %v", root, files)
		}
	}

	makefile := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, "Makefile")))
	for _, forbidden := range []string{
		"setup-context-check",
		"PYTHONDONTWRITEBYTECODE",
		"python3",
	} {
		if strings.Contains(makefile, forbidden) {
			t.Fatalf("post-cutover Makefile invokes retired runtime marker %q", forbidden)
		}
	}

	for _, relative := range []string{
		"README.md",
		filepath.Join("docs", "user-guide", "commands.md"),
		filepath.Join("docs", "user-guide", "context-driven-development.md"),
		filepath.Join(".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(".agents", "skills", "roundfix", "SKILL.md"),
	} {
		content := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, relative)))
		for _, forbidden := range []string{
			"context_" + "setup.py",
			"context_" + "baseline.py",
			"python3",
			"Python fallback",
		} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s references retired Baseline runtime marker %q", relative, forbidden)
			}
		}
	}
}

func TestThinSetupSkill(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	var setupFiles []string
	for _, file := range files {
		if file.Skill == "setup-context-driven" {
			setupFiles = append(setupFiles, file.Path)
		}
	}
	sort.Strings(setupFiles)
	if !reflect.DeepEqual(setupFiles, []string{"setup-context-driven/SKILL.md"}) {
		t.Fatalf("embedded setup skill files = %v", setupFiles)
	}
	for _, diagnostic := range Check() {
		if strings.HasPrefix(diagnostic.Path, "setup-context-driven/") {
			t.Fatalf("thin setup skill diagnostic: %s: %s", diagnostic.Path, diagnostic.Message)
		}
	}

	setup := string(readBaselineSkillContractFile(
		t,
		filepath.Join("..", ".agents", "skills", "setup-context-driven", "SKILL.md"),
	))
	for _, required := range []string{
		"recipes and interpretation only",
		"no independent",
		"behavioral fallback",
		"Instruction Hierarchy",
		"exact source bytes",
		"Repository-Specific Normative Rules",
		"Only `accepted` is active",
		"`pending`, `partial`, `deferred`, and `done`",
		"Profile alignment",
		"Change Baseline",
		"repository-owned Profile adaptation",
		"Decline",
		"without writing",
		"--profile-file",
		"mutually exclusive",
		"Generate a fresh plan",
	} {
		if !strings.Contains(setup, required) {
			t.Errorf("thin setup skill missing %q", required)
		}
	}
	for _, forbidden := range []string{
		"classify rules itself",
		"render managed guidance itself",
		"write Profile files directly",
		"mutate repository files itself",
	} {
		if strings.Contains(setup, forbidden) {
			t.Errorf("thin setup skill claims independent behavior %q", forbidden)
		}
	}
}

func TestAuthorialSkillSync(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded authorial skills: %v", err)
	}
	embeddedByPath := make(map[string][]byte, len(embeddedFiles))
	for _, file := range embeddedFiles {
		embeddedByPath[file.Path] = file.Data
	}

	for _, skillName := range Names() {
		t.Run(skillName, func(t *testing.T) {
			canonicalRoot := filepath.Join(repoRoot, ".agents", "skills", skillName)
			distributedRoot := filepath.Join(repoRoot, "skills", skillName)
			assertSkillTreesEqual(t, canonicalRoot, distributedRoot)

			canonicalDigest, err := SkillFolderHash(t.Context(), canonicalRoot)
			if err != nil {
				t.Fatalf("hash canonical skill: %v", err)
			}
			distributedDigest, err := SkillFolderHash(t.Context(), distributedRoot)
			if err != nil {
				t.Fatalf("hash distributed skill: %v", err)
			}
			if distributedDigest != canonicalDigest {
				t.Fatalf(
					"canonical digest = %q, distributed digest = %q",
					canonicalDigest,
					distributedDigest,
				)
			}

			canonicalFiles := 0
			if err := filepath.WalkDir(
				canonicalRoot,
				func(path string, entry os.DirEntry, walkErr error) error {
					if walkErr != nil {
						return walkErr
					}
					if entry.IsDir() {
						return nil
					}
					relative, err := filepath.Rel(canonicalRoot, path)
					if err != nil {
						return err
					}
					canonical := readBaselineSkillContractFile(t, path)
					embeddedPath := filepath.ToSlash(filepath.Join(skillName, relative))
					if !bytes.Equal(canonical, embeddedByPath[embeddedPath]) {
						t.Errorf("%s canonical and embedded bytes differ", embeddedPath)
					}
					canonicalFiles++
					return nil
				},
			); err != nil {
				t.Fatalf("walk canonical skill: %v", err)
			}
			embeddedCount := 0
			for embeddedPath := range embeddedByPath {
				if strings.HasPrefix(embeddedPath, skillName+"/") {
					embeddedCount++
				}
			}
			if embeddedCount != canonicalFiles {
				t.Fatalf(
					"canonical files = %d, embedded files = %d",
					canonicalFiles,
					embeddedCount,
				)
			}
		})
	}

	setupPaths, err := filepath.Glob(
		filepath.Join(repoRoot, "internal", "baseline", "assets", "setups", "*.json"),
	)
	if err != nil {
		t.Fatalf("list canonical setup snapshots: %v", err)
	}
	if len(setupPaths) == 0 {
		t.Fatal("canonical setup snapshots are missing")
	}
	knownOwned := make(map[string]struct{}, len(Names()))
	for _, skillName := range Names() {
		knownOwned[skillName] = struct{}{}
	}
	for _, setupPath := range setupPaths {
		t.Run(filepath.Base(setupPath), func(t *testing.T) {
			var setup struct {
				Skills []struct {
					Name          string `json:"name"`
					ContentDigest string `json:"contentDigest"`
					Source        struct {
						Type string `json:"type"`
						Name string `json:"name"`
					} `json:"source"`
				} `json:"skills"`
			}
			if err := json.Unmarshal(readBaselineSkillContractFile(t, setupPath), &setup); err != nil {
				t.Fatalf("decode setup snapshot: %v", err)
			}
			for _, skill := range setup.Skills {
				if skill.Source.Type != "repo" {
					continue
				}
				if skill.Source.Name != "roundfix" {
					t.Errorf("%s repository source = %q, want roundfix", skill.Name, skill.Source.Name)
					continue
				}
				if _, ok := knownOwned[skill.Name]; !ok {
					t.Errorf("%s is repository-sourced but not owned by Roundfix", skill.Name)
					continue
				}
				want := baselineSnapshotSkillDigest(
					t,
					filepath.Join(repoRoot, ".agents", "skills", skill.Name),
				)
				if skill.ContentDigest != want {
					t.Errorf(
						"%s contentDigest = %q, want canonical %q",
						skill.Name,
						skill.ContentDigest,
						want,
					)
				}
			}
		})
	}

	lockBytes := readBaselineSkillContractFile(t, filepath.Join(repoRoot, "skills-lock.json"))
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	const wantUpstreamDigest = "d61310f840938a57edeae639ad44e5b0140b2bada06bac460ae59216e1a790e7"
	if got := upstreamManagedSkillDigest(t, repoRoot, lock.Skills); got != wantUpstreamDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, wantUpstreamDigest)
	}
}

func TestWritePRDProjectConstraints(t *testing.T) {
	testAuthoringProjectConstraints(t, "write-prd", "prd-template.md")
}

func TestWriteTechSpecProjectConstraints(t *testing.T) {
	testAuthoringProjectConstraints(t, "write-techspec", "techspec-template.md")
}

func TestProjectConstraintTaskGate(t *testing.T) {
	testWorkflowProjectConstraintContract(t, "write-tasks", []string{
		"Project Constraint preflight",
		"active, non-archived, and not already completed",
		"complete `Project Constraints` sections",
		"`Identifier strategy`",
		"`Authentication and HTTP`",
		"`Active ADR obligations`",
		"`Tooling authority`",
		"`docs/agents/` source path",
		"refuse to create or update `_tasks.md`",
		"express maintainer authorization",
		"exact bounded repository-relative files",
		"Dependencies remain owned only by `_tasks.md`",
		"Task status remains owned only by each Task file",
	})
}

func TestProjectConstraintImplementationGate(t *testing.T) {
	testWorkflowProjectConstraintContract(t, "implement-task", []string{
		"Project Constraint preflight",
		"before changing the Task status to `in_progress`",
		"Tooling authorization is not implied by Task assignment",
		"exact bounded repository-relative file list",
		"assigned Task file itself",
		"capture the pre-existing changed paths",
		"before every mutation",
		"changed-file postflight",
		"`git status --short` and `git diff --name-only`",
		"set `status: failed`",
		"Never edit `_tasks.md` or any other Task file",
	})
}

func TestProjectConstraintQAGate(t *testing.T) {
	testWorkflowProjectConstraintContract(t, "qa-gate", []string{
		"Project Constraint audit",
		"applicability",
		"source path under `docs/agents/`",
		"express maintainer authorization",
		"exact bounded files",
		"actual changed paths",
		"git diff-tree --no-commit-id --name-only -r",
		"missing authorization",
		"out-of-scope tooling changes",
		"Task status or Task Graph dependencies",
	})
}

func TestProjectConstraintJourney(t *testing.T) {
	fixtureRoot := filepath.Join("testdata", "project-constraint-journey")
	fixtures := []string{"_prd.md", "_techspec.md"}
	for _, fixture := range fixtures {
		t.Run(fixture, func(t *testing.T) {
			content := string(readBaselineSkillContractFile(
				t,
				filepath.Join(fixtureRoot, fixture),
			))
			if missing := missingProjectConstraintFixtureRows(content); len(missing) != 0 {
				t.Fatalf("Project Constraint fixture missing %v", missing)
			}
			for _, row := range []string{
				"Identifier strategy",
				"Authentication and HTTP",
				"Active ADR obligations",
				"Tooling authority",
			} {
				mutated := strings.Replace(content, "- "+row+":", "- Removed row:", 1)
				if missing := missingProjectConstraintFixtureRows(mutated); !slicesContain(missing, row) {
					t.Fatalf("removing %q did not fail the fixture contract: %v", row, missing)
				}
			}
		})
	}

	matrix := string(readBaselineSkillContractFile(
		t,
		filepath.Join(fixtureRoot, "qa-matrix.md"),
	))
	for _, required := range []string{
		"fresh disposable Fluxus greenfield clone",
		"separate fresh disposable Fluxus update clone",
		"Keep-defaults reuses the persisted Better Auth HTTP reason",
		"formatter and repository Verification",
		"audit and empty reapply",
		"final `qa-gate`",
	} {
		if !strings.Contains(matrix, required) {
			t.Errorf("final QA matrix missing %q", required)
		}
	}
	if rows := strings.Count(matrix, "| pending |"); rows != 2 {
		t.Errorf("final QA matrix pending rows = %d, want 2", rows)
	}
}

func TestToolingAuthorizationJourney(t *testing.T) {
	fixtureRoot := filepath.Join("testdata", "project-constraint-journey")
	prd := string(readBaselineSkillContractFile(t, filepath.Join(fixtureRoot, "_prd.md")))
	techspec := string(readBaselineSkillContractFile(t, filepath.Join(fixtureRoot, "_techspec.md")))
	prdPaths := projectConstraintAuthorizedPaths(prd)
	techspecPaths := projectConstraintAuthorizedPaths(techspec)
	wantPaths := []string{".golangci.yml", "scripts/verify.sh"}
	if !reflect.DeepEqual(prdPaths, wantPaths) ||
		!reflect.DeepEqual(techspecPaths, wantPaths) {
		t.Fatalf(
			"bounded tooling files differ: PRD=%v TechSpec=%v want=%v",
			prdPaths,
			techspecPaths,
			wantPaths,
		)
	}

	tests := []struct {
		name    string
		content string
		changed []string
		want    bool
	}{
		{
			name:    "exact bounded files are authorized",
			content: prd,
			changed: wantPaths,
			want:    true,
		},
		{
			name: "missing express authorization refuses decomposition",
			content: strings.Replace(
				prd,
				"expressly authorizes changes to exactly",
				"proposes changes to",
				1,
			),
			changed: wantPaths,
			want:    false,
		},
		{
			name:    "changed path outside bounds refuses settlement",
			content: prd,
			changed: []string{".golangci.yml", "scripts/release.sh"},
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := projectConstraintAllowsToolingChanges(test.content, test.changed); got != test.want {
				t.Fatalf("tooling authorization = %v, want %v", got, test.want)
			}
		})
	}

	t.Run("decomposition contract", func(t *testing.T) {
		testWorkflowProjectConstraintContract(t, "write-tasks", []string{
			"refuse to create or update `_tasks.md`",
			"express maintainer authorization",
			"exact bounded repository-relative files",
		})
	})
	t.Run("execution contract", func(t *testing.T) {
		testWorkflowProjectConstraintContract(t, "implement-task", []string{
			"mutation allowlist",
			"before every mutation",
			"Every newly changed path",
			"set `status: failed`",
		})
	})
	t.Run("QA contract", func(t *testing.T) {
		testWorkflowProjectConstraintContract(t, "qa-gate", []string{
			"exact bounded files",
			"actual changed paths",
			"missing authorization",
			"out-of-scope tooling changes",
		})
	})
}

func TestLegacySpecConstraintExemption(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	for _, skillName := range []string{"write-tasks", "implement-task", "qa-gate"} {
		t.Run(skillName, func(t *testing.T) {
			testWorkflowProjectConstraintContract(t, skillName, []string{
				"completed or archived legacy Spec",
				"byte-identical",
			})
		})
	}

	module := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, "internal", "baseline", "assets", "modules", "spec-workflow.json"),
	))
	golden := string(readBaselineSkillContractFile(
		t,
		filepath.Join(
			repoRoot,
			"internal",
			"baseline",
			"assets",
			"formatter-fixtures",
			"standard-typescript-monorepo",
			"golden",
			"docs",
			"agents",
			"spec-routing.md",
		),
	))
	for _, required := range []string{
		"Project Constraints",
		"express maintainer authorization",
		"bounded repository-relative files",
		"actual changed-file scope",
		"completed or archived legacy Specs byte-identical",
		"Dependencies remain owned only by the Task Graph",
		"status remains owned only by each Task file",
	} {
		if !strings.Contains(module, required) {
			t.Errorf("spec-workflow module missing %q", required)
		}
		if !strings.Contains(golden, required) {
			t.Errorf("generated spec-routing fixture missing %q", required)
		}
	}
}

func TestAuthoringConstraintOwnership(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	lockBytes := readBaselineSkillContractFile(t, filepath.Join(repoRoot, "skills-lock.json"))
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}

	for _, name := range []string{"write-prd", "write-techspec"} {
		t.Run(name, func(t *testing.T) {
			if _, managedUpstream := lock.Skills[name]; managedUpstream {
				t.Fatalf("%s must remain repository-owned and absent from skills-lock.json", name)
			}
			assertSkillTreesEqual(
				t,
				filepath.Join(repoRoot, ".agents", "skills", name),
				filepath.Join(repoRoot, "skills", name),
			)
		})
	}

	const wantUpstreamDigest = "d61310f840938a57edeae639ad44e5b0140b2bada06bac460ae59216e1a790e7"
	if got := upstreamManagedSkillDigest(t, repoRoot, lock.Skills); got != wantUpstreamDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, wantUpstreamDigest)
	}
}

func TestUpstreamADRFormatUnchanged(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	lockBytes := readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, "skills-lock.json"),
	)
	var lock struct {
		Version int `json:"version"`
		Skills  map[string]struct {
			ComputedHash string `json:"computedHash"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(lockBytes, &lock); err != nil {
		t.Fatalf("decode skills lock: %v", err)
	}
	if lock.Version != 1 || len(lock.Skills) == 0 {
		t.Fatalf("unexpected skills lock contract: version=%d skills=%d", lock.Version, len(lock.Skills))
	}

	names := make([]string, 0, len(lock.Skills))
	for name := range lock.Skills {
		names = append(names, name)
	}
	sort.Strings(names)
	upstreamDigest := sha256.New()
	for _, name := range names {
		folderDigest, err := SkillFolderHash(t.Context(), filepath.Join(repoRoot, ".agents", "skills", name))
		if err != nil {
			t.Fatalf("hash upstream-managed skill %q: %v", name, err)
		}
		_, _ = upstreamDigest.Write([]byte(name))
		_, _ = upstreamDigest.Write([]byte(folderDigest))
	}
	const wantUpstreamDigest = "d61310f840938a57edeae639ad44e5b0140b2bada06bac460ae59216e1a790e7"
	if got := hex.EncodeToString(upstreamDigest.Sum(nil)); got != wantUpstreamDigest {
		t.Fatalf("upstream-managed skill tree digest = %q, want %q", got, wantUpstreamDigest)
	}

	const wantADRFormatSHA256 = "f1f36cd3f8d3b6474ddd5855da4e233bfc4ae1a1c5024909ccf11871819a41b2"
	adrFormat := readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", "domain-modeling", "ADR-FORMAT.md"),
	)
	sum := sha256.Sum256(adrFormat)
	if got := hex.EncodeToString(sum[:]); got != wantADRFormatSHA256 {
		t.Fatalf("domain-modeling/ADR-FORMAT.md digest = %q, want %q", got, wantADRFormatSHA256)
	}
}

func testAuthoringProjectConstraints(t *testing.T, skillName, templateName string) {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join(".."))
	skill := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", skillName, "SKILL.md"),
	))
	template := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", skillName, "references", templateName),
	))

	for _, required := range []string{
		"Project Constraints",
		"docs/agents/agent-instructions.md",
		"docs/agents/domain.md",
		"docs/agents/backend.md",
		"docs/agents/spec-routing.md",
		"applicable or not applicable with a reason",
		"express maintainer authorization",
		"bounded files",
		"MUST NOT report completion",
	} {
		if !strings.Contains(skill, required) {
			t.Errorf("%s skill missing %q", skillName, required)
		}
	}

	for _, area := range []string{
		"Identifier strategy",
		"Authentication and HTTP",
		"Active ADR obligations",
		"Tooling authority",
	} {
		if !strings.Contains(template, area) {
			t.Errorf("%s template missing %q", skillName, area)
		}
	}
	if count := strings.Count(template, "<applicable | not applicable>"); count != 4 {
		t.Errorf("%s template applicability placeholders = %d, want 4", skillName, count)
	}
	if count := strings.Count(template, "Source: `docs/agents/"); count != 4 {
		t.Errorf("%s template docs/agents source placeholders = %d, want 4", skillName, count)
	}
	for _, required := range []string{
		"express maintainer authorization",
		"bounded files",
		"no protected tooling mutation",
	} {
		if !strings.Contains(template, required) {
			t.Errorf("%s template missing tooling authorization guidance %q", skillName, required)
		}
	}

	frontmatter := authoringTemplateFrontmatter(t, template)
	if strings.Contains(strings.ToLower(frontmatter), "authoriz") {
		t.Errorf("%s template adds authorization frontmatter: %q", skillName, frontmatter)
	}
}

func testWorkflowProjectConstraintContract(t *testing.T, skillName string, required []string) {
	t.Helper()
	repoRoot := filepath.Clean(filepath.Join(".."))
	canonicalPath := filepath.Join(repoRoot, ".agents", "skills", skillName, "SKILL.md")
	distributedPath := filepath.Join(repoRoot, "skills", skillName, "SKILL.md")
	canonical := readBaselineSkillContractFile(t, canonicalPath)
	distributed := readBaselineSkillContractFile(t, distributedPath)
	if !bytes.Equal(canonical, distributed) {
		t.Fatalf("%s canonical and distributed guidance differ; run make skills-sync", skillName)
	}

	for _, missing := range missingWorkflowContractClauses(string(canonical), required) {
		t.Errorf("%s skill missing %q", skillName, missing)
	}
	normalized := normalizeWorkflowContract(string(canonical))
	for _, clause := range required {
		mutated := strings.ReplaceAll(normalized, normalizeWorkflowContract(clause), "")
		if missing := missingWorkflowContractClauses(mutated, required); !slicesContain(missing, clause) {
			t.Errorf("%s contract still accepts removal of %q; missing = %v", skillName, clause, missing)
		}
	}
}

func missingWorkflowContractClauses(content string, required []string) []string {
	content = normalizeWorkflowContract(content)
	var missing []string
	for _, clause := range required {
		if !strings.Contains(content, normalizeWorkflowContract(clause)) {
			missing = append(missing, clause)
		}
	}
	return missing
}

func normalizeWorkflowContract(content string) string {
	return strings.Join(strings.Fields(content), " ")
}

func slicesContain(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func missingProjectConstraintFixtureRows(content string) []string {
	required := []struct {
		row    string
		source string
	}{
		{row: "Identifier strategy", source: "docs/agents/domain.md"},
		{row: "Authentication and HTTP", source: "docs/agents/backend.md"},
		{row: "Active ADR obligations", source: "docs/agents/domain.md"},
		{row: "Tooling authority", source: "docs/agents/agent-instructions.md"},
	}
	var missing []string
	for _, requirement := range required {
		bullet := projectConstraintFixtureBullet(content, requirement.row)
		if bullet == "" ||
			!strings.Contains(bullet, "applicable") ||
			!strings.Contains(bullet, "Source: `"+requirement.source+"`") {
			missing = append(missing, requirement.row)
		}
	}
	return missing
}

func projectConstraintFixtureBullet(content, row string) string {
	sectionStart := strings.Index(content, "## Project Constraints")
	if sectionStart == -1 {
		return ""
	}
	section := content[sectionStart:]
	if next := strings.Index(section[len("## Project Constraints"):], "\n## "); next != -1 {
		section = section[:len("## Project Constraints")+next]
	}
	marker := "- " + row + ":"
	start := strings.Index(section, marker)
	if start == -1 {
		return ""
	}
	bullet := section[start:]
	if next := strings.Index(bullet[len(marker):], "\n- "); next != -1 {
		bullet = bullet[:len(marker)+next]
	}
	return strings.Join(strings.Fields(bullet), " ")
}

func projectConstraintAuthorizedPaths(content string) []string {
	bullet := projectConstraintFixtureBullet(content, "Tooling authority")
	if !strings.Contains(bullet, "expressly authorizes changes to exactly") {
		return nil
	}
	var paths []string
	for {
		start := strings.Index(bullet, "`")
		if start == -1 {
			break
		}
		bullet = bullet[start+1:]
		end := strings.Index(bullet, "`")
		if end == -1 {
			return nil
		}
		value := bullet[:end]
		bullet = bullet[end+1:]
		if !strings.HasPrefix(value, "docs/agents/") {
			paths = append(paths, value)
		}
	}
	return paths
}

func projectConstraintAllowsToolingChanges(content string, changed []string) bool {
	authorized := projectConstraintAuthorizedPaths(content)
	if len(authorized) == 0 {
		return false
	}
	allowlist := make(map[string]struct{}, len(authorized))
	for _, path := range authorized {
		allowlist[path] = struct{}{}
	}
	for _, path := range changed {
		if _, ok := allowlist[path]; !ok {
			return false
		}
	}
	return true
}

func authoringTemplateFrontmatter(t *testing.T, template string) string {
	t.Helper()
	const fence = "---"
	start := strings.Index(template, fence)
	if start == -1 {
		t.Fatal("template is missing frontmatter")
	}
	rest := template[start+len(fence):]
	end := strings.Index(rest, fence)
	if end == -1 {
		t.Fatal("template has unterminated frontmatter")
	}
	return rest[:end]
}

func assertSkillTreesEqual(t *testing.T, canonicalRoot, distributedRoot string) {
	t.Helper()
	if err := filepath.WalkDir(canonicalRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(canonicalRoot, path)
		if err != nil {
			return err
		}
		canonical := readBaselineSkillContractFile(t, path)
		distributed := readBaselineSkillContractFile(t, filepath.Join(distributedRoot, relative))
		if !bytes.Equal(canonical, distributed) {
			t.Errorf("%s canonical and distributed bytes differ", filepath.ToSlash(relative))
		}
		return nil
	}); err != nil {
		t.Fatalf("walk canonical skill tree %s: %v", canonicalRoot, err)
	}
}

func upstreamManagedSkillDigest(
	t *testing.T,
	repoRoot string,
	managed map[string]json.RawMessage,
) string {
	t.Helper()
	names := make([]string, 0, len(managed))
	for name := range managed {
		names = append(names, name)
	}
	sort.Strings(names)
	digest := sha256.New()
	for _, name := range names {
		folderDigest, err := SkillFolderHash(t.Context(), filepath.Join(repoRoot, ".agents", "skills", name))
		if err != nil {
			t.Fatalf("hash upstream-managed skill %q: %v", name, err)
		}
		_, _ = digest.Write([]byte(name))
		_, _ = digest.Write([]byte(folderDigest))
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func baselineSnapshotSkillDigest(t *testing.T, root string) string {
	t.Helper()
	type digestFile struct {
		path string
		data []byte
	}
	var files []digestFile
	if err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, digestFile{path: filepath.ToSlash(relative), data: data})
		return nil
	}); err != nil {
		t.Fatalf("hash Baseline snapshot skill %q: %v", root, err)
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})
	digest := sha256.New()
	for _, file := range files {
		_, _ = digest.Write([]byte(file.path))
		_, _ = digest.Write(file.data)
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func readBaselineSkillContractFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func baselineSkillBody(content string) string {
	const delimiter = "---"
	parts := strings.SplitN(content, delimiter, 3)
	if len(parts) != 3 {
		return content
	}
	return parts[2]
}
