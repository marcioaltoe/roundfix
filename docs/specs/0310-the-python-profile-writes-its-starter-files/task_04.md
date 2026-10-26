---
task: task_04
spec: 0310-the-python-profile-writes-its-starter-files
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec lets a Baseline module
declare Starter Files, gives the `python` module a Makefile, a Project Config
and a CI workflow, adds the confirmation-gated `roundfix baseline starters`
command that writes them only where absent, names it in the profile
alignment's next action, and documents it with the glossary term **Starter
File**.

Every row runs on temporary copies of the audited head and temporary homes
and repositories under the Run's temporary directory. No row reads or writes
the real `~/.roundfix`, reads or writes the Secondbrain, reads another
repository, runs uv, pytest, ruff, a starter Makefile target or a workflow,
reaches the network, sends a request to a provider, the judge or a reviewer,
or reads a credential. The gate never runs `roundfix spec judge` and never
runs an archive plan without `--offline`. Every command that runs
`bin/roundfix` or `go run ./cmd/roundfix` sets `HOME` and `TMPDIR` to fresh
empty directories created for that row under the Run's temporary directory
(the CLI loads User Config from `$HOME/.roundfix/config.yml`).

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 (the declarations):
   - Run `TestThePythonModuleDeclaresItsStarterFiles` and
     `TestAStarterFileDeclarationIsValidated` in `./internal/baseline`, and
     record the three paths, their identities and the four codes.
   - On temporary copies of the audited head, one sabotage per copy, record
     that the named case fails:
     - remove the template kind check of the Starter File validation:
       `TestAStarterFileDeclarationIsValidated/guide_template` fails;
     - delete the line `verify-incremental: verify` from
       `internal/baseline/assets/templates/starters/python/Makefile`:
       `TestThePythonModuleDeclaresItsStarterFiles` fails naming it;
     - remove ` typecheck` from the `verify:` line, or `--offline` from the
       `lock-check` recipe, of the same file:
       `TestThePythonModuleDeclaresItsStarterFiles` fails naming the line.
3. MUST verify Success Metrics 3 and 4 (the operation):
   - Run `TestStartersPreviewsAndCreatesOnlyAbsentFiles`,
     `TestStartersRefuseADigestThatNoLongerMatches`,
     `TestStartersRemoveWhatTheyCreatedWhenAWriteFails` and
     `TestTheUndeclaredGateNamesTheStartersCommand`, and record that they pass.
   - On temporary copies, one sabotage per copy, record that the named case
     fails:
     - skip the digest comparison in `Starters`:
       `TestStartersRefuseADigestThatNoLongerMatches` fails;
     - skip the removal of created files on a write failure:
       `TestStartersRemoveWhatTheyCreatedWhenAWriteFails` fails;
     - make the preview mark every path `create`:
       `TestStartersPreviewsAndCreatesOnlyAbsentFiles` fails on the Makefile
       entry.
4. MUST verify Success Metric 5 (the command), with `bin/roundfix` and `HOME`
   and `TMPDIR` set to fresh temporary directories, in a temporary Git
   repository with one commit:
   - run `bin/roundfix baseline starters --profile python --repo <repo>`
     and record the exit code, every stdout line and an empty stderr against
     Surface Transcript 1, and that no file was written;
   - run it with `--confirm-plan <digest>` and record the result against
     Surface Transcript 2, the three files' bytes against the templates, and a
     third run's `State: current` with exit 0;
   - in a fresh repository holding a `Makefile` of its own, record that the
     preview lists `- keep Makefile`, and that after the apply the Makefile is
     byte-identical;
   - run with a digest of another plan and record exit 3, empty stdout and the
     stale line on stderr, and that nothing was written;
   - run `bin/roundfix baseline starters --help` and `bin/roundfix --help` and
     record the synopsis and help sentences;
   - run `bin/roundfix baseline plan --profile python` with the release
     journey's decisions in a repository without a Makefile and record the
     `verification.command.undeclared` next action;
   - run the five tests of `internal/cli/baseline_starters_test.go`,
     `TestBaselineDocumentationContract`, `TestBaselineExamplesParse`,
     `TestEveryCommandPathPrintsItsHelp` and `TestTheSkillCoverageMapIsCurrent`,
     and record that they pass.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST compare the written files with `_prd.md` → Acceptance evidence and
   record, line by line, each quoted Quaestor line the starters reproduce and
   each deliberate difference the PRD's Validation names (`--frozen` in the
   Makefile, `lock-check` and `typecheck` in `verify`, `--offline` on the lock
   check, `--locked` in CI, no `rtk` in `verification.tools`). Because
   this Spec's authorization forbids network access, record the uv
   documentation as `blocked (environment: network prohibited by Spec
   authorization)`, quoting the passages the PRD cites; such a row never
   decides a qualifying `partial`.
7. MUST record the real-repository criterion as exactly one row, `blocked
   (declared unreachable)`, against the PRD's single Unreachable Acceptance
   entry, and count it once in `rows_blocked_declared`.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Starter File** with its ADR-0323
     link, that the glossary contract tests pass, and whether any other term
     the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record every command and flag the guide and the Roundfix reference name
     for `baseline starters`, and that the binary's help, run with `HOME` and
     `TMPDIR` set to fresh temporary directories, lists each; a named
     command or flag the binary lacks is a finding;
   - record that `make skills-sync-check` passes, the Roundfix skill's version
     at the Delivery Base and at the audited head, and that
     `TestNoPythonBaselineRuntime` passes.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     other than `python.json`, every profile and setup, everything under
     `internal/baseline/assets/source-baselines/`,
     `internal/baseline/assets/retention/` and
     `internal/baseline/assets/formatter-fixtures/`, every file under
     `docs/history/`, every owned skill other than `roundfix`, and every user
     guide page other than `docs/user-guide/commands/baseline.md`;
   - in `python.json`, only the module version and the `starterFiles` array
     changed;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0310-the-python-profile-writes-its-starter-files --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    and `TMPDIR` each set to a fresh empty temporary directory created for this row (the CLI
    loads User Config from `$HOME/.roundfix/config.yml`, so the default home
    would read `~/.roundfix`), and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the command journeys in temporary repositories, and each sabotage on temporary copies of the audited head.
- [ ] Compare the starters with the quoted adopter files, and record the outside sources and the declared row.
- [ ] Record the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The `python` Starter Files resolve and validate, and each sabotage makes its named case fail.
- [ ] The command previews, applies only absent paths, keeps existing files, refuses a stale digest and matches both transcripts.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 9 names.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0310-the-python-profile-writes-its-starter-files/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Surface Transcripts; Invariants 1-6; Testing Approach; Build Order 4
- ADR-0323
- ADR-0068
- ADR-0322
