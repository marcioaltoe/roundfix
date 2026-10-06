---
spec: 0234-an-openrouter-key-per-stage
prd: _prd.md
created: 2026-10-05
---

# An OpenRouter key per stage — Technical Spec

## Executive Summary

A new leaf package, `internal/openrouterkey`, names each OpenRouter stage's
environment variables in preference order and picks the first one that is
set, returning its name only. The judge maps its OpenRouter transport onto
the judge stage's list, so it reads `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` before
the shared key and still falls back to TypeSafe direct; the implementation
helper Spec 0233 adds delegates to the implementation stage's list. The
judge's report, its text summary and every Judge Log line gain
`key_variable`, and so does 0233's implementation spend record; the Doctor
Command lists both stages' variables by name. The primary trade-off is a new
package for three constants and two functions: it is the one place both
stages and the doctor can import without a cycle, which keeps the order from
being restated in three readers (ADR-0239).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the new
  names are environment variable names and one JSON field, `key_variable`.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the judge still sends the selected
  key only in the authorization header of its transport's endpoint, and the
  implementation stage hands OpenCode the selected variable as Spec 0233
  already does for the shared one. `openrouterkey.Select` returns a name,
  never a value; the judge reads the value of that name from the keys the
  command collected, and its redaction covers every variable of every
  transport. No test or Verification command reaches OpenRouter or TypeSafe:
  the judge's tests use a fake `http.RoundTripper` and the implementation
  tests use the fakes Spec 0233's tests use. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0239 (this Spec): "Each
  OpenRouter stage now reads its own Roundfix variable first and the shared
  key second", and "The generic `OPENROUTER_API_KEY` is never read". It
  supersedes ADR-0201 in part, whose key scope says a later feature "reads
  these two names, never the generic" key; the generic key stays unread.
  The judge's boundary is unchanged, ADR-0201: "No other service receives a
  request, and a run never switches recipient". The judge's ceiling stays a
  User Config value, ADR-0231: "Project Config cannot set it". The
  subscription rule stands and the implementation key changes no selectable
  model, ADR-0235: "OpenAI and Anthropic models run only through the codex
  and claude subscriptions". ADR-0200 and ADR-0209 stay as they are: the
  judge stays advisory and asks the same questions. ADR-0208 is honored by
  adopting nothing, since Spec 0233 owns the one Backlog Entry that shares
  this context. ADR-0229 and ADR-0237 do not apply,
  because no Delivery Queue park or archive behavior changes. ADR-0184 binds Surface Transcripts 1 and 2. ADR-0179
  bounds the Governed Paths. ADR-0189 and ADR-0233 bind each owned skill's
  version raise. The QA gate follows ADR-0080: "QA verdicts distinguish
  environment-blocked rows", and ADR-0091: "required to be terminal and to
  depend on every leaf"; ADR-0096, ADR-0097, ADR-0104, ADR-0167, ADR-0182,
  ADR-0194, ADR-0195 and ADR-0210 bind its stage, row carry, outside
  evidence, Pull Request row, settlement, row record, re-observation and
  evidence snapshot. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and
  ADR-0183 check consistency by citation and receipt. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 changes the Governed Paths
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`; task_02
  changes `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`,
  `skills/write-prd/SKILL.md`, `skills/write-techspec/SKILL.md` and both
  copies of the Roundfix Skill's `SKILL.md`; task_03 changes none. Express
  maintainer authorization: the standing grant of 2026-09-30, "considere
  autorizado a ajustar todas as skills se necessário", and "Concedo" of
  2026-10-05; bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`,
  `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-techspec/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0234-an-openrouter-key-per-stage/_authorization.md`.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Stage key list | new `internal/openrouterkey` | Names, stages, preference order, first-set selection by name |
| Judge transport keys | `Transport.KeyVariables`, `Questions.KeyVariables` in `internal/judge/questions.go` | The OpenRouter transport reads the judge stage's list |
| Judge selection | `selectTransport` in `internal/judge/client.go`, `Run` in `internal/judge/judge.go` | Returns and records the variable used; skip reason lists every variable |
| Judge Log | `logLine` in `internal/judge/log.go` | Gains `key_variable` |
| `spec judge` | `internal/cli/spec_judge.go` | Collects every transport variable, prints the variable in the summary and the help |
| Doctor | `environmentReadiness` in `internal/cli/readiness_toolchain.go` | Lists judge and implementation variables |
| Implementation key | the helper and spend record Spec 0233 adds | Delegate to the implementation stage's list; record `key_variable` |
| Records | `docs/user-guide/commands/spec.md`, `doctor.md`, `configuration.md`, the Roundfix Skill, write-prd, write-techspec | Name the stage keys, the order and the fallback |

## Implementation Design

### Interfaces

