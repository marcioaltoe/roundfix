//go:build !unix && !windows

package store

import "context"

func processAbsent(_ int) (bool, error) {
	return false, ErrOwnerProcessUnsupported
}

func signalOwnerProcess(_ int, _ bool) error {
	return ErrOwnerProcessUnsupported
}

func processStartIdentity(_ context.Context, _ int) (string, error) {
	return "", ErrOwnerProcessUnsupported
}
