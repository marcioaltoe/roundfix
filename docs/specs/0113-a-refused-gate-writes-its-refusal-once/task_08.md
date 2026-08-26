---
status: completed
type: backend
---

# Task: Give the refusal report a writer in production

QA finding F-1 (`Blocks-Completion`): `spec.WritePreconditionRefusalReport` and
`speccheck.PreconditionRefusal` are complete and correct, and nothing calls
either one. The only non-test reference to the writer in the tree is its own
definition. The Spec's stated goal is not delivered on any path a Supervisor
reaches.

## Work

- Wire the deriver into the Daemon's mechanical request so `Precondition`
  carries the refusing check and its reason instead of its zero value, and the
  report the Daemon writes emits `rows_blocked_precondition`,
  `precondition_check`, and `precondition_reason` when a precondition refused.
  The three pre-existing counts keep their current meaning.
- Give the refusing gate its contract: `.agents/skills/qa-gate/SKILL.md` states
  what a gate writes when it stops at a precondition — one terminal row, its
  frontmatter, and the check and reason that caused the refusal — so a gate that
  refuses leaves an artifact the mechanical stage accepts. This path is the one
  the PRD describes, and it is the reason the Tooling authority row was amended
  on 2026-08-26.
- Regenerate the skill mirror with the declared `make skills-sync`; never edit
  the mirror by hand.
- A gate that did not refuse at a precondition writes exactly what it writes
  today. The refusal path is additive.

## References

- User Story 1: Gate writes valid report
- User Story 2: Refusal recorded and auditable
- Core Feature 1: Terminal Row Writing
- Core Feature 3: Precondition recorded in the report

## Verification
- `grep -rq "WritePreconditionRefusalReport" internal/daemon/ && grep -q "rows_blocked_precondition" internal/daemon/task_engine.go && grep -qi "precondition" .agents/skills/qa-gate/SKILL.md && diff -q .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && go test -count=1 ./internal/daemon ./internal/spec 2>&1 | grep -q "^ok"`

## Result

A gate that refuses at a precondition writes a report its own contract accepts,
on the path a Supervisor actually reaches. F-1 is closed at both ends: the
deriver now has a production caller and the writer now has a production caller.

### What changed

**The precondition reaches the stage.** `qaGatePrecondition`
(`internal/daemon/task_engine.go:2095`) runs the gate's own precondition — the
Spec's `spec check --strict`, which is what `speccheck.GatePreconditionCheck`
names — and `qaMechanicalRequest` puts the classified result on
`MechanicalRequest.Precondition` instead of leaving it zero. The three calls are
the command spelled out: `speccheck.Check`, then `speccheck.PromoteGaps`, then
`speccheck.GatePrecondition`. `RunMechanicalStage` already derived
`PreconditionRefusal` from that field; it was reading a zero value on every run
before this Task.

The consequence a Supervisor sees is the one ADR 0117 asks for: a Spec that
still contradicts itself stops the gate before an Agent turn is spent on it,
because the mechanical stage blocks and `runQAGate` withholds the Agent.

**The refusal is what the report says.** `mechanicalQAReportContent`
(`internal/daemon/task_engine.go:2244`) branches on `PreconditionRefused` and
hands the whole report to `spec.WritePreconditionRefusalReport`: verdict `fail`,
`rows_blocked_precondition: 1` beside the three typed counts at zero,
`precondition_check` and `precondition_reason` naming what refused and why, and
one terminal row `| 0 | blocked | precondition |` under `## Results`.

The mechanical sections stay below that report under their own headings. Since
task_07 bound row collection to the table under `## Results`, a table under
`## Mechanical rows` yields no rows, so the refusal costs the artifact none of
the observations the stage already made while still presenting exactly one
terminal row as the matrix. The Daemon's own `| QA-PRECONDITION | fail | … |`
row is gone from the refusal path only; every other gate reaches the unchanged
assembly below the branch and writes byte-identical output.

**One rule, one copy.** `speccheck.PromoteGaps`
(`internal/speccheck/report.go:70`) is the CLI's former `promoteSpecCheckGaps`,
exported and called from both ends. The Daemon's precondition has to be the same
verdict `roundfix spec check <slug> --strict` gives, because that command is what
the refusal names; a second copy of the promotion rule is exactly how the two
would come to disagree.

