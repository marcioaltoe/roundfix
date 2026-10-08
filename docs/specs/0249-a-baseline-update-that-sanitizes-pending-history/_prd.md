---
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: active
created: 2026-10-08
surfaces: [backend, cli, docs]
---

# A Baseline update that sanitizes pending history

An adopter that upgrades Roundfix keeps every Legacy Archive Folder,
retired Finding, Backlog Entry, Review Artifact and handoff its history held
before Archive Records existed. Today the only way out is the History Sanitize
Command, run as a manual procedure from release notes: create the History Full
Tag, push it, then apply batches. On 2026-10-08 Fluxus, on Roundfix 0.59, still
held 57 pending units (23.9 MB), and two of them could not convert at all. Their
legacy `unproven` front matter is a list of maps, and the refusal reason printed
the YAML error over several lines
([the adopted Backlog Entry](references/2026-10-08-legacy-unproven-maps-refuse-a-folder.md)).

On 2026-10-08 the maintainer decided "No baseline update (Recommended)". After
`roundfix upgrade`, the repository's `roundfix baseline update` detects its
**Pending History** and includes the sanitize in the plan the operator already
reviews. Apply creates the History Full Tag when it is absent and converts the
history together with the rest of the Baseline Plan. `roundfix upgrade` only
says that this is waiting. The same Spec lets the two Fluxus folders convert and
prints each refusal on one line.

## Prerequisites

None. No other Spec is active in this worktree. Specs 0250 and 0251 are being
authored in parallel and deliver after this one. Neither one's artifacts are
read by this Spec's Verification. All three Specs touch the Roundfix Skill, so
each raises its version to the next free one when it records it.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The update
  result gains a `history` object with lowercase JSON keys and the existing tag
  name `history-full`, and the command gains the lowercase flag `--no-history`.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The update and the upgrade notice read local files and Git.
  Apply creates a local tag and never pushes it. Every test uses temporary
  repositories and a temporary home. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0254 (this Spec) decides that the
  update plans and applies the sanitize, batch size, Refused Units, the tag and
  the opt-out. The tag rule of ADR-0248: "holds every path the
  batch removes or rewrites"; the update keeps that rule per unit. ADR-0251: "It never
  applies to an active Spec", which holds for the list-of-maps `unproven` too.
  ADR-0100: "every byte
  outside a managed marker is identical before and after"; that proof still
  covers the instruction carriers. ADR-0173: "A Baseline Plan moves retired
  repository documents under the History Root", the precedent for a plan that
  changes the History Root. ADR-0184: "A TechSpec states a command surface as a
  transcript", answered by the TechSpec's Surface Transcripts. ADR-0177 cites ADR-0173 and
  does not apply: this Spec does not change the Relocation Citation scan. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths. Express maintainer
  authorization covers this: "considere autorizado a ajustar todas as skills se
  necessário", the standing grant "Concedo" for the Governed Paths each Spec
  declares, and the 2026-10-08 grant of the governed paths this Spec declares.
  No Baseline module, Baseline guide, Makefile, `.roundfixrc.yml`, CI workflow,
  formatter or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`. The
  Spec-contained authorization record is
  `docs/specs/0249-a-baseline-update-that-sanitizes-pending-history/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`, `skills/roundfix/SKILL.md`.

## Goals

- In a repository with Pending History, the update's preview lists every unit
  it would convert and every Refused Unit with its reason, and names the tag it
  would create. It writes nothing, and its Plan Digest binds the history it
  lists.
- The confirmed update converts every convertible unit in the same change as
  the Baseline Plan, creates the annotated History Full Tag when it is absent,
  and says the tag must be pushed. It never commits or pushes. A second update
  reports the history current.
- Refused Units never block the Baseline Plan and never change the update's
  state or exit code. `--no-history` leaves the history out entirely.
- After an upgrade, the notice names the Pending History of the working
  directory's repository, or says to run `roundfix baseline update` when it
  runs outside a repository.
- A Legacy Archive Folder whose `unproven` is a list of maps converts, with one
  stable text line per map. Every Refused Unit reason prints on one line.

## User Stories

1. As an adopter who just upgraded, I want `roundfix baseline update` to show
   my Pending History in the plan I already review, so that I do not have to
   learn a separate procedure from release notes.
2. As that adopter, I want one confirmation to apply the Baseline and the
   history together and to create the tag, so that I review and commit one
   change.
3. As a maintainer who prefers reviewed batches, I want `--no-history`, so that
   I can keep using `roundfix history sanitize --apply --batch <n>`.
4. As an operator running `roundfix upgrade`, I want the notice to tell me
   which history is waiting, so that I know to run the update.
5. As the Fluxus maintainer, I want the folders with the old `unproven` form to
   convert, and every refusal to fit on one line, so that the plan stays
   readable.

## Core Features

1. **Pending History in the update plan.** The Managed Refresh plans every
   pending unit through the History Sanitize Command's own planning. It lists
   the units it can convert, each Refused Unit with its reason, and the History
   Full Tag it would create, and it binds that section into its Plan Digest
   (ADR-0254).
2. **One confirmed change.** `--yes` or `--confirm-plan <digest>` applies the
   Baseline Plan, creates the tag when it is absent and converts every
   convertible unit. A second update reports the history current. The update
   never commits or pushes (ADR-0254).
3. **Refusals that never block.** A unit with uncommitted changes, a unit the
   existing tag does not hold, and a unit the conversion refuses are Refused
   Units. A lightweight or non-ancestor tag blocks only the history section
   (ADR-0254).
4. **The opt-out.** `--no-history` skips the detection and the sanitize
   (ADR-0254).
5. **The upgrade notice.** After a release outcome, `roundfix upgrade` names
   the Pending History, or says to run the update when it runs outside a
   repository. Its output and exit code do not change (ADR-0254).
6. **Legacy `unproven` maps.** Under Lenient Legacy Reading a list-of-maps
   `unproven` becomes one stable text line per map. An active Spec stays
   string-only, and every Refused Unit reason prints on one line (ADR-0254).
7. **Docs and glossary.** `CONTEXT.md` gains **Pending History** and revises
   **Managed Refresh**, **History Sanitize Command**, **Refused Unit**,
   **Lenient Legacy Reading**, **Sanitize Batch** and **History Full Tag**. The
   command references and the Roundfix Skill describe the behavior.

## Non-Goals / Out of Scope

- Changing the History Sanitize Command's selection, batch counting,
  promotion, advice, tag rules or exit codes. Only its Refused Unit reasons
  change, which now print on one line.
- Promotions or Jev advice from the update. An operator who needs them passes
  `--no-history` and uses the History Sanitize Command.
- Moving, replacing or pushing an existing `history-full` tag, or creating a
  second tag for history retired after it.
- Running the sanitize from `roundfix upgrade`, from the interactive
  `roundfix baseline` workflow, or from `roundfix baseline plan` and
  `roundfix baseline apply`.
- Requiring a clean working tree for the Baseline part of the update.
- Editing archived bytes, existing Archive Records or any adopter repository.

## Success Metrics

1. Success Metric: in a temporary adopted repository with two convertible
   Legacy Archive Folders, one with a list-of-maps `unproven`, one folder that
   cannot convert, and one retired Finding, the preview exits 3, lists three
   units and one Refused Unit, names the tag it would create, and leaves the
   tree and the refs byte-identical.
2. Success Metric: `--confirm-plan` with that preview's digest exits 0. It
   writes two Archive Records and one Reduced History Entry, leaves the refused
   folder byte-identical and creates an annotated `history-full` at the
   previous `HEAD`. No commit is created and no remote ref changes. A second
   update reports the history current.
3. Success Metric: in a repository with no Pending History, the update's Plan
   Digest equals the Baseline Plan Digest, and every existing
   `baseline update` test passes unchanged.
4. Success Metric: `--no-history` over the repository of Success Metric 1
   reports the history skipped and changes no history file and no tag.
5. Success Metric: `roundfix upgrade` in that repository prints one notice line
   naming three Legacy Archive Folders and the findings kind, the refused
   folder included because it is still pending; outside a repository it prints the line that names `roundfix baseline update`;
   its stdout and exit code equal today's.
6. Success Metric: the list-of-maps `unproven` record reads one line per map in
   a stable key order, and an active Spec with the same front matter still
   fails as it does today.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Fluxus report, triaged in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-07-history-sanitize-recusa-unproven-em-formato-antigo.md`.
  It carries the exact multi-line refusal (`yaml: unmarshal errors:` then
  `line 10: cannot unmarshal !!map into string`) and the two folders' key sets:
  `row`, `goal`, `claim` and `satisfied-by`, and `row`, `claim` and `reason`.
  The adopted Backlog Entry records the 0.59 measurement: 57 units, 23.9 MB,
  2 refused.
