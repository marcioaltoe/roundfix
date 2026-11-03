---
status: approved
granted: 2026-11-03
action: add the Project Config key pre_pr_review.omit_paths that the pre-PR review reads at the base commit and records as omit-path omissions, let a Claude round-1 candidate over its token budget be reviewed in up to four Review Segments, make every size refusal name the three largest files of its diff and the key, park a blocked review with its reason, and describe them in the review, configuration and deliver guides, the Roundfix skill's review and deliver references, CONTEXT.md and the ADRs
consuming: 0342-the-claude-review-segments-and-omits-declared-paths
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0342

On 2026-11-03 the maintainer answered the operator's cycle AW questions with
AskUserQuestion, in the `/roundfix --follow-up` mode, after Vortex reported
that the Claude pre-PR review refuses a candidate by size. The answers this
Spec rests on are, verbatim:

- 0342: "Segmentar + omit_paths (Recommended)" — option text: "Claude também
  divide em até 4 Review Segments; nova chave pre_pr_review.omit_paths
  (globs) deixa arquivos gerados fora da revisão, sem usar derived_paths;
  quando ainda assim bloquear, a mensagem nomeia os 3 maiores arquivos e a
  chave que os tiraria."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em
  falha técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados
  quando o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Those answers keep the cycle AV rules of 2026-11-02 (Rules "Mesmas do AU
(Recommended)", option text: "Sem supervisão: push, PR do plano, entrega e
merge pela fila; Makefile/CI/go.mod/lint liberados se o desenho pedir; falha
técnica = consertar e seguir; findings do review corrigidos ou dispensados com
evidência; mensagem do Pantheon consumida quando a Spec estiver no main."),
and through them the cycle AU answer of 2026-11-01 "Mesmas do AT
(Recommended)", the cycle AS answers of 2026-10-31 ("Consertar e seguir
(Recommended)" and the Tooling answer "Makefile, .github/workflows, go.mod /
go.sum, Config de lint/formatter"), the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.130.0, after Specs 0342,
0343 and 0344, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-03 on the branch `docs/plan-0342`,
based on `4e61bf6d`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted reference and ADR-0352.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair,
  whose version fields rise, and its canonical review and deliver references.

## What is not governed

- `internal/config/config.go`, `internal/config/pre_pr_review.go`,
  `internal/config/pre_pr_review_omit_paths_test.go`.
- `internal/cli/review.go`, `internal/cli/review_segments.go`,
  `internal/cli/review_size_refusal.go` and the review tests the Tasks create
  or change under `internal/cli/`.
- `internal/delivery/engine.go`, `internal/delivery/engine_test.go`,
  `internal/delivery/review_blocked_park_test.go`.
- The mirrors `skills/roundfix/references/review.md` and
  `skills/roundfix/references/deliver.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/review.md`, `docs/user-guide/configuration.md`,
  `docs/user-guide/commands/deliver.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`, ADR-0227, ADR-0282
  and ADR-0352.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  any skill other than the Roundfix skill's `SKILL.md` pair and its review
  and deliver references, or the `### QA settlement` section of any skill.
- No change to conflict regeneration, the `delivery.derived_paths` grammar,
  the 917,504-byte review bound, the 500,000-token Claude budget, the
  four-segment limit or the two-round Reviewer Lineage.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository, Vortex included;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
