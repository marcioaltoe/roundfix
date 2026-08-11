---
spec: 0024-context-efficient-runs
prd: _prd.md
created: 2026-07-10
---

# Context-Efficient Runs - Technical Spec

## Executive Summary

Roundfix will separate lossless evidence from the compact information needed by
an Agent, a Supervisor, and a Detached Run caller. The Daemon becomes the only
owner of authoritative Verification: successful output is discarded, while a
failure creates a diagnostic artifact and one path-based Verification Feedback
prompt in the same Agent Session before a final verdict. A new `events` command
projects four stable JSONL categories from the existing Run Event Journal, and
the Console Log renders structured ACP reads and edits as one-line summaries
without changing raw journal payloads. Spec Task prompts gain a bounded Spec
Context Bundle whose only embedded document is the assigned Task. The primary
trade-off is retaining failed verification logs as Run artifacts instead of
embedding command output; this costs small local storage but prevents another
unbounded context stream and is reclaimed by existing Run GC.

## System Architecture

The change extends existing ownership boundaries:

- `internal/daemon` owns the two-verdict Verification cycle, diagnostic
  artifact lifecycle, independent Work Item settlement, and Task Context Bundle
  assembly.
- `internal/agent` owns initial and repair prompt contracts plus structured ACP
  console summaries. Its Run Event payloads remain byte-preserving per
  ADR-0008.
- `internal/spec` parses an optional `## Context` path list from Task files so
  `write-tasks` can name relevant interfaces without prose/path guessing.
- `internal/cli` adds `events`, sharing the current attach journal replay and
  follow machinery after it is renamed to a command-neutral component.
- `internal/runevent` owns the public category projection from structured
  Daemon event payloads. The store schema and journal schema do not change.

Verification flow for a Task or review Batch:

```text
initial Agent prompt -> Daemon Verification attempt 1
    pass -> settle
    fail -> retain diagnostic artifact -> Verification Feedback
          -> same Agent Session repairs -> Daemon Verification attempt 2
             -> final settlement (no further repair)
```

If the Agent settles a Work Item failed before Verification because a required
credential or prerequisite is missing, the existing failure path remains
authoritative and no transport retry is inferred. The scheduler continues with
ready independent Tasks and blocks only dependents.

## Implementation Design

### Interfaces

Verification capture:

```go
type VerifyRequest struct {
    WorkDir, Command, OutputPath string
}

type VerifyResult struct {
    OutputPath string
}

type Verifier interface {
    Verify(context.Context, VerifyRequest) (VerifyResult, error)
}

type VerificationCommandError struct {
    Command, OutputPath string
    Err error
}
```

Repair prompt input:

```go
type VerificationFeedback struct {
    Command, DiagnosticPath, Failure string
    Attempt int
}

func BuildVerificationRepairPrompt(
    workItem string, feedback VerificationFeedback,
) (string, error)
```

Spec Context Bundle:

```go
type SpecContextBundle struct {
    PRD, TechSpec, TaskGraph string
    Instructions, Interfaces, PriorChangedFiles []string
    OmittedPriorFiles int
}
```

Task-authored references:

```go
type TaskContextRef struct {
    Kind ContextKind // instruction | interface
    Path string
}
```

Supervisor projection:

```go
type StreamRecord struct {
    Schema, RunID, Category, Time string
    Cursor int64
    Batch int `json:"batch,omitempty"`
    Attempt int `json:"attempt,omitempty"`
    WorkItem, Phase, Status, Verdict, Outcome, Summary string
}
```

The real JSON tags omit empty category-specific fields; `schema` is always
`roundfix-events/v1`. Projection accepts only Daemon events and never exposes
raw `Payload` bytes.

### Data Models

No SQLite migration is required. Failed Verification output is written under
the existing Run artifact tree:

```text
runs/<run-id>/verification/batch-<NNN>-attempt-<1|2>.log
```

`ExecVerifier` writes combined stdout/stderr to a temporary sibling file,
closes it, and removes it on success. On command failure it atomically keeps the
attempt path and returns `VerificationCommandError`. Context cancellation,
process-start failure, and filesystem failure remain infrastructure errors and
never enter the repair policy. The Agent prompt receives the failed command,
wrapped exit failure, and repository-relative diagnostic path; it receives no
log body. Thus passing-package lines can exist in the failure artifact without
entering Agent context automatically. Both failed attempt artifacts remain
lossless until normal Run retention/GC removes the artifact directory.

Every `daemon.verification` payload gains `attempt` (`1` or `2`) and `phase`
(`started`, `command-passed`, `failed`, or `verdict`). A failed command event
includes `diagnostic_path`; it does not embed command output. After the command
sequence, exactly one `phase: verdict` event records `passed` or `failed`; the
attempt-2 verdict, when present, is authoritative. Attempt 2 reruns the complete
configured sequence, not only the command that failed. Existing `daemon.task`,
`daemon.batch`, and `daemon.outcome` payloads remain settlement sources.

