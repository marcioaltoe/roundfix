---
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
status: active
created: 2026-10-25
surfaces: [backend, cli]
---

# The Baseline reports a required skill that is not installed

Minimal PRD for a bug fix, minted by `write-techspec`. When an adopted
repository's Baseline modules require and dispatch an external skill that is
not installed, the `roundfix baseline update` preview reports a clean
`current` with exit `0`, so an adopter that removed a skill the dispatch still
names is told nothing is wrong. After this Spec the preview lists each such
skill under `skills.missing` with the command that restores it and reports
`plan_ready` with exit `3`, as it already does for an owned skill older than
the binary's. `roundfix skills check`, run inside an adopted repository, fails
on the same skills, and `roundfix doctor`, which already fails on them, is
pinned by a test for both ways a skill disappears.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md`](references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md),
   `type: fix`, from Fluxus (ai-memory message `01a16ea1`, 2026-10-24 12:57
   BRT): after `bunx skills remove ... handoff ...`, `roundfix baseline update
   --repo . --no-history --format json` returned `ok:true`, `state:current`,
   `warnings:[]` and `skills.differing:[]` while
   `docs/agents/skill-dispatch.md` still required `handoff`; the request
   includes "verificar também que a validação detecte obrigatoriedade sem
   instalação".

No open Backlog Entry or unresolved Finding shares this context: at the
authoring base `docs/backlog/` and `docs/findings/` hold no entry. Spec 0317
answers the other half of the same message (the catalog stops requiring
`handoff`) and is not adopted here.

## Prerequisites

None in `requires`. This Spec is authored against `f3c39140` (v0.119.0) and
delivers right after Spec 0317, by the maintainer's answer quoted in
`_authorization.md`. Its tests never name the skill they remove: they read the
required set from the fixture's Setup Manifest, so they hold whether or not
Spec 0317 has removed `handoff` from the `core` module. Spec 0317 also
re-records `skills/testdata/owned-skill-versions.json`,
`docs/references/behavior-surfaces.json` and `CONTEXT.md`; each Task here
re-records those files on its own delivery base, so the order needs no
`requires`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the preview's JSON gains one optional list and one status value.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every inspection reads local files.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — no Baseline module, clause,
  trigger or Setup Snapshot changes. The decision records that govern the
  Setup Snapshot comparison (`skills.differing`, which never restores a
  lock-matching copy) and the Retired Skill report keep their behavior: a
  Missing Required Skill is absent or lock-mismatched, never a lock-matching
  copy, and a Retired Skill is required by no module, so it never enters the
  list. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to one test file
  and the Roundfix skill. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0318-the-baseline-reports-a-required-skill-that-is-not-installed/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- The `baseline update` preview never reports `current` while the adopted
  modules require an external skill that is absent from `.agents/skills` or
  `skills-lock.json`, or no longer matches its lock entry.
- Each such skill is named with the command that restores it, in JSON and
  text, and the state, category and exit code follow the existing precedent
  for an outdated owned skill.
- `roundfix skills check`, inside an adopted Git repository, fails on the
  same skills; outside one it keeps its bundle-only result.
- `roundfix doctor` keeps failing on them, now pinned for both removal shapes.
- The guides, the Roundfix skill references and `CONTEXT.md` describe it.

## Core Features

1. The preview list. Without `--yes` or `--confirm-plan`, and unless
   `--no-skills`, the preview adds `skills.missing` (each entry `skill` and
   `reason`, the reason ending with `restore with roundfix baseline skills
   restore --profile <id> --skill <name>`), sets `skills.status` to `missing`,
   and with unchanged guidance reports `plan_ready`, category `approval`,
   exit `3`, a message naming the skills, and a next action naming `--yes`
   and the restore command. Text prints `Skills missing: <n>` and one
   `- missing <name>: <reason>` line each.
2. The repository section of `skills check`. Inside a Git repository with a
   Setup Manifest, `roundfix skills check` adds one diagnostic per Missing
   Required Skill to its existing failure list on stderr and exits `1`.
3. Doctor characterization. A test pins `roundfix doctor`'s existing
   `skills: failed` line and exit `1` when the skill's folder and lock entry
   are both gone and when only the folder is gone.
4. Documentation. `docs/user-guide/commands/baseline.md`,
   `docs/user-guide/commands/skills.md`, the Roundfix skill's
   `references/baseline.md` and `references/setup.md`, and `CONTEXT.md`.

## User Experience

An adopter who runs `bunx skills remove handoff -y` while the Baseline still
dispatches `handoff` and then runs `roundfix baseline update --repo . --format
json` reads `"state":"plan_ready"`, `"skills":{"status":"missing","missing":
[{"skill":"handoff","reason":"required by the Baseline and not installed;
restore with roundfix baseline skills restore --profile <id> --skill
handoff"}]}` and exit `3`. They either restore the skill or, when the upstream
retired it, report the Baseline defect instead of keeping a local exception.
`roundfix skills check` in the same checkout exits `1` with the same reason.

## Non-Goals / Out of Scope

- Removing `handoff` or any other skill from the catalog: Spec 0317.
- Reporting owned Roundfix skills that are absent from the repository in the
  preview. Spec 0195 pins that an absent owned skill leaves the preview
  `current` (`TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch`), and
  `--yes` installs the owned bundle unconditionally; Doctor already fails on
  it. This Spec keeps that contract.
- Reporting a repository whose `skills-lock.json` is absent or unreadable in
  the preview: the inspection names no skill then, Doctor already fails on it,
  and treating it as every skill missing would move seven existing preview
  tests built on lock-less fixtures (measured at authoring time).
- Changing `--yes`: it already restores a missing or lock-mismatched required
  skill, and reports `skills.drifted` when the source is unreachable.
- Changing Doctor's bytes, its next action, or the `skills.differing`
  comparison.
- Changing this repository's Makefile; its `skills-check` target keeps its
  description.

## Success Metrics

1. Success Metric: in a temporary adopted repository with every required
   external skill installed and locked, removing one required skill's folder
   and lock entry, removing only its folder, or editing it so it no longer
   matches its lock makes the JSON preview exit `3` with `state:plan_ready`,
   `skills.status:missing` and exactly that skill under `skills.missing`; the
   text preview prints its `Skills missing` lines; `--no-skills` stays
   `current`; and the repository bytes are unchanged.
2. Success Metric: with every required skill installed, the preview stays
   `current`, exit `0`, with no `missing` key.
3. Success Metric: `roundfix skills check` exits `1` with the skill's
   diagnostic in that repository and exits `0` with its unchanged pass line in
   a complete adopted repository, in a repository without a Setup Manifest
   and outside a Git repository.
4. Success Metric: `roundfix doctor` exits `1` with a `skills: failed` line
   naming the skill for both removal shapes.
5. Success Metric: the existing preview, retired-skill and owned-skill tests
   pass, with the retired-skill fixture given its required skills.
6. Success Metric: the documentation contracts pass with the Behavior Surface
   record re-recorded last, and the owned-skill version is recorded.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-25.

- The vercel-labs `skills` CLI issue #577, "skills remove does not clean up
  skills-lock.json entries" (https://github.com/vercel-labs/skills/issues/577):
  "`npx skills remove <skill>` removes the skill files from `.claude/skills/`
  and `.agents/skills/` but does not remove the corresponding entry from
  `skills-lock.json`." This is why both removal shapes are covered.
- The Fluxus report in the adopted Backlog Entry, written before this Spec
  existed: `roundfix baseline update --repo . --no-history --format json`
  returned `ok:true, state:current, fileChanges:[], warnings:[],
  skills.differing:[]` at catalog
  `sha256:a5c687aa60ad35f3d967cb4bdc3094f404c8ff103f39447c4c5e83c8d4fa22ba`.

## Unreachable Acceptance

- criterion: after release, Fluxus's `roundfix baseline update` names a
  required skill it removed instead of reporting `current`.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the operator's reply to `default/gesttione-solutions-fluxus`
  through ai-memory after the release, and Fluxus's next managed refresh.

## Validation

The decisions were checked on 2026-10-25 against these sources.

- Reproduction at `f3c39140` with `bin/roundfix` built in a scratch clone,
  `HOME` and `TMPDIR` at fresh temporary directories: a temporary repository
  adopted `go-cli-tui` through `baseline plan` and `baseline apply`, installed
  the owned bundle, and gave every required external skill a folder and a
  matching `skills-lock.json` entry. Then `handoff`'s folder and lock entry
  were removed. `roundfix baseline update --repo . --no-history --format json`
  exited `0` with `ok:true`, `state:current`, `fileChanges:[]`,
  `warnings:[]`, `skills.differing:[]` at catalog
  `sha256:a5c687aa60ad35f3d967cb4bdc3094f404c8ff103f39447c4c5e83c8d4fa22ba`,
  the same as Fluxus. `roundfix doctor` exited `1` with `skills: failed (...
  external skills are missing: handoff; next: bunx skills add
  marcioaltoe/skills@handoff)`. `roundfix skills check` exited `0` with its
  bundle pass line. Removing only the folder gave the same preview, Doctor
  `skills: failed (missing: handoff; ...)`, and the same `skills check`.
- Root cause, read at `f3c39140`: the preview branch of
  `internal/cli/baseline_update.go` (from line 285) inspects the external
  Repository Skill Set only to compute `skills.differing`, and
  `snapshotComparableSkills` (`internal/cli/doctor.go:435`) drops the missing
  and lock-mismatched names; nothing reads `readiness.MissingExternal`,
  `readiness.OutdatedExternal` or the lock-entry error's `MissingExternal`
  there, so the state falls through to `current` (line 337). The applied
  path restores exactly those names (`baselineUpdateExternalDrift`, line 600),
  so the preview and `--yes` disagree. `skills check` (`internal/cli/cli.go`
  `case "check"`) validates only the binary's embedded bundle.
- A prototype in a scratch worktree implemented the TechSpec. The three
  removal shapes reported `plan_ready`, exit `3`; `skills check` exited `1`;
  the full `./internal/cli` suite moved only
  `TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState`, whose
  fixture lock lacks the required skills; with the fixture given them it
  passed. Counting an absent lock as every skill missing moved seven tests
  instead, which is why that state stays out of scope.
- Secondbrain, read only: `qmd query "roundfix baseline skill dispatch
  required skill not installed" --all --files --min-score 0.3` returned the
  Roundfix mirror's decision records on upstream setup names, owned-skill
  versions, Setup Snapshots and vendored skills, and
  `docs/agents/skill-dispatch.md`; none records a decision to let the preview ignore a
  missing required skill. No Inbox Entry was written.
- Exa: the `skills` CLI documentation for `skills remove` (removes the
  canonical `.agents/skills/<name>` directory and agent symlinks) and issue
  #577 quoted under Acceptance evidence.

## Glossary

- adds: **Missing Required Skill**

## Decisions

- 2026-10-25: Detection, asked through AskUserQuestion at about 09:05 BRT —
  the maintainer chose "Spec 0318 logo após a 0317 (Recomendado)": "`baseline update`, `doctor` e `skills check` acusam
  quando o dispatch exige uma skill que não está instalada. 1–2 Tasks + QA.
  Não atrasa a 0317."

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 7 briefing for Spec 0318 of 2026-10-25 and the characterization above.

Taken from those sources, with no further interview:

- Severity follows the existing precedent for an owned skill older than the
  binary's (Spec 0195, `TestBaselineUpdatePreviewReportsAnOlderOwnedSkill`):
  `plan_ready`, category `approval`, exit `3`, because `--yes` would act on
  it. A `warnings` entry with exit `0` was rejected: it would keep
  `ok:true` and the preview would still disagree with what `--yes` does.
- `skills check` gains no flag; it inspects the repository that contains the
  work directory, the way `skills install --target project` finds it, and
  reports through its existing failure list and exit `1`.
- Doctor already reports the class; it gains a characterization test only.
- No new ADR: the decision follows an existing contract.

## Open Questions

None.
