---
spec: 0251-skills-keep-up-with-the-behavior-they-describe
prd: _prd.md
created: 2026-10-08
---

# Skills keep up with the behavior they describe — Technical Spec

## Executive Summary

A new pure package, `internal/skillcoverage`, owns the Skill Coverage Map
(authored, `docs/references/skill-coverage.json`), the Behavior Surface Record
(generated, `docs/references/behavior-surfaces.json`) and the comparison of
two revisions. A `docscontract` test computes every Behavior Surface
fingerprint from the in-process command help, the configuration source, the
exit-code constants and the user-guide bytes, and keeps the record and the map
current. The Release Plan Command gains a `skill-coverage` check that reads
both files at the plan's base and target through Git and blocks a release
with a Lagging Surface by exit code 3 and its next action. The Spec
Consistency Check gains `SC-SKILLS-UNTASKED` and `SC-SKILLS-MALFORMED`, driven
by the map's source paths. The trade-off accepted: every change to a command's
help or a guide now re-records one generated file in the same change, and the
authoring rule sees only declared paths, so a change made only in the shared
usage file is caught at release, not at authoring (ADR-0256).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. Surface
  identifiers are readable names built from command paths and repository
  paths, the check statuses are lowercase words, and the new codes follow the
  `SC-` family. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added. The release check runs local Git reads, the record is computed
  in process, and every test uses temporary repositories and fixture Specs.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0256 (this Spec) governs the whole
  design. ADR-0189: "The step rests on checks that need no model and no
  network"; every new check is local and deterministic. ADR-0192: "conflict
  confined to declared derived paths is resolved by regeneration", so the
  record is declared derived. ADR-0233: "a conflict confined to declared paths
  takes the default branch's whole file", so the authored map is not. ADR-0222:
  "A Baseline guide says only what holds for the repository that reads it", so
  no Baseline module changes. ADR-0187: "A Task declares the one command file
  it changes", which the command coverage follows. ADR-0252: a `//verify:`
  directive "adds declared inputs to the package directory"; the new contract
  is `//verify:always`. ADR-0253: "A repository test asserts that a change to
  any declared output or line path selects the test", so the regeneration
  directive gains the new paths.
  ADR-0184: "A TechSpec states a command surface as a transcript", answered in
  Surface Transcripts. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-10-08: "Check no release + regra de autoria (Recommended)", the grant of
  the Governed Paths this Spec declares including skills, and the standing
  `qa_override` for an environment-only `partial`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0251-skills-keep-up-with-the-behavior-they-describe/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/release.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `.roundfixrc.yml`,
  `docs/agents/specific-repository.md`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/baseline_skill_contract_test.go`, `skills/roundfix/SKILL.md`,
  `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

Measured on 2026-10-08 at `c73a92e0`:

- `collectReleasePlanChecks` returns the `skills` and `baseline` checks, and
  `releasePlanExitCode` reads only the plan state. The `release plan` help says
  they "never change the decision state, the proposed version, or the exit
  code", and `TestReleasePlanChecksNeverChangeTheDecision` pins that over seven
  ranges whose fixture repositories hold no map.
- The root help yields 56 command paths through the `commandPaths` reader of
  `internal/docscontract`. Every `<path> --help` exits 0 except
  `qa-report accept --help`, which exits 1 and treats `--help` as a report
  path. `docs/user-guide` holds 38 Markdown files.
- The Roundfix skill's reference index maps every top-level command to one
  reference file. `TestEveryCommandIsNamedInTheRoundfixSkill` checks only that
  each command path is named somewhere in the skill.
- `SC-CLI-UNDOCUMENTED` accepts any file under `.agents/skills`, `skills`,
  `docs/user-guide` or `docs/agents` as a guide, so a CLI Task that names only
  a user-guide page passes without a skill.
- The operator's queue log, entries 193, 231 and 232, records QA failures in
  Specs 0235 and 0244 where the settle reference and the reopen guide trailed
  the shipped behavior.

## System Architecture

```text
internal/skillcoverage   (new, pure: parse, validate, fingerprint, match, compare)
  ^            ^                       ^
  |            |                       |
docscontract  internal/cli (release plan)   internal/speccheck
computes the  reads map+record at base and  reads the worktree map and the
record, keeps target via git show, git diff  Spec's declared paths
map current   --name-only for skill changes
```

