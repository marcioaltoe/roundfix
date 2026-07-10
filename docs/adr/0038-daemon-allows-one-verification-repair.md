# Daemon allows one Verification repair

ADR-0014 already makes the Daemon authoritative for Task Verification, but Agents still run the same commands and consume their full successful output. After initial Agent work, the Daemon therefore runs Verification, returns only Verification Feedback to the same Agent Session on failure, allows one repair, and settles the Work Item from the second verdict. Repeating until success and adding a configurable attempt count were rejected to keep resource use bounded and the failure policy deterministic.
