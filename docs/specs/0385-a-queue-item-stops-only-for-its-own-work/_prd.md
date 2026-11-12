---
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: active
created: 2026-11-12
surfaces: [backend, cli, docs]
type: fix
---

# A queue item stops only for its own work

Reported by quaestor and recorded by the operator on 2026-11-12: a Delivery Queue item stood at `gating` for about 8.5 hours because its Queue Owner panicked writing the Machine Verification Capacity wait line to a nil writer, and `roundfix deliver status` named no owner because the dying owner released its own record. Recorded by the operator on 2026-11-11: two real review findings anchored `path:line.` blocked an item as `findings name no file and line`. Reported by pantheon on 2026-11-12: every run of the review validator became a stray ai-memory project, because its sealed session works in a bare temporary directory. After this Spec the owner logs the wait and continues, an owner that dies with an item left stays visible in status, an anchor may end a sentence, and the validator's session carries the reviewed repository's own `.ai-memory.toml`.

## Sources

Two Backlog Entries and two cross-project messages, adopted by this Spec's
planning change (see `references/_index.md`):

1. [`references/2026-11-12-the-delivery-owner-crashes-when-the-repository-gate-waits-for-machine-verification.md`](references/2026-11-12-the-delivery-owner-crashes-when-the-repository-gate-waits-for-machine-verification.md),
   `type: fix`, recorded by the operator on 2026-11-12 from the 0382 delivery.
2. [`references/2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md`](references/2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md),
   ai-memory message `01a1d01d-ac6c-71a0-a548-0398a33c1ba6`, reported by
   quaestor on 2026-11-12: the same panic on Roundfix 0.141.1, with "o `deliver
   status` continua mostrando o item em `gating` e não mostra nenhum processo
   dono".
3. [`references/2026-11-11-a-review-finding-anchored-as-path-line-and-a-period-blocks-the-item.md`](references/2026-11-11-a-review-finding-anchored-as-path-line-and-a-period-blocks-the-item.md),
   `type: fix`, recorded by the operator on 2026-11-11 from the 0376 review.
4. [`references/2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md`](references/2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md),
   ai-memory message `01a1d01a-fa8c-7211-bd88-5023f39563e1`, reported by
   pantheon on 2026-11-12. This Spec owns only its item 2, the validator's
   working directory; the operator commits this repository's `.ai-memory.toml`
   (item 1) in the planning change.

## Prerequisites

None in `requires`. Authored against `d79b30b4` (v0.142.0 plus the Backlog
Entry merge). Specs 0383 and 0384 deliver first, in that order. 0383 changes the
skills-lock reading, `doctor` and `roundfix skills install`; 0384 changes the
`Makefile` cache block and amends its decision record. Neither touches a file this Spec changes
except these shared ones, each in its own section:

- `CONTEXT.md`: this Spec revises only the **Queue Owner** entry.
- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md` and `skills/testdata/owned-skill-versions.json`:
  task_04 raises the version with the record command from the tree it runs on
  and never writes a number.
