package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

func TestRunUpgradeFixtureMatrix(t *testing.T) {
	// Sequential: overrides package-level test seams.
	t.Run("newer release replaces binary", func(t *testing.T) {
		fake := newUpgradeFake(t)
		newBinary := []byte("#!/bin/sh\necho upgraded\n")
		fake.releaseTag = "v1.1.0"
		fake.assets = []app.ReleaseAsset{
			fake.platformAsset("roundfix_darwin_arm64", newBinary),
			fake.checksumAsset("roundfix_checksums.txt", checksumLine("roundfix_darwin_arm64", newBinary)),
		}
		withUpgradeFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected upgrade exit 0, got %d (stderr %q)", code, stderr.String())
		}
		if got := stdout.String(); got != "upgraded 1.0.0 → 1.1.0\n" {
			t.Fatalf("unexpected stdout %q", got)
		}
		content, err := os.ReadFile(fake.executablePath)
		if err != nil {
			t.Fatalf("read executable: %v", err)
		}
		if string(content) != string(newBinary) {
			t.Fatalf("expected binary replacement, got %q", string(content))
		}
		if stderr.Len() != 0 {
			t.Fatalf("expected no stderr, got %q", stderr.String())
		}
	})

	t.Run("current version is no-op", func(t *testing.T) {
		fake := newUpgradeFake(t)
		fake.releaseTag = "v1.0.0"
		fake.assets = []app.ReleaseAsset{fake.platformAsset("roundfix_darwin_arm64", []byte("same"))}
		withUpgradeFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected upgrade exit 0, got %d", code)
		}
		if got := stdout.String(); got != "already current 1.0.0\n" {
			t.Fatalf("unexpected stdout %q", got)
		}
		assertFileContent(t, fake.executablePath, "old binary\n")
		if fake.downloads != 0 {
			t.Fatalf("expected no downloads for current version, got %d", fake.downloads)
		}
		if stderr.Len() != 0 {
			t.Fatalf("expected no stderr, got %q", stderr.String())
		}
	})

	t.Run("empty releases are clean", func(t *testing.T) {
		fake := newUpgradeFake(t)
		fake.releaseErr = app.ErrNoReleases
		withUpgradeFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)

		if code != exitOK {
			t.Fatalf("expected upgrade exit 0, got %d", code)
		}
		if got := stdout.String(); got != "no releases published\n" {
			t.Fatalf("unexpected stdout %q", got)
		}
		assertFileContent(t, fake.executablePath, "old binary\n")
		if stderr.Len() != 0 {
			t.Fatalf("expected no stderr, got %q", stderr.String())
		}
	})

	t.Run("verify failure leaves binary untouched", func(t *testing.T) {
		fake := newUpgradeFake(t)
		newBinary := []byte("short")
		asset := fake.platformAsset("roundfix_darwin_arm64", newBinary)
		asset.Size = int64(len(newBinary) + 1)
		fake.releaseTag = "v1.1.0"
		fake.assets = []app.ReleaseAsset{asset}
		withUpgradeFakeDeps(t, fake)
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		code := runCLI(t, []string{"upgrade"}, &stdout, &stderr)

		if code != exitRunFailed {
			t.Fatalf("expected upgrade exit 1, got %d", code)
		}
		if stdout.Len() != 0 {
			t.Fatalf("expected no stdout, got %q", stdout.String())
		}
		assertFileContent(t, fake.executablePath, "old binary\n")
		for _, want := range []string{"upgrade failed", "manual fallback", fake.assets[0].DownloadURL, fake.executablePath} {
			if !strings.Contains(stderr.String(), want) {
				t.Fatalf("expected stderr to contain %q, got %q", want, stderr.String())
			}
		}
	})
}

