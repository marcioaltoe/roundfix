# Suite: Source Baseline contracts
# Invariant: corpus, manifest, identity, and independent index describe one exact 0.0.1 source.
# Boundary IN: local Source Baseline parsing, normalization, and integrity validation.
# Boundary OUT: Baseline Readoption inventory, planning, writes, subprocesses, and network access.

import json
import shutil
import subprocess
import sys
import tempfile
import unittest
import urllib.request
from dataclasses import FrozenInstanceError
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
FIXTURE_ROOT = Path(__file__).resolve().parent / "fixtures" / "source-baselines"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_baseline import (  # noqa: E402
    SOURCE_BASELINE_VERSION,
    SourceBaselineValidationError,
    load_source_baseline,
    load_source_baselines,
)


class SourceBaselineContractTests(unittest.TestCase):
    def test_valid_source_baseline_loads_deterministically(self):
        first = load_source_baselines(FIXTURE_ROOT)
        second = load_source_baselines(FIXTURE_ROOT)

        self.assertEqual(first, second)
        self.assertEqual(len(first), 1)
        baseline = first[0]
        self.assertEqual(baseline.version, SOURCE_BASELINE_VERSION)
        self.assertEqual(
            tuple(entry.entry_id for entry in baseline.entries),
            ("clause.root-cause", "contract.verification"),
        )
        self.assertEqual(
            baseline.entries[0].path,
            Path("corpus/AGENTS.md"),
        )
        with self.assertRaises(FrozenInstanceError):
            baseline.version = "changed"

    def test_one_source_baseline_loads_by_stable_id(self):
        baseline = load_source_baseline(FIXTURE_ROOT, "fixture")

        self.assertEqual(baseline.baseline_id, "fixture")

    def test_unknown_source_baseline_id_is_rejected(self):
        with self.assertRaises(SourceBaselineValidationError) as captured:
            load_source_baseline(FIXTURE_ROOT, "missing")

        self.assertIn("source-baseline.id.unknown", str(captured.exception))

    def test_schema_and_version_mutations_fail_with_stable_diagnostics(self):
        cases = [
            ("index schema", self._malformed_index_schema, "source-baseline.schema.invalid"),
            ("identity schema", self._malformed_identity_schema, "source-baseline.schema.invalid"),
            ("manifest schema", self._malformed_manifest_schema, "source-baseline.schema.invalid"),
            (
                "integer index version",
                self._integer_index_version,
                "source-baseline.version.invalid",
            ),
            (
                "integer identity version",
                self._integer_identity_version,
                "source-baseline.version.invalid",
            ),
            (
                "integer manifest version",
                self._integer_manifest_version,
                "source-baseline.version.invalid",
            ),
            ("mixed generation", self._mixed_manifest_version, "source-baseline.version.mismatch"),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                self._assert_stable_failure(mutator, expected_code)

    def test_membership_mutations_fail_with_stable_diagnostics(self):
        cases = [
            (
                "corpus entry absent from manifest",
                self._remove_manifest_entry,
                "source-baseline.manifest.entry.missing",
            ),
            (
                "manifest entry absent from corpus",
                self._remove_corpus_entry,
                "source-baseline.corpus.entry.missing",
            ),
            (
                "manifest entry absent from index",
                self._remove_index_entry,
                "source-baseline.index.entry.missing",
            ),
            (
                "entry removed from corpus and manifest remains required by index",
                self._remove_corpus_and_manifest_entry,
                "source-baseline.manifest.entry.missing",
            ),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                self._assert_stable_failure(mutator, expected_code)

    def test_structural_and_integrity_mutations_fail_with_stable_diagnostics(self):
        cases = [
            ("unsafe entry path", self._unsafe_entry_path, "source-baseline.entry.path.invalid"),
            (
                "duplicate manifest id",
                self._duplicate_manifest_id,
                "source-baseline.entry.id.duplicate",
            ),
            ("duplicate index id", self._duplicate_index_id, "source-baseline.entry.id.duplicate"),
            (
                "unknown identity field",
                self._unknown_identity_field,
                "source-baseline.identity.fields.invalid",
            ),
            (
                "invalid entry collection",
                self._invalid_entry_collection,
                "source-baseline.manifest.document.invalid",
            ),
            ("invalid entry kind", self._invalid_entry_kind, "source-baseline.entry.kind.invalid"),
            ("wrong byte range", self._wrong_entry_range, "source-baseline.entry.range.mismatch"),
            (
                "wrong entry digest",
                self._wrong_entry_digest,
                "source-baseline.entry.digest.mismatch",
            ),
            (
                "wrong identity count",
                self._wrong_identity_count,
                "source-baseline.entry.count.mismatch",
            ),
            ("wrong index count", self._wrong_index_count, "source-baseline.entry.count.mismatch"),
            (
                "wrong corpus digest",
                self._wrong_corpus_digest,
                "source-baseline.corpus.digest.mismatch",
            ),
            (
                "wrong manifest digest",
                self._wrong_manifest_digest,
                "source-baseline.manifest.digest.mismatch",
            ),
        ]

        for name, mutator, expected_code in cases:
            with self.subTest(name=name):
                self._assert_stable_failure(mutator, expected_code)

    def test_source_baseline_loading_has_no_write_command_or_network_side_effects(self):
        before = self._tree_bytes(FIXTURE_ROOT)

        with (
            mock.patch.object(Path, "write_text", side_effect=AssertionError("write attempted")),
            mock.patch.object(Path, "write_bytes", side_effect=AssertionError("write attempted")),
            mock.patch.object(Path, "mkdir", side_effect=AssertionError("write attempted")),
            mock.patch.object(Path, "unlink", side_effect=AssertionError("write attempted")),
            mock.patch.object(Path, "rename", side_effect=AssertionError("write attempted")),
            mock.patch.object(subprocess, "run", side_effect=AssertionError("command attempted")),
            mock.patch.object(
                urllib.request,
                "urlopen",
                side_effect=AssertionError("network attempted"),
            ),
        ):
            load_source_baselines(FIXTURE_ROOT)

        self.assertEqual(self._tree_bytes(FIXTURE_ROOT), before)

    def _assert_stable_failure(self, mutator, expected_code):
        first = self._load_invalid(mutator)
        second = self._load_invalid(mutator)

        self.assertEqual(first, second)
        self.assertTrue(
            any(expected_code in diagnostic for diagnostic in first),
            f"expected {expected_code!r} in {first!r}",
        )

    def _load_invalid(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            shutil.copytree(FIXTURE_ROOT, temp_root)
            mutator(temp_root)
            with self.assertRaises(SourceBaselineValidationError) as captured:
                load_source_baselines(temp_root)
        return captured.exception.diagnostics

    def _tree_bytes(self, root):
        return {
            path.relative_to(root).as_posix(): path.read_bytes()
            for path in sorted(root.rglob("*"))
            if path.is_file()
        }

    def _index_path(self, root):
        return root / "assets" / "source-baselines" / "index.json"

    def _identity_path(self, root):
        return root / "assets" / "source-baselines" / "fixture" / "baseline.json"

    def _manifest_path(self, root):
        return root / "assets" / "source-baselines" / "fixture" / "manifest.json"

    def _corpus_path(self, root):
        return root / "assets" / "source-baselines" / "fixture" / "corpus" / "AGENTS.md"

    def _read_json(self, path):
        return json.loads(path.read_text(encoding="utf-8"))

    def _write_json(self, path, value):
        path.write_text(json.dumps(value, indent=2) + "\n", encoding="utf-8")

    def _mutate_json(self, path, mutation):
        document = self._read_json(path)
        mutation(document)
        self._write_json(path, document)

    def _malformed_index_schema(self, root):
        self._mutate_json(
            self._index_path(root),
            lambda document: document.update(schemaVersion="wrong"),
        )

    def _malformed_identity_schema(self, root):
        self._mutate_json(
            self._identity_path(root),
            lambda document: document.update(schemaVersion="wrong"),
        )

    def _malformed_manifest_schema(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document.update(schemaVersion="wrong"),
        )

    def _integer_index_version(self, root):
        self._mutate_json(self._index_path(root), lambda document: document.update(version=1))

    def _integer_identity_version(self, root):
        self._mutate_json(self._identity_path(root), lambda document: document.update(version=1))

    def _integer_manifest_version(self, root):
        self._mutate_json(self._manifest_path(root), lambda document: document.update(version=1))

    def _mixed_manifest_version(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document.update(version="0.0.2"),
        )

    def _remove_manifest_entry(self, root):
        self._mutate_json(self._manifest_path(root), lambda document: document["entries"].pop())

    def _remove_corpus_entry(self, root):
        path = self._corpus_path(root)
        content = path.read_text(encoding="utf-8")
        start = content.index("<!-- source-baseline-entry: contract.verification -->")
        path.write_text(content[:start].rstrip() + "\n", encoding="utf-8")

    def _remove_index_entry(self, root):
        self._mutate_json(
            self._index_path(root),
            lambda document: document["baselines"][0]["entryIds"].pop(),
        )

    def _remove_corpus_and_manifest_entry(self, root):
        self._remove_corpus_entry(root)
        self._remove_manifest_entry(root)

    def _unsafe_entry_path(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document["entries"][0].update(path="../AGENTS.md"),
        )

    def _duplicate_manifest_id(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document["entries"].append(document["entries"][0]),
        )

    def _duplicate_index_id(self, root):
        self._mutate_json(
            self._index_path(root),
            lambda document: document["baselines"][0]["entryIds"].append(
                document["baselines"][0]["entryIds"][0]
            ),
        )

    def _unknown_identity_field(self, root):
        self._mutate_json(self._identity_path(root), lambda document: document.update(extra=True))

    def _invalid_entry_collection(self, root):
        self._mutate_json(self._manifest_path(root), lambda document: document.update(entries={}))

    def _invalid_entry_kind(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document["entries"][0].update(kind="summary"),
        )

    def _wrong_entry_range(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document["entries"][0].update(start=51),
        )

    def _wrong_entry_digest(self, root):
        self._mutate_json(
            self._manifest_path(root),
            lambda document: document["entries"][0].update(digest="0" * 64),
        )

    def _wrong_identity_count(self, root):
        self._mutate_json(self._identity_path(root), lambda document: document.update(entryCount=3))

    def _wrong_index_count(self, root):
        self._mutate_json(
            self._index_path(root),
            lambda document: document["baselines"][0].update(entryCount=3),
        )

    def _wrong_corpus_digest(self, root):
        self._mutate_json(
            self._identity_path(root),
            lambda document: document.update(corpusDigest="0" * 64),
        )

    def _wrong_manifest_digest(self, root):
        self._mutate_json(
            self._index_path(root),
            lambda document: document["baselines"][0].update(manifestDigest="0" * 64),
        )


if __name__ == "__main__":
    unittest.main()