- The adopter release notes Roundfix sent through the Secondbrain inbox on
  2026-10-07 (`inbox/conexus/2026-10-07-roundfix-v0-57-0-o-que-muda.md`). They
  show the four-step manual procedure this Spec folds into the update: run the
  update, plan the sanitize, create and push the tag, then apply batches.
- Django's development server, which reports "You have N unapplied
  migration(s)" and "Run 'python manage.py migrate' to apply them". The
  framework only notifies, and the apply is a separate explicit command.
- The Git documentation for `git push` (git-scm.com/docs/git-push). The default
  `push.default=simple` pushes the current branch only, so a locally created tag
  is not on the remote until it is pushed.

## Glossary

- adds: **Pending History**
- changes: **Managed Refresh**
- changes: **History Sanitize Command**
- changes: **Refused Unit**
- changes: **Lenient Legacy Reading**
- changes: **Sanitize Batch**
- changes: **History Full Tag**

## Decisions

- The update plans and applies the Pending History in the change it plans,
  bound by one Plan Digest; see ADR-0254.
- All convertible units at once, with no batch flag on the update, and
  `--no-history` for ADR-0248's reviewed batches; see ADR-0254.
- Refused Units are listed on one line each and never block or change the
  state; see ADR-0254.
- The update creates the annotated tag at `HEAD` when it is absent, never moves
  one, and says it must be pushed; see ADR-0254.
- A legacy list-of-maps `unproven` becomes one stable line per map; see
  ADR-0254.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
