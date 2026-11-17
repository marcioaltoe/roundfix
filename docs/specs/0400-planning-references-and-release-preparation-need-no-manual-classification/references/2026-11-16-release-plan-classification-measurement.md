# Release plan classification across this repository's releases

Measured by the author on 2026-11-16 in a clone of this repository at `c2d1f77c` (branch `docs/plan-0399`), with the installed Roundfix 0.147.0. Every command is read-only.

## Method

For each pair of consecutive stable tags from `v0.125.0` to `v0.147.0` (51 ranges), the author ran two plans:

1. `roundfix release plan --from <previous tag> --to '<tag>^'`, the range a release session plans before it writes the release-preparation commit. In this repository the tag points at that commit (`chore: release vX.Y.Z`), which changes `CHANGELOG.md`, `dist/npm/roundfix/package.json` and `internal/app/version.go`.
2. `roundfix release plan --from <previous tag> --to <tag>`, the same range with the release-preparation commit, as a plan rerun after that commit, or on the head of a release Pull Request, reads it.

A script then read, for each range, every commit that deletes `docs/specs/<slug>/_authorization.md` and adds `<slug>.md` in the same change (a delivery, as the plan finds one). It read the record at the commit's first parent and kept it when the record is `status: approved`, lists `release` and carries `release_ceiling` and `release_paths`. It then checked whether every path the release-preparation commit changes is one of the kept records' `release_paths`.

## Results

- Before the release-preparation commit, 48 of 51 ranges were conclusive or asked only the approval question, and one (`v0.130.0..v0.130.1^`) is not a range, because its tag holds only the preparation. Two ended in `manual_classification_required`:
  - `v0.142.0..v0.142.1^`: `73421db24d57 docs: plan 0383, 0384 and 0385 and add the ai-memory project marker`, blocked by `.ai-memory.toml`, a root file outside every planning path.
  - `v0.143.3..v0.143.4^`: `576ae51d9d4a docs: plan 0393 a feature is judged for the idea stage before its prd (#871)`, blocked only by `docs/references/idea-stage-judgment-measurement.md`. Every other path of that commit is under `docs/specs/` or `docs/adr/`. The Spec 0393 planning files of the same commit cite the reference (`_authorization.md`, `_prd.md`, `task_04.md`).
- With the release-preparation commit, all 51 ranges ended in `manual_classification_required`, each with the `chore: release` commit listed as `[impact=none, type=chore, manual]`.
- In 46 of the 51 ranges, every path the release-preparation commit changes is a `release_paths` entry of a Spec the range delivers with a granted ceiling. The five others (the four ranges from `v0.125.0` to `v0.129.0`, and `v0.130.0..v0.130.1`, whose only commit is the preparation) deliver no such Spec.

## Table

