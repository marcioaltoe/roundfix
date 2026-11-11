// Suite: User Config Machine Verification Capacity.
// Boundary: layered YAML and templates in temporary homes only.
package config

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestRunsMaxVerifyingDefaultsToOne(t *testing.T) {
	t.Parallel()
	loaded, err := Load(LoadOptions{HomeDir: t.TempDir(), WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Runs.MaxVerifying != 1 || Builtin().Runs.MaxVerifying != 1 {
		t.Fatalf("default capacity = %d, builtin = %d", loaded.Config.Runs.MaxVerifying, Builtin().Runs.MaxVerifying)
	}
}

func TestRunsMaxVerifyingAcceptsZeroAndRefusesANegativeValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, value, wantError string
		want                   int
	}{
		{name: "off", value: "0", want: 0},
		{name: "positive", value: "3", want: 3},
		{name: "negative", value: "-1", wantError: "runs.max_verifying must be greater than or equal to 0"},
		{name: "not an integer", value: "invalid", wantError: "cannot unmarshal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  max_verifying: "+tc.value+"\n")
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
			if tc.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantError) {
					t.Fatalf("error = %v, want %q", err, tc.wantError)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Runs.MaxVerifying != tc.want {
				t.Fatalf("capacity = %d, want %d", loaded.Config.Runs.MaxVerifying, tc.want)
			}
		})
	}
}

func TestRunsMaxVerifyingInProjectConfigIsIgnoredWithAWarning(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "9", "-1", "invalid"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			gittest.MarkRepository(t, work)
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  max_verifying: 3\n")
			mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "runs:\n  max_verifying: "+value+"\n")
			var stderr bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work, Stderr: &stderr})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Runs.MaxVerifying != 3 {
				t.Fatalf("User Config overridden: %d", loaded.Config.Runs.MaxVerifying)
			}
			want := "config: runs.max_verifying in Project Config is ignored; set runs.max_verifying in User Config\n"
			if stderr.String() != want {
				t.Fatalf("warning = %q, want %q", stderr.String(), want)
			}
		})
	}
}

func TestTheUserConfigTemplateWritesRunsMaxVerifying(t *testing.T) {
	t.Parallel()
	lines := []string{"  # Maximum Verifications at once across repositories; 0 disables the bound.", "  max_verifying: 1"}
	for _, line := range lines {
		if !strings.Contains(DefaultConfigYAML(), line+"\n") {
			t.Errorf("user template missing %q", line)
		}
		if strings.Contains(DefaultProjectConfigYAML(), line) {
			t.Errorf("project template contains %q", line)
		}
	}
	if !strings.Contains(DefaultConfigYAML(), fmt.Sprintf("  max_active: %d\n", Builtin().Runs.MaxActive)+strings.Join(lines, "\n")+"\n") {
		t.Error("machine capacity must follow max_active")
	}
	home := t.TempDir()
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), DefaultConfigYAML())
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Runs.MaxVerifying != 1 {
		t.Fatalf("template capacity = %d", loaded.Config.Runs.MaxVerifying)
	}
}
