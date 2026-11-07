// Suite: release workflow step scripts.
// Invariant: preflight distinguishes published content from a used version; non-identity failures never retry with a token.
// Boundary IN: release workflow YAML, copied npm manifests, fixture packuments and tarballs, and isolated command stubs.
// Boundary OUT: network access, real npm publication, builds, GitHub Releases, and process signaling.
package verifyselect_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"roundfix/internal/testfixture"
)

const releaseTestVersion = "0.132.1"

type releaseWorkflowCase struct {
	targetVersion, refName, eventName, fallbackWindow string
	failPublishDirectory, failPublishOutput           string
}

type releaseWorkflowHarness struct {
	root, runner, summary, fixtures, bash string
	env                                   []string
	coordinates                           []string
}

func newReleaseWorkflowHarness(t *testing.T, scenario releaseWorkflowCase) *releaseWorkflowHarness {
	t.Helper()
	bash := testfixture.RequireTool(t, "bash")
	tools := map[string]string{
		"jq":   testfixture.RequireTool(t, "jq"),
		"gzip": testfixture.RequireTool(t, "gzip"),
		"cmp":  testfixture.RequireTool(t, "cmp"),
	}
	root, err := filepath.Abs(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	h := &releaseWorkflowHarness{
		root: root, runner: filepath.Join(root, "runner"), summary: filepath.Join(root, "summary"),
		fixtures: filepath.Join(root, "fixtures"), bash: bash,
	}
	for _, directory := range []string{"bin", "home", "tmp", "runner", "fixtures/packuments", "fixtures/local", "fixtures/registry"} {
		if err := os.MkdirAll(filepath.Join(root, directory), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	repository := findRepositoryRoot(t)
	manifests, err := filepath.Glob(filepath.Join(repository, "dist/npm/packages/*/package.json"))
	if err != nil {
		t.Fatal(err)
	}
	manifests = append(manifests, filepath.Join(repository, "dist/npm/platforms.json"), filepath.Join(repository, "dist/npm/roundfix/package.json"))
	for _, source := range manifests {
		contents, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		relative, err := filepath.Rel(repository, source)
		if err != nil {
			t.Fatal(err)
		}
		writeFile(t, root, relative, string(contents))
	}
	var platforms []struct {
		Package string `json:"package"`
	}
	if err := json.Unmarshal(h.read(t, "dist/npm/platforms.json"), &platforms); err != nil {
		t.Fatal(err)
	}
	for _, platform := range platforms {
		h.coordinates = append(h.coordinates, platform.Package)
	}
	h.coordinates = append(h.coordinates, "roundfix")
	for name, path := range tools {
		if err := os.Symlink(path, filepath.Join(root, "bin", name)); err != nil {
			t.Fatal(err)
		}
	}
	for name, script := range map[string]string{"curl": releaseCurlStub, "npm": releaseNPMStub} {
		writeFile(t, root, "bin/"+name, script)
		if err := os.Chmod(filepath.Join(root, "bin", name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, root, "summary", "")
	writeFile(t, root, "publish.log", "")
	writeFile(t, root, "fixtures/publish-output", scenario.failPublishOutput)
	h.env = []string{
		"PATH=" + filepath.Join(root, "bin") + string(os.PathListSeparator) + os.Getenv("PATH"),
		"HOME=" + filepath.Join(root, "home"), "TMPDIR=" + filepath.Join(root, "tmp"),
		"RUNNER_TEMP=" + h.runner, "GITHUB_STEP_SUMMARY=" + h.summary,
		"TARGET_VERSION=" + scenario.targetVersion, "GITHUB_REF_NAME=" + scenario.refName,
		"GITHUB_EVENT_NAME=" + scenario.eventName, "FALLBACK_WINDOW=" + scenario.fallbackWindow,
		"NPM_FALLBACK_TOKEN=" + "fake-" + "release-fixture",
		"RELEASE_FIXTURES=" + h.fixtures, "RELEASE_PUBLISH_LOG=" + filepath.Join(root, "publish.log"),
		"RELEASE_FAIL_PUBLISH=" + scenario.failPublishDirectory,
	}
	return h
}

func (h *releaseWorkflowHarness) read(t *testing.T, relative string) []byte {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(h.root, relative))
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func (h *releaseWorkflowHarness) packument(t *testing.T, coordinate, body string) {
	t.Helper()
	writeFile(t, h.fixtures, "packuments/"+url.QueryEscape(coordinate)+".json", body)
}

func (h *releaseWorkflowHarness) step(t *testing.T, name string) (string, int) {
	t.Helper()
	contents, err := os.ReadFile(filepath.Join(findRepositoryRoot(t), ".github/workflows/release.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var workflow struct {
		Jobs map[string]struct {
			Steps []struct {
				Name string `yaml:"name"`
				Run  string `yaml:"run"`
			} `yaml:"steps"`
		} `yaml:"jobs"`
	}
	if err := yaml.Unmarshal(contents, &workflow); err != nil {
		t.Fatal(err)
	}
	for _, step := range workflow.Jobs["release"].Steps {
		if step.Name == name {
			return h.script(t, step.Run)
		}
	}
	t.Fatalf("release step %q not found", name)
	return "", -1
}

func (h *releaseWorkflowHarness) script(t *testing.T, script string) (string, int) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "step.sh")
	if err := os.WriteFile(path, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	command := exec.CommandContext(t.Context(), h.bash, "-e", path)
	command.Dir = h.root
	command.Env = h.env
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if err != nil {
		exit, ok := err.(*exec.ExitError)
		if !ok || t.Context().Err() != nil {
			t.Fatalf("step execution: %v\n%s\n%s", err, &stdout, &stderr)
		}
		if stderr.Len() > 0 {
			t.Logf("step stderr: %s", &stderr)
		}
		return stdout.String(), exit.ExitCode()
	}
	if stderr.Len() != 0 {
		t.Fatalf("unexpected step stderr: %s", &stderr)
	}
	return stdout.String(), 0
}

const releaseCurlStub = `#!/bin/sh
set -eu
output=
url=
while [ "$#" -gt 0 ]; do
  case "$1" in
    --output) output="$2"; shift ;;
    https://registry.npmjs.org/*) url="$1" ;;
  esac
  shift
done
fixture="$RELEASE_FIXTURES/packuments/${url##*/}.json"
if [ -f "$fixture" ]; then
  cat "$fixture" > "$output"
  printf 200
else
  printf 404
fi
`

const releaseNPMStub = `#!/bin/sh
set -eu
operation="$1"
shift
case "$operation" in
  pack)
    destination="$PWD"
    fixture="$RELEASE_FIXTURES/local/${PWD##*/}"
    while [ "$#" -gt 0 ]; do
      case "$1" in
        --pack-destination) destination="$2"; shift ;;
        --json) ;;
        *) fixture="$RELEASE_FIXTURES/registry/$1" ;;
      esac
      shift
    done
    [ -f "$fixture/package.tgz" ] && [ -f "$fixture/integrity" ] || exit 1
    mkdir -p "$destination"
    cp "$fixture/package.tgz" "$destination/package.tgz"
    integrity="$(cat "$fixture/integrity")"
    printf '[{"filename":"package.tgz","integrity":"%s"}]\n' "$integrity"
    ;;
  publish)
    directory="${PWD##*/}"
    printf '%s token=%s\n' "$directory" "${NODE_AUTH_TOKEN+set}" >> "$RELEASE_PUBLISH_LOG"
    if [ "$directory" = "$RELEASE_FAIL_PUBLISH" ]; then
      cat "$RELEASE_FIXTURES/publish-output"
      exit 1
    fi
    ;;
  *) printf 'unexpected npm operation: %s\n' "$operation" >&2; exit 2 ;;
esac
`

func TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct {
		name          string
		published     int
		launcherState string
		wantExit      int
		wantError     string
	}{
		{name: "all eligible"},
		{name: "five published platforms", published: 5},
		{name: "launcher version only under time", launcherState: "used", wantExit: 1,
			wantError: "::error::registry: roundfix@0.132.1 is already used\n"},
		{name: "all published", published: 6, wantExit: 1,
			wantError: "::error::registry: every coordinate of the Release Set is already published at 0.132.1\n"},
		{name: "published platform does not hide cooldown", published: 1, launcherState: "cooldown", wantExit: 1,
			wantError: "::error::registry: roundfix@0.132.1 is in cooldown until 2026-11-07T00:00:00Z\n"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			for _, event := range []string{"push", "workflow_dispatch"} {
				t.Run(event, func(t *testing.T) {
					t.Parallel()
					h := newReleaseWorkflowHarness(t, releaseWorkflowCase{
						targetVersion: "v" + releaseTestVersion, refName: "v" + releaseTestVersion, eventName: event,
					})
					// A stale list from an earlier step invocation must be truncated.
					writeFile(t, h.runner, "roundfix-release-published", "stale-coordinate\n")
					var wantOutput, wantSummary strings.Builder
					fmt.Fprintf(&wantOutput, "%-36s %-12s %s\n", "PACKAGE", "VERSION", "STATE")
					wantSummary.WriteString("### Publication preflight\n\n| Package | Version | State |\n| --- | --- | --- |\n")
					var published []string
					for i, coordinate := range h.coordinates {
						state, body := "eligible", `{"versions":{},"time":{}}`
						if i < scenario.published {
							state, body = "published", `{"versions":{"0.132.1":{}},"time":{"0.132.1":"2026-11-05T00:00:00Z","unpublished":"2026-11-06T00:00:00Z"}}`
							published = append(published, coordinate)
						} else if coordinate == "roundfix" && scenario.launcherState != "" {
							state = scenario.launcherState
							if state == "used" {
								body = `{"versions":{},"time":{"0.132.1":"2026-11-05T00:00:00Z","unpublished":"2026-11-06T00:00:00Z"}}`
							} else {
								body = `{"versions":{},"time":{"unpublished":{"time":"2026-11-06T00:00:00Z"}}}`
							}
						}
						h.packument(t, coordinate, body)
						fmt.Fprintf(&wantOutput, "%-36s %-12s %s\n", coordinate, releaseTestVersion, state)
						fmt.Fprintf(&wantSummary, "| `%s` | `%s` | `%s` |\n", coordinate, releaseTestVersion, state)
						if state == "published" {
							fmt.Fprintf(&wantOutput, "::notice::registry: %s@%s is already published; a tag run compares it with its own build\n", coordinate, releaseTestVersion)
						}
					}
					wantOutput.WriteString(scenario.wantError)
					output, code := h.step(t, "Publication preflight")
					if code != scenario.wantExit || output != wantOutput.String() {
						t.Fatalf("preflight exit = %d, want %d\nstdout:\n%swant:\n%s", code, scenario.wantExit, output, &wantOutput)
					}
					wantPublished := ""
					if len(published) > 0 {
						wantPublished = strings.Join(published, "\n") + "\n"
					}
					if got := string(h.read(t, "runner/roundfix-release-published")); got != wantPublished {
						t.Errorf("published list = %q, want %q", got, wantPublished)
					}
					if got := string(h.read(t, "summary")); got != wantSummary.String() {
						t.Errorf("summary = %q, want %q", got, wantSummary.String())
					}
				})
			}
		})
	}
}

func TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry(t *testing.T) {
	t.Parallel()
	const failure = "npm error code TLOG_CREATE_ENTRY_ERROR\nnpm error error creating tlog entry - (409) an equivalent entry already exists in the transparency log\n"
	h := newReleaseWorkflowHarness(t, releaseWorkflowCase{
		targetVersion: "v" + releaseTestVersion, refName: "v" + releaseTestVersion, eventName: "push", fallbackWindow: "1",
		failPublishDirectory: "roundfix", failPublishOutput: failure,
	})
	for _, coordinate := range h.coordinates {
		h.packument(t, coordinate, `{"versions":{},"time":{}}`)
	}
	if output, code := h.step(t, "Publication preflight"); code != 0 {
		t.Fatalf("eligible preflight exit = %d\n%s", code, output)
	}
	output, code := h.step(t, "Publish to npm")
	want := "::error::publish: roundfix failed for a non-identity reason\n" + failure
	if code != 1 || output != want {
		t.Fatalf("publish exit = %d, want 1; stdout = %q, want %q", code, output, want)
	}
	var wantLog strings.Builder
	for _, coordinate := range h.coordinates {
		fmt.Fprintf(&wantLog, "%s token=\n", strings.TrimPrefix(coordinate, "@roundfix/"))
	}
	if got := string(h.read(t, "publish.log")); got != wantLog.String() {
		t.Errorf("publish log = %q, want %q", got, wantLog.String())
	}
	if got := h.read(t, "runner/roundfix-npm-token-fallback"); len(got) != 0 {
		t.Errorf("non-identity failure used token fallback: %q", got)
	}
}

// releaseTarball builds valid fixture content without invoking npm or a build tool.
func releaseTarball(t *testing.T, body string, compression int) ([]byte, string) {
	t.Helper()
	var output bytes.Buffer
	compressed, err := gzip.NewWriterLevel(&output, compression)
	if err != nil {
		t.Fatal(err)
	}
	archive := tar.NewWriter(compressed)
	if err := archive.WriteHeader(&tar.Header{Name: "package/content", Mode: 0o644, Size: int64(len(body))}); err != nil {
		t.Fatal(err)
	}
	if _, err := archive.Write([]byte(body)); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := compressed.Close(); err != nil {
		t.Fatal(err)
	}
	digest := sha512.Sum512(output.Bytes())
	return output.Bytes(), "sha512-" + base64.StdEncoding.EncodeToString(digest[:])
}

func TestReleaseWorkflowHarnessPacksOnlyFixtureTarballs(t *testing.T) {
	t.Parallel()
	for _, scenario := range []struct {
		name, source, arguments, destination string
		missing                              bool
	}{
		{name: "local default destination", source: "local/roundfix", destination: "dist/npm/roundfix/package.tgz"},
		{name: "local explicit destination", source: "local/roundfix", arguments: "--pack-destination \"$RUNNER_TEMP/packed\"", destination: "runner/packed/package.tgz"},
		{name: "registry default destination", source: "registry/@roundfix/cli-linux-x64@" + releaseTestVersion, arguments: "@roundfix/cli-linux-x64@" + releaseTestVersion, destination: "dist/npm/roundfix/package.tgz"},
		{name: "registry explicit destination", source: "registry/roundfix@" + releaseTestVersion, arguments: "roundfix@" + releaseTestVersion + " --pack-destination \"$RUNNER_TEMP/packed\"", destination: "runner/packed/package.tgz"},
		{name: "missing local fixture", missing: true},
		{name: "missing registry fixture", arguments: "roundfix@" + releaseTestVersion, missing: true},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			h := newReleaseWorkflowHarness(t, releaseWorkflowCase{
				targetVersion: "v" + releaseTestVersion, refName: "v" + releaseTestVersion, eventName: "push",
			})
			tarball, integrity := releaseTarball(t, "fixture content", gzip.BestSpeed)
			if !scenario.missing {
				writeFile(t, h.fixtures, scenario.source+"/package.tgz", string(tarball))
				writeFile(t, h.fixtures, scenario.source+"/integrity", integrity)
			}
			output, code := h.script(t, "cd dist/npm/roundfix\nnpm pack --json "+scenario.arguments+"\n")
			if scenario.missing {
				if code != 1 || output != "" {
					t.Fatalf("missing fixture: exit = %d, stdout = %q", code, output)
				}
				return
			}
			want := fmt.Sprintf("[{\"filename\":\"package.tgz\",\"integrity\":\"%s\"}]\n", integrity)
			if code != 0 || output != want {
				t.Fatalf("pack: exit = %d, stdout = %q, want %q", code, output, want)
			}
			if got := h.read(t, scenario.destination); !bytes.Equal(got, tarball) {
				t.Fatal("pack did not copy the fixture tarball")
			}
			if log := h.read(t, "publish.log"); len(log) != 0 {
				t.Fatalf("pack reached publish: %q", log)
			}
		})
	}
}
