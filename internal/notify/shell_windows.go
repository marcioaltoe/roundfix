//go:build windows

package notify

func shellCommand(command string) (string, []string) {
	return "cmd", []string{"/C", command}
}
