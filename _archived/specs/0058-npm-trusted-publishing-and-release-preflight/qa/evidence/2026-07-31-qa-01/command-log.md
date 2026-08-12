# QA command log

Build: `b411b30b6ad4c55ff0198e6c39ae56a7b4e4a8b4`

## Repository Verification

Command:

```text
rtk make verify
```

Result: exit 0.

```text
Go test: 2941 passed in 24 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd,
write-techspec, write-tasks, setup-context-driven, implement-task,
implement-spec, brainstorming, council, business-analyst, archive-spec,
qa-gate, evidence-gate
go build completed for ./cmd/roundfix
```

## Workflow behavior matrix

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/workflow_probe.rb
```

Result: exit 0. The probe executed the YAML-embedded runtime guard, tag and
dispatch validation, classifier, Publication Preflight, and publish helper.
It passed the npm 11.5.0/11.5.1/12.0.0 boundaries; used, single-version
unpublish, cooldown, eligible, malformed, absent, HTTP, transport, and
multi-coordinate cases; all-OIDC, enabled-fallback, and closed-fallback
publish paths; token non-exposure; all-six summary coverage; and five-platform
before launcher ordering.

## Live npm read path

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/live_preflight.rb
```

The managed sandbox run could not resolve `registry.npmjs.org`; the approved
full-access rerun reached npm and exited the observation wrapper 0 with
`PREFLIGHT_EXIT=1`. The exact workflow script classified all six `0.0.2`
coordinates as `used`, printed six `registry:` errors, wrote all six summary
rows, and stopped without mutation.

## Governance, docs, and scope

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/docs_scope_probe.rb
```

Result: exit 1 at the deliberate PRD coverage assertion after passing the Task
completion, Project Constraint, authorization ancestry, exact `diff-tree`
scope, unchanged package/command paths, preserved Verification/build/release
scripts, release coordinate and asset, trusted-publisher setup, fallback,
failure-vocabulary, glossary, no-retry, and Release Plan non-goal checks.

```text
FAIL: runbook never instructs the maintainer to disallow token publication for the owned packages
```

## Publish failure attribution

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/failure_attribution_probe.rb
```

Result: exit 1. The exact publish helper received a simulated registry network
timeout, emitted `::warning::identity:` for
`@roundfix/cli-darwin-arm64`, and retried with the bounded token fallback.

```text
FAIL: a registry transport failure is classified as identity and retried with the token
PUBLISH_EXIT=1
```

## Pull Request fact

The Roundfix QA prompt states: `Pull Request: none open; Pull Request journeys
are environment-blocked.` The Run Worktree branch is never pushed and was not
queried as though it had a Pull Request.
