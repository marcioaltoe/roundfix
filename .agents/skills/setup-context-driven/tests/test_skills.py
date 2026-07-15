import copy
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SKILL_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = SKILL_ROOT / "scripts" / "context_setup.py"
sys.path.insert(0, str(SKILL_ROOT / "scripts"))
sys.path.insert(0, str(Path(__file__).resolve().parent))

from context_assets import load_asset_catalog  # noqa: E402
from context_setup import validate_profile_skill_references  # noqa: E402
from test_audit import (  # noqa: E402
    install_profile_skills,
    snapshot_files,
    write_compliant_repository,
    write_skill,
)


class SkillAuditTests(unittest.TestCase):
    def test_missing_required_skill_blocks_compliance(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli", install_skills=False)
            install_profile_skills(repo, "rust-cli", omit={"agentic-cli-design"})

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 1)
            finding = self.finding(result, "skills.required.missing")
            self.assertEqual(finding["severity"], "error")
            self.assertEqual(finding["path"], ".agents/skills/agentic-cli-design")
            self.assertIn("Install the rust-cli canonical skill setup", finding["action"])

    def test_extra_locked_skills_are_informational_only_when_requested(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            write_skill(repo, "autoresearch")
            write_lockfile(repo, ["autoresearch"])
            before = snapshot_files(repo)

            hidden = run_audit(repo, "--format", "json")
            shown = run_audit(repo, "--format", "json", "--show-extra-skills")

            self.assertEqual(hidden.returncode, 0, hidden.stderr)
            self.notFinding(hidden, "skills.extra.installed")
            self.assertEqual(shown.returncode, 0, shown.stderr)
            finding = self.finding(shown, "skills.extra.installed")
            self.assertEqual(finding["severity"], "info")
            self.assertNotIn("rm ", finding["action"])
            self.assertNotIn("delete", finding["action"].lower())
            self.assertEqual(snapshot_files(repo), before)

    def test_local_and_untracked_skills_are_not_removal_candidates(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            write_skill(repo, "repo-local")
            write_skill(repo, "local-review")
            manifest_path = repo / "docs" / "agents" / "setup-context.json"
            manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
            manifest["localSkills"] = ["repo-local"]
            manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")

            result = run_audit(repo, "--format", "json", "--show-extra-skills")

            self.assertEqual(result.returncode, 0, result.stderr)
            self.notFinding(result, "skills.extra.installed")
            untracked = self.finding(result, "skills.local.untracked")
            self.assertEqual(untracked["managedId"], "local-review")
            self.assertNotIn("repo-local", json.dumps(json.loads(result.stdout)["findings"]))
            self.assertNotIn("rm ", untracked["action"])
            self.assertNotIn("delete", untracked["action"].lower())

    def test_malformed_lockfile_is_invalid_input_without_writes(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "rust-cli")
            lock_path = repo / "skills-lock.json"
            lock_path.write_text("{not json", encoding="utf-8")
            before = snapshot_files(repo)

            result = run_audit(repo, "--format", "json", "--show-extra-skills")

            self.assertEqual(result.returncode, 2)
            finding = self.finding(result, "skills.lockfile.invalid")
            self.assertEqual(finding["path"], "skills-lock.json")
            self.assertEqual(snapshot_files(repo), before)

    def test_module_skill_reference_outside_setup_is_blocking(self):
        catalog = load_asset_catalog(SKILL_ROOT)
        mutated = copy.deepcopy(catalog)
        setup_id = mutated.profiles["rust-cli"]["setup"]
        mutated.setups[setup_id]["skills"] = [
            skill
            for skill in mutated.setups[setup_id]["skills"]
            if skill["name"] != "agentic-cli-design"
        ]
        findings = []

        validate_profile_skill_references(
            mutated,
            "rust-cli",
            mutated.ordered_modules_by_profile["rust-cli"],
            findings,
        )

        matches = [finding for finding in findings if finding.code == "skills.reference.outside-setup"]
        self.assertGreaterEqual(len(matches), 1)
        self.assertEqual(matches[0].severity, "error")
        self.assertEqual(matches[0].managed_id, "agentic-cli-design")

    def test_audit_is_portable_without_canonical_setups_directory(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            repo = Path(temp_dir)
            write_compliant_repository(repo, "go-cli-tui")

            result = run_audit(repo, "--format", "json")

            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual(json.loads(result.stdout)["findings"], [])

    def finding(self, result, code):
        payload = json.loads(result.stdout)
        matches = [finding for finding in payload["findings"] if finding["code"] == code]
        self.assertGreater(len(matches), 0, payload)
        return matches[0]

    def notFinding(self, result, code):
        payload = json.loads(result.stdout)
        self.assertEqual(
            [finding for finding in payload["findings"] if finding["code"] == code],
            [],
            payload,
        )


def write_lockfile(repo, skill_names):
    payload = {
        "version": 1,
        "skills": {
            name: {
                "source": "marcioaltoe/skills",
                "sourceType": "github",
                "skillPath": f"skills/{name}/SKILL.md",
                "computedHash": "0" * 64,
            }
            for name in skill_names
        },
    }
    (repo / "skills-lock.json").write_text(
        json.dumps(payload, indent=2) + "\n",
        encoding="utf-8",
    )


def run_audit(repo, *args):
    return subprocess.run(
        [sys.executable, str(SCRIPT), "audit", "--repo", str(repo), *args],
        text=True,
        capture_output=True,
        check=False,
    )


if __name__ == "__main__":
    unittest.main()
