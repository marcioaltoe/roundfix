---
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
prd: _prd.md
created: 2026-11-16
---

# Each phase of a Spec is recorded and starts in a fresh session — Technical Spec

## Executive Summary

A new package appends JSON lines to day files under Roundfix Home, and three writers use it: a new `roundfix phase start` command that the skills run, the Daemon at the point where it already records a prompt's token usage, and `roundfix review` after each review prompt. The Run Retention Sweep removes old day files. The boundary rule is text: the four authoring skills, the Roundfix skill's `modes` reference and one clause of the Baseline's autonomous-work module. Installation learns one rule from the tree, that a skill whose source sits under `skills/<name>/` at the repository root receives no header, and this repository then declares its authoring models. The primary trade-off: the Supervisor's phase is recorded by a command a skill must remember to run, not inferred by Roundfix, because every command the Supervisor already runs at a phase start is read-only or runs too late; the price is a gap wherever an agent skips the step, which the first criterion of the PRD's Unreachable Acceptance measures.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Phase Record stores identifiers other components made and names its files by date. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the record is a local file that holds no secret. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0407 and ADR-0408 (new) record the decisions, and ADR-0408 amends ADR-0401. ADR-0402, ADR-0405, ADR-0338 and ADR-0255 apply as the PRD's row states, each with its receipt there. Not applicable: ADR-0290, ADR-0296, ADR-0301, ADR-0353, ADR-0383, ADR-0390 and ADR-0386. `## ADR Alternatives` relates the design to each. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answers of 2026-11-14 to 2026-11-16 quoted in the authorization record, on the cycle BF terms. No Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/_authorization.md`; bounded files: `.roundfixrc.yml`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/storage.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Characterized at `2caa1d17`.

- Nothing under Roundfix Home names a phase. Roundfix Home has no resolver of its own: each caller joins `.roundfix` to a home directory it was given (`store.DatabasePath`, `internal/store/store.go:307`; `newJudgeLog`, `internal/judge/log.go:16`), and tests pass a temporary home.
- The Judge Log is the existing pattern for an append-only file there: directory `0700`, file `0600`, `O_CREATE|O_APPEND|O_WRONLY`, one JSON line per call (`internal/judge/log.go:111-118`).
- The Daemon records a prompt's usage in one place, `Engine.recordPromptUsage` (`internal/daemon/agent_session_owner.go:1124`), called from `runAgentSession` (`:214`, scope `session`) and from `runPrepared` (`:666`, scope `task`, `qa` or `review`). It holds the Run, the scope, the runtime, the working directory and the result, whose `ACPSessionID` (`internal/agent/agent.go:222`) the Daemon never reads. It fires once per prompt, repair prompts included.
- `roundfix review` keeps the reported session identifier of each prompt in its lineage (`internal/cli/review.go:933`) and runs with no Run.
- The Run Retention Sweep is `sweepRunDatabase` (`internal/cli/run_retention.go:32`), which holds the cutoff and the home directory and already removes artifact directories.
- A top-level command is three edits in `internal/cli/cli.go` (the usage text, the dispatch switch, `commandUsage`) and one file of its own, as `reopen` shows (`internal/cli/reopen.go`).
- Installation applies the header in `prepareAuthoringInstall` (`skills/authoring.go:256`, the call at `:312`), which knows the repository root; `InspectAuthoring` (`skills/authoring_inspect.go:27`) reports a header that differs from the declared one. Measured in a throwaway clone of this repository with the declaration added: four headers written, the mirror then different, five subagent files generated, the skills lock unchanged.
- The clause `clause.autonomous.loop-07-outlive-the-turn` of `internal/baseline/assets/modules/autonomous-work.json` names "a scheduled wake-up" as a way to keep working; no test pins its text.
- `TestTheAuthoringSkillsDelegateADeclaredPhase` (`skills/authoring_delegation_test.go:9`) requires thirteen phrases exactly once in the `## Declared model` section of each authoring skill, and `skills/roundfix_modes_test.go` pins the order and several sentences of the follow-up steps.

