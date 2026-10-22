// Suite: User Config Delivery Check Wait.
// Invariant: only the machine owner controls a positive bounded wait.
// Boundary IN: layered YAML configuration in temporary homes.
// Boundary OUT: no host configuration or external tools.
package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDeliveryCheckTimeoutIsReadFromUserConfig(t *testing.T) {
	t.Parallel()
	if DefaultDeliveryCheckTimeout != 60*time.Minute {
		t.Fatalf("default = %s", DefaultDeliveryCheckTimeout)
	}
	for _, tc := range []struct {
		value string
		want  time.Duration
	}{{"", 60 * time.Minute}, {"45m", 45 * time.Minute}, {"24h", MaxDeliveryCheckTimeout}} {
		t.Run(tc.value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			if tc.value != "" {
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  check_timeout: "+tc.value+"\n")
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.CheckTimeout != tc.want {
				t.Fatalf("CheckTimeout = %v, want %v", loaded.Config.Delivery.CheckTimeout, tc.want)
			}
		})
	}
}

func TestProjectConfigCannotSetTheDeliveryCheckTimeout(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "5m", "-1h", "invalid"} {
		t.Run(value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustMkdir(t, filepath.Join(work, ".git"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  check_timeout: 45m\n")
			mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "delivery:\n  check_timeout: "+value+"\n")
			var stderr bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work, Stderr: &stderr})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.CheckTimeout != 45*time.Minute {
				t.Fatalf("CheckTimeout = %v", loaded.Config.Delivery.CheckTimeout)
			}
			want := "config: delivery.check_timeout in Project Config is ignored; set delivery.check_timeout in User Config\n"
			if stderr.String() != want {
				t.Fatalf("warning = %q, want %q", stderr.String(), want)
			}
		})
	}
}

func TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds(t *testing.T) {
	t.Parallel()
	const rangeError = "delivery.check_timeout must be greater than 0 and at most 24h"
	for _, value := range []string{"0", "0s", "-1m", "24h1s", "soon"} {
		t.Run(value, func(t *testing.T) {
			home := t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  check_timeout: "+value+"\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
			if value == "soon" {
				if err == nil || !strings.Contains(err.Error(), "invalid duration") {
					t.Fatalf("error = %v", err)
				}
			} else if err == nil || err.Error() != rangeError {
				t.Fatalf("error = %v, want %q", err, rangeError)
			}
		})
	}
}
