---
spec: 0388-a-completion-notice-skips-a-record-the-merge-adopted
status: active
created: 2026-11-13
surfaces: [docs]
type: fix
---

# A completion notice skips a record the merge adopted

The `## Completion notices` section of the roundfix skill's release reference,
shipped in v0.143.0 with skill 0.1.228, tells the session that follows a merge
to send a completion notice for every Backlog Entry or Finding that
`git diff --name-only --diff-filter=D <merge commit>^ <merge commit> -- docs/backlog docs/findings`
lists. That listing also names every record the merge adopted, because the
adoption contract moves the record into a Spec's `references/` folder and a
diff limited to `docs/backlog` and `docs/findings` never pairs the move with
its destination. A Notice Request adopted by a planning merge would get its one
notice early, as if its work were closed, and a second one when the Spec ships.
After this Spec the section sets aside every record the merge adopted, tells
how to recognize one with Git alone, asks the user when the evidence
disagrees, restores the two guarantees the 0386 rewrites lost, and pins every
guarantee of the section with its own contract-test phrase.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md`](references/2026-11-13-a-completion-notice-reads-an-adopted-record-as-closed.md),
   recorded by the operator from ai-memory message `01a1d483` of 2026-11-13;
   reported by argus. Argus's pre-PR review of its Baseline refresh to
   Roundfix 0.143.0 found the closure rule at `release.md:145-153`; the
   operator confirmed the cause in the 0386 item commits `3ed725d7` (round 1)
   and `a7c2b447` (round 2). The message asks for no reply and no notice, so
   this entry carries no Notice Request.

At the authoring base the other open Backlog Entries are the retry defects
that Spec 0389 adopts, the Conexus status report and the Homebrew tap;
`docs/findings/` holds none. None of them touches
completion notices, so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `f953f77f` (v0.143.0 plus Backlog Entries
#856 and #857) and rebased on `753ac77d` (#858, which closed the stray
ai-memory projects entry and touches no file of this Spec). The queue delivers this Spec first and Spec 0389 after it. The
only files the two share are the `roundfix` skill's `SKILL.md` pair and
`skills/testdata/owned-skill-versions.json`: both raise the `roundfix` skill
version. task_01 writes "the next version" with the repository's record command
on the tree it runs on and never types a number, so whichever Spec runs second
raises past the first. Spec 0389 owns `references/deliver.md`; this Spec
changes only the `## Completion notices` section of `references/release.md`
and its contract test.

## Project Constraints

- Identifier strategy: not applicable — no identifier, record or schema is
  minted; the skill keeps its name and gains a version. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the change is skill text and tests that read the
  repository's own files and a temporary Git repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0396 applies unchanged: the request
  lives in its record, a Backlog Entry, a Finding or the adopting PRD's
  `## Sources`, and earns one notice once its work is delivered; this Spec
  makes the closure rule read that decision for an adopted record, and the
  Decision's wording is not affected, so no amendment is due. ADR-0383,
  ADR-0390, ADR-0330 and ADR-0319, which ADR-0396 cites, do not apply: this
  Spec changes no message clause, precedence, follow-up order or inbox rule.
  Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answers of 2026-11-13 quoted in the authorization record, which keep the
  grant of every owned skill; the Governed Paths are the `roundfix` skill's
  `SKILL.md` pair and its canonical release reference. No Makefile, workflow,
  Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0388-a-completion-notice-skips-a-record-the-merge-adopted/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/release.md`, `skills/roundfix/SKILL.md`.

## Goals

- The closure rule sends no notice for a Backlog Entry or Finding the merge
  adopted; its request follows the Spec that adopted it, and only a record the
  merge closed without a Spec sends one.
- The section says how to recognize an adopted record with Git alone, from what
  the merge adds under the Spec root, and asks the user when that evidence is
  unclear instead of sending.
- Every guarantee the section states is pinned by its own contract-test phrase,
  including the two the 0386 rewrites lost, and the Git commands it prescribes
  are run against a temporary repository that holds an adopted, a closed and an
  unclear record.
- The `roundfix` skill's version rises and is recorded, so adopters on 0.143.0
  receive the correction through their next `roundfix baseline update`.

## Core Features

1. The closure rule. The closure bullet sets aside each record the merge
   adopted; a new paragraph after the list says why the deletion listing names
   adopted records, why the record at `<merge commit>^` cannot tell, the two Git
   commands that find the adoption, the frontmatter that confirms it, the
   `## Sources` route that round 1 of 0386 stated, and the ask-the-user case.
