---
spec: 0209-sources-that-share-a-context-share-a-spec
prd: _prd.md
created: 2026-10-01
---

# Sources that share a context share a Spec — Technical Spec

## Executive Summary

Two independent changes meet in the authoring skills. The Baseline gains
three mandatory clauses, two in the Spec workflow module and one in the
CONTEXT workflow module, each with its Source Baseline row; the guides are
regenerated, never edited. The judge Spec 0205 builds gains a third judgment,
`source-grouping`, which pairs every Finding or Backlog Entry the Spec adopted
with every open Backlog Entry and unresolved Finding, asks the measured Noul
question, and raises a pair as `suggested` at 0.3 or more. The command prints
the suggestions and counts them in its summary. The `write-idea`,
`write-prd` and `write-techspec` skills teach the rules and how to answer a
suggestion.

The trade-off this design accepts is fidelity to the measurement over reach.
The question is asked exactly as measured, on text prepared exactly as
measured, and only between an adopted source and an open one, because that is
the shape the 0.3 threshold was measured on. Pairs of two open sources, pairs
inside one Spec and Inbox Entries are not asked, so the judge cannot propose a
new grouping from scratch or a split. The cost is low recall, which the
skills state.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Clauses keep fixed
  Baseline identifiers, a grouping judgment is named by its kind and the two
  source paths, and a request by the SHA-256 of its state as Spec 0205 names
  it. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — no new endpoint, header, key or
  field. A grouping request is the System One request of Spec 0205, sent
  through the transport its Invariant 2 selects: OpenRouter with
  `ROUNDFIX_OPENROUTER_API_KEY`, or else TypeSafe directly with
  `ROUNDFIX_TYPESAFE_API_KEY`, each read only from the command's environment
  and sent only in its own endpoint's authorization header. Only the state
  differs: two source texts. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0208 and ADR-0209 (this Spec),
  ADR-0200, ADR-0201, ADR-0083, ADR-0092, ADR-0186, ADR-0058, ADR-0060,
  ADR-0099, ADR-0203, ADR-0189, ADR-0193, ADR-0081, ADR-0149, ADR-0089,
  ADR-0184, ADR-0187, ADR-0182, ADR-0178, ADR-0166 and ADR-0167 hold as the
  PRD states, and ADR-0097, ADR-0116, ADR-0168, ADR-0176, ADR-0183,
  ADR-0192, ADR-0194 and ADR-0195 do not apply, for the reasons the PRD records.
  The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104,
  ADR-0117, ADR-0155 and ADR-0156, and by ADR-0093 and ADR-0094. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source, its guides and the skills, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/source-baselines/index.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`,
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`,
  `docs/agents/setup-context.json`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md`,
  `.agents/skills/write-prd/SKILL.md`, `skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `skills/write-techspec/SKILL.md`.
  Sanctioned regeneration:
  `make baseline-digests`, `make skills-sync` and the Managed Refresh
  `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Source of truth | Derived by |
| --- | --- | --- |
| Clause text, force, versions | `internal/baseline/assets/modules/spec-workflow.json`, `context-workflow.json` | hand edit under the grant |
| Source Baseline rows | the Standard TypeScript Monorepo corpus files, `manifest.json` row, `index.json` `entryIds` | hand edit of text, identity, force and carrier; the regeneration fills offsets, digests and the identity record |
| Goldens, profile digest pin, catalog snapshots, plan goldens | the sources above | `make baseline-digests` |
| This repository's guides and Setup Manifest | the sources above | the Managed Refresh |
| The grouping question, threshold and preparation | `internal/judge/questions.json` (Spec 0205) | hand edit |
| The grouping reader | `internal/judge/sources.go` (new) | — |
| Grouping pairs | `internal/judge/pairs.go` (Spec 0205) | — |
| Asking, the `suggested` outcome, the log line | `internal/judge/judge.go`, `internal/judge/log.go` (Spec 0205) | — |
| Output | `internal/cli/spec_judge.go` (Spec 0205) | — |
| Skills | `.agents/skills/{write-idea,write-prd,write-techspec,roundfix}` | `make skills-sync` |

A module edit keeps the file's formatting: insert clause objects and raise
version numbers in place, never re-encode a module. A disposable rehearsal of
task_01 on `c3be3bc9` measured the derived outputs: `make baseline-digests`
rewrote the two Standard TypeScript Monorepo goldens `docs-layout.md` and
`spec-routing.md`, the profile digest pin
`internal/baseline/assets/profiles/standard-typescript-monorepo.json`, the
Source Baseline `baseline.json`, `manifest.json` and `index.json`,
`catalog.diagnostics.golden.json`, `catalog.digest`,
`catalog.normalized.json` and the four plan goldens
`advisory-only-divergences`, `clean-adoption`,
`idempotent-replan-after-verified-apply` and
`same-baseline-changed-profile-and-catalog-digests`; the Managed Refresh
rewrote `docs/agents/docs-layout.md`, `docs/agents/spec-routing.md` and
`docs/agents/setup-context.json`, and a second refresh reported
`File changes: 0`. The whole `internal/baseline`, `skills` and Baseline CLI
test suites and `make docs-test repo-test` passed on the rehearsed tree. A
rehearsal of skill edits with raised versions showed that
`make baseline-digests` rewrites nothing after a skill edit.

## Implementation Design

### Exact clause texts

Each clause is `mandatory`. The two Spec-workflow clauses are appended to
`rule.spec.routing` after `clause.spec.verification-two-tiers`; the
CONTEXT-workflow clause is appended to `rule.context.docs-layout` after
`clause.context.inbox-02-fleet-flow`. The guides render a rule's clauses in
identifier order, so the two Spec-workflow clauses render between the routing
clauses and the Verification clause, and the extend clause after the
fleet-flow clause.

`clause.spec.sources-01-group-by-shared-context` — mandatory — "One Spec may
adopt several Inbox Entries, Backlog Entries, and Findings whose context is
similar or complementary: they change the same component or contract, share a
root cause, or one completes the other. One Spec per source is neither
required nor preferred, so before minting a Spec for a source, look among the
open Backlog Entries and unresolved Findings for others that share its
context. A grouping suggestion from `roundfix spec judge` is advisory: adopt
the suggested source or state why it stays apart."

`clause.spec.sources-02-bound-the-group` — mandatory — "Group sources into one
Spec only while the Spec fits four implementation Tasks plus its QA gate. When
the grouped scope needs more, split it into Specs that each fit and give each
source exactly one owning Spec; never add a source that shares no context with
the Spec only to save a Spec."

`clause.context.inbox-03-extend-before-minting` — mandatory — "Before minting a
Finding or Backlog Entry, including in Triage, look for an existing one that
is not yet implemented and that the new observation or intent fits: an `open`
Backlog Entry in `docs/backlog/` or an unresolved Finding in `docs/findings/`.
Revise a fitting Backlog Entry in place under its existing name; extend a
fitting Finding with a dated addendum, because a Finding is immutable history.
Mint a new file only when none fits. A Triage that extends an existing entry
resolves its Inbox Entry into that entry, which cites the Inbox Entry's
provenance."

No existing clause changes. The extend clause sits beside the Triage clause,
which still resolves one Inbox Entry into one Finding, one Backlog Entry or
one discard; the extend clause names which existing one that may be. None of
the three repeats an existing clause's text, cites a Spec or ADR number, or
names a mechanism the product lacks: `roundfix spec judge` ships with Spec
0205, which this Spec requires.

### Source Baseline rows

The catalog refuses a profile clause without a Source Baseline row
(`catalog.sourceBaseline.required-clause.missing`, observed in the
rehearsal), so each clause gains one. Each row is a corpus entry
`<!-- source-baseline-entry: <id> -->`, one line `- <text>`, and the closing
marker, after the named anchor entry with a blank line between entries; a
`manifest.json` row after the anchor's row with `kind` `normative-clause`,
`enforcement` `mandatory`, `carrier` `docs/agents/<file>`, `structure` null and
placeholder offsets and digest; and the identifier after the anchor in
`index.json` `entryIds`. The regeneration fills the rest.

| Clause | After | Corpus file | Text |
| --- | --- | --- | --- |
| `clause.spec.sources-01-group-by-shared-context` | `clause.spec.verification-two-tiers` | `spec-routing.md` | MUST look among the open Backlog Entries and unresolved Findings for sources that share a context before minting a Spec for one source, and let one Spec adopt every source whose context is similar or complementary. |
| `clause.spec.sources-02-bound-the-group` | `clause.spec.sources-01-group-by-shared-context` | `spec-routing.md` | MUST group sources into one Spec only while the Spec fits four implementation Tasks plus its QA gate, and split a larger grouped scope into Specs that each fit. |
| `clause.context.inbox-03-extend-before-minting` | `clause.context.inbox-02-fleet-flow` | `docs-layout.md` | MUST revise a fitting open Backlog Entry in place, or append a dated addendum to a fitting unresolved Finding, before minting a new Finding or Backlog Entry. |

Retention disposition: an adopter whose Source Baseline holds these rows has
each clause in the selected Baseline with the same identity and force, so
`classifySourceClauseTransition` records each `retained`. No clause is removed
or renamed, so no `replaces` entry is needed and no clause becomes
`replaced` or `unaccounted`.

### Version changes

Raise each by one from its value on the Task's starting commit:
`spec-workflow` module, `rule.spec.routing`, `guide.spec-routing`;
`context-workflow` module, `rule.context.docs-layout`, `guide.docs-layout`.
On `c3be3bc9` these are 13, 5, 9 and 20, 15, 14. Spec 0206 raises some of the
same numbers; whichever Spec lands second starts from the other's values.

### Existing tests that change

- `internal/baseline/clause_characterization_test.go`: the force record gains
  the three clauses as `mandatory`.
- `internal/baseline/preservation_test.go`: the maintained Source Baseline
  entry count rises by three from its value on the starting commit (138 on
  `c3be3bc9`).
- From Spec 0205: `internal/judge/questions_test.go` accepts the third
  judgment and its fifth pattern; `internal/cli/spec_judge_test.go` expects
  `0 suggested` in the summary of Transcripts 1, 4 and 5 of Spec 0205. No other
  existing assertion changes.

### The grouping question

`internal/judge/questions.json` gains one member of `judgments`, after
`goal-mechanism`, whose bytes are exactly the block below (the member
`goal-mechanism` gains the separating comma). The question text is the
measurement's, byte for byte; the threshold is its operating point; the
scrub pattern, the cut and the statuses are how its dataset was prepared.

```text
    "source-grouping": {
      "question_id": "same_spec",
      "question": {
        "type": "noul",
        "instructions": "Two work items from one software repository are given as `first` and `second` (findings or backlog entries). Should both be resolved by the same specification, delivered together? Answer yes when they change the same component or contract, share a root cause, or one item's fix depends on or completes the other's. Answer no when they are independent concerns that only share vocabulary or a date."
      },
      "suggest_when_noul_at_least": 0.3,
      "adopted_source_types": ["finding", "backlog"],
      "open_backlog_statuses": ["open"],
      "unresolved_finding_statuses": ["pending", "partial"],
      "source_scrub_pattern": "\\b0\\d{3}\\b|ADR-\\d{4}",
      "source_max_chars": 1500
    }
