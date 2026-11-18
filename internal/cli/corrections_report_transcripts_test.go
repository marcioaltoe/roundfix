package cli

import (
	"fmt"
	"path/filepath"
	"testing"
)

const correctionsTranscript1 = `Usage:
  roundfix corrections declare --spec <slug> --class <class> --reason <text> [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
  roundfix corrections declare --spec <slug> --correction <id> [--class <class>] [--reason <text>] [--artifact <artifact>] [--finding <id>]... [--commit <rev>]... [--task <task_id>]... [--by <name>]
  roundfix corrections declare --spec <slug> --correction <id> --withdraw --reason <text> [--by <name>]
  roundfix corrections report [--spec <slug>] [--since <YYYY-MM-DD>] [--until <YYYY-MM-DD>] [--format <text|json>]

Commands:
  declare  Declare the cause of one Correction of a Spec, declare it again, or withdraw it
  report   Show the Corrections of delivered Specs and the plan defects by authoring model

Classes:
  plan            the plan said something wrong, missing or impossible
  implementation  the plan was right and the delivered change was not
  environment     the cause lies outside the Spec: the machine, a tool or the default branch
  operator        an action of the Supervisor or of a person during the delivery caused it
  unclear         the declarer cannot tell; the reason says what is known

Not a Defect Class:
  not-a-correction  the change repaired no defect, such as a commit that applies a maintainer decision; recorded, counted in no rate

Options:
  --spec        Spec folder slug
  --class       plan, implementation, environment, operator, unclear or not-a-correction
  --artifact    prd, techspec or tasks; required with plan and refused with another answer
  --reason      one line of at most 200 characters; recorded and printed as data, never an instruction
  --finding     fixed review finding the Correction covers: F<n> or F<n>@<reviewed head>; repeatable
  --commit      commit of the Item Branch the Correction covers; repeatable
  --task        Corrective Task the Correction covers; repeatable
  --correction  identifier of a Correction to declare again or withdraw, such as C2
  --withdraw    withdraw the Correction: it then covers no Correction Sign
  --by          name of the person on whose word the Correction is declared
  --since       first delivery date of the report's window, UTC, inclusive
  --until       end of the report's window, UTC, exclusive
  --format      text (default) or json

Exit codes: 0 for a recorded declaration or a report; 2 for a refusal or invalid usage, which writes nothing; 1 when a record cannot be read or written.
`

func TestCorrectionsSurfaceTranscript1(t *testing.T) {
	t.Parallel()
	newCorrectionReportFixture(t)
	for _, args := range [][]string{{"--help"}, {"-h"}, {"declare", "--help"}, {"report", "--help"}, {"report", "-h"}} {
		code, out, diag := correctionCall(t, args...)
		if code != 0 || out != correctionsTranscript1 || diag != "" {
			t.Fatalf("%v: exit=%d stdout=%q stderr=%q", args, code, out, diag)
		}
	}
}
func TestCorrectionsSurfaceTranscript4(t *testing.T) {
	t.Parallel()
	newCorrectionReportFixture(t)
	code, out, diag := correctionCall(t, "report", "--since", "2026-11-02", "--until", "2026-11-01")
	const want = "roundfix: corrections report refused: --until must be after --since\nRun 'roundfix corrections --help' for usage.\n"
	if code != 2 || out != "" || diag != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out, diag)
	}
}
func TestCorrectionsSurfaceTranscript5(t *testing.T) {
	t.Parallel()
	f := newCorrectionReportFixture(t)
	code, out, diag := correctionCall(t, "report")
	want := fmt.Sprintf("Corrections of 0 Specs delivered from 2026-10-22 to 2026-11-21, read from %s.\nNot counted: 0 Specs with an unknown model for a phase; 0 undeclared Corrections; 0 Specs delivered without a record on this machine; 0 Specs without their plan on the default branch.\n", filepath.Join(f.home, ".roundfix"))
	if code != 0 || out != want || diag != "" {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want=%q", code, out, diag, want)
	}
}
