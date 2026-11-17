---
task: task_02
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
status: completed
type: backend
complexity: high
---

# Task 02: The Baseline lists each Misplaced Repository Rule from the Setup Manifest's carriers

## Overview

A read-only detector in the `baseline` package lists every **Misplaced
Repository Rule** of a repository: each passage outside a Managed Region of
a Markdown carrier its Setup Manifest lists, and each Unrecorded Managed
Region whose bytes do not match the manifest, by path, line range and first
line, with `docs/agents/specific-repository.md` as its destination. A block
another tool delimits with its own `BEGIN:`/`END:` pair is cut out. A table
test built from the shapes of the adopter inventory proves each entry
exactly. No command calls the detector yet; task 03 does. The glossary
gains the term.

## Requirements

1. MUST add `internal/baseline/misplaced_rules.go` with the declarations of
   `_techspec.md` → Interfaces, with those names, fields, JSON tags and
   signature: `SpecificRepositoryPath`, `MisplacedRepositoryRule`,
   `MisplacedRepositoryRuleReport` and `MisplacedRepositoryRules(root string)`.
2. MUST implement invariants 1 to 7 of `_techspec.md` → Interfaces:
   - it reads the Setup Manifest with the parser `classifyManagedRegions`
     uses, splits each listed `.md` carrier with `partitionRootSource` and
     takes the Unrecorded Managed Regions with reason `digest-mismatch`
     from `classifyManagedRegions`; it reuses those functions and does not
     copy them;
   - it reads a root `CLAUDE.md` only when it is a regular file, never
     follows a symbolic link, and never reads
     `docs/agents/specific-repository.md`, a nested carrier or a file the
     manifest does not list;
   - a passage is an `unmarked-span` with each foreign block cut out, then
     trimmed of leading and trailing blank lines; blank lines inside stay;
     an empty span is no passage. A foreign block runs from a line matching
     `^\s*<!--\s*BEGIN:<name>\s*-->\s*$` to the first later line with
     `END:` and the same name; an unclosed `BEGIN:` is not foreign;
   - a cut that leaves text on both sides of a foreign block gives two
     passages, one on each side;
   - `StartLine` and `EndLine` are 1-based, inclusive, and count lines of
     the bytes on disk; a `digest-mismatch` entry carries `ManagedID`, the
     marker pair's line range and the first non-blank line after the begin
     marker; a `marker-absent` region is no entry;
   - `FirstLine` is the first line trimmed of surrounding whitespace, cut
     at 120 bytes on a rune boundary with `...` appended when cut;
     `MoveTo` is always `SpecificRepositoryPath`;
   - entries are sorted by path, then start line;
   - `ExtensionEnabled` is true only when the manifest's decision
     `repository.extension.enabled` is `true`;
   - without a manifest it returns `ManifestFound` false and no error; an
     invalid manifest returns an error that wraps the code
     `baseline.preservation.managed-refresh.manifest.invalid`; an
     unreadable listed carrier returns an error that names its path.
3. MUST write, create, rename, change the mode of or remove no file, and
   invoke no semantic classification, Git command or network call.
4. MUST add `internal/baseline/misplaced_rules_test.go` with
   `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest`, calling
   `t.Parallel()` first, one subtest per case below. Each case builds its
   repository in `t.TempDir()` with a Setup Manifest and its carriers,
   through the helpers of `preservation_test.go` such as
   `writeManagedRefreshManifest` and `managedRegionFixture`, in neutral
   text that names no adopter, and asserts the whole
   `MisplacedRepositoryRuleReport` or error, never a substring of it:
   - `after_a_guide_managed_region`: a heading and two rules after a
     guide's region give one entry from the heading to the last rule;
   - `before_a_guide_managed_region`: a heading and three lines before the
     first region give one entry starting at line 1;
   - `between_root_managed_regions_beside_a_foreign_block`: in `AGENTS.md`,
     a section between two regions followed by a closed
     `<!-- BEGIN:tool-rules -->` block gives one entry for the section and
     none for the block;
   - `several_sections_after_one_managed_region`: two headed sections, each
     under its own HTML comment, after one region give one entry;
   - `digest_mismatch_region`: a region whose body differs from the
     recorded digest gives one entry with its `ManagedID` and the marker
     pair's lines;
   - `marker_absent_region_is_not_an_entry`;
   - `unclosed_begin_stays_in_the_passage`;
   - `first_line_is_cut_at_120_bytes`: a first line of 130 bytes with a
     multi-byte rune across byte 120 is cut before that rune and ends with
     `...`;
   - `clean_carrier`: only regions and blank lines give no entry;
   - `symbolic_link_claude_md`: a `CLAUDE.md` symbolic link to an
     `AGENTS.md` with a passage gives that passage once, under `AGENTS.md`;
   - `regular_claude_md`: a regular `CLAUDE.md` with a passage gives one
     entry under `CLAUDE.md`;
   - `specific_repository_is_not_read`: text in
     `docs/agents/specific-repository.md` gives no entry;
   - `no_manifest`; `invalid_manifest`;
   - `extension_decision_off` and `extension_decision_on`;
   - `reads_only`: the paths, bytes and modes of every file under the
     fixture are the same before and after a call.

   Removing the foreign-block cut, the trim or the symbolic-link skip each
   makes a named subtest fail.
