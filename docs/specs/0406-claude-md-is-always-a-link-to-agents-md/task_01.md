---
task: task_01
spec: 0406-claude-md-is-always-a-link-to-agents-md
status: completed
type: backend
complexity: high
---

# Task 01: The Baseline reads the shape of the root CLAUDE.md and the passages it holds

## Overview

A new read-only reader in `internal/baseline` tells whether the root `CLAUDE.md` is the Claude Instruction Link, is absent, is a regular file, is a link elsewhere, is the link checked out by a host that does not create symbolic links, or sits beside an `AGENTS.md` that is not a regular file. It also lists the passages of its own the file holds. The Misplaced Repository Rule detector stops reporting the `@AGENTS.md` import line and the bytes of an unlinked checkout. The glossary gains the term and ADR-0411 records the amendment.

## Requirements

1. MUST add `internal/baseline/claude_link.go` with the interface of `_techspec.md` → Interfaces: `ClaudeInstructionLinkPath`, `ClaudeInstructionLinkTarget`, `ClaudeInstructionLinkShape` with the six values `link`, `absent`, `file`, `link-elsewhere`, `unlinked-checkout` and `agents-not-regular`, `ClaudeInstructionLink`, its method `Convertible` and `ReadClaudeInstructionLink(ctx context.Context, root string, git GitRunner) (ClaudeInstructionLink, error)`.
2. MUST decide the shape in the order of `_techspec.md` → invariant 2, word for word:
   - `link` compares `path.Clean` of the link's target with `AGENTS.md`, so `./AGENTS.md` is `link`, and `Target` holds the target as read;
   - `unlinked-checkout` requires both bytes exactly `AGENTS.md`, with no newline, and the mode `120000` that `git ls-files -s -- CLAUDE.md` prints;
   - the Git command runs only when the bytes are exactly `AGENTS.md`;
   - an absolute link, a link that escapes the repository, a missing target, a directory or a special file returns an error naming `CLAUDE.md`.
3. MUST compute passages as `_techspec.md` → invariant 3 says:
   - with `partitionRootSource` and `misplacedUnmarkedPassages`, cutting foreign `BEGIN:`/`END:` blocks as the detector does;
   - adding, when the Setup Manifest lists `CLAUDE.md`, its Unrecorded Managed Regions with a digest mismatch, as the detector reports them;
   - leaving out a passage whose trimmed text is exactly `@AGENTS.md`;
   - returning no passage when the bytes equal those of the root `AGENTS.md`;
   - for `link-elsewhere`, reading the target's bytes and reporting at the target's path.

   `Convertible` follows invariant 4.
4. MUST NOT write any byte, take any lock or run any Git command other than the one of requirement 2 (invariant 1).
5. MUST narrow `MisplacedRepositoryRules` in `internal/baseline/misplaced_rules.go`:
   - a regular root `CLAUDE.md` whose bytes are exactly `AGENTS.md` yields no entry;
   - a passage of a regular root `CLAUDE.md` whose trimmed text is exactly `@AGENTS.md` yields no entry;
   - every other behaviour and every existing case of `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest` stays as it is.
6. MUST add `internal/baseline/claude_link_test.go` with `TestReadClaudeInstructionLinkClassifiesEveryShape`, calling `t.Parallel()` first. Each subtest calls `t.Parallel()`, builds a temporary Git repository with `t.TempDir()`, a Setup Manifest and a root `AGENTS.md`, and asserts the whole `ClaudeInstructionLink` value and `Convertible`. The subtests are:
   - `link`, `dot_slash_target`, `absent`, `empty_file`, `import_line` and `copy_of_agents`;
   - `managed_regions_only`, `file_with_a_rule` and `listed_digest_mismatch`;
   - `link_to_a_clean_file` and `link_to_a_file_with_a_rule`;
   - `unlinked_checkout`, built by committing the link and checking it out again with `git -c core.symlinks=false checkout`;
   - `agents_bytes_tracked_as_a_file`, a regular file with the bytes `AGENTS.md` that Git records as `100644`, which is `file` with no passage;
   - `agents_missing` and `agents_is_a_link`;
   - `absolute_link_is_an_error`;
   - `reads_only`, which asserts that every fixture file's bytes and mode are unchanged afterwards.

   Helper names MUST NOT reuse `assertClaudeLink` or `keepClaudeLinks`, which belong to the Claude Skill Link tests of the same package.
