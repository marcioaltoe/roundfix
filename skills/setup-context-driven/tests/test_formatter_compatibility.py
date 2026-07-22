"""Formatter compatibility for generated setup Markdown.

Suite: Formatter-Stable Output composition
Invariant: the declared formatter leaves every generated Markdown byte unchanged.
Boundary IN: profile formatter contracts, setup apply/audit, and checked-in golden bytes.
Boundary OUT: downloading or installing Oxfmt; final QA owns the real formatter probe.
"""

import json
import os
import shlex
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPTS_ROOT = SKILL_ROOT / "scripts"
TESTS_ROOT = Path(__file__).resolve().parent
FORMATTER_FIXTURE_ROOT = (
    TESTS_ROOT
    / "fixtures"
    / "formatter-compatibility"
    / "typescript-bun-monorepo"
)
GOLDEN_PREFIX = Path(
    "tests/fixtures/formatter-compatibility/typescript-bun-monorepo/golden"
)
PROVENANCE_PATH = FORMATTER_FIXTURE_ROOT / "provenance.json"
VERIFICATION_SOURCE = FORMATTER_FIXTURE_ROOT / "verify_fixture.py"

sys.path.insert(0, str(SCRIPTS_ROOT))
sys.path.insert(0, str(TESTS_ROOT))

import context_setup  # noqa: E402
from context_assets import load_asset_catalog, portable_file_digest  # noqa: E402
from test_apply import BASE_DECISIONS, run_apply_in_process  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_audit as run_fixture_audit,
    snapshot_files,
)


PROFILE_ID = "typescript-bun-monorepo"


def formatter_profile_decisions():
    replacements = {
        "domain.layout": "multi-context",
        "triage.external": "true",
        "secondbrain.enabled": "true",
        "repository.extension.enabled": "true",
        "verification.gate": "python3 -B .formatter-fixture-verify.py",
    }
    decisions = []
    for decision in BASE_DECISIONS:
        decision_id, _, value = decision.partition("=")
        decisions.append(f"{decision_id}={replacements.get(decision_id, value)}")
    return decisions


def load_provenance():
    return json.loads(PROVENANCE_PATH.read_text(encoding="utf-8"))


def formatter_target_path(fixture_path):
    try:
        return fixture_path.relative_to(GOLDEN_PREFIX)
    except ValueError as error:
        raise AssertionError(
            f"formatter fixture path is outside the golden corpus: {fixture_path}"
        ) from error


def formatter_corpus_digest(contract):
    files = []
    for fixture_path in contract.fixture_paths:
        target_path = formatter_target_path(fixture_path)
        files.append(
            (
                target_path.as_posix().encode("utf-8"),
                (SKILL_ROOT / fixture_path).read_bytes(),
            )
        )
    return portable_file_digest(files)


def assert_profile_formatter_canonical(repo, catalog, profile_id):
    contract = catalog.formatter_by_profile[profile_id]
    if contract.kind == "none":
        return

    expected_paths = {
        formatter_target_path(fixture_path)
        for fixture_path in contract.fixture_paths
    }
    generated_paths = {
        path.relative_to(repo)
        for path in [
            repo / "AGENTS.md",
            *sorted((repo / "docs" / "agents").glob("*.md")),
        ]
        if path.is_file()
    }
    if generated_paths != expected_paths:
        missing = sorted(path.as_posix() for path in expected_paths - generated_paths)
        extra = sorted(path.as_posix() for path in generated_paths - expected_paths)
        raise AssertionError(
            f"formatter corpus path mismatch: missing={missing or '-'} extra={extra or '-'}"
        )

    mismatches = []
    for fixture_path in contract.fixture_paths:
        target_path = formatter_target_path(fixture_path)
        if (repo / target_path).read_bytes() != (SKILL_ROOT / fixture_path).read_bytes():
            mismatches.append(target_path.as_posix())
    if mismatches:
        raise AssertionError(
            "generated Markdown differs from the pinned formatter corpus: "
            + ", ".join(sorted(mismatches))
        )


