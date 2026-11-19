---
task: task_05
spec: 0410-the-baseline-owns-the-documentation-rules-and-the-secondbrain-export
status: pending
type: docs
complexity: high
---

# Task 05: This repository marks its decisions and publishes only its README, its user guide and its lessons

## Overview

This repository applies the rules to itself. The Agent marks every ADR in the tree with the Knowledge Mark by the test of the export rule, writes the dated notices that ADR-0438 and ADR-0440 give ADR-0248, ADR-0278 and ADR-0412, and runs the Managed Refresh so the update fills their Amendment Lists and names the ADRs marked true in the Standard Export List. The export file ends with this repository's own lines, which drop the glossary, the references, the design notes and every ADR, so that only the README, the user guide and the lessons are published; its contract test follows the new rule, and the coverage record is recorded again.

## Requirements

1. MUST give every file under `docs/adr/` in the tree when the Task runs, the ADRs of other Specs that merged before it included, the Knowledge Mark `knowledge: true` or `knowledge: false` on the line right after `status` in its header, by the test of `_techspec.md` → API Contract 24 and `clause.secondbrain.export-03-knowledge-mark`: for this repository the product is Roundfix, so a decision about what a command, the Daemon, the Delivery Queue, the Baseline or a skill does for the people and repositories that use Roundfix is `true`, and a decision about this repository's tests, build, CI, release procedure or internal structure is `false`. The mark is a judgment made one ADR at a time; the Task adds exactly one line to each file for it, never reflows or re-saves a file, never moves `updated_at`, and changes no other line (`_techspec.md` → Risks & Considerations). The ADR files it marks, other than the three of requirement 2, are recorded by the Daemon.
2. MUST write the three dated notices of `_techspec.md` → API Contract 23, word for word, each after the text it changes: ADR-0248's after its paragraph that ends "A documentation contract holds both directions.", ADR-0278's after its paragraph that begins "**Lessons leave before the text.**", and ADR-0412's at the end of its decision text. No other line of the three changes by hand, and the decision text of no accepted ADR changes.
3. MUST replace every line of `.secondbrain-export` outside the Standard Export List with the own lines of API Contract 22: after the end marker, one blank line and exactly these seven lines, in this order, so that `CHANGELOG.md` and the lines the conversion kept leave the file (`_prd.md` → Decisions, "What does this repository's own variation drop"):
   - `# This repository publishes its README, its user guide and its lessons (maintainer, 2026-11-18).`
   - `# The glossary, the references and the design notes describe how the tool is built.`
   - `!CONTEXT.md`
   - `!docs/references/`
   - `!docs/design/`
   - `# Its decision records are about the tool; the user guide carries what other projects use.`
   - `!docs/adr/*`
   The lines between the two markers are never edited by hand.
4. MUST then run `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` until a run reports nothing to do: the run fills the Amendment Lists of ADR-0248, ADR-0278 and ADR-0412 from the new notices and rewrites the Standard Export List so it names the ADRs marked true, and changes no line outside the markers. Where an ADR's header already holds `amended_by`, as ADR-0248's and ADR-0278's do since task_02's refresh, the run replaces that one line with the longer list that keeps every identifier it held (`_techspec.md` → ADR header and the lesson form, "writes the Amendment List by inserting or replacing one line"); that is the only line of an ADR the Task may remove.
5. MUST rewrite `internal/docscontract/secondbrain_export_test.go` as API Contract 25 says, keeping its build tag `docscontract` and its `//verify:always` line: it holds `TestThisRepositoryExportsNoProcessPath`, which fails when the file's own lines publish a process path or lack a reason, through `ExportVariationDepartures`, which task_03 adds, and `TestTheExportRuleRefusesAProcessPath`, which proves on a temporary directory, in the subtests `docs_whole` and `agents_md`, that the same rule refuses `docs/` whole and `AGENTS.md`; and it no longer holds `TestSecondbrainExportFollowsTheHistoryForm` or `TestSecondbrainExportRuleRefusesBothMismatches`. Neither new test asserts that every ADR carries the mark or that the Standard Export List is current (`_prd.md` → Decisions).
6. MUST record the coverage record again, never by hand, with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`, because the two replaced tests are named in it.
7. MUST NOT change any Go file other than `internal/docscontract/secondbrain_export_test.go`, any skill, any guide under `docs/agents/`, any Baseline module, template, profile, golden or Source Baseline, any file under `docs/user-guide/`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `AGENTS.md`, `CHANGELOG.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration, and creates no file under `docs/lessons/`. Of the Governed Paths it changes only `docs/references/coverage-record.json`, listed in `_authorization.md` `paths`.

## Subtasks

- [ ] Mark every ADR in the tree, one line each.
- [ ] Write the three amendment notices.
- [ ] Write this repository's own lines of the export file and run the Managed Refresh until it reports nothing to do.
- [ ] Rewrite the export contract test and record the coverage record again.

## Acceptance Criteria

