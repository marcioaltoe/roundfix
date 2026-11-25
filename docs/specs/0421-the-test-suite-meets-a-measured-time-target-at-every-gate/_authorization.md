---
status: approved
granted: 2026-11-25
action: give each verification gate of this repository a measured Gate Time Target, settle a Task on the light checks, run the whole suite once per Spec, run continuous integration only for a pull request and in parallel jobs on the standard macOS runner, end the override of the release and gate its publication on the Release Provenance of its commit, remove the work the tests repeat, make each gate report its time, retire ADR-0387 and three things ADR-0452 fixes by ADR-0453, and write the Gate Time Rule last as a new ADR
consuming: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
paths:
  - Makefile
  - .github/workflows/ci-verify.yml
  - .github/workflows/release.yml
  - .roundfixrc.yml
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
  - docs/references/coverage-record.json
  - internal/cli/baseline_plan_test.go
  - internal/spec/coverage_test.go
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

# Authority for Spec 0421

This record is `approved`, granted on 2026-11-25. The maintainer authorized the delivery on that day, and the one bound his authorization left open, what `.github/workflows/release.yml` verifies, was answered the same day by the Praetor applying his decision; `## The release` records that answer and how it is held. The record must be on the default branch before the Run starts (`docs/agents/spec-routing.md`: "A new or widened grant must land in the delivery target's ancestry before the consuming squash delivery").

## What the maintainer has decided

Four messages of 2026-11-25, each delivered as `DECISÃO DE MARCIO` blocks to this project's ai-memory inbox by the Praetor, the session that coordinates the software projects, from `default/praetor`. The times are the ones the blocks state.

Message `01a2121b-9124-7571-8f8f-de0f40876364`, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a2109e-f31b-74e3-a162-d0f700280f5c
pergunta: A verificação completa do main e do release falha em macOS porque o pacote internal/cli passa do tempo-limite padrão de 10 min. Opções: -timeout 30m nos alvos test e test-budget; voltar main e release ao Linux; runner macOS maior.
resposta (literal): "Essa verificação do roundfix de 10 minutos já seria um absurdo, ainda mais se rodar em todas as runs, 30 minutos é ainda mais absurdo. É primordial uma revisão dos testes e conjuntos de testes para que tenhamos performance decente."
alcance: NÃO aplicar -timeout 30m nem subir o SUITE_BUDGET_SECONDS. Abrir uma Spec de performance da suíte de testes do Roundfix (o que roda em cada Run, em PR, em main e em release; testes lentos, duplicados, que esperam timeout, que não são herméticos; paralelismo; o que pode sair do gate de cada Run). Autoria da Spec sem Run; mande ao Praetor o plano com uma meta de tempo medida (proposta sua, por gate: Run, PR, main, release) e, separado, o que fazer com o main vermelho e a v0.153.3 enquanto a Spec não sai, como PEDIDO DE DECISÃO com opções. Nada de mudança em Makefile ou workflow sem nova decisão.
quando: 25/11/2026 06:40 BRT
```

Message `01a21256`, 07:58 BRT, as the brief of the Idea stage relays it; the author did not read the message itself:

- Does the Idea stage run? "Rodar a Ideia".
- Are the Gate Time Targets the right ones? "Como proposto".
- What does a Task settlement run? "Só o leve", with the scope "nenhuma chave nova de configuração e nenhuma mudança no Daemon".
- Does the rule of ADR-0452 change? "Retirar a regra e ao final, quando chegarmos ao resultado, produzimos uma nova regra. Pois o que aconteceu até agora é que toda vez que o tempo é estourado, a regra é modificada seja de tempo ou qualquer outra coisa. Então dessa vez quero inverter e no final criar uma regra para que só seja modificar em extrema necessidade por testes de grandes novas features."

Message `01a2127e-c61d-7a93-b702-ccbc6ebd2d45`, which answers the request `01a2126b-9f73-7791-983f-0b6a4cf191b1`, verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a2126b-9f73-7791-983f-0b6a4cf191b1 (pergunta 2)
pergunta: Onde é o "final" em que a regra nova é escrita? Opções: final é a 0421, em 2 Specs; final é a segunda Spec; uma Spec só.
resposta: "Final é a 0421, em 2 Specs (Recomendado)"
alcance: a regra nova é a última entrega da 0421, com as metas medidas em produção; a parte 2 (processos Git, divisão do pacote da CLI) vira backlog e, sob a regra, só pode baixar números.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 3)
pergunta: A ADR-0387 (orçamento único de 420 s) recebe o mesmo tratamento da ADR-0452? Opções: retirar e reescrever no fim; mudar agora; manter.
resposta (literal): "Retirar e criar uma nova no fim, não use o mesmo número de referencia." Esclarecido em seguida ("Qual número não deve ser reaproveitado?"): "O número da ADR".
alcance: a ADR-0387 é retirada agora; a regra final é uma ADR nova, com número próprio, e não uma reescrita da 0387 nem da 0452. No intervalo valem os limites de 25/11 (nada acima de 420 s, sem -timeout maior, sem runner pago), com cada gate mostrando o seu tempo contra a meta.
quando: 25/11/2026 09:17 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 4)
pergunta: Que número a regra final fixa em cada gate? Opções: medido + folga declarada; as metas como estão; medido sem folga.
resposta: "Medido + folga declarada (Recomendado)"
alcance: o resultado medido mais uma folga declarada no PRD, nunca acima da meta aprovada.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 5)
pergunta: Como se julga a "extrema necessidade"? Opções: decisão de Marcio com dossiê medido, só na autoria da feature; exceção com prazo; caso a caso.
resposta: "Sua decisão, com dossiê (Recomendado)" (a opção era: decisão de Marcio, com dossiê medido e só na autoria da feature)
alcance: o pedido só nasce na Ideia ou no PRD da Spec da feature grande, nunca durante um Run, numa corretiva ou no dia do estouro, com o tempo do gate antes e depois, o que já foi cortado, juntado ou movido e o número novo medido; sem o dossiê, o trabalho precisa caber.
quando: 25/11/2026 09:15 BRT

DECISÃO DE MARCIO
pedido: 01a2126b-9f73 (pergunta 6)
pergunta: "Leve por Task, completo uma vez por Spec": só no repositório do Roundfix ou capacidade do produto? Opções: só no Roundfix; produto já na 0421; produto em Spec própria.
resposta: "Só no Roundfix (Recomendado)"
alcance: sem mudar o produto; a TechSpec prova onde a suíte inteira roda uma vez por Spec; se não houver como sem chave nova, volta PEDIDO antes de qualquer Task.
quando: 25/11/2026 09:17 BRT
```

