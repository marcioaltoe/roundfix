---
status: approved
granted: 2026-11-17
action: record each Correction of a delivery with a declared Defect Class and reason, make the declaration required before the Delivery Queue's merge only where a repository's Project Config turns it on, record the model and reasoning effort of an Authoring Phase in its phase entry, add the read-only Corrections Report, and describe all of it in the Roundfix skill, the four authoring skills, the guides, the glossary and the amended ADR-0407
consuming: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/phase.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0404

The maintainer asked for this Spec on 2026-11-17, answered the questions of its Idea and approved the Idea the same day. The request and the answers are recorded whole in [`references/2026-11-17-maintainer-request-and-answers.md`](references/2026-11-17-maintainer-request-and-answers.md). The answers this record rests on are, verbatim as the operator recorded them:

- The route of the request, the record of the cause of each correction (2026-11-17): "Spec já".
- Authorization (2026-11-17): "Mesmas do BE (Recommended)": push, PR, merge, entra na release única se pronta a tempo; até 5 correções; autoria e entrega depois da 0403.
- Whether the class is required (AskUserQuestion, 2026-11-17 ~14:45 BRT), in his own words: "Só onde é ligado e por default é desligado, nos projetos da gesttione vai ser configurado para ligar".
- This repository and the fifth class (AskUserQuestion, 2026-11-17 ~15:30 BRT): "Ligar e manter unclear (Recommended)".
- The approval of the Idea (AskUserQuestion, 2026-11-17 ~15:30 BRT): "Aprovar e seguir (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it. These include the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário") and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)." The records of Specs 0398, 0402 and 0403 quote the same rules.

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them, and it changes no Baseline module and no guide under `docs/agents/`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: feat`: it adds two commands, a Project Config key, keys of the Phase Record and a blocker of the Delivery Queue. Its Release Ceiling is `minor`, the ceiling the cycle BE rules set for a release after a Spec, and the release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, alone or as the single release of the cycle when this Spec is ready in time, with up to five corrections included.

## The setting in this repository

The maintainer's answer "Ligar e manter unclear (Recommended)" turns the new setting on in this repository. No Task of this Spec changes `.roundfixrc.yml`, so the file is not in `paths`. The installed Roundfix that runs this repository's Delivery Queue refuses a Project Config key it does not know (ADR-0408), so the key cannot arrive in the change that teaches the binary that key. The operator adds the one key, in a commit of its own, once the release that carries this Spec is the installed binary on the machine that runs this repository's Roundfix. That answer is the express authorization for that one-key change and for nothing else in the file.

## How the governed set was measured

The author of the PRD measured the governed set on 2026-11-17 in a clone on the branch `docs/plan-0404`, based on `74d0e3a9`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. The probe file never existed in the tree, and `git status --porcelain` printed nothing afterwards.
- It covered 64 paths: the skill files, guides, records, source files and documents the PRD expects the Tasks to change, this Spec's planning files, ADR-0407 and the three new ADRs, and controls. The result is recorded in [`references/2026-11-17-item-branch-commits-and-author-session-measurement.md`](references/2026-11-17-item-branch-commits-and-author-session-measurement.md).
- The fourteen paths above are the governed ones the PRD's design changes: the `SKILL.md` pair of the Roundfix skill and of each of the four authoring skills, and the Roundfix skill's canonical `deliver`, `phase`, `review` and `spec-delivery` references.
- These governed paths are not granted, and no Task may change them:
  - `.roundfixrc.yml`, as the section above says;
  - `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go` and `docs/references/coverage-record.json`; new tests go in new files;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.

This is the measurement a PRD can make. The TechSpec names every file the design changes, and `write-tasks` measures the paths the Tasks declare on a scratch copy. When either finds a governed path outside `paths`, such as another reference of the Roundfix skill, the operator amends this record before the plan merges, outside the change it authorizes.

## What is not governed

- Every Go source and test file the Tasks declare, except the governed tests above.
- The mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json` and the tests under `skills/`.
- The generated subagents under `.claude/agents/` and `skills-lock.json`.
- The user guide, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, ADR-0407, ADR-0417, ADR-0418 and ADR-0419.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited. The TechSpec confirms this list against the design, and adds the installation that regenerates the subagents under `.claude/agents/` if the design changes what installation writes into them.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot, and any file under `docs/agents/`;
  - any other skill, Vendored Skill or Archive Record;
  - the `### QA settlement` section of any skill;
  - the Classification Signatures, the Cause Classes and the text and JSON of `roundfix runs causes`.
- No change to any file of another repository, and no read of one. No adopter's Project Config is changed.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex` or `claude`;
  - reads or writes the real `~/.roundfix`, `~/.claude` or the real home directory, the real Phase Record included;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
