package cli

import (
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/spec"
)

// folderRemovalFinishCommand quotes an absolute folder for a POSIX shell.
func folderRemovalFinishCommand(folder string) string {
	quoted := "'" + strings.ReplaceAll(folder, "'", "'\\''") + "'"
	return "chmod -R u+w " + quoted + " && rm -rf " + quoted
}

func printFolderRemovalRecovery(command, record, folder string, removal *spec.FolderRemovalError, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: %s wrote %s but could not remove the Spec folder %s: %v\n", app.Name, command, record, folder, removal.Err)
	fmt.Fprintf(stderr, "%s: finish the removal with: %s\n", app.Name, folderRemovalFinishCommand(removal.Folder))
}