## System Architecture

### Phase Record store

A new package, `internal/phaserecord`, owns the directory, the schema, the two entry kinds, the boundary rule and the removal of old day files. It depends on the standard library only and knows nothing of the Daemon, the CLI or a harness. Every writer receives a home directory and a clock, as the Judge Log does.

### Phase command

`internal/cli/phase.go` adds `roundfix phase start`. It resolves the Git root, reads the session identity through the command environment, calls the store and prints the result and the Boundary Notice. It loads no Project Config.

### Prompt recording

The Daemon gains one optional collaborator in `Dependencies`, a prompt recorder, wired by the CLI with the home directory. `recordPromptUsage` receives the time the prompt began and writes the entry beside the usage row. `roundfix review` writes the same entry after each prompt through the same store.

### Retention

`sweepRunDatabase` asks the store to remove day files older than its cutoff after its candidate loop, and its report gains one count.

### Published-skill rule

`prepareAuthoringInstall` skips the header of a declared authoring skill when `<root>/skills/<skill>/SKILL.md` is a regular file, and reports it as withheld. `InspectAuthoring` applies the mirror rule: for such a skill the declaration is applied when the installed file holds no `model` and no `effort`. The reader of user-level copies stops saying that a copy of such a skill hides a header.

### Skill and guide text

The four authoring skills gain a `## Phase boundary` section and the recording step; `write-techspec` and `write-tasks` gain the fix and refactor rule in `## Declared model`. The Roundfix skill gains a `phase` reference and changes its `modes` reference. The Baseline clause above is extended, which regenerates the guide and its derived files.

## Implementation Design

### Interfaces

The store:

```go
package phaserecord

const Schema = "roundfix/phase-record/v1"

type Recorder struct{ /* home directory, clock */ }

func New(homeDir string, now func() time.Time) *Recorder

// StartPhase appends a phase entry and reports whether the session crossed a boundary.
func (r *Recorder) StartPhase(ctx context.Context, entry PhaseEntry) (Boundary, error)

// RecordPrompt appends a prompt entry.
func (r *Recorder) RecordPrompt(ctx context.Context, entry PromptEntry) error

// RemoveBefore removes, or with dryRun counts, the day files dated before cutoff.
func (r *Recorder) RemoveBefore(ctx context.Context, cutoff time.Time, dryRun bool) (int, error)

func ParsePhase(name string) (Phase, error)
```

The entries:

```go
type PhaseEntry struct {
	Phase      Phase  // follow-up | idea | prd | techspec | tasks | delivery | release
	Spec       string // folder slug; empty only for follow-up and release
	Repository string // Git root the command ran in
	Harness    string // "claude" or empty
	Session    string // harness session identifier or empty
}

type PromptEntry struct {
	Started, Ended time.Time
	RunID, Spec    string // RunID empty for a pre-PR review
	ScopeKind      string // task | qa | review | session | pre-pr-review
	ScopeID        string
	Repository     string
	WorkDir        string
	Harness        string // the runtime identifier
	Session        string // the ACP session identifier the adapter reported
}

type Boundary struct {
	Crossed       bool
	PreviousPhase Phase
	PreviousSpec  string
}
```

The Daemon's collaborator and the session identity:

```go
// internal/daemon
type PromptRecorder interface {
	RecordPrompt(ctx context.Context, entry phaserecord.PromptEntry) error
}

// Dependencies gains: PhaseRecord PromptRecorder (nil records nothing).

// internal/cli: the harness session of the process, from the command environment.
func harnessSession(lookup func(string) string) (harness, session string)
```

`harnessSession` returns `claude` and the value of `CLAUDE_CODE_SESSION_ID` when it is set, and two empty strings otherwise.

### Data Models

