---
task: task_04
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec counts
`docs/references/` as planning evidence, recognises a
**Release-Preparation Commit** from the `release_paths` of the Specs a
range delivers, reports it in the release plan's JSON change evidence, and
describes both rules and both publication triggers in the release guides,
the Roundfix skill and the glossary. A Release Ceiling still never answers
a classification by hand.

Every row runs on this repository's own files, temporary repositories,
temporary homes and clean temporary worktrees or shared clones of this
repository, read-only toward its refs. No row starts an Agent harness, an
Agent session or an ACP adapter, starts a Run, a Daemon or a Delivery
Owner, calls ai-memory, reads or writes the real `~/.roundfix` or the real
home directory, runs the real `gh`, reaches the network, sends a request to
a provider, the judge or a reviewer, reads another repository, or asserts a
wall-clock duration. Every `go test` it runs carries an explicit
`-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary
`make build` made from the audited head, and never a `roundfix` found on
`PATH`; requirement 6 alone also runs a binary the gate builds from the
Delivery Base. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`,
`CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command,
reads standard input from `/dev/null`, and removes the temporary
directories afterwards. The release plan refuses a dirty worktree, so every
row that plans this repository's history runs in a clean shared clone of
it made under a temporary directory and removed afterwards; the gate never
creates, moves or deletes a tag, a ref or a remote of this repository. The
QA never writes into the checkout while a test suite runs. The gate keeps
its evidence under names the repository does not ignore, `.txt` for
command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running, with `-v`, the
   `internal/releaseplan` tests task_01 names in its Verification,
   recording every test and subtest by name, and the whole
   `./internal/releaseplan` package.
3. MUST verify Success Metric 2 and Success Metric 3 by running, with `-v`,
   the two `internal/cli` tests task_03 names in its Verification,
   recording every test and subtest by name, and the whole `./internal/cli`
   package. Then, through `bin/roundfix`, on a temporary Git repository the
   gate builds with a planning commit adding `docs/references/evidence.md`,
   a `feat:` delivery whose approved record grants `release` with
   `release_ceiling: minor` and `release_paths: [VERSION, CHANGELOG.md]`,
   and a `chore: release v0.2.0` commit changing only those two files,
   record in text and JSON: decision `ready`, `Release authorization:
   answered`, no blocking commit, and `releasePreparation: true` on the
   preparation commit only; and that the same plan with
   `--impact minor --reason x` reports
   `Release authorization: not answered — the plan was classified by hand`.
   Record a listing of every file of that repository with its size and
   SHA-256 before and after, unchanged.
4. MUST reproduce Surface Transcript 1 of `_techspec.md` through
   `bin/roundfix release plan --from v0.143.3 --to v0.143.4` in a clean
   shared clone of this repository checked out at `v0.143.4`, recording
   standard output, standard error and the exit code, compared line by line
   with the transcript, where `...` stands for the check lines.
5. MUST verify Success Metric 4 by searching and recording that both
   copies of the Roundfix release reference, both copies of its `SKILL.md`,
   the release plan usage text, `docs/user-guide/release-runbook.md` and
   `docs/user-guide/usage.md` carry the text task_02 prescribes; that
   neither copy of the release reference contains "Preserve the existing
   tag-triggered workflow"; that each copy under `skills/` equals its
   `.agents/skills/` source; and that the section `## Completion notices`
   of the release reference is byte-identical to the Delivery Base's (the
   merge base of the audited head with the default branch). Then run, with
   `-v`, the test task_02 names, the whole `./skills` package and
   `-tags docscontract ./internal/docscontract`, recording each pass.
6. MUST compare the audited head with the Delivery Base through the
   product. Build the Delivery Base's binary in a temporary copy made with
   `git archive` and committed in a fresh Git repository, outside the
   checkout. On a temporary repository with a `fix:` commit and a `docs:`
   commit changing only `docs/specs/` and `docs/backlog/` files, record
   that both binaries print the same text and the same JSON, byte for
   byte, and exit with the same code; and on the repository of
   requirement 3 without its reference file and preparation commit, that
   both reach the same decision, proposed version and authorization answer.
