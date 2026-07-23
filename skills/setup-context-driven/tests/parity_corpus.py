"""Generate the standalone Python-to-Go Baseline parity corpus.

Suite: Baseline cutover characterization data
Invariant: every maintained Python test has one classified Go destination and
every representative state regenerates the same standalone bytes.
Boundary IN: the current Python setup runtime, its 240 pre-cutover tests, and
the checked-in Baseline assets.
Boundary OUT: the future Go implementation that consumes these JSON files.
"""

from __future__ import annotations

import ast
import base64
from contextlib import contextmanager
from datetime import date
import hashlib
import json
import os
from pathlib import Path
import shutil
import stat
import tempfile
from unittest import mock


SKILL_ROOT = Path(__file__).resolve().parents[1]
TESTS_ROOT = Path(__file__).resolve().parent
CORPUS_ROOT = SKILL_ROOT / "assets" / "parity-corpus" / "v1"
FIXED_DATE = "2026-07-23"
SCHEMA_VERSION = "setup-context-driven/parity-corpus/v1"
MATRIX_SCHEMA_VERSION = "setup-context-driven/parity-matrix/v1"
FIXTURE_SCHEMA_VERSION = "setup-context-driven/parity-fixture/v1"


import sys

sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(TESTS_ROOT))

import context_setup  # noqa: E402
from context_assets import clone_assets_to, load_asset_catalog  # noqa: E402
from context_baseline import (  # noqa: E402
    SourceInventoryError,
    inventory_incompatible_source_baseline,
)
import test_apply  # noqa: E402
import test_audit  # noqa: E402
import test_formatter_compatibility  # noqa: E402
import test_readoption_apply  # noqa: E402
import test_restore_skills  # noqa: E402
import test_source_inventory  # noqa: E402
import test_sync_setups  # noqa: E402


ALLOWED_CLASSIFICATIONS = (
    "exact",
    "semantic",
    "designed-delta",
    "ancillary",
    "retired",
)

SEMANTIC_FILES = {
    "test_audit.py",
    "test_autonomous_secondbrain_decisions.py",
    "test_decision_rendering.py",
    "test_delegation.py",
    "test_legacy_rule_ledger.py",
    "test_macro_profiles.py",
    "test_operational_guidance.py",
    "test_secondbrain.py",
    "test_spec_triage_decisions.py",
    "test_standard_typescript_monorepo.py",
    "test_workflow.py",
}

ANCILLARY_TESTS = {
    "test_workflow.py::SetupWorkflowTests::test_asset_maintenance_doc_publishes_catalog_and_source_boundaries",
}

DESIGNED_DELTA_TESTS = {
    "test_audit.py::AuditCliTests::test_audit_help_describes_implemented_optional_checks":
        "The public Baseline Command replaces Python audit help; safety semantics remain, but command names and help bytes intentionally change.",
    "test_audit.py::AuditCliTests::test_default_subcommand_is_read_only_audit":
        "The public root command is an interactive adoption or update state machine; automation uses explicit plan and apply subcommands.",
    "test_preview.py::PreviewCliTests::test_top_level_help_names_every_exit_three_condition":
        "The Go CLI keeps exit category 3 but publishes it through the Baseline Command help instead of the Python setup script.",
    "test_restore_skills.py::RestoreSkillsCliTests::test_help_exposes_non_interactive_restore_contract":
        "The maintained restoration responsibility moves to roundfix baseline skills restore with the public CLI contract.",
    "test_restore_skills.py::RestoreSkillsCliTests::test_text_preview_is_deterministic_and_names_confirmation":
        "The Go operation preserves digest-bound confirmation while intentionally replacing Python text and command naming.",
    "test_workflow.py::SetupWorkflowTests::test_skill_workflow_requires_preview_and_confirmation_before_apply":
        "The setup skill becomes guidance over the public Baseline Command and no longer describes an independent Python engine.",
    "test_workflow.py::SetupWorkflowTests::test_skill_names_canonical_setup_and_optional_extra_skill_report":
        "The thin setup skill delegates to the public Baseline Command, so legacy Python workflow prose is replaced.",
    "test_workflow.py::SetupWorkflowTests::test_operator_docs_publish_schema_exit_and_confirmation_contract":
        "Public roundfix/baseline schemas and commands replace Python audit-v1 and restore-v1 documentation while retaining safety semantics.",
}

DESIGNED_DELTA_ROWS = (
    {
        "id": "delta.public-baseline-state-machine",
        "behavior": "Interactive adoption and update share one public Baseline Command state machine while automation uses explicit plan and apply.",
        "goDestination": "internal/cli/baseline.RunContext and internal/baseline.Workflow",
        "rationale": "The Python skill exposed separate audit, apply, and Decision File steps; the accepted product contract consolidates them without weakening decisions or confirmation.",
        "fixtureIds": ["greenfield-rust-cli", "update-rust-cli", "readoption-preservation"],
        "contractDimensions": ["input", "state", "action"],
    },
    {
        "id": "delta.public-plan-document-schema",
        "behavior": "Portable automation exchanges roundfix/baseline-plan/v1 and roundfix/baseline-result/v1 documents.",
        "goDestination": "internal/baseline.PlanDocument and internal/baseline.Result",
        "rationale": "The public Go API intentionally replaces Python audit-v1 output while retaining normalized decisions, exact bytes, refusals, and confirmation evidence.",
        "fixtureIds": ["greenfield-go-cli-tui", "stale-plan-refusal"],
        "contractDimensions": ["input", "refusal", "digest", "planned-byte-sequence"],
    },
    {
        "id": "delta.file-change-projection",
        "behavior": "Every Change Plan leads with a fileChanges projection while the complete managed-entry ledger remains canonical.",
        "goDestination": "internal/baseline.PlanDocument.FileChanges",
        "rationale": "Python plannedChanges is managed-entry oriented; the new file projection is an additive review surface derived from the retained exact ledger.",
        "fixtureIds": ["greenfield-standard-typescript-monorepo", "readoption-preservation"],
        "contractDimensions": ["state", "planned-byte-sequence", "digest"],
    },
)

