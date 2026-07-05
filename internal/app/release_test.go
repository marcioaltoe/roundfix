package app

import (
	"context"
	"errors"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name  string
		left  string
		right string
		want  int
	}{
		{name: "newer stable", left: "v1.2.0", right: "1.1.9", want: 1},
		{name: "equal with tag prefix", left: "v1.2.3", right: "1.2.3", want: 0},
		{name: "older patch", left: "1.2.2", right: "1.2.3", want: -1},
		{name: "dev prerelease lower than stable", left: "0.0.0-dev", right: "0.0.0", want: -1},
		{name: "build metadata hyphen equals stable", left: "1.2.3+build-info", right: "1.2.3", want: 0},
		{name: "prerelease with build metadata lower than stable", left: "1.2.3-rc.1+build-info", right: "1.2.3", want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareVersions(tt.left, tt.right)
			if got != tt.want {
				t.Fatalf("CompareVersions(%q, %q) = %d, want %d", tt.left, tt.right, got, tt.want)
			}
		})
	}
}

func TestLatestReleaseParsesGitHubCLIOutput(t *testing.T) {
	withGitHubAPIRunner(t, func(ctx context.Context, endpoint string) ([]byte, error) {
		if endpoint != "repos/"+ReleaseRepository+"/releases/latest" {
			t.Fatalf("unexpected gh api endpoint %q", endpoint)
		}
		return []byte(`{
			"tag_name": "v1.2.3",
			"assets": [
				{
					"name": "roundfix_darwin_arm64",
					"url": "https://api.github.test/assets/1",
					"browser_download_url": "https://github.test/download/roundfix_darwin_arm64",
					"size": 123
				}
			]
		}`), nil
	})

	tag, assets, err := LatestRelease(context.Background())

	if err != nil {
		t.Fatalf("LatestRelease returned error: %v", err)
	}
	if tag != "v1.2.3" {
		t.Fatalf("expected tag v1.2.3, got %q", tag)
	}
	if len(assets) != 1 {
		t.Fatalf("expected one asset, got %d", len(assets))
	}
	if assets[0].Name != "roundfix_darwin_arm64" ||
		assets[0].APIURL != "https://api.github.test/assets/1" ||
		assets[0].DownloadURL != "https://github.test/download/roundfix_darwin_arm64" ||
		assets[0].Size != 123 {
		t.Fatalf("unexpected asset: %#v", assets[0])
	}
}

func TestLatestReleaseMapsGitHub404ToNoReleases(t *testing.T) {
	withGitHubAPIRunner(t, func(context.Context, string) ([]byte, error) {
		return nil, errors.New("gh: Not Found (HTTP 404)")
	})

	_, _, err := LatestRelease(context.Background())

	if !errors.Is(err, ErrNoReleases) {
		t.Fatalf("expected ErrNoReleases, got %v", err)
	}
}

func withGitHubAPIRunner(t *testing.T, runner func(context.Context, string) ([]byte, error)) {
	t.Helper()
	old := runGitHubAPI
	runGitHubAPI = runner
	t.Cleanup(func() {
		runGitHubAPI = old
	})
}
