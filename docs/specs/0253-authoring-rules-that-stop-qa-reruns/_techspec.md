---
spec: 0253-authoring-rules-that-stop-qa-reruns
prd: _prd.md
created: 2026-10-08
---

# Authoring rules that stop QA reruns — Technical Spec

## Executive Summary

This Spec changes Baseline content, two owned skills and two repository
documents, not Roundfix code. Three serial Tasks edit the `spec-workflow`,
`autonomous-work`, `go` and `context-workflow` modules. Each then runs the
Module Version Record step, `make baseline-digests` and this repository's
Managed Refresh. A fourth Task, independent of the modules, changes
`write-tasks`, the `write-techspec` concrete contracts guide and
`docs/agents/specific-repository.md`. The glossary and the user guide follow
the retirement change.

The trade-off accepted is ADR-0257's rule that a reworded clause keeps its
identity. Every Baseline rule here extends a clause in place, even where a
separate clause would read better, because a new identity in a rule the
Standard TypeScript profile requires needs a Source Baseline row that the
regenerator never creates. The one new identity is the Go hermetic-test
clause, because no Source Baseline holds the Go module. An authoring prototype
applied all three module changes to a scratch clone at `b16ca020`, ran the
sanctioned commands, and ran `internal/baseline`, `internal/cli`,
`internal/spec`, `internal/speccheck`, `internal/delivery` and `skills`, with
and without the `docscontract` and `repocontract` tags.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. Every
  reworded clause keeps its `id`; `clause.go.keep-tests-hermetic` is the one
  new identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Tests use the embedded catalog and temporary repositories,
  and the Managed Refresh runs locally against this repository. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0258 (this Spec) governs every
  decision below. ADR-0257: "A clause whose wording changes keeps its
  identity". ADR-0186: "Every clause keeps its enforcement level when its
  wording changes", and "It never cites a Spec number or an ADR number".
  ADR-0250: "A Baseline module's version names one content, and the record step
  chooses it". ADR-0248: "Retired Review Artifacts and handoffs have no reader,
  so they are removed". ADR-0254: "So is a unit whose paths an existing History
  Full Tag does not hold". ADR-0244: "A Spec declares the domain terms it
  introduces".  ADR-0251 reads a Legacy Archive Folder leniently and keeps a
  failed QA's verdict; this Spec changes neither. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois",
  the standing "Concedo", and the 2026-10-08 grants for the governed paths this
  Spec declares. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0253-authoring-rules-that-stop-qa-reruns/_authorization.md`.
  Bounded files: `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/concrete-contracts.md`,
  `docs/agents/autonomous-work.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Measured on 2026-10-08 at `b16ca020` (the Spec 0252 plan on `55de2a73`):

- `clause.spec.project-constraints-06-outside-evidence` asks for an outside
  source and lets a row stay blocked, but says nothing about where the source
  must live. Specs 0241–0245 cited operator-only logs and forbidden lookups.
- No clause or skill line asks the implementing test to assert every line of
  a Surface Transcript. `write-tasks` says "asserts the transcript's command,
  output, and exit text in a test".
- `_tasks.md` accepts `requires` (`internal/spec/requires.go`), and the
  Delivery Queue honors it, but no Baseline clause or skill tells an author to
  use it.
- `clause.autonomous.loop-01-qa-once` says to reopen the gate after a
  corrective Task without adding the Task to the gate's `needs`, and has no
  parallel re-check. `clause.autonomous.loop-04-verify-the-class` has no
  wording sweep.
- The `go` module has no hermetic-test clause.
- `clause.context.docs-one-job-per-directory` lists `docs/history/handoffs/`
  and `docs/history/reviews/` as destinations, asks "a byte-identical move for
  Review Artifacts and handoffs", and `clause.context.findings-09-archive` asks
  to "record the true disposition and its reason in a dated addendum" before
  the whole Finding moves. `docs/history/` holds only `adr/ backlog/ findings/
  specs/`, and `spec.IsReducedHistoryEntry` recognizes an entry by its final
  line `Full text in Git at <commit>: <path>.`.
- Spec 0252 leaves every clause this Spec changes byte-identical. It changes
  `rule.autonomous.loop` (loop-05), the root block of `context-workflow`, the
  spec-routing guide and `docs/agents/specific-repository.md` (the skill-sync
  bullet).

## System Architecture

