---
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
status: active
created: 2026-10-22
surfaces: [backend, cli]
---

# Recorded paths are read before carry-forward provenance

Minimal PRD for a bug fix, minted by `write-techspec`. The Daemon writes a
Task's `## Recorded paths` section at the end of the Task file, and Task
Carry-Forward later appends `## Carry-forward provenance` after it. The reader
in `internal/spec` accepts the recorded section only while it is the last
level-2 section, so every carried Task loses its recorded paths. A Delivery
Retry of a `corrective-spec-required` item then refuses a correction to a path
the archived Spec's Task recorded, `roundfix reconcile` keeps a Run Worktree
dirty over such a path, and the writer appends a second recorded section after
the provenance. After this Spec the reader finds the recorded section when only
other Daemon-written settlement sections follow it, and the writer replaces it
in place. The TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-22-retry-corrective-spec-required-ignora-recorded-paths.md`](references/2026-10-22-retry-corrective-spec-required-ignora-recorded-paths.md),
   origin oraculum, type hint `fix`, 2026-10-22. With roundfix 0.108.0, after
   the documented correction recipe for Spec 0092 there, `roundfix deliver
   retry` refused with `changed path
   packages/backend/src/application/use-cases/vendas/comparacao-resumo-vendas.use-case.ts
   is outside the named archived Specs and the paths their Tasks declared;
   author a corrective Spec with its own authorization and QA gate`, although
   that path is listed under `## Recorded paths` of the archived Spec's
   `task_02.md` at the recorded revision. The entry asks to confirm and fix
   the defect; oraculum's release of that Spec is due on 2026-10-25.

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` holds two entries, about the QA precondition
refusal report and the documentation contract suite's cost, and
`docs/findings/` holds none.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0) and is
delivered before Spec 0300 (maintainer decision of 2026-10-22, quoted in
`_authorization.md`). Spec 0299, still active, changes the Delivery Queue's
resume of an archived item; this Spec touches none of its files, and the
corrective-scope rule both rely on keeps its text.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads and rewrites an existing Task file section.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and a local fake reviewer. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0166: "The Daemon records the paths
  a Task changed without declaring them", which makes `## Recorded paths` a
  Daemon-owned section and is the decision this fix restores for carried
  Tasks. ADR-0176: "Citation checks read only what a Spec's authors wrote"
  cites it and does not apply: this Spec changes no citation check, and the
  settlement headings it skips stay non-authorial. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no Task creates, edits, renames, moves
  or deletes protected tooling configuration, scripts, ignore files, plugin
  declarations or version pins, and no Task changes a Governed Path (`paths:
  []` in the Spec-contained authorization record
  `docs/specs/0308-recorded-paths-are-read-before-carry-forward-provenance/_authorization.md`).
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A carried Task's recorded paths are read, so a Delivery Retry of a
  `corrective-spec-required` item accepts a correction to a path the archived
  Spec's Task recorded.
- A changed path that no Task of the named archived Specs declares or records
  is still refused, with today's exact text.
- `roundfix reconcile` counts a carried Task's recorded paths when it proves an
  archived Spec's Run Worktree superseded.
- Recording paths on a carried Task replaces its section in place, keeps the
  provenance bytes, and never leaves two recorded sections.
- An authored `## Recorded paths` heading followed by an authored section, and
  one inside a fenced block, are still not read as the Daemon's.

## Core Features

1. The reader finds the recorded section among the trailing settlement
   sections. A `## Recorded paths` section counts as the Daemon's when every
   level-2 section after it is a Daemon-written settlement section (`##
   Result`, `## Carry-forward provenance`, or another `## Recorded paths`),
   never when an authored section follows it and never inside a fence.
2. The writer replaces that section in place. Recording paths rewrites the
   found section where it stands and keeps every byte before and after it;
   with no section it appends one at the end, as today.
3. Every consumer follows. The corrective-retry proof, reconcile's archived
   Spec supersession, the Spec Consistency Check's Verification path
   provision, the Daemon's Task commit and `roundfix settle` read or write
   through the same reader, and their results change only for the carried
   shape.

## Non-Goals / Out of Scope

- Changing the corrective-scope rule: a changed path still needs a named
  archived Spec, its Tasks' declarations or recorded paths, and a correction
  outside them still needs a corrective Spec.
- Changing Task Carry-Forward, the provenance section's bytes, or where the
  Daemon writes a new recorded section.
- Changing the other readers that classify settlement headings by name (the
  pre-PR review's settlement region, the citation projection, the judge's
  planning pairs); they do not depend on section order.
- Rewriting Task files already committed in any repository.

## Success Metrics

