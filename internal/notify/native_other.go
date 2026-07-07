//go:build !darwin && !linux

package notify

func platformNativeNotifier(deps dependencies) Notifier {
	return noopNotifier{}
}
