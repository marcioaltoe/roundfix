//go:build !unix

package cli

// deliveryTextFileBusy is false where the kernel does not refuse to execute a
// file open for writing.
func deliveryTextFileBusy(error) bool {
	return false
}
