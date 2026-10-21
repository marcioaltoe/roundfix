---
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
prd: _prd.md
created: 2026-10-21
---

# The queue's archive subject and pull request title fit 100 characters — Technical Spec

## Executive Summary

`Archive` in `internal/cli/deliver_workflow.go` commits with `"docs: archive
"+specSlug`, and `Publication` returns the title `"feat: deliver "+specSlug`;
both grow with the slug, and a squash merge appends ` (#<number>)` to the
title. Two pure functions replace the concatenations: `archiveCommitSubject`
keeps the slug form while it is at most 100 UTF-16 code units and otherwise
names the Spec Number; `pullRequestTitle` keeps the slug form while it leaves
9 units for ` (#99999)` and otherwise names the Spec Number and the Spec
title, cut at a word. The Spec Number comes from one helper, extracted from
`pullRequestRecordSubject`, whose output keeps its bytes. `Publication` reads
the Spec title from the Archive Record only when the number form is needed.

The trade-off accepted: the title reserves a fixed width for a pull request
number it cannot know, so a repository past 99,999 pull requests can still
pass 100 by one unit per extra digit, in exchange for no forge write after
the pull request opens.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Spec Number is read from the existing slug (ADR-0302) and the
  Spec title from the existing Archive Record. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the title reaches the same `gh pr create`
  call, and every test uses temporary repositories, homes and fake forges.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1, 2 and 3 follow
  ADR-0302 and ADR-0309, named with their titles in `_prd.md` → Project
  Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21 ("the same rules as cycle AM"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-21 at `956066a7` and rechecked at
`8fb2da1f`, after Spec 0291 merged. Line numbers below are at `8fb2da1f` and
move when other Specs merge; the function names do not.

1. `commandDeliveryWorkflow.Archive` (`internal/cli/deliver_workflow.go:1220`)
   stages the archive and calls `workflow.commitQueueChange(ctx, gitRoot,
   specSlug, "docs: archive "+specSlug, false)` (`:1281`): 14 characters plus
   the slug. Behind a `commit-msg` hook that refuses a header longer than 100
   characters, Git 2.56.0 refused `docs: archive <Spec 0278's slug>`, 110
   characters, with exit status 1 and left the change staged; `docs: archive
   the 0278 Spec` committed. `commitQueueChange` (Spec 0290) then parks the
   item `delivery-error` with the hook and the hand commit.
2. When the item worktree's head is past the reviewed head, `Archive` calls
   `reconcileArchiveCommit` (`:1291`), which proves the head with
   `archiveCommitIsExact` (parent, paths and bytes) and never reads the
   commit message. A hand commit under any subject is adopted.
3. `commandDeliveryWorkflow.Publication` (`:1686`) returns `Title: "feat:
   deliver "+specSlug` (`:1707`) and `Body: "Delivers Spec "+specSlug+"."`
   (`:1708`). The Delivery Engine appends the Linux Lane paragraph
   (`LaneSkipBody`, `internal/delivery/engine.go:112`) to that body when the
   item carries a Lane Environment Skip, then `Engine.createPullRequest`
   passes title and body to `FindOrCreatePullRequest`
   (`engine.go:1356`), which looks an open pull request up by head and base
   branch (`internal/delivery/github.go:222`) and uses the title only in `gh
   pr create`. `MergePullRequest` runs `gh pr merge <n> --squash
   --match-head-commit <head>` with no subject, so the repository's default
   squash message applies: the title followed by ` (#<number>)`. On `main`,
   104 `feat: deliver` squash subjects exist; `9e92ac39` has 117 characters
   and `433a7c5a` 102.
4. `pullRequestRecordSubject` (`:1316`) reads the Spec Number inline: the
   text before the first `-`, when it is one to ten ASCII digits.
   `TestPullRequestRecordSubjectNamesTheSpecNumber`
   (`internal/cli/deliver_commit_hook_test.go`) pins its six cases.
5. The Spec title is the Archive Record's `title`, which `BuildArchiveRecord`
   (`internal/spec/archive_record.go:177`) takes from the PRD's first `# `
   line; `spec.ReadArchivedSpec(root, slug)` returns it for a record and for
   a legacy archive folder. `archivePaths` (`deliver_workflow.go:1897`)
   resolves the archive destination for the built-in and configured Spec
   roots.

Sweep of every reader of the two messages, at `8fb2da1f` (`rg -n 'docs:
archive |feat: deliver ' internal -g '!*_test.go'` finds only the two writers;
`rg -n -F '(#' internal -g '!*_test.go'` finds the readers):

- `spec.PullRequestFromSubject` (`internal/spec/source_revision.go:27`), used
  by Source Revision Recovery, matches only `\(#([1-9][0-9]{0,9})\)$`.
- `FindLegacyDelivery` (`internal/spec/history_sanitize.go:139`) matches only
  `\(#([0-9]+)\)$`.
- `releaseplan.ClassifyCommit` (`internal/releaseplan/classifier.go:10`)
  reads the type before the first `:` and a `!` in that header.
- The archive-commit adoption (item 2) and the pull request lookup (item 3)
  read no message.
- Text: `docs/user-guide/commands/deliver.md` and the Roundfix `deliver`
  reference describe neither message; `.agents/skills/archive-spec/SKILL.md`
  (`## Steps`) says "The Delivery Queue writes the commit subject `docs:
  archive <slug>`". `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`
  quotes historical squash subjects and stays as written.

## System Architecture

- `internal/cli/deliver_workflow.go`: `specNumber`, extracted from
  `pullRequestRecordSubject`; `queueHeaderLength`; `archiveCommitSubject`;
  `pullRequestTitle`; `Archive` and `Publication` call them.
- `docs/user-guide/commands/deliver.md`, `.agents/skills/roundfix/references/deliver.md`
  and its mirror, `.agents/skills/archive-spec/SKILL.md` and its mirror: the
  two messages.
- `CONTEXT.md`: **Spec Number**.

No package, import edge or file outside `internal/cli` changes in production.

## Implementation Design

### Interfaces

```go
// specNumber returns the one to ten ASCII digits that start specSlug when a
// "-" follows them (ADR-0302).
func specNumber(specSlug string) (string, bool)

// queueHeaderLength counts UTF-16 code units, as commitlint does.
func queueHeaderLength(header string) int

// archiveCommitSubject returns the queue's archive commit subject.
func archiveCommitSubject(specSlug string) string

// pullRequestTitle returns the queue's pull request title; specTitle is read
// only when the slug form does not fit.
func pullRequestTitle(specSlug string, specTitle func() (string, error)) (string, error)
```

### Data Models

None. The Archive Record, the Delivery Queue item, `delivery.Publication`
and the Pull Request body keep their shapes.

### API Contracts

1. API Contract: `specNumber`. Returns the text before the slug's first `-`
   when it is one to ten ASCII digits, else `"", false`.
   `pullRequestRecordSubject` uses it and returns exactly today's bytes for
   every slug, so `TestPullRequestRecordSubjectNamesTheSpecNumber` passes
   unchanged.
2. API Contract: `queueHeaderLength` and the limits. Length is the number of
   UTF-16 code units (`unicode/utf16`), the JavaScript `length` commitlint
   reads. The header limit is 100; the squash reserve is ` (#99999)`, 9
   units, so a title is at most 91.
3. API Contract: `archiveCommitSubject(slug)`. `docs: archive <slug>` when its
   length is at most 100; otherwise `docs: archive the <spec number> Spec`;
   otherwise `docs: archive the Spec`. `Archive` commits with it through
   `commitQueueChange`, so a hook refusal still names this subject.
4. API Contract: `pullRequestTitle(slug, specTitle)`. `feat: deliver <slug>`
   when its length is at most 91, without calling `specTitle`. Otherwise the
   prefix is `feat: deliver <spec number>: `, or `feat: deliver ` without a
   Spec Number, and the title text is `specTitle()` with each whitespace run
   collapsed to one space and the ends trimmed. When prefix plus text passes
   91, the text is cut to the longest run of whole runes that fits; if the
   next rune is not a space and the cut text holds a space, it is cut again
   before its last space. Trailing spaces, `.`, `,`, `;` and `:` are then
   removed, cut or not. A non-empty text gives prefix plus text; an empty one
   gives `feat: deliver <spec number>`, or `feat: deliver the Spec`. An error
   from `specTitle` is returned wrapped as `read Spec title: %w`.
5. API Contract: `Publication`. It builds the title with `pullRequestTitle`,
   passing a reader that resolves the archive destination with
   `archivePaths` and returns `spec.ReadArchivedSpec(<repository>/<archive
   root>, slug).Record.Title`; a reader error fails planning as `plan
   publication: read Spec title: <error>`, before any forge call. The body
   stays `Delivers Spec <slug>.`. The Linux Lane paragraph the Delivery
   Engine appends to the body is untouched.
6. API Contract: unchanged readers. `reconcileArchiveCommit`,
   `archiveCommitIsExact`, `FindOrCreatePullRequest`,
   `spec.PullRequestFromSubject`, `FindLegacyDelivery` and
   `releaseplan.ClassifyCommit` keep their bytes; tests prove each still
   reads the new messages.
7. API Contract: documentation. The `deliver` guide and the Roundfix
   `deliver` reference state both rules with their bounds; the
   `archive-spec` skill's sentence names both archive subject forms; the
   glossary defines **Spec Number**.

### Surface Transcripts

None. No command's flags, help text, standard output or exit code changes.
The archive subject and the title reach an operator only in Git history and
on the forge, which only a queue running against a forge produces; the tests
of API Contracts 3-5 assert the exact strings instead.

### Invariants

```text
1. The archive subject is at most 100 UTF-16 code units for every slug and equals "docs: archive <slug>" whenever that fits.
2. The title plus " (#99999)" is at most 100 units for every slug and title, and equals "feat: deliver <slug>" whenever that is at most 91.
3. The Spec Number is read by one helper, and the Pull Request Record Commit subject keeps its bytes.
4. A short slug's publication reads no Archive Record.
5. An archive commit is adopted by its change alone; its message is never read.
6. Every squash subject built on a title still yields its trailing number to Source Revision Recovery and the type feat to the release plan.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1-2, Invariant 1 → `archiveCommitSubject`, `Archive` (API Contracts 2-3)
- Goal 2, Core Feature 2, Success Metrics 3 and 6, Invariants 2 and 4 → `pullRequestTitle`, `Publication` (API Contracts 2, 4-5)
- Core Feature 3, Invariant 3 → `specNumber` (API Contract 1)
- Goal 3, Core Feature 4, Success Metrics 4-5, Invariants 5-6 → the unchanged readers and their tests (API Contract 6)
- Goal 4, Core Feature 5 → the `deliver` guide, the Roundfix reference, the `archive-spec` skill and **Spec Number** (API Contract 7)
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `specNumber`, `queueHeaderLength`, `archiveCommitSubject`, `pullRequestTitle` — command: `grep -rn -e 'func specNumber' -e queueHeaderLength -e archiveCommitSubject -e 'func pullRequestTitle' -e utf16 internal/cli`; result: no match at `8fb2da1f`; the Spec Number is read inline in `pullRequestRecordSubject`, and both messages are inline concatenations in `Archive` and `Publication`.
2. A length rule for the archive subject or the title — command: `grep -rn -e '"docs: archive "' -e '"feat: deliver "' internal -g '!*_test.go'`; result: only `deliver_workflow.go:1281` and `:1707`, both plain concatenations with no bound.

## ADR Alternatives

- ADR-0302: no rejected alternative is taken; the slug is never cut, no `--no-verify` is passed, and the Spec Number rule and the hook refusal are reused as ADR-0302 defines them.
- ADR-0309: no rejected alternative is taken; short slugs keep the slug form, the reserve is ` (#99999)` rather than ten digits, and no pull request is renamed after it opens.

## Integration Points

Git, for the archive commit and the hooks it runs in the item worktree, and
the forge, through the unchanged `gh pr create` call that receives the title.
No test, Verification command or QA row reads the real `~/.roundfix`,
another repository or a forge.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, homes and Git configuration, writes hooks with
the package's script fixture helper, and waits on no real time. Tests that
run real Git commits re-run themselves under an isolated home with the
existing `queueCommitTestHome` helper.

- `internal/cli/deliver_queue_subject_test.go` (new):
  - `TestQueueArchiveSubjectKeepsTheSlugWhileItFits`: an 86-unit slug (slug
    form, 100), an 87-unit slug (`docs: archive the <n> Spec`), a 230-unit
    slug, a ten-digit Spec Number, an eleven-digit prefix and a slug without
    digits (`docs: archive the Spec`); each result at most 100 units.
  - `TestQueuePullRequestTitleKeepsTheSlugWhileItFits`: a 77-unit slug (slug
    form, 91, and the title reader is never called); a 78-unit slug; Spec
    0278's slug and title (the 90-unit title of Success Metric 3); a title
    that fits whole and ends with `.` (dropped); a title with runs of
    whitespace; a first word longer than the budget (cut inside it); a title
    holding an astral rune counted as two units; an empty title with and
    without a Spec Number; a slug without a Spec Number; a reader error
    (wrapped `read Spec title:`). Each title plus ` (#99999)` is at most 100
    units.
  - `TestQueueMessagesStayReadableBySquashReaders`: for each title of the
    previous test, `spec.PullRequestFromSubject(title+" (#575)")` returns
    `575` and `releaseplan.ClassifyCommit` gives type `feat`, a minor impact
    and no breaking change; each archive subject form classifies as `docs`.
  - `TestDeliveryWorkflowArchiveCommitPassesAHeaderLimitHook`: the archived
    item fixture of Spec 0290's hook test with a 100-character slug starting
    `0296-` and a `commit-msg` hook refusing a first line longer than 100;
    `Archive` returns an exact head whose subject is `docs: archive the 0296
    Spec`, and the worktree is clean.
  - `TestDeliveryWorkflowAdoptsAnArchiveCommitWhateverItsSubject`: the same
    fixture with a refusing `pre-commit` hook; after the refusal the hook is
    made to pass, the staged archive is committed by hand under `docs:
    archive <slug>` (the old long subject), and a second `Archive` returns
    that head as an exact move with no further commit.
  - `TestDeliveryPublicationTitleNamesTheSpecNumberForALongSlug`: an
    archived item fixture whose slug has 100 characters and whose Archive
    Record title is known; `Publication` returns `feat: deliver <n>: <cut
    title>` and `Delivers Spec <slug>.`; with the record removed it fails
    with `plan publication: read Spec title:`.
- Unchanged and run by name: `TestPullRequestRecordSubjectNamesTheSpecNumber`,
  `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery` and
  `TestDeliveryPublicationRequestsNoCodeRabbitReview` (a short slug plans
  with no Archive Record, proving Invariant 4).
- No existing test, golden or fixture changes: every slug the existing tests
  deliver is shorter than 78 characters.

## Build Order

1. The `deliver` guide, the Roundfix `deliver` reference, the `archive-spec`
   skill, both skill versions, the Behavior Surface Record and **Spec
   Number** (no dependency: the text states the contracts of
   `_techspec.md`, and its Verification reads only documents).
2. `specNumber`, `queueHeaderLength`, `archiveCommitSubject`,
   `pullRequestTitle`, their use in `Archive` and `Publication`, and their
   tests (depends on: 1, which writes the guide and references for the
   surface this step changes, so the code and its text land in one chain).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. Spec 0291 has merged; this Spec delivers
after the active Specs 0292 and 0293 and the cycle AN Specs 0294 and 0295,
which also edit Roundfix skill references and versions, and it rebases on
them, changing only the `Archive` subject and the `Publication` title in
`internal/cli/deliver_workflow.go`.

## Risks & Considerations

- A slug whose leading digits are not a Spec number (`2026-10-…`) names
  `2026`, as the record subject already does; the body still names the slug.
- A Spec title is repository-authored text that reaches the forge as the
  pull request title; it is already public in the Archive Record, and the cut
  never adds text.
- A single-commit pull request takes GitHub's default summary from its
  commit, not its title; the queue's candidates carry several commits.
- A repository whose squash setting uses only the commit details is outside
  this Spec; the title still fits.

## Glossary

- adds: **Spec Number**

## Decisions

- Both messages keep the slug while it fits and name the Spec Number
  otherwise; see ADR-0309.
- One Spec Number helper serves the record subject, the archive subject and
  the title; see ADR-0302.
- The title reserves ` (#99999)` and counts UTF-16 code units, and "fits"
  for the title includes that reserve.
- The Spec title is read lazily from the Archive Record, so a short slug's
  publication is unchanged in bytes and in reads.
