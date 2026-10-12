from pathlib import Path
import re,subprocess,json,hashlib
root=Path.cwd();s=root/'docs/specs/0256-an-archive-record-that-survives-a-squash-merge';e=s/'qa/evidence/2026-10-12-local';base='9001173f391108413013cb5cd31b28bcf3d6a579'
def git(*args):return subprocess.check_output(['git',*args],text=True)
def norm(t):return ' '.join(t.split())
checks=[]
def has(p,phrase):
 ok=norm(phrase) in norm(Path(p).read_text());checks.append({'path':p,'phrase':phrase,'pass':ok})
for stem,phrases in {'archive':['`--pull-request <n>` stores the delivery pull request in the Archive Record','warns when `source_revision` is not an ancestor of the default branch'],'review':['Only a file named `NNNN-<slug>.md` that declares the Archive Record schema is read as a record','fetches `refs/pull/<N>/head` when the commit is absent'],'runs':['fetches `refs/pull/<N>/head` when the commit is absent'],'deliver':['records the pull request number in the Archive Record in a Pull Request Record Commit']}.items():
 for p in [f'.agents/skills/roundfix/references/{stem}.md',f'skills/roundfix/references/{stem}.md',f'docs/user-guide/commands/{stem}.md']:
  for phrase in phrases:has(p,phrase)
for p in ['.agents/skills/roundfix/references/spec-delivery.md','skills/roundfix/references/spec-delivery.md']:has(p,'after the pull request opens, run `roundfix archive <slug> --pull-request <n>`, commit the record and push')
for p in ['.agents/skills/archive-spec/SKILL.md','skills/archive-spec/SKILL.md']:
 has(p,'### Squash merges');has(p,"A squash merge keeps `source_revision` only in the pull request's head ref; the record names that pull request in `pull_request`")
for phrase in ['**Archive Record**:', 'any other Markdown file in that root is not a record','after a squash merge through Source Revision Recovery','**Archive Command**:','stores the delivery pull request in the record, also in a record the current branch already added','**Source Revision Recovery**:',"A reader's fetch of `refs/pull/<N>/head`",'It writes no ref, never prompts', '**Pull Request Record Commit**:',"The Delivery Queue's commit, made after it opens a Spec's pull request",'proof instead of a second repository gate']:has('CONTEXT.md',phrase)
(e/'phrase-checks.json').write_text(json.dumps(checks,indent=2));print('Phrase checks',len(checks),'failed',[x for x in checks if not x['pass']])
old=['If the commit is absent, context is skipped','If that commit is absent, the Spec stays archived in the report with its disposition']
files=[]
for r in ['docs','skills','.agents/skills']:
 files += [p for p in Path(r).rglob('*.md') if not str(p).startswith('docs/history/') and not str(p).startswith(str(s.relative_to(root))+'/')]
files += [Path('CONTEXT.md'),Path('AGENTS.md')]
hits=[]
for p in sorted(set(files)):
 text=p.read_text()
 for phrase in old:
  if phrase in norm(text):hits.append({'path':str(p),'phrase':phrase,'allowed':False})
 for m in re.finditer(r'(?:commit is absent|commit is missing|source_revision.{0,80}(?:absent|unavailable)|without writes or network requests|writes nothing and makes no network request)',norm(text),re.I):
  context=norm(text)[max(0,m.start()-140):m.end()+210];bad=m.group(0) in ['without writes or network requests','writes nothing and makes no network request'];hits.append({'path':str(p),'hit':m.group(0),'context':context,'allowed':not bad,'reason':'stale no-network/no-write claim contradicts Source Revision Recovery' if bad else 'recovery-first wording or unrelated reader/context; inspected in wording-sweep.json'})
