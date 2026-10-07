// Suite: Repository Contract Test declarations and selective gate wiring.
// Invariant: real contracts declare their relevance and the gate runs selected invocations in order.
// Boundary IN: repository test headers and the Makefile with temporary command stubs.
// Boundary OUT: contract test bodies and the full repository Verification.
package verifyselect_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/config"
	"roundfix/internal/verifyselect"
)

func TestRepositoryContractTestsDeclareTheirRelevance(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	contracts, err := verifyselect.DiscoverContracts(root)
	if err != nil {
		t.Fatal(err)
	}
	regeneration := []string{"internal/baseline/assets/", "internal/baseline/testdata/", ".agents/skills/", "skills/"}
	declaredRegeneration := []string{".roundfixrc.yml", "internal/baseline/", ".agents/skills/", "skills/", "docs/agents/", "docs/references/coverage-record.json", "internal/spec/", "internal/cli/baseline_*", "cmd/roundfix/", "internal/suiteguard/", "internal/suiteguardcontract/"}
	relevant := map[string][]string{
		"TestRegenerationIsDeclared":                            declaredRegeneration,
		"TestMeasuredSanctionedOwnershipMatchesRecords":         regeneration,
		"TestDeclaredStepRegenerationAndFrozenBoundaries":       regeneration,
		"TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical": regeneration,
		"TestRepositoryGateRunsTheAnalyzer":                     {"internal/baseline/analyzer/"},
	}
	governed := map[string]bool{
		"TestGovernedSetCoversOwnedShippedTemplates": true,
		"TestGovernedSetOnlyGrows":                   true,
		"TestCleanupHistoricalGrantEvidence":         true,
		"TestEveryBoundedPathIsGoverned":             true,
	}
	seen := make(map[string]bool)
	foundRelevant := make(map[string]bool)
	foundGoverned := make(map[string]bool)
	var foundInstallation, foundPackage bool
	docsFiles := make(map[string]bool)
	for _, contract := range contracts {
		seen[contract.Name] = true
		want := verifyselect.ContractPackage
		switch {
		case contract.Package == "internal/docscontract" && contract.Tag == "docscontract":
			want = verifyselect.ContractAlways
			docsFiles[contract.File] = true
		case contract.Package == "internal/suiteguard" && contract.Name == "TestEverySpawningPackageInstallsTheSuiteGuard":
			want = verifyselect.ContractAlways
			foundInstallation = true
		case contract.Package == "internal/speccheck" && governed[contract.Name]:
			want = verifyselect.ContractAlways
			foundGoverned[contract.Name] = true
		case relevant[contract.Name] != nil:
			want = verifyselect.ContractRelevant
			foundRelevant[contract.Name] = true
			if !reflect.DeepEqual(contract.Paths, relevant[contract.Name]) {
				t.Errorf("%s paths = %v, want %v", contract.Name, contract.Paths, relevant[contract.Name])
			}

		}
		if contract.Class == verifyselect.ContractBoundary {
			t.Errorf("%s/%s is boundary", contract.Package, contract.Name)
		}
		if contract.Class != want {
			t.Errorf("%s/%s class = %s, want %s", contract.Package, contract.Name, contract.Class, want)
		}
		if contract.Package == "internal/config" && contract.Name == "TestEverySpawningPackageInstallsTheSuiteGuard" {
			foundPackage = true
		}
	}
	if !foundInstallation || !foundPackage || len(docsFiles) != 8 || len(foundRelevant) != len(relevant) || len(foundGoverned) != len(governed) {
		t.Fatalf("missing declarations: installation=%t package=%t docs files=%d relevant=%v governed=%v", foundInstallation, foundPackage, len(docsFiles), foundRelevant, foundGoverned)
	}
	contents, err := os.ReadFile(filepath.Join(root, "Makefile"))
	if err != nil {
		t.Fatal(err)
	}
	assignment := regexp.MustCompile(`(?m)^REPO_CONTRACT_TESTS := (.+)$`).FindStringSubmatch(string(contents))
	if len(assignment) != 2 {
		t.Fatal("REPO_CONTRACT_TESTS assignment missing")
	}
	for _, name := range strings.Split(assignment[1], "|") {
		if !seen[name] {
			t.Errorf("Makefile contract %s was not discovered", name)
		}
	}

	for _, scenario := range []struct {
		name, path string
		extra      int
	}{
		{"user guide selects only always", "docs/user-guide/example.md", 0},
		{"owned skill selects regeneration", ".agents/skills/roundfix/SKILL.md", 4},
		{"project config selects declared regeneration", ".roundfixrc.yml", 1},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			selected := verifyselect.SelectContractPaths(contracts, []string{scenario.path})
			always, extra := 0, 0
			for _, contract := range contracts {
				if contract.Class == verifyselect.ContractAlways {
					always++
				}
			}
			for _, contract := range selected {
				if contract.Class == verifyselect.ContractAlways {
					continue
				}
				if contract.Class != verifyselect.ContractRelevant || !reflect.DeepEqual(contract.Paths, relevant[contract.Name]) ||
					(contract.Name != "TestRegenerationIsDeclared" && !reflect.DeepEqual(contract.Paths, regeneration)) ||
					(scenario.path == ".roundfixrc.yml" && contract.Name != "TestRegenerationIsDeclared") {
					t.Errorf("unexpected selected contract: %+v", contract)
				}
				extra++
			}
			if extra != scenario.extra || len(selected) != always+scenario.extra {
				t.Errorf("selection has %d contracts, %d by change; want %d, %d", len(selected), extra, always+scenario.extra, scenario.extra)
			}
		})
	}

}

