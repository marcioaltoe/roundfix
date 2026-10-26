---
status: approved
granted: 2026-10-26
action: move every bundled upstream Setup Snapshot to the published marcioaltoe/skills commit 8d71fade through the asset sync, regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands, and make the tests that pin the upstream revision read one constant with a new test that pins the republished onioncry tree
consuming: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
paths:
  - docs/agents/setup-context.json
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0320

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in the chat, in his own words:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` or later), then the
Spec Check probe fix, then the AP triage (the shared build cache first), then
the skill and tool studies. This Spec is the Baseline snapshot step; it
delivers after Spec 0319 and is independent of it.

The request it answers came from the upstream project
`default/marcioaltoe-skills` through ai-memory (messages `01a17807`,
2026-10-26 08:45 BRT, and `01a1781b`, 09:07 BRT); the adopted Backlog Entry
records them. Upstream's published `main` at authoring time is
`8d71fade68c54693bedbd917371f0514311af157` (pull request #135).

The cycle AQ rules are the same as cycle AP's: the Governed Paths each Spec
declares, and still not this repository's Makefile, `.github/workflows`,
`go.mod`, `go.sum`, or lint and formatter configuration. Those rules carry
the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push", and the 2026-09-30 grant of Baseline sources,
guides and `.roundfixrc.yml` ("Autorizar os dois").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on 2026-10-26 on the branch `docs/plan-0320`,
based on `86cb6b73`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Task declares under `interface:` or `creates:`,
  every planning file of this Spec and its adopted Backlog Entry copy.
- The five `internal/baseline/assets/setups/` files and
  `docs/agents/setup-context.json` are governed and listed in `paths`; a
  prototype that ran the sanctioned commands below in a scratch worktree
  changed exactly those six governed files.

## Why each governed path is unavoidable

- The five Setup Snapshots are rewritten by the asset sync at the new
  commit; the sync is their only writer (ADR-0191).
- `docs/agents/setup-context.json` records the catalog digest, which the new
  snapshots move; the managed refresh rewrites that one entry.

## What is not governed

- `internal/baseline/upstream_setup_revision_test.go`,
  `internal/baseline/ai_memory_setup_test.go` and
  `internal/baseline/python_setup_test.go`.
- `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest`, `catalog.normalized.json` and the four plan
  characterization goldens the Task declares.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files. This
declaration records the regeneration that follows the import and adds no
source path.

The asset sync writes the Setup Snapshots from a repository fetched,
read-only, from the local checkout `/Users/marcio/dev/skills` at
`8d71fade68c54693bedbd917371f0514311af157`, with the remote `origin`
`https://github.com/marcioaltoe/skills.git` added and never contacted; its
`--check` must then exit `0`:

```yaml
command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
```

```yaml
command: make baseline-digests
```

The managed refresh renders this repository's Setup Manifest, run from a
binary built under the Run's environment with a fresh `HOME`; a second run
reports no file change:

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

No snapshot, digest, golden or Setup Manifest value is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Go
  production file, any Baseline module, profile, template, coverage entry,
  Source Baseline or retention file, any installed skill under
  `.agents/skills/` or `skills/`, `skills-lock.json`,
  `skills/recommended.txt`, any rendered guide under `docs/agents/` other
  than `setup-context.json`, any user guide page, ADR or `CONTEXT.md`, or
  the `### QA settlement` section of any skill.
- No write to `/Users/marcio/dev/skills`: the import only fetches the
  published object from it.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads `~/dev/skills` (only task_01's import does) or another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
