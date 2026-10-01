package cli

import (
	"bytes"
	"roundfix/internal/app"
	"testing"
)

func TestUpgradeStdoutAndExitCodesAreCharacterized(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name                       string
		args                       []string
		tag                        string
		noRelease, downloadFailure bool
		want                       string
		code                       int
	}{
		{name: "no release", noRelease: true, want: "no releases published\n"},
		{name: "current", tag: "v1.0.0", want: "already current 1.0.0\n"},
		{name: "check newer", args: []string{"--check"}, tag: "v1.1.0", want: "upgrade available 1.0.0 → 1.1.0\n"},
		{name: "check current", args: []string{"--check"}, tag: "v1.0.0", want: "already current 1.0.0\n"},
		{name: "installed", tag: "v1.1.0", want: "upgraded 1.0.0 → 1.1.0\n"},
		{name: "download failed", tag: "v1.1.0", downloadFailure: true, code: exitRunFailed},
		{name: "unknown flag", args: []string{"--unknown"}, code: exitPreflight},
		{name: "help", args: []string{"--help"}, want: commandUsage("upgrade")},
	} {
		t.Run(test.name, func(t *testing.T) {
			fake := newUpgradeFake(t)
			fake.releaseTag = test.tag
			fake.assets = []app.ReleaseAsset{fake.platformAsset("roundfix_darwin_arm64", []byte("new binary"))}
			if test.noRelease {
				fake.releaseErr = app.ErrNoReleases
			}
			if test.downloadFailure {
				fake.contentByURL = map[string][]byte{}
			}
			withUpgradeFakeDeps(t, fake)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, append([]string{"upgrade"}, test.args...), &stdout, &stderr)
			if code != test.code || stdout.String() != test.want {
				t.Fatalf("exit=%d stdout=%q; want exit=%d stdout=%q", code, stdout.String(), test.code, test.want)
			}
		})
	}
}
