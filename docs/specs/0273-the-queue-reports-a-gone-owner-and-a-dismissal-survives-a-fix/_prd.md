---
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
status: active
created: 2026-10-16
surfaces: [backend, cli, docs]
---

# The queue reports a gone owner, and a dismissal survives a fix

Three small frictions of the Delivery Queue cost the operator time in cycle
AG. After a reboot on 2026-10-13 the Queue Owner was gone, yet
`roundfix deliver status` showed the item as `checking` with no warning for
about eleven hours
([the Backlog Entry](references/2026-10-14-deliver-status-does-not-say-the-queue-owner-is-gone.md)).
On 2026-10-14 a dismissal recorded after a committed fix was refused because
HEAD had moved, so the operator reset the branch to dismiss and reset forward
again; and `roundfix deliver retry` run inside the item worktree said that no
Delivery Queue existed instead of naming the checkout that owns it
([the Backlog Entry](references/2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md)).
After this Spec, status names a recorded Queue Owner that is proven not
running and the command that resumes it, a dismissal is accepted after a fix
that leaves the finding's anchor unchanged, and a `deliver` command run in an
item worktree names the checkout to run it from.

## Prerequisites

None. Specs 0272, 0274 and 0275 of cycle AH are authored at the same time and
touch other files; no Verification here depends on their artifacts. This Spec
is first in the cycle's order (B → A → C → D).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, Run Database
  table or column is added; the item worktree lookup reads the `worktree` the
  queue already records, and the new ledger field holds a commit name Git
  already reports. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; Git and the Run Database are read locally and tests run on
  temporary repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0279 (this Spec) decides the three
  behaviors: "Status stays read-only and never reclaims the owner record".
  ADR-0044: "anything short of proof keeps the block", so a live
  PID whose identity cannot be read is not reported gone. ADR-0266: "a
  mismatch sends nothing and counts as proof that the recorded process ended"
  is the same start-time proof, and status sends no signal beyond the
  existence probe the owner commands already use. ADR-0174: "The disposition ledger stays shared, because
  every disposition already names its repository and reviewed head", which the
  dismissal after a fix keeps. ADR-0196: "Each finding must open with a
  `path:line` anchor", the anchor the dismissal proves unchanged.
  ADR-0272: "The current record keeps its role", and its rejected
  "`review dispose --head` flag" stays rejected. ADR-0268: "Path comparison
  resolves the longest existing prefix of a path", which the item worktree
  lookup uses. ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0244: "A Spec declares the domain terms it
  introduces".
  ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts"; no new test
  waits on a clock. ADR-0280: "A guide Task records its own surface
  fingerprint", which each guide Task here does, and "a PRD records who
  decided", which the Decisions section does. ADR-0281: "A test replaces an
  executable it rebuilds instead of writing into it"; not applicable, because
  no test here rebuilds an executable.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados", under the cycle AE
  "Concedo" for the Governed Paths each Spec declares) and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"); bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.
  No Makefile, workflow, lint, formatter, module manifest or Baseline module
  change.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix/_authorization.md`.

## Goals

- An operator who reads `deliver status` learns, from that output alone, that
  the recorded Queue Owner is gone and that `roundfix deliver resume` recovers
  the queue, without a false report when the PID was reused.
- A finding can be dismissed after another finding's fix is committed, without
  resetting the branch, as long as the dismissed finding's anchored lines are
  unchanged.
- A `deliver` command run in an item worktree says where to run it instead of
  saying that no queue exists.

## User Stories

1. As an operator returning to a machine that rebooted, I want `deliver
   status` to say that the Queue Owner is not running and to name `roundfix
   deliver resume`, so that a green Pull Request does not wait for hours behind
   a dead process.
2. As an operator disposing a parked item's findings, I want to dismiss one
   finding after committing the fix of another, so that I do not reset the
   branch to the reviewed head and forward again.
3. As an operator working inside an item worktree, I want `deliver retry` and
   the other queue commands to name the checkout that owns the Delivery Queue,
   so that I know where to run them.

## Core Features

1. Owner report. `deliver status` prints one `Owner:` line when the queue
   records an owner, at least one item can still advance, and that owner is
   proven not running: its PID is absent, or alive under another start-time
   identity. The line names the PID, the state and `roundfix deliver resume`.
   A proven live owner, a queue without a recorded owner and a queue whose
   items are all merged or parked add no line. An owner identity that cannot
   be read adds no line and writes one diagnostic to stderr. Status never
   reclaims the owner record, and its exit code and every other line stay.
2. Dismissal after a fix. `review dispose <id> --dismiss` is accepted when the
   current HEAD equals the reviewed head, as today, or descends from it while
   the finding's recorded anchor is unchanged at the same line numbers. The
   disposition stays tied to the reviewed head and records the current HEAD in
   a new optional field. Otherwise it refuses with a reason that says whether
   the head does not descend, the finding records no anchor, or the anchor
   changed. Every other dispose rule stays.
3. Owning checkout. `deliver status`, `resume`, `retry` and `stop`, run where
   no Delivery Queue exists for the checkout, add to the refusal the item and
   the checkout whose queue records this checkout as its item worktree, and
   say to run the command there. Elsewhere the refusal keeps its bytes.
4. Documentation. The `deliver` and `review` user guides, the Roundfix skill's
   `deliver` and `review` references and `CONTEXT.md` describe the three
   behaviors.

## User Experience

The `Owner:` line follows the item rows and any `Warning:` lines and precedes
the `Park:` lines, so it is read before the queue's limits and usage. It
reuses the words `deliver resume` and `deliver retry` already print when they
reclaim an owner record. The refusals keep their exit code `2`, their stderr
placement and their first words, so a reader or a script that matched the old
text still matches it.

## Non-Goals / Out of Scope

- Resuming a queue automatically at login or through launchd.
- Reporting a queue with no recorded owner, or the window inside `deliver
  start` and `deliver retry` before the new owner claims the queue.
- A JSON form of `deliver status`; it has none today and gains none.
- Resolving the owning queue from an item worktree and acting on it.
- Changing the review correction proof, the Reviewer Lineage, review reuse or
  the shape of a `fixed` disposition.
- Dismissing a finding whose anchored lines a later commit changed or moved.

## Success Metrics

1. Success Metric: through the built `bin/roundfix`, a queue whose recorded
   owner PID cannot exist prints the `Owner:` line naming `roundfix deliver
   resume`, a live PID recorded with a different identity prints the
   different-identity line, and a live PID recorded with its own identity
   prints none.
2. Success Metric: through the built `bin/roundfix`, a dismissal after a
   committed fix that leaves the anchor unchanged is appended once with the
   reviewed head and the current head, and a fix that changes or moves the
   anchored lines is refused with the ledger unchanged.
3. Success Metric: through the built `bin/roundfix`, `deliver retry` and
   `deliver status` run in an item worktree name the owning checkout, and run
   in an unrelated repository keep the delivery base's refusal byte for byte.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa on 2026-10-16:

- psutil issue 2400, comment by its maintainer on 2024-04-10
  (https://github.com/giampaolo/psutil/issues/2400): "I think the `pid +
  create_time` pair strategy used by psutil is the best we can do." The
  r-lib/ps documentation (https://github.com/r-lib/ps/blob/main/R/low-level.R)
  says the same: "The pid and the start time pair serves as the identifier of
  the process, as process ids might be reused". This is why status reports a
  live PID with another start-time identity as gone and a live PID with the
  recorded identity as running.
- GitHub Changelog, "Precision review comment outdating", 2018-09-25
  (https://github.blog/changelog/2018-09-24-precision-review-comment-outdating/):
  "the exact line of code that anchors the comment must change to outdate the
  conversation." The GitHub REST reference for pull request review comments
  (https://docs.github.com/en/rest/pulls/comments) adds: "Not using the latest
  commit SHA may render your comment outdated if a subsequent commit modifies
  the line you specify as the position." A review comment therefore stays
  current across later commits while its anchored line is unchanged, which is
  the rule the dismissal after a fix adopts.
- The git-merge-base manual (https://git-scm.com/docs/git-merge-base) for
  `--is-ancestor`: "exit with status 0 if true, or with status 1 if not.
  Errors are signaled by a non-zero status that is not 1." The dismissal reads
  status 1 as a head that does not descend and any other failure as an error.

## Validation

The Secondbrain was consulted through `wiki/index.md` and two queries,
`qmd query "roundfix deliver status queue owner died reboot resume"` and
`qmd query "roundfix review dispose dismiss after fix reviewed head ancestor anchor"`.
They returned the mirrors of both adopted Backlog Entries, the `deliver` user
guide, the mirror of the decision that operational Runs detach from their
caller (`projects/roundfix/mirror/docs/adr/0028-operational-runs-can-detach-from-their-caller.md`),
ADR-0196, the mirror of the reviewer lineage decision (`projects/roundfix/mirror/docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`), the mirror of the review-only correction decision
(`projects/roundfix/mirror/docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`)
and the Archive Record of Spec 0187
(`projects/roundfix/mirror/docs/history/specs/0187-a-queue-that-recovers-without-a-supervisor.md`),
which added the owner-older-than-main warning but no liveness report; and
Spec 0266 (`projects/roundfix/mirror/docs/history/specs/0266-the-review-record-survives-a-fix.md`),
which kept the reviewed head's record after a fix but left the dismissal rule
as it was. They confirmed that no earlier Spec decided either behavior and
that the start-time proof already exists for `resume` and `retry`. The pending
inbox of the repository is empty. Exa returned the sources quoted above; they
are standard references, not new research, so no Inbox Entry was captured.

The current behavior was characterized on 2026-10-16 at `9e723b82` with
`go test -overlay` probes in `internal/cli` that wrote nothing to the
repository: a queue whose item is `checking` and whose recorded owner PID
does not exist printed only the item row, `Limits:` and `Usage:`; `deliver
retry` and `deliver status` inside a recorded item worktree refused with
`Delivery Queue for repository "<item worktree>" does not exist`; and a
dismissal after a commit that left the anchored file untouched refused with
`dismissal requires current HEAD "<head>" to equal reviewed head "<reviewed>"`.

## Glossary

- changes: **Queue Owner**
- changes: **Review Finding Disposition**

## Decisions

- 2026-10-16 — the maintainer, through AskUserQuestion, chose cycle AH's scope
  and order (B first) and gave this Spec's scope: "deliver status avisa quando
  o dono da fila morreu e indica deliver resume; review dispose aceita dismiss
  depois de um fix commitado (head revisado é ancestral, âncora intacta); a
  recusa do retry nomeia o checkout principal."
- Autonomous route: no person present; decisions taken from the maintainer's
  scope above, the adopted Backlog Entries and the characterization of
  2026-10-16.
- Status proves the recorded owner as `resume` and `retry` do and never
  reclaims it; see ADR-0279.
- A dismissal after a fix requires an unchanged anchor at the same line
  numbers and records the current head; see ADR-0279.
- An item worktree refusal names the owning checkout instead of acting on its
  queue; see ADR-0279.
- `deliver status` gains no JSON form, and the command help keeps its text.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the characterization, the contracts
and the build order.