- `internal/skillcoverage` is the only new package. It reads bytes and never
  runs Git or a process.
- `internal/docscontract/skill_coverage_test.go` is the record's generator and
  its contract, reusing `commandPaths`.
- `internal/cli` wires the check into `collectReleasePlanChecks`, the text and
  JSON printers and the exit code.
- `internal/speccheck/skills_declaration.go` holds the two detectors, wired
  next to the glossary detectors.

## Implementation Design

### Interfaces

```go
package skillcoverage

const (
	MapPath, RecordPath     = "docs/references/skill-coverage.json", "docs/references/behavior-surfaces.json"
	MapSchema, RecordSchema = "roundfix/skill-coverage/v1", "roundfix/behavior-surfaces/v1"
)

type Surface struct {
	ID        string   `json:"id"`
	Skills    []string `json:"skills,omitempty"`
	Uncovered string   `json:"uncovered,omitempty"`
	Sources   []string `json:"sources,omitempty"`
	Review    string   `json:"review,omitempty"`
}
type Map struct{ SchemaVersion string `json:"schemaVersion"`; Surfaces []Surface `json:"surfaces"` }
type Record struct{ SchemaVersion string `json:"schemaVersion"`; Surfaces map[string]string `json:"surfaces"` }
type Snapshot struct{ Map *Map; Record *Record } // nil when absent at that revision

func ParseMap(data []byte) (Map, error)
func ParseRecord(data []byte) (Record, error)
func EncodeRecord(record Record) []byte
func Fingerprint(content []byte) string // "sha256:" + 64 lowercase hex digits
func (m Map) SurfacesForPath(path string) []Surface
func Compare(base, target Snapshot, changedPaths map[string]bool) []Change
```

```go
type Change struct {
	ID      string   // surface identifier
	Kind    string   // "added", "changed" or "removed"
	Skills  []string // covering files of the entry that judged it
	Outcome string   // "described", "reviewed", "uncovered" or "lagging"
}
```

### Data Models

1. **Skill Coverage Map**, authored JSON. `schemaVersion` is `MapSchema`;
   `surfaces` is sorted by `id`, ids are non-empty and unique. Each entry has
   exactly one of a non-empty `skills` list or a non-blank `uncovered` reason.
   `skills` and `sources` hold clean, relative, slash-separated paths without
   `..`; a source ending in `/` matches everything under it and any other
   source is a `path.Match` pattern over the whole path, the grammar of
   ADR-0252. `review` is an optional Coverage Review, written
   `YYYY-MM-DD — <reason>`. Unknown fields are refused.
2. **Behavior Surface Record**, generated JSON. `schemaVersion` is
   `RecordSchema`; `surfaces` maps each id to `Fingerprint` of its bytes.
   `EncodeRecord` writes two-space indentation, keys in byte order and one
   trailing newline, so equal records are equal bytes.
3. **Behavior Surface identifiers and fingerprint bytes.**
   - `command <path>` for every path `commandPaths` reads from the root help:
     the stdout of `cli.Run(<path words> --help)`; when that exits non-zero,
     the stdout of the nearest ancestor path whose `--help` exits 0.
   - `config keys`: the unique `yaml` tag names of every struct field in the
     non-test Go files directly in `internal/config`, read with `go/parser`,
     excluding `-`, sorted, one per line.
   - `exit codes`: the constants in `internal/cli/cli.go` whose names start
     with `exit`, one `<name> <value>` line each, sorted by name.
   - `guide <path>` for every `.md` file under `docs/user-guide`: its bytes.
4. **Coverage of this repository.** Each `command` surface lists the
   Roundfix skill reference file that the skill's reference index assigns to
   its first word. A `guide docs/user-guide/commands/<family>.md` surface
   lists the same reference as its command. `SKILL.md` is listed only for a
   surface no reference covers, because its version line changes with every
   reference edit. `sources` of a `guide` surface include its own path; a
   `command` surface lists its user-guide page and its command-specific
   `internal/cli` files, never `internal/cli/cli.go`, which holds every
   command's usage text. `exit codes` may have no source.

