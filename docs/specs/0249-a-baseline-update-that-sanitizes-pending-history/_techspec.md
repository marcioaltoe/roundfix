---
spec: 0249-a-baseline-update-that-sanitizes-pending-history
prd: _prd.md
created: 2026-10-08
---

# A Baseline update that sanitizes pending history — Technical Spec

## Executive Summary

`roundfix baseline update` gains a history section. It plans the Pending
History with the History Sanitize Command's own per-unit planning, binds that
section into the update's Plan Digest, and on confirmation creates the History
Full Tag when it is absent and converts every unit it can. Refused Units are
listed on one line each and never change the update's state or exit code.
`--no-history` leaves the section out. `roundfix upgrade` appends one notice
line that names the Pending History. Under Lenient Legacy Reading a legacy
list-of-maps `unproven` becomes one stable line per map.

The trade-off accepted: the update converts all convertible units in one
change, which for a large history is a large diff of removals and new records.
An operator who wants reviewed batches opts out and keeps ADR-0248's
procedure. The design adds no package and no new seam. It extends
`internal/cli` (`history.go`, `baseline_update.go`, `upgrade.go`) and one
reader in `internal/spec`.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The JSON
  result gains the lowercase `history` object and its lowercase keys, the
  command gains the lowercase flag `--no-history`, and the tag keeps the
  existing name `history-full`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The update reads files and runs local Git, and it creates a
  local annotated tag that it never pushes. The upgrade notice reads the
  filesystem only. Tests use temporary repositories built with Git and a
  temporary home. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0254 (this Spec) governs the history
  section, the batch size, Refused Units, the tag, the order of writes, the
  opt-out and the notice. The tag rule of ADR-0248: "holds every path the batch
  removes or rewrites"; each selected unit is checked against it. ADR-0251:
  "Malformed and duplicate rows still refuse", and the list-of-maps reading
  adds nothing else to what a Legacy Archive Folder tolerates. ADR-0100: "every byte outside a managed marker is
  identical before and after"; that proof keeps covering the instruction
  carriers, and the history section changes only the History Root. ADR-0184:
  "A TechSpec states a command surface as a transcript", answered in Surface
  Transcripts. ADR-0177 cites ADR-0173 and
  does not apply: this Spec does not change the Relocation Citation scan. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits the Roundfix Skill. Express
  maintainer authorization: "considere autorizado a ajustar todas as skills se
  necessário", the standing "Concedo", and the 2026-10-08 grant of the governed
  paths this Spec declares. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0249-a-baseline-update-that-sanitizes-pending-history/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured on 2026-10-08 at `c73a92e0` with the worktree's `bin/roundfix` 0.59.0:

- `roundfix history sanitize` in this repository prints
  `history sanitize plan: 0 unit(s) pending; nothing pending`. Its tag exists,
  so the update's history section here reports `current`.
- `BuildArchiveRecord` decodes the PRD front matter `unproven` into
  `[]string`. A list of maps fails with `yaml: unmarshal errors:` followed by
  `line <n>: cannot unmarshal !!map into string` on the next line, and
  `printHistoryRefused` prints that error verbatim. The refusal line breaks in
  two (Fluxus, 2026-10-07).
- `roundfix baseline update` never inspects the working tree's cleanliness.
  `ApplyPlan` binds each preimage and refuses a stale one with exit 3. It does
  not refuse a dirty tree, so this Spec adds no tree-wide refusal to it.
- The update's text output prints `History moves` and `Verified history moves`
  for History Relocations and nothing about the History Sanitize Command.
  `roundfix upgrade` prints only the recommendation notice.

## System Architecture

```text
roundfix baseline update [--no-history]
  baseline.BuildPlan                         -> Baseline Plan (unchanged)
  planBaselineUpdateHistory (new, cli)       -> history section
    historyInventory + planHistoryUnit       (shared with history sanitize)
    per unit: uncommitted-change check, tag coverage
  Plan Digest = Baseline Plan Digest | combined digest when a unit is selected
  apply: ApplyPlan -> create tag if absent -> convert units -> skills stage
roundfix history sanitize                    (same planHistoryUnit; one-line reasons)
roundfix upgrade
  printUpgradeRecommendationNotice -> printHistoryNotice (new, filesystem only)
internal/spec BuildArchiveRecord (Legacy) -> list-of-maps unproven -> lines
```

