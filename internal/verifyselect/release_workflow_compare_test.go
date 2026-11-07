// Suite: release workflow content comparison and partial publication recovery.
// Invariant: only byte-identical decompressed tar streams may be skipped; every comparison precedes publication.
// Boundary IN: release workflow YAML, fixture tarballs and packuments, and the isolated step-script harness.
// Boundary OUT: network access, real npm publication, builds, GitHub Releases, and process signaling.
package verifyselect_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReleaseCompareAcceptsOnlyTheSameContent(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct {
		name, failure string
		launcher      bool
		empty         bool
	}{
		{name: "same tar stream with different gzip bytes"},
		{name: "launcher uses its own package directory", launcher: true},
		{name: "different file body", failure: "different"},
		{name: "missing registry tarball", failure: "registry pack"},
		{name: "empty published list", empty: true},
		{name: "missing local tarball", failure: "local pack"},
		{name: "unreadable local pack fields", failure: "local fields"},
		{name: "unreadable registry pack fields", failure: "registry fields"},
		{name: "corrupt local gzip", failure: "local gzip"},
		{name: "corrupt registry gzip", failure: "registry gzip"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			h := newReleaseWorkflowHarness(t, releaseWorkflowCase{
				targetVersion: "v" + releaseTestVersion, refName: "v" + releaseTestVersion, eventName: "push",
			})
			coordinate := h.coordinates[0]
			if scenario.launcher {
				coordinate = "roundfix"
			}
			// A following equal coordinate proves every failure continues the loop.
			following := h.coordinates[1]
			published := coordinate + "\n" + following + "\n"
			if scenario.empty {
				published = ""
			}
			writeFile(t, h.runner, "roundfix-release-published", published)
			writeFile(t, h.runner, "roundfix-release-complete", "stale-coordinate\n")
			body := "build for " + coordinate
			registryBody := body
			if scenario.failure == "different" {
				registryBody += " changed"
			}
			localIntegrity, registryIntegrity := releaseCompareFixtures(t, h, coordinate, releaseTestVersion, registryBody)
			releaseCompareFixtures(t, h, following, releaseTestVersion, "build for "+following)
			local := "local/" + strings.TrimPrefix(coordinate, "@roundfix/")
			registry := "registry/" + coordinate + "@" + releaseTestVersion
			switch scenario.failure {
			case "local pack", "registry pack":
				directory := local
				if scenario.failure == "registry pack" {
					directory = registry
				}
				if err := os.Remove(filepath.Join(h.fixtures, directory, "package.tgz")); err != nil {
					t.Fatal(err)
				}
			case "local fields", "registry fields":
				directory := local
				if scenario.failure == "registry fields" {
					directory = registry
				}
				// The stub emits invalid JSON, so jq's field read must fail.
				writeFile(t, h.fixtures, directory+"/integrity", `"`)
			case "local gzip", "registry gzip":
				directory := local
				if scenario.failure == "registry gzip" {
					directory = registry
				}
				writeFile(t, h.fixtures, directory+"/package.tgz", "not gzip")
			}

			wantExit := 0
			wantOutput := fmt.Sprintf("%-36s %-12s %s\n", coordinate, releaseTestVersion, "same content")
			wantComplete := published
			wantSummary := "\n### Already published\n\n"
			if scenario.failure != "" {
				wantExit = 1
				wantComplete = following + "\n"
				switch scenario.failure {
				case "different":
					wantOutput = fmt.Sprintf("::error::registry: %s@%s is already published with different content (registry %s, this build %s)\n", coordinate, releaseTestVersion, registryIntegrity, localIntegrity)
				case "local pack", "local fields":
					wantOutput = fmt.Sprintf("::error::undetermined: %s@%s: this build could not be packed\n", coordinate, releaseTestVersion)
				case "registry pack", "registry fields":
					wantOutput = fmt.Sprintf("::error::undetermined: %s@%s: the published tarball could not be fetched\n", coordinate, releaseTestVersion)
				case "local gzip", "registry gzip":
					wantOutput = fmt.Sprintf("::error::undetermined: %s@%s: a tarball could not be decompressed\n", coordinate, releaseTestVersion)
				}
			} else {
				wantSummary += fmt.Sprintf("- `%s@%s`: same content as this build\n", coordinate, releaseTestVersion)
			}
			wantOutput += fmt.Sprintf("%-36s %-12s %s\n", following, releaseTestVersion, "same content")
			wantSummary += fmt.Sprintf("- `%s@%s`: same content as this build\n", following, releaseTestVersion)
			if scenario.empty {
				wantOutput = ""
				wantSummary = "\n### Already published\n\nNo coordinate was already published.\n"
			}
			output, code := h.step(t, "Compare published packages")
			if code != wantExit || output != wantOutput {
				t.Fatalf("compare exit = %d, want %d; stdout = %q, want %q", code, wantExit, output, wantOutput)
			}
			if got := string(h.read(t, "runner/roundfix-release-complete")); got != wantComplete {
				t.Errorf("complete list = %q, want %q", got, wantComplete)
			}
			if got := string(h.read(t, "summary")); got != wantSummary {
				t.Errorf("summary = %q, want %q", got, wantSummary)
			}
			if got := h.read(t, "publish.log"); len(got) != 0 {
				t.Errorf("comparison published: %q", got)
			}
			entries, err := os.ReadDir(filepath.Join(h.root, "tmp"))
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("comparison left temporary files: %v", entries)
			}
		})
	}
}

