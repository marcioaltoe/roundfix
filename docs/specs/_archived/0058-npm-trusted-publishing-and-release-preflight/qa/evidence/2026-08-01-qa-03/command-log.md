# Command log — 2026-08-01 QA rerun 03

Build: `171f6a378c9e640a8a10c9382e28b501b21ff5a0`

## Repository Verification

- `rtk make verify` — exit 0. `rtk go test ./...` reported 2,941 passing
  tests in 24 packages; the four focused Skill/Baseline tests passed;
  `roundfix skills check` passed; the `roundfix` CLI build completed.

## Pending command groups

- Local workflow and documentation probes.

## Project Constraint and history audit

- `rtk ruby .../current_scope_probe.rb` — exit 0. All Task/constraint,
  authorization, exact `diff-tree`, ancestry, review-artifact settlement, and
  current-delta assertions passed.

## Live surfaces

- `rtk gh workflow run release.yml --ref
  ma/npm-trusted-publishing-and-release-preflight -f version=v0.0.2` — exit 0;
  created publish-free run `30703974453` on the exact head.
- `rtk gh run watch 30703974453 --exit-status` — exit 1 as expected from the
  used-version preflight barrier. Runtime, validation, and Verification passed;
  preflight named all six used coordinates; build, publish, and Release were
  skipped.
- Six `rtk npm view <coordinate>@0.0.2 version` reads — exit 0, each returned
  `0.0.2`.
- `rtk gh pr view 61 --json ...` — exit 0; open exact head, approved,
  MERGEABLE/CLEAN, both checks successful.
- `rtk gh api graphql ...reviewThreads...` — exit 0; two resolved threads,
  zero unresolved.
- `rtk gh run list --workflow release.yml --event push` and `rtk gh release
  list` — exit 0; newest real release predates this implementation.
- `rtk npx --yes ctx7@latest docs /websites/npmjs ...` — exit 0; current npm
  documentation confirms the official Trusted Publishing example uses
  `setup-node` with `registry-url`, a token-free publish command, and OIDC
  detection before traditional-token fallback.
- Final report scalar/row audit — exit 0: `closed`, `fail`, 29 rows, 27 pass,
  1 fail, 1 environment block, 0 finding blocks.
