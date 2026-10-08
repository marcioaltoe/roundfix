# External evidence read during QA

Shell curl-cffi lookup of sqlite.org was blocked by the sandbox allowlist. The web reader subsequently reached all three named pages; these are actual external observations, not Spec-created substitutes.

[SQLite PRAGMA documentation](https://sqlite.org/pragma.html): auto_vacuum section and incremental_vacuum section read. Mode 0 reuses freed pages without shrinking the file. Mode 2 needs explicit incremental_vacuum(N), which releases up to N free pages. Moving from mode 0 to mode 2 requires the new pragma value followed by VACUUM. This supports bounded page slices and explicit conversion; it does not guarantee a wall clock bound for each slice.

[SQLite VACUUM documentation](https://sqlite.org/lang_vacuum.html): Description section read. VACUUM repacks a database and may change its auto_vacuum property even in WAL mode. This supports the conversion; the full operation is not an automatic start step.

[Jim Nelson, SQLite, VACUUM, and auto_vacuum](https://blogs.gnome.org/jnelson/2015/01/06/sqlite-vacuum-and-auto_vacuum/): auto_vacuum and auto-vacuum strategies sections read. The article proposes a small incremental page count at startup to gradually release free pages. It also notes remaining fragmentation and states that the author had not personally used incremental mode. This is a design rationale, not a measured latency guarantee.

The original live-database measurement is historical, not rerun: see historical-measurement.txt, extracted from accepted ADR-0255 at delivery base 317f54ef. It predates all Task commits. The TechSpec records the exact 1,180,217,344-byte, 832,085-event inventory and detailed copy measurements; the ADR independently retains the 1.1 GB / 832k / 313-Run inventory and rounded timing results. No live home was accessed. The copy measurement supports explicit conversion being slower than a bounded start and slicing incremental compaction. It cannot establish current live size.
