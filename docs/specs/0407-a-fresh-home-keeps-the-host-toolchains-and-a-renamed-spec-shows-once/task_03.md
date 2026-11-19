---
task: task_03
spec: 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec adds a section on a fresh home to the `qa-gate` skill and a QA environment section to the `write-tasks` skill, pinned by a contract test, and makes the status view of the Follow-up Command drop a phase entry of a renamed Spec's earlier slug, with the `follow-up` guide, the glossary entry **Spec Queue Table**, ADR-0410, ADR-0434 and the Skill Coverage Map stating the rule.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session or an ACP adapter;
- starts a Run, a Daemon or a Delivery Owner;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- sends a request to a provider, the judge or a reviewer;
- asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the section this Spec adds to the `qa-gate` skill. Every `go test` package or test run, and the repository Verification, keeps the session's own `HOME` and `TMPDIR`, except the rows of requirement 5 that run under a fresh `HOME` on purpose. The `qa-gate` skill installed for the gate may predate that section, because this Spec is not merged yet, so this Task states it: when the gate isolates `HOME` for any row that runs `./internal/cli` or the repository Verification, it keeps the host Rust toolchain readable, with `RUSTUP_HOME` and `CARGO_HOME` pointing read-only at the host's and `PATH` unchanged; only the first run of requirement 5 leaves them unset, on purpose. Every row that runs the product runs `bin/roundfix`, the binary `make build` made from the audited head, and never a `roundfix` found on `PATH`. Each such row prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running, with `-v`, `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome` and recording every subtest by name, then the whole `./skills` package and `make skills-sync-check`. Record that both copies of `qa-gate` hold the section `### A fresh home keeps the host toolchains` of `_techspec.md` → API Contract 2 once, before `### QA settlement`, that both copies of `write-tasks` hold the section `## QA environment` once, that `TestSettlementGuidanceIsOneTable` passes, and that `skills/testdata/owned-skill-versions.json` records the version each `SKILL.md` of the two skills declares.
3. MUST verify Success Metric 2 by running, with `-v`, `TestARenamedSpecLeavesNoAuthoringRow`, `TestASpecIsBeingAuthored` and `TestTheQueueRowsComeInTableOrder`, recording every test and subtest by name, then the whole `./internal/cli` package and `-tags docscontract ./internal/docscontract`, recording each pass.
4. MUST verify Success Metric 3 through `bin/roundfix`, in a fresh home prepared as above, on a temporary Git repository whose `origin` is `https://github.com/owner/widget.git` and whose commit holds `docs/specs/0401-the-new-slug/_prd.md` with `status: active`:
   - run `bin/roundfix phase start tasks --spec 0401-the-old-slug` there, then `bin/roundfix follow-up --status --format json` and `bin/roundfix follow-up --status`, and record that the `queue` holds one row whose `spec` starts with `0401`, `0401-the-new-slug`, and that the text table prints one `0401` row;
   - run `bin/roundfix phase start prd --spec 0402-other` there, and record that the queue then also holds `0402-other` with the status `authoring: prd`;
   - run the same three commands with the binary the gate builds from the Delivery Base, in a temporary copy made with `git archive` and committed in a fresh Git repository outside the checkout, and record that it lists `0401-the-old-slug` as `authoring: tasks`, which is the defect;
   - record that both binaries exit with the same code and print the same table header and `Forecast basis:` line.
