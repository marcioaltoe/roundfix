package cli

import "io"

// terminalOutput is the writer a terminal program draws on: the command's
// environment wrapper removed, so Bubble Tea receives the terminal file it
// can size. Color and TUI decisions keep reading the wrapper.
func terminalOutput(writer io.Writer) io.Writer {
	_, output := environmentForWriter(writer)
	return output
}