## Implementation Design

### Interfaces

```go
// internal/cli/history.go: one inventoried unit planned as history sanitize plans it.
// A returned error is environmental (Git delivery history); a unit refusal is u.refusal.
func planHistoryUnit(ctx context.Context, repo, specRel, archive, archiveRel, revision string, u historyUnit) (historyUnit, error)

// historyRefusalLine renders a Refused Unit reason on one line.
func historyRefusalLine(err error) string

// internal/cli/baseline_update_history.go (new)
type baselineUpdateHistoryPlan struct {
	Report   baselineUpdateHistory // the JSON "history" object
	selected []historyUnit
}
func planBaselineUpdateHistory(ctx context.Context, repo string, environment commandEnvironment) baselineUpdateHistoryPlan
func (p baselineUpdateHistoryPlan) digest(baselinePlanDigest string) (string, error)
func applyBaselineUpdateHistory(ctx context.Context, repo string, p baselineUpdateHistoryPlan) (baselineUpdateHistory, error)

// internal/cli/upgrade.go
func printHistoryNotice(environment commandEnvironment, stderr io.Writer)
```

`historyInventory` keeps its signature and behavior. `runHistoryCommand` calls
`planHistoryUnit` in its loop and keeps every other step.

### Data Models

The update result `roundfix/baseline-update-result/v1` gains one always-present
field:

```go
type baselineUpdateHistory struct {
	Status   string                        `json:"status"` // skipped|current|pending|applied|blocked
	Message  string                        `json:"message,omitempty"`
	Revision string                        `json:"revision,omitempty"` // HEAD, 40-hex
	Tag      *baselineUpdateHistoryTag     `json:"tag,omitempty"`
	Units    []baselineUpdateHistoryUnit   `json:"units"`
	Refused  []baselineUpdateHistoryRefusal `json:"refused"`
	Applied  *baselineUpdateHistoryApplied `json:"applied,omitempty"`
}
type baselineUpdateHistoryTag struct{ Name, Action, Commit, PushCommand string } // action create|created|present
type baselineUpdateHistoryUnit struct{ Unit, Action, Record, Disposition string; Files int; Bytes, BytesAfter int64 }
type baselineUpdateHistoryRefusal struct{ Unit, Reason string }
type baselineUpdateHistoryApplied struct{ Units, Records, Reduced, Removed int; Bytes int64 }
```

The JSON keys are `name`, `action`, `commit`, `pushCommand`, `unit`, `record`,
`disposition`, `files`, `bytes`, `bytesAfter` (reductions only), `reason`,
`units`, `records`, `reduced` and `removed`. The Archive Record schema is
unchanged: `unproven` stays a list of strings.

### API Contracts

1. API Contract: the history section of the plan. Unless `--no-history` is
   given, every invocation of `roundfix baseline update` plans the history of
   the `--repo` Git root. It loads that root's configuration with the command
   environment's home, resolves the Spec Root as `roundfix history sanitize`
   does, inventories with `historyInventory` and plans each unit with
   `planHistoryUnit` at `HEAD`. Each unit is selected unless one of these makes
   it a Refused Unit, with the reason in this order:
   - its conversion refusal;
   - `has uncommitted changes under <path>; commit or restore them and rerun`,
     when `git status --porcelain=v1 --untracked-files=all` over its paths is
     not empty;
   - `history-full does not hold <path>`, when the tag exists and its tree
     lacks a path the unit removes or rewrites.

   No unit counts against a batch: all selectable units are selected. Status
   is `skipped` with `--no-history`, `pending` when a unit is selected and
   nothing was applied, and `current` when none is selected. Status is
   `blocked`, with `message` holding a one-line reason and no unit selected,
   when any of these fails: the configuration load, an external Spec Root, the
   inventory, `HEAD` resolution, a Git failure reading delivery history or
   status, or the tag check (`history-full` exists but is not annotated or is
   not an ancestor of `HEAD`). The tag object is `create` with `commit` =
   `HEAD` when the tag is absent and a unit is selected, and `present` with the
   tag's commit when it exists. `pushCommand` is
   `git push origin history-full` for `create` and `created`.