### API Contracts

1. API Contract: `skillcoverage.Compare`. For every id in either record whose
   fingerprints differ (present in one only counts as `added` or `removed`),
   it returns one Change, ordered by id. The judging entry is the target
   map's for `added` and `changed`, the base map's for `removed`. The outcome
   is `uncovered` when the entry is uncovered; else `described` when any of
   its `skills` is in `changedPaths`; else `reviewed` when the kind is not
   `removed`, the target entry's `review` is non-empty and differs from the
   base entry's (absent counts as empty); else `lagging`. An id with no
   judging entry is `lagging` with no skills.
2. API Contract: the repository contract
   `TestTheSkillCoverageMapIsCurrent` (`//verify:always`, build tag
   `docscontract`). It computes every surface. With `-record-skill-coverage`
   it first writes the record. It then fails, naming each id, when the record
   misses a surface, holds an extra one or holds a stale fingerprint; when a
   computed surface has no map entry or an entry names no computed surface;
   when a covering file does not exist or lies outside
   `.agents/skills/<name>/` for an owned skill (a key of
   `skills/testdata/owned-skill-versions.json`); when a source matches no
   existing file or a guide entry omits its own path; and when a `command`
   entry omits the reference its first word maps to in the Roundfix skill's
   reference index.
3. API Contract: the derived declaration. `.roundfixrc.yml` gains
   `paths: [docs/references/behavior-surfaces.json]` with `regenerate:`
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   The `//verify:relevant` directive of
   `internal/config/regeneration_declared_test.go` gains
   `docs/references/behavior-surfaces.json` and
   `docs/references/skill-coverage.json`.
4. API Contract: the `skill-coverage` check of a range Release Plan. It reads
   `MapPath` and `RecordPath` at the base and target commits with
   `git cat-file -e` and `git show`, and the changed paths with
   `git diff --name-only --no-renames <base> <target>`. Status values:
   - `not_declared` when the target has no map: detail
     `the repository has no Skill Coverage Map`;
   - `introduced` when the target has one and the base has none: detail
     `the Skill Coverage Map is new since <base tag>; no surface is compared`;
   - `failed` when a present map or record cannot be read or parsed, the
     target has a map but no record, or Git fails: detail is the error;
   - `behind` when any Change is `lagging`: detail
     `<n> lagging surface(s) since <base tag>`;
   - `current` otherwise: detail
     `<c> changed surface(s) since <base tag>: <d> described, <r> reviewed, <u> uncovered`.
   `behind` and `failed` carry the next action
   `update a covering skill or record a Coverage Review in docs/references/skill-coverage.json for each lagging surface, then rerun roundfix release plan`
   and are blocking unless the state is `no_release`. The other statuses carry
   no next action and never block.
5. API Contract: text and JSON. The text prints, after the `baseline:` line,
   `skill-coverage: <status>: <detail>` with `; next: <action>` when present,
   then one line `- lagging: <id> (<kind>; covering skills unchanged: <skills joined by ", ">)`
   per Lagging Surface, or `(<kind>; no covering skill)` without skills. When
   the check is blocking, `Release blocked: skill-coverage` precedes the
   `Next action:` line, which reads
   `Next action: resolve the blocking skill-coverage check, then rerun roundfix release plan before any release mutation.`
   JSON gains `checks.skillCoverage` with `status`, `detail`, `nextAction`
   (omitted when empty), `blocking` and `lagging` (omitted when empty; items
   `surface`, `change`, `skills`). The schema version, the other checks, the
   decision state and the proposed version are unchanged.
6. API Contract: exit code. A blocking check turns exit 0 into exit 3; every
   other exit is unchanged. The `release plan` help keeps its sentence about
   the skills and baseline checks, adds that a blocking `skill-coverage` check
   exits 3 and names the check in its next action, and lists exit 3 as
   `approval_required, manual_classification_required, or a blocking skill-coverage check`.
7. API Contract: `SC-SKILLS-MALFORMED` (error, from the PRD stage). Under the
   PRD's optional `## Skills` section every non-blank line must be
   `- unchanged: <surface id> — <reason>` with a non-blank reason and an id the
   worktree map holds. Any other line, an unknown id or an empty reason is one
   finding at that line.
