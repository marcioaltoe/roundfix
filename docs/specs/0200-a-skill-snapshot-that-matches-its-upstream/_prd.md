---
spec: 0200-a-skill-snapshot-that-matches-its-upstream
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# A skill snapshot that matches its upstream

The Baseline ships setup snapshots that name the upstream skills a profile
installs, pinned to one commit of `marcioaltoe/skills`. The pin is `14fdf46`,
from 2026-07-25; upstream is at `a4e18e4`. Upstream renamed `context7` to
`context7-cli`, `feature-systems-pattern` to `app-renderer-systems` and `rust`
to `rust-expert`, and moved `bubbletea` and `tui-design` out of `go-cli` into a
new `go-tui` list. The catalog, the rendered guides and one universal
capability still use the old names, so a new adopter is told to restore a
skill no snapshot can supply. The asset sync cannot take the refresh: run with
`--check` against the upstream checkout it exits `2` with
`catalog.profile.skill.outside-setup` for `bubbletea`, `tui-design`, `context7`,
`rust` and `feature-systems-pattern`. The `go` module dispatches three skills
no Roundfix setup lists, because the catalog checks only required skills.

This Spec refreshes the four setup snapshots to `a4e18e4`, follows every rename
by name, gives the Go CLI/TUI profile the `go-tui` setup, and adds the catalog
check that would have caught the gap. It also moves this repository to the
state an adopter reaches after its next update. See ADR-0191.

## Project Constraints

