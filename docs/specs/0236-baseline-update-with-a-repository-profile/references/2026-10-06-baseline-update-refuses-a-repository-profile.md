---
type: fix
status: promoted
created: 2026-10-06
spec: 0236-baseline-update-with-a-repository-profile
---

# `baseline update` refuses every repository Baseline Profile

## Problem

Two adopters report the same failure on Roundfix 0.43.0 (`69462a0c`):
Pantheon with `.roundfix/baseline/profiles/pantheon-devops.json` and Oraculum
with `.roundfix/baseline/profiles/oraculum-backend.json`. Both profiles pass
`roundfix baseline profile validate` (`valid (repository)`). `roundfix
baseline update`, even as a preview, computes a valid plan and then exits 1:

```
roundfix: baseline update failed: load profile for snapshot comparison: Unknown built-in Baseline Profile "pantheon-devops".
```

`TrailingSetupSkills` in `internal/baseline/skills_trailing.go` (Spec 0215,
`dd0e5de9`) loads the profile with `loadRestoreProfile(catalog, profileID)`,
which resolves only built-in profiles. `internal/cli/baseline_update.go` calls
it; `roundfix doctor` (`internal/cli/doctor.go`) uses the same comparison and
shows `skills: ok` with the note "snapshot comparison unavailable". The
workaround is `baseline update --no-skills` plus `roundfix skills install
--target project` for the owned skills, which every release note's "run
`baseline update` until current" step does not mention.

## Expected

The snapshot comparison resolves a repository profile the same way the rest of
`update` does, so a repository-profile adopter reaches `state: current` and
receives owned-skill updates. A profile that cannot be resolved names the
profile path, not "built-in".

## Sources

Secondbrain inbox, triaged 2026-10-06:
`inbox/roundfix/_triaged/2026-10-05-baseline-update-falha-com-perfil-de-repositorio.md`
(Pantheon) and
`inbox/roundfix/_triaged/2026-10-05-baseline-update-recusa-perfil-proprio-na-checagem-de-snapshot.md`
(Oraculum).
