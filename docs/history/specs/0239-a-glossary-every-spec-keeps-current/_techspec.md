---
spec: 0239-a-glossary-every-spec-keeps-current
prd: _prd.md
created: 2026-10-06
---

# A glossary every Spec keeps current — Technical Spec

## Executive Summary

The Spec Consistency Check gains one detector, in a new file of
`internal/speccheck`, that reads a Spec's Glossary Declaration, the bolded
phrases of its PRD and TechSpec, the Tasks that write the glossary, and the
glossary itself. It reports three errors: an undeclared term, a declared term
no Task plans to write, and a declared term still missing after its Tasks
completed. The full check, its staged form, the QA gate precondition, the
Delivery Queue plan and the Archive Command already act on Spec Consistency
Check errors or are given this detector, so one detector closes the loop at
every point. The CONTEXT workflow's two domain clauses are reworded in place,
the authoring, QA and Roundfix skills follow the rule, and one docs Task adds
the dropped terms. The trade-off accepted is a detector that sees only bolded
multi-word terms and declared ones, over a model that could guess at lowercase
prose but would not answer the same way twice (ADR-0244).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the three
  finding codes follow the existing `SC-<FAMILY>-<STATE>` form. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the detector and the refusal read
  local files and Git history; no request, credential or provider call is
  added and every test uses a temporary repository. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0244 (this Spec) governs the
  declaration, the bold rule, the binding Task, the three findings and the
  horizon. ADR-0094: "a Spec that passes today's gates must not be blocked by
  a new check that merely disagrees with its style", so a Spec authored before
  the horizon without the section is skipped, not failed. ADR-0117: "A defect is checked by the stage that can produce it", so the
  declaration is read at the PRD stage and the binding at the Tasks stage. ADR-0168's
  commit-ancestry horizon is the model the glossary horizon follows, and
  ADR-0183's receipts are unchanged. ADR-0200: "Spec authoring gets an
  advisory judge that never gates", so the judge gains no question. ADR-0186:
  "It never cites a Spec number or an ADR number", which the reworded clauses
  obey. ADR-0187 and ADR-0189 govern the skill edits and their versions, and
  ADR-0233 regenerates a raised version at merge. ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts", answered below with the reason none
  applies. ADR-0240 decides when this Spec's QA partial qualifies. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167; ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; ADR-0166, ADR-0178 and ADR-0182 bind each Task commit.
  ADR-0130: "every path any authorization has ever bounded must be matched by
  the declared set". The catch-up definitions come from ADR-0174, ADR-0196,
  ADR-0231, ADR-0232, ADR-0236, ADR-0238, ADR-0239 and ADR-0240, which are
  cited and not changed. ADR-0201 fixed the judge's ceiling that ADR-0231 made configurable, ADR-0237
  reads the Delivery Commit when a Delivery Retry records a merge made outside
  the queue, and ADR-0235 retired the Jev Router, which is why it is not
  added; these are sources too and do not change. ADR-0096 and ADR-0097 cite
  ADR-0080 but decide the gate's machine stage and row carry; ADR-0179 cites
  ADR-0130 but decides an empty paths list; ADR-0192 cites ADR-0178 but
  decides how a conflict confined to declared derived paths is resolved;
  ADR-0209 cites ADR-0200 but decides the grouping suggestion; ADR-0229 cites
  ADR-0167 but decides how an operator archive resumes a park; ADR-0194,
  ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when
  it is observed again and its evidence snapshot; ADR-0208, which ADR-0209
  cites, decides that sources sharing a context share a Spec; this Spec
  changes none of them, so none applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Spec Consistency Check
  sources and the corpus golden with its characterization, task_02 the CONTEXT
  workflow module and its derived files, and task_03 the owned skills; express
  maintainer authorization: "considere autorizado a ajustar todas as skills se
  necessário", "Autorizar os dois", and the maintainer's answer "Concedo" of
  2026-10-06 for the Governed Paths this Spec declares; bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/glossary.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`. No other Governed
  Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0239-a-glossary-every-spec-keeps-current/_authorization.md`.