1. Success Metric: through the real corrective-retry proof
   (`ProveReviewCorrection`) on a temporary repository, an archived Spec whose
   Task lists `review.txt` under `## Recorded paths` followed by `##
   Carry-forward provenance` accepts a correction that changes `review.txt`,
   for an Archive Record and for a legacy archive folder.
2. Success Metric: in the same fixture, a correction that also changes a path
   no Task declares or records is refused with `changed path other.txt is
   outside the named archived Specs and the paths their Tasks declared`.
3. Success Metric: the reader returns the recorded paths of a Task whose
   recorded section is followed by `## Carry-forward provenance`, by `##
   Result` and then provenance, or by a second recorded section, and returns
   none when an authored section follows or the heading is fenced.
4. Success Metric: recording paths on a carried Task writes exactly the
   expected bytes for the same paths (no write), new paths and no paths, keeps
   the provenance section byte for byte, and leaves exactly one recorded
   section when the file held two.
5. Success Metric: reconcile classifies an archived Spec's Run Worktree whose
   only dirty path is a carried Task's recorded path as superseded.
6. Success Metric: after release, oraculum's `roundfix deliver retry` of the
   Spec 0092 item is no longer refused for
   `comparacao-resumo-vendas.use-case.ts`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-22.

- The oraculum report in the adopted Inbox Entry, written before this Spec
  existed: the refusal text `changed path
  packages/backend/src/application/use-cases/vendas/comparacao-resumo-vendas.use-case.ts
  is outside the named archived Specs and the paths their Tasks declared`,
  whose generic form the corrective-retry test asserts.
- CommonMark Spec, 4.5 Fenced code blocks (https://spec.commonmark.org/spec):
  "The content of a code fence is treated as literal text, not parsed as
  inlines", the rule that keeps a fenced `## Recorded paths` line from being a
  heading.

## Unreachable Acceptance

- criterion: Success Metric 6, after release, oraculum's `roundfix deliver
  retry` of the Spec 0092 item is no longer refused for
  `comparacao-resumo-vendas.use-case.ts`.
  reason: the item lives in oraculum's Delivery Queue and Roundfix Home, which
  no Verification or QA row may read or write, and this Spec's authorization
  reads no other repository.
  satisfied-by: oraculum's own `roundfix deliver retry` with the released
  binary, reported back through the Secondbrain inbox.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2` with three `go test -overlay` probes that
  wrote nothing to the repository (`git status --porcelain` identical before
  and after). In `internal/spec`, a synthetic Task ending in `## Recorded
  paths` then `## Carry-forward provenance` gave `RecordedTaskPaths` no path,
  `DeclaredTaskPaths` only the `interface:` path, and `RecordTaskPaths`
  appended a second recorded section after the provenance. In
  `internal/cli`, the real `ProveReviewCorrection` refused a correction to the
  carried Task's recorded path with the Inbox Entry's text, for an Archive
  Record and a legacy folder. In `internal/worktree`, reconcile kept such a
  Run Worktree `dirty` (`recorded.txt is outside the archived Spec's declared
  and recorded paths`). The shape matches oraculum's Spec 0092 `task_02.md`
  and `task_03.md` at `01cf384` and `b12d75b`, read only for their heading
  order; no oraculum content was copied.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "Recorded paths
  carry-forward provenance Task file section reader corrective retry" --all
  --files --min-score 0.3`, which returned the Roundfix mirror's ADRs on
  carry-forward, review corrections inside the declared paths and dirty Run
  supersession, its glossary and its `review` guide. None records a decision
  about the recorded section's position; no Inbox Entry was written.
- Exa: the CommonMark Spec sections 4.2 ATX headings and 4.5 Fenced code
  blocks, which the reader's existing fence handling follows. No outside
  source bears on Roundfix's own settlement section order.

## Glossary

None.

## Decisions

- 2026-10-22 (about 19:50 BRT): the maintainer chose "Hotfix antes da 0300
  (Recomendado)" and "Hotfix + Oraculum move o teste (Recomendado)": fix only
  the defect, keep the corrective-scope rule unchanged, deliver before Spec
  0300.

Autonomous route: no person present; decisions taken from the operator's cycle AP hotfix briefing for Spec 0308 of 2026-10-22 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- The settlement sections that may follow the recorded section are the ones
  convention C2 names as Daemon-written: `## Result` and `## Carry-forward
  provenance`.
- A file that already holds two recorded sections is read as their union and
  rewritten with one, at the first one's place.
- No guide changes: the guides already promise that recorded paths count, and
  the Daemon still writes a new section at the end of the Task file.
- No new ADR: the section is already Daemon-owned (Project Constraints, Active
  ADR obligations); this Spec restores that ownership for carried Tasks.

## Open Questions

None.
