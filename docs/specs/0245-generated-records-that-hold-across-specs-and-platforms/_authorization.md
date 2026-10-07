---
status: approved
granted: 2026-10-07
action: make the Baseline module version and the coverage record outputs of their record steps, so parallel Specs and other platforms regenerate them instead of colliding, and declare both as derived paths
consuming: 0245-generated-records-that-hold-across-specs-and-platforms
paths:
  - .roundfixrc.yml
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
  - docs/references/coverage-record.json
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/cli-surface.json
  - internal/baseline/assets/modules/external-triage.json
  - internal/baseline/assets/modules/monorepo.json
  - internal/baseline/assets/modules/rust.json
  - internal/baseline/assets/modules/tui-surface.json
  - internal/spec/coverage_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0245

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of Baseline sources, guides and `.roundfixrc.yml` the
maintainer said "Autorizar os dois". For the Governed Paths each Spec of the
cycle declares, the maintainer said "Concedo".

On 2026-10-07 the maintainer approved the order of the cycle, with this Spec
second after Spec 0244: "Pode seguir nessa ordem".

The governed set was measured with `GovernedPath` on the authoring branch at
`66f2d85b`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. 11 of them are governed.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/backend.json`, `cli-surface.json`,
  `external-triage.json`, `monorepo.json`, `rust.json` and `tui-surface.json`
  — each puts its top-level `version` on one line with four other keys. The
  line-scoped derived merge and the record step rewrite that line, so it must
  stand alone. Only the header line is split; the content digest is unchanged.
- `docs/agents/setup-context.json` — its catalog digest is the sanctioned
  derived fallout of the module bytes.
- `internal/spec/coverage_test.go` and `docs/references/coverage-record.json`
  — the collection and the record that follow the host today.
- `.roundfixrc.yml` — the derived declarations that run the record step at a
  merge conflict and re-record the coverage record.
- `docs/agents/specific-repository.md` — the repository rule that tells an
  author to run the module record step, beside the owned-skill rule.

## What is not governed

These are ordinary: `internal/baseline/module_versions_test.go`,
`internal/baseline/module-versions.json`, `internal/baseline/testdata/**`,
`internal/spec/coverage_platform_test.go`,
`internal/config/verification_tools_test.go`, `CONTEXT.md` and
`docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files. This
declaration records the regeneration that follows the approved edits and adds
no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes a changed module's version line and the
record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/bun.json
  - internal/baseline/assets/modules/cli-surface.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/external-triage.json
  - internal/baseline/assets/modules/frontend.json
  - internal/baseline/assets/modules/go.json
  - internal/baseline/assets/modules/monorepo.json
  - internal/baseline/assets/modules/repository-extension.json
  - internal/baseline/assets/modules/rust.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/tui-surface.json
  - internal/baseline/assets/modules/typescript.json
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden or generated guide is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to any module's content beyond splitting the six header lines;
  every module keeps its version number in this Spec.
- No test, Verification command or QA row opens a network connection or
  reads or writes the real `~/.roundfix`.
- No change to the `### QA settlement` section of any skill.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
