---
origin: quaestor
destination: roundfix
type-hint: fix
created_at: 2026-10-20
capture: memory-message
---

# Duas contradições nas referências `review.md` e `events.md` da skill roundfix

## Observação

Relato do repositório `marcioaltoe/quaestor`, autorizado por Marcio em 20/10/2026 e recebido como
mensagem do ai-memory (transcrito aqui pelo roundfix em 21/10/2026). Na pré-revisão do baseline
(`roundfix review`, provider codex, 19/10/2026), o Codex apontou duas contradições nas referências da
skill roundfix que o setup copia para `.agents/skills/roundfix/references/` (Roundfix 0.91.0, tag
`v0.91.0` = `f382480b`). O quaestor não corrigiu nada, porque os arquivos são gerenciados pelo setup.

1. `references/review.md`, linhas 128 a 134 na 0.91.0: o prompt do revisor descreve as convenções C3/C5
   de arquivamento antigas, que exigem a pasta da Spec movida e campos de override no `_prd.md`
   arquivado. A entrega atual grava um Archive Record e remove a pasta (`roundfix archive`). Com isso,
   arquivamentos válidos podem gerar achados falsos ou dispensas de validação de convenção incorretas.
2. `references/events.md`, linha 24: a lista de filtros permitidos não inclui `usage`, mas a linha 67
   manda usar `--filter usage` e diz que ele vem incluído por padrão. O contrato público da CLI fica
   contraditório e pode levar um supervisor a rodar um comando inválido.

## Conferência no roundfix (main `956066a7`, 0.101.0)

- `skills/roundfix/references/review.md:157-164` ainda descreve C3 como mover a pasta e C5 com os campos
  no `_prd.md` arquivado; `internal/cli/review_conventions.go:27` e `:29` já dizem Archive Record com a
  pasta removida e os campos de override no record.
- `skills/roundfix/references/events.md:24-26` lista só `task-status`, `batch`, `verification`,
  `outcome` e `agent-selection`; `roundfix events --help` e `internal/runevent/stream.go:23` incluem
  `usage`, e o teste `internal/runevent/stream_usage_test.go:66` exige `usage` no filtro padrão.

## Bom resultado

As duas referências coerentes com o comportamento atual da CLI, e o baseline republicado. O quaestor
não precisa de resposta.