5. MUST add to `CONTEXT.md` the glossary entry **Misplaced Repository
   Rule**, through `domain-modeling`, beside **Unrecorded Managed Region**,
   in the form of its neighbours: a passage of an instruction carrier the
   Setup Manifest lists, outside every Managed Region and not inside a
   block another tool delimits, or an Unrecorded Managed Region, which
   belongs in `docs/agents/specific-repository.md`; Roundfix names it and
   never moves it (ADR-0411); with an `_Avoid_` line. MUST change no other
   glossary entry.
6. MUST keep the package building for every platform:
   `GOOS=windows go build -buildvcs=false ./internal/baseline`.
7. MUST change no file outside the three this Task declares and no
   Governed Path.

## Subtasks

- [ ] Add the report types and the detector over the manifest's carriers.
- [ ] Cut foreign blocks, trim passages, convert offsets to lines and shorten the first line.
- [ ] Add the table test of the inventory's shapes with the read-only check.
- [ ] Add the glossary entry.

## Acceptance Criteria

- [ ] For every shape of the adopter inventory the detector returns exactly the expected path, line range and first line, and for a clean carrier, a symbolic-link `CLAUDE.md`, a foreign block and `docs/agents/specific-repository.md` it returns nothing.
- [ ] An Unrecorded Managed Region with a digest mismatch is one entry with its managed identity; a region without markers is none.
- [ ] Without a manifest the report says so and carries no error; an invalid manifest is an error with the manifest code.
- [ ] No file under the fixture changes, and the package builds for Windows.
- [ ] `CONTEXT.md` defines **Misplaced Repository Rule**.

## Context

- instruction: `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md`
- instruction: `internal/baseline/preservation.go`
- instruction: `internal/baseline/preservation_test.go`
- creates: `internal/baseline/misplaced_rules.go`
- creates: `internal/baseline/misplaced_rules_test.go`
- interface: `CONTEXT.md`

## Consumers

