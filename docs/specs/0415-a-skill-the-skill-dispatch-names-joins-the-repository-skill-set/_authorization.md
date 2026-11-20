---
status: approved
granted: 2026-11-20
action: make the typescript Baseline module require the skill solid, refuse in catalog validation a Skill Activation whose bundle names a skill its owner does not require, read the required skills through one function, and record the rule in the glossary and in ADR-0448
consuming: 0415-a-skill-the-skill-dispatch-names-joins-the-repository-skill-set
paths:
  - docs/agents/setup-context.json
  - internal/baseline/assets/modules/typescript.json
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

# Approved authority for Spec 0415

The maintainer granted this authority on 2026-11-20. The Supervisor session asked for it in ai-memory message `01a1f96f-a274-7033-9765-4c60054e8d82`, sent to the Jarvis project, which put the question to the maintainer and returned his answer in ai-memory message `01a1f976-6dbf-7a43-93ee-89cc6806e754`. The Supervisor session popped that message from this project's inbox on 2026-11-20 at 11:57 BRT (14:57 UTC). The answer this record rests on is, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a1f96f-a274-7033-9765-4c60054e8d82
pergunta: Autoriza implementar, commit, push, PR, merge e soltar a versão da Spec 0415?
resposta: opção 1, "Tudo, com patch v0.152.1".
alcance: Spec 0415; caminhos protegidos internal/baseline/assets/modules/typescript.json e docs/agents/setup-context.json; regeneração pelos comandos do repositório; implement, commit, push, pull_request, merge e release com teto patch (CHANGELOG.md, dist/npm/roundfix/package.json, internal/app/version.go).
quando: 20/11/2026 13:12 BRT
```

Both of the maintainer's answers of that message, and what the planning author could check about it, are in [`references/2026-11-20-maintainer-answers.md`](references/2026-11-20-maintainer-answers.md). The block's `quando` reads 13:12 BRT while the message that carries it was created and popped at 11:57 BRT; the date is the same, and `granted` records it.

## What the grant covers

- The operations `implement`, `commit`, `push`, `pull_request`, `merge` and `release`.
- Two Governed Paths:
  - `internal/baseline/assets/modules/typescript.json`: one line, `"solid"` in `requiredSkills`, and the version line the record command writes;
  - `docs/agents/setup-context.json`: this repository's Setup Manifest, whose `catalogDigest` and `roundfixVersion` the Managed Refresh rewrites after the catalog changes.
- The regeneration by the repository's own commands, declared under `## Sanctioned regeneration`.
- A release at a `patch` Release Ceiling. The Spec is a `type: fix`: it adds no command, flag, field or output line. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge.

The grant names no number of corrections and no path beyond these. The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are not granted, and this Spec's design needs none of them. Anything outside this scope needs a new decision request.

## How the governed set was measured

The author measured it on 2026-11-20 on the branch `fix/dispatch-named-skills-join-the-skill-set` at `851024fa`, with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` that wrote nothing to the repository, over every file the TechSpec names, the files the regeneration commands write on a scratch clone, and five controls ([`references/2026-11-20-request-and-characterization.md`](references/2026-11-20-request-and-characterization.md), sections 5 and 6).

- The two paths above are the governed ones this Spec changes.
- These governed paths are not granted, and no Task may change them:
  - `internal/baseline/assets/skill-activations.json` and every other Baseline module, the `frontend` module included;
  - every template, profile, Source Baseline, Setup Snapshot and formatter golden under `internal/baseline/assets/`;
  - every guide under `docs/agents/`;
  - every owned skill, canonical or mirror, and every Vendored Skill;
  - the controls: the `Makefile`, `.github/workflows`, `go.mod` and `.roundfixrc.yml`.

## What is not governed

- `internal/baseline/catalog.go`, `internal/baseline/catalog_validate.go`, `internal/baseline/skills_restore.go`, `internal/cli/doctor.go` and the tests under `internal/baseline` and `internal/cli`.
- `internal/baseline/module-versions.json` and the six derived files under `internal/baseline/testdata/`.
- `CONTEXT.md`, ADR-0448 and this Spec's planning files.
- `docs/references/skill-coverage.json`, whose `command doctor` entry gets the Coverage Review that `_prd.md` → Skills declares, because task_01 changes `internal/cli/doctor.go`, a source of that Behavior Surface; and `docs/references/behavior-surfaces.json`, which no Task of this Spec is expected to change, since no command help, configuration key, exit code or user-guide page changes.

## Sanctioned regeneration

Repository-owned commands write the module's version line and the derived records. This declaration records the regeneration that follows the approved edit and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/module-versions.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

No version line, digest, golden or manifest field is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any owned or Vendored Skill, any guide under `docs/agents/`, any user-guide page;
  - `docs/references/coverage-record.json`, and any entry of `docs/references/skill-coverage.json` other than the Coverage Review of `command doctor`;
  - any Archive Record.
- No test or Task Verification:
  - reads another repository, except a temporary one the test itself created;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