GO_DESTINATIONS = {
    "test_apply.py": "internal/baseline.Workflow.Plan and internal/baseline.Transaction",
    "test_assets.py": "internal/baseline.Catalog",
    "test_audit.py": "internal/baseline.Audit and internal/cli/baseline",
    "test_autonomous_secondbrain_decisions.py": "internal/baseline.DecisionPlan",
    "test_capabilities.py": "internal/baseline.RepositoryCapabilities",
    "test_decision_plan_contracts.py": "internal/baseline.DecisionPlan",
    "test_decision_rendering.py": "internal/baseline.Renderer",
    "test_delegation.py": "internal/baseline.Inventory",
    "test_documentation_contract.py": "internal/cli/baseline documentation contracts",
    "test_formatter_compatibility.py": "internal/baseline.FormatterContract",
    "test_governed_corpus.py": "internal/baseline.Catalog",
    "test_legacy_rule_ledger.py": "internal/baseline.UpgradeRetentionContract",
    "test_macro_profiles.py": "internal/baseline.Workflow",
    "test_manifest_migration.py": "internal/baseline.ManifestReader",
    "test_operational_guidance.py": "internal/baseline.Renderer",
    "test_preview.py": "internal/baseline.PlanDocument and internal/cli/baseline",
    "test_profile_alignment.py": "internal/baseline.Catalog",
    "test_readoption_apply.py": "internal/baseline.Readoption and internal/baseline.Transaction",
    "test_readoption_decisions.py": "internal/baseline.Readoption",
    "test_repository_extension.py": "internal/baseline.RepositoryOwnedExtension",
    "test_restore_skills.py": "internal/baseline/maintenance.SkillsRestore",
    "test_secondbrain.py": "internal/baseline.Renderer",
    "test_skill_dispatch.py": "internal/baseline.SkillActivation",
    "test_skills.py": "internal/baseline.RepositorySkillSet",
    "test_source_baselines.py": "internal/baseline.SourceBaseline",
    "test_source_inventory.py": "internal/baseline.Inventory",
    "test_spec_triage_decisions.py": "internal/baseline.DecisionPlan",
    "test_standard_typescript_monorepo.py": "internal/baseline.Catalog and internal/baseline.Workflow",
    "test_sync_setups.py": "internal/baseline/maintenance.AssetsSync",
    "test_upgrade_contracts.py": "internal/baseline.ManifestReader",
    "test_upgrade_retention.py": "internal/baseline.UpgradeRetentionContract",
    "test_version_contract.py": "internal/baseline.ManifestReader",
    "test_workflow.py": "internal/baseline.Workflow and internal/cli/baseline",
}

FILE_FIXTURES = {
    "test_apply.py": ["greenfield-rust-cli", "update-rust-cli", "stale-plan-refusal", "atomic-rollback"],
    "test_assets.py": ["greenfield-go-cli-tui", "greenfield-rust-cli", "greenfield-standard-typescript-monorepo"],
    "test_audit.py": ["update-rust-cli", "unsafe-carrier-refusal"],
    "test_capabilities.py": ["missing-capability-refusal", "greenfield-standard-typescript-monorepo"],
    "test_formatter_compatibility.py": ["formatter-composition"],
    "test_macro_profiles.py": ["greenfield-go-cli-tui", "greenfield-rust-cli", "greenfield-standard-typescript-monorepo"],
    "test_profile_alignment.py": ["profile-change-rust-to-go-cli-tui"],
    "test_readoption_apply.py": ["readoption-preservation", "atomic-rollback", "stale-plan-refusal"],
    "test_readoption_decisions.py": ["readoption-preservation"],
    "test_restore_skills.py": ["skill-restoration", "skill-restoration-rollback"],
    "test_source_inventory.py": ["readoption-preservation", "unsafe-carrier-refusal"],
    "test_sync_setups.py": ["asset-sync"],
}


class FrozenDate(date):
    @classmethod
    def today(cls):
        return cls(2026, 7, 23)


def canonical_json(document):
    return (
        json.dumps(document, indent=2, sort_keys=True, ensure_ascii=False)
        + "\n"
    ).encode("utf-8")


def sha256_bytes(content):
    return hashlib.sha256(content).hexdigest()


def python_test_inventory(skill_root=SKILL_ROOT):
    tests_root = skill_root / "tests"
    rows = []
    for path in sorted(tests_root.glob("test_*.py")):
        if path.name == "test_parity_corpus.py":
            continue
        tree = ast.parse(path.read_text(encoding="utf-8"), filename=str(path))
        for class_node in (
            node for node in tree.body if isinstance(node, ast.ClassDef)
        ):
            for method in (
                node
                for node in class_node.body
                if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef))
                and node.name.startswith("test_")
            ):
                key = f"{path.name}::{class_node.name}::{method.name}"
                classification, rationale = classify_test(key, path.name)
                rows.append(
                    {
                        "id": f"python.{path.stem}.{class_node.name}.{method.name}",
                        "sourceKind": "python-test",
                        "sourcePath": f"tests/{path.name}",
                        "sourceSuite": class_node.name,
                        "sourceTest": method.name,
                        "behavior": method.name.removeprefix("test_").replace("_", " "),
                        "classification": classification,
                        "rationale": rationale,
                        "goDestination": GO_DESTINATIONS[path.name],
                        "fixtureIds": fixture_ids_for(path.name),
                        "contractDimensions": contract_dimensions(key),
                    }
                )
    return sorted(rows, key=lambda row: row["id"])