None. The Task adds an exported function and types that no code calls yet; task 03 adds their callers and runs the suites that consume them.

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in after_a_guide_managed_region before_a_guide_managed_region between_root_managed_regions_beside_a_foreign_block several_sections_after_one_managed_region digest_mismatch_region marker_absent_region_is_not_an_entry unclosed_begin_stays_in_the_passage first_line_is_cut_at_120_bytes clean_carrier symbolic_link_claude_md regular_claude_md specific_repository_is_not_read no_manifest invalid_manifest extension_decision_off extension_decision_on reads_only; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; if grep -qE 'os[.](WriteFile|Create|CreateTemp|Remove|RemoveAll|Mkdir|MkdirAll|Rename|Chmod|Symlink|OpenFile)[(]|exec[.]Command' internal/baseline/misplaced_rules.go; then echo 'the detector writes or runs a command' >&2; exit 1; fi; GOOS=windows go build -buildvcs=false ./internal/baseline && go test -count=1 -timeout 30m ./internal/baseline` — expected: exit 0. Before this Task the test does not exist, so the command fails; afterwards every shape of the inventory, the foreign block, the symbolic link, the cut first line and the read-only check pass by name, the detector calls no writing function, the package builds for Windows and the whole Baseline suite passes.

- `grep -F -- '**Misplaced Repository Rule**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Misplaced Repository Rule\*\*(:| —( |$))'` — expected: exit 0. Before this Task `CONTEXT.md` has no such entry; afterwards it defines the term.

## References

- `_prd.md` → Goals (second and fourth bullets); Core Feature 2; Success Metric 2; Non-Goals; Glossary
- `_techspec.md` → Misplaced Repository Rule detector; Interfaces (invariants 1 to 7); Testing Approach (Task 2); Build Order 2
- [`references/2026-11-16-adopter-rule-placement-inventory.md`](references/2026-11-16-adopter-rule-placement-inventory.md), sections Method and Passages outside every Managed Region
- ADR-0411, ADR-0145, ADR-0203

## Result

Skills loaded: implement-task, coding-guidelines, domain-modeling, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Added the read-only detector, exact-report table tests, and the glossary entry beside Unrecorded Managed Region. No command calls the new exported API yet. Changes are limited to `internal/baseline/misplaced_rules.go`, `internal/baseline/misplaced_rules_test.go`, `CONTEXT.md`, and this Result. The pre-existing Daemon status transition is preserved; no Task status, Task Graph, other Task, or Governed Path was edited.

The verification feedback artifact was inspected. It identified the absent required table test; no diagnostic log body is copied here. The detector and all 17 required named cases now exist, with 8 additional cases for foreign-block splitting and name matching, excluded and unlisted carriers, missing carriers, sorting, CRLF, and a non-boolean extension decision.

Implementation detail for the earlier helper-boundary conflict: `classifyManagedRegions` is reused when every manifest artifact is an eligible carrier. It has no filtering parameter and otherwise would read excluded files or reject a symbolic-link `CLAUDE.md`. When its carrier set includes exclusions, the detector instead reuses `managedEntryMatchesManifest` on eligible partitions, preserving the classifier's duplicate-ID exclusion. The manifest parser, partitioner, and digest matcher are not copied, and `preservation.go` is unchanged. Thus classifier reuse is conditional rather than universal; the excluded-carrier case includes a digest mismatch to exercise this path. Root `CLAUDE.md` is read only when manifest-listed and regular, following this Task's explicit manifest-only boundary; the TechSpec and ADR describe a broader exception for an unlisted regular root file.

Acceptance evidence:

- Inventory shapes: whole-report assertions check the exact path, inclusive line range, trimmed first line, managed identity, destination, and report flags. The foreign-block cases assert both exclusion and separate passages on either side. UTF-8 truncation uses a 130-byte line with a three-byte rune crossing byte 120. Clean, unlisted, nested, non-Markdown, destination-guide, and symbolic-link cases produce exactly the expected absence of entries.
- Unrecorded regions: `digest_mismatch_region` asserts the marker pair's lines and `ManagedID`; `marker_absent_region_is_not_an_entry` asserts no entry. `excluded_carriers_are_not_opened` lists absent excluded files alongside an eligible digest mismatch, proving those paths are not opened by the detector or by an unconditional classifier call.
- Manifest and errors: `no_manifest` asserts the zero report and no error; `invalid_manifest` asserts the entire wrapped manifest-code error and report. `missing_listed_carrier` asserts the entire path-bearing error. Decision cases cover absent, false, true, and string-valued decisions.
- Read-only behavior: every case compares `snapshotInspectionTree` before and after the call, including paths, bytes, modes, and link targets. `reads_only` also gives the manifest and carrier files read-only modes. The detector contains no write, mutation, subprocess, Git, network, or semantic-classification operation.
- Windows: the baseline package and its tests cross-compiled for `windows/amd64` using the focused compile-only command below. The Task's declared Windows build remains for the Daemon.
- Glossary: only the Misplaced Repository Rule entry was added to `CONTEXT.md`, including the ADR-0411 link and `_Avoid_` line. Existing glossary contracts passed with their required build tag.

Focused checks, all exit 0:

- `go test -count=1 -timeout 3m ./internal/baseline -run '^TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/(after_|before_|between_|several_|digest_|marker_|unclosed_|first_|clean_|foreign_|crlf_)'`
- `go test -count=1 -timeout 3m ./internal/baseline -run '^TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/(symbolic_|regular_|specific_|no_manifest|invalid_|extension_|reads_only|excluded_|unlisted_|missing_|paths_)'`
- `go test -tags docscontract -count=1 -timeout 3m ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|AvoidLinesArePlainLists|StatesEachRuleOnce)$'`. An earlier invocation without the build tag selected no tests; the tagged run above supplied the glossary evidence.
- `env GOOS=windows GOARCH=amd64 go test -c -o /private/tmp/roundfix-0399-task02-windows.test.exe ./internal/baseline`
- `git diff --check` and changed-path inspection.

Negative probes used temporary Go overlays, without editing repository source, running `go test -json -overlay <temporary-overlay> -count=1 -timeout 2m ./internal/baseline -run '^TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/<named-subtest>$'`. Each returned exit 1 with the intended named subtest failing, and each actual assertion was inspected:

- Removing foreign-block cutting made `between_root_managed_regions_beside_a_foreign_block` include foreign lines in its reported end line.
- Removing blank-line trimming made `after_a_guide_managed_region` report blank outer lines and an empty first line.
- Removing the root alias skip made `symbolic_link_claude_md` return a regular-file error instead of its expected report.

No declared Verification command, repository Verification, or full test suite was run. No commit, push, or pull request was created. Daemon Verification and settlement remain outstanding.

## Carry-forward provenance

- Source Run: `run_20261117T004735Z_ff97830eca75451a`
- Source commit: `e929c1b2091c45462e8aac16fd7914690bf193f4`
