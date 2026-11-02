---
status: approved
granted: 2026-11-01
action: make an approved roundfix baseline update delete each Retired Skill whose copies it proves unedited against skills-lock.json, bind the deletions to the Plan Digest, keep every other skill, retire ADR-0246 for ADR-0350, and update the glossary, the baseline guide and the Roundfix skill's baseline reference
consuming: 0339-an-approved-baseline-update-deletes-the-retired-skills
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0339

On 2026-11-01 the maintainer answered the operator's cycle AU questions with
AskUserQuestion, in the `/roundfix --unattended` mode. The answers this Spec
rests on are, verbatim:

- Scope: "3 Specs (Recommended)" — option text: "0338 release: tipo do PR
  segue a Spec + teto de release no _authorization.md (os dois mexem no
  release plan e no write-*); 0339 apply aprovado apaga as aposentadas (nova
  ADR, context7-cli/exa-web-search ficam); 0340 links de evidência da QA
  existem no commit. Fila nessa ordem."
- Retired skills (asked earlier the same day, in the Backlog Entry): "Apply
  apaga, valide se as skills context7 e exa estão fora da baseline e sendo
  removidas, deveriam ficar." Option text: "Nova Spec substitui a ADR-0246:
  o apply aprovado (--confirm-plan) apaga .agents/skills/<x>, o link em
  .claude/skills e a entrada no skills-lock.json."
- Release of the cycle: "Uma no fim, v0.129.0 (Recommended)" — option text:
  "Depois das três Specs, com aviso aos adotantes; o apply que apaga
  aposentadas e o teto de release chegam juntos pela Baseline."
- Authorization: "Mesmas do AT (Recommended)" — option text: "Push e PR dos
  planos, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, mudança da Baseline
  (agent-instructions, release, spec-routing, skills de autoria) que chega
  aos adotantes, release v0.129.0 no fim e aviso aos adotantes."
- Start: "Já (Recommended)".

Those answers keep the cycle AT rules of 2026-11-01 ("Mesmas do AS
(Recommended)"), and through them the cycle AS answers of 2026-10-31
"Consertar e seguir (Recommended)" (option text: "Diagnostico, emendo o plano
ou crio uma Spec bloqueadora e entrego, sem teto fixo de corretivas. Só pulo
a Spec se o conserto exigir uma decisão de produto ou arquitetura que
contradiga o PRD aprovado; aí ela sai do ciclo e a release vai com as
demais, ou espera se for a 0333.") and the Tooling answer "Makefile,
.github/workflows, go.mod / go.sum, Config de lint/formatter", the answer of
2026-10-30 "Mesmas regras do AR (Recommended)" (option text: "Push de
branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
ponto de release e aviso os adotantes por ai-memory."), the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para
os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere
autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release, v0.129.0 after the three Specs of
the cycle, is the operator's, and the operator corrects the v0.127.0
CHANGELOG line in it.

The governed set was measured on 2026-11-01 on the branch `docs/plan-0339`,
based on `dd3b2b62`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, ADR-0350,
  every planning file of this Spec and its adopted reference.
- The three paths above are governed: the canonical Roundfix `baseline`
  reference and the `SKILL.md` pair whose version fields rise.

## What is not governed

- `internal/baseline/retired_skills_installed.go`,
  `internal/baseline/retired_skills_installed_test.go`,
  `internal/baseline/retired_skills_test.go`,
  `internal/cli/baseline_update.go`,
  `internal/cli/baseline_update_retired_skills_test.go`,
  `internal/cli/repository_profile_docscontract_test.go`.
- The mirror `skills/roundfix/references/baseline.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`,
  `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`,
  `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`,
  `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, `CHANGELOG.md`, any Baseline module, profile, Setup
  Snapshot, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  or the `### QA settlement` section of any skill.
- No deletion of any skill outside the Retired Skills set, in code or in this
  repository.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
