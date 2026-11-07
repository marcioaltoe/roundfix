---
task: task_04
spec: 0364-a-rerun-finishes-a-partial-npm-publication
status: completed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec gives the release
workflow's Publication Preflight a `published` state, adds the tag-run step
"Compare published packages", which accepts an already published coordinate
only when its decompressed tar stream equals this run's build, makes the
publish step skip such a coordinate, and documents the recovery in the
release runbook, the glossary and ADR-0146's amendment notice.

Every row runs the workflow's own step scripts under `bash` with stub `curl`
and `npm` commands in temporary directories. No row reaches the network,
reads the npm registry, publishes a package, pushes a tag, creates a GitHub
Release, starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix`, sends a request to a provider, the judge or a reviewer, reads
a credential, or sends a signal to any process. Every `go test` carries an
explicit `-timeout` and keeps the session's own `HOME` and `TMPDIR`. Every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, with `GOMODCACHE` and `GOCACHE` kept at their `go env` values,
and removes them afterwards. Every stub directory on a `PATH` is an absolute
path. The QA never writes into the checkout while a test suite runs. The gate
keeps its evidence under names the repository does not ignore, `.txt` for
command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v`, the tests tasks 02
   and 03 name over `./internal/verifyselect`, then the whole
   `./internal/verifyselect` package and
   `-tags docscontract ./internal/docscontract`, recording each pass and
   that none of the five `TestRelease…` tests reports `SKIP`.
3. MUST record, from the `-v` output of
   `TestReleaseRerunPublishesOnlyTheMissingCoordinates`, Success Metrics 1
   and 2: the v0.132.1 shape publishes only `roundfix`, and the v0.130.0
   shape publishes the four other platform coordinates and then `roundfix`.
   It MUST also run `TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne`
   on a `git archive` copy of the audited head whose
   `.github/workflows/release.yml` is replaced by the Delivery Base's
   (`git show <delivery-base>:.github/workflows/release.yml`), and record
   that it fails on the v0.132.1 shape with five `is already used` lines,
   the refusal runs 41750342366 and 42056741468 met.
4. MUST verify Success Metric 5 from
   `TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry`: one
   `roundfix` publish, the `publish:` line, and no publish with
   `NODE_AUTH_TOKEN` set.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - the classifier puts `used` back for a version under `versions`
     (`elif (.versions | has($tag)) or (.time | has($tag)) then "used"`):
     `TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne`
     fails on the v0.132.1 shape;
   - the classifier calls a version under `time` only `published`:
     the same test fails on the `time`-only case;
   - the all-published refusal removed: the same test fails on the
     all-published case;
   - the compare step compares `$local_integrity` with
     `$registry_integrity` instead of the tar streams:
     `TestReleaseCompareAcceptsOnlyTheSameContent` fails on case (a), whose
     two gzip levels give different integrities for the same tar stream;
   - the compare step appends every coordinate to the complete file:
     the same test fails on case (b);
   - the publish skip removed from `publish_coordinate`:
     `TestReleaseRerunPublishesOnlyTheMissingCoordinates` fails on the
     v0.132.1 shape, publishing six coordinates.
6. MUST execute the tests every non-QA Task names in its Verification.
7. MUST record, as an outside-evidence row, a source this Spec did not
   write, with its date and not blocked: the npm `publish` documentation
   passage `_prd.md` → Acceptance evidence quotes ("Once a package is
   published with a given name and version, that specific name and version
   combination can never be used again, even if it is removed with `npm
   unpublish`."), compared offline with the v0.132.1 shape of
   `TestReleaseRerunPublishesOnlyTheMissingCoordinates`, which never
   publishes a coordinate already published, and with the `time`-only case
   of `TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne`,
   which still refuses an unpublished version.
8. MUST record the PRD's one Unreachable Acceptance entry (a real rerun of a
   failed release job on GitHub Actions against the npm registry) as exactly
   one row, `blocked (environment: publishing to npm is forbidden to this
   Spec)`, counted once in `rows_blocked_declared`, supported by rows 3
   and 5.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `- changes:` for **Release Set** and
   **Publication Preflight**, that task_01's commit changes exactly those
   two entries of `CONTEXT.md`, and that
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs)$'`
   passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - ADR-0146 carries the `Amended by [ADR-0374]` paragraph and its
      `updated_at` is `2026-11-06T23:30:00Z`;
    - the Task commit range changes none of these: the Makefile, any
      workflow other than `.github/workflows/release.yml`, `go.mod`,
      `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/**`, any file under
      `internal/baseline/` or `internal/suiteguardcontract/`, any guide under
      `docs/agents/`, any skill, or any ADR other than ADR-0146;
    - in `.github/workflows/release.yml` only the steps "Publication
      preflight", "Compare published packages" (new) and "Publish to npm"
      changed, and in "Publish to npm" only the lines task_03 names.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0364-a-rerun-finishes-a-partial-npm-publication --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run `npm publish`, `npm pack <coordinate>@<version>` against the
      registry, `gh release` or `git push`;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the Delivery Base characterization.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside source, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A rerun of a partial publication publishes only the missing coordinates, and the Delivery Base refuses it.
- [ ] A published coordinate with different content stops the job before any publish, and a gzip-only difference compares equal.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-3; Testing Approach; Build Order 4
