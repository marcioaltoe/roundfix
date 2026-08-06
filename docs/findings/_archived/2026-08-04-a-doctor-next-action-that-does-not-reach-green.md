---
date: 2026-08-04
surface: internal/cli, skills
status: done
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-agent-selection-and-execution-environments.md
---

# A Doctor next action that does not reach green

`roundfix doctor` reported one failing check and gave its next action:

```text
skills: failed (validate required skills lock entry ".../skills-lock.json":
external skills are missing: knowledge-workspace;
next: bunx skills add marcioaltoe/skills@knowledge-workspace)
```

The command was run exactly as printed. It succeeded, and it left the
repository in a state where `make verify` — the repository's only gate —
exits 2. Reaching green needed two further changes the Doctor never mentioned:

1. **`skills/recommended.txt`.** `TestRecommendedSkillsMatchLock` asserts the
   recommended manifest equals the lock's entry set. Adding to one half of that
   contract without the other fails the suite.
2. **The upstream skill tree digest.** Three tests pin it as a literal `const`,
   and the declared set moved from 25 members to 26.

The Doctor's remediation is the first of three steps and is presented as the
whole remedy.

## The digest pin is worse than a missing step

`make baseline-digests` is the sanctioned regeneration command, and the
repository's HARD RULE says every derived pin it rewrites is deterministic
fallout needing no separate authorization, while *"a hand-edited pin value
remains an unauthorized mutation"*.

Running it here fails:

```text
baseline-digests: regeneration failed at ./skills:TestAuthorialSkillSync
nextSteps: "Read the failing test output above, fix the canonical source it
validates, then rerun make baseline-digests."
```

The canonical source it validates is `skills-lock.json`, which is now correct.
Only the pin is stale — and the pin is a `const` inside a `_test.go` file, not
a derived artifact under `DERIVED_DIGEST_PATHS`, so the command cannot reach
it and its advice cannot be followed.

That leaves a legitimate, maintainer-directed change with no sanctioned path:
the command that owns pins refuses, and editing the pin by hand is what the
rule calls unauthorized. The rule and the tooling disagree about which pins are
derived.

## And the pin was duplicated three times

The same repository states:

> **an assertion reads the constant it means** — when a value must be
> duplicated, change every occurrence in the same commit; fixing one of three
> is the most repeated defect in this repository's history.

`const wantUpstreamDigest` appeared at three call sites in
`skills/baseline_skill_contract_test.go`, at lines 460, 895, and 935. Three, and
the rule's own example is "one of three". The repair hoisted them to one
package-level constant, which is what the rule asks for; the observation worth
keeping is that the rule existed, was written down, and did not prevent the
duplication.

## Why it matters for autonomous work

A Doctor check is a done-check: an autonomous Supervisor treats "run this next
action" as a closed instruction. A next action that reaches one third of green
converts a two-minute fix into an investigation, and — because the remaining
two thirds land on a Verification pin — into a `stop-and-ask`, which is where
this one ended.

The second brain's loop material makes the general form of this explicit
(`wiki/concepts/agent-workflows-e-loop-engineering.md`): a loop needs a **done
check** that is objective and complete. A remediation that reports success
while the gate stays red is a done check that lies, and it is the failure mode
that stops an autonomous loop most expensively — not by erroring, but by
appearing finished.

## Evidence

- `roundfix doctor` output before and after, 2026-08-04.
- `make verify` exit 2 after the printed remediation; exit 0 after
  `skills/recommended.txt` and the digest constant were also changed.
- `make baseline-digests` refusal, stage `./skills:TestAuthorialSkillSync`,
  `errorCode: regeneration_failed`.
- `skills/baseline_skill_contract_test.go` lines 460, 895, 935 before the
  repair; PR #100.
- `make verify` exit 0, 3,136 tests in 24 packages, after all three changes;
  `roundfix doctor` skills `ok (38 required: 14 Roundfix-owned, 24 external)`.
