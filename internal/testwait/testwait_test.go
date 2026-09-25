package testwait_test

// Suite: deadline-bound test waits
// Invariant: readiness wins, ended work is reported immediately, and only the test deadline bounds an unfinished wait.
// Boundary IN: the public testwait helpers and their observable channel, diagnostic, and timing behavior
// Boundary OUT: caller-owned work cancellation and process lifecycle

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"roundfix/internal/testwait"
)

type recordingTB struct {
	testing.TB
	deadline     time.Time
	hasDeadline  bool
	fatalMessage string
}

func (tb *recordingTB) Deadline() (time.Time, bool) {
	return tb.deadline, tb.hasDeadline
}

func (tb *recordingTB) Fatalf(format string, args ...any) {
	tb.fatalMessage = fmt.Sprintf(format, args...)
}

func blockForStackDump(started chan<- struct{}, release <-chan struct{}) {
	close(started)
	<-release
}

func TestBoundFollowsTheTestDeadline(t *testing.T) {
	t.Parallel()
	deadline := time.Now().Add(10 * time.Minute)
	tb := &recordingTB{TB: t, deadline: deadline, hasDeadline: true}
	want := 10*time.Minute - testwait.Margin

	got := testwait.Bound(tb)

	if difference := got - want; difference < -time.Second || difference > time.Second {
		t.Fatalf("Bound() = %s, want %s within one second", got, want)
	}
}

func TestBoundFallsBackWithoutATestDeadline(t *testing.T) {
	t.Parallel()
	tb := &recordingTB{TB: t}

	got := testwait.Bound(tb)

	if want := 10 * time.Minute; got != want {
		t.Fatalf("Bound() = %s, want %s", got, want)
	}
}

func TestUntilReturnsTheReadyValue(t *testing.T) {
	t.Parallel()
	ready := make(chan string, 2)
	ready <- "first ready value"
	ready <- "second ready value"
	var ended <-chan struct{}

	got := testwait.Until(t, "ready value", ready, ended)

	if want := "first ready value"; got != want {
		t.Fatalf("Until() = %q, want %q", got, want)
	}
}

func TestUntilPrefersTheReadyValueOverAnEndedWork(t *testing.T) {
	t.Parallel()
	ready := make(chan string, 1)
	ended := make(chan string, 1)
	ready <- "ready value"
	ended <- "ended value"
	tb := &recordingTB{TB: t}

	got := testwait.Until(tb, "simultaneous result", ready, ended)

	if want := "ready value"; got != want {
		t.Fatalf("Until() = %q, want %q", got, want)
	}
	if tb.fatalMessage != "" {
		t.Fatalf("Until() recorded an unexpected failure: %q", tb.fatalMessage)
	}
}

func TestUntilFailsAtOnceWhenTheWorkEnds(t *testing.T) {
	t.Parallel()
	ready := make(chan string)
	ended := make(chan string, 1)
	ended <- "exit status 23"
	tb := &recordingTB{TB: t}
	started := time.Now()

	testwait.Until(tb, "Agent start", ready, ended)

	elapsed := time.Since(started)
	if elapsed >= time.Second {
		t.Fatalf("Until() took %s after work ended, want under one second", elapsed)
	}
	if !strings.Contains(tb.fatalMessage, "exit status 23") {
		t.Fatalf("Until() failure = %q, want ended work value", tb.fatalMessage)
	}
}

func TestUntilFailsAtTheDeadlineAndNotBefore(t *testing.T) {
	const wait = 300 * time.Millisecond
	started := time.Now()
	tb := &recordingTB{
		TB:          t,
		deadline:    started.Add(testwait.Margin + wait),
		hasDeadline: true,
	}
	stackStarted := make(chan struct{})
	release := make(chan struct{})
	go blockForStackDump(stackStarted, release)
	<-stackStarted
	t.Cleanup(func() { close(release) })
	var ready <-chan string
	var ended <-chan struct{}

	testwait.Until(tb, "deadline diagnostic", ready, ended)

	elapsed := time.Since(started)
	if elapsed < wait {
		t.Fatalf("Until() failed after %s, want no earlier than %s", elapsed, wait)
	}
	if !strings.Contains(tb.fatalMessage, "timed out waiting for deadline diagnostic") {
		t.Fatalf("Until() failure = %q, want named wait", tb.fatalMessage)
	}
	if !strings.Contains(tb.fatalMessage, "blockForStackDump") {
		t.Fatalf("Until() failure omitted another goroutine's stack:\n%s", tb.fatalMessage)
	}
}