def classify_test(key, filename):
    if key in DESIGNED_DELTA_TESTS:
        return "designed-delta", DESIGNED_DELTA_TESTS[key]
    if key in ANCILLARY_TESTS:
        return (
            "ancillary",
            "This test protects maintainer-only asset guidance rather than public Baseline runtime behavior; the Go maintenance operation retains the contract.",
        )
    if filename == "test_documentation_contract.py":
        return (
            "designed-delta",
            "The public Baseline Command replaces Python-backed setup documentation; the Go documentation contract retains the behavior without preserving legacy prose.",
        )
    if filename in SEMANTIC_FILES:
        return (
            "semantic",
            "The user-visible invariant and normalized outcome must match; incidental prose and Python implementation structure are not byte contracts.",
        )
    return "exact", None


def fixture_ids_for(filename):
    return FILE_FIXTURES.get(filename, ["greenfield-rust-cli"])


def contract_dimensions(key):
    lowered = key.lower()
    patterns = {
        "input": ("input", "decision", "schema", "flag", "profile", "source", "contract"),
        "state": ("state", "manifest", "inventory", "repository", "snapshot", "audit", "idempot", "reapply"),
        "action": ("apply", "restore", "sync", "render", "create", "remove", "update", "write"),
        "refusal": ("reject", "invalid", "missing", "unsafe", "stale", "failure", "block", "prevent", "mismatch", "unsupported", "drift", "limit"),
        "digest": ("digest", "hash", "identity", "confirmation", "snapshot", "provenance"),
        "planned-byte-sequence": ("plan", "preview", "apply", "render", "artifact", "output", "restore"),
        "rollback-outcome": ("rollback", "restore_every", "failure_before", "postwrite", "tamper", "atomic", "preimage"),
    }
    dimensions = [
        name for name, terms in patterns.items() if any(term in lowered for term in terms)
    ]
    return dimensions or ["state"]


@contextmanager
def deterministic_runtime():
    git_dates = {
        "GIT_AUTHOR_DATE": "2026-07-23T12:00:00+00:00",
        "GIT_COMMITTER_DATE": "2026-07-23T12:00:00+00:00",
        "TZ": "UTC",
    }
    with (
        mock.patch.object(context_setup, "date", FrozenDate),
        mock.patch.dict(os.environ, git_dates),
    ):
        yield


def capture_state(root):
    state = {}
    if not root.exists():
        return state
    for path in sorted(root.rglob("*")):
        relative = path.relative_to(root).as_posix()
        metadata = path.lstat()
        mode = stat.S_IMODE(metadata.st_mode)
        if stat.S_ISLNK(metadata.st_mode):
            target = os.readlink(path)
            identity = "symlink-sha256:" + sha256_bytes(target.encode("utf-8"))
            state[relative] = {
                "kind": "symlink",
                "mode": mode,
                "identity": identity,
                "target": target,
            }
        elif stat.S_ISREG(metadata.st_mode):
            content = path.read_bytes()
            state[relative] = {
                "kind": "file",
                "mode": mode,
                "identity": "sha256:" + sha256_bytes(content),
                "bytes": content,
            }
    return state


def state_records(state, blobs):
    records = []
    for path, item in sorted(state.items()):
        record = {
            "path": path,
            "kind": item["kind"],
            "mode": item["mode"],
            "identity": item["identity"],
        }
        if item["kind"] == "file":
            content = item["bytes"]
            record["size"] = len(content)
            blobs[item["identity"]] = base64.b64encode(content).decode("ascii")
        else:
            record["target"] = item["target"]
        records.append(record)
    return records


def state_digest(state):
    normalized = []
    for path, item in sorted(state.items()):
        normalized.append(
            {
                "path": path,
                "kind": item["kind"],
                "mode": item["mode"],
                "identity": item["identity"],
                "target": item.get("target"),
            }
        )
    return "sha256:" + sha256_bytes(canonical_json(normalized))


def content_state_digest(state):
    normalized = []
    for path, item in sorted(state.items()):
        normalized.append(
            {
                "path": path,
                "kind": item["kind"],
                "identity": item["identity"],
                "target": item.get("target"),
            }
        )
    return "sha256:" + sha256_bytes(canonical_json(normalized))


def content_state_equal(before, after):
    return content_state_digest(before) == content_state_digest(after)


def normalize_paths(value, replacements):
    if isinstance(value, dict):
        return {
            key: normalize_paths(item, replacements)
            for key, item in sorted(value.items())
        }
    if isinstance(value, list):
        return [normalize_paths(item, replacements) for item in value]
    if isinstance(value, str):
        for source, target in replacements:
            value = value.replace(source, target)
    return value


def process_result(result, replacements=()):
    stdout = result.stdout.strip()
    try:
        output = json.loads(stdout) if stdout else None
    except json.JSONDecodeError:
        output = stdout
    return {
        "exitCode": result.returncode,
        "stdout": normalize_paths(output, replacements),
        "stderr": normalize_paths(result.stderr, replacements),
    }


