# Sources of Spec 0420

Sections 1 and 2 hold the maintainer's two decisions, verbatim. Section 3 holds what triggered the first one. Section 4 holds what the author measured on 2026-11-24 with the GitHub CLI and on two commits of the default branch: `808f9001` in the first round of the PRD phase, and `cefd4e82` in the second, after Spec 0418 and Spec 0419 merged and v0.153.2 was released. Section 5 holds the published passages the PRD quotes, each read on 2026-11-24. Section 6 is the inventory the codebase sweep returned, kept for the TechSpec.

## 1. DECISÃO DE MARCIO, message 01a20d73-109f-7bb2-a2b6-718854f6b25d from default/praetor (2026-11-24 09:10 BRT)

Relayed whole by the operator's brief for Spec 0420. The Praetor is the session that coordinates the software projects and puts their decision requests to the maintainer.

```text
DECISÃO DE MARCIO
pedido: PEDIDO de 24/11 09:05 BRT (Linux Lane da 0418)
pergunta: Roundfix, portão de testes em Linux que roda no Docker do mini antes do PR: Marcio não usa o Roundfix em Linux. Como fica, a começar pela 0418 que está parada nele? (Tirar do local, só no CI; tirar de vez; subir a memória por ora)
resposta: "Tirar de vez" — com o contexto que Marcio escreveu antes: "Eu estou quase pedindo para remover esses testes do linux, não uso o roundfix no linux, o quaestor-legacy só roda no mini e o quaestor vai rodar em container. Isso está gerando tanto problema. Talvez só deixar algum tipo de teste no CI antes do merge na main."
alcance: (1) Spec própria no repositório roundfix para remover a verificação Linux, tanto a Linux Lane local quanto o job Linux do CI, ficando só macOS; inclui os caminhos de tooling necessários (Makefile, CI, configuração da Lane), implementação, PR, merge e release patch. (2) A 0418 pode seguir sem o portão Linux: se a Lane não puder ser desligada sem a Spec nova, priorize essa Spec antes de retomar a 0418, dentro das filas novas já autorizadas (01a20d32). Se a remoção afetar outros repositórios que usam o Roundfix (Lane em quaestor, vortex etc.), descreva o efeito no PRD e avise o Praetor.
quando: 24/11/2026 09:10 BRT
```

## 2. DECISÃO DE MARCIO, message 01a20f32-07af-7b11-a675-a2cbc7a64743 from default/praetor (2026-11-24 17:18 BRT, 20:18 UTC)

Relayed whole by the operator's addendum to the brief. It answers the Reserved Question the author of the PRD returned in the first round: the first decision names the tooling paths by kind, and eight Governed Paths of the design are of other kinds.

```text
DECISÃO DE MARCIO
pedido: 01a20d89-cf32
pergunta: Roundfix, Spec 0420 (tirar de vez os testes Linux): ela precisa mudar 8 arquivos protegidos que a decisão não citou — os textos da skill do Roundfix que descrevem o portão Linux, o registro de cobertura dos testes e uma frase do guia que chama o CI de "portão Linux". Pode? (Os 8 arquivos; menos a frase do guia; não, só desligar)
resposta: "Os 8 arquivos (Recomendado)"
alcance: o alcance do pedido 01a20d89-cf32, até o release patch: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/{deliver,setup,implement,runtime}.md, skills/roundfix/SKILL.md, docs/references/coverage-record.json (só pelo comando próprio; só saem testes da Lane, do comando dela, da checagem do doctor e das duas chaves), docs/agents/specific-repository.md (só a frase "CI is the Linux gate"), com as regenerações listadas no pedido; mais .roundfixrc.yml, Makefile, .github/workflows/ci-verify.yml e release.yml (CI para macOS). Sem go.mod, lint, Baseline nem outra skill. Chaves antigas aceitas com aviso.
quando: 24/11/2026 17:18 BRT
```

The scope of that decision points at the decision request `01a20d89-cf32`, which the Supervisor session sent to the Praetor. The author did not read the message as sent. The text below is the request as the author drafted it for the session in the first round, and it lists the regeneration commands the scope refers to.

