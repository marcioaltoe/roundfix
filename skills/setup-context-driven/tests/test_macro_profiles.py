"""Macro profile flows for setup-context-driven.

Suite: supported profile repository flows
Invariant: every bundled profile applies, audits cleanly, and re-applies without file changes.
Boundary IN: context_setup.py CLI, bundled assets, temporary repository files.
Boundary OUT: Makefile orchestration and embedded skill synchronization checks.
"""

import hashlib
import json
import os
import shlex
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(Path(__file__).resolve().parent))

from test_apply import BASE_DECISIONS, run_apply, run_audit, run_context_setup  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    run_audit as run_fixture_audit,
    snapshot_files,
    write_skill,
)
from test_skills import write_lockfile  # noqa: E402
from test_formatter_compatibility import (  # noqa: E402
    VERIFICATION_SOURCE as FORMATTER_VERIFICATION_SOURCE,
    assert_profile_formatter_canonical,
    formatter_profile_decisions,
)
sys.path.insert(0, str(SKILL_ROOT / "scripts"))

from context_assets import (  # noqa: E402
    build_standard_profile_plan,
    load_asset_catalog,
    render_standard_profile_snapshot,
)
from context_setup import render_skill_dispatch  # noqa: E402


SUPPORTED_PROFILES = [
    "go-cli-tui",
    "rust-cli",
    "standard-typescript-monorepo",
]
MACRO_VERIFICATION_COMMAND = "python3 -B .macro-profile-verify.py"
MACRO_VERIFICATION_SOURCE = (
    Path(__file__).resolve().parent
    / "fixtures"
    / "macro-profiles"
    / "verify_fixture.py"
)


