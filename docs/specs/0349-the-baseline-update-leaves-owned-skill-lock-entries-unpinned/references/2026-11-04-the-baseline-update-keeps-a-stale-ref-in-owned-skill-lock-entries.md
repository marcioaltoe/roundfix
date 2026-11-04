---
type: fix
status: promoted
created: 2026-11-04
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
reason: null
---

# The baseline update keeps a stale ref in owned skill lock entries

## Symptom

`roundfix baseline update` (0.130.1) installs the Roundfix-owned skills from
the binary and rewrites only the `computedHash` of their existing
`skills-lock.json` entries. The `ref` the upstream skills tool recorded stays,
so each entry names new content with an old catalog commit. The next
`bunx skills update -p` clones the catalog at that `ref`, sees a different
folder hash, reinstalls the old copy, and `roundfix doctor` fails again until
the change is undone. Reported by marcioaltoe-skills (the skills catalog
repository) on 2026-11-04 from the gesttione-solutions/quaestor repository.

## Report

ai-memory message `01a1a5f7-7406-76d3-b98e-659e265b27db` from
marcioaltoe-skills, 2026-11-04 09:50:20 UTC (06:50 BRT), nature `measured`,
subject "baseline update mantém o ref antigo no skills-lock.json e o skills
update desfaz a instalação". Body, verbatim:

> Relato de defeito e sugestão para você avaliar, enviado a pedido do Marcio a partir do repositório marcioaltoe/skills. Não precisa responder.
>
> O QUE ACONTECEU (no repo gesttione-solutions/quaestor, 04/11)
> 1. O `roundfix baseline update` (Roundfix 0.130.1) reinstalou as skills do Roundfix abaixo do mínimo (versão 0.1.179 da skill `roundfix`). O `roundfix doctor` passou a mostrar `skills: ok`.
> 2. No `skills-lock.json`, o baseline update trocou só o `computedHash` dessas entradas. O campo `ref` continuou em `8d71fade68c54693bedbd917371f0514311af157`, um commit de marcioaltoe/skills de 26/10 em que a skill `roundfix` estava na 0.1.141. As 12 entradas fixadas nesse commit são archive-spec, brainstorming, business-analyst, evidence-gate, implement-task, qa-gate, roundfix, setup-context-driven, write-idea, write-prd, write-tasks e write-techspec.
> 3. Em seguida, `bunx skills update -p` (skills 1.7.6) clonou marcioaltoe/skills nesse `ref` e comparou o hash da pasta com o `computedHash`. Os hashes eram diferentes, então ele reinstalou o conteúdo de 26/10. O `roundfix doctor` voltou a acusar falha até a mudança ser desfeita.
>
> Ou seja: a entrada do lock fica incoerente (conteúdo e hash novos, `ref` antigo). Qualquer `skills update` ou `skills experimental_install --force` desfaz o baseline update.
>
> SUGESTÕES (escolha a que fizer sentido para o modelo de Baseline)
> - A. Ao reinstalar uma skill gerenciada, gravar também um `ref` coerente com o conteúdo: o commit de marcioaltoe/skills que tem exatamente aquela versão. Hoje a main, c8b2e85, tem a 0.1.179.
> - B. Remover o `ref` das entradas gerenciadas, se a Baseline não depende dele.
> - C. No mínimo, o `roundfix doctor` detectar entradas do lock em que `ref` e `computedHash` não correspondem e avisar que o `skills update` vai rebaixar essas skills.
>
> Nenhuma skill foi renomeada ou removida em marcioaltoe/skills. Nada muda nos snapshots da Baseline.
>
> Evidence: Quaestor commit 742877d (04/11 06:17 BRT, "atualizar as skills do Roundfix para a 0.130.1"): diff de skills-lock.json só altera computedHash; ref segue 8d71fade68c54693bedbd917371f0514311af157. skills@1.7.6 dist/cli.mjs ~linhas 8684-8740: cloneRepo(source, entry.ref); se computedHash != entry.computedHash, reinstala. marcioaltoe/skills 8d71fad: roundfix 0.1.141; c8b2e85 (main): 0.1.179.

