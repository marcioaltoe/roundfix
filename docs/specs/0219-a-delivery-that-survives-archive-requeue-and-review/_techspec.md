---
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
prd: _prd.md
created: 2026-10-03
---

# A delivery that survives archive, requeue and review — Technical Spec

## Executive Summary

Four existing seams carry the whole change; no package, store table or flag
is added. The Spec Consistency Check gains two findings, a file that pins an
active Spec's directory and a truncated Verification span, both run from the
reference and Verification detectors in `internal/speccheck/citations.go`, so
the governed `constraints.go` and `coherence.go` stay untouched, and the
Daemon's Settlement Check (ADR-0182) turns the first into a repair the Task
that caused it makes. A Task Context path follows its Spec into the archive
root by resolution, not by rewrite. The Delivery Retry and `deliver start`
keep the item's work: a descendant head after a post-archive correction goes
back to review, and a new queue adopts the one item branch that holds work.
The review gains a fifth Delivery Convention, and the one Verification prober
(ADR-0148) asks `sh -n` before it runs a command. The trade-off accepted is
that the archive refuses instead of repairing: an operator who commits a pin
outside a Run still fixes it by hand, but the archive commit stays the exact
Spec move the queue already proves, and no code is ever rewritten for a move.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. The two finding
  codes join the existing `SC-` family, and items keep their Spec slug and
  item branch. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every change reads local files,
  local Git and, for the review, the configured runtime through the existing
  sealed validator; no credential and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — the PRD's row holds. ADR-0223:
  "Roundfix never rewrites code to follow an archive" and ADR-0226:
  "A command the shell rejects is reported as malformed" set the design.
  ADR-0148 keeps one prober, ADR-0148: "the vacuous/failing/unknown classification lives in one extracted prober",
  so the parse check lives in that prober and both callers inherit it.
  ADR-0182 carries the pin finding to the Task, ADR-0182: "the refusing Spec Consistency findings that were".
  ADR-0165 keeps a blocking review after archive parked for a corrective Spec,
  ADR-0165: "Publication parks as".
  ADR-0196 keeps validation fail-closed, ADR-0196: "Every other finding stands, and only a standing finding parks".
  ADR-0154's override, ADR-0197's two rounds, ADR-0170, ADR-0192, ADR-0193,
  ADR-0199, ADR-0211 and ADR-0120's history root hold unchanged; ADR-0187
  and ADR-0189 govern the skill edits; ADR-0184 asks for the transcripts
  below. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill files are Governed Paths
  under the maintainer's skill authorization ("considere autorizado a ajustar
  todas as skills se necessário"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## System Architecture

- **Spec Consistency Check** (`internal/speccheck`): a new file
  `active_spec_paths.go` holds the pin detector and the archive resolution of
  a Task Context path; `citations.go` calls both. `verification.go` gains the
  truncated-span detector, also called from `citations.go`.
- **Archive Command** (`internal/cli/archive.go`): calls the pin detector
  before `spec.Archive`, so a refusal moves nothing. `internal/spec/archive.go`
  is not changed.
- **Delivery Engine** (`internal/delivery/engine.go`): the archived branch of
  `Retry` accepts a descendant head for every blocker other than an
  override-less `qa-environment-partial`.
- **Delivery workflow** (`internal/cli/deliver_workflow.go`, `deliver.go`):
  one helper lists a Spec's item branches with work; `deliver start` refuses
  on two or more, and `CreateItemBranch` adopts the single one.
- **Pre-PR review** (`internal/cli/review_conventions.go`): convention C5 and
  its region.
- **Verification prober** (`internal/daemon/verification_probe.go`) and the
  **spec check command** (`internal/cli/spec_check.go`): the parse check, the
  `malformed` verdict and the uncommitted-source report.

## Implementation Design

### Interfaces

```go
// internal/speccheck/active_spec_paths.go
const CodeSpecPathPinned = "SC-SPEC-PATH-PINNED"

type SpecPathPin struct {
	Path string // repository-relative, slash-separated
	Line int
}

// ActiveSpecPathPins lists non-Markdown files outside the Spec Root, its
// archive root and docs/history that name <spec-root>/<slug>.
func ActiveSpecPathPins(repoRoot, specsRoot, slug string) ([]SpecPathPin, error)

// internal/speccheck/verification.go
const CodeVerifyTruncated = "SC-VERIFY-TRUNCATED"
func TruncatedVerification(taskFile string, content []byte) []Finding

// internal/daemon/verification_probe.go
var ErrVerificationMalformed = errors.New("the shell cannot parse the command")

// internal/cli/deliver_workflow.go
func existingItemBranches(ctx context.Context, git preflight.GitRunner, gitRoot, base, specSlug string) ([]string, error)
```