```

`Load` parses the member into `Questions.Grouping` and compiles
`source_scrub_pattern`, so it now compiles five patterns. No other file of the
package repeats one of these values.

### Interfaces

```go
// internal/judge/sources.go
// readGroupingSource reads one Finding or Backlog Entry the grouping question
// may send; ok is false when the file is not one (The grouping reader).
func readGroupingSource(dir, name string) (src Source, ok bool, err error)

// adoptedSources reads <specDir>/references/_index.md; openSources lists
// docs/backlog and docs/findings of repoRoot. Both return skipped artifacts.
func adoptedSources(q Questions, specDir string) ([]Source, []SkippedArtifact, error)
func openSources(q Questions, repoRoot string) ([]Source, []SkippedArtifact, error)

// prepareSource strips front matter, scrubs and cuts, as measured.
func prepareSource(q Questions, src Source) string
```

`Source` keeps the unexported fields Spec 0205 gave it, and only
`readSpecArtifact`, `readADR` and `readGroupingSource` build one.
`SkippedArtifact` is whatever value Spec 0205's `Plan` already holds for a
skipped artifact, with its `artifact` and `reason`.

### The grouping reader

`readGroupingSource` accepts a file only when every condition holds:

1. Its name ends in `.md`, is not `_index.md`, and it is a regular file
   directly in `dir` (`Lstat`, never following a symbolic link) of at most
   1 MiB.
2. Adopted sources: `dir` is `<specDir>/references`, the file is the `path` of
   a row of `references/_index.md` whose header is the adoption table's
   (`source | type | owner | adopted date | path`), and the row's `type` is in
   `adopted_source_types`. A row of another type, such as `inbox`, is a
   skipped artifact with reason `only Findings and Backlog Entries are sent`.
3. Open sources: `dir` is `<repoRoot>/docs/backlog` and the front matter
   `status` is in `open_backlog_statuses`, or `dir` is
   `<repoRoot>/docs/findings` and the front matter `status` is in
   `unresolved_finding_statuses`. A file with any other status, or without
   front matter, is not read and is not reported.
4. Its text after front matter passes Spec 0205's language gate; otherwise it
   is a skipped artifact with reason `not English`, and no pair includes it.

A refused file under rule 1 that a row or a directory listing names is a
skipped artifact with reason `not a regular file in its directory`. A Spec
without `references/_index.md` has no adopted sources and plans no grouping
judgment, silently. `docs/_inbox/`, `docs/history/`, other Specs' directories
and every other path are never listed.

### Preparing a source

`prepareSource` reproduces the measurement's preparation, in this order:

1. Remove a leading front matter block: from a first line `---` through the
   first following line `---` and its newline (the measurement's
   `(?s)^---.*?---\n`).
2. Remove every match of `source_scrub_pattern` (four-digit Spec numbers
   starting with `0`, and `ADR-` numbers), replacing it with nothing.
3. Keep the first `source_max_chars` Unicode code points. Nothing is trimmed.

### Grouping pairs

At every stage (`prd`, `techspec` and both), after the citation and goal
judgments Spec 0205 plans:

1. Anchors are the adopted sources in `_index.md` row order; candidates are
   the open Backlog Entries by path, then the unresolved Findings by path.
2. Each anchor and each candidate form one pending `source-grouping`
   judgment with state `{"first": <prepared anchor>, "second": <prepared
   candidate>}`, artifact the anchor's repository path, target the
   candidate's repository path, and no line.
3. Identical states are asked once, as Spec 0205 does for every judgment.

### Asking and outcomes

A grouping judgment is asked as Spec 0205's Asking section asks any
judgment, with the same transport, retries, stops, skips, Judge Log and
ceiling. Its answer is read as a Noul under `same_spec`:

- An answer whose reported model is not pinned is `skipped` with Spec 0205's
  reason and never suggested or cleared.
- A pinned answer with `noul` at or above `suggest_when_noul_at_least` is
  `suggested`; any other pinned answer is `clear`.
- A malformed answer is `skipped` as `unreadable answer`.

The Judge Log line of a grouping judgment has `"judgment":"source-grouping"`,
`"artifact"` the anchor, `"line":null`, `"target"` the candidate,
`"question_id":"same_spec"`, null `answer`, `probabilities` and `confidence`,
its `noul`, and `"outcome"` `suggested`, `clear` or `skipped`. Its schema stays
`roundfix/judge-log/v1`, and the month's cost sums it like any other line.

### Surface Transcripts

The fixture repository's Spec `0300-example` is English, has no ADR
citation and adopted one Finding,
`references/2026-09-20-run-events-grow-without-bound.md`, through its
`references/_index.md`. `docs/backlog/` holds two `open` entries,
`2026-09-28-prune-run-events-after-archive.md` and
`2026-09-29-color-the-tui-header.md`, and one `declined` entry; `docs/_inbox/`
holds one note; `docs/findings/` holds one `done` Finding. The command's
environment holds only `ROUNDFIX_OPENROUTER_API_KEY`. The fake transport
answers as OpenRouter does, `model` `typesafe/jev-1.13-20260917`, with
`input_tokens` 920 and a `usage.cost` of 0.00003864 per request, and a
`noul` of 0.81 for the pruning entry and 0.04 for the header entry.

1. Surface Transcript: one suggestion at the PRD stage.

   ```transcript
   $ roundfix spec judge 0300-example --stage prd
   stdout:
   suggested source-grouping docs/specs/0300-example/references/2026-09-20-run-events-grow-without-bound.md → docs/backlog/2026-09-28-prune-run-events-after-archive.md: P(same Spec) 0.81
   Judge: 0 advisory, 1 suggested, 1 clear, 0 skipped; 2 call(s), 1840 input tokens, US$0.0001; month US$0.0001 of US$5.00; model jev-1.13 via openrouter
   stderr:
   exit: 0
   ```

2. Surface Transcript: an open Backlog Entry written in Portuguese, added to
   the fixture as `docs/backlog/2026-09-30-cor-do-cabecalho.md`.

   ```transcript
   $ roundfix spec judge 0300-example --stage techspec
   stdout:
   skipped docs/backlog/2026-09-30-cor-do-cabecalho.md: not English
   suggested source-grouping docs/specs/0300-example/references/2026-09-20-run-events-grow-without-bound.md → docs/backlog/2026-09-28-prune-run-events-after-archive.md: P(same Spec) 0.81
   Judge: 0 advisory, 1 suggested, 1 clear, 0 skipped; 2 call(s), 1840 input tokens, US$0.0001; month US$0.0001 of US$5.00; model jev-1.13 via openrouter
   stderr:
   exit: 0
   ```

3. Surface Transcript: no Jev key; the environment holds only the generic
   `OPENROUTER_API_KEY`.

   ```transcript
   $ roundfix spec judge 0300-example --stage prd
   stdout:
   Judge: skipped: ROUNDFIX_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY); 2 judgment(s) not asked
   stderr:
   exit: 0
   ```

Transcript 2 runs on a fixture whose TechSpec has neither a citation nor a
Coverage Map line, so only grouping judgments are planned. Skipped artifact
lines precede judgment lines, as in Spec 0205's Transcript 4. A skipped
grouping judgment prints `skipped source-grouping <anchor> → <candidate>:
<reason>`; a clear one prints nothing in text form.

### API Contracts

1. API Contract: `roundfix spec judge <slug> [--stage <prd|techspec>] [--format <text|json>]`
   — unchanged flags, exit codes and help synopsis; the summary line becomes
   `Judge: <a> advisory, <s> suggested, <c> clear, <k> skipped; …` in every
   form that counts judgments, and a suggested pair prints the line of
   Surface Transcript 1.
2. API Contract: `--format json`, schema `roundfix/spec-judge/v1` — a
   grouping judgment is one element of `judgments` with `kind`
   `source-grouping`, `artifact` the anchor, `line` null, `target` the
   candidate, `text` and `section_title` null, `outcome` `suggested`, `clear`
   or `skipped`, `reason`, `answer`, `probabilities` and `confidence` null,
   `noul` and `model`. The document gains no field.
3. API Contract: the System One request — Spec 0205's API Contract 3 with the
   question `same_spec` and the state `{"first", "second"}`.
4. API Contract: the Judge Log line — Spec 0205's API Contract 4, with the
   grouping fields of Asking and outcomes.
5. API Contract: `roundfix spec judge --help` names the grouping question and
   that a suggestion never gates.
6. API Contract: rendered `docs/agents/spec-routing.md` and
   `docs/agents/docs-layout.md` of this repository and of every built-in
   profile that selects the Spec and CONTEXT workflows — the three clauses
   with `mandatory` force.
7. API Contract: `roundfix baseline update` for a Standard TypeScript
   Monorepo adopter of the Source Baseline — a ready plan that records each
   new clause `retained`.

## Coverage Map

- Goal 1 → Exact clause texts; API Contract 6.
- Goal 2 → Exact clause texts; API Contract 6.
- Goal 3 → Exact clause texts; API Contract 6.
- Goal 4 → Source Baseline rows; Version changes; API Contract 7.
- Goal 5 → The grouping question; Grouping pairs; Asking and outcomes; API Contract 1.
- User Story 1 → Exact clause texts; Testing Approach 5.
- User Story 2 → Exact clause texts; Testing Approach 5.
- User Story 3 → Exact clause texts; Testing Approach 5.
- User Story 4 → Grouping pairs; Surface Transcript 1.
- User Story 5 → The grouping reader; Testing Approach 3.
- User Story 6 → Source Baseline rows; API Contract 7.
- Core Feature 1 → Exact clause texts; API Contract 6.
- Core Feature 2 → Exact clause texts; API Contract 6.
- Core Feature 3 → Exact clause texts; API Contract 6.
- Core Feature 4 → Source Baseline rows; API Contract 7.
- Core Feature 5 → Testing Approach 5.
- Core Feature 6 → The grouping question; Grouping pairs; API Contract 3.
- Core Feature 7 → Asking and outcomes; API Contracts 1 and 2; Surface Transcript 1.
- Core Feature 8 → The grouping reader; Preparing a source; Surface Transcript 2.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 4.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 2.
- Success Metric 6 → Testing Approach 5.

## Integration Points

- **The Baseline catalog.** The new clauses ride the existing module, rule
  and guide records; no profile's `requiredRules` changes, because both rules
  are already selected by all three built-in profiles.
- **Spec 0205's judge.** The grouping judgment reuses its transport, model
  pin, Judge Log, ceiling, language gate and output. Nothing new reaches the
  network: tests and the QA gate inject the transport or run where nothing is
  sent.
- **The adoption table.** The reader parses `references/_index.md` with the
  header the `write-prd` skill writes and the Spec Consistency Check reads; it
  does not import `internal/speccheck`, so the check keeps no network
  dependency.
- **The authoring skills.** `write-idea`, `write-prd` and `write-techspec`
  gain one section each, `## Sources that share a context`, immediately
  before `## Process`, so the author reads it before choosing what a Spec
  adopts.

