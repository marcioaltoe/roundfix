"""Behavior tests for byte-exhaustive Baseline Readoption inventory.

Suite: incompatible Source Baseline inventory
Invariant: every byte in every bounded carrier appears in exactly one stable entry.
Boundary IN: local carrier discovery, structural partitioning, and audit rendering.
Boundary OUT: Readoption classification, dispositions, planning, and apply.
"""

import json
import os
import tempfile
import unittest
from pathlib import Path

from test_audit import run_audit, snapshot_files

from context_baseline import (  # noqa: E402
    SourceInventoryError,
    SourceInventoryLimits,
    inventory_incompatible_source_baseline,
)


class SourceInventoryTests(unittest.TestCase):
    def test_inventory_is_byte_exhaustive_structural_and_stable(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            self.write_incompatible_repository(repo)
            (repo / "packages" / "api").mkdir(parents=True)
            (repo / "packages" / "api" / "CLAUDE.md").write_bytes(
                b"# Nested\r\n\r\nKeep this byte-for-byte.\r\n"
            )
            (repo / "docs" / "agents" / "opaque.bin").write_bytes(b"\x00\xff\n")

            first = inventory_incompatible_source_baseline(
                repo, "baseline.pre-0.0.1"
            )
            second = inventory_incompatible_source_baseline(
                repo, "baseline.pre-0.0.1"
            )

            self.assertEqual(first, second)
            self.assertEqual(
                {entry.kind for entry in first.entries},
                {"file", "managed-block", "manifest-record", "unmarked-span"},
            )
            self.assertEqual(
                list(first.entries),
                sorted(
                    first.entries,
                    key=lambda entry: (
                        entry.path.as_posix(),
                        entry.start,
                        entry.end,
                        entry.kind,
                    ),
                ),
            )
            for path in first.carriers:
                entries = [entry for entry in first.entries if entry.path == path]
                source = (repo / path).read_bytes()
                self.assertEqual(
                    b"".join(entry.source_bytes for entry in entries),
                    source,
                    path,
                )
                self.assertEqual(
                    [(entry.start, entry.end) for entry in entries],
                    self.contiguous_ranges(entries, len(source)),
                )
                self.assertTrue(all(entry.carrier_digest for entry in entries))
                self.assertTrue(
                    all(entry.entry_id.startswith("source-entry.") for entry in entries)
                )

            unmarked = [entry for entry in first.entries if entry.kind == "unmarked-span"]
            self.assertTrue(any(b"Repository-owned" in entry.source_bytes for entry in unmarked))
            self.assertTrue(all("disposition" not in entry.to_json() for entry in first.entries))

    def test_inventory_ignores_declared_trees_and_rejects_unsafe_carriers(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            repo = root / "repo"
            self.write_incompatible_repository(repo)
            ignored = repo / "node_modules" / "dependency"
            ignored.mkdir(parents=True)
            (ignored / "AGENTS.md").write_text("ignored\n", encoding="utf-8")

            inventory = inventory_incompatible_source_baseline(
                repo, "baseline.pre-0.0.1"
            )
            self.assertNotIn(
                "node_modules/dependency/AGENTS.md",
                {path.as_posix() for path in inventory.carriers},
            )

            outside = root / "outside.md"
            outside.write_text("outside\n", encoding="utf-8")
            unsafe = repo / "docs" / "agents" / "linked.md"
            unsafe.symlink_to(outside)
            with self.assertRaises(SourceInventoryError) as captured:
                inventory_incompatible_source_baseline(repo, "baseline.pre-0.0.1")
            self.assertDiagnostic(captured.exception, "source-inventory.carrier.symlink")
            unsafe.unlink()

            fifo = repo / "docs" / "agents" / "special"
            os.mkfifo(fifo)
            with self.assertRaises(SourceInventoryError) as captured:
                inventory_incompatible_source_baseline(repo, "baseline.pre-0.0.1")
            self.assertDiagnostic(captured.exception, "source-inventory.carrier.special")
            fifo.unlink()

            limits = SourceInventoryLimits(max_file_bytes=8, max_total_bytes=64)
            with self.assertRaises(SourceInventoryError) as captured:
                inventory_incompatible_source_baseline(
                    repo, "baseline.pre-0.0.1", limits=limits
                )
            self.assertDiagnostic(captured.exception, "source-inventory.limit.file-bytes")

    def test_inventory_rejects_unsafe_manifest_record_paths(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            self.write_incompatible_repository(repo)
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["managedArtifacts"][0]["path"] = "../escape.md"
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")

            with self.assertRaises(SourceInventoryError) as captured:
                inventory_incompatible_source_baseline(repo, "baseline.pre-0.0.1")

            self.assertDiagnostic(captured.exception, "source-inventory.manifest.path.unsafe")

    def test_audit_reports_complete_unresolved_inventory_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            self.write_incompatible_repository(repo)
            before = snapshot_files(repo)

            first = run_audit(repo, "--format", "json")
            second = run_audit(repo, "--format", "json")
            text = run_audit(repo, "--format", "text")

            self.assertEqual(first.returncode, 3, first.stderr)
            self.assertEqual(second.returncode, 3, second.stderr)
            self.assertEqual(text.returncode, 3, text.stderr)
            self.assertEqual(json.loads(first.stdout), json.loads(second.stdout))
            payload = json.loads(first.stdout)
            self.assertEqual(payload["sourceBaseline"]["compatibility"], "incompatible")
            self.assertEqual(
                payload["sourceBaseline"]["entryCount"], len(payload["sourceEntries"])
            )
            self.assertTrue(payload["sourceEntries"])
            self.assertTrue(
                all("disposition" not in entry for entry in payload["sourceEntries"])
            )
            self.assertEqual(
                sum(
                    entry["end"] - entry["start"]
                    for entry in payload["sourceEntries"]
                ),
                payload["sourceBaseline"]["byteCount"],
            )
            self.assertIn("source baseline:", text.stdout.lower())
            self.assertIn("source entries:", text.stdout.lower())
            self.assertEqual(snapshot_files(repo), before)

    @staticmethod
    def write_incompatible_repository(repo):
        agents = repo / "docs" / "agents"
        agents.mkdir(parents=True)
        (repo / "AGENTS.md").write_bytes(
            b"# Repository-owned\n\n"
            b"<!-- setup-context-driven:begin id=root.old version=3 -->\n"
            b"Managed bytes.\n"
            b"<!-- setup-context-driven:end id=root.old -->\n"
            b"\nTrailing repository bytes.\n"
        )
        (agents / "guide.md").write_bytes(b"# Guide\n\nUnmarked guidance.\n")
        manifest = {
            "schemaVersion": 1,
            "generator": {
                "skill": "setup-context-driven",
                "version": 3,
                "baseline": "baseline.pre-0.0.1",
            },
            "profile": "old-profile",
            "modules": ["old"],
            "decisions": {"verification.gate": {"value": "make verify"}},
            "managedArtifacts": [
                {
                    "id": "root.old",
                    "path": "AGENTS.md",
                    "kind": "block",
                    "module": "old",
                    "template": "old",
                    "version": 3,
                    "digest": "0" * 64,
                }
            ],
            "localSkills": [],
        }
        (agents / "setup-context.json").write_text(
            json.dumps(manifest, indent=2) + "\n", encoding="utf-8"
        )

    @staticmethod
    def contiguous_ranges(entries, size):
        cursor = 0
        ranges = []
        for entry in entries:
            if entry.start != cursor:
                return []
            ranges.append((entry.start, entry.end))
            cursor = entry.end
        return ranges if cursor == size else []

    def assertDiagnostic(self, error, code):
        self.assertIn(code, [diagnostic.code for diagnostic in error.diagnostics])


if __name__ == "__main__":
    unittest.main()
