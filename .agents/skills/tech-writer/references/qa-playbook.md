# QA playbook

Test plans, test cases, exploratory charters, QA bug reports, and verification reports. QA documents trade in evidence: every verdict cites what was observed, on which build, by which path — and states what was _not_ tested.

## Principles

- **User impact first.** Rank everything by what the failure costs a user: Blocks-Completion > Data-Loss > Trust-Damage > Friction > Cosmetic. Severity taxonomies serve this ordering, not the reverse.
- **Evidence before verdicts.** No PASS without fresh evidence from the current build. A screenshot path, a log excerpt, or an observed output — never "looked fine".
- **Traceability both ways.** Test cases cite the persona and journey they exercise; bugs cite the test case or charter that found them; the report rolls all of it up.
- **Document what was not tested.** Skipped personas, blocked flows, and out-of-scope areas appear explicitly with reasons. Silence reads as coverage.
- **Real-user conditions, not lab conditions.** Network quality, device, locale, autofill state, input modality — the conditions a test ran under are part of its result.

## Core concepts

- **Persona** — a named user archetype with attributes. A common canonical set: New User, Power User, Casual User, Mobile User, Accessibility-Reliant, Recovering User.
- **Journey** — entry point → actions (with a time budget per step) → goal → exit, plus branches and abandonment paths. Numbered `J-NN`.
- **Charter** — one exploratory session: mission + persona + surface + exactly one tour + time-box (30/60/90 min). Numbered `CH-NN`. Shape: `Explore <area> / With <tools, data> / To discover <risks>`.
- **Tour** — a themed exploration lens, one per charter. A useful catalog: Feature, Money, Garbage (invalid input), Back-Button, Multi-Tab, Network, Locale, Paste, Autofill, Interrupt.
- **CFR** — cross-functional requirement categories swept on every release: usability, accessibility (WCAG AA quick check), perceived performance, compatibility, error recoverability, production parity.

## Test plan skeleton

```markdown
# <feature> — Test Plan

## Executive summary

<user value delivered and the journey-level risks>

## Personas covered

<each persona exercised, with source citation; skipped personas with reasons>

## Journeys mapped

<each J-NN: entry, goal, abandonment paths>

## Charters planned

<each CH-NN: mission + persona + tour + time-box>

## CFR scope

<which of the six CFR categories this change touches>

## Test strategy

<approach; what is manual, what is automated>

## Automation strategy

<which journeys become E2E specs, which stay manual, which are blocked and why>

## Entry criteria

- Build reachable in a production-parity environment
- CI gate green (run separately)
- Test data matches journey preconditions
- Personas, journeys, and charters documented

## Exit criteria

- Every P0 journey reaches its goal, observably
- Zero open Blocks-Completion or Data-Loss bugs on P0 journeys
- CFR pass on at least 2 journeys with no critical findings open
- Automation follow-up registered for every Missing/Blocked annotation

## Risks

| Risk | Probability | User impact | Mitigation |
| ---- | ----------- | ----------- | ---------- |

## Timeline and deliverables

<dates and the artifacts this plan produces>
```

Distinguish retesting (re-validating the fix of one reported defect) from regression (validating the change broke nothing else — a journey-driven suite).

## Test case skeleton

ID by category: `TC-FUNC-###`, `TC-UI-###`, `TC-REG-###`, `SMOKE-###`, `TC-PERSONA-###`, `TC-JOURNEY-###`, `TC-TOUR-###`, `TC-CFR-###`.

```markdown
# TC-FUNC-001: <what is validated>

Priority: <P0 | P1 | P2 | P3>
Persona: <from the canonical set>
Objective: <what this validates, from the user's perspective>

## Preconditions

<setup, test data, environment state>

## Real-user conditions

Network: <wifi-fast | wifi-slow | 4g | 3g | flaky> · Device: <desktop | laptop | tablet | phone>
Browser: <name+version> · Locale: <locale> · Timezone: <IANA>
Autofill: <empty | stale | current> · Modality: <mouse-keyboard | touch | screen-reader | keyboard-only>

## Steps

1. <action with input> — Expected: <observable result>
2. <action> — Expected: <observable>

## Edge cases

<boundary user behaviors to also try>

## Automation

Target: <E2E | Manual-only> · Status: <Existing | Missing | Blocked | N/A>
Spec/command: <path or command> · Notes: <why automated, manual, or blocked>

## Post-conditions

<resulting system state, cleanup, data verification>

## Execution history

| Date | Tester | Build | Result | Bug ID | Notes |
| ---- | ------ | ----- | ------ | ------ | ----- |
```

