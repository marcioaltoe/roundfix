---
spec: 0211-a-delivery-queue-that-finishes-without-intervention
prd: _prd.md
created: 2026-10-01
---

# A delivery queue that finishes without intervention — Technical Spec

## Executive Summary

Each of the four interventions of 2026-10-01 has one local cause. The archived
retry reads a Run's start head through a lookup that compares the Run's
working root with the queue's checkout root, and a Run the queue owner starts
always has the item worktree as its working root. The checking stage moves to
merge when every check `gh pr checks` lists passes, and a re-run that has not
registered yet is simply absent from that list. A merge GitHub refuses becomes
a generic error, which the owner parks as `delivery-error`. The agent
environment is the process environment with two variables removed, so a
`NODE_OPTIONS` preload that no longer exists reaches every Node process.

The fix matches the Run by its repository, reads GitHub's merge state beside
the checks and waits while it is blocked, turns a base-branch policy refusal
into a return to checking bounded to one per head per owner, prints a refused
retry as a refusal, and filters missing preloads out of the agent environment
with a notice. The trade-off is one extra JSON field on a read the queue
already makes, and an agent environment that can differ from the user's by a
preload that could not have loaded.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Items keep their
  Spec slug and Runs their identifiers. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — GitHub is still read only through the
  user's authenticated `gh`; the change adds `mergeStateStatus` to an existing
  `gh pr view --json` read. Tests replace `gh` with the existing scripted
  command runner and never reach GitHub. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0211 (this Spec) governs the agent
  environment. ADR-0140, ADR-0154, ADR-0161, ADR-0170, ADR-0192, ADR-0193 and
  ADR-0199 hold unchanged, as the PRD records, and ADR-0187 and ADR-0189
  govern the Roundfix Skill edit. ADR-0184 applies: the changed command
  surface is stated as Surface Transcripts. The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155 and ADR-0167, and ADR-0212 does not
  apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths; express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário"; bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0211-a-delivery-queue-that-finishes-without-intervention/_authorization.md`.

## System Architecture

No new package, command or flag.

| Component | Where | Change |
| --- | --- | --- |
| Run start lookup | `commandDeliveryWorkflow.RunStart` in `internal/cli/deliver_workflow.go` | Matches the Run by repository root |
| Pull Request read | `GitHubCLI.pullRequest`, `pullRequestJSONFields` and `parsePullRequest` in `internal/delivery/github.go` | Reads `mergeStateStatus` into `PullRequest.MergeState` |
| Check report | `GitHubCLI.CurrentHeadChecks` | Carries the after-read `MergeState` in `CheckReport` |
| Merge | `GitHubCLI.MergePullRequest` | Types a base-branch policy refusal as `MergePolicyRefusalError` |
| Checking stage | `Engine.checkCandidate` in `internal/delivery/engine.go` | Waits while the merge state is `BLOCKED` or `UNKNOWN` |
| Merging stage | `Engine.mergeCandidate` | Returns a first policy refusal per head to `checking` |
| Retry refusal output | `runDeliverRetry` in `internal/cli/deliver.go` | Prints `Retry refused` for an `Engine.Retry` refusal |
| Agent environment | `ACPXRunner.baseEnv` in `internal/agent/codex_spawn.go`, new `internal/agent/node_options.go` | Drops missing preloads and prints a notice |
| Skill and guides | Roundfix Skill `deliver` and `runtime` references; `deliver` and `doctor` command guides | Describe the new behavior |

## Implementation Design

### Interfaces

```go
// internal/delivery/github.go
type PullRequest struct {
	// existing fields unchanged
	MergeState string // GitHub mergeStateStatus, upper case; empty when gh omits it
}

type CheckReport struct {
	// existing fields unchanged
	MergeState string
}

// MergePolicyRefusalError reports that GitHub refused the merge because the
// base branch policy still prohibits it.
type MergePolicyRefusalError struct{ Stderr string }

func (err MergePolicyRefusalError) Error() string

// internal/agent/node_options.go
// agentNodeOptions returns the NODE_OPTIONS value to give an agent process and
// the preload paths it removed. exists reports whether a path names an
// existing file system entry.
func agentNodeOptions(value string, exists func(string) bool) (kept string, dropped []string, ok bool)

// internal/agent/acpx_runner.go
type ACPXRunner struct {
	// existing fields unchanged
	Notices io.Writer // nil writes notices to os.Stderr
}
```

