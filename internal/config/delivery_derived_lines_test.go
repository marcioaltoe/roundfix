// Boundary: config proposal parsing and line-scoped declaration validation.
package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestDerivedLineDeclarationsAreReadAndValidated(t *testing.T) {
	const prefix = "delivery:\n  derived_paths:\n    - paths: [record.json]\n      regenerate: record\n      lines: "
	cfg, err := ResolveConfigProposal(nil, []byte(prefix+"{paths: ['skills/*/SKILL.md', cache/], match: '^ *version: '}\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := &DerivedLineDeclaration{Paths: []string{"skills/*/SKILL.md", "cache/"}, Match: "^ *version: "}
	if !reflect.DeepEqual(cfg.Delivery.DerivedPaths[0].Lines, want) {
		t.Fatalf("lines=%+v", cfg.Delivery.DerivedPaths[0].Lines)
	}
	for name, expected := range map[string]bool{"skills/example/SKILL.md": true, "cache/nested/file": true, "skills/a/b/SKILL.md": false, "record.json": false} {
		if cfg.Delivery.DerivedPaths[0].MatchesLines(name) != expected {
			t.Errorf("match %q", name)
		}
	}
	for _, entry := range []string{
		"null", "[]", "version", "{paths: skill, match: version}", "{paths: [skill], match: []}",
		"{}", "{paths: [skill]}", "{match: version}", "{paths: [], match: version}",
		"{paths: [skill], match: ''}", "{paths: [skill], match: ' '}", "{paths: [skill], match: '['}",
		"{paths: [/tmp/file], match: version}", "{paths: [../file], match: version}",
		"{paths: [a/../file], match: version}", "{paths: [a/./file], match: version}",
		"{paths: [''], match: version}", "{paths: ['a\\b'], match: version}", "{paths: ['['], match: version}",
	} {
		t.Run(entry, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte(prefix+entry+"\n"))
			if err == nil || !strings.Contains(err.Error(), "delivery.derived_paths[0].lines") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