func TestPollReturnsWhenTheConditionHolds(t *testing.T) {
	t.Parallel()
	checks := 0
	started := time.Now()
	var ended <-chan struct{}

	testwait.Poll(t, "third observation", ended, func() (bool, string) {
		checks++
		return checks == 3, fmt.Sprintf("check %d", checks)
	})

	elapsed := time.Since(started)
	if checks != 3 {
		t.Fatalf("Poll() evaluated the condition %d times, want 3", checks)
	}
	if elapsed >= time.Second {
		t.Fatalf("Poll() took %s after the condition became true, want under one second", elapsed)
	}
}

func TestPollFailsAtOnceWhenTheWorkEnds(t *testing.T) {
	t.Parallel()
	ended := make(chan string, 1)
	ended <- "exit status 29"
	tb := &recordingTB{TB: t}
	checks := 0
	started := time.Now()

	testwait.Poll(tb, "published event", ended, func() (bool, string) {
		checks++
		return false, fmt.Sprintf("check %d", checks)
	})

	elapsed := time.Since(started)
	if elapsed >= time.Second {
		t.Fatalf("Poll() took %s after work ended, want under one second", elapsed)
	}
	if checks != 2 {
		t.Fatalf("Poll() evaluated the condition %d times, want an initial and final check", checks)
	}
	if !strings.Contains(tb.fatalMessage, "exit status 29") {
		t.Fatalf("Poll() failure = %q, want ended work value", tb.fatalMessage)
	}
}

func TestPollReturnsWhenEndedWorkPublishedItsLastEffect(t *testing.T) {
	t.Parallel()
	ended := make(chan string, 1)
	ended <- "clean exit"
	tb := &recordingTB{TB: t}
	checks := 0
	started := time.Now()

	testwait.Poll(tb, "published event", ended, func() (bool, string) {
		checks++
		return checks == 2, fmt.Sprintf("check %d", checks)
	})

	elapsed := time.Since(started)
	if elapsed >= time.Second {
		t.Fatalf("Poll() took %s after the final effect, want under one second", elapsed)
	}
	if checks != 2 {
		t.Fatalf("Poll() evaluated the condition %d times, want an initial and final check", checks)
	}
	if tb.fatalMessage != "" {
		t.Fatalf("Poll() recorded an unexpected failure: %q", tb.fatalMessage)
	}
}

func TestPollFailsAtTheDeadlineWithTheLastObservation(t *testing.T) {
	const wait = 300 * time.Millisecond
	started := time.Now()
	tb := &recordingTB{
		TB:          t,
		deadline:    started.Add(testwait.Margin + wait),
		hasDeadline: true,
	}
	checks := 0
	lastObservation := ""
	var ended <-chan struct{}

	testwait.Poll(tb, "Run state", ended, func() (bool, string) {
		checks++
		lastObservation = fmt.Sprintf("state observation %d", checks)
		return false, lastObservation
	})

	elapsed := time.Since(started)
	if elapsed < wait {
		t.Fatalf("Poll() failed after %s, want no earlier than %s", elapsed, wait)
	}
	if checks < 2 {
		t.Fatalf("Poll() evaluated the condition %d times, want repeated observations", checks)
	}
	if !strings.Contains(tb.fatalMessage, "timed out waiting for Run state") {
		t.Fatalf("Poll() failure = %q, want named wait", tb.fatalMessage)
	}
	if !strings.Contains(tb.fatalMessage, lastObservation) {
		t.Fatalf("Poll() failure = %q, want last observation %q", tb.fatalMessage, lastObservation)
	}
}
