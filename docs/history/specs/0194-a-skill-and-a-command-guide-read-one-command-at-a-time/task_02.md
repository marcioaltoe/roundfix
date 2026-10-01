---
task: task_02
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: completed
type: backend
complexity: high
---

# Task 02: The Roundfix skill is an entry file plus one reference per command family

## Overview

`.agents/skills/roundfix/SKILL.md` is one file that every Agent loads whole and that almost every CLI Task edits. This Task splits it as the TechSpec's "The skill layout" section states: `SKILL.md` keeps what every reader needs plus an index, and every other section moves, byte for byte, into one of eighteen files under `references/`. It splits whatever content the file holds on this Task's starting commit, which may differ from the content this Spec was written against. The Verification proves, against that starting commit, that no line was lost, duplicated or changed.

## Requirements

1. MUST cut the body of `.agents/skills/roundfix/SKILL.md`, after its front matter, at every line that starts with `## ` outside a fenced code block, and move each section whole, with its heading, to the file that the TechSpec's section table names. A file that receives several sections keeps them in their original order.
2. MUST place a section the table does not name by the TechSpec's three-step rule and its command table, and record each such placement in its Result.
3. MUST keep in `SKILL.md`, in the TechSpec's order: the front matter, the `# Roundfix` introduction, the reference index, the `### QA settlement` block, and the sections `Context-Efficient Evidence Boundaries`, `Assigned Review Issue Batches`, `Assigned Task Batches`, `Forbidden Actions` and `Completion Report`.
4. MUST move the `### QA settlement` block, from its heading to the next heading, out of whatever section holds it and into `SKILL.md`, with every byte of the block unchanged.
5. MUST write the reference index between the lines `<!-- roundfix:reference-index:begin -->` and `<!-- roundfix:reference-index:end -->`. It holds one table row per reference file: a relative link such as `references/deliver.md`, the commands the file covers, and when to read it. It also states that a change to a command edits that command's reference file. The index is the only new text this Task writes in the skill.
6. MUST NOT change, add, remove or re-wrap any other line of the skill. It MUST perform the move with a throwaway script or tool that copies lines, never by retyping. That script is not part of the change and stays out of the repository.
7. MUST keep `SKILL.md` at or under 20,000 bytes. The layout test reads that budget from one named constant.
8. MUST raise the skill's version in both front-matter fields by one minor step from the value on its starting commit, resetting the patch component to zero, and follow whatever version rule is in force on that commit. When a test in `skills/skills_test.go` finds the Roundfix skill's version line by the floor's literal and no longer finds it, it MUST change that test to find the line by the declared version, and name it in its Result. On `9e439dbb`, `TestOwnedSkillContractRejectsSetAndVersionDisagreement` and `TestOwnedSkillBundleReadinessKeepsStatesDistinct` do this and fail once the version moves.
9. MUST leave `agents/openai.yaml` unchanged, then run the sanctioned `make skills-sync` and `make baseline-digests`, and name in its Result every file either command rewrote.
10. MUST put the new tests in `skills/roundfix_layout_test.go`. They read the canonical skill under `.agents/skills/roundfix` and the embedded bundle, and write only to temporary directories.
11. MUST edit no governed file other than the two `SKILL.md` files and the eighteen canonical reference files, rename or remove no top-level test, and change no exported function signature.
12. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`.

## Subtasks

- [ ] Move every section to its reference file with a script, and keep the entry sections.
- [ ] Write the reference index and raise the version.
- [ ] Add the layout tests.
- [ ] Regenerate the mirror and the digests with the sanctioned commands.

## Acceptance Criteria

- [ ] Every non-blank line of the skill's body on the starting commit is present, with the same count, in `SKILL.md` and the eighteen reference files together, outside the index block.
- [ ] The index names every file under `references/`, names no file that does not exist, and every index link resolves.
- [ ] `SKILL.md` is at most 20,000 bytes, and its `### QA settlement` section is byte-identical to the one in the `qa-gate` and `archive-spec` skills.
- [ ] `roundfix skills check` reports no diagnostic, and an install writes every reference file next to `SKILL.md`.
- [ ] Every contract that pins the skill's text passes while reading the entry file with its references.
- [ ] The mirror under `skills/roundfix` is byte-identical to the canonical skill.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/archive.md`
- creates: `.agents/skills/roundfix/references/baseline.md`
- creates: `.agents/skills/roundfix/references/deliver.md`
- creates: `.agents/skills/roundfix/references/events.md`
- creates: `.agents/skills/roundfix/references/implement.md`
- creates: `.agents/skills/roundfix/references/profiles.md`
- creates: `.agents/skills/roundfix/references/reconcile.md`
- creates: `.agents/skills/roundfix/references/release.md`
- creates: `.agents/skills/roundfix/references/review.md`
- creates: `.agents/skills/roundfix/references/review-runs.md`
- creates: `.agents/skills/roundfix/references/runs.md`
- creates: `.agents/skills/roundfix/references/runtime.md`
- creates: `.agents/skills/roundfix/references/settle.md`
- creates: `.agents/skills/roundfix/references/setup.md`
- creates: `.agents/skills/roundfix/references/spec.md`
- creates: `.agents/skills/roundfix/references/spec-delivery.md`
- creates: `.agents/skills/roundfix/references/stop.md`
- creates: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `skills/roundfix/references/archive.md`
- creates: `skills/roundfix/references/baseline.md`
- creates: `skills/roundfix/references/deliver.md`
- creates: `skills/roundfix/references/events.md`
- creates: `skills/roundfix/references/implement.md`
- creates: `skills/roundfix/references/profiles.md`
- creates: `skills/roundfix/references/reconcile.md`
- creates: `skills/roundfix/references/release.md`
- creates: `skills/roundfix/references/review.md`
- creates: `skills/roundfix/references/review-runs.md`
- creates: `skills/roundfix/references/runs.md`
- creates: `skills/roundfix/references/runtime.md`
- creates: `skills/roundfix/references/settle.md`
- creates: `skills/roundfix/references/setup.md`
- creates: `skills/roundfix/references/spec.md`
- creates: `skills/roundfix/references/spec-delivery.md`
- creates: `skills/roundfix/references/stop.md`
- creates: `skills/roundfix/references/storage.md`
- creates: `skills/roundfix_layout_test.go`
- interface: `skills/skills_test.go`
- instruction: `skills/skills.go`
- instruction: `skills/settlement_guidance_repocontract_test.go`
- instruction: `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRoundfixSkillIndexNamesEveryReference|TestRoundfixSkillIndexLinksResolve|TestRoundfixSkillEntryFileStaysWithinItsBudget|TestInstallWritesEveryRoundfixReference|TestSettlementGuidanceIsOneTable|TestCheckValidatesRoundfixSkillArtifacts|TestNoPythonBaselineRuntime|TestOwnedSkillContractRejectsSetAndVersionDisagreement|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestBaselineExamplesParse)$" ./skills ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRoundfixSkillIndexNamesEveryReference TestRoundfixSkillIndexLinksResolve TestRoundfixSkillEntryFileStaysWithinItsBudget TestInstallWritesEveryRoundfixReference TestSettlementGuidanceIsOneTable TestCheckValidatesRoundfixSkillArtifacts TestNoPythonBaselineRuntime TestOwnedSkillContractRejectsSetAndVersionDisagreement TestOwnedSkillBundleReadinessKeepsStatesDistinct TestBaselineExamplesParse; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- .agents/skills/roundfix/references/deliver.md)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'NR==1 && $0=="---" {fm=1; next} fm && $0=="---" {fm=0; next} !fm && NF {print}' "$tmp/before.md" | LC_ALL=C sort > "$tmp/old"; awk 'FNR==1 {fm=0} FNR==1 && $0=="---" {fm=1; next} fm && $0=="---" {fm=0; next} fm {next} $0=="<!-- roundfix:reference-index:begin -->" {skip=1; next} $0=="<!-- roundfix:reference-index:end -->" {skip=0; next} !skip && NF {print}' .agents/skills/roundfix/SKILL.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/deliver.md .agents/skills/roundfix/references/events.md .agents/skills/roundfix/references/implement.md .agents/skills/roundfix/references/profiles.md .agents/skills/roundfix/references/reconcile.md .agents/skills/roundfix/references/release.md .agents/skills/roundfix/references/review.md .agents/skills/roundfix/references/review-runs.md .agents/skills/roundfix/references/runs.md .agents/skills/roundfix/references/runtime.md .agents/skills/roundfix/references/settle.md .agents/skills/roundfix/references/setup.md .agents/skills/roundfix/references/spec.md .agents/skills/roundfix/references/spec-delivery.md .agents/skills/roundfix/references/stop.md .agents/skills/roundfix/references/storage.md | LC_ALL=C sort > "$tmp/new"; cmp "$tmp/old" "$tmp/new" || { printf 'the split lost, duplicated or changed a line\n' >&2; exit 1; }; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; make skills-sync-check` — expected: exit 0; before this Task the four new named tests do not exist, so the command fails.
- `test -f .agents/skills/roundfix/references/profiles.md || { printf 'the references are missing\n' >&2; exit 1; }; docs="$(go test -count=1 -tags docscontract -v -run "^(TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$docs"; exit 1; }; for name in TestBaselineDocumentationContract TestProfilesDocumentationContractMatchesPublicGuidance TestReleasePlanDocumentationContract; do printf "%s\\n" "$docs" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no reference file exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 3; Core Feature 1; Success Metric 1, Success Metric 2, Success Metric 3
- [_techspec.md](_techspec.md) — The skill layout; The move proof; API Contract 1; API Contract 2; Testing Approach 3; Build Order 2
- ADR-0187

## Result

Split the starting commit's skill with `/tmp/roundfix-task02-split.py`, a
throwaway line-copying script outside the repository. The entry keeps the
introduction, index, unchanged QA settlement block, and five shared Batch
sections in the specified order. Eighteen reference files retain every moved
section's bytes and original order within each file. Both version fields move
from `0.0.4` to `0.1.0`; both `agents/openai.yaml` files remain unchanged.

Placements decided by the three-step rule:

- `Config initialization` → `references/setup.md`, from its first shown
  command, `roundfix init`.
- `Run Window` → `references/deliver.md`, from the heading's `window` command.
- `QA Report acceptance` → `references/settle.md`, from its first shown
  command, `roundfix qa-report`.

Added `skills/roundfix_layout_test.go`. Its four named layout tests exercise
both the canonical tree and embedded bundle, require eighteen indexed
references with exactly one row each, resolve every indexed link, read the
20,000-byte budget from `roundfixSkillEntryByteBudget`, and compare the real
installer's temporary output to canonical bytes.

The two version-mutation tests,
`TestOwnedSkillContractRejectsSetAndVersionDisagreement` and
`TestOwnedSkillBundleReadinessKeepsStatesDistinct`, already locate the version
by the declared bundle version on this starting commit. Their shared
`skillVersionBelow` helper in `skills/skills_test.go` assumed a nonzero patch;
it now finds a lower valid version when a minor-step version resets the patch
to zero. The sanctioned version-recording test adds `0.1.0` to
`skills/testdata/owned-skill-versions.json`; this ordinary file is declared as
additional version-rule fallout of this Task.

Regeneration:

- `make skills-sync` exited 0. Its byte-changed outputs are
  `skills/roundfix/SKILL.md` and these eighteen files under
  `skills/roundfix/references/`: `archive.md`, `baseline.md`, `deliver.md`,
  `events.md`, `implement.md`, `profiles.md`, `reconcile.md`, `release.md`,
  `review.md`, `review-runs.md`, `runs.md`, `runtime.md`, `settle.md`,
  `setup.md`, `spec.md`, `spec-delivery.md`, `stop.md`, and `storage.md`.
  The command also recopies other owned skill mirrors; their bytes did not
  change.
- `GOCACHE=/tmp/roundfix-task02-gocache make baseline-digests` exited 0,
  reported `changed:false`, and rewrote no derived pin.
- `GOCACHE=/tmp/roundfix-task02-gocache go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` exited 0
  and updated only the owned-skill version ledger named above.

Acceptance evidence from focused implementation checks:

| Criterion | Evidence |
| --- | --- |
| Original nonblank body lines retain their counts | The split script and a separate Python audit compared line Counters against `HEAD` with front matter and the new index excluded; both comparisons matched. Whole sections and the extracted QA block were copied directly. |
| Index covers every reference and links resolve | All four new layout tests failed before the split. `TestRoundfixSkillIndexNamesEveryReference` and `TestRoundfixSkillIndexLinksResolve` now pass for canonical and embedded sources. |
| Entry budget and identical QA settlement | Entry is 12,629 bytes. `TestRoundfixSkillEntryFileStaysWithinItsBudget` and `TestSettlementGuidanceIsOneTable` pass; the script also compared the untrimmed QA block byte for byte with the other two canonical skills. |
| Shipped check and installation | `go run ./cmd/roundfix skills check` exits 0 with no diagnostic. `TestCheckValidatesRoundfixSkillArtifacts` and `TestInstallWritesEveryRoundfixReference` pass; the latter installs to `t.TempDir()` and compares every reference and entry with canonical bytes. |
| Text-pinning contracts | Focused settlement, wording, version, and documentation contracts pass. The broader skills run exposes the out-of-slice reader omission described below; this criterion is not fully evidenced. |
| Canonical and mirror match | An independent recursive Python audit compared both file sets and every file's bytes, including references and the unchanged manifests; all match. |

Focused commands (RTK proxy was used to preserve output):

- `GOCACHE=/tmp/roundfix-task02-gocache go test ./skills -run
  'TestRoundfixSkill|TestInstallWritesEveryRoundfixReference' -count=1`:
  initial red signal, all four new tests failed on missing references/index
  and the 148,315-byte entry.
- `GOCACHE=/tmp/roundfix-task02-gocache go test ./skills -run
  'TestRoundfixSkill|TestInstallWritesEveryRoundfixReference|TestSettlementGuidance|TestCheckValidatesRoundfix|TestOwnedSkillContract|TestOwnedSkillBundle|TestEveryOwnedSkillVersion'
  -count=1 -v`: exited 0, all selected tests passed.
- `GOCACHE=/tmp/roundfix-task02-gocache go test -tags docscontract
  ./internal/docscontract -run
  'Test(BaselineDocumentation|ProfilesDocumentation|ReleasePlanDocumentation)'
  -count=1`: exited 0.
- `GOCACHE=/tmp/roundfix-task02-gocache go run ./cmd/roundfix skills check`:
  exited 0.
- `git diff --check`: exited 0.
- `GOCACHE=/tmp/roundfix-task02-gocache go test ./skills -count=1`:
  exited 1 on `TestReviewRequestContract`.
- `GOCACHE=/tmp/roundfix-task02-gocache make verify-incremental`:
  exited 2. `go vet` passed. The test phase reported the reader omission
  below, two process-owner stop tests denied process-table access by the
  sandbox, and a suiteguard diagnostic because the Agent edited this Result
  while the CLI tests were running. The log is
  `/tmp/roundfix-task02-incremental.log`. No test or guard was weakened.
- `GOCACHE=/tmp/roundfix-task02-gocache go test ./internal/cli -run
  '^TestRunForceStop(LegacyRunWithoutOwnerIdentityStillStopsOwner|OwnerProcessIntegrationProvesExitBeforeStoreCompletion)$'
  -count=1`, with process-table access and no concurrent edits: exited 0.
  This focused rerun clears those two environment failures and produces no
  suiteguard diagnostic; it does not turn the earlier incremental run into
  a passing gate.

Follow-up for the reader Task: `TestReviewRequestContract` calls
`testWorkflowProjectConstraintContract` in
`skills/baseline_skill_contract_test.go:1338`, which still reads only the
canonical and mirror entry files. Seven required phrases moved unchanged to
`references/review-runs.md`, so the broader skills run fails. That helper
needs to read the entry together with references, as the reader Task specifies.
This governed file is outside this Task's permitted edit set and was left
unchanged. No wording was restored or duplicated to hide the diagnostic.

The authored Verification commands were not run. Task status, checkboxes,
the Task Graph, and other Task files were left to their owners; no commit,
push, or Pull Request was created.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `skills/testdata/owned-skill-versions.json`

## Carry-forward provenance

- Source Run: `run_20261001T010622Z_8bfe4bebad66f971`
- Source commit: `7f2133f5eb425f5a0b335457f43ecf70d33001c4`