func TestReleaseRerunPublishesOnlyTheMissingCoordinates(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct {
		name, version, failure, wantLog string
		published                       int
		launcher                        bool
	}{
		{name: "v0.132.1 five platforms", version: releaseTestVersion, published: 5,
			wantLog: "roundfix token=\n"},
		{name: "v0.130.0 one platform", version: "0.130.0", published: 1,
			wantLog: "cli-darwin-x64 token=\ncli-linux-arm64 token=\ncli-linux-x64 token=\ncli-win32-x64 token=\nroundfix token=\n"},
		{name: "different content stops before publishing", version: releaseTestVersion, published: 2, failure: "different"},
		{name: "missing tarball stops before publishing", version: releaseTestVersion, published: 2, failure: "missing"},
		{name: "published launcher is skipped last", version: releaseTestVersion, launcher: true,
			wantLog: "cli-darwin-arm64 token=\ncli-darwin-x64 token=\ncli-linux-arm64 token=\ncli-linux-x64 token=\ncli-win32-x64 token=\n"},
		{name: "all eligible publish in order", version: releaseTestVersion,
			wantLog: "cli-darwin-arm64 token=\ncli-darwin-x64 token=\ncli-linux-arm64 token=\ncli-linux-x64 token=\ncli-win32-x64 token=\nroundfix token=\n"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			h := newReleaseWorkflowHarness(t, releaseWorkflowCase{
				targetVersion: "v" + scenario.version, refName: "v" + scenario.version, eventName: "push", fallbackWindow: "1",
			})
			var wantSkips strings.Builder
			for i, coordinate := range h.coordinates {
				body := `{"versions":{},"time":{}}`
				if i < scenario.published || scenario.launcher && coordinate == "roundfix" {
					body = fmt.Sprintf(`{"versions":{%q:{}},"time":{}}`, scenario.version)
					registryBody := "build for " + coordinate
					if scenario.failure == "different" && i == 0 {
						registryBody += " changed"
					}
					releaseCompareFixtures(t, h, coordinate, scenario.version, registryBody)
					if scenario.failure == "missing" && i == 0 {
						if err := os.Remove(filepath.Join(h.fixtures, "registry", coordinate+"@"+scenario.version, "package.tgz")); err != nil {
							t.Fatal(err)
						}
					}
					fmt.Fprintf(&wantSkips, "skipping %s@%s: already published with the same content\n", coordinate, scenario.version)
				}
				h.packument(t, coordinate, body)
			}
			if output, code := h.step(t, "Publication preflight"); code != 0 {
				t.Fatalf("preflight exit = %d\n%s", code, output)
			}
			output, code := h.step(t, "Compare published packages")
			if scenario.failure != "" {
				if code != 1 {
					t.Fatalf("compare exit = %d, want 1\n%s", code, output)
				}
				if got := h.read(t, "publish.log"); len(got) != 0 {
					t.Errorf("failed comparison reached publish: %q", got)
				}
				return
			}
			if code != 0 {
				t.Fatalf("compare exit = %d\n%s", code, output)
			}
			// Substrings and regex patterns in the complete file must not skip a coordinate.
			complete := string(h.read(t, "runner/roundfix-release-complete"))
			complete += "prefix-roundfix-suffix\n@roundfix/cli-.*\n"
			writeFile(t, h.runner, "roundfix-release-complete", complete)
			output, code = h.step(t, "Publish to npm")
			if code != 0 || output != wantSkips.String() {
				t.Fatalf("publish exit = %d, stdout = %q, want %q", code, output, wantSkips.String())
			}
			if got := string(h.read(t, "publish.log")); got != scenario.wantLog {
				t.Errorf("publish log = %q, want %q", got, scenario.wantLog)
			}
			if got := h.read(t, "runner/roundfix-npm-token-fallback"); len(got) != 0 {
				t.Errorf("rerun used token fallback: %q", got)
			}
		})
	}
}

