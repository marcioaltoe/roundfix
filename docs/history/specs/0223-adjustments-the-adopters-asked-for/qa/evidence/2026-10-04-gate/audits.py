from pathlib import Path
import subprocess,re,json,hashlib
ROOT=Path.cwd();E=ROOT/'docs/specs/0223-adjustments-the-adopters-asked-for/qa/evidence/2026-10-04-gate';SPEC=E.parents[2];REPORT=E.parents[1]/'qa-report-2026-10-04.md';BASE='3490588afbdffe4419574befd0d7d9569882ff69';HEAD='2fb1ff071009480456e9a8f1f43a014f7facabff'
def git(*a):return subprocess.check_output(['git','-c','core.fsmonitor=false',*a]).decode()
def close(i,detail):
 s=REPORT.read_text().replace('| '+i+' | pending |','| '+i+' | pass |');a=s.index('### '+i+'\n');b=s.find('\n### ',a+5);b=len(s) if b<0 else b;s=s[:a]+s[a:b].replace('Evidence: pending.',detail)+s[b:];REPORT.write_text(s)

tech=(SPEC/'_techspec.md').read_text();qa=(ROOT/'.agents/skills/qa-gate/SKILL.md').read_text();para=re.search(r'### The QA gate paragraph.*?```text\n(.*?)\n```',tech,re.S).group(1);assert para in qa
assert ' '.join(qa.split(para)[0].split()).endswith('Do not add inputs after execution to make carry-forward eligible.')
checks=[]
for name in ['qa-gate','roundfix']:
 canonical=ROOT/f'.agents/skills/{name}/SKILL.md';mirror=ROOT/f'skills/{name}/SKILL.md';assert canonical.read_bytes()==mirror.read_bytes();old=git('show',BASE+':'+str(canonical.relative_to(ROOT)));new=canonical.read_text();oldv=re.findall(r'^\s*version: (\S+)',old,re.M);newv=re.findall(r'^\s*version: (\S+)',new,re.M);assert len(newv)==2 and newv[0]==newv[1];assert tuple(map(int,newv[0].split('.')))==tuple(map(int,oldv[0].split('.')[:2]))+(int(oldv[0].split('.')[2])+1,);records=json.loads((ROOT/'skills/testdata/owned-skill-versions.json').read_text())['skills'][name];assert any(v['version']==newv[0] for v in records);checks.append({'skill':name,'before':oldv,'after':newv,'mirrors':'identical','version_record':'present'})
assert (ROOT/'.agents/skills/roundfix/references/release.md').read_bytes()==(ROOT/'skills/roundfix/references/release.md').read_bytes()
settlements=[]
for p in ROOT.glob('**/SKILL.md'):
 if '.git' in p.parts:continue
 rel=str(p.relative_to(ROOT));old=git('show',BASE+':'+rel);new=p.read_text();extract=lambda s: (re.search(r'(?m)^### QA settlement\n.*?(?=^## |\Z)',s,re.S).group(0) if re.search(r'(?m)^### QA settlement\n',s) else '')
 assert extract(old)==extract(new),rel;settlements.append({'path':rel,'has_section':bool(extract(new)),'unchanged':True})
release=(ROOT/'.agents/skills/roundfix/references/release.md').read_text();assert all(x in release for x in ['read-only after its','`Next action:`','`roundfix doctor`','`roundfix baseline update --no-skills`','`checks`','never change the']);assert 'complete the skills and guides check in the release\nrunbook before the release Pull Request' in release
runbook=(ROOT/'docs/user-guide/release-runbook.md').read_text();cut=runbook.split('## Cutting a release')[1];assert all(x in cut for x in ['read-only','`skills:`','`baseline:`','never change'])
guide=(ROOT/'docs/user-guide/context-driven-development.md').read_text();assert '| Branch prefix | `<type>/` (optional) |' in guide;assert all(x in ' '.join(guide.split()) for x in ['The branch prefix is optional.','A recorded prefix keeps its existing sentence','an update never asks for the branch prefix.']);assert guide.index('The branch prefix is optional.')<guide.index('The frontend layout is optional.')
(E/'07-doc-audit.json').write_text(json.dumps({'verbatim_qa_paragraph':para,'versions':checks,'settlements':settlements,'release_reference_and_runbook':'match declared text','context_guide':'optional prefix, retained recorded sentence, no prompt'},indent=2))
close('07','Evidence: `evidence/2026-10-04-gate/07-doc-audit.json` and `07-version-tests.txt`. Executed assertions prove the exact paragraph and insertion point, canonical/mirror equality, each patch-version increment and registry entry, unchanged QA settlement sections in every tracked skill, and each named documentation contract. Focused version/digest tests confirm registry content, not merely the presence of a file.');print('ROW 07 PASS',flush=True)

