---
spec: 0184-baseline-plans-that-show-what-a-history-move-breaks
status: archived
created: 2026-09-29
surfaces: [backend, cli, docs]
archived: "2026-09-29"
source_slug: 0184-baseline-plans-that-show-what-a-history-move-breaks
qa_override: true
qa_override_approval: 'Maintainer standing authorization of 2026-08-09 (agent project memory ''qa_override por ambiente''): archive with override when the gate closes partial only because of environment, never with rows_blocked_finding above zero. Reconfirmed for this cycle on 2026-09-29 (''Autonomia ampla'').'
qa_override_reason: QA closed partial with rows_blocked_environment 2 and rows_blocked_finding 0. Row 4, the Fluxus outside-evidence replay, cannot run because the exported Fluxus Setup Manifest has no maintained transition to the go-cli-tui profile. Row 9 is the pre-PR Pull Request row, which records equivalent evidence. Every other row passed.
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 2423efe2af34390f6575e649e90dbcd94500091e
---


# Baseline plans that show what a History Relocation breaks

A Baseline Plan that brings an adopted repository to the current layout carries
its History Relocations in the same approval as the managed refresh. The plan
names each moved file, but not the citations the move breaks. On 2026-09-17,
Fiscus's `baseline update` proposed a one-entry managed refresh together with
six ADRs moving from the decision directory into the History Root. The adopter
knew those paths were cited in Specs, references and the glossary, but could
not tell which citations the move would break, so the whole plan stayed
unapplied. Fluxus shows the other outcome. One ADR consolidation moved sixteen
ADRs into history, and a later audit found sixteen relative links broken there,
one in each moved file. This Spec makes every plan that carries History
Relocations list, before approval, each tracked file whose citations stop
resolving, and where each citation points before and after the move.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a warning is keyed
  by the existing Baseline finding code and the citing file's
  repository-relative path. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — planning reads local Git state and
  tracked files only; no credential is read and no network call is added.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0173 (this Spec) records that a
  Baseline Plan reports the citations its History Relocations break, as
  warnings the Plan Digest binds, and rejects a deferral flag and a separate
  digest for now. ADR-0073 keeps apply one recoverable multi-file transaction,
  and this Spec adds no second transaction and no partial approval. ADR-0103
  makes an applied refresh republish the Setup Manifest so the update converges;
  a plan without History Relocations reads no additional file and reports
  nothing new, so a converged repository still reports `current`. ADR-0071 keeps
  plans portable and preimage-bound; citing files never become preimages, and
  `baseline apply` of a portable plan performs no scan. ADR-0068 keeps one
  confirmation-gated workflow, and the warnings reach its review unchanged.
  ADR-0120 places retired documentation under one History Root, which is where
  every relocation lands. ADR-0064 inventories the repository byte-exhaustively
  for readoption, and ADR-0070 audits every carrier but preserves root
  instructions; neither changes, and the scan writes nothing. ADR-0100 makes a
  managed refresh prove preservation instead of backing it up; citing files
  are not preimages, so that proof is unchanged. ADR-0074 gives repository
  rules hybrid semantic ownership, ADR-0078 moves confirmed root rules to
  semantic owners, and ADR-0075 lets profile divergence use a confirmed
  repository-owned adaptation, and ADR-0118 makes a repository constant that
  varies become a decision; this Spec classifies and moves no rule and adds no
  decision, so they do not apply. ADR-0150 makes work branches express purpose
  and Run branches express ownership; this Spec names no branch, so it does not
  apply. ADR-0097 carries a QA row forward only on declared, unmoved
  evidence, and ADR-0165 parks publication for a corrective Spec after a
  blocking review; neither governs Baseline planning, so they do not apply. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and
  ADR-0156. ADR-0093 checks Spec consistency by citation, and ADR-0094 makes
  that check artifact-presence-aware. ADR-0166, ADR-0167 and ADR-0168 (Spec
  0181) record undeclared Task paths, keep the pre-PR Pull Request row from
  deciding a qualifying partial, and give the related-ADR check a horizon; every
  Task here declares its paths and this Spec's gate aims at `pass`. ADR-0169 and
  ADR-0170 (Spec 0182) govern the pre-PR review base and Task Carry-Forward.
  ADR-0171 and ADR-0172 (Spec 0183) govern Run storage reporting. None of these
  governs Baseline planning, so they do not apply. ADR-0176 (Spec 0181) narrows
  only which Spec text the citation checks read, and this Spec does not depend
  on it, so it does not apply. All the others hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved this plan in chat on
  2026-09-29. Answering structured questions, they chose "Duas filas" (Onda 5
  delivers this Spec in a second queue) and, for this Backlog Entry, "Mostrar o
  impacto" (show each History Relocation's citation impact as warnings inside
  the same Plan Digest). The Roundfix skill files ride the standing grant of
  2026-09-18 for keeping the shipped skills true to the CLI, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Before approving a Baseline Plan, an adopter sees every tracked file whose
  citations the plan's History Relocations would break, with the line and the
  cited text.