## System Architecture

No new package, command or flag. Four existing components change:

- **Spec Consistency Check** (`internal/speccheck`): a new `glossary.go` holds
  the detector. `Check` in `constraints.go` calls it with the Task Graph when
  one is present; `CheckStage` in `coherence.go` calls it at the PRD and
  TechSpec stages, and its `stagedDetectors` table lists the three codes. The
  QA gate precondition (`specGatePrecondition`) and the Delivery Queue plan
  (`strictSpecFindings`) call `Check`, so they inherit the findings unchanged.
- **Archive Command** (`internal/cli/archive.go`): after the active-path-pin
  preflight and before `spec.Archive`, it calls `speccheck.GlossaryFindings`
  and refuses on any finding.
- **CONTEXT workflow module** (`internal/baseline/assets/modules/context-workflow.json`):
  two clauses of `rule.context.domain-docs` are reworded in place.
- **Owned skills**: `write-prd` gains the glossary guide that starts the
  horizon; `write-techspec`, `write-tasks`, `qa-gate` and `roundfix` gain the
  text below.

`CONTEXT.md` gains the catch-up entries.

## Implementation Design

### Interfaces

```go
// internal/speccheck/glossary.go
const (
	CodeGlossaryUndeclared = "SC-GLOSSARY-UNDECLARED"
	CodeGlossaryUnplanned  = "SC-GLOSSARY-UNPLANNED"
	CodeGlossaryMissing    = "SC-GLOSSARY-MISSING"
	// GlossaryGuidePath is the guide whose adding commit starts the glossary horizon.
	GlossaryGuidePath = ".agents/skills/write-prd/references/glossary.md"
)

// GlossaryFindings runs only the glossary detector, with the Task Graph, for one
// active Spec. The Archive Command reads it.
func GlossaryFindings(specsRoot, repoRoot, slug string) ([]Finding, error)

type glossaryEntry struct {
	Kind   string // "adds", "changes" or "not a term"
	Term   string
	Reason string
	Where  Location
}

func detectGlossary(result *Result, repoRoot string, artifacts []string, stage Stage, graph *spec.Graph) error
```

### Invariants

1. Glossary files are the root `CONTEXT.md` and `GLOSSARY.md` that exist, and
   every existing file named `CONTEXT.md` or `GLOSSARY.md` that a relative
   Markdown link in the root `CONTEXT-MAP.md` or `GLOSSARY-MAP.md` targets. A
   defined term is the text between a line-leading `**` and the following
   `**:`. Terms compare with case folded and whitespace collapsed. With no
   glossary file, the three codes are skipped with the missing input
   `CONTEXT.md or GLOSSARY.md`.
2. A line that is exactly `## Glossary` in the PRD or the TechSpec starts a
   Glossary Declaration that ends at the next line starting with `## `. Each
   non-blank line inside is `None.`, `- adds: **<term>**`,
   `- changes: **<term>**`, or `- not a term: **<phrase>** — <reason>` with a
   non-empty reason; `adds` and `changes` may be followed by ` — <note>`.
   `None.` is valid only as the section's single entry. Any other line is
   malformed and reports `SC-GLOSSARY-UNDECLARED` at its line.
3. When neither artifact declares the section, the detector applies the
   glossary horizon: the guide at `GlossaryGuidePath` exists and its adding
   commit is an ancestor of the commit that added the PRD, or the PRD is
   uncommitted, or history cannot be read, exactly as the contract horizon
   decides. Inside the horizon, one `SC-GLOSSARY-UNDECLARED` error names the
   PRD; outside it, the three codes are skipped with the missing input
   `a PRD committed at or after .agents/skills/write-prd/references/glossary.md`.
   A Spec that declares the section is checked at any age.
