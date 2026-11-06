package app

import (
	"context"
	"errors"
	"reflect"
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

func TestListReleasesReadsTitlesAndSkipsDraftsAndPrereleases(t *testing.T) {
	withGitHubAPIRunner(t, func(ctx context.Context, endpoint string) ([]byte, error) {
		if endpoint != "repos/"+ReleaseRepository+"/releases?per_page=30" {
			t.Fatalf("endpoint = %q", endpoint)
		}
		return []byte(`[
   {"tag_name":"v1.8.0","body":"No headings here","html_url":"https://github.test/v1.8.0"},
   {"tag_name":"v2.0.0","body":"### First title  \ntext\n## Ignored\n ### Also ignored\n### Second title","html_url":"https://github.test/v2.0.0"},
   {"tag_name":"v3.0.0","body":"### Draft","draft":true},
   {"tag_name":"v2.1.0-rc.1","body":"### Prerelease","prerelease":true},
   {"tag_name":"v1.10.0","body":"### Middle title","html_url":"https://github.test/v1.10.0"}
  ]`), nil
	})
	got, err := ListReleases(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	want := []Release{
		{Version: "2.0.0", Titles: []string{"First title", "Second title"}, URL: "https://github.test/v2.0.0"},
		{Version: "1.10.0", Titles: []string{"Middle title"}, URL: "https://github.test/v1.10.0"},
		{Version: "1.8.0", Titles: []string{}, URL: "https://github.test/v1.8.0"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("releases = %#v, want %#v", got, want)
	}
}

func TestListReleasesMapsGitHub404ToNoReleases(t *testing.T) {
	for _, tt := range []struct {
		name, payload string
		err           error
	}{
		{name: "404", err: errors.New("gh: Not Found (HTTP 404)")},
		{name: "empty", payload: "[]"},
		{name: "only unpublished", payload: `[{"tag_name":"v2.0.0","draft":true},{"tag_name":"v3.0.0","prerelease":true}]`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			withGitHubAPIRunner(t, func(context.Context, string) ([]byte, error) { return []byte(tt.payload), tt.err })
			_, err := ListReleases(t.Context())
			if !errors.Is(err, ErrNoReleases) {
				t.Fatalf("error = %v, want ErrNoReleases", err)
			}
		})
	}
}