func TestReleaseComparesPublishedPackagesAfterTheBuildAndBeforePublishing(t *testing.T) {
	t.Parallel()
	contents, err := os.ReadFile(filepath.Join(findRepositoryRoot(t), ".github/workflows/release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var workflow struct {
		Jobs map[string]struct {
			Steps []struct {
				Name string `yaml:"name"`
				If   string `yaml:"if"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(contents, &workflow); err != nil {
		t.Fatal(err)
	}
	steps := workflow.Jobs["release"].Steps
	count := 0
	for i, step := range steps {
		if step.Name != "Compare published packages" {
			continue
		}
		count++
		if step.If != "${{ github.event_name == 'push' }}" {
			t.Errorf("compare condition = %q, want push only", step.If)
		}
		if i == 0 || steps[i-1].Name != "Cross-compile and stage" || i+1 >= len(steps) || steps[i+1].Name != "Publish to npm" {
			t.Errorf("compare at index %d must be directly between build and publish", i)
		}
	}
	if count != 1 {
		t.Errorf("compare step count = %d, want 1", count)
	}
}

// releaseCompareFixtures uses the harness's archive/tar and compress/gzip builder.
func releaseCompareFixtures(t *testing.T, h *releaseWorkflowHarness, coordinate, version, registryBody string) (string, string) {
	t.Helper()
	localBody := "build for " + coordinate
	local, localIntegrity := releaseTarball(t, localBody, gzip.BestSpeed)
	registry, registryIntegrity := releaseTarball(t, registryBody, gzip.BestCompression)
	if bytes.Equal(local, registry) || localIntegrity == registryIntegrity {
		t.Fatal("fixture must have different gzip bytes and integrity strings")
	}
	directory := strings.TrimPrefix(coordinate, "@roundfix/")
	writeFile(t, h.fixtures, "local/"+directory+"/package.tgz", string(local))
	writeFile(t, h.fixtures, "local/"+directory+"/integrity", localIntegrity)
	writeFile(t, h.fixtures, "registry/"+coordinate+"@"+version+"/package.tgz", string(registry))
	writeFile(t, h.fixtures, "registry/"+coordinate+"@"+version+"/integrity", registryIntegrity)
	return localIntegrity, registryIntegrity
}
