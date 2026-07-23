"""Alignment tests for maintained profiles and the Repository Skill Set.

Suite: maintained Context-Driven Baseline profiles
Invariant: every current profile uses generation 0.0.1, binds the universal
Repository Capabilities, and agrees exactly with its deterministic skill snapshot.
Boundary IN: profile assets, setup snapshots, activation contracts, and legacy
profile routing.
Boundary OUT: external Git acquisition and atomic restoration details
(test_restore_skills.py).
"""

import copy
import hashlib
import json
import os
import re
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]


def find_repository_root():
    for candidate in (SKILL_ROOT, *SKILL_ROOT.parents):
        if (candidate / "skills-lock.json").is_file() and (
            candidate / ".agents" / "skills"
        ).is_dir():
            return candidate
    raise RuntimeError("repository root is not discoverable from the skill tree")


REPOSITORY_ROOT = find_repository_root()
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    setup_snapshot_digest,
    write_json,
)
from context_capabilities import (  # noqa: E402
    EvidenceKind,
    RequirementStrength,
    UNIVERSAL_CAPABILITIES,
)
from test_source_inventory import SourceInventoryTests  # noqa: E402


MAINTAINED_PROFILES = {
    "go-cli-tui": "go-cli",
    "rust-cli": "rust-cli",
    "standard-typescript-monorepo": "typescript-bun",
}
OWNED_SKILLS = {
    "archive-spec",
    "brainstorming",
    "business-analyst",
    "council",
    "evidence-gate",
    "implement-spec",
    "implement-task",
    "qa-gate",
    "roundfix",
    "setup-context-driven",
    "write-idea",
    "write-prd",
    "write-tasks",
    "write-techspec",
}
IMMUTABLE_REF = re.compile(r"^(?:[0-9a-f]{40}|[0-9a-f]{64})$")
LOCK_FILE_SHA256 = "a76dbc94b59b8b05f87bf4d739e439ebfb9629f0b5137bc57f90ab648dbfeceb"
LOCK_FIXTURE_SHA256 = "4020da4a35468cafb9091dad688be99f4bec87385be348a37d06d2e5192e60ec"


