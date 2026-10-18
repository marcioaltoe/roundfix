package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const (
	detachTestChildSleepBeforeLiveness = "sleep-before-liveness"
	detachTestChildSlowRunCreation     = "slow-run-creation"
	detachTestChildRunCreated          = "run-created"
	detachTestChildInvalidLiveness     = "invalid-liveness"
	detachTestChildMalformedRunCreated = "malformed-run-created"
	detachTestChildExitTwoWithOutput   = "exit-two-with-output"
	detachTestChildExitOneSilently     = "exit-one-silently"
	detachTestChildIgnoreSentinel      = "ignore-sentinel"
	detachTestRunID                    = "run-detach-test"
	detachTestOwnerIdentityEnv         = "ROUNDFIX_DETACH_TEST_OWNER_IDENTITY"
	detachTestOwnerPIDEnv              = "ROUNDFIX_DETACH_TEST_OWNER_PID"
	detachTestPIDPathEnv               = "ROUNDFIX_DETACH_TEST_PID_PATH"
	detachTestSentinelPathEnv          = "ROUNDFIX_DETACH_TEST_SENTINEL_PATH"
	detachTestRunCreationDelayEnv      = "ROUNDFIX_DETACH_TEST_RUN_CREATION_DELAY"
)

func TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.detachTimeouts = detachPhaseTimeouts{
					liveness:    detachTestBound(t),
					runCreation: detachTestBound(t),
				}
			})
			artifactParent := t.TempDir()
			artifactDir := filepath.Join(artifactParent, "artifacts")
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			environment := commandEnvironmentForTest(t)

			code := runDetachedCommand(
				[]string{command, "--detach"},
				commandRequest{name: command, detach: true, artifactDir: artifactDir},
				roundconfig.Loaded{GitRoot: t.TempDir(), HomeDir: t.TempDir()},
				&stdout,
				&stderr,
				withEnvValue(environment.environ, detachTestChildModeEnv, detachTestChildExitTwoWithOutput),
				t.TempDir(),
				environment.dependencies.detachTimeouts,
			)

			if code != exitPreflight {
				t.Fatalf("expected child refusal exit %d, got %d stdout=%q stderr=%q", exitPreflight, code, stdout.String(), stderr.String())
			}
			if _, err := os.Stat(artifactDir); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("Artifact Directory exists after detached refusal: stat error = %v", err)
			}
			entries, err := os.ReadDir(artifactParent)
			if err != nil {
				t.Fatalf("read Artifact Directory parent after refusal: %v", err)
			}
			if len(entries) != 0 {
				t.Fatalf("detached refusal left filesystem entries: %v", entries)
			}
		})
	}
}

func runDetachTestChild(mode string) int {
	switch mode {
	case detachTestChildSleepBeforeLiveness:
		return waitForDetachTestSentinel()
	case detachTestChildSlowRunCreation:
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		return waitForDetachTestSentinel()
	case detachTestChildRunCreated:
		if !recordDetachTestSurvivor() {
			return exitRunFailed
		}
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		delay := 750 * time.Millisecond
		if value := os.Getenv(detachTestRunCreationDelayEnv); value != "" {
			var err error
			delay, err = time.ParseDuration(value)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse detached test Run creation delay: %v\n", err)
				return exitRunFailed
			}
		}
		time.Sleep(delay)
		consoleLog := filepath.Join(filepath.Dir(os.Getenv(detachConsoleTempEnv)), detachTestRunID, "console.log")
		if _, err := fmt.Fprintf(handshake, "%s\t%s\n", detachTestRunID, consoleLog); err != nil {
			return exitRunFailed
		}
		return waitForDetachTestSentinel()
	case detachTestChildInvalidLiveness:
		handshake, ok := openDetachTestHandshake()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		if _, err := handshake.Write([]byte{'x'}); err != nil {
			return exitRunFailed
		}
		return exitRunFailed
	case detachTestChildMalformedRunCreated:
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		if _, err := fmt.Fprintln(handshake, "malformed"); err != nil {
			return exitRunFailed
		}
		return exitRunFailed
	case detachTestChildExitTwoWithOutput:
		fmt.Fprintln(os.Stderr, "detached child wrote a pre-handshake failure")
		return exitPreflight
	case detachTestChildExitOneSilently:
		return exitRunFailed
	case detachTestChildIgnoreSentinel:
		ownerPID, err := strconv.Atoi(os.Getenv(detachTestOwnerPIDEnv))
		if err != nil || ownerPID <= 0 {
			return exitRunFailed
		}
		if !recordDetachTestSurvivor() {
			return exitRunFailed
		}
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			if !detachTestOwnerAlive(ownerPID) {
				return exitRunFailed
			}
			<-ticker.C
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown detach test child mode %q\n", mode)
		return exitRunFailed
	}
}

func recordDetachTestSurvivor() bool {
	pidPath := strings.TrimSpace(os.Getenv(detachTestPIDPathEnv))
	if pidPath == "" {
		fmt.Fprintln(os.Stderr, "detached test child pid path is empty")
		return false
	}
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "record detached test child pid: %v\n", err)
		return false
	}
	return true
}

func waitForDetachTestSentinel() int {
	ownerPID, err := strconv.Atoi(os.Getenv(detachTestOwnerPIDEnv))
	if err != nil || ownerPID <= 0 {
		return exitRunFailed
	}
	sentinelPath := strings.TrimSpace(os.Getenv(detachTestSentinelPathEnv))
	if sentinelPath == "" {
		fmt.Fprintln(os.Stderr, "detached test child sentinel path is empty")
		return exitRunFailed
	}
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if !detachTestOwnerAlive(ownerPID) {
			return exitRunFailed
		}
		if _, err := os.Stat(sentinelPath); err == nil {
			return exitOK
		} else if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "read detached test child sentinel: %v\n", err)
			return exitRunFailed
		}
		<-ticker.C
	}
}

func writeDetachTestLiveness() (*os.File, bool) {
	handshake, ok := openDetachTestHandshake()
	if !ok {
		return nil, false
	}
	if _, err := handshake.Write([]byte{detachLivenessMarker}); err != nil {
		fmt.Fprintf(os.Stderr, "write detach test liveness marker: %v\n", err)
		_ = handshake.Close()
		return nil, false
	}
	return handshake, true
}

func openDetachTestHandshake() (*os.File, bool) {
	fd, err := strconv.Atoi(strings.TrimSpace(os.Getenv(detachHandshakeFDEnv)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid detach test handshake fd: %v\n", err)
		return nil, false
	}
	handshake := os.NewFile(uintptr(fd), "roundfix-detach-test-handshake")
	if handshake == nil {
		fmt.Fprintf(os.Stderr, "open detach test handshake fd %d failed\n", fd)
		return nil, false
	}
	return handshake, true
}

func detachTestOwnerEnvironment(t *testing.T, environment []string) []string {
	t.Helper()
	identity, err := store.OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("read detach test owner identity: %v", err)
	}
	environment = withEnvValue(environment, detachTestOwnerPIDEnv, strconv.Itoa(os.Getpid()))
	return withEnvValue(environment, detachTestOwnerIdentityEnv, identity)
}

func detachTestOwnerAlive(pid int) bool {
	recorded := os.Getenv(detachTestOwnerIdentityEnv)
	if recorded == "" || !store.ProcessAlive(pid) {
		return false
	}
	identity, err := store.OwnerProcessIdentity(context.Background(), pid)
	return err == nil && identity == recorded
}
