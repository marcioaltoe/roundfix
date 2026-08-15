// Suite: governed repository paths
// Invariant: only repository paths bound by a named tooling kind are governed
// Boundary IN: the public governed-path predicate and repository-relative paths
// Boundary OUT: authorization history and changed-path audit integration
package speccheck_test

import (
	"testing"

	"roundfix/internal/speccheck"
)

func TestGovernedPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "linter configuration", path: ".golangci.yml", want: true},
		{name: "formatter configuration", path: ".prettierrc", want: true},
		{name: "typechecker configuration", path: "tsconfig.json", want: true},
		{name: "test-runner configuration", path: "vitest.config.ts", want: true},
		{name: "architecture-checker configuration", path: ".dependency-cruiser.js", want: true},
		{name: "build-tool configuration", path: "Makefile", want: true},
		{name: "package-manager configuration", path: "go.mod", want: true},
		{name: "code-generator script", path: "scripts/generate.sh", want: true},
		{name: "ignore file", path: ".gitignore", want: true},
		{name: "plugin declaration", path: ".codex-plugin/plugin.json", want: true},
		{name: "version pin", path: ".tool-versions", want: true},
		{name: "ordinary Go source", path: "internal/app/metadata.go", want: false},
		{name: "ordinary test", path: "internal/app/metadata_test.go", want: false},
		{name: "Spec document", path: "docs/specs/example/_prd.md", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := speccheck.GovernedPath(tt.path); got != tt.want {
				t.Fatalf("GovernedPath(%q) = %t, want %t", tt.path, got, tt.want)
			}
		})
	}
}
