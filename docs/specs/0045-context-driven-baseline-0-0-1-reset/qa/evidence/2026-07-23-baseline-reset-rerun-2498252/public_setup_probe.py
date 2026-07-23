"""Public setup journeys for the Spec 0045 QA rerun."""

from __future__ import annotations

import hashlib
import json
import sys
import tempfile
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[6]
TESTS = REPO_ROOT / ".agents/skills/setup-context-driven/tests"
sys.path.insert(0, str(TESTS))

from test_audit import install_profile_skills, run_context_setup, snapshot_files  # noqa: E402
from test_readoption_apply import (  # noqa: E402
    MANIFEST,
    PROFILE,
    REPOSITORY_RULES,
    ReadoptionApplyTests,
    run_apply,
)


def finding_codes(payload: dict) -> list[str]:
    return [item["code"] for item in payload["findings"]]


def make_clean_fixture(root: Path) -> tuple[Path, Path, ReadoptionApplyTests]:
    case = ReadoptionApplyTests()
    repo = root / "repo"
    repo.mkdir(parents=True)
    install_profile_skills(repo, PROFILE)
    (repo / "DESIGN.md").write_text("# Design contract\n", encoding="utf-8")
    http_path = repo / "docs/architecture/http-contract.json"
    http_path.parent.mkdir(parents=True)
    http_path.write_text('{"mode":"REST"}\n', encoding="utf-8")
    decision_file = root / "decisions.json"
    decision_file.write_text(
        json.dumps(
            {
                "schemaVersion": "setup-context-driven/decisions/0.0.1",
                "version": "0.0.1",
                "decisions": case.profile_decisions(),
            },
            indent=2,
        )
        + "\n",
        encoding="utf-8",
    )
    return repo, decision_file, case


def clean_adoption_journey(root: Path) -> dict:
    repo, decision_file, case = make_clean_fixture(root)
    initial = snapshot_files(repo)
    blocked = run_apply(repo, decision_file)
    blocked_payload = json.loads(blocked.stdout)
    assert blocked.returncode == 1
    assert blocked.stderr == ""
    assert "capability.required.missing" in finding_codes(blocked_payload)
    assert "plan.confirmation.required" not in finding_codes(blocked_payload)
    assert snapshot_files(repo) == initial

    case.write_capability_evidence(repo)
    before_preview = snapshot_files(repo)
    preview = run_apply(repo, decision_file)
    preview_payload = json.loads(preview.stdout)
    assert preview.returncode == 3
    assert preview.stderr == ""
    assert "plan.confirmation.required" in finding_codes(preview_payload)
    assert snapshot_files(repo) == before_preview
    assert (
        preview_payload["setupSnapshot"]["schemaVersion"]
        == "setup-context-driven/profile-snapshot/0.0.1"
    )

    invalid = run_apply(repo, decision_file, "not-a-digest")
    invalid_payload = json.loads(invalid.stdout)
    assert invalid.returncode == 2
    assert invalid.stderr == ""
    assert "plan.confirmation.invalid" in finding_codes(invalid_payload)
    assert snapshot_files(repo) == before_preview

    applied = run_apply(repo, decision_file, preview_payload["planDigest"])
    applied_payload = json.loads(applied.stdout)
    assert applied.returncode == 0
    assert applied.stderr == ""
    assert applied_payload["planDigest"] == preview_payload["planDigest"]
    assert applied_payload["plannedOutputs"] == preview_payload["plannedOutputs"]

    manifest = json.loads((repo / MANIFEST).read_text(encoding="utf-8"))
    frontend = (repo / "docs/agents/frontend.md").read_text(encoding="utf-8")
    backend = (repo / "docs/agents/backend.md").read_text(encoding="utf-8")
    architecture = {
        "frontendDomainSystems": "organize frontend feature code by domain system"
        in frontend,
        "frontendPublicBoundary": "one public boundary" in frontend,
        "backendLayers": "domain, application, and infrastructure layers" in backend,
        "thinHandlers": "keep HTTP handlers thin" in backend,
        "httpIndependentUseCases": "independent of HTTP" in backend,
        "infrastructurePersistence": "persistence implementation in infrastructure"
        in backend,
    }
    assert all(architecture.values())
    assert manifest["schemaVersion"] == "setup-context-driven/manifest/0.0.1"
    assert manifest["version"] == "0.0.1"
    assert manifest["generator"]["version"] == "0.0.1"
    assert (
        manifest["generator"]["baseline"]
        == "baseline.standard-typescript-monorepo-0.0.1"
    )
    assert manifest["setupSnapshot"] == preview_payload["setupSnapshot"]

    settled = snapshot_files(repo)
    audit_one = run_context_setup("audit", "--repo", str(repo), "--format", "json")
    audit_two = run_context_setup("audit", "--repo", str(repo), "--format", "json")
    audit_text = run_context_setup("audit", "--repo", str(repo), "--format", "text")
    reapplied = run_apply(repo, decision_file)
    reapplied_payload = json.loads(reapplied.stdout)
    assert audit_one.returncode == audit_two.returncode == audit_text.returncode == 0
    assert audit_one.stderr == audit_two.stderr == audit_text.stderr == ""
    assert audit_text.stdout.strip() == "setup-context-driven audit: ok"
    assert reapplied.returncode == 0
    assert reapplied.stderr == ""
    assert reapplied_payload["plannedChanges"] == []
    assert snapshot_files(repo) == settled

    return {
        "blocked": {
            "exit": blocked.returncode,
            "findings": finding_codes(blocked_payload),
            "bytesUnchanged": True,
        },
        "preview": {
            "exit": preview.returncode,
            "findings": finding_codes(preview_payload),
            "schema": preview_payload["setupSnapshot"]["schemaVersion"],
            "capabilityCount": len(preview_payload["capabilities"]),
            "plannedOutputCount": len(preview_payload["plannedOutputs"]),
            "bytesUnchanged": True,
        },
        "invalidConfirmation": {
            "exit": invalid.returncode,
            "findings": finding_codes(invalid_payload),
            "bytesUnchanged": True,
        },
        "apply": {
            "exit": applied.returncode,
            "digest": applied_payload["planDigest"],
            "previewApplyParity": True,
            "manifestSchema": manifest["schemaVersion"],
            "manifestVersion": manifest["version"],
            "generatorVersion": manifest["generator"]["version"],
            "baseline": manifest["generator"]["baseline"],
            "architecture": architecture,
        },
        "persistence": {
            "auditExits": [audit_one.returncode, audit_two.returncode],
            "textAudit": audit_text.stdout.strip(),
            "reapplyExit": reapplied.returncode,
            "plannedChanges": reapplied_payload["plannedChanges"],
            "bytesUnchanged": True,
        },
    }