One directory, `<home>/.roundfix/phase-record`, mode `0700`. One file per UTC date, `YYYY-MM-DD.jsonl`, mode `0600`, opened with `O_CREATE|O_APPEND|O_WRONLY` for each entry and closed. One JSON object per line, with the keys in this order.

A phase entry: `schema`, `kind` (`phase`), `time`, `phase`, `spec`, `repository`, `harness`, `session`, `boundary`, and, only when `boundary` is true, `previous_phase` and `previous_spec`.

A prompt entry: `schema`, `kind` (`prompt`), `time`, `started`, `ended`, `run_id`, `spec`, `scope_kind`, `scope_id`, `repository`, `work_dir`, `harness`, `session`.

Every time is RFC 3339 in UTC with seconds. An absent value is the empty string, never a missing key, so a reader needs no schema branch. A line is bounded: each string field is cut at 1,024 bytes, so a line stays far below the size at which an appending write could interleave with another.

### API Contracts

1. API Contract: `roundfix phase start <phase> [--spec <slug>]`. Phases: `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery`, `release`. `--spec` is a Spec folder slug matching `^[0-9]{4}(-[a-z0-9]+)+$`; it is required for `idea`, `prd`, `techspec`, `tasks` and `delivery`. The command resolves the Git root of the working directory, reads the session identity, appends one phase entry and prints one line on standard output. Exit `0`: recorded. Exit `2`: an unknown phase, a missing or malformed `--spec`, an extra argument, or no Git repository; nothing is written. Exit `1`: the record could not be written. It reads no Project Config and no User Config.
2. API Contract: the boundary rule. Before appending, `StartPhase` reads the day files of the current UTC date and the six dates before it and takes the latest phase entry whose `session` equals the new entry's non-empty `session`. A boundary is crossed when that entry exists, its phase is not `follow-up`, and its phase or its Spec differs from the new entry's, except when it is `delivery` and the new one is `release` of the same Spec or of none. An entry without a session never crosses and is never matched.
3. API Contract: the Boundary Notice, one line on standard error after the result line: `roundfix: notice: this session already recorded <phase> of <spec>; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.` When the earlier entry names no Spec, ` of <spec>` is left out. The exit code and standard output are unchanged.
4. API Contract: the Daemon's prompt entry. `recordPromptUsage` writes one entry per prompt that returned, with the scope kind and identifier it already records, the Spec of the owning scope or the empty string, the Run's recorded repository root, the request's working directory, the runtime identifier of the request that ran and `result.ACPSessionID`. A failure prints `roundfix: warning: phase record not written for <scope_kind> <scope_id>: <error>` on the Run's progress output and changes nothing else. A nil recorder writes nothing.
5. API Contract: the review's prompt entry. After each prompt of `roundfix review`, the command writes one entry with an empty `run_id`, `scope_kind` `pre-pr-review`, the review session's name as `scope_id`, the single Spec its record names or the empty string when it names none or several, the reviewed worktree as `repository` and `work_dir`, the runtime and the reported session identifier. A failure prints the same warning on standard error and the review's result and exit code are unchanged.
6. API Contract: retention. In a live sweep `sweepRunDatabase` removes the day files whose date is before the cutoff's UTC date; in a dry run it counts them. The Run Retention section of the GC report gains one line, indented as its neighbours, `Phase Record day files removed: <n>` or, in a dry run, `Phase Record day files removable: <n>`, after the line that counts the Runs. A file whose name is not a date is left alone. A failure is reported as the sweep reports its other failures and removes no Run less.
7. API Contract: installation of a published skill. For each declared authoring phase whose skill has a regular file at `<root>/skills/<skill>/SKILL.md`, installation writes the shipped file with no header and prints `Withheld declared header: <skill> (this repository publishes skills/<skill>/SKILL.md)` where it prints `Applied declared header:` for the others. The result of `roundfix baseline update` gains a `withheld` list beside `headers`; no existing field is renamed.
8. API Contract: inspection of a published skill. For such a skill `InspectAuthoring` reports no problem when the installed `SKILL.md` holds neither a top-level `model` nor a top-level `effort`, and otherwise the installable problem `<skill> holds a header, and this repository publishes skills/<skill>/SKILL.md`. `roundfix skills check`, the Doctor Command and the Update Check print it through their existing lines. A user-level copy of such a skill is reported by its version only: `hidesDeclaredHeader` is false for it.
9. API Contract: the skill text. Each authoring skill's `## Phase boundary` section says: record the phase with `roundfix phase start <phase> --spec <slug>` once the slug is known, and say so and continue when the installed Roundfix does not know the command; when the phase ends, an Author Subagent returns its report, and a session writes the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, tells the user what stays armed, and starts no further phase. In `write-techspec` and `write-tasks`, the bullet of `## Declared model` for a skill the agent activated gains: on the fix or refactor route, run the phase in the session unless the user's request or one of the user's answers asks for the Author Subagent, and say in the report where it ran.
10. API Contract: the Roundfix skill. The `modes` reference records `follow-up` in step 1, `delivery` when step 7 starts to follow a Spec and `release` before a release; its step 8 and its `--unattended` section state the boundary rule, that a scheduled wake-up keeps a phase going and never crosses a boundary, and that a fix or a refactor is authored in the session with no question about it. A new `phase` reference documents the command and the record, and the reference index gains its row.
11. API Contract: the Baseline clause. `clause.autonomous.loop-07-outlive-the-turn` keeps its identity and its enforcement and gains the boundary rule after its existing text: at a phase boundary of a Spec, after each authoring phase and after a Spec is delivered and released, write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed, tell the user what stays armed, and start no further phase in that session; a scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and is never the way across a boundary.
12. API Contract: this repository's declaration. `.roundfixrc.yml` gains `skills.authoring.claude` with `idea`, `prd` and `techspec` at `model: fable`, `reasoning_effort: xhigh`; `tasks` at `model: fable`, `reasoning_effort: high`; and `research` at `model: sonnet`, `reasoning_effort: medium`, `tools: [WebSearch, WebFetch, mcp__exa__web_search_exa, mcp__exa__web_fetch_exa]`. The five files under `.claude/agents/` are the output of `roundfix skills install --target project` run with a binary built after `make skills-sync`.