## Where

- `internal/baseline/owned_skills_lock.go` `RecordInstalledOwnedSkills`
  replaces only the `computedHash` value of each existing lock entry that
  names an owned skill (ADR-0189, addendum of 2026-10-27).
- Its callers: `runBaselineUpdateSkillsStageWith` in
  `internal/cli/baseline_update.go` (the skills stage of an approved
  `roundfix baseline update`) and the project target of
  `roundfix skills install` in `internal/cli/cli.go`.
- `roundfix doctor` (`repositorySkillsCheck` in `internal/cli/doctor.go`)
  judges owned skills by version only and never reads their lock entries.

## Measurement of the published skills CLI (2026-11-04)

Read offline from the published package, never run against a repository.

- `npm view skills version` printed `1.7.6` (`time.modified`
  `2026-11-03T00:17:11.892Z`); `npm pack skills@1.7.6` gave
  `skills-1.7.6.tgz`, SHA-256
  `7ab1d2a52f1e0058d2f7d4f425821e6346ce90526304594f2275f008b816563e`.
- `package/dist/cli.mjs`, `updateProjectSkills`: entries are grouped by
  `` `${skill.entry.sourceUrl || skill.entry.source}\n${skill.entry.ref ?? ""}` ``
  and each group is cloned with
  `cloneRepo(cloneSource, skillsForSource[0].entry.ref)`. A skill is left
  alone only when
  `skillPath === skill.entry.skillPath && computedHash === skill.entry.computedHash && onDisk.has(skill.name)`;
  otherwise it is reinstalled from the clone and the lock entry is rewritten
  as `{ ...skill.entry, skillPath, computedHash }`, keeping `ref`.
- `cloneRepo(url, ref)`: with a `ref` it clones `--depth 1 --branch <ref>`
  (falling back to a fetch of the commit for a SHA); without one it clones
  `--depth 1`, the catalog's default branch.
- `projectLockEntry` writes `ref: parsed.ref`, which `JSON.stringify` omits
  when the install named no ref, so an entry without `ref` is the tool's own
  form for a default-branch install.
- `runInstallFromLock` (`experimental_install`) skips a skill whose
  `SKILL.md` exists unless `--force`, and installs from the entry's source,
  folder and `ref` (`appendFolderAndRef`).
- `computeSkillFolderHash`: SHA-256 over each file's relative path and bytes,
  sorted by relative path, skipping `.git` and `node_modules`; Roundfix's
  external skills digest reproduces it (`internal/baseline/assets/lock-hash-compatibility-v1.json`).

## Characterization with roundfix 0.130.1 (2026-11-04)

A temporary repository with a `skills-lock.json` in the tool's format: three
owned entries (`archive-spec`, `roundfix`, `write-prd`) with
`source: marcioaltoe/skills`, `ref: 8d71fade68c54693bedbd917371f0514311af157`,
`sourceType: github`, the catalog `skillPath` and a zero `computedHash`, plus
one external entry (`testing-boss`). `roundfix skills install` printed
`Installed Roundfix skill for project: <dir> (41 file(s))` and one
`Updated skills-lock.json computedHash: <name>` line per owned entry, exited
0, and changed exactly the three `computedHash` values; each `ref` and the
external entry stayed byte for byte.

The local catalog checkout at `c8b2e85` holds
`skills/06-review-repair/roundfix` at version 0.1.179, byte-identical to this
repository's `.agents/skills/roundfix` at `6a0ed904` (`diff -rq`). The
catalog's `skills-registry.json` syncs that copy daily from
`marcioaltoe/roundfix` `main` (`ref: main`, `update: auto`, cron `0 1 * * *`).

## Expected

The lock entry of a skill Roundfix installs describes the installed copy and
pins no other copy, and `roundfix doctor` names every owned entry a
`skills update` would replace.
