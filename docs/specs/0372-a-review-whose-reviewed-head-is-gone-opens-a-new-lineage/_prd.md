---
spec: 0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage
status: active
created: 2026-11-10
surfaces: [cli]
type: fix
---

# A review whose reviewed head is gone opens a new lineage

`roundfix review` blocks with exit 2, before any reviewer call, when the commit an earlier review of the same checkout recorded as its reviewed head no longer exists in the repository. On 2026-11-10 Conexus hit it on Roundfix 0.137.7 after the usual squash merge: the reviewed branch was deleted, Git pruned its commits, and every later review of that checkout stopped at `check review lineage ancestry: git merge-base --is-ancestor <prior head> <head>: fatal: Not a valid commit name <prior head>: exit status 128`, with no supported way out. After this Spec a recorded reviewed head or merge base that a complete repository proves absent opens a new Reviewer Lineage at round 1 with its own restart reason naming the missing commit, `review dispose` and the Delivery Queue's correction proof refuse by name instead of failing on Git, and every other Git failure, including an absence a shallow or partial clone cannot prove, still blocks.

## Sources

One source, adopted (see `references/_index.md`):

1. [`references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md`](references/2026-11-10-conexus-review-blocks-when-the-reviewed-head-is-gone.md),
   the ai-memory message `01a1c53f-c3be-78a0-851d-cda2f4dd43a8` of
   2026-11-10, reported by conexus, nature `measured`, with the operator's
   note on `decideReviewLineage`.

At the authoring base `docs/backlog/` and `docs/findings/` hold no entry, so
no other source shares this context.

## Prerequisites

