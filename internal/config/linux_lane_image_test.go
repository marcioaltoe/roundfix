// Boundary: real config parsing and immutable default-branch Project Config.
package config

import (
	"fmt"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, user, project, want string }{
		{name: "default"},
		{name: "project", project: "delivery:\n  linux_lane_image: golang:1.26\n", want: "golang:1.26"},
		{name: "user", user: "delivery:\n  linux_lane_image: golang:1.26\n", want: "golang:1.26"},
		{name: "project over user", user: "delivery:\n  linux_lane_image: golang:1.25\n", project: "delivery:\n  linux_lane_image: golang:1.26\n", want: "golang:1.26"},
		{name: "explicit empty", user: "delivery:\n  linux_lane_image: golang:1.25\n", project: "delivery:\n  linux_lane_image: ''\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := itemBinaryLoadOptions(t, []byte(tc.user), []byte(tc.project))
			loaded, err := Load(options)
			if err != nil {
				t.Fatal(err)
			}
			if got := loaded.Config.Delivery.LinuxLaneImage; got != tc.want {
				t.Fatalf("image=%q want %q", got, tc.want)
			}
			gittest.Run(t, options.WorkDir, "add", ".")
			gittest.Run(t, options.WorkDir, "commit", "--allow-empty", "-m", "default config")
			config, err := DeliveryConfigAtCommit(t.Context(), preflight.ExecGitRunner{}, options.WorkDir, loaded.UserConfigPath, "HEAD")
			if err != nil || config.Delivery.LinuxLaneImage != tc.want {
				t.Fatalf("commit config=%+v err=%v", config.Delivery, err)
			}
		})
	}
}

func TestDeliveryLinuxLaneImageRefusesAnInvalidReference(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"'-v'", "'golang 1.26'", "'Golang:1.26'", "' golang:1.26'", "'golang:1.26 --privileged'", "[]", "{}", "12", "null", "true"} {
		t.Run("refuses "+value, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte("delivery:\n  linux_lane_image: "+value+"\n"))
			want := `parse config "Project Config proposal": delivery.linux_lane_image must be a container image reference, such as golang:1.26`
			if err == nil || err.Error() != want {
				t.Fatalf("error=%v want %s", err, want)
			}
		})
	}
	for _, value := range []string{"", "golang:1.26", "registry.example.com/team/go-ci:1.26", "localhost:5000/go-ci", "go-ci@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123"} {
		t.Run("accepts "+value, func(t *testing.T) {
			_, err := ResolveConfigProposal(nil, []byte(fmt.Sprintf("delivery:\n  linux_lane_image: %q\n", value)))
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
