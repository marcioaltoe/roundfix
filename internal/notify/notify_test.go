package notify

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"roundfix/internal/config"
)

func TestCommandNotifierEnvironmentAndFailedReceipt(t *testing.T) {
	runner := &fakeRunner{
		output: []byte("delivery failed\n"),
		err:    errors.New("exit status 12"),
	}
	notifier := commandNotifier{
		command: "send-notification",
		timeout: time.Second,
		runner:  runner,
	}

	before := time.Now()
	receipt, err := notifier.Notify(context.Background(), testOutcome())

	if err == nil {
		t.Fatal("expected command failure")
	}
	assertReceipt(t, receipt, RouteCommand, StatusFailed, before)
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
		"ROUNDFIX_RUN_ID":              "run-123",
		"ROUNDFIX_OUTCOME":             "Clean",
		"ROUNDFIX_KIND":                "implement",
		"ROUNDFIX_TARGET":              "spec:0019-run-outcome-notifications",
		"ROUNDFIX_REASON":              "verification failed",
		"ROUNDFIX_CONSOLE_LOG":         "/tmp/runs/run-123/console.log",
		"ROUNDFIX_ATTACH_COMMAND":      "roundfix attach run-123",
		"ROUNDFIX_REVIEW_ISSUES_KNOWN": "false",
		"ROUNDFIX_NEXT_ACTION":         "inspect the failed verification",
	}
	for key, value := range want {
		if env[key] != value {
			t.Fatalf("expected %s=%q, got %q", key, value, env[key])
		}
	}
}

func TestCommandNotifierSentReceipt(t *testing.T) {
	runner := &fakeRunner{}
	notifier := commandNotifier{
		command: "send-notification",
		timeout: time.Second,
		runner:  runner,
	}

	before := time.Now()
	receipt, err := notifier.Notify(context.Background(), testOutcome())

	if err != nil {
		t.Fatalf("send command notification: %v", err)
	}
	assertReceipt(t, receipt, RouteCommand, StatusSent, before)
}

func TestCommandNotifierTimesOutWithFailedReceipt(t *testing.T) {
	runner := &fakeRunner{waitForDone: true}
	notifier := commandNotifier{
		command: "slow-notification",
		timeout: 10 * time.Millisecond,
		runner:  runner,
	}

	before := time.Now()
	receipt, err := notifier.Notify(context.Background(), testOutcome())

	if err == nil {
		t.Fatal("expected timeout error")
	}
	assertReceipt(t, receipt, RouteCommand, StatusFailed, before)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", err)
	}
	if !strings.Contains(err.Error(), `notify command "slow-notification" timed out after 10ms`) {
		t.Fatalf("expected timeout error to name command and bound, got %q", err.Error())
	}
}

func TestNewSelectsNotifierAndDisabledReceipt(t *testing.T) {
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
				before := time.Now()
				receipt, err := got.Notify(context.Background(), testOutcome())
				if err != nil {
					t.Fatalf("disabled notification: %v", err)
				}
				assertReceipt(t, receipt, RouteDisabled, StatusSkipped, before)
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

func TestDesktopNotifierNativeReceipts(t *testing.T) {
	tests := []struct {
		name       string
		lookupErr  error
		runErr     error
		wantStatus Status
		wantCalls  int
	}{
		{
			name:       "unavailable is skipped",
			lookupErr:  errors.New("not found"),
			wantStatus: StatusSkipped,
		},
		{
			name:       "accepted is sent",
			wantStatus: StatusSent,
			wantCalls:  1,
		},
		{
			name:       "runner failure is failed",
			runErr:     errors.New("exit status 1"),
			wantStatus: StatusFailed,
			wantCalls:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeRunner{err: tt.runErr}
			notifier := desktopNotifier{
				tool: "native-notifier",
				args: func(Outcome) []string {
					return []string{"Roundfix", "Clean - spec:0019-run-outcome-notifications"}
				},
				lookPath: func(command string) (string, error) {
					if command != "native-notifier" {
						t.Fatalf("expected lookup for native-notifier, got %q", command)
					}
					return command, tt.lookupErr
				},
				runner: runner,
			}

			before := time.Now()
			receipt, err := notifier.Notify(context.Background(), testOutcome())

			if tt.wantStatus == StatusFailed {
				if err == nil {
					t.Fatal("expected native notification failure")
				}
			} else if err != nil {
				t.Fatalf("native notification: %v", err)
			}
			assertReceipt(t, receipt, RouteNative, tt.wantStatus, before)
			if len(runner.calls) != tt.wantCalls {
				t.Fatalf("native runner calls = %d, want %d", len(runner.calls), tt.wantCalls)
			}
		})
	}
}

