package cli

import (
	"testing"

	"roundfix/internal/app"
)

// TestSelectPlatformAssetResolvesReleaseNamingScheme pins the release-asset
// naming scheme (dist/npm/platforms.json and .github/workflows/release.yml) to
// the Upgrade Command's platform selection so the two cannot drift: every
// per-platform asset the release workflow uploads must resolve for its
// GOOS/GOARCH, and checksum-style assets must never be selected as the binary.
func TestSelectPlatformAssetResolvesReleaseNamingScheme(t *testing.T) {
	assets := []app.ReleaseAsset{
		{Name: "roundfix-darwin-arm64"},
		{Name: "roundfix-darwin-amd64"},
		{Name: "roundfix-linux-arm64"},
		{Name: "roundfix-linux-amd64"},
		{Name: "roundfix-windows-amd64.exe"},
		{Name: "SHA256SUMS"},
	}

	cases := []struct {
		goos, goarch, want string
	}{
		{"darwin", "arm64", "roundfix-darwin-arm64"},
		{"darwin", "amd64", "roundfix-darwin-amd64"},
		{"linux", "arm64", "roundfix-linux-arm64"},
		{"linux", "amd64", "roundfix-linux-amd64"},
		{"windows", "amd64", "roundfix-windows-amd64.exe"},
	}

	for _, c := range cases {
		got, ok := selectPlatformAsset(assets, c.goos, c.goarch)
		if !ok {
			t.Fatalf("no asset resolved for %s/%s", c.goos, c.goarch)
		}
		if got.Name != c.want {
			t.Fatalf("%s/%s: resolved %q, want %q", c.goos, c.goarch, got.Name, c.want)
		}
	}

	if _, ok := selectPlatformAsset([]app.ReleaseAsset{{Name: "SHA256SUMS"}}, "linux", "amd64"); ok {
		t.Fatal("a checksum-style asset must not be selected as the platform binary")
	}
}
