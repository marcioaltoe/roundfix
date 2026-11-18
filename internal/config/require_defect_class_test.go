package config

import (
	"bytes"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestRequireDefectClassIsAProjectConfigBoolean(t *testing.T) {
	t.Parallel()
	const warning = "config: delivery.require_defect_class in User Config is ignored; set delivery.require_defect_class in Project Config\n"
	for _, tc := range []struct {
		name, user, project string
		want                bool
		warned              bool
	}{
		{name: "absent"}, {name: "true", project: "delivery:\n  require_defect_class: true\n", want: true},
		{name: "false", project: "delivery:\n  require_defect_class: false\n"},
		{name: "user true", user: "delivery:\n  require_defect_class: true\n", warned: true},
		{name: "user invalid", user: "delivery:\n  require_defect_class: [not, boolean]\n", warned: true},
		{name: "project over user", user: "delivery:\n  require_defect_class: nope\n", project: "delivery:\n  require_defect_class: true\n", want: true, warned: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := itemBinaryLoadOptions(t, []byte(tc.user), []byte(tc.project))
			var diagnostics bytes.Buffer
			options.Stderr = &diagnostics
			loaded, err := Load(options)
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.RequireDefectClass != tc.want {
				t.Fatalf("delivery = %+v", loaded.Config.Delivery)
			}
			expected := ""
			if tc.warned {
				expected = warning
			}
			if diagnostics.String() != expected {
				t.Fatalf("warning = %q, want %q", diagnostics.String(), expected)
			}
			gittest.Run(t, options.WorkDir, "add", ".")
			gittest.Run(t, options.WorkDir, "commit", "--allow-empty", "-m", "default config")
			config, err := DeliveryConfigAtCommit(t.Context(), preflight.ExecGitRunner{}, options.WorkDir, loaded.UserConfigPath, "HEAD")
			if err != nil || config.Delivery.RequireDefectClass != tc.want {
				t.Fatalf("at commit = %+v, %v", config.Delivery, err)
			}
		})
	}
	for _, value := range []string{"'true'", "yes", "1", "[]", "{}", "null"} {
		t.Run("refuses "+value, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte("delivery:\n  require_defect_class: "+value+"\n"))
			const want = `parse config "Project Config proposal": delivery.require_defect_class must be boolean: cannot unmarshal non-boolean value`
			if err == nil || !strings.Contains(err.Error(), want) {
				t.Fatalf("error = %v, want %s", err, want)
			}
		})
	}
}