2. API Contract: the Plan Digest. When the history section selects no unit,
   the update's `planDigest` is the Baseline Plan Digest unchanged. When it
   selects one, `planDigest` is `sha256:` and the hex SHA-256 of the domain
   string `roundfix/baseline-update-history/v1`, a NUL, the Baseline Plan
   Digest, a NUL, and the canonical JSON of `revision`, the tag `action` and
   `commit`, and per selected unit, in plan order:
   - its name and action;
   - the sorted repository-relative paths it removes or rewrites;
   - for a folder, its record path and the SHA-256 of the rendered record;
   - for a kind, its `bytesAfter`.

   `--yes` approves that digest. `--confirm-plan <digest>` must equal it, or
   the update reports `action_required` in the `approval` category, exits 3,
   prints the current digest and writes nothing. Refused Units and a `blocked`
   message are reported, and the digest does not bind them, because nothing
   acts on them.
3. API Contract: state and messages. A selected unit makes a preview
   `plan_ready` (exit 3) even when the Baseline Plan has no change. Its
   message is then
   `guidance matches the current Baseline catalog; <n> history unit(s) pending sanitize`,
   checked before the existing drifted-skill and outdated-skill messages. The
   next action is the existing `--confirm-plan`/`--yes` sentence. `skipped`,
   `current`, `blocked` and Refused Units never change the state, the category
   or the exit code. A repository with no Pending History reports exactly
   what it reports today.
4. API Contract: apply. After the confirmation matches, the update runs, in
   order:
   1. `baseline.ApplyPlan` with the Baseline Plan Digest. Any refusal there
      writes nothing, the tag included.
   2. When the tag action is `create`,
      `git tag -a history-full -m "docs/history before roundfix baseline update sanitized it (ADR-0254)" <revision>`,
      after which the object type must be `tag`. The action becomes `created`.
   3. Each selected unit in plan order, through `spec.ApplyLegacyConversion`
      or `spec.ApplyHistoryKind`.
   4. The existing skills stage.

   On success the history status is `applied`, and `applied` counts units,
   records, reduced files, removed files and removed bytes. A failure in step 2
   or 3 exits 1 with state `failed`, category `history` and a message naming
   the unit, and keeps the Baseline result. The next action is
   `restore the History Root with git restore --staged --worktree -- <history-root> and git clean -fd -- <history-root>, then rerun roundfix baseline update`.
   The update never commits, never pushes and never moves an existing tag.
5. API Contract: text output. When the status is `pending`, `applied` or
   `blocked`, or a Refused Unit exists, the text output prints, directly
   before the `Plan Digest:` line (or before `Next action:` when no digest is
   printed):

   ```text
   History: <status>
   History result: <message>                                   (blocked only)
   History units: <n> (<f> file(s), <b> bytes leave docs/history)   (pending)
   - folder <folder>: removes <k> file(s) (<b> bytes) and writes <record> (<disposition>)
   - <kind>: reduces <k> file(s) from <b> to <a> bytes
   - <kind>: removes <k> file(s) (<b> bytes)
   History applied: <n> unit(s): wrote <r> Archive Record(s), reduced <f> file(s), removed <d> file(s) (<b> bytes) kept in Git at <12-hex> and tag history-full   (applied)
   History refused: <r>
   - refused <unit>: <reason>
   History tag: creates annotated history-full at <12-hex>; push it with git push origin history-full
   History tag: created annotated history-full at <12-hex>; push it with git push origin history-full
   History tag: history-full at <12-hex> holds every planned path
   ```

   Each line appears only when its value exists. `skipped`, and `current`
   without Refused Units, print nothing, so today's text output is unchanged
   for a repository with no Pending History. The JSON always carries
   `history`.