7. MUST add to `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest` the cases `import_line_claude_md`, a regular `CLAUDE.md` holding `@AGENTS.md` and a newline, and `unlinked_checkout_claude_md`, a regular `CLAUDE.md` holding exactly `AGENTS.md`. Each expects no entry.
8. MUST write, through `domain-modeling`, the glossary entry **Claude Instruction Link** in `CONTEXT.md`: the root `CLAUDE.md` of an adopted repository, a relative symbolic link whose target is exactly `AGENTS.md`, which the Baseline creates at adoption and restores in a Managed Refresh when nothing would be lost (ADR-0432, ADR-0433). It MUST also add to the entry **Misplaced Repository Rule** that the `@AGENTS.md` import line of a regular root `CLAUDE.md`, and the bytes of an unlinked checkout, are not one.
9. MUST append to `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md` one paragraph that begins `**Amended by [ADR-0432](0432-claude-md-is-always-a-symbolic-link-to-agents-md.md) (2026-11-18).**`. It MUST state the narrowing of requirement 5 and that every other passage is still named and never moved.
10. MUST NOT change any Governed Path, the Plan Document, a plan golden, `internal/cli` or any guide or skill.

## Subtasks

- [ ] Add the reader and its shape and passage rules.
- [ ] Narrow the detector and add its two cases.
- [ ] Add the table test of every shape.
- [ ] Write the glossary entries and the ADR-0411 amendment.

## Acceptance Criteria

- [ ] Every shape of the TechSpec is told apart, with the exact passages, and the reader leaves every fixture byte and mode as it was.
- [ ] An unlinked checkout is recognised only when Git records the link; the same bytes tracked as a file are a file with nothing to lose.
- [ ] The detector no longer names the import line or an unlinked checkout's bytes, and every earlier case still passes.
- [ ] Removing the Git check, the import-line exception or the copy rule makes a named subtest fail.

## Context

- instruction: `docs/adr/0432-claude-md-is-always-a-symbolic-link-to-agents-md.md`
- instruction: `docs/adr/0433-a-checkout-without-symbolic-links-keeps-the-link-git-records-for-claude-md.md`
- instruction: `internal/baseline/repository.go`
- instruction: `internal/baseline/preservation.go`
- creates: `internal/baseline/claude_link.go`
- creates: `internal/baseline/claude_link_test.go`
- interface: `internal/baseline/misplaced_rules.go`
- interface: `internal/baseline/misplaced_rules_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md`

## Consumers

- The Misplaced Repository Rule report of `baseline.MisplacedRepositoryRules` — consumers: `./internal/baseline`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReadClaudeInstructionLinkClassifiesEveryShape|TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReadClaudeInstructionLinkClassifiesEveryShape/link TestReadClaudeInstructionLinkClassifiesEveryShape/dot_slash_target TestReadClaudeInstructionLinkClassifiesEveryShape/absent TestReadClaudeInstructionLinkClassifiesEveryShape/empty_file TestReadClaudeInstructionLinkClassifiesEveryShape/import_line TestReadClaudeInstructionLinkClassifiesEveryShape/copy_of_agents TestReadClaudeInstructionLinkClassifiesEveryShape/managed_regions_only TestReadClaudeInstructionLinkClassifiesEveryShape/file_with_a_rule TestReadClaudeInstructionLinkClassifiesEveryShape/listed_digest_mismatch TestReadClaudeInstructionLinkClassifiesEveryShape/link_to_a_clean_file TestReadClaudeInstructionLinkClassifiesEveryShape/link_to_a_file_with_a_rule TestReadClaudeInstructionLinkClassifiesEveryShape/unlinked_checkout TestReadClaudeInstructionLinkClassifiesEveryShape/agents_bytes_tracked_as_a_file TestReadClaudeInstructionLinkClassifiesEveryShape/agents_missing TestReadClaudeInstructionLinkClassifiesEveryShape/agents_is_a_link TestReadClaudeInstructionLinkClassifiesEveryShape/absolute_link_is_an_error TestReadClaudeInstructionLinkClassifiesEveryShape/reads_only TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/import_line_claude_md TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/unlinked_checkout_claude_md TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest/regular_claude_md; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline && go test -count=1 -timeout 30m -run 'Misplaced|BaselineRules' ./internal/cli` — expected: exit 0. Before this Task the reader and its test do not exist and the two detector cases are missing, so their pass lines are missing; afterwards every shape and both cases pass, the earlier detector cases still pass, and the Baseline suite and the CLI's misplaced-rule tests pass.

- `grep -F -- '**Claude Instruction Link**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Claude Instruction Link\*\*(:| —( |$))' && grep -F -- '**Misplaced Repository Rule**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Misplaced Repository Rule\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < docs/adr/0411-a-repository-rule-outside-the-specific-repository-guide-is-named-never-moved.md | grep -qF -- 'Amended by [ADR-0432]'` — expected: exit 0. Before this Task the glossary has no **Claude Instruction Link** entry and ADR-0411 has no amendment, so the command fails; afterwards both are present.

