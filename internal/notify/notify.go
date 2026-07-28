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

const (
	commandTimeout  = 30 * time.Second
	nativeTextLimit = 256
)

type Route string
type Status string

const (
	RouteCommand  Route = "command"
	RouteNative   Route = "native"
	RouteDisabled Route = "disabled"

	StatusSent    Status = "sent"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
)

// Outcome carries the Run context a notification names.
type Outcome struct {
	RunID             string
	State             string
	Kind              string
	Target            string
	Reason            string
	ConsoleLog        string
	AttachCommand     string
	ReviewIssuesKnown *bool
	NextAction        string
}

// NotificationReceipt records the completed delivery attempt independently
// from the Run outcome.
type NotificationReceipt struct {
	Route       Route
	Status      Status
	CompletedAt time.Time
}

// Notifier sends one notification; implementations are best-effort.
type Notifier interface {
	Notify(ctx context.Context, outcome Outcome) (NotificationReceipt, error)
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
		return noopNotifier{route: RouteDisabled}
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

type noopNotifier struct {
	route Route
}

func (notifier noopNotifier) Notify(context.Context, Outcome) (NotificationReceipt, error) {
	route := notifier.route
	if route == "" {
		route = RouteDisabled
	}
	return completedReceipt(route, StatusSkipped), nil
}

type commandNotifier struct {
	command string
	timeout time.Duration
	runner  commandRunner
}

func (notifier *commandNotifier) Notify(ctx context.Context, outcome Outcome) (NotificationReceipt, error) {
	if ctx == nil {
		return completedReceipt(RouteCommand, StatusFailed), errors.New("notify command: context is required")
	}
	if notifier == nil {
		return completedReceipt(RouteCommand, StatusFailed), errors.New("notify command: notifier is required")
	}
	if notifier.command == "" {
		return completedReceipt(RouteCommand, StatusFailed), errors.New("notify command: command is required")
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
		return completedReceipt(RouteCommand, StatusSent), nil
	}
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return completedReceipt(RouteCommand, StatusFailed), fmt.Errorf("notify command %q timed out after %s: %w", notifier.command, timeout, context.DeadlineExceeded)
	}
	if errors.Is(runCtx.Err(), context.Canceled) && ctx.Err() != nil {
		return completedReceipt(RouteCommand, StatusFailed), fmt.Errorf("notify command %q canceled: %w", notifier.command, ctx.Err())
	}
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		detail = err.Error()
	}
	return completedReceipt(RouteCommand, StatusFailed), fmt.Errorf("notify command %q failed: %s: %w", notifier.command, detail, err)
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
		"ROUNDFIX_REASON="+outcome.Reason,
		"ROUNDFIX_CONSOLE_LOG="+outcome.ConsoleLog,
		"ROUNDFIX_ATTACH_COMMAND="+outcome.AttachCommand,
		"ROUNDFIX_REVIEW_ISSUES_KNOWN="+reviewIssuesKnownEnvironmentValue(outcome.ReviewIssuesKnown),
		"ROUNDFIX_NEXT_ACTION="+outcome.NextAction,
	)
	return env
}

func reviewIssuesKnownEnvironmentValue(known *bool) string {
	if known == nil {
		return ""
	}
	if *known {
		return "true"
	}
	return "false"
}

type desktopNotifier struct {
	tool     string
	args     func(Outcome) []string
	lookPath func(string) (string, error)
	runner   commandRunner
}

func (notifier desktopNotifier) Notify(ctx context.Context, outcome Outcome) (NotificationReceipt, error) {
	if ctx == nil {
		return completedReceipt(RouteNative, StatusFailed), errors.New("native notification: context is required")
	}
	if notifier.tool == "" || notifier.args == nil {
		return completedReceipt(RouteNative, StatusSkipped), nil
	}
	lookPath := notifier.lookPath
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	if _, err := lookPath(notifier.tool); err != nil {
		return completedReceipt(RouteNative, StatusSkipped), nil
	}
	runner := notifier.runner
	if runner == nil {
		runner = execRunner{}
	}
	output, err := runner.Run(ctx, notifier.tool, notifier.args(outcome), nil)
	if err == nil {
		return completedReceipt(RouteNative, StatusSent), nil
	}
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		detail = err.Error()
	}
	return completedReceipt(RouteNative, StatusFailed), fmt.Errorf("native notification %q failed: %s: %w", notifier.tool, detail, err)
}

func notificationBody(outcome Outcome) string {
	state := strings.TrimSpace(outcome.State)
	target := strings.TrimSpace(outcome.Target)
	if state == "" {
		state = "Run"
	}
	body := state
	if target != "" {
		body += " - " + target
	}
	nextAction := strings.Join(strings.Fields(outcome.NextAction), " ")
	if nextAction == "" || strings.EqualFold(state, "Clean") {
		return boundNativeText(body, nativeTextLimit)
	}
	action := "Next: " + nextAction
	separator := ". "
	availableBody := nativeTextLimit - len([]rune(separator)) - len([]rune(action))
	if availableBody < len([]rune(state)) {
		return boundNativeText(state+separator+action, nativeTextLimit)
	}
	return boundNativeText(body, availableBody) + separator + action
}

func boundNativeText(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit <= 0 {
		return ""
	}
	if limit == 1 {
		return "…"
	}
	return string(runes[:limit-1]) + "…"
}

func completedReceipt(route Route, status Status) NotificationReceipt {
	return NotificationReceipt{
		Route:       route,
		Status:      status,
		CompletedAt: time.Now().UTC(),
	}
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
