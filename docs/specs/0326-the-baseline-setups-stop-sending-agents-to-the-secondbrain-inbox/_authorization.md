---
status: approved
granted: 2026-10-26
action: move every bundled upstream Setup Snapshot to the published marcioaltoe/skills commit b56fc5ff through the asset sync, restore this repository's knowledge-workspace skill and its lock entry from the same commit through the confirmed restore command, regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands, and move the shared revision test constant with a new test that pins the corrected knowledge-workspace tree
consuming: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
paths:
  - .agents/skills/knowledge-workspace/SKILL.md
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

# Approved authority for Spec 0326

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in the chat, in his own words:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The cycle AQ batch 4 (2026-10-26, about 20:20 BRT) runs under the same
approval and rules as batches 1 to 3. It answers Pantheon's report that the
vendored `knowledge-workspace` skill still sends agents to the retired
Secondbrain inbox: ai-memory message `01a179b6` from
`default/gesttione-solutions-pantheon` (2026-10-26 16:36 BRT), and item 3 of
message `01a1795f` from `default/marcioaltoe-skills` (2026-10-26 15:01 BRT);
the adopted Backlog Entry records both. Upstream fixed the skill in
marcioaltoe/skills#136, squash
`b56fc5ffcad1ab366fbc03839a6ce8c662401178`, upstream `main` at authoring
time. The Spec is enqueued after Spec 0322 and ships with the next release
the operator cuts.

The cycle AQ rules are the same as cycle AP's: the Governed Paths each Spec
declares, and still not this repository's Makefile, `.github/workflows`,
`go.mod`, `go.sum`, or lint and formatter configuration. Those rules carry
the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push", and the 2026-09-30 grants of Baseline sources,
guides and `.roundfixrc.yml` ("Autorizar os dois") and of every skill
("considere autorizado a ajustar todas as skills se necessário").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on 2026-10-26 on the branch `docs/plan-0326`,
based on `f06afdd7`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Task declares under `interface:`, every
  planning file of this Spec and its adopted Backlog Entry copy.
- The five `internal/baseline/assets/setups/` files,
  `docs/agents/setup-context.json` and
  `.agents/skills/knowledge-workspace/SKILL.md` are governed and listed in
  `paths`; a prototype that ran the sanctioned commands below in a scratch
  clone changed exactly those seven governed files.

## Why each governed path is unavoidable

- The five Setup Snapshots are rewritten by the asset sync at the new
  commit; the sync is their only writer (ADR-0191).
- `.agents/skills/knowledge-workspace/SKILL.md` is this repository's
  installed copy of a required external skill, held to the `go` snapshot
  (ADR-0221); only the confirmed restore writes it (ADR-0264).
- `docs/agents/setup-context.json` records the catalog digest, which the new
  snapshots move; the managed refresh rewrites that one entry.

## What is not governed

- `internal/baseline/upstream_setup_revision_test.go` and `skills-lock.json`.
- `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `catalog.digest`, `catalog.normalized.json` and the four plan
  characterization goldens the Task declares.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files. This
declaration records the regeneration that follows the import and adds no
source path.

The asset sync writes the Setup Snapshots from a repository fetched,
read-only, from the local clone
`/Users/marcio/.roundfix-operator/scratch/spec-0326/skills` at
`b56fc5ffcad1ab366fbc03839a6ce8c662401178`, with the remote `origin`
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

The restore takes this repository's `knowledge-workspace` from the same
fetched repository; its preview must list exactly the two outputs below
before the confirmed apply:

```yaml
command: roundfix baseline skills restore --profile go-cli-tui --skill knowledge-workspace --source-dir <fetched> --repo . --confirm-plan <digest> --format text
outputs:
  - .agents/skills/knowledge-workspace/SKILL.md
  - skills-lock.json
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

No snapshot, installed skill, lock entry, digest, golden or Setup Manifest
value is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Go
  production file, any Baseline module, profile, template, coverage entry,
  Source Baseline or retention file, any installed skill under
  `.agents/skills/` or `skills/` other than
  `.agents/skills/knowledge-workspace/SKILL.md`, any `skills-lock.json` entry
  other than `knowledge-workspace`'s, `skills/recommended.txt`, any rendered
  guide under `docs/agents/` other than `setup-context.json`, any user guide
  page, ADR or `CONTEXT.md`, or the `### QA settlement` section of any skill.
- No write to the operator's clone
  `/Users/marcio/.roundfix-operator/scratch/spec-0326/skills`, and no read or
  write of `/Users/marcio/dev/skills`: the import only fetches the published
  object from the clone.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads the operator's clone (only task_01's import does), `~/dev/skills`
    or another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
