---
status: approved
granted: 2026-11-12
action: make the Makefile use the Go build cache set with go env -w GOCACHE when it is an existing writable directory, under the 10240 MiB default bound, and otherwise warn once and fall back to its temporary cache under a 2048 MiB default bound without creating the set path; add the test that proves it and make the go stubs of the existing Makefile tests answer env queries; regenerate the setup manifest; rewrite the Makefile cache rule of docs/agents/specific-repository.md, the README paragraph and the Build Cache Bound entry; write the Amended by ADR-0394 notice into ADR-0370
consuming: 0384-the-makefile-follows-the-go-build-cache-set-with-go-env
paths:
  - Makefile
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0384

On 2026-11-12 gesttione-drive reported, in ai-memory message
`01a1d0ad-dd6f-7bb0-ad09-8a1e38911657`, that this repository's Go build cache
held 10 GB of the Mac mini's internal disk although the mini had set its
cache on an external disk with `go env -w GOCACHE`
(`references/2026-11-12-gesttione-drive-the-roundfix-go-build-cache-fills-the-mini-internal-disk.md`,
reported by gesttione-drive). The maintainer then answered the operator's
cycle BD questions with AskUserQuestion (via `/roundfix --follow-up`,
2026-11-12 ~17:05–17:20 BRT). The answers this Spec rests on are, verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o
  dono da fila não morre mais ao esperar a Machine Verification Capacity +
  âncora 'path:line.' aceita, C = cache Go do roundfix fora do disco interno
  do mini (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B
  porque o disco do mini está abaixo dos 80 GB." (C is this Spec.)
- First question "Spec C muda o Makefile (autoridade de ferramenta). Quando o
  GOCACHE do `go env` (mini: /Volumes/Brain-SSD/cache/go-build) aponta para um
  disco não montado, o que o make faz?" → free answer "Não é possível usar o
  cache no SSD?"; clarified by the second question "Sim, o cache vai para o
  SSD: no mini o make passa a usar /Volumes/Brain-SSD/cache/go-build, com
  limite de 10 GB. Falta decidir só o caso raro de o SSD não estar montado.
  Autoriza a mudança no Makefile, e o que o make faz nesse caso?" →
  "Autorizo; sem SSD, TMPDIR com 2 GB" — option text: "O build segue no disco
  interno com limite pequeno e avisa."
- Question "Depois que a Spec C trocar o cache para o SSD, apago os ~10,6 GB
  antigos do disco interno ($TMPDIR/roundfix-go-build-501 e
  ~/Library/Caches/go-build), numa janela sem build do Roundfix nem do
  Quaestor?" → "Sim, depois da troca (Recommended)" — option text: "Confiro
  que nenhum go build/test roda (ps) e apago só esses dois diretórios; é
  cache regenerável." This is an operator action after delivery; no Task
  deletes anything.

"Autorizo" grants the `Makefile` to this Spec. `docs/agents/setup-context.json`
records the `Makefile` digest and changes only through the regeneration below
(lessons 42 and 55 of the operator's brief). "Mesmas do BC" keeps the cycle BC
rules of 2026-11-11, and through them the cycle BB, BA, AZ, AY, AX, AW, AV,
AU, AT and AS answers and the cycle AJ grant of 2026-10-19 "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica)."), which covers
the repository-owned guide `docs/agents/specific-repository.md`: its rule
describes the Makefile's cache choice and would be wrong after the change.

This Spec needs no `.github/workflows`, `go.mod`, `go.sum`, lint, formatter
or `.roundfixrc.yml` change, and no Baseline module, template, profile or
skill change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No released code changes: the binary, the skills and the Baseline are
untouched; `roundfix release plan` still proposes a patch because the
`Makefile` and documentation are not maintenance-only paths, and whether to
cut it is the operator's call. No Task performs a release.

The governed set was measured on 2026-11-12 on the branch `docs/plan-0384`,
based on `d79b30b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0394 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its references, ADR-0394
  and `.roundfixrc.yml`.
- The three paths above are governed: the `Makefile` and the setup manifest
  task_01 changes, and the guide task_02 changes. The probe also named
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md` and
  `.roundfixrc.yml`, which no Task changes; none is granted.

## What is not governed

- The new and changed test files under `internal/verifyselect/`.
- `README.md`, `CONTEXT.md`, ADR-0370 and ADR-0394.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

The repository-owned command writes the setup manifest after the approved
`Makefile` edit. This declaration records that regeneration and adds no
source path.

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

The digest is never hand-edited. No top-level test is removed or renamed, so
the coverage record is not re-recorded, and no guide under
`docs/user-guide/` or skill changes, so the Behavior Surface Record is not
re-recorded.

## Limits

- No change to any `Makefile` line outside the
  `ifeq ($(origin GOCACHE),undefined)` block, the CI workflows, `go.mod`,
  `go.sum`, the lint, formatter or test-runner configuration,
  `.roundfixrc.yml`, `skills-lock.json`, any skill, any Baseline module,
  template, profile, Setup Snapshot or plan golden, any guide under
  `docs/agents/` other than `specific-repository.md`, any user-guide page,
  `docs/references/coverage-record.json`, any Archive Record, Go production
  code, or the `### QA settlement` section of any skill.
- No Task, test or QA row deletes, cleans or moves a Go build cache
  directory, edits the operator's `~/.roundfix-operator/release-checks.sh`,
  or writes the host's Go environment file; the operator deletes the old
  caches after delivery and removes the script's `GOCACHE` export.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix`;
  - creates files in `/tmp`, `/private/tmp`, `/Volumes` or `$HOME` outside its
    own temporary directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
