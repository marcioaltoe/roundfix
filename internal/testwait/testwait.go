// Package testwait provides test-only waits bounded by the enclosing test's
// deadline. A failing wait calls t.Fatalf, so callers must invoke these helpers
// on the test goroutine. The helpers cancel nothing; each test remains
// responsible for cancelling watched work through its own context and cleanups.
package testwait

import (
	"runtime"
	"testing"
	"time"
)

const (
	// Margin leaves enough of the test deadline for a wait failure to report.
	Margin = time.Second

	fallback     = 10 * time.Minute
	pollInterval = 10 * time.Millisecond
)

type deadliner interface {
	Deadline() (time.Time, bool)
}

// Bound returns the remaining test deadline minus Margin. Tests without a
// deadline use go test's default ten-minute timeout.
func Bound(t testing.TB) time.Duration {
	t.Helper()
	deadlineSource, ok := t.(deadliner)
	if !ok {
		return fallback
	}
	deadline, ok := deadlineSource.Deadline()
	if !ok {
		return fallback
	}
	return time.Until(deadline) - Margin
}

// Until returns the first value from ready. If ready and ended are both ready,
// the ready value wins. If ended delivers first, Until fails immediately with
// that value. A nil ended channel disables work-end observation. If neither
// channel delivers before Bound, Until fails with every goroutine's stack.
func Until[T, R any](t testing.TB, what string, ready <-chan T, ended <-chan R) T {
	t.Helper()
	timer := time.NewTimer(Bound(t))
	defer timer.Stop()

	select {
	case value := <-ready:
		return value
	default:
	}

	select {
	case value := <-ready:
		return value
	case result := <-ended:
		select {
		case value := <-ready:
			return value
		default:
		}
		t.Fatalf("stopped waiting for %s: watched work ended: %v", what, result)
	case <-timer.C:
		select {
		case value := <-ready:
			return value
		default:
		}
		failTimeout(t, what, nil)
	}

	var zero T
	return zero
}

// Poll evaluates condition until it reports ready. If ended delivers first,
// Poll evaluates condition once more before failing with the ended value. A
// nil ended channel disables work-end observation. If the bound expires, Poll
// fails with the condition's last observation.
func Poll[R any](t testing.TB, what string, ended <-chan R, condition func() (bool, string)) {
	t.Helper()
	ready, observation := condition()
	if ready {
		return
	}

	timer := time.NewTimer(Bound(t))
	defer timer.Stop()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case result := <-ended:
			ready, observation = condition()
			if ready {
				return
			}
			t.Fatalf("stopped waiting for %s: watched work ended: %v; last observation: %s", what, result, observation)
			return
		case <-ticker.C:
			ready, observation = condition()
			if ready {
				return
			}
		case <-timer.C:
			ready, observation = condition()
			if ready {
				return
			}
			failTimeout(t, what, &observation)
			return
		}
	}
}

func failTimeout(t testing.TB, what string, observation *string) {
	t.Helper()
	stack := make([]byte, 64<<10)
	for {
		n := runtime.Stack(stack, true)
		if n < len(stack) {
			if observation == nil {
				t.Fatalf("timed out waiting for %s\n%s", what, stack[:n])
				return
			}
			t.Fatalf("timed out waiting for %s; last observation: %s\n%s", what, *observation, stack[:n])
			return
		}
		stack = make([]byte, len(stack)*2)
	}
}
