---
status: approved
granted: 2026-10-26
action: trim every declared Shared Build Cache automatically at the start of every Implement Run and before the Delivery Queue's repository gate, under a thirty-second budget that never changes their outcome, record each trim as a Run Event or a Queue Owner log line, add a doctor disk-space warning below 10 GiB free, and document both in the guides, the Roundfix skill references, ADR-0334 and CONTEXT.md
consuming: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/storage.md
  - internal/cli/cli_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0321

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in chat. The answers this Spec rests on are, verbatim:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` and later) → the Spec
Check probe fix → the cycle AP triage, with the shared build cache in `gc`
first → skill and tool studies. This Spec is that first triage item. The
operator's triage of 2026-10-26, at about 15:30 BRT, validated with Jev,
proposed it as "0321 gc automático + doctor espaço baixo"; it delivers after
Specs 0319 and 0320 and is independent of the other cycle AQ Specs.

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

The governed set was measured on the authoring branch `docs/plan-0321`, based
on `d08f8557`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or `creates:`,
  ADR-0334, every planning file of this Spec and its adopted Backlog Entry
  copy.
- `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/storage.md`,
  `internal/cli/cli_test.go` and `skills/roundfix/SKILL.md` are governed and
  listed in `paths`.

## What is not governed

- `internal/sharedcache/sharedcache.go`,
  `internal/sharedcache/sharedcache_test.go`,
  `internal/cli/gc_shared_cache.go`, `internal/cli/gc.go`,
  `internal/cli/implement.go`, `internal/cli/deliver_workflow.go`,
  `internal/cli/shared_cache_trim_test.go`.
- `internal/store/available_bytes.go`,
  `internal/store/available_bytes_test.go`,
  `internal/cli/doctor_disk_space.go`,
  `internal/cli/doctor_disk_space_test.go`, `internal/cli/doctor.go`,
  `internal/cli/health.go`, `internal/cli/cli.go`.
- `docs/user-guide/commands/gc.md`, `docs/user-guide/commands/doctor.md`,
  `docs/user-guide/commands/implement.md`,
  `docs/user-guide/commands/deliver.md`, `docs/user-guide/configuration.md`.
- `skills/roundfix/references/deliver.md`,
  `skills/roundfix/references/implement.md`,
  `skills/roundfix/references/setup.md`,
  `skills/roundfix/references/storage.md`,
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical `roundfix` skill to
  `skills/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`,
  `.agents/skills/roundfix/SKILL.md`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module or derived Baseline file, `docs/agents/`, any skill other than
  `roundfix`, any ADR other than the new ADR-0334, or the
  `### QA settlement` section of any skill.
- No change to ADR-0327's trim rules, to `worktree.shared_caches` or its
  validation, to Worktree Build Caches, or to which Runs, Tasks or stages run.
- No new park class, configuration key, command, flag or exit code.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or the real shared Go build cache;
  - fills a disk or depends on the host's free space;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
