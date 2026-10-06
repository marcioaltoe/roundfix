---
spec: 0236-baseline-update-with-a-repository-profile
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# Baseline update with a repository profile

On Roundfix 0.43.0 `roundfix baseline update` exits 1 in every repository
adopted with its own Baseline Profile, after it has already computed a valid
plan:

```text
roundfix: baseline update failed: load profile for snapshot comparison: Unknown built-in Baseline Profile "pantheon-devops".
```

Pantheon and Oraculum reported it on 2026-10-05
([references/2026-10-06-baseline-update-refuses-a-repository-profile.md](references/2026-10-06-baseline-update-refuses-a-repository-profile.md)).
The snapshot comparison Spec 0215 added (`TrailingSetupSkills` in
`internal/baseline/skills_trailing.go`) loads the profile through the skill
restore loader, which resolves only built-in profiles, while the rest of
`update` and `baseline profile validate` also resolve
`.roundfix/baseline/profiles/<id>.json`. Doctor shows the same error as
"snapshot comparison unavailable".

This is a bug fix. During authoring it was reproduced with the built binary
(`bde616d1`) on a disposable repository adopted through
`baseline profile init --id scratch-go --from go-cli-tui`, `baseline plan
--profile scratch-go` and `baseline apply`:

1. `baseline update --format json` exits 1 with the message above for
   `scratch-go`; `--no-skills` reports `current`.
2. `baseline update --yes` installs the fourteen Roundfix-owned skills and
   then exits 1, `skills.status: failed`, on the same comparison.
3. `roundfix doctor` appends "snapshot comparison unavailable: load profile
   for snapshot comparison: Unknown built-in Baseline Profile "scratch-go"."
   to its `skills:` line.
4. `baseline plan --profile scratch-go` recommends `roundfix baseline skills
   restore --profile scratch-go --skill context7-cli`, and that command
   refuses with `restore.profile-unknown`, so the restore path the plan
   names is also closed. Lock reconciliation uses the same loader.

## Prerequisites

