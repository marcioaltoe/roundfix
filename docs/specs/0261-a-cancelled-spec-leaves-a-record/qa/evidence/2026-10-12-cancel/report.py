from pathlib import Path
import re
REPORT=Path(__file__).resolve().parents[2]/"qa-report-2026-10-12.md"
def settle(row,status,evidence):
 p=REPORT.read_text();p=re.sub(r"(?m)^\| "+str(row)+r" \| [^|]+ \|", "| "+str(row)+" | "+status+" |",p)
 start=p.index("\n### "+str(row)+"\n");end=p.find("\n### ",start+1);end=len(p) if end<0 else end
 block=p[start:end].replace("Evidence: pending.","Evidence: "+evidence)
 REPORT.write_text(p[:start]+block+p[end:])
