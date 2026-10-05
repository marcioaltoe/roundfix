---
spec: 0228-queue-items-that-stay-current-with-main
status: active
created: 2026-10-04
surfaces: [backend, cli, docs]
---

# Queue items that stay current with main

This is a bug fix with two causes, both recorded on 2026-10-04 in the
operator's intervention log and adopted from two Backlog Entries:
[queued Specs raise the same skill version](references/2026-10-04-queued-specs-raise-the-same-skill-version.md)
and [a docs fix after archive forces a manual Pull Request](references/2026-10-04-a-docs-fix-after-archive-forces-a-manual-pull-request.md).

1. **Two items record the same owned-skill version.** Specs 0223, 0226, 0225
   and 0224 each raised the Roundfix Skill version from the tree their Task
   started on. 0226 and 0225 both recorded `0.1.26`, and 0224 recorded
   `0.1.25` over the default branch's `0.1.27`. One item parked
   `gate-failed` and another `pull-request-conflict`; the operator merged the
   default branch, took its version lines and raised to the next patch by
   hand twice (entries 157 and 159). Earlier, a Codex Task agent wrote the
   version digest by hand instead of running the record command, and since
   the record never replaces a digest every later attempt failed (entries 150
   and 154).
2. **A review-only correction after archive needs a corrective Spec.** The
   pre-PR review of 0225's archived candidate raised two findings on the
   archived Spec's own records. The operator fixed one in a docs-only commit
   and dismissed the other with evidence, and `roundfix deliver retry`
   refused with `corrective-spec-required` because the head moved; the
   operator ran review round 2 and opened Pull Request #382 by hand (entry
   158). 0221 ended the same way (entries 137 and 138).

## Prerequisites

None. Specs 0227 and 0229 are authored in the same cycle and may also raise
the Roundfix Skill's version. Until this Spec merges, the operator orders the
queue so that a later Spec raises it from the earlier one's value; task_04
records this Spec's raise with the record command on the tree it starts from.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; queue
  items, Runs, blockers, review findings and skill names keep their names.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the record command, the derived
  merge and the retry read and write local files, local Git and the Run
  Database; the merge keeps the existing fetch of the default branch, and no
  credential, forge read or network call is added. No test reaches GitHub.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0233 (this Spec) decides the three
  rules, ADR-0233: "The owned-skill record command chooses the version", and
  ADR-0233: "Anything else still requires a corrective Spec". ADR-0189 ties a
  version to its content, ADR-0189: "A version names one content". ADR-0192 resolves a conflict
  confined to declared derived paths by regeneration, ADR-0192: "Any other
  conflicted path aborts the merge", and this Spec adds line-scoped paths
  beside it. ADR-0165 parks a finding after archive for a corrective Spec,
  ADR-0165: "an archived Spec is never edited to absorb a finding", which
  still holds for a finding about the delivered change. ADR-0197 bounds the
  Reviewer Lineage to two rounds, ADR-0197: "Round 2 reviews only the diff
  from the round-1 head to the current head". ADR-0223 keeps a delivery's work
  across archive and review, ADR-0223: "A park for a corrective Spec still
  refuses a moved head", which this Spec narrows for a review-only
  correction. ADR-0229 keeps the operator-archive anchor unchanged, ADR-0229: "The
  Run start head is a weak anchor". ADR-0187
  splits the Roundfix Skill by command and ADR-0189 governs both skill
  edits. ADR-0184: "A TechSpec now declares numbered Surface Transcripts". The
  gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156
  and ADR-0167, and ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  this Spec's consistency by citation and receipt. ADR-0178 and ADR-0179
  decide the grant each Task commit runs under, ADR-0182 settles a Task on the
  facts its gate checks, and ADR-0166 records undeclared paths; every Task
  declares its paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the
  gate's machine stage and row carry, ADR-0194, ADR-0195 and ADR-0210 cite
  ADR-0097 but decide what a QA row records, when it is observed again and its
  evidence snapshot, ADR-0169 cites ADR-0165 but decides the review's
  merge-base diff, and ADR-0196 cites ADR-0169 but decides when a review
  finding parks; this Spec changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  standing skills authorization of 2026-09-30, "considere autorizado a
  ajustar todas as skills se necessário", the authorization of the same day
  for guides and `.roundfixrc.yml`, "Autorizar os dois", and the approval of
  this cycle on 2026-10-04, "Aprovado". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0228-queue-items-that-stay-current-with-main/_authorization.md`;
  bounded files: `.agents/skills/implement-task/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`, `.roundfixrc.yml`,
  `docs/agents/specific-repository.md`,
  `skills/implement-task/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

- Two queued items that raise the same owned skill never record the same
  version, whatever order they merge in.
- A Task agent records an owned-skill version with the record command, and a
  stale or colliding record no longer fails every later attempt.