- A citation that keeps resolving after the plan is never reported, including
  a relative link between two documents that move together.
- An approval covers the citation impact the adopter saw.
- Planning without History Relocations costs nothing new, and apply writes
  exactly what it wrote before.

## User Stories

1. As an adopter reviewing `roundfix baseline update`, I want the plan to name
   each file whose citations its History Relocations break, so that I can
   decide whether to approve, fix the citations first, or wait, without
   grepping the repository by hand.
2. As an Agent running `roundfix baseline plan` or `baseline update` in JSON, I
   want each broken citation as a typed warning, so that I can report it or
   repair the citing file before asking for approval.
3. As a maintainer, I want the digest I confirm to cover the impact I reviewed,
   so that a plan whose impact changed after my review is not applied under my
   earlier confirmation.

## Core Features

1. **A plan with History Relocations reports its Relocation Citations.** When
   a Baseline Plan carries at least one History Relocation, planning reads the
   repository's tracked text files. For each file, it reports the citations
   that resolve to an existing tracked file or directory today and would not
   resolve after the plan.
2. **What counts as a citation.** Two forms count. The first is a
   repository-relative path written anywhere in a tracked text file. The second
   is a Markdown link destination in a tracked Markdown file, in an inline link,
   an image or a reference definition. A link destination resolves from the
   citing file's directory, or from the repository root when it begins with
   `/`. Fragments and queries are ignored, and URLs with a scheme are never
   citations. After the plan, a citation resolves from the citing file's own
   post-plan location, so a relocated file's relative link to a document that
   stays behind is reported. A link between two documents that move together
   keeps resolving and is not reported.
3. **One warning per citing file.** Each citing file yields one warning. Its
   code is `baseline.history.citation`, its path is the citing file's current
   repository-relative path, and its message names up to three citations with
   their line, cited text, and resolution before and after the plan, plus how
   many more there are. Warnings are ordered by path. After two hundred citing
   files, one further warning with the code
   `baseline.history.citation.omitted` states how many citing files were not
   listed.
4. **The Plan Digest binds the warnings.** The warnings are part of the
   Baseline Plan, so the Plan Digest covers them. A citing file edited between
   the preview and `--confirm-plan` so that the impact changes produces a new
   digest, and the earlier confirmation does not apply.
5. **The scan reads only what it may.** Only paths in the Git index are read.
   Untracked and ignored files are never opened. A symbolic link is never
   followed, and nothing outside the repository root is read. Files that look
   binary are skipped. A tracked text file too large to scan, or unreadable, is
   counted in one `baseline.history.citation.unscanned` warning rather than
   skipped silently. A warning carries only paths, line numbers and cited path
   text, never other file content.
6. **Apply is unchanged.** Citing files are not preimages. The warnings never
   add, remove or alter a file change or a History Relocation, and `baseline
   apply` of a portable plan performs no scan. A plan without History
   Relocations performs no scan, and a current repository still reports
   `current`.

## User Experience

The text output of `baseline plan`, `baseline update` and the interactive
review lists each warning in the existing form:

```text
Warning: baseline.history.citation: docs/references/layout.md: line 80 cites ../adr/0056-layout.md, which resolves to docs/adr/0056-layout.md; after this plan nothing exists there (it moves to docs/history/adr/0056-layout.md)
Warning: baseline.history.citation: docs/adr/0041-old.md: line 12 cites 0070-new.md, which resolves to docs/adr/0070-new.md; after this plan it resolves to docs/history/adr/0070-new.md, where nothing exists
```

JSON output adds the same entries to the existing `warnings` array with no new
field and no schema version change.

## Non-Goals / Out of Scope

- Approving a managed refresh without its History Relocations, whether through
  a deferral flag recorded as a Setup Manifest decision or through a separate
  digest and confirmation for history moves. Both were rejected for now in
  ADR-0173, because of ADR-0073's single recoverable transaction and ADR-0103's
  convergence to `current`. A later Spec may reopen either if the reported
  impact proves insufficient.
- Rewriting, repairing or suggesting a replacement for any citation. The citing
  documents belong to the repository.
- Reporting citations that are already broken before the plan, citations in
  untracked files, links in non-Markdown syntaxes (HTML anchors, wiki links),
  and file-relative paths such as `../adr/x.md` written as plain prose rather
  than as a link. Repository-relative path tokens are reported in every
  scanned file.
- A standalone link checker, a new command or flag, or any change to the
  Setup Manifest, the plan schema or the result schema.
- Scanning when a plan has no History Relocations.

## Success Metrics

1. Success Metric: in a disposable adopted repository whose plan relocates a
   retired ADR, the plan reports every citing file, whether it cites the ADR by
   repository path, by relative Markdown link or by reference definition. It
   reports a relocated file's relative link to a document that stays behind,
   and does not report a relative link between two documents that move
   together, a citation already broken, or a citation in an untracked file.
