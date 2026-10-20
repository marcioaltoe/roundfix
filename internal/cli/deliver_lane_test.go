// Boundary: disposable Git repositories and compiled Docker/Go stand-ins.
package cli

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
)

func newDeliverLaneFixture(t *testing.T) *linuxLaneFixture {
	t.Helper()
	docker, goBinary := linuxLaneTools(t)
	f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{
		"go.mod": "module fixture\n\ngo 1.26.1\n", "pkg/code.go": "package pkg\n",
		".roundfixrc.yml": "delivery:\n  linux_lane: false\n",
	})
	return f
}

func invokeDeliverLane(t *testing.T, f *linuxLaneFixture, args ...string) (string, string, int) {
	t.Helper()
	return invokeDeliverLaneContext(t, f.context(t.Context()), f, args...)
}

func invokeDeliverLaneContext(t *testing.T, ctx context.Context, f *linuxLaneFixture, args ...string) (string, string, int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	environment := commandEnvironment{homeDir: f.workflow.loaded.HomeDir, workDir: f.repo, environ: f.environment, dependencies: commandDependenciesForContext(ctx)}
	code := runWithContext(ctx, append([]string{"deliver", "lane"}, args...), &stdout, &stderr, environment)
	return stdout.String(), stderr.String(), code
}

func deliverLaneLog(f *linuxLaneFixture) string {
	artifactDir, err := roundconfig.ResolveArtifactDirectory("", f.repo, f.workflow.loaded.HomeDir)
	if err != nil {
		panic(err)
	}
	return filepath.Join(artifactDir, "delivery", "linux-lane.log")
}

func assertDeliverLaneTranscript(t *testing.T, f *linuxLaneFixture, line string, wantCode int) {
	t.Helper()
	stdout, stderr, code := invokeDeliverLane(t, f, "--all")
	if stdout != line+"\n" || stderr != "" || code != wantCode {
		t.Fatalf("stdout=%q stderr=%q exit=%d; want stdout=%q stderr empty exit=%d", stdout, stderr, code, line+"\n", wantCode)
	}
}

func TestDeliverLanePassedTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	assertDeliverLaneTranscript(t, f, "roundfix: Linux Lane passed 1 package(s) on golang:1.26.1", 0)
}

func TestDeliverLaneChangeFailureTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	f.environment = withEnvValue(f.environment, "LANE_EXIT", "1")
	assertDeliverLaneTranscript(t, f, "roundfix: Linux Lane change failure: go test exited 1; see "+deliverLaneLog(f), 1)
}

func TestDeliverLaneContainerFailureTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	f.environment = withEnvValue(f.environment, "LANE_EXIT", "3")
	f.environment = withEnvValue(f.environment, "LANE_OUTPUT", "roundfix-linux-lane: container: missing tool git\n")
	assertDeliverLaneTranscript(t, f, "roundfix: Linux Lane container failure: missing tool git; see "+deliverLaneLog(f), 3)
}

func TestDeliverLaneSkipTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	_, goBinary := linuxLaneTools(t)
	f.environment = withEnvValue(f.environment, "PATH", filepath.Dir(goBinary))
	assertDeliverLaneTranscript(t, f, "roundfix: notice: Linux Lane skipped: docker is not on PATH", 3)
}