Task files gain an optional `## Context` section containing labeled,
backticked, repository-relative paths: `- instruction: <path>` or
`- interface: <path>`. The write-tasks workflow uses it for relevant code
interfaces and task-specific instruction documents. `internal/spec` parses
those entries into `Task.Context []TaskContextRef` while preserving the current
frontmatter and Task Graph contracts. A Task may declare at most 50 unique
context entries; exceeding that bound is a typed Task-file validation error
rather than silent truncation.

The Daemon constructs the bundle deterministically:

1. Standard Spec paths: `_prd.md`, `_techspec.md`, and `_tasks.md` when present.
2. Instruction paths: root `AGENTS.md`, the canonical implement-task Skill,
   and valid Task `## Context` paths classified as instructions.
3. Relevant interfaces: remaining valid Task `## Context` paths.
4. Prior changed files: sorted output of `git diff --name-only <Run HEAD>..HEAD`
   in the Task's actual worktree, which naturally includes only Tasks already
   integrated before that Task Worktree was created.

Paths are cleaned, deduplicated, constrained to the repository, and never read
by the bundle builder. The complete manifest has a 200-path ceiling. Standard
paths and the at-most-50 explicit Task context entries are reserved first;
sorted prior changed files fill the remaining capacity, and the omitted prior
count is stated. The full current Task remains embedded exactly once after the
manifest. PRD, TechSpec, skills, source files, and diffs are path references
only. Replacement or resumed sessions receive the same builder output for
their next Task.

### API Contracts

`roundfix events <run-id>` is a new read-only command:

- Required positional Run ID; no implicit newest/active Run selection.
- Replay starts at journal cursor `0` and emits one JSON object per stdout line
  in cursor order.
- `--follow` continues after replay while the Run is Active and exits `0` after
  the Run becomes terminal and the final journal page has drained.
- `--filter <csv>` accepts a non-empty subset of `task-status`, `batch`,
  `verification`, and `outcome`. Omitted means all four. Duplicates are
  deduplicated; unknown or empty categories exit `2`.
- Missing/unknown Run ID exits `2`. Store or malformed relevant Daemon payload
  errors exit `1`; the command never guesses from human summary text.
- SIGINT/SIGTERM during follow exits `130` without a stdout trailer. Diagnostics
  and validation errors go only to stderr.

Stable category mapping:

| Category | Journal source | Projected fields |
| --- | --- | --- |
| `task-status` | `daemon.task` | `batch`, `work_item`, `phase`, `status`, `summary` |
| `batch` | `daemon.batch` | `batch`, `phase`, `summary` |
| `verification` | `daemon.verification` | `batch`, `work_item`, `attempt`, `phase`, `verdict`, `summary` |
| `outcome` | `daemon.outcome` | `outcome`, `summary` |

Each line also carries `schema`, full `run_id`, journal `cursor`, `category`,
and RFC3339Nano UTC `time`. `status` carries Task settlement; `verdict` is
`passed`/`failed` only on aggregate Verification verdict events. The raw command
and diagnostic path remain in the journal event but are not projected into the
Supervisor stream.

The current `attachFollower`/replay functions are renamed to a neutral
`journalFollower` and reused by attach and events. The follower keeps the
existing store `data_version` wake-up strategy and cursor paging; events adds a
JSON encoder sink and terminal-exit policy rather than a polling monitor.

Console rendering changes only the presentation derived from structured ACP
updates:

- ACP `kind: read` plus its location and output line count renders
  `read <path> (N lines)`.
- ACP diff content uses `path`, `oldText`, and `newText` to render
  `edit <path> (+N/-M)`.
- Tool updates lacking enough structured metadata render one marker naming the
  tool/path and state, never raw input, raw output, file bodies, ACP JSON, or a
  unified diff.

`StreamUpdate` therefore gains Tool Kind/location metadata and diff line
counts. `ConsoleText` and the Live Run View timeline use the same compact
formatter. `newAgentRunEvent` still stores the original ACP payload unchanged;
only its bounded human summary uses the compact line. Roundfix does not pass
acpx `--suppress-reads`, because doing so would destroy the journal's lossless
input before sinks can choose their representation.

Initial Task and review prompts stop requiring the Agent to run the final
configured Verification. They allow focused checks while working and state
that the Daemon owns the authoritative gate. On attempt-1 failure, the repair
prompt names only the assigned Work Item, failure, and diagnostic path, directs
the Agent to inspect only relevant sections, and forbids commit/push. It uses
the same `SessionRef`. The Daemon never emits a second Batch-start event for a
repair turn.

## Coverage Map

- Goal 1 and Stories 1-2 -> captured Daemon Verification, failed diagnostic
  artifact, same-session repair, ADR-0038 attempt bound
