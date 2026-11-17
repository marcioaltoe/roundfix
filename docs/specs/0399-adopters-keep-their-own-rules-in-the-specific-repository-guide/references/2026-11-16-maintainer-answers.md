# Maintainer answers for Spec 0399

The maintainer answered three questions through AskUserQuestion on 2026-11-16, about 18:25 BRT, in the Supervisor session that queued this Spec. The operator recorded them in Portuguese, and they are copied here verbatim.

- Alcance: "Texto + detecção (Recommended)" — regra e exemplo em agent-instructions, skill-dispatch e specific-repository; baseline update/doctor, ao achar edição local num arquivo do Baseline ou regra no AGENTS.md raiz, diz qual trecho mover para specific-repository.md.
- Inventário: "Ler os checkouts (Recommended)" — só leitura dos checkouts dos adotantes em ~/dev; sem commit nem mensagem a outro projeto; inventário vira referência da Spec e casos de teste.
- Autorização: "Mesmas do BE (Recommended)" — push, PR, merge, release logo após a Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir.

The operator's brief for the TechSpec added one instruction: no automatic move. The Baseline only names the passages to move.