func TestRegenerationContractRunsWhenADerivedInputChanges(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	contracts, err := verifyselect.DiscoverContracts(root)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(filepath.Join(root, ".roundfixrc.yml"))
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := config.ResolveConfigProposal(nil, contents)
	if err != nil {
		t.Fatal(err)
	}
	declarations := proposal.Delivery.DerivedPaths
	if len(declarations) < 4 {
		t.Fatalf(".roundfixrc.yml derived_paths has %d declarations, want at least four", len(declarations))
	}
	assertSelected := func(t *testing.T, entry, sample string, want bool) {
		t.Helper()
		selected := false
		for _, contract := range verifyselect.SelectContractPaths(contracts, []string{sample}) {
			if contract.Name == "TestRegenerationIsDeclared" {
				selected = true
			}
		}
		if selected != want {
			t.Errorf("entry %q (sample %q) selects TestRegenerationIsDeclared = %t, want %t", entry, sample, selected, want)
		}
	}
	for i, declaration := range declarations {
		scopes := []struct {
			name  string
			paths []string
		}{{"paths", declaration.Paths}}
		if declaration.Lines != nil {
			scopes = append(scopes, struct {
				name  string
				paths []string
			}{"lines.paths", declaration.Lines.Paths})
		}
		for _, scope := range scopes {
			for _, entry := range scope.paths {
				t.Run(fmt.Sprintf("declaration_%d/%s/%s", i, scope.name, entry), func(t *testing.T) {
					sample := strings.ReplaceAll(entry, "*", "sample")
					if strings.HasSuffix(entry, "/") {
						sample += "sample.json"
					}
					assertSelected(t, entry, sample, true)
				})
			}
		}
	}
	for _, entry := range []string{
		".roundfixrc.yml",
		"internal/baseline/plan.go",
		"skills/skills.go",
		"internal/spec/archive.go",
		"internal/cli/baseline_update.go",
		"cmd/roundfix/main.go",
		"internal/suiteguard/suiteguard.go",
		"internal/suiteguardcontract/regeneration.go",
	} {
		t.Run(entry, func(t *testing.T) {
			assertSelected(t, entry, entry, true)
		})
	}
	for _, entry := range []string{"docs/user-guide/example.md", "internal/daemon/daemon.go"} {
		t.Run(entry, func(t *testing.T) {
			assertSelected(t, entry, entry, false)
		})
	}
}

