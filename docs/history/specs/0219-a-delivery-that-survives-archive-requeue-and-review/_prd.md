---
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: archived
created: 2026-10-03
surfaces: [backend, cli, docs]
archived: "2026-10-03"
source_slug: 0219-a-delivery-that-survives-archive-requeue-and-review
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings; the only blocked row is the Pull Request row (R10)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 6d7bb6f1722526228a978b6133c72af9b61d6a1a
---


# A delivery that survives archive, requeue and review

Between 2026-10-01 and 2026-10-03 the operator finished five Delivery Queue
items by hand although their work was correct, and each stop came after the
Spec's own work was done. Spec 0205's test read its own active TechSpec and
failed the repository gate once the archive commit moved it. The Pull Request
then failed the corpus check, because Spec 0218's Task Context still named
Spec 0205's active path. The Delivery Retry after the operator's correction
refused because the archived item head had moved, so the operator opened the
Pull Request by hand. Spec 0205, and on 2026-10-03 Spec 0217, needed a new
queue, and `deliver start` created a fresh item from the default branch that
would have run every completed Task again, so the operator merged the old
item branch into it. The pre-PR review of Specs 0215 and 0218 parked each item
on a finding that the Spec was archived with a failed QA Task, although the
archive was an authorized QA Archive Override. And a Verification command
truncated by escaped backticks was reported honest by the authoring probe,
then failed a complete Task twice. This Spec closes each of those traps: a
Spec that archives breaks no other file, a retry and a requeue keep the item's
work, the review knows an authorized override, and a command the shell cannot
parse is caught when it is written.

## Prerequisites

None. This Spec raises the Roundfix Skill's version in each implementation
Task. Specs authored in the same cycle may also raise it; whichever merges
second takes the next version at its own implementation time, because each
Task raises the version found on its starting tree.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Specs keep their
  slugs, Delivery Queue items their Spec slug and item branch, and the two new
  Spec Consistency findings take codes in the existing `SC-` family. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the archive, the Spec Consistency
  Check, the retry, the queue start and the review validator read only local
  files, local Git and the review's configured runtime; no credential and no
  network call is added, and no test reaches GitHub. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0223 (this Spec) sets the five
  delivery rules, ADR-0223: "Roundfix never rewrites code to follow an archive".
  ADR-0226 (this Spec) decides the malformed verdict,
  ADR-0226: "A command the shell rejects is reported as malformed".
  ADR-0148 keeps one prober for the Daemon and authoring, ADR-0148: "the vacuous/failing/unknown classification lives in one extracted prober".
  ADR-0182 returns a new refusing finding to the Task that caused it,
  ADR-0182: "the refusing Spec Consistency findings that were".
  ADR-0154 keeps an override a record of user authority,
  ADR-0154: "It is not a successful Run or an authorization to".
  ADR-0165 keeps a blocking review after archive parked for a corrective Spec,
  ADR-0165: "Publication parks as".
  ADR-0196 keeps the review's validation, ADR-0196: "Every other finding stands, and only a standing finding parks".
  ADR-0197's two review rounds, ADR-0170, ADR-0192, ADR-0193, ADR-0199,
  ADR-0211 and ADR-0120's single history root hold unchanged. ADR-0187 and
  ADR-0189 govern the Roundfix Skill edits, and ADR-0184 has the TechSpec
  state the changed command surface as a transcript,
  ADR-0184: "A TechSpec now declares numbered Surface Transcripts". The gate
  is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167. ADR-0097's carry conditions, ADR-0194's recorded observations,
  ADR-0195's always-observed rows and ADR-0210's per-input digest hold for
  this Spec's own gate, while ADR-0096's mechanical stage keeps its role.
  ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt, and ADR-0178 and ADR-0179 decide the
  grant each Task commit runs under. ADR-0208 groups the three adopted
  sources, which share the delivery's archive and review context; ADR-0209's
  grouping suggestion was skipped at the Jev ceiling and changes nothing; and
  ADR-0169's merge-base diff holds for the review. ADR-0212 does not apply,
  because reconcile is unchanged, and ADR-0220 does not apply, because the
  Doctor Command is unchanged. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill under `.agents/skills/`
  and its `SKILL.md` mirror are Governed Paths, and the maintainer authorized
  skill edits ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A Spec's archive breaks no test, fixture or other active Spec: a file that
  would break is refused while the Task that wrote it can still fix it, and
  another Spec's Task Context keeps resolving.
