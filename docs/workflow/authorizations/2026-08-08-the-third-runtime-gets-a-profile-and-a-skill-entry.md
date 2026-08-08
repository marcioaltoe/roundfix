# Tooling authorization — the third runtime gets a profile and a skill entry (2026-08-08)

On 2026-08-08 the maintainer directed the OpenCode route to be made usable:

> vamos trabalhar nos ajustes necessários para que possamos usar o opencode para
> acesso aos modelos do opencode go e openrouter

Asked where Spec 0088's acceptance ends, the maintainer chose a real Roundfix
Run on an `opencode-go` model, with a working profile in this repository's
configuration. Asked whether the two protected paths that outcome requires are
authorized, the maintainer answered **"Autorizo os dois caminhos"**.

## What this covers

**The skill.** `docs/agents/specific-repository.md` carries a HARD RULE: a pull
request that changes CLI behavior ships the roundfix skill update with it. Spec
0088 changes observable CLI behavior twice. `roundfix doctor` stops reporting
`profiles: ok` while a configured optional-category profile is failing — it
resolves every Agent Work Category the effective configuration defines and
enumerates the ACP Runtimes those tuples reference, so its `profiles:` and
`adapter:` lines both change shape. Configuration validation gains a refusal: a
non-empty `reasoning_effort` on the `opencode` runtime is rejected with the
empty value named as the repair, per ADR-0106. An Agent reading today's skill
would know neither.

**The configuration.** The chosen acceptance requires a Roundfix Run to execute
a real Task on an `opencode-go` model. A Run resolves its Agent Selection from
an Agent Selection Profile, so no Run can reach `opencode` until a profile in
`.roundfixrc.yml` selects it. That file also carries `defaults.verification`,
which places it inside the verification-configuration rule even though this
grant does not touch that key.

## Authorized paths

- `.agents/skills/roundfix/SKILL.md`, limited to describing the widened Agent
  Selection Profile Readiness scope on the Doctor Command, the ACP Runtime
  enumeration that follows it, and the `opencode` model-managed reasoning
  refusal.
- `.roundfixrc.yml`, limited to adding or amending Agent Selection Profiles so
  that at least one Agent Work Category selects `runtime: opencode` with a
  model-managed `reasoning_effort`, together with the comment lines that explain
  the choice.

`.claude/skills/` is a symbolic link to `.agents/skills/`, which is the
authoritative source. The generated copies under `skills/` are rewritten by
`make skills-sync` and are sanctioned fallout, not separate targets.

## Bounded by purpose

This grant covers the two behaviors named above and the profile that makes the
runtime reachable. It does not authorize changing the skill's fetch, resolve,
watch, implement, settle, reconcile, archive, release, or baseline guidance, its
agent bundles under `.agents/skills/roundfix/agents/`, any other skill in the
repository, or any key of `.roundfixrc.yml` outside the `profiles` section.

## Sanctioned fallout — no separate grant

Skill digests and generated copies rewritten by `make skills-sync` and
`make baseline-digests` are deterministic consequences of the authorized source
edit, per ADR-0081. A hand-edited digest remains an unauthorized mutation.

## Consuming Spec

This authorization is consumed by Spec `0088-a-third-runtime-that-can-run`.

## Commit choreography

This record lands as its own commit, before the commits that change the skill
and the configuration.