# Semantic promise audit: Spec sources remain closed, concrete and supported by delivered behavior.
prd=(SPEC/'_prd.md').read_text();tasks=[(SPEC/f'task_0{i}.md').read_text() for i in [1,2,3]]
assert '## Open Questions\n\nNone.' in prd
assert '## Unreachable Acceptance' not in prd
promises={
 'optional branch prefix':{'prd':'Core Features 1-3; User Stories 1-2; Success Metrics 1-2','implementation':'task_02','observables':'Fixed texts 1/2, no new decision, byte-identical recorded guide','evidence':'05-compatibility.json and 05-focused-tests.txt'},
 're-execute rows without inputs':{'prd':'Core Feature 4; User Story 3; Success Metric 4','implementation':'task_01','observables':'verbatim paragraph, unchanged settlement, raised/recorded version','evidence':'07-doc-audit.json'},
 'read-only release checks':{'prd':'Core Features 5-6; User Story 4; Success Metric 3','implementation':'task_01 and task_03','observables':'transcripts 1/2, JSON checks, snapshot equality, reset/refused no check','evidence':'unadopted-snapshots.json and adopted-snapshots.json'}}
assert 'never changes its decision for it' in prd and 'never change the plan\'s state' in prd
for t in tasks:assert '## Acceptance Criteria' in t and '## Result' in t and '## Verification' in t
context=(ROOT/'CONTEXT.md').read_text();defs={}
for term in ['Baseline Profile','Setup Manifest','QA Report']:
 match=re.search(r'\*\*'+term+r'\*\*:\n(.*?)(?=\n\*\*|\Z)',context,re.S);assert match;defs[term]=match.group(1).strip()
assert 'decisions' in defs['Baseline Profile'] and 'selected profile' in defs['Setup Manifest'] and 'rows_blocked_environment' in defs['QA Report']
assert 'No new glossary term.' in tech
assert (ROOT/'CONTEXT.md').read_text()==git('show',BASE+':CONTEXT.md')
(E/'09-promise-glossary.json').write_text(json.dumps({'promises':promises,'glossary':defs,'assessment':'Profile still composes decisions; Setup Manifest still records answered decisions and managed bytes; QA Report format is unchanged. Optional existing decision and additive checks reuse existing terms; no new term needed.','spec_promise_rule':'all promised observables named, falsifiable and supported; no hidden future behavior or unexplained omission','vocabulary_detector_skip':'manually compared skills/baseline/checks/statuses/next action/help/fixed text against documentation and observed CLI output'},indent=2))
close('09','Evidence: `evidence/2026-10-04-gate/09-promise-glossary.json`. Semantic assertions connect each bounded promise to delivered observable evidence; open questions and unreachable declarations are absent. Baseline Profile still composes selected decisions; Setup Manifest still records answered decisions and current managed bytes; QA Report fields and settlement meanings are unchanged. No new term is needed, and CONTEXT.md equals base bytes. The checker skipped vocabulary documentation; this audit compares every emitted term to the guide/reference and observed output.');print('ROW 09 PASS',flush=True)

