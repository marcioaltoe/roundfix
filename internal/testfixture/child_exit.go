package testfixture

import "roundfix/internal/childexit"

// ChildExitCode returns a child's exit code and, when the child did not exit
// on its own, a note for the helper's caller: a Signal Exit or a child that
// did not run. The note is empty for an exited child.
func ChildExitCode(err error) (code int, note string) {
	exit, ok := childexit.From(err)
	if !ok {
		return -1, "child did not run: " + err.Error()
	}
	if exit.Signaled() {
		return -1, "child " + exit.String()
	}
	return exit.Code, ""
}