def planned_byte_sequence(planned_changes, before, after):
    sequence = []
    seen = set()
    for ordinal, change in enumerate(planned_changes):
        path = change.get("path")
        if not path:
            continue
        before_item = before.get(path)
        after_item = after.get(path)
        sequence.append(
            {
                "ordinal": ordinal,
                "path": path,
                "managedId": change.get("managedId") or change.get("skill"),
                "action": change.get("action"),
                "beforeIdentity": before_item["identity"] if before_item else None,
                "afterIdentity": after_item["identity"] if after_item else None,
            }
        )
        seen.add(path)
    for path in sorted(set(before) | set(after)):
        if path in seen or before.get(path) == after.get(path):
            continue
        before_item = before.get(path)
        after_item = after.get(path)
        sequence.append(
            {
                "ordinal": len(sequence),
                "path": path,
                "managedId": None,
                "action": "replace" if before_item and after_item else ("create" if after_item else "remove"),
                "beforeIdentity": before_item["identity"] if before_item else None,
                "afterIdentity": after_item["identity"] if after_item else None,
            }
        )
    return sequence


def fixture_document(
    *,
    fixture_id,
    classification,
    rationale,
    profile,
    adoption_state,
    contract_areas,
    source_tests,
    input_document,
    normalized_output,
    before,
    after,
    planned_changes,
    plan_digest,
    managed_entry_ledger,
    manifest,
    refusal,
    rollback,
    blobs,
):
    return {
        "schemaVersion": FIXTURE_SCHEMA_VERSION,
        "id": fixture_id,
        "classification": classification,
        "rationale": rationale,
        "profile": profile,
        "adoptionState": adoption_state,
        "contractAreas": contract_areas,
        "sourceTests": source_tests,
        "input": input_document,
        "normalizedOutput": normalized_output,
        "repositoryPreimage": state_records(before, blobs),
        "fileIdentities": {
            "preimage": state_digest(before),
            "postimage": state_digest(after),
        },
        "managedEntryLedger": managed_entry_ledger,
        "manifest": manifest,
        "planDigest": plan_digest,
        "plannedByteSequence": planned_byte_sequence(
            planned_changes, before, after
        ),
        "refusal": refusal,
        "postState": state_records(after, blobs),
        "rollback": rollback,
    }


def apply_args(repo, profile, decisions=(), decision_file=None, confirmation=None):
    args = [
        "apply",
        "--repo",
        str(repo),
        "--format",
        "json",
        "--profile",
        profile,
    ]
    if decision_file is not None:
        args.extend(["--decision-file", str(decision_file)])
    for decision in decisions:
        args.extend(["--decision", decision])
    if confirmation is not None:
        args.extend(["--confirm-plan", confirmation])
    return args


def run_in_process(args):
    return test_apply.run_context_setup_in_process(*args)


def successful_apply_fixture(
    blobs,
    *,
    fixture_id,
    profile,
    adoption_state,
    decisions=(),
    setup=None,
    decision_document=None,
    source_tests,
):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        repo.mkdir()
        if setup is not None:
            setup(repo)
        decision_file = None
        if decision_document is not None:
            decision_file = root / "decisions.json"
            decision_file.write_bytes(canonical_json(decision_document))
        before = capture_state(repo)
        preview = run_in_process(
            apply_args(repo, profile, decisions, decision_file)
        )
        preview_payload = json.loads(preview.stdout)
        digest = preview_payload["planDigest"]
        applied = run_in_process(
            apply_args(repo, profile, decisions, decision_file, digest)
        )
        after = capture_state(repo)
        manifest_path = repo / "docs" / "agents" / "setup-context.json"
        manifest = (
            json.loads(manifest_path.read_text(encoding="utf-8"))
            if manifest_path.is_file()
            else None
        )
        return fixture_document(
            fixture_id=fixture_id,
            classification="exact",
            rationale=None,
            profile=profile,
            adoption_state=adoption_state,
            contract_areas=[
                "input",
                "state",
                "action",
                "digest",
                "planned-byte-sequence",
                "manifest",
                "managed-entry-ledger",
                "post-state",
            ],
            source_tests=source_tests,
            input_document={
                "action": "apply",
                "profile": profile,
                "decisions": list(decisions),
                "decisionDocument": decision_document,
            },
            normalized_output={
                "preview": process_result(preview, ((str(root), "<fixture-root>"),)),
                "apply": process_result(applied, ((str(root), "<fixture-root>"),)),
            },
            before=before,
            after=after,
            planned_changes=preview_payload.get("plannedChanges", []),
            plan_digest=digest,
            managed_entry_ledger=(manifest or {}).get("managedArtifacts", []),
            manifest=manifest,
            refusal=None,
            rollback={
                "attempted": False,
                "succeeded": None,
                "restoredPreimage": None,
                "evidence": "successful post-state recorded",
            },
            blobs=blobs,
        )


def base_decision_document():
    return {
        "schemaVersion": "setup-context-driven/decisions/0.0.1",
        "version": "0.0.1",
        "decisions": test_readoption_apply.ReadoptionApplyTests.profile_decisions(),
    }


def setup_standard_repository(repo, with_capabilities=True):
    test_audit.install_profile_skills(repo, "standard-typescript-monorepo")
    if with_capabilities:
        test_audit.write_standard_profile_capability_evidence(repo)
    (repo / "DESIGN.md").write_text(
        "# Repository-authored design contract\n", encoding="utf-8"
    )
    http_path = repo / "docs" / "architecture" / "http-contract.json"
    http_path.parent.mkdir(parents=True, exist_ok=True)
    http_path.write_text('{"mode":"REST"}\n', encoding="utf-8")


