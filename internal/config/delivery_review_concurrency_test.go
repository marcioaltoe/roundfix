// Suite: User Config Delivery Review Capacity.
// Invariant: only the machine owner controls a positive Review Capacity.
// Boundary IN: layered YAML configuration in temporary homes and repositories.
// Boundary OUT: no host configuration or external tools.
package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestDeliveryReviewConcurrencyIsReadFromUserConfig(t *testing.T) {
	t.Parallel()
	if DefaultDeliveryReviewConcurrency != 2 {
		t.Fatalf("default = %d, want 2", DefaultDeliveryReviewConcurrency)
	}
	for _, tc := range []struct {
		name  string
		value string
		want  int
	}{
		{name: "default", want: 2},
		{name: "one", value: "1", want: 1},
		{name: "four", value: "4", want: 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			if tc.value != "" {
				mustMkdir(t, filepath.Join(home, ".roundfix"))
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  review_concurrency: "+tc.value+"\n")
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.ReviewConcurrency != tc.want {
				t.Fatalf("ReviewConcurrency = %d, want %d", loaded.Config.Delivery.ReviewConcurrency, tc.want)
			}
		})
	}
}

func TestProjectConfigCannotSetTheDeliveryReviewConcurrency(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "5", "-1", "two"} {
		t.Run(value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			gittest.MarkRepository(t, work)
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  review_concurrency: 3\n")
			mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "delivery:\n  review_concurrency: "+value+"\n")
			var stderr bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work, Stderr: &stderr})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Delivery.ReviewConcurrency != 3 {
				t.Fatalf("ReviewConcurrency = %d, want 3", loaded.Config.Delivery.ReviewConcurrency)
			}
			want := "config: delivery.review_concurrency in Project Config is ignored; set delivery.review_concurrency in User Config\n"
			if stderr.String() != want {
				t.Fatalf("warning = %q, want %q", stderr.String(), want)
			}
		})
	}
}

func TestDeliveryReviewConcurrencyMustBePositive(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "-1", "two"} {
		t.Run(value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "delivery:\n  review_concurrency: "+value+"\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err == nil {
				t.Fatal("expected invalid Review Capacity to fail Load")
			}
			if value != "two" && !strings.Contains(err.Error(), "delivery.review_concurrency must be greater than 0") {
				t.Fatalf("error = %v", err)
			}
		})
	}
}