5. MUST record the outside-evidence row of `_prd.md` → Success Metric 4 and Acceptance evidence:
   - read offline the `qa_override_reason` of the Archive Records of Specs 0400 and 0401 under `docs/history/specs/`, and, with `git show <source_revision>:<path>`, row Q04 of the 0399 QA Report `qa-report-2026-11-16.md` and row Q02 of the 0402 QA Report `qa-report-2026-11-17.md`, at the source revisions their Archive Records name, recording the quoted cause of each;
   - quote the rustup book passage the PRD quotes, as published literature;
   - when the host has rustup, run `go test -count=1 -timeout 15m -v -run '^TestGuidanceCompositionJourney$/^rust-cli$' ./internal/cli` with `HOME` and `TMPDIR` fresh temporary directories outside the checkout, the host's `GOMODCACHE` and `GOCACHE`, and `RUSTUP_HOME` and `CARGO_HOME` unset, and record that it fails naming rustup; run it again adding `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"` as taken above, and record that it passes; then run the whole `./internal/cli` package that second way, and record that it passes. When the host lacks rustup, record that and the cause from the Archive Records instead.
6. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Each sabotage keeps the package compiling, so the named test fails on its assertion; record that it fails on the named subtest. No sabotage regenerates or synchronizes anything.
   1. `followUpHeldSpecNumbers` no longer reads the Archive Root: `TestARenamedSpecLeavesNoAuthoringRow/renamed_to_an_archive_record`.
   2. `followUpHeldSpecNumbers` skips a merged Delivery Queue item: `TestARenamedSpecLeavesNoAuthoringRow/renamed_to_a_merged_queue_item`.
   3. `followUpRenamedSlug` ignores a later phase entry of the same Spec Number: `TestARenamedSpecLeavesNoAuthoringRow/renamed_while_authoring`.
   4. `followUpRenamedSlug` treats a slug of another Spec Number as the same Spec: `TestARenamedSpecLeavesNoAuthoringRow/another_number_stays`.
   5. `followUpQueueRows` applies the rule to a candidate that has a folder under the Spec Root: `TestARenamedSpecLeavesNoAuthoringRow/the_old_folder_stays`.
   6. The Rust sentence is removed from the canonical `qa-gate` section: `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/canonical/qa-gate/rust_locations`.
   7. The last sentence of the `qa-gate` section is removed from the mirror `skills/qa-gate/SKILL.md`: `TestTheQASkillsKeepTheHostToolchainsUnderAFreshHome/embedded/qa-gate/a_hidden_toolchain_is_no_environment_block`.
7. MUST execute the tests task_01 and task_02 name in their Verification.
8. MUST search and record:
   - the sentence of API Contract 1 in `docs/user-guide/commands/follow-up.md` and the one in the **Spec Queue Table** entry of `CONTEXT.md`, with its ADR-0434 link;
   - the "Amended by ADR-0434" paragraph of ADR-0410, right after its paragraph "Specs being authored", and its moved `updated_at`;
   - the two `review` values of API Contract 1 in `docs/references/skill-coverage.json`, and that `TestTheSkillCoverageMapIsCurrent` passes;
   - that `followUpHeldSpecNumbers` and `followUpRenamedSlug` call no function that writes, creates, renames or removes a file.
9. MUST perform the glossary check (`docs/agents/domain.md`). Record that the PRD declares `changes:` for **Spec Queue Table**, that task_01 wrote that entry and no other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `skills-lock.json`, `.roundfixrc.yml`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, any file under `internal/baseline/` or `docs/agents/`, any owned skill other than `qa-gate` and `write-tasks`, any Vendored Skill under `.agents/skills`, and any ADR other than ADR-0410;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once --plan --offline; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
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
- [ ] Reproduce the renamed Spec through both binaries.
- [ ] Record the outside evidence and the toolchain reproduction.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Both skills state the fresh-home rule in both copies, and each guarantee is pinned.
- [ ] A renamed Spec is one row of the Spec queue through the product, and another Spec Number keeps its row.
- [ ] The rust-cli journey fails under a fresh home without the toolchain locations and passes with them.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0407-a-fresh-home-keeps-the-host-toolchains-and-a-renamed-spec-shows-once/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Skills; Glossary
- `_techspec.md` → Interfaces; API Contracts 1-2; Testing Approach; Build Order 3
- [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md)
- ADR-0434, ADR-0410, ADR-0080
