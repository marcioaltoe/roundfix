# Maintainer answers for Spec 0403

The maintainer answered through AskUserQuestion on 2026-11-17 in the Supervisor session that planned Specs 0399 to 0403. The operator recorded the answers in Portuguese, and they are copied here verbatim.

- 17/11 ~12:40 BRT: Marcio corrigiu o bloco recomendado: "Fable xhigh para Ideia; Fable high para PRD e TechSpec; Opus high para Tasks; Sonnet medium para pesquisa" (aplica à 0402 e à declaração do Roundfix). Pediu "revisão e sugestão/atualização dos modelos já configurados para implementações e revisões. Já peguei configuração antiga sem atualizar." → "Spec 0403 separada (Recommended)": atualiza o Recommended Profile (pesquisa; tabela proposta passa pelo Marcio antes de aplicar) e o update check / --update mostra cada categoria que difere com o comando para corrigir; emenda ADR-0181; depois da 0402. CLAUDE.md global → "Alinhar (Recommended)". Autorização 0403 → "Mesmas do BE (Recommended)".
- 17/11 ~13:40 BRT (tabela da 0403): Frontend → "Manter opus high (Recommended)"; Review e chore → "gpt-6-luna max (Recommended)"; Docs → "gpt-6.1-terra high (Recommended)". Implementação/QA sem mudança (codex gpt-6.1-sol high + claude opus high). Fallbacks propostos: review gpt-6.1-sol high; docs/chore claude sonnet high; frontend codex gpt-6.1-sol xhigh. Snapshot 2026-11-17.

The operator's brief for the TechSpec of 2026-11-17 adds: no `ultra` effort anywhere; update the Model Catalog so every recommended model is in it, accepting `ultra` only if the runtime validation needs it; design the Update Check so that a pinned deviation is reported as pinned to an older snapshot with the command to re-pin or adopt; update this repository's `.roundfixrc.yml` profiles and its stale comments; OpenAI and Anthropic models run only through the subscriptions, never through OpenRouter.

The approved table, snapshot 2026-11-17:

| Agent Work Category | Preferred Selection | Fallback Selection | Change |
| --- | --- | --- | --- |
| `general`, `backend`, `data`, `infra`, `test`, `qa` | `codex / gpt-6.1-sol / high` | `claude / opus / high` | none |
| `frontend` | `claude / opus / high` | `codex / gpt-6.1-sol / xhigh` | none |
| `review` | `codex / gpt-6-luna / max` | `codex / gpt-6.1-sol / high` | preferred was `codex / gpt-5.6-luna / max` |
| `docs` | `codex / gpt-6.1-terra / high` | `claude / sonnet / high` | preferred was `codex / gpt-5.6-luna / max` |
| `chore` | `codex / gpt-6-luna / max` | `claude / sonnet / high` | preferred was `codex / gpt-5.6-luna / max` |