**The refusing gate has a contract.** `.agents/skills/qa-gate/SKILL.md` §1 now
states what a gate writes when the strict check fails — the full frontmatter,
the single row, and that the reason names every refusing `SC-*` code with its
sentence — and why an empty Results table is not an option. §6 states the closing
rule: a report carrying the refusal row carries `rows_blocked_precondition` set
to the exact number of `precondition`-provenance rows plus the two metadata keys,
and a gate that reached its matrix writes none of the three. That matches
`detectPreconditionCount`, which requires the fourth count only of a report that
records a refusal. The mirror was regenerated with `make skills-sync`; it was not
edited by hand.

### Rendered artifact

`mechanicalQAReportContent` for a refusal on `SC-CONSTRAINT-MISSING`, with one
mechanical finding and one skip still on the result:

```markdown
---
verdict: fail
rows_blocked_precondition: 1
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_blocked_declared: 0
precondition_check: "spec check --strict"
precondition_reason: "SC-CONSTRAINT-MISSING: _prd.md declares no Identifier strategy row"
---

# QA Report

## Results

| # | Status | Provenance |
| - | --- | --- |
| 0 | blocked | precondition |

## Precondition refusal

- check: spec check --strict
- reason: SC-CONSTRAINT-MISSING: _prd.md declares no Identifier strategy row

The gate stopped at this check before it built its QA matrix, so no requirement was measured and the row above records the refusal itself.

## Performed repairs

None.

## Assigned repair failures

None.

## Mechanical findings

### SC-CONSTRAINT-MISSING

- location: `docs/specs/0113/_prd.md:1`
- detail: _prd.md declares no Identifier strategy row
- fix: Declare the row.
- blocked row: `R01`

## Mechanical rows

| # | Status | Provenance |
| - | --- | --- |
| R01 | blocked (finding: SC-CONSTRAINT-MISSING — waits on no Identifier strategy row) | mechanical finding |

## Mechanical skips

| Detector | Missing artifact |
| --- | --- |
| consequent-fix commit order | declarations |
```

### Evidence per acceptance criterion

Red first: with the two new tests in place and `internal/daemon/task_engine.go`
stashed, `go test -count=1 -run
'TestQAMechanicalRequestCarriesTheGatePrecondition|TestWriteMechanicalQAReportWritesThePreconditionRefusal'
./internal/daemon` failed on four subtests — `Precondition =
speccheck.GatePreconditionResult{…, Blocking:false}, want the contradicted Spec
to block the gate`; the frontmatter and Results assertions; and `QAReport = {…
RowsBlockedPrecondition:0 Precondition:{CheckName: Reason:}}, want verdict "fail"
with one precondition-blocked row`. The two subtests that passed unstashed are
the regression guards.

Focused check after the last edit — the same `-run` selection with `-v` → `ok
roundfix/internal/daemon 0.788s`, both tests and all seven subtests passing.

| Work item | Evidence |
| --- | --- |
| `Precondition` carries the refusing check and reason | `TestQAMechanicalRequestCarriesTheGatePrecondition` builds a real fixture Spec, strips its Project Constraints, and requires `qaMechanicalRequest` to return a blocking precondition whose `speccheck.PreconditionRefusal` names `GatePreconditionCheck` and the refusing `SC-CONSTRAINT-MISSING` code — the reason is asserted against the code constant, not a copied sentence. Its sibling subtest requires the untouched fixture to refuse nothing, so the wiring cannot pass by refusing everything |
| The report emits `rows_blocked_precondition`, `precondition_check`, `precondition_reason` | `TestWriteMechanicalQAReportWritesThePreconditionRefusal/the refusal is recorded as one terminal row and its frontmatter` pins the exact frontmatter prefix and the exact Results block, comparing the row against `spec.QAPreconditionRowID`/`Status`/`Provenance` rather than literals, and requires the pre-contract `| QA-PRECONDITION |` row to be gone |
| The shape the refusal writes is one the stage accepts | `…/the mechanical stage accepts the shape the refusal wrote` runs the real `speccheck.RunMechanicalStage` over the written report in a real Git repository and fails on any `QA-REPORT-SHAPE` finding — the deadlock this Spec exists to end, measured against the writer's own bytes |
| The refusal is recoverable by a reader | `…/the refusal reads back as the refusal that was recorded` reads the written file through `spec.ReadQAReport` and requires verdict `fail`, `RowsBlockedPrecondition = 1`, and `Precondition` equal to the refusal that went in |
| The stage's observations survive the refusal | `…/the refusal costs the report none of the stage's own observations` requires the mechanical finding, its detail, the `## Mechanical rows` heading, and the blocked row to all still be present — while the subtest above proves the Results table still carries exactly one row |
| A gate that did not refuse writes what it wrote today | `…/a gate that refused nothing writes the report it wrote before` takes the same result with `PreconditionRefused` cleared and compares the whole file byte for byte against the pre-Task output. The four pre-existing subtests of `TestWriteMechanicalQAReportRecordsTheRefusal` and `TestRefusedReportDoesNotBlockItsSuccessor` are unchanged and still pass |
| The skill states the contract | `.agents/skills/qa-gate/SKILL.md` §1 (the refusal report, its row, and its reason) and §6 (the closing count rule). `diff -r .agents/skills/qa-gate skills/qa-gate` is clean via `make skills-sync-check`; `make skills-check` → `Roundfix skill check passed: …, qa-gate, …` |

