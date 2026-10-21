---
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
prd: _prd.md
created: 2026-10-21
---

# The Roundfix skill references agree with the CLI — Technical Spec

## Executive Summary

The reviewer's prompt takes its Delivery Conventions from `deliveryConventions`
in `internal/cli/review_conventions.go`, and the Run Event Stream takes its
categories from `runevent.AllStreamCategories` in `internal/runevent/stream.go`.
The Roundfix skill's `review` and `events` references were written against
earlier versions of both and were never checked against them, so they drifted:
C3 and C5 describe the folder-moving archive, and the category list leaves out
`usage`. The same drift sits in the `deliver` reference and user guide, which
name only an archived `_prd.md`, and in the glossary's Run Event Stream entry,
which names four filters.

task_01 corrects every such statement and quotes the conventions verbatim.
task_02 adds documentation contract tests that read the conventions and their
version by parsing the reviewer's source, and the categories from the
`runevent` package, and compare them with each document, and adds the Run
Event Stream entry to the glossary closed-list contract. No production code
changes.

The trade-off accepted: the convention guard parses a Go source file, so a
refactor that moves the literals out of `deliveryConventions` breaks the
guard with a parse failure until the test follows; the alternative, an
exported accessor, would be a symbol that exists only for tests, which the
repository forbids.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, schema or record
  shape is minted or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no decision is recorded and no
  behavior changes; the documents follow the code. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21 ("the same rules as cycle AM"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0294-the-roundfix-skill-references-agree-with-the-cli/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-21 at `956066a7`, and rechecked at
`8fb2da1f`.

1. `deliveryConventions()` (`internal/cli/review_conventions.go`) returns C1–C5
   as `{ID, Prompt}` literals under `deliveryConventionsVersion =
   "roundfix/delivery-conventions/v3"`; `buildReviewPrompt`
   (`internal/cli/review.go`) writes each as `<ID>. <Prompt>`. C3 reads "The
   archive commit adds an Archive Record at the archive root and removes the
   Spec folder, whose bytes stay in Git at source_revision; legacy archives
   move the folder and stamp its front matter." C5 ends "For an Archive Record,
   these fields are in the record and the unchanged QA Task and Report stay in
   Git at source_revision."
2. `.agents/skills/roundfix/references/review.md`, under "The prompt carries
   Delivery Conventions version", gives C3 as "The archive commit moves a
   completed Spec's directory to the archive root and stamps its archive front
   matter." and C5 without its Archive Record sentence. Its eligibility
   sentence reads "C5 applies to files under an archived Spec only when its
   `_prd.md` at the review head records `qa_override: true` and non-blank
   approval and reason." `docs/user-guide/commands/review.md` paraphrases C3
   and C5 in the Archive Record form but not verbatim (C3 is split into two
   sentences with `source_revision` in code font; C5 says "Archive Record or
   legacy archived `_prd.md` front matter").
3. `conventionRegions` (`internal/cli/review_conventions.go`) appends C3 for
   every anchor under the archive root, never under the active Spec root.
   It appends C5 when `reviewSpecQAOverride` holds for the anchor's Spec:
   for a direct `<archive-root>/<slug>.md` and for a file under
   `<archive-root>/<slug>/`. `reviewSpecQAOverride` reads the Spec through
   `readArchivedSpecAt`: an Archive Record qualifies when its QA override,
   which the record parser sets only for disposition `qa-override` with
   `qa_override: true`, has non-blank approval and reason; otherwise the
   folder's `_prd.md` front matter qualifies when it records `qa_override:
   true` and non-blank `qa_override_approval` and `qa_override_reason`. So a
   file under `<archive-root>/<slug>/` is judged by the record when
   `<archive-root>/<slug>.md` is one.
4. `runevent.AllStreamCategories()` holds `task-status`, `batch`,
   `verification`, `outcome`, `agent-selection` and `usage`;
   `ParseStreamCategoryFilter` accepts exactly those; `parseEventsArgs`
   (`internal/cli/events.go`) defaults to `AllStreamCategories()`. The built
   `bin/roundfix events --help` prints `--filter  Comma-separated categories:
   task-status,batch,verification,outcome,agent-selection,usage`.
   `TestDefaultFilterIncludesUsage` (`internal/runevent/stream_usage_test.go`)
   pins `usage` in the default.
5. `.agents/skills/roundfix/references/events.md` says "Default replay emits
   these public categories in journal cursor order: `task-status`, `batch`,
   `verification`, `outcome`, and `agent-selection`. `--filter` accepts a
   comma-separated subset of only those category names." Its stable-fields
   table has rows for the first four categories only. Its Token usage section
   says `usage` "is included in the default filter".
   `docs/user-guide/commands/events.md` already lists all six.
6. `ProjectStreamEvent` (`internal/runevent/stream.go`) sets, for
   `agent-selection`, `scope_kind`, `scope_id`, `scope_identity`,
   `selection_role`, `status`, `runtime`, `model`, `reasoning_effort` and
   `summary`, and for lifecycle events also `work_category`,
   `profile_source`, `attempt`, `fallback_index`, `next_runtime`,
   `next_model`, `next_reasoning_effort`, `reason_code` and `reason`;
   `projectUsageRecord` sets `scope_kind`, `scope_id`, `runtime`, `model`,
   `reasoning_effort`, `token_basis`, `summary` and, when recorded, `tokens`,
   `input_tokens`, `output_tokens`, `cached_read_tokens`,
   `cached_write_tokens`, `thought_tokens`, `cost_amount` and
   `cost_currency`. Every record may carry `batch` and `work_item` from the
   event.
7. The Delivery Queue's prerequisite check (the function that reads
   `requires` in `internal/cli/deliver_workflow.go`, after "One
   immutable-tree read checks every prerequisite") accepts a prerequisite
   whose `<archive-root>/<slug>/_prd.md` or Archive Record
   (`spec.ArchiveRecordPath`) is on the fetched default branch.
   `.agents/skills/roundfix/references/deliver.md` says it "requires each
   prerequisite's archived `_prd.md` there", and
   `docs/user-guide/commands/deliver.md` says "A prerequisite is met only when
   its archived `_prd.md` exists there."
8. `provenDeliveryEvidence` (`internal/worktree/merged_head.go`), used by
   `deliver retry`'s merged fallback, finds the commit that added the Archive
   Record, or the legacy `_prd.md` when the archive is a folder. The `deliver`
   reference and user guide say it "reads the Spec's archived `_prd.md`".
9. `CONTEXT.md`'s Run Event Stream entry says "Its stable Supervisor filters
   are task status, Batch boundary, Verification verdict, and terminal
   outcome."

Sweep, at `956066a7` and again over the paragraphs `8fb2da1f` added, of the
Roundfix skill (`SKILL.md` and every reference),
`docs/user-guide/commands/review.md`, `review-report-shape.md`, `events.md`,
the glossary and, because the sweep found the same fact there,
`docs/user-guide/commands/deliver.md`, with `rg -n -i` over `qa_override`,
`archived .?_prd`, moved folder or directory wording, `legacy archiv` and the
category names:

- Corrected by task_01: the statements in items 2, 5, 7, 8 and 9.
- Checked and left: `SKILL.md`'s settlement table row for `override`, which
  already names the Archive Record (inside the `### QA settlement` section,
  which no Spec may change); `archive.md`, which already describes the record
  and the Legacy Archive Folder conversion; `deliver.md`'s "outside the
  archived folder", which names the removed source folder; `deliver.md`'s
  retry rule "requires `qa_override: true`", which names no location;
  `baseline.md`, `setup.md` and `settle.md`, which describe Legacy Archive
  Folders or a Spec directory without claiming the archive moves one; the
  `events` user guide, which lists all six categories; and
  `review-report-shape.md`, which states no convention.

## System Architecture

- `.agents/skills/roundfix/references/review.md`, `events.md`, `deliver.md`
  and their mirrors: the corrected statements.
- `docs/user-guide/commands/review.md` and `deliver.md`: the same.
- `CONTEXT.md`: the Run Event Stream entry.
- `internal/docscontract/command_documentation_test.go`: the convention and
  category guards.
- `internal/docscontract/user_guide_contract_test.go`: the Run Event Stream
  closed list.

`internal/docscontract` already parses Go sources with `go/parser` (the QA
secret refusal guard) and imports internal packages; it gains an import of
`roundfix/internal/runevent`.

## Implementation Design

### Interfaces

```go
// deliveryConventionSource parses internal/cli/review_conventions.go and
// returns the conventions version and the ordered C1..Cn prompts.
func deliveryConventionSource(t *testing.T, root string) (string, []deliveryConventionText)

// checkDeliveryConventionDocumentation reports the first document that does
// not carry the version and every convention verbatim, or carries another.
func checkDeliveryConventionDocumentation(docs []documentedText, version string, conventions []deliveryConventionText) error

// checkEventCategoryDocumentation reports the first document whose anchored
// category sentence, or reference stable-fields table, differs from categories.
func checkEventCategoryDocumentation(docs []anchoredText, categories []string) error
```

### Data Models

None. No record, configuration or output shape changes.

### API Contracts

1. API Contract: the `review` documents. The canonical reference, its mirror
   and the user guide each keep "The prompt carries Delivery Conventions
   version `roundfix/delivery-conventions/v3`." and follow it with the five
   conventions as list items `- C<n>. <prompt>`, each prompt the exact text
   of `deliveryConventions()`, line breaks aside. The C5 eligibility sentence
   becomes the paragraph task_01 gives, which states item 3 of Current
   behavior.
2. API Contract: the `events` reference. Its default-replay sentence lists
   the six categories, its `--filter` sentence keeps "only those category
   names", and its stable-fields table gains an `agent-selection` row and a
   `usage` row with the fields of item 6 of Current behavior.
3. API Contract: the sweep. The `deliver` reference and user guide say a
   prerequisite is met by its Archive Record or, for a legacy archive, its
   archived `_prd.md` on the default branch, and that the merged fallback
   reads the Archive Record, or a legacy archive's `_prd.md`. The glossary's
   Run Event Stream entry names the six categories in code font and says
   replay includes all of them by default.
4. API Contract: the convention guard.
   `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem` parses
   `internal/cli/review_conventions.go` with `go/parser`, takes the value of
   the `deliveryConventionsVersion` constant and, from the body of
   `deliveryConventions`, every composite literal of two string literals
   whose first matches `^C[0-9]+$`, in source order; it fails when the
   version is empty or no convention is found. For each of the canonical
   `review` reference, its mirror and `docs/user-guide/commands/review.md`,
   on text whose whitespace runs are collapsed to one space, it requires
   ``Delivery Conventions version `<version>` ``, `- <ID>. <prompt>` for every
   convention, and exactly as many `- C<digits>. ` list items as conventions.
   A failure names the document and the convention ID.
   `TestTheDeliveryConventionContractReportsAStaleConvention` plants the old
   C3 sentence into the canonical reference's text and requires the check to
   fail naming that path and `C3`.
5. API Contract: the category guard.
   `TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem` takes the keys
   of `runevent.AllStreamCategories()`, sorted. For the canonical `events`
   reference and its mirror, anchored at "Default replay emits these public
   categories in journal cursor order:", and for
   `docs/user-guide/commands/events.md`, anchored at "The public Supervisor
   categories, in journal cursor order, are", it collects the code spans from
   the anchor to the first `.` or `;` and requires them, sorted, to equal the
   categories. In both reference copies it also requires a table row starting
   ``| `<category>` |`` for every category. A failure names the document and
   the missing or extra category.
   `TestTheEventCategoryContractReportsAMissingCategory` removes `usage` from
   the canonical reference's sentence and requires the check to fail naming
   that path and `usage`.
6. API Contract: the glossary closed list. `glossaryClosedLists` gains the
   term `Run Event Stream` with the values of `runevent.AllStreamCategories()`,
   so `TestGlossaryClosedListsEqualTheirConstants` requires the entry's code
   spans to be exactly the six categories, and
   `TestGlossaryContractReportsAPlantedViolationOfEachRule` keeps planting a
   conforming entry for it.
7. API Contract: no other change. The reviewer's conventions, the review
   validation, the Run Event Stream, the `events` help, the Delivery Queue and
   every other test keep their bytes and behavior.

### Surface Transcripts

None. No command's flags, help text, standard output or exit code changes.

### Invariants

```text
1. Each review document carries the conventions version and every convention exactly as the reviewer's prompt does, and no other convention item.
2. Each events document's category sentence names exactly the categories the stream replays by default and accepts in a filter.
3. The glossary's Run Event Stream entry names exactly those categories.
4. No production Go source, Archive Record or file under docs/history changes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariant 1 → the `review` documents (API Contract 1) and the convention guard (API Contract 4)
- Goal 2, Core Feature 2, Success Metric 2, Invariant 2 → the `events` reference (API Contract 2) and the category guard (API Contract 5)
- Goal 3, Core Feature 3, Success Metrics 3-4, Invariant 3 → the sweep (API Contract 3) and the glossary closed list (API Contract 6)
- Goal 4, Core Feature 4, Invariant 4 → API Contracts 4-7 and the QA scope audit
- Success Metric 5 → the PRD's Unreachable Acceptance

## Absence Evidence

1. The convention guard — command: `grep -rn -e deliveryConventions -e delivery-conventions internal/docscontract`; result: no match at `8fb2da1f`; only `internal/cli` tests read the conventions, and they compare the prompt with `deliveryConventions()`, never a document.
2. The category guard — command: `grep -rn -e runevent -e AllStreamCategories -e agent-selection internal/docscontract`; result: no match at `8fb2da1f`; no test compares a document with the stream's categories.
3. The Run Event Stream closed list — command: `grep -n 'term: "' internal/docscontract/user_guide_contract_test.go`; result: four closed lists at `8fb2da1f` (Task Type, Mechanical Refusal Code, Run Worktree Reconciliation, Selection Encoding), none for Run Event Stream.

## ADR Alternatives

None. The Active ADR obligations row lists no ADR: the Spec records no
decision and changes no behavior.

## Integration Points

None beyond the repository's own files: the tests read Markdown files and one
Go source file of the repository and import `internal/runevent`. No process,
provider or network endpoint.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, reads only repository
files from the root the package already resolves
(`baselineDocumentationRepoRoot`), and compiles for Windows. The new tests go
into the existing `//verify:always` files of `internal/docscontract`, so no
new file needs a Contract Relevance declaration.

- `internal/docscontract/command_documentation_test.go`: the four tests of
  API Contracts 4 and 5 and their helpers. Each guard has a planted-violation
  test, so removing a comparison makes a named test fail.
- `internal/docscontract/user_guide_contract_test.go`: the closed-list entry
  of API Contract 6; the existing glossary tests exercise it.
- Documentation: task_01 re-records `docs/references/behavior-surfaces.json`
  with the repository's record command, because the `review` and `deliver`
  guides' fingerprints change, raises and records the Roundfix skill version,
  and syncs the skill mirror.

## Build Order

1. The `review`, `events` and `deliver` references and user guides, the Run
   Event Stream glossary entry, the skill version, the mirror and the
   Behavior Surface Record (no dependency: the text states the contracts of
   `_techspec.md`, and its Verification reads only documents).
2. The convention guard, the category guard and the glossary closed list
   (depends on: 1, whose corrected documents the guards require; on the
   tree before 1 they fail).
3. QA gate (depends on: 1, 2).

No other Spec is a prerequisite. Spec 0291, merged in the authoring base,
changed other paragraphs of the `deliver` reference and user guide; task_01
finds its two sentences by their text.

## Risks & Considerations

- A refactor that builds the conventions differently, for example from a
  table outside `deliveryConventions`, makes the guard fail with "no
  convention found" until the test follows the source; that failure is the
  intended signal, not a false positive.
- The category guard compares sets, not order; the documents keep the
  declaration order of the constants by convention.
- Adopters keep the stale references until their next `roundfix baseline
  update`; the upstream skills repository's copy of the Roundfix skill is
  outside this repository and is not changed.

## Glossary

- changes: **Run Event Stream**

## Decisions

- Quote the conventions verbatim in the reference and the user guide, so the
  guard compares whole sentences.
- Parse the reviewer's source instead of exporting the conventions.
- Pin the glossary's Run Event Stream entry through the existing closed-list
  contract rather than a new test.
- Correct the `deliver` documents in this Spec, because they state the same
  stale archive form the sweep targets.
