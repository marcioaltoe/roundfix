---
task: task_01
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
status: completed
type: docs
complexity: medium
---

# Task 01: The guides, the Roundfix skill and the glossary describe the close of a profile proof

## Overview

The user guides, the Roundfix skill's `profiles` and `runtime` references and
`CONTEXT.md` describe what tasks 02 and 03 change: a proof whose setup
answered passes when only its close fails, the close shares the attempt's
Proof Setup Deadline, an Unclosed Proof Session is recorded and closed by a
later proof, and a setup timeout names its step. The owned-skill versions
and the Behavior Surface record follow through their record commands. This
documents the adopted Backlog Entry of 2026-10-27 (a profile proof fails when
closing its disposable session under load).

## Requirements

1. MUST update `docs/user-guide/configuration.md` section "Profile proof
   setup deadline", keeping its heading: the deadline bounds each attempt's
   setup and the close of its disposable session; a close that fails after a
   passing setup leaves the proof passed with a warning and an Unclosed Proof
   Session record under `~/.roundfix/unclosed-proof-sessions/`, which the
   first proof of a later command closes; a setup timeout names its step as
   `profile proof step "<step>"`. Remove the sentence "The cleanup deadline,
   the Fallback Chain probe after Run creation and Run sessions keep their
   existing deadlines." and replace it with what stays: the Fallback Chain
   probe keeps the built-in setup deadline, and Run sessions and sealed
   prompts keep their cleanup deadlines. In the `profiles configure`
   paragraph, say that a close failure after a passing proof is a warning,
   not a proof failure.
2. MUST update `docs/user-guide/commands/profiles.md` (the proof timeout
   paragraph and the `profiles validate` description) and the canonical
   `.agents/skills/roundfix/references/profiles.md` section "Profile proof
   setup deadline" and its `profiles validate` and `profiles configure`
   paragraphs with the same facts, the `warning` field of
   `roundfix/profiles-validate/v1` and the text line `warning: <text>`
   (`_techspec.md` → API Contracts 4, 5 and 7). Each of the three files MUST
   contain `Unclosed Proof Session`, `.roundfix/unclosed-proof-sessions/` and
   `profile proof step`, and none may still contain `The cleanup deadline,
   the Fallback Chain probe after Run creation and Run sessions keep their
   existing deadlines`.
3. MUST update the canonical `.agents/skills/roundfix/references/runtime.md`
   paragraph on the `profiles:` line so it says a close failure after a
   passing proof is a warning naming the recorded session; it MUST contain
   `Unclosed Proof Session`. New text goes under existing headings or a new
   heading of its own, never inside a `### QA settlement` section.
4. MUST copy the canonical skill to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/roundfix/SKILL.md` and its
   mirror with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`.
5. MUST, through `domain-modeling`, add **Unclosed Proof Session** to
   `CONTEXT.md` and revise **Exact Agent Selection Proof** (its verdict rests
   on the setup; the close follows and a failed close leaves an Unclosed
   Proof Session) and **Proof Setup Deadline** (it also bounds the close of
   that attempt's disposable session), each citing ADR-0344 and never a Spec;
   **Exact Agent Selection Proof** MUST no longer say "closes the Session
   successfully".
6. MUST re-record the Behavior Surface fingerprint as the last step, after
   the last guide edit: `go test -tags docscontract ./internal/docscontract
   -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
   -count=1`, recording `docs/references/behavior-surfaces.json`.

## Subtasks

- [ ] Update the two user guide pages.
- [ ] Update the canonical `profiles` and `runtime` references, sync the mirrors and record the version.
- [ ] Add and revise the glossary terms.
- [ ] Re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Each guide and reference names the Unclosed Proof Session record, the shared close deadline and the named step, and none keeps the old cleanup-deadline sentence.
- [ ] `CONTEXT.md` defines **Unclosed Proof Session** and the two revised terms cite ADR-0344.
- [ ] The owned-skill version and Behavior Surface contracts pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0344-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails-and-the-next-proof-closes-it.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/profiles.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The owned-skill version record and the Behavior Surface record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for phrase in 'Unclosed Proof Session' '.roundfix/unclosed-proof-sessions/' 'profile proof step'; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in the configuration guide: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/profiles.md | grep -qF -- "$phrase" || { printf 'missing phrase in the profiles guide: %s\n' "$phrase" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/profiles.md | grep -qF -- "$phrase" || { printf 'missing phrase in the profiles reference: %s\n' "$phrase" >&2; exit 1; }; done; if cat docs/user-guide/configuration.md docs/user-guide/commands/profiles.md .agents/skills/roundfix/references/profiles.md | tr -s '[:space:]' ' ' | grep -qF -- 'The cleanup deadline, the Fallback Chain probe after Run creation and Run sessions keep their existing deadlines'; then echo 'the old cleanup-deadline sentence remains' >&2; exit 1; fi; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- 'Unclosed Proof Session' || { echo 'runtime reference lacks Unclosed Proof Session' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Unclosed Proof Session**:' || { echo 'CONTEXT.md lacks **Unclosed Proof Session**' >&2; exit 1; }; awk 'index($0, "**Exact Agent Selection Proof**:") == 1 { f=1; next } f && index($0, "_Avoid_") == 1 { exit } f { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF 'ADR-0344' || { echo 'CONTEXT.md **Exact Agent Selection Proof** does not cite ADR-0344' >&2; exit 1; }; awk 'index($0, "**Proof Setup Deadline**:") == 1 { f=1; next } f && index($0, "_Avoid_") == 1 { exit } f { print }' CONTEXT.md | tr -s '[:space:]' ' ' | grep -qF 'ADR-0344' || { echo 'CONTEXT.md **Proof Setup Deadline** does not cite ADR-0344' >&2; exit 1; }; if tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'closes the Session successfully'; then echo 'Exact Agent Selection Proof still requires a successful close' >&2; exit 1; fi; cmp .agents/skills/roundfix/references/profiles.md skills/roundfix/references/profiles.md && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide names `Unclosed Proof Session`, so the first phrase check fails. Afterwards every phrase is present, the old sentence is gone, **Unclosed Proof Session** is defined, **Exact Agent Selection Proof** and **Proof Setup Deadline** cite ADR-0344, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-5; User Experience; Non-Goals / Out of Scope; Success Metric 6; Glossary
- `_techspec.md` → API Contracts 2-7; Glossary; Build Order 3
- ADR-0344

## Result

- Updated both user guides, the canonical `profiles` and `runtime` references,
  and the synchronized skill mirrors to describe setup-based proof verdicts,
  the shared setup/close deadline, named `profile proof step` timeouts,
  `Unclosed Proof Session` records, and close warnings. The validate contract
  now documents the optional `warning` field and `warning: <text>` output.
- Updated `CONTEXT.md` with **Unclosed Proof Session** and revised
  **Exact Agent Selection Proof** and **Proof Setup Deadline**, each citing
  ADR-0344; the former no longer requires a successful close.
- Focused checks run:
  - `make skills-sync` — passed; canonical skill references synchronized.
  - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1` — passed (1 package); both skill version fields raised to `0.1.156` and the owned-skill record updated.
  - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` — passed (1 package); Behavior Surface fingerprint re-recorded after the final guide edit.
- The Task's declared Verification section and repository settlement checks
  remain for the Daemon; no Task status or other Task file was changed.

## Carry-forward provenance

- Source Run: `run_20261027T180614Z_696bb6e78e893361`
- Source commit: `3230c04c0b8e8173f63a9fe654fc6be272d28ea1`
