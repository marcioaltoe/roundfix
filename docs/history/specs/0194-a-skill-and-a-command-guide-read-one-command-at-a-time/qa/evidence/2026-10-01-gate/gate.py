from pathlib import Path
from collections import Counter
import subprocess, os, re, tempfile, shutil, json, sys
ROOT=Path.cwd(); E=Path(__file__).resolve().parent; REPORT=E.parent.parent/'qa-report-2026-10-01.md'
SCRATCH=Path(tempfile.mkdtemp(prefix='roundfix-0194-qa-',dir='/private/tmp'))
ENV=os.environ.copy(); ENV['GOCACHE']=str(SCRATCH/'gocache')
# Explicit --dir installation never accesses a user skill directory or Run Database.
(SCRATCH/'roundfix-home').mkdir()
def run(args,log,cwd=ROOT,expected=0):
    p=subprocess.run(['rtk','proxy']+args,cwd=cwd,env=ENV,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
    with (E/log).open('a') as f: f.write('$ '+ ' '.join(args)+'\ncwd: '+str(cwd)+'\n'+p.stdout+'\nexit: '+str(p.returncode)+'\n')
    if expected is not None: assert (p.returncode==0 if expected==0 else p.returncode!=0), (args,p.returncode,p.stdout[-2000:])
    return p.stdout

def status(n,value='pass'):
    s=REPORT.read_text().replace(f'| {n} | pending |',f'| {n} | {value} |');REPORT.write_text(s)

def git(*args):
    return subprocess.check_output(['rtk','proxy','git','-c','core.fsmonitor=false',*args],cwd=ROOT,text=True)

def lines(text,skill):
    if skill and text.startswith('---\n'): text=text.split('---\n',2)[2]
    marker='reference' if skill else 'command'
    text=re.sub(r'<!-- roundfix:'+marker+r'-index:begin -->.*?<!-- roundfix:'+marker+r'-index:end -->','',text,flags=re.S)
    if not skill: text=re.sub(r'\]\([^)]*\)',']()',text)
    return Counter(x for x in text.splitlines() if x.strip())

def proof(commit,entry,comp,skill,tree=None):
    old=lines(git('show',f'{commit}^:{entry}'),skill)
    paths=[entry]+git('ls-tree','-r','--name-only',commit,comp).splitlines()
    new=Counter()
    for path in paths:
        new.update(lines((tree/path).read_text() if tree else git('show',f'{commit}:{path}'),skill))
    return old,new

run(['./bin/roundfix','--version'],'cli.txt')
run(['./bin/roundfix','skills','check'],'cli.txt')
for i in range(2):
    run(['./bin/roundfix','skills','install','--target','codex','--dir',str(SCRATCH/'roundfix-home/skills')],'cli.txt')
    dest=SCRATCH/'roundfix-home/skills/roundfix'; canonical=ROOT/'.agents/skills/roundfix'
    refs=sorted((canonical/'references').glob('*.md')); assert len(refs)==18
    for src in [canonical/'SKILL.md',*refs]: assert src.read_bytes()==(dest/src.relative_to(canonical)).read_bytes()
    with (E/'cli.txt').open('a') as f:f.write(f'Fresh independent read {i+1}: entry and all 18 references match canonical bytes.\n')
status(2)
MOVES=[('1e6fd5455cb88d5d76e7aaad19a3eedbb4661bf5','.agents/skills/roundfix/SKILL.md','.agents/skills/roundfix/references',True),('d1c5e28f','docs/user-guide/commands.md','docs/user-guide/commands',False)]
for commit,entry,comp,skill in MOVES:
    for tree in [None,ROOT]:
        old,new=proof(commit,entry,comp,skill,tree); assert old==new, (old-new,new-old)
        with (E/'move.txt').open('a') as f:f.write(f'{git("rev-parse",commit).strip()} parent {git("rev-parse",commit+"^").strip()}: {sum(old.values())} old / {sum(new.values())} new nonblank normalized lines; equality PASS; comparison: {"HEAD worktree" if tree else "splitting commit"}\n')
status(3)
# Disposable full tree from audited HEAD; no untracked file is imported.
COPY=SCRATCH/'copy';COPY.mkdir()
a=subprocess.Popen(['rtk','proxy','git','-c','core.fsmonitor=false','archive','HEAD'],cwd=ROOT,stdout=subprocess.PIPE)
b=subprocess.run(['rtk','proxy','tar','-xf','-','-C',str(COPY)],stdin=a.stdout,stderr=subprocess.PIPE);a.stdout.close();assert a.wait()==0 and b.returncode==0
(COPY/'bin').mkdir();shutil.copy2(ROOT/'bin/roundfix',COPY/'bin/roundfix')
commit,entry,comp,skill=MOVES[0]
ref=COPY/comp/'runtime.md'; saved=ref.read_text(); ref.write_text('\n'.join(saved.splitlines()[1:])+'\n')
old,new=proof(commit,entry,comp,skill,COPY); assert old!=new
(E/'sabotage.txt').write_text('Deleted nonblank runtime heading: move proof rejected; missing '+repr(dict(old-new))+'\n')
ref.write_text(saved)
other=COPY/comp/'setup.md'; original=other.read_text();other.write_text(original+'\n'+saved)
old,new=proof(commit,entry,comp,skill,COPY);assert old!=new
with (E/'sabotage.txt').open('a') as f:f.write('Copied entire runtime section to setup: move proof rejected; duplicated lines '+str(sum((new-old).values()))+'\n')
other.write_text(original)
idx=COPY/entry;original=idx.read_text(); changed=re.sub(r'^\| \[[^\]]+\]\(references/[^)]+\) \|.*\n','',original,count=1,flags=re.M);assert changed!=original;idx.write_text(changed)
run(['go','test','-count=1','-run','^TestRoundfixSkillIndexNamesEveryReference$','./skills'],'sabotage.txt',COPY,1);idx.write_text(original)
phrase='Prefer `roundfix` commands over manual GitHub scraping.'
modified={}
for path in [COPY/entry,*list((COPY/comp).glob('*.md')),*list((COPY/'skills/roundfix').rglob('*.md'))]:
    text=path.read_text()
    if phrase in text:modified[path]=text;path.write_text(text.replace(phrase,''))
assert modified
run(['go','build','-buildvcs=false','-o','bin/roundfix','./cmd/roundfix'],'sabotage.txt',COPY)
out=run(['./bin/roundfix','skills','check'],'sabotage.txt',COPY,1);assert 'missing required wording' in out and phrase in out
for path,text in modified.items():path.write_text(text)
# Restore the copied production binary before the next mutation.
shutil.copy2(ROOT/'bin/roundfix',COPY/'bin/roundfix')
idx.write_text(original.replace('### QA settlement\n','### QA settlement\n\nSabotage\n',1))
run(['go','test','-count=1','-run','^TestSettlementGuidanceIsOneTable$','./skills'],'sabotage.txt',COPY,1);idx.write_text(original)
status(4)
patterns={
'./skills':r'^(TestText.*|TestRoundfixWording.*|TestRoundfixSkill.*|TestInstallWritesEveryRoundfixReference|TestSettlementGuidanceIsOneTable|TestCheckValidatesRoundfixSkillArtifacts|TestNoPythonBaselineRuntime|TestOwnedSkillContractRejectsSetAndVersionDisagreement|TestOwnedSkillBundleReadinessKeepsStatesDistinct|TestReviewRequestContract|TestProjectConstraint.*Gate)$',
'./internal/cli':r'^(TestEventsHelpDocumentsAgentSelectionFilter|TestBaselineExamplesParse)$',
'./internal/mdtree':r'^TestText.*$'}
for package,pat in patterns.items():run(['go','test','-count=1','-v','-run',pat,package],'pins.txt')
run(['go','test','-count=1','-tags','docscontract','-v','-run','TestEveryCommandIsNamed|TestAnUndocumentedCommand|TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract|TestUserGuide|TestCommandIndexNamesEveryCommandFile|TestCommandReferenceLinksResolve','./internal/docscontract'],'pins.txt')
run(['grep','-rn','--include=*.go','-e',r'commands\.md','-e','"roundfix", "SKILL.md"','-e',r'roundfix/SKILL\.md','.'],'sweep.txt')
status(5)
run(['go','test','-count=1','-v','-run','TestRoundfixSkillIndex|TestRoundfixSkillEntry|TestInstallWritesEveryRoundfixReference','./skills'],'layout.txt')
run(['go','test','-count=1','-tags','docscontract','-v','-run','TestCommandIndexNamesEveryCommandFile|TestCommandReferenceLinksResolve','./internal/docscontract'],'layout.txt')
for path,comp,pattern in [('.agents/skills/roundfix/SKILL.md','.agents/skills/roundfix/references',r'\]\((references/[^)]+)\)'),('docs/user-guide/commands.md','docs/user-guide/commands',r'\]\((commands/[^)]+)\)')]:
    text=(ROOT/path).read_text();refs=re.findall(pattern,text); files=[str(Path(comp).name+'/'+p.name) for p in (ROOT/comp).glob('*.md')];assert Counter(refs)==Counter(files)
    with (E/'layout.txt').open('a') as f:f.write(f'{path}: {len(text.encode())} bytes, {len(text.splitlines())} lines, {len(refs)} resolving index links; exact file-set equality.\n')
status(6)
run(['go','test','-count=1','-v','-run','^(TestWaveCollisionAllowsDifferentCommandFiles|TestWaveCollisionRefusesTheSameCommandFile|TestWriteTasksSkillStatesTheDeclaredPathRules|TestTaskAuthoringGuidanceNamesDeclarations)$','./internal/speccheck','./skills'],'waves.txt')
for path in ['.agents/skills/write-tasks/SKILL.md','skills/write-tasks/SKILL.md','.agents/skills/write-tasks/references/task-template.md','skills/write-tasks/references/task-template.md']:
    text=' '.join((ROOT/path).read_text().split())
    required=['.agents/skills/roundfix/references/<command>.md','docs/user-guide/commands/<command>.md']
    if path.endswith('SKILL.md'):required+=['declares the one command file it changes','cannot share a Wave']
    for phrase in required:assert phrase in text,(path,phrase)
    with (E/'authoring.txt').open('a') as f:f.write(path+': all applicable exact phrases present\n')
run(['make','skills-sync-check'],'authoring.txt');status(7)
(E/'scratch.txt').write_text(str(SCRATCH)+'\n')
print('Rows 2–7 settled; scratch '+str(SCRATCH))
