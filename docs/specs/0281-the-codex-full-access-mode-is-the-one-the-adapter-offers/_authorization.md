---
status: approved
granted: 2026-10-19
action: select the full-access mode from the modes the Agent Session advertises, send the Codex sandbox preset when the session offers it, report a degraded effective access policy when no full-access mode is advertised, and describe the change in ADR-0011's amendment note, the configuration guide, the Roundfix runtime reference, the .roundfixrc.yml comment, the glossary and the qa-gate skill
consuming: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/profiles.md
  - .roundfixrc.yml
  - internal/cli/cli_test.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0281

On 2026-10-19, at about 16:20 BRT, the maintainer approved cycle AJ through
AskUserQuestion. The answers were:

- decomposition: "Três Specs (Recomendado)", with this Spec as J, delivered
  before K (0282) and L (0283), one minor each;
- design of J: "Descobrir pelos availableModes (Recomendado)" (option text:
  "Escolher o modo pela lista da sessão ACP (hoje 'agent-full-access'),
  manter o preset sandbox danger-full-access, atualizar o comentário do
  .roundfixrc gerado. NÃO ligo agent_full_access neste repo; o QA usa o
  chrome-devtools MCP.");
- authorizations, all four selected: "Religar delivery.linux_lane: true",
  "docker pull/build da imagem", "Mexer em Governed Paths declarados" (option
  text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
  registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
  pedem autorização específica).") and "Merge squash + release + push";
- costs "Regras do ciclo AH/AI (Recomendado)", failures "Regras AE–AH e parar
  na proposta (Recomendado)", authoring "Opus + strict check, enfileirar sem
  esperar você (Recomendado)".

The 2026-09-30 skills grant ("considere autorizado a ajustar todas as skills
se necessário") also stands. The operator, not this Spec, releases, so the
operations above do not include `release`. This Spec uses neither the Linux
Lane nor Docker.

The governed set was measured on the authoring branch `docs/plan-0281`, based
on `f382480b`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and the ADR the plan writes.
- Seven of those files are governed, and all seven are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md` and `skills/qa-gate/SKILL.md` are the
  `qa-gate` skill and its mirror. They gain the heading `### Browser driver
  inside a sandbox`, with the version fields they carry (task_02).
- `.agents/skills/roundfix/references/profiles.md` is the Roundfix reference
  for Agent Selection Profiles and their readiness proofs, which gains the
  full-access selection (task_01).
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with that content (task_01).
- `.roundfixrc.yml` holds the comment that cites codex-acp 1.1.9. Only the
  comment above `agent_full_access` changes; `agent_full_access: false` stays
  (task_02).
- `internal/cli/cli_test.go` holds `TestRunResolveAgentFullAccessIsExplicitOptIn`,
  which pins the Codex full-access name `full-access` as `wantAccess`; it
  moves to `agent-full-access`, the first Codex candidate (task_01).

## What is not governed

- The production and test sources the Tasks change or create in
  `internal/agent`, and `internal/cli/implement_test.go`.
- The reference mirror `skills/roundfix/references/profiles.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, `docs/user-guide/configuration.md`, ADR-0011 and ADR-0293.
- Every planning file of this Spec, including its adopted reference.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and
no platform set changes, so the coverage record is not re-recorded. No record
or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `docs/agents/setup-context.json`, the pinned adapter versions, or the
  `### QA settlement` section of any skill.
- `defaults.agent_full_access` stays `false` in `.roundfixrc.yml`, and
  `delivery.linux_lane` is not touched.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or `~/.acpx`;
  - starts acpx, codex-acp, claude-agent-acp or a browser;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
