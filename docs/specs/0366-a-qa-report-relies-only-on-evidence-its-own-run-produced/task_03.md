---
task: task_03
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
status: completed
type: docs
complexity: low
---

# Task 03: The qa-gate skill writes each pass's evidence in its own directory and lists its executables

## Overview

The `qa-gate` skill gains its own section on evidence that belongs to the
Run, so the gate regenerates evidence instead of learning it from a refusal:
each pass writes under a directory no earlier pass used, builds and runs
executables from its own Run Worktree, lists them in `user_flow_executables`,
and never copies, replays or cites what an earlier Run left. The evidence
directory of step 4 and the report template follow. This answers the message
reported by quaestor on 2026-11-06, whose gate reused
`qa/evidence/2026-11-06-gate/` and an earlier Run's `.venv/bin/quaestor`
([references](references/2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md)).

## Requirements

1. MUST add to the canonical `.agents/skills/qa-gate/SKILL.md` a section
   `### Evidence belongs to this Run`, after
   `### Evidence links resolve in the committed tree` and before
   `### Declared rows name their criterion`, that states
   (`_techspec.md` → API Contract 9):
   - each pass writes its evidence under a directory no earlier pass used,
     `qa/evidence/<report name>/`, the QA Report's file name without `.md`;
   - build and run every executable a row exercises inside your own Run
     Worktree, at the audited head, and list their paths, relative to the
     worktree or absolute inside it, in the `user_flow_executables`
     frontmatter list; for an interpreted entry point, list its script;
   - the sentence "Never copy, replay or cite evidence, a script or an
     executable an earlier Run left, even when it sits on your branch;
     regenerate it." appears verbatim, and only a Carried Row keeps an
     earlier pass's evidence;
   - evidence that breaks this is Borrowed Evidence: a linked file the
     audited head or the QA Report Import already held byte for byte, an
     executable outside your worktree, or a path inside another Run's Run
     Worktree; it gets one repair turn, after that the QA Task settles
     `failed`, and `roundfix qa-report accept` and the Settle Command refuse
     the same report.
2. MUST change, in step 4 of `## 4. Exercise real flows`, only the
   directory `qa/evidence/<YYYY-MM-DD>-<run-slug>/` to
   `qa/evidence/<report name>/`, and add to the report template of
   `## 6. Close the report`, right after the `user_flow_binary` line, the line
   `user_flow_executables: [] # paths, inside this Run Worktree, of the executables that ran user-flow rows`.
3. MUST NOT add or change any text inside the `### QA settlement` section or
   any other existing section of the skill.
4. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/qa-gate/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`,
   recording `skills/testdata/owned-skill-versions.json`. Never hand-edit a
   mirror or a version field.
5. MUST NOT touch Go code, the Roundfix skill, any guide, `CONTEXT.md`, any
   Baseline module, `docs/agents/`, the Makefile, CI workflows, `go.mod` or
   `go.sum`.

## Subtasks

- [ ] Add the section, change the evidence directory and add the template key.
- [ ] Sync the mirror and record the version.

## Acceptance Criteria

- [ ] Both copies of the `qa-gate` skill carry the new section with the verbatim sentence, the new evidence directory and the template key, and the `### QA settlement` section keeps its bytes.
- [ ] The mirror equals the canonical copy and the owned-skill version record and settlement guidance contracts pass.

## Context

- instruction: `docs/adr/0376-a-qa-report-relies-only-on-evidence-its-own-run-produced.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The `qa-gate` skill text and the owned-skill version record — consumers: `./skills`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do grep -qF -- '### Evidence belongs to this Run' "$f" || { printf 'missing section in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Never copy, replay or cite evidence, a script or an executable an earlier Run left, even when it sits on your branch; regenerate it.' || { printf 'missing rule in %s\n' "$f" >&2; exit 1; }; grep -qF -- 'user_flow_executables: []' "$f" || { printf 'missing template key in %s\n' "$f" >&2; exit 1; }; if grep -qF -- 'run-slug' "$f"; then printf 'old evidence directory in %s\n' "$f" >&2; exit 1; fi; done; cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'` — expected: exit 0. Before this Task neither copy has the section, so the first check fails; afterwards both copies carry it with the new directory and template key, the mirror equals the canonical copy, and the owned-skill version record, the settlement guidance contract and the glossary phrase test pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 4
- `_techspec.md` → Current behavior 4; API Contract 9; Declared breaks; Build Order 3
- ADR-0376

## Result

Skills loaded: implement-task, tech-writer, evidence-gate, ai-memory-handoff

Implemented the `qa-gate` evidence contract in the canonical skill. Each pass
now uses `qa/evidence/<report name>/`, records every user-flow executable from
the Run Worktree in `user_flow_executables`, and rejects Borrowed Evidence with
the required repair-and-settlement rule. The report template carries the new
key immediately after `user_flow_binary`. `make skills-sync` regenerated the
mirror, and the owned-skill recorder raised both skill version fields to
`0.0.25` and recorded the new digest in
`skills/testdata/owned-skill-versions.json`.

Focused checks:

- `make skills-sync`: passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed; `1 passed`.
- `cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md`: passed; mirror equal.
- Settlement-section SHA-256 before and after: `ff016a0621700e5fd9e4b7da10e98edc417671b9df49d0416a62edee76de70a4` in both copies.
- `git diff --check`: passed.

Acceptance evidence:

- Both skill copies contain `### Evidence belongs to this Run`, the verbatim
  regeneration sentence, `qa/evidence/<report name>/`, and the
  `user_flow_executables` template key; the settlement section remained byte
  identical.
- The canonical and mirror copies are equal, both carry version `0.0.25`, and
  the owned-skill version record contains the regenerated digest.

The Daemon-owned Task status and declared `## Verification` commands were not
changed or run.
