package agent

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

const CursorLoginRequired = "cursor_login_required"

// CursorLoginRequiredError refuses a Cursor selection without retaining account data.
type CursorLoginRequiredError struct{ Command string }

func (CursorLoginRequiredError) Error() string          { return "cursor-agent is not logged in" }
func (CursorLoginRequiredError) Classification() string { return CursorLoginRequired }
func (CursorLoginRequiredError) NextAction() string {
	return "run cursor-agent login in a terminal yourself"
}

// checkCursorLogin inspects status only in memory, under the caller's deadline
// and a ten-second ceiling. Command failures never expose their output.
func checkCursorLogin(ctx context.Context, executable string, environment []string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, executable, "status")
	command.Env = environment
	output, err := command.CombinedOutput()
	if err != nil || bytes.Contains(output, []byte("Not logged in")) {
		return CursorLoginRequiredError{Command: executable}
	}
	return nil
}