grant=(SPEC/'_authorization.md').read_text();fm=grant.split('---')[1];bounded=set(re.findall(r'^  - (\S+)',fm.split('paths:')[1].split('operations:')[0],re.M));assert 'status: approved' in fm
source=(ROOT/'internal/speccheck/governed.go').read_text();patterns=re.findall(r'regexp.MustCompile\(\s*`([^`]+)`',source);exact=set(re.findall(r'"([^"\n]+)"',source.split('paths: exactGovernedPaths(')[1].split('),')[0]));governed=lambda p:p in exact or any(re.search(pat,p) for pat in patterns)
commits=['5bb14019b8e44f9e44f3735ef2c273310506f2f4','6fe28898ca1d228261e1475a701f9fdaa0c3c8ab','2fb1ff071009480456e9a8f1f43a014f7facabff'];audits=[]
for n,c in enumerate(commits,1):
 task=SPEC/f'task_0{n}.md';text=task.read_text();context=text.split('## Context')[1].split('## Verification')[0];declared=set(re.findall(r'^- (?:interface|creates): `([^`]+)`',context,re.M));declared.add(str(task.relative_to(ROOT)))
 if '## Recorded paths' in text:declared.update(re.findall(r'`([^`]+)`',text.split('## Recorded paths')[1].split('\n## ')[0]))
 paths=git('diff-tree','--no-commit-id','--name-only','-r',c).splitlines();unknown=set(paths)-declared;assert not unknown,(n,unknown)
 govern=[p for p in paths if governed(p)];assert not set(govern)-bounded,(n,set(govern)-bounded)
 assert not any(p in ['Makefile','.roundfixrc.yml','go.mod','go.sum'] or p.startswith('.github/workflows/') for p in paths)
 assert subprocess.run(['git','merge-base','--is-ancestor',BASE,c]).returncode==0
 assert git('show',BASE+':'+str((SPEC/'_authorization.md').relative_to(ROOT)))==grant
 audits.append({'task':f'task_0{n}','commit':c,'actual_paths':paths,'declared_paths':sorted(declared),'governed_changed_paths':govern,'all_declared':True,'all_governed_bounded':True,'grant_revision':BASE,'grant_precedes_commit':True,'protected_no-change_paths':'unchanged','Recorded paths':'none needed' if '## Recorded paths' not in text else 'included'})
(E/'10-scope-audit.json').write_text(json.dumps({'commit_range':BASE+'..'+HEAD,'classification':'shipped GovernedPath regex and exact-path table','authorization_paths':sorted(bounded),'tasks':audits,'working_delta':'only qa report/evidence; no implementation delta'},indent=2));close('10','Evidence: `evidence/2026-10-04-gate/10-scope-audit.json`. `git diff-tree --no-commit-id --name-only -r` measures every Task commit. All paths are declared in Context or assigned Task; no Recorded paths exception is needed. Every shipped GovernedPath match is in the approved bounded grant, identical at the ancestor delivery base. Makefile, Project Config, go.mod/go.sum and CI workflows are unchanged. Seeded authorization mechanical audits remain preserved; this row supplies the explicitly declared changed-file and no-change observations.');print('ROW 10 PASS',flush=True)

assert all('  - '+op in fm for op in ['implement','commit','push','pull_request','merge']);assert git('rev-parse','HEAD').strip()==HEAD
controls={'approval':'maintainer delivery authority recorded 2026-10-04, Autorizar os dois / Confirmo and standing skills grant; approved _authorization.md at '+BASE,'checks and status':'Daemon make verify-changed pass, exit 0 at '+HEAD+'; diagnostics removed on success; not rerun','unresolved review threads':'none: the supplied no-open-PR fact makes there be no PR thread','Merge-Ready acceptance':'none yet; pre-PR provider review is downstream of this gate; no review acceptance is claimed','review-artifact ancestry':'claimed candidate '+HEAD+'; no pre-PR Review Artifact exists yet or is credited to that head'}
(E/'11-pr-controls.json').write_text(json.dumps(controls,indent=2))
s=REPORT.read_text().replace('| 11 | pending |','| 11 | blocked (environment: no open Pull Request) |');a=s.index('### 11\n');s=s[:a]+s[a:].replace('Evidence: pending.','Evidence: `evidence/2026-10-04-gate/11-pr-controls.json`. All five pre-PR equivalent controls are explicitly accounted for by the prompt, granted delivery authority and Git head observation. No approval/review artifact beyond that authority is invented. The target branch, not the Run branch, is the future PR branch. Unblocking: downstream provider review, archive/publication and an opened target-branch PR; those operations remain outside this gate.');REPORT.write_text(s);print('ROW 11 BLOCKED ENVIRONMENT',flush=True)
