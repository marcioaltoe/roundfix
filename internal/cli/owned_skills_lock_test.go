package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/baseline"
	"strings"
	"testing"

	roundskills "roundfix/skills"
)

func TestSkillsInstallRewritesTheOwnedSkillLockHash(t *testing.T) {
	t.Parallel()
	_, root := withCLIWorkspace(t)
	stale := strings.Repeat("0", 64)
	before := `{"version":1, "skills":{"roundfix":{"extra":[1,true],"computedHash":"` + stale + `","ref":"upstream"},"external":{"computedHash":"` + stale + `"}}}` + "\n"
	filename := filepath.Join(root, "skills-lock.json")
	if err := os.WriteFile(filename, []byte(before), 0o640); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"skills", "install"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	hash, err := roundskills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", "roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Replace(before, stale, hash, 1); string(after) != want {
		t.Fatalf("lock=%s want=%s", after, want)
	}
	if !strings.Contains(stdout.String(), "Updated skills-lock.json computedHash: roundfix\n") {
		t.Fatalf("stdout=%s", &stdout)
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, []string{"skills", "install"}, &stdout, &stderr); code != 0 || strings.Contains(stdout.String(), "Updated skills-lock.json") {
		t.Fatalf("second exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
}

func TestBaselineUpdateSkillsStageReportsTheRewrittenLockEntries(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	stale := strings.Repeat("0", 64)
	before := `{"version":1,"skills":{"roundfix":{"computedHash":"` + stale + `","extra":true},"external":{"computedHash":"` + stale + `"}}}`
	filename := filepath.Join(root, "skills-lock.json")
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	dependencies := ownedLockStageDependencies(root)
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, dependencies)
	if err != nil || !reflect.DeepEqual(result.LockUpdated, []string{"roundfix"}) {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	hash, err := roundskills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", "roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != strings.Replace(before, stale, hash, 1) {
		t.Fatalf("lock=%s", after)
	}
	payload := newBaselineUpdateResult()
	payload.Skills = result
	var output bytes.Buffer
	if err := writeBaselineUpdateResult(payload, false, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "- installed "+result.Installed[len(result.Installed)-1]+"\nSkills lock updated: 1\n- lock updated roundfix\n") {
		t.Fatalf("text=%s", &output)
	}
	output.Reset()
	if err := writeBaselineUpdateResult(payload, true, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"lockUpdated":["roundfix"]`) {
		t.Fatalf("json=%s", &output)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	result, err = runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, dependencies)
	if err != nil || len(result.LockUpdated) != 0 {
		t.Fatalf("second result=%+v err=%v", result, err)
	}
	afterInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(info, afterInfo) || !info.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatal("second stage wrote lock")
	}
	output.Reset()
	payload.Skills = result
	if err := writeBaselineUpdateResult(payload, true, &output); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "lockUpdated") {
		t.Fatalf("empty lockUpdated not omitted: %s", &output)
	}
}

func ownedLockStageDependencies(root string) baselineUpdateSkillsDependencies {
	return baselineUpdateSkillsDependencies{
		resolveProjectRoot: func(context.Context, string) (string, error) { return root, nil },
		install:            roundskills.Install,
		ownedNames:         roundskills.Names,
		recordLock:         baseline.RecordInstalledOwnedSkills,
		resolveExternal:    func(string) ([]string, bool, error) { return nil, true, nil },
	}
}

func TestBaselineUpdateSkillsStageFailsOnAnUnsafeLock(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"malformed", "symlink", "directory"} {
		t.Run(kind, func(t *testing.T) {
			_, root := withCLIWorkspace(t)
			filename := filepath.Join(root, "skills-lock.json")
			before := []byte("{invalid")
			switch kind {
			case "symlink":
				if err := os.WriteFile(filepath.Join(root, "target"), before, 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink("target", filename); err != nil {
					t.Fatal(err)
				}
			case "directory":
				if err := os.Mkdir(filename, 0o755); err != nil {
					t.Fatal(err)
				}
			default:
				if err := os.WriteFile(filename, before, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			info, err := os.Lstat(filename)
			if err != nil {
				t.Fatal(err)
			}
			result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, ownedLockStageDependencies(root))
			if err == nil || result.Status != baselineUpdateSkillsFailed || !strings.HasPrefix(err.Error(), "record installed Roundfix skills in skills-lock.json: ") {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"skills", "install"}, &stdout, &stderr); code != 2 || !strings.Contains(stderr.String(), "roundfix: skills install failed:") || !strings.Contains(stderr.String(), "skills-lock.json") {
				t.Fatalf("exit=%d stderr=%s", code, &stderr)
			}
			after, err := os.Lstat(filename)
			if err != nil {
				t.Fatal(err)
			}
			if !os.SameFile(info, after) || info.Mode() != after.Mode() {
				t.Fatal("unsafe lock replaced")
			}
			if kind != "directory" {
				data, err := os.ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(data, before) {
					t.Fatal("unsafe lock changed")
				}
			}
			// An explicit directory bypasses even an unsafe repository lock.
			stdout.Reset()
			stderr.Reset()
			if code := runCLI(t, []string{"skills", "install", "--dir", filepath.Join(root, "custom")}, &stdout, &stderr); code != 0 {
				t.Fatalf("override exit=%d stderr=%s", code, &stderr)
			}
		})
	}
}
