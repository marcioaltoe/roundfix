# Tooling authorization — Secondbrain research sources (2026-09-08)

The maintainer accepted the proposed research-source policy in this session:

> Vou seguir sua sugestão para ajustar a base canonica de secondbrain

The accepted proposal identified
`internal/baseline/assets/modules/secondbrain.json` as the canonical source and
`docs/agents/secondbrain.md` as its generated guide.

The maintainer then extended the same policy:

> Gostaria de adicionar mais uma regra ao secondbrain para que o resultado de
> pesquisas do exa mcp que sejam relevantes ou artigos relevantes encontrados
> sobre os assuntos que são importantes e fundamentais para os projetos sejam
> salvos no secondbrain para serem ingeridos e usados no futuro como base/fonte
> de consulta.

## Authorized scope

- `internal/baseline/assets/modules/secondbrain.json`: require both local
  Secondbrain and Exa MCP consultation before proposals, recommendations,
  product, design, or architecture decisions, and Idea, PRD, or TechSpec
  authoring. Define each source's purpose, require source paths and URLs with
  their influence on the decision, record failed or irrelevant consultations
  without stopping work, and prohibit external discovery of local code or
  disclosure of credentials, private client records, or proprietary source.
  Extend `clause.secondbrain.research-capture` to save relevant Exa findings
  and articles, including relevant articles found outside a dedicated research
  session, as sourced pending research Inbox Entries for later ingestion and
  reuse. Each digest carries source titles and URLs, a summary of the content
  read, project or topic relevance, and the evidence/inference distinction.
  Preserve the index-first workflow and the existing Inbox, citation, safety,
  and ownership obligations. Increment the affected module, guide, and rule
  versions.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`
  and its sibling `manifest.json`: add the four new clause bodies and inventory
  rows required by the catalog's exhaustive Source Baseline contract. Keep all
  existing source entries unchanged; let `make baseline-digests` calculate the
  new byte ranges and digests.
- `internal/baseline/preservation_test.go`: update only the named
  `maintainedSourceBaselineEntries` expectation from 134 to 138 because the
  accepted policy adds four inventoried clauses. The exact cardinality check
  remains intact; the catalog's rejection of missing clause rows is correct.
- `docs/agents/secondbrain.md` and `docs/agents/setup-context.json`: regenerate
  the managed guide and its manifest through the public Baseline update
  command using the edited catalog.
- Every derived artifact rewritten by `make baseline-digests` is sanctioned
  fallout of the canonical module edit under the repository's sanctioned
  digest regeneration rule in `docs/agents/specific-repository.md`.

This authorization covers the accepted policy and its generated artifacts.
It does not change repository Verification, other modules, or skill sources.

## Decision sources

The analysis started with the local Secondbrain's `wiki/index.md` and ran
`rtk proxy qmd query 'decisões arquitetura pesquisa fontes evidência Secondbrain Exa' --all --files --min-score 0.3`.
The query returned relevant results and warned that nine documents needed
embeddings; its coverage was therefore incomplete.

The consulted note was
`/Users/marcio/dev/secondbrain/wiki/concepts/fundamentals-of-software-architecture-capitulo-23-19-architecture-decisions.md`.
Its discussion of decision rationale and analysis paralysis informed the
requirements to record the evidence's influence and continue when a source
cannot contribute. The index consulted was
`/Users/marcio/dev/secondbrain/wiki/index.md`.

Exa MCP located and fetched Michael Nygard's primary article,
[Documenting Architecture Decisions](https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions).
Its treatment of context, rationale, and consequences informed the requirement
to explain how sources affected the decision. The obligation to consult both
Secondbrain and Exa is the maintainer's policy; neither reference establishes
that this particular pair of tools is universally necessary.

The existing `docs/agents/secondbrain.md` already required index-first lookup
and continuation after an unavailable or empty Secondbrain consultation. The
accepted policy retains those obligations and applies the explicit limitation
record to both sources. Exa supplements the authoritative documentation
workflow in `docs/agents/agent-instructions.md`.

## Delivery and verification

The first repository Verification stopped at `fmt-check` because
`internal/cli/baseline_skills_restore_test.go` and
`internal/cli/baseline_assets_sync_test.go` had indentation rejected by `gofmt`.
Both files were byte-identical to `HEAD` before this repair. Applying `gofmt`
to these two files is the prerequisite formatting repair; it changes no test
inputs or assertions and must land separately before the policy change.

Apply directly because the maintainer accepted the policy and its canonical
destination in the same session. Regenerate the derived artifacts, review the
public Baseline update plan, apply the bounded managed refresh, and confirm
that a second plan contains no changes. Run the repository Verification and
documentation contracts, and record their results before handoff.

When these changes are committed, this authorization record must land in its
own commit before the canonical module and generated-artifact changes, as
required by `docs/agents/agent-instructions.md`.
