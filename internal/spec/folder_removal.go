package spec

import "fmt"

// FolderRemovalError reports a failed final removal after every archive or
// cancellation output was written. The caller must preserve those writes.
type FolderRemovalError struct {
	Folder string // Absolute Spec folder.
	Err    error
}

func (e *FolderRemovalError) Error() string {
	return fmt.Sprintf("remove Spec folder %q: %v", e.Folder, e.Err)
}

func (e *FolderRemovalError) Unwrap() error { return e.Err }