None in `requires`. Authored against `9fda0b65` (v0.137.7). The queue
delivers this Spec first, then Specs 0371 and 0373 of the same cycle. This
Spec owns `internal/cli/review_lineage.go`, the ancestry checks in
`internal/cli/review.go` and `internal/cli/deliver_review_correction.go`, and
the `review` reference. It shares with them only the owned Roundfix skill
version fields (`.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
`skills/testdata/owned-skill-versions.json`) and
`docs/references/behavior-surfaces.json`, which task_01 regenerates with the
repository's record commands and never by hand; it also edits the `deliver`
reference and guide, which neither 0371 nor 0373 changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted or changed; review records and the disposition ledger keep their
  fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change reads the local Git
  repository only and never fetches; it refuses to probe a partial clone
  because the probe would fetch. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0197 starts a new lineage when the
  head does not descend from the reviewed head, and a reviewed head a
  complete repository no longer holds cannot be an ancestor of the head;
  ADR-0272 owns the Delivery Retry's correction proof and ADR-0279 the
  dismissal anchor proof, which gain a named refusal. ADR-0233, ADR-0286,
  ADR-0366 and ADR-0373 are not applicable: they cite those ADRs, and the
  retry path rule, the Correction Review, the merge-base rules and the owner
  liveness proof they decide are unchanged. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized: no tooling configuration, script, ignore file, plugin
  declaration or version pin changes. The Governed Paths task_01 touches are the owned Roundfix skill
  files and the canonical `review` and `deliver` references, measured with
  `speccheck.GovernedPath`, which the maintainer's cycle grants cover.
  Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0372-a-review-whose-reviewed-head-is-gone-opens-a-new-lineage/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A review of a checkout whose recorded reviewed head or merge base the
  repository no longer holds runs as round 1 of a new Reviewer Lineage, and
  `deliver status` and `deliver retry` name that reason.
- Absence is proven, never inferred from a Git failure: only a complete
  repository's unambiguous answer counts, so the two-round ceiling cannot be
  reset by a transient or unrelated failure.
- `review dispose` and the correction proof of a Delivery Retry name a
  reviewed head, parked candidate or fixing commit that is gone, instead of
  printing a raw Git failure.

## User Stories

1. As a developer whose reviewed branch was squash-merged, deleted and
   pruned, I want the next `roundfix review` of my checkout to start a new
   lineage, so that the checkout is not blocked forever.
2. As the maintainer, I want only a proven absence to restart the lineage, so
   that a failing Git command, a shallow clone or a partial clone never resets
   the two-round ceiling.
3. As an operator disposing a finding or retrying a parked item, I want a
   refusal that names the missing commit and what to do next.

## Core Features

1. `decideReviewLineage` restarts with `the reviewed head <prior head> is no
   longer in the repository` or `the recorded merge base <prior base> is no
   longer in the repository` when the ancestry check fails and the repository
   proves the recorded commit absent; `deliver status` and `deliver retry`
   print it through the existing notice.
2. One proof of absence: the recorded value is a full object ID, the
   repository is neither shallow nor partial, and `git cat-file -e <id>` exits
   1 with no output. A shallow or partial clone blocks with a named reason;
   every other failure keeps today's error.
3. `review dispose` (`--dismiss` and `--fixed-by`) refuses a reviewed head
   that is gone with a named reason; the correction proof refuses a parked
   candidate that is gone and treats a fixing commit that is gone as an
   invalid disposition.
4. The `review` and `deliver` guides and Roundfix skill references list the
   two new restart reasons and the refusals.

## User Experience

Before this Spec, in the Conexus shape, `roundfix review --base origin/main`
prints `roundfix: review blocked: check review lineage ancestry: git
merge-base --is-ancestor <prior head> <head>: fatal: Not a valid commit name
<prior head>: exit status 128` and exits 2 at every run. After it, the review
runs round 1 of a new lineage; `deliver status` of an item in that state
prints `Review: <slug> head <head> opens a new Reviewer Lineage at round 1, a
full review of the candidate: the reviewed head <prior head> is no longer in
the repository`. `roundfix review dispose F1 --dismiss --evidence <text>`
against that record prints `roundfix: review dispose refused: reviewed head
"<prior head>" is no longer in the repository; run roundfix review to open a
new Reviewer Lineage` and exits 2. In a shallow or partial clone the review
still exits 2, with `check review lineage ancestry: a shallow or partial
clone cannot prove that the recorded commit <commit> is gone; fetch the full
history and try again`.

## Non-Goals / Out of Scope

- The detached-HEAD refusal (`detect current branch: HEAD is detached`,
  `preflight.InspectGit`) that Conexus also met in a worktree. It is not a
  defect by the `review` reference's own words: the reference reviews the
  current candidate of a branch and never promises a detached HEAD, and the
  refusal belongs to the Git State every command shares through
  `InspectGit` (17 production callers at `9fda0b65`). Changing it is a product decision for another
  Spec.
- Fetching a missing commit back (from the remote, `refs/pull/<n>/head` or a
  promisor remote): the review record stores no Pull Request number, the
  review never reaches the network, and an absent reviewed head is a head
  that does not descend, which already restarts.
- A `roundfix review` option to reset the lineage by hand.
- Any change to the restart rules of ADR-0366, the Correction Review, the
  ceiling, the record and ledger formats, or the `--fixed-by` refusal texts
  for a reviewed head that is present.
- Ancestry checks that read no review record or ledger (release plan,
  reconcile, archive, worktree and branch integrity).
- Any change to the Makefile, workflows, `go.mod`, `go.sum`, lint or
  formatter configuration, the Baseline or `CONTEXT.md`.

## Success Metrics

1. Success Metric: in a temporary repository where a first review recorded a
   head whose branch was then deleted, its reflog expired and its objects
   pruned, the second `roundfix review` exits 0 at round 1 with an empty
   stderr; on the authoring base it exits 2 with
   `check review lineage ancestry`.
2. Success Metric: a shallow clone and a partial clone each still block with
   the named incomplete-history reason, and a broken alternate object store,
   an abbreviated recorded commit and a Git failure on a present commit still
   block with today's error.
3. Success Metric: `review dispose --dismiss`, `review dispose --fixed-by` and
   the correction proof each refuse a gone commit by name instead of printing
   or returning a Git failure.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Conexus message, measured on a repository this Spec did not
  build (`useconexus/conexus`, 2026-11-10, Roundfix 0.137.7): `roundfix
  review` "saiu com exit 2 antes de qualquer chamada ao reviewer" with
  `fatal: Not a valid commit name f96765ecdf99eb392e3225a5febf455d60e12d92:
  exit status 128`, and both recorded commits were absent ("nenhum dos dois
  presente no clone (git cat-file falha)").
- The Git manual, read through Exa on 2026-11-10: `git cat-file -e` "Exit
  with zero status if <object> exists and is a valid object. If <object> is
  of an invalid format, exit with non-zero status and emit an error on
  stderr." (https://git-scm.com/docs/git-cat-file); `git gc` removes
  "unreachable objects" (https://git-scm.com/docs/git-gc), so a commit an
  existing commit reaches is never pruned; and partial clone "requires that
  the user be online and the origin remote or other promisor remotes be
  available for on-demand fetching of missing objects"
  (https://git-scm.com/docs/partial-clone).

## Unreachable Acceptance

None.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Reading at `9fda0b65`: `decideReviewLineage`
  (`internal/cli/review_lineage.go`) treats `merge-base --is-ancestor` exit 1
  as non-ancestry and returns every other failure as `check review lineage
  ancestry: …`; `reviewLineageNotice` calls it for `deliver status` and
  `deliver retry`; `proveDismissalAnchor` and the `--fixed-by` branch of the
  `review dispose` command (`internal/cli/review.go`), `ProveReviewCorrection`
  (`internal/cli/deliver_review_correction.go`) and
  `reviewDispositionValidAtHead` run the same check on commits a review
  record or the ledger names.
- Git 2.56.0 probes in a scratch directory: `git cat-file -e` of a full
  object ID exits 0 when present and 1 with no output when absent; with
  `^{commit}` or an abbreviated or invalid name it exits 128, so the peeled
  form is not unambiguous; in a shallow clone a commit past the boundary
  exits 1; in a partial clone it starts a lazy fetch from the promisor
  remote and still exits 1 when that fetch fails; with an alternate object
  directory that no longer exists it exits 1 and prints
  `error: unable to normalize alternate object path`.
- A prototype of the design in a scratch worktree of `9fda0b65`: the five
  new tests fail on the authoring base (the Conexus error, a raw `merge-base`
  failure, the dispose and correction errors) and pass with the change; each
  listed sabotage fails its named subtest; every `Review|Lineage|Dispos|
  Dismiss|Correction|Deliver` test of `./internal/cli` passes once the one
  pinned assertion this Spec moves is updated, and
  `GOOS=windows go vet ./internal/cli` passes.
- Published sources (Exa): the three Git manual passages quoted above.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix review lineage missing commit squash merge gc" --all --files --min-score 0.3`,
  which returned this repository's mirrored CHANGELOG, the Spec 0175
  Archive Record on cleanup after a squash merge, the Correction Review
  decision and two decisions on cleaning up merged items; none decides how a review treats a recorded commit that is
  gone, and none contradicts this design.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB brief for Spec 0372 of 2026-11-10, the adopted Conexus message, and the maintainer's answer below.

The maintainer's answer (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~15:35 BRT), verbatim:

- Question "Com que autorizações eu toco sem supervisão as três Specs deste
  ciclo: 0372 (review do Conexus), 0371 (relatório de QA, que você já decidiu
  às 06:20 pelo Jarvis) e 0373 (regra de resposta)?" → "Mesmas do BA, teto
  minor (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas duas correções; a regra
  pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis, sem
  notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir."

The Conexus report is a measured defect with the reporter's own suggestion
("tratar o resultado \"not a valid commit\" de `git merge-base --is-ancestor`
(ou um `git cat-file -e` prévio do head anterior) como reinício de linhagem,
com um motivo como \"the reviewed head <prior head> is no longer in the
repository\""); no further question was asked. Taken from those and the
evidence, with no further interview:

- Absence is proven by `git cat-file -e <full object ID>` exiting 1 with no
  output, never by parsing `merge-base`'s message, and never with the
  `^{commit}` peel, which exits 128 for a missing object like any other
  failure.
- The proof runs only after the ancestry check fails with a status other
  than 1, so a review whose recorded commits are present runs no extra Git
  command and keeps today's error text when the commit is present.
- A shallow or partial clone cannot prove absence: a shallow clone can lack
  an ancestor that still is one, and asking a partial clone fetches from its
  promisor remote, where a failed fetch also exits 1. Both block with a named
  reason; a partial clone is never asked.
- A recorded merge base that is gone restarts like a reviewed head that is
  gone, with its own reason. Its check runs only after the reviewed head was
  proven an ancestor, and a merge base Roundfix recorded is an ancestor of
  that head, so in a complete repository it can be gone only when the record
  is inconsistent; a restart, never a pass, is the safe answer.
- `review dispose` refuses by name; the correction proof refuses a gone
  parked candidate by name and counts a gone fixing commit as an invalid
  disposition, the same answer as a commit the head does not contain.
- No ADR: ADR-0197 already restarts on a head that does not descend, and the
  proof only establishes that case; the reference records the rule.
- Two implementation Tasks and the QA gate: the guides and skill
  references, then the code with its tests, which names its guides through
  that dependency. The release is the operator's patch right
  after the merge; no Task releases.

## Open Questions

None.
