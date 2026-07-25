<!-- setup-context-driven:begin id=guide.backend version=0.0.1 -->

# Backend

This setup-owned guide defines portable backend rules. Repository-authored
architecture and service contracts remain authoritative.

- Keep blocking, network, process, database, and daemon boundaries explicit about ownership, cancellation, timeouts, and error reporting. Test the lowest real boundary that proves the repository-authored contract; do not invent authentication, database, or transport policy.

## HTTP contract

Application HTTP mode: **Post-only**.

Confirmed ordered exceptions:

1. **Better Auth** owns `GET` and `POST` for `/api/auth/*`: Session, OAuth redirect, callback, and related provider protocol routes require provider-owned GET and POST semantics.

### Better Auth

Better Auth owns the authentication protocol for `/api/auth/*`. Its confirmed `GET` and `POST` exception preserves this provider contract: Session, OAuth redirect, callback, and related provider protocol routes require provider-owned GET and POST semantics.

<!-- setup-context-driven:end id=guide.backend -->