- A Delivery Retry after a correction committed on top of the archive, and a
  new queue for a Spec whose item branch holds completed work, continue that
  work with no hand-made Pull Request and no hand merge.
- The pre-PR review does not park an item on an archive the maintainer
  authorized through the QA Archive Override.
- A Verification command the shell cannot parse is reported as malformed
  before any Run spends an attempt on it, and the authoring probe says which
  of the commands it ran come from uncommitted Spec files.

## User Stories

1. As a Spec author, I want the Spec Consistency Check to refuse a test or
   fixture that names my Spec's active directory, so that the Task that wrote
   it fixes it before the archive moves the directory.
2. As the operator, I want another active Spec's Task Context that names an
   archived Spec's former path to keep resolving, so that a delivery's Pull
   Request does not fail the corpus check on a Spec it did not touch.
3. As the operator, I want a Delivery Retry after I commit a correction on top
   of the archive to send the item back to review, so that I never open the
   Pull Request by hand.
4. As the operator, I want `roundfix deliver start` for a Spec whose item
   branch already holds completed Tasks to continue that branch, so that no
   proved Task runs again and I never merge branches by hand.
5. As the maintainer, I want the pre-PR review to know that an archive
   carrying my recorded QA Archive Override is authorized, so that it stops
   parking items on a decision I already made.
6. As a Spec author, I want a Verification command that the shell cannot
   parse to be reported as malformed by the Spec Consistency Check and by the
   probe, so that it never passes as honest and never fails a complete Task.

## Core Features

1. **A file that pins an active Spec's directory is refused.** For each Spec
   it checks, the Spec Consistency Check reports an error for every tracked or
   untracked, non-ignored file outside the Spec Root, the archive root and the
   history root, other than Markdown, that names the Spec's active directory,
   with each file and line. The Daemon's Settlement Check returns the finding
   to the Task that wrote the file (ADR-0182), and a file that existed before
   the Task started does not fail it.
2. **The archive refuses on such a file.** The Archive Command refuses to
   archive a Spec while such a file exists, names each file and line, and
   moves nothing.
3. **A Task Context path follows the archive.** A Task Context entry under an
   active Spec's directory that no longer exists resolves to the same relative
   path under the archive root when that path exists there, and is reported
   unresolved otherwise. No other Spec's file is rewritten.
4. **A retry after a post-archive correction resumes at review.** A Delivery
   Retry of an archived item whose head descends from its last candidate head
   records the head as a new candidate and returns the item to `reviewing`.
   The review, archive and gate stages then run as for any reviewed candidate,
   and the archive stage finds the Spec already archived. A head that does not
   descend still refuses with its reason, a `corrective-spec-required` park
   still refuses any moved head, and an operator-archived
   `qa-environment-partial` item still needs its override record.
5. **A new queue continues the existing item branch.** When `roundfix deliver
   start` queues a Spec that has exactly one local item branch with commits
   the default branch lacks, the new item continues that branch and its
   worktree, prints which branch it continues, and starts at the stage a new
   item starts at. With two or more such branches the command refuses before
   recording a queue, names every branch, and says how to keep one. A Spec
   with no such branch gets a new item branch as today.
6. **An authorized override archive is a Delivery Convention.** The review
   prompt carries a fifth Delivery Convention: an archived `_prd.md` that
   records `qa_override: true` with an approval source and a reason is an
   authorized QA Archive Override, and its failed QA Task and its QA verdict
   are recorded as observed by design. The validator may dismiss a finding
   anchored in that archived Spec that only restates it. An archive without
   that record is not covered, and every other finding still stands.
7. **A command the shell cannot parse is malformed.** The one Verification
   prober asks the shell to parse each command without running it. The probe
   of `roundfix spec check --run-verification` reports a command the shell
   rejects with the verdict `malformed`, its parse error, and exit `1`; the
   Daemon refuses the Task before the Agent starts. The Spec Consistency Check
   also reports, without running anything, a Verification bullet whose code
   span ends in a backslash or is followed by another backtick.
8. **The probe names uncommitted sources.** When `roundfix spec check
   --run-verification` runs commands whose Task file or Task Graph is
   untracked or differs from `HEAD`, its output names each such file, and the
   JSON report lists them. Nothing is refused for it.
9. **The skill and the guides say so.** The Roundfix Skill's `archive`, `spec`,
   `deliver` and `review` references and the matching command guides describe
   each change.