func TestDeliverLaneSelectsPackagesAndImage(t *testing.T) {
	t.Parallel()
	t.Run("help", func(t *testing.T) {
		f := newDeliverLaneFixture(t)
		stdout, stderr, code := invokeDeliverLane(t, f, "--help")
		if code != 0 || stderr != "" {
			t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
		}
		for _, line := range []string{
			"  roundfix deliver lane [--all] [<revision>]\n",
			"  lane    Run the Linux Lane on demand and exit with its outcome\n",
			"  --all  Test every package in the Linux Lane\n",
		} {
			if !strings.Contains(stdout, line) {
				t.Fatalf("help missing %q: %s", line, stdout)
			}
		}
	})
	t.Run("all without default branch and revision archive", func(t *testing.T) {
		f := newDeliverLaneFixture(t)
		const image = "registry.example.com/ci/go:custom"
		f.commit(t, ".roundfixrc.yml", "delivery:\n  linux_lane: false\n  linux_lane_image: "+image+"\n")
		revision := f.head
		f.commit(t, "pkg/code.go", "package pkg\n// later commit\n")
		gittest.Run(t, f.repo, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
		gittest.Run(t, f.repo, "update-ref", "-d", "refs/remotes/origin/main")
		gittest.Run(t, f.repo, "branch", "-m", "topic")
		stdout, stderr, code := invokeDeliverLane(t, f, "--all", revision)
		if code != 0 || stderr != "" || stdout != "roundfix: Linux Lane passed 1 package(s) on "+image+"\n" {
			t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
		}
		starts := f.starts(t)
		if len(starts) != 4 || !reflect.DeepEqual(starts[1].Args, []string{"image", "inspect", image}) {
			t.Fatalf("starts=%+v", starts)
		}
		args := starts[3].Args
		if !containsExact(args, "--init") || !containsExact(args, image) || args[len(args)-1] != "./..." || !containsExact(args, fmt.Sprintf("roundfix-linux-lane-%d", os.Getpid())) {
			t.Fatalf("run args=%q", args)
		}
		if log := mustRead(t, deliverLaneLog(f)); !strings.Contains(log, "docker 'run'") {
			t.Fatalf("log=%q", log)
		}
		archive, err := os.ReadFile(f.archive)
		if err != nil {
			t.Fatal(err)
		}
		reader := tar.NewReader(bytes.NewReader(archive))
		found := false
		for {
			header, err := reader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			if header.Name == "pkg/code.go" {
				content, err := io.ReadAll(reader)
				if err != nil {
					t.Fatal(err)
				}
				if string(content) != "package pkg\n" {
					t.Fatalf("archived later revision: %q", content)
				}
				found = true
			}
		}
		if !found {
			t.Fatal("revision's code absent from archive")
		}
	})
	t.Run("changed packages and checkout config over user config", func(t *testing.T) {
		f := newDeliverLaneFixture(t)
		f.commit(t, "pkg/new.go", "package pkg\n")
		mustMkdir(t, filepath.Join(f.workflow.loaded.HomeDir, ".roundfix"))
		mustWrite(t, filepath.Join(f.workflow.loaded.HomeDir, ".roundfix", "config.yml"), "delivery:\n  linux_lane_image: user/image:ignored\n")
		mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: false\n  linux_lane_image: checkout/image:selected\n")
		stdout, stderr, code := invokeDeliverLane(t, f)
		if code != 0 || stderr != "" || stdout != "roundfix: Linux Lane passed 1 package(s) on checkout/image:selected\n" {
			t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
		}
		starts := f.starts(t)
		args := starts[len(starts)-1].Args
		if args[len(args)-1] != "./pkg" {
			t.Fatalf("packages=%q", args)
		}
	})
	for _, tc := range []struct {
		name                          string
		args                          []string
		invalidConfig, unknownDefault bool
	}{
		{name: "extra argument", args: []string{"--all", "HEAD", "HEAD"}},
		{name: "not a commit", args: []string{"--all", "HEAD^{tree}"}},
		{name: "missing revision", args: []string{"--all", "missing"}},
		{name: "unknown flag", args: []string{"--unknown"}},
		{name: "invalid config", args: []string{"--all"}, invalidConfig: true},
		{name: "unknown default", unknownDefault: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newDeliverLaneFixture(t)
			if tc.invalidConfig {
				mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery: [invalid\n")
			}
			if tc.unknownDefault {
				gittest.Run(t, f.repo, "symbolic-ref", "--delete", "refs/remotes/origin/HEAD")
				gittest.Run(t, f.repo, "update-ref", "-d", "refs/remotes/origin/main")
				gittest.Run(t, f.repo, "branch", "-m", "topic")
			}
			stdout, stderr, code := invokeDeliverLane(t, f, tc.args...)
			if code != 2 || stdout != "" || !strings.HasPrefix(stderr, "Preflight failed\n\nReason:\n") || !strings.Contains(stderr, "Run 'roundfix deliver lane --help' for usage.") {
				t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
			}
			if starts := f.starts(t); len(starts) != 0 {
				t.Fatalf("preflight started tools: %+v", starts)
			}
			if _, err := os.Stat(deliverLaneLog(f)); !os.IsNotExist(err) {
				t.Fatalf("preflight wrote log: %v", err)
			}
		})
	}
}

func TestDeliverLaneCancellationRemovesContainer(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	f.environment = withEnvValue(f.environment, "LANE_READY", listener.Addr().String())
	ctx, cancel := context.WithCancel(f.context(t.Context()))
	ready, accepted, finished := make(chan struct{}), make(chan struct{}), make(chan struct{})
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	var stdout, stderr string
	var code int
	go func() {
		defer close(accepted)
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		close(ready)
		<-ctx.Done()
	}()
	go func() {
		defer close(finished)
		stdout, stderr, code = invokeDeliverLaneContext(t, ctx, f, "--all")
	}()
	t.Cleanup(func() {
		cancel()
		stop()
		_ = listener.Close()
		<-accepted
		<-finished
	})
	select {
	case <-ready:
		cancel()
	case <-finished:
		t.Fatalf("lane returned before the container started: exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	case <-t.Context().Done():
		t.Fatal(t.Context().Err())
	}
	<-finished
	if code != 130 || stdout != "" || stderr != "" {
		t.Fatalf("stdout=%q stderr=%q exit=%d", stdout, stderr, code)
	}
	starts := f.starts(t)
	want := []string{"rm", "-f", fmt.Sprintf("roundfix-linux-lane-%d", os.Getpid())}
	if len(starts) != 5 || !reflect.DeepEqual(starts[4].Args, want) {
		t.Fatalf("starts=%+v want cleanup=%q", starts, want)
	}
}

func TestDeliverLaneObjectStoreSkipTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	objects := linuxLaneFixtureObjects(t, f.repo)
	other := filepath.Join(t.TempDir(), "other")
	gittest.InitRepo(t, other, "-b", "main")
	mustWrite(t, filepath.Join(objects, "info", "alternates"), linuxLaneFixtureObjects(t, other)+"\n")
	assertDeliverLaneTranscript(t, f, "roundfix: notice: Linux Lane skipped: the object store "+objects+" borrows objects through info/alternates, which the Linux Lane does not mount", 3)
	if starts := f.starts(t); len(starts) != 0 {
		t.Fatalf("skip started tools: %+v", starts)
	}
}

func TestDeliverLanePackageCountTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	f.environment = withEnvValue(f.environment, "LANE_PACKAGES", "43")
	assertDeliverLaneTranscript(t, f, "roundfix: Linux Lane passed 43 package(s) on golang:1.26.1", 0)
}

func TestDeliverLaneNoPackageCountTranscript(t *testing.T) {
	t.Parallel()
	f := newDeliverLaneFixture(t)
	f.environment = withEnvValue(f.environment, "LANE_PACKAGES", "")
	assertDeliverLaneTranscript(t, f, "roundfix: Linux Lane passed on golang:1.26.1; the lane script reported no package count", 0)
}