### The Run start lookup

`RunStart(ctx, gitRoot, runID)` reads the Run as today and accepts it when
it is an Implement Run with a non-empty start head whose repository is the
queue's: `run.RepositoryRoot` equals `roundconfig.RepositoryRoot(gitRoot)`.
When the Run row has an empty `RepositoryRoot`, written before that column
existed, the lookup derives it with `roundconfig.RepositoryRoot(run.GitRoot)`.
The refusal message keeps its text,
`Run %q has no Implement start head for repository %q`, for an absent Run, a
non-Implement Run, an empty head or another repository. `Engine.Retry` is
unchanged: with the lookup fixed, the archived branch anchors at the start
head, proves descent and returns the item to `reviewing`.

### The merge state

`pullRequestJSONFields` gains `mergeStateStatus`, and `parsePullRequest`
stores it upper-cased and trimmed in `PullRequest.MergeState`.
`CurrentHeadChecks` copies the after-read's `MergeState` into the report.

In `checkCandidate`, after the bucket loop, an item that is otherwise ready
to merge stays pending when `report.MergeState` is `BLOCKED` or `UNKNOWN`.
Every other value, the empty value included, keeps today's decision, so a
`gh` that omits the field and every existing fake behave as before. The
existing deadline bounds the wait, and its expiry parks `checks-timeout` as
today. The log line `roundfix: checks: Delivery Queue item <slug>: merge state
<state>` is written once per change of state.

### The merge refusal

`MergePullRequest` returns `MergePolicyRefusalError` when the merge command
exits non-zero and its standard error contains
`base branch policy prohibits the merge`, case-insensitively. Every other
failure keeps its `commandFailure` error.

`mergeCandidate` handles that error before the generic path. The Engine keeps,
for the life of the owner process, the set of `<slug>@<head>` pairs that met a
policy refusal. On the first refusal for a pair it logs
`roundfix: merge refused by branch policy: Delivery Queue item <slug>; checking again`,
leaves the merge intent unmatched for the next attempt, and sets the stage to
`checking`. A second refusal for the same pair returns the error, which parks
`delivery-error` as today. A new owner process starts with an empty set.

### The retry of a `delivery-error` park

The retry's archived branch already returns an archived item whose head
equals its candidate and whose Pull Request is recorded to `checking`. The
2026-10-01 refusal's reason line was not preserved. task_03 therefore pins the
resumption with a test that parks an item `delivery-error` through a refused
merge and retries it, and pins that every refusal the retry raises for a
`delivery-error` item names its rule. Should the reproduction expose a refusal
on that path, task_03 removes it within this contract.

### The refusal output

`runDeliverRetry` prints an `Engine.Retry` error through a new
`printDeliverRetryRefusal`. Parse, configuration and store errors keep
`printDeliverFailure`. The refusal reads the queue item once, read-only, to
name its stage and blocker; when that read fails, the `Item:` block is
omitted. The exit code stays `2` (`exitPreflight`). Surface Transcript 1 is
the exact shape: the heading, `Reason:` with the refusal, `Item:` with
`stage: <stage>; blocker: <blocker>`, and `No side effects:`, separated by
blank lines, with no `Usage:` block.

### The agent environment

`agentNodeOptions` splits the value into words on unquoted spaces; a double
quote groups a word and a backslash escapes the next character inside
quotes. A value with an unbalanced quote returns `ok` false, and the caller
passes it through unchanged. A preload is the word pair `--require <v>`,
`-r <v>` or `--import <v>`, or the single word `--require=<v>` or
`--import=<v>`. A preload is dropped when `<v>` is an absolute path, or a
`file://` URL whose path is absolute, and `exists` reports it missing. A
relative path or a package name is kept. Kept words are joined by one space,
and a word holding a space is written back in double quotes. When nothing is
kept, the variable is removed.

