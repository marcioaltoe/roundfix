---
created: 2026-08-25
wave: 1
status: decomposed
author: autonomous-work
---

# Wave 1 — Desbloqueadores P0 (Decomposição Completa)

## Status: ✅ PRONTO PARA IMPLEMENTAÇÃO

Ambos specs foram totalmente decompostos com TechSpec e Task Graph.

---

## Spec 0098 — Hook Strictness

**Status**: ✅ Decomposição Completa

### Estrutura
- ✅ `_prd.md` — PRD existente (2026-08-12)
- ✅ `_techspec.md` — TechSpec escrita (2026-08-25)
- ✅ `_tasks.md` — Task Graph com 8 tasks

### Design Principal
1. Daemon é autoridade em verificação
2. Hook não pode ser mais strict que gate
3. Settle aceita completed-but-uncommitted
4. Invariant escrito em autonomous-work.md

### Tarefas (8 total)
- task_01: Detect hook refusal
- task_02: Extend settle acceptance
- task_03: Re-run verification in settle
- task_04: Integrate settled commits
- task_05: Handle deleted files
- task_06: Document invariant
- task_07: Acceptance verification (3 cases)
- task_08: QA Gate

### Bloqueadores Removidos
- ✅ Hook refusal é agora recuperável
- ✅ Settle não perde trabalho verificado

---

## Spec 0113 — Gate Report Shape

**Status**: ✅ Decomposição Completa

### Estrutura
- ✅ `_prd.md` — PRD criada (2026-08-25)
- ✅ `_techspec.md` — TechSpec escrita (2026-08-25)
- ✅ `_tasks.md` — Task Graph com 6 tasks

### Design Principal
1. Gate refusada escreve terminal row (blocked)
2. Mechanical stage lê apenas newest report
3. Precondition metadata capturado
4. Spec 0103 pode retry sem deadlock

### Tarefas (6 total)
- task_01: Write terminal row on refusal
- task_02: Detect precondition failure
- task_03: Store precondition metadata
- task_04: Update mechanical stage validation
- task_05: Read newest report only
- task_06: QA Gate

### Bloqueadores Removidos
- ✅ Spec 0103 desbloqueada
- ✅ Loop infinito de rejeições eliminado

---

## Próxima Ação

### Comando para Iniciar Implementação

```bash
# Começar 0098 (hook strictness)
roundfix implement --spec 0098-a-hook-that-cannot-outrank-the-gate --detach

# Acompanhar
roundfix attach <run-id>
# ou
roundfix events <run-id> --follow
```

### Sequência Recomendada
1. ✅ Specs decompostas
2. → `roundfix implement 0098` (Task Capacity: 2, Verification: 1)
3. → Acompanhar até Clean
4. → `roundfix implement 0113` (dependência removida após 0098)
5. → Retry Spec 0103 (agora desbloquedo)

---

## Decisões Gravadas em TechSpec

### 0098 — Hook Refusal Outcome
**Escolha**: Option C — Record and stop with recovery path
- Não sobrescreve hook (seguro)
- Requer ação explícita do Supervisor
- Deixa verificação como autoridade
- Recuperável via settle sem perda

### 0113 — Precondition Metadata
**Escolha**: Terminal blocked row
- Válido per QA Report contract
- Precondition check/reason capturado
- Newest-only reading implementado

---

## Autorização Confirmada

✅ **Tooling Authority** (2026-08-12-the-authoring-and-baseline-corrections.md):
- Makefile mutations (0098 settle, 0113 gate)
- docs/agents/autonomous-work.md (0098 invariant)
- internal/daemon (both specs)
- internal/cli/settle.go (0098)

✅ **No external acceptance required**:
- Ambos são fixes internos ao Roundfix
- Evidência: measured pattern (Specs 0078, 0094, 0103)

---

## Readiness Check

**Antes de começar implementação:**

```bash
# Verificar decomposição
ls -la docs/specs/0098-*/
ls -la docs/specs/0113-*/

# Verificar Task Graphs
wc -l docs/specs/0098-*/_tasks.md docs/specs/0113-*/_tasks.md

# Verificar especificações
grep -c "^## " docs/specs/0098-*/_techspec.md
grep -c "^## " docs/specs/0113-*/_techspec.md
```

---

## Paralelo Librado Após 0098+0113

Wave 2 (independente):
- 0104 (Gate cache cert)
- Vacuous detection

Wave 3 (independente):
- 0097 (Wave collision)
- 0101 (Terminal branch)

Funcionalidade (sempre):
- 0099–0112, 0115