2. The restored channel. The send paragraph names `memory_message_send` again,
   which the first version of the section stated and round 1 dropped.
3. The pins. `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest`
   pins every guarantee of the section outside the adoption rule, and the new
   `TestACompletionNoticeSkipsARecordTheMergeAdopted` pins the adoption rule's
   sentences and runs the section's own Git commands on a temporary repository
   under three rename settings, for a squashed merge and a merge commit.
4. The skill version, raised and recorded by the record command.

## Non-Goals / Out of Scope

- A Roundfix command that finds or sends completion notices. ADR-0396:
  "No Roundfix command reaches ai-memory, so only a session can send the
  notice."
- A Spec Root that lies in another Git repository (a knowledge workspace whose
  `docs/specs` points elsewhere): the merge of this repository then adds no
  copy under the Spec Root, and the rule falls back to the PRD's `## Sources`
  and, failing that, sends. No adopter of 0.143.0 is known to use that layout.
- The release and no-release events, the delivery rule for several Specs, the
  grouping by message, the cancel rule, the outbox check and `references/modes.md`
  steps 4 and 7: they stay as v0.143.0 wrote them, only pinned.
- Any notice due for work planned before this Spec: the operator sends those on
  the maintainer's answer.

## Success Metrics

1. Success Metric: in both copies of `references/release.md`, the closure
   bullet and the paragraph after it read as `_techspec.md` → Interfaces writes
   them, the send paragraph names `memory_message_send`, the two copies are
   byte-identical, and nothing outside `## Completion notices` changes.
2. Success Metric: `TestACompletionNoticeSkipsARecordTheMergeAdopted` passes
   for `canonical` and `embedded`, each with `squash` and `merge`: under
   `diff.renames` `false`, `true` and `copies`, the section's commands classify
   the two adopted Backlog Entries and the adopted Finding as adopted, the
   declined entry as closed and the entry whose copy records no adoption as one
   to ask about.
3. Success Metric: `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest`
   passes with one pinned phrase for each guarantee the TechSpec's guarantee
   inventory lists, and each sabotage of `_techspec.md` → Sabotages fails its
   named test.