4. A bold candidate is a `**…**` span on one line of the PRD or TechSpec,
   outside fenced blocks and outside a Glossary Declaration, whose text has two
   to five words; contains no backtick, digit or parenthesis; does not end in
   `.`, `:`, `;`, `,`, `?` or `!`; and whose every word starts with an
   uppercase letter, except the connectors `a`, `an`, `and`, `at`, `by`, `for`,
   `from`, `in`, `of`, `on`, `or`, `per`, `the`, `to` and `with` after the
   first word. A candidate is covered when it equals a defined term, that term
   followed by `s` or `es`, or any declared term or phrase. Each uncovered
   candidate reports one `SC-GLOSSARY-UNDECLARED` error per distinct phrase and
   artifact, at its first line.
5. A Task binds a declared `adds` or `changes` term when it is not the QA Task,
   its Context declares a glossary file of Invariant 1 under `interface:` or
   `creates:`, and its `## Verification` text contains `**<term>**` exactly as
   declared. A term with no binding Task reports `SC-GLOSSARY-UNPLANNED`.
6. A `changes` term that the glossary does not define reports
   `SC-GLOSSARY-UNPLANNED`, whose fix is to declare it as added.
7. When every binding Task of a term has status `completed` and the glossary
   does not define the term, it reports `SC-GLOSSARY-MISSING`.
8. `SC-GLOSSARY-UNDECLARED` runs from the PRD stage, reading the PRD only at
   `--stage prd` and both artifacts after; `SC-GLOSSARY-UNPLANNED` and
   `SC-GLOSSARY-MISSING` run at the Tasks stage and in the full check and need
   the Task Graph, and without one they are skipped with the Task Graph as the
   missing input. Every finding is an error, so `--strict` changes nothing for
   them and the QA gate precondition blocks on each.
9. The Archive Command calls `GlossaryFindings` after the active-path-pin
   preflight and before any file moves, with or without `--qa-override`. With
   one or more findings it prints the existing preflight failure with the
   reason `Spec "<slug>" cannot archive with a Glossary Gap: <code>: <summary>`
   joined by `; ` for each finding, exits `2` and changes no file.
10. The detector reads only the Spec's PRD and TechSpec, its Task files, the
    glossary files and Git history; it writes nothing and calls no network.

### Exact clause texts

Each replacement is the clause's whole `guidance`. Every other byte of the
module stays as it is, each clause keeps its `id` and `enforcement`
(`mandatory`), and no clause gets a `replaces` list.

`clause.context.read-domain-contract`:

```text
The repository's selected domain context and its accepted ADRs are the base of the CONTEXT-driven workflow. Read both before naming domain concepts or changing behavior, and flag conflicts instead of silently overriding repository decisions.
```

`clause.domain.glossary-currency`:

```text
A Spec adds each domain term it introduces to the domain context, and revises each term it changes, through `domain-modeling` in one of its own Tasks, whose Verification names the term. It lists those terms in a `## Glossary` section of its PRD or TechSpec, with any bolded phrase it declares not a domain term. `roundfix spec check` reports a bolded term the domain context lacks and the section does not cover, and the QA gate and `roundfix archive` refuse a Spec whose declared term is still missing. When `domain-modeling` names its glossary file differently, it writes to the selected domain context. Work outside a Spec checks at its close whether it introduced, changed, or retired a term and updates the domain context when it did. Neither the check nor the update waits for human interaction; reach for `grilling` only when a term is ambiguous enough to need sharpening before it is written down.
```

### Version changes

Raise by one, from the value on task_02's starting commit, the version of the
`context-workflow` module, of `rule.context.domain-docs` and of `guide.domain`
(21, 5 and 6 on the authoring base).

### Retention

Both clauses keep their identity and enforcement, so `classifySourceClauseTransition`
records each `retained` for an adopter whose Setup Manifest declares the
Standard TypeScript Source Baseline; the Source Baseline corpus, manifest and
index stay byte-identical and no `replaces` is needed. The force record of
`internal/baseline/clause_characterization_test.go` already lists both clauses
as `mandatory` and does not change.

### Catch-up glossary entries

task_04 appends these entries, in this order, at the end of `CONTEXT.md`, and
adds the sentence below to the end of the first paragraph of the existing
**Spec Consistency Check** entry. Before writing each catch-up entry the Task
reads the ADR named beside it; an entry the ADR contradicts is corrected to the
ADR's wording, and the Result names each correction.

```markdown
**Glossary Declaration**:
The `## Glossary` section of a Spec's PRD or TechSpec that lists each domain term the Spec adds to or changes in the glossary, and each bolded phrase it declares not a domain term, or `None.`. A Task that declares the glossary file and names the term in its Verification writes each listed term.
_Avoid_: Vocabulary Contract, term list, glossary candidates

