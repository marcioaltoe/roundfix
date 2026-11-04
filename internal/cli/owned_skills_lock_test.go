package cli

import (
	"bytes"
	"context"
	"fmt"
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
	before := `{"version":1, "skills":{"roundfix":{"extra":[1,true],"computedHash":"` + stale + `","sourceType":"upstream"},"external":{"computedHash":"` + stale + `"}}}` + "\n"
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

func TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries(t *testing.T) {
	t.Parallel()
	_, root := withCLIWorkspace(t)
	stale := strings.Repeat("0", 64)
	before := "{\n  \"version\": 1,\n  \"skills\": {\n"
	for _, name := range []string{"archive-spec", "roundfix", "write-prd"} {
		before += fmt.Sprintf(`    "%s": {
      "source": "marcioaltoe/skills",
      "ref": "8d71fade68c54693bedbd917371f0514311af157",
      "sourceType": "github",
      "skillPath": "skills/%s",
      "computedHash": "%s"
    },
`, name, name, stale)
	}
	before += `    "testing-boss": {"source":"external/skills","ref":"keep","sourceType":"github","computedHash":"` + stale + `"}
  }
}
`
	filename := filepath.Join(root, "skills-lock.json")
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"skills", "install"}, &stdout, &stderr)
	files, err := roundskills.Files()
	if err != nil {
		t.Fatal(err)
	}
	wantOutput := fmt.Sprintf("Installed Roundfix skill for project: %s (%d file(s))\n", filepath.Join(root, ".agents", "skills"), len(files)) +
		"Updated skills-lock.json computedHash: archive-spec\n" +
		"Updated skills-lock.json computedHash: roundfix\n" +
		"Updated skills-lock.json computedHash: write-prd\n" +
		"Removed skills-lock.json ref: archive-spec\n" +
		"Removed skills-lock.json ref: roundfix\n" +
		"Removed skills-lock.json ref: write-prd\n"
	if code != 0 || stderr.Len() != 0 || stdout.String() != wantOutput {
		t.Fatalf("exit=%d stderr=%s stdout=%s\nwant=%s", code, &stderr, &stdout, wantOutput)
	}
	wantLock := strings.ReplaceAll(before, "      \"ref\": \"8d71fade68c54693bedbd917371f0514311af157\",\n", "")
	for _, name := range []string{"archive-spec", "roundfix", "write-prd"} {
		hash, err := roundskills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", name))
		if err != nil {
			t.Fatal(err)
		}
		wantLock = strings.Replace(wantLock, stale, hash, 1)
	}
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != wantLock {
		t.Fatalf("lock=%s\nwant=%s", after, wantLock)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"skills", "install"}, &stdout, &stderr)
	wantSecond := strings.SplitN(wantOutput, "\n", 2)[0] + "\n"
	if code != 0 || stderr.Len() != 0 || stdout.String() != wantSecond {
		t.Fatalf("second exit=%d stderr=%s stdout=%s", code, &stderr, &stdout)
	}
	afterInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	after, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(info, afterInfo) || !info.ModTime().Equal(afterInfo.ModTime()) || string(after) != wantLock {
		t.Fatal("second install wrote lock")
	}
}

func TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	stale := strings.Repeat("0", 64)
	before := `{"version":1,"skills":{"roundfix":{"source":"marcioaltoe/skills","ref":"old","sourceType":"github","skillPath":"skills/roundfix","computedHash":"` + stale + `","extra":true},"external":{"ref":"keep","computedHash":"` + stale + `"}}}`
	filename := filepath.Join(root, "skills-lock.json")
	if err := os.WriteFile(filename, []byte(before), 0o644); err != nil {
		t.Fatal(err)
	}
	dependencies := ownedLockStageDependencies(root)
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, dependencies)
	if err != nil || !reflect.DeepEqual(result.LockUpdated, []string{"roundfix"}) || !reflect.DeepEqual(result.LockUnpinned, []string{"roundfix"}) {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	hash, err := roundskills.SkillFolderHash(t.Context(), filepath.Join(root, ".agents", "skills", "roundfix"))
	if err != nil {
		t.Fatal(err)
	}
	wantLock := strings.Replace(strings.Replace(before, `,"ref":"old"`, "", 1), stale, hash, 1)
	after, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != wantLock {
		t.Fatalf("lock=%s\nwant=%s", after, wantLock)
	}
	payload := newBaselineUpdateResult()
	payload.Skills = result
	var output bytes.Buffer
	if err := writeBaselineUpdateResult(payload, false, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Skills lock updated: 1\n- lock updated roundfix\n- lock ref removed roundfix\nSkills restored: ") {
		t.Fatalf("text=%s", &output)
	}
	output.Reset()
	if err := writeBaselineUpdateResult(payload, true, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"lockUpdated":["roundfix"],"lockUnpinned":["roundfix"]`) {
		t.Fatalf("json=%s", &output)
	}
	info, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	result, err = runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, dependencies)
	if err != nil || len(result.LockUpdated) != 0 || len(result.LockUnpinned) != 0 {
		t.Fatalf("second result=%+v err=%v", result, err)
	}
	afterInfo, err := os.Stat(filename)
	if err != nil {
		t.Fatal(err)
	}
	after, err = os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	if !os.SameFile(info, afterInfo) || !info.ModTime().Equal(afterInfo.ModTime()) || string(after) != wantLock {
		t.Fatal("second stage wrote lock")
	}
	payload.Skills = result
	for _, jsonOutput := range []bool{false, true} {
		output.Reset()
		if err := writeBaselineUpdateResult(payload, jsonOutput, &output); err != nil {
			t.Fatal(err)
		}
		for _, unwanted := range []string{"lockUpdated", "lockUnpinned", "- lock updated", "- lock ref removed"} {
			if strings.Contains(output.String(), unwanted) {
				t.Fatalf("second report=%s", &output)
			}
		}
	}
}