4. Success Metric: the `roundfix` skill carries its next version in both
   `SKILL.md` copies, recorded in `skills/testdata/owned-skill-versions.json`,
   and `go test ./skills` passes.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- This repository's own history, which no Spec of this cycle wrote: the
  planning merges `73421db2` (#842), `a99c17b2`, `dfb4dbac` and `62bc941a`
  each adopted a Backlog Entry; the v0.143.0 closure command lists each adopted
  entry as deleted, while `git diff --name-status` over the same commits without
  a path limit reports the same file as `R094` to `R096`, a rename into
  `docs/specs/<slug>/references/`, and the entry at `<merge commit>^` reads
  `status: open`.
- Linus Torvalds on the Git mailing list, "Re: [PATCH] Diff updates, fixing
  pathspec and rename/copy interaction", 2005-05-27
  (https://public-inbox.org/git/Pine.LNX.4.58.0505270848220.17402@ppc970.osdl.org/),
  read through Exa on 2026-11-13: "if you limit yourself to a directory, and
  there was a rename that moved a file from outside that directory into it,
  then it is not a rename at all, it is a "create". The stuff outside the
  pathspec limit simply doesn't exist." Git's gitdiffcore documentation states
  the same: "The pathspecs are used to limit the world diff operates in."
- Argus's report, under `references/`: a measured finding of another
  repository's `roundfix review`, not a premise this Spec designed.

## Unreachable Acceptance

- criterion: after the v0.143.1 release, an adopter's `roundfix baseline
  update` installs the `roundfix` skill at its new version, and its release
  reference sets aside a record the merge adopted.
  reason: it needs a released binary and another repository's managed
  refresh, which no Verification or QA row may run or read.
  satisfied-by: the operator's or Argus's first managed refresh after the
  release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-13 against these sources.

- Characterization at `f953f77f`: the section is the end of
  `.agents/skills/roundfix/references/release.md` (line 117 on), identical to
  the shipped copy; the closure bullet is lines 145 to 153.
  `TestTheRoundfixSkillSendsOneCompletionNoticePerNoticeRequest` pins 20
  phrases and none of them excludes an adopted record. `roundfix` is at 0.1.228.
- The 0386 history: the first version (`3ed725d7^`) sent "with
  `memory_message_send`"; round 1 (`3ed725d7`) dropped that and added "A record
  that a PRD's `## Sources` names was adopted by that Spec, and its notice
  follows the Spec instead."; round 2 (`a7c2b447`) dropped that sentence. Round
  2 also replaced round 1's notice for a cancelled Spec with "sends no notice",
  because the cancel returns the sources; that change was deliberate and stays.
- The four planning merges under Acceptance evidence, read with
  `git diff --name-only --diff-filter=D`, `git diff --name-status` and
  `git show <merge commit>^:<path>`: the defect is not rare rename-pairing
  failure, as the Backlog Entry supposed, but every adoption. The pre-merge
  frontmatter cannot discriminate, so the rule reads the copy the merge adds.
- A prototype on a worktree at `f953f77f`: the section text of
  `_techspec.md` → Interfaces, `make skills-sync`, the record command and the
  test changes; `go test ./skills`, `go test -tags repocontract ./skills`,
  `go test -tags docscontract ./internal/docscontract` and
  `GOOS=windows go vet ./skills` passed, the Behavior Surface Record did not
  move, and the sabotages of `_techspec.md` failed their named tests.
- Secondbrain, read only: `wiki/index.md` (entries `argus` and `roundfix`) and
  `qmd query "completion notice adopted backlog entry roundfix release" --all
  --files --min-score 0.3`, which returned no result. Nothing there changes the
  decision.
- Exa: the mailing-list passage under Acceptance evidence and a 2019 reply by
  Jeff King on the same list ("we limited our tree diff to just look at "b";
  hence we don't even consider "a" a source"), which confirm that no rename
  setting pairs a move across the path limit.

## Glossary

None.

## Decisions

The maintainer's answers (AskUserQuestion, 2026-11-13 ~08:20 BRT), verbatim:

- Question: "Quais trabalhos do backlog do Roundfix eu planejo e entrego
  agora, sem supervisão? (Cada item marcado vira Spec com plano validado e é
  entregue pela fila. Os desmarcados ficam no backlog.)" — selected "Aviso de
  conclusão (Recommended)"; option text: "Spec corretiva da regra publicada na
  v0.143.0: a regra lê como fechado um registro que uma Spec adotou e pode
  mandar o aviso cedo e de novo (achado da Argus; a falha foi minha na rodada
  2 da 0386)."
- Question: "Quando publico as releases desses trabalhos?" — answer: "Patch
  após cada Spec (Recommended)"; option text: "v0.143.1 assim que sair a
  correção do aviso; a próxima patch assim que sair o retry. Os adotantes que
  já trouxeram a 0.143.0 (a Argus, por exemplo) recebem logo a correção do
  texto."
- Question: "Quanta folga eu tenho para consertar falhas sem parar? (Na 0386 o
  review usou 3 correções.)" — answer: "Até 3 por Spec + Makefile/CI
  (Recommended)"; option text: "Até 3 correções por Spec, contando as rodadas
  do review. Fica liberado mexer em Makefile, workflows do CI, go.mod e lint
  quando a correção exigir; cada mudança vai descrita no relatório."

Autonomous route: no person present; decisions taken from the operator's cycle BE brief for Spec 0388 of 2026-11-13, the answers above and the characterization above.

Taken from those sources, with no further interview:

- An adopted record is recognized by the copy the same merge adds under a
  Spec's `references/` folder, confirmed by that copy's frontmatter, or by a
  PRD's `## Sources` that names it; not by the frontmatter at
  `<merge commit>^`, which the characterization shows reads `open` or
  `pending` for an adopted and a closed record alike.
- When a copy exists but its frontmatter does not confirm the adoption and no
  PRD names the record, the session sends nothing and asks the user: a notice
  sent early cannot be taken back once the recipient consumed it, and a missed
  one can still be sent on the answer.
- Keep the v0.143.0 listing command, so the deletion listing stays the same and
  only the classification after it changes.
- `type: fix`; the Release Ceiling is `patch` (v0.143.1). One implementation
  Task and a QA gate; no ADR, because ADR-0396's decision already keeps an
  adopted request in the adopting PRD's `## Sources`.

## Open Questions

None.