```text
PEDIDO DE DECISÃO
projeto: Roundfix
pergunta: A Spec 0420 (tirar de vez a verificação Linux, decidida em 24/11 09:10) precisa mudar oito arquivos protegidos que a decisão não cita pelo nome: os textos da skill do Roundfix que descrevem a Linux Lane, o registro de cobertura dos testes (os testes da Lane saem com ela) e uma frase do guia do repositório que chama o CI de "portão Linux". Sem a skill, o próprio Roundfix recusa o release; sem o registro, a suíte acusa regressão. Pode mudar?
opções: (1) Sim, os 8 arquivos com a regeneração (Recomendado); (2) Sim, menos a frase do guia (7 arquivos; o guia fica com a frase errada até outra mudança); (3) Não: a Spec encolhe para desligar a Lane e mover o CI para macOS, e o comando, as chaves e a documentação da Lane continuam no produto.
alcance: repositório roundfix, Spec 0420-roundfix-stops-verifying-on-linux, até o release patch. Caminhos: .agents/skills/roundfix/SKILL.md, .agents/skills/roundfix/references/deliver.md, .agents/skills/roundfix/references/setup.md, .agents/skills/roundfix/references/implement.md, .agents/skills/roundfix/references/runtime.md, skills/roundfix/SKILL.md, docs/references/coverage-record.json (só pelo comando próprio; só saem testes da Lane, do comando dela, da checagem do doctor e das duas chaves), docs/agents/specific-repository.md (só a frase "CI is the Linux gate"). Regeneração: make skills-sync; go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage; go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1. Lidos como já cobertos pela decisão de 09:10 (Makefile, CI, configuração da Lane): .roundfixrc.yml, Makefile, .github/workflows/ci-verify.yml e .github/workflows/release.yml (o release também roda a suíte em Linux antes de publicar). Sem go.mod, lint, Baseline nem outra skill.
urgência: o PRD da 0420 espera esta resposta; a 0418 e a 0419 seguem na fila sem depender dela.
```

## 3. What triggered the first decision (operator's brief, 2026-11-24)

- The Linux Lane of Spec 0418's candidate failed twice on the maintainer's machine, at 08:46 and at 08:57 BRT of 2026-11-24, with `compile: signal: killed` while it built the test binary that holds `TestJournalConsumerCorpusReplaysEveryConsumer`. The Docker virtual machine on that host has 7.7 GiB of memory.
- Pull request #950, `chore: turn off the Linux Lane in this repository`, sets `delivery.linux_lane: false` in this repository's Project Config so that Spec 0418 proceeds. It merged on 2026-11-24 as `0cea895a`, and on `cefd4e82` the line reads `linux_lane: false` under six comment lines.

## 4. Measured by the author on 2026-11-24

Repository and checks, read with the GitHub CLI:

- `gh repo view --json visibility,isPrivate` answered `"visibility":"PUBLIC"` and `"isPrivate":false`.
- The default branch's protection requires two status checks: `Verification gate` and `Validate PR title`.
- The 30 latest runs of the workflow `CI | Verify`, from 2026-11-20 to 2026-11-24, all on `ubuntu-latest`: the 13 pull request runs that completed took 8.0 to 15.6 minutes from creation to last update, and the 14 pushes to main took 10.5 to 15.9 minutes.

Workflows on `808f9001`, the same on `cefd4e82`:

| Workflow | Job | Runner | Runs a test |
| --- | --- | --- | --- |
| `CI \| Verify` | `Verification gate` | `ubuntu-latest` | `make verify-changed` and `make verify-docs` on a pull request; `make verify`, `make verify-docs`, `make verify-contracts` and `make suite-budget-check` on a push to main |
| `release` | `release` | `ubuntu-latest` | `make verify` and `make verify-contracts`, before the publication preflight, the cross-compile and the publish steps |
| `CI \| Conventions` | `Validate PR title` | `ubuntu-latest` | none: it validates the pull request title |
| `Secondbrain sync` | `sync` | `ubuntu-latest` | none: it mirrors exported files |

The operator's brief says that only `CI | Verify` runs the Go test suite. The release workflow runs it too, in its `Verify gate` step, so two workflows verify on Linux.

Configuration of the repositories on this machine: 14 files match `/Users/<user>/dev/*/.roundfixrc.yml`, and only this repository's names `linux_lane`.

The maintainer's machine, where the Daemon and the Delivery Queue verify: macOS 27.2.1 on arm64, 10 logical processors, 16 GiB of memory, `/usr/bin/make` is GNU Make 3.81 and `/bin/bash` is 3.2.57.

