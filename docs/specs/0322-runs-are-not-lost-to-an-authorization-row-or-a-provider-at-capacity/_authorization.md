---
status: approved
granted: 2026-10-26
action: make the shared QA partial policy accept an outside-evidence row recorded as blocked (environment <subject> prohibited by Spec authorization) with a bounded subject and count each row once, let the next Fallback Selection take any Task or QA turn after a Provider System Error and reopen a QA verdict the failed turn wrote, reword the two Baseline clauses that name the exemption and regenerate their derived files and this repository's managed guides, and document both rules in the qa-gate, archive-spec and roundfix skills, the user guides and CONTEXT.md
consuming: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/profiles.md
  - .agents/skills/roundfix/references/settle.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/archive-spec/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0322

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in chat. The answers this Spec rests on are, verbatim:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` and later) → the Spec
Check probe fix → the cycle AP triage (the shared build cache in `gc` first)
→ skill and tool studies. This Spec belongs to the cycle AP triage of cycle
AQ batch 2, validated with Jev (`follow-up --messages`,
`jev-1.13-20260917`): two sources whose stops cost Runs that were not the
work's fault. It delivers after Specs 0319 and 0320 and is independent of
them.

The cycle AQ rules carry the cycle AP rules unchanged: the Governed Paths each
Spec declares, under the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push", and the skills
authorization "considere autorizado a ajustar todas as skills se necessário".
Still not authorized: this repository's Makefile, `.github/workflows`,
`go.mod`, `go.sum`, and lint and formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on the authoring branch `docs/plan-0322`, based
on `d08f8557`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, `creates:` or
  `instruction:`, every planning file of this Spec, its two adopted Backlog
  Entry copies, ADR-0335 and the amended ADR-0274 and ADR-0276.
- The seventeen files in `paths` are governed and changed by the Tasks.
  `.agents/skills/domain-modeling/SKILL.md`, `internal/baseline/plan_test.go`,
  `internal/docscontract/publicdocs_test.go` and
  `internal/spec/coverage_test.go` are governed but only read
  (`instruction:`), so they are not granted.

## What is not governed

- `internal/spec/qa.go`, `internal/spec/qa_prohibited_subject_test.go`,
  `internal/cli/qa_partial_policy_test.go`,
  `internal/daemon/qa_partial_policy_test.go`,
  `internal/daemon/agent_session_owner.go`,
  `internal/daemon/provider_system_error_recovery_test.go`,
  `internal/daemon/provider_system_error_unfinished_qa_test.go`.
- `internal/baseline/module-versions.json`,
  `internal/baseline/network_prohibited_row_clause_test.go`, the three
  catalog files under `internal/baseline/testdata/` and the four plan
  characterization goldens under
  `internal/baseline/testdata/plan-characterization/`.
- `skills/roundfix/references/archive.md`,
  `skills/roundfix/references/settle.md`,
  `skills/roundfix/references/profiles.md`,
  `skills/testdata/owned-skill-versions.json`,
  `skills/acceptance_rules_test.go`.
- `docs/user-guide/commands/settle.md`, `docs/user-guide/commands/qa-report.md`,
  `docs/user-guide/commands/archive.md`,
  `docs/user-guide/context-driven-development.md`,
  `docs/user-guide/usage.md`, `docs/user-guide/configuration.md`,
  `CONTEXT.md`, `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec, its adopted Backlog Entry copies and the
  three ADRs.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

The Module Version Record step writes the two modules' version lines and the
record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/autonomous-work.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/spec-routing.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

The Managed Refresh runs with `HOME` at a temporary directory, and a second
refresh reports no file change. No digest pin, golden, generated guide,
mirror or record is hand-edited. No top-level test is removed or renamed, so
the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module other than `spec-workflow` and `autonomous-work`, the Source
  Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, the Setup Snapshots, any skill
  other than `qa-gate`, `archive-spec` and `roundfix`, `internal/agent`, any
  ADR during delivery, or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