8. API Contract: `SC-SKILLS-UNTASKED` (error, from the Tasks stage). For each
   map surface that is not uncovered and has a source matching a path a
   non-QA Task declares under `interface:`, `creates:` or `deletes:`, the
   Spec needs a non-QA Task that declares one of its `skills` under
   `interface:` or `creates:`, or an `- unchanged:` entry for it plus a non-QA
   Task that declares `MapPath`. Otherwise one finding names the surface, the
   first matching Task and path, and the covering files. Both detectors skip,
   with the map path as the reason, when the worktree has no map, and skip
   with a horizon reason when the commit that added the PRD does not descend
   from the oldest commit that added the map; an uncommitted PRD is held.

### Invariants

```text
1. Only the record command writes the Behavior Surface Record; the contract fails on any stale fingerprint.
2. Every computed Behavior Surface has exactly one map entry, and every entry names a computed surface.
3. The release check reads Git objects at the plan's base and target, never the working tree.
4. A blocking skill-coverage check changes only the exit code and the Next action line; state, version and the other checks are unchanged.
5. A repository without a map and a base without one never block.
6. A Coverage Review excuses a surface only in the range in which its text changed.
7. The authoring rule reads declared paths only and never the working tree's changes.
8. No check calls a model or the network.
```

### Surface Transcripts

1. Surface Transcript: in a temporary repository tagged `v0.4.0` whose map and
   record cover one command, a later `fix:` commit changes that surface's
   recorded fingerprint and no covering skill.

   ```transcript
   $ roundfix release plan
   stdout:
   Decision: ready
   ...
   Release blocked: skill-coverage
   Next action: resolve the blocking skill-coverage check, then rerun roundfix release plan before any release mutation.
   skills: <status>
   baseline: <status>
   skill-coverage: behind: 1 lagging surface(s) since v0.4.0; next: update a covering skill or record a Coverage Review in docs/references/skill-coverage.json for each lagging surface, then rerun roundfix release plan
   - lagging: command implement (changed; covering skills unchanged: .agents/skills/roundfix/references/implement.md)
   ...
   stderr:
   exit: 3
   ```

2. Surface Transcript: the same range with a commit that also changes the
   covering skill file.

   ```transcript
   $ roundfix release plan
   stdout:
   Decision: ready
   ...
   Next action: release may proceed for v0.4.1 after independent release verification.
   skills: <status>
   baseline: <status>
   skill-coverage: current: 1 changed surface(s) since v0.4.0: 1 described, 0 reviewed, 0 uncovered
   ...
   stderr:
   exit: 0
   ```

3. Surface Transcript: a repository without a map.

   ```transcript
   $ roundfix release plan
   stdout:
   Decision: ready
   ...
   skill-coverage: not_declared: the repository has no Skill Coverage Map
   ...
   stderr:
   exit: 0
   ```

4. Surface Transcript: a fixture Spec committed after the map, whose Task
   declares a command source and no covering skill file.

   ```transcript
   $ roundfix spec check <slug>
   stdout:
   Spec <slug>
   ...
   [error] SC-SKILLS-UNTASKED: <summary>
   ...
   stderr:
   exit: 1
   ```

## Coverage Map

- Goal "refused before any release mutation" → API Contracts 4-6, Invariants 3-5, Surface Transcripts 1-3.
- Goal "checked-in, tested artifact" → Data Models 1-4, API Contracts 2-3, Invariants 1-2.
- Goal "reviewed change recorded once" → API Contract 1, Invariant 6.
- Goal "refused while authored" → API Contracts 7-8, Invariant 7, Surface Transcript 4.
- Goal "never refused without a map" → API Contract 4, Invariant 5, Surface Transcript 3.
- User Story 1 → API Contract 4, API Contract 5.
- User Story 2 → API Contract 8.
- User Story 3 → API Contract 1.
- User Story 4 → API Contract 2.
- User Story 5 → API Contract 4.
- Core Feature 1 → Data Model 3.
- Core Feature 2 → Data Model 1, Data Model 4.
- Core Feature 3 → Data Model 2, API Contract 2, API Contract 3.
- Core Feature 4 → API Contract 1, API Contracts 4-6.
- Core Feature 5 → API Contract 7, API Contract 8.
- Core Feature 6 → Build Orders 2-4.
- Success Metric 1 → API Contracts 4-6, Surface Transcript 1.
- Success Metric 2 → API Contract 1, Surface Transcript 2.
- Success Metric 3 → API Contract 4, Surface Transcript 3.
- Success Metric 4 → API Contract 2.
- Success Metric 5 → API Contract 7, API Contract 8, Surface Transcript 4.
- Success Metric 6 → Build Order 3, Build Order 4.

