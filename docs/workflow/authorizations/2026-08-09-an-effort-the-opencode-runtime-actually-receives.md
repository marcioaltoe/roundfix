# Tooling authorization — an effort the OpenCode runtime actually receives (2026-08-09)

On 2026-08-09 the maintainer asked for a Spec that makes reasoning effort real on
the OpenCode runtime, after learning that three of four candidate models open at
the bottom of their own advertised range:

> Quero que o grok seja usado como preferido no lugar do deepseek v4 ou seja
> usado o deepseek v4 pro no max. […] Vamos a uma spec para alinhar o trabalho
> com effort no opencode.

Asked which protected paths Spec 0089 may mutate, the maintainer answered
**"Autorizo os dois"**, and asked which selection the profile should encode once
the capability exists, answered **"deepseek v4 pro em xhigh"**.

## What this covers

**The skill.** `docs/agents/specific-repository.md` carries a HARD RULE: a pull
request that changes CLI behavior ships the roundfix skill update with it. Spec
0089 removes a refusal the skill currently documents — that `opencode` requires
`reasoning_effort: ""` and that any other value is rejected at configuration
load and again at runtime validation. After this Spec that text is false. The
skill must instead describe the session warm-up, the `runtime_deferred`
encoding, and the proof that splits across preflight and Run.

**The configuration.** The capability is only demonstrated by a profile that
uses it. The maintainer chose `deepseek-v4-pro` at `xhigh` — the model's own
maximum, and the variant the published benchmarks measure.

## Authorized paths

- `.agents/skills/roundfix/SKILL.md`, limited to replacing the OpenCode
  model-managed reasoning refusal with the contract Spec 0089 delivers: the
  session warm-up, the `runtime_deferred` encoding, and what preflight proves
  against what the Run proves.
- `.roundfixrc.yml`, limited to the `profiles` section, setting the `opencode`
  selections to `openrouter/deepseek/deepseek-v4-pro` with
  `reasoning_effort: xhigh`, together with the comment lines that explain the
  choice.

`.claude/skills/` is a symbolic link to `.agents/skills/`, which is the
authoritative source. The generated copies under `skills/` are rewritten by
`make skills-sync`, and derived Baseline pins by `make baseline-digests`; both
are sanctioned fallout under ADR-0081, not separate targets.

## Bounded by purpose

This grant covers the OpenCode reasoning-effort contract and the profile that
proves it. It does not authorize other changes to the skill, to any other skill,
or to any key of `.roundfixrc.yml` outside the `profiles` section.

## Consuming Spec

This authorization is consumed by Spec `0089-an-effort-the-runtime-actually-receives`.

## Commit choreography

This record lands as its own commit, before the commits that change the skill
and the configuration.
