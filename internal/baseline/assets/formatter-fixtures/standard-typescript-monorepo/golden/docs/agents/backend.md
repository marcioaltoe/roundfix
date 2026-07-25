<!-- setup-context-driven:begin id=guide.backend version=0.0.1 -->

# Backend

This setup-owned guide defines portable backend rules. Repository-authored
architecture and service contracts remain authoritative.

- Keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting. Test the lowest real boundary that proves the repository-authored contract; do not invent authentication, database, or transport policy.

<!-- setup-context-driven:end id=guide.backend -->