**Glossary Gap**:
A domain term a Spec introduces that the glossary does not carry as required: a bolded phrase the glossary lacks and the Glossary Declaration does not cover, a declared term no Task plans to write, or a declared term still missing after its Tasks completed. The Spec Consistency Check reports it, and the QA gate and the Archive Command refuse a Spec that has one.
_Avoid_: Missing term, undocumented token, stale glossary

**Light Tier**:
The dispatch tier that runs a Task on an open model through OpenCode and OpenRouter when its complexity is low, its type is not qa, and it declares no Governed Path; every other Task runs on the standard tier. A light Task whose first Verification fails escalates its one feedback turn to its category's Preferred Selection.
_Avoid_: Cheap model, Jev routing, light model profile

**Light Spend Log**:
The per-UTC-month record in Roundfix Home of the cost each Light Tier session reports, summed across every repository and compared with the Light Tier's monthly ceiling. A month that reached the ceiling, or whose log cannot be read, runs light Tasks on their category's profile.
_Avoid_: Judge Log, OpenRouter usage, spend cache

**Jev Ceiling**:
The monthly spending ceiling of the Jev judge, read from User Config only and compared with the Judge Log's sum across every repository on the machine; when it is unset a built-in default applies. The judge's key must carry a monthly credit limit no higher than it.
_Avoid_: Judge budget, key limit, Project Config ceiling

**Stage Key**:
The Roundfix environment variable an OpenRouter stage reads first, the judge's or implementation's own key, before the shared Roundfix OpenRouter key; the generic OpenRouter variable is never read. Each stage records the variable it used by name, never by value.
_Avoid_: API key value, generic OpenRouter key, shared key

**Tested Base**:
The default-branch commit a failed required check actually tested, read from the check job's annotation. A failure whose Tested Base does not contain the current default-branch tip is stale, and the Delivery Queue re-runs it instead of parking the item.
_Avoid_: PR base, merge base, Delivery Base

**Merge Evidence**:
Proof that a Spec was delivered: the default-branch head holds its archived PRD, and the Delivery Commit that added that file is not reachable from the Run Branch or Item Branch being proven. With it, reconcile releases that Spec's terminal Runs and Item Branches without comparing content.
_Avoid_: Content proof, merged flag, merge commit

**Delivery Commit**:
The default-branch commit that added a Spec's archived PRD. Merge Evidence and a Delivery Retry that records a merge made outside the queue both read it.
_Avoid_: Merge commit, archive commit, release commit

**Item Branch**:
The branch a Delivery Queue item builds its candidate on, named `roundfix/deliver-<slug>-<16 hex>`. Once no live item names it, Merge Evidence releases it, with its worktree only when that worktree is clean.
_Avoid_: Run Branch, PR Head Branch, deliver branch

**Network-Denied Row**:
An outside-evidence QA row blocked only because the Run sandbox denied network access, recorded with the host it could not reach. Like the pre-PR Pull Request row, it never decides a qualifying partial; whoever needs that proof declares it under Unreachable Acceptance.
_Avoid_: Failed outside evidence, skipped row, environment pass