Message `01a213bf-b859-7f60-91a1-201ca006c98a`, which answers the requests `01a212a7-9fe6-7fa3-9b84-d251948359e0` and `01a212ca-e97f`, verbatim, as the Supervisor session saved it for the author:

```text
DECISÃO DE MARCIO
pedido: 01a212a7-9fe6-7fa3-9b84-d251948359e0 (bloco 1, autorização) e acréscimos (i), (ii) e (iii) de 01a212ca-e97f
pergunta: Autoriza a entrega da 0421: 9 arquivos protegidos com limite, regenerações pelos comandos próprios, retirada da ADR-0387 e de partes da ADR-0452 por ADR nova curta, CI em jobs paralelos no macOS, até o release patch; nada acima de 420 s, sem -timeout maior, sem runner pago?
resposta: "Autorizar (Recomendado)"
alcance: o bloco (a) a (d) do 01a212a7 como descrito, mais os acréscimos (i) registro de cobertura (saem os 2 testes do orçamento, entram os do código novo), (ii) a regeneração condicional do registro de superfícies e (iii) a nota de emenda na ADR-0252, AJUSTADOS pela decisão do bloco 4 abaixo. Fora: o que o 01a212a7 já punha fora.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212a7-9fe6 (bloco 2)
pergunta: Como a regra final é entregue? Opções: PR só de documentos; Spec curta 0422.
resposta: "PR só de documentos (Recomendado)"
alcance: PR de fechamento da própria 0421, sem Run, aberto quando as janelas fecharem, só com arquivos não protegidos, sem release.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212a7-9fe6 (bloco 3)
pergunta: O que acontece quando um gate passa do número? Opções: julgar pela janela e segurar o release; cobrar cada execução; só avisar.
resposta: "Julgar pela janela (Recomendado)"
alcance: nenhuma execução isolada falha por tempo; antes do release lê-se a janela, e um gate acima do número segura o release até o tempo voltar ou Marcio decidir com dossiê.
quando: 25/11/2026 ~14:25 BRT

DECISÃO DE MARCIO
pedido: 01a212ca-e97f (pergunta 4, contratos duas vezes no CI)
pergunta: A regra "repository contracts validate at the pull request boundary" muda no CI? Opções: mudar só no CI (uma vez por evento); manter; retirar.
resposta (literal): "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main."
alcance (leitura do Praetor; se discordar ou se faltar autorização, mande PEDIDO): o CI do repositório roda só no PR, como gate do squash merge; o push na main deixa de disparar a verificação. No PR, cada contrato roda uma vez. Fora do CI nada muda: make verify-docs segue obrigatório antes do PR. O "Full Contract Run da main" e as medições de push na main deixam de existir: a janela da regra final passa a ser PR e fechamento de Task (e release, se houver verificação nele). O release: proponha na TechSpec o que ele verifica sob essa regra (por exemplo, nada além do que o PR já provou no commit marcado, ou um conjunto mínimo de publicação); se isso mudar o alcance do release.yml além do já autorizado, mande PEDIDO. A ADR nova da regra final registra esta decisão.
quando: 25/11/2026 ~14:25 BRT
```