`ACPXRunner.baseEnv` applies it to the last `NODE_OPTIONS` entry of the base
environment, through `os.Stat`, for both the process environment and an
explicit `Environment`. Every path that reads the base environment, the
`acpx` probe, the adapter checks, the profile proof and every session, uses
`baseEnv`, so they all receive the same result. For each dropped path the
runner writes, once per process, to `Notices` or standard error:

```text
roundfix: notice: NODE_OPTIONS preload "<path>" does not exist; Roundfix left it out of the agent environment
```

The queue owner's child `roundfix implement` inherits the user's environment
as today and applies the same filter in its own runner.

### Data Models

No Run Database change. `PullRequest` and `CheckReport` gain one in-memory
field each.

### API Contracts

1. API Contract: `roundfix deliver retry <slug>` refusal — stderr as Surface
   Transcript 1, exit `2`. Argument errors keep the `Preflight failed` form
   with its usage hint.
2. API Contract: the checking stage merges only when no listed check is
   pending, failed or cancelled and the merge state is neither `BLOCKED` nor
   `UNKNOWN`; it otherwise waits until the existing checks timeout.
3. API Contract: a merge refused with GitHub's
   `base branch policy prohibits the merge` returns the item to `checking`
   once per head per owner process; a second refusal parks `delivery-error`.
4. API Contract: the agent environment notice — the line of "The agent
   environment", on standard error, once per dropped path per process; the
   agent process receives the remaining options in their original order.

### Surface Transcripts

1. Surface Transcript: a refused retry. The fixture queue holds item
   `0300-example` parked `delivery-error` on an archived candidate
   `1111111111111111111111111111111111111111`, and its item head was moved to
   `2222222222222222222222222222222222222222` by hand.

   ```transcript
   $ roundfix deliver retry 0300-example
   stdout:
   stderr:
   Retry refused

   Reason:
     retry Delivery Queue item "0300-example": archived item head "2222222222222222222222222222222222222222" differs from candidate head "1111111111111111111111111111111111111111"

   Item:
     stage: parked; blocker: delivery-error: merge pull request: read pull request before merge: gh failed

   No side effects:
     Roundfix did not change the Delivery Queue item, start a queue owner, commit, or push.
   exit: 2
   ```

The agent environment notice has no transcript of its own. No command prints
it alone: it is one extra standard-error line of whichever command starts an
agent process, such as `roundfix doctor` or `roundfix implement`, whose other
output depends on the machine. API Contract 4 fixes its bytes, task_04's tests
assert them through the runner seam, and the QA gate runs the built
`roundfix doctor` with a fake `acpx` to observe the line.

## Vocabulary Contract

No new glossary term. The emitted words are `Retry refused`, `Item:`,
`merge refused by branch policy`, `merge state` and the notice line above.
task_01 documents each in the `deliver` and `doctor` command guides and the
Roundfix Skill.

## Coverage Map

- Goal 1 → The Run start lookup; The merge refusal; The retry of a `delivery-error` park.
- Goal 2 → The merge state; API Contract 2.
- Goal 3 → The agent environment; API Contract 4.
- User Story 1 → The Run start lookup.
- User Story 2 → The merge state.
- User Story 3 → The merge refusal; The retry of a `delivery-error` park.
- User Story 4 → The refusal output; Surface Transcript 1.
- User Story 5 → The agent environment; API Contract 4.
- Core Feature 1 → The Run start lookup.
- Core Feature 2 → The merge state; API Contract 2.
- Core Feature 3 → The merge refusal; API Contract 3.
- Core Feature 4 → The retry of a `delivery-error` park.
- Core Feature 5 → The refusal output; API Contract 1; Surface Transcript 1.
- Core Feature 6 → The agent environment; API Contract 4.
- Core Feature 7 → Build Order 1.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 2; Testing Approach 3.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- **Delivery Queue store.** No schema change; the retry refusal reads the item
  through the existing read-only accessor.
- **ADR-0140's proof.** The profile proof runs on the filtered environment, so
  a preload that cannot load no longer causes a rejection; every other
  rejection still fails Preflight without fallback.
- **Spec 0212.** It edits the Roundfix Skill after this Spec and names this
  Spec in its Task Graph's `requires`.

## Testing Approach