## Integration Points

- Git, through the release plan's existing `releasePlanGitSource` runner:
  `cat-file -e`, `show` and `diff --name-only` at two commits.
- The in-process CLI, through `cli.Run` with `--help`, inside the contract.
- The Delivery Queue's derived-path regeneration (ADR-0192) through the new
  `.roundfixrc.yml` declaration, and `TestRegenerationIsDeclared`, which runs
  every declared command in a repository copy.

## Testing Approach

- `internal/skillcoverage/skillcoverage_test.go` (task_01): table tests over
  in-memory maps and records. No process, so the package needs no suite guard.
- `internal/docscontract/skill_coverage_test.go` (task_02): the real
  repository contract plus fixture cases that feed a stale record and an
  unmapped surface to the same comparison helpers and assert the named ids.
- `internal/cli/releaseplan_skill_coverage_test.go` (task_03): the existing
  `newReleasePlanCommandRepo` fixtures with map and record files committed at
  the base and the target; it asserts exact lines, JSON and exit codes. The
  package already installs `suiteguard.Main`.
- `internal/speccheck/skills_declaration_test.go` (task_04): fixture
  repositories in `t.TempDir()` with committed map, PRD and Task Graph, as the
  glossary tests build them; the package already installs `suiteguard.Main`.

## Build Order

1. The `internal/skillcoverage` package and its tests.
2. The map and record of this repository, the contract and record command,
   the derived declaration, the regeneration directive, and the glossary terms
   **Behavior Surface**, **Skill Coverage Map**, **Behavior Surface Record** and
   **Coverage Review** (depends on: 1).
3. The release plan check, its help, the release runbook, the usage guide, the
   Roundfix skill's release reference, the record re-recorded, and the glossary
   term **Lagging Surface** with the revised **Release Plan** (depends on: 1, 2).
4. The two detectors, the corpus characterization, the write-tasks and
   write-prd skills, the Roundfix skill's spec reference, the spec command
   guide, the repository rule, the record re-recorded and the glossary term
   **Skills Declaration** (depends on: 1, 2, 3; it shares the record, the
   Roundfix skill entry, the skill version record and the glossary with 3).
5. QA gate (depends on: 4).

## Risks & Considerations

- The first release after this Spec has a base without a map and reports
  `introduced`; the check compares from the following release.
- A covering `SKILL.md` changes its version line whenever any of its
  references changes, which would mark every surface it covers `described`.
  Data Model 4 keeps `SKILL.md` off every surface a reference covers.
- A change made only in `internal/cli/cli.go`'s usage text matches no
  command's sources, so `SC-SKILLS-UNTASKED` misses it; `SC-CLI-UNDOCUMENTED`
  still asks for a guide and the release check catches the help change.
- `qa-report accept --help` exits 1 today; the nearest-ancestor rule
  fingerprints `qa-report --help` for it. The defect itself is outside this
  Spec.
- Parallel Specs that change help or guides conflict on the record; the
  declaration regenerates it at merge.

## Glossary

None.

## Decisions

- A new pure package rather than code inside `internal/cli`, because the
  release plan, the contract and the Spec Consistency Check all read it; see
  ADR-0256.
- Fingerprints live in a separate generated record, declared derived; see
  ADR-0256.
- The release check reads Git objects at both commits, so `--to` plans the
  same way as `HEAD`; see ADR-0256.
- A blocking check changes only the exit code and the next action; see
  ADR-0256.
- The new contract is `//verify:always`, because it costs about a second, like
  the other `internal/docscontract` contracts (ADR-0252).