## References

- `_prd.md` → Goals (second and fourth bullets); Core Features 1 and 3; Success Metric 3; Decisions
- `_techspec.md` → Claude Instruction Link reader; Detector; Interfaces and invariants 1 to 4; Testing Approach (Task 1); Absence Evidence; Build Order 1
- ADR-0432, ADR-0433, ADR-0411

## Result

Skills loaded: implement-task, coding-guidelines, golang-context, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate.

Implemented the read-only Claude Instruction Link reader with all six shapes, the original link target, passage locations and `Convertible`. It reuses the inventory's safe alias resolver, the detector's foreign-block passage parser and the manifest digest matcher. Only a regular `CLAUDE.md` holding exactly `AGENTS.md` consults Git, through `ls-files -s -- CLAUDE.md`. Missing, absolute, escaping, cyclic and non-regular link targets return errors naming `CLAUDE.md`; cancellation and Git errors retain their error chains.

Narrowed the detector only for the exact unlinked-checkout bytes and the standalone import passage. Shared the existing managed-passage rendering so digest-mismatch locations and first lines retain the detector's behavior. Added the two required detector cases, the glossary definition and exclusions, and the appended ADR-0411 amendment.

Focused evidence by acceptance criterion:

1. Shapes, exact passages and no mutation: `go test ./internal/baseline -run '^(TestReadClaudeInstructionLink.*|TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest|TestMisplacedFirstLineNeverPanicsOnInvalidUTF8)$' -count=1 -timeout=3m` exited 0. The shape table covers all 17 required cases and additional safety, foreign-block, duplicate-region and import/copy cases. Every shape case asserts the whole reader value and `Convertible`, checks permitted Git arguments/call count/context, and compares the complete fixture tree's bytes, modes, links and paths before and after the read. `reads_only` also uses read-only carrier and manifest modes. The companion test checks cancellation and Git-error propagation.
2. Git provenance: the same focused command exercised `unlinked_checkout`, committing a symbolic link and recreating it with `git -c core.symlinks=false checkout -- CLAUDE.md`, alongside `agents_bytes_tracked_as_a_file` (a committed regular file) and `agents_bytes_untracked`. The former is `unlinked-checkout` and non-convertible; the latter two are `file`, without passages, and convertible. The newline variant remains a passage and runs no Git command.
3. Detector compatibility: the focused command exercised every existing case of `TestMisplacedRepositoryRulesAreNamedFromTheSetupManifest` as well as `import_line_claude_md` and `unlinked_checkout_claude_md`; all passed, including the prior `regular_claude_md` case. The invalid-UTF8 first-line cases passed too.
4. Regression sensitivity: three isolated `go test -overlay <overlay.json> ./internal/baseline -run <selector> -count=1 -timeout=2m -v` probes each exited 1 with an assertion failure in the intended named case, rather than a build or environment failure:
   - Removed the mode predicate `fields[0] == "120000"`: selector `^TestReadClaudeInstructionLinkClassifiesEveryShape$/^agents_bytes_tracked_as_a_file$` failed because the reader reported `unlinked-checkout` and non-convertible instead of `file` and convertible.
   - Disabled the import exception: selector `^TestReadClaudeInstructionLinkClassifiesEveryShape$/^import_line$` failed because an `@AGENTS.md` passage appeared and conversion was withheld.
   - Disabled the byte-identical copy rule: selector `^TestReadClaudeInstructionLinkClassifiesEveryShape$/^copy_of_agents$` failed because the root instructions appeared as a passage and conversion was withheld.
   The overlays and full failure logs are in the session scratch directory `roundfix-0406-mutations-2e_5h59l` under the host temporary directory. Production source files were never mutated by these probes.

Documentation evidence: `go test -tags=docscontract ./internal/docscontract -run '^TestGlossary(PhrasesAreDefinedTerms|LinksNameAcceptedADRs|AvoidLinesArePlainLists|StatesEachRuleOnce)$' -count=1 -timeout=3m` exited 0. The first attempt without `-tags=docscontract` reported `[no test files]`; it was not counted as glossary evidence. Before implementation, the focused `copy_of_agents` check failed to build because the reader interface was absent.

Scope: changed only this Task's declared implementation, test and documentation paths plus this Result. The pre-existing `status: in_progress` remains unchanged. No Task Graph, other Task, Governed Path, Plan Document, golden, CLI, guide or skill was edited. Declared Verification, repository Verification and the complete package suites were left to the Daemon. No commit, push or pull request. No follow-up implementation was added.
