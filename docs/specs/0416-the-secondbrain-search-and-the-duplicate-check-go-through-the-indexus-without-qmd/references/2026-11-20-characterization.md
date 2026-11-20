# Spec 0416 — characterization, prototype and consultations (2026-11-20)

Recorded by the author of the PRD on 2026-11-20 at the authoring base `851024fa` (v0.152.0). Sections 1 to 4 were measured with commands on the planning worktree or on a scratch clone of it under a temporary directory; nothing was written to the repository but this Spec folder. Sections 5 and 6 record what the Indexus, the Secondbrain and the published sources gave. "Measured" means the author ran the command; "stated" means a source says it and the author did not observe it.

## 1. Where `qmd` is named

Measured with a whole-word, case-insensitive search of the tree outside `.git`.

| Place | Lines | Kind | After this Spec |
| --- | --- | --- | --- |
| `internal/baseline/assets/modules/secondbrain.json` | 67, 115 | Baseline module source, the two clauses | changed |
| `docs/agents/secondbrain.md` | 7, 17 | this repository's guide, written by the Managed Refresh | changed by the refresh |
| `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md` | 7, 17 | formatter golden, derived | changed by regeneration |
| `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md` | 135 | Roundfix-owned skill, canonical copy and mirror | changed |
| `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md` | 32 | accepted ADR, the sentence "read `wiki/index.md` first, then `qmd`" | keeps the sentence, gains an amendment notice |
| `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md` | 20, 61 | Source Baseline corpus | unchanged |
| `internal/speccheck/testdata/receipt-characterization/0181-prd.md`, `0182-prd.md` | 258; 226, 227 | fixtures of the receipt check, copies of old PRDs | unchanged |
| `docs/references/archived-evidence-measurement.md` | 100 | a dated measurement that records an attempted `qmd` query | unchanged |
| `.agents/skills/knowledge-workspace/SKILL.md` | 3, 44, 47 to 50, 53 | Vendored Skill, pinned in `skills-lock.json` to the catalog commit `e45466d2` | unchanged here; the Catalog Skill Bump writes it |
| `docs/history/specs/0079-one-door-for-fleet-knowledge.md` | 18 | Archive Record | unchanged |

No Go file, no other Baseline module, no profile, no template, no user-guide page, no other owned skill, `CONTEXT.md`, `README.md` and `CHANGELOG.md` name it. Outside `docs/history/`, only `CHANGELOG.md` and four ADRs about other subjects name the Indexus, as an adopter, and no file names `pesquisar_acervos`.

## 2. The two clauses and the skill step, as they are

Module `secondbrain`, version 15.

- `clause.secondbrain.02-query-order`, `mandatory`, no `replaces`, in `rule.secondbrain.index-first`: "Read `wiki/index.md` first. Then run `qmd query "<question>" --all --files --min-score 0.3`. Inspect `projects/<project>/mirror/` only when the index and query point there, and open only the files required for the task; treat mirrors as references, not workspaces."
- `clause.secondbrain.research-to-the-knowledge-owner`, `mandatory`, `replaces: ["clause.secondbrain.research-capture"]`, in `rule.secondbrain.research-digests`. Its last sentence: "Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and name a strong verified match in the message so the owner extends existing knowledge instead of duplicating it; a score alone never decides."
- `clause.secondbrain.05-research-limitations`, `mandatory`, unchanged by this Spec, already says what an Agent does without a source: "If either source is unavailable or yields no relevant result, record the attempted consultation and its limitation explicitly, then continue with the available evidence. Do not claim validation from an unavailable, unread, or irrelevant source. Another research tool does not satisfy the Exa MCP requirement."
- `write-idea`, version 0.0.17, research step, in both copies: "When the session performed substantive external research, send a digest with its sources in one ai-memory message to the project that owns the Secondbrain before leaving this step. Run the advisory qmd duplicate check first through an authorized access path, verify that returned paths exist, and name a strong verified match in the message; a score alone never decides." The sentence wraps across lines 135 and 136.

All five built-in profiles list `rule.secondbrain.index-first`. `internal/baseline/clause_characterization_test.go` pins the force of both clauses, and `internal/baseline/cross_project_continuity_test.go` pins that `research-capture` is recorded `replaced` by `research-to-the-knowledge-owner`. The plan code, `classifySourceClauseTransition` in `internal/baseline/plan.go`, compares identity and force and never the text: a clause that keeps both is recorded `retained`.

## 3. Prototype on a scratch clone

The two clauses and the skill step were reworded with trial text. The trial text only made the commands run; the TechSpec decides the wording.

