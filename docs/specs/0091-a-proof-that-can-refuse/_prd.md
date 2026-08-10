---
spec: 0091-a-proof-that-can-refuse
status: active
created: 2026-08-09
surfaces: [backend, cli]
---

# A proof that can refuse

`CONTEXT.md` defines Exact Agent Selection Proof as a check that applies a
selection's exact model and reasoning assignment and **observes matching
effective state**. On the `claude` runtime it observes nothing it did not supply.

Measured on 2026-08-09 with the built binary:

```
$ roundfix profiles validate     # profile: claude / opus-9-does-not-exist / high
Profiles validate passed.
1. claude / opus-9-does-not-exist / high — passed
```

The same shape on `codex` is refused with `model_not_advertised`, and on
`opencode` an unadvertised reasoning effort is refused with
`reasoning_control_not_advertised`. Only `claude` accepts anything, and this
repository routes its `frontend` category to `claude` / `opus`, so the category
that most depends on model choice is the one running without a check.

The cause is that Roundfix does not decide what is advertised; it reads the
adapter's refusal out of stderr, matching the string `did not advertise that
model`. `codex-acp` emits it. `claude-agent-acp` does not, and goes further —
measured the same day, asking for a nonexistent model makes it report that model
as the session's `currentValue` **and** append it to the list of advertised
options:

```
model current= 'opus-9-does-not-exist'
    default, opus[1m], claude-fable-5[1m], sonnet, haiku, opus-9-does-not-exist
```

Ensuring the same session without the override returns the honest catalogue:
`default`, `opus[1m]`, `claude-fable-5[1m]`, `sonnet`, `haiku`. So the evidence
Roundfix needs exists and is free to obtain — it is just being read after being
contaminated by the question.

A second, smaller defect travels with this one. When a proof does fail, the real
diagnosis is followed by an unactionable line about closing a disposable session
that was never opened, so the message a maintainer must act on ends in noise
about a session they did not create.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity. The proof keys
  off the existing runtime, model and effort tuple.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — proof runs through `acpx` over ACP
  on a local adapter process and reaches no HTTP surface of this repository's.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0050 has a configured fallback activate
  only after its notification, so every tuple a Run may reach is proven before
  the Run rather than substituted during it; a proof that cannot refuse defeats
  that directly. ADR-0069
  keeps Baseline semantic analysis read-only and supervised; it cites ADR-0050
  but does not apply here, because this Spec changes Agent Selection proof and
  touches no Baseline analysis path. ADR-0114 also cites ADR-0050 and does not
  apply here either: it moves when a Fallback Chain becomes ineligible during a
  Run, while this Spec changes only what preflight will accept before one
  starts. ADR-0091 makes the QA gate a Task node of
  its own type, which is why this Spec's graph carries one. ADR-0096 has the QA
  gate prove machine facts before spending an Agent turn, and this Spec's gate
  follows it: the live refusals are cheap command runs, not Agent work.
  ADR-0104 requires an
  acceptance row on evidence this Spec did not author. This Spec adds ADR-0112 and ADR-0119.
  ADR-0117 places a check with the stage that can produce its defect; it does not change what this Spec delivers, and it moves where this Spec's gate rows run only once Spec 0093 ships. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the change lives in `internal/agent`,
  which is ordinary source. No `Makefile`, CI workflow, skill contract, ignore
  file or version pin is touched. Source: `docs/agents/agent-instructions.md`.

## Goals

1. An Agent Selection naming a model the runtime does not offer is refused at
   preflight, on every runtime, including `claude`.
2. The refusal names the advertised set wherever Roundfix owns the verdict. Where
   an adapter refuses first, its message reaches the operator unchanged and the
   failure travels to whoever is orchestrating, per ADR-0119: the QA gate
   measured all three live runtimes refusing before the membership check could
   speak, and pre-empting every adapter is a race Roundfix loses by
   construction.
3. Proof stays token-free. No prompt is sent to establish what a runtime offers.
4. A failed proof reports one actionable diagnosis, without trailing cleanup
   noise about a session that never opened.
5. Codex and OpenCode proof outcomes are unchanged.

## Core Features

- **A catalogue read before the request.** What a runtime offers is established
  from a session ensured without the requested model, so the answer cannot be
  contaminated by the question. Membership is then decided against that
  catalogue rather than against the adapter's willingness to complain.
- **Refusal owned by Roundfix where no adapter claims it.** An adapter that
  declines to refuse no longer means the selection is sound; where an adapter
  does refuse, its message stands (ADR-0119). The existing
  stderr-matched refusal stays as a fast path where an adapter does emit it.
- **A diagnosis that ends where it stops being useful.** A cleanup failure for a
  session that was never created is not part of the maintainer's problem and is
  recorded rather than appended.

## Non-Goals / Out of Scope

- Validating that a model exists at the vendor. The question is what this
  runtime advertises now, which is what a Run will actually receive.
- Changing the reasoning-effort proof on any runtime. `runtime_deferred` and
  `runtime_managed` keep the behaviour Spec 0089 established.
- Model recommendations or a picker. Selection stays configuration.
- Repairing an invalid profile automatically. Refusal names the problem; the
  maintainer chooses the replacement.

## Decisions

- The honest catalogue is obtained from the same disposable session the proof
  already creates, ensured once without the override. This costs one extra ACP
  round trip and no tokens.
- The contaminated read is not treated as an adapter bug to route around
  silently. Where an adapter echoes the request back into its own advertisement,
  that is recorded in the proof evidence, because a maintainer reading a passing
  proof deserves to know which runtime's word it rests on.
