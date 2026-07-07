//go:build darwin

package notify

func platformNativeNotifier(deps dependencies) Notifier {
	return desktopNotifier{
		tool: "osascript",
		args: func(outcome Outcome) []string {
			script := "display notification " + appleScriptQuote(notificationBody(outcome)) + " with title " + appleScriptQuote("Roundfix")
			return []string{"-e", script}
		},
		lookPath: deps.lookPath,
		runner:   deps.runner,
	}
}