1. **Run start**, in `internal/cli`, new file
   `internal/cli/deliver_archived_retry_test.go`, over a disposable repository
   with a linked item worktree and a disposable Run Database. A Run created
   with the linked worktree as `GitRoot` yields its start head to
   `RunStart(ctx, <checkout>, id)`; a Run of another repository is refused
   with the unchanged message. An end-to-end test drives `Engine.Retry` of a
   real `delivery.NewEngine` whose Workspace, Recovery and History are the
   command workflow, with a fake Pull Request boundary, from an item parked
   `qa-environment-partial` whose Spec the in-process archive command
   archived with `--qa-override`, to stage `reviewing`.
   `TestOperatorArchiveHistoryReadsTheRunStartAndAncestry` and every
   `internal/delivery` operator-archive test pass unedited.
2. **Checks and merge**, in `internal/delivery`, new file
   `internal/delivery/merge_state_test.go`, over the scripted command runner
   and the fake Pull Request boundary with the fake clock: the merge state is
   parsed from `gh pr view`; `BLOCKED` and `UNKNOWN` with every check passing
   issue no merge; `BLOCKED` until the deadline parks `checks-timeout`;
   `BLOCKED` then `CLEAN` merges exactly once; a policy-refusal stderr is typed
   and returns the item to `checking`, and the next clean read merges with no
   park; a second refusal at the same head parks `delivery-error`; an item
   parked `delivery-error` by a refused merge and retried resumes `checking`
   and merges once; an empty merge state behaves as today. Every existing
   `internal/delivery` test passes unedited.
3. **Retry output**, in `internal/cli/deliver_retry_test.go`, through the
   existing `retryCommandDeliveryEngine` fake: an engine refusal prints
   Surface Transcript 1's stderr with no `Usage:` line and exits `2`; an
   unexpected argument still prints `Preflight failed` and `Usage:`.
4. **Agent environment**, in `internal/agent`, new file
   `internal/agent/node_options_test.go`: the splitter's cases (equals and
   space forms, `-r`, `--import` with a `file://` URL, quoted paths with a
   space, a package name, a relative path, an existing path, an unbalanced
   quote, every option dropped); and through an `ACPXRunner` with an explicit
   `Environment`, a `Notices` buffer and a fake `acpx` script that prints its
   environment, the probe receives only `--max-old-space-size=4096`, the
   notice appears once across two probes, and `CODEX_PATH` and `CLAUDECODE`
   are still stripped. `TestACPXCommandEnvStripsGuardAndHygieneVariables`
   passes unedited.

## Build Order

1. The Roundfix Skill `deliver` and `runtime` references, their mirrors, the
   version record, and the `deliver` and `doctor` command guides, written
   from this TechSpec, task_01 (depends on: none). The guides come first
   because every Task that changes a command names its guide in itself or a
   dependency, and one Task raising the skill version once keeps the version
   record simple.
2. The Run start lookup and its tests, task_02 (depends on: 1).
3. The merge state, the merge refusal, the retry resumption and the refusal
   output with their tests, task_03 (depends on: 1).
4. The agent environment filter and notice with its tests, task_04 (depends
   on: none).
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A merge state that stays blocked for another reason.** A repository that
  requires an approving review reports `BLOCKED` until someone approves. The
  item then parks `checks-timeout` instead of failing its merge; the log line
  names the merge state, so the park is explainable.
- **The 2026-10-01 retry refusal is not reproduced from its own record.** Its
  reason line was lost, so the fix rests on the reproduced path and on the
  merge refusal no longer parking at all.
- **Path-only preload detection.** A relative preload that resolves to a
  missing file still fails its process, as it does today; only absolute paths
  and `file://` URLs are checked, which covers the observed case.
- **Notice frequency.** Once per process per path keeps a long-running owner
  from repeating the line, and every child `roundfix implement` prints it once
  in its own output.

## Decisions

- Match the Run by repository, as the Run Database's other listings do.
- Wait on GitHub's merge state rather than reading required-check
  configuration, which needs another API and permissions.
- Bound the policy-refusal return to one per head per owner process, so a
  refusal that persists still parks.
- Drop missing preloads with a notice (ADR-0211), the maintainer's preferred
  option of 2026-10-01.
