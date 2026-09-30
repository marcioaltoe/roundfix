---
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# Owned skills and a release step that follow the bundle

Roundfix ships 14 owned skills inside its binary and compares each installed
copy with a minimum version. The minimum is one literal for all of them, and
the bundle is already ahead of it. A repository can therefore run an older
`qa-gate` than the binary carries while Doctor reports `skills: ok` and the
`baseline update` preview reports `current`. A skill's content can also change
while its version stays the same, so the comparison says nothing about
content. Two setups omit the Roundfix skill, and no profile tells an Agent
when to load it. No release re-reads skills or guides before it ships.

This Spec makes the minimum follow the bundle, makes a version name one
content, lists and dispatches the Roundfix skill everywhere, and adds the
release step that checks all of it. The adopted entries in
[references/_index.md](references/_index.md) record each case with its
evidence.

## Project Constraints

- Identifier strategy: applicable — one new dispatch trigger,
  `trigger.autonomous-work.roundfix`, follows the existing
  `trigger.<module>.<skill>` form. Two optional result fields,
  `skills.outdated` and the status value `outdated`, join the
  `roundfix/baseline-update-result/v1` result. No other identifier is added
  or renamed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — embedded files, local files and
  Git only; no credential is read and no network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0189 (this Spec) makes the minimum
  version of an owned skill the version the binary carries, makes a version
  name one content, keeps the setup minimum as the catalog's lower bound, and
  requires the release step. ADR-0186 governs the wording of the sentence
  added to the release clause. ADR-0062 restarted owned versions at `0.0.1`
  and stays in force. ADR-0143 keeps release planning read-only, which is why
  the Release Plan Command does not gain a line here. ADR-0081 and ADR-0149
  keep sanctioned regeneration outputs inside a grant, and ADR-0073 and
  ADR-0103 require this repository's managed refresh to apply in one
  transaction and converge to `current`. This Spec's gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and
  ADR-0156, and ADR-0093 and ADR-0094 check its consistency by citation and
  artifact presence. These ADRs cite a listed ADR and do not apply, because
  this Spec touches none of their behaviors: ADR-0168 (the related-ADR gap)
  and ADR-0176 (citation checks read authored text) cite
  ADR-0093; ADR-0173 (citations a History Relocation breaks) cites ADR-0073.
  Reached only through those citations, and equally untouched: ADR-0097, ADR-0167, ADR-0177.
  All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization of
  2026-09-30 for the Baseline source, its generated guides and the owned
  skills ("Autorizar os dois"; "considere autorizado a ajustar todas as
  skills se necessário"), recorded in [_authorization.md](_authorization.md);
  bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `docs/agents/specific-repository.md`,
  `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `docs/agents/skill-dispatch.md`, `docs/agents/agent-instructions.md`,
  `docs/agents/setup-context.json`.
  Sanctioned regeneration: `make baseline-digests`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A repository whose installed owned skill is older than the one the binary
  carries is told so, by Doctor and by the `baseline update` preview, with the
  command that fixes it.
- Two copies of an owned skill that declare the same version have the same
  content.
- Every setup the Baseline ships lists the Roundfix skill, and every profile's
  guide tells an Agent when to load it.
- Each release confirms, before its Pull Request, that skills and guides
  describe what is being released, through checks that need no model and no
  network.

## Core Features

1. **An installed owned skill older than the bundle is reported.** The minimum
   version of each owned skill is the version the binary carries for it, read
   from the embedded skill; the separate list is removed.
   - Doctor's `skills:` line fails for an installed copy below that version
     and names the skill, the version required, the version found and
     `roundfix skills install --target project`.
   - `roundfix baseline update` without confirmation lists each installed
     owned skill that is older. When the guidance is otherwise unchanged it
     reports `plan_ready`, not `current`, and says that `--yes` refreshes the
     Repository Skill Set. `--no-skills` skips the read.
   - A repository whose owned skills match gets the output it gets today,
     byte for byte.
2. **A version names one content.** The repository keeps a record of every
   version each owned skill has shipped, with the digest of that version's
   folder. A test refuses a skill whose content differs from the digest
   recorded for its version, and a version that is not recorded. Recording
   adds a new, higher version and never replaces a recorded digest. The
   repository's own guide states the rule where a Spec author reads it.
3. **The Roundfix skill is in every setup and has a dispatch trigger.** The
   `go-cli` and `rust-cli` setups list the Roundfix skill, as `typescript-bun`
   already does. The `autonomous-work` module requires it and gives it one
   trigger, so every profile's skill-dispatch guide names it. The asset sync
   keeps an owned entry the snapshot records when the upstream list omits it.
4. **Every release checks skills and guides first.** The release runbook gains
   a mandatory step, run before the release Pull Request. It names the checks
   to run: Doctor's `skills:` line, the version record, the minimum check, the
   Baseline clause checks Spec 0193 adds, the command documentation checks
   Spec 0192 adds, and this repository's managed refresh. It ends with a
   reading pass over what the release changes. A test refuses a step that
   names a check the repository does not have. The release clause of the
   Baseline tells every adopter's Agent that the step exists.

## Non-Goals / Out of Scope

- The membership of upstream skills in a setup: the Go skills a module
  dispatches and no setup lists, a search skill that setups list and no module
  requires, the count of external skills the user guide prints, and skills
  renamed upstream. A separate Spec refreshes the upstream skill snapshot and
  owns them.
- The dispatch trigger for `cut-release`, a skill only a person can start. It
  is recorded as the open Backlog Entry
  `docs/backlog/2026-09-30-a-dispatch-trigger-names-a-skill-the-model-cannot-invoke.md`.
- A `skills:` or `baseline:` line in `roundfix release plan`. It is recorded
  as the open Backlog Entry
  `docs/backlog/2026-09-30-release-plan-does-not-report-skill-and-guide-checks.md`.
- Any change to the Makefile, including the `skills-link` target. See
  Recorded limits.
- Raising `minimumVersion` in a setup snapshot.
- The statements inside the owned skills, which Spec 0192 corrects, and the
  Baseline clause corrections, which Spec 0193 makes. This Spec edits two
  sentences of the Roundfix skill that describe the behavior it changes, and
  appends one sentence to one clause.
- A model-based scan of skills or clauses as a release step.

## Success Metrics

1. Success Metric: for each of the 14 owned skills the minimum version equals
   the version its embedded `SKILL.md` declares, and a minimum that differs is
   reported by the check.
2. Success Metric: for a repository whose installed `qa-gate` is one patch
   below the bundle, Doctor prints `skills: failed` with the skill, both
   versions and the install command, and the `baseline update` preview exits
   `3` with `plan_ready` and one `skills.outdated` entry. For a repository
   whose owned skills match, both outputs are byte-identical to today's.
3. Success Metric: an owned skill whose folder digest differs from the digest
   recorded for its version fails the record check, an unrecorded version
   fails it, and recording refuses to replace a recorded digest.
4. Success Metric: all three setups list the Roundfix skill, the
   skill-dispatch guide of every profile carries
   `trigger.autonomous-work.roundfix`, and an asset sync from an upstream list
   without the Roundfix skill keeps it while still dropping an external skill
   the list omits.
5. Success Metric: the release runbook lists the skills-and-guides step before
   the tag step, every check the step names exists in the repository, and a
   step that names a missing check is reported.
6. Success Metric: after the sanctioned regeneration and this repository's
   managed refresh, a second refresh reports no file change; every setup keeps
   its recorded `minimumVersion` values and every clause keeps its enforcement
   level.

## Declared breaks

- **Break — Doctor.** A repository with an installed owned skill older than
  the binary's now gets `skills: failed`. It got `skills: ok` whenever the
  copy was at least `0.0.2`. The line prints the command that fixes it.
- **Break — update preview.** `roundfix baseline update` without `--yes`,
  `--confirm-plan` or `--no-skills` exits `3` with `plan_ready` for such a
  repository when its guidance is unchanged. It exited `0` with `current`.
- **Break — owned skill edits.** A change to an owned skill's content without
  a version change now fails `go test ./skills`, and with it the repository
  gate. A later Spec that edits an owned skill raises its version and records
  it in the same Task.
- **Break — asset sync.** A refresh from an upstream list that omits an owned
  skill keeps the snapshot's entry. It dropped it before.
- **Addition — dispatch.** Every profile's skill-dispatch guide gains one
  trigger for the Roundfix skill, and adopters on the `go-cli` and `rust-cli`
  setups gain the skill in their setup.
- **Clarification — release clause.** One sentence is appended. Its
  enforcement level and every existing sentence stay.

## Recorded limits

- The preview reports an owned skill only when it is installed and older. A
  repository with no installed owned skill, or with a skill directory Roundfix
  cannot read, keeps today's preview result, and Doctor reports it. Counting a
  missing skill changed the result of three existing update tests, measured on
  2026-09-30.
- A setup snapshot's `minimumVersion` stays `0.0.2` while the binary requires
  more. The snapshot states the catalog's lower bound, and the binary states
  what it carries.
- The record starts with the versions present when this Spec is implemented.
  Versions shipped before that are not recorded.
- The `skills-link` Makefile target removes every entry under `.claude/skills`
  and links every directory under `.agents/skills`, including skills no setup
  lists. The Makefile is outside this Spec's authority, so the target is left
  as it is.
- The last part of the release step is a reading pass. The checks prove what
  code can see: versions, digests, clause order, duplicates, citations and
  command names.

## Prerequisites

- Spec 0192 is delivered first. It owns
  `TestEveryCommandIsNamedInTheRoundfixSkill` and
  `TestEveryCommandIsNamedInTheUserGuide`, which the release step names.
- Spec 0193 is delivered first. It owns the clause checks the release step
  names, and it edits `core.json` and `autonomous-work.json` before this Spec
  does.

The evidence owner for both is the default branch: task_04's check fails when
a named test is absent.

## Decisions

- **One source for the minimum.** The embedded skill states its version, and
  the minimum is that version. See ADR-0189.
- **The setup minimum does not follow.** Raising it with each skill release
  would change the catalog on every skill edit.
- **Report what is installed and older, not what is missing.** It keeps every
  existing preview result.
- **The record has its own flag.** The shared `-update` flag of the package's
  tests cannot record a version by accident.
- **Checks, not a model scan.** The release step names plain tests and two
  commands.

## Acceptance evidence

Each Core Feature requires positive and negative evidence in the Task Graph.
The outside-evidence row rests on sources this Spec did not produce:

- The upstream skill lists, read through the Secondbrain mirror on
  2026-09-30. `projects/skills/mirror/setups/go-cli.txt` and
  `setups/rust-cli.txt` omit `skills/06-review-repair/roundfix`, and
  `setups/typescript-bun.txt` lists it on line 107.
- The Semantic Versioning specification, item 3: "Once a versioned package
  has been released, the contents of that version MUST NOT be modified. Any
  modifications MUST be released as a new version."
  (<https://semver.org/spec/v2.0.0.html>).
- The adopter manifests in the Secondbrain mirrors: `conexus`, `fiscus`,
  `fluxus`, `tax-poc` and `vortex` record the `standard-typescript-monorepo`
  profile, so they receive the new trigger on their next refresh.

## Research basis

The Secondbrain was consulted through the upstream skill mirror
`projects/skills/mirror/setups/`, the adopter manifests
`projects/*/mirror/docs/agents/setup-context.json`, and the triaged entries
`inbox/roundfix/_triaged/2026-08-08-baseline-update-medido-contra-a-frota-adotada.md`
and
`inbox/roundfix/_triaged/2026-08-30-skills-nao-declara-ownership-e-forca-todo-grant-a-enumerar-espelhos.md`.
The query
`qmd query "skills owned versão mínima doctor desatualizada release checklist"`
returned nothing in its time limit. The mirror settled where the Roundfix
skill is listed upstream. The first triaged entry measured `baseline update`
against the adopted fleet and is why the preview keeps every existing result.
Exa located the Semantic Versioning specification, whose item 3 is the rule
the version record enforces.

## Technical candidate

The [_techspec.md](_techspec.md) records the interfaces, the fixed texts, the
measured regeneration outputs and the build order.
