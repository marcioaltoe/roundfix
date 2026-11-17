---
task: task_04
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec states in the
Baseline guides, the root pointer and the starter of
`docs/agents/specific-repository.md` that a repository's own rules go only
there, with an example; adds a read-only detector of each **Misplaced
Repository Rule**; reports the list in `roundfix baseline update`, outside
the Plan Document, and in one `baseline-rules` Doctor line that never
fails; and describes both in the user guides, the Roundfix and
`setup-context-driven` skills and the glossary.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head or of the Delivery
Base, except the outside-evidence row of requirement 7, which reads the
adopted checkouts the inventory names, read-only. No row starts an Agent
harness, an Agent session or an ACP adapter, starts a Run, a Daemon or a
Delivery Owner, calls ai-memory, reads or writes the real `~/.roundfix` or
the real home directory (apart from reading the Rust toolchain requirement 4
names), runs the real `gh`, reaches the network, sends a
request to a provider, the judge or a reviewer, or asserts a wall-clock
duration. Every `go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary
`make build` made from the audited head, and never a `roundfix` found on
`PATH`; requirement 6 alone also runs a binary the gate builds from the
Delivery Base. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`,
`CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command,
reads standard input from `/dev/null`, and removes the temporary
directories afterwards. The QA never writes into the checkout while a test
suite runs. The gate keeps its evidence under names the repository does not
ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running, with `-v`, the
   `internal/baseline` tests task_01 names in its Verification, recording
   every test and subtest by name; then by searching and recording that
   this repository's `docs/agents/agent-instructions.md`,
   `docs/agents/skill-dispatch.md`, `docs/agents/docs-layout.md` and
   `AGENTS.md`, and the three Standard TypeScript formatter goldens, carry
   the sentences task_01 prescribes, that the three clauses keep their ids
   and `mandatory` force, and that `bin/roundfix baseline update --repo .
   --no-skills --no-history --format json` reports `current`.
3. MUST verify Success Metric 2 by running, with `-v`, the test task_02
   names in its Verification, recording every subtest by name, then
   `GOOS=windows go build -buildvcs=false ./internal/baseline`.
4. MUST verify Success Metric 3 and Success Metric 4 by running, with `-v`,
   the three `internal/cli` tests task_03 names in its Verification,
   recording every test and subtest by name, and the whole `./internal/cli`
   package, `./skills` package and `-tags docscontract
   ./internal/docscontract`, recording each pass. For the whole
   `./internal/cli` package only, the fresh home keeps the host's existing
   Rust toolchain: take `RUSTUP_HOME` and `CARGO_HOME` before `HOME` changes
   (their values, or `$HOME/.rustup` and `$HOME/.cargo` when unset) and pass
   them to the command, read-only, so `TestGuidanceCompositionJourney/rust-cli`
   can select its toolchain.
5. MUST reproduce Surface Transcript 1 of `_techspec.md` through
   `bin/roundfix baseline update --repo . --no-skills --no-history --format
   text` in this repository, recording standard output, standard error and
   the exit code, compared line by line with the transcript, where
   `<digest>` is whatever the command prints and `...` the lines after
   `Declared decisions: 0`. Then, through `bin/roundfix`, on a temporary
   Git repository adopted by `bin/roundfix` itself in a fresh home, record:
   - with a heading and two lines appended after the Managed Region of
     `docs/agents/skill-dispatch.md`, the plan prints the section
     `Misplaced repository rules: 1` with that path, its line range and the
     heading, and `  Move to: docs/agents/specific-repository.md`; the JSON
     carries the same entry under `misplacedRepositoryRules`; the plan
     exits with the code it exits with before the lines are appended; a
     listing of every file with its size and SHA-256 is the same before and
     after the plan;
   - after `--yes`, the applied result names the passage at the lines of
     the bytes after the apply, and the passage's bytes are unchanged;
   - with a closed `<!-- BEGIN:tool-rules -->` block appended to
     `AGENTS.md`, no entry names it;
   - `bin/roundfix doctor` in that repository prints
     `baseline-rules: warn (DR-BASELINE-RULE-MISPLACED: 1 passage belongs in docs/agents/specific-repository.md: ...` and, after the
     passage is removed, `baseline-rules: ok (no repository rule outside docs/agents/specific-repository.md)`;
     in a directory with no Setup Manifest it prints
     `baseline-rules: skipped (no Setup Manifest at docs/agents/setup-context.json)`;
     in each case Doctor's exit code equals the one the same run gives on
     the Delivery Base's binary of requirement 6. Only the `baseline-rules:`
     line is compared; the other lines depend on the host.
6. MUST compare the audited head with the Delivery Base (the merge base of
   the audited head with the default branch) through the product. Build the
   Delivery Base's binary in a temporary copy made with `git archive` and
   committed in a fresh Git repository, outside the checkout. On the fixture
   of requirement 5 with the passage, record that every line the Delivery
   Base's `baseline update` plan prints is printed by the audited head's in
   the same order apart from the new section and, when the catalog moved,
   the lines derived from the catalog change (`Current catalog:`, the
   `File changes:` count and its planned-file lines, `Plan Digest:` and
   `Next action:`); that every key of the Delivery Base's JSON result is in
   the audited head's with the same value, except `priorCatalog`,
   `currentCatalog`, `fileChanges`, `planDigest` and the `nextAction`
   string that embeds that digest when the catalog moved; and that both
   exit with the same code. On each binary, the exit code with the passage
   equals the one without it. The Plan Digest binds the source file's bytes
   (ADR-0071), so removing the passage changes it; instead, on the audited
   head, two runs over identical repository bytes print the same Plan
   Digest, and `misplacedRepositoryRules` is absent from the plan document
   the digest hashes.
7. MUST record the outside-evidence row of `_prd.md` → Success Metric 5 and
   Acceptance evidence. For each adopted checkout that
   [`references/2026-11-16-adopter-rule-placement-inventory.md`](references/2026-11-16-adopter-rule-placement-inventory.md)
   names and the QA host holds, record its commit and its
   `git status --porcelain`, run
   `bin/roundfix baseline update --repo <checkout> --no-skills --no-history --format json`
   without confirmation, record `misplacedRepositoryRules`, and record the
   same `git status --porcelain` afterwards. When the commit equals the
   inventory's, the entries must equal the inventory's table for that
   repository, path, line range and first line; argus and this repository
   report none, and the Turborepo block of fiscus `AGENTS.md` is not
   named. When the commit moved, compare with a fresh read made by the
   inventory's method and record the method's output. Name each checkout
   the host lacks. This row never passes `--yes`, `--confirm-plan` or
   `--promote`, and never writes in another repository.
8. MUST run each sabotage below on its own temporary copy of the audited
   head made with `git archive` and committed in a fresh Git repository
   (the suite guard refuses a tree outside Git), with `-timeout` on every
   `go test`, and record that the named test fails on the named subtest. No
   sabotage regenerates or synchronizes anything:
   1. The detector keeps a foreign block in its passage:
      `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/between_root_managed_regions_beside_a_foreign_block`.
   2. The detector does not trim blank lines:
      `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/after_a_guide_managed_region`.
   3. The detector follows a symbolic-link `CLAUDE.md`:
      `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/symbolic_link_claude_md`.
   4. A `marker-absent` region becomes an entry:
      `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/marker_absent_region_is_not_an_entry`.
   5. The applied result keeps the line numbers of the plan:
      `TestBaselineUpdateNamesMisplacedRepositoryRules/applied_json`.
   6. The Doctor line returns `failed` for a passage:
      `TestDoctorReportsMisplacedRepositoryRules/warn_one_passage`.
   7. The old sentence "Preserve repository-authored extensions outside
      setup markers." is restored in `context-workflow.json`:
      `TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide/docs_one_job_per_directory`.
   8. The starter template loses its paragraph:
      `TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide/starter`.
9. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
10. MUST search and record:
    - the user guides of `baseline` and `doctor`, both copies of the
      Roundfix `baseline` and `setup` references and both copies of the
      `setup-context-driven` skill carry the text task_03 prescribes, and
      each copy under `skills/` equals its `.agents/skills/` source;
    - `internal/baseline/misplaced_rules.go` calls no function that writes,
      creates, renames, changes the mode of or removes a file, and runs no
      command;
    - the Plan Document type gained no field from this Spec.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD declares `adds:` for **Misplaced Repository Rule**, that task_02
    wrote that entry and no other, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `internal/cli/cli_test.go`,
      `.roundfixrc.yml`, `docs/agents/specific-repository.md`, any Vendored
      Skill under `.agents/skills`, any Source Baseline under
      `internal/baseline/assets/source-baselines/`, any Setup Snapshot under
      `internal/baseline/assets/setups/`, any profile other than
      `standard-typescript-monorepo.json`, any ADR, and any owned skill
      other than `roundfix` and `setup-context-driven`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix` or the real home directory, apart from
      reading the Rust toolchain of requirement 4;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - write in an adopted checkout, or pass `--yes`, `--confirm-plan` or
      `--promote` there;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the three Tasks and the packages.
- [ ] Reproduce the Surface Transcript and the product rows, and compare with the Delivery Base.
- [ ] Re-read the adopted checkouts read-only and compare with the inventory.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The rendered guides, the root pointer and a new `docs/agents/specific-repository.md` send every repository rule there, with the example, and this repository is current.
- [ ] `baseline update` and `doctor` name each Misplaced Repository Rule with its destination, move nothing, and keep the Plan Digest and exit codes of the Delivery Base.
- [ ] The adopted checkouts report exactly the inventory's passages, the Turborepo block and the clean repositories report none, and no checkout changes.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0399-adopters-keep-their-own-rules-in-the-specific-repository-guide/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-3; Surface Transcript 1; Interfaces; Testing Approach; Build Order 4
- [`references/2026-11-16-adopter-rule-placement-inventory.md`](references/2026-11-16-adopter-rule-placement-inventory.md)
- ADR-0411, ADR-0145, ADR-0203, ADR-0220, ADR-0186, ADR-0257
