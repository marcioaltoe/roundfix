<!-- setup-context-driven:begin id=guide.tui-surface version=0.0.1 -->

# TUI surface

- Drive TUI model updates synchronously and assert rendered state, messages, and transitions. Use terminal emulation only when model-level tests cannot prove the behavior, and keep layout and interaction policy in repository-owned design guidance.

<!-- setup-context-driven:end id=guide.tui-surface -->

<!-- roundfix:repository-rule:begin id=rule.2658ee39ca2d77c945009bf95f8c5035770790be9cfad1e8e70079891e238ff1 -->
- TUI code uses **Bubble Tea v2 module paths** (`charm.land/bubbletea/v2`,
  `charm.land/lipgloss/v2`) and the v2 API (`tea.KeyPressMsg`, `tea.Key`).

<!-- roundfix:repository-rule:end id=rule.2658ee39ca2d77c945009bf95f8c5035770790be9cfad1e8e70079891e238ff1 -->

<!-- roundfix:repository-rule:begin id=rule.3495f2da07f3e9d988b5280f7ba7683eee66b54c72197294faf5cf2500c33927 -->
  Drive `model.Update(...)` synchronously in tests — no terminal emulation.

<!-- roundfix:repository-rule:end id=rule.3495f2da07f3e9d988b5280f7ba7683eee66b54c72197294faf5cf2500c33927 -->