Compiling every test for Linux from the maintainer's machine, on `cefd4e82`: `GOOS=linux go vet -tags "repocontract docscontract" ./...`, the command of the contract that compiles them for Windows with the other platform, exited 0 in 14 seconds, and in 1 second on a second run.

External Test Tools the tests resolve through `testfixture.RequireTool`: `bash`, `cmp`, `grep`, `gzip`, `jq`, `node`, `rustfmt` and `sh`.

Governed Paths, measured with `speccheck.GovernedPath` in a `go test -overlay` probe whose test file lives outside the repository, over 82 files the change is expected to touch on `808f9001`, and again over 74 files on `cefd4e82`. The source of the function did not change between the two commits, and both probes agree. Governed: `.roundfixrc.yml`, `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runtime.md`, `skills/roundfix/SKILL.md`, `docs/references/coverage-record.json` and `docs/agents/specific-repository.md`. Governed and not to be changed: `.github/workflows/ci-conventions.yml`, `.github/workflows/secondbrain-sync.yml`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `go.mod` and `internal/baseline/assets/templates/starters/python/verify.yml`, the Baseline's starter workflow, which names `ubuntu-latest` for an adopter's repository. Every other file of the probe is ordinary, the mirror `skills/roundfix/references/*.md`, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, the user guide, `README.md`, `CONTEXT.md`, `.secondbrain-export` and every ADR among them.

## 5. Published passages

GitHub Docs, "GitHub-hosted runners reference", https://docs.github.com/en/actions/reference/runners/github-hosted-runners:

- "Use of the standard GitHub-hosted runners is free and unlimited on public repositories."
- The table of standard runners for public repositories lists, for macOS on arm64: 3 processors (M1), 7 GB of memory, 14 GB of storage, and the workflow labels `macos-latest`, `macos-15`, `macos-26` and `xcode-27`. The research sweep read 4 processors and 16 GB for `ubuntu-latest` in the same table.
- Under the limitations of arm64 macOS runners: "Nested-virtualization is not supported due to the limitation of Apple's Virtualization Framework."

GitHub Docs, "GitHub Actions billing", https://docs.github.com/en/billing/concepts/product-billing/github-actions:

- "Larger runners are billed for public repositories."
- The table of per-minute rates lists $0.006 for Linux 2-core and $0.062 for macOS 3-core or 4-core. These rates apply to private repositories.

GitHub Docs, "Actions limits", https://docs.github.com/en/actions/reference/limits, read by the research sweep:

- The Free and Pro plans run at most 5 macOS jobs at the same time.
- "During periods of high demand, jobs that run on macOS runners may take longer to start."

`actions/runner-images`, README, https://github.com/actions/runner-images, read by the research sweep: `macos-latest` names the macOS 26 arm64 image since mid-2026, after naming macOS 15. The macOS 15 arm64 image lists Bash 3.2.57 and Go 1.25, 1.26 and 1.27 among its cached tools.

Rust, "Platform Support", https://doc.rust-lang.org/rustc/platform-support.html, read by the research sweep:

- Tier 1: "automated testing ensures that each tier 1 target builds and passes tests after each change."
- Tier 2: "Automated tests are not always run so it's not guaranteed to produce a working build".

Go, "Porting Policy", https://go.dev/wiki/PortingPolicy, read by the research sweep: a first class port is one where "Broken builds block releases".

Go, "Installing Go from source", https://go.dev/doc/install/source, read by the research sweep: "Note that `$GOARCH` and `$GOOS` identify the target environment, not the environment you are running on. In effect, you are always cross-compiling."

## 6. Inventory of the Linux Lane (codebase sweep on 808f9001, measured again on cefd4e82)

Production code:

- `internal/cli/deliver_linux_lane.go`: the lane, its container script, its package selection, its outcome classification and its opt-in read.
- `internal/cli/deliver_lane.go`: the `roundfix deliver lane` command.
- `internal/cli/deliver.go` and `internal/cli/cli.go`: the usage text, the dispatch and the root help line.
- `internal/cli/deliver_workflow.go`: the call from the repository gate.
- `internal/delivery/engine.go` and `internal/delivery/park_class.go`: the blocker `linux-lane-container`, the warning `linux-lane-skipped`, the pull request paragraph and the Park Class of the blocker.
- `internal/config/config.go`: the two keys, their validation and their merge.
- `internal/cli/doctor.go` and `internal/cli/health.go`: the `linux-lane` check and its four codes.