## User Experience

The operator sees fewer parks and no new command or flag. A refused archive
prints each file and line that names the Spec's active directory. `roundfix
deliver start` prints one line for each item that continues an existing
branch, or refuses with every branch it found. A retry after a post-archive
correction reports the stage `reviewing`. `roundfix spec check
--run-verification` prints `malformed` beside a command the shell rejects and
a line naming each uncommitted source of the commands it ran.

## Non-Goals / Out of Scope

- Rewriting tests, code or another Spec's files when a Spec archives.
- Markdown outside the Spec roots that names an active Spec's path, and the
  relative links inside a Spec that the archive's move breaks; a pending Inbox
  Entry from another repository reports the second and awaits triage.
- Merging the default branch into a continued item branch, or continuing an
  item branch whose Spec is already archived on it.
- A correction for a `corrective-spec-required` park, which stays a
  corrective Spec's work (ADR-0165).
- Refusing an uncommitted Verification source, or an approval record that
  makes one executable.
- The queue owner's binary, which another Spec of this cycle owns.
- Editing `CONTEXT.md`; the glossary check is the QA gate's.

## Success Metrics

1. Success Metric: a Spec whose test names the Spec's active directory reports
   the pin finding, a Task that adds such a test fails its Settlement Check,
   and the Archive Command refuses that Spec and leaves its directory in
   place; before this Spec all three pass.
2. Success Metric: another active Spec whose Task Context names an archived
   Spec's former path reports no `SC-REF-UNRESOLVED`, and a path missing in
   both places still does.
3. Success Metric: a retry of an archived `gate-failed` item whose head is one
   correction commit past its candidate returns the item to `reviewing` with
   that head as its newest candidate; before this Spec it refuses with
   `differs from candidate head`.
4. Success Metric: `deliver start` for a Spec with one item branch two commits
   ahead of the default branch records that branch on the new item, and with
   two such branches it refuses, names both, and records no queue.
5. Success Metric: a review finding anchored in an archived Spec whose
   `_prd.md` records the override is eligible for the fifth convention, and
   one anchored in an archived Spec without the record is not.
6. Success Metric: the command recorded in Spec 0205's delivery is reported
   `malformed` by the probe and by the Spec Consistency Check, where before
   this Spec the probe reports it `honest`.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The operator's intervention log, read-only on 2026-10-03, entries 70, 74,
  76 to 78, 108, 118, 124 and 127, and the three Backlog Entries this Spec
  adopts, which record each stop and its hand answer.
- GitHub's record of Pull Request #329, the Spec 0205 delivery opened by hand
  after the retry refused on a moved archived head.
- The CommonMark specification, code spans
  (<https://spec.commonmark.org/0.31.2/>): "Backslash escapes do not work in
  code blocks, code spans, autolinks, or raw HTML".
- The POSIX shell's `set -n`
  (<https://pubs.opengroup.org/onlinepubs/9799919799/utilities/V3_chap02.html>):
  "The shell shall read commands but does not execute them; this can be used
  to check for shell script syntax errors."
- Git's `git grep` manual (<https://git-scm.com/docs/git-grep>), `--untracked`:
  "In addition to searching in the tracked files in the working tree, search
  also in untracked files."
- A reproduction during authoring on 2026-10-03: the truncated command from
  Spec 0205's delivery exits `2` under `sh -c` with `unexpected EOF while
  looking for matching`, and `sh -n -c` rejects it with the same status
  without running anything.

## Decisions

- An archive refuses rather than rewrites a file that pins the Spec's active
  directory, and another Spec's Task Context resolves through the archive
  instead of being rewritten, so the archive commit stays the Spec's exact
  move. See ADR-0223.
- A post-archive correction is reviewed again before publication rather than
  published on the strength of the earlier review. See ADR-0223.
- A requeue continues the one item branch that holds work and refuses when
  there are several, rather than choosing one. See ADR-0223.
- An authorized override archive becomes a Delivery Convention rather than a
  deterministic dismissal, so the validator's fail-closed rule still applies.
  See ADR-0223.
- Uncommitted Verification sources are named, not refused, and no approval
  record is added, because authoring runs the probe on uncommitted Specs by
  design. See ADR-0226.
- Re-measured against `main` at `305211a3`: Spec 0211 already resumes an
  operator-archived `qa-environment-partial` retry from a descendant head
  (entry 107), so that part of the archived-retry stop is not repeated here.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
