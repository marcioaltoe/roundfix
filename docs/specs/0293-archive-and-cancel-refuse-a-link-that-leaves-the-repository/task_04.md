---
task: task_04
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `roundfix archive`
and `roundfix cancel` refuse, before writing anything, every path whose
resolved target lies outside the repository, while a link that stays inside,
a checkout reached through a link and an external Spec Root keep working
with today's bytes.

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
   - In a temporary directory resolved with `pwd -P`, with `HOME` and
     `GIT_CONFIG_GLOBAL` set to files under it, build the repository
     `_techspec.md` → Current behavior and Surface Transcripts describe: one
     commit holding `docs/agents/agent-instructions.md`,
     `docs/agents/domain.md` and the Spec `0001-demo` (a `_prd.md` with its
     Project Constraints, an approved `_authorization.md`, a `_tasks.md` with
     `qa: task_02`, a completed `task_01.md` and a pending QA `task_02.md`,
     each with a Verification bullet). Make one fresh copy per case.
   - With `docs/history` committed as a link to `<tmp>/outside`, run
     `bin/roundfix archive 0001-demo --qa-override --approval probe --reason
     probe` from the repository root and record that standard output,
     standard error and the exit code match Surface Transcript 1.
   - With an adopted `2026-10-20-entry.md` indexed as `backlog` in the Spec's
     `references/_index.md` and `docs/backlog` committed as a link to
     `<tmp>/outside`, run `bin/roundfix cancel 0001-demo --reason probe` and
     record that it matches Surface Transcript 2.
   - Record the refusal, exit `2`, for: `archive … --promote
     references/notes.md` with `docs/references` linked outside; `archive …`
     with the Spec folder itself a committed link to an outside directory;
     `archive 0001-demo --plan --offline` with `docs/history` linked outside;
     `cancel …` with `docs/history` linked outside.
   - For every case above, record that the outside directory holds no file
     the command wrote, that the Spec folder's files are byte-identical to
     the commit, and that `git status --porcelain --untracked-files=all` is
     empty.
3. MUST verify Success Metric 3 through the built product, each against a
   control copy without links and the same archive command: `docs/history` a
   committed link to `docs/hist-real` inside the repository; the command run
   from a link to the repository root (`cd` through the link, so `PWD` names
   it); and an external Spec Root in its own Git repository, configured in
   `.roundfixrc.yml` as `specs.root` through a link spelling of that
   repository. Record that each exits 0, prints the control's standard output
   and standard error with only the record path's spelling allowed to differ
   for the external root, and writes a record byte-identical (`cmp`) to the
   control's.
4. MUST verify Success Metric 4 and the sabotage proofs:
   - Run `TestRepositoryPathRefusesALinkThatLeavesTheRepository`,
     `TestRepositoryPathKeepsLinksThatStayInside`,
     `TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting`,
     `TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting`,
     `TestArchiveAndCancelKeepLinksThatStayInside`,
     `TestArchiveRefusesALinkThatLeavesTheRepositoryTranscript`,
     `TestCancelRefusesALinkThatLeavesTheRepositoryTranscript` and
     `TestArchiveKeepsLinksThatStayInsideTranscript`, and record each
     subtest's outcome (Invariants 1-4).
   - On temporary copies of the audited head made with `git archive`, one
     sabotage per copy, record that each named test fails:
     - make `RepositoryPath` return the lexical relative path before the
       containment test, as before this Spec:
       `TestRepositoryPathRefusesALinkThatLeavesTheRepository` and
       `TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting` fail;
     - resolve only the path, comparing it with the unresolved root:
       `TestRepositoryPathKeepsLinksThatStayInside` fails;
     - remove the archive's promotion destination check:
       `TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting` fails;
     - remove the cancel's returned-source destination check:
       `TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting` fails.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - run `go doc os.Root` and `go doc path/filepath.Rel` and record that the
     passages the PRD quotes occur in their whitespace-normalized output;
   - because this Spec's authorization forbids network access, record the Go
     blog post as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`;
   - build a binary from a `git archive` copy of the Delivery Base, run the
     `docs/history` case of Requirement 2 with it on a fresh copy, and record
     that it exits 0 and writes `specs/0001-demo.md` into the outside
     directory, the symptom the adopted Backlog Entry reports.
7. MUST record that the PRD declares no Unreachable Acceptance, so
   `rows_blocked_declared` is 0 and no row is recorded `blocked (declared
   unreachable)`.
8. MUST perform the glossary and documentation checks:
   - record that the PRD and TechSpec declare no term, that the glossary
     contract tests pass, and whether any term the Spec introduced needs an
     entry (`docs/agents/domain.md`);
   - record the paragraph each of `docs/user-guide/commands/archive.md` and
     `docs/user-guide/commands/cancel.md` gained, the paragraph both copies of
     the Roundfix `archive` reference gained, the raised and recorded skill
     version, the line ADR-0269 gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says either command accepts or follows a symbolic link whose target
     lies outside the repository, and record what it finds.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every file under `internal/config/`,
     `docs/references/skill-coverage.json`, every file under `docs/adr/` other
     than ADR-0269, every file under `skills/` and `.agents/skills/` other
     than the Roundfix `SKILL.md` pair, the Roundfix `references/archive.md`
     pair and `skills/testdata/owned-skill-versions.json`, and every file
     under `docs/user-guide/` other than `commands/archive.md` and
     `commands/cancel.md`;
   - `internal/cli/archive.go` and `internal/cli/cancel.go` are byte-identical
     to the Delivery Base, unless task_02's Result names why one changed;
   - every `_test.go` file under `internal/spec` and `internal/cli` that
     exists at the Delivery Base is byte-identical to it;
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
    `bin/roundfix archive 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository --plan --offline`
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
- [ ] Reproduce Surface Transcripts 1 and 2 and every outside and inside case through the built binary.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Record the outside evidence, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every outside case exits 2 with its refusal naming the path and writes nothing, and each removed guard makes its named test fail.
- [ ] Every inside, linked-checkout and external Spec Root case archives with the control's output and record bytes.
- [ ] No Archive Record already written changed, and the scope audit finds no undeclared Governed Path change.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Declared break; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Testing Approach; Build Order 4
- ADR-0306; ADR-0269; ADR-0247; ADR-0260; ADR-0262; ADR-0268; ADR-0104; ADR-0240; ADR-0276
