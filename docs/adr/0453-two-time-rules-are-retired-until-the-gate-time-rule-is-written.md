---
status: accepted
knowledge: false
created_at: 2026-11-26T04:16:40Z
updated_at: 2026-11-26T04:16:40Z
deprecated_at: null
superseded_by: null
---

# Two time rules are retired until the Gate Time Rule is written

Supersedes ADR-0387 (2026-11-10); it is retired, and Git keeps its text.

Amends ADR-0452 (2026-11-24), ADR-0253 (2026-10-07), ADR-0252 (2026-10-07), ADR-0259 (2026-10-08), ADR-0257 (2026-10-08) and ADR-0182 (2026-09-30).

The maintainer decided on 2026-11-25 to retire the suite-time budget and the grouping, runner-type and 420-second provisions of ADR-0452, then write the Gate Time Rule from measured results. During the interval there is no suite-time budget above 420 seconds, no longer `-timeout` on a test invocation, and no paid, larger or self-hosted runner. Each gate reports its elapsed time; a single execution never fails for its time.

## The interval

The approved Gate Time Targets live in `docs/references/gate-time-targets.json`. On the mini, a Task settlement targets median 120 seconds and p95 240 seconds in the first 30 passing settlements; its lab targets are 45 seconds, or 150 seconds when regeneration contracts are selected. The Complete Suite Run targets 360 seconds and the check before a pull request opens targets 90 seconds in the lab. On the standard macOS runner, the first 20 passing pull request runs target median 480 seconds and p95 600 seconds; each test binary under Go's default limit targets 300 seconds. Runner waiting time is reported separately.

A Task settlement runs format, vet, build, selected Repository Contract Tests and the coverage record comparison, with no package tests. The Complete Suite Run tests every package with fresh results once per Spec, at the Delivery Queue repository gate after the archive commit, or on the archived branch for a Run delivered outside the queue. Before a pull request opens, `make verify-docs` remains mandatory.

Continuous integration runs only for a pull request into main, as the maintainer decided on 2026-11-25 (ai-memory `01a213bf-b859-7f60-91a1-201ca006c98a`). Five parallel parts on `macos-26` run the packages the change reaches, two CLI shards, and every Repository Contract Test once in the two contract parts. Each part merges the current base; the gathering job retains the required check name `Verification gate` and passes only when every part succeeds. Nothing verifies a push to main after the merge.

A release runs no test. It first proves the Release Provenance of the commit that last changed the launcher version: the pushed tag points at it, it is the squash merge of one pull request into main, its parent is an ancestor of the PR head with the same tree, and the latest GitHub Actions `Verification gate` of that head completed successfully. A manual dispatch evaluates that release commit and publishes nothing. An unreadable GitHub answer refuses publication. This is the Praetor's application of the maintainer's decision (answer `01a213da-6946-7920-bae3-fedc57a48923`), confirmed by the maintainer in `01a2140b-ac0f-7261-8134-ffeea8fa78a8`. The publication job stays as it is.

Before each release, the operator reads the windows with `make gate-times`. A gate over its number holds the release until its time returns or the maintainer decides on a Measured Dossier. The dossier is a Rule Challenge only in the Idea or PRD of a large new feature, with the gate's time before and after on the named machine, work already cut, merged or moved, and the measured new number; it is never requested during a Run, a corrective Spec or the day of an overrun.

## Consequences

Two pull requests can each pass and meet on main in a tree neither verified. Nothing verifies a direct push to main either. The release proof refuses a release merged behind main; update the release branch before its merge and rehearse with manual dispatch before pushing the tag. A refused release commit is never published and is cut again as the next patch version.

The Gate Time Rule comes later, after both measurement windows close, as a new accepted ADR with its own number. It amends this record, which stays accepted as the retirement record; it reuses none of ADR-0387, ADR-0452 or ADR-0453. Its numbers are measured results plus the declared 25 percent allowance, rounded up to ten seconds, never above the approved targets. This interval does not pre-approve those measured numbers.