func TestRunUpgradeCheckReportsAvailableWithoutInstalling(t *testing.T) {
	// Sequential: overrides package-level test seams.
	fake := newUpgradeFake(t)
	newBinary := []byte("#!/bin/sh\necho upgraded\n")
	fake.releaseTag = "v1.1.0"
	fake.assets = []app.ReleaseAsset{fake.platformAsset("roundfix_darwin_arm64", newBinary)}
	withUpgradeFakeDeps(t, fake)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"upgrade", "--check"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected upgrade --check exit 0, got %d", code)
	}
	if got := stdout.String(); got != "upgrade available 1.0.0 → 1.1.0\n" {
		t.Fatalf("unexpected stdout %q", got)
	}
	assertFileContent(t, fake.executablePath, "old binary\n")
	if fake.downloads != 0 {
		t.Fatalf("expected --check to avoid downloads, got %d", fake.downloads)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestVersionFreshnessCachesDailyAndReportsBehind(t *testing.T) {
	// Sequential: overrides package-level test seams.
	homeDir := t.TempDir()
	checkedAt := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	calls := 0
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{
		now: func() time.Time {
			return checkedAt
		},
		currentVersion: func() string {
			return "1.0.0"
		},
		latestRelease: func(context.Context) (string, []app.ReleaseAsset, error) {
			calls++
			return "v1.1.0", nil, nil
		},
	})
	var first bytes.Buffer
	var second bytes.Buffer
	loaded := roundconfig.Loaded{HomeDir: homeDir}

	maybeReportVersionFreshness(context.Background(), loaded, &first)
	checkedAt = checkedAt.Add(time.Hour)
	maybeReportVersionFreshness(context.Background(), loaded, &second)

	if calls != 1 {
		t.Fatalf("expected one release lookup inside 24h, got %d", calls)
	}
	wantLine := "roundfix 1.0.0 is behind latest 1.1.0; run roundfix upgrade\n"
	if first.String() != wantLine {
		t.Fatalf("expected first freshness line %q, got %q", wantLine, first.String())
	}
	if second.String() != wantLine {
		t.Fatalf("expected cached freshness line %q, got %q", wantLine, second.String())
	}
	cache := readFreshnessCacheFile(t, homeDir)
	if !strings.Contains(cache, `"latest_version": "1.1.0"`) || !strings.Contains(cache, `"checked_at"`) {
		t.Fatalf("expected cache to record latest version and check time, got %s", cache)
	}
}

func TestVersionFreshnessNetworkFailureIsSilentAndCachesAttempt(t *testing.T) {
	// Sequential: overrides package-level test seams.
	homeDir := t.TempDir()
	now := time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
	calls := 0
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{
		now: func() time.Time {
			return now
		},
		currentVersion: func() string {
			return "1.0.0"
		},
		latestRelease: func(context.Context) (string, []app.ReleaseAsset, error) {
			calls++
			return "", nil, errors.New("network unavailable")
		},
	})
	var first bytes.Buffer
	var second bytes.Buffer
	loaded := roundconfig.Loaded{HomeDir: homeDir}

	maybeReportVersionFreshness(context.Background(), loaded, &first)
	now = now.Add(time.Hour)
	maybeReportVersionFreshness(context.Background(), loaded, &second)

	if calls != 1 {
		t.Fatalf("expected failed attempt to be cached for 24h, got %d calls", calls)
	}
	if first.Len() != 0 || second.Len() != 0 {
		t.Fatalf("expected network failures to stay silent, first=%q second=%q", first.String(), second.String())
	}
	cache := readFreshnessCacheFile(t, homeDir)
	if !strings.Contains(cache, `"checked_at"`) || strings.Contains(cache, `"latest_version": "1.1.0"`) {
		t.Fatalf("expected cache to record only the attempt, got %s", cache)
	}
}