def update_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        repo.mkdir()
        initial = test_apply.run_apply_in_process(
            repo, "rust-cli", test_apply.BASE_DECISIONS
        )
        if initial.returncode != 0:
            raise AssertionError(initial.stdout + initial.stderr)
        test_audit.install_profile_skills(repo, "rust-cli")
        before = capture_state(repo)
        decisions = ["secondbrain.enabled=true"]
        preview = run_in_process(apply_args(repo, "rust-cli", decisions))
        payload = json.loads(preview.stdout)
        applied = run_in_process(
            apply_args(repo, "rust-cli", decisions, confirmation=payload["planDigest"])
        )
        after = capture_state(repo)
        manifest = json.loads(
            (repo / "docs" / "agents" / "setup-context.json").read_text(
                encoding="utf-8"
            )
        )
        return fixture_document(
            fixture_id="update-rust-cli",
            classification="exact",
            rationale=None,
            profile="rust-cli",
            adoption_state="update",
            contract_areas=["state", "action", "digest", "planned-byte-sequence", "manifest", "post-state"],
            source_tests=["tests/test_workflow.py", "tests/test_apply.py"],
            input_document={"action": "update", "profile": "rust-cli", "decisions": decisions},
            normalized_output={"preview": process_result(preview), "apply": process_result(applied)},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=manifest["managedArtifacts"],
            manifest=manifest,
            refusal=None,
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "successful update post-state recorded"},
            blobs=blobs,
        )


def profile_change_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        repo.mkdir()
        initial = test_apply.run_apply_in_process(
            repo, "rust-cli", test_apply.BASE_DECISIONS
        )
        if initial.returncode != 0:
            raise AssertionError(initial.stdout + initial.stderr)
        test_audit.install_profile_skills(repo, "go-cli-tui")
        before = capture_state(repo)
        preview = run_in_process(apply_args(repo, "go-cli-tui"))
        payload = json.loads(preview.stdout)
        applied = run_in_process(
            apply_args(repo, "go-cli-tui", confirmation=payload["planDigest"])
        )
        after = capture_state(repo)
        manifest = json.loads(
            (repo / "docs" / "agents" / "setup-context.json").read_text(
                encoding="utf-8"
            )
        )
        return fixture_document(
            fixture_id="profile-change-rust-to-go-cli-tui",
            classification="exact",
            rationale=None,
            profile="go-cli-tui",
            adoption_state="profile-change",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "manifest", "post-state"],
            source_tests=["tests/test_macro_profiles.py", "tests/test_profile_alignment.py"],
            input_document={"action": "profile-change", "fromProfile": "rust-cli", "profile": "go-cli-tui", "decisions": []},
            normalized_output={"preview": process_result(preview), "apply": process_result(applied)},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=manifest["managedArtifacts"],
            manifest=manifest,
            refusal=None,
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "successful profile-change post-state recorded"},
            blobs=blobs,
        )


def readoption_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        case = test_readoption_apply.ReadoptionApplyTests()
        repo, decision_file, _, _ = case.fixture(root)
        decision_document = json.loads(decision_file.read_text(encoding="utf-8"))
        before = capture_state(repo)
        preview = test_readoption_apply.run_apply_in_process(repo, decision_file)
        payload = json.loads(preview.stdout)
        applied = test_readoption_apply.run_apply_in_process(
            repo, decision_file, payload["planDigest"]
        )
        after = capture_state(repo)
        manifest = json.loads(
            (repo / "docs" / "agents" / "setup-context.json").read_text(
                encoding="utf-8"
            )
        )
        return fixture_document(
            fixture_id="readoption-preservation",
            classification="exact",
            rationale=None,
            profile="standard-typescript-monorepo",
            adoption_state="preservation",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "manifest", "managed-entry-ledger", "post-state"],
            source_tests=["tests/test_source_inventory.py", "tests/test_readoption_decisions.py", "tests/test_readoption_apply.py"],
            input_document={"action": "readoption", "profile": "standard-typescript-monorepo", "decisionDocument": decision_document},
            normalized_output={"preview": process_result(preview, ((str(root), "<fixture-root>"),)), "apply": process_result(applied, ((str(root), "<fixture-root>"),))},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=manifest["managedArtifacts"],
            manifest=manifest,
            refusal=None,
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "successful preservation post-state recorded"},
            blobs=blobs,
        )


def stale_plan_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        repo.mkdir()
        preview = run_in_process(
            apply_args(repo, "rust-cli", test_apply.BASE_DECISIONS)
        )
        payload = json.loads(preview.stdout)
        (repo / "AGENTS.md").write_text(
            "changed after preview\n", encoding="utf-8"
        )
        before = capture_state(repo)
        stale = run_in_process(
            apply_args(
                repo,
                "rust-cli",
                test_apply.BASE_DECISIONS,
                confirmation=payload["planDigest"],
            )
        )
        after = capture_state(repo)
        stale_payload = json.loads(stale.stdout)
        return fixture_document(
            fixture_id="stale-plan-refusal",
            classification="exact",
            rationale=None,
            profile="rust-cli",
            adoption_state="stale-input",
            contract_areas=["input", "state", "refusal", "digest", "planned-byte-sequence", "post-state"],
            source_tests=["tests/test_apply.py", "tests/test_readoption_apply.py"],
            input_document={"action": "apply", "profile": "rust-cli", "confirmPlan": payload["planDigest"], "mutationAfterPreview": {"path": "AGENTS.md", "identity": before["AGENTS.md"]["identity"]}},
            normalized_output={"preview": process_result(preview), "apply": process_result(stale)},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=[],
            manifest=None,
            refusal={"exitCode": stale.returncode, "findingCodes": [item["code"] for item in stale_payload["findings"]]},
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": True, "evidence": "stale input was never mutated"},
            blobs=blobs,
        )