- Identifier strategy: applicable — three dispatch trigger identifiers follow
  their renamed skill in the existing `trigger.<module>.<skill>` form
  (`trigger.core.context7-cli`, `trigger.rust.rust-expert`,
  `trigger.typescript.app-renderer-systems`), one trigger is added
  (`trigger.core.exa-web-search`), one catalog diagnostic code is added
  (`catalog.profile.skill.dispatch-outside-setup`), one setup snapshot
  identifier is added (`go-tui`), and the capability probe of
  `capability.context7` gains a `priorSkills` list. The capability identifier
  `capability.context7` is kept. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded assets, local files and
  Git only. The Tasks read a local checkout of the upstream skills repository
  and fall back to cloning it; no Verification command opens a network
  connection and no credential is read. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0191 (this Spec) governs the
  snapshot, the rename, the profile's setup, the dispatch check and the prior
  skill names. ADR-0072 keeps the maintained Python parity contract and asks
  for a new behavior to be "recorded explicitly as a designed delta rather than
  being mistaken for parity"; the sync's validation change and the new setup
  row are recorded that way. ADR-0073 and ADR-0103 require the managed refresh
  to apply in one transaction and converge to `current`. ADR-0081 and ADR-0149
  keep the sanctioned regeneration outputs inside the grant, and ADR-0130 has
  the audit judge only Governed Paths. ADR-0058 and ADR-0060 stay satisfied
  because no clause is added, removed or reworded, and ADR-0099 because no
  clause identifier changes. ADR-0067 keeps repository-owned profiles on the
  built-in profile they derive from, so they follow its new setup. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155,
  ADR-0156 and ADR-0166, and ADR-0093 and ADR-0094 check the Spec by citation
  and artifact presence. ADR-0178 and ADR-0179 decide how the audit reads the
  grant. ADR-0180 to ADR-0184 are present in this tree and do not apply: they
  decide model profiles, Settlement Checks, Claim Receipts and Surface
  Transcripts, none of which this Spec touches, and the guide that would bind
  ADR-0183 and ADR-0184 is not in this tree. ADR-0074 (hybrid semantic ownership) cites ADR-0067, ADR-0097 (QA row
  carry-forward) and ADR-0167 (the pre-PR Pull Request row) cite ADR-0080,
  ADR-0168 (the related-ADR gap) and ADR-0176 (citation checks read authored
  text) cite ADR-0093, ADR-0173 (citations a History Relocation breaks) cites
  ADR-0073, ADR-0181 (comparison with the Recommended Profile) cites ADR-0180,
  and ADR-0182 (Settlement Checks) cites ADR-0096; none applies, because this
  Spec changes none of their behaviors. ADR-0177 is reached only through
  ADR-0173 and is equally untouched. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source, its generated guides and every skill
  ("Autorizar os dois"; "considere autorizado a ajustar todas as skills se
  necessário"), and the program decision that this Spec refreshes the snapshot
  to `a4e18e4` with the renames, the `go-tui` setup and the dispatch check,
  recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/go-tui.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/go.json`,
  `internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/modules/typescript.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/derived_ownership_test.go`,
  `docs/agents/skill-dispatch.md`, `docs/agents/setup-context.json`,
  `skills/baseline_skill_contract_test.go`,
  `.agents/skills/context7/SKILL.md`,
  `.agents/skills/context7-cli/SKILL.md`,
  `.agents/skills/context7-cli/references/docs.md`,
  `.agents/skills/context7-cli/references/setup.md`,
  `.agents/skills/context7-cli/references/skills.md`,
  `.agents/skills/golang-dependency-management/SKILL.md`,
  `.agents/skills/golang-safety/SKILL.md`,
  `.agents/skills/golang-safety/references/nil-safety.md`,
  `.agents/skills/golang-safety/references/slice-map-safety.md`,
  `.agents/skills/golang-structs-interfaces/SKILL.md`,
  `.agents/skills/golang-structs-interfaces/references/struct-fields.md`,
  `.agents/skills/golang-structs-interfaces/references/type-assertions.md`.
  Sanctioned regeneration: `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- Every setup snapshot names the skills its upstream list names at one
  commit, and a refresh that follows an upstream rename is one sync run.
- No catalog entry, rendered guide or remediation names a skill that the
  profile's setup cannot supply.
- The catalog refuses a profile whose guides would dispatch a skill outside its
  setup.
- An adopter's next `baseline update` installs the renamed skills, and a
  second update reports `current`.

## User Stories

1. As a Roundfix maintainer, I want the asset sync to accept a refresh that
   repairs a rename, so that following upstream is one command and not a
   sequence of hand edits that each leave the catalog invalid.
2. As an adopter's Agent, I want every skill `docs/agents/skill-dispatch.md`
   names to be a skill the repository can install, so that a trigger never
   points at nothing.
3. As an adopter, I want a repository that still holds the old `context7`
   skill to keep passing profile alignment, so that the update that installs
   `context7-cli` is not refused first.
4. As a Roundfix maintainer, I want this repository to hold what an adopter
   holds after the update, so that Doctor's `skills:` line here tests the
   shipped catalog.

## Core Features

1. **The sync validates what it produces.** The asset sync checks the catalog
   it would write, with the refreshed snapshots in place. A catalog that is
   invalid only because its snapshots are stale is refreshed; a catalog the
   refresh leaves invalid is refused as today, and a catalog with no drift that
   is invalid is refused as today.
2. **The snapshots follow upstream.** The four snapshots (`go-cli`, `go-tui`,
   `rust-cli`, `typescript-bun`) pin `a4e18e4` and list exactly the upstream
   entries plus the Roundfix-owned entry the upstream list omits. The Go
   CLI/TUI profile takes `go-tui`. The `core`, `rust` and `typescript` modules
   name the renamed skills, and the `go` module requires the three skills it
   already dispatches.
3. **The catalog refuses a dispatch outside the setup.** A profile whose
   selected modules dispatch a skill, or whose owned activation bundles name a
   skill, that its setup does not list fails catalog validation. The `core`
   module requires and dispatches `exa-web-search`, the skill a universal
   required capability already probes.
4. **Context7 evidence follows the new name.** The universal Context7
   capability probes `context7-cli`, still accepts an installed `context7`, and
   remediates with `--skill context7-cli`, which every built-in profile's setup
   can restore.
5. **This repository takes its own update.** This repository's managed guides
   are refreshed, its required upstream skills are restored through the public
   update, the obsolete lock entry is removed through `baseline skills
   reconcile`, and the old skill directory is deleted. The user guide explains
   the same three steps to adopters, and the documented Doctor line matches
   this repository.

## User Experience

A maintainer refreshes the snapshots with
`roundfix baseline assets sync --source-dir <checkout>/setups` after editing
the modules to the new names. An adopter runs `roundfix baseline update --yes`:
`docs/agents/skill-dispatch.md` names the new skills, the skills stage restores
them, and the old directory stays until the adopter runs `roundfix baseline
skills reconcile` and deletes it. The user guide section on Repository Skill
Set restoration shows those commands.

## Non-Goals / Out of Scope

- Editing any upstream skill. Upstream skills are vendored third-party content;
  Spec 0199 states that repository rules take precedence over skill defaults.
- Requiring any skill an upstream list adds, other than the three `go` module
  skills and `exa-web-search`. `typesafe-ai`, the new Go skills, the new web
  quality skills and the Firecrawl variants join the snapshots and are
  restorable, and no module requires them.
- Refreshing a required upstream skill whose installed bytes match its lock
  entry but not the snapshot. `baseline update` compares with the lock, and
  ten of this repository's installed upstream skills are in that state. It is
  recorded as the open Backlog Entry
  [2026-09-30-an-installed-skill-that-matches-its-lock-can-trail-the-snapshot.md](../../backlog/2026-09-30-an-installed-skill-that-matches-its-lock-can-trail-the-snapshot.md).
- Deleting an installed skill tree from Roundfix. The old directory is the
  repository's to delete.
- The dispatch trigger of `cut-release`, a skill only a person can start. Its
  Backlog Entry stays open.
- The Roundfix skill's setup membership and the owned-skill minimum, which
  Spec 0195 owns, and the stack clause wording, which Spec 0199 owns.
- Any Makefile, lint, formatter, CI or `go.mod` change.

## Success Metrics

1. Success Metric: an asset sync from a source whose skill was renamed, run
   against a catalog already naming the new skill, refreshes the snapshots and
   leaves a valid catalog; the same sync against a catalog naming a skill no
   source provides is refused and writes nothing.
2. Success Metric: no module, activation bundle or setup names `context7`,
   `feature-systems-pattern` or `rust`; every snapshot pins `a4e18e4`; the
   `go-cli-tui` profile's setup is `go-tui` and lists `bubbletea` and
   `tui-design`.
3. Success Metric: every module requires every skill it dispatches, and a
   catalog in which a profile's modules or owned bundles name a skill outside
   its setup reports `catalog.profile.skill.dispatch-outside-setup`.
4. Success Metric: a repository with only `context7-cli` or only `context7`
   installed satisfies `capability.context7`; one with neither is blocked with
   a remediation naming `--skill context7-cli`; that skill and `exa-web-search`
   are required by `core` and listed by every built-in profile's setup.
5. Success Metric: in this repository every required external skill is
   installed with its lock hash, no lock entry or directory names `context7`,
   Doctor reports `skills: ok (42 required: 14 Roundfix-owned, 28 external)`,
   and a second `roundfix baseline update --repo .` reports `current`.
6. Success Metric: after the sanctioned regeneration and this repository's
   managed refresh, a second refresh reports no file change; the Source
   Baseline corpus, the retention transitions and
   `lock-hash-compatibility-v1.json` are byte-identical.

## Declared breaks

- **Break — asset sync.** A sync whose current catalog is invalid but whose
  refreshed catalog is valid now succeeds; it exited `2` with "Go-owned
  canonical Baseline assets are invalid". When both the source and the catalog
  are invalid, the source error is now reported instead of the catalog error.
- **Break — catalog.** A catalog whose profile dispatches a skill outside its
  setup, or whose `go` module dispatches a skill it does not require, no
  longer loads. The embedded catalog satisfies both after this Spec.
- **Break — guides.** Every adopter's `docs/agents/skill-dispatch.md` changes
  on its next update: renamed skills, the Go skills of `go-cli-tui`, and the
  new `exa-web-search` trigger.
- **Break — required skills.** Go CLI/TUI adopters require three more Go
  skills and every adopter requires `exa-web-search`. Doctor's `skills:` line
  fails until the update restores them. `exa-web-search` was already needed to
  pass profile alignment.
- **Break — setup membership.** The snapshots gain every skill upstream added
  and lose the three renamed names; the `go-cli-tui` profile's Setup Snapshot
  becomes `go-tui`.
- **Break — capability remediation.** The Context7 remediation names
  `context7-cli`. An installed `context7` still satisfies the capability.

## Prerequisites

- Spec 0195 is delivered first. It adds the Roundfix skill entry to the
  `go-cli` and `rust-cli` snapshots and teaches the sync to keep an owned entry
  the upstream list omits; the `go-tui` snapshot is seeded from `go-cli` so it
  inherits that entry.
- Spec 0199 is delivered first. It edits `core.json`, `typescript.json` and
  the activation file before this Spec does, and this Spec raises module
  versions from what the starting main holds.

The evidence owner for both is the default branch. task_02's Verification
fails when a setup lacks the Roundfix skill, because the `autonomous-work`
module Spec 0195 changes requires it.

## Decisions

- **Take upstream lists whole.** The sync is the only writer of a snapshot;
  a hand-filtered snapshot would be undone by the next sync. See ADR-0191.
- **Add `go-tui`, keep `go-cli`.** The maintainer chose the `go-tui` setup.
  `go-cli` is kept as the upstream's Go setup without a terminal interface; no
  built-in profile selects it.
- **Validate the result, not the start.** It makes a rename one sync run.
- **Dispatch implies requirement.** The `go` module was the only one that
  dispatched skills it did not require.
- **Accept the prior Context7 name.** It keeps existing adopters and the frozen
  parity evidence satisfied while new adopters get the current skill.
- **Dogfood the adopter path.** This repository is updated with the commands
  an adopter uses, not by copying trees.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- The upstream skills repository at `a4e18e4`, read through the Secondbrain
  mirror `projects/skills/mirror/setups/`: `go-tui.txt` lists `bubbletea` and
  `tui-design`, `go-cli.txt` does not, and the three renamed paths exist under
  their new names.
- The Secondbrain inbox entry
  `inbox/skills/2026-09-30-skills-vendorizadas-contradizem-o-baseline-ou-estao-quebradas.md`,
  written by another session, which records the three renames the Baseline had
  not followed.
- The `skills` command-line tool's update documentation
  (<https://vercel-labs-skills.mintlify.app/commands/update>): a skill removed
  from its source repository "disappears after update", and the stated remedy is
  to remove it from the lock file. Roundfix's reconcile step is the same remedy.

## Research basis

The Secondbrain was consulted through its index, the upstream mirror
`projects/skills/mirror/`, and
`qmd query "skills snapshot upstream renomeado setup go-tui context7-cli"`,
whose top results were Roundfix's archived Specs 0043 and 0030 on skill
restoration and the mirrored `context7-cli` skill. The inbox entry above
settled the rename list. Exa located the `skills` tool's lock-file source
(`src/skill-lock.ts`, where `ref` is the "Branch or tag ref used for
installation") and its update page; they explain why this repository's lock
records `ref: main` while the snapshots pin a commit.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the measured outputs, the procedure
each Task follows and the build order.