### Surface Transcripts

Fixture A is a temporary home and, as the working directory, a temporary Git repository, with `CLAUDE_CODE_SESSION_ID` set to `11111111-2222-4333-8444-555555555555`. Fixture B is fixture A after Surface Transcript 2. Fixture C is fixture A with `CLAUDE_CODE_SESSION_ID` unset. Fixture D is a temporary home and a temporary Git repository with no `.claude` directory, whose `.roundfixrc.yml` declares `prd` and `techspec` with `model: opus` and `reasoning_effort: xhigh`, and which holds `skills/write-prd/SKILL.md`, a copy of the shipped file. Fixture E is fixture D after Surface Transcript 6, with a line `model: opus` added to the frontmatter of `.agents/skills/write-prd/SKILL.md`. `<names>` is the bundle's names joined by `, `.

1. Surface Transcript: the usage of the new command.

   ```transcript
   $ roundfix phase --help
   stdout:
   Usage:
     roundfix phase start <phase> [--spec <slug>]

   Commands:
     start  Record in the Phase Record that this session starts a phase

   Phases:
     follow-up, idea, prd, techspec, tasks, delivery, release

   Options:
     --spec  Spec folder slug the phase works on; required for every phase
             but follow-up and release
   stderr:
   exit: 0
   ```

2. Surface Transcript: a phase is recorded (fixture A).

   ```transcript
   $ roundfix phase start prd --spec 0412-example
   stdout:
   Recorded prd of 0412-example for claude session 11111111-2222-4333-8444-555555555555
   stderr:
   exit: 0
   ```

