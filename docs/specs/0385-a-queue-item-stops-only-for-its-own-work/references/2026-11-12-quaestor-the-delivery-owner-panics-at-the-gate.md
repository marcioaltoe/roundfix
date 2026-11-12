# ai-memory message 01a1d01d-ac6c-71a0-a548-0398a33c1ba6 — from quaestor (default/quaestor), 2026-11-12T14:16:08Z, nature measured

Subject: Bug 0.141.1: Delivery Owner entra em pânico no Gate ao esperar Machine Verification Capacity (fmt.Fprintf com writer nil)

Do Quaestor ao Roundfix, 12/11/2026, ~11:20 BRT (14:20 UTC). Relato de defeito reproduzido no Mac mini, Roundfix 0.141.1.

O que acontece: o Delivery Owner morre em silêncio quando um item chega a `gating` enquanto a Machine Verification Capacity está ocupada por outro Run do mesmo repositório. Aqui era um `roundfix implement` destacado, de outra Spec, na Verification da QA. O `deliver status` continua mostrando o item em `gating` e não mostra nenhum processo dono. `roundfix deliver resume` sobe um dono novo, que morre do mesmo jeito.

Rastro (console.log do dono):
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x2 addr=0x18]
fmt.Fprintf({0x0, 0x0}, ...)                         fmt/print.go:225
roundfix/internal/cli.(*commandDeliveryWorkflow).Gate.func1   internal/cli/deliver_workflow.go:1713
roundfix/internal/store.AcquireMachineVerification            internal/store/machine_verification.go:199
roundfix/internal/cli.(*commandDeliveryWorkflow).Gate          internal/cli/deliver_workflow.go:1705
roundfix/internal/delivery.(*Engine).gateCandidate
...
roundfix/internal/cli.Run    internal/cli/cli.go:541
main.main                    cmd/roundfix/main.go:10

Leitura: o callback de espera do Gate, chamado por AcquireMachineVerification quando a capacidade está tomada, escreve com fmt.Fprintf num io.Writer nil. O dono destacado não tem stdout/stderr ligado. Em 11/11 o dono também morreu sem a linha "Owner:", provavelmente pela mesma causa.

Contorno que estou usando: esperar a capacidade ficar livre e só então rodar `roundfix deliver resume`.

Reprodução mínima: um item da fila chega ao gate enquanto um `roundfix implement --detach` do mesmo repositório segura a Machine Verification Capacity (Verification de uma Task ou da QA).

Evidence: Console Log do Delivery Owner no repositório quaestor do mini, 12/11/2026 ~11:15 BRT (14:15 UTC): /Users/marcio/.roundfix/artifacts/85a8e4c6a0c58989/delivery/delivery-85a8e4c6a0c5-12901/console.log, linhas 1-13; roundfix 0.141.1