def unsafe_carrier_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        test_source_inventory.SourceInventoryTests.write_incompatible_repository(
            repo
        )
        outside = root / "outside.md"
        outside.write_text("outside\n", encoding="utf-8")
        unsafe = repo / "docs" / "agents" / "linked.md"
        unsafe.symlink_to(Path("../../../outside.md"))
        before = capture_state(repo)
        try:
            inventory_incompatible_source_baseline(
                repo, "baseline.pre-0.0.1"
            )
        except SourceInventoryError as error:
            diagnostics = [
                {
                    "code": item.code,
                    "path": item.path.as_posix(),
                    "message": item.message,
                }
                for item in error.diagnostics
            ]
        else:
            raise AssertionError("unsafe carrier unexpectedly loaded")
        after = capture_state(repo)
        return fixture_document(
            fixture_id="unsafe-carrier-refusal",
            classification="exact",
            rationale=None,
            profile=None,
            adoption_state="unsafe-carrier",
            contract_areas=["input", "state", "refusal", "post-state"],
            source_tests=["tests/test_source_inventory.py", "tests/test_audit.py"],
            input_document={"action": "inventory", "declaredIdentity": "baseline.pre-0.0.1", "unsafePath": "docs/agents/linked.md"},
            normalized_output={"exitCode": 2, "diagnostics": diagnostics},
            before=before,
            after=after,
            planned_changes=[],
            plan_digest=None,
            managed_entry_ledger=[],
            manifest=None,
            refusal={"exitCode": 2, "findingCodes": [item["code"] for item in diagnostics]},
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": True, "evidence": "inventory refusal is read-only"},
            blobs=blobs,
        )


def missing_capability_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        repo = root / "repo"
        repo.mkdir()
        setup_standard_repository(repo, with_capabilities=False)
        decision_file = root / "decisions.json"
        decision_file.write_bytes(canonical_json(base_decision_document()))
        before = capture_state(repo)
        result = run_in_process(
            apply_args(
                repo,
                "standard-typescript-monorepo",
                decision_file=decision_file,
            )
        )
        after = capture_state(repo)
        payload = json.loads(result.stdout)
        return fixture_document(
            fixture_id="missing-capability-refusal",
            classification="exact",
            rationale=None,
            profile="standard-typescript-monorepo",
            adoption_state="capability-refusal",
            contract_areas=["input", "state", "refusal", "post-state"],
            source_tests=["tests/test_capabilities.py", "tests/test_readoption_apply.py"],
            input_document={"action": "apply", "profile": "standard-typescript-monorepo", "decisionDocument": base_decision_document()},
            normalized_output=process_result(result, ((str(root), "<fixture-root>"),)),
            before=before,
            after=after,
            planned_changes=payload.get("plannedChanges", []),
            plan_digest=payload.get("planDigest"),
            managed_entry_ledger=[],
            manifest=None,
            refusal={"exitCode": result.returncode, "findingCodes": [item["code"] for item in payload["findings"]]},
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": True, "evidence": "capability refusal is read-only"},
            blobs=blobs,
        )


def formatter_fixture(blobs):
    fixture_root = (
        TESTS_ROOT
        / "fixtures"
        / "formatter-compatibility"
        / "typescript-bun-monorepo"
    )
    golden = capture_state(fixture_root / "golden")
    provenance = json.loads(
        (fixture_root / "provenance.json").read_text(encoding="utf-8")
    )
    catalog = load_asset_catalog(SKILL_ROOT)
    contract = catalog.formatter_by_profile[
        "standard-typescript-monorepo"
    ]
    return fixture_document(
        fixture_id="formatter-composition",
        classification="exact",
        rationale=None,
        profile="standard-typescript-monorepo",
        adoption_state="formatter-composition",
        contract_areas=["input", "state", "action", "planned-byte-sequence", "post-state"],
        source_tests=["tests/test_formatter_compatibility.py"],
        input_document={
            "action": "format-compose",
            "profile": "standard-typescript-monorepo",
            "formatter": {
                "kind": contract.kind,
                "id": contract.formatter_id,
                "version": contract.version,
                "goldenDigest": contract.golden_digest,
            },
        },
        normalized_output={
            "provenance": provenance,
            "corpusDigest": test_formatter_compatibility.formatter_corpus_digest(
                contract
            ),
            "changedFiles": [],
        },
        before=golden,
        after=golden,
        planned_changes=[],
        plan_digest=None,
        managed_entry_ledger=[],
        manifest=provenance,
        refusal=None,
        rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "formatter and verification compose without byte changes"},
        blobs=blobs,
    )


def skill_restoration_fixture(blobs):
    with test_restore_skills.RestoreFixture() as fixture:
        files = {
            "SKILL.md": b"# restored\n",
            "references/guide.md": b"guide\n",
        }
        fixture.add_source_skill("agentic-cli-design", files)
        fixture.commit_and_configure("agentic-cli-design")
        shutil.rmtree(
            fixture.repo / ".agents" / "skills" / "agentic-cli-design"
        )
        before = capture_state(fixture.repo)
        preview = fixture.restore("--skill", "agentic-cli-design")
        payload = json.loads(preview.stdout)
        applied = fixture.restore(
            "--skill",
            "agentic-cli-design",
            "--confirm-plan",
            payload["planDigest"],
        )
        after = capture_state(fixture.repo)
        lock = json.loads(
            (fixture.repo / "skills-lock.json").read_text(encoding="utf-8")
        )
        replacements = ((str(fixture.root), "<fixture-root>"),)
        return fixture_document(
            fixture_id="skill-restoration",
            classification="exact",
            rationale=None,
            profile="rust-cli",
            adoption_state="skill-restoration",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "manifest", "post-state"],
            source_tests=["tests/test_restore_skills.py"],
            input_document={"action": "skills-restore", "profile": "rust-cli", "skills": ["agentic-cli-design"], "sourceRevision": fixture.revision},
            normalized_output={"preview": process_result(preview, replacements), "apply": process_result(applied, replacements)},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=payload["skills"],
            manifest=lock,
            refusal=None,
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "restored post-state and portable lock recorded"},
            blobs=blobs,
        )


