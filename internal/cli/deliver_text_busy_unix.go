//go:build unix

package cli

import (
	"errors"
	"syscall"
)

// deliveryTextFileBusy reports a start refused because the executable was still
// open for writing in some process.
func deliveryTextFileBusy(err error) bool {
	return errors.Is(err, syscall.ETXTBSY)
}