- A correction that only answers the pre-PR review of an archived candidate
  returns the item to review round 2 through `roundfix deliver retry`, with
  no hand-run review and no hand-opened Pull Request.
- Every other change after an archived review still requires a corrective
  Spec, and every existing derived-merge and retry refusal holds.

## Core Features

1. **The record command chooses the version.** The owned-skill record command
   raises a skill whose content is not recorded under its declared version,
   and whose declared version is not above every recorded one, to one patch
   above the highest recorded version, in both version fields of the
   canonical `SKILL.md` and its mirror, then records it (ADR-0233).
2. **Line-scoped derived paths.** A derived-path declaration may name paths
   whose conflicts are resolved line by line: a conflict whose hunks hold
   only lines matching the declared pattern takes the default branch's side
   of each hunk, and the regeneration may change those paths only on matching
   lines. This repository declares the version fields of every `SKILL.md` for
   the owned-skill record command (ADR-0233).
3. **A review-only correction returns to round 2.** A Delivery Retry of a
   `corrective-spec-required` item whose head moved returns it to
   `reviewing` when the head descends from the parked candidate, every
   standing finding of the review at that candidate has one disposition, and
   every changed path lies under an archived Spec the blocker names
   (ADR-0233).
4. **The guides and skills say so.** The Roundfix Skill's `deliver` and
   `review` references, their user guides, the configuration guide, the
   repository's skill version rule and the `implement-task` skill describe
   the three rules; the `implement-task` skill tells the agent that a command
   a requirement names is part of the work even when Verification runs the
   same test.

## Non-Goals / Out of Scope

- Refreshing an item from the default branch before its Run, or serializing
  items that raise the same skill.
- Replacing a recorded digest, or deciding which versions a release shipped.
- Changing the whole-file resolution of declared derived paths, the
  Pre-PR Review Command's output, its corrective-Spec line, the two-round
  ceiling, Park Class names, blocker names or retry limits.
- Accepting a correction that changes anything outside the archived Spec's
  directory, or one with an undisposed finding.

## Success Metrics

1. Success Metric: in a disposable skill tree whose record holds `0.1.26` with
   other content, recording a skill declared `0.1.26` writes `0.1.27` into
   both version fields of the skill and its mirror and records it, where today
   it fails with "content changed under version".
2. Success Metric: in a disposable repository, two branches that raised the
   same version line to different values and appended to the same record
   merge through the derived merge without parking, and a conflict hunk
   holding any other line still parks `pull-request-conflict`.
3. Success Metric: a `corrective-spec-required` item whose head adds one
   commit under its archived Spec, with one finding fixed by that commit and
   one dismissed with evidence, resumes at `reviewing` with the head as its
   newest candidate; a change outside that directory, an undisposed finding or
   a non-descendant head still refuses.
4. Success Metric: every existing derived-merge, owned-skill version and
   corrective-Spec retry test passes, with two declared changes: the record
   mode cases of `TestRecordingNeverReplacesARecordedVersion` and the
   repository declaration of `TestThisRepositoryDeclaresItsToolsAndDerivedPaths`.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The operator's intervention log
  (`~/.roundfix-operator/queue-interventions.md`), entries 137 and 138 (0221),
  150 and 154 (0225's hand-written and early digests), 157 (0225 and 0226
  both at `0.1.26`), 158 (0225 review round 2 and Pull Request by hand) and
  159 (0224 raised to `0.1.28` by a hunk-only resolution).
- The live Run Database, read with `sqlite3 -readonly` on 2026-10-04: the
  0225 item is `parked`, blocker
  `corrective-spec-required: 0225-a-jev-ceiling-the-maintainer-sets`,
  candidate commits `1f745c31` and `814ca147`.
- Git's `git merge` documentation (<https://git-scm.com/docs/git-merge>),
  which defines the hunks the line-scoped resolution reads: "The part before
  the `=======` is typically your side, and the part afterwards is typically
  their side".
- The Changesets project's decisions
  (<https://github.com/changesets/changesets/blob/main/docs/decisions.md>),
  a published account of choosing the version when changes are combined
  rather than in each change: changesets "flatten the version bumps into one
  single bump" so that they "can be added and accumulated safely".
- Gerrit's review label documentation
  (<https://gerrit-review.googlesource.com/Documentation/config-labels.html>),
  a published rule for carrying a review across a new patch set whose change
  kind is bounded: `NO_CODE_CHANGE` means "only the commit message may be
  different".
- The Secondbrain holds no independent entry on either defect; its searches
  on 2026-10-04 returned mirrors of the two Backlog Entries, ADR-0223, Spec
  0219 and an upstream git-rebase skill reference.

## Decisions

- The record command chooses the version; a merge regenerates it. See
  ADR-0233.
- Line-scoped derived paths resolve only conflict hunks of matching lines.
  See ADR-0233.
- A review-only correction after archive returns to round 2. See ADR-0233.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
