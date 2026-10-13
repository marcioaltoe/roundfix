---
spec: 0266-the-review-record-survives-a-fix
status: active
created: 2026-10-13
surfaces: [backend, cli, docs]
---

# The review record survives a fix

On 2026-10-13 the Delivery Queue parked items 0259 and 0258 as
`corrective-spec-required` after their pre-PR review. The operator fixed each
finding, disposed it with `roundfix review dispose --fixed-by`, and ran
`roundfix deliver retry`; both retries were refused for a missing review
correction record, and both Pull Requests (#512 and #514) were opened by hand
(operator log entries 307, 308 and 312). Two defects compound, one per adopted
Backlog Entry:

- [the review correction proof misses its record under a symlinked root](references/2026-10-13-the-review-correction-proof-misses-its-record-under-a-symlinked-root.md):
  `roundfix review` names its record directory from the checkout root Git
  reports, the resolved path, while the retry's correction proof names it from
  the item worktree path the Delivery Queue recorded under the
  symbolic-linked `worktree.location`;
- [a second review round overwrites the parked findings record](references/2026-10-13-review-round-two-overwrites-the-parked-findings-record.md):
  the record is kept per checkout only, so round 2 at the fixed head replaces
  the record of the parked candidate's findings that the proof compares
  against.

Both fixes also changed code, and ADR-0233's retry rule refuses any change
outside the archived Spec's own records, so even a found record would not have
let either retry proceed: a characterization on `main` at v0.76.0 refused a
disposed code fix with `changed path review.txt is outside the named archived Specs`
although the archived Spec's Task declared `review.txt`.

No command, flag, output field or exit code changes. The review record
directory is found from either spelling of a checkout, the record of each
reviewed head stays readable after later rounds, and the retry's path rule
widens to the files the archived Spec's Tasks declared.

The maintainer decided on 2026-10-13 (cycle AG, final): the record writer and
the correction proof key by the same resolved path,
`config.CanonicalizeResolvedPath`, so a symbolic-linked `worktree.location`
works; a later review round keeps the earlier round's findings and
dispositions, with the design left to this Spec; every reader and writer of
the record is swept; and the fix is characterized end to end under a
temporary worktree root behind a symbolic link. Asked whether the retry should
also accept a code correction, the maintainer answered "Relaxar, com limites"
(AskUserQuestion, 2026-10-13, about 18:00 BRT): the retry also accepts a
changed path among the files the archived Spec's Tasks declared
(`creates:`, `interface:` and recorded paths, read at the candidate or the
record's source revision); review round 2 and the repository checks stay
mandatory; anything else still requires a corrective Spec, and the refusal
names the first path outside both sets.

## Prerequisites

None. Spec 0262 (merged, v0.74.0) exported `config.CanonicalizeResolvedPath`,
which this Spec reuses. The other cycle AG Specs (0267 to 0271) touch QA
evidence, ACP turn acceptance, test helpers, settlement wording and history
sanitize; none declares a file this Spec declares, and no Verification here
depends on their artifacts. The plan is authored against `main` at v0.76.0
and is rebased by the operator before its planning Pull Request.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes. A Reviewed Head Record has the existing Pre-PR Review Record
  schema; recorded repositories and ledger lines are read, never rewritten.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; tests and Verification use temporary homes, temporary
  repositories, a fake reviewer runner and the local `git` binary. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0272 (this Spec) decides the
  resolved checkout key and the Reviewed Head Record. ADR-0174: "Each checkout
  keeps its own review record and answer under the Artifact Directory, in a
  directory named from the checkout path", now from the resolved path.
  ADR-0233: the retry returns to review only when "every path changed from the
  candidate to the head lies under an archived Spec the blocker names", which
  ADR-0272 refines to also admit the archived Spec's declared paths. ADR-0165:
  "the correction is a new corrective Spec with its own `_authorization.md`
  and QA gate" for anything else, which stays. ADR-0261: the Archive Record's
  `source_revision` is recovered as every reader recovers it. ADR-0197:
  the ceiling "keeps the round-2 record in place", unchanged. ADR-0268: "Path
  comparison resolves the longest existing prefix of a path", the helper this
  Spec reuses. ADR-0259: "Every top-level test in a Parallel Test Package
  calls `t.Parallel()` as its first statement". ADR-0271: "A test bounds by
  real time only what it asserts". ADR-0256: "A release waits for the skills
  that describe a changed surface", answered by task_01's skills edit.
  ADR-0189 and ADR-0233 govern the owned skill's version raise. ADR-0244: "A
  Spec now carries a Glossary Declaration". ADR-0257: "Inside a Run the Daemon
  is the only full gate". ADR-0258: "A Spec is authored against the QA rerun
  classes". ADR-0265: "A Spec proves what it claims", answered by the
  TechSpec's Absence Evidence, ADR Alternatives and Coverage Cross-Read.
  ADR-0169: a record's `baseCommit` is the merge base, and a Reviewed Head
  Record keeps it as written. ADR-0196: a finding dismissed by validation
  stays outside the proof's disposition count, as today. ADR-0250 does not
  apply: no Baseline module version or coverage record changes.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths task_01 changes, the owned
  `roundfix` skill's review and deliver references and its `SKILL.md` pair,
  were measured with `speccheck.GovernedPath` through a `go test -overlay`
  probe that wrote nothing, and are granted by the cycle AG approval of
  2026-10-13 ("Pode iniciar", with "Concedo" for each Spec's declared Governed
  Paths and "considere autorizado a ajustar todas as skills se necessário").
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0266-the-review-record-survives-a-fix/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- The correction proof of a Delivery Retry finds the review record that
  `roundfix review` wrote for the item worktree, whichever spelling of the
  worktree path each side holds.
- A later review round never removes the record of a head an earlier round
  reviewed, so the proof can still read the parked candidate's findings.
- Every reader and writer of the Pre-PR Review Record and every repository
  comparison over it agree on one checkout identity.
- A disposed code correction confined to the files the archived Spec's Tasks
  declared returns the item to `reviewing`; any other path still parks it for
  a corrective Spec, named in the refusal.
- The review and deliver guides and the Roundfix skill say where each record
  lives, what the retry reads and which paths it admits.

## User Stories

1. As the operator of a Delivery Queue under a symbolic-linked worktree
   location, I want a retry after a disposed fix to return the item to
   `reviewing`, so that I stop opening Pull Requests by hand.
2. As the operator, I want to run review round 2 on a parked item before the
   retry without losing the round-1 findings and dispositions the retry
   proves, so that checking my fix never costs me the delivery.
3. As the operator, I want a code fix to a file the archived Spec's Tasks
   created or changed to return the item to review, and a fix anywhere else
   to be refused naming that file, so that only an unplanned change needs a
   corrective Spec.
4. As an agent reading the Roundfix skill, I want the review reference to
   name the record of each reviewed head, so that I know what a later round
   leaves in place.

## Core Features

1. **Resolved checkout key.** The record directory of a checkout is named
   from its resolved path, so the review command, `review dispose` and the
   correction proof meet in one directory for every spelling.
2. **Resolved comparisons.** Every comparison of a record's or a
   disposition's repository with a checkout path or another record compares
   resolved paths.
3. **Reviewed Head Record.** A review that ends in a verdict for a head, and
   is not a reuse, also keeps its record under that head; later rounds never
   replace it, and the proof reads it for the parked candidate.
4. **Declared-path correction.** The retry also admits a changed path that
   the archived Spec's Task files declared under `creates:` or `interface:`
   or recorded under `## Recorded paths`, read at the recorded revision, and
   its refusal names the first path outside both sets.
5. **Documented.** The review and deliver user guides and the Roundfix
   skill's review and deliver references describe the resolved key, the
   Reviewed Head Record, what the retry reads and which paths it admits.
6. **Glossary.** `CONTEXT.md` adds **Reviewed Head Record** and revises
   **Pre-PR Review Record**.

## User Experience

A repository whose worktree location is a real directory sees no change in
any command's output. Under a symbolic-linked location, a retry after a
disposed fix that changes only the archived Spec's records returns the item to
`reviewing` instead of refusing with
`read review correction record: open …: no such file or directory`. A retry
after a code fix to a file the archived Spec's Tasks declared also returns to
`reviewing`; one that changes any other file is refused with
`changed path <path> is outside the named archived Specs and the paths their Tasks declared`. Each checkout's
review directory gains `heads/<head commit>.json` files beside
`pre-pr-review.json`.

## Non-Goals / Out of Scope

- Accepting a correction that changes a path neither under the archived
  Specs nor among their declared paths; that still needs a corrective Spec
  (ADR-0165).
- Skipping review round 2, the repository gate or the required checks for an
  accepted correction.
- A flag to dispose a finding of an earlier round; `review dispose` keeps
  reading the current record.
- Rewriting recorded repositories, ledger lines or Delivery Queue items, or
  migrating the Run Database.
- Changing the Artifact Directory identity of a main checkout reached through
  a symbolic link; an item worktree's Artifact Directory already resolves
  through its Git common directory.
- Reclaiming Reviewed Head Records in `gc`.

## Success Metrics

1. Success Metric: in a temporary home with a worktree root behind a symbolic
   link, `roundfix review` with findings, a docs-only fix under the archived
   Spec, `roundfix review dispose --fixed-by` and the correction proof called
   with the link spelling of the item worktree end with the proof accepted,
   for an explicit and for the default Artifact Directory.
2. Success Metric: after round 2 at the fixed head in the same setup, the
   Reviewed Head Record of the candidate still holds the round-1 findings, and
   the proof still accepts.
3. Success Metric: no repository comparison over a review record or a
   disposition outside the resolved-path helper remains in the review and
   correction-proof sources.
4. Success Metric: a disposed code fix to a file the archived Spec's Task
   declared is accepted, for a Spec archived as a record and as a folder, and
   a fix that also changes an undeclared file is refused naming that file.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Published literature, Stack Overflow answer to "Why does git rev-parse
  --show-toplevel return different casing than actual path"
  (https://stackoverflow.com/questions/53396610, found and read through Exa):
  "When Git resolves a path, it canonicalizes the path. On Unix systems, that
  means it resolves all symlinks". That is why `roundfix review`, which keys
  by the Git root, holds the resolved path while the Delivery Queue holds the
  link spelling. The same need appears in git-for-windows issue 2114
  (https://github.com/git-for-windows/git/issues/2114), where Git LFS broke
  because a path relativized against the Git root did not "share the proper
  components in common". The installed Git 2.56.0, measured offline during
  authoring, reports a linked worktree added through a symbolic-linked
  directory under its resolved directory; the QA gate repeats it.
- The operator's queue log (entries 307, 308 and 312 of 2026-10-13) records
  the two refused retries this Spec answers; it lives on the operator's
  machine, so it is context, not the row's evidence.
- The Secondbrain, read during authoring through `wiki/index.md` and two
  `qmd query` searches: the mirrors of ADR-0174
  (`projects/roundfix/mirror/docs/history/specs/0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record.md`),
  ADR-0233 and ADR-0268 confirmed that the per-checkout record, the
  review-only correction rule and the prefix-resolving comparison are settled
  decisions this fix must keep; no page covered round history.

## Glossary

- adds: **Reviewed Head Record**
- changes: **Pre-PR Review Record**
- not a term: **Resolved checkout key** — a Core Feature heading
- not a term: **Resolved comparisons** — a Core Feature heading
- not a term: **Declared-path correction** — a Core Feature heading

## Decisions

- The record directory and every repository comparison use the resolved
  checkout path; see ADR-0272.
- Each reviewed head keeps its record, and the proof reads it; see ADR-0272.
- The retry also admits the archived Spec's declared paths, refining
  ADR-0233; see ADR-0272.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the sweep, the contracts and the
build order.