class ProfileMacroFlowTests(unittest.TestCase):
    def test_standard_profile_snapshot_is_deterministic_alongside_existing_profiles(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        decision = {
            "mode": "REST",
            "exceptions": [],
            "source": {
                "path": "docs/architecture/http-contract.json",
                "digest": "c" * 64,
            },
        }

        first = build_standard_profile_plan(catalog, decision)
        second = build_standard_profile_plan(catalog, decision)

        self.assertEqual(
            render_standard_profile_snapshot(first),
            render_standard_profile_snapshot(second),
        )
        self.assertEqual(
            set(SUPPORTED_PROFILES),
            set(catalog.profiles),
        )

    def test_profiles_declare_complete_rule_coverage_and_verification_entry_decision(self):
        catalog = load_asset_catalog(SKILL_ROOT)

        for profile_id in SUPPORTED_PROFILES:
            with self.subTest(profile=profile_id):
                profile = catalog.profiles[profile_id]
                declared_rules = set(profile["requiredRules"])
                module_rules = {
                    rule["id"]
                    for module_id in catalog.ordered_modules_by_profile[profile_id]
                    for rule in catalog.modules[module_id]["rules"]
                }
                covered_categories = {
                    coverage_id
                    for rule_id in declared_rules
                    for coverage_id in catalog.rule_contracts[rule_id].coverage
                }

                self.assertEqual(declared_rules, module_rules)
                self.assertIn("verification.gate", profile["entryDecisions"])
                self.assertIn(profile_id, catalog.formatter_by_profile)
                self.assertTrue(
                    {
                        "coverage.universal-safety",
                        "coverage.verification",
                        "coverage.verification-integrity",
                        "coverage.skill-dispatch",
                        "coverage.language",
                        "coverage.research",
                        "coverage.dependencies",
                        "coverage.git-delivery",
                        "coverage.security-configuration",
                    }.issubset(covered_categories)
                )

    def test_supported_profiles_apply_audit_clean_and_reapply_without_changes(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        for profile_id in SUPPORTED_PROFILES:
            with self.subTest(profile=profile_id):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    prepare_repository_owned_contracts(repo, profile_id)
                    formatter = catalog.formatter_by_profile[profile_id]
                    if formatter.kind == "selected":
                        (repo / ".formatter-fixture-verify.py").write_bytes(
                            FORMATTER_VERIFICATION_SOURCE.read_bytes()
                        )
                        decisions = formatter_profile_decisions()
                    else:
                        (repo / ".macro-profile-verify.py").write_bytes(
                            MACRO_VERIFICATION_SOURCE.read_bytes()
                        )
                        decisions = decisions_with(
                            autonomous=False,
                            verification_gate=MACRO_VERIFICATION_COMMAND,
                        )

                    first_apply = run_profile_apply(repo, profile_id, decisions)
                    install_profile_skills(repo, profile_id)
                    after_apply = snapshot_files(repo)
                    run_profile_formatter(repo, catalog, profile_id)
                    verification = run_persisted_verification(repo)
                    first_audit = run_audit(repo)
                    run_profile_formatter(repo, catalog, profile_id)
                    second_audit = run_audit(repo)
                    after_composition = snapshot_files(repo)
                    second_apply = run_profile_apply(repo, profile_id, [])

                    self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
                    self.assertEqual(verification.returncode, 0, verification.stderr)
                    self.assertEqual(first_audit.returncode, 0, first_audit.stderr)
                    self.assertEqual(second_audit.returncode, 0, second_audit.stderr)
                    self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
                    self.assertEqual(after_composition, after_apply)
                    self.assertEqual(snapshot_files(repo), after_composition)

    def test_supported_profiles_cover_representative_decision_combinations(self):
        cases = [
            (
                "standard-typescript-monorepo",
                decisions_with(
                    domain_layout="multi-context",
                    triage_external=True,
                    runtime_backend="codex macro-backend xhigh",
                    runtime_design="claude macro-design xhigh",
                    secondbrain=True,
                ),
                ["root.spec-workflow", "root.external-triage", "root.autonomous-work", "root.secondbrain"],
                [],
            ),
            (
                "go-cli-tui",
                decisions_with(
                    spec_scaffold=False,
                    autonomous=False,
                ),
                [],
                ["root.spec-workflow", "root.external-triage", "root.autonomous-work", "root.secondbrain"],
            ),
            (
                "rust-cli",
                decisions_with(
                    domain_layout="multi-context",
                    triage_external=True,
                    autonomous=False,
                    secondbrain=True,
                ),
                ["root.external-triage", "root.secondbrain"],
                ["root.autonomous-work"],
            ),
        ]
        for profile_id, decisions, present_markers, absent_markers in cases:
            with self.subTest(profile=profile_id):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    prepare_repository_owned_contracts(repo, profile_id)

                    first_apply = run_profile_apply(repo, profile_id, decisions)
                    install_profile_skills(repo, profile_id)
                    clean_audit = run_audit(repo)
                    after_audit = snapshot_files(repo)
                    second_apply = run_profile_apply(repo, profile_id, [])

                    self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
                    self.assertEqual(clean_audit.returncode, 0, clean_audit.stderr)
                    self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
                    self.assertEqual(snapshot_files(repo), after_audit)
                    generated = generated_text(repo)
                    for marker in present_markers:
                        self.assertIn(marker, generated)
                    for marker in absent_markers:
                        self.assertNotIn(marker, generated)
                    manifest = json.loads(
                        (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
                    )
                    for decision_id in [
                        "spec.scaffold",
                        "domain.layout",
                        "triage.external",
                        "autonomous.enabled",
                        "runtime.backend",
                        "runtime.design",
                        "verification.gate",
                        "language.generated",
                        "secondbrain.enabled",
                        "repository.extension.enabled",
                    ]:
                        self.assertIn(decision_id, manifest["decisions"])
                    if profile_id == "standard-typescript-monorepo":
                        self.assertIn("macro-backend", generated)
                        self.assertIn("macro-design", generated)

    def test_every_profile_renders_selected_verification_and_exact_active_skill_dispatch(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        verification = "task-02-profile-verification"

        for profile_id in SUPPORTED_PROFILES:
            with self.subTest(profile=profile_id):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    prepare_repository_owned_contracts(repo, profile_id)
                    decisions = decisions_with(
                        autonomous=False,
                        verification_gate=verification,
                    )

                    applied = run_profile_apply(repo, profile_id, decisions)

                    self.assertEqual(applied.returncode, 0, applied.stderr)
                    instructions = (repo / "docs" / "agents" / "agent-instructions.md").read_text(
                        encoding="utf-8"
                    )
                    dispatch = (repo / "docs" / "agents" / "skill-dispatch.md").read_text(
                        encoding="utf-8"
                    )
                    manifest = json.loads(
                        (repo / "docs" / "agents" / "setup-context.json").read_text(
                            encoding="utf-8"
                        )
                    )
                    setup = catalog.setups[catalog.profiles[profile_id]["setup"]]
                    installed = {skill["name"] for skill in setup["skills"]}
                    required = {
                        skill
                        for module_id in manifest["modules"]
                        for skill in catalog.modules[module_id]["requiredSkills"]
                    }
                    expected_dispatch = render_skill_dispatch(
                        catalog,
                        manifest["modules"],
                    )

                    self.assertIn(f"`{verification}`", instructions)
                    self.assertEqual(installed, required)
                    self.assertIn(expected_dispatch, dispatch)
                    self.assertEqual(
                        [
                            line
                            for line in dispatch.splitlines()
                            if line.startswith("- `")
                            and not line.startswith("- `trigger.")
                        ],
                        [f"- `{skill}`:" for skill in sorted(installed)],
                    )

    def test_every_profile_preserves_repository_extensions_and_rerenders_identically(self):
        for profile_id in SUPPORTED_PROFILES:
            with self.subTest(profile=profile_id):
                with tempfile.TemporaryDirectory() as temp_dir:
                    repo = Path(temp_dir)
                    prepare_repository_owned_contracts(repo, profile_id)
                    applied = run_profile_apply(repo, profile_id, BASE_DECISIONS)
                    self.assertEqual(applied.returncode, 0, applied.stderr)
                    root = repo / "AGENTS.md"
                    extension = "\nRepository-authored extension: keep this byte-for-byte.\n"
                    root.write_text(root.read_text(encoding="utf-8") + extension, encoding="utf-8")
                    before = snapshot_files(repo)

                    reapplied = run_profile_apply(repo, profile_id, [])

                    self.assertEqual(reapplied.returncode, 0, reapplied.stderr)
                    self.assertEqual(snapshot_files(repo), before)
                    self.assertTrue(root.read_text(encoding="utf-8").endswith(extension))

    def test_standard_typescript_single_context_monorepo_always_generates_monorepo_guide(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            prepare_repository_owned_contracts(repo, "standard-typescript-monorepo")

            applied = run_profile_apply(
                repo,
                "standard-typescript-monorepo",
                decisions_with(domain_layout="single-context", autonomous=False),
            )

            self.assertEqual(applied.returncode, 0, applied.stderr)
            root = (repo / "AGENTS.md").read_text(encoding="utf-8")
            self.assertIn("id=root.monorepo", root)
            self.assertTrue((repo / "docs" / "agents" / "monorepo.md").is_file())

    def test_required_skill_failure_and_extra_reporting_keep_exit_semantics(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            first_apply = run_apply(repo, "rust-cli", BASE_DECISIONS)
            install_profile_skills(repo, "rust-cli", omit={"agentic-cli-design"})
            missing_required = run_audit(repo)
            write_skill(repo, "agentic-cli-design")
            write_skill(repo, "autoresearch")
            write_lockfile(repo, ["autoresearch"])
            compliant_without_extras = run_audit(repo)
            extras_visible = run_audit_cli(repo, "--show-extra-skills")

            self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
            self.assertEqual(missing_required.returncode, 1)
            self.assertFinding(missing_required, "skills.required.missing", "error")
            self.assertEqual(compliant_without_extras.returncode, 0, compliant_without_extras.stderr)
            self.assertNoFinding(compliant_without_extras, "skills.extra.installed")
            self.assertEqual(extras_visible.returncode, 0, extras_visible.stderr)
            extra = self.finding(extras_visible, "skills.extra.installed")
            self.assertEqual(extra["severity"], "info")

    def test_secondbrain_opt_in_profile_applies_audits_and_reapplies_without_changes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)

            first_apply = run_apply(repo, "go-cli-tui", secondbrain_decisions())
            install_profile_skills(repo, "go-cli-tui")
            clean_audit = run_audit(repo)
            after_audit = snapshot_files(repo)
            second_apply = run_apply(repo, "go-cli-tui", [])

            self.assertEqual(first_apply.returncode, 0, first_apply.stderr)
            self.assertEqual(clean_audit.returncode, 0, clean_audit.stderr)
            self.assertEqual(second_apply.returncode, 0, second_apply.stderr)
            self.assertEqual(snapshot_files(repo), after_audit)
            manifest = json.loads(
                (repo / "docs" / "agents" / "setup-context.json").read_text(encoding="utf-8")
            )
            self.assertIn("secondbrain", manifest["modules"])
            self.assertIn("root.secondbrain", (repo / "AGENTS.md").read_text(encoding="utf-8"))
            self.assertTrue((repo / "docs" / "agents" / "secondbrain.md").is_file())

    def finding(self, result, code):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        return matches[0]

    def assertFinding(self, result, code, severity):
        match = self.finding(result, code)
        self.assertEqual(match["severity"], severity)

    def assertNoFinding(self, result, code):
        payload = json.loads(result.stdout)
        self.assertEqual(
            [finding for finding in payload["findings"] if finding["code"] == code],
            [],
            payload,
        )


def secondbrain_decisions():
    return [
        "secondbrain.enabled=true"
        if decision.startswith("secondbrain.enabled=")
        else decision
        for decision in BASE_DECISIONS
    ]


def decisions_with(
    *,
    spec_scaffold=True,
    domain_layout="single-context",
    triage_external=False,
    autonomous=True,
    runtime_backend="codex gpt-5.5 xhigh",
    runtime_design="claude opus xhigh",
    verification_gate="make verify",
    secondbrain=False,
):
    values = {
        "spec.scaffold": str(spec_scaffold).lower(),
        "domain.layout": domain_layout,
        "triage.external": str(triage_external).lower(),
        "autonomous.enabled": str(autonomous).lower(),
        "runtime.backend": runtime_backend,
        "runtime.design": runtime_design,
        "verification.gate": verification_gate,
        "language.generated": "English",
        "secondbrain.enabled": str(secondbrain).lower(),
        "repository.extension.enabled": "false",
    }
    return [f"{key}={value}" for key, value in values.items()]


def generated_text(repo):
    paths = [repo / "AGENTS.md", *sorted((repo / "docs" / "agents").glob("*.md"))]
    return "\n".join(path.read_text(encoding="utf-8") for path in paths if path.exists())


def run_profile_formatter(repo, catalog, profile_id):
    contract = catalog.formatter_by_profile[profile_id]
    if contract.kind == "none":
        return
    assert_profile_formatter_canonical(repo, catalog, profile_id)


def run_persisted_verification(repo):
    manifest = json.loads(
        (repo / "docs" / "agents" / "setup-context.json").read_text(
            encoding="utf-8"
        )
    )
    command = manifest["decisions"]["verification.gate"]["value"]
    return subprocess.run(
        shlex.split(command),
        cwd=repo,
        text=True,
        capture_output=True,
        check=False,
        env={
            **os.environ,
            "PYTHONDONTWRITEBYTECODE": "1",
        },
    )


def prepare_repository_owned_contracts(repo, profile_id):
    if profile_id != "standard-typescript-monorepo":
        return
    (repo / "DESIGN.md").write_text(
        "# Repository-authored design contract\n",
        encoding="utf-8",
    )
    http_path = repo / "docs" / "architecture" / "http-contract.json"
    http_path.parent.mkdir(parents=True)
    http_path.write_text('{"mode":"REST"}\n', encoding="utf-8")


def run_profile_apply(repo, profile_id, decisions):
    if profile_id != "standard-typescript-monorepo" or not decisions:
        return run_apply(repo, profile_id, decisions)

    http_bytes = (repo / "docs" / "architecture" / "http-contract.json").read_bytes()
    document_decisions = []
    for decision in decisions:
        decision_id, _, raw_value = decision.partition("=")
        value = {"true": True, "false": False}.get(raw_value, raw_value)
        document_decisions.append({"id": decision_id, "value": value})
    document_decisions.append(
        {
            "id": "http.contract",
            "value": {
                "mode": "REST",
                "exceptions": [],
                "source": {
                    "path": "docs/architecture/http-contract.json",
                    "digest": hashlib.sha256(http_bytes).hexdigest(),
                },
            },
        }
    )
    with tempfile.NamedTemporaryFile(mode="w", suffix=".json", encoding="utf-8") as handle:
        json.dump(
            {
                "schemaVersion": "setup-context-driven/decisions/0.0.1",
                "version": "0.0.1",
                "decisions": document_decisions,
            },
            handle,
        )
        handle.flush()
        args = [
            "apply",
            "--repo",
            str(repo),
            "--format",
            "json",
            "--profile",
            profile_id,
            "--decision-file",
            handle.name,
        ]
        preview = run_context_setup(*args)
        if preview.returncode != 3:
            return preview
        payload = json.loads(preview.stdout)
        if not any(
            finding["code"] == "plan.confirmation.required"
            for finding in payload["findings"]
        ):
            return preview
        return run_context_setup(*args, "--confirm-plan", payload["planDigest"])


def run_audit_cli(repo, *extra_args):
    return run_fixture_audit(repo, "--format", "json", *extra_args)


if __name__ == "__main__":
    unittest.main()
