---
task: task_03
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: pending
type: docs
complexity: low
---

# Task 03: The qa-gate skill keeps evidence under names the repository does not ignore

## Overview

The `qa-gate` skill gains its own section on evidence links, so the gate
saves evidence where the QA Report commit can hold it instead of learning it
from a refusal: names the repository does not ignore (a `.txt` extension for
logs and command output), never a citation of a file it did not keep, and the
rule the Daemon, `qa-report accept` and the Settle Command now apply. This
answers the adopted Argus report of 2026-11-01
([references](references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md)),
whose gate saved `.log` evidence that `.gitignore` excluded.

## Requirements

1. MUST add to the canonical `.agents/skills/qa-gate/SKILL.md` a section
   `### Evidence links resolve in the committed tree`, after
   `### Evidence carries no secret` and before
   `### Declared rows name their criterion`, that states
   (`_techspec.md` → API Contract 8):
   - save evidence under names the repository does not ignore, a `.txt`
     extension for logs and command output, and check a new name with
     `git check-ignore` before linking it;
   - never cite a file you did not keep; the sentence "Never cite a file you
     did not keep." appears verbatim;
   - every local link of the QA Report must resolve in the committed tree:
     an ignored or missing target gets one repair turn, then settles the QA
     Task `failed` with the link, its line and the reason `ignored` or
     `missing`, and `roundfix qa-report accept` and the Settle Command refuse
     the same report;
   - keep every verdict and row status the evidence supports; never change
     one to hide lost evidence.
2. MUST NOT add or change any text inside the `### QA settlement` section or
   any other existing section of the skill.
3. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/qa-gate/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. Never hand-edit a mirror or a
   version field.
4. MUST NOT touch Go code, the Roundfix skill, any guide, any Baseline module,
   `docs/agents/`, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Add the section to the canonical skill.
- [ ] Sync the mirror and record the version.

## Acceptance Criteria

- [ ] Both copies of the `qa-gate` skill carry the new section with the verbatim sentence, and the `### QA settlement` section keeps its bytes.
- [ ] The mirror equals the canonical copy and the owned-skill version record and settlement guidance contracts pass.

## Context

- instruction: `docs/adr/0351-a-qa-report-settles-only-when-its-local-links-resolve-in-the-committed-tree.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The `qa-gate` skill text and the owned-skill version record — consumers: `./skills`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do grep -qF -- '### Evidence links resolve in the committed tree' "$f" || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Never cite a file you did not keep.' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task neither copy has the section, so the first check fails; afterwards both copies carry it, the mirror equals the canonical copy, and the owned-skill version record and the settlement guidance contract pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 7
- `_techspec.md` → Current behavior 6; API Contract 8; Build Order 3