**Pre-PR Review Command**:
The command that obtains the Pre-PR Review Provider's review of the current candidate before a Pull Request opens, reads the verdict from the reviewer's final message and keeps the Pre-PR Review Record. A finding parks delivery only after it is validated.
_Avoid_: Review Run, PR feedback Watch, review command alias
```

The sentence added to **Spec Consistency Check**:

```text
Its Glossary Declaration gap begins at the glossary horizon, the commit that added the glossary guide to the write-prd skill, and a Spec that carries a Glossary Declaration is checked at any age.
```

Sources: Light Tier and Light Spend Log, ADR-0238; Jev Ceiling, ADR-0231;
Stage Key, ADR-0239; Tested Base, ADR-0236; Merge Evidence, Delivery Commit
and Item Branch, ADR-0232; Network-Denied Row, ADR-0240; Pre-PR Review
Command, ADR-0174 and ADR-0196. Five candidates from the adopted Backlog Entry
and the replay are left out: "Evidence Manifest" names nothing in the code or
ADRs; "refused permission" is a test-case label; "Merge Observer" and
"Delivery Engine" are code components and "ACPX Runner" an adapter, which the
glossary format keeps out; the Jev Router was retired.

### Skill text

Each skill keeps its existing text, gains the content below under its own
heading, and never changes its `### QA settlement` section. The quoted
sentences are exact.

- `write-prd` → new heading `## Glossary Declaration`, after `## Claim
  receipts`. It states: "Whenever the PRD names a concept the glossary does not
  define, activate domain-modeling before writing it, even on the autonomous
  route, and sharpen the term there." It states that the PRD carries the
  `## Glossary` section of
  [references/glossary.md](references/glossary.md), that every new term is
  bolded where it is introduced, and: "When a person drives the authoring,
  grilling or grill-with-docs stays the interactive entry." The template
  `references/prd-template.md` gains a `## Glossary` section before
  `## Decisions` with a comment naming the three entry forms.
- `write-prd` → new file `references/glossary.md`: the entry grammar of
  Invariant 2, the bold rule of Invariant 4, the binding rule of Invariant 5,
  the three codes and when each is reported, the horizon of Invariant 3, the
  glossary files of Invariant 1, and the sentence: "A domain term the Spec
  introduces is written by one of the Spec's own Tasks, never left for a later
  Spec."
- `write-techspec` → new heading `## Glossary terms`: "When the TechSpec names
  a concept the glossary does not define, activate domain-modeling, sharpen the
  term, and declare it in the Glossary section of the PRD or of the TechSpec."
  The template `references/techspec-template.md` gains an optional
  `## Glossary` section before `## Decisions`.
- `write-tasks` → new heading `### Glossary requirement`, after the closing
  node's glossary paragraph: for each term the declaration adds or changes, the
  docs Task, or a dedicated glossary Task when the Spec has none, "declares the
  glossary file, carries a requirement to write the term through
  domain-modeling, and its Verification checks the bolded term followed by a
  colon in the glossary on whitespace-normalized text."
- `qa-gate` → new heading `## Glossary at the gate`, after `## 3. Run static
  gates first`: the strict precondition reports every Glossary Gap, and "The
  report records, for each term the Glossary Declaration adds or changes,
  whether the glossary defines it after the work."
- `roundfix` → `references/spec.md` describes the three codes, the declaration
  and the horizon; `references/archive.md` states: "Archive refuses a Spec with
  a Glossary Gap, also under a QA Archive Override."

### Data Models

None. No record, schema or payload changes; findings use the existing
`Finding` shape.

### API Contracts

1. API Contract: Glossary Declaration. `roundfix spec check <slug>` at every
   stage parses the declaration per Invariant 2 and requires it per Invariant
   3, reporting `SC-GLOSSARY-UNDECLARED` with the artifact and line.
2. API Contract: undeclared terms. The check reports each uncovered bold
   candidate of Invariant 4 as `SC-GLOSSARY-UNDECLARED`, naming the phrase,
   whose fix offers defining it, declaring it as added, or declaring it not a
   term.