| range | before the release-preparation commit | at the tag | Specs delivered with a granted ceiling | preparation paths within their `release_paths` |
| --- | --- | --- | --- | --- |
| `v0.125.0..v0.126.0` | approval_required | manual_classification_required | - | no |
| `v0.126.0..v0.127.0` | approval_required | manual_classification_required | - | no |
| `v0.127.0..v0.128.0` | approval_required | manual_classification_required | - | no |
| `v0.128.0..v0.129.0` | approval_required | manual_classification_required | - | no |
| `v0.129.0..v0.129.1` | ready | manual_classification_required | 0341 | yes |
| `v0.129.1..v0.130.0` | ready | manual_classification_required | 0343,0344,0342 | yes |
| `v0.130.0..v0.130.1` | invalid range (the tag's only commit is the preparation) | manual_classification_required | - | no |
| `v0.130.1..v0.131.0` | ready | manual_classification_required | 0348,0347,0345,0346,0349,0350 | yes |
| `v0.131.0..v0.132.0` | ready | manual_classification_required | 0351 | yes |
| `v0.132.0..v0.132.1` | ready | manual_classification_required | 0352 | yes |
| `v0.132.1..v0.133.0` | ready | manual_classification_required | 0353 | yes |
| `v0.133.0..v0.133.1` | ready | manual_classification_required | 0356 | yes |
| `v0.133.1..v0.134.0` | ready | manual_classification_required | 0355 | yes |
| `v0.134.0..v0.134.1` | ready | manual_classification_required | 0354 | yes |
| `v0.134.1..v0.135.0` | ready | manual_classification_required | 0358 | yes |
| `v0.135.0..v0.136.0` | ready | manual_classification_required | 0357 | yes |
| `v0.136.0..v0.136.1` | ready | manual_classification_required | 0359 | yes |
| `v0.136.1..v0.136.2` | ready | manual_classification_required | 0360 | yes |
| `v0.136.2..v0.136.3` | ready | manual_classification_required | 0363 | yes |
| `v0.136.3..v0.136.4` | ready | manual_classification_required | 0362 | yes |
| `v0.136.4..v0.137.0` | ready | manual_classification_required | 0361 | yes |
| `v0.137.0..v0.137.1` | ready | manual_classification_required | 0364 | yes |
| `v0.137.1..v0.137.2` | ready | manual_classification_required | 0365 | yes |
| `v0.137.2..v0.137.3` | ready | manual_classification_required | 0366 | yes |
| `v0.137.3..v0.137.4` | ready | manual_classification_required | 0367 | yes |
| `v0.137.4..v0.137.5` | ready | manual_classification_required | 0368 | yes |
| `v0.137.5..v0.137.6` | ready | manual_classification_required | 0369 | yes |
| `v0.137.6..v0.137.7` | ready | manual_classification_required | 0370 | yes |
| `v0.137.7..v0.137.8` | ready | manual_classification_required | 0372 | yes |
| `v0.137.8..v0.137.9` | ready | manual_classification_required | 0371 | yes |
| `v0.137.9..v0.138.0` | ready | manual_classification_required | 0373 | yes |
| `v0.138.0..v0.139.0` | ready | manual_classification_required | 0374,0375 | yes |
| `v0.139.0..v0.139.1` | ready | manual_classification_required | 0377 | yes |
| `v0.139.1..v0.139.2` | ready | manual_classification_required | 0376 | yes |
| `v0.139.2..v0.139.3` | ready | manual_classification_required | 0378 | yes |
| `v0.139.3..v0.140.0` | ready | manual_classification_required | 0379 | yes |
| `v0.140.0..v0.141.0` | ready | manual_classification_required | 0380 | yes |
| `v0.141.0..v0.141.1` | ready | manual_classification_required | 0381 | yes |
| `v0.141.1..v0.142.0` | ready | manual_classification_required | 0382 | yes |
| `v0.142.0..v0.142.1` | manual_classification_required — 73421db24d57 | manual_classification_required | 0383 | yes |
| `v0.142.1..v0.142.2` | ready | manual_classification_required | 0385,0384 | yes |
| `v0.142.2..v0.143.0` | approval_required | manual_classification_required | 0387,0386 | yes |
| `v0.143.0..v0.143.1` | ready | manual_classification_required | 0388 | yes |
| `v0.143.1..v0.143.2` | ready | manual_classification_required | 0389 | yes |
| `v0.143.2..v0.143.3` | ready | manual_classification_required | 0390 | yes |
| `v0.143.3..v0.143.4` | manual_classification_required — 576ae51d9d4a | manual_classification_required | 0391 | yes |
| `v0.143.4..v0.144.0` | ready | manual_classification_required | 0392 | yes |
| `v0.144.0..v0.145.0` | ready | manual_classification_required | 0393 | yes |
| `v0.145.0..v0.145.1` | ready | manual_classification_required | 0396 | yes |
| `v0.145.1..v0.146.0` | ready | manual_classification_required | 0397 | yes |
| `v0.146.0..v0.147.0` | ready | manual_classification_required | 0398 | yes |