```text
internal/baseline/assets/modules/*.json -> Module Version Record step -> module-versions.json
                                           make baseline-digests     -> formatter goldens, profile digest, testdata
                                           baseline update --repo .  -> docs/agents/*, setup-context.json
.agents/skills/<owned>/** -> skill version record -> make skills-sync -> skills/<owned>/**
docs/agents/specific-repository.md, docs/user-guide/context-driven-development.md, CONTEXT.md (hand-written)
```

No Go source outside tests changes.

## Implementation Design

### Interfaces

No exported or unexported production symbol changes. Each new test is a
package `baseline` test that reads the embedded catalog through
`mustEmbeddedCatalog` and `catalog.Asset`, builds a Source Baseline adopter
through `newClauseReplacementAdopter` and a Go adopter through
`buildTestPlan(t, newPlanRepository(t))`, as `glossary_clauses_test.go` and
`stack_force_go_cli_tui_test.go` do. The retirement test calls
`spec.IsReducedHistoryEntry` and `spec.PlanHistoryKind`.

### Data Models

None. No schema, record or front-matter field changes. The reduced entry is the
existing form of ADR-0248.

### API Contracts

1. API Contract: the Baseline clauses. The embedded catalog carries each
   change of Clause changes under its identity and enforcement. No reworded
   clause gains a `replaces` list, no identity is removed, and
   `clause.go.keep-tests-hermetic` is added as `mandatory` in
   `rule.go.observable-tests`.
2. API Contract: the rendered guidance. A Managed Refresh renders each added
   sentence once in its guide: `spec-routing.md`, `autonomous-work.md`,
   `go.md` and `docs-layout.md`.
3. API Contract: the Managed Refresh of an adopter. For a Source Baseline
   adopter it is ready, records each reworded clause `retained` and has no
   `unaccounted` clause. For a `go-cli-tui` plan the Go guide carries the new
   clause.
4. API Contract: the reduced retirement form. An entry written in the form
   `clause.context.findings-09-archive` gives, with a real 40-hex commit, is
   one `spec.IsReducedHistoryEntry` recognizes and `spec.PlanHistoryKind`
   plans no change for.
5. API Contract: the authoring guidance. `write-tasks`, the `write-techspec`
   concrete contracts guide and `docs/agents/specific-repository.md` carry the
   texts of Skill and document texts.

### Surface Transcripts

None. No command gains, loses or reorders an output line, a flag or an exit
code. The observable changes are the bytes of managed guidance, skills and
repository documents, which tests and phrase checks assert.

### Clause changes

Each change below edits one clause's `guidance` string in place; every other
byte of the module stays as it is, except the version lines in Version
changes. Each clause keeps its `id` and its enforcement.

task_01, `spec-workflow`.

`clause.spec.project-constraints-06-outside-evidence` (rule
`rule.spec.project-constraints`) inserts, after "so a later reader can tell it
apart from a rehearsal of the Spec's own premise.", these two sentences:

```text
Name a source the QA gate can check without the network before the Run starts: a document the planning change commits under `docs/references/`, or published literature whose cited passage the row quotes. Never rest the row on an artifact that exists only on one machine, on another repository's manifest, or on a lookup the Spec's authorization forbids, and declare a criterion that needs the network or an open Pull Request under Unreachable Acceptance when the row is authored.
```

`clause.spec.routing-05-task-graph` (rule `rule.spec.routing`) appends:

```text
The Task that implements a Surface Transcript asserts the whole transcript in one named test: every line of its standard output and standard error, in order, and its exit code, copied from the transcript with only its declared controlled variations; the Task's Verification runs that test by name.
```

`clause.spec.sources-02-bound-the-group` (rule `rule.spec.routing`) appends:

```text
When a Spec's Verification or premise depends on another Spec that has not merged, name that Spec in the Task Graph manifest's `requires` and state the dependency in the TechSpec Build Order.
```

task_03, `autonomous-work`, `rule.autonomous.loop`.

`clause.autonomous.loop-01-qa-once` inserts, before "Run the Implement Command
while implementation work is pending":

```text
Before a Spec authored beside another one enters a Run or the Delivery Queue, rebase it on the default branch and run `roundfix spec check <slug> --strict --run-verification` again, because a decision record or glossary term the other Spec merged can make it fail.
```

and replaces "When a corrective Task is added after the gate settled, reopen the
gate with `roundfix reopen --spec <slug>`; never edit the QA Task by hand."
with:

```text
When a corrective Task is added after the gate settled, add it to the QA Task's `needs` in the Task Graph and then reopen the gate with `roundfix reopen --spec <slug>`, because reopen returns the gate to `pending` only over a Task inside the gate's dependency closure; never edit the QA Task by hand.
```

