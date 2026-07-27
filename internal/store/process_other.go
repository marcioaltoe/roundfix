//go:build !unix && !windows

package store

func processAbsent(_ int) (bool, error) {
	return false, ErrOwnerProcessUnsupported
}

func signalOwnerProcess(_ int, _ bool) error {
	return ErrOwnerProcessUnsupported
}
