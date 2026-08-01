# QA command log

Build: `4493addcc79c4e6d40a2b95471536082e9c5d1b3`

## Repository Verification

Command: `rtk make verify`

Result: exit 0.

```text
Go test: 2941 passed in 24 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed
go build completed for ./cmd/roundfix
```

## Current workflow behavior

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-01/current_workflow_probe.rb
```

Result: exit 0. The probe executed the current YAML-embedded runtime guard,
tag and dispatch validation, classifier, Publication Preflight, and publish
helper. It passed the npm version boundaries; used, single-unpublish,
cooldown, eligible, malformed, absent, HTTP, transport, and multi-coordinate
cases; OIDC-first, enabled-fallback, and closed-fallback paths; token
non-exposure; all-six summary coverage; and platform-before-launcher ordering.

## Publish failure attribution

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/failure_attribution_probe.rb
```

Result: exit 0. The exact current publish helper emitted
`::error::publish: @roundfix/cli-darwin-arm64 failed for a non-identity reason`,
surfaced `npm ERR! network timeout while contacting registry`, invoked no
identity fallback, and exited the publish script 1.

## Governance, remediation, docs, and scope

Commands:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/docs_scope_probe.rb
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-01/remediation_scope_probe.rb
```

Result: both exited 0. The probes cover all seven completed Tasks and Results,
both Project Constraint sections, pre-Task authorization, chronological
ancestry, exact `git diff-tree` path scope for all seven Task/remediation
commits, package and asset compatibility, preserved release stages, amended
identity-preflight scope, all-six trusted-publisher setup, fallback close
evidence, registry-side token shutdown, confirmation, glossary, Release Plan
non-goal, and absence of cooldown retry.

## Failure-vocabulary finding

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-08-01-qa-01/failure_vocabulary_probe.rb
```

Result: exit 1, the expected reproduction of QA-004.

```text
FAIL: runbook omits workflow failure prefix(es): publish:
WORKFLOW_PREFIXES=identity,publish,registry,runtime,undetermined
RUNBOOK_PREFIXES=identity,registry,runtime,undetermined
```

## Live npm read path

Command:

```text
rtk ruby docs/specs/0058-npm-trusted-publishing-and-release-preflight/qa/evidence/2026-07-31-qa-01/live_preflight.rb
```

The managed sandbox could not resolve `registry.npmjs.org`. The exact current
script still iterated all six coordinates, classified each as `undetermined`,
printed six coordinate-specific errors, wrote six summary rows, stopped with
`PREFLIGHT_EXIT=1`, and the observation wrapper exited 0. The current
Publication Preflight step is byte-equivalent to build `b411b30`, whose
approved full-access run on 2026-07-31 reached npm, classified all six `0.0.2`
coordinates as `used`, and stopped before mutation.

## Pull Request fact

The Roundfix QA prompt states: `Pull Request: none open; Pull Request journeys
are environment-blocked.` The Run Worktree branch is never pushed and was not
queried as though it had a Pull Request.
