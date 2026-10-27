---
status: approved
granted: 2026-10-26
action: change the Baseline plan, update, apply, reconcile and restore code and tests so a profile change carries the recorded decisions and refuses once, the update names the plan command, the apply names each Root Backup, reconcile explains a no-change result, and restore keeps the repository's Claude Skill Links; update the guides, the roundfix and setup-context-driven skills, the glossary and the Behavior Surface record
consuming: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
paths:
  - internal/cli/baseline_plan_test.go
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0328

On 2026-10-26, at about 22:35 BRT, the maintainer answered the cycle AQ
batch 5 questions (AskUserQuestion) the operator asked in the
`/roundfix --follow-up` mode about Quaestor's report, ai-memory message
`01a17abd-20fe-7213-ab5c-f1a1bdb043fe` from
`default/gesttione-solutions-quaestor` (2026-10-26 21:22 BRT). Verbatim:

- Grouping: "Duas Specs (Recomendado)" — 0327 Run lifecycle (executable
  scripts, settle exit 128, archive with `qa: declined`); 0328 Baseline
  profile-change frictions.
- Executables: "_authorization.md nomeia os executáveis (Recomendado)" —
  option text: "O commit da Task aceita 0755 só nos caminhos que a
  autorização da Spec lista; todo o resto continua recusado. O write-tasks e
  a recusa passam a citar essa regra." This answer governs Spec 0327; this
  Spec adds no executable file.
- Authorization: "Mesmas regras do ciclo AQ (Recomendado)" — option text:
  "implement, commit, push, PR e merge nos Governed Paths que cada Spec
  declarar; fora Makefile, .github/workflows, go.mod/go.sum e lint. A release
  continua a meu critério, pelos critérios gravados."

The cycle AQ rules are the ones of 2026-10-26, about 15:00 BRT, when the
maintainer wrote "Sim, pode fazer a higiene." and "Vamos seguir a sua
sugestão de sequencia para a spec e os demais itens. Talvez seja importante
usar o jev para validar.": the Governed Paths each Spec declares, and still
not this repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, or
lint and formatter configuration. They carry the cycle AJ grant of
2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grants of every skill
("considere autorizado a ajustar todas as skills se necessário") and of
Baseline sources and guides ("Autorizar os dois").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release stays at the maintainer's
discretion, by the recorded criteria.

The governed set was measured on 2026-10-26 on the branch `docs/plan-0328`,
based on `6a2630e5`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or `creates:`,
  the new ADR and every planning file of this Spec.
- The six paths above are governed and listed in `paths`. Every other
  declared file is not governed, among them the production files under
  `internal/baseline/` and `internal/cli/`, the new and existing tests other
  than `internal/cli/baseline_plan_test.go`, both user guide pages,
  `skills/roundfix/references/baseline.md`,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and
  `docs/references/behavior-surfaces.json`.

## Why each governed path is needed

- `internal/cli/baseline_plan_test.go` holds the plan characterization
  corpus whose missing-decision cases read the refusal message task_02
  extends; task_02 may only add to it.
- The four `SKILL.md` files and `.agents/skills/roundfix/references/baseline.md`
  carry the Roundfix and `setup-context-driven` owned skills, whose text must
  describe the changed commands (`docs/agents/cli.md`) and whose versions
  rise with it.

## Sanctioned regeneration

The owned-skill version record is rewritten only by its record command,
after `make skills-sync` copies the canonical skills to their mirrors:

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
```

The Behavior Surface record is rewritten only by its record command, run
after the last guide edit:

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No version field, fingerprint or digest is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template, Setup Snapshot, coverage or retention asset,
  `docs/agents/`, `skills-lock.json`, any installed external skill, or the
  `### QA settlement` section of any skill.
- No file of another Spec's folder.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository or `~/dev/skills`;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
