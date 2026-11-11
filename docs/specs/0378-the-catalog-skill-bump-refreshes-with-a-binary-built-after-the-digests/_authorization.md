---
status: approved
granted: 2026-11-10
action: make the Catalog Skill Bump build its Roundfix binary again after its make baseline-digests step, so the Managed Refresh it runs records the catalog digest of the tree on disk, prove the order in the Catalog Skill Bump Makefile contract, and regenerate the Setup Manifest that records the Makefile digest
consuming: 0378-the-catalog-skill-bump-refreshes-with-a-binary-built-after-the-digests
paths:
  - Makefile
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0378

On 2026-11-10 the operator recorded that the Catalog Skill Bump leaves a
stale catalog digest in the Setup Manifest
(`references/2026-11-10-the-catalog-skill-bump-leaves-a-stale-catalog-digest-in-the-setup-manifest.md`,
recorded by the operator). This Spec reproduced it with a local catalog
repository and adopts the entry.

The maintainer's answers this Spec rests on (AskUserQuestion via
`/roundfix --follow-up`, 2026-11-10 ~20:55–21:00 BRT) are, verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among four, "Digest
  do bump de skills" — option text: "O bump de skills deixa um digest velho
  no Setup Manifest e exige um comando manual. Só aparece quando alguém roda
  o bump, o que é raro."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

"Mesmas do BB" keeps the cycle BB rules of 2026-11-10 ("Mesmas do BA, teto
minor (Recommended)" — "Escrever, entregar pela fila, merge e release de
cada Spec assim que pronta (patch nas duas correções; a regra pode sair como
minor), teto de 2 Tasks corretivas, eventos ao Jarvis, sem notas às caixas,
e Makefile/CI/go.mod/lint só onde a Spec exigir."), and through them the
cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).").

The design needs the `Makefile`, which the answer grants to this Spec ("Makefile na 0378 só nos arquivos que elas exigirem"): one line in the
`skills-third-party-bump` recipe. `docs/agents/setup-context.json` records
the `Makefile` digest in its verification entries, so the `Makefile` edit
makes it stale and the Managed Refresh regenerates it in the same Task. No
other tooling, workflow, Go module, lint or formatter file changes.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0378`,
based on `b233761e` and rebased on `b6cd0f14`, which changes none of the
probed files, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  the Backlog Entry's former path.
- The two paths above are governed. The Catalog Skill Bump Makefile
  contract `internal/verifyselect/skills_bump_makefile_test.go` and every
  planning file are not.

## Sanctioned regeneration

The Managed Refresh writes the Setup Manifest. This declaration records the
regeneration that follows the approved `Makefile` edit and adds no source
path.

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

The command runs twice; the second run reports `File changes: 0`. No digest
is hand-edited. No Baseline module, profile or derived Baseline asset
changes, so `make baseline-digests` is not run, and no top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile outside the `skills-third-party-bump` recipe, to
  the CI workflows, `go.mod`, `go.sum`, the lint, formatter or test-runner
  configuration, `.roundfixrc.yml`, `skills-lock.json`, any file under
  `.agents/skills/` or `skills/`, any Baseline module, profile, template or
  derived Baseline file other than the Setup Manifest, any guide under
  `docs/agents/`, `CONTEXT.md`, any ADR, any Archive Record, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - runs a real Catalog Skill Bump, or reads the skills catalog or any other
    repository;
  - reads or writes the real `~/.roundfix`;
  - creates files in `/tmp`, `/private/tmp` or `$HOME` outside its own
    temporary directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network, or fetches from any remote other than a local
    `file://` source it created;
  - sends a request to a provider, the judge or a reviewer;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