- Goal 2 and Story 3 -> `events`, journal projection, four stable categories
- Goal 3 and Story 4 -> structured ACP compact formatter, unchanged raw payload
- Goal 4 and Story 5 -> Task `## Context`, bundle builder, prior-file resolver
- Story 6 -> unchanged ready-Task scheduler and final Task settlement
- Story 7 -> usage guide, canonical Roundfix Skill, embedded Skill sync
- Core Feature 11 -> deterministic bundle regeneration for every Task prompt

## Integration Points

- **Shell Verification commands** remain verbatim `sh -c` execution in the
  worktree. The only boundary change is file capture instead of streaming.
- **Git** supplies prior changed paths through `diff --name-only` against the
  Run's persisted initial `HeadSHA`; no diff content enters the prompt.
- **acpx/ACP** supplies structured Tool Kind, locations, and diff old/new text.
  Raw JSON remains the Run Event Journal payload under ADR-0008.
- **SQLite Run Event Journal** remains the source for replay/follow. The new
  command is a projection, not a second event store.

## Testing Approach

Verifier tests use real shell commands and temporary artifact directories to
prove: success removes captured output; failure retains combined output;
filenames are deterministic per attempt; cancellation and write failures are
wrapped correctly. Engine fakes then cover both Task and review Batch flows:
pass on attempt 1; fail then pass; fail twice; repair Agent error; Stop Request
during either verify or repair; exactly one repair prompt; same SessionRef; and
no passing output in any prompt.

Scheduler tests preserve the missing-credential case: one Task settles failed,
its dependents do not start, and an independent ready Task does start. No test
reclassifies that result as a transport failure.

`internal/spec` table tests cover Context section parsing, absent sections,
deduplication, invalid/external paths, the 50-entry validation bound, and
byte-preserving status rewrites.
Bundle tests use a real temporary Git repository to verify initial Head range,
sorted prior files, a parallel Task Worktree base, the 200-path ceiling and
omitted count, exactly one full Task, and zero embedded larger documents.

Run Event Stream tests pin complete JSONL lines, schema/category mapping,
filter validation, replay order, replay-to-follow handoff without duplicates,
terminal drain, stdout/stderr separation, cancellation, unknown Runs, and
malformed relevant payloads. Shared follower tests continue to cover attach.

Agent stream fixtures include 31 structured edits and 330 reads and assert the
exact summary counts plus absence of old/new text, raw output, and file bodies.
A separate assertion round-trips the original journal payload byte-for-byte.

Documentation examples are executed through the CLI fixture and parsed one
line at a time with `encoding/json`. The final task runs `rtk make verify`,
`make skills-sync`, and the zero-drift Skill check.

## Build Order

1. Verification capture foundation: artifact-backed `Verifier`, attempt-aware
   Daemon events, and success/failure artifact tests.
2. Task and review one-repair cycles plus revised initial/repair prompts and
   independent-settlement regression tests (depends on: 1).
3. Stable event projection and shared journal follower, then the `events`
   command with replay/filter/follow tests (depends on: 1 for final event shape).
4. Structured ACP Console Log compaction with lossless-payload regression
   fixtures (depends on: none; can run in parallel with 1-3).
5. Task Context path contract in `internal/spec`, prior-change Git resolver,
   bounded Spec Context Bundle, and prompt tests (depends on: none; can run in
   parallel with 1-4).
6. Shipped guidance and workflow sync: update write-tasks guidance for Task
   `## Context`, update implement-task/roundfix guidance, README/usage recipes,
   run `make skills-sync`, and verify documented JSONL examples (depends on:
   2, 3, 4, 5).

## Risks & Considerations

- Failed command output can contain credentials. It already flowed to the
  Console Log; moving it to a Run artifact reduces exposure but does not make it
  secret. Existing artifact permissions and retention remain binding, and the
  event stream never emits the body.
- ACP adapters may omit structured Tool Kind or diff text. The formatter must
  degrade to metadata-only lines, never fall back to raw payload echoing.
- A 200-path bundle can omit useful prior files in very broad Specs. Explicit
  Task Context paths have a separate 50-entry validation bound and are reserved
  first; the omitted count tells the Agent when targeted Git inspection may be
  needed.
- `git diff <initial>..HEAD` includes all prior integrated changes, including
  generated files. Sorting and the ceiling keep the prompt deterministic; the
  write-tasks author must name indispensable interfaces explicitly.
- The events schema becomes public automation API. Field names, category
  meanings, and exit behavior require byte-pinned CLI tests before release.

## Decisions

- The Daemon owns Verification and permits one same-session repair. See
  ADR-0038.
- Failed output is a path-addressed Run artifact; no command body is embedded
  in Verification Feedback.
- `events` is a stable JSONL projection of one explicit Run, not raw journal
  access.
- Console compaction happens after lossless journal capture; acpx read
  suppression is not enabled. See ADR-0008.
- Task `## Context` supplies explicit interfaces; Git supplies prior changed
  files; prose is never scraped for paths.
- The bundle ceiling applies only to prior changed files, never the assigned
  Task or explicit context.
- Canonical and embedded Roundfix Skills ship in sync with the new behavior.
