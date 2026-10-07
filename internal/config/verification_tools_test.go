package config

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestVerificationToolsAreBareExecutableNames(t *testing.T) {
	for _, entry := range []string{`[go, go]`, `["go test"]`, `[./go]`, `[/bin/go]`, `[" go"]`, `["-go"]`, `[""]`, `["go;bad"]`, `go`, `[1]`} {
		t.Run(entry, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte("verification:\n  tools: "+entry+"\n"))
			if err == nil || !strings.Contains(err.Error(), "verification.tools") {
				t.Fatalf("invalid tools error = %v", err)
			}
		})
	}
	user := []byte("verification:\n  tools: [go, tool_1.2+-]\n")
	for _, tc := range []struct {
		name, project string
		want          []string
	}{
		{"user", "", []string{"go", "tool_1.2+-"}},
		{"replace", "verification:\n  tools: [git]\n", []string{"git"}},
		{"empty replaces", "verification:\n  tools: []\n", []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := ResolveConfigProposal(user, []byte(tc.project))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg.Verification.Tools, tc.want) {
				t.Fatalf("tools = %v, want %v", cfg.Verification.Tools, tc.want)
			}
		})
	}
	if len(Builtin().Verification.Tools) != 0 {
		t.Fatal("default tools is not empty")
	}
	if !strings.Contains(DefaultConfigYAML(), "\n  tools: []\n") {
		t.Fatal("generated User Config must expose the default empty tools list")
	}
}

func TestThisRepositoryDeclaresItsToolsAndDerivedPaths(t *testing.T) {
	data, err := os.ReadFile("../../.roundfixrc.yml")
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := ResolveConfigProposal(nil, data)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Verification.Tools, []string{"go", "gofmt", "git", "make"}) {
		t.Fatalf("tools = %v", cfg.Verification.Tools)
	}
	want := []DerivedPathDeclaration{
		{Paths: []string{"internal/baseline/module-versions.json", "internal/baseline/testdata/catalog.digest", "internal/baseline/testdata/catalog.normalized.json", "internal/baseline/testdata/catalog.diagnostics.golden.json", "internal/baseline/testdata/plan-characterization/*.golden.json"}, Lines: &DerivedLineDeclaration{Paths: []string{"internal/baseline/assets/modules/*.json"}, Match: `^  "version": [0-9]+,$`}, Regenerate: "go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 && make baseline-digests"},
		{Paths: []string{"internal/baseline/assets/profiles/standard-typescript-monorepo.json"}, Lines: &DerivedLineDeclaration{Paths: []string{"internal/baseline/assets/profiles/standard-typescript-monorepo.json"}, Match: `^    "goldenDigest": "[0-9a-f]+",$`}, Regenerate: "go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 && make baseline-digests"},
		{Paths: []string{"internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/"}, Regenerate: "go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 && make baseline-digests"},
		{Paths: []string{"docs/agents/setup-context.json", "docs/agents/"}, Regenerate: "go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text"},
		{Paths: []string{"skills/testdata/owned-skill-versions.json"}, Lines: &DerivedLineDeclaration{Paths: []string{".agents/skills/*/SKILL.md", "skills/*/SKILL.md"}, Match: "^ *version: "}, Regenerate: "go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions"},
		{Paths: []string{"docs/references/coverage-record.json"}, Regenerate: "go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1"},
	}
	if !reflect.DeepEqual(cfg.Delivery.DerivedPaths, want) {
		t.Fatalf("derived declarations = %#v", cfg.Delivery.DerivedPaths)
	}
}
