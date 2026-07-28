# TUI interaction evidence

The built 100×30 PTY flow exercised the owning Live Run View, live event
updates, text phase markers, terminal Clean transition, and `q` terminal
restoration.

Focused current-build Bubble Tea v2 checks passed for:

- Tab/arrow focus and movement;
- Work Queue Task detail and Escape restoration;
- follow suspension while scrolling and resume at tail;
- attach-mode quit leaving the Run untouched;
- too-short terminal fallback;
- review/spec key and viewport regressions;
- styled/no-color Detail and Work Item marker parity.

Model updates were synchronous. No terminal sleep or production-only test hook
was used in these checks.
