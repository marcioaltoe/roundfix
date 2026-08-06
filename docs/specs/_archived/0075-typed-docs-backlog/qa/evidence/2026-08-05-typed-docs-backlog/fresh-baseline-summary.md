# Fresh Baseline journey

Build: `c24df276`
Scratch repository: detached sparse clone rooted at `914da344`
Profile: `go-cli-tui`

## Plan and apply

- `roundfix baseline plan ... --format json`: exit 0
- Approved plan digest:
  `sha256:da9d3e533b53e118a30a7de97139bb653a6f4f0c913aacb4fbe7105711488b4f`
- `roundfix baseline apply ... --confirm-plan <digest> --format json`: exit 0
- Apply result: `state: verified`; approved postimages and Profile alignment
  verified; no warnings
- Generated `docs/agents/docs-layout.md` SHA-256:
  `000ae9d811d4c28d43cd5af3f13ae8b599c95da8b0194e898d110fa6845dfb91`
- Immediate re-plan: exit 0, `fileChanges: []`, no warnings

The generated guide documented `docs/backlog/`, the five-field frontmatter,
all four type templates, promotion and decline, the finding boundary,
`write-idea` distinction, canonical `refactor` token, and the extension
rule.

## Corpus convergence

In a separate full detached clone at `c24df276`:

- First `rtk make baseline-digests`: exit 0,
  `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}`
- Direct `TestBaselinePlanCharacterization` re-record: environment-blocked
  opening an inherited cache object under
  `/Users/marcio/Library/Caches/go-build`
- Direct `TestCatalogDiagnosticCharacterization` re-record:
  environment-blocked by the same cache
- Second `rtk make baseline-digests`: exit 0, identical
  `changed:false` result
- Final detached-clone `git status --short --branch`: clean,
  `## HEAD (no branch)`

`go env GOCACHE` confirmed the inherited out-of-sandbox cache path. The gate
did not retry the sandbox denial.