The first three messages authorize the authoring of this Spec, without a Run, and decide its shape. The fourth authorizes the delivery. Its first block grants the nine paths with their bounds, the regenerations, the changes of the decision records, the operations up to a patch release and the three additions the TechSpec had found, all adjusted by its fourth block. Its second block grants the closing pull request. Its third decides what follows a gate over its number. Its fourth changes where continuous integration runs, and leaves the release to a proposal and, if the proposal goes beyond the granted bound, to a new decision request, which `## The release` records as answered.

## The release

The fourth block of message `01a213bf-b859-7f60-91a1-201ca006c98a` left what a release verifies to a proposal and, if the proposal went beyond the bound of `release.yml`, to a decision request. The request, `01a213d8-c4d8`, offered three options: the Release Provenance (recommended), a minimal publication set, and the whole verification as already granted. It was answered on 2026-11-25, at about 15:12 BRT, by ai-memory message `01a213da-6946-7920-bae3-fedc57a48923` from `default/praetor`. That answer is not a `DECISÃO DE MARCIO` block. It is recorded, as the Praetor asks, as the aplicação da DECISÃO de Marcio 01a213bf-b859 pelo Praetor, resposta 01a213d8-c4d8: the Praetor applied the maintainer's decision of about 14:25 BRT, "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main.", without a new question. Verbatim, as the Supervisor session saved it for the author:

```text
Praetor → Roundfix, 25/11/2026 ~15:12 BRT, resposta ao 01a213d8-c4d8. Decidido pelo Praetor aplicando a decisão de Marcio de 14:25 BRT (01a213bf-b859), sem nova pergunta: "O CI só deve rodar para liberar o squash merge do pr, não roda depois do merge na main." A opção 3 roda verificação depois do merge, contra a regra; a opção 2 publica uma tag em qualquer commit sem prova de PR. A opção 1 é a única que cumpre a regra sem abrir esse buraco.

Vale a opção 1, com o alcance exatamente como no seu pedido: no release.yml, só o job verify: tirar make verify, make verify-contracts e o setup do Go; pôr o passo de procedência pela API do GitHub com o token do workflow (a tag aponta para o commit que mudou a versão; é o squash merge de um PR na main; o branch do PR continha a ponta da main e tem a mesma árvore; o "Verification gate" desse PR passou; qualquer outra situação recusa, inclusive GitHub fora do ar); runner ubuntu-latest sem o comentário do override; permissões só de leitura (contents, checks, pull-requests). O job de publicação fica igual. Notas de emenda na ADR-0253 e na ADR-0452 pela ADR-0453; tirar dos não-objetivos do PRD o reaproveitamento do veredito. Runbook: atualizar o branch do PR de release antes do merge e disparo manual antes da tag. Nenhuma configuração do GitHub.

Registre no _authorization.md como: aplicação da DECISÃO de Marcio 01a213bf-b859 pelo Praetor, resposta 01a213d8-c4d8. Se Marcio discordar depois, vira emenda. O risco que você registrou (dois PRs verdes que se encontram numa árvore nunca verificada) fica escrito na ADR nova da regra final.
```

