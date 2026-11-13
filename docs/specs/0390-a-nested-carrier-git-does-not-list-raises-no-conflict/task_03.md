---
task: task_03
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Baseline
inventory keep a nested `AGENTS.md` or `CLAUDE.md` only when Git lists it as
a file of the repository, so an ignored copy and a carrier inside a linked
worktree, a nested repository or a submodule raise no
`baseline.inventory.nested-carrier-conflict`; documents it in the `baseline`
guide, the Roundfix skill's `baseline` reference and the ADR-0070 amendment
notice; and records ADR-0399.

Every row runs on this repository's own files, temporary repositories, a
temporary local clone of the audited head, temporary homes and temporary
copies of the audited head. No row starts an Agent session or an ACP adapter,
calls ai-memory, reads or writes the real `~/.roundfix` or the real home
directory, runs the real `gh`, reaches the network, reads another
repository, sends a request to a provider, the judge or a reviewer, or reads
a credential. Every `go test` and `go build` it runs carries an explicit
`-timeout` where the command takes one and `GOPROXY=off`. Every row that
runs `bin/roundfix` or the delivery-base binary sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and
every Judge Key Variable for that command, puts a fake `gh` script that
prints `[]` first on `PATH`, and removes the temporary directories
afterwards. The QA never writes into the checkout while a test suite runs.
The gate keeps its evidence under names the repository does not ignore,
`.txt` for command output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes and repositories, and some, such as
`TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that
needs the session's toolchain configuration under `HOME`. Only rows that run
a built binary use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 and Invariants 1 to 3 by running
   `go test -count=1 -timeout 40m -v` over `./internal/baseline` with
   `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists`,
   `TestNestedCarrierInventoryAsksGitOnceWithNarrowArguments`,
   `TestBoundedInventoryIncludesAllCarriersAndIgnoresUnboundedPaths`,
   `TestRepositoryInspectionUsesNarrowReadOnlyGitCommands` and
   `TestRepositorySnapshotDigestChangesWithBoundedBytesOnly`, recording every
   test by name; then the whole `./internal/baseline` and `./skills`
   packages, the `Baseline` tests of `./internal/cli` (`-run 'Baseline'`)
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST verify Success Metric 3 through the built binaries, recording each
   command, its exit code and the extracted values:
   - build the delivery-base binary: `git archive` the merge base of the
     audited head with `main` into a temporary directory and run `go build
     -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix` there with
     `GOPROXY=off`;
   - `git clone -q` the audited head into a temporary directory, append
     `.storage/` to the clone's `.git/info/exclude`, and add a linked
     worktree with `git -C <clone> worktree add -q --detach .storage/sim/wt
     HEAD`;
   - run `baseline update --repo <clone> --no-skills --no-history --format
     json` with `bin/roundfix` and with the delivery-base binary, then remove
     the worktree with `git -C <clone> worktree remove .storage/sim/wt` and
     run both again;
   - from each JSON output extract `.planDigest`, `.state` and the sorted
     unique `path` of every object whose `code` is
     `baseline.inventory.nested-carrier-conflict`;
   - record that `bin/roundfix` reports the same Plan Digest and the same
     paths with and without the worktree, that its paths without the
     worktree equal the delivery-base binary's paths without the worktree,
     and that the delivery-base binary with the worktree reports at least one
     path under `.storage/sim/wt/`, the reproduction of the adopted message.
   The Spec adds no command, flag, exit code or output line, so no Surface
   Transcript applies.
4. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive` and committed
   in a fresh Git repository (the suite guard refuses a tree outside Git),
   with `-timeout` on every `go test`, and record that the named test fails
   on the named case: sabotages 1 to 6 as listed there. A whole-package
   failure of `TestOutputsForCommand` in such a copy, which reads this
   repository's history, is not the observation and is not counted.
5. MUST execute the tests tasks 01 and 02 name in their Verification.
6. MUST search and record:
   - the `baseline` guide and both copies of the Roundfix `baseline`
     reference carry Text A of API Contract 5, and the two copies hold the
     same bytes;
   - ADR-0070 carries Text B of API Contract 6 and is `accepted`, and
     ADR-0399 is `accepted`;
   - `internal/baseline` outside tests holds exactly one
     `"--exclude-standard"` argument list for `ls-files`, in
     `listRepositoryCarriers`.
7. MUST record, as outside-evidence rows, sources this Spec did not write,
   each with its date (read 2026-11-13) and none blocked, through the
   passages `_prd.md` → Acceptance evidence quotes:
   - the adopted quaestor message, a measurement on a repository this Spec
     did not build: compare its two warnings with the worktree and none
     after `git worktree remove` with the delivery-base binary's paths of
     requirement 3;
   - the git-ls-files and gitignore documentation: compare the quoted
     `--exclude-standard` sources and "Files already tracked by Git are not
     affected" with `forced/AGENTS.md` and `.storage/copy/CLAUDE.md` in
     `TestNestedCarrierInventoryKeepsOnlyCarriersGitLists`;
   - the installed `git`, a measurement this Spec did not design beyond its
     commands: in a temporary repository with a linked worktree and a
     submodule, `git ls-files -z --cached --others --exclude-standard` lists
     no file of either, and `git check-ignore` of a path inside the
     submodule exits `128`; record `git --version` and both results.
8. MUST record that `_prd.md` → Unreachable Acceptance declares none, so
   `rows_blocked_declared` is `0`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `None.` under `## Glossary`, that no Task added
   or changed a term in `CONTEXT.md`, and that
   `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `CONTEXT.md`,
      `docs/agents/setup-context.json`, any file under `docs/agents/` or
      `internal/baseline/assets/`, any file under `internal/baseline/testdata/`,
      `docs/references/coverage-record.json`, any Vendored Skill under
      `.agents/skills`, any skill other than the Roundfix skill, and any
      Roundfix reference other than `baseline`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0390-a-nested-carrier-git-does-not-list-raises-no-conflict --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent session or an ACP adapter;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the whole packages and the clone row.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An ignored linked worktree changes neither the warnings nor the Plan Digest of `baseline update`, and a tracked nested carrier still warns.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0390-a-nested-carrier-git-does-not-list-raises-no-conflict/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-6; Invariants 1-3; Declared breaks; Testing Approach; Build Order 3
