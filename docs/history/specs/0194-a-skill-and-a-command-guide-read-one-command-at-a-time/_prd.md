---
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: archived
created: 2026-09-30
surfaces: [backend, docs]
archived: "2026-10-01"
source_slug: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
---


# A skill and a command guide read one command at a time

The Roundfix Skill is one file, `.agents/skills/roundfix/SKILL.md`: 141,699
bytes, 2,670 lines and 33 sections on 2026-09-30. The command reference,
`docs/user-guide/commands.md`, is 79,467 bytes. Two costs were measured in the
two days before this Spec:

- **Every reader loads everything.** An Agent that needs one command loads the
  whole skill. Each Spec authoring session of 2026-09-29 and 09-30 used between
  400,000 and 750,000 tokens.
- **Every CLI Task collides.** Almost every Task that changes a command
  declares both files, so `SC-WAVE-COLLISION` forces those Tasks into a serial
  chain. Spec 0187 is a four-Task chain for this reason alone. Three QA gates
  refused because a Task wrote into a section all commands share: Spec 0181
  twice and Spec 0187 once, each on `TestSettlementGuidanceIsOneTable`.

This Spec splits both documents by command without changing a word. It is the
last Spec of Wave 7, so it splits whatever content the earlier Specs of the
wave left in the two files.

## Project Constraints