func TestNativeNotificationTextIsBoundedAndKeepsNextAction(t *testing.T) {
	outcome := testOutcome()
	outcome.State = "Failed"
	outcome.Target = "spec:" + strings.Repeat("large-target-", 40)
	outcome.NextAction = "inspect the Console Log"

	body := notificationBody(outcome)

	if got := utf8.RuneCountInString(body); got > nativeTextLimit {
		t.Fatalf("native notification body length = %d, want <= %d: %q", got, nativeTextLimit, body)
	}
	if !strings.Contains(body, "Next: "+outcome.NextAction) {
		t.Fatalf("native notification body omitted next action: %q", body)
	}
}

func TestCommandEnvironmentPreservesExistingVariablesAndAddsTerminalContext(t *testing.T) {
	known := true
	outcome := testOutcome()
	outcome.ReviewIssuesKnown = &known
	base := []string{"EXISTING=unchanged"}

	got := commandEnvironment(base, outcome)

	want := []string{
		"EXISTING=unchanged",
		"ROUNDFIX_RUN_ID=run-123",
		"ROUNDFIX_OUTCOME=Clean",
		"ROUNDFIX_KIND=implement",
		"ROUNDFIX_TARGET=spec:0019-run-outcome-notifications",
		"ROUNDFIX_REASON=verification failed",
		"ROUNDFIX_CONSOLE_LOG=/tmp/runs/run-123/console.log",
		"ROUNDFIX_ATTACH_COMMAND=roundfix attach run-123",
		"ROUNDFIX_REVIEW_ISSUES_KNOWN=true",
		"ROUNDFIX_NEXT_ACTION=inspect the failed verification",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("command environment mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
	if base[0] != "EXISTING=unchanged" {
		t.Fatalf("base environment mutated: %#v", base)
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

	receipt, err := notifier.Notify(context.Background(), testOutcome())
	if err != nil {
		t.Fatalf("expected missing native tool to no-op, got %v", err)
	}
	if receipt.Status != StatusSkipped {
		t.Fatalf("missing native tool receipt = %#v, want skipped", receipt)
	}
	if len(runner.calls) != 0 {
		t.Fatalf("expected missing tool to skip command run, got %d calls", len(runner.calls))
	}
}

func testOutcome() Outcome {
	known := false
	return Outcome{
		RunID:             "run-123",
		State:             "Clean",
		Kind:              "implement",
		Target:            "spec:0019-run-outcome-notifications",
		Reason:            "verification failed",
		ConsoleLog:        "/tmp/runs/run-123/console.log",
		AttachCommand:     "roundfix attach run-123",
		ReviewIssuesKnown: &known,
		NextAction:        "inspect the failed verification",
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

func (sentinelNotifier) Notify(context.Context, Outcome) (NotificationReceipt, error) {
	return NotificationReceipt{Route: RouteNative, Status: StatusSent, CompletedAt: time.Now().UTC()}, nil
}

func assertReceipt(t *testing.T, receipt NotificationReceipt, route Route, status Status, notBefore time.Time) {
	t.Helper()
	if receipt.Route != route || receipt.Status != status {
		t.Fatalf("receipt = %#v, want route %q status %q", receipt, route, status)
	}
	if receipt.CompletedAt.IsZero() || receipt.CompletedAt.Before(notBefore) || receipt.CompletedAt.After(time.Now().Add(time.Second)) {
		t.Fatalf("receipt completion time = %s, want current non-zero UTC time", receipt.CompletedAt)
	}
	if receipt.CompletedAt.Location() != time.UTC {
		t.Fatalf("receipt completion location = %s, want UTC", receipt.CompletedAt.Location())
	}
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
