---
task: task_08
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
type: backend
complexity: high
---

# Task 08: Guarantee Formatter-Stable Output

## Overview

Make the final generated managed corpus stable under each profile's declared
Markdown formatter contract. The ordinary gate must prove the pinned contract
hermetically, while final QA retains a reproducible real-formatter probe.

## Requirements

1. MUST declare formatter compatibility for every bundled profile and pin the
   TypeScript/Bun profile to the reviewed Oxfmt contract.
2. MUST use one formatter-stable managed-block and shared-guide framing
   convention across all generated Markdown paths.
3. MUST include a checked-in golden corpus and provenance digest produced by
   the exact declared formatter version.
4. MUST prove confirmed apply, formatter-canonical comparison, selected
   fixture Verification, fresh audit, and second apply leave no diff and no
   Change Plan.
5. MUST keep ordinary Verification network-free and free of package
   installation or ambient formatter dependencies.
6. MUST provide the exact pinned real-formatter probe for final QA without
   running a formatter as an audit or apply mutation.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [x] Declare profile formatter contracts and pinned provenance.
- [x] Correct managed root and shared-guide framing.
- [x] Build the hermetic TypeScript/Bun golden corpus.
- [x] Add the full apply-to-reapply composition fixture.
- [x] Record the exact real Oxfmt QA probe.
- [x] Exercise formatter contracts across every supported profile.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] Every bundled profile explicitly selects a formatter contract or records
      that no Markdown formatter is selected.
- [x] Generated TypeScript/Bun Markdown is byte-identical to the pinned Oxfmt
      golden corpus.
- [x] The composition fixture finishes with an empty repository diff and empty
      `plannedChanges` after formatter comparison, Verification, audit, and
      reapply.
- [x] Changing generated framing or golden bytes without updating valid
      provenance fails the fixture.
- [x] Ordinary formatter tests run successfully without network, package
      installation, or an installed Oxfmt executable.
- [x] Audit and apply execute no formatter process and retain existing exit and
      output contracts.
- [x] Every profile's apply, audit, and reapply macro flow remains clean.
- [x] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `docs/adr/0059-generated-output-is-formatter-stable-in-the-target-repository.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/assets/templates/index.json`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `Makefile`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_formatter_compatibility.py'` — expected: the hermetic formatter composition ends with byte-identical golden output, a clean audit, and no reapply changes.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py` — expected: every bundled formatter contract composes with apply, audit, and reapply.
- `rtk make verify` — expected: the full repository gate passes without network or ambient formatter setup.

## References

- `_prd.md` → Goal 4; User Story 5; Core Feature 10; User Experience; Success
  Metrics.
- `_techspec.md` → System Architecture; Integration Points; Testing Approach;
  Risks & Considerations; Build Order 5.
- ADR-0059 → formatter compatibility as generated-output behavior.

## Result

Implemented Formatter-Stable Output as a profile contract and generated-byte
invariant. All bundled profiles now use profile schema v3: TypeScript/Bun pins
Oxfmt 0.59.0, while Go CLI/TUI and Rust CLI explicitly select no Markdown
formatter. Managed markers and rendered clause lists use one formatter-stable
blank-line convention.

The TypeScript/Bun fixture contains the maximal generated Markdown corpus and
a portable-file digest bound to checked-in provenance. The provenance records
the final-QA probe exactly as
`rtk bunx oxfmt@0.59.0 --check AGENTS.md docs/agents`; ordinary Verification
does not execute or install Oxfmt.

Acceptance evidence:

1. The asset contract test matched all three bundled formatter declarations.
   `test_selected_contract_binds_exact_oxfmt_provenance_and_golden_digest`
   matched Oxfmt 0.59.0, the exact QA probe, every corpus path, and digest
   `f001aebf81530abb9d5069145db4fe3f3c562306d7030503d65b87687dca5fbb`.
2. The focused formatter suite generated the maximal TypeScript/Bun corpus and
   compared every Markdown byte with the pinned golden corpus.
3. The composition test completed confirmed apply, hermetic formatter
   comparison, selected fixture Verification, fresh audit, and second apply
   with an unchanged repository snapshot and empty `plannedChanges`.
4. The provenance test changed one in-memory golden artifact and observed a
   different portable-file digest; generated framing remains protected by the
   byte-for-byte corpus comparison.
5. The focused suite passed with Oxfmt absent from the test `PATH`, without
   network access or package installation.
6. Apply subprocess calls reject Oxfmt execution in the fixture, audit passes
   with a formatter-free `PATH`, and the full gate preserved existing command,
   output, and exit-code coverage.
7. `test_macro_profiles.py` passed all 8 tests, including apply, declared
   formatter comparison, clean audit, and unchanged reapply for every profile.
8. `rtk diff -r .agents/skills/setup-context-driven
   skills/setup-context-driven` exited 0 after synchronization.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_formatter_compatibility.py'`
  passed: 2 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B .agents/skills/setup-context-driven/tests/test_macro_profiles.py`
  passed: 8 tests.
- `rtk make verify` passed: 1,694 Go tests in 20 packages, 170 canonical
  setup-context-driven tests, 170 distributed setup-context-driven tests,
  asset validation, skill synchronization/checks, and build.

Follow-ups: none.