- `docs/references/behavior-surfaces.json`: each guide Task re-records it after
  its own last guide edit.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the marker copy is a local file. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0395 (new) records the decisions:
  "A queue item stops only for its own work". ADR-0279 applies unchanged:
  status reports a recorded owner proven gone, and a queue with no recorded
  owner still adds no line. ADR-0373 applies unchanged: an identity mismatch
  proves the kept owner gone. ADR-0384 applies unchanged: the repository gate
  still waits for the Machine Verification Capacity and prints the same line.
  ADR-0196 applies, refined: a finding still opens with a `path:line` anchor,
  which may now end a sentence. ADR-0069 applies unchanged: Baseline analysis
  still works in an empty private directory. Not applicable: ADR-0056 and
  ADR-0385, which cite ADR-0384; ADR-0286 and ADR-0317, which cite ADR-0196;
  ADR-0180 and ADR-0344, which cite ADR-0069 (reasons in the TechSpec's ADR
  Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BD answers of 2026-11-12 quoted in the authorization record; the
  Governed Paths are the Roundfix skill files task_04 changes. No Makefile,
  workflow, Go module, lint or formatter change, and no Baseline module,
  template or profile change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0385-a-queue-item-stops-only-for-its-own-work/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Queue Owner whose gate waits for the Machine Verification Capacity logs
  the wait in its console log and continues when the capacity frees, instead
  of dying.
- An owner that dies while an item can still advance stays visible: `deliver
  status` names it and says to run `roundfix deliver resume`.
- A finding whose anchor ends a sentence is validated, never blocked as
  unanchored.
- A review validator session is recorded by session tooling under the
  reviewed repository, while it still sees no repository file.

## User Stories

1. As the operator, I want the owner to survive a capacity wait, so that an
   item reaches its merge without my intervention.
2. As the operator or an adopter's agent, I want `deliver status` to say that
   the owner is gone, so that I run `roundfix deliver resume` instead of
   waiting for hours.
3. As the maintainer, I want a reviewer's well-formed finding that ends with a
   period to be judged on the diff, so that a real defect parks the item for
   the right reason.
4. As a user of ai-memory, I want the validator's sessions under the reviewed
   repository's project, so that no stray project appears per review.

## Core Features

1. The production Delivery Queue workflow and engine share the owner's log, and
   no workflow message reads a nil writer.
2. The Queue Owner releases its record on exit only when no item can advance.
3. The anchor parser accepts `.` or `;` after the line when a space or the end
   of the finding follows.
4. The validator copies the reviewed checkout's regular `.ai-memory.toml` into
   its private working directory, and a sealed working directory admits that
   one file.
5. The `deliver` and `review` guides, the Roundfix skill's `deliver` and
   `review` references and the **Queue Owner** glossary entry say so.

## User Experience

The owner's console log of the 0382 delivery ended in a panic at
`deliver_workflow.go:1713`. With this Spec the same wait reads, in that log:

```text
roundfix: Delivery Queue item 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote: waiting for the Machine Verification Capacity; 1 of 1 held by Run <run id> of 0096-<slug> in <quaestor checkout>, Task task_qa Verification attempt 1
```

and an owner that dies anyway leaves this in `deliver status`:

```text
0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote	gating	-	<worktree>
Owner: Delivery Queue owner PID 49802 is not running; run roundfix deliver resume
```

## Non-Goals / Out of Scope

- Reporting a queue that records no owner: ADR-0279 rejected it, and the
  owner now keeps its record when it dies with an item left.
- Recovering a panic inside the owner or parking the item for it.
- Baseline semantic analysis and profile proofs, whose sealed sessions also
  run in bare temporary directories: their interface takes no repository by
  design (ADR-0069), they run rarely, and the operator's `/var/folders` route
  covers them; this Spec leaves them as they are and its agent rule already
  admits the marker for a later change.
- Run Worktrees and item worktrees: they are checkouts of the repository, so
  the `.ai-memory.toml` the operator commits at the root in the planning
  change is present in every worktree created from a commit that holds it,
  and ai-memory reads it there; worktrees created before it are not covered.
- Writing an ai-memory marker the repository does not have, or deriving its
  workspace or project.
- Messages to quaestor or pantheon: the operator sends them.

## Success Metrics

1. Success Metric: a repository gate built by the production dependencies,
   waiting behind a capacity unit another Run holds, writes the wait line
   naming that Run to the owner's log, passes when the unit is released, and
   leaves no capacity entry; at `d79b30b4` the same gate panics.
2. Success Metric: a Queue Owner whose engine panics or fails while an item is
   at `gating` keeps its owner PID and identity; one whose last item merged
   releases them; at `d79b30b4` both failures clear the record.
3. Success Metric: `path:2.`, `path:2;`, backticked `path:2-4`. and
   `notes.v2.md:3.` followed by a space or the end are read with their path
   and lines; `path:2.5`, `path:2.md`, `path:2;x` and `path:2..` are not; a
   review answer of two such findings anchored in the candidate diff ends
   with standing findings, not `blocked`.
4. Success Metric: a validator session for a checkout whose root commits a
   `.ai-memory.toml` sees a private working directory holding exactly that
   file with the same bytes; without it, an empty one; a sealed working
   directory holding any other entry is refused.
5. Success Metric: the `deliver` and `review` guides, both copies of the
   Roundfix `deliver` and `review` references and the **Queue Owner** entry
   carry the new text; the documentation, skills and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published documents this Spec did not
write, read through Exa on 2026-11-12.

- The Go Programming Language Specification (https://go.dev/ref/spec),
  section "Selectors", rule 6: "If `x` is of interface type and has the value
  `nil`, calling or evaluating the method `x.f` causes a run-time panic." This
  is the panic of the source trace, `fmt.Fprintf({0x0, 0x0}, ...)`.
- The same specification, section "Handling panics": "While executing a
  function `F`, an explicit call to `panic` or a run-time panic terminates the
  execution of `F`. Any functions deferred by `F` are then executed as usual."
  This is why the dying owner released its record and status showed none.
- ai-memory, "Marker file: `.ai-memory.toml`"
  (https://github.com/akitaonrails/ai-memory/blob/main/docs/marker-file.md),
  section "What the marker file does NOT do": "A non-git directory outside
  `$HOME` needs a marker in its exact cwd." This is why the marker goes in the
  validator's working directory itself and not in a parent.

## Unreachable Acceptance

- criterion: a review validator session that runs after the release is recorded by ai-memory under this repository's project, and no new `roundfix-review-validator-<n>` project appears
  reason: it needs a real Agent session, the ai-memory server and the operator's machine, which this Spec's authorization keeps out of every test and QA row
  satisfied-by: the operator, after the release, lists `ai-memory list-projects` before and after the first delivery of this repository whose pre-PR review runs the validator, with the `.ai-memory.toml` the planning change commits at the root

## Validation

The decisions were checked on 2026-11-12 against these sources.

- Characterization at `d79b30b4`, with focused tests in a scratch copy:
  - `commandDeliveryEngineDependencies` builds the workflow without `log` and
    gives the engine `os.Stderr`; a gate whose capacity is held calls its wait
    callback, and `fmt.Fprintf` on the nil log panics with `invalid memory
    address or nil pointer dereference`;
  - `runDeliveryOwner` with an engine that panics, or that returns an error,
    while an item is at `gating` leaves owner PID `0`, and `deliver status`
    prints the item and no `Owner:` line;
  - `parseReviewFindingAnchor` reads none of `skills.md:34. Failure`,
    `x.md:34;`, `` `x.md:34`. ``, `x.md:34.` and `x.md:3-4.`;
  - the validator creates its directory with `os.MkdirTemp("",
    "roundfix-review-validator-")`, and `validateSealedPromptRequest` refuses
    any entry in a sealed working directory; Baseline analysis creates its own
    empty directories the same way.
- The owner's console log of the 0382 delivery and quaestor's trace name the
  same frame, `Gate.func1` at `deliver_workflow.go:1713`.
- A prototype of tasks 01 to 04 on a scratch copy at `d79b30b4`: the new tests
  and the eight Task Verification commands passed; `./internal/agent`,
  `./internal/store`, `./internal/delivery` and `./internal/testfixture`
  passed whole; the one whole run of `./internal/cli` failed only
  `TestHistorySanitizePlanTranscript`, whose suite guard caught the author
  editing the scratch tree during the run, and that test passed alone
  afterwards; each sabotage of the TechSpec made its named test fail.
- Secondbrain, read only: `wiki/index.md`, then `qmd query --all --files
  --min-score 0.3` for the owner panic at the capacity gate and for the
  validator's stray ai-memory projects; neither query returned a record above
  the threshold, so the decisions rest on the adopted sources, the maintainer's
  answers and the characterization.
- Exa: the three published passages above.

## Glossary

- changes: **Queue Owner**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BD brief for Spec 0385 of 2026-11-12, the maintainer's answers below and the adopted sources.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-12 ~17:05–17:20 BRT), verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o dono
  da fila não morre mais ao esperar a Machine Verification Capacity + âncora
  'path:line.' aceita, C = cache Go do roundfix fora do disco interno do mini
  (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B porque
  o disco do mini está abaixo dos 80 GB." (B is this Spec.)
- Question "Pantheon achou que o validador de revisão roda o Codex em
  /var/folders/…/roundfix-review-validator-<n>, criando projetos soltos no
  ai-memory. Já há rota em ~/.ai-memory.toml mandando /var/folders para o
  Roundfix. O que faço com o achado?" → "Juntar à Spec B" — option text: "Vira
  uma quarta Task da Spec B (fila e revisão)."

Taken from those answers, the sources and the characterization, with no
further interview:

- `type: fix`: a crash, a missing status line, a refused anchor and a misplaced
  session become correct behavior; the Release Ceiling is `patch`.
- The owner and the engine share one log passed by the constructor, and every
  workflow write goes through one accessor that falls back to standard error.
- The owner keeps its record when it stops with an item left, instead of
  status reporting a queue without a recorded owner (ADR-0279's rejected
  option stays rejected).
- Punctuation ends an anchor only before a space or the end.
- The validator copies the repository's own `.ai-memory.toml` into its
  working directory, the only place ai-memory reads a marker outside a
  checkout and outside the home directory; the agent's sealed rule admits that
  one regular file. Baseline analysis is unchanged.
- ADR-0395 records the decisions. Three code Tasks that each document their own
  guide (task_01 also revises the glossary), one Task for the Roundfix skill's
  references, and the QA gate.

## Open Questions

None.
