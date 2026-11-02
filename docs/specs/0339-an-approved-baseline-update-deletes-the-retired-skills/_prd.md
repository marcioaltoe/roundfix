---
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
status: active
created: 2026-11-01
surfaces: [backend, docs]
---

# An approved Baseline update deletes the Retired Skills

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix baseline
update` lists each Retired Skill an adopter still holds under `Skills retired`
and deletes nothing, while the v0.127.0 release note said the update removes
them; Argus removed seven copies by hand on 2026-11-01. After this Spec an
approved update deletes each Retired Skill whose copies it proves unedited
against `skills-lock.json`: the `.agents/skills/<skill>` folder, the
`.claude/skills/<skill>` link or copy and the lock entry. The preview lists
every deletion and the Plan Digest binds them, an edited or unrecorded copy is
kept with its reason, and nothing outside the Retired Skills set is ever
deleted, so `context7-cli`, `exa-web-search` and a repository's own skills
survive. ADR-0350 supersedes ADR-0246.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md`](references/2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md),
   `type: fix`, reported by Argus (ai-memory message `01a19718`,
   2026-11-01) with the maintainer's decision of the same day.

At the authoring base `docs/backlog/` holds no other entry and
`docs/findings/` no unresolved Finding on Retired Skills or on deleting
installed skills.

## Prerequisites

None in `requires`. Authored against `dd3b2b62` (v0.128.0), where no other
Spec is active. The cycle queue delivers 0338 first and 0340 after this
Spec; neither changes a file this Spec's code or tests read. They share only
the Roundfix skill's `SKILL.md` pair, the owned-skill version record,
`CONTEXT.md` and the Behavior Surface Record, which the regeneration commands
rewrite on rebase (`_techspec.md` → Build Order).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the `roundfix/baseline-update-result/v1` JSON gains optional
  fields inside its existing schema. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the update reads and writes a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0350, "An approved Baseline update
  deletes the Retired Skills it proves unedited", which this Spec's planning
  change records and every Task follows; ADR-0246, "A retired skill leaves
  the Baseline and the adopter removes its copy", which ADR-0350 supersedes
  and task_02 retires by deletion; ADR-0191, "A setup snapshot follows its
  upstream by name", whose rule that Roundfix never deletes an installed
  skill tree ADR-0350 narrows for Retired Skills only and task_02 amends; and
  ADR-0206, "The Baseline takes upstream setup names and drops a removed
  skill", which ADR-0246 cites and this Spec leaves unchanged; ADR-0204 and
  ADR-0322, which cite ADR-0191 for composed and Python setups, ADR-0219,
  which cites ADR-0204 for profile drafts, and ADR-0241, which cites ADR-0219
  for repository profiles' skill contracts: none applies, because this Spec
  changes no setup composition, import, draft or snapshot. Source:
  `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AU answers of 2026-11-01 quoted in the authorization record, which keep the
  cycle AT and AS rules and the cycle AJ grant for the Governed Paths each
  Spec declares; the governed files are the Roundfix skill's canonical
  `baseline` reference and its `SKILL.md` pair. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0339-an-approved-baseline-update-deletes-the-retired-skills/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- An approved `roundfix baseline update` (`--confirm-plan` or `--yes`)
  deletes each Retired Skill whose copies it proves unedited: the
  `.agents/skills/<skill>` folder, the `.claude/skills/<skill>` entry and the
  `skills-lock.json` entry, and reports each deletion in text and JSON.
- The preview lists every deletion before approval, the Plan Digest binds the
  deletion set, and a pending deletion alone makes the preview `plan_ready`.
- A retired copy that differs from its recorded hash, has no recorded hash,
  or is not a plain folder or `.agents/skills` link is kept, listed with its
  reason and the paths to delete by hand.
- Nothing outside the Retired Skills set is deleted: `context7-cli`,
  `exa-web-search`, which every Setup Snapshot requires, and any skill the
  repository installed on its own survive an approved update byte for byte.
- ADR-0350 replaces ADR-0246, and the glossary, the baseline guide and the
  Roundfix skill's baseline reference say what the update deletes.

## Core Features

1. Retired Skill dispositions. `baseline.InstalledRetiredSkills` gives each
   held Retired Skill an action, `delete` or `keep`, a keep reason
   (`edited`, `unrecorded`, `unsafe`) and the lock entry's hash.
2. Deletion. `baseline.DeleteRetiredSkills` deletes exactly the planned
   `delete` entries after re-inspecting them, and refuses, deleting nothing,
   when one changed since the plan or is not a Retired Skill.
3. The update. The preview inspects the Retired Skills before the digest,
   folds the deletions into the Plan Digest, reports `plan_ready` when a
   deletion is the only pending change, and the approved apply runs the
   deletion after the skills stage and reports `Skills deleted`.
4. The rules in writing. ADR-0350, the ADR-0191 amendment, the retirement of
   ADR-0246, **Retired Skill** in `CONTEXT.md`, the baseline guide and both
   copies of the Roundfix skill's baseline reference.

## Non-Goals / Out of Scope

- Deleting any skill outside the Retired Skills set, including a renamed
  skill's old folder, which ADR-0191 keeps the repository's, ADR-0191: "the
  old directory is the repository's to delete".
- A flag to keep an unedited Retired Skill; a repository that wants one keeps
  it under another name or edits it.
- Changing the Retired Skills set, the asset sync, the Setup Snapshots, the
  Doctor checks or `baseline skills reconcile`.
- Correcting the v0.127.0 CHANGELOG line; the operator corrects it in the
  release that ships this Spec, and no Task edits `CHANGELOG.md`.
- Any Baseline module or guide under `docs/agents/`, any text inside a
  `### QA settlement` section, the Makefile, the CI workflows, `go.mod` or
  `go.sum`.