Tests named after the lane: `internal/cli/deliver_linux_lane_test.go`, `deliver_lane_test.go`, `deliver_linux_lane_go_list_test.go`, `deliver_lane_skip_visible_test.go`, `doctor_linux_lane_test.go`, `internal/config/linux_lane_test.go`, `linux_lane_image_test.go` and `internal/delivery/lane_skip_warning_test.go`. Those eight files declare 35 tests, and the coverage record on `cefd4e82` lists all 35 and one more test of the lane's command, `TestDeliverLanePassesItsVerificationCoreBudget`, in `internal/cli/deliver_core_budget_test.go`: 36 names that contain `Lane`, of which 18 contain `LinuxLane`. Other tests mention the lane in a fixture, a stub or a pinned line and keep their names: `internal/cli/deliver_retry_test.go`, which pins the root help line, and tests in `deliver_owner_log_test.go`, `doctor_disk_space_test.go`, `doctor_verification_cores_test.go`, `machine_verification_status_test.go`, `releaseplan_release_preparation_test.go` and `shared_cache_trim_test.go` of `internal/cli`, and in `gate_blocker_test.go`, `low_disk_warning_test.go` and `park_class_test.go` of `internal/delivery`.

Tests that pin the workflows: `internal/verifyselect/full_contract_run_test.go`, `suite_budget_test.go`, `verify_docs_makefile_test.go`, `release_workflow_test.go` and `release_workflow_compare_test.go`. `TestRepositoryContractTestsDeclareTheirRelevance`, in `internal/verifyselect/contracts_repository_test.go`, pins the number of `docscontract` files: 14 on `cefd4e82`, of which 3 are live documentation contracts. The contract that compiles every test for Windows is `TestEveryTestCompilesForWindows`, in `internal/verifyselect/windows_test_compile_repocontract_test.go`; the Makefile names it in `REPO_CONTRACT_TESTS`, and the coverage record does not list it.

Documents:

- User guide: `docs/user-guide/commands/deliver.md` (the sections `Linux Lane` and `Run the Linux Lane on demand`, and the Park Class table), `docs/user-guide/commands/doctor.md` and `docs/user-guide/configuration.md`.
- Roundfix skill, both copies: the references `deliver`, `setup`, `implement` and `runtime`.
- Glossary: Linux Lane, Lane Environment Skip, Lane Container Failure and Lane Change Failure; Base Test Tools and External Test Tool cite ADR-0294.
- Decision records: ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 are the lane's. ADR-0295, ADR-0384, ADR-0385 and ADR-0423 mention it. ADR-0252 calls CI the Linux gate, and so do `docs/agents/specific-repository.md` and a comment of the `Makefile`.
- Records: `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json` hold the surface `command deliver lane`. `.secondbrain-export` lists ADR-0288, ADR-0303 and ADR-0316.

Existing mechanism for a retired key: Roundfix removes a key it registers as deprecated before the strict decode and prints one line on standard error, in the form `config: <key> is deprecated and ignored; <reason>, so remove it`. Five keys are registered today, three with a removal reason and two with a replacement.

Storage: a lane blocker lives in the item's `blocker` text and a lane skip in the item's `warning` text. No column, table or stage belongs to the lane.

Platform files: production code for Linux only is in `internal/notify/native_linux.go`, `internal/store/process_linux.go` and `internal/suiteguard/residue_linux.go`. Tests that build only for Linux are in `internal/cli/deliver_item_binary_fork_linux_test.go`, `internal/cli/detach_fixture_group_linux_test.go`, `internal/store/process_linux_test.go` and `internal/suiteguard/residue_linux_test.go`. The release platforms of `dist/npm/platforms.json` are darwin and linux on arm64 and amd64, and windows on amd64.

No Baseline module, template or profile names the lane, and `roundfix init` writes neither key.

State of the shared records on `cefd4e82`: the Roundfix skill is at version 0.1.286 in both copies, and the Makefile's suite-time budget is 420 seconds, with a comment that states the Linux runs of 2026-11-06 to 2026-11-10 it was measured on.
