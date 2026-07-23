"""Characterization tests for the standalone Baseline parity corpus.

Suite: Python-to-Go compatibility corpus
Invariant: all maintained Python tests and transition states are represented
once by byte-stable, self-contained artifacts with explicit Go destinations.
Boundary IN: checked-in parity JSON, its generator, and both distributed skill
surfaces.
Boundary OUT: future Go parity execution, which consumes the same JSON corpus.
"""

from __future__ import annotations

import base64
import hashlib
import json
from pathlib import Path
import unittest

from parity_corpus import (
    ALLOWED_CLASSIFICATIONS,
    CORPUS_ROOT,
    FIXTURE_SCHEMA_VERSION,
    MATRIX_SCHEMA_VERSION,
    SCHEMA_VERSION,
    generate_corpus,
    python_test_inventory,
)
from test_support import setup_skill_roots


class ParityCorpusTests(unittest.TestCase):
    _generated = None

    @classmethod
    def generated(cls):
        if cls._generated is None:
            cls._generated = generate_corpus()
        return cls._generated

    def test_matrix_covers_every_python_test_once_with_go_destinations(self):
        matrix = self.read_json("matrix.json")
        self.assertEqual(matrix["schemaVersion"], MATRIX_SCHEMA_VERSION)
        source_rows = [
            row
            for row in matrix["rows"]
            if row["sourceKind"] == "python-test"
        ]
        expected = python_test_inventory()

        self.assertEqual(
            [row["id"] for row in source_rows],
            [row["id"] for row in expected],
        )
        self.assertEqual(len(source_rows), 240)
        self.assertEqual(len({row["id"] for row in matrix["rows"]}), matrix["rowCount"])
        self.assertEqual(matrix["classifications"], list(ALLOWED_CLASSIFICATIONS))
        self.assertGreater(matrix["classificationCounts"]["designed-delta"], 0)
        for row in matrix["rows"]:
            with self.subTest(row=row["id"]):
                self.assertIn(row["classification"], ALLOWED_CLASSIFICATIONS)
                self.assertTrue(row["goDestination"])
                self.assertTrue(row["fixtureIds"])
                self.assertTrue(row["contractDimensions"])
                if row["classification"] != "exact":
                    self.assertTrue(row["rationale"])
                if row["classification"] in {"designed-delta", "retired"}:
                    self.assertTrue(row["rationale"])

    def test_fixtures_cover_profiles_states_digests_and_rollback(self):
        manifest = self.read_json("manifest.json")
        self.assertEqual(manifest["schemaVersion"], SCHEMA_VERSION)
        self.assertEqual(
            manifest["profiles"],
            [
                "go-cli-tui",
                "rust-cli",
                "standard-typescript-monorepo",
            ],
        )
        fixtures = {
            fixture_id: self.read_json(f"fixtures/{fixture_id}.json")
            for fixture_id in manifest["fixtureIds"]
        }
        required_states = {
            "greenfield",
            "update",
            "preservation",
            "profile-change",
            "stale-input",
            "unsafe-carrier",
            "capability-refusal",
            "formatter-composition",
            "skill-restoration",
            "asset-sync",
            "rollback",
        }
        self.assertTrue(
            required_states.issubset(
                {fixture["adoptionState"] for fixture in fixtures.values()}
            )
        )
        profile_fixtures = {
            fixture["profile"]
            for fixture in fixtures.values()
            if fixture["adoptionState"] == "greenfield"
        }
        self.assertEqual(profile_fixtures, set(manifest["profiles"]))

        blobs_document = self.read_json("blobs.json")
        blobs = blobs_document["blobs"]
        for fixture_id, fixture in fixtures.items():
            with self.subTest(fixture=fixture_id):
                self.assertEqual(
                    fixture["schemaVersion"], FIXTURE_SCHEMA_VERSION
                )
                for field in (
                    "input",
                    "normalizedOutput",
                    "repositoryPreimage",
                    "fileIdentities",
                    "managedEntryLedger",
                    "manifest",
                    "planDigest",
                    "plannedByteSequence",
                    "refusal",
                    "postState",
                    "rollback",
                ):
                    self.assertIn(field, fixture)
                for record in [
                    *fixture["repositoryPreimage"],
                    *fixture["postState"],
                ]:
                    if record["kind"] != "file":
                        continue
                    encoded = blobs[record["identity"]]
                    content = base64.b64decode(encoded)
                    self.assertEqual(record["size"], len(content))
                    self.assertEqual(
                        record["identity"],
                        "sha256:" + hashlib.sha256(content).hexdigest(),
                    )

        exact_digest_fixtures = (
            "greenfield-rust-cli",
            "greenfield-go-cli-tui",
            "greenfield-standard-typescript-monorepo",
            "readoption-preservation",
            "stale-plan-refusal",
            "atomic-rollback",
            "skill-restoration",
            "skill-restoration-rollback",
        )
        for fixture_id in exact_digest_fixtures:
            self.assertRegex(fixtures[fixture_id]["planDigest"], r"^[0-9a-f]{64}$")
        for fixture_id in ("atomic-rollback", "skill-restoration-rollback"):
            rollback = fixtures[fixture_id]["rollback"]
            self.assertTrue(rollback["attempted"])
            self.assertTrue(rollback["succeeded"])
            self.assertTrue(rollback["restoredPreimage"])
        self.assertTrue(fixtures["readoption-preservation"]["managedEntryLedger"])
        self.assertTrue(fixtures["readoption-preservation"]["plannedByteSequence"])
        self.assertIsNotNone(fixtures["readoption-preservation"]["manifest"])

    def test_regeneration_is_byte_identical_and_manifest_hashes_every_artifact(self):
        generated = self.generated()
        actual = {
            path.relative_to(CORPUS_ROOT).as_posix(): path.read_bytes()
            for path in sorted(CORPUS_ROOT.rglob("*"))
            if path.is_file()
        }
        self.assertEqual(set(actual), set(generated))
        for path, expected in generated.items():
            with self.subTest(path=path):
                self.assertEqual(actual[path], expected)

        manifest = self.read_json("manifest.json")
        for artifact in manifest["artifacts"]:
            content = actual[artifact["path"]]
            self.assertEqual(artifact["bytes"], len(content))
            self.assertEqual(
                artifact["sha256"], hashlib.sha256(content).hexdigest()
            )

    def test_canonical_and_distributed_suites_validate_identical_corpus(self):
        canonical, distributed = setup_skill_roots(Path(__file__))
        canonical_root = canonical / "assets" / "parity-corpus" / "v1"
        distributed_root = distributed / "assets" / "parity-corpus" / "v1"
        canonical_files = {
            path.relative_to(canonical_root).as_posix(): path.read_bytes()
            for path in sorted(canonical_root.rglob("*"))
            if path.is_file()
        }
        distributed_files = {
            path.relative_to(distributed_root).as_posix(): path.read_bytes()
            for path in sorted(distributed_root.rglob("*"))
            if path.is_file()
        }
        self.assertEqual(distributed_files, canonical_files)

    @staticmethod
    def read_json(relative):
        return json.loads((CORPUS_ROOT / relative).read_text(encoding="utf-8"))


if __name__ == "__main__":
    unittest.main()
