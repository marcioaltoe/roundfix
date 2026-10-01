package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestProfilesConfigurePreviewsAndWritesADeviation(t *testing.T) {
	t.Parallel()
	_, repo := withCLIWorkspace(t)
	path := filepath.Join(repo, "fragment.yml")
	configPath := filepath.Join(repo, ".roundfixrc.yml")
	fragment := profilesConfigureTestFragment("chosen-model", "backup-model")
	mustWrite(t, path, fragment+"  deviation: {from: 2026-09-30, reason: Keep the validated model}\n")
	withSuccessfulProfilesConfigureProof(t)
	var stdout, stderr bytes.Buffer
	args := []string{"profiles", "configure", "--scope", "project", "--file", path, "--yes"}
	if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	if !strings.Contains(stdout.String(), "Deviation: from 2026-09-30 — Keep the validated model\n") {
		t.Fatalf("preview missing deviation: %s", &stdout)
	}
	persisted, err := roundconfig.ParseProfilesFragment([]byte(mustRead(t, configPath)))
	if err != nil {
		t.Fatal(err)
	}
	want := roundconfig.ProfileDeviation{From: "2026-09-30", Reason: "Keep the validated model"}
	if got := persisted[roundconfig.CategoryBackend].Deviation; got == nil || *got != want {
		t.Fatalf("persisted deviation = %+v", got)
	}
	if content := mustRead(t, configPath); strings.Index(content, "deviation:") < strings.Index(content, "fallbacks:") {
		t.Fatal("deviation must follow fallbacks")
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, append(args, "--json"), &stdout, &stderr); code != exitOK {
		t.Fatalf("JSON exit=%d stderr=%s", code, &stderr)
	}
	var response struct {
		Profiles []struct {
			Deviation *roundconfig.ProfileDeviation `json:"deviation"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if len(response.Profiles) != 1 || response.Profiles[0].Deviation == nil || *response.Profiles[0].Deviation != want {
		t.Fatalf("JSON missing deviation: %s", &stdout)
	}
	mustWrite(t, path, fragment)
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
		t.Fatalf("replacement exit=%d stderr=%s", code, &stderr)
	}
	if content := mustRead(t, configPath); strings.Contains(content, "deviation:") {
		t.Fatalf("replacement retained deviation: %s", content)
	}
}

func TestProfilesConfigureOutputIsUnchangedWithoutADeviation(t *testing.T) {
	t.Parallel()
	_, repo := withCLIWorkspace(t)
	path := filepath.Join(repo, "fragment.yml")
	configPath := filepath.Join(repo, ".roundfixrc.yml")
	mustWrite(t, path, profilesConfigureTestFragment("chosen-model", "backup-model"))
	withSuccessfulProfilesConfigureProof(t)
	var stdout, stderr bytes.Buffer
	args := []string{"profiles", "configure", "--scope", "project", "--file", path, "--dry-run"}
	if code := runCLI(t, args, &stdout, &stderr); code != exitOK {
		t.Fatalf("exit=%d stderr=%s", code, &stderr)
	}
	wantText := fmt.Sprintf("Profile Configure Preview\nScope: project\nPath: %s\nChanges:\nadded: backend\nCategory: backend\nPreferred Selection: codex / chosen-model / high\nFallback Chain:\n  1. codex / backup-model / xhigh\nProfile configuration dry run: %s\n", configPath, configPath)
	if stdout.String() != wantText {
		t.Fatalf("text changed:\nwant %q\ngot  %q", wantText, stdout.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI(t, append(args, "--json"), &stdout, &stderr); code != exitOK {
		t.Fatalf("JSON exit=%d stderr=%s", code, &stderr)
	}
	encodedPath, err := json.Marshal(configPath)
	if err != nil {
		t.Fatal(err)
	}
	wantJSON := fmt.Sprintf(`{"schema":%q,"changed":false,"refused":false,"scope":"project","path":%s,"profiles":[{"category":"backend","preferred":{"runtime":"codex","model":"chosen-model","reasoning_effort":"high"},"fallbacks":[{"runtime":"codex","model":"backup-model","reasoning_effort":"xhigh"}]}],"changes":[{"category":"backend","kind":"added"}]}`+"\n", profilesConfigureSchema, encodedPath)
	if stdout.String() != wantJSON {
		t.Fatalf("JSON changed:\nwant %q\ngot  %q", wantJSON, stdout.String())
	}
}