2. Success Metric: replaying Fluxus commit `986f9dc` gives two results. The
   replay restores, in a disposable repository, the sixteen ADRs that commit
   moved into `docs/history/adr/` to their `docs/adr/` paths. First, each of
   the sixteen relocated ADRs yields exactly one `baseline.history.citation`
   warning, for its link to the ADR that consolidated it. These are the sixteen
   links Fluxus Spec 0067 reported broken in `docs/history/adr`. Second, every
   other citing file the plan reports matches a recount over the same History
   Relocations that is independent of Roundfix, with none missing and none
   extra. A recount made while authoring this Spec found fifteen: one active
   ADR and fourteen archived Spec files.
3. Success Metric: two plans of one repository that differ only in a citing
   file have different Plan Digests, identical file changes and History
   Relocations, and identical apply results.
4. Success Metric: on a clone of this repository (about 5,300 tracked files and
   46 MB) with one added History Relocation, `baseline update` completes within
   two seconds of the v0.20.0 binary's time on the same clone. A clone with no
   History Relocation reports the same result as v0.20.0.

## Decisions

- **Show the impact; keep one approval.** The maintainer chose this option on
  2026-09-29 over a deferral flag and a separate digest. See ADR-0173.
- **Report only what the plan breaks.** A citation counts when it resolves
  today and would not resolve after the plan, judged from where the citing file
  will be. This one criterion covers a citation of a moved document, a moved
  document's link to one that stays, and a pair that moves together, without
  special cases.
- **Bind the warnings into the digest.** The confirmation covers what the
  adopter reviewed. The cost is that an edit changing the impact requires a
  fresh review, which is the point.
- **Tracked files only, content never echoed.** The scan works from the Git
  index so that two clean clones report the same impact. It reads no untracked
  or ignored file, where secrets usually live, and reports only paths, line
  numbers and path text.

## Recorded limits

- A History Relocation whose destination is already occupied, which apply
  refuses and leaves in place, is not counted as moving, so its citations are
  not reported. A collision that only appears between planning and apply is
  outside what the plan can see.
- A citation is recognized only as an exact repository-relative path, or as a
  Markdown link destination. A bare path in prose that contains a space or a
  character outside `letters, digits and ._~/@+-` is not recognized; the same
  path written as a Markdown link destination is. A partial path in prose, an HTML anchor or a
  generated link is not recognized.
- When a tracked file has uncommitted edits, the scan reads its working-tree
  bytes, as the rest of planning does.
- A directory counts as a cited target only when it lies inside a relocated
  unit, such as one legacy Spec folder or Review Artifact. A family or legacy
  root, such as the decision directory, is never reported. Guidance describes
  those roots by convention, and they are expected to exist again.

## Acceptance evidence

Each Core Feature requires positive and negative evidence. The negative cases
carry the weight. Each of these would pass a happy-path test:

- a relative link between co-moved documents reported as broken;
- an already-broken citation reported;
- an untracked file read or reported;
- a symbolic link followed;
- a plan without relocations that scans;
- a digest that ignores a changed citing file;
- an apply that differs because of a warning.

The outside-evidence rows rest on sources this Spec did not produce:

- **The Fluxus repository**, read-only at commit `986f9dc`, which moved sixteen
  superseded ADRs into `docs/history/adr/`. Fluxus Spec 0067 later measured
  sixteen broken relative links there. It is mirrored in the Secondbrain at
  `projects/fluxus/mirror/docs/history/specs/0067-os-links-de-documentacao-resolvem/_prd.md`.
  When the Fluxus repository or that commit is unavailable, the row is recorded
  as blocked with its reason.
- **The Fiscus evidence**, captured on 2026-09-17 in the Secondbrain entry
  `inbox/roundfix/_triaged/2026-09-17-um-digest-so-prende-refresh-gerenciado-e-mudanca-de-documento-do-repositorio.md`.
  The Fiscus mirror as of 2026-09-28 still holds the six superseded ADRs in
  `docs/adr/`. A live reference document links one of them twice with a
  relative link, and an archived QA report cites another by repository path.
- **The Markdown link forms**, from the VS Code documentation for updating
  links on file moves. It names relative links, root-relative links, fragments,
  images and reference definitions as the link forms a move affects
  (<https://github.com/microsoft/vscode-docs/blob/main/docs/languages/markdown.md>,
  <https://github.com/microsoft/vscode/issues/158416>).

## Research basis

The adopted Backlog Entry is indexed in
[references/_index.md](references/_index.md).

**Secondbrain.** Consulted through `wiki/index.md` and
`qmd query "mover arquivo quebra citação por caminho links relativos markdown detectar referências antes de mover"`.
It returned the Fluxus Spec 0067 mirror and the Fiscus mirror used above. A
recount run for this Spec at Fluxus `986f9dc` found exactly sixteen relative
links that resolved before the move and not after it, one in each moved file.
It found fifteen more files citing the moved ADRs: one active ADR by relative
link, and fourteen archived Spec files by repository path.
That result grounds Core Feature 2's criterion, which judges a citation from
the citing file's post-plan location.

**Exa.** Located the VS Code link-update documentation. It confirms which link
forms a file move affects, and that editors treat link repair on move as opt-in
("never" by default), which supports reporting instead of rewriting.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
