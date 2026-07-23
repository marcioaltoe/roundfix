# Suite: upgrade compatibility asset contracts
# Invariant: supported upgrade assets normalize deterministically or fail locally with stable diagnostics.
# Boundary IN: versioned asset loading, normalization, and semantic validation.
# Boundary OUT: repository inspection, planning, rendering, process execution, and network access.

import dataclasses
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
FIXTURE_ROOT = Path(__file__).resolve().parent / "fixtures" / "upgrade-contracts"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_support import setup_skill_roots  # noqa: E402

from context_assets import (  # noqa: E402
    AssetValidationError,
    clone_assets_to,
    load_asset_catalog,
    read_json_copy,
    write_json,
)


class UpgradeContractTests(unittest.TestCase):
    def test_valid_fixture_loads_immutable_normalized_values_deterministically(self):
        first = load_asset_catalog(FIXTURE_ROOT)
        second = load_asset_catalog(FIXTURE_ROOT)

        self.assertEqual(first, second)
        rule = first.rule_contracts["rule.core.upgrade-safety"]
        self.assertEqual(
            tuple(clause.clause_id for clause in rule.clauses),
            (
                "clause.core.confirm-version-transition",
                "clause.core.retain-hard-rules",
            ),
        )
        with self.assertRaises(dataclasses.FrozenInstanceError):
            rule.clauses[0].guidance = "changed"

        transition = first.upgrade_transitions["transition.fixture-v1-to-v2"]
        self.assertEqual(
            transition.legacy_manifest_fingerprints,
            ("c" * 64,),
        )
        self.assertEqual(
            tuple(mapping.from_clause for mapping in transition.mappings),
            (
                "clause.legacy.project-rules",
                "clause.legacy.retain-hard-rules",
            ),
        )
        self.assertEqual(
            first.formatter_by_profile["fixture"].fixture_paths,
            (Path("tests/fixtures/formatter/input.md"),),
        )
        self.assertEqual(
            first.repository_extensions["extension.repository-rules"].target_path,
            Path("docs/agents/repository.md"),
        )
        self.assertEqual(
            first.coverage_contracts["coverage.upgrade-safety"].delegation_aliases,
            ("hard rules", "upgrade safety"),
        )
        self.assertEqual(
            tuple(entry.trigger_id for entry in first.skill_dispatch_by_skill["coding-guidelines"]),
            ("trigger.implementation.change", "trigger.implementation.review"),
        )

    def test_duplicate_and_enum_mutations_fail_with_stable_diagnostics(self):
        cases = [
            ("duplicate clause", self._duplicate_clause, "clause.id.duplicate"),
            ("duplicate trigger", self._duplicate_trigger, "skill.dispatch.trigger.id.duplicate"),
            ("invalid enforcement", self._invalid_enforcement, "clause.enforcement.invalid"),
            ("invalid disposition", self._invalid_disposition, "transition.disposition.invalid"),
        ]

        for name, mutator, diagnostic in cases:
            with self.subTest(name=name):
                self.assertIn(diagnostic, self._invalid_diagnostics(mutator))

    def test_unsafe_malformed_and_incomplete_mutations_fail_locally(self):
        cases = [
            ("unsafe extension path", self._unsafe_extension_path, "extension.path.invalid"),
            ("malformed formatter", self._malformed_formatter, "profile.formatter.version.invalid"),
            ("incomplete transition", self._incomplete_transition, "transition.mapping.incomplete"),
            ("invalid delegation alias", self._invalid_delegation_alias, "coverage.delegationAliases.invalid"),
            ("missing clause guidance", self._missing_clause_guidance, "clause.field.missing"),
            ("unknown transition target", self._unknown_transition_target, "transition.target.unknown"),
            ("invalid trigger structure", self._invalid_trigger_structure, "skill.dispatch.triggers.invalid"),
            (
                "invalid legacy manifest fingerprint",
                self._invalid_legacy_manifest_fingerprint,
                "transition.legacyManifestFingerprint.invalid",
            ),
        ]

        for name, mutator, diagnostic in cases:
            with self.subTest(name=name):
                with (
                    mock.patch("subprocess.run", side_effect=AssertionError("command attempted")),
                    mock.patch("urllib.request.urlopen", side_effect=AssertionError("network attempted")),
                ):
                    self.assertIn(diagnostic, self._invalid_diagnostics(mutator))

    def test_canonical_catalog_loads_reviewed_ledgers_and_v2_stays_compatible(self):
        canonical = load_asset_catalog(SKILL_ROOT)
        v2 = load_asset_catalog(
            Path(__file__).resolve().parent / "fixtures" / "asset-contracts-v2"
        )

        self.assertEqual(
            set(canonical.upgrade_transitions),
            {
                "transition.managed-v2-to-portable-v3",
                "transition.legacy-typescript-bun-to-portable-v3",
            },
        )
        self.assertEqual(
            set(canonical.repository_extensions),
            {"extension.repository-rules"},
        )
        self.assertEqual(
            {
                profile_id: formatter.kind
                for profile_id, formatter in canonical.formatter_by_profile.items()
            },
            {
                "go-cli-tui": "none",
                "rust-cli": "none",
                "standard-typescript-monorepo": "selected",
                "typescript-bun-monorepo": "selected",
            },
        )
        self.assertEqual(v2.upgrade_transitions, {})
        self.assertEqual(v2.repository_extensions, {})
        self.assertEqual(v2.formatter_by_profile, {})

    def test_loading_is_read_only_process_free_and_network_free(self):
        before = self._asset_bytes(FIXTURE_ROOT)

        with (
            mock.patch.object(Path, "write_text", side_effect=AssertionError("write attempted")),
            mock.patch("subprocess.run", side_effect=AssertionError("command attempted")),
            mock.patch("urllib.request.urlopen", side_effect=AssertionError("network attempted")),
        ):
            load_asset_catalog(FIXTURE_ROOT)

        self.assertEqual(self._asset_bytes(FIXTURE_ROOT), before)

    def test_canonical_and_distributed_setup_skill_trees_are_byte_identical(self):
        canonical_root, distributed_root = setup_skill_roots(SKILL_ROOT)
        self.assertEqual(
            self._tree_bytes(canonical_root),
            self._tree_bytes(distributed_root),
        )

    def _invalid_diagnostics(self, mutator):
        with tempfile.TemporaryDirectory() as temp_dir:
            temp_root = Path(temp_dir) / "setup-context-driven"
            clone_assets_to(FIXTURE_ROOT, temp_root)
            mutator(temp_root)
            with self.assertRaises(AssetValidationError) as captured:
                load_asset_catalog(temp_root)
        return "\n".join(captured.exception.diagnostics)

    def _asset_bytes(self, root):
        return {
            path.relative_to(root).as_posix(): path.read_bytes()
            for path in sorted((root / "assets").rglob("*"))
            if path.is_file()
        }

    def _tree_bytes(self, root):
        return {
            path.relative_to(root).as_posix(): path.read_bytes()
            for path in sorted(root.rglob("*"))
            if path.is_file()
        }

    def _module_path(self, root):
        return root / "assets" / "modules" / "core.json"

    def _profile_path(self, root):
        return root / "assets" / "profiles" / "fixture.json"

    def _coverage_path(self, root):
        return root / "assets" / "coverage.json"

    def _transition_path(self, root):
        return root / "assets" / "retention" / "transition.fixture-v1-to-v2.json"

    def _mutate(self, path, mutation):
        data = read_json_copy(path)
        mutation(data)
        write_json(path, data)

    def _duplicate_clause(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["rules"][0]["clauses"].append(data["rules"][0]["clauses"][0]),
        )

    def _duplicate_trigger(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["skillDispatch"][0]["triggers"].append(
                data["skillDispatch"][0]["triggers"][0]
            ),
        )

    def _invalid_enforcement(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["rules"][0]["clauses"][0].update(enforcement="optional"),
        )

    def _invalid_disposition(self, root):
        self._mutate(
            self._transition_path(root),
            lambda data: data["mappings"][0].update(disposition="ignored"),
        )

    def _unsafe_extension_path(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["repositoryExtensions"][0].update(path="../repository.md"),
        )

    def _malformed_formatter(self, root):
        self._mutate(
            self._profile_path(root),
            lambda data: data["formatter"].update(version=""),
        )

    def _incomplete_transition(self, root):
        self._mutate(self._transition_path(root), lambda data: data["mappings"].pop())

    def _invalid_legacy_manifest_fingerprint(self, root):
        self._mutate(
            self._transition_path(root),
            lambda data: data.update(legacyManifestFingerprints=["nearest-match"]),
        )

    def _invalid_delegation_alias(self, root):
        self._mutate(
            self._coverage_path(root),
            lambda data: data["coverage"][0]["delegationAliases"].append("unsafe\nalias"),
        )

    def _missing_clause_guidance(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["rules"][0]["clauses"][0].pop("guidance"),
        )

    def _unknown_transition_target(self, root):
        self._mutate(
            self._transition_path(root),
            lambda data: data["mappings"][0].update(targets=["clause.missing"]),
        )

    def _invalid_trigger_structure(self, root):
        self._mutate(
            self._module_path(root),
            lambda data: data["skillDispatch"][0].update(triggers={}),
        )


if __name__ == "__main__":
    unittest.main()