(e/'wording-sweep.json').write_text(json.dumps({'replaced_sentences':old,'files_searched':len(files),'hits':hits},indent=2));print('Wording sweep files',len(files),'hits',len(hits),'disallowed',[x['path'] for x in hits if not x['allowed']])
commits={'task_01':'ef8fe2fa6cb74a040a4b4ed23a0fe969ac8c2c25','task_02':'003f2a8769e196cae168af15aa9920bc0db0f632','task_03':'219a9b9ed461a37ef11e239f0d1f6583d826dba5','task_04':'fccf7bcfe7d0274bea01376a657d609d0f523994'}
auth=(s/'_authorization.md').read_text();bounded=set(re.findall(r'^  - ([.\w/-]+)$',auth.split('operations:')[0],re.M));audits=[]
for task,commit in commits.items():
 text=(s/(task+'.md')).read_text();context=text.split('## Context')[1].split('\n## ')[0];declared=set(re.findall(r'^- (?:instruction|interface|creates): `([^`]+)`',context,re.M));recorded=set(re.findall(r'^- `([^`]+)`',text.split('## Recorded paths')[1] if '## Recorded paths' in text else '',re.M));declared|=recorded;declared.add(str(s.relative_to(root)/(task+'.md')))
 changed=git('diff-tree','--no-commit-id','--name-only','-r',commit).splitlines();governed=[p for p in changed if p.startswith('.agents/skills/') or (p.startswith('skills/') and p.endswith('/SKILL.md'))]
 audits.append({'task':task,'commit':commit,'command':['git','diff-tree','--no-commit-id','--name-only','-r',commit],'changed':changed,'declared':sorted(declared),'recorded_paths':sorted(recorded),'undeclared':sorted(set(changed)-declared),'governed':governed,'unbounded':sorted(set(governed)-bounded),'authorization_ancestor':subprocess.run(['git','merge-base','--is-ancestor',base,commit]).returncode==0,'authorization_unchanged_in_task':str(s.relative_to(root)/'_authorization.md') not in changed})
allpaths=git('diff','--name-only',base+'..HEAD').splitlines();forbidden=[p for p in allpaths if p in ['go.mod','go.sum','.roundfixrc.yml','Makefile','internal/spec/archive.go','docs/agents/setup-context.json'] or p.startswith(('.github/workflows/','internal/baseline/','docs/history/specs/'))];settlement=[]
def section(t):
 start=t.index('### QA settlement\n');tail=t[start:];m=re.search(r'^#{1,3} (?!QA settlement)',tail,re.M);return tail[:m.start()] if m else tail
for folder in ['skills','.agents/skills']:
 for p in Path(folder).rglob('SKILL.md'):
  content=p.read_text()
  if '### QA settlement\n' not in content:continue
  current=section(content);previous=section(git('show',base+':'+str(p)));settlement.append({'path':str(p),'sha256':hashlib.sha256(current.encode()).hexdigest(),'unchanged':current==previous})
mirror=[]
for p in Path('.agents/skills/roundfix').rglob('*.md'):
 dst=Path('skills/roundfix')/p.relative_to('.agents/skills/roundfix');mirror.append({'canonical':str(p),'mirror':str(dst),'equal':p.read_bytes()==dst.read_bytes()})
p=Path('.agents/skills/archive-spec/SKILL.md');mirror.append({'canonical':str(p),'mirror':'skills/archive-spec/SKILL.md','equal':p.read_bytes()==Path('skills/archive-spec/SKILL.md').read_bytes()})
# Every pre-existing Avoid line in the four touched glossary entries is retained.
prior=git('show',base+':CONTEXT.md');current=Path('CONTEXT.md').read_text();avoid=[l for l in prior.splitlines() if '_Avoid_' in l];avoid_missing=[l for l in avoid if l not in current]
res={'range':base+'..HEAD','task_audits':audits,'forbidden_changed_paths':forbidden,'settlement_sections':settlement,'settlement_equal':len(set(x['sha256'] for x in settlement))==1,'mirrors':mirror,'prior_avoid_lines_missing':avoid_missing}
(e/'scope-audit.json').write_text(json.dumps(res,indent=2));print('Scope audit undeclared',[(x['task'],x['undeclared'],x['unbounded']) for x in audits],'forbidden',forbidden,'settlement equal',res['settlement_equal'],'Avoid missing',len(avoid_missing))
assert all(x['pass'] for x in checks);assert not forbidden and not avoid_missing;assert all(not x['undeclared'] and not x['unbounded'] and x['authorization_ancestor'] and x['authorization_unchanged_in_task'] for x in audits);assert res['settlement_equal'] and all(x['unchanged'] for x in settlement) and all(x['equal'] for x in mirror)