```go
// internal/openrouterkey/openrouterkey.go
const (
	Shared    = "ROUNDFIX_OPENROUTER_API_KEY"
	Judge     = "ROUNDFIX_OPENROUTER_JUDGE_API_KEY"
	Implement = "ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY"
)
type Stage string
const (StageJudge Stage = "judge"; StageImplement Stage = "implement")
func Variables(stage Stage) []string                      // stage variable, then Shared
func Select(environ []string, stage Stage) (string, bool)  // name only

// internal/judge
func (t Transport) KeyVariables() []string
func (q Questions) KeyVariables() []string
func selectTransport(q Questions, keys map[string]string) (t Transport, variable, key string, ok bool)
// Report gains:  KeyVariable *string `json:"key_variable"`
// logLine gains: KeyVariable string  `json:"key_variable"`
```

```text
1. Variables(StageJudge) is [Judge, Shared]; Variables(StageImplement) is [Implement, Shared]; any other stage is empty. Each call returns a new slice.
2. Select scans environ entries "NAME=value"; for a name that appears more than once the last entry wins, as readinessKeyIsSet already decides; an empty value is unset. It returns the first of Variables(stage) that is set, and never a value. OPENROUTER_API_KEY and TYPESAFE_API_KEY are never consulted.
3. Transport.KeyVariables returns openrouterkey.Variables(StageJudge) for a transport whose KeyVariable is openrouterkey.Shared, and [KeyVariable] otherwise. questions.json does not change.
4. Questions.KeyVariables concatenates every transport's KeyVariables in transport order, dropping repeats: Judge, Shared, ROUNDFIX_TYPESAFE_API_KEY.
5. selectTransport returns the first transport, and within it the first variable, whose collected key is non-empty. A run never switches transport or variable.
6. Run sets Report.KeyVariable to the selected variable whenever it sets Report.Transport, and writes it to every Judge Log line of the run; a skipped run without a key has key_variable null.
7. Without a key, Report.Skipped is "<first> is not set (nor <second>, nor <third>)" over Questions.KeyVariables.
8. Redaction snapshots the value of every variable in Questions.KeyVariables.
9. runSpecJudgeCommand collects exactly the variables in Questions.KeyVariables from its own environment; no other name is read.
10. The text summary ends "model <model> via <transport> on <variable>", before any "; stopped: ..." suffix.
11. The doctor's environment detail lists "spec judge keys: " over Questions.KeyVariables, then "; implementation keys: " over openrouterkey.Variables(StageImplement), each "<name> set" or "<name> not set".
12. Spec 0233's implementation-key helper returns openrouterkey.Select(environ, openrouterkey.StageImplement), and every place that hands OpenCode the key or checks openrouter.implement_monthly_ceiling_usd uses that one selected variable.
13. Each implementation spend record Spec 0233 writes carries key_variable, the selected variable's name; a record never contains a key value.
```

### Data Models

The Judge Log line (`roundfix/judge-log/v1`) gains one string field,
`key_variable`; older lines without it still parse and still count toward the
month's ceiling. The `roundfix/spec-judge/v1` report gains `key_variable`,
null when the run skipped without a key. Spec 0233's implementation spend
record gains `key_variable`. No Run Database schema changes.

### API Contracts

1. API Contract: `roundfix spec judge` reads, first set wins,
   `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY` and
   `ROUNDFIX_TYPESAFE_API_KEY`; its text summary ends
   `model <model> via <transport> on <variable>`, its JSON report carries
   `key_variable`, and without a key it prints Surface Transcript 1.
2. API Contract: every Judge Log line carries `key_variable`, the variable
   whose key was sent; no line carries a key value.
3. API Contract: `roundfix doctor`'s `environment:` detail ends with
   `spec judge keys: ROUNDFIX_OPENROUTER_JUDGE_API_KEY <state>, ROUNDFIX_OPENROUTER_API_KEY <state>, ROUNDFIX_TYPESAFE_API_KEY <state>; implementation keys: ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY <state>, ROUNDFIX_OPENROUTER_API_KEY <state>`,
   where `<state>` is `set` or `not set`, before any `; next: ...` suffix.
   Its status and other findings are unchanged; doctor's exit code depends on
   the machine, so this contract is proven by the doctor tests rather than a
   transcript.
4. API Contract: implementation on an open model through OpenCode uses
   `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` when set and otherwise
   `ROUNDFIX_OPENROUTER_API_KEY`, and each implementation spend record names
   the variable in `key_variable`.

### Surface Transcripts

1. Surface Transcript: the judge without any key skips and names the judge
   key first.

   ```transcript
   $ roundfix spec judge <slug> --stage prd
   stdout:
   ...
   Judge: skipped: ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set (nor ROUNDFIX_OPENROUTER_API_KEY, nor ROUNDFIX_TYPESAFE_API_KEY); <n> judgment(s) not asked
   stderr:
   exit: 0
   ```

