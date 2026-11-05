---
task: task_04
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Git root walk of
configuration resolution and of the skills project scope accept only a `.git`
entry that Git accepts (ADR-0362), moves the test fixtures onto
`gittest.MarkRepository`, and states the rule in the Roundfix skill, the
configuration and `runs` guides and ADR-0357's amendment notice.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, reads or writes the real `~/.roundfix` or the real
home directory, reads another repository, reaches the network, sends a
request to a provider, the judge or a reviewer, or reads a credential. Every
`go test` it runs carries an explicit `-timeout`. Every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, because the CLI
loads `$HOME/.roundfix/config.yml`, and removes them afterwards. The `go test`
package runs of requirement 2 and of the sabotages do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. The QA never
writes into the checkout while a test suite runs. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts`,
   `TestFindGitRootAgreesWithGit` and
   `TestLoadUnderARejectedGitDirectoryFindsNoGitRoot` over
   `./internal/config`, and `TestRejectedGitDirectoryTranscripts`,
   `TestSkillsProjectRootSkipsARejectedGitDirectory`,
   `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory` and
   `TestStrayGitDirectoryTranscripts` over `./internal/cli`, recording every
   row of the shape table by name; then the whole `./internal/config`,
   `./internal/gittest`, `./internal/judge`, `./internal/cli` and `./skills`
   packages, `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce Surface Transcripts 1, 2 and 3 of `_techspec.md` through
   the built binary, each with its own fresh temporary `HOME` (`<home>`) and
   `TMPDIR`, from `<home>/drive` with `<home>/.git/gk/` and
   `<home>/.git/info/` and nothing else in `<home>/.git`. Compare stdout,
   stderr and exit code with the transcript under its matching conventions,
   record that `<home>/.roundfixrc.yml` does not exist afterwards, and remove
   every temporary directory afterwards.
4. MUST reproduce the maintainer's report through the built binary and Git:
   with a fresh temporary `HOME` holding `<home>/.git/gk/` and
   `<home>/.git/info/`, from `<home>/drive`, record what `git rev-parse
   --show-toplevel` prints (`not a git repository`) and that `bin/roundfix
   runs list` prints no `worktree.location` warning; then, in
   `<home>/drive/repo` made with `git init`, record that `bin/roundfix init
   --scope project` exits `0` and writes `<home>/drive/repo/.roundfixrc.yml`,
   the root `git rev-parse --show-toplevel` prints there.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `gitEntryAccepted` returning true for any existing entry:
     `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts` fails on `only gk and
     info`, and `TestRejectedGitDirectoryTranscripts` fails on `Surface
     Transcript 1`;
   - `refs` no longer required:
     `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts` fails on `HEAD and
     objects without refs`;
   - `HEAD` content not checked: it fails on `HEAD that names nothing`;
   - `commondir` ignored: it fails on `common directory`;
   - any regular `.git` file accepted: it fails on `file without gitdir`;
   - a `gitdir:` file accepted only when its target exists: it fails on
     `gitdir file whose target is missing`;
   - `defaultResolveSkillsProjectRoot` back on a walk that accepts any `.git`:
     `TestSkillsProjectRootSkipsARejectedGitDirectory` fails;
   - `fakeWorktree.Snapshot` no longer writing its file:
     `TestRunOperationalCommandAcceptsMVPFlags` fails on `resolve`.
6. MUST execute the tests tasks 01 to 03 name in their Verification.
7. MUST search both copies of the Roundfix `review-runs` and `runs`
   references, `docs/user-guide/configuration.md`,
   `docs/user-guide/commands/runs.md` and ADR-0357, and record that the first
   six state the rule sentence of API
   Contract 5, that none still says that repository detection finds any
   `.git` entry in an ancestor, and that ADR-0357 carries the amendment
   notice naming ADR-0362 and its `updated_at` of 2026-11-05.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the maintainer report quoted in
     `references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md`
     (2026-11-05): compare offline its `~/.git` contents and Git's answer
     with the `only gk and info` row and with the results of requirement 4;
   - Git's own answer, through the `git` binary: record the rows of
     `TestFindGitRootAgreesWithGit` and the `git rev-parse` results of
     requirement 4, and compare them with the gitrepository-layout(5)
     passage `_prd.md` → Acceptance evidence quotes ("a valid Git
     repository must have the HEAD file").
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `None.` under `## Glossary`, and that the
    Task commit range changes no glossary entry in `CONTEXT.md`.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `CONTEXT.md`,
      `docs/references/skill-coverage.json`,
      `docs/references/coverage-record.json`, anything under
      `internal/baseline/`, `internal/store/`, `internal/worktree/`,
      `internal/spec/`, `internal/speccheck/` or `docs/agents/`, any skill
      other than the Roundfix skill files task_03 declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0352-only-a-git-entry-that-git-accepts-marks-a-git-root --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcripts and the report through the binary and Git.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A `.git` that Git rejects marks no Git root, every shape Git accepts still does, and Roundfix agrees with `git rev-parse` on every shape Git answers.
- [ ] Under the MacBook's `~/.git` the commands print what the transcripts say, and no Project Config lands in the home directory.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0352-only-a-git-entry-that-git-accepts-marks-a-git-root/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 1-4; Declared breaks; Testing Approach; Build Order 4