## Success Metrics

1. Success Metric: in a temporary adopted repository holding an unedited
   `tdd` copy with its `.claude/skills/tdd` link and lock entry, the preview
   exits 3 `plan_ready` with the message "guidance matches the current
   Baseline catalog; the approved update deletes 1 Retired Skill(s)", and
   `--confirm-plan <digest>` deletes the folder, the link and the lock entry
   and reports `Skills deleted: 1`; at `dd3b2b62` the preview reports
   `current` and the apply deletes nothing.
2. Success Metric: in the same repository `context7-cli`, `exa-web-search`,
   a repository-installed `local-only` skill with its link, and an edited
   `design-patterns` copy keep their folders, links and lock entries after the
   approved apply, and the edited copy is listed as kept with reason
   `edited`.
3. Success Metric: the preview's Plan Digest differs from the `--no-skills`
   digest when a deletion is pending, a confirmation with that digest after
   the planned copy was edited is refused with no byte changed, and a preview
   without a deletion keeps its digest bytes.
4. Success Metric: `baseline.InstalledRetiredSkills` gives `delete` to an
   unedited folder, an unedited `.claude/skills` copy and a lock-only entry,
   and `keep` with `edited`, `unrecorded` and `unsafe` to the three keep
   cases; `baseline.DeleteRetiredSkills` refuses a name outside the set.
5. Success Metric: no Setup Snapshot under
   `internal/baseline/assets/setups/` lacks `context7-cli` or
   `exa-web-search`, and neither name is a Retired Skill.
6. Success Metric: ADR-0246 is deleted, ADR-0191 carries its amendment, and
   `CONTEXT.md`, the baseline guide and both copies of the Roundfix baseline
   reference describe the deletion and no longer say the update never deletes
   a retired copy.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Argus report quoted in the adopted reference, written by the Argus
  project before this Spec existed: the approved apply of v0.127.0 printed
  `Skills retired: 7` and deleted nothing.
- The FreeBSD `pkg` FAQ ("What happens if I delete a package where I've
  modified one of the files managed by the package?",
  https://github.com/freebsd/pkg/blob/main/FAQ.md, read through Exa on
  2026-11-01): "The package is removed, but modified file is not", with the
  tool's message "fails original SHA256 checksum, not removing". The keep
  reason `edited` follows that published rule: a removal deletes only what
  still matches the recorded checksum.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-01 against these sources.