How it is held:

- The bound of `release.yml` below is the one this answer grants. The bound the first block granted, the verification back on `macos-26` with its steps regrouped, is replaced by it and is not done.
- The two changes of decision records it names, the notices in ADR-0253 and ADR-0452 for the release, are granted with it.
- If the maintainer disagrees later, his answer becomes an amendment of this Spec, as the Praetor says. The Supervisor session asks the Praetor for the maintainer's own block on the release before the Run starts; until it arrives, this application is the operative answer.

Nothing else is missing. The Task Graph phase measured the Governed Paths again over every file the Tasks declare and needs none beyond the nine.

## What the grant covers

The nine paths above are the Governed Paths among the files the plan changes or may change. The author measured them on 2026-11-25 at `33ece0cf` with `speccheck.GovernedPath`, in a `go test -overlay` probe that wrote nothing to the repository; `references/sources.md`, section 8, holds that measurement. The TechSpec measured twice more at `3f621ac6`, the last time over the 122 paths its revised design creates, edits or deletes: the same nine, and six governed paths it leaves unchanged. Section 9.3 of the same file holds it. The Task Graph phase measured again on 2026-11-25, at `3f621ac6`, over every file the five Tasks declare, in the same kind of probe: the nine paths above are the only governed ones the Tasks declare, and `_tasks.md` names which Task declares each. A path no Task needs stays untouched.

- `Makefile`: what the Incremental Verification selects at a Task settlement and for a pull request, and the lines it prints; the targets the verification jobs call to run their share of the tests and of the contracts; the Complete Suite Run, if the TechSpec places it in a target; the Gate Time Report; the targets and the variable of the suite-time budget, which are removed or stop failing, with no number above 420 seconds; the value of the in-package parallel setting with the comment that states its measurement; the comments that say what each gate runs. No longer `-timeout` on any invocation.
- `.github/workflows/ci-verify.yml`: the verification of a pull request in parallel jobs on the standard macOS runner, `macos-26`, within the five macOS jobs the account runs at once; the required check keeps its name and passes only when every verification job passed; the per-test times and the job summary; the steps of the suite-time budget leave with ADR-0387. By the fourth block the trigger for a push to main leaves, with every step that ran only for a push, and each contract runs once for each pull request event. The pull request trigger, the permissions, the checkout and setup actions and the workflow's name stay, and the job timeout does not rise.
- `.github/workflows/release.yml`, by the answer `## The release` records: only the job `verify` changes. The override of 2026-11-25 leaves with its comment. The job keeps its id and stops running tests: `make verify` and `make verify-contracts` leave it, with its Go setup, and one step proves the Release Provenance of the release commit through the GitHub CLI with the workflow's own token: the tag points at the commit that changed the version, that commit is the squash merge of a pull request into main, the pull request's branch held main's tip and has the same tree, and that pull request's `Verification gate` passed; any other state, GitHub unreadable included, refuses. The job runs on `ubuntu-latest`, because no test runs in it, and its permissions become `contents: read`, `checks: read` and `pull-requests: read`, all read-only. The publication job, its runner, its permissions, its steps and their order stay, and it still needs `verify`. No setting of the repository on GitHub changes.
- `.roundfixrc.yml`: only the comment above `defaults.verification`, and only if what it says about that command stops being true. No key is added, removed or changed in value. The path may stay untouched.
- `docs/agents/specific-repository.md`: the sentences that say what each gate runs, among them, by the fourth block, the sentence of the contracts rule on where the Full Contract Run runs; two new rules, that a code test starts neither `go test` nor `go vet` and that a package builds a binary in one place; and the sentences of the Gate Time Rule, where the numbers live, what follows a gate over its number, as the third block decides, and the change clause. The requirement of the rule "repository contracts validate at the pull request boundary" stays as it is: `make verify-docs` passes before a pull request opens. No other rule changes.
- `internal/spec/coverage_test.go`: the collection that runs the Go tool over the whole module inside a code test, and where the comparison of the coverage record runs. No change to what the coverage record holds or to how it is recorded.
- `internal/cli/baseline_plan_test.go`: only the build of the product binary inside its test, which takes the one build of its package. No assertion changes.
- `docs/agents/setup-context.json` and `docs/references/coverage-record.json`: only through the commands below.

