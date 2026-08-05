# QA-13 — CLI contract

Status: pass.

Fresh built-CLI evidence:

- Clean text and JSON exited 0.
- Attention text and JSON exited 1 and agreed on result data.
- Missing slug, unknown slug, and `--format yaml` each exited 2 with an
  actionable stderr diagnostic.
- `--format json` before the slug parsed and produced the same one-object
  attention result.
- Help lists the command, formats, read-only promise, and exit codes.

The fresh `TestRunSpecAudit*` selection passed seven command tests, including
decode-to-EOF JSON, holder naming, discovery, and `spec check` non-regression.
