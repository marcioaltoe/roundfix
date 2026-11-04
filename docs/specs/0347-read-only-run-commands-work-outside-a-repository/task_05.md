---
task: task_05
spec: 0347-read-only-run-commands-work-outside-a-repository
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the five Read-Only Run
Commands (bare `runs`, `runs list`, `runs show`, `attach`, `events`)
load without failing on a `worktree.location` refusal and print it as one
warning, makes `runs list` outside a Git repository list every repository,
and makes the refusal name the detected Git root and its `.git` path, with
the Roundfix skill, the guides and `CONTEXT.md` stating the contract.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, reads or writes the real `~/.roundfix`, reads
another repository, reaches the network, sends a request to a provider, the
judge or a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE`
kept at their `go env` values, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never writes
into the checkout while a test suite runs. The gate keeps its evidence under
names the repository does not ignore, `.txt` for command output, and links
only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` and
   `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal` over
   `./internal/config`, and
   `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository`,
   `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory`,
   `TestStrayGitDirectoryTranscripts`,
   `TestRunRunsListOutsideRepositoryListsEveryRepository` and
   `TestRunsListOutsideARepositoryTranscript` over `./internal/cli`, then the
   whole `./internal/config`, `./internal/cli`, `./skills` packages,
   `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce Surface Transcripts 1, 2 and 3 of `_techspec.md` through
   the built binary, each with its own fresh temporary `HOME` (`<home>`) and
   `TMPDIR`: Transcript 1 from a fresh directory outside any repository;
   Transcripts 2 and 3 from `<home>/drive` with an empty `<home>/.git`
   directory. Compare stdout, stderr and exit code with the transcript under
   its matching conventions and record each result. Then repeat Transcript 2
   with an empty `<home>/.git` file instead of a directory and record the same
   result. Remove every temporary directory afterwards.
4. MUST reproduce the maintainer's report through the built binary: with a
   fresh temporary `HOME` holding an empty `.git` directory, from
   `<home>/drive`, run `bin/roundfix runs list`, `bin/roundfix runs list
   --all` and `bin/roundfix attach run_x` and record that the first two exit
   `0` and the third fails only on the absent Run Database, each with the
   warning as its first stderr line; then run `bin/roundfix runs list` from a
   fresh directory outside any repository with that same `HOME` minus its
   `.git` and record exit `0`.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - the refusal reverted to `worktree.location must not be inside the
     repository tree: %q` alone:
     `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` fails on `git
     directory`;
   - `Load` failing even with `WorktreeLocationOptional`:
     `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal` fails;
   - `findGitRoot` accepting only a `.git` directory:
     `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` fails on `git
     file`;
   - `runAttachCommand` back on `loadCommandConfig`:
     `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory` fails on `directory`;
   - the warning line not written:
     `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory` fails;
   - `loadCommandConfig` setting `WorktreeLocationOptional` for every
     command: `TestStrayGitDirectoryTranscripts` fails on `reconcile under a
     stray .git`;
   - `runs list` without a Git root left repository-scoped:
     `TestRunRunsListOutsideRepositoryListsEveryRepository` fails;
   - the old refusal restored for a missing Git root:
     `TestRunsListOutsideARepositoryTranscript` fails.
6. MUST execute the tests tasks 01 to 04 name in their Verification.
7. MUST search both copies of the Roundfix `runs`, `events` and
   `review-runs` references, `docs/user-guide/commands/runs.md`,
   `attach.md`, `events.md`, `docs/user-guide/configuration.md`, the `runs`
   help and `CONTEXT.md`, and record that they state the contract and that no
   copy of the `runs` reference still says that `runs list` without `--all`
   exits `2` outside a repository.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the maintainer report quoted in
     `references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md`
     (2026-11-03): compare offline its commands, directory and refusal with
     the results of requirement 4;
   - the git(1) passage `_prd.md` → Acceptance evidence quotes ("Git tries to
     find such a directory in the parent directories to find the top of the
     working tree"): compare it with the `git directory` and `git file` rows
     of `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit`, whose
     working directory lies below the root that holds the `.git`.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds: **Read-Only Run Command**`, that
    `CONTEXT.md` defines it in the `**Term**:` form, and that the Task commit
    range changes no other glossary entry.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`,
      anything under `internal/baseline/assets/` or `docs/agents/`, any skill
      other than the Roundfix skill files task_04 declares;
    - no `### QA settlement` section changed.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0347-read-only-run-commands-work-outside-a-repository --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcripts and the report through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The five commands work outside a repository and under a stray `.git`, with the warning, and keep their output inside a repository.
- [ ] The refusal names the Git root and its `.git`, and the commands that need the location still refuse.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0347-read-only-run-commands-work-outside-a-repository/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