- [ ] Every ADR in the tree carries the Knowledge Mark on the line after `status`, and no ADR lost a line other than an `amended_by` line the refresh replaced with a list that keeps every identifier it held.
- [ ] ADR-0248, ADR-0278 and ADR-0412 carry their dated notices and list the ADRs that amend them in their headers.
- [ ] This repository's update reports `current` with a current documentation part and no unmarked ADR, and its export file is the Standard Export List followed by exactly the own lines of API Contract 22.
- [ ] The export contract tests hold that this repository exports no process path and that the rule refuses one, and the coverage record is current.

## Context

- instruction: `docs/adr/0439-an-adrs-header-says-whether-the-secondbrain-wants-it-and-which-adrs-amend-it.md`
- instruction: `docs/agents/secondbrain.md`
- interface: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- interface: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- interface: `docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md`
- interface: `.secondbrain-export`
- interface: `internal/docscontract/secondbrain_export_test.go`
- interface: `docs/references/coverage-record.json`

## Consumers

- This repository's export file and its contract tests — consumers: `./internal/docscontract`
- The coverage record — consumers: `./internal/spec`

## Verification

- `for f in docs/adr/*.md; do awk 'NR == 1 && $0 == "---" { fm = 1; next } fm && $0 == "---" { exit } fm && previous ~ /^status:/ { found = 1; ok = ($0 ~ /^knowledge: (true|false)$/) } fm { previous = $0 } END { exit((found && ok) ? 0 : 1) }' "$f" || { echo "no Knowledge Mark on the line after status: $f" >&2; exit 1; }; done; changed="$(git diff --numstat HEAD -- docs/adr | awk '$2 != 0 { print $3 }')" || exit 1; for f in $changed; do old="$(git diff -U0 HEAD -- "$f" | awk '/^--- / { next } /^-/ { print }')" || exit 1; test "$(printf '%s\n' "$old" | grep -c '^-amended_by: \[.*\]$')" -eq 1 && test "$(printf '%s\n' "$old" | wc -l | tr -d ' ')" -eq 1 || { printf 'ADR lines removed: %s\n' "$f" >&2; exit 1; }; new="$(grep '^amended_by: \[' "$f")" || { printf 'Amendment List removed: %s\n' "$f" >&2; exit 1; }; for id in $(printf '%s\n' "$old" | grep -oE 'ADR-[0-9]+'); do printf '%s\n' "$new" | grep -qE -- "$id[],]" || { printf 'Amendment List lost %s: %s\n' "$id" "$f" >&2; exit 1; }; done; done` — expected: exit 0. Before this Task no ADR carries the mark; afterwards every ADR carries it on the line after `status`, and no ADR lost a line other than an `amended_by` line the refresh replaced with a list that keeps every identifier it held.
- `for pair in 'docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md|The export no longer mirrors' 'docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md|A lesson is no longer sent as a message:' 'docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md|Extended by [ADR-0438](0438-'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; grep -q '^amended_by: \[.*ADR-0438' docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md && grep -q '^amended_by: \[.*ADR-0440' docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md && grep -q '^amended_by: \[.*ADR-0438.*ADR-0440' docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md || { echo 'an Amendment List lacks its new identifier' >&2; exit 1; }; actual="$(awk 'after { print } /^# roundfix:standard-export:end$/ { after = 1 }' .secondbrain-export)"; expected="$(printf '%s\n' '' '# This repository publishes its README, its user guide and its lessons (maintainer, 2026-11-18).' '# The glossary, the references and the design notes describe how the tool is built.' '!CONTEXT.md' '!docs/references/' '!docs/design/' '# Its decision records are about the tool; the user guide carries what other projects use.' '!docs/adr/*')"; test "$actual" = "$expected" || { printf 'unexpected own lines:\n%s\n' "$actual" >&2; exit 1; }; bin_dir="$(mktemp -d)" || exit 1; home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; for want in '"state":"current"' '"documentation":{"status":"current"' '"unmarked":0}'; do printf '%s\n' "$out" | grep -qF -- "$want" || { printf 'missing %s in:\n%s\n' "$want" "$out" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three notices and the own lines are absent and every ADR is unmarked; afterwards the notices are written and listed in the headers, the export file ends with exactly the own lines of API Contract 22, and this repository's update under a fresh home reports `current` with a current documentation part and no unmarked ADR.
- `if grep -q 'func TestSecondbrainExportFollowsTheHistoryForm' internal/docscontract/secondbrain_export_test.go; then echo 'the history-form test is still present' >&2; exit 1; fi; out="$(go test -count=1 -timeout 20m -v -tags docscontract -run '^(TestThisRepositoryExportsNoProcessPath|TestTheExportRuleRefusesAProcessPath)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThisRepositoryExportsNoProcessPath TestTheExportRuleRefusesAProcessPath/docs_whole TestTheExportRuleRefusesAProcessPath/agents_md; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the file holds the history-form tests and not the new ones; afterwards the two new tests pass, the coverage record matches the tests, and the documentation contracts pass.

## References

- `_prd.md` → Goals; User Stories 1, 5, 6 and 7; Core Features 6, 7, 8 and 19; Success Metric 6; Paths that leave a mirror
- `_techspec.md` → This repository; API Contract 22; API Contract 23; API Contract 24; API Contract 25; Risks & Considerations; Build Order 5
- ADR-0438, ADR-0439, ADR-0440, ADR-0248, ADR-0278, ADR-0412