Order that worked: edit the module and the canonical skill; `make skills-sync`; the skill record command; the module record command; `make baseline-digests`; the Managed Refresh. Run before the skill was synced, `make baseline-digests` stopped at `./skills:TestAuthorialSkillSync` with "SKILL.md canonical and distributed bytes differ" and left the formatter golden, the profile pin and the plan goldens unwritten.

Commands, as run:

- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`: wrote module version 16 into `secondbrain.json` and recorded it in `internal/baseline/module-versions.json`.
- `make skills-sync`, then `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: raised both version fields of both `SKILL.md` copies to 0.0.18 and recorded it in `skills/testdata/owned-skill-versions.json`.
- `make baseline-digests`: ended `"ok":true,"changed":true`.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`: "File changes: 2", `docs/agents/secondbrain.md` and `docs/agents/setup-context.json`, one managed entry each. The Setup Manifest's `roundfixVersion` also moved from `0.150.0` to `0.152.0`, which this tree owes whatever the Spec does.

The sixteen files the prototype changed:

```text
.agents/skills/write-idea/SKILL.md
docs/agents/secondbrain.md
docs/agents/setup-context.json
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
internal/baseline/assets/modules/secondbrain.json
internal/baseline/assets/profiles/standard-typescript-monorepo.json
internal/baseline/module-versions.json
internal/baseline/testdata/catalog.diagnostics.golden.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json
internal/baseline/testdata/plan-characterization/clean-adoption.golden.json
internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json
internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json
skills/testdata/owned-skill-versions.json
skills/write-idea/SKILL.md
```

The Source Baseline corpus, manifest and index did not change. After the prototype the word `qmd` was left only in the six unchanged places of section 1.

`go test ./internal/baseline ./skills ./internal/speccheck -count=1` passed on the prototype tree, as it does on the unchanged tree. So no existing test pins the words of the two clauses or of the skill step, and a Verification that must fail before the change needs a new test.

## 4. Governed Paths

Measured with `speccheck.GovernedPath` in a `go test -overlay` probe of `internal/speccheck` whose test file lived outside the repository.

Governed, and changed by this Spec:

```text
.agents/skills/write-idea/SKILL.md
docs/agents/secondbrain.md
docs/agents/setup-context.json
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
internal/baseline/assets/modules/secondbrain.json
internal/baseline/assets/profiles/standard-typescript-monorepo.json
skills/write-idea/SKILL.md
```

Governed, and not changed: the Source Baseline corpus file, `.agents/skills/knowledge-workspace/SKILL.md`, `docs/agents/specific-repository.md`, `.roundfixrc.yml`, `Makefile`, `go.mod`, `.github/workflows/ci-verify.yml`.

Not governed: `internal/baseline/module-versions.json`, the three catalog files and the four plan goldens under `internal/baseline/testdata/`, `skills/testdata/owned-skill-versions.json`, ADR-0330 and a new ADR file, `.secondbrain-export`, `CONTEXT.md`, `internal/baseline/clause_characterization_test.go`, a new test file under `internal/baseline/`, `skills/` or `internal/docscontract/`, `docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json`, `skills-lock.json`, `docs/references/archived-evidence-measurement.md` and the receipt fixtures.

The probe covered the files the prototype changed and the candidates the author expected. The Tasks author measures again every file the Tasks declare.

## 5. The Indexus, measured and stated

Calls from this repository's session on 2026-11-20, through the MCP tools the session exposes:

- `pesquisar_acervos` answered four questions, each with `contrato: "indexus-acervos/1"` and three groups: `geral` (level `equipe`), `gesttione` and `legal` (level `atendimento`). No answer held a group `pessoal`.
- A question about the Baseline and the Secondbrain guide returned relevant pages of the group `geral` under `projects/<project>/mirror/` and `raw/files/`, paths relative to the Secondbrain root.
- Two questions restricted to pages of the source kind `wiki` returned no result in any group (`sem_resultado_relevante`, and `fora_do_filtro` for `legal`), one of them the single word "Indexus", for which `wiki/entities/indexus.md` exists.
- `abrir_conteudo` opened `projects/indexus/mirror/README.md` of the group `geral` and returned its whole text.
- The required `qmd query "<question>" --all --files --min-score 0.3` was attempted once: it ran for more than 120 seconds and ended with `SqliteError: database is locked` (`SQLITE_BUSY`). No result of it was used.

Stated by the Indexus's own documents, read in the Secondbrain at commit `21b715993`:

- `projects/indexus/mirror/docs/references/contrato-indexus-acervos-1.md`, revision 2026-11-18: "A resposta tem `contrato: "indexus-acervos/1"` e `grupos`. Há um grupo por acervo alcançado e pedido."; the levels grow "atendimento, equipe, pessoal"; "a chave de atendimento ou equipe não alcança o Pessoal"; a group with no item of the required relevance carries `motivo_vazio: "sem_resultado_relevante"`; the tools are `pesquisar_acervos` and `abrir_conteudo`.
- `wiki/entities/indexus.md`, updated 2026-11-18: "o wiki inteiro está no Acervo Pessoal até existir uma separação própria por página"; and, of 2026-11-13, "o Indexus respondeu em menos de 1 s e o `qmd` em 49 a 166 s, mas o Indexus acertou só uma", which the search redesign of the days after answered.
- `projects/indexus/mirror/docs/references/2026-11-18-teste-do-jarvis-da-v0-2-0.md`: with the maintainer's personal key, the right page among the first five in 15 of 20 questions, mean time 1.2 s; wiki pages named as `entities/<page>.md`; and "`pesquisar_acervos` falha no cliente MCP do Claude Code com "Structured content does not match the tool's output schema": a sessão guardou o esquema anterior".

Stated by the Secondbrain's message and not observed: "A `wiki` result's `caminho` is relative to `wiki/`." This session received no `wiki` result. The test of 2026-11-18 names wiki pages that way; the contract's example of `abrir_conteudo` shows a path that begins with `wiki/`.

What follows for this Spec:

1. The two measured facts together, no group `pessoal` and no `wiki` result, agree with the two stated ones: this session's key is below the personal level, and the wiki sits in the personal acervo. An Agent of a project session can then find mirrors and raw sources through the Indexus and no wiki page. The clause keeps "Read `wiki/index.md` first", which reads the local file, and the duplicate check stays advisory. Whether this is a use in which the Indexus loses to qmd is for the side-by-side test the maintainer's decision names.
2. "Does not answer" and "not in the session" are ordinary states: a stale tool schema made the tool fail in one client on 2026-11-18, and qmd itself failed today.

## 6. Consultations

Secondbrain files read, all under the Secondbrain root:

- `wiki/index.md`: read first; it lists the pages `indexus` and `indexus-historico` and four Indexus digests.
- `wiki/entities/indexus.md`: the acervos, the levels, the tools, where the wiki sits and the comparison with qmd of 2026-11-13. It gave section 5 its stated facts and the risk about the wiki.
- `projects/indexus/mirror/README.md`, `projects/indexus/mirror/docs/references/contrato-indexus-acervos-1.md` and `projects/indexus/mirror/docs/references/2026-11-18-teste-do-jarvis-da-v0-2-0.md`: the contract the guide's sentence must agree with, and a measurement the Indexus did not design.

The Research Subagent swept the codebase and the web. It holds no shell and no Indexus tool, so the author ran the Indexus calls and the `qmd` attempt.

Published sources, read through Exa on 2026-11-20. Quotes are verbatim from the page text the search returned.

- Jacob Gabrielson, "Avoiding fallback in distributed systems", Amazon Builders' Library, https://aws.amazon.com/builders-library/avoiding-fallback-in-distributed-systems/ : "One of the worst things about fallback is that it isn't exercised regularly and is likely to fail or increase the scope of impact when it triggers during an outage."; "Instead, we favor code paths that are exercised in production continuously rather than rarely." Supports one search path. The article is about services; applying it to a rarely used second tool in an instruction guide is the author's inference.
- Google, Site Reliability Engineering, "Addressing Cascading Failures", https://sre.google/sre-book/addressing-cascading-failures/ : "Remember that the code path you never use is the code path that (often) doesn't work." The same chapter values a degraded answer from the same system, which challenges removing every second path; the limitation clause is that degraded answer here.
- Model Context Protocol specification 2025-06-18, Tools, https://modelcontextprotocol.io/specification/2025-06-18/server/tools : "`listChanged` indicates whether the server will emit notifications when the list of available tools changes." And Authorization, https://modelcontextprotocol.io/specification/2025-06-18/basic/authorization : "403 | Forbidden | Invalid scopes or insufficient permissions". The set of tools and what a credential reaches vary by session, so the guide needs one rule for a tool that is absent or silent.
- Claude Code documentation, "Connect Claude Code to tools via MCP", https://code.claude.com/docs/en/mcp : "The scope you choose controls which projects the server loads in and whether the configuration is shared with your team."
- Claude Code documentation, memory files, https://code.claude.com/docs/en/claude-md : "if two rules contradict each other, Claude may pick one arbitrarily." Supports removing the old command in the same change that names the new tool.
- tobi/qmd README, https://github.com/tobi/qmd : "QMD combines BM25 full-text search, vector semantic search, and LLM re-ranking—all running locally via node-llama-cpp with GGUF models." It explains the latency measured against it: three local models are loaded for `qmd query`.
- Pinecone, "Introducing cascading retrieval", https://www.pinecone.io/blog/cascading-retrieval/ , a vendor's own measurement that hybrid search with reranking beats either method alone on average. It says nothing about this corpus; the side-by-side test does.

How they changed the PRD: the fallback the first request proposed was already dropped by the maintainer, and the two reliability sources say why a second path that is never exercised is a cost; the MCP sources made "not in the session or does not answer" a named case of a Core Feature; the Indexus contract fixed the names the guide may state; the measurement of section 5 became a risk and an open question for the test.

## 7. Measured at the TechSpec stage

Recorded by the author of the TechSpec on 2026-11-20 at `851024fa`, on two scratch clones of the planning worktree under a temporary directory. Nothing was written to the repository but this Spec folder.

The wording. The contract `indexus-acervos/1`, revision 2026-11-18, read in the Secondbrain at commit `760cbc204`, names `caminho` as a field of a page and never says what it is relative to; its one example of `abrir_conteudo` passes `"caminho":"wiki/fiscal/exemplo.md"`. The Indexus's answer (messages and decisions, section 9) says a `wiki` result's `caminho` is relative to `wiki/`, that `abrir_conteudo` takes the path as the search returned it, and that the example will be corrected. The contract gives relevance as the labels `alta`, `media` and `baixa` on groups and items and exposes no number. Its input of `abrir_conteudo` is `{acervo, urn | caminho, cliente?}`.

Two more live calls of `pesquisar_acervos` from the author's session returned the groups `geral`, `gesttione` and `legal`, no group `pessoal`, and, restricted to the source kind `wiki`, no result. Page results of `geral` carried paths that begin with `projects/<project>/mirror/` and `raw/files/`. The `qmd query` the guide still requires was attempted once more and ended with `SqliteError: database is locked`.

Law. A whole-word, case-insensitive search of `internal/baseline/assets/modules/secondbrain.json`, `docs/agents/secondbrain.md` and `.agents/skills/write-idea/SKILL.md` for `leis`, `lei`, `law`, `legal`, `legislation` and `Fiscus` found nothing, and no Baseline module or template names `qmd -c`.

The prototype, with the wording of the TechSpec's API Contracts 1 to 3. The order of section 3 ran clean and changed the same sixteen files. The search clause grew from 264 to 649 characters and the duplicate-check sentence from 255 to 283. A second clone that changed only `write-idea` moved three files, `.agents/skills/write-idea/SKILL.md`, `skills/write-idea/SKILL.md` and `skills/testdata/owned-skill-versions.json`, and `make baseline-digests` then ended `"ok":true,"changed":false`: the skill edit moves no derived Baseline file.

The decision record. With a new ADR whose header says `knowledge: true` and a notice `**Amended by [ADR-NNNN](<file>) (2026-11-20).**` after the first paragraph of ADR-0330's decision, a second run of `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` printed "File changes: 2": `.secondbrain-export` gained the new ADR's path as its last ADR line, and ADR-0330's header gained the new identifier at the end of `amended_by`. It did not change `updated_at`.

The contracts. On the whole prototype tree, nineteen files, `go test -count=1 ./internal/baseline ./skills` passed and `make verify-docs` exited 0. So the skill coverage record, the Behavior Surface Record and the coverage record do not move.

The clause accounting. A test on the scratch clone adopted the Standard TypeScript profile, with the Secondbrain turned on, from a catalog whose search clause carried the earlier sentence, and then planned a Managed Refresh with the catalog of the prototype. The plan was `ready`, changed `docs/agents/secondbrain.md` and `docs/agents/setup-context.json` and nothing else, and recorded `clause.secondbrain.02-query-order` as `retained` with `mandatory` force, `clause.secondbrain.research-capture` as `replaced` with the target `clause.secondbrain.research-to-the-knowledge-owner` and `mandatory` force, 114 clauses `retained`, 16 `replaced` and none `unaccounted`. The accounting walks the clauses of the Source Baseline, so `clause.secondbrain.research-to-the-knowledge-owner`, which the Source Baseline does not hold, has no disposition of its own: it stays the target of the replacement it already was.

How MCP tools are named to a model, read through the Research Subagent on 2026-11-20. Claude Code: "To match one tool, write its full name, `mcp__<server>__<tool>`" (https://code.claude.com/docs/en/permissions). OpenCode: "OpenCode names a tool `<server>_<tool>`." (https://opencode.ai/docs/mcp-servers/). The Model Context Protocol specification: "Tool name uniqueness is scoped to a single server." (https://modelcontextprotocol.io/specification/latest/server/tools). The identifier a model sees differs by harness, so the guide names the server and the tool separately and no harness's joined form. That reading is the author's inference from the three pages.