3. API Contract: planned and present. With a Task Graph, the check reports
   `SC-GLOSSARY-UNPLANNED` per Invariants 5 and 6 and `SC-GLOSSARY-MISSING` per
   Invariant 7; the QA gate precondition blocks on each.
4. API Contract: archive refusal. `roundfix archive <slug>` refuses per
   Invariant 9, with or without `--qa-override`, and archives as before when
   the detector reports nothing or is skipped.
5. API Contract: glossary files. The detector reads the files of Invariant 1
   and skips with the named missing input when none exists.
6. API Contract: Baseline clauses. `docs/agents/domain.md` and the Standard
   TypeScript Monorepo golden render the two clauses of Exact clause texts as
   `mandatory`, and an adopter's Managed Refresh plan records both `retained`.
7. API Contract: skills. The skills carry Skill text, their mirrors equal their
   canonical files, and each raised version is recorded.
8. API Contract: catch-up. `CONTEXT.md` carries every entry of Catch-up
   glossary entries and the sentence on **Spec Consistency Check**.

### Surface Transcripts

None. The changed surfaces are new finding lines in the existing Spec
Consistency Check renderer and a new reason in the existing archive preflight
failure, whose surrounding bytes do not change. They are proved at the command
entry points with temporary repositories in task_01's tests, and the QA gate
runs both through the built binary on a temporary repository.

## Coverage Map

- Goal 1 → API Contracts 1, 2, 3 and 4.
- Goal 2 → API Contract 7.
- Goal 3 → API Contract 6.
- Goal 4 → API Contract 8.
- User Story 1 → API Contracts 3 and 8.
- User Story 2 → API Contracts 1 and 2.
- User Story 3 → API Contracts 3 and 4.
- User Story 4 → API Contract 6.
- Core Feature 1 → API Contract 1.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → API Contracts 3 and 4.
- Core Feature 4 → API Contract 5.
- Core Feature 5 → API Contract 6.
- Core Feature 6 → API Contract 7.
- Core Feature 7 → API Contract 8.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1 and 2.
- Success Metric 3 → Testing Approach 1 and 3.
- Success Metric 4 → Testing Approach 4.
- Success Metric 5 → Testing Approach 6.
- Success Metric 6 → Testing Approach 5.

## Integration Points

- Git history: the horizon reads the commit that added the PRD and the commit
  that added `GlossaryGuidePath` with the same helpers the contract horizon
  uses (`prdAddingCommit`, `adrHorizonGitOutput`).
- The upstream `domain-modeling` skill now names its file `GLOSSARY.md`; the
  check accepts both names, and the Baseline clause points the skill at the
  selected domain context.

Research record. The Secondbrain was read first (`wiki/index.md`, then a `qmd`
query on keeping the glossary current): it returned this repository's mirror
of the adopted Backlog Entry, the ubiquitous-language record adopted by Spec
0040, and the triaged changelog of the upstream skills' v1.3, which renamed
`CONTEXT.md` to `GLOSSARY.md` on 2026-10-05. Exa found Evans' Ubiquitous
Language pattern and Shore's chapter, recorded in `_prd.md` → Acceptance
evidence. The replay of the bold rule over archived Specs and the read-only
adopter measurement decided the detector's shape.

## Testing Approach

1. Detector seam: a new `internal/speccheck/glossary_test.go` builds temporary
   Git repositories with a glossary, a Spec, a Task Graph and, where needed,
   the guide at `GlossaryGuidePath`, and calls `Check` and `CheckStage`:
   `TestGlossaryDeclarationIsRequiredAfterTheHorizon`,
   `TestGlossaryMalformedEntryIsReported`,
   `TestGlossaryUndeclaredBoldTermIsReported`,
   `TestGlossaryBoldRuleIgnoresLabelsDigitsAndFences`,
   `TestGlossaryNotATermCoversAPhrase`,
   `TestGlossaryAddedTermNeedsATaskThatWritesIt`,
   `TestGlossaryChangedTermMustExist`,
   `TestGlossaryCompletedTaskMustLeaveTheTermDefined`,
   `TestGlossaryReadsMappedAndRenamedGlossaries`,
   `TestGlossaryStagesReportOnlyTheirCodes` and
   `TestGlossaryLegacySpecIsSkipped`.