### Measured outside evidence

On 2026-10-01 the maintainer's session ran the grouping benchmark, outside
this Spec's Tasks: 107 sources adopted by 43 archived Specs, read from their
`references/` copies with front matter removed and Spec and ADR numbers
scrubbed, cut at 1,500 characters; 94 positive pairs, at most six per Spec,
and 94 hard negatives, each pairing a positive's first source with a source a
different Spec adopted within three days. The question of The grouping
question was sent as a Noul through OpenRouter, `jev-1.13`, which reported
`typesafe/jev-1.13-20260917`, for US$0.0073 in total. Jev reached AUROC 0.783
and TF-IDF cosine 0.753; the difference's 95% interval was [-0.032, 0.093]. At
a probability of 0.3 or more Jev's precision was 0.93 and its recall 0.27; the
TF-IDF baseline's 27 highest pairs had a precision of 0.85. The script and its
log are session evidence of 2026-10-01, outside this repository.

## Testing Approach

1. **The clauses.** New `internal/baseline/grouped_sources_clauses_test.go`
   requires each of the three clauses in its module and rule with `mandatory`
   force and the exact text of Exact clause texts; each text in the rendered
   Standard TypeScript Monorepo golden of its guide; each Source Baseline row
   with its force and carrier; and that a Standard TypeScript Monorepo
   adopter's plan is ready and records each `retained`. The existing force
   record, duplicate-text check, record-citation check, catalog validation,
   formatter composition, compatibility corpus and plan characterization
   stay green, and a second Managed Refresh changes no file.