3. Surface Transcript: the same session starts the next phase (fixture B).

   ```transcript
   $ roundfix phase start techspec --spec 0412-example
   stdout:
   Recorded techspec of 0412-example for claude session 11111111-2222-4333-8444-555555555555
   stderr:
   roundfix: notice: this session already recorded prd of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
   exit: 0
   ```

4. Surface Transcript: no session identity (fixture C).

   ```transcript
   $ roundfix phase start follow-up
   stdout:
   Recorded follow-up without a harness session: the environment names none
   stderr:
   exit: 0
   ```

5. Surface Transcript: an unknown phase is refused and nothing is written (fixture A).

   ```transcript
   $ roundfix phase start plan --spec 0412-example
   stdout:
   stderr:
   roundfix: phase start failed: unknown phase "plan"; use follow-up, idea, prd, techspec, tasks, delivery or release
   Run 'roundfix phase --help' for usage.
   exit: 2
   ```

6. Surface Transcript: installation withholds the header of a published skill (fixture D).

   ```transcript
   $ roundfix skills install --target project
   stdout:
   Installed Roundfix skill for project: <path> (<count> file(s))
   Withheld declared header: write-prd (this repository publishes skills/write-prd/SKILL.md)
   Applied declared header: write-techspec (model opus, effort xhigh)
   Generated subagent: .claude/agents/roundfix-prd-author.md
   Generated subagent: .claude/agents/roundfix-techspec-author.md
   stderr:
   exit: 0
   ```

7. Surface Transcript: a header in a published skill is reported (fixture E).

   ```transcript
   $ roundfix skills check
   stdout:
   Roundfix skill check passed: <names>
   stderr:
   roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.claude):
     write-prd holds a header, and this repository publishes skills/write-prd/SKILL.md
     next: roundfix skills install --target project
   exit: 0
   ```

8. Surface Transcript: a dry run counts the day files past the cutoff (a temporary Git repository and a temporary home with no Run Database, whose `phase-record` directory holds two day files dated more than 30 days ago and one of the current date).

   ```transcript
   $ roundfix gc --dry-run
   stdout:
   ...
     Phase Record day files removable: 2
   ...
   stderr:
   exit: 0
   ```

### Invariants

```text
1. An entry holds only the keys of its kind, in their order, and no value read from a prompt, a transcript or the environment other than the session identifier.
2. A line is never rewritten or removed except with its whole day file by retention.
3. A failed write of a prompt entry changes no Task, Run, review or exit code.
4. The Boundary Notice changes no exit code and no standard output.
5. No read-only command writes the record: follow-up, its status view, runs, events and attach are untouched.
6. Without a declaration, and in a repository that publishes no skill, installation and the three checks print the bytes they print at the delivery base.
7. No shipped SKILL.md, canonical or mirrored, holds a model or an effort.
8. The Baseline clause keeps its identifier and its enforcement.
```

### Declared breaks

- The root usage text and the dispatch gain `phase`; the tests that pin the root usage and `TestEveryCommandPathPrintsItsHelp` move to the new text.
- `TestTheAuthoringSkillsDelegateADeclaredPhase` keeps its thirteen phrases once each and gains the fix and refactor phrases for `write-techspec` and `write-tasks` and their absence in `write-idea` and `write-prd`.
- The follow-up tests of `skills/roundfix_modes_test.go` keep their order and sentences; steps 1, 7 and 8 and the `--unattended` section gain sentences that new tests pin.
- `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json` gain the `phase` command and its guide; the surfaces of `skills install`, `skills check`, `doctor`, `update check` and `gc` are re-recorded.
- The version of the `autonomous-work` module rises, and with it the module record, the catalog records, the profile's golden digest, the formatter golden of the guide, the generated guide and the Setup Manifest.
- The versions of the five owned skills rise with their record.
- `GC`'s Run Retention section gains one line in both modes; the tests that pin that section move to it.

## Coverage Map

