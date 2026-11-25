// Boundary: CLI dispatch and complete streams; no external commands.
package cli

import (
	"bytes"
	"testing"
)

const retiredLaneDeliverHelp = `Usage:
  roundfix deliver plan [--json] [<slug>...]
  roundfix deliver start <slug>...
  roundfix deliver status
  roundfix deliver resume
  roundfix deliver retry <slug>
  roundfix deliver stop

Records and advances an ordered queue of Specs. The detached owner survives the
calling terminal. A blocked item is parked and the owner continues with the
next item.

Commands:
  plan    Report which Specs are approved to run
  start   Validate and record a new queue, then start its detached owner
  status  Print every queued Spec's stage, blocker and item worktree
  resume  Start a detached owner for the persisted queue
  retry   Return one parked Spec to the queue and reach its owner
  stop    Prove and terminate the persisted queue owner

Flags:
  --max-duration <duration>  Set a positive queue duration
  --max-retries <n>          Set the per-item retry limit to at least 1
  --max-tokens <n>           Set the queue token ceiling to at least 1
`

func TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name            string
		args            []string
		out, diagnostic string
		code            int
	}{
		{"no_arguments", []string{"deliver", "lane"}, "", "roundfix: unknown deliver command \"lane\"\nRun 'roundfix deliver --help' for usage.\n", 2},
		{"with_arguments", []string{"deliver", "lane", "--all", "HEAD"}, "", "roundfix: unknown deliver command \"lane\"\nRun 'roundfix deliver --help' for usage.\n", 2},
		{"help_flag", []string{"deliver", "lane", "--help"}, retiredLaneDeliverHelp, "", 0},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			var out, diagnostic bytes.Buffer
			code := runCLI(t, row.args, &out, &diagnostic)
			if code != row.code || out.String() != row.out || diagnostic.String() != row.diagnostic {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
			}
		})
	}
}
func TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane(t *testing.T) {
	t.Parallel()
	var out, diagnostic bytes.Buffer
	code := runCLI(t, []string{"deliver", "--help"}, &out, &diagnostic)
	if code != 0 || out.String() != retiredLaneDeliverHelp || diagnostic.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}