```text
1. A pin is a line of a file that contains "<spec-root-relative>/<slug>" followed by "/" or the end of a word.
2. A Markdown file, any file under the Spec Root, its archive root or docs/history, and an ignored file are never pins.
3. A Task Context path resolves through the archive only when its active form is missing and its archived form exists.
4. An archived retry moves to reviewing only when its head descends from its newest candidate head.
5. A corrective-spec-required park and an override-less qa-environment-partial park never accept a moved head.
6. deliver start adopts a branch only when exactly one item branch of the slug has commits the default branch lacks.
7. A command the shell rejects is never executed by the prober.
```

### The pin detector

In a Git work tree, `ActiveSpecPathPins` runs `git grep -n -I -F --untracked
-e <spec-root-relative>/<slug>` with exclude pathspecs for `*.md`, the Spec
Root, its archive root and `docs/history`, so ignored files are never read
(Git's `--untracked` searches untracked files as well as tracked ones). Outside
a work tree it walks the directory, skipping `.git` and the same exclusions.
A Spec Root outside the repository root yields no pins. `citations.go` adds
one `SC-SPEC-PATH-PINNED` error per pin for the checked Spec, with the pin's
path and line as its location and the fix "Read a fixture or an exported
constant instead of the Spec's file; a Spec archives and may be deleted."
The detector runs for every Spec the full Check runs, whatever its Task
statuses. It is not added to the staged detector list in `coherence.go`, so a
`--stage prd` or `--stage techspec` run neither runs nor lists it.

At settlement the Daemon already compares refusing findings with those
present when the Task started (ADR-0182), so a pin the Task adds fails the
attempt with the finding's location as Verification Feedback, and a pin that
existed before does not.

### The archive refusal

`runArchiveCommand` resolves the Spec Root, calls `ActiveSpecPathPins`, and on
any pin prints a Preflight failure with exit `2` whose reason is
`Spec "<slug>" cannot archive while another file names its active directory <rel>/: <path>:<line>[, <path>:<line>...]`
before `spec.Archive` runs. In the Delivery Queue the refusal reaches the
archive stage as a failed Archive Command, as any other archive refusal does.

### Task Context through the archive

`detectTaskContextReferences` keeps every path that exists. For a missing path
of the form `<spec-root-relative>/<other-slug>/<rest>`, it checks
`<archive-root-relative>/<other-slug>/<rest>`, where the archive root is the
one the Archive Command resolves for that Spec Root (the history root for the
built-in `docs/specs`, `<spec-root>/_archived` otherwise). If that path
exists the entry is resolved and no finding is raised.

### The archived retry

In `Engine.Retry`, the `state.Archived` branch computes `anchor` and
`accepted` for every blocker, except that a `qa-environment-partial` blocker
still requires `state.QAOverride`, and only that blocker may fall back to the
Run start head when no candidate is recorded. When `accepted`, the head is
appended to `CandidateCommits` and the stage becomes `reviewing`, as today for
the operator archive. The `corrective-spec-required` branch, which runs first,
is unchanged. An unchanged head keeps today's `gating` or `checking` stage.

### The continued item branch

`existingItemBranches` lists `refs/heads/roundfix/deliver-<slug>-*`, keeps
names matching `^roundfix/deliver-<slug>-[0-9a-f]{16}$`, and keeps those for
which `git rev-list --count <base>..<branch>` is above zero, sorted. `base` is
`<remote>/<default>` as `CreateItemBranch` already resolves it, read from the
local remote-tracking ref without fetching. `runDeliverStart` calls it for
each slug after the prerequisite check and before the authorization and
readiness checks; two or more branches refuse with exit `2`. Exactly one
prints `Continuing item branch <branch> for <slug>` on stdout after the queue
is recorded. `CreateItemBranch` calls it again after its fetch and, with one
branch, passes that branch to `RecordDeliveryQueueItemWorktree` instead of a
new name, which reaches the existing branch path of `useAndProvisionItem`.

### Convention C5

`deliveryConventionsVersion` becomes `roundfix/delivery-conventions/v2`.
`deliveryConventions()` gains C5: "A Spec archived through the QA Archive
Override records `qa_override: true`, `qa_override_approval` and
`qa_override_reason` in its archived `_prd.md` front matter; its QA Task keeps
its observed status and its QA Report its observed verdict." In
`conventionRegions`, an anchor under an archive root (index above zero) adds
`C5` when that Spec's `_prd.md` at the review head has front matter with
`qa_override: true` and non-blank approval and reason.

### The parse check and the probe report

`ProbeCommands` runs `sh -n -c <command>` before `Verify`. A non-zero exit
records `Unknown: true` with a `VerificationUnknownError` wrapping
`ErrVerificationMalformed` and the parser's stderr, and the command is not
run. If `sh` cannot be started, the prober runs the command as today. The
Daemon's pre-work refusal for an unknown verdict therefore covers a malformed
command. `specCheckVerificationCommand` maps that cause to the verdict
`malformed`, which counts like `unknown` for exit `1`.
`probeSpecVerifications` runs `git status --porcelain --untracked-files=all`
over each probed Spec's `_tasks.md` and Task files in the repository root and
reports each listed file as `untracked` or `modified`.

### Data Models

`specCheckVerificationReport` gains `Uncommitted []{path, state}` in JSON as
`uncommitted`, always present as an array. No store schema changes.

### API Contracts

1. API Contract 1: `roundfix spec check` reports `SC-SPEC-PATH-PINNED` (error)
   per pin and `SC-VERIFY-TRUNCATED` (error) per truncated bullet of a Task
   that is not completed; exit `1` when either is present.
2. API Contract 2: `roundfix archive <slug>` exits `2` with the reason above
   and changes nothing while a pin exists.
3. API Contract 3: `roundfix deliver retry <slug>` of an archived item whose
   head descends from its newest candidate prints the stage `reviewing`.
4. API Contract 4: `roundfix deliver start` continues one item branch with
   work and prints `Continuing item branch <branch> for <slug>`; with two or
   more it exits `2` naming each and records no queue.
5. API Contract 5: the review prompt and record carry conventions version
   `roundfix/delivery-conventions/v2` and the eligible rule `C5`.
6. API Contract 6: `roundfix spec check --run-verification` prints
   `malformed` with the parser's message for a command `sh -n` rejects, and
   one `Uncommitted Verification source: <path> (<state>)` line per
   uncommitted source after `Verification tree: HEAD`.

### Surface Transcripts

1. Surface Transcript: a refused archive. The fixture repository holds an
   active Spec `0300-example` and `internal/example/example_test.go`, whose
   line 7 reads `docs/specs/0300-example/_techspec.md`.

   ```transcript
   $ roundfix archive 0300-example
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "0300-example" cannot archive while another file names its active directory docs/specs/0300-example/: internal/example/example_test.go:7

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   exit: 2
   ```

2. Surface Transcript: a refused queue start. The fixture repository's
   `origin/main` lacks commits held by two branches
   `roundfix/deliver-0300-example-<16 hex>`.

   ```transcript
   $ roundfix deliver start 0300-example
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "0300-example" has 2 item branches with commits origin/main lacks: roundfix/deliver-0300-example-<hex>, roundfix/deliver-0300-example-<hex>; delete every branch but the one to continue, then run roundfix deliver start again

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   exit: 2
   ```

3. Surface Transcript: a malformed, uncommitted command. The fixture Spec
   `0300-example` is untracked, and its task_01 Verification bullet is the
   truncated command recorded in Spec 0205's delivery.

   ```transcript
   $ roundfix spec check 0300-example --run-verification
   stdout:
   ...
   Verification tree: HEAD
   Uncommitted Verification source: docs/specs/0300-example/_tasks.md (untracked)
   Uncommitted Verification source: docs/specs/0300-example/task_01.md (untracked)
   - task_01: malformed — <command> (<parse error>)
   stderr:
   exit: 1
   ```

A continued item branch has no transcript: its line prints only after the
delivery readiness checks pass and a queue owner starts, which needs an
authenticated `gh`. task_02's tests assert the line through the command seam.

## Vocabulary Contract

No new glossary term. The emitted words are `SC-SPEC-PATH-PINNED`,
`SC-VERIFY-TRUNCATED`, `cannot archive while another file names its active
directory`, `Continuing item branch`, `item branches with commits`,
`delivery-conventions/v2`, `C5`, `malformed` and `Uncommitted Verification
source`. task_01 documents the pin in the `archive` and `spec` guides,
task_02 the retry and start words in `deliver`, task_03 C5 in `review`, and
task_04 the probe words in `spec`, each with the matching Roundfix Skill
reference.

## Coverage Map

- Goal 1 → The pin detector; The archive refusal; Task Context through the archive.
- Goal 2 → The archived retry; The continued item branch.
- Goal 3 → Convention C5.
- Goal 4 → The parse check and the probe report.
- User Story 1 → The pin detector; API Contract 1.
- User Story 2 → Task Context through the archive.
- User Story 3 → The archived retry; API Contract 3.
- User Story 4 → The continued item branch; API Contract 4.
- User Story 5 → Convention C5; API Contract 5.
- User Story 6 → The parse check and the probe report; API Contract 6.
- Core Feature 1 → The pin detector.
- Core Feature 2 → The archive refusal; API Contract 2.
- Core Feature 3 → Task Context through the archive.
- Core Feature 4 → The archived retry.
- Core Feature 5 → The continued item branch.
- Core Feature 6 → Convention C5.
- Core Feature 7 → The parse check and the probe report; API Contract 1.
- Core Feature 8 → The parse check and the probe report; Data Models.
- Core Feature 9 → Build Order 1; Build Order 2; Build Order 3; Build Order 4.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 2.
- Success Metric 5 → Testing Approach 3.
- Success Metric 6 → Testing Approach 4.

## Integration Points

- Git, through the existing runners: `git grep`, `git for-each-ref`,
  `git rev-list` and `git status`, all local and read-only.
- The POSIX shell, through `sh -n -c`, which reads and does not execute.
- The review's configured runtime, only through the existing sealed validator.

## Testing Approach

1. `internal/speccheck/active_spec_paths_test.go` builds temporary Git
   repositories: a Go test naming the active directory is an error, a Markdown
   file, a file under the archive root and an ignored file are not, and a
   Task Context path resolves through the archive only when the archived
   file exists. `internal/daemon/settlement_active_spec_path_test.go` proves
   the finding is refusing through `SpecCheckSettlementChecker`.
   `internal/cli/archive_active_spec_path_test.go` runs the Archive Command
   and asserts the refusal, exit `2` and an unmoved directory.
2. `internal/delivery/archived_head_retry_test.go` drives `Retry` with the
   existing fakes for `gate-failed`, a non-descendant head and a
   corrective-Spec park. `internal/cli/deliver_item_branch_test.go` builds a
   bare remote and a clone, creates item branches with and without commits,
   and asserts the refusal, the recorded branch on the new item and the
   continuation line through the command's seams; no test reaches GitHub.
3. `internal/cli/review_override_convention_test.go` asserts the C5 region
   through `conventionRegions` over a temporary repository with and without
   the override front matter, and the prompt text and version.
4. `internal/daemon/verification_probe_malformed_test.go` asserts the verdict,
   that the command never ran (a marker file stays absent) and the pre-work
   refusal; `internal/speccheck/verification_truncated_test.go` the two
   truncated shapes and a clean bullet; `internal/cli/spec_check_provenance_test.go`
   the `malformed` verdict, the uncommitted lines, the JSON field and their
   absence for a committed Spec.

## Build Order

1. The pin detector, Task Context through the archive, the archive refusal,
   and the `archive` and `spec` guide and skill text.
2. The archived retry, the continued item branch, and the `deliver` guide and
   skill text (depends on: 1).
3. Convention C5 and the `review` guide and skill text (depends on: 2).
4. The parse check, the truncated-span finding, the probe report, and the
   `spec` guide and skill text (depends on: 3).

Each step raises the Roundfix Skill's version and re-records it, so the four
run in series; steps 1 and 4 also share `citations.go` and the `spec` guide.

## Risks & Considerations

- `git grep` over a large repository at every settlement costs one process;
  the exclusions keep the history root out. A directory walk outside Git is
  only for test fixtures.
- A pin in a Go comment is refused too; the comment would dangle after the
  archive, which is the same defect in prose.
- A continued branch may conflict with the default branch at publication; the
  existing conflict park handles it, and merging is out of scope.
- `sh -n` on a shell other than POSIX `sh` is not used: the Daemon already runs
  every command through `sh -c`.

## Decisions

- Refuse a pin rather than rewrite it, resolve another Spec's Task Context
  through the archive, re-review a post-archive correction, and continue one
  item branch with work. See ADR-0223.
- Name uncommitted Verification sources without refusing them, and report a
  command the shell cannot parse as malformed through the shared prober. See
  ADR-0226.
- Keep `constraints.go`, `coherence.go` and `internal/spec/archive.go`
  unchanged, so no Governed Path beyond the Roundfix Skill is needed.
