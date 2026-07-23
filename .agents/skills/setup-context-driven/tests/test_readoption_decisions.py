"""Behavior tests for explicit Baseline Readoption decisions.

Suite: source-bound Readoption Decision Plan preview
Invariant: every Source Baseline Entry has one explicit typed disposition before mutation.
Boundary IN: decision-file parsing, disposition validation, preview, and digest binding.
Boundary OUT: confirmed Readoption mutation and rollback (test_readoption_apply.py).
"""

import base64
import hashlib
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_baseline import inventory_incompatible_source_baseline  # noqa: E402
from test_audit import snapshot_files  # noqa: E402
from test_source_inventory import SourceInventoryTests  # noqa: E402


DECISION_SCHEMA = "setup-context-driven/decisions/0.0.1"
REPOSITORY_RULES = Path("docs/agents/repository-rules.md")


class ReadoptionDecisionTests(unittest.TestCase):
    def test_complete_decisions_preview_exact_repository_rules_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            inventory = self.inventory(repo)
            proposed = b"# Repository-specific normative rules\n\nKeep local policy.\n"
            document = self.complete_document(inventory, proposed=proposed)
            decision_file = self.write_decision_file(repo, document)
            before = snapshot_files(repo)

            first = run_audit(repo, decision_file)
            second = run_audit(repo, decision_file)

            self.assertEqual(first.returncode, 3, first.stderr)
            self.assertEqual(second.returncode, 3, second.stderr)
            self.assertEqual(snapshot_files(repo), before)
            first_payload = json.loads(first.stdout)
            self.assertEqual(first_payload, json.loads(second.stdout))
            self.assertRegex(first_payload["planDigest"], r"^[0-9a-f]{64}$")
            self.assertNotIn(
                "readoption.disposition.required",
                {item["code"] for item in first_payload["findings"]},
            )
            self.assertIn(
                "plan.confirmation.required",
                {item["code"] for item in first_payload["findings"]},
            )
            change = self.repository_rules_change(first_payload)
            self.assertEqual(change["path"], REPOSITORY_RULES.as_posix())
            self.assertEqual(change["afterDigest"], hashlib.sha256(proposed).hexdigest())
            normalized = first_payload["decisionDocument"]
            rule_destination = normalized["readoption"]["dispositions"][0]["destination"]
            self.assertEqual(rule_destination["proposedBytes"], self.encoded(proposed))
            self.assertEqual(rule_destination["digest"], hashlib.sha256(proposed).hexdigest())

    def test_incomplete_duplicate_unknown_stale_and_invalid_entries_fail_specifically(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            inventory = self.inventory(repo)
            complete = self.complete_document(inventory)
            first = complete["readoption"]["dispositions"][0]
            cases = []

            missing = json.loads(json.dumps(complete))
            missing["readoption"]["dispositions"].pop(0)
            cases.append(("missing", missing, "readoption.disposition.missing", first["entryId"]))

            duplicate = json.loads(json.dumps(complete))
            duplicate["readoption"]["dispositions"].append(dict(first))
            cases.append(("duplicate", duplicate, "readoption.disposition.duplicate", first["entryId"]))

            unknown = json.loads(json.dumps(complete))
            extra = dict(first)
            extra["entryId"] = "source-entry." + "f" * 64
            unknown["readoption"]["dispositions"].append(extra)
            cases.append(("unknown", unknown, "readoption.disposition.unknown", extra["entryId"]))

            stale = json.loads(json.dumps(complete))
            stale["readoption"]["dispositions"][0]["entryDigest"] = "0" * 64
            cases.append(("stale", stale, "readoption.disposition.stale", first["entryId"]))

            invalid = json.loads(json.dumps(complete))
            invalid["readoption"]["dispositions"][0]["classification"] = "inferred-rule"
            cases.append(("invalid", invalid, "readoption.disposition.invalid", first["entryId"]))

            for name, document, code, entry_id in cases:
                with self.subTest(name=name):
                    decision_file = self.write_decision_file(repo, document, name)
                    before = snapshot_files(repo)
                    result = run_audit(repo, decision_file)
                    self.assertEqual(result.returncode, 2, result.stderr)
                    self.assertEqual(snapshot_files(repo), before)
                    finding = self.finding(json.loads(result.stdout), code)
                    self.assertEqual(finding["managedId"], entry_id)

    def test_typed_document_and_repository_rules_destinations_fail_closed(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            inventory = self.inventory(repo)
            typed_path = repo / "docs" / "agents" / "owned-policy.md"
            typed_bytes = b"# Owned policy\n"
            typed_path.write_bytes(typed_bytes)
            inventory = self.inventory(repo)
            complete = self.complete_document(inventory)
            disposition = complete["readoption"]["dispositions"][0]
            disposition["classification"] = "normative-clause"
            disposition["disposition"] = "repository-document"
            disposition["destination"] = {
                "documentType": "agent-guide",
                "path": "docs/agents/owned-policy.md",
                "digest": hashlib.sha256(typed_bytes).hexdigest(),
            }
            disposition["reason"] = ""
            accepted = self.write_decision_file(repo, complete, "typed")

            result = run_audit(repo, accepted)

            self.assertEqual(result.returncode, 0, result.stderr)
            payload = json.loads(result.stdout)
            normalized = payload["decisionDocument"]["readoption"]["dispositions"][0]
            self.assertEqual(normalized["destination"], disposition["destination"])
            self.assertFalse(
                any(change["path"] == disposition["destination"]["path"] for change in payload["plannedChanges"])
            )

            mutations = []
            unsafe = json.loads(json.dumps(complete))
            unsafe["readoption"]["dispositions"][0]["destination"]["path"] = "../escape.md"
            mutations.append((unsafe, "readoption.destination.path.unsafe"))
            untyped = json.loads(json.dumps(complete))
            untyped["readoption"]["dispositions"][0]["destination"]["documentType"] = "misc"
            mutations.append((untyped, "readoption.destination.document-type.invalid"))
            changed = json.loads(json.dumps(complete))
            changed["readoption"]["dispositions"][0]["destination"]["digest"] = "0" * 64
            mutations.append((changed, "readoption.destination.stale"))
            for index, (document, code) in enumerate(mutations):
                with self.subTest(code=code):
                    path = self.write_decision_file(repo, document, f"unsafe-{index}")
                    rejected = run_audit(repo, path)
                    self.assertEqual(rejected.returncode, 2, rejected.stderr)
                    self.finding(json.loads(rejected.stdout), code)

    def test_existing_repository_rules_are_preserved_without_confirmation(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            target = repo / REPOSITORY_RULES
            target.write_bytes(b"repository-owned bytes\n")
            inventory = self.inventory(repo)
            document = self.complete_document(
                inventory,
                proposed=b"a proposal that must not replace existing bytes\n",
            )
            decision_file = self.write_decision_file(repo, document)
            before = snapshot_files(repo)

            result = run_audit(repo, decision_file)

            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual(snapshot_files(repo), before)
            payload = json.loads(result.stdout)
            self.assertNotIn(
                "plan.confirmation.required",
                {item["code"] for item in payload["findings"]},
            )
            self.assertFalse(
                any(change["path"] == REPOSITORY_RULES.as_posix() for change in payload["plannedChanges"])
            )

    def test_normalized_disposition_and_target_change_plan_digest(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            SourceInventoryTests.write_incompatible_repository(repo)
            first_path = repo / "docs" / "agents" / "first-policy.md"
            second_path = repo / "docs" / "agents" / "second-policy.md"
            target_bytes = b"# Existing typed policy\n"
            first_path.write_bytes(target_bytes)
            second_path.write_bytes(target_bytes)
            inventory = self.inventory(repo)
            base = self.complete_document(inventory)
            disposition = base["readoption"]["dispositions"][0]
            disposition.update(
                {
                    "classification": "normative-clause",
                    "disposition": "repository-document",
                    "destination": {
                        "documentType": "agent-guide",
                        "path": "docs/agents/first-policy.md",
                        "digest": hashlib.sha256(target_bytes).hexdigest(),
                    },
                    "reason": "",
                }
            )
            first = run_audit(repo, self.write_decision_file(repo, base, "first"))

            classification_changed = json.loads(json.dumps(base))
            classification_changed["readoption"]["dispositions"][0]["classification"] = "recommendation"
            second = run_audit(
                repo,
                self.write_decision_file(repo, classification_changed, "classification"),
            )

            target_changed = json.loads(json.dumps(base))
            target_changed["readoption"]["dispositions"][0]["destination"]["path"] = (
                "docs/agents/second-policy.md"
            )
            third = run_audit(repo, self.write_decision_file(repo, target_changed, "target"))

            self.assertEqual(first.returncode, 0, first.stderr)
            self.assertEqual(second.returncode, 0, second.stderr)
            self.assertEqual(third.returncode, 0, third.stderr)
            digests = {
                json.loads(first.stdout)["planDigest"],
                json.loads(second.stdout)["planDigest"],
                json.loads(third.stdout)["planDigest"],
            }
            self.assertEqual(len(digests), 3)

    def inventory(self, repo):
        return inventory_incompatible_source_baseline(repo, "baseline.pre-0.0.1")

    def complete_document(self, inventory, proposed=None):
        dispositions = []
        for index, entry in enumerate(inventory.entries):
            if proposed is not None and index == 0:
                destination = {
                    "documentType": "repository-rules",
                    "path": REPOSITORY_RULES.as_posix(),
                    "proposedBytes": self.encoded(proposed),
                    "digest": hashlib.sha256(proposed).hexdigest(),
                }
                dispositions.append(
                    {
                        "entryId": entry.entry_id,
                        "entryDigest": entry.digest,
                        "classification": "normative-clause",
                        "disposition": "repository-rules",
                        "destination": destination,
                        "reason": "",
                    }
                )
                continue
            dispositions.append(
                {
                    "entryId": entry.entry_id,
                    "entryDigest": entry.digest,
                    "classification": "non-governed",
                    "disposition": "rejected",
                    "destination": None,
                    "reason": "Structural evidence is not a governed instruction.",
                }
            )
        return {
            "schemaVersion": DECISION_SCHEMA,
            "version": "0.0.1",
            "decisions": [],
            "readoption": {
                "sourceBaseline": {
                    "id": inventory.baseline_id,
                    "digest": inventory.digest,
                },
                "dispositions": dispositions,
            },
        }

    @staticmethod
    def write_decision_file(repo, document, suffix="decisions"):
        path = repo.parent / f"{repo.name}-{suffix}.json"
        path.write_text(json.dumps(document, indent=2) + "\n", encoding="utf-8")
        return path

    @staticmethod
    def encoded(content):
        return base64.b64encode(content).decode("ascii")

    @staticmethod
    def finding(payload, code):
        for item in payload["findings"]:
            if item["code"] == code:
                return item
        raise AssertionError(f"missing finding {code}: {payload['findings']}")

    @staticmethod
    def repository_rules_change(payload):
        for change in payload["plannedChanges"]:
            if change["path"] == REPOSITORY_RULES.as_posix():
                return change
        raise AssertionError(f"missing repository-rules preview: {payload['plannedChanges']}")


def run_audit(repo, decision_file):
    env = os.environ.copy()
    env["PYTHONDONTWRITEBYTECODE"] = "1"
    return subprocess.run(
        [
            sys.executable,
            str(SCRIPT),
            "audit",
            "--repo",
            str(repo),
            "--format",
            "json",
            "--decision-file",
            str(decision_file),
        ],
        text=True,
        capture_output=True,
        check=False,
        env=env,
    )


if __name__ == "__main__":
    unittest.main()
