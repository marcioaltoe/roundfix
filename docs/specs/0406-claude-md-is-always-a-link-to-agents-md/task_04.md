---
task: task_04
spec: 0406-claude-md-is-always-a-link-to-agents-md
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds a read-only reader of the **Claude Instruction Link**. Every adoption plan writes `CLAUDE.md` as a symbolic link to `AGENTS.md`; a managed refresh does so when the file holds nothing of its own, and otherwise warns and leaves it. The file transaction replaces and restores that one path. An unlinked checkout is left and reported. `roundfix doctor` gains a `baseline-claude-link` line that never fails. The guide, the user guides, the Roundfix and `setup-context-driven` skills and the glossary state the rule.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head or of the Delivery Base. The one exception is the outside-evidence row of requirement 7, which reads the adopted checkouts the inventory names, read only. No row:

- starts an Agent harness, an Agent session or an ACP adapter;
- starts a Run, a Daemon or a Delivery Owner;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix` or the real home directory;
- sends a request to a provider, the judge or a reviewer;
- asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`; requirement 6 alone also runs a binary the gate builds from the Delivery Base. Each such row prepares its environment the same way:

- it takes `m="$(go env GOMODCACHE)"` and `c="$(go env GOCACHE)"` first;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"` and `GOCACHE="$c"`;
- a fresh `HOME` that runs the whole `./internal/cli` package keeps the host Rust toolchain: `PATH`, `RUSTUP_HOME` and `CARGO_HOME` point at the host's toolchain, read only;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards.

The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Success Metric 2. Run, with `-v`, the `internal/baseline` tests task_01 and task_02 name in their Verification, and record every test and subtest by name. Then, through `bin/roundfix`, in a fresh home, on temporary Git repositories each given this repository's `.agents/skills`, `skills-lock.json` and a `Makefile` with `verify` and `verify-changed` targets, as the inventory's measurement did:
   - adopt one whose `CLAUDE.md` is a regular file of three lines with `bin/roundfix baseline plan --profile go-cli-tui --decision preservation.mode=greenfield` and the other decisions the profile asks for, then `bin/roundfix baseline apply` with the printed Plan Digest. Record that the plan lists `update CLAUDE.md`, that `CLAUDE.md` is afterwards a symbolic link whose target is exactly `AGENTS.md`, that the backup `CLAUDE.<digest>.md` holds the three lines, and that `bin/roundfix baseline update --no-skills --no-history --format text` reports `current`;
   - in that adopted repository, replace `CLAUDE.md` with a regular file holding `@AGENTS.md` and a newline, and record that `bin/roundfix baseline update --no-skills --no-history --yes --format text` replaces it with the link, writes the backup `CLAUDE.<digest>.md` holding those bytes, and that a second run reports `current`;
   - remove `CLAUDE.md`, and record that the same command creates the link;
   - replace it with a regular file holding a heading and a bullet, and record the whole `baseline.root.claude-link.withheld` warning of `_techspec.md` → API Contract 2, that the file is byte-identical after a run with `--yes`, and that the exit code equals the run with the link in place.
3. MUST verify Success Metric 3: on a temporary Git repository adopted as in requirement 2, commit the link, check it out again with `git -c core.symlinks=false checkout`, and record that the plan prints the whole `baseline.root.claude-link.unlinked-checkout` warning, leaves `CLAUDE.md` byte-identical, and reports no Misplaced Repository Rule for it.
4. MUST verify Success Metric 4 and Success Metric 5 by running, with `-v`, the `internal/cli` tests task_03 names in its Verification and `TestClaudeInstructionLinkRollbackRestoresTheFile`, recording every test and subtest by name. Then run the whole `./internal/baseline`, `./internal/cli` and `./skills` packages and `-tags docscontract ./internal/docscontract`, recording each pass, and `GOOS=windows go build -buildvcs=false ./internal/baseline ./internal/cli`. Through `bin/roundfix doctor`, in the repository of requirement 2, record only the `baseline-claude-link:` line: `ok` with the link, the `warn` line for the regular file with a passage, and `skipped` in a directory with no Setup Manifest. In each case Doctor's exit code equals the one the Delivery Base's binary of requirement 6 gives in the same place.
5. MUST reproduce Surface Transcript 1 of `_techspec.md` through `bin/roundfix baseline update --repo . --no-skills --no-history --format text` in this repository. Record standard output, standard error and the exit code, compared line by line with the transcript, where `<digest>` is whatever the command prints and `...` the lines after `Declared decisions: 0`. Record that no line names `baseline.root.claude-link` and that `CLAUDE.md` is a symbolic link to `AGENTS.md`.
6. MUST compare the audited head with the Delivery Base, the merge base of the audited head with the default branch, through the product. Build the Delivery Base's binary in a temporary copy made with `git archive` and committed in a fresh Git repository, outside the checkout. On an adopted fixture whose `CLAUDE.md` is already the link, record that every line the Delivery Base's `baseline update` plan prints is printed by the audited head's in the same order, apart from a catalog digest or a guide file change that the catalog change of task_02 explains, and that both exit with the same code.
7. MUST record the outside-evidence row of `_prd.md` → Success Metric 6 and Acceptance evidence. For each adopted checkout that [`references/2026-11-18-claude-link-inventory-and-measurements.md`](references/2026-11-18-claude-link-inventory-and-measurements.md) names and the QA host holds:
   - record its commit, `ls -l CLAUDE.md`, `git ls-files -s -- CLAUDE.md` and its `git status --porcelain`;
   - run `bin/roundfix baseline update --repo <checkout> --no-skills --no-history --format json` without confirmation, and record that no file change and no warning names `CLAUDE.md`;
   - record the same `git status --porcelain` afterwards.

   When a checkout's `CLAUDE.md` is no longer the link, record the method's output and the warning or file change the plan prints for it. Name each checkout the host lacks. This row never passes `--yes`, `--confirm-plan` or `--promote`, and never writes in another repository.
8. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository, because the suite guard refuses a tree outside Git. Every `go test` carries `-timeout`. Record that the named test fails on the named subtest. No sabotage regenerates or synchronizes anything.
   1. The reader skips the Git check and calls every file with the bytes `AGENTS.md` an unlinked checkout: `TestReadClaudeInstructionLinkClassifiesEveryShape/agents_bytes_tracked_as_a_file`.
   2. The reader counts the `@AGENTS.md` line as a passage: `TestReadClaudeInstructionLinkClassifiesEveryShape/import_line`.
   3. `validateManagedRefreshPreservation` loses its exception for the link: `TestBaselinePlansTheClaudeInstructionLink/refresh_import_line_converted`.
   4. A managed refresh converts a file with a passage: `TestBaselinePlansTheClaudeInstructionLink/refresh_file_with_a_rule_withheld`.
   5. The rollback leaves the link in place: `TestClaudeInstructionLinkRollbackRestoresTheFile/regular_file`.
   6. The Doctor line returns `failed` for a regular file: `TestDoctorReportsTheClaudeInstructionLink/regular_file`.
9. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
10. MUST search and record:
    - the guide sentence of `_techspec.md` → API Contract 4 in this repository's `docs/agents/agent-instructions.md` and the Standard TypeScript formatter golden;
    - the text task_03 prescribes in the user guides of `baseline` and `doctor`, both copies of the Roundfix `baseline`, `modes` and `setup` references and both copies of the `setup-context-driven` skill, and that each copy under `skills/` equals its `.agents/skills/` source;
    - that `internal/baseline/claude_link.go` calls no function that writes, creates, renames, changes the mode of or removes a file, and runs no command other than `git ls-files -s -- CLAUDE.md`;
    - the "Amended by ADR-0432" paragraphs of ADR-0100 and ADR-0411.
11. MUST perform the glossary check (`docs/agents/domain.md`). Record that the PRD declares `adds:` for **Claude Instruction Link** and `changes:` for **Managed Refresh** and **Misplaced Repository Rule**, that task_01 and task_02 wrote those entries and no other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, `internal/cli/cli_test.go`, `internal/baseline/plan_test.go`, `.roundfixrc.yml`, `docs/agents/specific-repository.md`, `AGENTS.md`, `CLAUDE.md`, any Vendored Skill under `.agents/skills`, any Source Baseline under `internal/baseline/assets/source-baselines/`, any Setup Snapshot under `internal/baseline/assets/setups/`, any profile other than `standard-typescript-monorepo.json`, any ADR other than ADR-0100 and ADR-0411, and any owned skill other than `roundfix` and `setup-context-driven`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0406-claude-md-is-always-a-link-to-agents-md --plan --offline; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a reviewer;
    - write in an adopted checkout, or pass `--yes`, `--confirm-plan` or `--promote` there;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the three Tasks and the packages.
- [ ] Reproduce the adoption, refresh, unlinked-checkout and Doctor rows, the Surface Transcript and the Delivery Base comparison through the product.
- [ ] Re-read the adopted checkouts read only and compare with the inventory.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Adoption always ends with `CLAUDE.md` the link and the replaced bytes backed up; a managed refresh converts only what holds nothing of its own and names each passage otherwise.
- [ ] An unlinked checkout is left and reported, and Doctor prints the line in every form without changing its exit code.
- [ ] The adopted checkouts plan nothing for `CLAUDE.md`, and no checkout changes.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0406-claude-md-is-always-a-link-to-agents-md/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-5; Surface Transcript 1; Interfaces; Testing Approach; Build Order 4
- [`references/2026-11-18-claude-link-inventory-and-measurements.md`](references/2026-11-18-claude-link-inventory-and-measurements.md)
- ADR-0432, ADR-0433, ADR-0100, ADR-0411, ADR-0070, ADR-0071, ADR-0203, ADR-0220, ADR-0285, ADR-0341
