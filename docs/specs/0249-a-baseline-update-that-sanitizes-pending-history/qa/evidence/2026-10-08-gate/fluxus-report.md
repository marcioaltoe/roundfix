---
origin: fluxus
destination: roundfix
type-hint: fix
created_at: 2026-10-07
capture: manual
---

# Roundfix 0.55: `history sanitize` recusa o repositório inteiro por um `unproven` em formato antigo

No fluxus, `roundfix history sanitize` (sem flags, só o plano) sai com exit 2:

```text
Preflight failed
Reason:
  build legacy record: parse archive PRD: yaml: unmarshal errors:
  line 10: cannot unmarshal !!map into string
```

A causa são duas Specs arquivadas com o frontmatter `unproven` como lista de mapas, formato de versões antigas do
Roundfix, em vez de lista de texto:

- `docs/history/specs/0032-catalog-full-verification-is-its-own-operation/_prd.md` — itens com `row`, `goal`, `claim`,
  `satisfied-by`.
- `docs/history/specs/0039-the-oraculum-reads-operational-facts/_prd.md` — itens com `row`, `claim`, `reason`, …

Spec arquivada não se edita no lugar, então o repositório fica sem caminho para migrar nenhuma pasta. Duas melhorias
possíveis: aceitar o formato antigo de `unproven` (convertendo cada mapa em texto no Archive Record), ou pular a pasta
com um aviso nomeando o arquivo, em vez de recusar o plano inteiro. A mensagem também não diz qual arquivo falhou.
