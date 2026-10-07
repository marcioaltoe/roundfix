---
status: done
created_at: 2026-08-14
updated_at: 2026-08-26
kind: finding
---

# A clean checker left eight vacuous commands standing (2026-08-14)

Three freshly decomposed Specs in `fluxus` passed `roundfix spec check <slug>` —
the full unscoped sweep, not a stage-scoped run — and all three reported
`No findings.` The first `roundfix implement` then refused its first Task at the
starting line. Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-spec-check-limpo-nao-significa-verification-nao-vacua.md`
in the Secondbrain.

The same shape was measured independently in `roundfix` on 2026-08-12/13 while
implementing Spec 0094: eleven authored Verification commands could not fail on
the pre-work tree, every one of them past a `spec check` that reported no
findings, and every one caught instead by the Daemon's pre-work probe at seconds
and zero token cost.

## 1. Eight vacuous commands survived a clean checker

- Symptom / evidence: after the refusal, auditing the three Specs by hand found
  eight commands that pass before any work exists — three naming a suite or a
  file that already exists, and five invoking the repository verification gate.
  The refusal that exposed them named one:

  ```text
  Task task_01 failed: Pre-work Verification vacuous: commands
  "bun run --cwd packages/backend test src/application/use-cases/planogram-publication"
  exited zero against the unchanged tree
  ```

- Root cause: `SC-VERIFY-VACUOUS-COMMAND` is real and works — in the same graph
  it refused four commands written as `git diff --exit-code -- <path>`. It fires
  on some shapes and is silent on others.
- Action / suggestion: execute each authored Verification line in a disposable
  checkout and report its exit status, which is what `spec check
  --run-verification` proposes. A static detector cannot enumerate the shapes
  that pass without work; running the command can.

## 2. A partial detector is read as a gate

- Symptom / evidence: the author corrected the four commands the checker refused
  and read the resulting clean verdict as complete coverage. That reading cost a
  Run.
- Root cause: a detector that fires on some shapes and is silent on others is
  not read as partial. It is read as a gate, and a clean gate ends the checking.
- Action / suggestion: a check that cannot be complete should say what it did
  not cover. The author's own words are the sharpest statement of the cost: had
  the checker carried no detector at all, they would have emulated the check by
  hand — which is exactly what they did after losing a Run.

## What worked — keep

The pre-work probe caught every one of them, in seconds, opening no Agent
Session and spending no tokens. Both repositories measured the same thing: the
probe is the cheapest place this defect class is ever caught, and it catches it
after the Spec is authored rather than while it is being authored.

---

Adopted by Spec `0116-a-verdict-that-states-its-own-scope`.
