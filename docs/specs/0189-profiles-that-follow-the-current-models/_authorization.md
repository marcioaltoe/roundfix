---
status: approved
granted: 2026-09-30
action: refresh the Model Catalog, the picker efforts and the adapter floors to what the installed adapters advertise, replace the recommendation ranking with one dated Recommended Profile, derive the built-in profiles, the generated config, the legacy Codex default and the Baseline analysis models from it, and refresh the model reference document and the Roundfix skill to match
consuming: 0189-profiles-that-follow-the-current-models
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/agents/openai.yaml
  - skills/roundfix/agents/openai.yaml
  - internal/cli/cli_test.go
  - internal/docscontract/publicdocs_test.go
  - .roundfixrc.yml
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0189

On 2026-09-30 the maintainer approved the program this Spec belongs to
("Três ondas") and decided to adopt the current models directly, fallbacks
included ("Trocar direto", then "Atualize os fallbacks também"). The same day
the maintainer expressly authorized:

- the skill files: "considere autorizado a ajustar todas as skills se
  necessário";
- `.roundfixrc.yml`, answering "Autorizar os dois" to a structured question
  that named the Baseline source and guides and `.roundfixrc.yml`.

The two governed test files ride the standing grant of 2026-09-21 for governed
paths a slice genuinely needs. The set was measured with `GovernedPath` on
`9e439dbb`, through a `go test -overlay` probe that wrote nothing to the
repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the skill
  states the adapter floors and their install commands (task_01), the
  `profiles show` schema and its recommendation rows (task_02), and the
  built-in profiles (task_03). `.agents/skills/` is canonical and `skills/` its
  mirror.
- `.agents/skills/roundfix/agents/openai.yaml`,
  `skills/roundfix/agents/openai.yaml` — the manifest names the Codex adapter
  version a legacy override migrates to, which is the floor (task_01).
- `internal/cli/cli_test.go` — it pins the `profiles show` output and schema
  and the rule that every ranked model is in the catalog (task_02), and the
  built-in fallback in seven profile and Doctor tests (task_03). `TestRunCommandHelp` is not
  changed, so `docs/references/coverage-record.json` does not move.
- `internal/docscontract/publicdocs_test.go` —
  `TestProfilesDocumentationContractMatchesPublicGuidance` requires the guides
  and the skill to contain `2026-08-07`, `roundfix/profiles/v1`,
  `category_specific: false` (task_02), and `gpt-5.5` and `claude-fable-5`
  (task_03). Those strings
  leave the documents with this Spec.
- `.roundfixrc.yml` — one comment says the reference document still holds the
  2026-08-07 snapshot. task_04 corrects that comment and changes no key.

## What is not governed

`internal/agent/catalog.go`, `internal/agent/acpx_runner.go`,
`internal/cli/cli.go`, `internal/cli/profiles.go`,
`internal/cli/profiles_configure.go`, `internal/config/recommendations.go`,
`internal/config/profiles.go`, `internal/config/config.go`,
`internal/baselineacp/analyzer.go`, every other test file the Tasks name, the
new test files, `docs/user-guide/commands.md`,
`docs/user-guide/configuration.md`, `docs/user-guide/usage.md`,
`docs/references/model-selection.md`, `CONTEXT.md` and
`docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md` are
ordinary. `internal/baseline/assets/**`, the Baseline fixtures and golden
files that name `gpt-5.5` as a decision value, `Makefile`, the CI workflows and
`go.mod` are governed and are not touched.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No profile key in `.roundfixrc.yml` changes; #292 already set them.
- No Baseline module, decision default, fixture or golden file changes.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `skills/_ownership.yml` or the `### QA settlement` section of any skill.
- No test reaches a real ACP adapter, a provider or the network, and the live
  Run Database under `~/.roundfix` is never opened for writing.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
