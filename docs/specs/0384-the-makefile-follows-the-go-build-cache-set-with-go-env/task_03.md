---
task: task_03
spec: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
status: pending
type: qa
complexity: low
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the
`Makefile` use the Go build cache set with `go env -w GOCACHE` when it is an
existing writable directory, under the 10,240 MiB default bound; when that
cache is set but unusable, `make` never creates it, warns once and falls back
to its temporary cache under a 2,048 MiB default bound; with nothing set, or
with `GOCACHE` exported, nothing changes. The repository guide, README,
glossary and ADR-0370 describe it, and ADR-0394 records it.

Every row runs in this repository or in temporary directories it creates
and removes. Only row 1's `make build` and repository Verification run in the
session's own environment, as every Run's Verification does; every other
`make` a row runs gets `GOENV` set to a file in a `mktemp -d` directory and
`TMPDIR` set to a `mktemp -d` directory, and either a stub `go` or
`GOCACHE_BOUND_MB=0`, so it neither reads the host's Go environment file nor
trims a real cache. No row creates, removes or changes any file in `/tmp`,
`/private/tmp`, `/Volumes` or `$HOME` outside a directory made by
`mktemp -d`, reaches the network, runs `gh` or `git push`, reads another
repository, starts an Agent session or an ACP adapter, reads or writes the
real `~/.roundfix`, sends a request to a provider, the judge or a reviewer,
records an environment value, or sends a signal to any process. Every
`go test` carries an explicit `-timeout` and keeps the session's own `HOME`
and `TMPDIR`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, and removes them afterwards. The QA never writes into the
checkout while a test suite runs. The gate keeps its evidence under names the
repository does not ignore, `.txt` for command output, and links only files
it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v -run '^TestMakeFollowsTheGoBuildCacheSetWithGoEnv$' ./internal/verifyselect`,
   the new test and record the pass line of each of its seven subtests
   (`set_and_writable`, `set_and_not_mounted`, `set_and_read-only`,
   `set_to_a_file`, `set_with_a_bound`, `environment_wins`, `nothing_set`);
   then the whole `./internal/verifyselect` package,
   `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract`
   and `GOOS=windows go vet ./internal/verifyselect`.
3. MUST record Success Metrics 1 to 3 against the starting behavior: in a
   `git archive` of the Delivery Base (the merge base of the audited head
   with `main`) unpacked in a temporary directory and made a Git repository
   there, copy the audited head's
   `internal/verifyselect/build_cache_go_env_makefile_test.go` over and run
   it with `-timeout`; record that `set_and_writable`, `set_and_not_mounted`,
   `set_and_read-only`, `set_to_a_file` and `set_with_a_bound` fail there,
   that `nothing_set` and `environment_wins` pass there and on the audited
   head (the kept paths), and that all seven pass on the audited head
   (row 2).
4. MUST record Success Metric 4 and the Declared breaks: on the audited head,
   the five changed test files of `internal/verifyselect` differ from the
   Delivery Base only by the `env` early return of their `go` stubs (one
   line each, two in `contracts_repository_test.go`), and
   `internal/verifyselect/build_cache_makefile_test.go` is byte-identical to
   the Delivery Base; the whole package passes (row 2).
5. MUST record the warning line by running the audited `Makefile` once
   outside the test: in a `mktemp -d` directory `d`, write `d/go.env` holding
   `GOCACHE=<d>/Volumes/missing/cache/go-build`, create `d/Volumes` and
   `d/tmp`, and run
   `env -u GOCACHE GOENV="$d/go.env" TMPDIR="$d/tmp/" GOCACHE_BOUND_MB=0 make --no-print-directory -f <checkout>/Makefile help`
   from `d`; record that standard error is exactly one line
   `gocache: warning: GOCACHE set with go env (<d>/Volumes/missing/cache/go-build) is not a writable directory; using <d>/tmp/roundfix-go-build-<id -u> with a 2048 MiB default bound`
   and that `<d>/Volumes/missing` does not exist; then remove `d`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and made a Git repository, change only the `Makefile`,
   run `TestMakeFollowsTheGoBuildCacheSetWithGoEnv` with `-timeout`, and
   record that the named subtest fails:
   - `[ -d "$$cache" ]` replaced by `mkdir -p "$$cache" 2>/dev/null` in
     `GOCACHE_FROM_GO_ENV`: `set_and_not_mounted`;
   - `printf 2048)` replaced by `printf 10240)`: `set_and_not_mounted`,
     `set_and_read-only` and `set_to_a_file`;
   - the warning `printf` turned into `:`: the same three;
   - `GOCACHE_BOUND_MB ?= 10240` after `GOCACHE := $(GOCACHE_FROM_GO_ENV)`
     replaced by `GOCACHE_BOUND_MB ?= 0`: `set_and_writable`;
   - the probe removed, so the check is only `[ -d "$$cache" ]`:
     `set_and_read-only`;
   - `ifeq ($(origin GOCACHE),undefined)` replaced by
     `ifneq ($(origin GOCACHE),command line)`: `environment_wins`.
7. MUST execute the tests the non-QA Tasks name in their Verification.
8. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted message
   (`references/`, measured by gesttione-drive on the mini on 2026-11-12:
   "du $TMPDIR/roundfix-go-build-501 = 10G" and
   "go env GOCACHE = /Volumes/Brain-SSD/cache/go-build"), the go command
   documentation quoted in `_prd.md` → Acceptance evidence ("The -changed
   flag prints only those settings whose effective value differs from the
   default value that would be obtained in an empty environment with no prior
   uses of the -w flag." and "Must be an absolute path."), and the 2026-10-12
   Secondbrain digest quoted there, compared with the `Makefile` block and
   rows 2 and 5. The row reads the files in this repository and the PRD,
   never the network or the Secondbrain.
9. MUST record the PRD's one Unreachable Acceptance entry, the operator's run
   on the Mac mini after the merge, as exactly one row
   `blocked (environment: the Mac mini's Go environment file, its external disk and the operator's deletion after the merge required)`
   with `declared: ` followed by the entry's criterion text, counted once in
   `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `- changes: **Build Cache Bound**` and that the
    Task commit range changes the definition line of that entry in
    `CONTEXT.md` and no other entry.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is listed in `_authorization.md` `paths`;
      `docs/agents/setup-context.json` changes only its two `Makefile`
      `declarationDigest` values and, when the Delivery Base has none, gains
      `roundfixVersion` equal to the `Version` of `internal/app/version.go`
      at the audited head; `docs/agents/specific-repository.md` changes only
      the one rule line;
    - the `Makefile` changes only inside the
      `ifeq ($(origin GOCACHE),undefined)` block;
    - ADR-0370 gains only the `updated_at` change, the amendment paragraph and
      the two-line option note;
    - the Task commit range changes none of these: `.roundfixrc.yml`, any
      workflow, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/**`, any Go production file, any
      skill, any Baseline module or other derived Baseline file, any
      user-guide page, or the `### QA settlement` section of any skill.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0384-the-makefile-follows-the-go-build-cache-set-with-go-env --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp`, `/Volumes` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - run `make` outside row 1 with the host's Go environment file, or
      clean or delete any Go build cache directory;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base with the new test, record the warning line, and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the declared unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A usable cache set with `go env` is used and bounded at 10240 on the audited head and ignored on the Delivery Base; an unusable one is never created, and `make` warns once and falls back at 2048.
- [ ] Nothing set and an exported `GOCACHE` behave as on the Delivery Base; each sabotage makes its named subtest fail.
- [ ] The one Unreachable Acceptance entry is one declared row, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0384-the-makefile-follows-the-go-build-cache-set-with-go-env/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Invariants; Declared breaks; Testing Approach; Build Order 3