class FormatterCompatibilityTests(unittest.TestCase):
    def test_selected_contract_binds_exact_oxfmt_provenance_and_golden_digest(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        contract = catalog.formatter_by_profile[PROFILE_ID]
        provenance = load_provenance()

        self.assertEqual(contract.kind, "selected")
        self.assertEqual(contract.formatter_id, "formatter.oxfmt-markdown")
        self.assertEqual(contract.version, "0.59.0")
        self.assertEqual(provenance["profile"], PROFILE_ID)
        self.assertEqual(provenance["formatter"]["id"], contract.formatter_id)
        self.assertEqual(provenance["formatter"]["version"], contract.version)
        self.assertEqual(
            provenance["realFormatterProbe"]["argv"],
            [
                "bunx",
                f"oxfmt@{contract.version}",
                "--check",
                "AGENTS.md",
                "docs/agents",
            ],
        )
        self.assertEqual(
            provenance["realFormatterProbe"]["command"],
            f"rtk bunx oxfmt@{contract.version} --check AGENTS.md docs/agents",
        )
        self.assertEqual(
            provenance["fixturePaths"],
            [path.as_posix() for path in contract.fixture_paths],
        )
        observed_digest = formatter_corpus_digest(contract)
        self.assertEqual(observed_digest, contract.golden_digest)
        self.assertEqual(observed_digest, provenance["goldenDigest"])

        first_path = contract.fixture_paths[0]
        target_path = formatter_target_path(first_path)
        mutated = [
            (
                formatter_target_path(path).as_posix().encode("utf-8"),
                (SKILL_ROOT / path).read_bytes()
                + (b"\n" if path == first_path else b""),
            )
            for path in contract.fixture_paths
        ]
        self.assertNotEqual(
            portable_file_digest(mutated),
            provenance["goldenDigest"],
            f"mutating {target_path} must invalidate formatter provenance",
        )

    def test_apply_formatter_verification_audit_and_reapply_leave_no_delta(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        provenance = load_provenance()
        real_subprocess_run = subprocess.run
        formatter_processes = []

        def reject_formatter_process(args, *positional, **kwargs):
            argv = [os.fspath(value) for value in args]
            if any("oxfmt" in value.lower() for value in argv):
                formatter_processes.append(argv)
                raise AssertionError(f"setup executed formatter process: {argv}")
            return real_subprocess_run(args, *positional, **kwargs)

        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            (repo / "DESIGN.md").write_text(
                "# Repository-authored design contract\n",
                encoding="utf-8",
            )
            (repo / ".formatter-fixture-verify.py").write_bytes(
                VERIFICATION_SOURCE.read_bytes()
            )
            decisions = formatter_profile_decisions()

            with mock.patch.object(
                context_setup.subprocess,
                "run",
                side_effect=reject_formatter_process,
            ):
                first_apply = run_apply_in_process(repo, PROFILE_ID, decisions)

            self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
            install_profile_skills(repo, PROFILE_ID)
            assert_profile_formatter_canonical(repo, catalog, PROFILE_ID)
            before_composition = snapshot_files(repo)

            verification = subprocess.run(
                shlex.split(provenance["fixtureVerification"]),
                cwd=repo,
                text=True,
                capture_output=True,
                check=False,
                env={
                    **os.environ,
                    "PATH": "/usr/bin:/bin",
                    "PYTHONDONTWRITEBYTECODE": "1",
                },
            )
            self.assertEqual(verification.returncode, 0, verification.stderr)

            audit = run_fixture_audit(
                repo,
                "--format",
                "json",
                extra_env={
                    "PATH": "/usr/bin:/bin",
                    "PYTHONDONTWRITEBYTECODE": "1",
                },
            )
            with mock.patch.object(
                context_setup.subprocess,
                "run",
                side_effect=reject_formatter_process,
            ):
                second_apply = run_apply_in_process(
                    repo,
                    PROFILE_ID,
                    [],
                )

            self.assertEqual(audit.returncode, 0, audit.stdout)
            self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
            self.assertEqual(json.loads(second_apply.stdout)["plannedChanges"], [])
            self.assertEqual(snapshot_files(repo), before_composition)
            self.assertEqual(formatter_processes, [])


if __name__ == "__main__":
    unittest.main()