Its first sentence, which `internal/delivery` parses, stays unchanged.

`clause.autonomous.loop-04-verify-the-class` appends:

```text
A Task that changes a message, field, flag, refusal, or rule searches the repository's documentation, agent guides, and skill copies for the old wording and declares every file that must change with it; the QA gate repeats the same search. Characterize current behavior before changing it, and declare each break: name every existing test, golden, or fixture the Task changes and the contract it moves to.
```

task_03, `go`, `rule.go.observable-tests`: a new clause
`clause.go.keep-tests-hermetic` (mandatory), placed after
`clause.go.test-observable-behavior`:

```text
Keep tests hermetic: set or clear with `t.Setenv` every environment variable the code under test reads, and never read the host's credentials, home directory, or tool state. Create a Unix socket under a short directory from `os.MkdirTemp("", ...)` rather than a deep `t.TempDir()`, because macOS refuses a socket path longer than 104 bytes, and never let a Verification depend on a test that can skip on the host.
```

task_04, `context-workflow`, `rule.context.docs-layout`.

`clause.context.docs-one-job-per-directory` replaces the text from "Each of
these directories holds live work only:" through "update dependent references
as part of that transition." with:

```text
Each of these directories holds live work only. The single history root `docs/history/` keeps only what a later reader needs, family by family (`docs/history/adr/`, `docs/history/backlog/`, `docs/history/findings/`, `docs/history/specs/`), and Git keeps every full text; `docs/history/handoffs/` and `docs/history/reviews/` hold only what an earlier layout retired there, which a history sanitize removes. Handoffs retire only on the user's explicit confirmation, and when confirmed, every handoff is deleted together and never enters history — the active `docs/handoffs/` directory is the capture door, never a shelf. Disposition is recorded in the front matter of the reduced entry for Findings and Backlog Entries, lifecycle front matter for ADRs, and the Archive Record or archive stamp for Specs. Rejected, deprecated and superseded ADRs retire whole to `docs/history/adr/`, while a `proposed` ADR stays in `docs/adr/`. A finished orphan Review Artifact retires by deletion and never enters history; a new Round never writes into history. Findings, Backlog Entries and Rollups with a terminal lifecycle status leave their active directory for the matching history family as reduced entries; preserve valid provenance and update dependent references as part of that transition.
```

The Archive Record sentence and the history sanitize sentence after it stay
byte-identical.

`clause.context.backlog-01-operational-contract` replaces "moves to
`docs/history/backlog/` in the same operation that records its terminal
status. Preserve its original intent, the true disposition, any consuming Spec
and the reason for closure." with:

```text
moves to `docs/history/backlog/` as a reduced entry in the same operation that records its terminal status. Its front matter carries the true disposition, any consuming Spec and the reason for closure, and the reduced entry keeps that front matter, the title, the first paragraph and the commit that holds the original intent, in the form the Findings archive rule gives.
```

`clause.context.findings-09-archive` makes three replacements. "Preserve the
original observation and record the true disposition and its reason in a dated
addendum." becomes the text below, whose middle lines are a fenced `markdown`
block inside the guidance string:

````text
Record the true disposition in the front matter, and write the history copy as a reduced entry: the front matter, the title, the first paragraph, and a final line naming a commit that already holds the full text at its active path, normally the commit before the retirement:

```markdown
Full text in Git at `<40-hex commit>`: `<active path>`.
```

An addendum written at retirement is committed in the active file first, so the named commit holds it; the original observation stays in Git and is never rewritten.
````

"Preserve its original observation and dated disposition; absence of a Spec"
becomes "Its original observation and dated disposition stay in Git; absence of
a Spec". "Preserve the prior routing in a dated addendum and validate" becomes
"Record the prior routing in a dated addendum committed before the Rollup is
reduced, and validate".

### Skill and document texts

task_02, `.agents/skills/write-tasks/SKILL.md`:

1. Under `## Surface Transcripts in Tasks`, "asserts the transcript's command,
   output, and exit text in a test." becomes "asserts the whole transcript in
   one named test: every line of its standard output and standard error, in
   order, and its exit code, copied from the TechSpec with only the
   transcript's declared controlled variations, and that Task's Verification
   runs the test by name. A line the reader might skip still counts: a final
   digest or summary line, a `Usage` block, indentation, and the
   `exit status <n>` line that `go run` writes to standard error when the
   program exits non-zero."
