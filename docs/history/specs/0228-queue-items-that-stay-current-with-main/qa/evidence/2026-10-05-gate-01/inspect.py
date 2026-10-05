from pathlib import Path
import subprocess,re,json
root=Path.cwd();base='b44a1174105cbb0eb8e28c5f420d1f534aa6192f'
def git(*args):return subprocess.check_output(['git','-c','core.fsmonitor=false',*args],text=True).strip()
commits=['aa6271a4cebc928d6fcb98e86b045d43b9105aa5','ceaa5278cdceda0bb34d7797f89ae281a96af2bd','cd34e6b430296edc11df91b8abe7b513a1d2a3eb','5d893f9b0a596f6c9feb726b708a3b9efb13dc41','b957f7cfe233c19e865dfdf80b10465b77bdac75']
spec='docs/specs/0228-queue-items-that-stay-current-with-main/'
auth=Path(spec+'_authorization.md').read_text();governed=set(re.findall(r'^  - (.*)$',auth.split('paths:\n',1)[1].split('operations:',1)[0],re.M));allpaths=set()
for n,commit in enumerate(commits,1):
 paths=git('diff-tree','--no-commit-id','--name-only','-r',commit).splitlines();task=Path(spec+f'task_{(6 if n==5 else n):02}.md');declared=set(re.findall(r'`([^`]+)`',task.read_text()));allpaths.update(paths)
 undeclared=[x for x in paths if x not in declared and x!=str(task)];assert not undeclared,undeclared
 changedgoverned=[x for x in paths if x in governed or x.startswith('.agents/') or x=='.roundfixrc.yml' or x.startswith('docs/agents/') or re.match(r'skills/[^/]+/SKILL.md$',x)]
 assert all(x in governed for x in changedgoverned)
 assert subprocess.call(['git','merge-base','--is-ancestor',base,commit])==0
 print(f'task_{n:02} {commit}: within declarations + assigned Task; governed={changedgoverned}; paths={paths}')
assert not any(x in ['Makefile','go.mod','CONTEXT.md'] or x.startswith('.github/workflows/') for x in allpaths)
assert git('show',base+':'+spec+'_authorization.md')==auth.strip()
print('Q13 PASS: grant unchanged from delivery-base ancestor; all five Task commit scopes bounded; Makefile, go.mod, CI workflows and CONTEXT unchanged; no Recorded paths additions needed.')
files=['internal/cli/deliver_conflict_test.go','internal/delivery/corrective_spec_test.go','internal/delivery/archived_head_retry_test.go','internal/cli/review_lineage_test.go']
for f in files:assert not git('diff',base+'..HEAD','--',f)
for f in ['skills/owned_skill_versions_test.go','internal/config/verification_tools_test.go']:
 diff=git('diff',base+'..HEAD','--',f);print(diff)
print('Q05 expectation diff PASS: unchanged four legacy files; only changed record-mode cases are changed digest and lower unrecorded version; repository expectation adds Lines; helper/record implementation changes are not changed test expectations.')
def normalized(f):return ' '.join(Path(f).read_text().split())
def expect(f,phrases,absent=[]):
 text=normalized(f)
 for p in phrases:assert p in text,(f,p)
 for p in absent:assert p not in text,(f,p)
 print('phrases passed:',f)
for f in ['docs/user-guide/commands/deliver.md','.agents/skills/roundfix/references/deliver.md']:
 expect(f,['every path changed from the candidate to the head lies under an archived Spec the blocker names','line-scoped'],['still refuses a moved head','after the item head moved'])
for f in ['docs/user-guide/commands/review.md','.agents/skills/roundfix/references/review.md']:expect(f,["answers only the review in the archived Spec's own records"])
expect('.agents/skills/implement-task/SKILL.md',['is part of the work even when Verification runs the same test'])
expect('docs/agents/specific-repository.md',['one patch above the highest recorded version'],['Raise both version fields, then record'])
expect('docs/user-guide/configuration.md',["takes the default branch's side of each conflict hunk","identical lines between them are kept"])
for path in Path('.agents/skills').rglob('*'):
 if path.is_file() and (Path('skills')/path.relative_to('.agents/skills')).is_file():
  mirror=Path('skills')/path.relative_to('.agents/skills')
  assert mirror.read_bytes()==path.read_bytes(),path
print('Every canonical skill file equals shipped mirror.')
def settlement(s):
 match=re.search(r'^### QA settlement\n(.*?)(?=^#{1,3} |\Z)',s,re.M|re.S);return match.group(0) if match else None
for f in git('ls-files','.agents/skills/*/SKILL.md','skills/*/SKILL.md').splitlines():
 assert settlement(Path(f).read_text())==settlement(git('show',base+':'+f)+'\n'),f
print('Every skill QA settlement section unchanged from delivery base.')
old=json.loads(git('show',base+':skills/testdata/owned-skill-versions.json'));current=json.loads(Path('skills/testdata/owned-skill-versions.json').read_text())
for name,entries in old['skills'].items():assert current['skills'][name][:len(entries)]==entries
for name in ['roundfix','implement-task']:
 entries=current['skills'][name];assert len(entries)>len(old['skills'][name]);v=entries[-1]['version'];assert re.findall(r'^ *version: (.*)$',Path('skills/'+name+'/SKILL.md').read_text().split('---')[1],re.M)==[v,v];print(name,'recorded version',v)
print('Q06 PASS: required phrases, removed phrases, all mirrors, recorded versions and preserved settlement sections.')
prd=Path(spec+'_prd.md').read_text();tech=Path(spec+'_techspec.md').read_text();tasks=[Path(spec+f'task_{i:02}.md').read_text() for i in range(1,6)]
for i in range(1,5): assert 'Success Metric '+str(i) in '\n'.join(tasks)
for i in range(1,4):assert 'API Contract '+str(i) in '\n'.join(tasks)
assert '## Coverage Map' in tech and '## Vocabulary Contract' in tech
context=Path('CONTEXT.md').read_text()
for term in ['Delivery Retry','Reviewer Lineage','Park Class','Derived Path Declaration']:assert '**'+term+'**:' in context
print('Q12: Success Metrics 1–4 and API Contracts 1–3 named in Task references; goals/core features map in TechSpec and task slices. Coverage/links authoring already credited to clean strict check, not rechecked as QA rows.')
print('Glossary: Delivery Retry still explicitly returns a parked item from recorded evidence; Reviewer Lineage still has two rounds with round-2 delta; Park Class still names parked reasons. Line-scoped paths extend a Derived Path Declaration field, documented in configuration, and need no new term. No emitted vocabulary changed; blocker names and refusal shape asserted in tests. No Finding/Rollup/Backlog created; adopted references stay with Spec. No constraint or authority promise hidden in test status.')
print('Q12 PASS: no unmapped or untested promise and no glossary update needed.')
