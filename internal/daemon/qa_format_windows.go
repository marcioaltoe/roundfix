package daemon

import "os/exec"

// CommandContext kills the formatter process on cancellation on Windows.
func configureQAFormatProcess(_ *exec.Cmd) {}
