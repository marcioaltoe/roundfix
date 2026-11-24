// Boundary: configuration loader, temporary homes and committed repository configuration.
package config

import (
	"bytes"
	"reflect"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"strings"
	"testing"
)

const retiredLaneWarnings = "config: delivery.linux_lane is deprecated and ignored; the Linux Lane was retired, so remove it\nconfig: delivery.linux_lane_image is deprecated and ignored; the Linux Lane was retired, so remove it\n"

func TestTheRetiredLaneKeysLoadWithOneWarningEach(t *testing.T) {
	t.Parallel()
	const legacy = "delivery:\n  linux_lane: true\n  linux_lane_image: golang:1.26\n"
	for _, name := range []string{"project_config", "user_config", "values_of_every_kind", "key_in_both_files_warns_once", "only_retired_keys_leave_an_empty_delivery", "unknown_key_still_fails", "config_at_a_commit_prints_nothing"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if name == "project_config" || name == "user_config" {
				for _, key := range []string{"linux_lane", "linux_lane_image"} {
					content := []byte("delivery:\n  " + key + ": {ignored: value}\n")
					var user, project []byte
					if name == "project_config" {
						project = content
					} else {
						user = content
					}
					opts := itemBinaryLoadOptions(t, user, project)
					var diagnostics bytes.Buffer
					opts.Stderr = &diagnostics
					loaded, err := Load(opts)
					want := "config: delivery." + key + " is deprecated and ignored; the Linux Lane was retired, so remove it\n"
					if err != nil || !reflect.DeepEqual(loaded.Config.Delivery, Builtin().Delivery) || diagnostics.String() != want {
						t.Fatalf("single key %s: delivery=%+v err=%v warnings=%q", key, loaded.Config.Delivery, err, diagnostics.String())
					}
				}
			}

			values := []string{"true"}
			if name == "values_of_every_kind" {
				values = []string{"true", `"yes"`, "{nested: value}", "[one, two]", ""}
			}
			for _, value := range values {
				project, user := legacy, ""
				if name == "user_config" {
					project, user = "", legacy
				}
				if name == "key_in_both_files_warns_once" {
					user = legacy
				}
				if name == "values_of_every_kind" {
					project = "delivery:\n  linux_lane: " + value + "\n  linux_lane_image: " + value + "\n"
				}
				if name == "unknown_key_still_fails" {
					project += "  bogus_key: true\n"
				}
				opts := itemBinaryLoadOptions(t, []byte(user), []byte(project))
				var diagnostics bytes.Buffer
				opts.Stderr = &diagnostics
				loaded, err := Load(opts)
				if name == "unknown_key_still_fails" {
					if err == nil || !strings.Contains(err.Error(), "field bogus_key not found in type config.deliveryOverlay") {
						t.Fatalf("error=%v", err)
					}
					if diagnostics.String() != retiredLaneWarnings {
						t.Fatalf("warnings=%q", diagnostics.String())
					}
					continue
				}
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(loaded.Config.Delivery, Builtin().Delivery) {
					t.Fatalf("delivery=%+v", loaded.Config.Delivery)
				}
				if diagnostics.String() != retiredLaneWarnings {
					t.Fatalf("warnings=%q", diagnostics.String())
				}
				if name == "config_at_a_commit_prints_nothing" {
					gittest.Run(t, opts.WorkDir, "add", ".")
					gittest.Run(t, opts.WorkDir, "commit", "-m", "legacy config")
					diagnostics.Reset()
					cfg, err := DeliveryConfigAtCommit(t.Context(), preflight.ExecGitRunner{}, opts.WorkDir, loaded.UserConfigPath, "HEAD")
					if err != nil || !reflect.DeepEqual(cfg.Delivery, Builtin().Delivery) || diagnostics.Len() != 0 {
						t.Fatalf("config=%+v err=%v warnings=%q", cfg.Delivery, err, diagnostics.String())
					}
				}
			}
		})
	}
}