def skill_restoration_rollback_fixture(blobs):
    with test_restore_skills.RestoreFixture() as fixture:
        fixture.add_source_skill(
            "agentic-cli-design", {"SKILL.md": b"# restored\n"}
        )
        fixture.commit_and_configure("agentic-cli-design")
        target = (
            fixture.repo
            / ".agents"
            / "skills"
            / "agentic-cli-design"
            / "SKILL.md"
        )
        target.write_bytes(b"# drifted\n")
        fixture.write_lock({"unrelated": {"computedHash": "c" * 64}})
        catalog = load_asset_catalog(fixture.skill_root)
        plan = context_setup.build_restore_plan(
            fixture.repo,
            catalog,
            "rust-cli",
            ["agentic-cli-design"],
            fixture.source,
        )
        before = capture_state(fixture.repo)
        try:
            context_setup.apply_restore_plan(
                fixture.repo,
                plan,
                filesystem=test_restore_skills.FailingFilesystem("lock"),
            )
        except context_setup.RestoreError as error:
            normalized_error = {
                "code": error.code,
                "message": error.message,
                "action": error.action,
                "exitCode": error.exit_code,
            }
        else:
            raise AssertionError("injected restore failure unexpectedly succeeded")
        after = capture_state(fixture.repo)
        plan_payload = plan.to_json(ok=False, applied=False)
        return fixture_document(
            fixture_id="skill-restoration-rollback",
            classification="exact",
            rationale=None,
            profile="rust-cli",
            adoption_state="skill-restoration-rollback",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "rollback-outcome"],
            source_tests=["tests/test_restore_skills.py"],
            input_document={"action": "skills-restore", "profile": "rust-cli", "skills": ["agentic-cli-design"], "failurePoint": "lock"},
            normalized_output={"error": normalized_error},
            before=before,
            after=after,
            planned_changes=plan_payload["plannedChanges"],
            plan_digest=plan_payload["planDigest"],
            managed_entry_ledger=plan_payload["skills"],
            manifest=None,
            refusal={"exitCode": normalized_error["exitCode"], "findingCodes": [normalized_error["code"]]},
            rollback={"attempted": True, "succeeded": content_state_equal(before, after), "restoredPreimage": content_state_equal(before, after), "evidence": {"before": content_state_digest(before), "after": content_state_digest(after), "contract": "paths, kinds, targets, and exact bytes"}},
            blobs=blobs,
        )


def asset_sync_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        skill_root = root / "skill"
        checkout = root / "canonical"
        clone_assets_to(SKILL_ROOT, skill_root)
        source_dir, revision = test_sync_setups.write_git_setup_sources(
            checkout, skill_root
        )
        before = capture_state(skill_root / "assets" / "setups")
        result, invalid_input = context_setup.sync_setup_snapshots(
            skill_root, source_dir, check=False
        )
        after = capture_state(skill_root / "assets" / "setups")
        snapshots = [
            json.loads(path.read_text(encoding="utf-8"))
            for path in sorted((skill_root / "assets" / "setups").glob("*.json"))
        ]
        changes = [
            {
                "action": "refresh",
                "path": path,
                "managedId": f"setup-snapshot.{Path(path).stem}",
            }
            for path in sorted(set(before) | set(after))
            if before.get(path) != after.get(path)
        ]
        return fixture_document(
            fixture_id="asset-sync",
            classification="exact",
            rationale=None,
            profile=None,
            adoption_state="asset-sync",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "manifest", "post-state"],
            source_tests=["tests/test_sync_setups.py"],
            input_document={"action": "assets-sync", "sourceRepository": "example/skills", "sourceRevision": revision},
            normalized_output={"invalidInput": invalid_input, "result": result.to_json()},
            before=before,
            after=after,
            planned_changes=changes,
            plan_digest=None,
            managed_entry_ledger=[{"id": item["id"], "digest": item["digest"]} for item in snapshots],
            manifest={"setups": snapshots},
            refusal=None,
            rollback={"attempted": False, "succeeded": None, "restoredPreimage": None, "evidence": "second sync is byte-idempotent in the owning Python test"},
            blobs=blobs,
        )


