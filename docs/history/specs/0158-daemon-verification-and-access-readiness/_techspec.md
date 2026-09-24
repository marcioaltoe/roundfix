---
spec: 0158-daemon-verification-and-access-readiness
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Daemon verification and access readiness

## Executive Summary

Let a Task declare its Verification commands independent so one repair turn sees
every failure, let the frozen authorization record name the Tasks allowed to
repair a red repository gate, and prove the requested access mode during
profile readiness.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes and ACP adapters
  only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0011, ADR-0038, ADR-0056, ADR-0057, ADR-0080,
  ADR-0091, ADR-0093, ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0130,
  ADR-0155 and ADR-0156 hold; this Spec adds ADR-0159 and ADR-0160. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection with
  `internal/speccheck/governed.go` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Independent Verification

Task frontmatter gains `verification: independent`. The parser in
`internal/spec/task.go` records it on `spec.Task`; any other value is refused
by name. For such a Task, `runVerificationAttempt` runs every command, publishes
each failure as it does today, and returns an outcome that carries every
failed command. `repairTaskVerification` hands all of them to the one repair
turn; the retry after it runs every command again. A temporary failure keeps
its current handling. Without the declaration the attempt still returns at the
first failure.

## Precondition repair

`_authorization.md` frontmatter gains `precondition_repairs`, a list of Task
identifiers of the consuming Spec. `internal/authorization` parses it; an
identifier that is not a Task in the Graph, or a named Task whose Verification
does not carry the configured repository command verbatim, is refused when the
Run is planned. In `verifyRepositoryPrecondition`, a red result for a named Task
publishes the known-red entry as a Daemon Verification event and lets the Task
start; settlement then requires every Verification command, the repository
command included, to pass. The list is read from the authorization the Run
resolved at start, never from the Run Worktree.

## Access policy in readiness

`agent.RuntimeSpec` reports whether its runtime has an access mode for a
requested policy. The disposable proof in `applyDisposableSelection` applies the
requested mode through the same calls `applyFullAccess` makes, and the proof
records the effective policy. Profile preflight refuses before Run creation when
a runtime has no mode or the adapter refuses it, naming the selection, the
failed predicate and the remedy: disable full access or select a runtime that
supports it. This covers the preferred selection and every fallback the
preflight already proves.

## API Contracts

1. Task frontmatter `verification: independent` makes every command run and
   every failure reach the repair turn.
2. Authorization frontmatter `precondition_repairs: [task_NN]` lets exactly the
   named Tasks enter a red repository gate.
3. Profile readiness applies and records the requested access policy and
   refuses one the runtime cannot honour.

## Coverage Map

- Goal 1 → Independent Verification; API Contract 1.
- Goal 2 → Precondition repair; API Contract 2.
- Goal 3 → Access policy in readiness; API Contract 3.
- Core Feature 1 → Independent Verification.
- Core Feature 2 → Precondition repair.
- Core Feature 3 → Access policy in readiness.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → Independent Verification, Precondition repair, Access
  policy in readiness.

## Integration Points

- **ADR-0038.** The one repair turn stays one; it receives more evidence.
- **Spec 0156.** The delivery loop parks an item whose readiness is refused
  instead of starting a Run that fails every Task.
- **The follow-on authoring Spec.** Documents both declarations in the task and
  authorization guidance.

## Testing Approach

1. **Independent failures.** A fake verifier fails two of three commands; with
   the declaration the repair request carries both failures and the third
   command ran, and without it only the first failure is reported and the later
   commands never ran.
2. **Red precondition.** With the repository command failing, a Task named in
   `precondition_repairs` reaches its Agent turn and settles completed once the
   commands pass; an unnamed Task settles failed with "repository not green on
   entry"; a named Task lacking the verbatim command is refused at planning.
3. **Access policy.** A runtime without a full-access mode and an adapter that
   refuses the mode both fail readiness with the named predicate, and no Run is
   created; a supported mode records its effective policy in the proof.
4. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Independent Verification (depends on: none).
2. Precondition repair (depends on: 1).
3. Access policy in readiness (depends on: 2).
4. Decisions recorded (depends on: 1, 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4, 6, 7).
6. A repair entry that the gate still closes (depends on: 4).
7. Keep every failure and name a degraded access (depends on: 6).

## Risks & Considerations

- **Masking a setup failure.** Running every command after a failure is only
  safe when none sets up another; that is why independence is declared by the
  author, never inferred.
- **A wider entry than granted.** The repair entry is read only from the
  authorization resolved at Run start and only for Tasks carrying the
  configured command verbatim.
- **Readiness cost.** Applying the mode adds one adapter call per proven
  selection, on sessions readiness already opens.
