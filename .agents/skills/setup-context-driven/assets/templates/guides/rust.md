# Rust

Use current crate, toolchain, and CLI documentation before changing Rust APIs or
configuration. Prefer focused `cargo check`, `cargo test`, and lint commands
while iterating.

CLI behavior is public API: flags, stdout, stderr, JSON fields, and exit codes
must be deliberate and tested through observable command behavior.
