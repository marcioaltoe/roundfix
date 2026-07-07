//go:build linux

package notify

func platformNativeNotifier(deps dependencies) Notifier {
	return desktopNotifier{
		tool: "notify-send",
		args: func(outcome Outcome) []string {
			return []string{"Roundfix", notificationBody(outcome)}
		},
		lookPath: deps.lookPath,
		runner:   deps.runner,
	}
}