2. **The question.** New `internal/judge/grouping_test.go` requires that the
   question file loads the `source-grouping` member with the measured text,
   threshold, statuses, pattern and cut, and task_02's Verification requires
   the member's bytes to equal the block under The grouping question.
3. **The boundary.** The same file requires that `readGroupingSource` refuses
   a symbolic link, a file over 1 MiB, `_index.md`, an `inbox` row, a
   `declined` Backlog Entry and a `done` Finding, and that a captured run whose
   repository also holds an Inbox Entry, a declined entry, a done Finding and
   a Go file with sentinel strings sends none of them: every string of every
   grouping state occurs in an accepted source's prepared text, once per
   transport.
4. **Pairs and outcomes.** The same file requires the preparation rules,
   one judgment per anchor and candidate in plan order, planning at every
   stage, no grouping without `references/_index.md`, `suggested` at 0.30 and
   `clear` at 0.29, a skip for an unpinned model, and the Judge Log fields of
   Asking and outcomes. New tests in `internal/cli/spec_judge_test.go`
   reproduce Surface Transcripts 1 to 3 and API Contract 2.
5. **Skills.** task_04's Verification requires the new section and its
   phrases in the three skills, each mirror equal to its canonical copy, and
   each raised version recorded.

Each new gate is proved to fail: the Task that adds it records, in its
Result, the sabotage it applied and the failing test, then restores the code.

