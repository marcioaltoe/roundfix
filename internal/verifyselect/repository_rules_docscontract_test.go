//go:build docscontract

// Suite: repository record contracts.
// Invariant: records and repository rules remain current.
// Boundary IN: live repository records and source declarations.
// Boundary OUT: product execution and publication.
//
//verify:always
package verifyselect_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"roundfix/internal/skillcoverage"
	"roundfix/internal/verifyselect"
	"strings"
	"testing"
)

func TestRepositoryContractTestsDeclareTheirRelevance(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	contracts, err := verifyselect.DiscoverContracts(root)
	if err != nil {
		t.Fatal(err)
	}
	regeneration := []string{"internal/baseline/assets/", "internal/baseline/testdata/", ".agents/skills/", "skills/"}
	declaredRegeneration := []string{".roundfixrc.yml", "internal/baseline/", ".agents/skills/", "skills/", "docs/agents/", "docs/references/coverage-record.json", "internal/spec/", "internal/cli/baseline_*", "cmd/roundfix/", "internal/suiteguard/", "internal/suiteguardcontract/", skillcoverage.RecordPath, skillcoverage.MapPath}
	relevant := map[string][]string{
		"TestCoverageCollectionMatchesTheToolchainOnThisPlatform": {"dist/npm/platforms.json", "go.mod", "go.sum"},
		"TestJournalConsumerCorpusReplaysEveryConsumer":           {"internal/store/testdata/", "internal/runevent/", "internal/tui/", "internal/cli/attach.go", "internal/cli/events.go", "internal/cli/gc*.go", "internal/cli/reconcile.go"},
		"TestRegenerationIsDeclared":                              declaredRegeneration,
		"TestMeasuredSanctionedOwnershipMatchesRecords":           regeneration,
		"TestDeclaredStepRegenerationAndFrozenBoundaries":         regeneration,
		"TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical":   regeneration,
		"TestRepositoryGateRunsTheAnalyzer":                       {"internal/baseline/analyzer/"},
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
	liveDocsFiles := make(map[string]bool)
	for _, contract := range contracts {
		seen[contract.Name] = true
		want := verifyselect.ContractPackage
		switch {
		case (contract.Package == "internal/testfixture" || contract.Package == "internal/verifyselect") && contract.Tag == "docscontract":
			want = verifyselect.ContractAlways
		case contract.Package == "internal/docscontract" && contract.Tag == "docscontract":
			want = verifyselect.ContractAlways
			docsFiles[contract.File] = true
		case (contract.Package == "internal/cli" || contract.Package == "internal/baseline") && contract.Tag == "docscontract":
			// ADR-0326: tests reading the live documentation trees are always selected.
			want = verifyselect.ContractAlways
			liveDocsFiles[contract.File] = true
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
	if len(liveDocsFiles) != 4 {
		t.Fatalf("Live documentation contract files = %d, want 4: %v", len(liveDocsFiles), liveDocsFiles)
	}
	if !foundInstallation || !foundPackage || len(docsFiles) != 17 || len(foundRelevant) != len(relevant) || len(foundGoverned) != len(governed) {
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

func shardCoverageError(names []string, patterns []*regexp.Regexp, skips []bool) error {
	for _, name := range names {
		count := 0
		for i, p := range patterns {
			match := p.MatchString(name)
			if skips[i] {
				match = !match
			}
			if match {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("%q is selected by %d shards, want exactly 1", name, count)
		}
	}
	return nil
}
func TestCLIShardsCoverEveryTestExactlyOnce(t *testing.T) {
	t.Parallel()
	root := findRepositoryRoot(t)
	names, err := verifyselect.CLITestNames(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) == 0 {
		t.Fatal("no CLI tests")
	}
	names = append(names, "TestNotListedByTheParser")
	for _, of := range []int{1, 2, 3, 7} {
		var patterns []*regexp.Regexp
		var skips []bool
		for shard := 1; shard <= of; shard++ {
			flag, pattern, e := verifyselect.CLIShard(root, shard, of)
			if e != nil {
				t.Fatal(e)
			}
			patterns = append(patterns, regexp.MustCompile(pattern))
			skips = append(skips, flag == "-skip")
		}
		if err := shardCoverageError(names, patterns, skips); err != nil {
			t.Fatal(err)
		}
	}
	for _, c := range []struct {
		name     string
		patterns []*regexp.Regexp
		skips    []bool
		count    string
	}{{"missing", []*regexp.Regexp{regexp.MustCompile("^$")}, []bool{false}, "0 shards"}, {"duplicate", []*regexp.Regexp{regexp.MustCompile("^TestSeed$"), regexp.MustCompile("^TestSeed$")}, []bool{false, false}, "2 shards"}} {
		t.Run(c.name, func(t *testing.T) {
			err := shardCoverageError([]string{"TestSeed"}, c.patterns, c.skips)
			if err == nil || !strings.Contains(err.Error(), c.count) || !strings.Contains(err.Error(), "TestSeed") {
				t.Fatalf("seed not detected: %v", err)
			}
		})
	}
}
