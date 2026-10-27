package store

// AvailableBytes reports the bytes an unprivileged user can write on the
// filesystem holding path, which must exist.
func AvailableBytes(path string) (int64, error) {
	return availableTemporaryCapacity(path)
}