## Build Order

1. The three clauses, their Source Baseline rows, the regeneration and the
   Managed Refresh, task_01 (depends on: none).
2. The grouping question, reader, preparation, pairs and outcome in
   `internal/judge`, task_02 (depends on: none). It extends the package Spec
   0205 creates.
3. The command's output, help, command reference and Roundfix Skill, task_03
   (depends on: 2).
4. The skill sections, task_04 (depends on: 1, 3). They teach the clauses
   task_01 ships and answer the suggestion task_03 prints, and they share
   `skills/testdata/owned-skill-versions.json` with task_03.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

This Spec is delivered after Spec 0205 (`requires` in the Task Graph).

## Risks & Considerations

- **Spec 0205 lands different from its TechSpec.** task_02 and task_03 extend
  code that does not exist yet. Each Task follows the shapes Spec 0205 shipped
  and stops to report when a named interface is missing rather than invent
  one.
- **Regex dialects.** Go's `\b` and `\d` are ASCII and Python's are Unicode,
  so a non-ASCII digit next to a Spec number scrubs differently. No archived
  source has one.
- **Cost per run.** One request per pair; 4 adopted sources and 15 open
  Backlog Entries are 60 requests, about US$0.002 and 25 seconds at the
  measured p50.
- **Repeated suggestions.** The question runs at both stages, so a declined
  suggestion returns at the TechSpec stage; the skills say to point at the
  earlier answer.
- **Spec 0206 overlap.** Both Specs edit the Spec workflow module, its guide
  and the Source Baseline files. Their clauses differ, the rehearsal of each
  starts from its own main, and whichever lands second rebases its version
  numbers and its entry count.

## Decisions

- Three new clauses, no change to an existing one. See ADR-0208.
- The suggestion, its threshold and its data boundary. See ADR-0209.
- Adopted-to-open pairs only, as measured.
- Grouping at every stage, because the bug-fix route never runs the PRD
  stage.
- A separate reader beside Spec 0205's two, so each refuses what the others
  accept.
