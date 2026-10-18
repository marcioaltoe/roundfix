// Boundary: real config loading and immutable default-branch Project Config.
package config

import (
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, user, project string
		want                bool
	}{
		{name: "default"},
		{name: "project enables", project: "delivery:\n  linux_lane: true\n", want: true},
		{name: "user enables", user: "delivery:\n  linux_lane: true\n", want: true},
		{name: "project disables user", user: "delivery:\n  linux_lane: true\n", project: "delivery:\n  linux_lane: false\n"},
		{name: "project enables over user", user: "delivery:\n  linux_lane: false\n", project: "delivery:\n  linux_lane: true\n", want: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := itemBinaryLoadOptions(t, []byte(tc.user), []byte(tc.project))
			loaded, err := Load(options)
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.LinuxLane != tc.want {
				t.Fatalf("linux_lane=%v want %v", loaded.Config.Delivery.LinuxLane, tc.want)
			}
			gittest.Run(t, options.WorkDir, "add", ".")
			gittest.Run(t, options.WorkDir, "commit", "--allow-empty", "-m", "default config")
			config, err := DeliveryConfigAtCommit(t.Context(), preflight.ExecGitRunner{}, options.WorkDir, loaded.UserConfigPath, "HEAD")
			if err != nil || config.Delivery.LinuxLane != tc.want {
				t.Fatalf("commit config=%+v err=%v", config.Delivery, err)
			}
		})
	}
	for _, value := range []string{"'true'", "yes", "1", "[]", "{}", "null"} {
		t.Run("refuses "+value, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte("delivery:\n  linux_lane: "+value+"\n"))
			if err == nil || !strings.Contains(err.Error(), "delivery.linux_lane must be boolean: cannot unmarshal non-boolean value") {
				t.Fatalf("error=%v", err)
			}
		})
	}
}
