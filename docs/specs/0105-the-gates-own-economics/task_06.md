---
status: pending
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running commands.

## Work

- A `qa` Task's Verification is supplied by Roundfix, rendered into the Task
  file, and an authored one is refused by name
- The derived command cannot accept a verdict outside the domain, exercised
  against the measured case that did
- A finding blocks the rows it names; unnamed rows in the same matrix are
  measured; a finding naming every row still blocks every row
- Withholding is unchanged: a blocking machine fact before a matrix exists still
  withholds the Agent Session
- The two measured citation forms parse, a bare number outside an obligations
  line is still not a citation, and an unrecognised citation names the form that
  is
- The authoring skill carries the characterization obligation and its ordering
- The gate applies the equivalent-evidence path to the Pull Request row and
  records the evidence; a row with no evidence stays blocked
- The skill changes stay inside the recorded authorization's bounded files,
  checked from Git evidence, with the authorization commit preceding the skill
  commit and the generated copies matching `make skills-sync`
- The glossary check: whether this Spec introduced, changed, or retired a term
  the domain context should carry

## Outside evidence

One acceptance row rests on evidence this Spec did not author. Of 201 failed
Tasks measured across five repositories, 123 were the QA gate returning a
verdict rather than code breaking, and one Spec paid six of its eight gate
executions for the Pull Request row alone — repositories this Spec did not
build, measured before this Spec existed. The row records that this measurement
is what establishes the requirement, rather than a rehearsal of the Spec's own
premise.

Read against the counter-number in the same evidence set: eleven non-passing
verdicts in one session all failed on contract rather than business logic, and
the same gate found four real defects no suite would catch. A change that made
the gate cheaper by making it weaker would satisfy the first number and betray
the second.

## References

- All user stories and core features

## Verification
- `newest="$(ls -1 docs/specs/0105-the-gates-own-economics/qa/qa-report-*.md 2>/dev/null | sort | tail -1)"; test -n "$newest" || exit 1; grep -q "^verdict: pass" "$newest" || exit 1; roundfix spec check 0105-the-gates-own-economics --strict && go test -count=1 ./internal/spec ./internal/speccheck ./internal/daemon`
