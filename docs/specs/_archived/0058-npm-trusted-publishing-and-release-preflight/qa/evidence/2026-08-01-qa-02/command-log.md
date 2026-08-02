# Command log — 2026-08-01 QA rerun 02

Build: `04d5bbb99f63f3343b29a8ed9e7d6b93923cd93d`

The full QA matrix was written with all 29 rows pending before executing the
static gate or public-surface probes.

## Static gate

- `rtk make verify` — exit 0. The gate reported 2,941 passing Go tests in 24
  packages, four passing focused Skill/Baseline tests, a passing
  `roundfix skills check`, and a completed CLI build.

## Current-build probes

- `rtk ruby qa/evidence/2026-08-01-qa-02/workflow_probe.rb` — exit 0. It
  exercised the exact YAML-embedded validation, runtime guard, classifier,
  preflight, and publish scripts. All assertions passed, including the six
  coordinate Release Set, blocked and undetermined states, OIDC-first calls,
  open and closed fallback behavior, token-sentinel non-disclosure, and
  platform-before-launcher order.
- `rtk ruby qa/evidence/2026-08-01-qa-02/remediation_scope_probe.rb` — exit 0.
  All seven Tasks, Project Constraints, authorization ancestry, bounded Task
  paths, amended identity boundary, and registry/repository close procedure
  passed.
- `rtk ruby qa/evidence/2026-08-01-qa-02/failure_vocabulary_probe.rb` — exit 0.
  Both the workflow and runbook expose `identity`, `publish`, `registry`,
  `runtime`, and `undetermined`; this closes prior QA-004.
- `rtk ruby qa/evidence/2026-08-01-qa-02/failure_attribution_probe.rb` — exit
  0. A simulated npm network timeout emitted `publish:` plus the raw npm
  cause, exited nonzero after one call, and never entered the identity/token
  fallback.

The historical `qa/evidence/2026-07-31-qa-01/workflow_probe.rb` exited 1 at its
old fallback assertion after all earlier assertions passed. That probe injects
a generic first-publish failure and predates Task 06, which intentionally
requires `ENEEDAUTH`, `E401`, or `Unable to authenticate` evidence before a
fallback. The current-build successor above injects the required
authentication signal and passes the open-window, closed-window, and
non-authentication branches.

## Tooling and scope audit

- Fresh `git diff-tree --no-commit-id --name-only -r` checks show each tooling
  commit changes only `.github/workflows/release.yml`, its assigned Task file,
  and Task 02's Spec-owned fixtures where applicable: `21bc4bf`, `8d14a67`,
  `b0052e9`, `47de307`, and `ab34e03`.
- Authorization commit `397227ff` is an ancestor of the first tooling Task.
  The later attribution repair `ab34e03` follows the fallback Task `47de307`.
  The QA-004 documentation repair `04d5bbb` follows `ab34e03`, changes only
  `docs/user-guide/release-runbook.md`, and is not folded into a tooling Task.
- `origin/ma/npm-trusted-publishing-and-release-preflight` resolves to the
  exact tested build `04d5bbb99f63f3343b29a8ed9e7d6b93923cd93d`.
- The full local gate and the current probes confirm the package names, assets,
  Upgrade Command compatibility, immutable tag authority, release ordering,
  and Non-Goals.

## Live publish-free rehearsal

- `rtk gh workflow view release.yml --ref
  ma/npm-trusted-publishing-and-release-preflight --yaml` — exit 0 with full
  network access. The published workflow at the target ref matches the tested
  OIDC permission, runtime guard, dispatch input, preflight, and push-only
  mutation guards.
- `rtk gh workflow run release.yml --ref
  ma/npm-trusted-publishing-and-release-preflight -f version=v0.0.2` — exit 0
  and created GitHub Actions run `30699898430`.
- `rtk gh run watch 30699898430 --exit-status` — exit 1, the expected workflow
  outcome for an ineligible used version. Setup, checkout, Go/Node setup,
  runtime guard, dispatch validation, and remote `make verify` passed.
  Publication Preflight reported all six `0.0.2` coordinates as `used` with a
  `registry:` error for each. Cross-compilation, npm publication, and GitHub
  Release were all skipped.
- `rtk gh run view 30699898430 --json ...` and `--log-failed` — confirmed
  `event=workflow_dispatch`, head branch
  `ma/npm-trusted-publishing-and-release-preflight`, exact head SHA `04d5bbb`,
  completed status, the six live npm reads, and the three skipped mutating
  stages. See `live-dispatch.md`.

GitHub emitted one external runner warning that `actions/checkout@v4`,
`actions/setup-go@v5`, and `actions/setup-node@v4` target the deprecated Node
20 action runtime and are being forced onto Node 24. The workflow's own npm
runtime guard passed, and the warning did not alter any Spec journey.

## Environment boundaries

- Existing tag-run history has no OIDC release on this implementation. The
  newest successful tag run is `v0.0.2`, run `30477190935`, at commit
  `755e9fba`, which predates the Spec's tooling commits. QA did not create a
  new tag, publish packages, or create a GitHub Release. A fresh complete OIDC
  release with an empty fallback record remains unproven.
- The prompt states that no Pull Request is open. The Run Worktree branch was
  not queried as though it had a Pull Request, so approval, review threads,
  checks, and Merge-Ready evidence remain unavailable.
