package cli

import (
	"fmt"
	"io"

	"roundfix/internal/app"
)

const usage = `Roundfix

Usage:
  roundfix --help
  roundfix --version
  roundfix fetch
  roundfix resolve
  roundfix watch

Commands:
  fetch      Download review issues for an open pull request (not implemented)
  resolve    Resolve downloaded unresolved review issues (not implemented)
  watch      Fetch and resolve in a watched loop (not implemented)

Options:
  -h, --help      Show help
  --version       Show version
`

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "--version", "version":
		fmt.Fprintf(stdout, "%s %s\n", app.Name, app.Version)
		return 0
	case "fetch", "resolve", "watch":
		fmt.Fprintf(stderr, "%s: command %q is not implemented yet\n", app.Name, args[0])
		return 2
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s --help' for usage.\n", app.Name)
		return 2
	}
}
