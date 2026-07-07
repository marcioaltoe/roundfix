package notify

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"roundfix/internal/config"
)

func TestCommandNotifierPassesEnvironmentAndReportsFailure(t *testing.T) {
	runner := &fakeRunner{
		output: []byte("delivery failed\n"),
		err:    errors.New("exit status 12"),
	}
	notifier := commandNotifier{
		command: "send-notification",
		timeout: time.Second,
		runner:  runner,
	}

	err := notifier.Notify(context.Background(), testOutcome())

	if err == nil {
		t.Fatal("expected command failure")
	}
	if !strings.Contains(err.Error(), `notify command "send-notification" failed`) {
		t.Fatalf("expected error to name command, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "delivery failed") {
		t.Fatalf("expected error to include command output, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "exit status 12") {
		t.Fatalf("expected error to include command error, got %q", err.Error())
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected one command call, got %d", len(runner.calls))
	}
	env := envMap(runner.calls[0].env)
	want := map[string]string{
		"ROUNDFIX_RUN_ID":  "run-123",
		"ROUNDFIX_OUTCOME": "Clean",
		"ROUNDFIX_KIND":    "implement",
		"ROUNDFIX_TARGET":  "spec:0019-run-outcome-notifications",
	}
	for key, value := range want {
		if env[key] != value {
			t.Fatalf("expected %s=%q, got %q", key, value, env[key])
		}
	}
}

func TestCommandNotifierTimesOut(t *testing.T) {
	runner := &fakeRunner{waitForDone: true}
	notifier := commandNotifier{
		command: "slow-notification",
		timeout: 10 * time.Millisecond,
		runner:  runner,
	}

	err := notifier.Notify(context.Background(), testOutcome())

	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if !strings.Contains(err.Error(), `notify command "slow-notification" timed out after 10ms`) {
		t.Fatalf("expected timeout error to name command and bound, got %q", err.Error())
	}
}

func TestNewSelectsNotifier(t *testing.T) {
	native := sentinelNotifier{}
	tests := []struct {
		name string
		cfg  config.Config
		want string
	}{
		{
			name: "disabled returns no-op",
			cfg: config.Config{
				Notify: config.Notify{
					Enabled: false,
					Command: "send-notification",
				},
			},
			want: "noop",
		},
		{
			name: "command returns command notifier",
			cfg: config.Config{
				Notify: config.Notify{
					Enabled: true,
					Command: "send-notification",
				},
			},
			want: "command",
		},
		{
			name: "enabled without command returns native",
			cfg: config.Config{
				Notify: config.Notify{
					Enabled: true,
				},
			},
			want: "native",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWithDeps(tt.cfg, dependencies{
				native: func(dependencies) Notifier {
					return native
				},
			})

			switch tt.want {
			case "noop":
				if _, ok := got.(noopNotifier); !ok {
					t.Fatalf("expected noopNotifier, got %T", got)
				}
			case "command":
				if _, ok := got.(*commandNotifier); !ok {
					t.Fatalf("expected *commandNotifier, got %T", got)
				}
			case "native":
				if got != native {
					t.Fatalf("expected native notifier, got %T", got)
				}
			default:
				t.Fatalf("unsupported expectation %q", tt.want)
			}
		})
	}
}

func TestDesktopNotifierMissingToolIsNoop(t *testing.T) {
	runner := &fakeRunner{}
	notifier := desktopNotifier{
		tool: "missing-notifier",
		args: func(Outcome) []string {
			return []string{"Roundfix", "Clean - spec:0019-run-outcome-notifications"}
		},
		lookPath: func(command string) (string, error) {
			if command != "missing-notifier" {
				t.Fatalf("expected lookup for missing-notifier, got %q", command)
			}
			return "", errors.New("not found")
		},
		runner: runner,
	}

	if err := notifier.Notify(context.Background(), testOutcome()); err != nil {
		t.Fatalf("expected missing native tool to no-op, got %v", err)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected missing tool to skip command run, got %d calls", len(runner.calls))
	}
}

func testOutcome() Outcome {
	return Outcome{
		RunID:  "run-123",
		State:  "Clean",
		Kind:   "implement",
		Target: "spec:0019-run-outcome-notifications",
	}
}

type fakeRunner struct {
	calls       []runCall
	output      []byte
	err         error
	waitForDone bool
}

func (runner *fakeRunner) Run(ctx context.Context, name string, args []string, env []string) ([]byte, error) {
	runner.calls = append(runner.calls, runCall{
		name: name,
		args: append([]string(nil), args...),
		env:  append([]string(nil), env...),
	})
	if runner.waitForDone {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return runner.output, runner.err
}

type runCall struct {
	name string
	args []string
	env  []string
}

type sentinelNotifier struct{}

func (sentinelNotifier) Notify(context.Context, Outcome) error {
	return nil
}

func envMap(env []string) map[string]string {
	values := map[string]string{}
	for _, entry := range env {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	return values
}
