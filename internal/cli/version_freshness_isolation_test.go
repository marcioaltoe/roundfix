// Suite: version freshness isolation.
// Invariant: operational CLI tests never reach the live release service by default.
// Boundary IN: in-process commands, the CLI helper process, and built-binary cache behavior.
// Boundary OUT: production release lookup and upgrade installation behavior.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/app"
)

const offlineTestReleaseTag = "v0.0.0-offline-test-lookup"

func init() {
	versionFreshnessDeps.latestRelease = func(context.Context) (string, []app.ReleaseAsset, error) {
		return offlineTestReleaseTag, nil, nil
	}
}

func seedFreshVersionCache(t *testing.T, homeDir string) {
	t.Helper()
	content, err := json.MarshalIndent(versionFreshnessCache{
		CheckedAt: time.Now().UTC(),
	}, "", "  ")
	if err != nil {
		t.Fatalf("marshal fresh version cache: %v", err)
	}
	content = append(content, '\n')
	path := versionFreshnessCachePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create version cache directory: %v", err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write fresh version cache: %v", err)
	}
}

func TestSuiteDefaultReleaseLookupStaysOffline(t *testing.T) {
	deps := defaultCommandDependencies().versionFreshness

	tag, assets, err := deps.latestRelease(context.Background())

	if err != nil {
		t.Fatalf("default release lookup error = %v", err)
	}
	if tag != offlineTestReleaseTag {
		t.Fatalf("default release lookup tag = %q, want %q", tag, offlineTestReleaseTag)
	}
	if len(assets) != 0 {
		t.Fatalf("default release lookup assets = %+v, want none", assets)
	}
	if app.CompareVersions(app.NormalizeVersion(tag), "0.0.0") >= 0 {
		t.Fatalf("offline tag %q must sort below every release", tag)
	}
}

func TestOperationalCommandsRecordTheOfflineLookup(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "fetch", args: []string{"fetch", "--offline-probe"}},
		{name: "resolve", args: []string{"resolve", "--offline-probe"}},
		{name: "watch", args: []string{"watch", "--offline-probe"}},
		{name: "implement", args: []string{"implement", "--offline-probe"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			homeDir, _ := withCLIWorkspace(t)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := runCLI(t, test.args, &stdout, &stderr)

			if code != exitPreflight {
				t.Fatalf("%s probe exit = %d, want %d; stdout=%q stderr=%q", test.name, code, exitPreflight, stdout.String(), stderr.String())
			}
			assertOfflineVersionCache(t, homeDir)
		})
	}
}

func TestHelperProcessRecordsTheOfflineLookup(t *testing.T) {
	homeDir, repoDir := withCLIWorkspace(t)

	stdout, stderr, code := runCLIHelper(t, repoDir, "", nil, "implement", "--offline-probe")

	if code != exitPreflight {
		t.Fatalf("helper implement probe exit = %d, want %d; stdout=%q stderr=%q", code, exitPreflight, stdout, stderr)
	}
	assertOfflineVersionCache(t, homeDir)
}

func TestBuiltBinaryLeavesTheSeededVersionCacheUntouched(t *testing.T) {
	binary := buildRoundfixBinaryForMacro(t)
	homeDir, repoDir := withCLIWorkspace(t)
	seedFreshVersionCache(t, homeDir)
	cachePath := versionFreshnessCachePath(homeDir)
	before, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read seeded version cache: %v", err)
	}
	command := exec.Command(binary, "implement", "--offline-probe")
	command.Dir = repoDir
	command.Env = withEnvValue(isolatedGitEnvForTest(), "HOME", homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Run()

	if code := exitCodeFromWait(err); code != exitPreflight {
		t.Fatalf("built implement probe exit = %d, want %d; stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	after, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read version cache after built binary: %v", err)
	}
	if !bytes.Equal(after, before) {
		t.Fatalf("built binary changed seeded version cache\nbefore: %s\nafter: %s", before, after)
	}
	if strings.Contains(stderr.String(), "is behind latest") {
		t.Fatalf("built binary printed an upgrade warning: %q", stderr.String())
	}
}

func TestPerTestReleaseLookupOverridesTheSuiteDefault(t *testing.T) {
	homeDir, _ := withCLIWorkspace(t)
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{
		now:            time.Now,
		currentVersion: func() string { return "1.0.0" },
		latestRelease: func(context.Context) (string, []app.ReleaseAsset, error) {
			return "v1.1.0", nil, nil
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"fetch", "--offline-probe"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("per-test override probe exit = %d, want %d; stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix 1.0.0 is behind latest 1.1.0; run roundfix upgrade") {
		t.Fatalf("per-test override did not print the upgrade warning: %q", stderr.String())
	}
	cache := readVersionCacheForIsolationTest(t, homeDir)
	if cache.LatestVersion != "1.1.0" {
		t.Fatalf("per-test override cache latest_version = %q, want %q", cache.LatestVersion, "1.1.0")
	}
}

func assertOfflineVersionCache(t *testing.T, homeDir string) {
	t.Helper()
	cache := readVersionCacheForIsolationTest(t, homeDir)
	want := app.NormalizeVersion(offlineTestReleaseTag)
	if cache.LatestVersion != want {
		t.Fatalf("version cache latest_version = %q, want %q", cache.LatestVersion, want)
	}
}

func readVersionCacheForIsolationTest(t *testing.T, homeDir string) versionFreshnessCache {
	t.Helper()
	path := versionFreshnessCachePath(homeDir)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read version cache %s: %v", path, err)
	}
	var cache versionFreshnessCache
	if err := json.Unmarshal(content, &cache); err != nil {
		t.Fatalf("decode version cache %s: %v", path, err)
	}
	if cache.CheckedAt.IsZero() {
		t.Fatalf("version cache %s has no checked_at: %s", path, content)
	}
	return cache
}
