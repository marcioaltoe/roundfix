# QA-13 — CLI contract

Status: pass.

Fresh built-CLI evidence from `qa-public-fixtures.sh` and the real repository:

- Clean text and JSON exited 0.
- Attention text and JSON exited 1 and agreed on survivors and undelivered
  data.
- Missing slug, unknown slug, and `--format yaml` exited 2 with actionable
  diagnostics.
- `--format json` before the slug parsed successfully.
- JSON parsed as exactly one `roundfix-specaudit/v1` document.
- Top-level and command help list the formats, read-only promise, and exit-code
  meanings.

The fresh `TestRunSpecAudit*` and `TestRunSpecCheck*` selection passed 15 named
command tests, including decode-to-EOF JSON and `spec check` non-regression.
