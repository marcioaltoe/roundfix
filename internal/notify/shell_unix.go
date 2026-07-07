//go:build !windows

package notify

func shellCommand(command string) (string, []string) {
	return "sh", []string{"-c", command}
}
