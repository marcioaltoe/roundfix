# Tooling authorization — HTTP Contract default becomes REST (2026-08-07)

On 2026-08-07 the maintainer directed the Baseline catalog's HTTP Contract
Decision default to change from `Post-only` to `REST`:

> Vamos mudar o http contract default para REST — o default do catálogo (vale
> pra frota).

This is a fleet-wide policy change: the value becomes the proposal every
repository adopting or re-prompting a profile that selects `http.contract`
receives. It is recorded here because the decision asset is protected tooling
and the 2026-08-07 Spec 0082 grant does not reach it — that grant is bounded to
`internal/baseline/assets/modules/*.json` for the purpose of teaching the
update command, and a decision default is neither that file set nor that
purpose.

## Authorized paths

- `internal/baseline/assets/decisions.json`, limited to changing the
  `http.contract` decision's `default.mode` value from `Post-only` to `REST`.

No other decision, no other field of the `http.contract` declaration, and no
other asset is authorized by this record. The `modes` list is unchanged and
both modes remain selectable.

## Sanctioned fallout — no separate grant

Derived pins rewritten by `make baseline-digests` are deterministic
consequences of the authorized asset edit, per ADR-0081 and the sanctioned
digest-regeneration rule in `docs/agents/specific-repository.md`. Hand-edited
pin values remain unauthorized mutations.

## Not authorized

- Changing any repository's stored `http.contract` value. Existing manifests
  keep the value they recorded; this changes only what is proposed when the
  decision is answered.
- Repairing the interactive prompt's handling of the typed value. Choosing
  "Change" currently returns only `{"mode": ...}` and discards the recorded
  `exceptions` and `source`. That defect is filed as a finding and needs its own
  Spec and its own authorization; this record does not license a fix.

## Commit choreography

This record lands as its own commit, before the commit that changes the asset.