Supporting checks, all after the last edit: `gofmt -l .` → no output; `go build
./...` → clean; `go vet ./internal/daemon ./internal/speccheck ./internal/cli
./internal/spec` → clean (`go vet ./...` reports only the pre-existing `passes
lock by value` findings in `internal/agent`). `go test -count=1 ./internal/daemon
./internal/spec ./internal/speccheck ./internal/cli` → `ok`, `ok`, `ok`, `ok`.
`go test -count=1 ./...` → every package `ok`. `go test -count=1 -tags
repocontract ./internal/speccheck ./internal/baseline` → `ok`, `ok`.
`make skills-sync-check` and `make skills-check` → pass.

### Fixture corrections this behavior forced

Wiring the precondition means the QA gate now refuses on a Spec that fails
`spec check --strict`, and the `internal/cli` implement fixtures were Specs that
do. Two corrections, both to the fixture and not to the rule:

- `writeImplementSpecAtRoot` writes the four Project Constraint rows, and
  `writeImplementSpec` creates the `docs/agents/` guides those rows cite — the
  refusal was four `SC-CONSTRAINT-SOURCE` findings against citations that
  resolved to nothing. The `internal/daemon` fixtures already did both, which is
  why no daemon test moved.
- `TestRunImplementHasNoSpecCheckPrecondition` proves the implement command has
  no consistency precondition, and its premise was that the shared fixture
  carried an error. It now strips its own Spec's Project Constraints and commits
  that, so the premise is a fact the test creates rather than one it inherits
  from a fixture that is deliberately clean.

### Follow-ups for later Tasks

- `go test -count=1 -tags docscontract ./internal/docscontract` fails on
  `TestCheckCorpusGolden` with one `SC-ADR-RELATED` finding in the active Spec
  corpus. Reproduced with the entire working tree stashed, so it predates this
  Task. task_09 runs that suite in its Verification and will have to settle
  whether the golden or the corpus is wrong.
- Identifier strategy: this Task introduced no glossary term. It exported one
  existing Go function (`speccheck.PromoteGaps`) and emits the frontmatter keys
  task_01 minted. Whether `CONTEXT.md` should name
  `rows_blocked_precondition`, `precondition_check`, and `precondition_reason`
  is task_09's slice, as F-2 records.
- The qa-gate skill's `version` is unchanged at `0.0.2`. `make
  skills-version-check` requires exactly one non-empty version and no bump on
  content change; if the repository wants a bump policy, it is not written down
  anywhere this Task could cite.

Changed paths: `internal/daemon/task_engine.go`,
`internal/daemon/task_engine_test.go`, `internal/daemon/daemon_test.go`,
`internal/speccheck/report.go`, `internal/cli/spec_check.go`,
`internal/cli/implement_test.go`, `.agents/skills/qa-gate/SKILL.md`, its
`make skills-sync` mirror `skills/qa-gate/SKILL.md`, and this Task file. No
tooling configuration, derived artifact, or other Spec file was touched.