Every step has a specific input and a user-observable expected result. Priority by user impact: P0 = failure blocks completion or loses data on a paying path; P1 = trust damage or repeated friction; P2/P3 = secondary surfaces.

## Exploratory charter skeleton

```markdown
# CH-01: <mission as a short phrase>

Explore: <area or feature>
With: <persona, tools, data>
To discover: <risks, in priority order>

Tour: <one from the catalog> · Time-box: <30 | 60 | 90> min · Surface: <pages/flows>

## Session notes

<real-time observations, questions, ideas>

## Debrief

Bugs filed: <BUG-NNN list> · Coverage: <thorough | nominal | compromised>
Not covered: <explicitly out of this session's reach>
Suggested next charter: <follow-up mission, or "none">
```

## Bug report skeleton

Numbered `BUG-NNN`. Title: component, then observed behavior.

```markdown
# BUG-001: <component>: <observed behavior>

Impact: <Blocks-Completion | Data-Loss | Trust-Damage | Friction | Cosmetic>
Severity: <Critical | High | Medium | Low> · Priority: <P0 | P1 | P2 | P3>
Type: <Functional | UI | Accessibility | Usability | Data | Crash> · Status: pending
Persona affected: <name> · Journey step: <J-NN, step N>

## Environment

Build: <version/commit> · OS: <os> · Browser: <name+version>
Viewport: <size> · Network: <condition> · Locale: <locale> · URL: <where>

## Summary

<one paragraph, from the affected persona's perspective>

## Reproduction

Found via: <CH-NN | TC-ID | off-script>
1. <step>
2. <step>
Observed: <what actually happened>

## Expected

<the correct user-side behavior>

## Root cause

<engineering fills in>

## Fix

<engineering fills in>

## Verification

- [ ] Narrow reproduction rerun on the fixed build
- [ ] Affected journey rerun end to end

## Impact

Users affected: <who> · Frequency: <how often> · Workaround: <if any>

## Related

<test case, design link, related journeys or charters>
```

A bug that cannot be reproduced cannot be fixed: numbered steps, exact inputs, and evidence (screenshot path or log excerpt) are non-negotiable. Never assert a root cause as fact — a hypothesis is fine when evidence points somewhere; omit it otherwise.

## Verification report skeleton

```markdown
# Verification report — <feature>

Build: <version/commit> · Generated: <ISO timestamp> · Verdict: <PASS | FAIL>

## Persona coverage

<personas exercised with counts; skipped personas with reasons>

## Journey execution log

<per journey: J-NN, persona, entry URL, goal, step-by-step verb | observed | screenshot | verdict, goal reached?, bugs filed>

## Charter log

<per charter: CH-NN, mission, tour, time-box, findings, bugs, surprises, suggested next charter>

## Off-script findings

<edge cases tried outside planned journeys, results>

## CFR findings

<per category — usability, accessibility, perceived performance, compatibility, error recoverability, production parity: pass | friction | fail>

## Evidence

<screenshot paths, flows tested, viewports, blocked flows with reasons>

## Issues filed

<totals by impact tier: Blocks-Completion / Data-Loss / Trust-Damage / Friction / Cosmetic; release blockers listed individually>
```

PASS requires: every P0 journey reached its goal, zero open Blocks-Completion or Data-Loss bugs. FAIL when either fails. Anything conditional (P1 with documented workaround) is stated, not hidden.

## Regression suites

Journey-driven tiers, not test-case piles:

| Suite    | Duration  | Frequency      | Scope                                 |
| -------- | --------- | -------------- | ------------------------------------- |
| Smoke    | 15–30 min | per build      | 2–4 P0 journeys                       |
| Targeted | 30–60 min | per change     | journeys touching the changed surface |
| Full     | 2–4 h     | weekly/release | all P0 + P1 journeys, every persona   |
| Sanity   | 10–15 min | after hotfix   | the single journey the hotfix affects |

Execution order: Smoke → P0 → P1 → P2 → exploratory charters.

## Quality gates

1. Every verdict backed by fresh evidence from the current build.
2. Every finding cites persona, journey/charter, and step.
3. Every bug independently reproducible from its own steps.
4. Untested areas listed explicitly with reasons.
5. Severity assigned by user impact, not by technical surprise.