2. Archive seam: a new `internal/cli/archive_glossary_test.go` with
   `TestArchiveRefusesASpecWithAGlossaryGap` and
   `TestArchiveWithQAOverrideStillRefusesAGlossaryGap`, each asserting exit
   `2`, the reason, and that the Spec directory did not move; and
   `TestArchiveAcceptsASpecWhoseGlossaryIsCurrent`.
3. Corpus: `TestStageScopeDefaultSweepIsUnchanged`, `TestCheckCorpusGolden`,
   `TestCheckActiveCorpusHasNoErrors` and
   `TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095` pass with
   the three codes added to the characterized set at zero.
4. Baseline: a new `internal/baseline/glossary_clauses_test.go` with
   `TestTheGlossaryClausesCarryTheirForceAndText`,
   `TestTheGlossaryClausesRenderInTheDomainGuide` and
   `TestAnAdopterRetainsTheGlossaryClauses`, beside the existing force,
   duplicate-text, record-citation, catalog, formatter and plan tests.
5. Skills: phrase checks on whitespace-normalized text, mirror equality,
   `make skills-sync-check` and `TestEveryOwnedSkillVersionIsRecorded`.
6. Catch-up: one phrase check per entry heading and the added sentence, on
   whitespace-normalized `CONTEXT.md`; the QA gate compares each catch-up
   definition with its ADR.

### Existing tests that change

- `internal/docscontract/corpus_test.go`: `corpusFindingCodes` gains the three
  codes.
- `internal/docscontract/testdata/corpus-golden.json`: `active` gains the three
  codes at `0`, and `update` gains one sentence naming them.
- `internal/spec/archive_layout_characterization_test.go`: the wanted golden
  gains the same keys and sentence.

No other existing assertion changes. A test the Task finds pinning the full
staged skip list may gain only the three new skip entries, and the Task names
it in its Result.

## Build Order

1. The glossary detector, its stage wiring, the corpus code set and the
   archive refusal, with the `spec` and `archive` command guides (API
   Contracts 1 to 5).
2. The reworded domain clauses, their versions, regeneration and Managed
   Refresh (API Contract 6).
3. The skill text, mirrors and recorded versions (API Contract 7) (depends on:
   1).
4. The catch-up glossary entries and this Spec's own terms (API Contract 8).
5. The final QA gate (depends on: 2, 3, 4).

## Risks & Considerations

- A Spec authored after the horizon that bolds a heading-like phrase sees an
  error until it declares the phrase not a term; the fix text says how, and
  the replay found no such phrase in thirty-eight Specs.
- A lowercase term is still missed unless declared; the skills and the
  Baseline clause carry that part, and the QA gate records the declared terms.
- Specs 0241 and 0242 may edit the same module and command; each version rises
  from its Task's starting commit, and the operator orders the queue.

## Vocabulary Contract

- emits: `internal/cli/archive.go`
  pattern: `Glossary Gap`
  documented-in: `docs/user-guide/commands/archive.md`

## Glossary

None.

## Decisions

- Detect only bolded phrases of two to five capitalized words, plus declared
  terms; see ADR-0244.
- Bind a declared term to the Task that declares the glossary file and names
  the term in its Verification; see ADR-0244.
- Require the section from the glossary guide's adding commit; see ADR-0244.
- Reword the two existing clauses in place rather than add a clause, so
  adopters retain both without a Source Baseline row.
- Accept `GLOSSARY.md` beside `CONTEXT.md` instead of renaming this
  repository's glossary.
