---
type: perf
status: promoted
created: 2026-09-30
spec: 0214-measure-before-changing
reason: null
---

# Archived Specs keep QA evidence and binaries nobody reads

## Slow

`docs/history` holds 77% of the repository's tracked files. Every repository-wide search and every Secondbrain mirror sync scans it, and the authoring skills' overlap checks read it.

## Measured

On 2026-09-30: `docs/history` was 39 MB and 4,167 of 5,405 tracked files. QA evidence directories were 9.1 MB in 1,571 files, and three binaries (one PDF and two PNGs) were 7 MB. Measured with `du` and `git ls-files` on main.

## Target

About 16 MB less, by dropping QA evidence directories and binaries at or after archive, while PRDs, TechSpecs, Task files, QA reports and authorization records stay. The archive command, the pre-PR review, reconcile, the citation and authorization checks and the tests that read archived Specs keep passing; the tests that read evidence need checking first.
