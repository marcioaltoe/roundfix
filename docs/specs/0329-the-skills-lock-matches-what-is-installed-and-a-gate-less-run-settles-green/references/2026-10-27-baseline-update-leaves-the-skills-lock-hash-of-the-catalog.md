---
type: fix
status: promoted
created: 2026-10-27
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
reason: null
---

# Baseline update leaves the skills lock hash of the catalog

## Symptom

`roundfix baseline update --confirm-plan <digest>` installs the
Roundfix-owned skills bundled in the binary (on 0.124.0: `qa-gate` 0.0.20 and
`roundfix` 0.1.142), which are newer than the copies in the pinned catalog
snapshot (`8d71fade`). The repository's `skills-lock.json` keeps the catalog's
`computedHash` for those skills, so a repository test that compares the lock
with the installed trees fails after the update. The adopter fixed the hashes
by hand and kept the `ref`.

## Where

`skills/repository.go` and `internal/baseline/skills_restore_lock.go` (lock
writing), `internal/baseline/skills_reconcile.go` (lock classification).

## Expected

After an update the lock describes what is installed: the update rewrites the
`computedHash` of every skill it installed from the binary, or the catalog
snapshot carries the same versions as the binary.

## Evidence

Quaestor (ai-memory message `01a17b1a-38e7-7c82-b268-971c34cfcf4d`,
2026-10-26 23:04 BRT), friction 1, PR gesttione-solutions/quaestor#52:

```
1. O baseline update instala as versões de qa-gate e roundfix que vêm no Roundfix (0.0.20 e 0.1.142), mais novas que as do catálogo 8d71fade. Mas o skills-lock.json continua com os computedHash do catálogo, e o teste do repositório que confere o lock contra as árvores instaladas falha. Corrigimos à mão o computedHash, mantendo a ref. Esperado: o update atualizar o lock, ou o catálogo trazer as mesmas versões.
```

Repeated on 0.125.0. Quaestor (ai-memory message
`01a17d0a-cd60-7191-a724-2988affea2d0`, 2026-10-27 08:07 BRT), item 1 under
"Os dois atritos do backlog se repetiram na 0.125.0", PR
gesttione-solutions/quaestor#61:

```
1. Depois do `baseline update`, o skills-lock.json ficou com o computedHash do catálogo nas cinco skills atualizadas. Corrigimos à mão, mantendo as refs.
```

The same message reports the update installing archive-spec 0.0.16, qa-gate
0.0.21, roundfix 0.1.151, setup-context-driven 0.0.7 and write-tasks 0.0.21.
Its "Outro ponto da 0024" paragraph belongs to Spec 0330 and is out of scope
here.