2. The bullet "**One acceptance row rests on evidence the Spec did not
   author.**" gains, after "and must record where that evidence came from.":
   "Name a source the QA gate can check without the network: a document the
   planning change commits under `docs/references/`, or published literature
   whose cited passage the row quotes. Never cite an artifact that exists only
   on the operator's machine, another repository's manifest, or a lookup the
   Spec's authorization forbids; a criterion that needs the network or an open
   Pull Request goes under Unreachable Acceptance when the row is authored."
3. A new bullet after "**Temporal prerequisites do not grant authority.**":
   "**A cross-Spec prerequisite is declared.** When a Task's Verification or
   premise depends on another Spec that has not merged, list that Spec in the
   manifest's `requires` frontmatter and state the dependency in the TechSpec
   Build Order; the Delivery Queue does not start the Spec until each
   prerequisite has merged."
4. The bullet "**Verification must be hermetic, portable, effect-proving, and
   Daemon-owned.**" gains, after its first sentence about fresh worktrees: "A
   Verification command never contains a backtick: the inline-code span ends
   at the first one and cuts the command, so match with a regular-expression
   class such as `.` instead."

No `### QA settlement` section of any skill changes.

task_02, `.agents/skills/write-techspec/references/concrete-contracts.md`,
under `## Surface transcript blocks`, a paragraph after the one that ends
"standard output and standard error are compared separately.":

```text
Copy each line from a run of the real command, not from memory. Keep the lines a reader might skip: a final digest or summary line, a `Usage` block, leading indentation and, for a `go run` command that exits non-zero, the `exit status <n>` line that `go run` adds to standard error. The implementing Task's test asserts every line, so a line left out here becomes a QA rerun later.
```

The paragraph holds no source-and-quote pair, so the guide keeps exactly one
Claim Receipt.

task_02, `docs/agents/specific-repository.md`: seven bullets inserted before
the bullet that begins "The release plan the baseline requires":

