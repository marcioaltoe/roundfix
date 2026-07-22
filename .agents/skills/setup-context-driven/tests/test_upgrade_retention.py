"""Upgrade Retention Contract behavior for setup-context-driven.

Suite: spec 0044 upgrade retention
Invariant: a baseline upgrade is authorizable only when every prior clause has ordered, equivalent accounting in the selected future artifact graph.
Boundary IN: context_setup.py audit/apply commands, Setup Manifest identities, transition assets, and repository bytes.
Boundary OUT: asset-schema mutation diagnostics owned by test_upgrade_contracts.py and final repository Verification.
"""

import io
import json
import shutil
import subprocess
import tempfile
import unittest
from contextlib import redirect_stderr, redirect_stdout
from dataclasses import replace
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
FIXTURES = Path(__file__).resolve().parent / "fixtures" / "upgrade-contracts"

import sys

sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

import context_setup  # noqa: E402
from context_assets import load_asset_catalog  # noqa: E402
from test_apply import (  # noqa: E402
    run_apply_preview,
    run_apply_with_confirmation,
    run_context_setup,
)
from test_audit import snapshot_files  # noqa: E402


PROFILE = "typescript-bun-monorepo"
TRANSITION_ID = "transition.legacy-typescript-bun-to-portable-v3"