class ProfileAlignmentTests(unittest.TestCase):
    def test_maintained_profiles_use_only_the_current_identity_and_generation(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        self.assertEqual(set(catalog.profiles), set(MAINTAINED_PROFILES))
        self.assertFalse(
            (SKILL_ROOT / "assets" / "profiles" / "typescript-bun-monorepo.json").exists()
        )
        for profile_id, setup_id in MAINTAINED_PROFILES.items():
            with self.subTest(profile=profile_id):
                profile = catalog.profiles[profile_id]
                self.assertEqual(
                    profile["schemaVersion"], "setup-context-driven/profile/0.0.1"
                )
                self.assertEqual(profile["version"], "0.0.1")
                self.assertEqual(profile["markerVersion"], "0.0.1")
                self.assertEqual(profile["setup"], setup_id)
                self.assertEqual(profile["capabilitySets"], ["universal"])

    def test_universal_capability_contract_is_exact_for_every_profile(self):
        expected = {
            "capability.context7": (
                RequirementStrength.REQUIRED,
                EvidenceKind.INSTALLED_SKILL,
            ),
            "capability.exa": (
                RequirementStrength.REQUIRED,
                EvidenceKind.INSTALLED_SKILL,
            ),
            "capability.firecrawl": (
                RequirementStrength.RECOMMENDED,
                EvidenceKind.INSTALLED_SKILL,
            ),
            "capability.rg": (
                RequirementStrength.RECOMMENDED,
                EvidenceKind.EXECUTABLE,
            ),
            "capability.rtk": (
                RequirementStrength.RECOMMENDED,
                EvidenceKind.EXECUTABLE,
            ),
        }
        actual = {
            item.capability_id: (item.strength, item.evidence_kind)
            for item in UNIVERSAL_CAPABILITIES
        }

        self.assertEqual(actual, expected)
        catalog = load_asset_catalog(SKILL_ROOT)
        self.assertTrue(
            all(
                catalog.profiles[profile_id]["capabilitySets"] == ["universal"]
                for profile_id in MAINTAINED_PROFILES
            )
        )

    def test_every_profile_snapshot_matches_skills_activations_and_trusted_bytes(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        for profile_id, setup_id in MAINTAINED_PROFILES.items():
            with self.subTest(profile=profile_id):
                profile = catalog.profiles[profile_id]
                setup = catalog.setups[setup_id]
                self.assertEqual(
                    setup["schemaVersion"],
                    "setup-context-driven/setup-snapshot/0.0.1",
                )
                self.assertEqual(setup["version"], "0.0.1")
                names = [item["name"] for item in setup["skills"]]
                self.assertEqual(len(names), len(set(names)))
                required = {
                    skill
                    for module_id in catalog.ordered_modules_by_profile[profile_id]
                    for skill in catalog.modules[module_id]["requiredSkills"]
                }
                self.assertEqual(set(names), required)
                repo_owned = {
                    item["name"]
                    for item in setup["skills"]
                    if item["source"]["type"] == "repo"
                }
                self.assertEqual(repo_owned, OWNED_SKILLS)
                self.assertTrue(
                    {"context7", "exa-web-search", "firecrawl"}.issubset(names)
                )
                expected_bundles = [
                    {
                        "id": bundle.bundle_id,
                        "skills": list(bundle.skills),
                    }
                    for bundle in catalog.activation_bundles_by_setup[setup_id]
                ]
                self.assertEqual(setup["activationBundles"], expected_bundles)
                self.assertEqual(
                    setup["digest"],
                    setup_snapshot_digest(
                        setup["skills"], setup["activationBundles"]
                    ),
                )
                for skill in setup["skills"]:
                    if skill["source"]["type"] == "repo":
                        trusted = (
                            REPOSITORY_ROOT
                            / ".agents"
                            / "skills"
                            / skill["name"]
                            / "SKILL.md"
                        ).read_bytes()
                        self.assertEqual(
                            skill["contentDigest"], hashlib.sha256(trusted).hexdigest()
                        )
                    else:
                        self.assertEqual(
                            set(skill["source"]),
                            {"type", "repository", "ref", "path"},
                        )
                        self.assertRegex(skill["source"]["ref"], IMMUTABLE_REF)
                        self.assertRegex(skill["treeDigest"], r"^[0-9a-f]{64}$")

    def test_skill_snapshot_mutations_fail_closed(self):
        mutations = {
            "missing": self.remove_owned_skill,
            "unexpected": self.add_unexpected_skill,
            "duplicate": self.duplicate_owned_skill,
            "bundle mismatch": self.change_bundle_membership,
        }
        for name, mutate in mutations.items():
            with self.subTest(name=name), tempfile.TemporaryDirectory() as temp_dir:
                root = Path(temp_dir) / "setup-context-driven"
                clone_assets_to(SKILL_ROOT, root)
                setup_path = root / "assets" / "setups" / "go-cli.json"
                setup = json.loads(setup_path.read_text(encoding="utf-8"))
                mutate(setup)
                setup["digest"] = setup_snapshot_digest(
                    setup["skills"], setup["activationBundles"]
                )
                write_json(setup_path, setup)

                with self.assertRaises(AssetValidationError) as raised:
                    load_asset_catalog(root)

                diagnostics = "\n".join(raised.exception.diagnostics)
                self.assertRegex(
                    diagnostics,
                    r"profile\.skill\.set\.mismatch|setup\.skill\.name\.duplicate|setup\.activationBundle\.members\.mismatch",
                )

    def test_former_typescript_profile_is_rejected_and_readopted_as_evidence(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["profile"] = "typescript-bun-monorepo"
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")

            readoption = self.run_audit(repo, "standard-typescript-monorepo")
            rejected = self.run_audit(repo, "typescript-bun-monorepo")

            self.assertEqual(readoption.returncode, 3, readoption.stderr)
            payload = json.loads(readoption.stdout)
            self.assertEqual(
                payload["sourceBaseline"]["declaredIdentity"],
                "baseline.pre-0.0.1",
            )
            self.assertTrue(payload["sourceEntries"])
            self.assertEqual(rejected.returncode, 3, rejected.stderr)
            rejected_payload = json.loads(rejected.stdout)
            self.assertTrue(
                {"readoption.baseline.incompatible", "readoption.disposition.required"}
                .issubset(
                    {item["code"] for item in rejected_payload["findings"]}
                )
            )

    def test_external_lock_contract_fixtures_remain_byte_identical(self):
        lock_file = REPOSITORY_ROOT / "skills-lock.json"
        fixture = SKILL_ROOT / "assets" / "lock-hash-compatibility-v1.json"

        self.assertEqual(hashlib.sha256(lock_file.read_bytes()).hexdigest(), LOCK_FILE_SHA256)
        self.assertEqual(
            hashlib.sha256(fixture.read_bytes()).hexdigest(), LOCK_FIXTURE_SHA256
        )
        lock = json.loads(lock_file.read_text(encoding="utf-8"))
        self.assertEqual(set(lock), {"version", "skills"})
        self.assertEqual(lock["version"], 1)
        self.assertIsInstance(lock["skills"], dict)

    @staticmethod
    def remove_owned_skill(setup):
        setup["skills"] = [
            item
            for item in setup["skills"]
            if item["name"] != "setup-context-driven"
        ]

    @staticmethod
    def add_unexpected_skill(setup):
        item = copy.deepcopy(
            next(
                skill
                for skill in setup["skills"]
                if skill["name"] == "setup-context-driven"
            )
        )
        item["name"] = "unexpected-owned"
        item["path"] = "skills/unexpected-owned"
        setup["skills"].append(item)

    @staticmethod
    def duplicate_owned_skill(setup):
        setup["skills"].append(
            copy.deepcopy(
                next(
                    skill
                    for skill in setup["skills"]
                    if skill["name"] == "setup-context-driven"
                )
            )
        )

    @staticmethod
    def change_bundle_membership(setup):
        setup["activationBundles"][0]["skills"].append("evidence-gate")

    @staticmethod
    def run_audit(repo, profile):
        environment = os.environ.copy()
        environment["PYTHONDONTWRITEBYTECODE"] = "1"
        return subprocess.run(
            [
                sys.executable,
                str(SCRIPT),
                "audit",
                "--repo",
                str(repo),
                "--profile",
                profile,
                "--format",
                "json",
            ],
            text=True,
            capture_output=True,
            check=False,
            env=environment,
        )


if __name__ == "__main__":
    unittest.main()