None. Specs 0235 and 0237 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the record command picks a free version
and the operator orders the queue (ADR-0233).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  only new names are two restore finding codes,
  `restore.profile-unresolved` and `restore.snapshot-conflict`, in the
  existing `restore.` family. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the comparison and the profile
  resolution read only local files and the embedded catalog; restore keeps
  its existing source acquisition, and no test reaches the network. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0241 (this Spec): "The skill
  readers now resolve a profile exactly as the rest of `baseline update` and
  `baseline profile validate` do". Repository profiles stay
  repository-owned, ADR-0067: "Custom profiles live at
  `.roundfix/baseline/profiles/<id>.json`". The comparison keeps its
  meaning, ADR-0221: "Doctor now compares each required external skill that
  is installed and matches its lock". This repository's own skills stay
  held to the snapshot through that comparison, ADR-0224: "every required
  external skill of this repository's Baseline Profile must match the
  `treeDigest` its Setup Snapshot pins", and a composed snapshot already
  refuses disagreement, ADR-0204: "would merge two different entries under
  one skill name". ADR-0219 is left unchanged: draft binding keeps choosing
  a plan's source. ADR-0189 and ADR-0233 bind the Roundfix Skill's version
  raise. The QA gate follows ADR-0080: "QA verdicts distinguish
  environment-blocked rows", and ADR-0091: "required to be terminal and to
  depend on every leaf"; ADR-0096, ADR-0097, ADR-0104 and ADR-0167 bind its
  machine stage, row carry, outside evidence and pre-PR Pull Request row,
  and ADR-0117, ADR-0182, ADR-0194, ADR-0195 and ADR-0210 its stage
  ownership, settlement, row record, re-observation and evidence snapshot.
  ADR-0093, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt, and ADR-0184 binds the TechSpec's
  three Surface Transcripts. ADR-0229 and ADR-0237 cite ADRs listed here but
  decide the Delivery Retry of an operator archive or an outside merge,
  which this Spec does not change, so neither applies. ADR-0203 and
  ADR-0074 cite ADRs listed here but decide repository-rule ownership in the
  Baseline, which this Spec does not change, so neither applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's baseline reference
  and both copies of its `SKILL.md` are Governed Paths; the maintainer
  authorized skill edits ("considere autorizado a ajustar todas as skills
  se necessário") and this Spec's Governed Paths ("Concedo", 2026-10-06).
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0236-baseline-update-with-a-repository-profile/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A repository adopted with its own profile reaches `state: current` through
  `roundfix baseline update` once its guidance and skills match, and its
  preview lists an external skill that trails its Setup Snapshot, as a
  built-in-profile repository does.
- Roundfix-owned skills refresh through `baseline update` in such a
  repository: the preview reports outdated owned skills and `--yes`
  installs them without the run failing afterwards.
- Doctor, skill restore and lock reconciliation resolve the same profile the
  same way.
- A profile that resolves neither as built-in nor from the repository is
  reported with the path that was searched, never as an unknown built-in
  profile.

## Core Features

1. One skill-contract resolution for built-in and repository profiles,
   used by the snapshot comparison, skill restore and lock reconciliation.
2. A repository profile requires the skills of its own modules, with each
   external contract taken from the embedded Setup Snapshots that agree on
   it; disagreement is refused by name.
3. Findings that name `.roundfix/baseline/profiles/<id>.json` for an
   unresolvable profile, and an update next action that repeats the
   finding's action.
4. Regression tests through the real `baseline update`, Doctor and restore
   commands on a repository adopted with a repository profile.
5. The baseline and doctor references and the Roundfix Skill's baseline
   reference describe repository profiles in these commands.

## Non-Goals / Out of Scope

- The repository profile schema, `baseline profile init`, draft binding
  (ADR-0219) and profile alignment.
- The embedded catalog, its Setup Snapshots and every derived Baseline file.
- Skipping the comparison with a warning: a repository profile is compared,
  not exempted.
- Adopters' repositories; they upgrade by running the fixed binary.
- The Makefile, `go.mod`, CI workflows, `CONTEXT.md` and `CHANGELOG.md`.

## Success Metrics

1. Success Metric: on a repository adopted with a repository profile whose
   installed, locked external skills match the snapshot, `roundfix baseline
   update --format json` exits 0 with `state: current` (before: exit 1 with
   "Unknown built-in Baseline Profile").
2. Success Metric: on the same repository with one locked external skill
   edited, the preview exits 3 with `plan_ready` and lists that skill under
   `skills.drifted` with "trails its Setup Snapshot".
3. Success Metric: the skills stage of `baseline update --yes` for a
   repository profile installs the owned skills, asks restore for the
   trailing skill with the repository profile's ID, and reports
   `verified`.
4. Success Metric: Doctor's `skills:` line for that repository carries
   `DR-SKILL-TRAILS-SNAPSHOT` for the edited skill and no "snapshot
   comparison unavailable".
5. Success Metric: `baseline skills restore --profile <repository-id>`
   proceeds past profile resolution, and an unresolvable ID prints
   `restore.profile-unresolved` naming
   `.roundfix/baseline/profiles/<id>.json` and not "built-in".
6. Success Metric: every existing baseline, restore, reconcile, Doctor and
   update test passes with only the declared changes.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Pantheon's report of 2026-10-05, Secondbrain
  `inbox/roundfix/_triaged/2026-10-05-baseline-update-falha-com-perfil-de-repositorio.md`:
  `baseline update` exits 1 on 0.43.0 (`69462a0c`) with
  `.roundfix/baseline/profiles/pantheon-devops.json`, after a valid plan with
  four outdated owned skills, while `baseline profile validate` answers
  `valid (repository)`.
- Oraculum's report of 2026-10-05, Secondbrain
  `inbox/roundfix/_triaged/2026-10-05-baseline-update-recusa-perfil-proprio-na-checagem-de-snapshot.md`:
  the same failure with `oraculum-backend`, the `--no-skills` plus
  `roundfix skills install --target project` workaround, and Doctor's
  "snapshot comparison unavailable".

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "baseline update repository profile snapshot comparison"`
(`--all --files --min-score 0.3`); it returned this repository's mirrors of
ADR-0224, ADR-0203 and earlier Baseline Specs, none of which resolves a
repository profile's skills, and the two adopter reports above. Exa found
Gentoo's custom profiles guide
(<https://wiki.gentoo.org/wiki/Portage/Profiles/Custom_profiles>), where a
custom profile refers to the distribution's profiles "to avoid recreating all
profile definitions", and Yarn's workspace profiles
(<https://v6.yarnpkg.com/concepts/profiles.html>), where a profile that
extends others inherits their pinned dev dependencies: in both, a local
profile takes its pinned set from the shipped definitions rather than
restating it, which is the shape ADR-0241 chooses. No other open Backlog
Entry or Finding shares this context: the two other entries of 2026-10-06
belong to Specs 0235 and 0237.

## Decisions

- Resolve the profile in the skill readers as `update` does; take a
  repository profile's contracts from the agreeing embedded snapshots. See
  ADR-0241.
- Keep the comparison for repository profiles rather than skipping it.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
