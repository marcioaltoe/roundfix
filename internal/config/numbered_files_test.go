// Boundary: numbered-file config parsing, precedence, validation, and matching.
package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestNumberedFileDeclarationsAreReadAndValidated(t *testing.T) {
	t.Parallel()
	const user = "delivery:\n  numbered_files: [a/*.sql]\n"
	for _, test := range []struct {
		name    string
		user    []byte
		project []byte
		want    []string
	}{
		{name: "absent"},
		{name: "project_replaces_user", user: []byte(user), project: []byte("delivery:\n  numbered_files: [drizzle/*.sql]\n"), want: []string{"drizzle/*.sql"}},
		{name: "user_inherited", user: []byte(user), project: []byte("delivery: {}\n"), want: []string{"a/*.sql"}},
		{name: "explicit_empty", user: []byte(user), project: []byte("delivery:\n  numbered_files: []\n"), want: []string{}},
	} {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := ResolveConfigProposal(test.user, test.project)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(cfg.Delivery.NumberedFiles, test.want) {
				t.Fatalf("numbered files = %#v, want %#v", cfg.Delivery.NumberedFiles, test.want)
			}
		})
	}
	t.Run("unsafe", func(t *testing.T) {
		for _, pattern := range []string{"[../x]", "[/abs]", `['a\b']`} {
			t.Run(pattern, func(t *testing.T) {
				_, err := ResolveConfigProposal(nil, []byte("delivery:\n  numbered_files: "+pattern+"\n"))
				if err == nil || !strings.Contains(err.Error(), "delivery.numbered_files has unsafe path") {
					t.Fatalf("error = %v, want unsafe numbered_files path", err)
				}
			})
		}
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := ResolveConfigProposal(nil, []byte("delivery:\n  numbered_files: ['[']\n"))
		if err == nil || !strings.Contains(err.Error(), "delivery.numbered_files has invalid pattern") {
			t.Fatalf("error = %v, want invalid numbered_files pattern", err)
		}
	})
}

func TestIsNumberedFileNeedsAMatchAndALeadingDigit(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name     string
		patterns []string
		path     string
		want     bool
	}{
		{name: "directory_file", patterns: []string{"drizzle/"}, path: "drizzle/0052_a.sql", want: true},
		{name: "directory_nested_file", patterns: []string{"drizzle/"}, path: "drizzle/meta/0052_snapshot.json", want: true},
		{name: "directory_journal", patterns: []string{"drizzle/"}, path: "drizzle/meta/_journal.json"},
		{name: "glob_file", patterns: []string{"drizzle/*.sql"}, path: "drizzle/0052_a.sql", want: true},
		{name: "glob_digit_inside_name", patterns: []string{"drizzle/*.sql"}, path: "drizzle/x0052.sql"},
		{name: "glob_other_directory", patterns: []string{"drizzle/*.sql"}, path: "db/0052_a.sql"},
		{name: "glob_nested_file", patterns: []string{"drizzle/*.sql"}, path: "drizzle/meta/0052_snapshot.json"},
		{name: "no_patterns", path: "drizzle/0052_a.sql"},
		{name: "zero", patterns: []string{"drizzle/"}, path: "drizzle/0.sql", want: true},
		{name: "nine", patterns: []string{"drizzle/"}, path: "drizzle/9.sql", want: true},
		{name: "before_zero", patterns: []string{"drizzle/"}, path: "drizzle/-1.sql"},
		{name: "after_nine", patterns: []string{"drizzle/"}, path: "drizzle/:1.sql"},
		{name: "unicode_digit", patterns: []string{"drizzle/"}, path: "drizzle/５２.sql"},
		{name: "empty_name", patterns: []string{"*"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			delivery := Delivery{NumberedFiles: test.patterns}
			if got := delivery.IsNumberedFile(test.path); got != test.want {
				t.Fatalf("IsNumberedFile(%q) with %v = %t, want %t", test.path, test.patterns, got, test.want)
			}
		})
	}
}

func TestPushRemoteStartingWithADashIsRefused(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"--upload-pack=touch injected", "-o", " -x"} {
		config := Builtin()
		config.Watch.PushRemote = value
		err := Validate(config)
		if err == nil || !strings.Contains(err.Error(), "watch.push_remote") || !strings.Contains(err.Error(), "must not start with -") {
			t.Fatalf("Validate(push_remote=%q) = %v, want the leading-dash refusal", value, err)
		}
	}
	for _, value := range []string{"", "origin", "upstream-fork"} {
		config := Builtin()
		config.Watch.PushRemote = value
		if err := Validate(config); err != nil {
			t.Fatalf("Validate(push_remote=%q) = %v, want nil", value, err)
		}
	}
}