6. API Contract: `--no-history`. It is a boolean flag of `baseline update`
   only. It skips the configuration load, the inventory and every history Git
   call, and reports `history.status: skipped`. The usage text lists
   `[--no-history]` and says that the update otherwise plans the Pending
   History, converts it on approval and creates `history-full` when absent.
7. API Contract: the upgrade notice. After the recommendation notice, or its
   `recommendations not checked` line, on every release outcome that prints
   one, `roundfix upgrade` loads the working directory's configuration. A
   load failure prints nothing more. Outside a Git repository it prints
   `roundfix: history: outside a repository; run roundfix baseline update in each adopted repository to plan its pending history`.
   Inside one it resolves the Spec Root and runs `historyInventory`, with no
   Git subprocess. With pending units it prints
   `roundfix: history: <u> unit(s) pending sanitize (<parts>); run roundfix baseline update to plan them`.
   `<parts>` is `<f> Legacy Archive Folder(s)` when `<f>` is positive,
   followed by each pending kind name, joined by `, `. It prints nothing when
   none is pending or the Spec Root is external, and it prints
   `roundfix: history not checked: <one-line reason>` when the inventory
   fails. The notice never changes stdout or the exit code, and help, usage
   errors and failed upgrades still print no notice.
8. API Contract: legacy `unproven` maps and one-line reasons. With
   `ArchiveRecordInput.Legacy`, `BuildArchiveRecord` reads `unproven` as a
   sequence of nodes:
   - A scalar item keeps its text.
   - A mapping item becomes one line: `row <row>: ` when `row` exists, then
     `claim`'s value, then the remaining keys in lexical order as
     `<key>: <value>` joined by `; `. Those are wrapped as ` (<...>)` after a
     claim, or stand alone without one.
   - Every value is whitespace-normalized, and a sequence of scalars is joined
     by `, `.
   - An empty mapping, or a value that is a mapping or holds one, refuses with
     `legacy unproven item <i> cannot be read as text`.

   Without `Legacy`, `unproven` stays `[]string` and fails exactly as today.
   `historyRefusalLine` joins the reason's whitespace-separated fields with
   single spaces. It renders every `refused <unit>: <reason>` line of the
   History Sanitize Command and of the update.

### Invariants

```text
1. With no Pending History, the update's stdout, state, exit code and planDigest equal today's for every existing fixture.
2. A preview, with or without Pending History, changes no file, no index entry and no ref.
3. Only units whose every removed or rewritten path is tracked and equal to HEAD are selected.
4. A selected unit's paths are held by the tag that exists, or by the tag the apply creates at the planned revision.
5. An existing history-full is never moved, replaced or deleted; no command of this Spec pushes.
6. A refusal in ApplyPlan writes nothing, the tag included.
7. A Refused Unit's bytes are unchanged after apply, and a Refused Unit never changes state, category or exit code.
8. A second update after apply, without committing, reports history current and the Baseline current.
9. --no-history runs no history Git command and changes no history byte.
10. The upgrade notice never changes stdout or the exit code and runs no Git subprocess.
11. A list-of-maps unproven is accepted only under Lenient Legacy Reading; an active Spec reads it as today.
12. Every refused line, in both commands, is one line.
```

### Surface Transcripts

Each runs from the built tree with `go run -buildvcs=false ./cmd/roundfix`,
against a temporary adopted repository whose Baseline is current. It holds
the committed Legacy Archive Folders `0001-maps-unproven` (a list-of-maps
`unproven`), `0002-plain` and `0003-broken` (an unparsable `_prd.md`), and one
unreduced retired Finding. It has no `history-full` tag and runs under a
temporary home.

