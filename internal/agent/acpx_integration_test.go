package agent

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const spawnedFixtureDeathTestEnv = "ROUNDFIX_TEST_SPAWNED_FIXTURE_DEATH"

// spawnedFixtureWaitPrefix marks the wait duration the nested fixture run
// reports, so the parent measures the wait rather than the nested `go test`
// process that contains it. A loaded runner can spend more than a second
// starting and tearing that process down while the wait itself stays in
// milliseconds, which is the only number this test claims anything about.
const spawnedFixtureWaitPrefix = "roundfix-fixture-wait="

func spawnedFixtureWaitDuration(output string) (time.Duration, bool) {
	for _, line := range strings.Split(output, "\n") {
		_, value, found := strings.Cut(strings.TrimSpace(line), spawnedFixtureWaitPrefix)
		if !found {
			continue
		}
		measured, err := time.ParseDuration(strings.TrimSpace(value))
		if err != nil {
			continue
		}
		return measured, true
	}
	return 0, false
}

// spawnedFixtureDeathBudget bounds how long the wait may take once the fixture
// has died. It is a tenth of the wait's own budget, so the assertion measures
// fast failure rather than machine speed.
const spawnedFixtureDeathBudget = agentWaitBudget / 10

func TestSpawnedFixtureDeathFailsFast(t *testing.T) {
	t.Parallel()

	const fixtureName = "immediate-exit acpx fixture"
	if os.Getenv(spawnedFixtureDeathTestEnv) == "fixture" {
		fmt.Println("fixture started")
		os.Exit(23)
	}
	if os.Getenv(spawnedFixtureDeathTestEnv) == "wait" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestSpawnedFixtureDeathFailsFast$")
		cmd.Env = environmentForTest(spawnedFixtureDeathTestEnv + "=fixture")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			t.Fatalf("read %s output: %v", fixtureName, err)
		}
		if err := cmd.Start(); err != nil {
			t.Fatalf("start %s: %v", fixtureName, err)
		}
		// Exclude fixture startup from the asserted death-detection span.
		if _, err := bufio.NewReader(stdout).ReadString('\n'); err != nil {
			waitErr := cmd.Wait()
			t.Fatalf("read %s first output: %v (exit: %v)", fixtureName, err, waitErr)
		}
		waitStarted := time.Now()
		fixtureExit := make(chan error, 1)
		go func() {
			fixtureExit <- cmd.Wait()
		}()
		t.Cleanup(func() {
			fmt.Printf("%s%s\n", spawnedFixtureWaitPrefix, time.Since(waitStarted))
		})
		waitForFile(t, fixtureName, filepath.Join(t.TempDir(), "never-created"), fixtureExit)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestSpawnedFixtureDeathFailsFast$")
	cmd.Env = environmentForTest(spawnedFixtureDeathTestEnv + "=wait")
	output, err := cmd.CombinedOutput()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("dead fixture wait exceeded its containment timeout:\n%s", output)
	}
	if err == nil {
		t.Fatalf("dead fixture wait unexpectedly passed:\n%s", output)
	}
	elapsed, ok := spawnedFixtureWaitDuration(string(output))
	if !ok {
		t.Fatalf("dead fixture wait did not report its measured duration:\n%s", output)
	}
	// The claim is that the wait notices the fixture's death instead of sitting
	// out its budget, so the bound is a fraction of that budget rather than a
	// fixed second: a loaded machine can spend more than a second and still
	// prove the point, and an absolute bound turns load into a failure.
	if elapsed >= spawnedFixtureDeathBudget {
		t.Fatalf("dead fixture wait took %s, want under %s (a tenth of the %s wait budget):\n%s", elapsed, spawnedFixtureDeathBudget, agentWaitBudget, output)
	}
	if text := string(output); !strings.Contains(text, fixtureName) || !strings.Contains(text, "exit status 23") {
		t.Fatalf("dead fixture failure must name %q and exit status 23:\n%s", fixtureName, output)
	}
}
