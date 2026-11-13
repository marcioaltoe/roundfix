---
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: active
created: 2026-11-13
surfaces: [backend, cli, docs]
type: fix
---

# A parked item always has a way back through a Delivery Retry

Reported by fiscus on 2026-11-13: a Delivery Queue item parked by a network timeout before it had an Item Branch was refused by `roundfix deliver retry` with `use item branch: branch is required`, so it had no way back inside its queue. Reported by quaestor on 2026-11-12: after the operator resolved a Pull Request conflict by removing the last file under `docs/specs/`, the retry stopped reading the delivery authorization of the archived Spec because the Spec Root directory no longer existed, and two later retries parked the item `review-stale` because the resolution commit was never pushed. After this Spec a branchless item returns to `queued`, the authorization read needs no Spec Root directory, and a Pull Request that holds an earlier candidate sends the item back through the gate to be pushed.

## Sources

Two Backlog Entries, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md`](references/2026-11-13-a-retry-refuses-an-item-parked-before-it-has-a-branch.md),
   `type: fix`, recorded by the operator from ai-memory message
   `01a1d42a-28f1-72d3-9957-76f106877d2e`; reported by fiscus on 2026-11-13.
2. [`references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md`](references/2026-11-12-a-retry-after-resolving-a-pull-request-conflict-finds-no-specs-root-and-no-pushed-head.md),
   `type: fix`, recorded by the operator from ai-memory message
   `01a1d231-03da-79f1-9834-51487a593986`; reported by quaestor on
   2026-11-12. The Queue Owner panic the same message repeats belongs to Spec
   0385, already delivered.

## Prerequisites

None in `requires`. Authored against `f953f77f` (v0.143.0 plus the Backlog
Entry merges #856 and #857). Spec 0388, authored in parallel, delivers first.
It changes the `## Completion notices` section of the Roundfix skill's
`references/release.md` and its contract test. The two Specs share only these
files, each in its own part:

- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md`, and `skills/testdata/owned-skill-versions.json`:
  task_04 raises the next version with the record command on the tree it runs
  on and never writes a number.
- `docs/references/behavior-surfaces.json`: each guide or skill Task of this
  Spec re-records it after its own last edit.
- `CONTEXT.md`, only if 0388 edits it: this Spec revises only the **Delivery
  Retry** entry.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed;
  the Item Branch keeps its name rule. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the push the owner makes is the existing publication push.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0398 (new) records the decisions:
  "A parked item always has a way back through a Delivery Retry". ADR-0193
  applies, extended: a Delivery Retry already returns a prerequisite park to
  `queued`, and now does the same for any item parked before it has an Item
  Branch. ADR-0223 applies unchanged: a retry of an archived item whose head
  descends from its last candidate still returns it to review, and the archive
  commit stays an exact move. Not applicable: ADR-0229, ADR-0237 and ADR-0312,
  which cite ADR-0223, and ADR-0262, which cites ADR-0193 (reasons in the
  TechSpec's ADR Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BE answers of 2026-11-13 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes. No Makefile,
  workflow, Go module, lint or formatter change, and no Baseline module,
  template or profile change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- An item parked before it has an Item Branch returns to `queued` through
  `roundfix deliver retry`, and its owner then creates the branch as on a
  first start.
- The delivery authorization of an archived Spec is read whether or not the
  Spec Root directory exists in the item worktree.
- An item whose Pull Request holds an earlier candidate of the item is gated
  again and its current candidate is pushed, instead of parking
  `review-stale`.
- A `review-stale` park names the action that answers it.

## User Stories

1. As an operator whose queue parked an item on a network timeout before its
   branch existed, I want `deliver retry` to return it to the queue, so that I
   do not have to wait for the owner to exit and record a new queue.
2. As an operator who resolved a Pull Request conflict that emptied the Spec
   Root, I want the retry to read the archived Spec's authorization, so that
   publication and merge continue without a `mkdir`.
3. As an operator whose resolution commit never reached the Pull Request, I
   want the owner to push it after the gate, so that the item merges without a
   manual `git push`.
4. As an operator who meets a `review-stale` park, I want status to tell me
   what to do, so that I do not have to read the code.

## Core Features

1. The Delivery Retry returns a parked item with no recorded Item Branch to
   `queued`, counting a retry, and clears a worktree recorded without its
   branch.
2. The delivery authorization resolves the Spec Root with the absent-tolerant
   resolver, and the Spec authorization reader resolves an absent Spec Root
   inside the repository through its nearest existing ancestor.
3. At `checking`, a Pull Request head that is an earlier candidate of the item
   returns the item to `gating`; any other head parks `review-stale`.
4. The `review-stale` Park Class answer names its next action.
5. The `deliver` guide, both copies of the Roundfix skill's `deliver`
   reference and the **Delivery Retry** entry say so.

## User Experience

Fiscus's refusal of 2026-11-13 read:

```text
Retry refused
Reason: retry Delivery Queue item "0092-...": use item branch "": use item branch: branch is required
```

With this Spec the same retry prints `Retried <slug>: delivery-error: read
prerequisites: ... -> queued` and hands the item to the owner. A
`review-stale` park prints in `deliver status`:

```text
Park: <slug> review: in <worktree>, merge into the item branch any commit its Pull Request's head branch holds that the item branch lacks, then run roundfix deliver retry <slug>
```

## Non-Goals / Out of Scope

- Retrying a failed network command inside the Queue Owner. The owner runs
  several network commands (the two fetches of the default branch, the push
  and every `gh` call); a bounded retry of only the prerequisite fetch would
  give a partial guarantee, and a uniform policy is a separate decision
  (ADR-0398). After this Spec the park is recovered by `deliver retry`.
- Pushing from `roundfix deliver retry`, which stays a local command.
- Changing how the owner resolves a Pull Request conflict, or the
  `delivery-error` park that a failed `git restore` produces.
- Removing a worktree directory that a failed worktree creation left on disk;
  the requeued item records a new Item Branch and worktree beside it.
- Messages to fiscus or quaestor: the operator sends them.

## Success Metrics

1. Success Metric: an item parked `delivery-error` with no recorded Item
   Branch, through the production Delivery Engine and workflow over a real
   temporary repository, retries to `queued` with one counted retry and no
   recorded workspace, and `CreateItemBranch` then creates a branch at the
   default branch; the same holds when a worktree was recorded without its
   branch. At `f953f77f` the retry refuses with `use item branch: branch is
   required`.
2. Success Metric: `Authorization` of an archived Spec whose Spec Root left
   the working tree returns `push`, `pull_request` and `merge` and creates no
   directory; `ReadSpecAuthorization` grants at the revision with the Spec
   Root removed. At `f953f77f` the first fails with `specs.root resolved to
   ..., which does not exist` and the second is `unresolved`.
3. Success Metric: an item at `checking` whose Pull Request holds its earlier
   candidate runs the gate, pushes its current candidate once and merges; one
   whose Pull Request holds a head the item never had parks `review-stale`
   with no push. At `f953f77f` both park `review-stale`.
4. Success Metric: status and the Pending Question answer a `review-stale`
   park with the next action of User Experience.
5. Success Metric: the `deliver` guide, both copies of the Roundfix `deliver`
   reference and the **Delivery Retry** entry carry the new text; the
   documentation, skills and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read through
Exa on 2026-11-13, and on Git itself.

- Junio C Hamano, git mailing list, "Re: Unexpected empty directory removal"
  (https://public-inbox.org/git/7vehrs6bt1.fsf@alter.siamese.dyndns.org/): "As
  Git does not track directories at all, but merely uses directories as a means
  to instantiate files (which it tracks), when the last file is removed as the
  result of a merge in repository B, it notices that the directory is no
  longer needed to hold anything it cares about, and removes it." This is why
  quaestor's `git rm -r` of the last Spec folder left no `docs/specs/`.
- The installed `git`, which this Spec did not build: in a temporary
  repository whose `docs/specs/` holds one tracked file, `git rm -r -q
  docs/specs/<slug>` leaves no `docs/specs` directory.
- actions/checkout, `src/retry-helper.ts` at commit
  `11bd71901bbe5b1630ceea73d27597364c9af683`
  (https://github.com/actions/checkout/blob/11bd71901bbe5b1630ceea73d27597364c9af683/src/retry-helper.ts):
  `const defaultMaxAttempts = 3`, `const defaultMinSeconds = 10`, `const
  defaultMaxSeconds = 20`, used for its `fetch`, `getDefaultBranch` and
  `lfsFetch` commands alike. This is the practice ADR-0398 weighs and defers:
  a retry policy that covers every network command, not one fetch.

## Unreachable Acceptance

None. Every criterion is proved with temporary repositories, the Run Database
and fake Pull Request boundaries; no criterion needs the network, an open
Pull Request or another repository.

## Validation

The decisions were checked on 2026-11-13 against these sources.

- Characterization at `f953f77f`, with focused tests in a scratch copy:
  - the production engine's `Retry` of an item parked `delivery-error: read
    prerequisites: ...` with no branch returned `retry Delivery Queue item
    "<slug>": use item branch "": use item branch: branch is required` and
    left the item parked with no retry counted;
  - `Engine.park` writes the in-memory item, so a park after
    `CreateItemBranch` recorded its pair but failed writes an empty branch
    over the recorded one and keeps the recorded worktree;
  - `commandDeliveryWorkflow.Authorization` on an archived Spec with
    `docs/specs` removed returned `specs.root resolved to "<repo>/docs/specs",
    which does not exist`, and `spec.ReadSpecAuthorization` with that root
    returned `unresolved` (`unreadable_record`, field `spec_root`: `git
    rev-parse --show-toplevel: fatal: cannot change to '<repo>/docs/specs'`);
  - an engine item at `checking` with candidates `pushed-head` and
    `resolved-head` whose Pull Request holds `pushed-head` parked
    `review-stale` with no push.
- A prototype of tasks 01 to 04 on a scratch copy at `f953f77f`: the new tests
  passed, `./internal/delivery`, `-tags docscontract ./internal/docscontract`
  and `./skills` passed whole, and each sabotage of the TechSpec made its named
  test fail on the named case.
- Secondbrain, read only: `wiki/index.md`, then `qmd query --all --files
  --min-score 0.3` for a parked queue item's retry after a network failure and
  an unpushed head; it returned this repository's own mirrored ADRs and
  `deliver` guide, and the concept
  `wiki/concepts/agendas-e-reexecucao-em-motores-de-workflow-durable.md`, whose
  rule that a resumed schedule restarts "a partir de agora" without replaying
  the exhausted fire matches requeueing an item that never started rather than
  replaying its failed step.
- Exa: the two published sources above.

## Glossary

- changes: **Delivery Retry**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BE brief for Spec 0389 of 2026-11-13, the maintainer's answers below and the adopted sources.

The maintainer's answers (AskUserQuestion, 2026-11-13 ~08:20 BRT), verbatim:

- Question "Quais trabalhos do backlog do Roundfix eu planejo e entrego agora,
  sem supervisão? (Cada item marcado vira Spec com plano validado e é entregue
  pela fila. Os desmarcados ficam no backlog.)" → selected "Retry da fila
  (Recommended)" — option text: "Uma Spec com os dois defeitos do retry: item
  parado antes de ter branch não aceita retry (Fiscus), e o retry depois de
  resolver conflito no PR não acha a pasta de Specs nem envia o head
  (Quaestor)." This Spec is that item.
- Question "Quando publico as releases desses trabalhos?" → "Patch após cada
  Spec (Recommended)" — option text: "v0.143.1 assim que sair a correção do
  aviso; a próxima patch assim que sair o retry. Os adotantes que já trouxeram
  a 0.143.0 (a Argus, por exemplo) recebem logo a correção do texto."
- Question "Quanta folga eu tenho para consertar falhas sem parar? (Na 0386 o
  review usou 3 correções.)" → "Até 3 por Spec + Makefile/CI (Recommended)" —
  option text: "Até 3 correções por Spec, contando as rodadas do review. Fica
  liberado mexer em Makefile, workflows do CI, go.mod e lint quando a correção
  exigir; cada mudança vai descrita no relatório."

Taken from those answers, the sources and the characterization, with no
further interview:

- `type: fix`: a refused retry, a failed read and a wrong park become correct
  behavior; the Release Ceiling is `patch`.
- A branchless item is requeued like a prerequisite park, and the retry is
  counted, because the park was the item's own failure to start.
- The authorization read tolerates an absent Spec Root at both layers, because
  it reads Git history and never the working tree.
- The owner pushes, never the retry, and only after the repository gate runs
  again on the candidate the Pull Request lacks.
- No automatic retry of a failed fetch in this Spec (Non-Goals, ADR-0398).
- ADR-0398 records the decisions. Three code Tasks, each documenting its own
  guide text (task_01 also revises the glossary), one Task for the Roundfix
  skill's `deliver` reference and version, and the QA gate.

## Open Questions

None.