- Identifier strategy: not applicable — no identifier is added; the change
  moves Markdown and adds one reader. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; no credential
  and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0187 (this Spec) makes both
  documents an entry file plus one file per command family. ADR-0130 keeps
  every bounded path governed, so the grant lists the canonical reference
  files and not their ordinary mirror copies. ADR-0081 makes the digest
  regeneration that follows a skill edit sanctioned fallout, and ADR-0149 has
  the grant name each regeneration command while the tree names its outputs,
  which is how the mirror copies of the reference files are covered. ADR-0166
  records a Task's undeclared paths, and every Task here declares its paths.
  ADR-0167
  keeps the pre-PR Pull Request row from deciding a qualifying partial; this
  Spec's gate aims at `pass`. ADR-0176 reads citations only from authored
  text. This Spec's gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096,
  ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and ADR-0093 and ADR-0094 check
  its consistency by citation and artifact presence. ADR-0097 cites ADR-0080
  but carries a QA row forward, ADR-0168 cites ADR-0093 but narrows the
  related-ADR check, and ADR-0179 cites ADR-0130 but grants operations through
  an empty list; this Spec changes none of them and its grant lists paths, so
  they do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer authorized on 2026-09-30
  adjusting every skill ("considere autorizado a ajustar todas as skills se
  necessário") and stated that architecture and structures may change freely;
  the governed test files ride the standing grant of 2026-09-21 for governed
  paths a slice genuinely needs. Recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/release.md`,
  `.agents/skills/roundfix/references/review.md`,
  `.agents/skills/roundfix/references/review-runs.md`,
  `.agents/skills/roundfix/references/runs.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/references/spec-delivery.md`,
  `.agents/skills/roundfix/references/stop.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `internal/docscontract/publicdocs_test.go`, `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go`,
  `skills/baseline_skill_contract_test.go`. Sanctioned regeneration:
  `make skills-sync`, `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- An Agent that needs one command reads the skill's entry file and one
  reference file, not the whole skill.
- Two Tasks that change different commands can run in the same Wave.
- The split loses no line and rewords none.
- Every contract that read either document still holds, with the same
  strength.

## Core Features

1. **The Roundfix Skill is an entry file plus one reference per command
   family.** `SKILL.md` keeps the introduction, the contract of an assigned
   Batch, the evidence boundaries, the shared `### QA settlement` section and
   an index that names each file under `references/` and when to read it.
   Every other section moves, byte for byte, into the reference file of its
   command family (ADR-0187).
2. **The command reference is an index plus one file per command.**
   `docs/user-guide/commands.md` keeps its introduction, the global contract,
   its group headings and an index. Each command's section moves, byte for
   byte, into `docs/user-guide/commands/<command>.md`. Only relative link
   targets change, so that each link still reaches the same document.
3. **A contract reads the entry file with its companion files as one text.**
   Every check that pinned a phrase in either document, including
   `roundfix skills check`, reads the entry file followed by its companion
   files. A phrase may live in any of them, and a forbidden phrase is refused
   in all of them.
4. **A Task declares the one command file it changes.** The `write-tasks`
   skill and its Task template say so, and state the consequence: two Tasks
   that declare the same file cannot share a Wave.

## Non-Goals / Out of Scope

- Rewording, condensing or deleting any prose of the skill or the guide. A
  later wave slims the authoring skills.
- A generator, a build step or any change to the `Makefile`. The existing
  mirror copy already carries nested files.
- Changing a command's behavior, flags, output or exit codes.
- Splitting any other skill, or editing the setup-owned guides under
  `docs/agents/`.
- Making `SC-CLI-UNDOCUMENTED` require a specific command file. It already
  accepts any file under the guide roots.
- Keeping the old `commands.md#<command>` anchors alive for readers outside
  this repository.

## Success Metrics

1. The skill's `SKILL.md` is at most 20,000 bytes, its index names every file
   under `references/`, and every index link resolves.
2. Every non-blank line of the skill's body and of the command reference, as
   they stood on the splitting Task's starting commit, is present with the
   same count in the new files. For the command reference the comparison
   ignores link targets.
3. Every test that pinned text of either document passes while reading the
   entry file with its companion files, `TestSettlementGuidanceIsOneTable`
   passes unchanged, and `roundfix skills check` reports no diagnostic.
4. Two Tasks in one Wave that declare different command files raise no
   `SC-WAVE-COLLISION`, and two that declare the same command file still do.

## Recorded limits

- A moved section keeps its heading level, so a reference file starts with a
  second-level heading. Correcting the levels would reword the file and is
  left to a later change.
- An adopting repository keeps the old single-file skill until its Roundfix
  skill is installed or restored again. Both layouts are valid skills.
- An Agent loads a reference through its own file read. Roundfix does not
  inject one.
- A single file is easier to paste or print than a directory. That is given up
  for conditional loading.

## Decisions

- **No generator.** The files are the source. A generated monolith would keep
  one shared file that every Task edits.
- **Move, then nothing else.** The split is provable only if it rewords
  nothing, so every improvement to the text waits for a later change.
- **Pins read the whole document.** Moving a phrase between files of one
  document must not break or weaken a contract.
- See ADR-0187.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence rows rest on sources this Spec did not produce:

- Published guidance for agent skills. Microsoft's Agent Skills documentation
  (<https://learn.microsoft.com/en-us/agent-framework/agents/skills>) states
  "Keep `SKILL.md` under 500 lines and move detailed reference material to
  separate files", describes reference files "loaded on demand", and
  recommends under 5,000 tokens for the loaded body. The Roundfix skill was
  2,670 lines.
- This repository's recorded history: Spec 0187's Task Graph is a serial chain
  of four Tasks that each declare the Roundfix skill, and the QA reports of
  Specs 0181 and 0187 record three refusals on
  `TestSettlementGuidanceIsOneTable`.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "progressive disclosure agent skills SKILL.md references carregadas sob demanda tamanho"`.

- `wiki/concepts/arquitetura-de-instrucoes-e-progressive-disclosure.md`
  records the three loading tiers of a skill (name and description, body,
  reference files) and that a short entry file is what keeps context small.
- `inbox/secondbrain/2026-09-30-orquestradores-e-ecossistema-jev-o-que-se-transfere-ao-roundfix.md`
  records that comparable projects keep skills between 1 and 14 KB with
  references loaded on demand. It informed the choice of one file per command.
- `inbox/secondbrain/2026-09-30-modelos-de-codificacao-ago-set-2026-e-runtimes-cursor-grok.md`
  was read and adds nothing to this design.

Exa returned the Microsoft documentation above and three community authoring
guides. One of them argues the opposite for a skill that documents a CLI: keep
one file, so it can be printed by a subcommand. It also sets a hard ceiling of
50,000 characters, past which a skill must be condensed or split. The Roundfix
skill is almost three times that ceiling, the maintainer ruled out rewording in
this wave, and each reader needs one command, so the split is the option left.
That source changed one decision: the entry file keeps everything an assigned
Agent needs on every run, so the common path costs no extra read.

## Technical candidate

The [_techspec.md](_techspec.md) records the file layout, the move rules, the
coverage and the build order.