def stale_confirmation_journey(root: Path) -> dict:
    repo, decision_file, case = make_clean_fixture(root)
    case.write_capability_evidence(repo)
    preview = run_apply(repo, decision_file)
    preview_payload = json.loads(preview.stdout)
    assert preview.returncode == 3
    package_path = repo / "package.json"
    package_path.write_text(
        package_path.read_text(encoding="utf-8") + "\n", encoding="utf-8"
    )
    changed = snapshot_files(repo)
    stale = run_apply(repo, decision_file, preview_payload["planDigest"])
    stale_payload = json.loads(stale.stdout)
    assert stale.returncode == 3
    assert stale.stderr == ""
    assert "plan.confirmation.stale" in finding_codes(stale_payload)
    assert snapshot_files(repo) == changed
    return {
        "exit": stale.returncode,
        "findings": finding_codes(stale_payload),
        "bytesUnchanged": True,
        "originalDigest": preview_payload["planDigest"],
        "recomputedDigest": stale_payload["planDigest"],
        "digestChanged": preview_payload["planDigest"] != stale_payload["planDigest"],
    }


def readoption_journey(root: Path) -> dict:
    case = ReadoptionApplyTests()
    repo, decision_file, original_guide, proposed_rules = case.fixture(root)
    before = snapshot_files(repo)
    preview = run_apply(repo, decision_file)
    preview_payload = json.loads(preview.stdout)
    assert preview.returncode == 3
    assert preview.stderr == ""
    assert "plan.confirmation.required" in finding_codes(preview_payload)
    assert snapshot_files(repo) == before

    applied = run_apply(repo, decision_file, preview_payload["planDigest"])
    assert applied.returncode == 0
    assert applied.stderr == ""
    assert (repo / REPOSITORY_RULES).read_bytes() == proposed_rules
    assert (repo / "docs/agents/guide.md").read_bytes() == original_guide

    repository_owned = proposed_rules + b"\nMaintainer-owned follow-up.\n"
    (repo / REPOSITORY_RULES).write_bytes(repository_owned)
    owned_digest = hashlib.sha256(repository_owned).hexdigest()
    before_reapply = snapshot_files(repo)
    reapplied = run_apply(repo, decision_file)
    reapplied_payload = json.loads(reapplied.stdout)
    audit = run_context_setup("audit", "--repo", str(repo), "--format", "json")
    assert reapplied.returncode == 0
    assert reapplied.stderr == ""
    assert reapplied_payload["plannedChanges"] == []
    assert audit.returncode == 0
    assert audit.stderr == ""
    assert (repo / REPOSITORY_RULES).read_bytes() == repository_owned
    assert snapshot_files(repo) == before_reapply
    return {
        "previewExit": preview.returncode,
        "inventoryCount": preview_payload["sourceBaseline"]["entryCount"],
        "dispositionCount": len(
            preview_payload["decisionDocument"]["readoption"]["dispositions"]
        ),
        "applyExit": applied.returncode,
        "auditExit": audit.returncode,
        "reapplyExit": reapplied.returncode,
        "plannedChanges": reapplied_payload["plannedChanges"],
        "repositoryOwnedDigest": owned_digest,
        "repositoryOwnedBytesPreserved": True,
        "typedDocumentBytesPreserved": True,
    }


def main() -> int:
    with tempfile.TemporaryDirectory(prefix="roundfix-qa0045-public-") as temp_dir:
        root = Path(temp_dir)
        result = {
            "cleanAdoption": clean_adoption_journey(root / "clean"),
            "staleConfirmation": stale_confirmation_journey(root / "stale"),
            "readoption": readoption_journey(root / "readoption"),
        }
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
