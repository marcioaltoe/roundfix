---
status: completed
type: backend
---

# Task: Detect Precondition Failure

Detect which precondition failed and why.

## Work
- Parse `spec check --strict` output
- Extract error codes
- Extract check name and reason
- Store for report writing

## Verification
- `grep -q "PreconditionRefusal\|preconditionRefusal" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck 2>&1 | grep -q "^ok"`


## References

- User Story 2: Precondition captured
- Core Feature 1: Terminal Row Writing

## Result

The mechanical stage now names the precondition that stopped it. `speccheck.PreconditionRefusal`
(`internal/speccheck/mechanical.go:131`) turns the gate's already-classified
strict Spec check result into the `spec.PreconditionRefusal` task_01 writes
reports from: the refusing check in `GatePreconditionCheck`
(`internal/speccheck/mechanical.go:87`, `spec check --strict`) and, as the
reason, every blocking finding rendered as its durable code followed by the
sentence that explains it — `SC-REQUIREMENT-CONTRADICTORY: _prd.md requires one
report per run and forbids writing one`, joined with `; ` when more than one
check refused. `RunMechanicalStage` stores the pair on its result
(`internal/speccheck/mechanical.go:423`), so the writer never has to
reconstruct why the gate stopped from whichever findings the result happens to
carry.

Four decisions are deliberate:

- **The check result is read where it is produced, not re-parsed from rendered
  text.** `spec check --strict` renders the same `speccheck.Result` the gate
  already holds, so parsing its text back would only add a way for a reworded
  line to change what the report records.
- **The code leads each reason.** It is the durable name a later reader and
  detector share; the sentence beside it may be reworded, the token may not
  (`CONTEXT.md`, Mechanical Refusal Code).
- **Every distinct cause survives; only an exact repeat is dropped.** The
  reason is the refusal report's sole record of why the gate stopped, so
  truncating it would delete the evidence the row exists to carry.
- **A refusal whose cause has no name is still a refusal.** `PreconditionRefused`
  carries presence separately from content, so an unnamed cause degrades to
  task_01's recorded placeholder rather than to a run that looks like it never
  refused — a refusal that cannot be written is the deadlock this Spec exists to
  end.

`Blocking` and `PreconditionRefused` stay distinct: the first says this
mechanical run refuses for any cause, the second says the gate's precondition
is what refused.

### Evidence per acceptance criterion

Red first: with the tests written and the contract absent, `go test -run
'TestGateRefusalNamesThePreconditionThatStoppedIt|TestMechanicalStageStoresThePreconditionRefusalForTheReport'
-count=1 ./internal/speccheck` failed to build — `undefined:
speccheck.PreconditionRefusal`, `undefined: speccheck.GatePreconditionCheck`,
and `MechanicalResult has no field or method PreconditionRefused`.

Focused check after the last edit — the same `-run` selection → `ok
roundfix/internal/speccheck 0.765s`, all seven subtests passing.

| Work item | Evidence |
| --- | --- |
| Parse `spec check --strict` output | `a strict Spec check refusal names its check and every refusing code` runs the real checker over the `constraint-missing` fixture Spec through `checkFixture`, classifies it with `GatePrecondition`, and requires the refusal to name `GatePreconditionCheck` |
| Extract error codes | the same test asserts every blocking finding's code reaches the reason, and `every distinct cause survives and a repeat is recorded once` pins the exact joined reason for two distinct codes given three findings |
| Extract check name and reason | the same two tests compare against `speccheck.GatePreconditionCheck` and the findings' own summaries rather than copied literals; `a refusal reason stays on one line` proves a summary carrying a newline and a tab still yields one frontmatter value; `a refusal whose cause has no name still records the refusing check` proves an unnamed cause records the check and invents no reason |
| Store for report writing | `a refusing precondition reaches the refusal report writer` runs `RunMechanicalStage` against a real temporary Git repository and feeds `result.PreconditionRefusal` straight into `spec.WritePreconditionRefusalReport`, then reads `precondition_check: "spec check --strict"`, the refusing code, its sentence, and `\| 0 \| blocked \| precondition \|` back out of the rendered report; `a passed precondition stores no refusal` holds the zero value when nothing refused, and `a gate input alone refuses nothing` proves a Spec's own declared term stays a repair input rather than becoming a refusal |

Supporting checks, all after the last edit: `gofmt -l internal/speccheck/
internal/spec/` → no output; `go build ./...` → clean; `go vet
./internal/speccheck ./internal/spec` → clean; `go test -count=1
-tags repocontract ./internal/speccheck` → `ok`. Dependents of the changed
`MechanicalResult`: `go test -count=1 ./internal/spec ./internal/daemon
./internal/cli` → `ok`, `ok`, `ok`.

Changed paths: `internal/speccheck/mechanical.go`,
`internal/speccheck/report.go`, `internal/speccheck/mechanical_test.go`, and
this Task file. No tooling configuration, derived artifact, or other Spec file
was touched.

### Follow-ups for later Tasks

- No production caller feeds the gate precondition into the stage yet:
  `qaMechanicalRequest` (`internal/daemon/task_engine.go:2044`) leaves
  `MechanicalRequest.Precondition` at its zero value, and
  `writeMechanicalQAReport` still writes its own `| QA-PRECONDITION | fail |
  mechanical refusal: … |` row instead of calling
  `spec.WritePreconditionRefusalReport`. Both ends of that wiring are Daemon
  report-lifecycle work, outside this Task's detection slice; no Task in
  `_tasks.md` currently claims them, so the Spec needs one before the refusal
  path reaches a real run.
- `detectMechanicalReportShape` (`internal/speccheck/mechanical.go:1419`) would
  still reject the `blocked`/`precondition` row this refusal produces, and its
  `mechanicalBlockedCounts` does not track `rows_blocked_precondition`. That is
  task_04, which the graph already sequences next.
- Identifier strategy: this Task introduced no glossary term. `spec check
  --strict` names an existing CLI surface and the finding codes it records are
  the established `SC-*` tokens, so the `QA Report` entry in `CONTEXT.md` is
  unchanged by this slice; the frontmatter fields task_01 added remain the
  closing node's check.
