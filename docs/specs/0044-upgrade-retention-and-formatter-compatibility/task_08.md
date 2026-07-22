---
task: task_08
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
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

- [ ] Declare profile formatter contracts and pinned provenance.
- [ ] Correct managed root and shared-guide framing.
- [ ] Build the hermetic TypeScript/Bun golden corpus.
- [ ] Add the full apply-to-reapply composition fixture.
- [ ] Record the exact real Oxfmt QA probe.
- [ ] Exercise formatter contracts across every supported profile.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] Every bundled profile explicitly selects a formatter contract or records
      that no Markdown formatter is selected.
- [ ] Generated TypeScript/Bun Markdown is byte-identical to the pinned Oxfmt
      golden corpus.
- [ ] The composition fixture finishes with an empty repository diff and empty
      `plannedChanges` after formatter comparison, Verification, audit, and
      reapply.
- [ ] Changing generated framing or golden bytes without updating valid
      provenance fails the fixture.
- [ ] Ordinary formatter tests run successfully without network, package
      installation, or an installed Oxfmt executable.
- [ ] Audit and apply execute no formatter process and retain existing exit and
      output contracts.
- [ ] Every profile's apply, audit, and reapply macro flow remains clean.
- [ ] Canonical and distributed setup skill trees are byte-identical.

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