- Reading at `dd3b2b62`: `InstalledRetiredSkills`
  (`internal/baseline/retired_skills_installed.go`) only lists paths and the
  lock entry; `runBaselineUpdateCommandWithSkillsStage`
  (`internal/cli/baseline_update.go`) calls it in the preview and after the
  skills stage, outside the Plan Digest, and `TestBaselineUpdateAppliesAndListsRetiredSkills`
  asserts that an applied update leaves the copies and the lock unchanged.
- Read-only survey of two adopters' checkouts (Fluxus at `043368ac`, Conexus
  at `6170544`): each holds the eight Retired Skills of 0.127.0 with a
  `.claude/skills` link to `../../.agents/skills/<skill>` and a lock entry
  whose `computedHash` equals the folder's hash, so every one would be
  deleted; their content was not copied into this repository or sent anywhere.
- A prototype of task_01 in a second clone passed the tests
  `_techspec.md` → Testing Approach names, the `./internal/baseline` and
  `./internal/cli` suites, `-tags docscontract ./internal/docscontract` and
  `GOOS=windows go vet`; each sabotage that section lists failed its named
  test.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix retired skill baseline update delete skills-lock" --all --files --min-score 0.3`.
  It returned only project mirrors (ADR-0246, the owned-skill version ADR,
  Spec 0241, the CHANGELOG and adopters' tooling records); none records a decision beyond
  ADR-0246.
- Exa: the FreeBSD `pkg` FAQ quoted above, and the GSD installer migration
  guide (https://github.com/open-gsd/gsd-core/blob/main/docs/installer-migrations.md),
  whose `remove-managed` action removes "a path only when it is known to be
  GSD-managed and unchanged from the previous manifest" and never "an unknown
  file", the same two limits this Spec sets.

## Glossary

- changes: **Retired Skill** — an approved update deletes a copy it proves unedited

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AU brief for Spec 0339 of 2026-11-01 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-01), verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- Retired skills (earlier the same day, in the Backlog Entry): "Apply apaga,
  valide se as skills context7 e exa estão fora da baseline e sendo
  removidas, deveriam ficar." Option text: "Nova Spec substitui a ADR-0246: o
  apply aprovado (--confirm-plan) apaga .agents/skills/<x>, o link em
  .claude/skills e a entrada no skills-lock.json."
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega
  aos adotantes, release v0.129.0 no fim e aviso aos adotantes."

Taken from those answers, ADR-0246, ADR-0191 and the characterization, with
no further interview:

- `context7-cli` and `exa-web-search` are not retired: every Setup Snapshot
  requires both, and the Retired Skills set never names them. The Vortex plan
  the operator read listed only the seven retired skills and
  `implement-spec`. A regression test pins both facts.
- An edited copy is kept with a warning, never deleted and never a refusal of
  the whole update: deleting it would lose the repository's own change, and
  refusing would block every Baseline change on one kept skill.
- A copy without a recorded hash is kept too: nothing proves it unedited.
  Both adopters surveyed record a hash for every retired copy.
- The deletion runs after the skills stage, on the approved digest, and
  re-inspects first; the Plan Digest is unchanged when nothing is deleted, so
  every existing plan and golden keeps its bytes.
- `--yes` approves the digest computed in one invocation, so it deletes too,
  and `--no-skills` skips the inspection and the deletion.
- The history-sanitize ADR's remark that a Refused Unit leaves the state
  unchanged "as a Retired Skill does not" stays true for a kept copy and is
  not amended.
- Two implementation Tasks in series (the code with its tests and the
  baseline guide, then the ADRs, the glossary and the Roundfix skill) and a QA gate: the disposition fields change the CLI's JSON and
  text in the same commit, so the baseline and CLI changes cannot settle
  apart.

## Open Questions

None.