The decision records the maintainer's answers decide are not Governed Paths. The record names them so the grant is whole:

- ADR-0387 is superseded and its file deleted.
- ADR-0452 stays accepted and is amended: the grouping of the verification in one job, the runner type and the sentence about the 420 seconds are withdrawn. Every other decision of it stays.
- ADR-0259, ADR-0257 and ADR-0182 receive the amendment notices their changed sentences need.
- ADR-0453, the next free number, is the short ADR that records the retirement, the three limits of the interval and, by the fourth block, that continuous integration runs only for a pull request.
- The Gate Time Rule is a later ADR with the next free number when it is written, never 0387, 0452 or 0453. It amends ADR-0453, which stays accepted, and it records the decision of the fourth block, as that block asks.
- ADR-0252 receives an amendment notice from ADR-0453 for one consequence, that continuous integration runs the selected contracts twice. This is addition (iii) of the first block.
- ADR-0253 receives an amendment notice from ADR-0453 for the run on every push to main, which the fourth block ends. The block does not name that ADR; the record names it because the block's scope ends what that ADR decided.
- By the answer `## The release` records: the same notice also withdraws the release clause of ADR-0253, and the notice in ADR-0452 also withdraws its sentences on the release's verification job.

## The last delivery

The Gate Time Rule is written after the two measurement windows close, about three days after the first delivery merges. By the second block it is one closing pull request of this Spec, by the Supervisor, without a Run and without a release, that changes documents only and no Governed Path: the new ADR, the reference document with the measured results and the numbers, the glossary's link and the notices of the lifecycle. Its commit, push, pull request and merge are part of what that block grants, and that change quotes the id of the message.

## Sanctioned regeneration

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

This command runs after the `Makefile` change and may change only the `declarationDigest` of the `Makefile` in `docs/agents/setup-context.json`.

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

The only test names the record may lose or gain are the ones the Task Graph lists. Each is a test that moved to a contract, was merged into a named test with its assertions or was renamed; one of the two tests of the suite-time budget, which leave with ADR-0387; or a test of code this Spec adds. The last two kinds are addition (i) of the first block. The TechSpec keeps the comparison of the coverage record in its test, so this command does not change.

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

This third command is conditional, and it is addition (ii) of the first block. It runs only when a Task changes a source of a covered Behavior Surface: in task_02, if the repair of the test that waits for the detach handshake changes `internal/cli/detach.go`, a source of the surfaces `command resolve` and `command watch` (TechSpec, Build Order step 2); in task_04, because the release runbook the Praetor's answer changes, `docs/user-guide/release-runbook.md`, is a source of the surfaces `command release plan` and `guide docs/user-guide/release-runbook.md` (Build Order step 4). The record of command surfaces is not a Governed Path; the command is listed so that its test may rewrite the record. If no such source changes, the command does not run.

No digest or record is hand-edited.

## Limits

- No change to any of these:
  - a longer `-timeout` on any `go test` invocation, and a suite-time budget above 420 seconds;
  - a larger, paid or self-hosted runner, and any setting of the repository on GitHub;
  - any key or value of the Project Config or of the User Config, and what the Daemon runs at a settlement;
  - any other change of the product: no command, flag, message or Run behavior changes;
  - `go.mod` or `go.sum`, and the lint or formatter configuration;
  - the workflows that validate a pull request title and mirror files to the Secondbrain;
  - any Baseline module, template, profile, Setup Snapshot or plan golden, and any skill;
  - any guide under `docs/agents/` other than the sentences and rules named above;
  - any Archive Record;
  - any accepted ADR other than the ones named above.
- No run of continuous integration after a merge, and no workflow that verifies a push to main.
- No time rule is written between the retirement and the Gate Time Rule, and no number of that rule is above an approved Gate Time Target.
- No single run fails for its time. Before a release the windows are read, and a gate over its number holds the release until its time returns or the maintainer decides on a Measured Dossier.
- A release at the `patch` ceiling, with the three version-only paths above, is part of the first block. It is operative only when this record is approved.
