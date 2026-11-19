---
task: task_03
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes the status view of the Follow-up Command count a Delivery Queue item's Tasks on its Item Branch through a Task Graph loader that takes a reader, and makes the Release Plan count every Adopter Maintenance Path as maintenance-only in an adopter, with the `follow-up` guide, the `release plan` help, the release runbook, the usage guide, both copies of the Roundfix release reference, the glossary entries **Spec Queue Table** and **Adopter Maintenance Path**, ADR-0447 and the Skill Coverage Map stating the rules.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session or an ACP adapter;
- starts a Run, a Daemon or a Delivery Owner;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- reads another repository, except a temporary one the row itself created;
- sends a request to a provider, the judge or a reviewer;
- asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` skill's section "A fresh home keeps the host toolchains". Every `go test` package or test run, and the repository Verification, keeps the session's own `HOME` and `TMPDIR`. Every row that runs the product runs `bin/roundfix`, the binary `make build` made from the audited head, or the binary the gate builds from the Delivery Base, and never a `roundfix` found on `PATH`. Each such row prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running, with `-v`, `TestLoadAtReadsThroughTheGivenReader` in `./internal/spec` and `TestADeliveryItemCountsTasksOnItsItemBranch`, `TestTheQueueRowsComeInTableOrder`, `TestAQueueItemNamesItsWait` and `TestFollowUpStatusListsPendingSpecsInDeliveryOrder` in `./internal/cli`, recording every test and subtest by name, then the whole `./internal/spec`, `./internal/cli` and `./internal/daemon` packages, recording each pass.
3. MUST verify Success Metrics 2 and 3 by running, with `-v`, `TestAnAdopterManagedPathIsMaintenanceOnly` in `./internal/releaseplan`, `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget`, `TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths` and `TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary` in `./internal/cli`, and `TestTheReleaseDocumentsNameTheAdopterMaintenancePaths` in `./skills`, recording every subtest by name, then the whole `./internal/releaseplan` and `./skills` packages, `make skills-sync-check` and `-tags docscontract ./internal/docscontract`. Record that `skills/testdata/owned-skill-versions.json` records the version the Roundfix `SKILL.md` declares in both copies, and that `bin/roundfix release plan --help` prints the paragraph of `_techspec.md` → API Contract 2 and still says it reads no Roundfix configuration.
4. MUST record the outside-evidence row of `_prd.md` → Success Metric 4 and Acceptance evidence:
   - build the Delivery Base binary from a temporary copy of the Delivery Base made with `git archive` and committed in a fresh Git repository outside the checkout, with `make build` there;
   - in a fresh temporary Git repository outside the checkout, rebuild the path shape of the four Fiscus commits from `references/2026-11-19-characterization.md` section 2, with a tag `v0.1.0` on a seed commit holding a Setup Manifest that records `AGENTS.md` and the four `docs/agents/` guides that section names, a `skills-lock.json` recording `roundfix`, `setup-context-driven`, `write-idea`, `write-prd`, `write-tasks` and `write-techspec`, and the repository-local skill `.agents/skills/revisao-tributaria-cadastro/SKILL.md` that the lock does not record; then one commit per Fiscus commit with its subject shape (`chore:` or `docs:`) and its changed paths as the section lists them, using placeholder content;
   - run `release plan --from v0.1.0 --format json` there with `bin/roundfix` and with the Delivery Base binary, in a fresh home prepared as above, and record that the Delivery Base binary blocks four commits and `bin/roundfix` blocks only the fourth, whose crossing paths are `CONTEXT-MAP.md`, the `docs/user-guide/` files and `packages/backend/.env.example`;
   - add one commit there changing only `.agents/skills/revisao-tributaria-cadastro/SKILL.md` and record that `bin/roundfix` reports it as crossing the boundary;
   - run `follow-up --status --format json` there with both binaries and record each one's release state, the field the `Next release` line prints;
   - in this repository's checkout, run `release plan --from <tag> --to HEAD --format json` with both binaries for `v0.140.0`, `v0.147.0` and `v0.148.0`, and record that each pair prints the same `state`, the same blocking commits and the same `crossesMaintenanceOnlyBoundary` for every commit; a clone without one of those tags records that and uses the earliest stable tag it holds;
   - quote googleapis/release-please#1875 as published literature: "Idea is to exclude commits if all files in those commits are from specific paths."
5. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Each sabotage keeps the package compiling, so the named test fails on its assertion; record that it fails on the named subtest. No sabotage regenerates or synchronizes anything.
   1. `followUpItemTaskGraph` always reads the Spec Root: `TestADeliveryItemCountsTasksOnItsItemBranch/a_running_item_counts_its_branch`.
   2. `followUpItemTaskGraph` skips the Archive Record: `TestADeliveryItemCountsTasksOnItsItemBranch/an_archived_item_counts_the_source_revision`.
   3. `followUpItemTaskGraph` loads the item's worktree with `spec.Load` before Git: `TestADeliveryItemCountsTasksOnItsItemBranch/the_worktree_is_not_read`.
   4. `followUpPendingSpecs` keeps `spec.Load` on the Spec Root: `TestADeliveryItemCountsTasksOnItsItemBranch/the_pending_specs_line_counts_the_same`.
   5. `ManagedPaths.Contains` ignores `Adopter`: `TestAnAdopterManagedPathIsMaintenanceOnly/not_an_adopter_changes_nothing`.
   6. `ManagedPaths.Contains` drops `docs/agents/specific-repository.md`: `TestAnAdopterManagedPathIsMaintenanceOnly/specific_repository`.
   7. `ManagedPaths.Contains` accepts every `.claude/agents/*.md`: `TestAnAdopterManagedPathIsMaintenanceOnly/another_claude_agent_crosses`.
   8. The boundary treats a commit without changed paths as maintenance-only: `TestAnAdopterManagedPathIsMaintenanceOnly/an_empty_commit_still_crosses`.
   9. `releasePlanGitSource.ManagedPaths` reads the records at the target only: `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget/a_skill_removed_in_the_range`.
   10. The sentence "Roundfix's own repository authors the `roundfix` skill" is removed from the mirror `skills/roundfix/references/release.md` only: `TestTheReleaseDocumentsNameTheAdopterMaintenancePaths/embedded/roundfix_itself`.
   11. The phrase "a generated .claude/agents/roundfix-*.md" is removed from the `release plan` help: `TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths/generated_subagents`.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST search and record:
   - the text of API Contract 1 in `docs/user-guide/commands/follow-up.md` and in the **Spec Queue Table** entry of `CONTEXT.md`;
   - the text of API Contract 2 in both copies of the Roundfix release reference, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md`, and the **Adopter Maintenance Path** entry of `CONTEXT.md` right after the Release Plan entry, with its ADR-0447 link;
   - the two `review` values of API Contract 1 in `docs/references/skill-coverage.json`, and that `TestTheSkillCoverageMapIsCurrent` passes;
   - that `followUpItemTaskGraph` and `releasePlanGitSource.ManagedPaths` call no function that writes, creates, renames or removes a file, and that `releasePlanGitSource.ManagedPaths` names no `.roundfixrc.yml` and reads no working-tree file.
8. MUST perform the glossary check (`docs/agents/domain.md`). Record that the PRD declares `adds:` for **Adopter Maintenance Path** and `changes:` for **Spec Queue Table**, that task_01 wrote the second and task_02 the first and no other entry changed, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `## Recorded paths`;
   - every changed Governed Path is listed in `_authorization.md` `paths`;
   - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `skills-lock.json`, `.roundfixrc.yml`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, any file under `internal/baseline/` or `docs/agents/`, any owned skill file other than the Roundfix skill's `SKILL.md` and `references/release.md` in both copies, any Vendored Skill under `.agents/skills`, and any ADR;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release --plan --offline; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of both Tasks and the packages.
- [ ] Reproduce the Fiscus range shape and the Roundfix ranges through both binaries.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Delivery Queue item's Tasks are counted on its Item Branch, archived or not, and every fallback reads the Spec Root.
- [ ] In the Fiscus-shaped adopter, only the commit with product paths blocks, and a repository-local skill still blocks.
- [ ] Roundfix's own ranges classify identically with both binaries.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Skills; Glossary
- `_techspec.md` → Interfaces; API Contracts 1-2; Testing Approach; Build Order 3
- [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md)
- ADR-0447, ADR-0412
