---
task: task_03
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Spec Consistency
Check resolve an `interface:` or `instruction:` Task Context path that the
checked tree no longer holds when `HEAD` or the Delivery Base tracks it,
changes the Task Context `SC-REF-UNRESOLVED` fix to suggest removing the
entry or declaring the deletion, and teaches `write-tasks`, its Task
template, the Roundfix `spec` reference, the spec command guide and
`CONTEXT.md` that a removed file is listed only under `deletes:`.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network or fetches from a remote, sends a request to
a provider, the judge or a reviewer, or reads a credential. Every `go test`
it runs carries an explicit `-timeout`. Every row that runs `bin/roundfix`
sets `HOME` and `TMPDIR` to fresh temporary directories, with `GOMODCACHE`
and `GOCACHE` kept at their `go env` values, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestTaskContextPathTrackedAtTheBaseResolves` over `./internal/speccheck`,
   `TestAGatedTaskWhoseInterfacePathItDeletesSettles` and
   `TestSettlementStillRefusesAContextPathNoRevisionTracks` over
   `./internal/daemon`, and
   `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes` over
   `./skills`, then the whole `./internal/speccheck`, `./internal/daemon`,
   `./internal/spec`, `./internal/cli` and `./skills` packages (with
   `-timeout 40m` for `./internal/cli`) and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the Fiscus case through the built binary in a temporary
   Git repository the row creates and removes, with commits made under a
   throwaway identity passed by `-c user.name=… -c user.email=…`: on branch
   `main`, commit `src/gone.txt` and a Spec `docs/specs/9001-fixture/` whose
   `_prd.md` has `status: active`, whose `_tasks.md` declares
   `qa: declined` with a `qa_reason`, and whose one pending Task lists
   ``- interface: `src/gone.txt` ``; write `refs/remotes/origin/main` with
   `git update-ref`, then check out a branch `work`. Run
   `bin/roundfix spec check 9001-fixture` (with the temporary `HOME` and
   `TMPDIR`) and record, counting the lines that start with
   `[error] SC-REF-UNRESOLVED:` and name `declares unresolved Task Context path`:
   - none after removing `src/gone.txt` from the working tree only;
   - none after `git rm` and a commit of `src/gone.txt` on `work`;
   - one, naming `src/typo.txt`, after the Context entry is changed to
     `src/typo.txt` instead, whose `fix:` line starts with
     `Remove the Task Context entry for src/typo.txt` and does not contain
     `Create`.
   Remove the temporary directories and the repository afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `taskContextTrackedAtBase` made to return false:
     `TestTaskContextPathTrackedAtTheBaseResolves` fails on
     `working tree removal` and `committed removal`;
   - the Delivery Base left out (only `HEAD` read): it fails on
     `committed removal`;
   - `HEAD` left out (only the Delivery Base read): it fails on
     `branch file removed in the working tree`;
   - the top-level guard removed: it fails on
     `nested below another work tree`;
   - a non-empty `git log --all -- <path>` accepted as tracked: it fails on
     `removed on main before the branch` and
     `created and removed on the branch`;
   - the fix restored to "Create <path> or update the Task Context entry in
     <task>.": it fails on `never tracked`, and
     `TestSettlementStillRefusesAContextPathNoRevisionTracks` fails;
   - "still lists it once" removed from the mirror
     `skills/write-tasks/SKILL.md`:
     `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes` fails.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST search `docs/user-guide/commands/spec.md`, `CONTEXT.md` and every copy
   of the `write-tasks` and Roundfix skills, and record that both copies of
   `write-tasks` and of its Task template say a removed file is listed only
   under `deletes:` and no longer carry "declared under `deletes:`, never
   `instruction:`"; that both copies of the Roundfix `spec` reference and the
   spec command guide state the `HEAD` and Delivery Base rule and the new
   fix; that `CONTEXT.md` changes **Delivery Base** as task_02 requires; and
   that `grep -rn 'or update the Task Context entry' internal skills docs .agents`
   finds nothing outside this Spec's folder.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Fiscus report quoted in
     `references/2026-11-04-a-task-context-path-the-task-deletes-fails-its-own-settlement.md`
     (2026-11-04): compare offline its sequence (an `interface:` entry for a
     file the Task deleted, then `SC-REF-UNRESOLVED`, then a swap to another
     path) with the results of requirements 2 and 3;
   - the Git manual passages `_prd.md` → Acceptance evidence quotes for
     `git cat-file -e` and `<rev>:<path>`: compare them with the
     `nested below another work tree` row and with the revisions API
     Contract 1 reads.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes: **Delivery Base**`, that `CONTEXT.md`
   still defines it in the `**Term**:` form with the added sentence, and that
   the Task commit range changes no other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/references/skill-coverage.json`,
      anything under `internal/baseline/assets/` or `docs/agents/`, the
      `domain-modeling` skill, any skill other than the Roundfix and
      `write-tasks` files the Tasks declare, and any Go source of
      `internal/speccheck` other than `citations.go` and the new
      `task_context_base.go`;
    - no `### QA settlement` section changed.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0348-a-task-context-path-the-task-deletes-is-not-unresolved --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named tests and reproduce the Fiscus case through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A removed Task Context path that `HEAD` or the Delivery Base tracks reports nothing, and every other missing path reports one finding with the new fix, in the tests and through the binary.
- [ ] A gated Task that deletes its `interface:` path settles.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0348-a-task-context-path-the-task-deletes-is-not-unresolved/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-3; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3