```text
- **HARD RULE — spawning test packages install the suite guard**: a Go test package whose tests start a process installs `suiteguard.Main` in its `TestMain`, carries `TestEverySpawningPackageInstallsTheSuiteGuard` in a `suiteguard_repocontract_test.go`, and is listed in `guardedSpawningPackages` in `internal/suiteguardcontract/contract.go`. A Task that adds such a package or its first spawning test declares all three files.
- Never add or change text inside the `### QA settlement` section of a skill: `TestSettlementGuidanceIsOneTable` requires it byte-identical in `qa-gate`, `archive-spec` and `roundfix`. New skill text goes under its own heading.
- A reworded Baseline clause keeps its identity. A removed clause is named in the `replaces` list of the clause that absorbs it, with the same enforcement, so the Source Baseline transition (`classifySourceClauseTransition` in `internal/baseline/plan.go`) records it `replaced`. A new clause in a rule the Standard TypeScript profile requires needs a Source Baseline row, so prefer extending an existing clause. A Task that removes a clause declares `internal/cli/baseline_update_test.go` and `internal/baseline/plan_test.go` when their fleet fixtures pin it.
- The canonical copy of an owned skill is `.agents/skills/<name>/`, and `make skills-sync` copies it to `skills/<name>/`. A Task that changes an owned skill declares the canonical file, its mirror, the `SKILL.md` pair whose version fields rise, and `skills/testdata/owned-skill-versions.json`.
- Derived Baseline files (the profile digest pin, the catalog snapshots and plan goldens under `internal/baseline/testdata/`, the formatter goldens and `docs/agents/setup-context.json`) change only through `make baseline-digests` and the Managed Refresh. A Task that edits module, template or profile source declares every derived file those commands rewrite, measured on a scratch copy.
- A Spec that edits a Baseline module declares two commands in the Sanctioned regeneration section of its `_authorization.md`: `make baseline-digests`, and the module record command with an explicit `outputs:` list holding `internal/baseline/module-versions.json` and each changed `internal/baseline/assets/modules/<name>.json`, without globs. Suiteguard refuses the module and record writes of an undeclared command.
- Measure the `paths:` of `_authorization.md`, never guess them: run `speccheck.GovernedPath` over every file the Tasks declare in a `go test -overlay` probe that writes nothing to the repository, and record `paths: []` when none is governed.
```

task_04, `docs/user-guide/context-driven-development.md`: in the paragraph
that begins "Move terminal Findings and Rollups", "Preserve the reason and any
valid absorption pointer, and update dependent references." becomes "Write each
as a reduced entry: its front matter, which keeps the reason and any valid
absorption pointer, its title, its first paragraph, and a final line naming a
commit that holds its full text. Delete a finished Review Artifact or a
confirmed handoff instead of moving it into history; Git keeps it. Update
dependent references."

task_04, `CONTEXT.md`, the **Reduced History Entry** entry becomes:

```text
A retired Finding or Backlog Entry cut to its front matter, title, first paragraph and the revision holding its full text. A retirement writes it directly, and the History Sanitize Command reduces the entries retired before (ADR-0248, ADR-0258).
```

Its `_Avoid_` line is unchanged.

### Version changes

Each Task raises by one, from the value on its starting commit, every version
it touches in this list. It then runs the Module Version Record step, which
chooses the module's own version (ADR-0250), so no Task names a module
version number.

- task_01: `rule.spec.routing`, `rule.spec.project-constraints` and
  `guide.spec-routing` in `spec-workflow`.
- task_02: the `write-tasks` and `write-techspec` skills through the skill
  version record.
- task_03: `rule.autonomous.loop` and `guide.autonomous-work` in
  `autonomous-work`; `rule.go.observable-tests` and `guide.go` in `go`.
- task_04: `rule.context.docs-layout` and `guide.docs-layout` in
  `context-workflow`.

### Retention

For an adopter whose Setup Manifest declares the Standard TypeScript Source
Baseline, the Managed Refresh records each reworded clause `retained`, because
its identity and enforcement are unchanged, and no clause `unaccounted`. The Go
module is not in that profile. The Source Baseline corpus, manifest,
accounting and index, and `internal/baseline/assets/retention/**`, stay
byte-identical.

### Derived files

The prototype measured the files the sanctioned commands rewrite:

- every module Task: `internal/baseline/module-versions.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest`, `catalog.normalized.json` and the four
  `testdata/plan-characterization/*.golden.json`, plus
  `docs/agents/setup-context.json`;
- task_01: the Standard TypeScript golden `docs/agents/spec-routing.md` and
  this repository's `docs/agents/spec-routing.md`;
- task_03: the golden `docs/agents/autonomous-work.md`, and this repository's
  `docs/agents/autonomous-work.md` and `docs/agents/go.md`;
- task_04: the golden `docs/agents/docs-layout.md` and this repository's
  `docs/agents/docs-layout.md`.

### Invariants

```text
1. No clause identity is removed and no reworded clause gains replaces; clause.go.keep-tests-hermetic is the only new identity.
2. No two Baseline clauses share text, and the force record lists exactly the catalog's clauses.
3. A Source Baseline adopter's Managed Refresh is ready and has no unaccounted clause.
4. The Source Baseline assets and the retention transition are byte-identical.
5. After each module Task, roundfix baseline update --repo . --no-skills reports current, and a second refresh changes no file.
6. The network-denied sentence, the Delivery Queue order sentence, the Archive Record sentence and the history sanitize sentence stay byte-identical.
7. The rendered docs-layout guide still names every history family directory spec.ArchiveDir returns.
8. Every module version is the one the record step chose; no hand-written pin, golden or guide.
9. Every owned-skill change raises both version fields through the record command, and skills/ mirrors .agents/skills/.
10. No Makefile, .roundfixrc.yml, CI workflow, go.mod, Source Baseline asset or production Go file changes.
```

## Coverage Map

- Goal "reachable outside evidence" → API Contracts 1, 2 and 5.
- Goal "whole transcripts" → API Contracts 1, 2 and 5.
- Goal "requires" → API Contracts 1, 2 and 5.
- Goal "loop and Go rules" → API Contracts 1, 2 and 3.
- Goal "repository rules" → API Contract 5.
- Goal "reduced retirement" → API Contracts 1, 2 and 4.
- User Story 1 → API Contracts 1 and 5.
- User Story 2 → API Contracts 1 and 5.
- User Story 3 → API Contracts 1 and 5.
- User Story 4 → API Contracts 1 and 2.
- User Story 5 → API Contract 5.
- User Story 6 → API Contracts 1 and 4.
- Core Feature 1 → API Contract 1, Build Order 1.
- Core Feature 2 → API Contract 1, Build Order 1.
- Core Feature 3 → API Contract 5, Build Order 2.
- Core Feature 4 → API Contracts 1 and 3, Build Order 3.
- Core Feature 5 → API Contract 5, Build Order 2.
- Core Feature 6 → API Contracts 1 and 4, Build Order 4.
- Core Feature 7 → Build Order 4.
- Success Metric 1 → API Contracts 2 and 3, Invariants 3 and 5.
- Success Metric 2 → API Contracts 2 and 3.
- Success Metric 3 → API Contract 5, Invariant 9.
- Success Metric 4 → API Contracts 2 and 4, Invariant 7.

## Integration Points

None. No external system is touched.

## Testing Approach

- `internal/baseline/authoring_evidence_clauses_test.go` (task_01, new) holds
  the three added task_01 texts as literals. It requires each once in its
  clause, with `mandatory` enforcement and no `replaces`, and once in the
  whitespace-normalized Standard TypeScript golden and this repository's
  `docs/agents/spec-routing.md`. It requires an adopter's refresh to record the
  three clauses `retained` and none `unaccounted`.
- `internal/baseline/loop_and_go_clauses_test.go` (task_03, new) holds the
  task_03 texts as literals. It requires the two loop clauses to carry their
  added sentences and no longer the replaced reopen sentence, the golden and
  this repository's `autonomous-work.md` to carry each once, the Go clause in
  the catalog as `mandatory`, a `go-cli-tui` plan's `docs/agents/go.md`
  postimage and this repository's `docs/agents/go.md` to carry it once as a
  `mandatory` bullet, and an adopter's refresh to retain both loop clauses.
- `internal/baseline/reduced_retirement_clauses_test.go` (task_04, new)
  requires the three context clauses to carry their new texts and no longer
  the replaced phrases, and the golden and this repository's `docs-layout.md`
  to carry them. It reads the fenced provenance line from
  `clause.context.findings-09-archive`, fills it with a 40-hex commit and a
  `docs/findings/` path, writes an entry of front matter, title, first
  paragraph and that line under a temporary `docs/history/findings/`, and
  requires `spec.IsReducedHistoryEntry` to accept it and `spec.PlanHistoryKind`
  to plan no file. The same entry without the line must be planned. It
  requires an adopter's refresh to retain the three clauses.
- Declared breaks, each changed only to the new contract:
  `internal/baseline/grouped_sources_clauses_test.go` (task_01, the
  `clause.spec.sources-02-bound-the-group` literal) and
  `internal/baseline/clause_characterization_test.go` (task_03, one force row
  for `clause.go.keep-tests-hermetic`). The prototype found no other failure
  in the packages named in the Executive Summary, including
  `TestDocsLayoutGuideNamesEveryHistoryFamily`,
  `TestLifecycleClausesCarryTheScopedWording`,
  `TestHistorySanitizeClauseIsAppended` and
  `TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow`.
- task_02 has no Go test of its own; its Verification checks each added
  phrase on whitespace-normalized text, the mirrors and the skill version
  record, and `TestTheGuideExamplesAreReadByTheCheck` proves the guide keeps
  its one receipt and transcript.

## Build Order

1. Outside evidence, whole transcripts and prerequisites in `spec-workflow`.
   It starts from Spec 0252's merge, because both change `spec-workflow`, the
   spec-routing guide and the derived files.
2. Authoring skills and repository rules. It starts from Spec 0252's merge,
   because 0252 also changes `docs/agents/specific-repository.md`; it shares no
   file with step 1.
3. Loop and Go rules (depends on: 1, because both change
   `module-versions.json`, `setup-context.json` and the derived files).
4. Reduced retirement (depends on: 3, for the same shared files).
5. QA gate (depends on: 2, 4).

## Risks & Considerations

- Spec 0252 must merge first. The Task Graph names it in `requires`, and every
  version this Spec raises starts from the value on the Task's starting
  commit.
- History Relocation still relocates a finished legacy Review Artifact under
  `docs/specs/_reviews/` or `docs/specs/reviews/` to `docs/history/reviews/`,
  where the update's Pending History removes it. The guidance governs new
  retirements; the code change is left for a Backlog Entry.
- A reduced entry names the commit before its retirement. An entry never
  committed must be committed whole first, which the findings clause says.
- Stale wording: the Task files name every file that states the old wording
  of each changed rule, found with the sweep the new loop-04 text requires.

## Glossary

- changes: **Reduced History Entry**

## Decisions

- Each rule extends a clause in place, and only the Go clause is new; see
  ADR-0258.
- A retirement writes a reduced entry or deletes, and the history sanitize
  sentence stays for older history; see ADR-0258.
- The docs-layout guide keeps naming `docs/history/handoffs/` and
  `docs/history/reviews/` as legacy locations a sanitize empties, because the
  history families the archive code knows are still listed; see ADR-0258.
