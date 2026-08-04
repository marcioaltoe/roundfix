# QA-08 — Failed-report refusal

Status: pass

Without override, `roundfix archive qa-case` exited 2 with the established
diagnostic:

```text
no passing QA verdict: newest QA Report verdict is "fail"; expected "pass"
```

A fresh Spec Root listing retained the active Spec and contained no archived
destination.