func TestVersionFreshnessFetchWiringDoesNotAffectOutcome(t *testing.T) {
	// Sequential: overrides package-level test seams.
	homeDir, repoDir := withCLIWorkspace(t)
	withSuccessfulPreflight(t, repoDir)
	withVersionFreshnessFakeDeps(t, versionFreshnessDependencies{
		now: func() time.Time {
			return time.Date(2026, 7, 5, 12, 0, 0, 0, time.UTC)
		},
		currentVersion: func() string {
			return "1.0.0"
		},
		latestRelease: func(context.Context) (string, []app.ReleaseAsset, error) {
			return "v1.1.0", nil, nil
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLI(t, []string{"fetch", "--source", "coderabbit", "--pr", "123", "--no-input"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("expected fetch exit 0, got %d (stderr %q)", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "roundfix 1.0.0 is behind latest 1.1.0; run roundfix upgrade") {
		t.Fatalf("expected freshness note in stderr, got %q", stderr.String())
	}
	if !strings.Contains(readFreshnessCacheFile(t, homeDir), `"latest_version": "1.1.0"`) {
		t.Fatalf("expected fetch to update freshness cache")
	}
}

type upgradeFake struct {
	executablePath string
	currentVersion string
	goos           string
	goarch         string
	releaseTag     string
	releaseErr     error
	assets         []app.ReleaseAsset
	contentByURL   map[string][]byte
	downloads      int
}

func newUpgradeFake(t *testing.T) *upgradeFake {
	t.Helper()
	dir := t.TempDir()
	executablePath := filepath.Join(dir, "roundfix")
	if err := os.WriteFile(executablePath, []byte("old binary\n"), 0o755); err != nil {
		t.Fatalf("write fake executable: %v", err)
	}
	return &upgradeFake{
		executablePath: executablePath,
		currentVersion: "1.0.0",
		goos:           "darwin",
		goarch:         "arm64",
		contentByURL:   map[string][]byte{},
	}
}

func (fake *upgradeFake) platformAsset(name string, content []byte) app.ReleaseAsset {
	asset := app.ReleaseAsset{
		Name:        name,
		APIURL:      "https://api.github.test/assets/" + name,
		DownloadURL: "https://github.test/download/" + name,
		Size:        int64(len(content)),
	}
	fake.contentByURL[asset.APIURL] = content
	return asset
}

func (fake *upgradeFake) checksumAsset(name string, content string) app.ReleaseAsset {
	asset := app.ReleaseAsset{
		Name:        name,
		APIURL:      "https://api.github.test/assets/" + name,
		DownloadURL: "https://github.test/download/" + name,
		Size:        int64(len(content)),
	}
	fake.contentByURL[asset.APIURL] = []byte(content)
	return asset
}

func withUpgradeFakeDeps(t *testing.T, fake *upgradeFake) {
	t.Helper()
	old := upgradeDeps
	upgradeDeps = upgradeDependencies{
		latestRelease: func(context.Context) (string, []app.ReleaseAsset, error) {
			return fake.releaseTag, fake.assets, fake.releaseErr
		},
		downloadAsset: func(_ context.Context, asset app.ReleaseAsset) ([]byte, error) {
			fake.downloads++
			content, ok := fake.contentByURL[asset.APIURL]
			if !ok {
				return nil, fmt.Errorf("missing fixture for %s", asset.APIURL)
			}
			return content, nil
		},
		executablePath: func() (string, error) {
			return fake.executablePath, nil
		},
		currentVersion: func() string {
			return fake.currentVersion
		},
		goos:   fake.goos,
		goarch: fake.goarch,
	}
	t.Cleanup(func() {
		upgradeDeps = old
	})
}

func withVersionFreshnessFakeDeps(t *testing.T, deps versionFreshnessDependencies) {
	t.Helper()
	old := versionFreshnessDeps
	versionFreshnessDeps = deps
	t.Cleanup(func() {
		versionFreshnessDeps = old
	})
}

func checksumLine(name string, content []byte) string {
	sum := sha256.Sum256(content)
	return fmt.Sprintf("%x  %s\n", sum, name)
}

func assertFileContent(t *testing.T, path string, want string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(content) != want {
		t.Fatalf("expected %s to contain %q, got %q", path, want, string(content))
	}
}

func readFreshnessCacheFile(t *testing.T, homeDir string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(homeDir, ".roundfix", "version-check.json"))
	if err != nil {
		t.Fatalf("read freshness cache: %v", err)
	}
	return string(content)
}