func TestVerifyChangedRunsTheSelectedContracts(t *testing.T) {
	t.Parallel()
	makefile := filepath.Join(findRepositoryRoot(t), "Makefile")
	contents, err := os.ReadFile(makefile)
	if err != nil {
		t.Fatal(err)
	}
	recipe := regexp.MustCompile(`(?ms)^verify-changed:.*?\n(.*?)(?:\n\n|\z)`).FindStringSubmatch(string(contents))
	if len(recipe) != 2 {
		t.Fatal("verify-changed recipe missing")
	}
	loopEnd := strings.LastIndex(recipe[1], "done")
	invocation := strings.Index(recipe[1], "$(MAKE) --no-print-directory verify-changed-contracts")
	if loopEnd < 0 || invocation <= loopEnd {
		t.Fatal("verify-changed must invoke verify-changed-contracts after its set loop")
	}
	const lines = "docscontract ^(TestDocs)$ ./internal/docscontract\nrepocontract ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
	const first = "test -count=1 -tags docscontract -run ^(TestDocs)$ ./internal/docscontract\n"
	const second = "test -count=1 -tags repocontract -run ^(TestOne|TestTwo)$ ./internal/baseline ./skills\n"
	for _, scenario := range []struct {
		name, output                      string
		selectorFail, goFail, wantFailure bool
		want                              string
	}{
		{name: "ordered invocations", output: lines, want: first + second},
		{name: "first test failure stops execution", output: lines, goFail: true, wantFailure: true, want: first},
		{name: "selector failure prevents execution", output: lines, selectorFail: true, wantFailure: true},
		{name: "empty selection runs no tests"},
		{name: "empty package list runs no tests", output: "docscontract ^(TestDocs)$\n"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "selection", scenario.output)
			writeFile(t, dir, "calls", "")
			selector := "#!/bin/sh\n[ \"$*\" = '-contracts -base fixture-base' ] || exit 92\ncat selection\n"
			if scenario.selectorFail {
				selector += "exit 17\n"
			}
			goStub := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> calls\n"
			if scenario.goFail {
				goStub += "exit 19\n"
			}
			writeFile(t, dir, "selector", selector)
			writeFile(t, dir, "go", goStub)
			for _, name := range []string{"selector", "go"} {
				if err := os.Chmod(filepath.Join(dir, name), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "verify-changed-contracts", "VERIFY_BASE=fixture-base", "VERIFY_SELECT="+filepath.Join(dir, "selector"), "GO="+filepath.Join(dir, "go"))
			command.Dir = dir
			output, err := command.CombinedOutput()
			if (err != nil) != scenario.wantFailure {
				t.Fatalf("make error = %v, want failure %t\n%s", err, scenario.wantFailure, output)
			}
			calls, err := os.ReadFile(filepath.Join(dir, "calls"))
			if err != nil {
				t.Fatal(err)
			}
			if string(calls) != scenario.want {
				t.Fatalf("go invocations = %q, want %q\n%s", calls, scenario.want, output)
			}
		})
	}
	for _, fail := range []bool{false, true} {
		name := "empty sets still run contracts"
		if fail {
			name = "contract failure fails verify-changed with empty sets"
		}
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, "calls", "")
			writeFile(t, dir, "selector", "#!/bin/sh\ncase \"$1\" in\n-base) exit 0 ;;\n-contracts) printf '%s\\n' 'docscontract ^(TestDocs)$ ./internal/docscontract' ;;\n*) exit 92 ;;\nesac\n")
			goStub := "#!/bin/sh\nprintf '%s\\n' \"$*\" >> calls\n"
			if fail {
				goStub += "exit 19\n"
			}
			writeFile(t, dir, "go", goStub)
			for _, script := range []string{"selector", "go"} {
				if err := os.Chmod(filepath.Join(dir, script), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			// Override only the prerequisites so this exercises the real parent
			// recipe without invoking host formatting, analysis, or builds.
			writeFile(t, dir, "prerequisites.mk", "fmt-check: ;\nvet: ;\nbuild: ;\n")
			command := exec.CommandContext(t.Context(), "make", "--no-print-directory", "-f", makefile, "-f", filepath.Join(dir, "prerequisites.mk"), "verify-changed", "VERIFY_SELECT="+filepath.Join(dir, "selector"), "GO="+filepath.Join(dir, "go"), "MAKE=make --no-print-directory -f "+makefile)
			command.Dir = dir
			output, err := command.CombinedOutput()
			if (err != nil) != fail {
				t.Fatalf("make verify-changed error = %v, want failure %t\n%s", err, fail, output)
			}
			calls, err := os.ReadFile(filepath.Join(dir, "calls"))
			if err != nil {
				t.Fatal(err)
			}
			if string(calls) != first {
				t.Fatalf("go invocations = %q, want %q\n%s", calls, first, output)
			}
		})
	}
}
