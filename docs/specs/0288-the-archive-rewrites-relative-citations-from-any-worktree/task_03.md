---
task: task_03
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `roundfix archive`
and `roundfix cancel` compute every repository-relative path from one
repository root, so both work from any spelling of the checkout: a linked
worktree or the main checkout, its root or a subdirectory, through a symbolic
link or not.

Every row runs on temporary repositories, homes and Git configuration under
the Run's temporary directory. No row reads or writes the real `~/.roundfix`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. The gate never runs
`roundfix spec judge`, never runs `roundfix review` against a real provider,
and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 through the built product:
   - In a temporary directory, with `HOME` and `GIT_CONFIG_GLOBAL` set to
     files under it, build the repository `_techspec.md` → Current behavior
     and Surface Transcripts describe: a main checkout at `<tmp>/real/main`
     whose one commit holds `docs/agents/agent-instructions.md`,
     `docs/agents/domain.md`, the Spec `0001-demo` (a `_prd.md` whose first
     paragraph links `../../backlog/2026-10-20-follow-up.md`, an approved
     `_authorization.md`, a `_tasks.md` with `qa: task_02`, a completed
     `task_01.md` and a pending QA `task_02.md`, each with a `##
     Verification` bullet) and a Backlog Entry linking
     `../specs/0001-demo/_prd.md`; a bare `origin` holding `main` with
     `origin/HEAD` set; a linked worktree added with `git worktree add -b
     item <tmp>/real/item main`; and `ln -s <tmp>/real <tmp>/link`.
   - Run `bin/roundfix archive 0001-demo --qa-override --approval probe
     --reason probe` with the work directory `<tmp>/link/item` (`cd` through
     the link, so `PWD` names it), and record that standard output, standard
     error and the exit code match Surface Transcript 1.
   - On a fresh copy of that repository, run the same command from
     `<tmp>/real/item`, `<tmp>/link/item/docs` and `<tmp>/link/main`, one
     fresh copy each, and record that every run prints the same lines and
     that its `docs/history/specs/0001-demo.md` and Backlog Entry bytes equal
     the resolved run's (`cmp`).
   - On a fresh copy, run `bin/roundfix cancel 0001-demo --reason probe` from
     `<tmp>/link/item` and from `<tmp>/real/item`; record that the linked run
     matches Surface Transcript 2 and that both records and Backlog Entries
     are byte-identical.
3. MUST verify Success Metrics 3 and 4 and the sabotage proofs:
   - Run `TestRepositoryPathComparesResolvedSpellings`,
     `TestArchiveFromALinkedSpellingRewritesAndRelinks`,
     `TestCancelFromALinkedSpellingRelinksAndReturnsSources`,
     `TestArchiveFromEverySpellingOfTheCheckoutTranscript`,
     `TestCancelFromALinkedSpellingTranscript`,
     `TestArchiveRefusalsNameRepositoryPathsFromALinkedSpelling` and
     `TestArchivePlanWithALinkedExternalSpecRootListsTheRelink`, and record
     each subtest's outcome (Invariants 1-4).
   - On a temporary copy of the audited head made with `git archive`, make
     `RepositoryPath` return the lexical relative path without the resolved
     comparison, and record that
     `TestArchiveFromALinkedSpellingRewritesAndRelinks` fails with `rebase
     archive outcome` and that
     `TestArchiveFromEverySpellingOfTheCheckoutTranscript` fails. On another
     copy, make it return the lexical path even when it starts with `../`,
     and record that `TestRepositoryPathComparesResolvedSpellings` fails.
4. MUST execute the tests task_01 and task_02 name in their Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote from the adopted Inbox Entry under `references/` the refusal it
     recorded and the item worktree path under `~/.roundfix/worktrees`; then
     build a binary from a `git archive` copy of the Delivery Base, run the
     archive command of Requirement 2 with it from `<tmp>/link/item` on a
     fresh copy of the repository, and record that it exits 2 with `rebase
     archive outcome: rebase citation`, the class the Inbox Entry reported;
   - run `go doc os.Getwd` and `go doc path/filepath.Rel` and record that the
     passages the PRD quotes occur in their output;
   - because this Spec's authorization forbids network access, record the
     Git source and the Visual Studio Code issue as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the passages the PRD
     cites. Such a row never decides a qualifying `partial`.
6. MUST record Success Metric 5 as one row, `blocked (declared unreachable)`,
   against the PRD's single Unreachable Acceptance entry, counted once in
   `rows_blocked_declared`: no row may read the real `~/.roundfix`.
7. MUST perform the glossary and documentation checks:
   - record that the PRD and TechSpec declare no term, that the glossary
     contract tests pass, and whether any term the Spec introduced needs an
     entry (`docs/agents/domain.md`);
   - record the sentence each of `docs/user-guide/commands/archive.md` and
     `docs/user-guide/commands/cancel.md` gained, that
     `docs/references/skill-coverage.json` carries the dated Coverage Review
     on the four surfaces `_prd.md` → Skills lists, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says either command must run from the main checkout or the
     repository root, and record what it finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every file under `skills/` and
     `.agents/skills/`, every file under `docs/user-guide/` other than
     `commands/archive.md` and `commands/cancel.md`, `internal/spec/git_locator.go`,
     `internal/spec/legacy_link_rebase.go`, `internal/spec/live_citations.go`
     and `internal/cli/deliver_workflow.go`;
   - every `_test.go` file under `internal/spec` and `internal/cli` that
     exists at the Delivery Base is byte-identical to it;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0288-the-archive-rewrites-relative-citations-from-any-worktree --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Reproduce Surface Transcripts 1 and 2 and compare every spelling's bytes through the built binary.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Record the outside evidence, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The archive and the cancel print and write the same bytes from every spelling of the checkout, and each removed guard makes its named test fail.
- [ ] Refusals from a linked spelling name repository-relative paths, and a path outside the repository is refused by name.
- [ ] No Archive Record already written changed, and the scope audit finds no undeclared Governed Path change.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0288-the-archive-rewrites-relative-citations-from-any-worktree/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Skills; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0247; ADR-0260; ADR-0262; ADR-0268; ADR-0104; ADR-0240; ADR-0276
