---
status: approved
granted: 2026-10-27
action: change the disposable-session profile proof and Fallback Chain candidate probe so the setup decides the verdict, the close shares the Proof Setup Deadline, an unclosed session is recorded and closed by the next proof, and a timed-out step is named; keep Baseline semantic analysis's stop; make the CLI count only a setup deadline as a timeout and report the warning; update the guides, the Roundfix skill references, the glossary and the Behavior Surface record
consuming: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
paths:
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0331

On 2026-10-27, at about 11:40 BRT, the maintainer answered the cycle AR
questions (AskUserQuestion) the operator asked in the `/roundfix follow-up`
mode, over an inventory read with `--messages` and validated with Jev
`typesafe/jev-1.13-20260917`. Verbatim:

- Scope: "A, B e C (Recommended)" — option text: "Lock + Clean vermelho;
  família do settle + aceite da QA; profile proof sob carga. Todas são health
  e duas se repetem em adotantes." This Spec is C, the profile proof under
  load.
- QA acceptance: "Casar por linha (Recommended)" — option text: "Uma
  declaração cobre cada criterion que nomeia; o accept confere linha a linha,
  não a contagem. write-prd/qa-gate passam a dizer isso." This answer governs
  Spec 0330.
- Authorization: "Mesmas regras do AQ (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."

The cycle AQ rules are the ones of 2026-10-26, about 15:00 BRT, when the
maintainer wrote "Sim, pode fazer a higiene." and "Vamos seguir a sua
sugestão de sequencia para a spec e os demais itens. Talvez seja importante
usar o jev para validar.": the Governed Paths each Spec declares, and still
not this repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, or
lint and formatter configuration. They carry the cycle AJ grant of
2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release stays at the maintainer's
discretion, by the recorded criteria.

The governed set was measured on 2026-10-27 on the branch `docs/plan-0331`,
based on `5e9475bf`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, `creates:` or
  `instruction:`, the new ADR and every planning file of this Spec.
- The four paths above are governed and listed in `paths`. The
  `domain-modeling` skill is governed too, but task_01 only reads it as an
  instruction and no Task changes it. Every other declared file is not
  governed, among them the production files and tests under
  `internal/agent/`, `internal/baselineacp/` and `internal/cli/`,
  `internal/testfixture/real_time_budget_test.go`, both user guide pages,
  the two `skills/roundfix/references/` mirrors,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and
  `docs/references/behavior-surfaces.json`.

## Why each governed path is needed

- `.agents/skills/roundfix/references/profiles.md` and
  `.agents/skills/roundfix/references/runtime.md` describe the profile proof
  and the `profiles:` readiness line, whose behavior this Spec changes
  (`docs/agents/cli.md`).
- The two Roundfix `SKILL.md` files carry the owned skill's version, which
  rises with its content.

## Sanctioned regeneration

The owned-skill version record is rewritten only by its record command,
after `make skills-sync` copies the canonical skill to its mirror:

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
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
  module, profile, template or snapshot asset, `docs/agents/`,
  `skills-lock.json`, any installed external skill, or the `### QA
  settlement` section of any skill.
- No file of another Spec's folder.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or `~/.acpx`;
  - starts acpx or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository or `~/dev/skills`;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
