Source: https://pkg.go.dev/golang.org/x/mod/sumdb/dirhash#Hash1
Observed: 2026-10-01, published Go package documentation, Hash1 section.

Hash1 uses SHA-256 over a summary sorted by filename. Each summary entry contains the hexadecimal content digest, two ASCII spaces, filename and newline; filenames containing newline are rejected. Hash1 encodes the final digest in base64 with an h1 prefix. ADR-0210 uses the same summary and renders its final digest as lowercase hexadecimal.