7. MUST record the outside-evidence row of `_prd.md` → Success Metric 5
   and Acceptance evidence. In one clean shared clone of this repository,
   for each of the 51 ranges of
   [`references/2026-11-16-release-plan-classification-measurement.md`](references/2026-11-16-release-plan-classification-measurement.md),
   run `bin/roundfix release plan --from <previous tag> --to <tag>` and
   `--to '<tag>^'` (except `v0.130.0..v0.130.1`, whose tag holds only the
   preparation) and record each decision. The row passes when: every range
   the measurement marks `yes` reaches the same decision at the tag as at
   the commit before it; `v0.143.3..v0.143.4` is `ready` with `v0.143.4` at
   both; `v0.142.0..v0.142.1` still blocks on `.ai-memory.toml` at both;
   and the five ranges marked `no` stay `manual_classification_required`
   at the tag. Record the clone's `git status --porcelain` empty afterwards.
8. MUST run each sabotage below on its own temporary copy of the audited
   head made with `git archive` and committed in a fresh Git repository
   (the suite guard refuses a tree outside Git), with `-timeout` on every
   `go test`, and record that the named test fails on the named subtest. No
   sabotage regenerates or synchronizes anything:
   1. `docs/references/` is removed from the planning-evidence paths:
      `TestMaintenanceOnly/docs/references/evidence.md`.
   2. The impact condition is dropped, so a `fix:` commit on release paths
      becomes a preparation:
      `TestReleasePreparationCommitSetsNoImpact/fix_on_release_paths`.
   3. A commit with no path counts as a preparation:
      `TestReleasePreparationCommitSetsNoImpact/no_changed_path`.
   4. One listed path is enough instead of every path:
      `TestReleasePreparationCommitSetsNoImpact/one_unlisted_path`.
   5. `Build` passes no release paths to the classifier:
      `TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease/preparation_is_answered`.
   6. `AnswerApproval` answers a plan classified by hand:
      `TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/classified_by_hand`.
   7. The JSON field is left unfilled:
      `TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/preparation_keeps_the_decision`.
   8. The sentence "Preserve the existing tag-triggered workflow" is
      restored in the canonical release reference:
      `TestTheReleaseReferenceKeepsTheRepositoryPublicationTrigger`.
9. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and the TechSpec declare `adds:` for
    **Release-Preparation Commit**, that task_02 wrote that entry and no
    other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`,
      `docs/references/skill-coverage.json`, `internal/cli/cli_test.go`,
      `internal/docscontract/publicdocs_test.go`, `.roundfixrc.yml`, any
      file under `docs/agents/` or `internal/baseline/assets/`, any
      Vendored Skill under `.agents/skills`, any ADR, and any owned skill
      other than `roundfix`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent
    of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0400-planning-references-and-release-preparation-need-no-manual-classification --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write
    the report's Outcome as one paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - create, move or delete a tag, a ref or a remote of this repository;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the three Tasks and their packages.
- [ ] Reproduce the product rows, the Surface Transcript and the comparison with the Delivery Base.
- [ ] Rerun the measured ranges of this repository's history in a clean clone.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A planning reference and a release-preparation commit need no classification by hand, and the Release Ceiling answers such a range as it answers the range without them.
- [ ] A Release Ceiling still never answers a plan classified by hand, and plans without either commit keep the Delivery Base's bytes.
- [ ] The release guides describe both rules and both publication triggers.
- [ ] This repository's measured ranges reach the decisions Success Metric 5 names, and each sabotage makes its named test fail.

## Verification

- `newest="$(find 'docs/specs/0400-planning-references-and-release-preparation-need-no-manual-classification/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-5; Surface Transcript 1; Interfaces; Testing Approach; Build Order 4
- [`references/2026-11-16-release-plan-classification-measurement.md`](references/2026-11-16-release-plan-classification-measurement.md)
- ADR-0412, ADR-0349, ADR-0143, ADR-0256
