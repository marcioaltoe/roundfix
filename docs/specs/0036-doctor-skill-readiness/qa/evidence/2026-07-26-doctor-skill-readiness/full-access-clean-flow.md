# Full-access clean Doctor flow

Build: `6e5618dba5059666c891f9078806631d30502d5b`

This rerun closes the environment-only block recorded in
`qa-report-2026-07-26.md`. The sandboxed QA Agent could observe
`skills: ok`, but its runtime reported the configured Codex binary as
unsigned. A full-access supervisor reproduction proved that
`/Users/marcio/.local/bin/codex` has a valid signature and no quarantine
attribute.

The original clean fixture still inherited this repository's unrelated
frontend profile, whose Claude adapter currently lacks valid capability
evidence. A second disposable fixture copied the same Git repository,
Repository Skill Set, and `skills-lock.json`, then changed only its temporary
project configuration so `frontend` used the already-proven Codex selections:

```yaml
frontend:
  preferred:
    runtime: codex
    model: gpt-5.6-sol
    reasoning_effort: high
  fallbacks:
    - runtime: codex
      model: gpt-5.6-terra
      reasoning_effort: xhigh
```

This fixture change supplies the acceptance criterion's "otherwise ready"
precondition. It does not change the product build, repository configuration,
installed skills, or lock authority.

## Before fingerprint

Command:

```text
rtk proxy find . -type f -not -path '*/.git/*' -print0 | rtk sort -z | rtk xargs -0 shasum -a 256 | rtk shasum -a 256
```

Exit: `0`

```text
30ef13d9aff07f90da8eaacdead275b7c00ddc486da9e24417a027776d0712ba  -
```

## Root invocation

Command:

```text
rtk env CODEX_PATH=/Users/marcio/.local/bin/codex /Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260726T184338Z_c754fa2057c5fc07/bin/roundfix doctor
```

Working directory:
`/tmp/roundfix-qa-0036.KrAsVB/clean-ready`

Exit: `0`

```text
node: ok (v26.5.0 >= 22.13.0)
acpx: ok (0.12.1 >= 0.12.0)
adapter: ok (command="codex-acp"; package=@agentclientprotocol/codex-acp; version=1.1.7)
profiles: ok (2 distinct tuples; 10 category references)
skills: ok (39 required: 14 Roundfix-owned, 25 external)
codex: ok (/Users/marcio/.local/bin/codex has a valid code signature and no com.apple.quarantine attribute)
```

## Nested invocation

The same command ran from
`/tmp/roundfix-qa-0036.KrAsVB/clean-ready/nested/deeper`.

Exit: `0`

The complete output was byte-identical to the root invocation, independently
confirming Git-root discovery and the same `14/25/39` Repository Skill Set.

## After fingerprint

The fingerprint command exited `0` after both Doctor runs:

```text
30ef13d9aff07f90da8eaacdead275b7c00ddc486da9e24417a027776d0712ba  -
```

The before and after fingerprints match. Both public invocations were
read-only and left the complete disposable fixture unchanged.