def atomic_rollback_fixture(blobs):
    with tempfile.TemporaryDirectory() as temp_dir:
        root = Path(temp_dir)
        case = test_readoption_apply.ReadoptionApplyTests()
        repo, decision_file, _, _ = case.fixture(root)
        preview = test_readoption_apply.run_apply_in_process(repo, decision_file)
        payload = json.loads(preview.stdout)
        before = capture_state(repo)
        original_replace = context_setup.Path.replace
        writes = 0

        def replace_with_failure(source, target):
            nonlocal writes
            result = original_replace(source, target)
            writes += 1
            if writes == 2:
                raise OSError("injected second-write failure")
            return result

        with mock.patch.object(
            context_setup.Path,
            "replace",
            autospec=True,
            side_effect=replace_with_failure,
        ):
            failed = test_readoption_apply.run_apply_in_process(
                repo, decision_file, payload["planDigest"]
            )
        after = capture_state(repo)
        failed_payload = json.loads(failed.stdout)
        return fixture_document(
            fixture_id="atomic-rollback",
            classification="exact",
            rationale=None,
            profile="standard-typescript-monorepo",
            adoption_state="rollback",
            contract_areas=["input", "state", "action", "digest", "planned-byte-sequence", "refusal", "rollback-outcome"],
            source_tests=["tests/test_apply.py", "tests/test_readoption_apply.py"],
            input_document={"action": "readoption-apply", "profile": "standard-typescript-monorepo", "confirmPlan": payload["planDigest"], "failurePoint": "second-replace"},
            normalized_output={"preview": process_result(preview, ((str(root), "<fixture-root>"),)), "apply": process_result(failed, ((str(root), "<fixture-root>"),))},
            before=before,
            after=after,
            planned_changes=payload["plannedChanges"],
            plan_digest=payload["planDigest"],
            managed_entry_ledger=[],
            manifest=None,
            refusal={"exitCode": failed.returncode, "findingCodes": [item["code"] for item in failed_payload["findings"]]},
            rollback={"attempted": True, "succeeded": content_state_equal(before, after), "restoredPreimage": content_state_equal(before, after), "evidence": {"before": content_state_digest(before), "after": content_state_digest(after), "contract": "paths, kinds, targets, and exact bytes"}},
            blobs=blobs,
        )


def matrix_document():
    python_rows = python_test_inventory()
    delta_rows = [
        {
            **row,
            "sourceKind": "designed-delta",
            "sourcePath": "docs/specs/0046-public-context-driven-baseline-command/_techspec.md",
            "sourceSuite": None,
            "sourceTest": None,
            "classification": "designed-delta",
        }
        for row in DESIGNED_DELTA_ROWS
    ]
    rows = sorted([*python_rows, *delta_rows], key=lambda row: row["id"])
    counts = {
        classification: sum(
            row["classification"] == classification for row in rows
        )
        for classification in ALLOWED_CLASSIFICATIONS
    }
    return {
        "schemaVersion": MATRIX_SCHEMA_VERSION,
        "sourceTestCount": len(python_rows),
        "rowCount": len(rows),
        "classifications": list(ALLOWED_CLASSIFICATIONS),
        "classificationCounts": counts,
        "rows": rows,
    }


def generate_corpus():
    blobs = {}
    with deterministic_runtime():
        fixtures = [
            successful_apply_fixture(
                blobs,
                fixture_id="greenfield-rust-cli",
                profile="rust-cli",
                adoption_state="greenfield",
                decisions=test_apply.BASE_DECISIONS,
                source_tests=["tests/test_apply.py", "tests/test_macro_profiles.py"],
            ),
            successful_apply_fixture(
                blobs,
                fixture_id="greenfield-go-cli-tui",
                profile="go-cli-tui",
                adoption_state="greenfield",
                decisions=test_apply.BASE_DECISIONS,
                source_tests=["tests/test_apply.py", "tests/test_macro_profiles.py"],
            ),
            successful_apply_fixture(
                blobs,
                fixture_id="greenfield-standard-typescript-monorepo",
                profile="standard-typescript-monorepo",
                adoption_state="greenfield",
                setup=setup_standard_repository,
                decision_document=base_decision_document(),
                source_tests=["tests/test_standard_typescript_monorepo.py", "tests/test_macro_profiles.py"],
            ),
            update_fixture(blobs),
            profile_change_fixture(blobs),
            readoption_fixture(blobs),
            stale_plan_fixture(blobs),
            unsafe_carrier_fixture(blobs),
            missing_capability_fixture(blobs),
            formatter_fixture(blobs),
            skill_restoration_fixture(blobs),
            skill_restoration_rollback_fixture(blobs),
            asset_sync_fixture(blobs),
            atomic_rollback_fixture(blobs),
        ]
    matrix = matrix_document()
    files = {
        "matrix.json": canonical_json(matrix),
        "blobs.json": canonical_json(
            {
                "schemaVersion": "setup-context-driven/parity-blobs/v1",
                "encoding": "base64",
                "blobs": dict(sorted(blobs.items())),
            }
        ),
    }
    for fixture in fixtures:
        files[f"fixtures/{fixture['id']}.json"] = canonical_json(fixture)
    artifact_records = [
        {
            "path": path,
            "sha256": sha256_bytes(content),
            "bytes": len(content),
        }
        for path, content in sorted(files.items())
    ]
    files["manifest.json"] = canonical_json(
        {
            "schemaVersion": SCHEMA_VERSION,
            "frozenDate": FIXED_DATE,
            "sourceSuite": {
                "root": "tests",
                "testCount": matrix["sourceTestCount"],
                "inventoryDigest": "sha256:"
                + sha256_bytes(
                    canonical_json(
                        [
                            row["id"]
                            for row in matrix["rows"]
                            if row["sourceKind"] == "python-test"
                        ]
                    )
                ),
            },
            "profiles": [
                "go-cli-tui",
                "rust-cli",
                "standard-typescript-monorepo",
            ],
            "adoptionStates": sorted(
                {fixture["adoptionState"] for fixture in fixtures}
            ),
            "fixtureIds": [fixture["id"] for fixture in fixtures],
            "retiredBehavior": [],
            "artifacts": artifact_records,
        }
    )
    return dict(sorted(files.items()))


def write_corpus():
    files = generate_corpus()
    for relative, content in files.items():
        target = CORPUS_ROOT / relative
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(content)
    return files


if __name__ == "__main__":
    generated = write_corpus()
    print(f"wrote {len(generated)} parity corpus files to {CORPUS_ROOT}")