- Goal 1, Stories 1, 2 and 3, Core Features 1, 2, 4 and 5, Success Metrics 1, 3 and 4 → Phase Record store, Phase command, Prompt recording (API Contracts 1, 4 and 5; Surface Transcripts 1, 2, 4 and 5).
- Goal 2, Stories 4 and 5, Core Features 8 and 9, Success Metric 8 → Skill and guide text (API Contracts 9, 10 and 11).
- Goal 3, Story 6, Core Feature 3, Success Metric 2 → Phase command, which prints the Boundary Notice on standard error and keeps its standard output and exit code (API Contracts 2 and 3; Surface Transcript 3; Invariant 4).
- Core Feature 6, Success Metric 4 → Prompt recording (API Contracts 4 and 5; Invariant 3).
- Core Feature 7, Success Metric 5 → Retention (API Contract 6; Surface Transcript 8).
- Goal 4, Stories 7 and 8, Core Features 10 and 11, Success Metrics 6 and 7 → Published-skill rule (API Contracts 7, 8 and 12; Surface Transcripts 6 and 7; Invariants 6 and 7).
- Goal 5, Story 9, Core Feature 12, Success Metric 8 → Skill and guide text (API Contracts 9 and 10).
- Core Feature 13, Success Metric 8 → the documentation each Build Order step names, the glossary terms and the notice in ADR-0401.

Deviations from the PRD: None.

## Absence Evidence

1. The Phase Record store — command: `rg -n -i 'phase-record|phaserecord' --glob '*.go' .`; result: no match, and `internal/phaserecord` does not exist.
2. The `phase` command — command: `rg -n 'case "phase"|roundfix phase' internal/cli/cli.go docs/user-guide .agents/skills/roundfix`; result: no match.
3. A reader of the harness session — command: `rg -n 'CLAUDE_CODE_SESSION_ID' --glob '*.go' .`; result: no match.
4. The published-skill rule — command: `rg -n -i 'withheld|publish' skills/authoring.go skills/authoring_inspect.go`; result: no match; `prepareAuthoringInstall` applies the header to every declared phase.
5. The boundary text — command: `rg -n -i 'phase boundary|boundary notice' .agents/skills internal/baseline/assets/modules/autonomous-work.json`; result: no match.
6. The fix and refactor rule in the delegation — command: `rg -n -i 'fix or a refactor|refactor route|fix route' .agents/skills/write-techspec/SKILL.md .agents/skills/write-tasks/SKILL.md`; result: no match.
7. This repository's declaration and generated subagents — command: `rg -n '^skills:' .roundfixrc.yml; ls .claude/agents`; result: no match, and the directory does not exist.

## ADR Alternatives

- ADR-0407: the design takes its decision whole. Of its rejected options, the record written only by a Run, the Run Database with a read command, the Artifact Directory, recording from read-only commands and a refusal at a boundary are all absent from the design.
- ADR-0408: the design takes its decision whole: the rule is read from the tree, no configuration key or marker file exists, the mirror's synchronization is untouched, and the Skill Modes ask no question about a fix.
- ADR-0401: amended by ADR-0408. The design keeps its rejected options rejected: no hand-written header is preserved, no model name is shipped, the section stays in Project Config, and a check never fails for an unapplied declaration. It adds no entry per skill.
- ADR-0402: unchanged. This repository lists only tools of the accepted set, and the design gives the Research Subagent no shell and no file-writing tool.
- ADR-0405: unchanged in its commands. The design narrows where its "hides the declared header" words appear, as ADR-0408 states, and takes none of its rejected options.
- ADR-0338: unchanged. The boundary rule relies on the mode claiming the handoff first and changes no step of it.
- ADR-0255: unchanged in which Runs it removes. The sweep also removes day files of the record, under the same cutoff and the same budget.
- ADR-0290: not applicable — it decides when Item Binaries and Run Worktrees leave; the record is neither.
- ADR-0296: not applicable — the storage report accounts for the pages of the Run Database, and the record is files outside it.
- ADR-0301: not applicable — it classifies Artifact Roots, and the record is under Roundfix Home, outside every Artifact Root.
- ADR-0353: not applicable — the follow-up list and its Waiting Items are unchanged.
- ADR-0383: not applicable — no rule about a message or its reply changes.
- ADR-0390: not applicable — the precedence of a blocking message is unchanged.
- ADR-0386: not applicable — its comparison by version and its commands are untouched; ADR-0408 changes only when ADR-0405's words about a hidden header appear.

