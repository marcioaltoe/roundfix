# TUI surface

Terminal UI tests drive model updates synchronously and assert rendered state,
messages, and transitions. Do not use terminal emulation when model-level tests
can catch the behavior.

Keep layout, keyboard, mouse, and accessibility decisions in design or agent
guides instead of the root instruction file.
