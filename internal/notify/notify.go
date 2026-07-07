// Package notify fires one best-effort notification per terminal Run outcome.
package notify

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"roundfix/internal/config"
)

const commandTimeout = 30 * time.Second

// Outcome carries the Run context a notification names.
type Outcome struct {
	RunID  string
	State  string
	Kind   string
	Target string
}

// Notifier sends one notification; implementations are best-effort.
type Notifier interface {
	Notify(ctx context.Context, outcome Outcome) error
}

// New picks the notifier implementation from config.
func New(cfg config.Config) Notifier {
	return newWithDeps(cfg, dependencies{})
}

type dependencies struct {
	runner   commandRunner
	lookPath func(string) (string, error)
	native   func(dependencies) Notifier
	timeout  time.Duration
}

func (deps dependencies) withDefaults() dependencies {
	if deps.runner == nil {
		deps.runner = execRunner{}
	}
	if deps.lookPath == nil {
		deps.lookPath = exec.LookPath
	}
	if deps.native == nil {
		deps.native = platformNativeNotifier
	}
	if deps.timeout <= 0 {
		deps.timeout = commandTimeout
	}
	return deps
}

func newWithDeps(cfg config.Config, deps dependencies) Notifier {
	deps = deps.withDefaults()
	if !cfg.Notify.Enabled {
		return noopNotifier{}
	}
	if cfg.Notify.Command != "" {
		return &commandNotifier{
			command: cfg.Notify.Command,
			timeout: deps.timeout,
			runner:  deps.runner,
		}
	}
	return deps.native(deps)
}

type noopNotifier struct{}

func (noopNotifier) Notify(context.Context, Outcome) error {
	return nil
}

type commandNotifier struct {
	command string
	timeout time.Duration
	runner  commandRunner
}

func (notifier *commandNotifier) Notify(ctx context.Context, outcome Outcome) error {
	if ctx == nil {
		return errors.New("notify command: context is required")
	}
	if notifier == nil {
		return errors.New("notify command: notifier is required")
	}
	if notifier.command == "" {
		return errors.New("notify command: command is required")
	}
	runner := notifier.runner
	if runner == nil {
		runner = execRunner{}
	}
	timeout := notifier.timeout
	if timeout <= 0 {
		timeout = commandTimeout
	}

	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	name, args := shellCommand(notifier.command)
	output, err := runner.Run(runCtx, name, args, commandEnvironment(os.Environ(), outcome))
	if err == nil {
		return nil
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("notify command %q timed out after %s: %w", notifier.command, timeout, context.DeadlineExceeded)
	}
	if errors.Is(runCtx.Err(), context.Canceled) && ctx.Err() != nil {
		return fmt.Errorf("notify command %q canceled: %w", notifier.command, ctx.Err())
	}
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		detail = err.Error()
	}
	return fmt.Errorf("notify command %q failed: %s: %w", notifier.command, detail, err)
}

type commandRunner interface {
	Run(ctx context.Context, name string, args []string, env []string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, name string, args []string, env []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	if env != nil {
		cmd.Env = env
	}
	return cmd.CombinedOutput()
}

func commandEnvironment(base []string, outcome Outcome) []string {
	env := append([]string(nil), base...)
	env = append(env,
		"ROUNDFIX_RUN_ID="+outcome.RunID,
		"ROUNDFIX_OUTCOME="+outcome.State,
		"ROUNDFIX_KIND="+outcome.Kind,
		"ROUNDFIX_TARGET="+outcome.Target,
	)
	return env
}

type desktopNotifier struct {
	tool     string
	args     func(Outcome) []string
	lookPath func(string) (string, error)
	runner   commandRunner
}

func (notifier desktopNotifier) Notify(ctx context.Context, outcome Outcome) error {
	if ctx == nil {
		return errors.New("native notification: context is required")
	}
	if notifier.tool == "" || notifier.args == nil {
		return nil
	}
	lookPath := notifier.lookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	if _, err := lookPath(notifier.tool); err != nil {
		return nil
	}
	runner := notifier.runner
	if runner == nil {
		runner = execRunner{}
	}
	output, err := runner.Run(ctx, notifier.tool, notifier.args(outcome), nil)
	if err == nil {
		return nil
	}
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		detail = err.Error()
	}
	return fmt.Errorf("native notification %q failed: %s: %w", notifier.tool, detail, err)
}

func notificationBody(outcome Outcome) string {
	state := strings.TrimSpace(outcome.State)
	target := strings.TrimSpace(outcome.Target)
	if state == "" {
		state = "Run"
	}
	if target == "" {
		return state
	}
	return state + " - " + target
}

func appleScriptQuote(value string) string {
	var buffer bytes.Buffer
	buffer.WriteByte('"')
	for _, r := range value {
		switch r {
		case '\\', '"':
			buffer.WriteByte('\\')
			buffer.WriteRune(r)
		default:
			buffer.WriteRune(r)
		}
	}
	buffer.WriteByte('"')
	return buffer.String()
}