## Integration Points

- Claude Code: the only contact is the environment variable `CLAUDE_CODE_SESSION_ID`, read by one function. No transcript, setting or hook is read or written.
- ACP adapters: the session identifier comes from the result the runner already returns. Measured for the Codex adapter: the 23 identifiers this repository recorded for pre-PR reviews equal the session identifiers of Codex transcripts. Not measured for the Claude adapter.
- The operator's usage script: it reads the day files. It is outside this repository and is not changed here; the user guide is the contract.
- ai-memory: the Session Handoff is written by the skill through the tools the session exposes, as today. No Go code calls it.

## Testing Approach

- The store is tested at its public functions with a temporary home and a fixed clock: the bytes of each entry kind, the mode of the directory and the file, the boundary rule case by case (each exception, a different session, an entry without a session, a previous entry six days old and one eight days old), the cut of an oversized field, the removal by date and its dry run. A test fails when an entry gains a key that is not declared.
- The command is tested through the CLI entry point with a temporary home, a temporary Git repository and a session supplied through the command environment's lookup, never the process environment: each Surface Transcript 1 to 5 in one named test, each refusal writing nothing, and an unwritable record directory exiting `1`.
- The Daemon is tested with the existing fake runner and a recording `PromptRecorder`: one entry per prompt that returned, for a Task, the QA gate and a review Batch, with the reported session identifier; a recorder that fails leaves the Run's outcome and events unchanged and prints the warning once per prompt; a nil recorder changes nothing.
- `roundfix review` is tested with its fake runner and a temporary home in the same two ways.
- Retention is tested through `sweepRunDatabase` with a temporary home: live, dry run, and a file that is not a date.
- Installation and inspection are tested in temporary repositories: fixtures D and E, a repository that publishes every authoring skill, and the unchanged output of a repository that publishes none. The user-level copy reader is tested for a published skill.
- Skill text is tested by contract tests that require each new sentence once, in the canonical tree and in the mirror, and the Baseline clause by a contract test that requires each kept and each new phrase and fails when one is removed.
- This repository's adoption is tested by running the built product in the checkout: `skills check` prints no Authoring Model Declaration warning, the five generated files carry their marker, and the mirror check and `TestNoShippedSkillDeclaresAModelOrEffort` pass.
- No test starts a harness, reads a transcript, reaches the network or touches the real home. Non-test code builds for Windows with `GOOS=windows go build ./...`.

## Build Order

1. The Phase Record store and the `phase` command: `internal/phaserecord`, `internal/cli/phase.go`, the three edits of `internal/cli/cli.go`, `docs/user-guide/commands/phase.md`, the command index, the Phase Record section of `docs/user-guide/usage.md`, the Roundfix skill's `phase` reference with its index row and mirror, the skill coverage map and the behavior surfaces, and the terms **Phase Record**, **Supervisor Phase**, **Phase Boundary** and **Boundary Notice** in `CONTEXT.md`. Covers API Contracts 1 to 3.
2. Prompt entries and retention (depends on: 1): the Daemon's collaborator and its wiring, `recordPromptUsage`, `roundfix review`, `sweepRunDatabase`, the GC and Run Database guides, and the Roundfix skill's `storage` reference with its mirror. Covers API Contracts 4 to 6.
3. The published-skill rule and this repository's declaration (no dependency on 1 or 2): `skills/authoring.go`, `skills/authoring_inspect.go`, the user-level copy reader, the CLI's reports, the configuration, skills, doctor and update guides, the Roundfix skill's `setup` reference with its mirror, the changed **Declared Skill Header** in `CONTEXT.md`, the "Amended by ADR-0408" notice in ADR-0401, then `.roundfixrc.yml` and the five generated files. Covers API Contracts 7, 8 and 12.
4. The boundary rule and the fix route in the skills and the guide (depends on: 1): the four authoring skills, the Roundfix skill's `modes` reference, the Baseline clause with its derived files, and the changed **Author Subagent** and **Session Handoff** in `CONTEXT.md`. Covers API Contracts 9 to 11.
5. The QA gate (depends on: 1, 2, 3, 4).