class UpgradeRetentionTests(unittest.TestCase):
    def test_complete_transition_has_identical_ordered_text_and_json_accounting(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))

            json_result = run_apply_preview(repo, PROFILE, [])
            text_result = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "text",
                "--profile",
                PROFILE,
            )

            self.assertEqual(json_result.returncode, 3, json_result.stderr)
            self.assertEqual(text_result.returncode, 3, text_result.stderr)
            accounting = json.loads(json_result.stdout)["retentionAccounting"]
            self.assertEqual(
                {entry["disposition"] for entry in accounting},
                {"retained", "moved", "replaced", "rejected"},
            )
            self.assertTrue(all(entry["reason"] for entry in accounting))
            for entry in accounting:
                self.assertIn(f"reason: {entry['reason']}", text_result.stdout)
            clauses = [entry["fromClause"] for entry in accounting]
            self.assertEqual(clauses, sorted(clauses))
            self.assertEqual(
                [text_result.stdout.index(clause) for clause in clauses],
                sorted(text_result.stdout.index(clause) for clause in clauses),
            )

    def test_unknown_legacy_fingerprint_blocks_preview_and_apply_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            manifest_path = repo / context_setup.MANIFEST_PATH
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["managedArtifacts"][0]["digest"] = "2" * 64
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
            before = snapshot_files(repo)

            preview = run_context_setup(
                "audit", "--repo", str(repo), "--format", "json"
            )
            apply = run_apply_preview(repo, PROFILE, [])

            self.assertEqual(preview.returncode, 1, preview.stderr)
            self.assertEqual(apply.returncode, 1, apply.stderr)
            self.assertFinding(preview, "retention.baseline.unknown")
            self.assertFinding(apply, "retention.baseline.unknown")
            self.assertEqual(snapshot_files(repo), before)

    def test_unaccounted_clause_blocks_preview_and_apply_without_writes(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        transition = catalog.upgrade_transitions[TRANSITION_ID]
        missing_clause = transition.prior_clauses[0].clause_id
        incomplete = replace(
            transition,
            mappings=tuple(
                mapping
                for mapping in transition.mappings
                if mapping.from_clause != missing_clause
            ),
        )
        catalog = replace(
            catalog,
            upgrade_transitions={
                **catalog.upgrade_transitions,
                TRANSITION_ID: incomplete,
            },
        )

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            before = snapshot_files(repo)

            preview = run_with_catalog(
                catalog, "audit", "--repo", str(repo), "--format", "json"
            )
            apply = run_with_catalog(
                catalog,
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                PROFILE,
            )

            self.assertEqual(preview.returncode, 1, preview.stderr)
            self.assertEqual(apply.returncode, 1, apply.stderr)
            self.assertFinding(preview, "retention.clause.unaccounted", missing_clause)
            self.assertFinding(apply, "retention.clause.unaccounted", missing_clause)
            self.assertEqual(snapshot_files(repo), before)

    def test_weaker_target_and_missing_selected_carrier_block(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        transition = catalog.upgrade_transitions[TRANSITION_ID]
        accepted = next(
            mapping
            for mapping in transition.mappings
            if mapping.disposition != "rejected" and mapping.targets
        )
        target = accepted.targets[0]
        rules = dict(catalog.rule_contracts)
        for rule_id, rule in rules.items():
            if target not in {clause.clause_id for clause in rule.clauses}:
                continue
            rules[rule_id] = replace(
                rule,
                clauses=tuple(
                    replace(clause, enforcement="prohibited")
                    if clause.clause_id == target
                    else clause
                    for clause in rule.clauses
                ),
            )
            break
        weaker_catalog = replace(catalog, rule_contracts=rules)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            before = snapshot_files(repo)

            weaker = run_with_catalog(
                weaker_catalog,
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                PROFILE,
            )
            unreachable = run_with_catalog(
                catalog,
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                "rust-cli",
            )

            self.assertEqual(weaker.returncode, 1, weaker.stderr)
            self.assertFinding(weaker, "retention.target.enforcement-mismatch")
            self.assertEqual(unreachable.returncode, 1, unreachable.stderr)
            self.assertFinding(unreachable, "retention.target.unreachable")
            self.assertEqual(snapshot_files(repo), before)

    def test_retention_only_change_invalidates_old_confirmation(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            before = snapshot_files(repo)
            original, _, _ = context_setup.plan_apply(repo, catalog, PROFILE, [])
            revised_entries = (
                replace(
                    original.retention[0],
                    reason=original.retention[0].reason + " Reviewed.",
                ),
                *original.retention[1:],
            )
            with mock.patch.object(
                context_setup,
                "evaluate_retention",
                return_value=(revised_entries, []),
            ):
                revised, _, _ = context_setup.plan_apply(repo, catalog, PROFILE, [])

            self.assertEqual(original.findings, [])
            self.assertEqual(revised.findings, [])
            self.assertNotEqual(original.plan_digest, revised.plan_digest)

            with mock.patch.object(
                context_setup,
                "evaluate_retention",
                return_value=(revised_entries, []),
            ):
                stale = run_with_catalog(
                    catalog,
                    "apply",
                    "--repo",
                    str(repo),
                    "--format",
                    "json",
                    "--profile",
                    PROFILE,
                    "--confirm-plan",
                    original.plan_digest,
                )

            self.assertEqual(stale.returncode, 3, stale.stderr)
            self.assertFinding(stale, "plan.confirmation.stale")
            self.assertEqual(snapshot_files(repo), before)

    def test_confirmed_apply_records_current_baseline_and_repreview_is_empty(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            preview = run_apply_preview(repo, PROFILE, [])
            payload = json.loads(preview.stdout)

            applied = run_apply_with_confirmation(
                repo,
                PROFILE,
                [],
                payload["planDigest"],
            )
            manifest = json.loads(
                (repo / context_setup.MANIFEST_PATH).read_text(encoding="utf-8")
            )
            second = run_apply_preview(repo, PROFILE, [])

            self.assertEqual(preview.returncode, 3, preview.stderr)
            self.assertEqual(applied.returncode, 0, applied.stderr)
            self.assertEqual(
                json.loads(applied.stdout)["retentionAccounting"],
                payload["retentionAccounting"],
            )
            self.assertEqual(
                manifest["generator"]["baseline"],
                "baseline.portable-v3",
            )
            self.assertEqual(manifest["generator"]["fixtureMetadata"], "preserved")
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(json.loads(second.stdout)["plannedChanges"], [])

    def test_exit_code_meanings_remain_stable(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = write_legacy_fixture(Path(temp_dir))
            confirmation_required = run_apply_preview(repo, PROFILE, [])
            invalid = run_context_setup(
                "apply",
                "--repo",
                str(repo),
                "--format",
                "json",
                "--profile",
                PROFILE,
                "--confirm-plan",
                "invalid",
            )
            manifest_path = repo / context_setup.MANIFEST_PATH
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["managedArtifacts"][0]["digest"] = "2" * 64
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")
            blocking = run_apply_preview(repo, PROFILE, [])

            self.assertEqual(confirmation_required.returncode, 3)
            self.assertEqual(invalid.returncode, 2)
            self.assertEqual(blocking.returncode, 1)

    def assertFinding(self, result, code, managed_id=None):
        matches = [
            finding
            for finding in json.loads(result.stdout)["findings"]
            if finding["code"] == code
            and (managed_id is None or finding["managedId"] == managed_id)
        ]
        self.assertTrue(matches, result.stdout)


def write_legacy_fixture(repo: Path) -> Path:
    repo.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(FIXTURES / "pre-0.9-AGENTS.md", repo / "AGENTS.md")
    shutil.copyfile(FIXTURES / "pre-0.9-DESIGN.md", repo / "DESIGN.md")
    manifest_path = repo / context_setup.MANIFEST_PATH
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(FIXTURES / "legacy-manifest.json", manifest_path)
    return repo


def run_with_catalog(catalog, *args):
    stdout = io.StringIO()
    stderr = io.StringIO()
    with mock.patch.object(context_setup, "load_asset_catalog", return_value=catalog):
        with redirect_stdout(stdout), redirect_stderr(stderr):
            returncode = context_setup.main(list(args))
    return subprocess.CompletedProcess(
        args,
        returncode,
        stdout.getvalue(),
        stderr.getvalue(),
    )


if __name__ == "__main__":
    unittest.main()