1. Surface Transcript: the preview lists the history and writes nothing.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --format text
   stdout:
   Baseline update: plan ready
   Category: approval
   Result: guidance matches the current Baseline catalog; 3 history unit(s) pending sanitize
   ...
   History: pending
   History units: 3 (<files> file(s), <bytes> bytes leave docs/history)
   - folder docs/history/specs/0001-maps-unproven: removes <n> file(s) (<b> bytes) and writes docs/history/specs/0001-maps-unproven.md (no-qa)
   - folder docs/history/specs/0002-plain: removes <n> file(s) (<b> bytes) and writes docs/history/specs/0002-plain.md (no-qa)
   - findings: reduces 1 file(s) from <before> to <after> bytes
   History refused: 1
   - refused docs/history/specs/0003-broken: <reason>
   History tag: creates annotated history-full at <hex>; push it with git push origin history-full
   Plan Digest: <digest>
   Next action: <next>
   stderr:
   exit status 3
   exit: 1
   ```

2. Surface Transcript: the confirmed update converts, tags and stays
   uncommitted.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --yes --format text
   stdout:
   Baseline update: verified
   ...
   History: applied
   History applied: 3 unit(s): wrote 2 Archive Record(s), reduced 1 file(s), removed <n> file(s) (<b> bytes) kept in Git at <hex> and tag history-full
   History refused: 1
   - refused docs/history/specs/0003-broken: <reason>
   History tag: created annotated history-full at <hex>; push it with git push origin history-full
   ...
   stderr:
   exit: 0
   ```