Steps 1, 2, 3 and 4 each raise the version of the Roundfix skill and write `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json` and `CONTEXT.md`; steps 3 and 4 both touch the authoring skills' installed and mirrored trees. They share those files, so they run in series, in this order. Each is one Task that pays for its own Verification: a package with a command, a Daemon change, an installer rule with its adoption, and a body of pinned text with a Baseline regeneration. No other Spec is active, so `requires` stays empty.

## Risks & Considerations

- The installed binary and this repository. After the merge and until the installed binary is the release that carries this Spec, `roundfix skills check`, Doctor and the Update Check run with the older binary report four headers as not the declared ones and name the installation command, and that installation would write the headers. The release and `roundfix upgrade` come before any installation here; the mirror check and the shipped-skill test refuse the tree if one is run early.
- A stale bundle. `roundfix skills install --target project` writes every shipped file over the installed one, and here the installed skills are the source. Step 3 runs it only with a binary built after `make skills-sync`, and its Verification holds the mirror equal to the canonical tree.
- The generated subagents follow the binary's template. When a later release changes it, `skills check` reports the five files as different until the installation is run again here, as in any adopter.
- A skipped recording step leaves a phase unrecorded. The record cannot detect it; the first Unreachable Acceptance criterion does.
- The notice on a small fix. A fix authored and delivered in one session prints it at each phase; it is one line and names that nothing was refused.
- A subagent's entry carries its parent's session. Two Author Subagents of one session that record different phases print the notice in the second; that is the session crossing a boundary, as decided.
- Concurrent appends come from several processes. Each entry is one write of one bounded line to a file opened in append mode, the Judge Log's pattern.
- The record grows by a few hundred bytes per entry. Retention removes it with the Runs; a home whose retention never runs keeps it.
- Windows: the store uses no build constraint and no permission call beyond the modes passed at creation.

## Glossary

- adds: **Phase Record**
- adds: **Supervisor Phase**
- adds: **Phase Boundary**
- adds: **Boundary Notice**
- changes: **Declared Skill Header**
- changes: **Author Subagent**
- changes: **Session Handoff**

## Decisions

- The record is files under Roundfix Home and the phase is recorded by a command. See ADR-0407.
- A published skill receives no header by a rule read from the tree, and a fix or a refactor is authored in the session. See ADR-0408.
- One store package with the boundary rule inside it, so that the command and its tests share one definition of a boundary.
- The prompt entry is written where the usage row is written, so that every scope the Run already distinguishes is covered by one call site and a new scope is covered when it gains usage.
- The command takes no `--format`: its reader is the agent that ran it, one line is the whole result, and the record is the machine-readable surface.
- The command is `phase start` and not `phase`: a later subcommand can be added without changing this one.
- The boundary looks back seven days: the longest Supervisor session measured lasted 30.6 hours.
- The Baseline clause is extended, not joined by a new one: a new clause in a rule the Standard TypeScript profile requires needs a Source Baseline row, and the extended text belongs with the rule about what outlives a turn.
- This repository's Research Subagent lists the two web tools and the two Exa tools; a sweep that needs a command, the Secondbrain query among them, stays with the author.
