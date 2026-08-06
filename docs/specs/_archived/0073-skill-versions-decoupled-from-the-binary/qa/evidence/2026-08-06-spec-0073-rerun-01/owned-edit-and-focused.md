# Owned-Skill edit and focused contracts

In an isolated tracked clone, QA appended the same harmless content marker to
`.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, then ran the
real unpiped `rtk make verify` without any regeneration command. It exited 0:
3,522 Go tests passed across 26 packages, the isolated corpus-budget and four
Skill contract tests passed, Skill Check passed, the binary built, and Spec
Check found no Spec 0073 issue.

Postflight `git status --short --untracked-files=all` and `git diff
--name-only` listed only the two edited Skill files. No derived Baseline path,
characterization corpus, or archived Spec changed. The adjacent direct
artifact-stability test independently passed.

The fresh focused command covered
`TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`, readiness comparison,
missing minimum, generated-guide strict and regeneration modes, catalog
content independence, legacy/readoption compatibility, characterization
digest absence, third-party call-ledger exclusion, repository classifications,
and shared Doctor/Skills Check output. Every selected test and subtest passed.

`rtk grep -n 'make verify|exec\\.Command|os/exec' skills/*_test.go` returned no
matches, so no test nests the repository gate on this build.