3. Surface Transcript: the second update is current and still names the
   Refused Unit.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --format text
   stdout:
   Baseline update: current
   Result: the repository already matches the current Baseline catalog
   ...
   History: current
   History refused: 1
   - refused docs/history/specs/0003-broken: <reason>
   ...
   stderr:
   exit: 0
   ```

4. Surface Transcript: the opt-out prints no history and plans none.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --no-history --format text
   stdout:
   Baseline update: current
   Result: the repository already matches the current Baseline catalog
   ...
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal "preview lists every unit and refusal and names the tag" → API
  Contracts 1, 3 and 5, Invariant 2, Surface Transcript 1.
- Goal "one confirmed change, tag created, never committed or pushed" → API
  Contracts 2 and 4, Invariants 4 to 6 and 8, Surface Transcripts 2 and 3.
- Goal "refusals never block; opt-out" → API Contracts 1, 3 and 6, Invariants
  7 and 9, Surface Transcript 4.
- Goal "upgrade names the Pending History" → API Contract 7, Invariant 10.
- Goal "list-of-maps unproven converts; one-line reasons" → API Contract 8,
  Invariants 11 and 12.
- User Story 1 → API Contracts 1, 3 and 5.
- User Story 2 → API Contracts 2 and 4.
- User Story 3 → API Contract 6.
- User Story 4 → API Contract 7.
- User Story 5 → API Contract 8.
- Core Feature 1 → API Contracts 1, 2 and 5.
- Core Feature 2 → API Contract 4, Invariants 5, 6 and 8.
- Core Feature 3 → API Contracts 1 and 3, Invariant 7.
- Core Feature 4 → API Contract 6, Invariant 9.
- Core Feature 5 → API Contract 7, Invariant 10.
- Core Feature 6 → API Contract 8, Invariants 11 and 12.
- Core Feature 7 → Build Order 1.
- Success Metric 1 → API Contracts 1, 3 and 5, Surface Transcript 1.
- Success Metric 2 → API Contract 4, Surface Transcripts 2 and 3.
- Success Metric 3 → API Contract 2, Invariant 1.
- Success Metric 4 → API Contract 6, Surface Transcript 4.
- Success Metric 5 → API Contract 7.
- Success Metric 6 → API Contract 8.

## Integration Points

- Git, through `gitOutput`: `rev-parse HEAD`, `status --porcelain=v1
  --untracked-files=all -- <paths>`, `cat-file -t`, `merge-base
  --is-ancestor`, `ls-tree -r --name-only -z`, the existing delivery lookup
  `spec.FindLegacyDelivery`, and `tag -a` at apply. The tag uses the
  repository's tagger identity and configuration; a Git failure there is a
  history write failure (API Contract 4).
- The configuration loader, for the Spec Root of the `--repo` root.
- The filesystem only, for the upgrade notice.

## Testing Approach

- `internal/spec/history_sanitize_unproven_test.go` (task_02, new) builds
  synthetic legacy folders in `t.TempDir()`, as
  `history_sanitize_legacy_test.go` does. It uses the two key sets of the
  Fluxus report, the string form, a sequence value, a nested mapping and an
  empty mapping, and an active Spec with the same front matter.
- `internal/cli/history_refusal_test.go` (task_02) gains a test whose refusal
  error holds a newline, through the existing `historyLegacyFolder` fixtures.
- `internal/cli/baseline_update_history_test.go` (task_03, new) combines the
  existing adopted fixture `newBaselineUpdateRepository` with legacy folders
  and a retired Finding written and committed under `docs/history/`. It calls
  `runBaselineUpdateCommandWithSkillsStage` with an environment from
  `commandEnvironmentForTest`, whose home is temporary, and a successful fake
  skills stage. It asserts on the result, the tree snapshot, `git tag`,
  `git cat-file -t`, `git rev-list --count HEAD` and the refs. The package
  already installs `suiteguard.Main`, so the new tests need no wiring.
- `internal/cli/upgrade_notice_test.go` (task_04) extends the existing notice
  fixtures (`prepareNoticeUpgrade`, `withCLIWorkspace`) with history written
  into the workspace, and with a work directory outside any `.git`.
- Existing tests are characterization: every `baseline update`, `history`,
  `upgrade` and archive record test passes unchanged.

## Build Order

1. `CONTEXT.md` terms, the three command references and the Roundfix Skill.
   They describe ADR-0254 and this TechSpec, so every code step names its
   guide.
2. Legacy `unproven` maps in `BuildArchiveRecord` and one-line Refused Unit
   reasons in `roundfix history sanitize` (`internal/spec`, `internal/cli`)
   (depends on: 1, which documents them).
3. The history section of `roundfix baseline update`: `planHistoryUnit`
   extracted from `runHistoryCommand`, the plan, the digest, `--no-history`,
   apply, text and JSON, and usage (depends on: 2, because both change
   `internal/cli/history.go` and the update prints step 2's one-line reasons).
4. The upgrade notice (depends on: 3, because it reads the inventory step 3
   keeps stable and shares the package's history helpers).
5. QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- A large Pending History makes one large change. Fluxus's 57 units remove
  23.9 MB. The diff is removals plus small records, and `--no-history` keeps
  the batch path.
- History retired after an existing tag can never be selected, because the
  tag does not hold it. It is listed as a Refused Unit on every update. A
  second tag, or a rule that lets later retirements cite `HEAD`, is a
  follow-up and is out of scope.
- An update run inside a Run worktree shares refs with the main checkout, so
  a created tag is visible there. This repository already has its tag and no
  Pending History, so its own managed refresh and module regenerations stay
  unaffected.
- `git tag -a` follows the user's configuration: signing, a missing identity
  or a hook can fail it. That failure is a history write failure before any
  unit is converted.
- Existing tests run the update with the process environment. A malformed
  real User Config can only block the history section, never the Baseline
  part, and the new tests use a temporary home.

## Glossary

- adds: **Pending History**
- changes: **Managed Refresh**
- changes: **History Sanitize Command**
- changes: **Refused Unit**
- changes: **Lenient Legacy Reading**
- changes: **Sanitize Batch**
- changes: **History Full Tag**

## Decisions

- The history section reuses the History Sanitize Command's per-unit
  planning, not a new planner; see ADR-0254.
- The digest stays the Baseline Plan Digest when nothing is selected, so every
  repository without Pending History is unaffected; see ADR-0254.
- A unit with uncommitted changes is refused rather than the whole update;
  see ADR-0254.
- The Baseline Plan applies before the tag and the units; see ADR-0254.
- The upgrade notice counts in process from the filesystem; see ADR-0254.
- Legacy `unproven` maps become one line each in a stable key order; see
  ADR-0254.