2. Surface Transcript: the judge's help names the keys in order.

   ```transcript
   $ roundfix spec judge --help
   stdout:
   ...
   Keys, first set wins: ROUNDFIX_OPENROUTER_JUDGE_API_KEY, then the shared
   ROUNDFIX_OPENROUTER_API_KEY (both OpenRouter), then ROUNDFIX_TYPESAFE_API_KEY
   (TypeSafe directly). The generic OPENROUTER_API_KEY is not read. The summary
   names the variable used, never its value.
   ...
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal 1 → `Transport.KeyVariables`, `selectTransport`, `Run` (Invariants 3, 5, 7)
- Goal 2 → Spec 0233's helper over `openrouterkey.Select` (Invariants 2, 12)
- Goal 3 → `Report.KeyVariable`, `logLine.KeyVariable`, the summary, the implementation spend record (Invariants 6, 8, 10, 13)
- Goal 4 → `environmentReadiness` (Invariant 11)
- Goal 5 → unchanged ceiling sums; one selected variable per stage (Invariants 5, 12)
- Story 1 → `openrouterkey`, judge selection, implementation helper
- Story 2 → the shared key as second choice (Invariants 1, 2)
- Story 3 → `environmentReadiness` (API Contract 3)
- Story 4 → `key_variable` in the Judge Log and the spend record (API Contracts 2, 4)
- Core Feature 1 → judge selection (API Contract 1; Surface Transcripts 1, 2)
- Core Feature 2 → implementation helper (API Contract 4)
- Core Feature 3 → `openrouterkey` (Invariants 1 to 4, 9)
- Core Feature 4 → `key_variable` fields and the summary (API Contracts 1, 2, 4)
- Core Feature 5 → doctor (API Contract 3)
- Core Feature 6 → guides and skills
- Success Metric 1 → judge and `spec judge` tests with both OpenRouter variables set
- Success Metric 2 → shared-key and no-key tests (Surface Transcript 1)
- Success Metric 3 → implementation helper and spend-record tests
- Success Metric 4 → doctor tests
- Success Metric 5 → the repository Verification at each Task's settlement

## Integration Points

- OpenRouter's System One API and TypeSafe's API, unchanged destinations,
  reached only by a real `spec judge` with a key; never in tests.
- OpenCode, which receives the selected implementation variable exactly as
  Spec 0233 passes the shared one today.

## Testing Approach

- `internal/openrouterkey` gets table tests over environ slices: order,
  empty values, repeated names, and the generic names ignored.
- The judge's selection and log are tested in a new
  `internal/judge/stage_key_test.go` with the existing fixtures and a fake
  transport that asserts the bearer value; existing tests that pin the skip
  reason and the Judge Log's exact field set change with the declared output.
- `spec judge` is tested in `internal/cli/spec_judge_test.go` with
  `commandEnvironment.environ`, never the process environment; the summary
  goldens that end `via openrouter` or `via typesafe` change to name the
  variable.
- The doctor is tested through `environmentReadiness` and the five-line
  doctor test with scripted readiness dependencies.
- The implementation stage is tested beside Spec 0233's helper and spend
  record, with the fakes 0233's own tests use; no test starts OpenCode
  against OpenRouter.

## Build Order

1. Add the stage key list; make the judge read the judge stage first, record
   and print the variable, and describe it in the `spec judge` reference and
   the Roundfix Skill's spec reference.
2. Report each stage's variables in the Doctor Command, and name the stage
   keys in the doctor reference, the configuration guide, the Roundfix
   Skill's runtime reference and the write-prd and write-techspec skills
   (depends on: 1).
3. Make Spec 0233's implementation helper read the implementation stage first
   and its spend record name the variable (depends on: 2, and on Spec 0233
   merged).
4. Final QA gate (depends on: 1, 2, 3).

Steps 1 to 3 run in series because each raises the Roundfix Skill's version
or follows one that does, and step 3's tests live beside Spec 0233's code.

## Risks & Considerations

- Spec 0233 is authored in parallel. Its helper, spend record and ceiling
  check are named here by role, not by file; task_03 locates them in the
  tree it runs on. If 0233 already reports the implementation key in the
  doctor under another label, task_02's line replaces that report.
- With only the shared key set, the export still shows one total. That is
  the documented fallback, not a defect.
- The text summary's new suffix changes output other tools may parse; JSON
  is the stable form, and it only gains a field.
- A ceiling that 0233 checks against OpenRouter's key usage reads the
  selected key's usage, so with the shared key it includes the judge's spend
  as before, and with the implementation key only implementation spend.
- The owned skills' versions are raised by the record command, which picks a
  free version, so a concurrent raise in Spec 0233 resolves at merge
  (ADR-0233).

## Vocabulary Contract

- emits: `internal/cli/spec_judge.go`
  pattern: `Keys, first set wins`
  documented-in: `docs/user-guide/commands/spec.md`
- emits: `internal/cli/readiness_toolchain.go`
  pattern: `implementation keys`
  documented-in: `docs/user-guide/commands/doctor.md`

No glossary term is adopted; "stage key" is used in its plain sense beside
**Judge Log**.

## Decisions

- One leaf package names the stage variables for every reader. See ADR-0239.
- Map the judge's existing OpenRouter transport onto the judge stage's list
  instead of adding a transport to `questions.json`, so the measured
  settings file and its positional tests stay as they are.
- Record the variable in a new `key_variable` field rather than renaming
  `transport`.
